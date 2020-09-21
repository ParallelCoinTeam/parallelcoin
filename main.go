// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	_ "net/http/pprof"

	"github.com/stalker-loki/pod/cmd"
)

func main() {
	cmd.Main()
}
