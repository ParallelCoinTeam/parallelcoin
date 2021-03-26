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
	AddCheckpoints         *Strings `json:"add-checkpoints,omitempty"`
	AddPeers               *Strings `json:"add-peers,omitempty"`
	AddrIndex              *Bool `json:"addr-index,omitempty"`
	AutoListen             *Bool `json:"auto-listen,omitempty"`
	AutoPorts              *Bool `json:"auto-ports,omitempty"`
	BanDuration            *Duration `json:"ban-duration,omitempty"`
	BanThreshold           *Int `json:"ban-threshold,omitempty"`
	BlockMaxSize           *Int `json:"block-max-size,omitempty"`
	BlockMaxWeight         *Int `json:"block-max-weight,omitempty"`
	BlockMinSize           *Int `json:"block-min-size,omitempty"`
	BlockMinWeight         *Int `json:"block-min-weight,omitempty"`
	BlockPrioritySize      *Int `json:"block-priority-size,omitempty"`
	BlocksOnly             *Bool `json:"blocks-only,omitempty"`
	CAFile                 *String `json:"cafile,omitempty"`
	ConfigFile             *String `json:"config-file,omitempty"`
	ConnectPeers           *Strings `json:"connect-peers,omitempty"`
	Controller             *Bool `json:"controller,omitempty"`
	CPUProfile             *String `json:"cpuprofile,omitempty"`
	DarkTheme              *Bool `json:"dark-theme,omitempty"`
	DataDir                *String `json:"data-dir,omitempty"`
	DbType                 *String `json:"db-type,omitempty"`
	DisableBanning         *Bool `json:"disable-banning,omitempty"`
	DisableCheckpoints     *Bool `json:"disable-checkpoints,omitempty"`
	DisableDNSSeed         *Bool `json:"disable-dnsseed,omitempty"`
	DisableListen          *Bool `json:"disable-listen,omitempty"`
	DisableRPC             *Bool `json:"disable-rpc,omitempty"`
	Discovery              *Bool `json:"discovery,omitempty"`
	ExternalIPs            *Strings `json:"external-ips,omitempty"`
	FreeTxRelayLimit       *Float `json:"free-tx-relay-limit,omitempty"`
	Generate               *Bool `json:"generate,omitempty"`
	GenThreads             *Int `json:"gen-threads,omitempty"`
	Hilite                 *Strings `json:"hilite,omitempty"`
	LAN                    *Bool `json:"lan,omitempty"`
	Language               *String `json:"language,omitempty"`
	LimitPass              *String `json:"limit-pass,omitempty"`
	LimitUser              *String `json:"limit-user,omitempty"`
	LogDir                 *String `json:"log-dir,omitempty"`
	LogFilter              *Strings `json:"log-filter,omitempty"`
	LogLevel               *String `json:"log-level,omitempty"`
	MaxOrphanTxs           *Int `json:"max-orphan-txs,omitempty"`
	MaxPeers               *Int `json:"max-peers,omitempty"`
	MulticastPass          *String `json:"multicast-pass,omitempty"`
	MiningAddrs            *Strings `json:"mining-addrs,omitempty"`
	MinRelayTxFee          *Float `json:"min-relay-tx-fee,omitempty"`
	Network                *String `json:"network,omitempty"`
	NoCFilters             *Bool `json:"no-cfilters,omitempty"`
	NodeOff                *Bool `json:"node-off,omitempty"`
	NoInitialLoad          *Bool `json:"no-initial-load,omitempty"`
	NoPeerBloomFilters     *Bool `json:"no-peer-bloom-filters,omitempty"`
	NoRelayPriority        *Bool `json:"no-relay-priority,omitempty"`
	OneTimeTLSKey          *Bool `json:"one-time-tlskey,omitempty"`
	Onion                  *Bool `json:"onion,omitempty"`
	OnionProxy             *String `json:"onion-proxy,omitempty"`
	OnionProxyPass         *String `json:"onion-proxy-pass,omitempty"`
	OnionProxyUser         *String `json:"onion-proxy-user,omitempty"`
	P2PConnect             *Strings `json:"p-2-pconnect,omitempty"`
	P2PListeners           *Strings `json:"p-2-plisteners,omitempty"`
	Password               *String `json:"password,omitempty"`
	PipeLog                *Bool `json:"pipe-log,omitempty"`
	Profile                *String `json:"profile,omitempty"`
	Proxy                  *String `json:"proxy,omitempty"`
	ProxyPass              *String `json:"proxy-pass,omitempty"`
	ProxyUser              *String `json:"proxy-user,omitempty"`
	RejectNonStd           *Bool `json:"reject-non-std,omitempty"`
	RelayNonStd            *Bool `json:"relay-non-std,omitempty"`
	RPCCert                *String `json:"rpccert,omitempty"`
	RPCConnect             *String `json:"rpcconnect,omitempty"`
	RPCKey                 *String `json:"rpckey,omitempty"`
	RPCListeners           *Strings `json:"rpclisteners,omitempty"`
	RPCMaxClients          *Int `json:"rpcmax-clients,omitempty"`
	RPCMaxConcurrentReqs   *Int `json:"rpcmax-concurrent-reqs,omitempty"`
	RPCMaxWebsockets       *Int `json:"rpcmax-websockets,omitempty"`
	RPCQuirks              *Bool `json:"rpcquirks,omitempty"`
	RunAsService           *Bool `json:"run-as-service,omitempty"`
	ServerPass             *String `json:"server-pass,omitempty"`
	ServerTLS              *Bool `json:"server-tls,omitempty"`
	ServerUser             *String `json:"server-user,omitempty"`
	SigCacheMaxSize        *Int `json:"sig-cache-max-size,omitempty"`
	Solo                   *Bool `json:"solo,omitempty"`
	TLS                    *Bool `json:"tls,omitempty"`
	TLSSkipVerify          *Bool `json:"tlsskip-verify,omitempty"`
	TorIsolation           *Bool `json:"tor-isolation,omitempty"`
	TrickleInterval        *Duration `json:"trickle-interval,omitempty"`
	TxIndex                *Bool `json:"tx-index,omitempty"`
	UPNP                   *Bool `json:"upnp,omitempty"`
	UserAgentComments      *Strings `json:"user-agent-comments,omitempty"`
	Username               *String `json:"username,omitempty"`
	UUID                   *Int `json:"uuid,omitempty"`
	Wallet                 *Bool `json:"wallet,omitempty"`
	WalletFile             *String `json:"wallet-file,omitempty"`
	WalletOff              *Bool `json:"wallet-off,omitempty"`
	WalletPass             *String `json:"wallet-pass,omitempty"`
	WalletRPCListeners     *Strings `json:"wallet-rpclisteners,omitempty"`
	WalletRPCMaxClients    *Int `json:"wallet-rpcmax-clients,omitempty"`
	WalletRPCMaxWebsockets *Int `json:"wallet-rpcmax-websockets,omitempty"`
	WalletServer           *String `json:"wallet-server,omitempty"`
	Whitelists             *Strings `json:"whitelists,omitempty"`
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

func EmptyConfig() (c *Config, conf map[string]interface{}) {
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
			},
			[]string{},
		),
		AddrIndex: NewBool(
			metadata{
				Name:        "addrindex",
				Group:       "node",
				Label:       "Address Index",
				Description: "maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
				Widget:      "toggle",
				// Hook:        "dropaddrindex",
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
			},
			[]string{},
		),
		LAN: NewBool(
			metadata{
				Name:        "LAN",
				Group:       "debug",
				Label:       "LAN Testnet Mode",
				Description: "run without any connection to nodes on the internet (does not apply on mainnet)",
				Widget:      "toggle",
				// Hook:        "restart",
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
				Options:     []string{"off", "fatal", "error", "info", "check", "debug", "trace"},
				// Hook:        "loglevel",
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
				Options:     []string{"mainnet", "testnet", "regtestnet", "simnet"},
				// Hook:        "restart",
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
			},
			[]string{net.JoinHostPort("0.0.0.0", chaincfg.MainNetParams.DefaultPort)},
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
			},
			[]string{net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.DefaultPort)},
		),
		RPCMaxClients: NewInt(
			metadata{
				Name:        "rpcmaxclients",
				Group:       "rpc",
				Label:       "Maximum RPC Clients",
				Description: "maximum number of clients for regular RPC",
				Widget:      "integer",
				// Hook:        "restart",
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
			},
			"username",
		),
		UUID: &Int{
			Int64: uberatomic.NewInt64(rand.Int63()),
			metadata: metadata{
				Name:        "uuid",
				Label:       "UUID",
				Description: "instance unique id (64bit random value)",
				Widget:      "string",
			},
		},
		Wallet: NewBool(
			metadata{
				Name:        "walletconnect",
				Group:       "debug",
				Label:       "Connect to Wallet",
				Description: "set ctl to connect to wallet instead of chain server",
				Widget:      "toggle",
			},
			false,
		),
		WalletFile: NewString(
			metadata{
				Name:        "wallet-file",
				Aliases:     []string{"WF"},
				Group:       "config",
				Label:       "Wallet File",
				Description: "wallet database file",
				Type:        "path",
				Widget:      "string",
				// Hook:        "restart",
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
			},
			[]string{net.JoinHostPort("0.0.0.0", chaincfg.MainNetParams.WalletRPCServerPort)},
		),
		WalletRPCMaxClients: NewInt(
			metadata{
				Name:        "walletrpcmaxclients",
				Group:       "wallet",
				Label:       "Legacy RPC Max Clients",
				Description: "maximum number of RPC clients allowed for wallet RPC",
				Widget:      "integer",
				// Hook:        "restart",
			},
			DefaultRPCMaxClients,
		),
		WalletRPCMaxWebsockets: NewInt(
			metadata{
				Name:        "wallet-rpc-max-websockets",
				Group:       "wallet",
				Label:       "Legacy RPC Max Websockets",
				Description: "maximum number of websocket clients allowed for wallet RPC",
				Widget:      "integer",
				// Hook:        "restart",
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
			},
			net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.WalletRPCServerPort),
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
			},
			[]string{},
		),
	}
	conf = make(map[string]interface{})
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
			case *Strings:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
			case *Float:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
			case *Int:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
			case *String:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
			case *Duration:
				im := ii.metadata
				// sanity check for programmers
				if _, ok := conf[im.Name]; ok {
					panic("duplicate configuration item name")
				}
				conf[im.Name] = ii
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
	}
	Bool struct {
		metadata
		hook func(b bool)
		value *uberatomic.Bool
		def bool
	}
	Strings struct {
		metadata
		hook func(s []string)
		value *atomic.Value
		def []string
	}
	Float struct {
		metadata
		hook func(f float64)
		value *uberatomic.Float64
		def float64
	}
	Int struct {
		metadata
		hook func(i int64)
		*uberatomic.Int64
		def int64
	}
	String struct {
		metadata
		hook func(s string)
		value *atomic.Value
		def string
	}
	Duration struct {
		metadata
		hook func(d time.Duration)
		value *uberatomic.Duration
		def time.Duration
	}
)

func NewBool(m metadata, def bool) *Bool {
	return &Bool{value: uberatomic.NewBool(def), metadata: m, def: def}
}
func (x *Bool) True() bool {
	return x.value.Load()
}
func (x *Bool) False() bool {
	return !x.value.Load()
}
func (x *Bool) Flip() {
	x.value.Toggle()
}
func (x *Bool) Set(b bool) *Bool {
	x.value.Store(b)
	return x
}
func (x *Bool) T() *Bool {
	x.value.Store(true)
	return x
}
func (x *Bool) F() *Bool {
	x.value.Store(false)
	return x
}
func (x *Bool) String() string {
	return fmt.Sprint(x.value.Load())
}
func (x *Bool) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	if v == x.def {
		return json.Marshal(nil)
	}
	return json.Marshal(&v)
}
func (x *Bool) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

func NewStrings(m metadata, def []string) *Strings {
	as := &atomic.Value{}
	v := cli.StringSlice(def)
	as.Store(&v)
	return &Strings{value: as, metadata: m, def: def}
}
func (x *Strings) V() *cli.StringSlice {
	return x.value.Load().(*cli.StringSlice)
}
func (x *Strings) Len() int {
	return len(x.S())
}
func (x *Strings) Set(ss []string) *Strings {
	sss := cli.StringSlice(ss)
	x.value.Store(&sss)
	return x
}
func (x *Strings) S() []string {
	return *x.value.Load().(*cli.StringSlice)
}
func (x *Strings) MarshalJSON() (b []byte, e error) {
	xs := x.value.Load().(*cli.StringSlice)
	return json.Marshal(xs)
}
func (x *Strings) UnmarshalJSON(data []byte) (e error) {
	v := &cli.StringSlice{}
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

func NewFloat(m metadata, def float64) *Float {
	return &Float{value: uberatomic.NewFloat64(def), metadata: m, def: def}
}
func (x *Float) V() float64 {
	return x.value.Load()
}
func (x *Float) Set(f float64) *Float {
	x.value.Store(f)
	return x
}
func (x *Float) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}
func (x *Float) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

func NewInt(m metadata, def int64) *Int {
	return &Int{Int64: uberatomic.NewInt64(def), metadata: m, def: def}
}
func (x *Int) V() int {
	return int(x.Load())
}
func (x *Int) Set(i int) *Int {
	x.Store(int64(i))
	return x
}
func (x *Int) MarshalJSON() (b []byte, e error) {
	v := x.Load()
	return json.Marshal(&v)
}
func (x *Int) UnmarshalJSON(data []byte) (e error) {
	v := x.Load()
	e = json.Unmarshal(data, &v)
	x.Store(v)
	return
}

func NewString(m metadata, def string) *String {
	v := &atomic.Value{}
	v.Store([]byte{})
	return &String{value: v, metadata: m, def: def}
}
func (x *String) V() string {
	return string(x.value.Load().([]byte))
}
func (x *String) Empty() bool {
	return len(x.value.Load().([]byte)) == 0
}
func (x *String) Bytes() []byte {
	return x.value.Load().([]byte)
}
func (x *String) Set(s string) *String {
	x.value.Store([]byte(s))
	return x
}
func (x *String) SetBytes(s []byte) *String {
	x.value.Store(s)
	return x
}
func (x *String) String() string {
	return string(x.value.Load().([]byte))
}
func (x *String) MarshalJSON() (b []byte, e error) {
	v := x.value.Load().([]byte)
	return json.Marshal(&v)
}
func (x *String) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load().([]byte)
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

func NewDuration(m metadata, def time.Duration) *Duration {
	return &Duration{value: uberatomic.NewDuration(def), metadata: m, def: def}
}
func (x *Duration) V() time.Duration {
	return x.value.Load()
}
func (x *Duration) Set(d time.Duration) *Duration {
	x.value.Store(d)
	return x
}
func (x *Duration) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}
func (x *Duration) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

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
