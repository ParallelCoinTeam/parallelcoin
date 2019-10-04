package kopach

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/broadcast"
	"github.com/p9c/pod/pkg/log"
	"net"
	"sync"

	"github.com/p9c/pod/pkg/conte"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.WARN("starting kopach standalone miner worker")
	out:
		for {
			cancel := broadcast.Listen(broadcast.DefaultAddress, msgHandler)
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

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.INFO(n, " bytes read from ", src)
	spew.Dump(b[:n])
}