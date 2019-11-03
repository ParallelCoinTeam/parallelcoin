package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	gomod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	s := string(gomod)
	S := strings.Split(s, "\n")
	S = strings.Split(S[0], " ")
	args := `-X `
	args += S[1]
	args += `/app.Stamp=`
	args += time.Now().UTC().Format("v06.01.02.15")
	var dir string
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dir += "/bin/pod"
	verbose := ""
	if len(os.Args) > 1 {
		for i := range os.Args {
			if i == 0 {
				continue
			}
			if os.Args[i] == "-h" || os.Args[i] == "--help" {
				fmt.Println(`bld - builds and stamps builds with custom variables
usage: 	bld [-v] [-h]
	
	-v, --verbose
		prints compiler verbose output to stdout
	-h, --help
		show this help message`)
				os.Exit(0)
			}
			if os.Args[i] == "-v" || os.Args[i] == "--verbose" {
				verbose = "-v"
			}
		}
	}
	cmd := exec.Command("go", "build", "-o", dir, "-ldflags", args, verbose)
	if verbose != "" {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERR", err)
	}
}
