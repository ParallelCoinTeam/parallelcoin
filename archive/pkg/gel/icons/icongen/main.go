package main

import (
	"go/format"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/p9c/gel/apputil"
)

func main() {
	src := getSourceCode("icons", getIcons())
	filename := "../icons/icons.go"
	apputil.EnsureDir(filename)
	var e error
	if e = ioutil.WriteFile(filename, src, 0600); E.Chk(e) {
		panic(e)
	}
}

func getIcons() (iconNames []string) {
	url := "https://raw.githubusercontent.com/golang/exp/54ebac48fca0f39f9b63e0112b50a168ee5b5c00/shiny/materialdesign/icons/data.go"
	var e error
	var r *http.Response
	if r, e = http.Get(url); E.Chk(e) {
		panic(e)
	}
	var b []byte
	if b, e = ioutil.ReadAll(r.Body); E.Chk(e) {
		panic(e)
	}
	if e = r.Body.Close(); E.Chk(e) {
		panic(e)
	}
	s := string(b)
	split := strings.Split(s, "var ")
	iconNames = make([]string, len(split))
	for i, x := range split[1:] {
		split2 := strings.Split(x, " ")
		iconNames[i] = split2[0]
	}
	return iconNames
}

func getSourceCode(packagename string, iconNames []string) []byte {
	o := `// Package icons bundles the entire set of several icon sets into one package as maps to allow iteration

`+`//go:generate go run ./icongen/.

package ` + packagename + `


import (
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var Material = map[string]*[]byte {
`
	for i := range iconNames {
		if iconNames[i] == "" {
			continue
		}
		o += "\t" + `"` + iconNames[i] + `": &icons.` + iconNames[i] + ",\n"
	}
	o += "}\n"
	// I.Ln(o)
	var e error
	var out []byte
	if out, e = format.Source([]byte(o)); e != nil {
		panic(e)
	}
	return out
}
