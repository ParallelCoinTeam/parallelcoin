package controller

import (
	"context"
	"math/rand"
	"sync/atomic"
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

type Blocks []*mining.BlockTemplate

// Run starts a controller instance
func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	log.WARN("starting controller")
	var mh codec.MsgpackHandle
	var ctx context.Context
	var lastBlock atomic.Value // Blocks
	ctx, cancel = context.WithCancel(context.Background())
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	outConn, err := broadcast.New(*cx.Config.BroadcastAddress)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	bytes := make([]byte, 0, broadcast.MaxDatagramSize)
	enc := codec.NewEncoderBytes(&bytes, &mh)

	blockChan := make(chan Blocks)
	go func() {
		for {
			// work loop
			select {
			case lb := <-blockChan:
				lastBlock.Store(lb)
				// send out block broadcast
				log.DEBUG("sending out block broadcast")
				// serialize blocks
				err := enc.Encode(lb)
				if err != nil {
					log.ERROR(err)
					break
				}
				err = broadcast.Send(outConn, bytes, ciph, broadcast.Template)
				if err != nil {
					log.ERROR(err)
				}
				// reset the bytes for next round
				bytes = bytes[:0]
			case <-ctx.Done():
				// cancel has been called
				return
				// default:
			}
		}
	}()
	connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
	// if out of sync or disconnected,
	// once a second send out empty blocks
	for (connCount < 1 && !*cx.Config.Solo) || !current {
		time.Sleep(time.Second)
		connCount = cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
		current = cx.RPCServer.Cfg.SyncMgr.IsCurrent()
		log.WARN("waiting for sync/peers", connCount, current)
		select {
		case <-ctx.Done():
			log.WARN("cancelled before initial connection/sync")
			return
		default:
		}
	}
	// load initial blocks and send them out
	blocks := Blocks{}
	lastBlock.Store(blocks)
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
		blocks = append(blocks, template)
	}
	lastBlock.Store(blocks)
	// create subscriber for new block event
	cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.
	Notification) {
		switch n.Type {
		case chain.NTBlockConnected:
			log.WARN("new block found")
			blocks := Blocks{}
			lastBlock.Store(blocks)
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
				blocks = append(blocks, template)
			}
			lastBlock.Store(blocks)
			blockChan <- lastBlock.Load().(Blocks)
		}
	})
	// goroutine loop checking for connection and sync status
	go func() {
		for {
			time.Sleep(time.Second)
			connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
			current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
			// if out of sync or disconnected,
			// once a second send out empty blocks
			if connCount < 1 || !current {
				lastBlock.Store(Blocks{})
				blockChan <- lastBlock.Load().(Blocks)
			}
			select {
			case <-ctx.Done():
				break
			}
		}
	}()
	return
}
