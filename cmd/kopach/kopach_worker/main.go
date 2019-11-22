package main

import (
	"net/rpc"

	"github.com/davecgh/go-spew/spew"

	"github.com/p9c/pod/cmd/kopach/worker"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func main() {
	log.L.SetLevel("trace", true)
	log.DEBUG("miner worker starting")
	w, conn := worker.New(mine, sem.NewSemaphore(1))
	interrupt.AddHandler(func(){
		close(w.Quit)
	})
	err := rpc.Register(w)
	if err != nil {
		log.DEBUG(err)
		return
	}
	go rpc.ServeConn(conn)
	<-w.Quit
	log.DEBUG("finished")
}

func mine(blk *wire.MsgBlock, sem sem.T) {
out:
	for {
		log.DEBUG("mining on new block\n", spew.Sdump(blk))
		select {
		case <-sem.Release():
			break out
		}
	}
}
