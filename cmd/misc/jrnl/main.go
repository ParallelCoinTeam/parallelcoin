// This is just a convenient cli command to automatically generate a new file
// for a journal entry with names based on unix timestamps
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type jrnlCfg struct {
	Root string
}

func printErrorAndDie(stuff ...interface{}) {
	fmt.Fprintln(os.Stderr, stuff)
	os.Exit(1)
}

func main() {
	var home string
	var e error
	if home, e = os.UserHomeDir(); e != nil {
		os.Exit(1)
	}
	var configFile []byte
	if configFile, e = ioutil.ReadFile(
		filepath.Join(home, ".jrnl")); e != nil {
		printErrorAndDie(e, "~/.jrnl configuration file not found")
	}
	var cfg jrnlCfg
	if e = json.Unmarshal(configFile, &cfg); e != nil {
		printErrorAndDie(e, "~/.jrnl config file did not unmarshal")
	}
	filename := filepath.Join(cfg.Root, fmt.Sprintf("jrnl%d.txt",
		time.Now().Unix()))
	if e = ioutil.WriteFile(filename,
		[]byte(time.Now().Format(time.RFC1123Z)+"\n\n"),
		0600,
	); e != nil {
		printErrorAndDie(e,
			"unable to create file (is your keybase filesystem mounted?")
		os.Exit(1)
	}
	exec.Command("gedit", filename).Run()
}
