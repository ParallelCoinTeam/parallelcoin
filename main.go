// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	_ "net/http/pprof"

	_ "github.com/p9c/pod/pkg"

	"github.com/p9c/pod/cmd"
)

func main() {
	cmd.Main()
}
