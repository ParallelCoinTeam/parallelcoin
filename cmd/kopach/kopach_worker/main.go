package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/cmd/kopach/worker"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/sem"
	"net/rpc"
	"os"
)

func main() {
	w := worker.New(mine, sem.NewSemaphore(1))
	err := rpc.Register(w)
	if err != nil {
		printlnE(err)
		return
	}
	go rpc.ServeConn(w)
	<-w.Quit
	printlnE("finished")
}

func mine(blk *wire.MsgBlock, sem sem.T) {
out:
	for {
		printlnE("mining on new block\n", spew.Sdump(blk))
		select {
		case <-sem.Release():
			break out
		//default:
		}
	}
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[Worker]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
}
