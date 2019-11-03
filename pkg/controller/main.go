package controller

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/p9c/pod/cmd/node/rpc"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"go.uber.org/atomic"
	"math/rand"
	"time"
)

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	var busy, active atomic.Bool
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		// send out initial block broadcast
		blks := getBlocks(cx, true)
		log.SPEW(blks)
		//log.SPEW(blks.Serialize())
		log.SPEW(DeserializeBlocks(blks.Serialize()))
		// There is no unsubscribe but we can use an atomic to disable the function instead
		// This also ensures that new work doesn't start once the context is cancelled below
		active.Store(true)
		log.DEBUG("miner controller starting")
		cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.Notification) {
			if !busy.Load() && active.Load() {
				// first to arrive locks out any others while processing
				busy.Store(true)
				log.DEBUG("received new chain notification")
				switch n.Type {
				case chain.NTBlockAccepted:
					blocks := getBlocks(cx, true)
					log.SPEW(blocks)
					//log.SPEW(blocks.Serialize())
					log.SPEW(DeserializeBlocks(blocks.Serialize()))
				}
				busy.Store(false)
			} else {
				// drop the job
				log.DEBUG("busy processing prior notification")
			}
		})
		select {
		case <-ctx.Done():
			log.DEBUG("miner controller shutting down")
			active.Store(false)
			break
		}
	}()
	return
}

type Blocks struct {
	New       uint16
	Count     uint16
	Lengths   []uint32
	Templates [][]byte
}

// Serialize efficiently converts the structure into bytes for sending on a HMAC and FEC protected transport
// thus no checks here, it might go boom!
func (b *Blocks) Serialize() (out []byte) {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint16(buffer, b.New)
	out = append(out, buffer[:2]...)
	binary.BigEndian.PutUint16(buffer, b.Count)
	out = append(out, buffer[:2]...)
	for i := range b.Lengths {
		binary.BigEndian.PutUint32(buffer, b.Lengths[i])
		out = append(out, buffer[:4]...)
	}
	for i := range b.Templates {
		out = append(out, b.Templates[i]...)
	}
	return out
}

// DeserializeBlocks takes a byte slice assumed to be a blocks message and decodes it into a Blocks struct
func DeserializeBlocks(blockBytes []byte) (b *Blocks) {
	b = &Blocks{Templates: [][]byte{}}
	b.New = binary.BigEndian.Uint16(blockBytes[:2])
	b.Count = binary.BigEndian.Uint16(blockBytes[2:4])
	counts := blockBytes[4 : b.Count*4+4]
	var totalLength uint32
	for i := uint16(0); i < b.Count; i++ {
		length := binary.BigEndian.Uint32(counts[i*4 : i*4+4])
		b.Lengths = append(b.Lengths, length)
		totalLength += length
	}
	blocks := blockBytes[b.Count*4+4 : (uint32(b.Count)*4+4)+totalLength]
	for i := range b.Lengths {
		b.Templates = append(b.Templates, blocks[:b.Lengths[i]])
		blocks = blocks[b.Lengths[i]:]
	}
	return
}

func getBlocks(cx *conte.Xt, isNew bool) *Blocks {
	var newi uint16
	if isNew {
		newi = ^newi
	}
	blocks := Blocks{New: newi, Templates: [][]byte{}}
	for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.BestSnapshot().Height+1)].Algos {
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		log.TRACE("len active mining addrs", len(cx.StateCfg.ActiveMiningAddrs))
		payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.StateCfg.ActiveMiningAddrs))]
		template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0, payToAddr, algo)
		if err != nil {
			log.ERROR("failed to create new block template:", err)
			continue
		}
		//log.SPEW(template.Block)
		var blkBuf bytes.Buffer
		err = template.Block.BtcEncode(&blkBuf, rpc.MaxProtocolVersion, wire.BaseEncoding)
		if err != nil {
			log.ERROR(err)
			return nil
		} else {
			blkBytes := blkBuf.Bytes()
			blocks.Templates = append(blocks.Templates, blkBytes)
			blocks.Lengths = append(blocks.Lengths, uint32(len(blkBytes)))
			blocks.Count++
		}
	}
	return &blocks
}
