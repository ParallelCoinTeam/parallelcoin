package control

import (
	"errors"
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"github.com/p9c/pod/cmd/walletmain"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/comm/transport"
	rav "github.com/p9c/pod/pkg/data/ring"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
	qu "github.com/p9c/pod/pkg/util/quit"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"time"
)

const (
	MaxDatagramSize      = 8192
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

// State stores the state of the controller
type State struct {
	cx                *conte.Xt
	uuid              uint64
	start, stop, quit qu.C
	blockUpdate       chan *util.Block
	generator         *mining.BlkTmplGenerator
	nextAddress       util.Address
	walletClient      *rpcclient.Client
	msgBlockTemplate  *templates.Message
	templateShards    [][]byte
	multiConn         *transport.Channel
	otherNodes        map[uint64]*nodeSpec
	hashSampleBuf     *rav.BufferUint64
	hashCount         atomic.Uint64
}

type nodeSpec struct {
	time.Time
	addr string
}

// New creates a new controller
func New(cx *conte.Xt) (s *State) {
	var err error
	if *cx.Config.DisableController {
		Warn("controller is disabled")
		return
	}
	quit := qu.T()
	var mc *transport.Channel
	if mc, err = transport.NewBroadcastChannel(
		"controller",
		s,
		*cx.Config.MinerPass,
		transport.DefaultPort,
		MaxDatagramSize,
		handlersMulticast,
		quit,
	); Check(err) {
		return
	}
	rand.Seed(time.Now().UnixNano())
	s = &State{
		cx:            cx,
		quit:          quit,
		uuid:          rand.Uint64(),
		start:         qu.Ts(1),
		stop:          qu.Ts(1),
		blockUpdate:   make(chan *util.Block, 1),
		generator:     getBlkTemplateGenerator(cx),
		multiConn:     mc,
		hashSampleBuf: rav.NewBufferUint64(100),
	}
	go func() {
		Debug("starting shutdown signal watcher")
		select {
		case <-cx.KillAll:
			Debug("received killall signal, signalling to quit controller")
			s.Shutdown()
		case <-cx.NodeKill:
			Debug("received nodekill signal, signalling to quit controller")
			s.Shutdown()
		case <-s.quit:
			Debug("received quit signal, breaking out of shutdown signal watcher")
		}
	}()
	s.cx.RealNode.Chain.Subscribe(
		func(n *blockchain.Notification) {
			switch n.Type {
			case blockchain.NTBlockConnected:
				Debug("received block connected notification")
				if b, ok := n.Data.(*util.Block); !ok {
					Warn("block notification is not a block")
					break
				} else {
					s.blockUpdate <- b
				}
			}
		},
	)
	return
}

// Start up the controller
func (s *State) Start() {
	Debug("calling start controller")
	s.start.Signal()
}

// Stop the controller
func (s *State) Stop() {
	Debug("calling stop controller")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	Debug("sending shutdown signal to controller")
	s.quit.Q()
}

func (s *State) startWallet() (err error) {
	Debug("getting configured TLS certificates")
	certs := walletmain.ReadCAFile(s.cx.Config)
	Debug("establishing wallet connection")
	if s.walletClient, err = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         *s.cx.Config.WalletServer,
			Endpoint:     "ws",
			User:         *s.cx.Config.Username,
			Pass:         *s.cx.Config.Password,
			TLS:          *s.cx.Config.TLS,
			Certificates: certs,
		}, nil, s.quit,
	); Check(err) {
	}
	return
}

// Run must be start as a goroutine, central routing for the business of the
// controller
//
// For increased simplicity, every type of work runs in one thread, only signalling
// from background goroutines to trigger state changes.
func (s *State) Run() {
	Debug("starting controller server")
	var err error
	if *s.cx.Config.DisableController {
		Warn("controller is disabled")
		return
	}
	ticker := time.NewTicker(time.Second)
out:
	for {
		Debug("controller now pausing")
	pausing:
		for {
			select {
			case bu := <-s.blockUpdate:
				Debug("received new block update while paused")
				if err = s.doBlockUpdate(bu); Check(err) {
				}
			case <-ticker.C:
				Debug("controller ticker running")
				s.doTicker()
			case <-s.start.Wait():
				Debug("received start signal while paused")
				if s.walletClient.Disconnected() {
					Debug("wallet client is disconnected, retrying")
					if err = s.startWallet(); !Check(err) {
						Debug("wallet client is connected, switching to running")
						break pausing
					}
				}
			case <-s.stop.Wait():
				Debug("received stop signal while paused")
			case <-s.quit.Wait():
				Debug("received quit signal while paused")
				break out
			}
		}
		Debug("controller now running")
	running:
		for {
			select {
			case bu := <-s.blockUpdate:
				Debug("received new block update while running")
				if err = s.doBlockUpdate(bu); Check(err) {
				}
				Debug("sending out templates...")
				if err = s.multiConn.SendMany(job.Magic, s.templateShards); Check(err) {
				}
			case <-ticker.C:
				Debug("controller ticker running")
				s.doTicker()
				Debug("checking if wallet is connected")
				if s.walletClient.Disconnected() {
					Debug("wallet client has disconnected, switching to pausing")
					break running
				}
				Debug("resending current templates...")
				if err = s.multiConn.SendMany(job.Magic, s.templateShards); Check(err) {
				}
			case <-s.start.Wait():
				Debug("received start signal while running")
			case <-s.stop.Wait():
				Debug("received stop signal while running")
				break running
			case <-s.quit.Wait():
				Debug("received quit signal while running")
				break out
			}
		}
		Debug("disconnecting wallet client if it was connected")
		if !s.walletClient.Disconnected() {
			s.walletClient.Disconnect()
		}
	}
}

func (s *State) doTicker() {
	Debug("sending out advertisment")
	var err error
	if err = s.multiConn.SendMany(p2padvt.Magic, transport.GetShards(p2padvt.Get(s.cx))); Check(err) {
	}
}

func (s *State) doBlockUpdate(prev *util.Block) (err error) {
	if s.nextAddress == nil {
		Debug("getting new address for templates")
		if s.nextAddress, err = s.GetNewAddressFromMiningAddrs(); Check(err) {
			if s.nextAddress, err = s.GetNewAddressFromWallet(); Check(err) {
				return
			}
		}
	}
	Debug("getting templates...")
	if s.msgBlockTemplate, err = s.GetMsgBlockTemplate(prev, s.nextAddress); Check(err) {
		return
	}
	Debug("caching error corrected message shards...")
	s.templateShards = transport.GetShards(s.msgBlockTemplate.Serialize())
	return
}

func getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
	Debug("getting a block template generator")
	s := cx.RealNode
	c := cx.Config
	return mining.NewBlkTmplGenerator(
		&mining.Policy{
			BlockMinWeight:    uint32(*c.BlockMinWeight),
			BlockMaxWeight:    uint32(*c.BlockMaxWeight),
			BlockMinSize:      uint32(*c.BlockMinSize),
			BlockMaxSize:      uint32(*c.BlockMaxSize),
			BlockPrioritySize: uint32(*c.BlockPrioritySize),
			TxMinFreeFee:      cx.StateCfg.ActiveMinRelayTxFee,
		},
		s.ChainParams,
		s.TxMemPool,
		s.Chain,
		s.TimeSource,
		s.SigCache,
		s.HashCache,
	)
}

// GetMsgBlockTemplate gets a Message building on given block paying to a given
// address
func (s *State) GetMsgBlockTemplate(prev *util.Block, addr util.Address) (mbt *templates.Message, err error) {
	mbt = &templates.Message{
		UUID:      s.uuid,
		PrevBlock: prev.MsgBlock().BlockHash(),
		Height:    prev.Height() + 1,
		Bits:      make(map[int32]uint32),
		Merkles:   make(map[int32]chainhash.Hash),
	}
	mbt.ResetCoinbases()
	for next, curr, more := fork.AlgoVerIterator(mbt.Height); more(); next() {
		var templateX *mining.BlockTemplate
		if templateX, err = s.generator.NewBlockTemplate(addr, fork.GetAlgoName(curr(), mbt.Height)); Check(err) {
		} else {
			newB := templateX.Block
			newH := newB.Header
			mbt.SetCoinbase(curr(), newB.Transactions[len(newB.Transactions)-1])
			mbt.Bits[curr()] = newH.Bits
			mbt.Merkles[curr()] = newH.MerkleRoot
			mbt.Timestamp = newH.Timestamp.Add(time.Second)
			mbt.SetTxs(newB.Transactions[:len(newB.Transactions)-1])
		}
	}
	return
}

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (s *State) GetNewAddressFromWallet() (addr util.Address, err error) {
	if s.walletClient != nil {
		if !s.walletClient.Disconnected() {
			Debug("have access to a wallet, generating address")
			if addr, err = s.walletClient.GetNewAddress("default"); Check(err) {
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
func (s *State) GetNewAddressFromMiningAddrs() (addr util.Address, err error) {
	if s.cx.Config.MiningAddrs == nil {
		err = errors.New("mining addresses is nil")
		Debug(err)
		return
	}
	if len(*s.cx.Config.MiningAddrs) < 1 {
		err = errors.New("no mining addresses")
		Debug(err)
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*s.cx.Config.MiningAddrs))
	addr = s.cx.StateCfg.ActiveMiningAddrs[p2a]
	// remove the address from the state
	if p2a == 0 {
		s.cx.StateCfg.ActiveMiningAddrs = s.cx.StateCfg.ActiveMiningAddrs[1:]
	} else {
		s.cx.StateCfg.ActiveMiningAddrs = append(
			s.cx.StateCfg.ActiveMiningAddrs[:p2a],
			s.cx.StateCfg.ActiveMiningAddrs[p2a+1:]...,
		)
	}
	// update the config
	var ma cli.StringSlice
	for i := range s.cx.StateCfg.ActiveMiningAddrs {
		ma = append(ma, s.cx.StateCfg.ActiveMiningAddrs[i].String())
	}
	*s.cx.Config.MiningAddrs = ma
	save.Pod(s.cx.Config)
	return
}

var handlersMulticast = transport.Handlers{
	string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	Debug("processing advertisment message", src, dst)
	s := ctx.(*State)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	Trace(j.IPs)
	uuid := j.UUID
	if uuid == s.uuid {
		Debug("ignoring own advertisment message")
		return
	}
	if _, ok := s.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		Info("connecting to lan peer with same PSK", j.IPs, j.UUID)
		// try all IPs
		if *s.cx.Config.AutoListen {
			s.cx.Config.P2PConnect = &cli.StringSlice{}
		}
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			if err = s.cx.RPCServer.Cfg.ConnMgr.Connect(
				peerIP,
				false,
			); Check(err) {
				continue
			}
			Debug("connected to peer via address", peerIP)
			s.otherNodes[uuid] = &nodeSpec{}
			s.otherNodes[uuid].addr = addr
			break
		}
	}
	// update last seen time for uuid for garbage collection of stale disconnected
	// nodes
	s.otherNodes[uuid].Time = time.Now()
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range s.otherNodes {
		if time.Now().Sub(s.otherNodes[i].Time) > time.Second*9 {
			// also remove from connection manager
			if err = s.cx.RPCServer.Cfg.ConnMgr.RemoveByAddr(s.otherNodes[i].addr); Check(err) {
			}
			Debug("deleting", s.otherNodes[i])
			delete(s.otherNodes, i)
		}
	}
	on := int32(len(s.otherNodes))
	s.cx.OtherNodes.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (err error) {
	Debug("received solution", src, dst)
	s := ctx.(*State)
	var so sol.Solution
	gotiny.Unmarshal(b, &so)
	if s.uuid != s.msgBlockTemplate.UUID {
		Debug("solution not from current controller", s.uuid, s.msgBlockTemplate.UUID)
		return
	}
	var newHeader *wire.BlockHeader
	if newHeader, err = so.Decode(); Check(err) {
		return
	}
	if newHeader.PrevBlock != s.msgBlockTemplate.PrevBlock {
		Debug("block submitted by kopach miner worker is stale")
		return
	}
	var msgBlock *wire.MsgBlock
	if msgBlock, err = s.msgBlockTemplate.Reconstruct(newHeader); Check(err) {
		return
	}
	Debug("sending pause to workers")
	if err = s.multiConn.SendMany(pause.Magic, transport.GetShards(p2padvt.Get(s.cx))); Check(err) {
		return
	}
	block := util.NewBlock(msgBlock)
	var isOrphan bool
	Debug("submitting block for processing")
	if isOrphan, err = s.cx.RealNode.SyncManager.ProcessBlock(block, blockchain.BFNone); Check(err) {
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
	Debug("clearing address used for block")
	s.nextAddress = nil
	Debug("the block was accepted, new height", block.Height())
	Tracec(
		func() string {
			bmb := block.MsgBlock()
			coinbaseTx := bmb.Transactions[0].TxOut[0]
			prevHeight := block.Height() - 1
			prevBlock, _ := s.cx.RealNode.Chain.BlockByHeight(prevHeight)
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
	c.hashCount.Add(uint64(hr.Count))
	return
}

func (s *State) hashReport() float64 {
	s.hashSampleBuf.Add(s.hashCount.Load())
	av := ewma.NewMovingAverage()
	var i int
	var prev uint64
	if err := s.hashSampleBuf.ForEach(
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
