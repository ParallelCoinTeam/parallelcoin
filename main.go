// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	_ "github.com/p9c/pod/pkg"
	_ "net/http/pprof"

	"github.com/p9c/pod/cmd"
)

func main() {
	cmd.Main()
}
