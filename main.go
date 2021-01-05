// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	
	_ "github.com/p9c/pod/pkg"
	
	"github.com/p9c/pod/cmd"
)

func main() {
	Print()
	cmd.Main()
}

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
)

// Print the version stored in the version library
func Print() {
	_, _ = fmt.Fprintln(os.Stderr, sprintVersion())
}

func sprintVersion() string {
	return fmt.Sprintf("%s %s %s %s %s", URL, GitRef, GitCommit, BuildTime, Tag)
}
