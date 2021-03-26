package podcfg

import (
	"github.com/p9c/pod/pkg/appdata"
	"github.com/p9c/pod/pkg/chaincfg"
	uberatomic "go.uber.org/atomic"
	"math/rand"
	"net"
	"path/filepath"
	"sync/atomic"
	"time"
)

// New creates a fresh Config with default values stored in its fields
func New() (c *Config) {
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
