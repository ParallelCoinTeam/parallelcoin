package podcfgs

import (
	"math"
	"math/rand"
	"net"
	"path/filepath"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/p9c/matrjoska/pkg/appdata"
	"github.com/p9c/matrjoska/pkg/base58"
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/util/hdkeychain"
	"github.com/p9c/matrjoska/pod/config"
	"github.com/p9c/matrjoska/pod/podcmds"
	"github.com/p9c/matrjoska/pod/podconfig/checkpoints"
	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/duration"
	"github.com/p9c/opts/float"
	"github.com/p9c/opts/integer"
	"github.com/p9c/opts/list"
	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/sanitizers"
	"github.com/p9c/opts/text"
)

// GetDefaultConfig returns a Config struct pristine factory fresh
func GetDefaultConfig() (c *config.Config) {
	T.Ln("getting default config")
	c = &config.Config{
		Commands: podcmds.GetCommands(),
		Map:      GetConfigs(),
	}
	c.RunningCommand = c.Commands[0]

	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := range c.Map {
		tf := t.FieldByName(i)
		if tf.IsValid() && tf.CanSet() && tf.CanAddr() {
			val := reflect.ValueOf(c.Map[i])
			tf.Set(val)
		}
	}
	c.ForEach(func(ifc opt.Option) bool {

		return true
	})
	return
}

// GetConfigs returns configuration options for ParallelCoin Pod
func GetConfigs() (c config.Configs) {
	tags := func(s ...string) []string {
		return s
	}
	network := "mainnet"
	rand.Seed(time.Now().UnixNano())
	var datadir = &atomic.Value{}
	datadir.Store([]byte(appdata.Dir(constant.Name, false)))
	c = config.Configs{
		"AddCheckpoints": list.New(meta.Data{
			Aliases: []string{"AC"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Add Checkpoints",
			Description:
			"add custom checkpoints",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
			func(chkPts []string) (e error) {
				// todo: this closure should be added by the node and assign the output to its correct location
				var cpts []chaincfg.Checkpoint
				cpts, e = checkpoints.Parse(chkPts)
				_ = cpts
				return
			},
		),
		"AddPeers": list.New(meta.Data{
			Aliases: []string{"AP"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Add Peers",
			Description:
			"manually adds addresses to try to connect to",
			// Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"AddrIndex": binary.New(meta.Data{
			Aliases: []string{"AI"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Address Index",
			Description:
			"maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"AutoPorts": binary.New(meta.Data{
			Group: "debug",
			Label: "Automatic Ports",
			Tags:  tags("node", "wallet"),
			Description:
			"RPC and controller ports are randomized, use with controller for automatic peer discovery",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"AutoListen": binary.New(meta.Data{
			Aliases: []string{"AL"},
			Group:   "node",
			Tags:    tags("node", "wallet"),
			Label:   "Automatic Listeners",
			Description:
			"automatically update inbound addresses dynamically according to discovered network interfaces",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"BanDuration": duration.New(meta.Data{
			Aliases: []string{"BD"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Ban Opt",
			Description:
			"how long a ban of a misbehaving peer lasts",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			time.Hour*24,
			time.Second, time.Hour*24*365,
		),
		"BanThreshold": integer.New(meta.Data{
			Aliases: []string{"BT"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Ban Threshold",
			Description:
			"ban score that triggers a ban (default 100)",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultBanThreshold,
			1, 10000,
		),
		"BlockMaxSize": integer.New(meta.Data{
			Aliases: []string{"BMXS"},
			Group:   "mining",
			Tags:    tags("node"),
			Label:   "Block Max Size",
			Description:
			"maximum block size in bytes to be used when creating a block",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			blockchain.MaxBlockBaseSize-1000,
			constant.BlockMaxSizeMin, constant.BlockMaxSizeMax,
		),
		"BlockMaxWeight": integer.New(meta.Data{
			Aliases: []string{"BMXW"},
			Group:   "mining",
			Tags:    tags("node"),
			Label:   "Block Max Weight",
			Description:
			"maximum block weight to be used when creating a block",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.BlockMaxWeightMax,
			constant.BlockMaxWeightMin, constant.BlockMaxWeightMax,
		),
		"BlockMinSize": integer.New(meta.Data{
			Aliases: []string{"BMS"},
			Group:   "mining",
			Tags:    tags("node"),
			Label:   "Block Min Size",
			Description:
			"minimum block size in bytes to be used when creating a block",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.BlockMaxSizeMin,
			constant.BlockMaxSizeMin, constant.BlockMaxSizeMax,
		),
		"BlockMinWeight": integer.New(meta.Data{
			Aliases: []string{"BMW"},
			Group:   "mining",
			Tags:    tags("node"),
			Label:   "Block Min Weight",
			Description:
			"minimum block weight to be used when creating a block",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.BlockMaxWeightMin,
			constant.BlockMaxWeightMin, constant.BlockMaxWeightMax,
		),
		"BlockPrioritySize": integer.New(meta.Data{
			Aliases: []string{"BPS"},
			Group:   "mining",
			Tags:    tags("node"),
			Label:   "Block Priority Size",
			Description:
			"size in bytes for high-priority/low-fee transactions when creating a block",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultBlockPrioritySize,
			constant.BlockMaxSizeMin, constant.BlockMaxSizeMax,
		),
		"BlocksOnly": binary.New(meta.Data{
			Aliases: []string{"BO"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Blocks Only",
			Description:
			"do not accept transactions from remote peers",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"CAFile": text.New(meta.Data{
			Aliases: []string{"CA"},
			Group:   "tls",
			Tags:    tags("node", "wallet"),
			Label:   "Certificate Authority File",
			Description:
			"certificate authority file for TLS certificate validation",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "ca.cert"),
		),
		"ConfigFile": text.New(meta.Data{
			Aliases: []string{"CF"},
			Label:   "Configuration File",
			Description:
			"location of configuration file, cannot actually be changed",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			filepath.Join(string(datadir.Load().([]byte)), constant.PodConfigFilename),
		),
		"ConnectPeers": list.New(meta.Data{
			Aliases: []string{"CPS"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Connect Peers",
			Description:
			"connect ONLY to these addresses (disables inbound connections)",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
			DefaultPort:   constant.DefaultP2PPort,
		},
			[]string{},
		),
		"Controller": binary.New(meta.Data{
			Aliases: []string{"CN"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Enable Controller",
			Description:
			"delivers mining jobs over multicast",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"CPUProfile": text.New(meta.Data{
			Aliases: []string{"CPR"},
			Group:   "debug",
			Tags:    tags("node", "wallet", "kopach", "worker"),
			Label:   "CPU Profile",
			Description:
			"write cpu profile to this file",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"DarkTheme": binary.New(meta.Data{
			Aliases: []string{"DT"},
			Group:   "config",
			Tags:    tags("gui"),
			Label:   "Dark Theme",
			Description:
			"sets dark theme for GUI",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"DataDir": text.New(meta.Data{
			Aliases: []string{"DD"},
			Label:   "Data Directory",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Description:
			"root folder where application data is stored",
			Type:          sanitizers.Directory,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			appdata.Dir(constant.Name, false),
		),
		"DbType": text.New(meta.Data{
			Aliases: []string{"DB"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Database Type",
			Description:
			"type of database storage engine to use for node (" +
				"only one right now, ffldb)",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
			Options:       database.SupportedDrivers(),
		},
			constant.DefaultDbType,
		),
		"DisableBanning": binary.New(meta.Data{
			Aliases: []string{"NB"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Disable Banning",
			Description:
			"disables banning of misbehaving peers",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"DisableCheckpoints": binary.New(meta.Data{
			Aliases: []string{"NCP"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Disable Checkpoints",
			Description:
			"disables all checkpoints",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"DisableDNSSeed": binary.New(meta.Data{
			Aliases: []string{"NDS"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Disable DNS Seed",
			Description:
			"disable seeding of addresses to peers",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"DisableListen": binary.New(meta.Data{
			Aliases: []string{"NL"},
			Group:   "node",
			Tags:    tags("node", "wallet"),
			Label:   "Disable Listen",
			Description:
			"disables inbound connections for the peer to peer network",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"DisableRPC": binary.New(meta.Data{
			Aliases: []string{"NRPC"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "Disable RPC",
			Description:
			"disable rpc servers, as well as kopach controller",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"Discovery": binary.New(meta.Data{
			Aliases: []string{"DI"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Disovery",
			Description:
			"enable LAN peer discovery in GUI",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"ExternalIPs": list.New(meta.Data{
			Aliases: []string{"EI"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "External IP Addresses",
			Description:
			"extra addresses to tell peers they can connect to",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"FreeTxRelayLimit": float.New(meta.Data{
			Aliases: []string{"LR"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "Free Tx Relay Limit",
			Description:
			"limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultFreeTxRelayLimit,
			0, math.MaxFloat64,
		),
		"Generate": binary.New(meta.Data{
			Aliases: []string{"GB"},
			Group:   "mining",
			Tags:    tags("node", "kopach"),
			Label:   "Generate Blocks",
			Description:
			"turn on Kopach CPU miner",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"GenThreads": integer.New(meta.Data{
			Aliases: []string{"GT"},
			Group:   "mining",
			Tags:    tags("kopach"),
			Label:   "Generate Threads",
			Description:
			"number of threads to mine with",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			-1,
			-math.MaxInt64, runtime.NumCPU(),
		),
		"Hilite": list.New(meta.Data{
			Aliases: []string{"HL"},
			Group:   "debug",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Hilite",
			Description:
			"list of packages that will print with attention getters",
			Type:          "",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"LAN": binary.New(meta.Data{
			Group: "debug",
			Tags:  tags("node"),
			Label: "LAN Testnet Mode",
			Description:
			"run without any connection to nodes on the internet (does not apply on mainnet)",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"Locale": text.New(meta.Data{
			Aliases: []string{"LC"},
			Group:   "config",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Language",
			Description:
			"user interface language i18 localization",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
			Options:       []string{"en"},
		},
			"en",
		),
		"LimitPass": text.New(meta.Data{
			Aliases: []string{"LP"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "Limit Password",
			Description:
			"limited user password",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			genPassword(),
		),
		"LimitUser": text.New(meta.Data{
			Aliases: []string{"LU"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "Limit Username",
			Description:
			"limited user name",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			"limit",
		),
		"LogDir": text.New(meta.Data{
			Aliases: []string{"LD"},
			Group:   "config",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Log Directory",
			Description:
			"folder where log files are written",
			Type:          sanitizers.Directory,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			string(datadir.Load().([]byte)),
		),
		"LogFilter": list.New(meta.Data{
			Aliases: []string{"LF"},
			Group:   "debug",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Log Filter",
			Description:
			"list of packages that will not print logs",
			Type:          "",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"LogLevel": text.New(meta.Data{
			Aliases: []string{"LL"},
			Group:   "config",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Log Level",
			Description:
			"maximum log level to output",
			Options: []string{
				"off",
				"fatal",
				"error",
				"info",
				"check",
				"debug",
				"trace",
			},
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"info",

		),
		"MaxOrphanTxs": integer.New(meta.Data{
			Aliases: []string{"MO"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "Max Orphan Txs",
			Description:
			"max number of orphan transactions to keep in memory",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMaxOrphanTransactions,
			0, math.MaxInt64,
		),
		"MaxPeers": integer.New(meta.Data{
			Aliases: []string{"MP"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Max Peers",
			Description:
			"maximum number of peers to hold connections with",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMaxPeers,
			1, 256,
		),
		"MulticastPass": text.New(meta.Data{
			Aliases: []string{"PM"},
			Group:   "config",
			Tags:    tags("node", "kopach"),
			Label:   "Multicast Pass",
			Description:
			"password that encrypts the connection to the mining controller",
			Type:          sanitizers.Password,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"pa55word",
		),
		"MinRelayTxFee": float.New(meta.Data{
			Aliases: []string{"MRTF"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "Min Relay Transaction Fee",
			Description:
			"the minimum transaction fee in DUO/kB to be considered a non-zero fee",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMinRelayTxFee.ToDUO(),
			0, math.MaxFloat64,

		),
		"Network": text.New(meta.Data{
			Aliases: []string{"NW"},
			Group:   "node",
			Tags:    tags("node", "wallet"),
			Label:   "Network",
			Description:
			"connect to this network:",
			Options: []string{
				"mainnet",
				"testnet",
				"regtestnet",
				"simnet",
			},
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			network,
		),
		"NoCFilters": binary.New(meta.Data{
			Aliases: []string{"NCF"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "No CFilters",
			Description:
			"disable committed filtering (CF) support",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"NodeOff": binary.New(meta.Data{
			Aliases: []string{"NO"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Node Off",
			Description:
			"turn off the node backend",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"NoInitialLoad": binary.New(meta.Data{
			Aliases: []string{"NIL"},
			Group:   "wallet",
			Tags:    tags("wallet"),
			Label:   "No Initial Load",
			Description:
			"do not load a wallet at startup",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"NoPeerBloomFilters": binary.New(meta.Data{
			Aliases: []string{"NPBF"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "No Peer Bloom Filters",
			Description:
			"disable bloom filtering support",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"NoRelayPriority": binary.New(meta.Data{
			Aliases: []string{"NRPR"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "No Relay Priority",
			Description:
			"do not require free or low-fee transactions to have high priority for relaying",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"OneTimeTLSKey": binary.New(meta.Data{
			Aliases: []string{"OTK"},
			Group:   "wallet",
			Tags:    tags("node", "wallet"),
			Label:   "One Time TLS Key",
			Description:
			"generate a new TLS certificate pair at startup, but only write the certificate to disk",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"OnionEnabled": binary.New(meta.Data{
			Aliases: []string{"OE"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Onion Enabled",
			Description:
			"enable tor proxy",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"OnionProxyAddress": text.New(meta.Data{
			Aliases: []string{"OPA"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Onion Proxy Address",
			Description:
			"address of tor proxy you want to connect to",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"OnionProxyPass": text.New(meta.Data{
			Aliases: []string{"OPW"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Onion Proxy Password",
			Description:
			"password for tor proxy",
			Type:          sanitizers.Password,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"OnionProxyUser": text.New(meta.Data{
			Aliases: []string{"OU"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Onion Proxy Username",
			Description:
			"tor proxy username",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"P2PConnect": list.New(meta.Data{
			Aliases: []string{"P2P"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "P2P Connect",
			Description:
			"list of addresses reachable from connected networks",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"P2PListeners": list.New(meta.Data{
			Aliases: []string{"LA"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "P2PListeners",
			Description:
			"list of addresses to bind the node listener to",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{
				net.JoinHostPort("0.0.0.0",
					chaincfg.MainNetParams.DefaultPort,
				),
			},
		),
		"Password": text.New(meta.Data{
			Aliases: []string{"PW"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "Password",
			Description:
			"password for client RPC connections",
			Type:          sanitizers.Password,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			genPassword(),
		),
		"PipeLog": binary.New(meta.Data{
			Aliases: []string{"PL"},
			Label:   "Pipe Logger",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Description:
			"enable pipe based logger IPC",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"Profile": text.New(meta.Data{
			Aliases: []string{"HPR"},
			Group:   "debug",
			Tags:    tags("node", "wallet", "ctl", "kopach", "worker"),
			Label:   "Profile",
			Description:
			"http profiling on given port (1024-40000)",
			// Type:        "",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"ProxyAddress": text.New(meta.Data{
			Aliases: []string{"PA"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Proxy",
			Description:
			"address of proxy to connect to for outbound connections",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			"",
		),
		"ProxyPass": text.New(meta.Data{
			Aliases: []string{"PPW"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Proxy Pass",
			Description:
			"proxy password, if required",
			Type:          sanitizers.Password,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			genPassword(),
		),
		"ProxyUser": text.New(meta.Data{
			Aliases: []string{"PU"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "ProxyUser",
			Description:
			"proxy username, if required",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			"proxyuser",
		),
		"RejectNonStd": binary.New(meta.Data{
			Aliases: []string{"REJ"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Reject Non Std",
			Description:
			"reject non-standard transactions regardless of the default settings for the active network",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"RelayNonStd": binary.New(meta.Data{
			Aliases: []string{"RNS"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Relay Nonstandard Transactions",
			Description:
			"relay non-standard transactions regardless of the default settings for the active network",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"RPCCert": text.New(meta.Data{
			Aliases: []string{"RC"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "RPC Cert",
			Description:
			"location of RPC TLS certificate",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.cert"),
		),
		"RPCConnect": text.New(meta.Data{
			Aliases: []string{"RA"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "RPC Connect",
			Description:
			"full node RPC for wallet",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.RPCClientPort),
		),
		"RPCKey": text.New(meta.Data{
			Aliases: []string{"RK"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "RPC Key",
			Description:
			"location of rpc TLS key",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.key"),
		),
		"RPCListeners": list.New(meta.Data{
			Aliases: []string{"RL"},
			Group:   "rpc",
			Tags:    tags("node"),
			Label:   "Node RPC Listeners",
			Description:
			"addresses to listen for RPC connections",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{
				net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.RPCClientPort),
			},
		),
		"RPCMaxClients": integer.New(meta.Data{
			Aliases: []string{"RMXC"},
			Group:   "rpc",
			Tags:    tags("node"),
			Label:   "Maximum Node RPC Clients",
			Description:
			"maximum number of clients for regular RPC",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMaxRPCClients,
			0, 256,
		),
		"RPCMaxConcurrentReqs": integer.New(meta.Data{
			Aliases: []string{"RMCR"},
			Group:   "rpc",
			Tags:    tags("node"),
			Label:   "Maximum Node RPC Concurrent Reqs",
			Description:
			"maximum number of requests to process concurrently",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMaxRPCConcurrentReqs,
			0, 4096,
		),
		"RPCMaxWebsockets": integer.New(meta.Data{
			Aliases: []string{"RMWS"},
			Group:   "rpc",
			Tags:    tags("node"),
			Label:   "Maximum Node RPC Websockets",
			Description:
			"maximum number of websocket clients to allow",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultMaxRPCWebsockets,
			0, 4096,
		),
		"RPCQuirks": binary.New(meta.Data{
			Aliases: []string{"RQ"},
			Group:   "rpc",
			Tags:    tags("node"),
			Label:   "Emulate Bitcoin Core RPC Quirks",
			Description:
			"enable bugs that replicate bitcoin core RPC's JSON",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"RunAsService": binary.New(meta.Data{
			Aliases: []string{"RS"},
			Label:   "Run As Service",
			Description:
			"shuts down on lock timeout",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"Save": binary.New(meta.Data{
			Aliases: []string{"SV"},
			Label:   "Save Configuration",
			Description:
			"save opts given on commandline",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"ServerTLS": binary.New(meta.Data{
			Aliases: []string{"ST"},
			Group:   "wallet",
			Tags:    tags("node", "wallet"),
			Label:   "Server TLS",
			Description:
			"enable TLS for the wallet connection to node RPC server",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			true,
		),
		"SigCacheMaxSize": integer.New(meta.Data{
			Aliases: []string{"SCM"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Signature Cache Max Size",
			Description:
			"the maximum number of entries in the signature verification cache",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultSigCacheMaxSize,
			constant.DefaultSigCacheMaxSize, constant.BlockMaxSizeMax,
		),
		"Solo": binary.New(meta.Data{
			Group: "mining",
			Label: "Solo Generate",
			Tags:    tags("node"),
			Description:
			"mine even if not connected to a network",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"ClientTLS": binary.New(meta.Data{
			Aliases: []string{"CT"},
			Group:   "tls",
			Tags:    tags("node", "wallet"),
			Label:   "TLS",
			Description:
			"enable TLS for RPC client connections",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			true,
		),
		"TLSSkipVerify": binary.New(meta.Data{
			Aliases: []string{"TSV"},
			Group:   "tls",
			Tags:    tags("node", "wallet"),
			Label:   "TLS Skip Verify",
			Description:
			"skip TLS certificate verification (ignore CA errors)",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"TorIsolation": binary.New(meta.Data{
			Aliases: []string{"TI"},
			Group:   "proxy",
			Tags:    tags("node"),
			Label:   "Tor Isolation",
			Description:
			"makes a separate proxy connection for each connection",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"TrickleInterval": duration.New(meta.Data{
			Aliases: []string{"TKI"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "Trickle Interval",
			Description:
			"minimum time between attempts to send new inventory to a connected peer",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultTrickleInterval,
			time.Second, time.Second*30,
		),
		"TxIndex": binary.New(meta.Data{
			Aliases: []string{"TXI"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "Tx Index",
			Description:
			"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"UPNP": binary.New(meta.Data{
			Aliases: []string{"UP"},
			Group:   "node",
			Tags:    tags("node"),
			Label:   "UPNP",
			Description:
			"enable UPNP for NAT traversal",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"UserAgentComments": list.New(meta.Data{
			Aliases: []string{"UA"},
			Group:   "policy",
			Tags:    tags("node"),
			Label:   "User Agent Comments",
			Description:
			"comment to add to the user agent -- See BIP 14 for more information",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
		"Username": text.New(meta.Data{
			Aliases: []string{"UN"},
			Group:   "rpc",
			Tags:    tags("node", "wallet"),
			Label:   "Username",
			Description:
			"password for client RPC connections",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			"username",
		),
		"UUID": integer.New(meta.Data{
			Label: "UUID",
			Description:
			"instance unique id (32bit random value) (json mangles big 64 bit integers due to float64 numbers)",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			int64(rand.Uint32()),
			-math.MaxInt64, math.MaxInt64,
		),
		"UseWallet": binary.New(meta.Data{
			Aliases: []string{"WC"},
			Group:   "debug",
			Tags:    tags("ctl"),
			Label:   "Connect to Wallet",
			Description:
			"set ctl to connect to wallet instead of chain server",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"WalletFile": text.New(meta.Data{
			Aliases: []string{"WF"},
			Group:   "config",
			Tags:    tags("wallet"),
			Label:   "Wallet File",
			Description:
			"wallet database file",
			Type:          sanitizers.FilePath,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "mainnet", constant.DbName),
		),
		"WalletOff": binary.New(meta.Data{
			Aliases: []string{"WO"},
			Group:   "debug",
			Tags:    tags("wallet"),
			Label:   "Wallet Off",
			Description:
			"turn off the wallet backend",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			false,
		),
		"WalletPass": text.New(meta.Data{
			Aliases: []string{"WPW"},
			Label:   "Wallet Pass",
			Tags:    tags("wallet"),
			Description:
			"password encrypting public data in wallet - only hash is stored" +
				" so give on command line or in environment POD_WALLETPASS",
			Type:          sanitizers.Password,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     false,
		},
			"",
		),
		"WalletRPCListeners": list.New(meta.Data{
			Aliases: []string{"WRL"},
			Group:   "wallet",
			Tags:    tags("wallet"),
			Label:   "Wallet RPC Listeners",
			Description:
			"addresses for wallet RPC server to listen on",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{
				net.JoinHostPort("0.0.0.0",
					chaincfg.MainNetParams.WalletRPCServerPort,
				),
			},
		),
		"WalletRPCMaxClients": integer.New(meta.Data{
			Aliases: []string{"WRMC"},
			Group:   "wallet",
			Tags:    tags("wallet"),
			Label:   "Legacy RPC Max Clients",
			Description:
			"maximum number of RPC clients allowed for wallet RPC",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultRPCMaxClients,
			0, 4096,
		),
		"WalletRPCMaxWebsockets": integer.New(meta.Data{
			Aliases: []string{"WRMWS"},
			Group:   "wallet",
			Tags:    tags("wallet"),
			Label:   "Legacy RPC Max Websockets",
			Description:
			"maximum number of websocket clients allowed for wallet RPC",
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			constant.DefaultRPCMaxWebsockets,
			0, 4096,
		),
		"WalletServer": text.New(meta.Data{
			Aliases: []string{"WS"},
			Group:   "wallet",
			Tags:    tags("wallet"),
			Label:   "Wallet Server",
			Description:
			"node address to connect wallet server to",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
		),
		"Whitelists": list.New(meta.Data{
			Aliases: []string{"WL"},
			Group:   "debug",
			Tags:    tags("node"),
			Label:   "Whitelists",
			Description:
			"peers that you don't want to ever ban",
			Type:          sanitizers.NetAddress,
			Documentation: "<placeholder for detailed documentation>",
			OmitEmpty:     true,
		},
			[]string{},
		),
	}
	for i := range c {
		c[i].SetName(i)
	}
	return
}

func genPassword() string {
	s, e := hdkeychain.GenerateSeed(16)
	if e != nil {
		panic("can't do nothing without entropy! " + e.Error())
	}
	return base58.Encode(s)
}
