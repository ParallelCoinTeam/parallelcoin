package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/pod"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/l0k1"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "b5f1cb2187a5d87e08c9fc6d52c11d84a461788a"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-03-29T15:30:32+02:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v1.9.25+"
	// PathBase is the path base returned from runtime caller
	PathBase = "/home/loki/src/github.com/p9c/pod/"
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"ParallelCoin Pod\n"+
		"	git repository: "+URL+"\n",
		"	branch: "+GitRef+"\n"+
		"	commit: "+GitCommit+"\n"+
		"	built: "+BuildTime+"\n"+
		"	Tag: "+Tag+"\n",
	)
}
