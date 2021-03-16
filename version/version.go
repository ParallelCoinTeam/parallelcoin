package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/pod"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/l0k1"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "3f0f307e7fe0384092e7a10034e62e88702a6eac"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-03-16T02:15:41+01:00"
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
