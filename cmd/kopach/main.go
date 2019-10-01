package kopach

import (
	"sync"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	cancel := controller.Run(cx)
	go func() {
		wg.Add(1)
	out:
		for {
			select {
			case <-quit:
				cancel()
				break out
			}
		}
		wg.Done()
	}()

}
