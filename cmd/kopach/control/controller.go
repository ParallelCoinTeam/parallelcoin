package control

import (
	"container/ring"
	"errors"
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"github.com/p9c/pod/cmd/walletmain"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	"github.com/urfave/cli"
	"math/rand"
	"net"
	"time"
	
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util/quit"
	
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/comm/transport"
	rav "github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const (
	MaxDatagramSize      = 8192
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

type Controller struct {
	multiConn              *transport.Channel
	active                 atomic.Bool
	quit                   qu.C
	cx                     *conte.Xt
	isMining               atomic.Bool
	height                 atomic.Int32
	blockTemplateGenerator *mining.BlkTmplGenerator
	msgBlockTemplate       *templates.Message
	oldBlocks              atomic.Value
	lastTxUpdate           atomic.Value
	lastGenerated          atomic.Value
	pauseShards            [][]byte
	sendAddresses          []*net.UDPAddr
	buffer                 *ring.Ring
	began                  time.Time
	otherNodes             map[uint64]*nodeSpec
	uuid                   uint64
	hashCount              atomic.Uint64
	hashSampleBuf          *rav.BufferUint64
	lastNonce              int32
	walletClient           *rpcclient.Client
}

type nodeSpec struct {
	time.Time
	addr string
}

// Run starts up a controller
func Run(cx *conte.Xt) (quit qu.C) {
	if *cx.Config.DisableController {
		Info("controller is disabled")
		return
	}
	cx.Controller.Store(true)
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		Warn("not running controller without RPC enabled")
		return
	}
	if len(*cx.Config.P2PListeners) < 1 || *cx.Config.DisableListen {
		Warn("not running controller without p2p listener enabled", *cx.Config.P2PListeners)
		return
	}
	nS := make(map[uint64]*nodeSpec)
	c := &Controller{
		quit:                   qu.T(),
		cx:                     cx,
		sendAddresses:          []*net.UDPAddr{},
		blockTemplateGenerator: getBlkTemplateGenerator(cx),
		buffer:                 ring.New(BufferSize),
		began:                  time.Now(),
		otherNodes:             nS,
		uuid:                   cx.UUID,
		hashSampleBuf:          rav.NewBufferUint64(100),
	}
	c.isMining.Store(true)
	// maintain connection to wallet if it is available
	var err error
	go c.walletRPCWatcher()
	// c.prevHash.Store(&chainhash.Hash{})
	quit = c.quit
	c.lastTxUpdate.Store(time.Now().UnixNano())
	c.lastGenerated.Store(time.Now().UnixNano())
	c.height.Store(0)
	c.active.Store(false)
	if c.multiConn, err = transport.NewBroadcastChannel(
		"controller", c, *cx.Config.MinerPass, transport.DefaultPort, MaxDatagramSize, handlersMulticast,
		quit,
	); Check(err) {
		c.quit.Q()
		return
	}
	if c.pauseShards = transport.GetShards(p2padvt.Get(cx)); Check(err) {
	} else {
		c.active.Store(true)
	}
	interrupt.AddHandler(
		func() {
			Debug("miner controller shutting down")
			c.active.Store(false)
			if err = c.multiConn.SendMany(pause.Magic, c.pauseShards); Check(err) {
			}
			if err = c.multiConn.Close(); Check(err) {
			}
			c.quit.Q()
		},
	)
	Debug("sending broadcasts to:", UDP4MulticastAddress)
	
	go c.advertiserAndRebroadcaster()
	return
}

func (c *Controller) chainNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		switch n.Type {
		case blockchain.NTBlockConnected:
			Trace("received new chain notification")
			// construct work message
			if b, ok := n.Data.(*util.Block); !ok {
				Warn("chain accepted notification is not a block")
				break
			} else {
				c.height.Store(b.Height())
			}
			
			if err := c.updateAndSendWork(); Check(err) {
			}
		}
	}
}

func (c *Controller) hashReport() float64 {
	c.hashSampleBuf.Add(c.hashCount.Load())
	av := ewma.NewMovingAverage()
	var i int
	var prev uint64
	if err := c.hashSampleBuf.ForEach(
		func(v uint64) error {
			if i < 1 {
				prev = v
			} else {
				interval := v - prev
				av.Add(float64(interval))
				prev = v
			}
			i++
			return nil
		},
	); Check(err) {
	}
	return av.Value()
}

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (c *Controller) GetNewAddressFromWallet() (addr util.Address, err error) {
	if c.walletClient != nil {
		if !c.walletClient.Disconnected() {
			Debug("have access to a wallet, generating address")
			if addr, err = c.walletClient.GetNewAddress("default"); Check(err) {
			} else {
				Debug("-------- found address", addr)
			}
		}
	} else {
		err = errors.New("no wallet available for new address")
		Debug(err)
	}
	return
}

// GetNewAddressFromMiningAddrs tries to get an address from the mining
// addresses list in the configuration file
func (c *Controller) GetNewAddressFromMiningAddrs() (addr util.Address, err error) {
	if c.cx.Config.MiningAddrs == nil {
		err = errors.New("mining addresses is nil")
		Debug(err)
		return
	}
	if len(*c.cx.Config.MiningAddrs) < 1 {
		err = errors.New("no mining addresses")
		Debug(err)
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*c.cx.Config.MiningAddrs))
	addr = c.cx.StateCfg.ActiveMiningAddrs[p2a]
	// remove the address from the state
	if p2a == 0 {
		c.cx.StateCfg.ActiveMiningAddrs = c.cx.StateCfg.ActiveMiningAddrs[1:]
	} else {
		c.cx.StateCfg.ActiveMiningAddrs = append(
			c.cx.StateCfg.ActiveMiningAddrs[:p2a],
			c.cx.StateCfg.ActiveMiningAddrs[p2a+1:]...,
		)
	}
	// update the config
	var ma cli.StringSlice
	for i := range c.cx.StateCfg.ActiveMiningAddrs {
		ma = append(ma, c.cx.StateCfg.ActiveMiningAddrs[i].String())
	}
	*c.cx.Config.MiningAddrs = ma
	save.Pod(c.cx.Config)
	return
}


func (c *Controller) walletRPCWatcher() {
	Debug("starting wallet rpc connection watcher for mining addresses")
	var err error
	backoffTime := time.Second
	certs := walletmain.ReadCAFile(c.cx.Config)
totalOut:
	for {
	trying:
		for {
			select {
			case <-c.cx.KillAll.Wait():
				break totalOut
			default:
			}
			Debug("trying to connect to wallet for mining addresses...")
			// If we can reach the wallet configured in the same datadir we can mine
			if c.walletClient, err = rpcclient.New(
				&rpcclient.ConnConfig{
					Host:         *c.cx.Config.WalletServer,
					Endpoint:     "ws",
					User:         *c.cx.Config.Username,
					Pass:         *c.cx.Config.Password,
					TLS:          *c.cx.Config.TLS,
					Certificates: certs,
				}, nil, c.cx.KillAll,
			); Check(err) {
				Debug("failed, will try again")
				c.isMining.Store(false)
				select {
				case <-time.After(backoffTime):
				case <-c.quit.Wait():
					c.isMining.Store(false)
					break totalOut
				}
				if backoffTime <= time.Second*5 {
					backoffTime += time.Second
				}
				continue
			} else {
				Debug("<<<controller has wallet connection>>>")
				c.isMining.Store(true)
				backoffTime = time.Second
				break trying
			}
		}
		Debug("<<<connected to wallet>>>")
		retryTicker := time.NewTicker(time.Second)
	connected:
		for {
			select {
			case <-retryTicker.C:
				if c.walletClient.Disconnected() {
					c.isMining.Store(false)
					break connected
				}
			case <-c.quit.Wait():
				c.isMining.Store(false)
				break totalOut
			}
		}
		Debug("disconnected from wallet")
	}
}

func (c *Controller) advertiserAndRebroadcaster() {
	if !c.active.Load() {
		Info("ready to send out jobs!")
		c.active.Store(true)
	}
	ticker := time.NewTicker(time.Second)
	const countTick = 10
	counter := countTick / 2
	once := false
	var err error
out:
	for {
		select {
		case <-ticker.C:
			c.height.Store(c.cx.RPCServer.Cfg.Chain.BestSnapshot().Height)
			if c.isMining.Load() {
				if !once {
					c.cx.RealNode.Chain.Subscribe(c.chainNotifier())
					once = true
					c.active.Store(true)
				}
			}
			if counter%countTick == 0 {
				j := p2padvt.GetAdvt(c.cx)
				if *c.cx.Config.AutoListen {
					*c.cx.Config.P2PConnect = cli.StringSlice{}
					_, addresses := routeable.GetAllInterfacesAndAddresses()
					Traces(addresses)
					for i := range addresses {
						addrS := net.JoinHostPort(addresses[i].IP.String(), fmt.Sprint(j.P2P))
						*c.cx.Config.P2PConnect = append(*c.cx.Config.P2PConnect, addrS)
					}
					save.Pod(c.cx.Config)
				}
			}
			counter++
			// send out advertisment
			if err = c.multiConn.SendMany(p2padvt.Magic, transport.GetShards(p2padvt.Get(c.cx))); Check(err) {
			}
			if c.isMining.Load() {
				Debug("updating and sending out new work")
				if err = c.updateAndSendWork(); Check(err) {
				}
			}
		case <-c.quit.Wait():
			Debug("quitting on close quit channel")
			break out
		case <-c.cx.NodeKill.Wait():
			Debug("quitting on NodeKill")
			c.quit.Q()
			break out
		case <-c.cx.KillAll.Wait():
			Debug("quitting on KillAll")
			c.quit.Q()
			break out
		}
	}
	c.active.Store(false)
	Debug("controller exiting")
}

var handlersMulticast = transport.Handlers{
	string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	Debug("processing advertisment message", src, dst)
	c := ctx.(*Controller)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	Trace(j.IPs)
	uuid := j.UUID
	if _, ok := c.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		Debug("uuid", j.UUID, "P2P", j.P2P)
		Info("connecting to lan peer with same PSK", j.IPs, j.UUID)
		// try all IPs
		if *c.cx.Config.AutoListen {
			c.cx.Config.P2PConnect = &cli.StringSlice{}
		}
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			_, addresses := routeable.GetAllInterfacesAndAddresses()
			for i := range addresses {
				addrS := net.JoinHostPort(addresses[i].IP.String(), fmt.Sprint(j.P2P))
				if addrS == peerIP {
					Debug("not connecting to self")
					continue
				}
				if *c.cx.Config.AutoListen {
					*c.cx.Config.P2PConnect = append(*c.cx.Config.P2PConnect, addrS)
				}
			}
			if *c.cx.Config.AutoListen {
				save.Pod(c.cx.Config)
			}
			if err = c.cx.RPCServer.Cfg.ConnMgr.Connect(
				peerIP,
				false,
			); Check(err) {
				continue
			}
			Debug("connected to peer via address", peerIP)
			c.otherNodes[uuid] = &nodeSpec{}
			c.otherNodes[uuid].addr = addr
			break
		}
	}
	// update last seen time for uuid for garbage collection of stale disconnected
	// nodes
	c.otherNodes[uuid].Time = time.Now()
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range c.otherNodes {
		if time.Now().Sub(c.otherNodes[i].Time) > time.Second*9 {
			// also remove from connection manager
			if err = c.cx.RPCServer.Cfg.ConnMgr.RemoveByAddr(c.otherNodes[i].addr); Check(err) {
			}
			Debug("deleting", c.otherNodes[i])
			delete(c.otherNodes, i)
		}
	}
	on := int32(len(c.otherNodes))
	Trace("other nodes", on)
	c.cx.OtherNodes.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (err error) {
	Debug("received solution", src, dst)
	c := ctx.(*Controller)
	if !c.active.Load() { // || !c.cx.Node.Load() {
		Debug("not active yet")
		return
	}
	var s sol.Solution
	gotiny.Unmarshal(b, &s)
	if c.uuid != c.msgBlockTemplate.Nonce {
		Debug("solution not from current controller", c.uuid, c.msgBlockTemplate.Nonce)
		return
	}
	var newHeader *wire.BlockHeader
	if newHeader, err = s.Decode(); Check(err) {
		return
	}
	var msgBlock *wire.MsgBlock
	if msgBlock, err = c.msgBlockTemplate.Reconstruct(newHeader); Check(err) {
		return
	}
	// msgBlock := wire.NewMsgBlock(newHeader)
	Debug("-------------------------------------------------------")
	Debugs(msgBlock)
	if msgBlock.Header.PrevBlock != c.msgBlockTemplate.PrevBlock {
		Debug("block submitted by kopach miner worker is stale")
		if err := c.updateAndSendWork(); Check(err) {
		}
		return
	}
	// // Warn(msgBlock.Header.Version)
	// // cb, ok := c.coinbases.Load().(map[int32]*util.Tx)[msgBlock.Header.Version]
	// Debug("copying over transactions")
	// // copy merkle root
	// txs := append(c.msgBlockTemplate.GetTxs(), c.msgBlockTemplate.GetCoinbase(newHeader.Version))
	// for i := range txs {
	// 	msgBlock.Transactions = append(msgBlock.Transactions, txs[i])
	// }
	// set old blocks to stop and send stop directly as block is probably a
	// solution
	Debug("sending stop to workers")
	if err = c.multiConn.SendMany(pause.Magic, c.pauseShards); Check(err) {
		return
	}
	block := util.NewBlock(msgBlock)
	var isOrphan bool
	Debug("submitting block for processing")
	if isOrphan, err = c.cx.RealNode.SyncManager.ProcessBlock(block, blockchain.BFNone); Check(err) {
		// Anything other than a rule violation is an unexpected error, so log that
		// error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			Warnf(
				"Unexpected error while processing block submitted via kopach miner:", err,
			)
			return
		} else {
			Warn("block submitted via kopach miner rejected:", err)
			if isOrphan {
				Debug("block is an orphan")
				return
			}
			return
		}
	}
	Trace("the block was accepted")
	c.height.Store(block.Height())
	Tracec(
		func() string {
			bmb := block.MsgBlock()
			coinbaseTx := bmb.Transactions[0].TxOut[0]
			prevHeight := block.Height() - 1
			prevBlock, _ := c.cx.RealNode.Chain.BlockByHeight(prevHeight)
			prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
			since := bmb.Header.Timestamp.Unix() - prevTime
			bHash := bmb.BlockHashWithAlgos(block.Height())
			return fmt.Sprintf(
				"new block height %d %08x %s%10d %08x %v %s %ds since prev",
				block.Height(),
				prevBlock.MsgBlock().Header.Bits,
				bHash,
				bmb.Header.Timestamp.Unix(),
				bmb.Header.Bits,
				util.Amount(coinbaseTx.Value),
				fork.GetAlgoName(
					bmb.Header.Version,
					block.Height(),
				), since,
			)
		},
	)
	return
}

// hashrate reports from workers
func processHashrateMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	c := ctx.(*Controller)
	var hr hashrate.Hashrate
	gotiny.Unmarshal(b, &hr)
	if c.lastNonce == hr.Nonce {
		return
	}
	c.lastNonce = hr.Nonce
	// add to total hash counts
	c.hashCount.Store(c.hashCount.Load() + uint64(hr.Count))
	return
}

// GetMsgBlockTemplate gets a Message for the current chain paying to a
// given address
func (c *Controller) GetMsgBlockTemplate(addr util.Address) (mbt *templates.Message, err error) {
	best := c.cx.RealNode.Chain.BestChain.Tip().Header()
	prev := &best
	c.height.Store(c.cx.RealNode.Chain.BestChain.Height() + 1)
	mbt = &templates.Message{
		Nonce:     c.uuid,
		PrevBlock: prev.BlockHash(),
		Height:    c.height.Load(),
		Bits:      make(map[int32]uint32),
		Merkles:   make(map[int32]chainhash.Hash),
	}
	mbt.ResetCoinbases()
	next, curr, more := fork.AlgoVerIterator(mbt.Height)
	for ; more(); next() {
		var templateX *mining.BlockTemplate
		if templateX, err = c.blockTemplateGenerator.NewBlockTemplate(
			0, addr, fork.GetAlgoName(curr(), c.height.Load()),
		); Check(err) {
		} else {
			mbt.SetCoinbase(curr(), templateX.Block.Transactions[len(templateX.Block.Transactions)-1])
			mbt.Bits[curr()] = templateX.Block.Header.Bits
			mbt.Merkles[curr()] = templateX.Block.Header.MerkleRoot
			// Debugf(
			// 	"))))))))))))))))))) %d %d %0.8f %08x %v",
			// 	mbt.Height,
			// 	curr(),
			// 	util.Amount(mbt.GetCoinbase(curr()).TxOut[0].Value).ToDUO(),
			// 	mbt.Bits[curr()],
			// 	mbt.Merkles[curr()],
			// )
			mbt.Timestamp = templateX.Block.Header.Timestamp.Add(time.Second)
			mbt.SetTxs(templateX.Block.Transactions[:len(templateX.Block.Transactions)-1])
			// Debugs(mbt.GetTxs())
			// Debugs(mbt.GetCoinbase(curr()))
		}
	}
	// Debugs(mbt)
	return
}

// GetTemplateMessageShards gets a new address, template message and returns FEC
// shards for the template, and saves the template
func (c *Controller) GetTemplateMessageShards() (o [][]byte, err error) {
	var addr util.Address
	if addr, err = c.GetNewAddressFromMiningAddrs(); Check(err) {
		if addr, err = c.GetNewAddressFromWallet(); Check(err) {
			return
		}
	}
	if c.msgBlockTemplate == nil {
		Debug("getting msgblocktemplate")
		if c.msgBlockTemplate, err = c.GetMsgBlockTemplate(addr); Check(err) {
			return
		}
	}
	o = transport.GetShards(c.msgBlockTemplate.Serialize())
	return
}

func (c *Controller) SendShards(magic []byte, data [][]byte) (err error) {
	if err = c.multiConn.SendMany(magic, data); Check(err) {
	}
	return
}

func (c *Controller) updateAndSendWork() (err error) {
	var getNew bool
	// The current block is stale if the best block has changed.
	oB, ok := c.oldBlocks.Load().([][]byte)
	switch {
	case !ok:
		Trace("cached template is nil")
		getNew = true
	case len(oB) == 0:
		Trace("cached template is zero length")
		getNew = true
	// case c.msgBlockTemplate.PrevBlock != prev.BlockHash():
	// 	Debug("new best block hash")
	// 	getNew = true
	case c.lastTxUpdate.Load() != c.blockTemplateGenerator.GetTxSource().LastUpdated() &&
		time.Now().After(time.Unix(0, c.lastGenerated.Load().(int64)+int64(time.Minute))):
		Trace("block is stale, regenerating")
		getNew = true
		c.lastTxUpdate.Store(time.Now().UnixNano())
		c.lastGenerated.Store(time.Now().UnixNano())
	}
	if getNew {
		if oB, err = c.GetTemplateMessageShards(); Check(err) {
			return
		}
	}
	if err = c.SendShards(job.Magic, oB); Check(err) {
	}
	c.oldBlocks.Store(oB)
	return
}
