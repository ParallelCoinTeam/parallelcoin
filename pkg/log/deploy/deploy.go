package main

import (
	"github.com/p9c/log/version"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var loggerSourceCode string
	var e error
	var b []byte
	if b, e = ioutil.ReadFile("./deploy/template.go"); e != nil {
		panic(e)
	}
	loggerSourceCode = string(b)
	I.Ln(loggerSourceCode)
	if e = filepath.Walk(
		".",
		func(path string, info os.FileInfo, E error) (e error) {
			targetFilename := string(filepath.Separator) + "log.go"
			if strings.HasSuffix(path, targetFilename) {
				pkgName := strings.Split(path, targetFilename)[0]
				split := strings.Split(pkgName, string(filepath.Separator))
				_ = "github.com/p9c/log"
				editedFile := strings.Replace(loggerSourceCode, "package main", "package "+split[len(split)-1], 1)
				editedFile = strings.Replace(loggerSourceCode, "github.com/p9c/log", version.URL, 1)
				if e := ioutil.WriteFile(path, []byte(
					editedFile,
				), 0666,
				); e != nil {
					I.Ln(e.Error())
				}
			}
			return nil
		},
	); E.Chk(e) {
	}
}
