package main

import (
	"os"
	"time"
	
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func main() {
	// var err error
	logi.L.SetLevel("trace", false, "pod")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	consume.Start(w)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	consume.Kill(w)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// Debug(interrupt.GoroutineDump())
	// if err = w.Wait(); Check(err) {
	// }
	// time.Sleep(time.Second * 3)
}
