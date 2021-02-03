// +build !windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/app/apputil"
	
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

func populateVersionFlags() bool {
	// `-X 'package_path.variable_name=new_value'`
	BuildTime = time.Now().Format(time.RFC3339)
	var cwd string
	var err error
	if cwd, err = os.Getwd(); Check(err) {
		return false
	}
	var repo *git.Repository
	if repo, err = git.PlainOpen(cwd); Check(err) {
		return false
	}
	var rr []*git.Remote
	if rr, err = repo.Remotes(); Check(err) {
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
	// if rbr, err = repo.Branch("l0k1"); Check(err) {
	// }
	// var rbr storer.ReferenceIter
	// if rbr, err = repo.Branches(); Check(err){
	// 	return false
	// }
	// spew.Dump(rbr)
	// if rl, err = repo.Log(&git.LogOptions{
	// 	From:     plumbing.Hash{},
	// 	Order:    0,
	// 	FileName: nil,
	// 	All:      false,
	// }); Check(err) {
	// 	return false
	// }
	// if err = rl.ForEach(func(cmt *object.Commit) error {
	// 	spew.Dump(cmt)
	// 	return nil
	// }); Check(err) {
	// }
	var rh *plumbing.Reference
	if rh, err = repo.Head(); Check(err) {
		return false
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	// fmt.Println(rhs)
	// var rhco *object.Commit
	// if rhco, err = repo.CommitObject(rh.Hash()); Check(err) {
	// }
	// // var dateS string
	// rhcoS := rhco.String()
	// sS := strings.Split(rhcoS, "Date:")
	// sSs := strings.TrimSpace(strings.Split(sS[1], "\n")[0])
	// fmt.Println(sSs)
	// var ti time.Time
	// if ti, err = time.Parse("Mon Jan 02 15:04:05 2006 -0700", sSs); Check(err) {
	// }
	// fmt.Printf("time %v\n", ti)
	// fmt.Println(sSs)
	// fmt.Println(dateS)
	// Info(rh.Type(), rh.Target(), rh.Strings(), rh.String(), rh.Name())
	// var rb storer.ReferenceIter
	// if rb, err = repo.Branches(); Check(err) {
	// 	return false
	// }
	// if err = rb.ForEach(func(pr *plumbing.Reference) error {
	// 	Info(pr.String(), pr.Hash(), pr.Name(), pr.Strings(), pr.Target(), pr.Type())
	// 	return nil
	// }); Check(err) {
	// 	return false
	// }
	var rt storer.ReferenceIter
	if rt, err = repo.Tags(); Check(err) {
		return false
	}
	// latest := time.Time{}
	// biggest := ""
	// allTags := []string{}
	var maxVersion int
	var maxString string
	var maxIs bool
	if err = rt.ForEach(
		func(pr *plumbing.Reference) error {
			// var rcoh *object.Commit
			// if rcoh, err = repo.CommitObject(pr.Hash()); Check(err) {
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
	); Check(err) {
		return false
	}
	if !maxIs {
		maxString += "+"
	}
	// fmt.Println(maxVersion, maxString)
	Tag = maxString
	// sort.Ints(versionsI)
	// if runtime.GOOS == "windows" {
	
	ldFlags = []string{
		`"-X main.URL=` + URL + ``,
		`-X main.GitCommit=` + GitCommit + ``,
		`-X main.BuildTime=` + BuildTime + ``,
		`-X main.GitRef=` + GitRef + ``,
		`-X main.Tag=` + Tag + `"`,
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
	var err error
	var ok bool
	var home string
	if home, ok = os.LookupEnv("HOME"); !ok {
		panic(err)
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
				// Info(list[i])
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
				// 		Infof("'%s'", split[i])
				// }
				fmt.Printf(
					`executing item %d of list '%v' '%v' '%v'

`, i, os.Args[1],
					split[0], split[1:],
				)
				// Info(split)
				var cmd *exec.Cmd
				scriptPath := filepath.Join(appdata.Dir("stroy", false), "stroy.sh")
				apputil.EnsureDir(scriptPath)
				if err = ioutil.WriteFile(
					scriptPath,
					[]byte(strings.Join(split, " ")),
					0700,
				); Check(err) {
				} else {
					cmd = exec.Command("sh", scriptPath)
					cmd.Stdout = os.Stdout
					cmd.Stdin = os.Stdin
					cmd.Stderr = os.Stderr
				}
				if cmd == nil {
					panic("cmd is nil")
				}
				if err := cmd.Start(); Check(err) {
					Infos(err)
					os.Exit(1)
				}
				if err := cmd.Wait(); Check(err) {
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
		"app information: repo: %s branch: %s commit: %s built"+
			": %s tag: %s...\n", URL, GitRef, GitCommit, BuildTime, Tag,
	)
}

type command struct {
	name string
	args []string
}

var ldFlags []string

