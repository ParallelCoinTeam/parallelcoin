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
	"github.com/p9c/pod/pkg/appdata"
	"github.com/p9c/pod/pkg/base58"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/urfave/cli"
	uberatomic "go.uber.org/atomic"
	"io/ioutil"
	"math/rand"
	"net"
	"path/filepath"
	"reflect"
	"sync/atomic"
	"time"
)

const (
	Name              = "pod"
	confExt           = ".json"
	appLanguage       = "en"
	PodConfigFilename = Name + confExt
	PARSER            = "json"
)

type Config struct {
	// ShowAll is a flag to make the json encoder explicitly define all fields and not just the ones different to the
	// defaults
	ShowAll bool
	// Map is the same data but addressible using its name as found inside the various configuration types, the key is
	// the same as the .Name field field in the various data types
	Map map[string]interface{}
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

// ForEach iterates the configuration items in their defined order, running a
// function with the configuration item in the field
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

func ifcToStrings(ifc []interface{}) (o []string) {
	for i := range ifc {
		o = append(o, ifc[i].(string))
	}
	return
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

// EmptyConfig creates a fresh Config with default values stored in its fields
func EmptyConfig() (c *Config) {
	network := "mainnet"
	rand.Seed(time.Now().Unix())
	var datadir = &atomic.Value{}
	datadir.Store([]byte(appdata.Dir(Name, false)))
	c = &Config{
		AddCheckpoints: NewStrings(
			metadata{
				Name:        "addcheckpoint",
				Group:       "debug",
				Label:       "Add Checkpoints",
				Description: "add custom checkpoints",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		AddPeers: NewStrings(
			metadata{
				Name:        "addpeer",
				Group:       "node",
				Label:       "Add Peers",
				Description: "manually adds addresses to try to connect to",
				Type:        "ipaddress",
				Widget:      "multi",
				// Hook:        "addpeer",
				OmitEmpty: true,
			},
			[]string{"127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12344"},
		),
		AddrIndex: NewBool(
			metadata{
				Name:        "addrindex",
				Group:       "node",
				Label:       "Address Index",
				Description: "maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
				Widget:      "toggle",
				// Hook:        "dropaddrindex",
				OmitEmpty: true,
			},
			true,
		),
		AutoPorts: NewBool(
			metadata{
				Name:        "autoports",
				Group:       "debug",
				Label:       "Automatic Ports",
				Description: "RPC and controller ports are randomized, use with controller for automatic peer discovery",
				Widget:      "toggle",
				// Hook: "restart",
				OmitEmpty: true,
			},
			false,
		),
		AutoListen: NewBool(
			metadata{
				Name:        "autolisten",
				Group:       "node",
				Label:       "Manual Listeners",
				Description: "automatically update inbound addresses dynamically according to discovered network interfaces",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			true,
		),
		BanDuration: NewDuration(
			metadata{
				Name:        "banduration",
				Group:       "debug",
				Label:       "Ban Duration",
				Description: "how long a ban of a misbehaving peer lasts",
				Widget:      "duration",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			time.Hour*24,
		),
		BanThreshold: NewInt(
			metadata{
				Name:        "banthreshold",
				Group:       "debug",
				Label:       "Ban Threshold",
				Description: "ban score that triggers a ban (default 100)",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultBanThreshold,
		),
		BlockMaxSize: NewInt(
			metadata{
				Name:        "blockmaxsize",
				Group:       "mining",
				Label:       "Block Max Size",
				Description: "maximum block size in bytes to be used when creating a block",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			BlockMaxSizeMax,
		),
		BlockMaxWeight: NewInt(
			metadata{
				Name:        "blockmaxweight",
				Group:       "mining",
				Label:       "Block Max Weight",
				Description: "maximum block weight to be used when creating a block",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			BlockMaxWeightMax,
		),
		BlockMinSize: NewInt(
			metadata{
				Name:        "blockminsize",
				Group:       "mining",
				Label:       "Block Min Size",
				Description: "minimum block size in bytes to be used when creating a block",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			BlockMaxSizeMin,
		),
		BlockMinWeight: NewInt(
			metadata{
				Name:        "blockminweight",
				Group:       "mining",
				Label:       "Block Min Weight",
				Description: "minimum block weight to be used when creating a block",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			BlockMaxWeightMin,
		),
		BlockPrioritySize: NewInt(
			metadata{
				Name:        "blockprioritysize",
				Group:       "mining",
				Label:       "Block Priority Size",
				Description: "size in bytes for high-priority/low-fee transactions when creating a block",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultBlockPrioritySize,
		),
		BlocksOnly: NewBool(
			metadata{
				Name:        "blocksonly",
				Group:       "node",
				Label:       "Blocks Only",
				Description: "do not accept transactions from remote peers",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		CAFile: NewString(
			metadata{
				Name:        "cafile",
				Group:       "tls",
				Label:       "Certificate Authority File",
				Description: "certificate authority file for TLS certificate validation",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			filepath.Join(string(datadir.Load().([]byte)), "ca.cert"),
		),
		ConfigFile: NewString(
			metadata{
				Name:        "configfile",
				Label:       "Configuration File",
				Description: "location of configuration file, cannot actually be changed",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			filepath.Join(string(datadir.Load().([]byte)), PodConfigFilename),
		),
		ConnectPeers: NewStrings(
			metadata{
				Name:        "connect",
				Group:       "node",
				Label:       "Connect Peers",
				Description: "connect ONLY to these addresses (disables inbound connections)",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		Controller: NewBool(
			metadata{
				Name:        "controller",
				Group:       "node",
				Label:       "Enable Controller",
				Description: "delivers mining jobs over multicast",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		CPUProfile: NewString(
			metadata{
				Name:        "cpuprofile",
				Group:       "debug",
				Label:       "CPU Profile",
				Description: "write cpu profile to this file",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		DarkTheme: NewBool(
			metadata{
				Name:        "darktheme",
				Group:       "config",
				Label:       "Dark Theme",
				Description: "sets dark theme for GUI",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		DataDir: &String{
			value: datadir,
			metadata: metadata{
				Name:        "datadir",
				Aliases:     []string{"D"},
				Label:       "Data Directory",
				Description: "root folder where application data is stored",
				Type:        "directory",
				Widget:      "string",
				OmitEmpty:   true,
			},
			def: appdata.Dir(Name, false),
		},
		DbType: NewString(
			metadata{
				Name:        "dbtype",
				Group:       "debug",
				Label:       "Database Type",
				Description: "type of database storage engine to use (only one right now, ffldb)",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultDbType,
		),
		DisableBanning: NewBool(
			metadata{
				Name:        "nobanning",
				Group:       "debug",
				Label:       "Disable Banning",
				Description: "disables banning of misbehaving peers",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		DisableCheckpoints: NewBool(
			metadata{
				Name:        "nocheckpoints",
				Group:       "debug",
				Label:       "Disable Checkpoints",
				Description: "disables all checkpoints",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		DisableDNSSeed: NewBool(
			metadata{
				Name:        "nodnsseed",
				Group:       "node",
				Label:       "Disable DNS Seed",
				Description: "disable seeding of addresses to peers",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		DisableListen: NewBool(
			metadata{
				Name:        "nolisten",
				Group:       "node",
				Label:       "Disable Listen",
				Description: "disables inbound connections for the peer to peer network",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		DisableRPC: NewBool(
			metadata{
				Name:        "norpc",
				Group:       "rpc",
				Label:       "Disable RPC",
				Description: "disable rpc servers, as well as kopach controller",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		Discovery: NewBool(
			metadata{
				Name:        "discover",
				Group:       "node",
				Label:       "Disovery",
				Description: "enable LAN peer discovery in GUI",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		ExternalIPs: NewStrings(
			metadata{
				Name:        "externalip",
				Group:       "node",
				Label:       "External IP Addresses",
				Description: "extra addresses to tell peers they can connect to",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		FreeTxRelayLimit: NewFloat(
			metadata{
				Name:        "limitfreerelay",
				Group:       "policy",
				Label:       "Free Tx Relay Limit",
				Description: "limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
				Widget:      "float",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultFreeTxRelayLimit,
		),
		Generate: NewBool(
			metadata{
				Name:        "generate",
				Aliases:     []string{"g"},
				Group:       "mining",
				Label:       "Generate Blocks",
				Description: "turn on Kopach CPU miner",
				Widget:      "toggle",
				// Hook:        "generate",
				OmitEmpty: true,
			},
			false,
		),
		GenThreads: NewInt(
			metadata{
				Name:        "genthreads",
				Group:       "mining",
				Label:       "Generate Threads",
				Description: "number of threads to mine with",
				Widget:      "integer",
				// Hook:        "genthreads",
				OmitEmpty: true,
			},
			-1,
		),
		Hilite: NewStrings(
			metadata{
				Name:        "highlight",
				Group:       "debug",
				Label:       "Hilite",
				Description: "list of packages that will print with attention getters",
				Type:        "string",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		LAN: NewBool(
			metadata{
				Name:        "lan",
				Group:       "debug",
				Label:       "LAN Testnet Mode",
				Description: "run without any connection to nodes on the internet (does not apply on mainnet)",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		Language: NewString(
			metadata{
				Name:        "language",
				Group:       "config",
				Label:       "Language",
				Description: "user interface language i18 localization",
				Widget:      "string",
				// Hook:        "language",
				OmitEmpty: true,
			},
			"en",
		),
		LimitPass: NewString(
			metadata{
				Name:        "limitpass",
				Group:       "rpc",
				Label:       "Limit Password",
				Description: "limited user password",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			genPassword(),
		),
		LimitUser: NewString(
			metadata{
				Name:        "limituser",
				Group:       "rpc",
				Label:       "Limit Username",
				Description: "limited user name",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"limit",
		),
		LogDir: NewString(
			metadata{
				Name:        "logdir",
				Group:       "config",
				Label:       "Log Directory",
				Description: "folder where log files are written",
				Type:        "directory",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			string(datadir.Load().([]byte)),
		),
		LogFilter: NewStrings(
			metadata{
				Name:        "logfilter",
				Group:       "debug",
				Label:       "Log Filter",
				Description: "comma-separated list of packages that will not print logs",
				Type:        "string",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		LogLevel: NewString(
			metadata{
				Name:        "loglevel",
				Aliases:     []string{"l"},
				Group:       "config",
				Label:       "Log Level",
				Description: "maximum log level to output\n(fatal error check warning info debug trace - what is selected includes all items to the left of the one in that list)",
				Widget:      "radio",
				Options: []string{"off",
					"fatal",
					"error",
					"info",
					"check",
					"debug",
					"trace",
				},
				// Hook:        "loglevel",
				OmitEmpty: true,
			},
			"info",
		),
		MaxOrphanTxs: NewInt(
			metadata{
				Name:        "maxorphantx",
				Group:       "policy",
				Label:       "Max Orphan Txs",
				Description: "max number of orphan transactions to keep in memory",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMaxOrphanTransactions,
		),
		MaxPeers: NewInt(
			metadata{
				Name:        "maxpeers",
				Group:       "node",
				Label:       "Max Peers",
				Description: "maximum number of peers to hold connections with",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMaxPeers,
		),
		MulticastPass: NewString(
			metadata{
				Name:        "minerpass",
				Group:       "config",
				Label:       "Multicast Pass",
				Description: "password that encrypts the connection to the mining controller",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"pa55word",
		),
		MiningAddrs: NewStrings(
			metadata{
				Name:        "miningaddrs",
				Label:       "Mining Addresses",
				Description: "addresses to pay block rewards to (not in use)",
				Type:        "base58",
				Widget:      "multi",
				// Hook:        "miningaddr",
				OmitEmpty: true,
			},
			[]string{},
		),
		MinRelayTxFee: NewFloat(
			metadata{
				Name:        "minrelaytxfee",
				Group:       "policy",
				Label:       "Min Relay Transaction Fee",
				Description: "the minimum transaction fee in DUO/kB to be considered a non-zero fee",
				Widget:      "float",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMinRelayTxFee.ToDUO(),
		),
		Network: NewString(
			metadata{
				Name:        "network",
				Group:       "node",
				Label:       "Network",
				Description: "connect to this network: (mainnet, testnet)",
				Widget:      "radio",
				Options: []string{"mainnet",
					"testnet",
					"regtestnet",
					"simnet",
				},
				// Hook:        "restart",
				OmitEmpty: true,
			},
			network,
		),
		NoCFilters: NewBool(
			metadata{
				Name:        "nocfilters",
				Group:       "node",
				Label:       "No CFilters",
				Description: "disable committed filtering (CF) support",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		NodeOff: NewBool(
			metadata{
				Name:        "nodeoff",
				Group:       "debug",
				Label:       "Node Off",
				Description: "turn off the node backend",
				Widget:      "toggle",
				// Hook:        "node",
				OmitEmpty: true,
			},
			false,
		),
		NoInitialLoad: NewBool(
			metadata{
				Name:        "noinitialload",
				Label:       "No Initial Load",
				Description: "do not load a wallet at startup",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		NoPeerBloomFilters: NewBool(
			metadata{
				Name:        "nopeerbloomfilters",
				Group:       "node",
				Label:       "No Peer Bloom Filters",
				Description: "disable bloom filtering support",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		NoRelayPriority: NewBool(
			metadata{
				Name:        "norelaypriority",
				Group:       "policy",
				Label:       "No Relay Priority",
				Description: "do not require free or low-fee transactions to have high priority for relaying",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		OneTimeTLSKey: NewBool(
			metadata{
				Name:        "onetimetlskey",
				Group:       "wallet",
				Label:       "One Time TLS Key",
				Description: "generate a new TLS certificate pair at startup, but only write the certificate to disk",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		Onion: NewBool(
			metadata{
				Name:        "onion",
				Group:       "proxy",
				Label:       "Onion Enabled",
				Description: "enable tor proxy",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		OnionProxy: NewString(
			metadata{
				Name:        "onionproxy",
				Group:       "proxy",
				Label:       "Onion Proxy Address",
				Description: "address of tor proxy you want to connect to",
				Type:        "address",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		OnionProxyPass: NewString(
			metadata{
				Name:        "onionproxypass",
				Group:       "proxy",
				Label:       "Onion Proxy Password",
				Description: "password for tor proxy",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		OnionProxyUser: NewString(
			metadata{
				Name:        "onionproxyuser",
				Group:       "proxy",
				Label:       "Onion Proxy Username",
				Description: "tor proxy username",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		P2PConnect: NewStrings(
			metadata{
				Name:        "p2pconnect",
				Group:       "node",
				Label:       "P2P Connect",
				Description: "list of addresses reachable from connected networks",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		P2PListeners: NewStrings(
			metadata{
				Name:        "listen",
				Group:       "node",
				Label:       "P2PListeners",
				Description: "list of addresses to bind the node listener to",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		Password: NewString(
			metadata{
				Name:        "password",
				Group:       "rpc",
				Label:       "Password",
				Description: "password for client RPC connections",
				Type:        "password",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			genPassword(),
		),
		PipeLog: NewBool(
			metadata{
				Name:        "pipelog",
				Label:       "Pipe Logger",
				Description: "enable pipe based logger IPC",
				Widget:      "toggle",
				// Hook:        "",
				OmitEmpty: true,
			},
			false,
		),
		Profile: NewString(
			metadata{
				Name:        "profile",
				Group:       "debug",
				Label:       "Profile",
				Description: "http profiling on given port (1024-40000)",
				// Type:        "",
				Widget: "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		Proxy: NewString(
			metadata{
				Name:        "proxy",
				Group:       "proxy",
				Label:       "Proxy",
				Description: "address of proxy to connect to for outbound connections",
				Type:        "url",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		ProxyPass: NewString(
			metadata{
				Name:        "proxypass",
				Group:       "proxy",
				Label:       "Proxy Pass",
				Description: "proxy password, if required",
				Type:        "password",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			genPassword(),
		),
		ProxyUser: NewString(
			metadata{
				Name:        "proxyuser",
				Group:       "proxy",
				Label:       "ProxyUser",
				Description: "proxy username, if required",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"proxyuser",
		),
		RejectNonStd: NewBool(
			metadata{
				Name:        "rejectnonstd",
				Group:       "node",
				Label:       "Reject Non Std",
				Description: "reject non-standard transactions regardless of the default settings for the active network",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		RelayNonStd: NewBool(
			metadata{
				Name:        "relaynonstd",
				Group:       "node",
				Label:       "Relay Nonstandard Transactions",
				Description: "relay non-standard transactions regardless of the default settings for the active network",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		RPCCert: NewString(
			metadata{
				Name:        "rpccert",
				Group:       "rpc",
				Label:       "RPC Cert",
				Description: "location of RPC TLS certificate",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.cert"),
		),
		RPCConnect: NewString(
			metadata{
				Name:        "rpcconnect",
				Group:       "wallet",
				Label:       "RPC Connect",
				Description: "full node RPC for wallet",
				Type:        "address",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.DefaultPort),
		
		),
		RPCKey: NewString(
			metadata{
				Name:        "rpckey",
				Group:       "rpc",
				Label:       "RPC Key",
				Description: "location of rpc TLS key",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.key"),
		),
		RPCListeners: NewStrings(
			metadata{
				Name:        "rpclisten",
				Group:       "rpc",
				Label:       "RPC Listeners",
				Description: "addresses to listen for RPC connections",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		RPCMaxClients: NewInt(
			metadata{
				Name:        "rpcmaxclients",
				Group:       "rpc",
				Label:       "Maximum RPC Clients",
				Description: "maximum number of clients for regular RPC",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMaxRPCClients,
		),
		RPCMaxConcurrentReqs: NewInt(
			metadata{
				Name:        "rpcmaxconcurrentreqs",
				Group:       "rpc",
				Label:       "Maximum RPC Concurrent Reqs",
				Description: "maximum number of requests to process concurrently",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMaxRPCConcurrentReqs,
		),
		RPCMaxWebsockets: NewInt(
			metadata{
				Name:        "rpcmaxwebsockets",
				Group:       "rpc",
				Label:       "Maximum RPC Websockets",
				Description: "maximum number of websocket clients to allow",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultMaxRPCWebsockets,
		),
		RPCQuirks: NewBool(
			metadata{
				Name:        "rpcquirks",
				Group:       "rpc",
				Label:       "RPC Quirks",
				Description: "enable bugs that replicate bitcoin core RPC's JSON",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		RunAsService: NewBool(
			metadata{
				Name:        "runasservice",
				Label:       "Run As Service",
				Description: "shuts down on lock timeout",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		ServerPass: NewString(
			metadata{
				Name:        "serverpass",
				Group:       "rpc",
				Label:       "Server Pass",
				Description: "password for server connections",
				Type:        "password",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			genPassword(),
		),
		ServerTLS: NewBool(
			metadata{
				Name:        "servertls",
				Group:       "wallet",
				Label:       "Server TLS",
				Description: "enable TLS for the wallet connection to node RPC server",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			true,
		),
		ServerUser: NewString(
			metadata{
				Name:        "serveruser",
				Group:       "rpc",
				Label:       "Server User",
				Description: "username for chain server connections",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"client",
		),
		SigCacheMaxSize: NewInt(
			metadata{
				Name:        "sigcachemaxsize",
				Group:       "node",
				Label:       "Signature Cache Max Size",
				Description: "the maximum number of entries in the signature verification cache",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultSigCacheMaxSize,
		),
		Solo: NewBool(
			metadata{
				Name:        "solo",
				Group:       "mining",
				Label:       "Solo Generate",
				Description: "mine even if not connected to a network",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		TLS: NewBool(
			metadata{
				Name:        "clienttls",
				Group:       "tls",
				Label:       "TLS",
				Description: "enable TLS for RPC client connections",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			true,
		),
		TLSSkipVerify: NewBool(
			metadata{
				Name:        "tlsskipverify",
				Group:       "tls",
				Label:       "TLS Skip Verify",
				Description: "skip TLS certificate verification (ignore CA errors)",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			false,
		),
		TorIsolation: NewBool(
			metadata{
				Name:        "torisolation",
				Group:       "proxy",
				Label:       "Tor Isolation",
				Description: "makes a separate proxy connection for each connection",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			true,
		),
		TrickleInterval: NewDuration(
			metadata{
				Name:        "trickleinterval",
				Group:       "policy",
				Label:       "Trickle Interval",
				Description: "minimum time between attempts to send new inventory to a connected peer",
				Widget:      "duration",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultTrickleInterval,
		),
		TxIndex: NewBool(
			metadata{
				Name:        "txindex",
				Group:       "node",
				Label:       "Tx Index",
				Description: "maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
				Widget:      "toggle",
				// Hook:        "droptxindex",
				OmitEmpty: true,
			},
			true,
		),
		UPNP: NewBool(
			metadata{
				Name:        "upnp",
				Group:       "node",
				Label:       "UPNP",
				Description: "enable UPNP for NAT traversal",
				Widget:      "toggle",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			true,
		),
		UserAgentComments: NewStrings(
			metadata{
				Name:        "uacomment",
				Group:       "policy",
				Label:       "User Agent Comments",
				Description: "comment to add to the user agent -- See BIP 14 for more information",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
		Username: NewString(
			metadata{
				Name:        "username",
				Group:       "rpc",
				Label:       "Username",
				Description: "password for client RPC connections",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"username",
		),
		UUID: &Int{
			value: uberatomic.NewInt64(rand.Int63()),
			metadata: metadata{
				Name:        "uuid",
				Label:       "UUID",
				Description: "instance unique id (64bit random value)",
				Widget:      "string",
				OmitEmpty:   true,
			},
		},
		Wallet: NewBool(
			metadata{
				Name:        "walletconnect",
				Group:       "debug",
				Label:       "Connect to Wallet",
				Description: "set ctl to connect to wallet instead of chain server",
				Widget:      "toggle",
				OmitEmpty:   true,
			},
			false,
		),
		WalletFile: NewString(
			metadata{
				Name:        "walletfile",
				Aliases:     []string{"WF"},
				Group:       "config",
				Label:       "Wallet File",
				Description: "wallet database file",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			filepath.Join(string(datadir.Load().([]byte)), "mainnet", DbName),
		),
		WalletOff: NewBool(
			metadata{
				Name:        "walletoff",
				Group:       "debug",
				Label:       "Wallet Off",
				Description: "turn off the wallet backend",
				Widget:      "toggle",
				// Hook:        "wallet",
				OmitEmpty: true,
			},
			false,
		),
		WalletPass: NewString(
			metadata{
				Name:        "walletpass",
				Label:       "Wallet Pass",
				Description: "password encrypting public data in wallet - hash is stored so give on command line",
				Type:        "password",
				Widget:      "password",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			"",
		),
		WalletRPCListeners: NewStrings(
			metadata{
				Name:        "walletrpclisten",
				Group:       "wallet",
				Label:       "Wallet RPC Listeners",
				Description: "addresses for wallet RPC server to listen on",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
			},
		),
		WalletRPCMaxClients: NewInt(
			metadata{
				Name:        "walletrpcmaxclients",
				Group:       "wallet",
				Label:       "Legacy RPC Max Clients",
				Description: "maximum number of RPC clients allowed for wallet RPC",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultRPCMaxClients,
		),
		WalletRPCMaxWebsockets: NewInt(
			metadata{
				Name:        "walletrpcmaxwebsockets",
				Group:       "wallet",
				Label:       "Legacy RPC Max Websockets",
				Description: "maximum number of websocket clients allowed for wallet RPC",
				Widget:      "integer",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			DefaultRPCMaxWebsockets,
		),
		WalletServer: NewString(
			metadata{
				Name:        "walletserver",
				Aliases:     []string{"ws"},
				Group:       "wallet",
				Label:       "Wallet Server",
				Description: "node address to connect wallet server to",
				Type:        "address",
				Widget:      "string",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
		),
		Whitelists: NewStrings(
			metadata{
				Name:        "whitelists",
				Group:       "debug",
				Label:       "Whitelists",
				Description: "peers that you don't want to ever ban",
				Type:        "address",
				Widget:      "multi",
				// Hook:        "restart",
				OmitEmpty: true,
			},
			[]string{},
		),
	}
	c.Map = make(map[string]interface{})
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Strings:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Float:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Int:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *String:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			case *Duration:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := c.Map[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				c.Map[im.Name] = ii
			default:
			}
			return true
		},
	)
	return
}

type (
	metadata struct {
		Name        string
		Aliases     []string
		Group       string
		Label       string
		Description string
		Type        string
		Widget      string
		Options     []string
		OmitEmpty   bool
	}
	Bool struct {
		metadata
		hook  []func(b bool)
		value *uberatomic.Bool
		def   bool
	}
	Strings struct {
		metadata
		hook  []func(s []string)
		value *atomic.Value
		def   []string
	}
	Float struct {
		metadata
		hook  []func(f float64)
		value *uberatomic.Float64
		def   float64
	}
	Int struct {
		metadata
		hook  []func(i int64)
		value *uberatomic.Int64
		def   int64
	}
	String struct {
		metadata
		hook  []func(s Strice)
		value *atomic.Value
		def   string
	}
	Duration struct {
		metadata
		hook  []func(d time.Duration)
		value *uberatomic.Duration
		def   time.Duration
	}
	// Strice is a wrapper around byte slices to enable optional security features and possibly better performance for
	// bulk comparison and editing. There isn't any extensive editing primitives for this purpose,
	Strice []byte
)

// S returns the underlying bytes converted into string
func (s *Strice) S() string {
	return string(*s)
}

// E returns the byte at the requested index in the string
func (s *Strice) E(elem int) byte {
	if s.Len() > elem {
		return (*s)[elem]
	}
	return 0
}

// Len returns the length of the string in bytes
func (s *Strice) Len() int {
	return len(*s)
}

// Equal returns true if two Strices are equal in both length and content
func (s *Strice) Equal(sb *Strice) bool {
	if s.Len() == sb.Len() {
		for i := range *s {
			if s.E(i) != sb.E(i) {
				return false
			}
		}
		return true
	}
	return false
}

// Cat two Strices together
func (s *Strice) Cat(sb *Strice) *Strice {
	*s = append(*s, *sb...)
	return s
}

// Find returns true if a match of a substring is found and if found, the position in the first string that the second
// string starts, the number of matching characters from the start of the search Strice, or -1 if not found.
//
// You specify a minimum length match and it will trawl through it systematically until it finds the first match of the
// minimum length.
func (s *Strice) Find(sb *Strice, minLengthMatch int) (found bool, extent, pos int) {
	// can't be a substring if it's longer
	if sb.Len() > s.Len() {
		return
	}
	for pos = range *s {
		// if we find a match, grab onto it
		if s.E(pos) == sb.E(pos) {
			extent++
			// this exhaustively searches for a match between the two strings, but we do not restrict the match to the
			// minimum, maximising the ways this function can be used for simple position tests and editing
			for srchPos := 1; srchPos < sb.Len() || srchPos+pos < s.Len(); srchPos++ {
				// the first element is skipped
				if s.E(srchPos+pos) != sb.E(srchPos) {
					break
				}
				extent++
			}
			// the above loop ends when the bytes stop matching, then if it is under the minimum length requested, it
			// continues. Note that we are not mutating `i` so it iterates for a match comprehensively.
			if extent < minLengthMatch {
				// reset the extent
				extent = 0
			} else {
				break
			}
		}
	}
	return
}

// HasPrefix returns true if the given string forms the beginning of the current string
func (s *Strice) HasPrefix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == 0 {
			return true
		}
	}
	return false
}

// HasSuffix returns true if the given string forms the ending of the current string
func (s *Strice) HasSuffix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == s.Len()-sb.Len()-1 {
			return true
		}
	}
	return false
}

// Dup copies a string and returns it
func (s *Strice) Dup() *Strice {
	ns := make(Strice, s.Len())
	copy(ns, *s)
	return &ns
}

// Wipe zeroes the bytes of a string
func (s *Strice) Wipe() {
	for i := range *s {
		(*s)[i] = 0
	}
}

// Split the string by a given cutset
func (s *Strice) Split(cutset string) (out []*Strice) {
	// convert immutable string type to Strice bytes
	c := Strice(cutset)
	// need the pointer to call the methods
	cs := &c
	// copy the bytes so we can guarantee the original is unmodified
	cp := s.Dup()
	for {
		// locate the next instance of the cutset
		found, _, pos := s.Find(cp, cp.Len())
		if found {
			// add the found section to the return slice
			before := (*s)[:pos+cp.Len()]
			out = append(out, &before)
			// trim off the prefix and cutslice from the working copy
			*cs = (*cs)[pos+cp.Len():]
			// continue to search for more instances of the cutset
			continue
		} else {
			// once we get not found, the searching is over and whatever we have, we return
			break
		}
	}
	return
}

// NewBool creates a new podcfg.Bool with default values set
func NewBool(m metadata, def bool, hook ...func(b bool)) *Bool {
	return &Bool{value: uberatomic.NewBool(def), metadata: m, def: def, hook: hook}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Bool) AddHooks(hook ...func(b bool)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Bool) SetHooks(hook ...func(b bool)) {
	x.hook = hook
}

// True returns whether the value is set to true (it returns the value)
func (x *Bool) True() bool {
	return x.value.Load()
}

// False returns whether the value is false (it returns the inverse of the value)
func (x *Bool) False() bool {
	return !x.value.Load()
}

// Flip changes the value to its opposite
func (x *Bool) Flip() {
	x.value.Toggle()
}

// Set changes the value currently stored
func (x *Bool) Set(b bool) *Bool {
	x.value.Store(b)
	return x
}

// String returns a string form of the value
func (x *Bool) String() string {
	return fmt.Sprint(x.True())
}

// T sets the value to true
func (x *Bool) T() *Bool {
	x.value.Store(true)
	return x
}

// F sets the value to false
func (x *Bool) F() *Bool {
	x.value.Store(false)
	return x
}

// MarshalJSON returns the json representation of a Bool
func (x *Bool) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of a Bool
func (x *Bool) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewStrings  creates a new podcfg.Strings with default values set
func NewStrings(m metadata, def []string, hook ...func(s []string)) *Strings {
	as := &atomic.Value{}
	v := cli.StringSlice(def)
	as.Store(&v)
	return &Strings{value: as, metadata: m, def: def, hook: hook}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Strings) AddHooks(hook ...func(b []string)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Strings) SetHooks(hook ...func(b []string)) {
	x.hook = hook
}

// V returns the stored value
func (x *Strings) V() *cli.StringSlice {
	return x.value.Load().(*cli.StringSlice)
}

// Len returns the length of the slice of strings
func (x *Strings) Len() int {
	return len(x.S())
}

// Set the slice of strings stored
func (x *Strings) Set(ss []string) *Strings {
	sss := cli.StringSlice(ss)
	x.value.Store(&sss)
	return x
}

// S returns the value as a slice of string
func (x *Strings) S() []string {
	return *x.value.Load().(*cli.StringSlice)
}

// String returns a string representation of the value
func (x *Strings) String() string {
	return fmt.Sprint(x.S())
}

// MarshalJSON returns the json representation of
func (x *Strings) MarshalJSON() (b []byte, e error) {
	xs := x.value.Load().(*cli.StringSlice)
	return json.Marshal(xs)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Strings) UnmarshalJSON(data []byte) (e error) {
	v := &cli.StringSlice{}
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewFloat returns a new Float value set to a default value
func NewFloat(m metadata, def float64) *Float {
	return &Float{value: uberatomic.NewFloat64(def), metadata: m, def: def}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Float) AddHooks(hook ...func(f float64)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Float) SetHooks(hook ...func(f float64)) {
	x.hook = hook
}

// V returns the value stored
func (x *Float) V() float64 {
	return x.value.Load()
}

// Set the value stored
func (x *Float) Set(f float64) *Float {
	x.value.Store(f)
	return x
}

// String returns a string representation of the value
func (x *Float) String() string {
	return fmt.Sprintf("%0.8f", x.V())
}

// MarshalJSON returns the json representation of
func (x *Float) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Float) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewInt creates a new Int with a given default value
func NewInt(m metadata, def int64) *Int {
	return &Int{value: uberatomic.NewInt64(def), metadata: m, def: def}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Int) AddHooks(hook ...func(f int64)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Int) SetHooks(hook ...func(f int64)) {
	x.hook = hook
}

// V returns the stored int
func (x *Int) V() int {
	return int(x.value.Load())
}

// Set the value stored
func (x *Int) Set(i int) *Int {
	x.value.Store(int64(i))
	return x
}

// String returns the string stored
func (x *Int) String() string {
	return fmt.Sprintf("%d", x.V())
}

// MarshalJSON returns the json representation of
func (x *Int) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Int) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewString creates a new String with a given default value set
func NewString(m metadata, def string) *String {
	v := &atomic.Value{}
	v.Store([]byte(def))
	return &String{value: v, metadata: m, def: def}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *String) AddHooks(hook ...func(f Strice)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *String) SetHooks(hook ...func(f Strice)) {
	x.hook = hook
}

// V returns the stored string
func (x *String) V() string {
	return string(x.value.Load().([]byte))
}

// Empty returns true if the string is empty
func (x *String) Empty() bool {
	return len(x.value.Load().([]byte)) == 0
}

// Bytes returns the raw bytes in the underlying storage
func (x *String) Bytes() []byte {
	return x.value.Load().([]byte)
}

// Set the value stored
func (x *String) Set(s string) *String {
	x.value.Store([]byte(s))
	return x
}

func (x *String) SetBytes(s []byte) *String {
	x.value.Store(s)
	return x
}

// String returns a string representation of the value
func (x *String) String() string {
	return x.V()
}

// MarshalJSON returns the json representation
func (x *String) MarshalJSON() (b []byte, e error) {
	v := string(x.value.Load().([]byte))
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *String) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load().([]byte)
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewDuration creates a new Duration with a given default value set
func NewDuration(m metadata, def time.Duration) *Duration {
	return &Duration{value: uberatomic.NewDuration(def), metadata: m, def: def}
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Duration) AddHooks(hook ...func(d time.Duration)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Duration) SetHooks(hook ...func(d time.Duration)) {
	x.hook = hook
}

// V returns the value stored
func (x *Duration) V() time.Duration {
	return x.value.Load()
}

// Set the value stored
func (x *Duration) Set(d time.Duration) *Duration {
	x.value.Store(d)
	return x
}

// String returns a string representation of the value
func (x *Duration) String() string {
	return fmt.Sprint(x.V())
}

// MarshalJSON returns the json representation
func (x *Duration) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *Duration) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
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
