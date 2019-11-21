package kopach

import "github.com/p9c/pod/pkg/conte"

func Main(cx *conte.Xt, quit chan struct{}) {
out:
	for {

		select {
		case <-quit:
			break out
		}
	}
}
