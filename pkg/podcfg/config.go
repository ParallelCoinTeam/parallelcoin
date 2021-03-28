// Package podcfg is a configuration system to fit with the all-in-one philosophy guiding the design of the parallelcoin
// pod.
//
// The configuration is stored by each component of the connected applications, so all data is stored in concurrent-safe
// atomics, and there is a facility to invoke a function in response to a new value written into a field by other
// threads.
//
// There is a custom JSON marshal/unmarshal for each field type and for the whole configuration that only saves values
// that differ from the defaults, similar to 'omitempty' in struct tags but where 'empty' is the default value instead
// of the default zero created by Go's memory allocator.
//
//
package podcfg

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/base58"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"io/ioutil"
	"os"
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

// Config defines the configuration items used by pod along with the various components included in the suite
type Config struct {
	// ShowAll is a flag to make the json encoder explicitly define all fields and not just the ones different to the
	// defaults
	ShowAll bool
	// Map is the same data but addressible using its name as found inside the various configuration types, the key is
	// the same as the .Name field field in the various data types
	Map      map[string]interface{}
	Commands Commands
	// These are just the definitions, the things put in them are more useful than doc comments
	AddCheckpoints         *Strings
	AddPeers               *Strings
	AddrIndex              *Bool
	AutoListen             *Bool
	AutoPorts              *Bool
	BanDuration            *Duration
	BanThreshold           *Int
	BlockMaxSize           *Int
	BlockMaxWeight         *Int
	BlockMinSize           *Int
	BlockMinWeight         *Int
	BlockPrioritySize      *Int
	BlocksOnly             *Bool
	CAFile                 *String
	ConfigFile             *String
	ConnectPeers           *Strings
	Controller             *Bool
	CPUProfile             *String
	DarkTheme              *Bool
	DataDir                *String
	DbType                 *String
	DisableBanning         *Bool
	DisableCheckpoints     *Bool
	DisableDNSSeed         *Bool
	DisableListen          *Bool
	DisableRPC             *Bool
	Discovery              *Bool
	ExternalIPs            *Strings
	FreeTxRelayLimit       *Float
	Generate               *Bool
	GenThreads             *Int
	Hilite                 *Strings
	LAN                    *Bool
	Language               *String
	LimitPass              *String
	LimitUser              *String
	LogDir                 *String
	LogFilter              *Strings
	LogLevel               *String
	MaxOrphanTxs           *Int
	MaxPeers               *Int
	MulticastPass          *String
	MiningAddrs            *Strings
	MinRelayTxFee          *Float
	Network                *String
	NoCFilters             *Bool
	NodeOff                *Bool
	NoInitialLoad          *Bool
	NoPeerBloomFilters     *Bool
	NoRelayPriority        *Bool
	OneTimeTLSKey          *Bool
	Onion                  *Bool
	OnionProxy             *String
	OnionProxyPass         *String
	OnionProxyUser         *String
	P2PConnect             *Strings
	P2PListeners           *Strings
	Password               *String
	PipeLog                *Bool
	Profile                *String
	Proxy                  *String
	ProxyPass              *String
	ProxyUser              *String
	RejectNonStd           *Bool
	RelayNonStd            *Bool
	RPCCert                *String
	RPCConnect             *String
	RPCKey                 *String
	RPCListeners           *Strings
	RPCMaxClients          *Int
	RPCMaxConcurrentReqs   *Int
	RPCMaxWebsockets       *Int
	RPCQuirks              *Bool
	RunAsService           *Bool
	ServerPass             *String
	ServerTLS              *Bool
	ServerUser             *String
	SigCacheMaxSize        *Int
	Solo                   *Bool
	TLS                    *Bool
	TLSSkipVerify          *Bool
	TorIsolation           *Bool
	TrickleInterval        *Duration
	TxIndex                *Bool
	UPNP                   *Bool
	UserAgentComments      *Strings
	Username               *String
	UUID                   *Int
	Wallet                 *Bool
	WalletFile             *String
	WalletOff              *Bool
	WalletPass             *String
	WalletRPCListeners     *Strings
	WalletRPCMaxClients    *Int
	WalletRPCMaxWebsockets *Int
	WalletServer           *String
	Whitelists             *Strings
}

// ForEach iterates the configuration items in their defined order, running a function with the configuration item in
// the field
func (c *Config) ForEach(fn func(ifc interface{}) bool) bool {
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		if !fn(t.Field(i).Interface()) {
			return false
		}
	}
	return true
}

func (c *Config) GetOption(input string) (opt Option, value string, e error) {
	I.Ln("checking arg for option:", input)
	if c.ForEach(func(ifc interface{}) bool {
		var md *Metadata
		switch ii := ifc.(type) {
		case *Bool:
			opt = ii
			md = &ii.Metadata
		case *Strings:
			opt = ii
			md = &ii.Metadata
		case *Float:
			opt = ii
			md = &ii.Metadata
		case *Int:
			opt = ii
			md = &ii.Metadata
		case *String:
			opt = ii
			md = &ii.Metadata
		case *Duration:
			opt = ii
			md = &ii.Metadata
		}
		if md != nil {
			I.Ln(input, md.Option, input == md.Option, len(md.Option))
			if strings.HasPrefix(input, md.Option) {
				value = input[len(md.Option):]
				I.Ln("value", value)
				return false
			}
			I.Ln(md.Aliases)
			for i := range md.Aliases {
				if strings.HasPrefix(input, md.Aliases[i]) {
					value = input[len(md.Aliases[i]):]
					I.Ln("value", value)
					return false
				}
			}
		}
		return true
	},
	) {
		e = fmt.Errorf("option not found")
	}
	return
}

// MarshalJSON implements the json marshaller for the config. This marshaller only saves what is different from the
// defaults, and when it is unmarshalled, only the fields stored are altered, thus allowing stacking several sources
// such as environment variables, command line flags and the config file itself.
func (c *Config) MarshalJSON() (b []byte, e error) {
	outMap := make(map[string]interface{})
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				if ii.True() == ii.def && ii.Metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.True()
			case *Strings:
				v := ii.S()
				if len(v) == len(ii.def) && ii.Metadata.OmitEmpty && !c.ShowAll {
					foundMismatch := false
					for i := range v {
						if v[i] != ii.def[i] {
							foundMismatch = true
							break
						}
					}
					if !foundMismatch {
						return true
					}
				}
				outMap[ii.Option] = v
			case *Float:
				if ii.value.Load() == ii.def && ii.Metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.value.Load()
			case *Int:
				if ii.value.Load() == ii.def && ii.Metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = ii.value.Load()
			case *String:
				v := string(ii.value.Load().([]byte))
				// fmt.Printf("def: '%s'", v)
				// spew.Dump(ii.def)
				if v == ii.def && ii.Metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = v
			case *Duration:
				if ii.value.Load() == ii.def && ii.Metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Option] = fmt.Sprint(ii.value.Load())
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
	c.ForEach(func(iii interface{}) bool {
		switch ii := iii.(type) {
		case *Bool:
			if i, ok := ifc[ii.Option]; ok {
				if i.(bool) != ii.def {
					// I.Ln(ii.Option+":", i.(bool), "default:", ii.def, "prev:", c.Map[ii.Option].(*Bool).True())
					c.Map[ii.Option].(*Bool).Set(i.(bool))
				}
			}
		case *Strings:
			matched := true
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.([]interface{}); ok2 {
					for i := range ds {
						if ds[i] != ii.def[i] {
							matched = false
							break
						}
					}
					if matched {
						return true
					}
					// I.Ln(ii.Option+":", ds, "default:", ii.def, "prev:", c.Map[ii.Option].(*Strings).S())
					c.Map[ii.Option].(*Strings).Set(ifcToStrings(ds))
				}
			}
		case *Float:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", d.(float64), "default:", ii.def, "prev:", c.Map[ii.Option].(*Float).V())
				c.Map[ii.Option].(*Float).Set(d.(float64))
			}
		case *Int:
			if d, ok := ifc[ii.Option]; ok {
				// I.Ln(ii.Option+":", int64(d.(float64)), "default:", ii.def, "prev:", c.Map[ii.Option].(*Int).V())
				c.Map[ii.Option].(*Int).Set(int(d.(float64)))
			}
		case *String:
			if d, ok := ifc[ii.Option]; ok {
				if ds, ok2 := d.(string); ok2 {
					if ds != ii.def {
						// I.Ln(ii.Option+":", d.(string), "default:", ii.def, "prev:", c.Map[ii.Option].(*String).V())
						c.Map[ii.Option].(*String).Set(d.(string))
					}
				}
			}
		case *Duration:
			if d, ok := ifc[ii.Option]; ok {
				var parsed time.Duration
				parsed, e = time.ParseDuration(d.(string))
				// I.Ln(ii.Option+":", parsed, "default:", ii.def.String(), "prev:", c.Map[ii.Option].(*Duration).V())
				c.Map[ii.Option].(*Duration).Set(parsed)
			}
		default:
		}
		return true
	},
	)
	return
}

// Initialize loads in configuration from disk and from environment on top of the default base
//
// the several places configuration is sourced from are overlaid in the following order:
// default -> config file -> environment variables -> commandline flags
func (c *Config) Initialize() (e error) {
	// first process the commandline
	var cm *Command
	if cm, e = c.processCommandlineArgs(); E.Chk(e) {
		return
	}
	_ = cm
	return
}

func (c *Config) processCommandlineArgs() (cm *Command, e error) {
	// first we will locate all the commands specified to mark the 3 sections, options, commands, and the remainder is
	// arbitrary for the app
	var commands map[int]Command
	commands = make(map[int]Command)
	var commandsStart, commandsEnd int
	for i := range os.Args {
		if i == 0 {
			continue
		}
		var depth, dist int
		var found bool
		if found, depth, dist, cm, e = c.Commands.Find(os.Args[i], depth, dist); E.Chk(e) || !found {
			continue
		}
		if found {
			if commandsStart == 0 {
				commandsStart = i
			}
			commandsEnd = i
			if oc, ok := commands[depth]; ok {
				e = fmt.Errorf("second command found at same depth '%s' and '%s'", oc.Name, cm.Name)
				return
			}
			D.Ln("found command", cm.Name, "argument number", i, "at depth", depth, "distance", dist)
			commands[depth] = *cm
		} else {
			T.Ln("argument", os.Args[i], "is not a command")
		}
	}
	commandsEnd++
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
				found := false
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
	}
	var options []Option
	if commandsStart > 1 {
		// we have options to check
		for i := range os.Args {
			if i == 0 {
				continue
			}
			if i == commandsStart {
				break
			}
			var val string
			var opt Option
			if opt, val, e = c.GetOption(os.Args[i]); E.Chk(e) {
				e = fmt.Errorf("argument %d: '%s' lacks a valid option prefix", i, os.Args[i])
				return
			}
			I.Ln("found option:", opt.GetMetadata().Option, "with value", val)
			if opt, e = opt.ReadInput(val); E.Chk(e) {
				return
			}
			options = append(options, opt)
		}
	}
	I.S(commands[cmds[len(cmds)-1]], options, os.Args[commandsEnd:])
	return
}

// ReadCAFile reads in the configured Certificate Authority for TLS connections
func ReadCAFile(config *Config) []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if config.TLS.True() {
		var e error
		if certs, e = ioutil.ReadFile(config.CAFile.V()); E.Chk(e) {
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
