package control

import (
	"container/ring"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

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
	"github.com/p9c/pod/pkg/coding/simplebuffer/Uint16"
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
	multiConn *transport.Channel
	// uniConn                *transport.Channel
	active                 atomic.Bool
	quit                   chan struct{}
	cx                     *conte.Xt
	Ready                  atomic.Bool
	height                 atomic.Uint64
	blockTemplateGenerator *mining.BlkTmplGenerator
	coinbases              map[int32]*util.Tx
	transactions           []*util.Tx
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

func Run(cx *conte.Xt) (quit chan struct{}) {
	mining := true
	if len(cx.StateCfg.ActiveMiningAddrs) < 1 {
		// Warn("no mining addresses, not starting controller")
		// return
		mining = false
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
		quit:                   make(chan struct{}),
		cx:                     cx,
		sendAddresses:          []*net.UDPAddr{},
		submitChan:             make(chan []byte),
		blockTemplateGenerator: getBlkTemplateGenerator(cx),
		coinbases:              make(map[int32]*util.Tx),
		buffer:                 ring.New(BufferSize),
		began:                  time.Now(),
		otherNodes:             make(map[string]time.Time),
		listenPort:             int(Uint16.GetActualPort(*cx.Config.Controller)),
		hashSampleBuf:          rav.NewBufferUint64(1000),
	}
	quit = ctrl.quit
	ctrl.lastTxUpdate.Store(time.Now().UnixNano())
	ctrl.lastGenerated.Store(time.Now().UnixNano())
	ctrl.height.Store(0)
	ctrl.active.Store(false)
	var err error
	ctrl.multiConn, err = transport.NewBroadcastChannel("controller",
		ctrl, *cx.Config.MinerPass,
		transport.DefaultPort, MaxDatagramSize, handlersMulticast,
		ctrl.quit)
	if err != nil {
		Error(err)
		close(ctrl.quit)
		return
	}
	pM := pause.GetPauseContainer(cx)
	var pauseShards [][]byte
	if pauseShards = transport.GetShards(pM.Data); Check(err) {
	} else {
		ctrl.active.Store(true)
	}
	// ctrl.oldBlocks.Store(pauseShards)
	interrupt.AddHandler(func() {
		Debug("miner controller shutting down")
		ctrl.active.Store(false)
		err := ctrl.multiConn.SendMany(pause.PauseMagic, pauseShards)
		if err != nil {
			Error(err)
		}
		if err = ctrl.multiConn.Close(); Check(err) {
		}
	})
	Debug("sending broadcasts to:", UDP4MulticastAddress)
	if mining {
		err = ctrl.sendNewBlockTemplate()
		if err != nil {
			Error(err)
		} else {
			ctrl.active.Store(true)
		}
		cx.RealNode.Chain.Subscribe(ctrl.getNotifier())
		go rebroadcaster(ctrl)
		go submitter(ctrl)
	}
	go advertiser(ctrl)
	ticker := time.NewTicker(time.Second * 3)
	cont := true
	for cont {
		select {
		case <-ticker.C:
			if !ctrl.Ready.Load() {
				if cx.IsCurrent() {
					Info("ready to mine!")
					ctrl.Ready.Store(true)
					ctrl.active.Store(true)
				}
			}
			Debugf("cluster hashrate %.2f", ctrl.HashReport())
		case <-ctrl.quit:
			cont = false
			ctrl.active.Store(false)
		case <-interrupt.HandlersDone:
			cont = false
		}
	}
	Trace("controller exiting")
	return
}

func (c *Controller) HashReport() float64 {
	c.hashSampleBuf.Add(c.hashCount.Load())
	av := ewma.NewMovingAverage(15)
	var i int
	var prev uint64
	if err := c.hashSampleBuf.ForEach(func(v uint64) error {
		if i < 1 {
			prev = v
		} else {
			interval := v - prev
			av.Add(float64(interval))
			prev = v
		}
		i++
		return nil
	}); Check(err) {
	}
	return av.Value()
}

var handlersMulticast = transport.Handlers{
	// Solutions submitted by workers
	string(sol.SolutionMagic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		Trace("received solution")
		c := ctx.(*Controller)
		if !c.active.Load() { // || !c.cx.Node.Load() {
			Debug("not active yet")
			return
		}
		j := sol.LoadSolContainer(b)
		senderPort := j.GetSenderPort()
		if int(senderPort) != c.listenPort {
			return
		}
		msgBlock := j.GetMsgBlock()
		if !msgBlock.Header.PrevBlock.IsEqual(&c.cx.RPCServer.Cfg.Chain.
			BestSnapshot().Hash) {
			Debug("block submitted by kopach miner worker is stale")
			// c.UpdateAndSendTemplate()
			return
		}
		// Warn(msgBlock.Header.Version)
		cb, ok := c.coinbases[msgBlock.Header.Version]
		if !ok {
			Debug("coinbases not found", cb)
			return
		}
		cbs := []*util.Tx{cb}
		msgBlock.Transactions = []*wire.MsgTx{}
		txs := append(cbs, c.transactions...)
		for i := range txs {
			msgBlock.Transactions = append(msgBlock.Transactions, txs[i].MsgTx())
		}
		// set old blocks to pause and send pause directly as block is
		// probably a solution
		err = c.multiConn.SendMany(pause.PauseMagic, c.pauseShards)
		if err != nil {
			Error(err)
			return
		}
		block := util.NewBlock(msgBlock)
		isOrphan, err := c.cx.RealNode.SyncManager.ProcessBlock(block,
			blockchain.BFNone)
		if err != nil {
			// Anything other than a rule violation is an unexpected error, so log
			// that error as an internal error.
			if _, ok := err.(blockchain.RuleError); !ok {
				Warnf(
					"Unexpected error while processing block submitted"+
						" via kopach miner:", err)
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
		coinbaseTx := block.MsgBlock().Transactions[0].TxOut[0]
		prevHeight := block.Height() - 1
		prevBlock, _ := c.cx.RealNode.Chain.BlockByHeight(prevHeight)
		prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
		since := block.MsgBlock().Header.Timestamp.Unix() - prevTime
		bHash := block.MsgBlock().BlockHashWithAlgos(block.Height())
		Warnf("new block height %d %08x %s%10d %08x %v %s %ds since prev",
			block.Height(),
			prevBlock.MsgBlock().Header.Bits,
			bHash,
			block.MsgBlock().Header.Timestamp.Unix(),
			block.MsgBlock().Header.Bits,
			util.Amount(coinbaseTx.Value),
			fork.GetAlgoName(block.MsgBlock().Header.Version,
				block.Height()), since)
		return
	},
	string(p2padvt.Magic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		c := ctx.(*Controller)
		if !c.active.Load() {
			// Debug("not active")
			return
		}
		j := p2padvt.LoadContainer(b)
		otherIPs := j.GetIPs()
		// Trace("otherIPs", otherIPs)
		otherPort := fmt.Sprint(j.GetP2PListenersPort())
		myPort := strings.Split((*c.cx.Config.Listeners)[0], ":")[1]
		// Trace("myPort", myPort,*c.cx.Config.Listeners)
		for i := range otherIPs {
			o := fmt.Sprintf("%s:%s", otherIPs[i], otherPort)
			if otherPort != myPort {
				if _, ok := c.otherNodes[o]; !ok {
					Debug("ctrl", j.GetControllerListenerPort(), "P2P",
						j.GetP2PListenersPort(), "rpc", j.GetRPCListenersPort())
					// because nodes can be set to change their port each launch this always reconnects (for lan,
					// autoports is recommended).
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
	},
	// hashrate reports from workers
	string(hashrate.HashrateMagic): func(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
		c := ctx.(*Controller)
		if !c.active.Load() {
			Debug("not active")
			return
		}
		hp := hashrate.LoadContainer(b)
		count := hp.GetCount()
		nonce := hp.GetNonce()
		if c.lastNonce == nonce {
			return
		}
		c.lastNonce = nonce
		// add to total hash counts
		c.hashCount.Store(c.hashCount.Load() + uint64(count))
		return
	},
}

func (c *Controller) sendNewBlockTemplate() (err error) {
	template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
	if template == nil {
		err = errors.New("could not get template")
		Error(err)
		return
	}
	msgB := template.Block
	c.coinbases = make(map[int32]*util.Tx)
	var fMC job.Container
	adv := p2padvt.Get(c.cx)
	Traces(adv)
	fMC, c.transactions = job.Get(c.cx, util.NewBlock(msgB), adv, &c.coinbases)
	jobShards := transport.GetShards(fMC.Data)
	shardsLen := len(jobShards)
	if shardsLen < 1 {
		Warn("jobShards", shardsLen)
		return fmt.Errorf("jobShards len %d", shardsLen)
	}
	err = c.multiConn.SendMany(job.Magic, jobShards)
	if err != nil {
		Error(err)
	}
	c.prevHash.Store(&template.Block.Header.PrevBlock)
	c.oldBlocks.Store(jobShards)
	c.lastGenerated.Store(time.Now().UnixNano())
	c.lastTxUpdate.Store(time.Now().UnixNano())
	return
}

func getNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator,
) (template *mining.BlockTemplate) {
	Trace("getting new block template")
	if len(*cx.Config.MiningAddrs) < 1 {
		Debug("no mining addresses")
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
		MiningAddrs))]
	Trace("calling new block template")
	template, err := bTG.NewBlockTemplate(0, payToAddr,
		fork.SHA256d)
	if err != nil {
		Error(err)
	} else {
		// Debug("got new block template")
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
	return mining.NewBlkTmplGenerator(&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache)
}

func advertiser(ctrl *Controller) {
	advertismentTicker := time.NewTicker(time.Second)
	advt := p2padvt.Get(ctrl.cx)
	ad := transport.GetShards(advt.CreateContainer(p2padvt.Magic).Data)
out:
	for {
		select {
		case <-advertismentTicker.C:
			err := ctrl.multiConn.SendMany(p2padvt.Magic, ad)
			if err != nil {
				Error(err)
			}
		case <-ctrl.quit:
			break out
		}
	}
}

func rebroadcaster(c *Controller) {
	rebroadcastTicker := time.NewTicker(time.Second)
out:
	for {
		select {
		case <-rebroadcastTicker.C:
			if !c.cx.IsCurrent() {
				break
			}
			// The current block is stale if the best block has changed.
			best := c.blockTemplateGenerator.BestSnapshot()
			if !c.prevHash.Load().(*chainhash.Hash).IsEqual(&best.Hash) {
				Debug("new best block hash")
				c.UpdateAndSendTemplate()
				break
			}
			// The current block is stale if the memory pool has been updated since the block template was generated and
			// it has been at least one minute.
			if c.lastTxUpdate.Load() != c.blockTemplateGenerator.GetTxSource().
				LastUpdated() && time.Now().After(time.Unix(0,
				c.lastGenerated.Load().(int64)+int64(time.Minute))) {
				Debug("block is stale")
				c.UpdateAndSendTemplate()
				break
			}
			oB, ok := c.oldBlocks.Load().([][]byte)
			if len(oB) == 0 {
				Warn("template is zero length")
			}
			if !ok {
				Debug("template is nil")
			}
			err := c.multiConn.SendMany(job.Magic, oB)
			if err != nil {
				Error(err)
			}
			c.oldBlocks.Store(oB)
			break
		case <-c.quit:
			break out
		}
	}
}

func submitter(c *Controller) {
out:
	for {
		select {
		case msg := <-c.submitChan:
			Traces(msg)
			decodedB, err := util.NewBlockFromBytes(msg)
			if err != nil {
				Error(err)
				break
			}
			Traces(decodedB)
		case <-c.quit:
			break out
		}
	}
}

func (c *Controller) getNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		if !c.active.Load() {
			// Debug("not active")
			return
		}
		if !c.Ready.Load() {
			// Debug("not ready")
			return
		}
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
			c.UpdateAndSendTemplate()
		}
	}
}

func (c *Controller) UpdateAndSendTemplate() {
	c.coinbases = make(map[int32]*util.Tx)
	template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
	if template != nil {
		c.transactions = []*util.Tx{}
		for _, v := range template.Block.Transactions[1:] {
			c.transactions = append(c.transactions, util.NewTx(v))
		}
		msgB := template.Block
		var mC job.Container
		mC, c.transactions = job.Get(c.cx, util.NewBlock(msgB),
			p2padvt.Get(c.cx), &c.coinbases)
		nH := mC.GetNewHeight()
		if c.height.Load() < uint64(nH) {
			Trace("new height", nH)
			c.height.Store(uint64(nH))
		}
		shards := transport.GetShards(mC.Data)
		c.oldBlocks.Store(shards)
		if err := c.multiConn.SendMany(job.Magic, shards); Check(err) {
		}
		c.prevHash.Store(&template.Block.Header.PrevBlock)
		c.lastGenerated.Store(time.Now().UnixNano())
		c.lastTxUpdate.Store(time.Now().UnixNano())
	} else {
		Debug("got nil template")
	}
}
