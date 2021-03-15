package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/pod"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/l0k1"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "c4315cdf272d5131c25a2dd9b21ec34f9832a0bf"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-03-15T20:32:59+01:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v1.9.22+"
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
