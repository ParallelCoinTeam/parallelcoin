// Package main is the root of the Parallelcoin Pod software suite
//
// It slices, it dices
//
package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	_ "net/http/pprof"
	
	"github.com/p9c/pod/version"
	
	"github.com/p9c/pod/cmd"
)

func main() {
	version.URL = url
	version.GitRef = gitRef
	version.GitCommit = gitCommit
	version.BuildTime = buildTime
	version.Tag = tag
	version.Get = GetVersion
	logg.SortSubsystemsList()
	cmd.Main()
}

var (
	// url is the git url for the repository
	url string
	// gitRef is the gitref, as in refs/heads/branchname
	gitRef string
	// gitCommit is the commit hash of the current HEAD
	gitCommit string
	// buildTime stores the time when the current binary was built
	buildTime string
	// tag lists the tag on the build, adding a + to the newest tag if the commit is
	// not that commit
	tag string
)

func GetVersion() string {
	return fmt.Sprintf(
		`ParallelCoin Pod
	repo: %s
	branch: %s
	commit: %s
	built: %s
	tag: %s
`,
		url, gitRef, gitCommit, buildTime, tag,
	)
}
