// This generator reads a podcfg.Configs map and spits out a podcfg.Config struct
package main

import (
	"fmt"
	"github.com/p9c/pod/pkg/podcfg"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

func main() {
	c := podcfg.GetConfigs()
	var o string
	var cc podcfg.ConfigSlice
	for i := range c {
		cc = append(cc, podcfg.ConfigSliceElement{Opt: c[i], Name: i})
	}
	sort.Sort(cc)
	for i := range cc {
		t := reflect.TypeOf(cc[i].Opt).String()
		split := strings.Split(t, "podcfg.")[1]
		o += fmt.Sprintf("\t%s\t*%s\n", cc[i].Name, split)
	}
	var e error
	var out []byte
	var wd string
	generated := fmt.Sprintf(configBase, o)
	if out, e = format.Source([]byte(generated)); e != nil {
		// panic(e)
		fmt.Println(e)
	}
	if wd, e = os.Getwd(); e != nil {
		// panic(e)
	}
	// fmt.Println(string(out), wd)
	if e = ioutil.WriteFile(filepath.Join(wd, "struct.go"), out, 0660); e != nil {
		panic(e)
	}
}

var configBase = `package podcfg

// Config defines the configuration items used by pod along with the various components included in the suite
//go:generate go run gen/main.go
type Config struct {
	// ShowAll is a flag to make the json encoder explicitly define all fields and not just the ones different to the
	// defaults
	ShowAll bool
	// Map is the same data but addressible using its name as found inside the various configuration types, the key is
	// converted to lower case for CLI args
	Map            map[string]Option
	Commands       Commands
	RunningCommand *Command
%s}
`
