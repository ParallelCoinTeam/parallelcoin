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
func (c *Config) ForEach(fn func(ifc interface{}) bool) {
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		if !fn(t.Field(i).Interface()) {
			return
		}
	}
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
				if ii.True() == ii.def && ii.metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Name] = ii.True()
			case *Strings:
				v := ii.S()
				if len(v) == len(ii.def) && ii.metadata.OmitEmpty && !c.ShowAll {
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
				outMap[ii.Name] = v
			case *Float:
				if ii.value.Load() == ii.def && ii.metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Name] = ii.value.Load()
			case *Int:
				if ii.value.Load() == ii.def && ii.metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Name] = ii.value.Load()
			case *String:
				v := string(ii.value.Load().([]byte))
				// fmt.Printf("def: '%s'", v)
				// spew.Dump(ii.def)
				if v == ii.def && ii.metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Name] = v
			case *Duration:
				if ii.value.Load() == ii.def && ii.metadata.OmitEmpty && !c.ShowAll {
					return true
				}
				outMap[ii.Name] = fmt.Sprint(ii.value.Load())
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
			if i, ok := ifc[ii.Name]; ok {
				if i.(bool) != ii.def {
					// I.Ln(ii.Name+":", i.(bool), "default:", ii.def, "prev:", c.Map[ii.Name].(*Bool).True())
					c.Map[ii.Name].(*Bool).Set(i.(bool))
				}
			}
		case *Strings:
			matched := true
			if d, ok := ifc[ii.Name]; ok {
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
					// I.Ln(ii.Name+":", ds, "default:", ii.def, "prev:", c.Map[ii.Name].(*Strings).S())
					c.Map[ii.Name].(*Strings).Set(ifcToStrings(ds))
				}
			}
		case *Float:
			if d, ok := ifc[ii.Name]; ok {
				// I.Ln(ii.Name+":", d.(float64), "default:", ii.def, "prev:", c.Map[ii.Name].(*Float).V())
				c.Map[ii.Name].(*Float).Set(d.(float64))
			}
		case *Int:
			if d, ok := ifc[ii.Name]; ok {
				// I.Ln(ii.Name+":", int64(d.(float64)), "default:", ii.def, "prev:", c.Map[ii.Name].(*Int).V())
				c.Map[ii.Name].(*Int).Set(int(d.(float64)))
			}
		case *String:
			if d, ok := ifc[ii.Name]; ok {
				if ds, ok2 := d.(string); ok2 {
					if ds != ii.def {
						// I.Ln(ii.Name+":", d.(string), "default:", ii.def, "prev:", c.Map[ii.Name].(*String).V())
						c.Map[ii.Name].(*String).Set(d.(string))
					}
				}
			}
		case *Duration:
			if d, ok := ifc[ii.Name]; ok {
				var parsed time.Duration
				parsed, e = time.ParseDuration(d.(string))
				// I.Ln(ii.Name+":", parsed, "default:", ii.def.String(), "prev:", c.Map[ii.Name].(*Duration).V())
				c.Map[ii.Name].(*Duration).Set(parsed)
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
			if oc, ok := commands[depth]; ok {
				e = fmt.Errorf("second command found at same depth '%s' and '%s'", oc.Name, cm.Name)
				return
			}
			I.Ln("found command", cm.Name, "argument number", i, "at depth", depth, "distance", dist)
			commands[depth] = *cm
		} else {
			I.Ln("argument", os.Args[i], "is not a command")
		}
	}
	I.S(commands)
	if len(commands) == 0 {
		commands[0] = c.Commands[0]
	} else {
		cmds := []int{}
		for i := range commands {
			cmds = append(cmds, i)
		}
		if len(cmds) > 0 {
			sort.Ints(cmds)
			I.S(cmds)
			var cms []string
			for j := range commands {
				cms = append(cms, commands[j].Name)
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
					var cms []string
					for j := range commands {
						cms = append(cms, commands[j].Name)
					}
					e = fmt.Errorf("multiple commands are not a path on the command tree %v", cms)
					return
				}
			}
		}
	}
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
