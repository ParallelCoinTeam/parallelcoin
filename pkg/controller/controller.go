package controller

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"go.uber.org/atomic"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/advertisment"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/controller/sol"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const (
	// MaxDatagramSize is the largest a packet could be,
	// it is a little larger but this is easier to calculate.
	// There is only one listening thread but it needs a buffer this size for
	// worst case largest block possible.
	// Note also this is why FEC is used on the packets in case some get lost it
	// has to puncture 6 of the 9 to fail.
	// This protocol is connectionless and stateless so if one misses,
	// the next one probably won't, usually a second or 3 later
	MaxDatagramSize = blockchain.MaxBlockBaseSize / 3
	//UDP6MulticastAddress = "ff02::1"
	UDP4MulticastAddress = "224.0.0.1:11049"
)

type Controller struct {
	conn                   *transport.Connection
	active                 *atomic.Bool
	ctx                    context.Context
	cx                     *conte.Xt
	mx                     *sync.Mutex
	blockTemplateGenerator *mining.BlkTmplGenerator
	oldBlocks              *atomic.Value
	prevHash               *chainhash.Hash
	lastTxUpdate           time.Time
	lastGenerated          time.Time
	pauseShards            [][]byte
	sendAddresses          []*net.UDPAddr
	subMx                  *sync.Mutex
	submitChan             chan []byte
}

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
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
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := transport.NewConnection(UDP4MulticastAddress,
		*cx.Config.Controller, *cx.Config.MinerPass, MaxDatagramSize,
		ctx, true)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	for !cx.RealNode.SyncManager.IsCurrent() {
		log.DEBUG("node is not synced, waiting 5 seconds to start controller")
		time.Sleep(time.Second*5)
	}
	ctrl := &Controller{
		conn:                   conn,
		active:                 &atomic.Bool{},
		ctx:                    ctx,
		cx:                     cx,
		mx:                     &sync.Mutex{},
		oldBlocks:              &atomic.Value{},
		lastTxUpdate:           time.Now(),
		lastGenerated:          time.Now(),
		pauseShards:            [][]byte{},
		sendAddresses:          []*net.UDPAddr{},
		subMx:                  &sync.Mutex{},
		submitChan:             make(chan []byte),
		blockTemplateGenerator: getBlkTemplateGenerator(cx),
	}
	ctrl.active.Store(false)
	pM := pause.GetPauseContainer(cx)
	pauseShards, err := ctrl.conn.CreateShards(pM.Data, pause.PauseMagic)
	if err != nil {
		log.ERROR(err)
	} else {
		ctrl.active.Store(true)
	}
	ctrl.oldBlocks.Store(pauseShards)
	defer func() {
		log.DEBUG("miner controller shutting down")
		ctrl.active.Store(false)
		err := ctrl.conn.SendShards(pauseShards)
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
	err = ctrl.conn.Listen(handlers, ctrl, nil, nil)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	cx.RealNode.Chain.Subscribe(ctrl.getNotifier())
	go rebroadcaster(ctrl)
	go submitter(ctrl)
	select {
	case <-ctx.Done():
	case <-interrupt.HandlersDone:
	}
	log.TRACE("controller exiting")
	ctrl.active.Store(false)
	return
}

// these are the handlers for specific message types.
// Controller only listens for submissions (currently)
var handlers = transport.HandleFunc{
	string(sol.SolutionMagic): func(ctx interface{}) func(b []byte) (
		err error) {
		return func(b []byte) (err error) {
			log.DEBUG("received solution")
			c := ctx.(*Controller)
			j := sol.LoadSolContainer(b)
			msgBlock := j.GetMsgBlock()
			log.SPEW(msgBlock)
			if !msgBlock.Header.PrevBlock.IsEqual(&c.cx.RPCServer.Cfg.Chain.
				BestSnapshot().Hash) {
				log.WARN("block submitted by kopach miner worker is stale")
				return
			}
			// set old blocks to pause and send pause directly as block is
			// probably a solution
			c.oldBlocks.Store(c.pauseShards)
			err = c.conn.SendShards(c.pauseShards)
			if err != nil {
				log.ERROR(err)
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
							" via CPU miner:", err,
					)
				} else {
					log.WARN("block submitted via kopach miner rejected:", err)
				}
				if isOrphan {
					log.WARN("block is an orphan")
				}
				// maybe something wrong with the network,
				// send current work again
				err = c.sendNewBlockTemplate()
				if err != nil {
					log.DEBUG(err)
				}
				return
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
		}
	},
}

func (c *Controller) sendNewBlockTemplate() (err error) {
	template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
	if template == nil {
		err = errors.New("could not get template")
		log.ERROR(err)
		return
	}
	msgB := template.Block
	fMC := job.Get(c.cx, util.NewBlock(msgB), advertisment.Get(c.cx))
	shards, err := c.conn.CreateShards(fMC.Data, job.WorkMagic)
	err = c.conn.SendShards(shards)
	if err != nil {
		log.ERROR(err)
	}
	c.prevHash = &template.Block.Header.PrevBlock
	c.oldBlocks.Store(shards)
	c.lastGenerated = time.Now()
	c.lastTxUpdate = time.Now()
	return
}

func getNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator,
) (template *mining.BlockTemplate) {
	if len(*cx.Config.MiningAddrs) < 1 {
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
		MiningAddrs))]
	template, err := bTG.NewBlockTemplate(0, payToAddr,
		"sha256d")
	if err != nil {
		log.ERROR(err)
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
				err := ctrl.sendNewBlockTemplate()
				if err != nil {
					log.ERROR(err)
				}
				break
			}
			// The current block is stale if the memory pool has been updated
			// since the block template was generated and it has been at least
			// one minute.
			if ctrl.lastTxUpdate != ctrl.blockTemplateGenerator.GetTxSource().
				LastUpdated() && time.Now().After(
				ctrl.lastGenerated.Add(time.Minute)) {
				err := ctrl.sendNewBlockTemplate()
				if err != nil {
					log.ERROR(err)
				}
				break
			}
			oB, ok := ctrl.oldBlocks.Load().([][]byte)
			if len(oB) == 0 || !ok {
				log.DEBUG("template is empty")
				break
			}
			err := ctrl.conn.SendShards(oB)
			if err != nil {
				log.ERROR(err)
			}
			ctrl.oldBlocks.Store(oB)
		case <-ctrl.ctx.Done():
			break out
			//default:
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
	// check if new transactions have arrived

	// send out new work
}

func (c *Controller) getNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		if c.active.Load() {
			// First to arrive locks out any others while processing
			switch n.Type {
			case blockchain.NTBlockAccepted:
				c.subMx.Lock()
				defer c.subMx.Unlock()
				//log.DEBUG("received new chain notification")
				// construct work message
				//log.SPEW(n)
				_, ok := n.Data.(*util.Block)
				if !ok {
					log.WARN("chain accepted notification is not a block")
					break
				}
				template := getNewBlockTemplate(c.cx, c.blockTemplateGenerator)
				if template != nil {
					//log.DEBUG("got new template")
					msgB := template.Block
					//log.DEBUG(*c.cx.Config.Controller)
					mC := job.Get(c.cx, util.NewBlock(msgB),
						advertisment.Get(c.cx))
					//log.SPEW(mC.Data)
					shards, err := c.conn.CreateShards(mC.Data, job.WorkMagic)
					if err != nil {
						log.TRACE(err)
					}
					c.oldBlocks.Store(shards)
					err = c.conn.SendShards(shards)
					if err != nil {
						log.ERROR(err)
					}
				}
			}
		}
	}
}
