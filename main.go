//go:generate go run -tags generate gen.go
// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	"github.com/p9c/pod/cmd"
)

func main() {
	cmd.Main()
}