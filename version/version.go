package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "github.com/p9c/pod"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "refs/heads/l0k1"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "b09f236b11c3ba65a73860d0aff5614f56592c5e"
	// BuildTime stores the time when the current binary was built
	BuildTime = "2021-03-11T00:16:58+01:00"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "v1.9.16+"
	// PathBase is the path base returned from runtime caller
	PathBase = "C:/Users/loki/GolandProjects/pod/"
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"ParallelCoin Pod\n"+
		"	URL: "+URL+"\n",
		"	branch: "+GitRef+"\n"+
		"	commit: "+GitCommit+"\n"+
		"	built: "+BuildTime+"\n"+
		"	Tag: "+Tag+"\n",
	)
}
