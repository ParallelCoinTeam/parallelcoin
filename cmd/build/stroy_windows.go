package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
	
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
	if cwd, e = os.Getwd(); dbg.Chk(e) {
		return false
	}
	var repo *git.Repository
	if repo, e = git.PlainOpen(cwd); dbg.Chk(e) {
		return false
	}
	var rr []*git.Remote
	if rr, e = repo.Remotes(); dbg.Chk(e) {
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
	// if rbr, e = repo.Branch("l0k1"); dbg.Chk(e) {
	// }
	// var rbr storer.ReferenceIter
	// if rbr, e = repo.Branches(); dbg.Chk(e){
	// 	return false
	// }
	// spew.Dump(rbr)
	// if rl, e = repo.Log(&git.LogOptions{
	// 	From:     plumbing.Hash{},
	// 	Order:    0,
	// 	FileName: nil,
	// 	All:      false,
	// }); dbg.Chk(e) {
	// 	return false
	// }
	// if e = rl.ForEach(func(cmt *object.Commit) (e error) {
	// 	spew.Dump(cmt)
	// 	return nil
	// }); dbg.Chk(e) {
	// }
	var rh *plumbing.Reference
	if rh, e = repo.Head(); dbg.Chk(e) {
		return false
	}
	rhs := rh.Strings()
	GitRef = rhs[0]
	GitCommit = rhs[1]
	// fmt.Println(rhs)
	// var rhco *object.Commit
	// if rhco, e = repo.CommitObject(rh.Hash()); dbg.Chk(e) {
	// }
	// // var dateS string
	// rhcoS := rhco.String()
	// sS := strings.Split(rhcoS, "Date:")
	// sSs := strings.TrimSpace(strings.Split(sS[1], "\n")[0])
	// fmt.Println(sSs)
	// var ti time.Time
	// if ti, e = time.Parse("Mon Jan 02 15:04:05 2006 -0700", sSs); dbg.Chk(e) {
	// }
	// fmt.Printf("time %v\n", ti)
	// fmt.Println(sSs)
	// fmt.Println(dateS)
	// inf.Ln(rh.Type(), rh.Target(), rh.Strings(), rh.String(), rh.Name())
	// var rb storer.ReferenceIter
	// if rb, e = repo.Branches(); dbg.Chk(e) {
	// 	return false
	// }
	// if e = rb.ForEach(func(pr *plumbing.Reference) (e error) {
	// 	inf.Ln(pr.String(), pr.Hash(), pr.Name(), pr.Strings(), pr.Target(), pr.Type())
	// 	return nil
	// }); dbg.Chk(e) {
	// 	return false
	// }
	var rt storer.ReferenceIter
	if rt, e = repo.Tags(); dbg.Chk(e) {
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
			// if rcoh, e = repo.CommitObject(pr.Hash()); dbg.Chk(e) {
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
	); dbg.Chk(e) {
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
	var e error
	var ok bool
	var home string
	if runtime.GOOS == "windows" {
		var homedrive string
		if homedrive, ok = os.LookupEnv("HOMEDRIVE"); !ok {
			panic(err)
		}
		var homepath string
		if homepath, ok = os.LookupEnv("HOMEPATH"); !ok {
			panic(err)
		}
		home = homedrive + homepath
	} else {
		if home, ok = os.LookupEnv("HOME"); !ok {
			panic(err)
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
				// inf.Ln(list[i])
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
				// 		inf.F("'%s'", split[i])
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
				if e := cmd.Start(); dbg.Chk(e) {
					Infos(err)
					os.Exit(1)
				}
				if e := cmd.Wait(); dbg.Chk(e) {
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
