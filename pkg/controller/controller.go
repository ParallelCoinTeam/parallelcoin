package controller

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"
	
	"github.com/VividCortex/ewma"
	"go.uber.org/atomic"
	
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/advertisment"
	"github.com/p9c/pod/pkg/controller/hashrate"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/controller/sol"
	"github.com/p9c/pod/pkg/log"
	rav "github.com/p9c/pod/pkg/ring"
	"github.com/p9c/pod/pkg/simplebuffer/Uint16"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const (
	// MaxDatagramSize is the largest a
	MaxDatagramSize = 8192
	// UDP6MulticastAddress = "ff02::1"
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

type Controller struct {
	multiConn              *transport.Channel
	uniConn                *transport.Channel
	active                 atomic.Bool
	ctx                    context.Context
	cx                     *conte.Xt
	Ready                  atomic.Bool
	height                 atomic.Uint64
	blockTemplateGenerator *mining.BlkTmplGenerator
	coinbases              map[int32]*util.Tx
	transactions           []*util.Tx
	oldBlocks              atomic.Value
	prevHash               *chainhash.Hash
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
}

func Run(cx *conte.Xt) (cancel context.CancelFunc, buffer *ring.Ring) {
	if len(cx.StateCfg.ActiveMiningAddrs) < 1 {
		log.WARN("no mining addresses, not starting controller")
		return
	}
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		log.WARN("not running controller without RPC enabled")
		return
	}
	if len(*cx.Config.Listeners) < 1 || *cx.Config.DisableListen {
		log.WARN("not running controller without p2p listener enabled")
		return
	}
	//	for !cx.RealNode.SyncManager.IsCurrent() {
	//		log.DEBUG("node is not synced, waiting 2 seconds to start controller")
	//		time.Sleep(time.Second * 2)
	//	}
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &Controller{
		ctx:                    ctx,
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
	ctrl.lastTxUpdate.Store(time.Now().UnixNano())
	ctrl.lastGenerated.Store(time.Now().UnixNano())
	ctrl.height.Store(0)
	ctrl.active.Store(false)
	var err error
	ctrl.multiConn, err = transport.NewBroadcastChannel("controller",
		ctrl, *cx.Config.MinerPass,
		transport.DefaultPort, MaxDatagramSize, handlersMulticast)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	// buffer = ctrl.buffer
	pM := pause.GetPauseContainer(cx)
	ll := pM.GetP2PListeners()
	for i := range ll {
		ctrl.otherNodes[ll[i]] = time.Now()
	}
	var pauseShards [][]byte
	if pauseShards = transport.GetShards(pM.Data); log.Check(err) {
	} else {
		// log.DEBUG(pauseShards)
		ctrl.active.Store(true)
	}
	ctrl.oldBlocks.Store(pauseShards)
	defer func() {
		log.DEBUG("miner controller shutting down")
		ctrl.active.Store(false)
		err := ctrl.multiConn.SendMany(pause.PauseMagic, pauseShards)
		if err != nil {
			log.ERROR(err)
		}
	}()
	log.DEBUG("sending broadcasts to:", UDP4MulticastAddress)
	err = ctrl.sendNewBlockTemplate()
	if err != nil {
		log.ERROR(err)
	} else {
		ctrl.active.Store(true)
	}
	// ctrl.uniConn, err = transport.NewUnicastChannel("controller", ctrl, *cx.Config.MinerPass,
	// 	pM.GetIPs()[0].String()+":14422", pM.GetControllerListener()[0], MaxDatagramSize, handlersUnicast)
	// if err != nil {
	// 	log.ERROR(err)
	// 	cancel()
	// 	return
	// }
	cx.RealNode.Chain.Subscribe(ctrl.getNotifier())
	go rebroadcaster(ctrl)
	go submitter(ctrl)
	ticker := time.NewTicker(time.Second)
	cont := true
	for cont {
		select {
		case <-ticker.C:
			log.DEBUGF("network hashrate %.2f", ctrl.HashReport())
		case <-ctx.Done():
			cont = false
		case <-interrupt.HandlersDone:
			cont = false
		}
	}
	log.TRACE("controller exiting")
	ctrl.active.Store(false)
	return
}

func (c *Controller) HashReport() float64 {
	c.hashSampleBuf.Add(c.hashCount.Load())
	av := ewma.NewMovingAverage()
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
	}); log.Check(err) {
	}
	// log.INFO("controller ",c.hashSampleBuf.Cursor, c.hashSampleBuf.Buf)
	// log.INFO("average hashrate", )
	return av.Value()
}

// var handlersUnicast = transport.Handlers{}
var handlersMulticast = transport.Handlers{
	// Solutions submitted by workers
	string(sol.SolutionMagic):
	func(ctx interface{}, src *net.UDPAddr, dst string, b []byte) (err error) {
		log.DEBUG("received solution")
		// log.SPEW(ctx)
		c := ctx.(*Controller)
		if !c.active.Load() || !c.cx.Node.Load().(bool) {
			log.DEBUG("not active yet")
			return
		}
		j := sol.LoadSolContainer(b)
		senderPort := j.GetSenderPort()
		if int(senderPort) != c.listenPort {
			log.DEBUG("not able to submit jobs created by other node peers")
			return
		}
		msgBlock := j.GetMsgBlock()
		// log.WARN(msgBlock.Header.Version)
		// msgBlock.Transactions = append(c.coinbases[msgBlock.Header.Version], c.)
		cb, ok := c.coinbases[msgBlock.Header.Version]
		if !ok {
			log.DEBUG("coinbases not found", cb)
			return
		}
		cbs := []*util.Tx{cb}
		msgBlock.Transactions = []*wire.MsgTx{}
		txs := append(cbs, c.transactions...)
		for i := range txs {
			msgBlock.Transactions = append(msgBlock.Transactions, txs[i].MsgTx())
		}
		// log.SPEW(msgBlock)
		// log.SPEW(c.coinbases)
		// log.SPEW(c.transactions)
		if !msgBlock.Header.PrevBlock.IsEqual(&c.cx.RPCServer.Cfg.Chain.
			BestSnapshot().Hash) {
			log.DEBUG("block submitted by kopach miner worker is stale")
			return
		}
		// set old blocks to pause and send pause directly as block is
		// probably a solution
		c.oldBlocks.Store(c.pauseShards)
		err = c.multiConn.SendMany(pause.PauseMagic, c.pauseShards)
		if err != nil {
			log.ERROR(err)
			return
		}
		block := util.NewBlock(msgBlock)
		isOrphan, err := c.cx.RealNode.SyncManager.ProcessBlock(block,
			blockchain.BFNone)
		if err != nil {
			// Anything other than a rule violation is an unexpected error, so log
			// that error as an internal error.
			if _, ok := err.(blockchain.RuleError); !ok {
				log.WARNF(
					"Unexpected error while processing block submitted"+
						" via kopach miner:", err)
				return
			} else {
				log.WARN("block submitted via kopach miner rejected:", err)
				if isOrphan {
					log.WARN("block is an orphan")
					return
				}
				return
			}
			// // maybe something wrong with the network,
			// // send current work again
			// err = c.sendNewBlockTemplate()
			// if err != nil {
			// 	log.DEBUG(err)
			// }
			// return
		}
		log.DEBUG("the block was accepted")
		coinbaseTx := block.MsgBlock().Transactions[0].TxOut[0]
		prevHeight := block.Height() - 1
		prevBlock, _ := c.cx.RealNode.Chain.BlockByHeight(prevHeight)
		prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
		since := block.MsgBlock().Header.Timestamp.Unix() - prevTime
		bHash := block.MsgBlock().BlockHashWithAlgos(block.Height())
		log.WARNF("new block height %d %08x %s%10d %08x %v %s %ds since prev",
			block.Height(),
			prevBlock.MsgBlock().Header.Bits,
			bHash,
			block.MsgBlock().Header.Timestamp.Unix(),
			block.MsgBlock().Header.Bits,
			util.Amount(coinbaseTx.Value),
			fork.GetAlgoName(block.MsgBlock().Header.Version, block.Height()), since)
		return
	},
	string(job.WorkMagic):
	func(ctx interface{}, src *net.UDPAddr, dst string,
		b []byte) (err error) {
		c := ctx.(*Controller)
		if !c.active.Load() {
			log.DEBUG("not active yet")
			return
		}
		// log.WARN("received job")
		j := job.LoadContainer(b)
		otherIPs := j.GetIPs()
		otherPort := j.GetP2PListenersPort()
		for i := range otherIPs {
			o := fmt.Sprintf("%s:%d", otherIPs[i], otherPort)
			if _, ok := c.otherNodes[o]; !ok {
				// because nodes can be set to change their port each launch this always reconnects (for lan, autoports is
				// recommended).
				log.WARN("connecting to lan peer with same PSK", o)
				c.otherNodes[o] = time.Now()
				go func() {
					if err = c.cx.RPCServer.Cfg.ConnMgr.Connect(o, true); log.Check(err) {
					}
				}()
			}
		}
		return
	},
	// hashrate reports from workers
	string(hashrate.HashrateMagic):
	func(ctx interface{}, src *net.UDPAddr, dst string, b []byte) (err error) {
		c := ctx.(*Controller)
		if !c.active.Load() {
			// log.DEBUG("not active yet")
			return
		}
		hp := hashrate.LoadContainer(b)
		report := hp.Struct()
		// log.DEBUG(report)
		// add to total hash counts
		// current :=
		// log.DEBUG(c.hashCount.Load(), report.Count)
		c.hashCount.Store(c.hashCount.Load() + uint64(report.Count))
		return
	},
}

func (c *Controller) sendNewBlockTemplate() (err error) {
	template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
	// c.coinbases = template.Block.Transactions
	if template == nil {
		err = errors.New("could not get template")
		log.ERROR(err)
		return
	}
	msgB := template.Block
	c.coinbases = make(map[int32]*util.Tx)
	var fMC job.Container
	fMC, c.transactions = job.Get(c.cx, util.NewBlock(msgB), advertisment.Get(c.cx), &c.coinbases)
	shards := transport.GetShards(fMC.Data)
	shardsLen := len(shards)
	if shardsLen < 1 {
		log.WARN("shards", shardsLen)
	}
	err = c.multiConn.SendMany(job.WorkMagic, shards)
	if err != nil {
		log.ERROR(err)
	}
	c.prevHash = &template.Block.Header.PrevBlock
	c.oldBlocks.Store(shards)
	c.lastGenerated.Store(time.Now().UnixNano())
	c.lastTxUpdate.Store(time.Now().UnixNano())
	return
}

func getNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator,
) (template *mining.BlockTemplate) {
	log.DEBUG("getting new block template")
	if len(*cx.Config.MiningAddrs) < 1 {
		log.DEBUG("no mining addresses")
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
		MiningAddrs))]
	log.DEBUG("calling new block template")
	template, err := bTG.NewBlockTemplate(0, payToAddr,
		fork.SHA256d)
	if err != nil {
		log.ERROR(err)
	} else {
		log.DEBUG("got new block template")
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
		s.SigCache, s.HashCache, s.Algo)
}

func rebroadcaster(ctrl *Controller) {
	rebroadcastTicker := time.NewTicker(time.Second)
out:
	for {
		select {
		case <-rebroadcastTicker.C:
			// The current block is stale if the best block has changed.
			best := ctrl.blockTemplateGenerator.BestSnapshot()
			if !ctrl.prevHash.IsEqual(&best.Hash) {
				log.DEBUG("new best block hash")
				ctrl.UpdateAndSendTemplate()
				break
			}
			// // The current block is stale if the memory pool has been updated
			// // since the block template was generated and it has been at least
			// // one minute.
			// if ctrl.lastTxUpdate.Load() != ctrl.blockTemplateGenerator.GetTxSource().
			// 	LastUpdated() && time.Now().After(time.Unix(0,
			// 	ctrl.lastGenerated.Load().(int64)+int64(time.Minute))) {
			// 	log.DEBUG("block is stale")
			// 	ctrl.UpdateAndSendTemplate()
			// 	break
			// }
			oB, ok := ctrl.oldBlocks.Load().([][]byte)
			if len(oB) == 0 {
				log.WARN("template is zero length")
				break
			}
			if !ok {
				log.DEBUG("template is nil")
				break
			}
			// log.DEBUG("rebroadcaster sending out blocks")
			err := ctrl.multiConn.SendMany(job.WorkMagic, oB)
			if err != nil {
				log.ERROR(err)
			}
			ctrl.oldBlocks.Store(oB)
		case <-ctrl.ctx.Done():
			break out
			// default:
		}
	}
}

func submitter(ctrl *Controller) {
out:
	for {
		select {
		case msg := <-ctrl.submitChan:
			log.SPEW(msg)
			decodedB, err := util.NewBlockFromBytes(msg)
			if err != nil {
				log.ERROR(err)
				return
			}
			log.SPEW(decodedB)
			//
		case <-ctrl.ctx.Done():
			break out
		}
	}
}

func updater(ctrl *Controller) {
	// check if new coinbases have arrived
	
	// send out new work
}

func (c *Controller) getNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		if c.active.Load() {
			// First to arrive locks out any others while processing
			switch n.Type {
			case blockchain.NTBlockAccepted:
				log.DEBUG("received new chain notification")
				// construct work message
				_, ok := n.Data.(*util.Block)
				if !ok {
					log.WARN("chain accepted notification is not a block")
					break
				}
				c.UpdateAndSendTemplate()
			}
		} else {
			log.DEBUG("notifier called but controller not active")
		}
	}
}

func (c *Controller) UpdateAndSendTemplate() {
	// log.DEBUG("updating and sending out template")
	c.coinbases = make(map[int32]*util.Tx)
	template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
	if template != nil {
		// log.DEBUG("got a template for sending")
		c.transactions = []*util.Tx{}
		for _, v := range template.Block.Transactions[1:] {
			c.transactions = append(c.transactions, util.NewTx(v))
		}
		// log.DEBUG("got new template")
		msgB := template.Block
		// log.DEBUG(*c.cx.Config.Controller)
		// c.coinbases = msgB.Transactions
		var mC job.Container
		mC, c.transactions = job.Get(c.cx, util.NewBlock(msgB),
			advertisment.Get(c.cx), &c.coinbases)
		nH := mC.GetNewHeight()
		if c.height.Load() < uint64(nH) {
			log.DEBUG("new height")
			c.height.Store(uint64(nH))
		} else {
			log.DEBUG("stale or orphan from being later, not sending out")
			// return
		}
		// log.SPEW(c.coinbases)
		// log.SPEW(mC.Data)
		// log.DEBUG("getting shards for message")
		shards := transport.GetShards(mC.Data)
		c.oldBlocks.Store(shards)
		if err := c.multiConn.SendMany(job.WorkMagic, shards); log.Check(err) {
		}
		c.prevHash = &template.Block.Header.PrevBlock
		c.lastGenerated.Store(time.Now().UnixNano())
		c.lastTxUpdate.Store(time.Now().UnixNano())
		// log.DEBUG("sent out template")
	} else {
		log.DEBUG("got nil template")
	}
}
