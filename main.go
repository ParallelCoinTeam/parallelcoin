// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

//go:generate go run -tags generate gen.go

import (
	"github.com/p9c/pod/cmd"
	_ "net/http/pprof"
)

func main() {
	cmd.Main()
}
