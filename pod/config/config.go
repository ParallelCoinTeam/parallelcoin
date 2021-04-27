// Package podopts is a configuration system to fit with the all-in-one philosophy guiding the design of the parallelcoin
// pod.
//
// The configuration is stored by each component of the connected applications, so all data is stored in concurrent-safe
// atomics, and there is a facility to invoke a function in response to a new value written into a field by other
// threads.
//
// There is a custom JSON marshal/unmarshal for each field type and for the whole configuration that only saves values
// that differ from the defaults, similar to 'omitempty' in struct tags but where 'empty' is the default value instead
// of the default zero created by Go's memory allocator. This enables easy compositing of multiple sources.
//
package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"lukechampine.com/blake3"

	"github.com/p9c/log"
	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/cmds"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/text"
)

// Configs is the source location for the Config items, which is used to generate the Config struct
type Configs map[string]opt.Option
type ConfigSliceElement struct {
	Opt  opt.Option
	Name string
}
type ConfigSlice []ConfigSliceElement

func (c ConfigSlice) Len() int           { return len(c) }
func (c ConfigSlice) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c ConfigSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

// Initialize loads in configuration from disk and from environment on top of the default base
func (c *Config) Initialize(hf func(ifc interface{}) error) (e error) {
	// the several places configuration is sourced from are overlaid in the following order:
	// default -> config file -> environment variables -> commandline flags
	T.Ln("initializing configuration...")
	// first lint the configuration
	var aos map[string][]string
	if aos, e = getAllOptionStrings(c); E.Chk(e) {
		return
	}
	// this function will panic if there is potential for ambiguity in the commandline configuration args.
	T.Ln("linting configuration items")
	if _, e = findConflictingItems(aos); E.Chk(e) {
	}
	// generate and add the help commands to the help tree
	c.GetHelp(hf)
	// process the commandline
	T.Ln("processing commandline arguments", os.Args[1:])
	var cm *cmds.Command
	var options []opt.Option
	var optVals []string
	if c.ExtraArgs, cm, options, optVals, e = c.processCommandlineArgs(os.Args[1:]); E.Chk(e) {
		return
	}
	// I.S(options)
	if cm != nil {
		c.RunningCommand = *cm
	}
	// if the user sets the configfile directly, or the datadir on the commandline we need to load it from that path
	T.Ln("checking from where to load the configuration file")
	datadir := c.DataDir.V()
	var configPath string
	for i := range options {
		if options[i].Name() == "configfile" {
			if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
				configPath = optVals[i]
			}
		}
		if options[i].Name() == "datadir" {
			I.Ln("datadir was set", optVals[i])
			if _, e = options[i].ReadInput(optVals[i]); !E.Chk(e) {
				datadir = options[i].Type().(*text.Opt).V()
				// reset all defaults that base on the datadir to apply hereafter
				// if the value is default, update it to the new datadir, and update the default field, otherwise assume
				// it has been set in the commandline args, and if it is different in environment or config file
				// it will be loaded in next, with the command line option value overriding at the end.
				if c.CAFile.V() == c.CAFile.Def {
					e = c.CAFile.Set(filepath.Join(datadir, "ca.cert"))
				}
				c.CAFile.Def = filepath.Join(datadir, "ca.cert")
				if c.ConfigFile.V() == c.ConfigFile.Def {
					e = c.ConfigFile.Set(filepath.Join(datadir, "pod.json"))
				}
				c.ConfigFile.Def = filepath.Join(datadir, constant.PodConfigFilename)
				if c.RPCKey.V() == c.RPCKey.Def {
					e = c.RPCKey.Set(filepath.Join(datadir, "rpc.key"))
				}
				c.RPCKey.Def = filepath.Join(datadir, "rpc.key")
				if c.RPCCert.V() == c.RPCCert.Def {
					e = c.RPCCert.Set(filepath.Join(datadir, "rpc.cert"))
				}
				c.RPCCert.Def = filepath.Join(datadir, "rpc.cert")
			}
		}
	}
	for i := range options {
		if options[i].Name() == "network" {
			I.Ln("network was set", optVals[i])
			if c.WalletFile.V() == c.WalletFile.Def {
				_, e = c.WalletFile.ReadInput(filepath.Join(datadir, optVals[i], "wallet.db"))
			}
			c.WalletFile.Def = filepath.Join(datadir, optVals[i], constant.DbName)
			if c.LogDir.V() == c.LogDir.Def {
				_, e = c.LogDir.ReadInput(filepath.Join(datadir, optVals[i]))
			}
			c.LogDir.Def = filepath.Join(datadir, optVals[i])
		}
	}
	// load the configuration file into the config
	resolvedConfigPath := c.ConfigFile.V()
	if configPath != "" {
		T.Ln("loading config from", configPath)
		resolvedConfigPath = configPath
	} else {
		if datadir != "" {
			if strings.HasPrefix(datadir, "~") {
				var homeDir string
				var usr *user.User
				var e error
				if usr, e = user.Current(); e == nil {
					homeDir = usr.HomeDir
				}
				// Fall back to standard HOME environment variable that works for most POSIX OSes if the directory from the Go
				// standard lib failed.
				if e != nil || homeDir == "" {
					homeDir = os.Getenv("HOME")
				}

				datadir = strings.Replace(datadir, "~", homeDir, 1)
			}

			if resolvedConfigPath, e = filepath.Abs(filepath.Clean(filepath.Join(datadir, constant.PodConfigFilename,
			),
			),
			); E.Chk(e) {
			}
			T.Ln("loading config from", resolvedConfigPath)
		}
	}
	var configExists bool
	if e = c.loadConfig(resolvedConfigPath); !D.Chk(e) {
		configExists = true
	}
	// read the environment variables into the config
	I.Ln("reading environment variables...")
	if e = c.loadEnvironment(); D.Chk(e) {
	}
	// read in the commandline options over top as they have highest priority
	I.Ln("decoding option inputs")
	for i := range options {
		if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
		}
	}
	if !configExists || c.Save.True() {
		c.Save.F()
		// save the configuration file
		I.Ln("saving configuration file...")
		var j []byte
		// c.ShowAll=true
		if j, e = json.MarshalIndent(c, "", "    "); !E.Chk(e) {
			I.F("saving config\n%s\n", string(j))
			apputil.EnsureDir(resolvedConfigPath)
			if e = ioutil.WriteFile(resolvedConfigPath, j, 0660); E.Chk(e) {
				panic(e)
			}
		}
		I.Ln("configuration file saved")
	}
	I.Ln("configuration initialised")
	return
}

// loadEnvironment scans the environment variables for values relevant to pod
func (c *Config) loadEnvironment() (e error) {
	env := os.Environ()
	c.ForEach(func(o opt.Option) bool {
		varName := "POD_" + strings.ToUpper(o.Name())
		for i := range env {
			if strings.HasPrefix(env[i], varName) {
				envVal := strings.Split(env[i], varName)[1]
				if _, e = o.LoadInput(envVal); D.Chk(e) {
				}
			}
		}
		return true
	},
	)

	return
}

// loadConfig loads the config from a file and unmarshals it into the config
func (c *Config) loadConfig(path string) (e error) {
	e = fmt.Errorf("no config found at %s", path)
	var cf []byte
	if !apputil.FileExists(path) {
		return
	} else if cf, e = ioutil.ReadFile(path); !D.Chk(e) {
		I.Ln("read in file from", path)
		if e = json.Unmarshal(cf, c); D.Chk(e) {
			panic(e)
		}
		I.Ln("unmarshalled", path)
		c.WalletPass.Set("")
	}
	return
}

// WriteToFile writes the current config to a file as json
func (c *Config) WriteToFile(filename string) (e error) {
	var j []byte
	I.S(c.MulticastPass.Bytes())
	wpp := c.WalletPass.Bytes()
	wp := make([]byte, len(wpp))
	copy(wp,wpp)
	if len(wp) > 0 {
		bhb := blake3.Sum256(wpp)
		bh := hex.EncodeToString(bhb[:])
		c.WalletPass.Set(bh)
	}
	if j, e = json.MarshalIndent(c, "", "  "); E.Chk(e) {
		return
	}
	if e = ioutil.WriteFile(filename, j, 0660); E.Chk(e) {
	}
	if len(wp) > 0 {
		c.WalletPass.SetBytes(wp)
	}
	return
}

// ForEach iterates the options in defined order with a closure that takes an opt.Option
func (c *Config) ForEach(fn func(ifc opt.Option) bool) bool {
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		// asserting to an Option ensures we skip the ancillary fields
		if iff, ok := t.Field(i).Interface().(opt.Option); ok {
			if !fn(iff) {
				return false
			}
		}
	}
	return true
}

// GetOption searches for a match amongst the podopts
func (c *Config) GetOption(input string) (
	op opt.Option, value string,
	e error,
) {
	T.Ln("checking arg for opt:", input)
	found := false
	if c.ForEach(func(ifc opt.Option) bool {
		aos := ifc.GetAllOptionStrings()
		for i := range aos {
			if strings.HasPrefix(input, aos[i]) {
				value = input[len(aos[i]):]
				found = true
				op = ifc
				return false
			}
		}
		return true
	},
	) {
	}
	if !found {
		e = fmt.Errorf("opt not found")
	}
	return
}

// MarshalJSON implements the json marshaller for the config. It only stores non-default values so can be composited.
func (c *Config) MarshalJSON() (b []byte, e error) {
	outMap := make(map[string]interface{})
	c.ForEach(
		func(ifc opt.Option) bool {
			switch ii := ifc.(type) {
			case *binary.Opt:
				if ii.True() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.True()
			case *list.Opt:
				v := ii.S()
				if len(v) == len(ii.Def) && ii.Data.OmitEmpty && !c.ShowAll {
					foundMismatch := false
					for i := range v {
						if v[i] != ii.Def[i] {
							foundMismatch = true
							break
						}
					}
					if !foundMismatch {
						return true
					}
				}
				outMap[ii.Option] = v
			case *float.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.Value.Load()
			case *integer.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.Value.Load()
			case *text.Opt:
				v := string(ii.Value.Load().([]byte))
				// fmt.Printf("def: '%s'", v)
				// spew.Dump(ii.def)
				if v == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = v
			case *duration.Opt:
				if ii.Value.Load() == ii.Def && ii.Data.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = fmt.Sprint(ii.Value.Load())
			default:
			}
			return true
		},
	)
	return json.Marshal(&outMap)
}

// UnmarshalJSON implements the Unmarshaller interface so it only writes to fields with those non-default values set.
func (c *Config) UnmarshalJSON(data []byte) (e error) {
	ifc := make(map[string]interface{})
	if e = json.Unmarshal(data, &ifc); E.Chk(e) {
		return
	}
	// I.S(ifc)
	c.ForEach(func(iii opt.Option) bool {
		switch ii := iii.(type) {
		case *binary.Opt:
			if i, ok := ifc[ii.Option]; ok {
				var ir bool
				if ir, ok = i.(bool); ir != ii.Def {
					// I.Ln(ii.Option+":", i.(binary), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).True())
					ii.Set(ir)
				}
			}
		case *list.Opt:
			matched := true
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.([]interface{}); ok2 {
					for i := range ds {
						if len(ii.Def) >= len(ds) {
							if ds[i] != ii.Def[i] {
								matched = false
								break
							}
						} else {
							matched = false
						}
					}
					if matched {
						return true
					}
					// I.Ln(ii.Option+":", ds, "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).S())
					ii.Set(ifcToStrings(ds))
				}
			}
		case *float.Opt:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", d.(float64), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(d.(float64))
			}
		case *integer.Opt:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", int64(d.(float64)), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(int(d.(float64)))
			}
		case *text.Opt:
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.(string); ok2 {
					if ds != ii.Def {
						// I.Ln(ii.Option+":", d.(string), "default:", ii.def, "prev:", c.Map[ii.Option].(*Opt).V())
						ii.Set(d.(string))
					}
				}
			}
		case *duration.Opt:
			if d, ok := ifc[ii.Option]; ok {
				var parsed time.Duration
				parsed, e = time.ParseDuration(d.(string))
				// I.Ln(ii.Option+":", parsed, "default:", ii.Opt(), "prev:", c.Map[ii.Option].(*Opt).V())
				ii.Set(parsed)
			}
		default:
		}
		return true
	},
	)
	return
}

func (c *Config) processCommandlineArgs(args []string) (
	remArgs []string, cm *cmds.Command, op []opt.Option,
	optVals []string, e error,
) {
	// first we will locate all the commands specified to mark the 3 sections, opt, commands, and the remainder is
	// arbitrary for the node
	commands := make(map[int]cmds.Command)
	var commandsStart, commandsEnd int = -1, -1
	var found, helpFound bool
	for i := range args {
		T.Ln("checking for commands:", args[i], commandsStart, commandsEnd, "current arg index:", i)
		var depth, dist int
		if found, depth, dist, cm, e = c.Commands.Find(args[i], depth, dist, false); E.Chk(e) {
			continue
		}
		if cm != nil {
			if cm.Name == "help" {
				helpFound = true
			}
			if cm.Parent != nil {
				if cm.Parent.Name == "help" && !helpFound {
					// I.S(commands)
					if commands[0].Name != "help" {
						found = false
					}
				}
			}
		}
		if found {
			if commandsStart == -1 {
				commandsStart = i
				commandsEnd = i + 1
			}
			if oc, ok := commands[depth]; ok {
				if !helpFound {
					e = fmt.Errorf("second command found at same depth '%s' and '%s'", oc.Name, cm.Name)
					return
				} else {
					// if this is a command match after help is found it is a
					// help command, resume the search
					if cmdd, ok := commands[depth]; ok {
						I.Ln("base level command is already occupied by",
							cmdd.Name+"; marking end of commands")
						// commandsEnd--
						T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
						break
					}
					I.Ln("found help subcommand:", cm.Name, depth)
					depth--
					dist++
					if found, depth, dist, cm, e = c.Commands.Find(args[i],
						depth, dist, true); E.Chk(e) {
						I.Ln("we didn't attach this help command")
						continue
					}
					// I.S(cm)
				}
			}
			commandsEnd = i + 1
			T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			T.Ln("found command", cm.Name, "argument number", i, "at depth", depth, "distance", dist)
			commands[depth] = *cm
		} else {
			// commandsStart=i+1
			// commandsEnd=i+1
			// T.Ln("not found:", args[i], "commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			T.Ln("argument", args[i], "is not a command", commandsStart, commandsEnd)
			c.FoundArgs = append(c.FoundArgs, args[i])
		}
	}
	// I.S(commands, cm)
	// commandsEnd++
	cmds := []int{}
	if len(commands) == 0 {
		I.Ln("setting default command")
		commands[0] = c.Commands[0]
		// I.S(commands[0])
		ibs := commands[0]
		cm = &ibs
		log.AppColorizer = commands[0].Colorizer
		log.App = commands[0].AppText
	} else {
		T.Ln("checking commands")
		// I.S(commands)
		for i := range commands {
			cmds = append(cmds, i)
		}
		I.S(cmds)
		if len(cmds) > 0 {
			sort.Ints(cmds)
			var cms []string
			for i := range commands {
				cms = append(cms, commands[i].Name)
			}
			if cmds[0] != 1 {
				e = fmt.Errorf("commands must include base level item for disambiguation %v", cms)
			}
			prev := cmds[0]
			for i := range cmds {
				if i == 0 {
					continue
				}
				if cmds[i] != prev+1 {
					e = fmt.Errorf("more than one command specified, %v", cms)
					return
				}
				found = false
				for j := range commands[cmds[i-1]].Commands {
					if commands[cmds[i]].Name == commands[cmds[i-1]].Commands[j].Name {
						found = true
					}
				}
				if !found {
					e = fmt.Errorf("multiple commands are not a path on the command tree %v", cms)
					return
				}
			}
		}
		T.Ln("commands:", commandsStart, commandsEnd)
		T.Ln("length of gathered commands", len(commands))
		if len(commands) == 1 {
			for _, x := range commands {
				cm = &x
				if cm.Colorizer != nil {
					log.AppColorizer = cm.Colorizer
					log.App = cm.AppText
				}
			}
		}
	}
	// if there was no command the commands start and end after all the args
	if commandsStart < 0 || commandsEnd < 0 {
		commandsStart = len(args)
		commandsEnd = commandsStart
	}
	T.Ln("commands section:", commandsStart, commandsEnd)
	if commandsStart > 0 {
		T.Ln("opts found", args[:commandsStart])
		// we have opts to check
		for i := range args {
			// if i == 0 {
			// 	continue
			// }
			if i >= commandsStart {
				break
			}
			var val string
			var o opt.Option
			if o, val, e = c.GetOption(args[i]); E.Chk(e) {
				e = fmt.Errorf("argument %d: '%s' lacks a valid opt prefix", i, args[i])
				return
			}
			if _, e = o.ReadInput(val); E.Chk(e) {
				return
			}
			T.Ln("found opt:", o.String())
			op = append(op, o)
			optVals = append(optVals, val)
		}
	}
	T.Ln("options to pass to child processes to match config:", args[:commandsStart])

	if len(cmds) < 1 {
		cmds = []int{0}
		commands[0] = c.Commands[0]
	}
	if commandsEnd > 0 && len(args) > commandsEnd {
		remArgs = args[commandsEnd:]
	}
	D.F("args that will pass to command: %v", remArgs)
	// I.S(commands, cm)
	// I.S(commands[cmds[len(cmds)-1]], op, args[commandsEnd:])
	return
}

// ReadCAFile reads in the configured Certificate Authority for TLS connections
func (c *Config) ReadCAFile() []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if c.ClientTLS.True() {
		var e error
		if certs, e = ioutil.ReadFile(c.CAFile.V()); E.Chk(e) {
			// If there's an error reading the CA file, continue with nil certs and without the client connection.
			certs = nil
		}
	} else {
		I.Ln("chain server RPC TLS is disabled")
	}
	return certs
}

func ifcToStrings(ifc []interface{}) (o []string) {
	for i := range ifc {
		o = append(o, ifc[i].(string))
	}
	return
}

type details struct {
	name, option, desc, def string
	aliases                 []string
	documentation           string
}

// GetHelp walks the command tree and gathers the options and creates a set of help functions for all commands and
// options in the set
func (c *Config) GetHelp(hf func(ifc interface{}) error) {
	helpCommand := cmds.Command{
		Name:       "help",
		Title:      "prints information about how to use pod",
		Entrypoint: hf,
		Commands:   nil,
	}
	// first add all the options
	c.ForEach(func(ifc opt.Option) bool {
		o := fmt.Sprintf("Parallelcoin Pod All-in-One Suite\n\n")
		var dt details
		switch ii := ifc.(type) {
		case *binary.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *list.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *float.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *integer.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *text.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		case *duration.Opt:
			dt = details{
				ii.GetMetadata().Name, ii.Option, ii.Description,
				fmt.Sprint(ii.Def), ii.Aliases,
				ii.Documentation,
			}
		}
		allNames := append([]string{dt.option}, dt.aliases...)
		for i := range allNames {
			helpCommand.Commands = append(helpCommand.Commands, cmds.Command{
				Name:  allNames[i],
				Title: dt.desc,
				Entrypoint: func(ifc interface{}) (e error) {
					o += fmt.Sprintf(
						"Help information about %s\n\n"+
							"\toption name:\n\t\t%s\n"+
							"\taliases:\n\t\t%s\n"+
							"\tdescription:\n\t\t%s\n"+
							"\tdefault:\n\t\t%v\n",
						dt.name, dt.option, dt.aliases, dt.desc, dt.def,
					)
					if dt.documentation != "" {
						o += "\tdocumentation:\n\t\t" + dt.documentation + "\n\n"
					}
					fmt.Fprint(os.Stderr, o)
					return
				},
				Commands: nil,
				Parent:   &helpCommand,
			},
			)
		}
		// for i := range
		return true
	},
	)
	// next add all the commands
	c.Commands.ForEach(func(cm cmds.Command) bool {
		helpCommand.Commands = append(helpCommand.Commands, cmds.Command{
			Name:        cm.Name,
			Title:       cm.Title,
			Description: cm.Description,
			Entrypoint: func(ifc interface{}) (e error) {
				o := fmt.Sprintf(
					"Help information about command '%s'\n\n"+
						"%s\n\n",
					cm.Name, cm.Title,
				)
				if cm.Description != "" {
					split := strings.Split(cm.Description, "\n")
					docs := "\t" + strings.Join(split, "\n\n\t")
					o += docs + "\n\n"
				}
				o += "Related options:\n\n"
				descs := make(map[string]string)
				c.ForEach(func(ifc opt.Option) bool {
					meta := ifc.GetMetadata()
					found := false
					for i := range meta.Tags {
						if meta.Tags[i] == cm.Name {
							found = true
						}
					}
					if !found {
						// skip item as it isn't tagged with this program name
						return true
					}
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
				fmt.Fprint(os.Stderr, o)
				return
			},
			Parent: &helpCommand,
		})
		return true
	}, 0, 0,
	)
	c.Commands = append(c.Commands, helpCommand)
	return
}

func getAllOptionStrings(c *Config) (s map[string][]string, e error) {
	s = make(map[string][]string)
	if c.ForEach(func(ifc opt.Option) bool {
		md := ifc.GetMetadata()
		if _, ok := s[ifc.Name()]; ok {
			e = fmt.Errorf("conflicting opt names: %v %v", ifc.GetAllOptionStrings(), s[ifc.Name()])
			return false
		}
		s[ifc.Name()] = md.GetAllOptionStrings()
		return true
	},
	) {
	}
	// s["commandslist"] = c.Commands.GetAllCommands()
	return
}

func findConflictingItems(valOpts map[string][]string) (o []string, e error) {
	var ss, ls string
	for i := range valOpts {
		for j := range valOpts {
			if i == j {
				continue
			}
			a := valOpts[i]
			b := valOpts[j]
			for ii := range a {
				for jj := range b {
					ss, ls = shortestString(a[ii], b[jj])
					if ss == ls[:len(ss)] {
						E.F("conflict between %s and %s, ", ss, ls)
						o = append(o, ss, ls)
					}
				}
			}
		}
	}
	if len(o) > 0 {
		panic(fmt.Sprintf("conflicts found: %v", o))
	}
	return
}

func shortestString(a, b string) (s, l string) {
	switch {
	case len(a) > len(b):
		s, l = b, a
	default:
		s, l = a, b
	}
	return
}
