// Package opts is a configuration system to fit with the all-in-one philosophy guiding the design of the parallelcoin
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
package opts

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/text"
	"github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/base58"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	Name              = "pod"
	confExt           = ".json"
	appLanguage       = "en"
	PodConfigFilename = Name + confExt
	PARSER            = "json"
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
//
// the several places configuration is sourced from are overlaid in the following order:
// default -> config file -> environment variables -> commandline flags
func (c *Config) Initialize() (e error) {
	T.Ln("initializing configuration...")
	// first lint the configuration
	var aos map[string][]string
	if aos, e = c.getAllOptionStrings(); E.Chk(e) {
		return
	}
	// this function will panic if there is potential for ambiguity in the commandline configuration args.
	T.Ln("linting configuration items")
	if _, e = findConflictingItems(aos); E.Chk(e) {
	}
	// generate and add the help commands to the help tree
	c.getHelp()
	// process the commandline
	T.Ln("processing commandline arguments", os.Args[1:])
	var cm *Command
	var options []opt.Option
	var optVals []string
	if cm, options, optVals, e = c.processCommandlineArgs(os.Args[1:]); E.Chk(e) {
		return
	}
	_ = options
	c.RunningCommand = cm
	// if the user sets the configfile directly, or the datadir on the commandline we need to load it from that path
	T.Ln("checking from where to load the configuration file")
	datadir := c.DataDir.V()
	var configPath string
	for i := range options {
		if options[i].Name() == "configfile" {
			if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
				configPath = optVals[i]
			}
		} else if options[i].Name() == "datadir" {
			if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
				datadir = optVals[i]
			}
		}
	}
	// load the configuration file into the config
	resolvedConfigPath := c.ConfigFile.V()
	if configPath != "" {
		T.Ln("loading config from", configPath)
		resolvedConfigPath = configPath
	} else {
		if datadir != "" {
			resolvedConfigPath = filepath.Join(datadir, PodConfigFilename)
			T.Ln("loading config from", resolvedConfigPath)
		}
	}
	var configExists bool
	if e = c.loadConfig(resolvedConfigPath); !D.Chk(e) {
		configExists = true
	}
	// read the environment variables into the config
	if e = c.loadEnvironment(); D.Chk(e) {
	}
	// read in the commandline options
	for i := range options {
		if _, e = options[i].ReadInput(optVals[i]); E.Chk(e) {
		}
	}
	if !configExists || c.Save.True() {
		// save the configuration file
		var j []byte
		// c.ShowAll=true
		if j, e = json.MarshalIndent(c, "", "    "); !E.Chk(e) {
			I.F("saving config\n%s\n", string(j))
			apputil.EnsureDir(resolvedConfigPath)
			if e = ioutil.WriteFile(resolvedConfigPath, j, 0660); E.Chk(e) {
			}
		}
		
	}
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
	} else if cf, e = ioutil.ReadFile(path); !D.Chk(e) {
		if e = json.Unmarshal(cf, c); D.Chk(e) {
		}
	}
	return
}

// ForEach iterates the configuration items in their defined order, running a function with the configuration item in
// the field
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

// GetOption searches for a match amongst the opts
func (c *Config) GetOption(input string) (op opt.Option, value string, e error) {
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

// MarshalJSON implements the json marshaller for the config. This marshaller only saves what is different from the
// defaults, and when it is unmarshalled, only the fields stored are altered, thus allowing stacking several sources
// such as environment variables, command line flags and the config file itself.
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

// UnmarshalJSON implements the Unmarshaller interface with a specific goal to be well suited to compositing multiple
// layers on top of the default base from multiple sources
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

func (c *Config) processCommandlineArgs(args []string) (cm *Command, op []opt.Option, optVals []string, e error) {
	// first we will locate all the commands specified to mark the 3 sections, opt, commands, and the remainder is
	// arbitrary for the app
	var commands map[int]Command
	commands = make(map[int]Command)
	var commandsStart, commandsEnd int
	var found bool
	for i := range args {
		if i == 0 {
			// commandsStart = i
			// commandsEnd = i
			continue
		}
		T.Ln("checking for commands:", args[i])
		T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
		var depth, dist int
		if found, depth, dist, cm, e = c.Commands.Find(args[i], depth, dist); E.Chk(e) {
			continue
		}
		if found {
			if commandsStart == 0 {
				commandsStart = i
				commandsEnd = i + 1
			}
			T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			if oc, ok := commands[depth]; ok {
				e = fmt.Errorf("second command found at same depth '%s' and '%s'", oc.Name, cm.Name)
				return
			}
			commandsEnd = i + 1
			T.Ln("found command", cm.Name, "argument number", i, "at depth", depth, "distance", dist)
			commands[depth] = *cm
		} else {
			T.Ln("not found:", args[i], "commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
			// commandsStart++
			// commandsEnd++
			T.Ln("argument", args[i], "is not a command")
		}
	}
	// commandsEnd++
	cmds := []int{}
	if len(commands) == 0 {
		commands[0] = c.Commands[0]
	} else {
		for i := range commands {
			cmds = append(cmds, i)
		}
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
		T.Ln("commandStart", commandsStart, commandsEnd, args[commandsStart:commandsEnd])
	}
	if commandsStart > 1 {
		T.Ln("opt found", args[:commandsStart])
		// we have opt to check
		for i := range args {
			// if i == 0 {
			// 	continue
			// }
			if i == commandsStart {
				break
			}
			var val string
			var o opt.Option
			if o, val, e = c.GetOption(args[i]); E.Chk(e) {
				e = fmt.Errorf("argument %d: '%s' lacks a valid opt prefix", i, args[i])
				return
			}
			// if _, e = opt.ReadInput(val); E.Chk(e) {
			// 	return
			// }
			T.Ln("found opt:", o.String())
			op = append(op, o)
			optVals = append(optVals, val)
		}
	}
	if len(cmds) < 1 {
		cmds = []int{0}
		commands[0] = c.Commands[0]
	}
	I.S(commands[cmds[len(cmds)-1]], op, args[commandsEnd:])
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

func genPassword() string {
	s, e := hdkeychain.GenerateSeed(16)
	if e != nil {
		panic("can't do nothing without entropy! " + e.Error())
	}
	return base58.Encode(s)
}

func ifcToStrings(ifc []interface{}) (o []string) {
	for i := range ifc {
		o = append(o, ifc[i].(string))
	}
	return
}
