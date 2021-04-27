package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var loggerSourceCode string
	var e error
	var b []byte
	if b, e = ioutil.ReadFile("./pkg/logg/deploy/template.go"); e != nil {
		panic(e)
	}
	loggerSourceCode = string(b)
	fmt.Fprintln(os.Stderr, loggerSourceCode)
	if e = filepath.Walk(
		".",
		func(path string, info os.FileInfo, E error) (e error) {
			targetFilename := string(filepath.Separator) + "log.go"
			if strings.HasSuffix(path, targetFilename) {
				pkgName := strings.Split(path, targetFilename)[0]
				split := strings.Split(pkgName, string(filepath.Separator))
				if e := ioutil.WriteFile(path, []byte(
						strings.Replace(loggerSourceCode, "package main", "package "+split[len(split)-1], 1),
					), 0666); e != nil {
					fmt.Fprintln(os.Stderr, e.Error())
				}
			}
			return nil
		},
	); e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
	}
}
