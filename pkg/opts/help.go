package opts

import (
	"fmt"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/cmds"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/text"
	"os"
	"strings"
	"unicode/utf8"
)

// getHelp walks the command tree and gathers the options and creates a set of help functions for all commands and
// options in the set
func (c *Config) getHelp() {
	cm := cmds.Command{
		Name:        "help",
		Description: "prints information about how to use pod",
		Entrypoint:  helpFunction,
		Commands:    nil,
	}
	c.Commands = append(c.Commands, cm)
	return
}

func helpFunction(ifc interface{}) error {
	c := assertToConfig(ifc)
	var o string
	o+=fmt.Sprintf( "Parallelcoin Pod All-in-One Suite\n\n")
	o+=fmt.Sprintf( "Usage:\n\t%s [options] [commands] [command parameters]\n\n", os.Args[0])
	o+=fmt.Sprintf("Commands:\n" )
	for i := range c.Commands {
		oo := fmt.Sprintf("\t%s", c.Commands[i].Name)
		nrunes := utf8.RuneCountInString(oo)
		o += oo + fmt.Sprintf(strings.Repeat(" ", 9-nrunes)+"%s\n", c.Commands[i].Description)
	}
	o += fmt.Sprintf(
		"\nOptions:\n\tset values on options concatenated against the option keyword or separated with '='\n",
	)
	o += fmt.Sprintf("\teg: addcheckpoints=deadbeefcafe,someothercheckpoint AP127.0.0.1:11047\n")
	o += fmt.Sprintf("\tfor items that take multiple string values, you can repeat the option with further\n")
	o += fmt.Sprintf("\tinstances of the option or separate the items with (only) commas as the above example\n\n")
	c.ForEach(func(ifc opt.Option) bool {
		meta := ifc.GetMetadata()
		oo := fmt.Sprintf("\t%s %v", meta.Option, meta.Aliases)
		nrunes := utf8.RuneCountInString(oo)
		var def string
		switch ii:=ifc.(type) {
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
		o += oo + fmt.Sprintf(strings.Repeat(" ", 32-nrunes)+"%s, default: %s\n", meta.Description, def)
		return true
	},
	)
	o += fmt.Sprintf("\nadd the name of the command or option after 'help' in the commandline to " +
		"get more detail - eg: %s help upnp\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, o)
	return nil
}

func assertToConfig(ifc interface{}) (c *Config) {
	var ok bool
	if c, ok = ifc.(*Config); !ok {
		panic("wth")
	}
	return
}
