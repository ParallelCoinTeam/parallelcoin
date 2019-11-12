package kopach

import "C"
import (
	"context"
	"fmt"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"net"
	"sync"
	"time"
)

type msgBuffer struct {
	buffers    [][]byte
	first      time.Time
	decoded    bool
	superseded bool
}

var nilKey = make([]byte, 32)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	log.DEBUG("kopach miner starting")
	wg.Add(1)
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	var cancel context.CancelFunc
	var err error
	buffers := make(map[string]*msgBuffer)
	//var lastHash atomic.Uint64
out:
	for _, j := range controller.MCAddresses {
		i := j
		var mx sync.Mutex
		cancel, err = controller.Listen(i, func(a *net.UDPAddr, n int,
			b []byte) {
			mx.Lock()
			defer mx.Unlock()
			if n < 16 {
				log.ERROR("received short broadcast message")
				return
			}
			magic := string(b[12:16])
			if magic == string(controller.WorkMagic[:]) {
				nonce := string(b[:12])
				if bn, ok := buffers[nonce]; ok {
					if !bn.decoded {
						payload := b[16:n]
						newP := make([]byte, len(payload))
						copy(newP, payload)
						bn.buffers = append(bn.buffers, newP)
						if len(bn.buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							//log.SPEW(bn.buffers)
							cipherText, err = fec.Decode(bn.buffers)
							if err != nil {
								log.ERROR(err)
								return
							}
							//log.SPEW(cipherText)
							msg, err := ciph.Open(nil, []byte(nonce), cipherText,
								nil)
							if err != nil {
								log.ERROR(err)
								return
							}
							//log.SPEW(msg)
							bn.decoded = true
							mC := controller.LoadMinerContainer(msg)
							for i := range buffers {
								if i!=nonce {
									if buffers[i].superseded {
										//defer func(){
											log.DEBUGF("deleting buffer %x", i)
											delete(buffers, i)
										//}()
										// it will be deleted once this
										// function completes,
										// so make sure it isn't seen again
										// by this loop
										//buffers[i].superseded = false
									}
								}
							}
							//log.SPEW(mC.Data)
							//h64 := highwayhash.Sum64(mC.Data, nilKey)
							//if h64 == lastHash.Load() {
							//	fmt.Printf("received rebroadcast of block %d"+
							//		" %v\r", mC.GetNewHeight(), time.Now())
							//	return
							//}
							//lastHash.Store(h64)
							log.DEBUG(mC.GetIPs())
							log.DEBUG("P2PListenersPort", mC.GetP2PListenersPort())
							log.DEBUG("RPCListenersPort", mC.GetRPCListenersPort())
							log.DEBUG("ControllerListenerPort", mC.GetControllerListenerPort())
							log.DEBUG("NewHeight", mC.GetNewHeight())
							log.DEBUG(mC.GetPrevBlockHash())
							log.DEBUG(mC.GetBitses())
							log.SPEW(mC.GetTxs())
						}
					} else {
						for i := range buffers {
							if i != nonce {
								// superseded blocks can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded
								buffers[i].superseded = true
							}
						}
						fmt.Printf("received rebroadcast of %x %v\r", nonce,time.Now())
					}
				} else {
					buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(),
						false, false}
					payload := b[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					buffers[nonce].buffers = append(buffers[nonce].buffers,
						newP)
					//log.DEBUGF("%x", payload)
				}
				//log.DEBUGF("%v %v %012x %s", i, a, nonce, magic)

			}
		})
		if err != nil {
			continue
		}
		// we only need to start one of them, ipv6 is preferred
		if cancel != nil {
			log.DEBUG("listener started", i.IP, i.Port, i.Zone, i.String(),
				i.Network())
			select {
			case <-quit:
				log.DEBUG("kopach miner shutting down")
				cancel()
				break out
			}
		}
		log.ERROR("failed to start listener on", i)
	}
	wg.Done()
}
