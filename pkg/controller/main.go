package controller

import (
	"bytes"
	"context"
	"github.com/p9c/pod/cmd/node/rpc"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/ugorji/go/codec"
	"math/rand"
	"sync"
	"time"
)

// Blocks is a block broadcast message for miners to mine from
type Blocks struct {
	// New is a flag that distinguishes a newly accepted/connected block from a rebroadcast
	New bool
	// Payload is a map of bytes indexed by block version number
	Payload map[int32][]byte
}

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		log.DEBUG("miner controller starting")
		// shortcuts to necessary modules
		chainHandle := cx.RPCServer.Cfg.Chain
		generator := cx.RPCServer.Cfg.Generator
		// create cipher for decoding relevant packets
		ciph := gcm.GetCipher(*cx.Config.MinerPass)
		// create new multicast address
		outAddr, err := broadcast.New(*cx.Config.BroadcastAddress)
		if err != nil {
			log.ERROR(err)
			cancel()
			return
		}
		var subscriberMutex sync.Mutex
		byts := make([]byte, 0, broadcast.MaxDatagramSize)
		var mh codec.MsgpackHandle
		// create subscriber for new block
		cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.Notification) {
			subscriberMutex.Lock()
			enc := codec.NewEncoderBytes(&byts, &mh)
			switch n.Type {
			case chain.NTBlockConnected, chain.NTBlockAccepted:
				bh := chainHandle.BestSnapshot().Height
				hf := fork.GetCurrent(bh + 1)
				blocks := Blocks{New: true, Payload: make(map[int32][]byte)}
				// generate Blocks
				for algo := range fork.List[hf].Algos {
					// Choose a payment address at random.
					rand.Seed(time.Now().UnixNano())
					payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.StateCfg.ActiveMiningAddrs))]
					template, err := generator.NewBlockTemplate(0, payToAddr, algo)
					if err != nil {
						log.ERROR("failed to create new block template:", err)
						cancel()
						subscriberMutex.Unlock()
						return
					}
					blk := template.Block
					log.SPEW(blk)
					var blkBuf bytes.Buffer
					err = blk.BtcEncode(&blkBuf, rpc.MaxProtocolVersion, wire.LatestEncoding)
					if err != nil {
						log.ERROR(err)
						cancel()
						subscriberMutex.Unlock()
						return
					}
					blocks.Payload[blk.Header.Version] = blkBuf.Bytes()
				}
				// create buffer and load into msgpack codec
				err := enc.Encode(&blocks)
				if err != nil {
					log.ERROR(err)
					cancel()
					subscriberMutex.Unlock()
					return
				}
				log.SPEW(byts)
				err = broadcast.Send(outAddr, byts, ciph, broadcast.Template)
				if err != nil {
					log.ERROR(err)
					cancel()
					subscriberMutex.Unlock()
					return
				}
				// reset the bytes for next round
				byts = byts[:0]
				enc.ResetBytes(&byts)
			}
			subscriberMutex.Unlock()
		})
		select {
		case <-ctx.Done():
			log.DEBUG("miner controller shutting down")
			break
		}
	}()
	return
}
