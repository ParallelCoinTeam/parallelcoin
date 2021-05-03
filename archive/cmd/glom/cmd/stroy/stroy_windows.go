package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var WindowsExec = func(split []string) (out *exec.Cmd) {
	out = exec.Command(split[0])
	out.SysProcAttr = &syscall.SysProcAttr{}
	out.SysProcAttr.CmdLine = strings.Join(split, " ")
	return
}

func populateVersionFlags() bool {
	// `-X 'package_path.variable_name=new_value'`
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
	// spew.Dump(rr)
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
	// var rl object.CommitIter
	// var rbr *config.Branch
	// if rbr, e = repo.Branch("l0k1"); e != nil {
	// }
	// var rbr storer.ReferenceIter
	// if rbr, e = repo.Branches(); e != nil{
	// 	return false
	// }
	// spew.Dump(rbr)
	// if rl, e = repo.Log(&git.LogOptions{
	// 	From:     plumbing.Hash{},
	// 	Order:    0,
	// 	FileName: nil,
	// 	All:      false,
	// }); e != nil {
	// 	return false
	// }
	// if e = rl.ForEach(func(cmt *object.Commit) (e error) {
	// 	spew.Dump(cmt)
	// 	return nil
	// }); e != nil {
	// }
	var rh *plumbing.Reference
	if rh, e = repo.Head(); e != nil {
		return false
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	// fmt.Println(rhs)
	// var rhco *object.Commit
	// if rhco, e = repo.CommitObject(rh.Hash()); e != nil {
	// }
	// // var dateS string
	// rhcoS := rhco.String()
	// sS := strings.Split(rhcoS, "Date:")
	// sSs := strings.TrimSpace(strings.Split(sS[1], "\n")[0])
	// fmt.Println(sSs)
	// var ti time.Time
	// if ti, e = time.Parse("Mon Jan 02 15:04:05 2006 -0700", sSs); e != nil {
	// }
	// fmt.Printf("time %v\n", ti)
	// fmt.Println(sSs)
	// fmt.Println(dateS)
	// fmt.Fprintln(os.Stderr,rh.Type(), rh.Target(), rh.Strings(), rh.String(), rh.Name())
	// var rb storer.ReferenceIter
	// if rb, e = repo.Branches(); e != nil {
	// 	return false
	// }
	// if e = rb.ForEach(func(pr *plumbing.Reference) (e error) {
	// 	fmt.Fprintln(os.Stderr,pr.String(), pr.Hash(), pr.Name(), pr.Strings(), pr.Target(), pr.Type())
	// 	return nil
	// }); e != nil {
	// 	return false
	// }
	var rt storer.ReferenceIter
	if rt, e = repo.Tags(); e != nil {
		return false
	}
	// latest := time.Time{}
	// biggest := ""
	// allTags := []string{}
	var maxVersion int
	var maxString string
	var maxIs bool
	if e = rt.ForEach(
		func(pr *plumbing.Reference) (e error) {
			// var rcoh *object.Commit
			// if rcoh, e = repo.CommitObject(pr.Hash()); e != nil {
			// }
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
				// allTags = append(allTags, prs)
			}
			// fmt.Println(pr.String(), pr.Hash(), pr.Name(), pr.Strings(),
			// 	pr.Target(), pr.Type())
			return nil
		},
	); e != nil {
		return false
	}
	if !maxIs {
		maxString += "+"
	}
	// fmt.Println(maxVersion, maxString)
	Tag = maxString
	// txsort.Ints(versionsI)
	// if runtime.GOOS == "windows" {
	_, file, _, _ := runtime.Caller(0)
	fmt.Fprintln(os.Stderr, "file", file)
	urlSplit := strings.Split(URL, "/")
	fmt.Fprintln(os.Stderr, "urlSplit", urlSplit)
	baseFolder := urlSplit[len(urlSplit)-1]
	fmt.Fprintln(os.Stderr, "baseFolder", baseFolder)
	splitPath := strings.Split(file, baseFolder)
	fmt.Fprintln(os.Stderr, "splitPath", splitPath)
	PathBase := filepath.Join(splitPath[0], baseFolder) + string(filepath.Separator)
	PathBase = strings.ReplaceAll(PathBase, "\\", "/")
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
	// Tag lists the Tag on the build, adding a + to the newest Tag if the commit is
	// not that commit
	Tag = "%s"
	// PathBase is the path base returned from runtime caller
	PathBase = "%s"
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
	// } else {
	// 	ldFlags = []string{
	// 		`"-X 'main.URL=` + URL + ``,
	// 		`-X 'main.GitCommit=` + GitCommit + `'`,
	// 		`-X 'main.BuildTime=` + BuildTime + `'`,
	// 		`-X 'main.GitRef=` + GitRef + `'`,
	// 		`-X 'main.Tag=` + Tag + `'"`,
	// 	}
	// }
	
	// Infos(ldFlags)
	return true
}

func main() {
	fmt.Println(GetVersion())
	var e error
	var ok bool
	var home string
	if runtime.GOOS == "windows" {
		var homedrive string
		if homedrive, ok = os.LookupEnv("HOMEDRIVE"); !ok {
			panic(e)
		}
		var homepath string
		if homepath, ok = os.LookupEnv("HOMEPATH"); !ok {
			panic(e)
		}
		home = homedrive + homepath
	} else {
		if home, ok = os.LookupEnv("HOME"); !ok {
			panic(e)
		}
	}
	if len(os.Args) > 1 {
		folderName := "test0"
		var datadir string
		if len(os.Args) > 2 {
			datadir = os.Args[2]
		} else {
			datadir = filepath.Join(home, folderName)
		}
		if list, ok := commands[os.Args[1]]; ok {
			populateVersionFlags()
			// Infos(list)
			for i := range list {
				// fmt.Fprintln(os.Stderr,list[i])
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
				// Infos(split)
				// add ldflags to commands that have this
				// for i := range split {
				// 	split[i] =
				// 		fmt.Fprintf(os.Stderr,"'%s'", split[i])
				// }
				fmt.Printf(
					`executing item %d of list '%v' '%v' '%v'

`, i, os.Args[1],
					split[0], split[1:],
				)
				cmd := WindowsExec(split)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				if e := cmd.Start(); e != nil {
					spew.Fdump(os.Stderr, e)
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

var (
	URL       string
	GitRef    string
	GitCommit string
	BuildTime string
	Tag       string
)

func GetVersion() string {
	return fmt.Sprintf(
		"node information: repo: %s branch: %s commit: %s built"+
			": %s tag: %s...\n", URL, GitRef, GitCommit, BuildTime, Tag,
	)
}

type command struct {
	name string
	args []string
}

var ldFlags []string
