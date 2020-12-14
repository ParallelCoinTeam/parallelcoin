package main

import (
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/pipe/consume"
	"strings"
	"time"
)

func main() {
	// var err error
	logi.L.SetLevel("trace", false, "pod")
	command := "pod -D test0 -n testnet -l trace --solo --lan --pipelog node"
	quit := make(chan struct{})
	w := consume.Log(quit, consume.SimpleLog, consume.FilterNone, strings.Split(command, " ")...)
	consume.Start(w)
	time.Sleep(time.Second * 3)
	consume.Kill(w)
	// time.Sleep(time.Second * 3)
	// if err = w.Kill(); Check(err) {
	// }
}
