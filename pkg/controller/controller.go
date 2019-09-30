package controller

import (
	"context"

	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
)

type Blocks []mining.BlockTemplate

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	var lastBlock Blocks
	ctx, cancel = context.WithCancel(context.Background())
	blockChan := make(chan Blocks)
	go func() {
		for {
			// work loop
			select {
			case lastBlock = <-blockChan:
				// send out block broadcast
			case <-ctx.Done():
				// cancel has been called
				return
			default:
			}
		}
	}()
	// create subscriber for new block event

	// goroutine loop checking for connection and sync status
	
	return
}
