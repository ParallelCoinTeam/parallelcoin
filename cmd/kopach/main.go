package kopach

import (
	"crypto/cipher"
	"encoding/hex"
	"github.com/p9c/pod/pkg/broadcast"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/ugorji/go/codec"
	"net"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/conte"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.WARN("starting kopach standalone miner worker")
		m := newMsgHandle(*cx.Config.MinerPass)
	out:
		for {
			cancel := broadcast.Listen(broadcast.DefaultAddress, m.msgHandler)
			select {
			case <-quit:
				log.DEBUG("quitting on killswitch")
				cancel()
				break out
			}
		}
		wg.Done()
	}()
}

type msgBuffer struct {
	buffers [][]byte
	first   time.Time
	decoded bool
}

type msgHandle struct {
	buffers map[string]*msgBuffer
	ciph    *cipher.AEAD
	dec     *codec.Decoder
	decBuf  []byte
	blockChan chan controller.Blocks
}

func newMsgHandle(password string) (out *msgHandle) {
	out = &msgHandle{}
	out.buffers = make(map[string]*msgBuffer)
	ciph := gcm.GetCipher(password)
	out.ciph = &ciph
	var mh codec.MsgpackHandle
	out.decBuf = make([]byte, 0, broadcast.MaxDatagramSize)
	out.dec = codec.NewDecoderBytes(out.decBuf, &mh)
	return
}

func (m *msgHandle) msgHandler(src *net.UDPAddr, n int, b []byte) {
	// remove any expired message bundles in the cache
	var deleters []string
	for i := range m.buffers {
		if time.Now().Sub(m.buffers[i].first) > time.Millisecond*50 {
			deleters = append(deleters, i)
		}
	}
	for i := range deleters {
		log.TRACE("deleting old message buffer")
		delete(m.buffers, deleters[i])
	}
	b = b[:n]
	//log.SPEW(b)
	if n < 16 {
		log.ERROR("received short broadcast message")
		return
	}
	// snip off message magic bytes
	msgType := string(b[:8])
	b = b[8:]
	log.TRACE(n, " bytes read from ", src, "message type", msgType)
	if msgType == string(broadcast.Template) {
		log.TRACE("got block template shard")
		buffer := b
		nonce := string(b[:8])
		if x, ok := m.buffers[nonce]; ok {
			log.TRACE("additional shard with nonce", hex.EncodeToString([]byte(nonce)))
			if !x.decoded {
				log.TRACE("adding shard")
				x.buffers = append(x.buffers, buffer)
				lb := len(x.buffers)
				log.TRACE("have", lb, "buffers")
				if lb > 2 {
					// try to decode it
					bytes, err := broadcast.Decode(*m.ciph, x.buffers)
					if err != nil {
						log.ERROR(err)
						return
					}
					log.SPEW(bytes)
					m.dec.ResetBytes(bytes)
					blocks := controller.Blocks{}
					err = m.dec.Decode(&blocks)
					if err != nil {
						log.ERROR(err)
					}
					log.INFO(len(blocks), "block templates received")
					x.decoded = true
				}
			} else if x.buffers != nil {
				log.TRACE("nilling buffers")
				x.buffers = nil
			} else {
				log.TRACE("ignoring already decoded message shard")
			}
		} else {
			log.TRACE("adding nonce", hex.EncodeToString([]byte(nonce)))
			m.buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(), false}
			m.buffers[nonce].buffers = append(m.buffers[nonce].buffers, b)
		}
	}
}
