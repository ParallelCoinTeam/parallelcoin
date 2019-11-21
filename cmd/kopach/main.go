package kopach

import (
	"context"
	"crypto/cipher"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"go.uber.org/atomic"
	"net"
	"os"
	"sync"
	"time"
)

type Worker struct {
	active      *atomic.Bool
	buffers     map[string]*controller.MsgBuffer
	ciph        cipher.AEAD
	conn        *net.UDPConn
	ctx         context.Context
	cx          *conte.Xt
	mx          *sync.Mutex
	receiveChan chan []byte
}

func Main(cx *conte.Xt, quit chan struct{}) {
	log.DEBUG("miner controller starting")
	ctx, cancel := context.WithCancel(context.Background())
	// resolve listen address
	conn, err := net.ListenUDP("udp", controller.MCAddress)
	if err != nil {
		log.ERROR(err)
		os.Exit(1)
	}
	wrkr := &Worker{
		active:      &atomic.Bool{},
		buffers:     make(map[string]*controller.MsgBuffer),
		ciph:        gcm.GetCipher(*cx.Config.MinerPass),
		conn:        conn,
		ctx:         ctx,
		cx:          cx,
		mx:          &sync.Mutex{},
		receiveChan: make(chan []byte),
	}
	// start up listener
	var stopListener context.CancelFunc
	stopListener, err = controller.Listen(conn, getMsgHandler(wrkr))
	if err != nil {
		log.DEBUG(err)
		return
	}
out:
	for {
		select {
		case <-ctx.Done():
			stopListener()
		case <-quit:
			cancel()
			break out
		}
	}
}

func getMsgHandler(c *Worker) func(a *net.UDPAddr, n int, b []byte) {
	return func(src *net.UDPAddr, n int, buffer []byte) {
		var err error
		if n < 16 {
			log.ERROR("received short broadcast message")
			return
		}
		magic := string(buffer[12:16])
		if magic == string(job.WorkMagic[:]) {
			nonce := string(buffer[:12])
			if bn, ok := c.buffers[nonce]; ok {
				if !bn.Decoded {
					payload := buffer[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					bn.Buffers = append(bn.Buffers, newP)
					if len(bn.Buffers) >= 3 {
						// try to decode it
						var cipherText []byte
						log.SPEW(bn.Buffers)
						cipherText, err = fec.Decode(bn.Buffers)
						if err != nil {
							log.ERROR(err)
							return
						}
						log.SPEW(cipherText)
						msg, err := c.ciph.Open(nil, []byte(nonce),
							cipherText, nil)
						if err != nil {
							log.ERROR(err)
							return
						}
						bn.Decoded = true
						c.receiveChan <- msg
					}
				} else {
					for i := range c.buffers {
						if i != nonce {
							// superseded blocks can be deleted from the
							// buffers,
							// we don't add more data for the already
							// decoded
							c.buffers[i].Superseded = true
						}
					}
				}
			} else {
				c.buffers[nonce] = &controller.MsgBuffer{[][]byte{},
					time.Now(), false, false}
				payload := buffer[16:n]
				newP := make([]byte, len(payload))
				copy(newP, payload)
				c.buffers[nonce].Buffers = append(c.buffers[nonce].
					Buffers, newP)
				//log.DEBUGF("%x", payload)
			}
			//log.DEBUGF("%v %v %012x %s", i, src, nonce, magic)
		}
	}
}
