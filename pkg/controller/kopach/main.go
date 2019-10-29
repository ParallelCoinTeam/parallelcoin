package kopach

import (
	"crypto/cipher"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
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
}

func newMsgHandle(password string) (out *msgHandle) {
	out = &msgHandle{}
	out.buffers = make(map[string]*msgBuffer)
	ciph := gcm.GetCipher(password)
	out.ciph = &ciph
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
		log.WARN("deleting old message buffer")
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
	log.INFO(n, " bytes read from ", src, "message type", msgType == string(broadcast.Template))
	if msgType == string(broadcast.Template) {
		log.WARN("got block template shard")
		buffer := b
		nonce := string(b[:8])
		if x, ok := m.buffers[nonce]; ok {
			log.WARN("additional shard with nonce", nonce)
			if !x.decoded {
				log.WARN("adding shard")
				x.buffers = append(x.buffers, buffer)
				lb := len(x.buffers)
				log.WARN("have",lb, "buffers")
				if lb > 2 {
					// try to decode it
					//spew.Dump(x.buffers)
					//fmt.Println()
					bytes, err := broadcast.Decode(*m.ciph, x.buffers)
					if err != nil {
						log.ERROR(err)
						return
					}
					log.WARN(bytes)
					x.decoded = true
				}
			} else if x.buffers != nil {
				log.WARN("nilling buffers")
				x.buffers = nil
			} else {
				log.WARN("ignoring already decoded message shard")
			}
		} else {
			log.WARN("adding nonce", nonce)
			m.buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(), false}
			m.buffers[nonce].buffers = append(m.buffers[nonce].buffers, b)
		}
	}
}
