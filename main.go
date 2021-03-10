// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
// To regenerate the loggers run go generate in this folder, which should be
// done before release as this file can be edited to hide or emphasise log
// entries per package
//go:generate go run ./pkg/logg/deploy/.
//
package main

import (
	_ "net/http/pprof"
	
	"github.com/p9c/pod/cmd"
)

func main() {
	cmd.Main()
}
