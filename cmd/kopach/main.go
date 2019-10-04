package kopach

import (
	"github.com/p9c/pod/pkg/log"
	"sync"

	"github.com/p9c/pod/pkg/conte"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		log.WARN("starting kopach standalone miner worker")
	out:
		for {
			select {
			case <-quit:
				log.DEBUG("quitting on killswitch")
				break out
			}
		}
		wg.Done()
	}()
}
