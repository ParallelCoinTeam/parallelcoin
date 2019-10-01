package kopach

import (
	"sync"

	"github.com/p9c/pod/pkg/conte"
)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	go func() {
		wg.Add(1)
	out:
		for {
			//

			select {
			case <-quit:
				break out
			default:
			}
		}
		wg.Done()
	}()

}
