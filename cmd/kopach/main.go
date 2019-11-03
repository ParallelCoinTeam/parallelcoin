package kopach

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"sync"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	log.DEBUG("kopach miner starting")
	select {
	case <-quit:
		log.DEBUG("kopach miner shutting down")
		break
	}
}
