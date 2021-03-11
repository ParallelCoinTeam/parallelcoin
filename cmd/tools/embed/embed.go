// +build ignore

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var files, symbols []string
	for _, x := range os.Args[2:] {
		globs, _ := filepath.Glob(x)
		if globs != nil {
			files = append(files, globs...)
		}
	}
	for i := range files {
		out := files[i]
		out = strings.ReplaceAll(out, ".", "_")
		out = strings.ReplaceAll(out, "/", "")
		for ; strings.HasPrefix(out, "_"); out = strings.TrimPrefix(out, "_") {
		}
		out = strings.ToUpper(out[0:1]) + out[1:]
		symbols = append(symbols, out)
	}

	for i := range files {
		fh, _ := os.Create(files[i] + ".go")
		fmt.Fprint(fh, "// AUTOGENERATED by cmd/tools/embed/embed.go; do not edit.\n\n")

		fmt.Fprint(fh, "package "+os.Args[1]+"\n\n")
		b, e := ioutil.ReadFile(files[i])
		if e ==  nil {
			fmt.Fprint(fh, fmt.Sprintf(`var %s = func() string {
				s, _ := base64.StdEncoding.DecodeString("`, base64.StdEncoding.EncodeToString(b)))
			fmt.Fprint(fh, "\")\n")
		}
		fmt.Println("writing file", files[i]+".go")
		fh.Close()
	}
}
