package main

import (
	"github.com/parallelcointeam/parallelcoin/pkg/log"
)

func main() {
	l := log.NewLogger("warn")
	_ = l
	// time.Sleep(time.Second * 5)
}
