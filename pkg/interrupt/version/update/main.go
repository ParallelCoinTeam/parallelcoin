package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
)

func main() {
	BuildTime = time.Now().Format(time.RFC3339)
	var cwd string
	var e error
	if cwd, e = os.Getwd(); e != nil {
		return
	}
	var repo *git.Repository
	if repo, e = git.PlainOpen(cwd); e != nil {
		return
	}
	var rr []*git.Remote
	if rr, e = repo.Remotes(); e != nil {
		return
	}
	for i := range rr {
		rs := rr[i].String()
		if strings.HasPrefix(rs, "origin") {
			rss := strings.Split(rs, "git@")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = strings.ReplaceAll(rsss[0], ":", "/")
				break
			}
			rss = strings.Split(rs, "https://")
			if len(rss) > 1 {
				rsss := strings.Split(rss[1], ".git")
				URL = rsss[0]
				break
			}

		}
	}
	var rh *plumbing.Reference
	if rh, e = repo.Head(); e != nil {
		return
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	var rt storer.ReferenceIter
	if rt, e = repo.Tags(); e != nil {
		return
	}
	var maxVersion int
	var maxString string
	var maxIs bool
	if e = rt.ForEach(
		func(pr *plumbing.Reference) (e error) {
			prs := strings.Split(pr.String(), "/")[2]
			if strings.HasPrefix(prs, "v") {
				var va [3]int
				_, _ = fmt.Sscanf(prs, "v%d.%d.%d", &va[0], &va[1], &va[2])
				vn := va[0]*1000000 + va[1]*1000 + va[2]
				if maxVersion < vn {
					maxVersion = vn
					maxString = prs
				}
				if pr.Hash() == rh.Hash() {
					maxIs = true
				}
			}
			return nil
		},
	); e != nil {
		return
	}
	if !maxIs {
		maxString += "+"
	}
	Tag = maxString
	_, file, _, _ := runtime.Caller(0)
	// fmt.Fprintln(os.Stderr, "file", file)
	urlSplit := strings.Split(URL, "/")
	// fmt.Fprintln(os.Stderr, "urlSplit", urlSplit)
	baseFolder := urlSplit[len(urlSplit)-1]
	// fmt.Fprintln(os.Stderr, "baseFolder", baseFolder)
	splitPath := strings.Split(file, baseFolder)
	// fmt.Fprintln(os.Stderr, "splitPath", splitPath)
	PathBase := filepath.Join(splitPath[0], baseFolder) + string(filepath.Separator)
	PathBase = strings.ReplaceAll(PathBase, "\\", "\\\\")
	// fmt.Fprintln(os.Stderr, "PathBase", PathBase)
	versionFile := `package version

import "fmt"

var (

	// URL is the git URL for the repository
	URL = "%s"
	// GitRef is the gitref, as in refs/heads/branchname
	GitRef = "%s"
	// GitCommit is the commit hash of the current HEAD
	GitCommit = "%s"
	// BuildTime stores the time when the current binary was built
	BuildTime = "%s"
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "%s"
	// PathBase is the path base returned from runtime caller
	PathBase = "%s"
)

// Get returns a pretty printed version information string
func Get() string {
	return fmt.Sprint(
		"Repository Information\n"+
		"	git repository: "+URL+"\n",
		"	branch: "+GitRef+"\n"+
		"	commit: "+GitCommit+"\n"+
		"	built: "+BuildTime+"\n"+
		"	Tag: "+Tag+"\n",
	)
}
`
	versionFileOut := fmt.Sprintf(
		versionFile,
		URL,
		GitRef,
		GitCommit,
		BuildTime,
		Tag,
		PathBase,
	)
	if e = ioutil.WriteFile("version/version.go", []byte(versionFileOut), 0666); E.Chk(e) {
	}
	return
}
