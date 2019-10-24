// Package controller implements the server side of a lan multicast mining work dispatcher that goes with the
// `kopach` miner with automatic network topology based failover and redundancy for new block and solution broadcast
//go:generate go run ../../cmd/tools/genmsghandle/main.go controller Solution broadcast.SolBlock . msghandle.go
package controller

import (
	"context"
	"fmt"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util"
	"go.uber.org/atomic"
	"math/rand"
	"sync"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/p9c/pod/pkg/broadcast"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
)

type Templates []*mining.BlockTemplate
type Blocks struct {
	sync.Mutex
	Templates
	New bool
}
type Solution wire.MsgBlock

// Run starts a controller instance
func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	for len(cx.StateCfg.ActiveMiningAddrs) < 1 {
		log.ERROR("no mining addresses configured, controller exiting")
		return
	}
	log.WARN("starting controller")
	// create context with canceller to cleanly shut down
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	// create cipher for decoding relevant packets
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	// create new multicast address
	outAddr, err := broadcast.New(*cx.Config.BroadcastAddress)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	// create buffer and load into msgpack codec
	var mh codec.MsgpackHandle
	bytes := make([]byte, 0, broadcast.MaxDatagramSize)
	enc := codec.NewEncoderBytes(&bytes, &mh)
	// create channel to trigger a broadcast,
	// unbuffered so the encoder buffer is not accessed concurrently - one
	// cpu thread is enough to handle all the traffic so lower overhead if it
	// is not multi-threaded
	blockChan := make(chan Blocks)
	ticker := time.NewTicker(time.Second * 3)
	oldBlocks := Blocks{New: false}
	var pauseRebroadcast atomic.Bool
	pauseRebroadcast.Store(true)
	// work dispatch loop
	go func() {
		for {
			select {
			case lb := <-blockChan:
				pauseRebroadcast.Store(true)
				log.WARN("saving old Blocks for old block rebroadcast")
				oldBlocks.Lock()
				oldBlocks.Templates = lb.Templates
				oldBlocks.Unlock()
				log.DEBUG("sending out block broadcast")
				err := enc.Encode(lb)
				if err != nil {
					log.ERROR(err)
					break
				}
				log.SPEW(bytes)
				err = broadcast.Send(outAddr, bytes, ciph, broadcast.Template)
				if err != nil {
					log.ERROR(err)
				}
				// reset the bytes for next round
				bytes = bytes[:0]
				enc.ResetBytes(&bytes)
				pauseRebroadcast.Toggle()
			case <-ticker.C:
				if !pauseRebroadcast.Load() {
					oldBlocks.Lock()
					fmt.Print("\rsending out old block broadcast ", time.Now(), oldBlocks.New)
					err := enc.Encode(oldBlocks)
					if err != nil {
						log.ERROR(err)
						break
					}
					//log.SPEW(bytes)
					err = broadcast.Send(outAddr, bytes, ciph, broadcast.Template)
					if err != nil {
						log.ERROR(err)
					}
					// reset the bytes for next round
					bytes = bytes[:0]
					enc.ResetBytes(&bytes)
					oldBlocks.Unlock()
				}
			case <-ctx.Done():
				// cancel has been called
				return
				// default:
			}
		}
	}()
	// generate initial Blocks
	initialBlocks := Blocks{New: true}
	for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
		BestSnapshot().Height+1)].Algos {
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		log.TRACE("len active mining addrs", len(cx.StateCfg.ActiveMiningAddrs))
		payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
			StateCfg.ActiveMiningAddrs))]
		template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
			payToAddr, algo)
		if err != nil {
			log.ERROR("failed to create new block template:", err)
			continue
		}
		initialBlocks.Templates = append(initialBlocks.Templates, template)
	}
	pauseRebroadcast.Store(true)
	// send out the block templates
	blockChan <- initialBlocks
	pauseRebroadcast.Toggle()
	// shortcuts to necessary modules
	chainHandle := cx.RPCServer.Cfg.Chain
	generator := cx.RPCServer.Cfg.Generator
	// create subscriber for new block
	cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.Notification) {
		switch n.Type {
		case chain.NTBlockConnected:
			pauseRebroadcast.Store(true)
			log.DEBUG("new block connected to chain")
			bh := chainHandle.BestSnapshot().Height
			hf := fork.GetCurrent(bh + 1)
			blocks := Blocks{New: true}
			// generate Blocks
			for algo := range fork.List[hf].Algos {
				// Choose a payment address at random.
				rand.Seed(time.Now().UnixNano())
				payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.StateCfg.ActiveMiningAddrs))]
				template, err := generator.NewBlockTemplate(0, payToAddr, algo)
				if err != nil {
					log.ERROR("failed to create new block template:", err)
					continue
				}
				blocks.Templates = append(blocks.Templates, template)
			}
			blockChan <- blocks
			pauseRebroadcast.Toggle()
		}
	})
	// goroutine loop checking for connection and sync status
	//if !*cx.Config.Solo {
	go func() {
		// allow a little time for all goroutines to fire up
		time.Sleep(time.Second * 5)
		for {
			time.Sleep(time.Second)
			connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
			current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
			// if out of sync or disconnected,
			// once a second send out empty initialBlocks
			if connCount < 1 {
				log.DEBUG("node is offline", current)
				blockChan <- Blocks{New: false}
			}
			if !current {
				log.DEBUG("node is not current", current)
				blockChan <- Blocks{New: false}
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}()
	//}
	// goroutine loop checking for updates to block template consist
	go func() {
		lastTxUpdate := cx.RPCServer.Cfg.Generator.GetTxSource().LastUpdated()
		time.Sleep(time.Second * 5)
		for {
			// this check is much more frequent as we want to ensure
			// transactions are cleared immediately they appear if possible,
			// while not disrupting mining progress excessively
			time.Sleep(time.Second / 10)
			// when new transactions are received the last updated timestamp
			// changes, when this happens a new dispatch needs to be made
			if lastTxUpdate != cx.RPCServer.Cfg.Generator.GetTxSource().LastUpdated() {
				pauseRebroadcast.Store(true)
				oldBlocks.Lock()
				log.WARN("new transactions added to block, dispatching new block templates")
				blocks := Blocks{New: true}
				// generate Blocks
				for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
					BestSnapshot().Height+1)].Algos {
					// Choose a payment address at random.
					rand.Seed(time.Now().UnixNano())
					payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
						StateCfg.ActiveMiningAddrs))]
					template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
						payToAddr, algo)
					if err != nil {
						log.ERROR("failed to create new block template:", err)
						continue
					}
					blocks.Templates = append(blocks.Templates, template)
				}
				oldBlocks.Templates = blocks.Templates
				blockChan <- blocks
				oldBlocks.Unlock()
				pauseRebroadcast.Toggle()
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}()
	returnChan := make(chan *Solution)
	m := newMsgHandle(*cx.Config.MinerPass, returnChan)
	// goroutine loop listening for solutions sent out from workers on the LAN and
	// then submitting them to the network
	go func() {
		cancel := broadcast.Listen(broadcast.DefaultAddress, m.msgHandler)
		var submitLock sync.Mutex
	outReceive:
		for {
			select {
			case sol := <-returnChan:
				log.SPEW(sol)
				msgBlock := wire.MsgBlock(*sol)
				block := util.NewBlock(&msgBlock)
				submitLock.Lock()
				// Ensure the block is not stale since a new block could have shown up while
				// the solution was being found.  Typically that condition is detected and
				// all work on the stale block is halted to start work on a new block, but
				// the check only happens periodically, so it is possible a block was found
				// and submitted in between.
				if !msgBlock.Header.PrevBlock.IsEqual(&cx.RPCServer.Cfg.Chain.BestSnapshot().Hash) {
					pauseRebroadcast.Store(true)
					oldBlocks.Lock()
					log.WARNF(
						"Block submitted via kopach miner with previous block %s is stale",
						msgBlock.Header.PrevBlock)
					log.WARN("updating block templates")
					blocks := Blocks{New: true}
					// generate Blocks
					for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
						BestSnapshot().Height+1)].Algos {
						// Choose a payment address at random.
						rand.Seed(time.Now().UnixNano())
						payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
							StateCfg.ActiveMiningAddrs))]
						template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
							payToAddr, algo)
						if err != nil {
							log.ERROR("failed to create new block template:", err)
							continue
						}
						blocks.Templates = append(blocks.Templates, template)
					}
					oldBlocks.Templates = blocks.Templates
					blockChan <- blocks
					oldBlocks.Unlock()
					pauseRebroadcast.Toggle()
					continue
				}
				pauseRebroadcast.Store(true)
				oldBlocks.Lock()
				// Process this block using the same rules as blocks coming from other
				// nodes.  This will in turn relay it to the network like normal.
				isOrphan, err := cx.RealNode.SyncManager.ProcessBlock(block, blockchain.BFNone)
				if err != nil {
					// Anything other than a rule violation is an unexpected error, so log
					// that error as an internal error.
					if _, ok := err.(blockchain.RuleError); !ok {
						log.WARNF(
							"Unexpected error while processing block submitted via CPU miner:", err,
						)
						oldBlocks.Unlock()
						pauseRebroadcast.Toggle()
						continue
					}
					log.WARNF("block submitted via kopach miner rejected:", err)
					oldBlocks.Unlock()
					pauseRebroadcast.Toggle()
					continue
				}
				if isOrphan {
					log.WARN("block is an orphan")
					oldBlocks.Unlock()
					pauseRebroadcast.Toggle()
					continue
				}
				log.TRACE("the block was accepted")
				coinbaseTx := block.MsgBlock().Transactions[0].TxOut[0]
				prevHeight := block.Height() - 1
				prevBlock, _ := cx.RealNode.Chain.BlockByHeight(prevHeight)
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
					fork.GetAlgoName(block.MsgBlock().Header.Version, block.Height()),
					since)
				submitLock.Unlock()
				oldBlocks.Unlock()
				pauseRebroadcast.Toggle()
			case <-ctx.Done():
				log.DEBUG("quitting on quit channel close")
				cancel()
				break outReceive
			}
		}
	}()
	return
}
