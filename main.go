// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	"github.com/p9c/pod/cmd"
	"github.com/p9c/pod/pkg/log"
)

func main() {
	log.L.SetLevel("trace", true)
	cmd.Main()
}
