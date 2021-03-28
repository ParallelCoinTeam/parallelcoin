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

// GetDefaultConfig returns a Config struct pristine factory freshaoeu
func GetDefaultConfig() (c *Config) {
	network := "mainnet"
	rand.Seed(time.Now().Unix())
	var datadir = &atomic.Value{}
	datadir.Store([]byte(appdata.Dir(Name, false)))
	c = &Config{
		Commands: Commands{
			{Name: "gui", Description:
			"ParallelCoin GUI Wallet/Miner/Explorer",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "version", Description:
			"print version and exit",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "ctl", Description:
			"command line wallet and chain RPC client",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "node", Description:
			"ParallelCoin blockchain node",
				Entrypoint: func(c *Config) error { return nil },
				Commands: []Command{
					{Name: "dropaddrindex", Description:
					"drop the address database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "droptxindex", Description:
					"drop the transaction database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "dropcfindex", Description:
					"drop the cfilter database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "dropindexes", Description:
					"drop all of the indexes",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "resetchain", Description:
					"deletes the current blockchain cache to force redownload",
						Entrypoint: func(c *Config) error { return nil },
					},
				},
			},
			{Name: "wallet", Description:
			"run the wallet server (requires a chain node to function)",
				Entrypoint: func(c *Config) error { return nil },
				Commands: []Command{
					{Name: "drophistory", Description:
					"reset the wallet transaction history",
						Entrypoint: func(c *Config) error { return nil },
					},
				},
			},
			{Name: "kopach", Description:
			"standalone multicast miner for easy mining farm deployment",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "worker", Description:
			"single thread worker process, normally started by kopach",
				Entrypoint: func(c *Config) error { return nil },
			},
		},
		AddCheckpoints: NewStrings(Metadata{
			Option: "addcheckpoint",
			Group:  "debug",
			Label:  "Add Checkpoints",
			Description:
			"add custom checkpoints",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		AddPeers: NewStrings(Metadata{
			Option:  "addpeer",
			Aliases: []string{"a"},
			Group:   "node",
			Label:   "Add Peers",
			Description:
			"manually adds addresses to try to connect to",
			Type:   "ipaddress",
			Widget: "multi",
			// Hook:        "addpeer",
			OmitEmpty: true,
		},
			[]string{},
			// []string{"127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12344"},
		),
		AddrIndex: NewBool(Metadata{
			Option: "addrindex",
			Group:  "node",
			Label:  "Address Index",
			Description:
			"maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
			Widget: "toggle",
			// Hook:        "dropaddrindex",
			OmitEmpty: true,
		},
			true,
		),
		AutoPorts: NewBool(Metadata{
			Option: "autoports",
			Group:  "debug",
			Label:  "Automatic Ports",
			Description:
			"RPC and controller ports are randomized, use with controller for automatic peer discovery",
			Widget: "toggle",
			// Hook: "restart",
			OmitEmpty: true,
		},
			false,
		),
		AutoListen: NewBool(Metadata{
			Option: "autolisten",
			Group:  "node",
			Label:  "Manual Listeners",
			Description:
			"automatically update inbound addresses dynamically according to discovered network interfaces",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		BanDuration: NewDuration(Metadata{
			Option: "banduration",
			Group:  "debug",
			Label:  "Ban Duration",
			Description:
			"how long a ban of a misbehaving peer lasts",
			Widget: "duration",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			time.Hour*24,
		),
		BanThreshold: NewInt(Metadata{
			Option: "banthreshold",
			Group:  "debug",
			Label:  "Ban Threshold",
			Description:
			"ban score that triggers a ban (default 100)",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultBanThreshold,
		),
		BlockMaxSize: NewInt(Metadata{
			Option: "blockmaxsize",
			Group:  "mining",
			Label:  "Block Max Size",
			Description:
			"maximum block size in bytes to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxSizeMax,
		),
		BlockMaxWeight: NewInt(Metadata{
			Option: "blockmaxweight",
			Group:  "mining",
			Label:  "Block Max Weight",
			Description:
			"maximum block weight to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxWeightMax,
		),
		BlockMinSize: NewInt(Metadata{
			Option: "blockminsize",
			Group:  "mining",
			Label:  "Block Min Size",
			Description:
			"minimum block size in bytes to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxSizeMin,
		),
		BlockMinWeight: NewInt(Metadata{
			Option: "blockminweight",
			Group:  "mining",
			Label:  "Block Min Weight",
			Description:
			"minimum block weight to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxWeightMin,
		),
		BlockPrioritySize: NewInt(Metadata{
			Option: "blockprioritysize",
			Group:  "mining",
			Label:  "Block Priority Size",
			Description:
			"size in bytes for high-priority/low-fee transactions when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultBlockPrioritySize,
		),
		BlocksOnly: NewBool(Metadata{
			Option: "blocksonly",
			Group:  "node",
			Label:  "Blocks Only",
			Description:
			"do not accept transactions from remote peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		CAFile: NewString(Metadata{
			Option: "cafile",
			Group:  "tls",
			Label:  "Certificate Authority File",
			Description:
			"certificate authority file for TLS certificate validation",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "ca.cert"),
		),
		ConfigFile: NewString(Metadata{
			Option:  "configfile",
			Aliases: []string{"C"},
			Label:   "Configuration File",
			Description:
			"location of configuration file, cannot actually be changed",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), PodConfigFilename),
		),
		ConnectPeers: NewStrings(Metadata{
			Option:  "connect",
			Aliases: []string{"c"},
			Group:   "node",
			Label:   "Connect Peers",
			Description:
			"connect ONLY to these addresses (disables inbound connections)",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		Controller: NewBool(Metadata{
			Option: "controller",
			Group:  "node",
			Label:  "Enable Controller",
			Description:
			"delivers mining jobs over multicast",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		CPUProfile: NewString(Metadata{
			Option: "cpuprofile",
			Group:  "debug",
			Label:  "CPU Profile",
			Description:
			"write cpu profile to this file",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		DarkTheme: NewBool(Metadata{
			Option: "darktheme",
			Group:  "config",
			Label:  "Dark Theme",
			Description:
			"sets dark theme for GUI",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DataDir: &String{
			value: datadir,
			Metadata: Metadata{
				Option:  "datadir",
				Aliases: []string{"D"},
				Label:   "Data Directory",
				Description:
				"root folder where application data is stored",
				Type:      "directory",
				Widget:    "string",
				OmitEmpty: true,
			},
			def: appdata.Dir(Name, false),
		},
		DbType: NewString(Metadata{
			Option: "dbtype",
			Group:  "debug",
			Label:  "Database Type",
			Description:
			"type of database storage engine to use (only one right now, ffldb)",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultDbType,
		),
		DisableBanning: NewBool(Metadata{
			Option: "nobanning",
			Group:  "debug",
			Label:  "Disable Banning",
			Description:
			"disables banning of misbehaving peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableCheckpoints: NewBool(Metadata{
			Option: "nocheckpoints",
			Group:  "debug",
			Label:  "Disable Checkpoints",
			Description:
			"disables all checkpoints",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableDNSSeed: NewBool(Metadata{
			Option: "nodnsseed",
			Group:  "node",
			Label:  "Disable DNS Seed",
			Description:
			"disable seeding of addresses to peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableListen: NewBool(Metadata{
			Option: "nolisten",
			Group:  "node",
			Label:  "Disable Listen",
			Description:
			"disables inbound connections for the peer to peer network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableRPC: NewBool(Metadata{
			Option: "norpc",
			Group:  "rpc",
			Label:  "Disable RPC",
			Description:
			"disable rpc servers, as well as kopach controller",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Discovery: NewBool(Metadata{
			Option: "discover",
			Group:  "node",
			Label:  "Disovery",
			Description:
			"enable LAN peer discovery in GUI",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		ExternalIPs: NewStrings(Metadata{
			Option: "externalip",
			Group:  "node",
			Label:  "External IP Addresses",
			Description:
			"extra addresses to tell peers they can connect to",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		FreeTxRelayLimit: NewFloat(Metadata{
			Option: "limitfreerelay",
			Group:  "policy",
			Label:  "Free Tx Relay Limit",
			Description:
			"limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
			Widget: "float",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultFreeTxRelayLimit,
		),
		Generate: NewBool(Metadata{
			Option:  "generate",
			Aliases: []string{"g"},
			Group:   "mining",
			Label:   "Generate Blocks",
			Description:
			"turn on Kopach CPU miner",
			Widget: "toggle",
			// Hook:        "generate",
			OmitEmpty: true,
		},
			false,
		),
		GenThreads: NewInt(Metadata{
			Option:  "genthreads",
			Aliases: []string{"G"},
			Group:   "mining",
			Label:   "Generate Threads",
			Description:
			"number of threads to mine with",
			Widget: "integer",
			// Hook:        "genthreads",
			OmitEmpty: true,
		},
			-1,
		),
		Hilite: NewStrings(Metadata{
			Option: "highlight",
			Group:  "debug",
			Label:  "Hilite",
			Description:
			"list of packages that will print with attention getters",
			Type:   "string",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		LAN: NewBool(Metadata{
			Option: "lan",
			Group:  "debug",
			Label:  "LAN Testnet Mode",
			Description:
			"run without any connection to nodes on the internet (does not apply on mainnet)",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Language: NewString(Metadata{
			Option: "language",
			Group:  "config",
			Label:  "Language",
			Description:
			"user interface language i18 localization",
			Widget: "string",
			// Hook:        "language",
			OmitEmpty: true,
		},
			"en",
		),
		LimitPass: NewString(Metadata{
			Option: "limitpass",
			Group:  "rpc",
			Label:  "Limit Password",
			Description:
			"limited user password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		LimitUser: NewString(Metadata{
			Option: "limituser",
			Group:  "rpc",
			Label:  "Limit Username",
			Description:
			"limited user name",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"limit",
		),
		LogDir: NewString(Metadata{
			Option: "logdir",
			Group:  "config",
			Label:  "Log Directory",
			Description:
			"folder where log files are written",
			Type:   "directory",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			string(datadir.Load().([]byte)),
		),
		LogFilter: NewStrings(Metadata{
			Option:  "logfilter",
			Aliases: []string{"L"},
			Group:   "debug",
			Label:   "Log Filter",
			Description:
			"list of packages that will not print logs",
			Type:   "string",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		LogLevel: NewString(Metadata{
			Option:  "loglevel",
			Aliases: []string{"l"},
			Group:   "config",
			Label:   "Log Level",
			Description:
			"maximum log level to output\n(fatal error check warning info debug trace - what is selected includes all items to the left of the one in that list)",
			Widget: "radio",
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
		MaxOrphanTxs: NewInt(Metadata{
			Option: "maxorphantx",
			Group:  "policy",
			Label:  "Max Orphan Txs",
			Description:
			"max number of orphan transactions to keep in memory",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxOrphanTransactions,
		),
		MaxPeers: NewInt(Metadata{
			Option: "maxpeers",
			Group:  "node",
			Label:  "Max Peers",
			Description:
			"maximum number of peers to hold connections with",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxPeers,
		),
		MulticastPass: NewString(Metadata{
			Option: "minerpass",
			Group:  "config",
			Label:  "Multicast Pass",
			Description:
			"password that encrypts the connection to the mining controller",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"pa55word",
		),
		MiningAddrs: NewStrings(Metadata{
			Option: "miningaddrs",
			Label:  "Mining Addresses",
			Description:
			"addresses to pay block rewards to (not in use)",
			Type:   "base58",
			Widget: "multi",
			// Hook:        "miningaddr",
			OmitEmpty: true,
		},
			[]string{},
		),
		MinRelayTxFee: NewFloat(Metadata{
			Option: "minrelaytxfee",
			Group:  "policy",
			Label:  "Min Relay Transaction Fee",
			Description:
			"the minimum transaction fee in DUO/kB to be considered a non-zero fee",
			Widget: "float",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMinRelayTxFee.ToDUO(),
		),
		Network: NewString(Metadata{
			Option: "network",
			Group:  "node",
			Label:  "Network",
			Description:
			"connect to this network: (mainnet, testnet)",
			Widget: "radio",
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
		NoCFilters: NewBool(Metadata{
			Option: "nocfilters",
			Group:  "node",
			Label:  "No CFilters",
			Description:
			"disable committed filtering (CF) support",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NodeOff: NewBool(Metadata{
			Option: "nodeoff",
			Group:  "debug",
			Label:  "Node Off",
			Description:
			"turn off the node backend",
			Widget: "toggle",
			// Hook:        "node",
			OmitEmpty: true,
		},
			false,
		),
		NoInitialLoad: NewBool(Metadata{
			Option: "noinitialload",
			Label:  "No Initial Load",
			Description:
			"do not load a wallet at startup",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NoPeerBloomFilters: NewBool(Metadata{
			Option: "nopeerbloomfilters",
			Group:  "node",
			Label:  "No Peer Bloom Filters",
			Description:
			"disable bloom filtering support",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NoRelayPriority: NewBool(Metadata{
			Option: "norelaypriority",
			Group:  "policy",
			Label:  "No Relay Priority",
			Description:
			"do not require free or low-fee transactions to have high priority for relaying",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		OneTimeTLSKey: NewBool(Metadata{
			Option: "onetimetlskey",
			Group:  "wallet",
			Label:  "One Time TLS Key",
			Description:
			"generate a new TLS certificate pair at startup, but only write the certificate to disk",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Onion: NewBool(Metadata{
			Option: "onion",
			Group:  "proxy",
			Label:  "Onion Enabled",
			Description:
			"enable tor proxy",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		OnionProxy: NewString(Metadata{
			Option: "onionproxy",
			Group:  "proxy",
			Label:  "Onion Proxy Address",
			Description:
			"address of tor proxy you want to connect to",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		OnionProxyPass: NewString(Metadata{
			Option: "onionproxypass",
			Group:  "proxy",
			Label:  "Onion Proxy Password",
			Description:
			"password for tor proxy",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		OnionProxyUser: NewString(Metadata{
			Option: "onionproxyuser",
			Group:  "proxy",
			Label:  "Onion Proxy Username",
			Description:
			"tor proxy username",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		P2PConnect: NewStrings(Metadata{
			Option: "p2pconnect",
			Group:  "node",
			Label:  "P2P Connect",
			Description:
			"list of addresses reachable from connected networks",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		P2PListeners: NewStrings(Metadata{
			Option:  "listen",
			Aliases: []string{"L"},
			Group:   "node",
			Label:   "P2PListeners",
			Description:
			"list of addresses to bind the node listener to",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		Password: NewString(Metadata{
			Option:  "password",
			Aliases: []string{"p"},
			Group:   "rpc",
			Label:   "Password",
			Description:
			"password for client RPC connections",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		PipeLog: NewBool(Metadata{
			Option: "pipelog",
			Label:  "Pipe Logger",
			Description:
			"enable pipe based logger IPC",
			Widget: "toggle",
			// Hook:        "",
			OmitEmpty: true,
		},
			false,
		),
		Profile: NewString(Metadata{
			Option: "profile",
			Group:  "debug",
			Label:  "Profile",
			Description:
			"http profiling on given port (1024-40000)",
			// Type:        "",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		Proxy: NewString(Metadata{
			Option: "proxy",
			Group:  "proxy",
			Label:  "Proxy",
			Description:
			"address of proxy to connect to for outbound connections",
			Type:   "url",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		ProxyPass: NewString(Metadata{
			Option: "proxypass",
			Group:  "proxy",
			Label:  "Proxy Pass",
			Description:
			"proxy password, if required",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		ProxyUser: NewString(Metadata{
			Option: "proxyuser",
			Group:  "proxy",
			Label:  "ProxyUser",
			Description:
			"proxy username, if required",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"proxyuser",
		),
		RejectNonStd: NewBool(Metadata{
			Option: "rejectnonstd",
			Group:  "node",
			Label:  "Reject Non Std",
			Description:
			"reject non-standard transactions regardless of the default settings for the active network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RelayNonStd: NewBool(Metadata{
			Option: "relaynonstd",
			Group:  "node",
			Label:  "Relay Nonstandard Transactions",
			Description:
			"relay non-standard transactions regardless of the default settings for the active network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RPCCert: NewString(Metadata{
			Option: "rpccert",
			Group:  "rpc",
			Label:  "RPC Cert",
			Description:
			"location of RPC TLS certificate",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.cert"),
		),
		RPCConnect: NewString(Metadata{
			Option:  "rpcconnect",
			Aliases: []string{"R"},
			Group:   "wallet",
			Label:   "RPC Connect",
			Description:
			"full node RPC for wallet",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.DefaultPort),
		
		),
		RPCKey: NewString(Metadata{
			Option: "rpckey",
			Group:  "rpc",
			Label:  "RPC Key",
			Description:
			"location of rpc TLS key",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.key"),
		),
		RPCListeners: NewStrings(Metadata{
			Option:  "rpclisten",
			Aliases: []string{"r"},
			Group:   "rpc",
			Label:   "RPC Listeners",
			Description:
			"addresses to listen for RPC connections",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		RPCMaxClients: NewInt(Metadata{
			Option: "rpcmaxclients",
			Group:  "rpc",
			Label:  "Maximum RPC Clients",
			Description:
			"maximum number of clients for regular RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCClients,
		),
		RPCMaxConcurrentReqs: NewInt(Metadata{
			Option: "rpcmaxconcurrentreqs",
			Group:  "rpc",
			Label:  "Maximum RPC Concurrent Reqs",
			Description:
			"maximum number of requests to process concurrently",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCConcurrentReqs,
		),
		RPCMaxWebsockets: NewInt(Metadata{
			Option: "rpcmaxwebsockets",
			Group:  "rpc",
			Label:  "Maximum RPC Websockets",
			Description:
			"maximum number of websocket clients to allow",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCWebsockets,
		),
		RPCQuirks: NewBool(Metadata{
			Option: "rpcquirks",
			Group:  "rpc",
			Label:  "RPC Quirks",
			Description:
			"enable bugs that replicate bitcoin core RPC's JSON",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RunAsService: NewBool(Metadata{
			Option: "runasservice",
			Label:  "Run As Service",
			Description:
			"shuts down on lock timeout",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		ServerPass: NewString(Metadata{
			Option: "serverpass",
			Group:  "rpc",
			Label:  "Server Pass",
			Description:
			"password for server connections",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		ServerTLS: NewBool(Metadata{
			Option: "servertls",
			Group:  "wallet",
			Label:  "Server TLS",
			Description:
			"enable TLS for the wallet connection to node RPC server",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		ServerUser: NewString(Metadata{
			Option: "serveruser",
			Group:  "rpc",
			Label:  "Server User",
			Description:
			"username for chain server connections",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"client",
		),
		SigCacheMaxSize: NewInt(Metadata{
			Option: "sigcachemaxsize",
			Group:  "node",
			Label:  "Signature Cache Max Size",
			Description:
			"the maximum number of entries in the signature verification cache",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultSigCacheMaxSize,
		),
		Solo: NewBool(Metadata{
			Option: "solo",
			Group:  "mining",
			Label:  "Solo Generate",
			Description:
			"mine even if not connected to a network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		TLS: NewBool(Metadata{
			Option: "clienttls",
			Group:  "tls",
			Label:  "TLS",
			Description:
			"enable TLS for RPC client connections",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		TLSSkipVerify: NewBool(Metadata{
			Option: "tlsskipverify",
			Group:  "tls",
			Label:  "TLS Skip Verify",
			Description:
			"skip TLS certificate verification (ignore CA errors)",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		TorIsolation: NewBool(Metadata{
			Option: "torisolation",
			Group:  "proxy",
			Label:  "Tor Isolation",
			Description:
			"makes a separate proxy connection for each connection",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		TrickleInterval: NewDuration(Metadata{
			Option: "trickleinterval",
			Group:  "policy",
			Label:  "Trickle Interval",
			Description:
			"minimum time between attempts to send new inventory to a connected peer",
			Widget: "duration",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultTrickleInterval,
		),
		TxIndex: NewBool(Metadata{
			Option: "txindex",
			Group:  "node",
			Label:  "Tx Index",
			Description:
			"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
			Widget: "toggle",
			// Hook:        "droptxindex",
			OmitEmpty: true,
		},
			true,
		),
		UPNP: NewBool(Metadata{
			Option: "upnp",
			Group:  "node",
			Label:  "UPNP",
			Description:
			"enable UPNP for NAT traversal",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		UserAgentComments: NewStrings(Metadata{
			Option: "uacomment",
			Group:  "policy",
			Label:  "User Agent Comments",
			Description:
			"comment to add to the user agent -- See BIP 14 for more information",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		Username: NewString(Metadata{
			Option:  "username",
			Aliases: []string{"u"},
			Group:   "rpc",
			Label:   "Username",
			Description:
			"password for client RPC connections",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"username",
		),
		UUID: &Int{Metadata: Metadata{
			Option: "uuid",
			Label:  "UUID",
			Description:
			"instance unique id (64bit random value)",
			Widget:    "string",
			OmitEmpty: true,
		},
			value: uberatomic.NewInt64(rand.Int63()),
		},
		Wallet: NewBool(Metadata{
			Option: "walletconnect",
			Group:  "debug",
			Label:  "Connect to Wallet",
			Description:
			"set ctl to connect to wallet instead of chain server",
			Widget:    "toggle",
			OmitEmpty: true,
		},
			false,
		),
		WalletFile: NewString(Metadata{
			Option:  "walletfile",
			Aliases: []string{"W"},
			Group:   "config",
			Label:   "Wallet File",
			Description:
			"wallet database file",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "mainnet", DbName),
		),
		WalletOff: NewBool(Metadata{
			Option: "walletoff",
			Group:  "debug",
			Label:  "Wallet Off",
			Description:
			"turn off the wallet backend",
			Widget: "toggle",
			// Hook:        "wallet",
			OmitEmpty: true,
		},
			false,
		),
		WalletPass: NewString(Metadata{
			Option: "walletpass",
			Label:  "Wallet Pass",
			Description:
			"password encrypting public data in wallet - hash is stored so give on command line",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		WalletRPCListeners: NewStrings(Metadata{
			Option: "walletrpclisten",
			Group:  "wallet",
			Label:  "Wallet RPC Listeners",
			Description:
			"addresses for wallet RPC server to listen on",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
			},
		),
		WalletRPCMaxClients: NewInt(Metadata{
			Option: "walletrpcmaxclients",
			Group:  "wallet",
			Label:  "Legacy RPC Max Clients",
			Description:
			"maximum number of RPC clients allowed for wallet RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultRPCMaxClients,
		),
		WalletRPCMaxWebsockets: NewInt(Metadata{
			Option: "walletrpcmaxwebsockets",
			Group:  "wallet",
			Label:  "Legacy RPC Max Websockets",
			Description:
			"maximum number of websocket clients allowed for wallet RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultRPCMaxWebsockets,
		),
		WalletServer: NewString(Metadata{
			Option:  "walletserver",
			Aliases: []string{"w"},
			Group:   "wallet",
			Label:   "Wallet Server",
			Description:
			"node address to connect wallet server to",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
		),
		Whitelists: NewStrings(Metadata{
			Option: "whitelists",
			Group:  "debug",
			Label:  "Whitelists",
			Description:
			"peers that you don't want to ever ban",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
	}
	return
}
