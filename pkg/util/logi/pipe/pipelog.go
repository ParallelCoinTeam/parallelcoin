package main

import (
	"os"
	"time"
	
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func main() {
	// var e error
	logi.L.SetLevel("trace", false, "pod")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	consume.Start(w)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	consume.Kill(w)
	dbg.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// dbg.Ln(interrupt.GoroutineDump())
	// if e = w.Wait(); dbg.Chk(e) {
	// }
	// time.Sleep(time.Second * 3)
}
