package podhelp

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/p9c/matrjoska/pod/state"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/text"
)

func HelpFunction(ifc interface{}) error {
	ps := assertToPodState(ifc)
	c := ps.Config
	var o string
	o += fmt.Sprintf("Parallelcoin Pod All-in-One Suite\n\n")
	o += fmt.Sprintf("Usage:\n\t%s [options] [commands] [command parameters]\n\n", os.Args[0])
	o += fmt.Sprintf("Commands:\n")
	for i := range c.Commands {
		oo := fmt.Sprintf("\t%s", c.Commands[i].Name)
		nrunes := utf8.RuneCountInString(oo)
		o += oo + fmt.Sprintf(strings.Repeat(" ", 9-nrunes)+"%s\n", c.Commands[i].Title)
	}
	o += fmt.Sprintf(
		"\nOptions:\n\tset values on options concatenated against the option keyword or separated with '='\n",
	)
	o += fmt.Sprintf("\teg: addcheckpoints=deadbeefcafe,someothercheckpoint AP127.0.0.1:11047\n")
	o += fmt.Sprintf("\tfor items that take multiple string values, you can repeat the option with further\n")
	o += fmt.Sprintf("\tinstances of the option or separate the items with (only) commas as the above example\n\n")
	// items := make(map[string][]opt.Option)
	descs := make(map[string]string)
	c.ForEach(func(ifc opt.Option) bool {
		meta := ifc.GetMetadata()
		oo := fmt.Sprintf("\t%s %v", meta.Option, meta.Aliases)
		nrunes := utf8.RuneCountInString(oo)
		var def string
		switch ii := ifc.(type) {
		case *binary.Opt:
			def = fmt.Sprint(ii.Def)
		case *list.Opt:
			def = fmt.Sprint(ii.Def)
		case *float.Opt:
			def = fmt.Sprint(ii.Def)
		case *integer.Opt:
			def = fmt.Sprint(ii.Def)
		case *text.Opt:
			def = fmt.Sprint(ii.Def)
		case *duration.Opt:
			def = fmt.Sprint(ii.Def)
		}
		descs[meta.Group] += oo + fmt.Sprintf(strings.Repeat(" ", 32-nrunes)+"%s, default: %s\n", meta.Description, def)
		return true
	},
	)
	var cats []string
	for i := range descs {
		cats = append(cats, i)
	}
	// I.S(cats)
	sort.Strings(cats)
	for i := range cats {
		if cats[i] != "" {
			o += "\n" + cats[i] + "\n"
		}
		o += descs[cats[i]]
	}
	// for i := range cats {
	// }
	o += fmt.Sprintf("\nadd the name of the command or option after 'help' or append it after "+
		"'help' in the commandline to get more detail - eg: %s help upnp\n\n", os.Args[0],
	)
	fmt.Fprintf(os.Stderr, o)
	return nil
}

func assertToPodState(ifc interface{}) (c *state.State) {
	var ok bool
	if c, ok = ifc.(*state.State); !ok {
		panic("wth")
	}
	return
}
