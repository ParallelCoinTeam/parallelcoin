package control

import (
	"container/ring"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
	
	"github.com/niubaoshu/gotiny"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/quit"
	
	"github.com/VividCortex/ewma"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/comm/transport"
	rav "github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/util"
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
	height                 atomic.Uint64
	blockTemplateGenerator *mining.BlkTmplGenerator
	coinbases              atomic.Value
	transactions           atomic.Value
	txMx                   sync.Mutex
	oldBlocks              atomic.Value
	prevHash               atomic.Value
	lastTxUpdate           atomic.Value
	lastGenerated          atomic.Value
	pauseShards            [][]byte
	sendAddresses          []*net.UDPAddr
	submitChan             chan []byte
	buffer                 *ring.Ring
	began                  time.Time
	otherNodes             map[string]time.Time
	listenPort             int
	hashCount              atomic.Uint64
	hashSampleBuf          *rav.BufferUint64
	lastNonce              int32
}

func Run(cx *conte.Xt) (quit qu.C) {
	im := true
	cx.Controller.Store(true)
	if len(cx.StateCfg.ActiveMiningAddrs) < 1 {
		im = false
	}
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		Warn("not running controller without RPC enabled")
		return
	}
	if len(*cx.Config.Listeners) < 1 || *cx.Config.DisableListen {
		Warn("not running controller without p2p listener enabled")
		return
	}
	ctrl := &Controller{
		quit:                   qu.T(),
		cx:                     cx,
		sendAddresses:          []*net.UDPAddr{},
		submitChan:             make(chan []byte),
		blockTemplateGenerator: getBlkTemplateGenerator(cx),
		// coinbases:              make(map[int32]*util.Tx),
		buffer:        ring.New(BufferSize),
		began:         time.Now(),
		otherNodes:    make(map[string]time.Time),
		listenPort:    int(util.GetActualPort(*cx.Config.Controller)),
		hashSampleBuf: rav.NewBufferUint64(100),
	}
	ctrl.prevHash.Store(&chainhash.Hash{})
	ctrl.isMining.Store(im)
	quit = ctrl.quit
	ctrl.lastTxUpdate.Store(time.Now().UnixNano())
	ctrl.lastGenerated.Store(time.Now().UnixNano())
	ctrl.height.Store(0)
	ctrl.active.Store(false)
	var err error
	if ctrl.multiConn, err = transport.NewBroadcastChannel(
		"controller",
		ctrl, *cx.Config.MinerPass,
		transport.DefaultPort, MaxDatagramSize, handlersMulticast,
		quit,
	); Check(err) {
		ctrl.quit.Q()
		return
	}
	// var pauseShards [][]byte
	if ctrl.pauseShards = transport.GetShards(p2padvt.Get(cx)); Check(err) {
	} else {
		ctrl.active.Store(true)
	}
	// ctrl.oldBlocks.Store(pauseShards)
	interrupt.AddHandler(
		func() {
			Debug("miner controller shutting down")
			ctrl.active.Store(false)
			if err = ctrl.multiConn.SendMany(pause.Magic, ctrl.pauseShards); Check(err) {
			}
			if err = ctrl.multiConn.Close(); Check(err) {
			}
			ctrl.quit.Q()
		},
	)
	Debug("sending broadcasts to:", UDP4MulticastAddress)
	
	// go advertiser(ctrl)
	factor := 1
	// if err = ctrl.sendNewBlockTemplate(); Check(err) {
	// } else {
	// 	ctrl.active.Store(true)
	// }
	ticker := time.NewTicker(time.Second * time.Duration(factor))
	advt := p2padvt.Get(cx)
	ad := transport.GetShards(advt)
	once := false
	go func() {
	out:
		for {
			select {
			case <-ticker.C:
				// qu.PrintChanState()
				Debug("controller ticker")
				if !ctrl.active.Load() {
					if cx.IsCurrent() {
						Info("ready to send out jobs!")
						ctrl.active.Store(true)
					}
				}
				if ctrl.isMining.Load() {
					if !once {
						cx.RealNode.Chain.Subscribe(ctrl.getNotifier())
						once = true
					}
					if err = ctrl.sendNewBlockTemplate(); Check(err) {
					} else {
						ctrl.active.Store(true)
					}
				}
				// send out advertisment
				// todo: big question: how to deal with change of IP address
				var err error
				if err = ctrl.multiConn.SendMany(p2padvt.Magic, ad); Check(err) {
				}
				if ctrl.isMining.Load() {
					Debug("rebroadcasting")
					// ctrl.rebroadcast()
				}
			case msg := <-ctrl.submitChan:
				Traces(msg)
				decodedB, err := util.NewBlockFromBytes(msg)
				if err != nil {
					Error(err)
					break
				}
				Traces(decodedB)
			case <-ctrl.quit:
				Debug("quitting on close quit channel")
				break out
			case <-ctrl.cx.NodeKill:
				Debug("quitting on NodeKill")
				ctrl.quit.Q()
				break out
			case <-ctrl.cx.KillAll:
				Debug("quitting on KillAll")
				ctrl.quit.Q()
				break out
			}
		}
		ctrl.active.Store(false)
		// panic("aren't we stopped???")
		Debug("controller exiting")
	}()
	return
}

func (c *Controller) rebroadcast() {
	Debug("rebroadcaster ticker")
	// if !c.cx.IsCurrent() {
	// 	Debug("is not current")
	// 	continue
	// } else {
	// 	Debug("is current")
	// }
	Debug("checking for new block")
	// The current block is stale if the best block has changed.
	best := c.blockTemplateGenerator.BestSnapshot()
	if !c.prevHash.Load().(*chainhash.Hash).IsEqual(&best.Hash) {
		c.prevHash.Store(&best.Hash)
		Debug("new best block hash")
		if err := c.sendNewBlockTemplate(); Check(err) {
		}
		return
	}
	Debug("checking for new transactions")
	// The current block is stale if the memory pool has been updated since the
	// block template was generated and it has been at least one minute.
	if c.lastTxUpdate.Load() != c.blockTemplateGenerator.GetTxSource().
		LastUpdated() && time.Now().After(
		time.Unix(
			0,
			c.lastGenerated.Load().(int64)+int64(time.Minute),
		),
	) {
		Trace("block is stale, regenerating")
		if err := c.sendNewBlockTemplate(); Check(err) {
		}
		return
	}
	Debug("checking that block contains payload")
	oB, ok := c.oldBlocks.Load().([][]byte)
	if len(oB) == 0 {
		Warn("template is zero length")
		return
	}
	if !ok {
		Debug("template is nil")
		return
	}
	Debug("sending out job")
	var err error
	if err = c.multiConn.SendMany(job.Magic, oB); Check(err) {
	}
	return
}

func (c *Controller) HashReport() float64 {
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

var handlersMulticast = transport.Handlers{
	string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	c := ctx.(*Controller)
	if !c.active.Load() {
		// Debug("not active")
		return
	}
	// Debug(src, dst)
	// j := p2padvt.LoadContainer(b)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	otherIPs := j.IPs
	// Debug(otherIPs)
	// Trace("otherIPs", otherIPs)
	otherPort := fmt.Sprint(j.P2P)
	// Debug(otherPort)
	myPort := strings.Split((*c.cx.Config.Listeners)[0], ":")[1]
	// Debug(myPort)
	// Trace("myPort", myPort,*c.cx.Config.Listeners)
	for i := range otherIPs {
		o := fmt.Sprintf("%s:%s", otherIPs[i], otherPort)
		if otherPort != myPort {
			if _, ok := c.otherNodes[o]; !ok {
				Debug("ctrl", j.Controller, "P2P", j.P2P, "rpc", j.RPC)
				// because nodes can be set to change their port each launch this always
				// reconnects (for lan, autoports is recommended).
				// TODO: readd autoports for GUI wallet
				Info("connecting to lan peer with same PSK", o, otherIPs)
				if err = c.cx.RPCServer.Cfg.ConnMgr.Connect(o, true); Check(err) {
				}
			}
			c.otherNodes[o] = time.Now()
		}
	}
	for i := range c.otherNodes {
		if time.Now().Sub(c.otherNodes[i]) > time.Second*9 {
			delete(c.otherNodes, i)
		}
	}
	c.cx.OtherNodes.Store(int32(len(c.otherNodes)))
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (err error) {
	Trace("received solution", src, dst)
	c := ctx.(*Controller)
	if !c.active.Load() { // || !c.cx.Node.Load() {
		Debug("not active yet")
		return
	}
	var s sol.Solution
	gotiny.Unmarshal(b, &s)
	// j := sol.LoadSolContainer(b)
	senderPort := s.Port
	if int(senderPort) != c.listenPort {
		return
	}
	msgBlock := s.MsgBlock
	if !msgBlock.Header.PrevBlock.IsEqual(
		&c.cx.RPCServer.Cfg.Chain.
			BestSnapshot().Hash,
	) {
		Debug("block submitted by kopach miner worker is stale")
		if err := c.sendNewBlockTemplate(); Check(err) {
		}
		return
	}
	// Warn(msgBlock.Header.Version)
	cb, ok := c.coinbases.Load().(map[int32]*util.Tx)[msgBlock.Header.Version]
	if !ok {
		Debug("coinbases not found", cb)
		return
	}
	cbs := []*util.Tx{cb}
	msgBlock.Transactions = []*wire.MsgTx{}
	txs := append(cbs, c.transactions.Load().([]*util.Tx)...)
	for i := range txs {
		msgBlock.Transactions = append(msgBlock.Transactions, txs[i].MsgTx())
	}
	// set old blocks to pause and send pause directly as block is probably a
	// solution
	err = c.multiConn.SendMany(pause.Magic, c.pauseShards)
	if err != nil {
		Error(err)
		return
	}
	block := util.NewBlock(msgBlock)
	var isOrphan bool
	if isOrphan, err = c.cx.RealNode.SyncManager.ProcessBlock(
		block,
		blockchain.BFNone,
	); Check(err) {
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
				Warn("block is an orphan")
				return
			}
			return
		}
	}
	Trace("the block was accepted")
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
	if !c.active.Load() {
		Debug("not active")
		return
	}
	var hr hashrate.Hashrate
	gotiny.Unmarshal(b, &hr)
	// hp := hashrate.LoadContainer(b)
	// count := hp.GetCount()
	// nonce := hp.GetNonce()
	
	// hp := hashrate.LoadContainer(b)
	// count := hp.GetCount()
	// nonce := hp.GetNonce()
	if c.lastNonce == hr.Nonce {
		return
	}
	c.lastNonce = hr.Nonce
	// add to total hash counts
	c.hashCount.Store(c.hashCount.Load() + uint64(hr.Count))
	return
}

func (c *Controller) sendNewBlockTemplate() (err error) {
	var template *mining.BlockTemplate
	if template, err = getNewBlockTemplate(c.cx, c.blockTemplateGenerator); Check(err) {
		return
	}
	// Debugs(template)
	if template == nil {
		Debug("template is nil")
		return
	}
	msgB := template.Block
	// c.coinbases = make(map[int32]*util.Tx)
	var txs []*util.Tx
	ccb := make(map[int32]*util.Tx)
	var fMC []byte
	fMC, txs = job.Get(c.cx, util.NewBlock(msgB), &ccb)
	// Debugs(fMC)
	// var jr job.Job
	// gotiny.Unmarshal(fMC, &jr)
	// Debugs(jr)
	c.coinbases.Store(ccb)
	jobShards := transport.GetShards(fMC)
	shardsLen := len(jobShards)
	if shardsLen < 1 {
		Warn("jobShards", shardsLen)
		return fmt.Errorf("jobShards len %d", shardsLen)
	}
	c.oldBlocks.Store(jobShards)
	err = c.multiConn.SendMany(job.Magic, jobShards)
	if err != nil {
		Error(err)
	}
	c.prevHash.Store(&template.Block.Header.PrevBlock)
	c.transactions.Store(txs)
	c.lastGenerated.Store(time.Now().UnixNano())
	c.lastTxUpdate.Store(time.Now().UnixNano())
	return
}

func getNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator) (
	template *mining.BlockTemplate, err error,
) {
	Trace("getting new block template")
	if cx.Config.MiningAddrs == nil {
		Debug("mining addresses is nil")
		return
	}
	if len(*cx.Config.MiningAddrs) < 1 {
		Debug("no mining addresses")
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*cx.Config.MiningAddrs))
	payToAddr := cx.StateCfg.ActiveMiningAddrs[p2a]
	// do this after returning, in the background
	defer func() {
		go func() {
			// remove the address from the state
			if p2a == 0 {
				cx.StateCfg.ActiveMiningAddrs = cx.StateCfg.ActiveMiningAddrs[1:]
			} else {
				cx.StateCfg.ActiveMiningAddrs = append(
					cx.StateCfg.ActiveMiningAddrs[:p2a],
					cx.StateCfg.ActiveMiningAddrs[p2a+1:]...,
				)
			}
			// update the config
			var ma []string
			for i := range cx.StateCfg.ActiveMiningAddrs {
				ma = append(ma, cx.StateCfg.ActiveMiningAddrs[i].String())
			}
			mma := cli.StringSlice(ma)
			cx.Config.MiningAddrs = &mma
			save.Pod(cx.Config)
			// TODO: trigger wallet to generate new ones at some point, if one is connected, when a mined
			// block uses a key and it is deleted here afterwards
		}()
	}()
	Trace("calling new block template")
	if template, err = bTG.NewBlockTemplate(0, payToAddr, fork.SHA256d); Check(err) {
	} else {
		Debug("********** got new block template")
		// Debugs(template)
	}
	return
}

func getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
	policy := mining.Policy{
		BlockMinWeight:    uint32(*cx.Config.BlockMinWeight),
		BlockMaxWeight:    uint32(*cx.Config.BlockMaxWeight),
		BlockMinSize:      uint32(*cx.Config.BlockMinSize),
		BlockMaxSize:      uint32(*cx.Config.BlockMaxSize),
		BlockPrioritySize: uint32(*cx.Config.BlockPrioritySize),
		TxMinFreeFee:      cx.StateCfg.ActiveMinRelayTxFee,
	}
	s := cx.RealNode
	return mining.NewBlkTmplGenerator(
		&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache,
	)
}

func (c *Controller) getNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		if !c.active.Load() {
			Debug("not active")
			return
		}
		// if !c.Ready.Load() {
		// 	Debug("not ready")
		// 	return
		// }
		// First to arrive locks out any others while processing
		switch n.Type {
		case blockchain.NTBlockConnected:
			Trace("received new chain notification")
			// construct work message
			_, ok := n.Data.(*util.Block)
			if !ok {
				Warn("chain accepted notification is not a block")
				break
			}
			if err := c.sendNewBlockTemplate(); Check(err) {
			}
		}
	}
}
