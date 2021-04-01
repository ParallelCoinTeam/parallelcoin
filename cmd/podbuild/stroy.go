// +build !windows

package main

import (
	"fmt"
	_ "github.com/p9c/pod/pkg/fork"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	
	"github.com/p9c/pod/pkg/appdata"
	"github.com/p9c/pod/pkg/apputil"
	
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

type command struct {
	name string
	args []string
}

var ldFlags []string

func main() {
	var e error
	var ok bool
	var home string
	var list []string
	if home, ok = os.LookupEnv("HOME"); !ok {
		panic(e)
	}
	if len(os.Args) > 1 {
		folderName := "test0"
		var datadir string
		if len(os.Args) > 2 {
			datadir = os.Args[2]
		} else {
			datadir = filepath.Join(home, folderName)
		}
		if list, ok = commands[os.Args[1]]; ok {
			writeVersionFile()
			for i := range list {
				// inject the data directory
				var split []string
				out := strings.ReplaceAll(list[i], "%datadir", datadir)
				split = strings.Split(out, " ")
				for i := range split {
					split[i] = strings.ReplaceAll(
						split[i], "%ldflags",
						fmt.Sprintf(
							`-ldflags=%s`, strings.Join(
								ldFlags,
								" ",
							),
						),
					)
				}
				fmt.Printf("executing item %d of list '%v' '%v' '%v'",
					i, os.Args[1], split[0], split[1:],
				)
				var cmd *exec.Cmd
				scriptPath := filepath.Join(appdata.Dir("stroy", false), "stroy.sh")
				apputil.EnsureDir(scriptPath)
				if e = ioutil.WriteFile(
					scriptPath,
					[]byte(strings.Join(split, " ")),
					0700,
				); e != nil {
				} else {
					cmd = exec.Command("sh", scriptPath)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
				}
				if cmd == nil {
					panic("cmd is nil")
				}
				var e error
				if e = cmd.Start(); e != nil {
					fmt.Fprintln(os.Stderr, e)
					os.Exit(1)
				}
				if e := cmd.Wait(); e != nil {
					os.Exit(1)
				}
			}
		} else {
			fmt.Println("command", os.Args[1], "not found")
		}
	} else {
		fmt.Println("no command requested, available:")
		for i := range commands {
			fmt.Println(i)
			for j := range commands[i] {
				fmt.Println("\t" + commands[i][j])
			}
		}
		fmt.Println()
		fmt.Println(
			"adding a second string to the commandline changes the name" +
				" of the home folder selected in the scripts",
		)
	}
}

func writeVersionFile() bool {
	BuildTime = time.Now().Format(time.RFC3339)
	var cwd string
	var e error
	if cwd, e = os.Getwd(); e != nil {
		return false
	}
	var repo *git.Repository
	if repo, e = git.PlainOpen(cwd); e != nil {
		return false
	}
	var rr []*git.Remote
	if rr, e = repo.Remotes(); e != nil {
		return false
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
		return false
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	var rt storer.ReferenceIter
	if rt, e = repo.Tags(); e != nil {
		return false
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
		return false
	}
	if !maxIs {
		maxString += "+"
	}
	Tag = maxString
	_, file, _, _ := runtime.Caller(0)
	fmt.Fprintln(os.Stderr, "file", file)
	urlSplit := strings.Split(URL, "/")
	fmt.Fprintln(os.Stderr, "urlSplit", urlSplit)
	baseFolder := urlSplit[len(urlSplit)-1]
	fmt.Fprintln(os.Stderr, "baseFolder", baseFolder)
	splitPath := strings.Split(file, baseFolder)
	fmt.Fprintln(os.Stderr, "splitPath", splitPath)
	PathBase := filepath.Join(splitPath[0], baseFolder) + string(filepath.Separator)
	PathBase = strings.ReplaceAll(PathBase, "\\", "\\\\")
	fmt.Fprintln(os.Stderr, "PathBase", PathBase)
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
	// Tag lists the Tag on the podbuild, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "%s"
	// PathBase is the path base returned from runtime caller
	PathBase = "%s"
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
	if e = ioutil.WriteFile("version/version.go", []byte(versionFileOut), 0666); e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
	return true
}

func GetVersion() string {
	return fmt.Sprintf(
		"app information: repo: %s branch: %s commit: %s built"+
			": %s tag: %s...\n", URL, GitRef, GitCommit, BuildTime, Tag,
	)
}
