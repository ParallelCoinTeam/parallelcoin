package controller

import (
	"context"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			log.DEBUG("context cancelled, killing controller")
			break
		}
	}()
	return
}
