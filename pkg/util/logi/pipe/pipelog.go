package main

import (
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	qu "github.com/p9c/pod/pkg/util/quit"
	"strings"
	"time"
)

func main() {
	var err error
	logi.L.SetLevel("trace", false, "pod")
	command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := make(qu.C)
	w := consume.Log(quit, consume.SimpleLog("node"), consume.FilterNone, strings.Split(command, " ")...)
	interrupt.AddHandler(
		func() {
			if err = w.Kill(); Check(err) {
			}
		},
	)
	consume.Start(w)
	time.Sleep(time.Second * 5)
	consume.Kill(w)
	// time.Sleep(time.Second * 5)
	// Debug(interrupt.GoroutineDump())
	// if err = w.Wait(); Check(err) {
	// }
	// time.Sleep(time.Second * 3)
}
