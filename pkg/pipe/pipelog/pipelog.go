package main

import (
	"github.com/p9c/log"
	"os"
	"time"
	
	"github.com/p9c/pod/pkg/pipe/consume"
	"github.com/p9c/qu"
)

func main() {
	// var e error
	log.SetLogLevel("trace")
	// command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := qu.T()
	// splitted := strings.Split(command, " ")
	splitted := os.Args[1:]
	w := consume.Log(quit, consume.SimpleLog(splitted[len(splitted)-1]), consume.FilterNone, splitted...)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> starting")
	consume.Start(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> started")
	time.Sleep(time.Second * 4)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopping")
	consume.Kill(w)
	D.Ln("\n\n>>> >>> >>> >>> >>> >>> >>> >>> >>> stopped")
	// time.Sleep(time.Second * 5)
	// D.Ln(interrupt.GoroutineDump())
	// if e = w.Wait(); E.Chk(e) {
	// }
	// time.Sleep(time.Second * 3)
}
