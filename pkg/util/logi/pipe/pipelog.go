package main

import (
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	qu "github.com/p9c/pod/pkg/util/quit"
	"os"
	"time"
)

func main() {
	// var err error
	logi.L.SetLevel("trace", false, "pod")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>>")
	consume.Start(w)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>>")
	time.Sleep(time.Second * 5)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>>")
	consume.Kill(w)
	Debug("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>>")
	// time.Sleep(time.Second * 5)
	// Debug(interrupt.GoroutineDump())
	// if err = w.Wait(); Check(err) {
	// }
	// time.Sleep(time.Second * 3)
}
