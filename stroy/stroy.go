package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type command struct {
	name string
	args []string
}

var commands = map[string][]string{
	"build": {
		"go build -v",
	},
	"windows": {
		`go build -v -ldflags="-H windowsgui"`,
	},
	"tests": {
		`go test ./...`,
	},
	"kopachgui": {
		"go install -v",
		"pod -D %datadir -n testnet -l debug --lan --solo --kopachgui kopach",
	},
	"testkopach": {
		"go install -v",
		"pod -D %datadir -n testnet -l trace -g -G 1 --lan kopach",
	},
	"testnode": {
		"go install -v",
		"pod -D %datadir -n testnet -l debug --solo --lan node",
	},
	"nodegui": {
		"go install -v",
		"pod -D %datadir -n testnet nodegui",
	},
	"gui": {
		"go install -v",
		"pod -D %datadir -n testnet --lan",
	},
	"guis": {
		"go install -v",
		"pod -D test1 --minerpass pa55word",
	},
	"guass": {
		"go install -v",
		"pod -D %datadir --minerpass pa55word",
	},
	"resetwallet0": {
		"pod -D %datadir -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"resetwallet1": {
		"pod -D test1 -l trace --walletpass aoeuaoeu wallet drophistory",
	},
	"guihttpprof": {
		"go install -v",
		"pod -D %datadir -n testnet --lan --solo --kopachgui --profile 6969",
	},
	"guiprof": {
		"go install -v",
		"pod -D %datadir -n testnet --lan --solo --kopachgui",
	},
	"mainnode": {
		"go install -v",
		"pod -D testmain -n mainnet -l info --connect seed3.parallelcoin." +
			"io:11047 node",
	},
	"testwallet": {
		"go install -v",
		"pod -D %datadir -n testnet -l trace --walletpass aoeuaoeu wallet",
	},
	"mainwallet": {
		"go install -v",
		"pod -D testmain -n mainnet -l trace wallet",
	},
	"teststopkopach": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan kopach",
	},
	"teststopnode": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan node",
	},
	"teststopwallet": {
		"go install -v",
		"go install -v ./pkg/util/logi/pipe",
		"pipe pod -D %datadir --pipelog -l trace --walletpass aoeuaoeu -g -G 1" +
			" --solo --lan wallet",
	},
}

func main() {
	var err error
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
		if len(os.Args) > 2 {
			folderName = os.Args[2]
		}
		datadir := filepath.Join(home, folderName)
		if list, ok := commands[os.Args[1]]; ok {
			for i := range list {
				fmt.Println("executing item", i, "of list", os.Args[1], list[i])
				out := strings.ReplaceAll(list[i], "%datadir", datadir)
				split := strings.Split(out, " ")
				cmd := exec.Command(split[0], split[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
				}
			}
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
		fmt.Println("adding a second string to the commandline changes the name" +
			" of the home folder selected in the scripts")
	}
}
