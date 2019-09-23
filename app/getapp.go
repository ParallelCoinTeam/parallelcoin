package app

import (
	"fmt"
	"time"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"

	"github.com/p9c/pod/app/util"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/base58"
	"github.com/p9c/pod/pkg/util/hdkeychain"
)

// getApp defines the pod app
func getApp(cx *conte.Xt) (a *cli.App) {
	return &cli.App{
		Name:    "pod",
		Version: "v0.0.1",
		Description: "Parallelcoin Pod Suite -- All-in-one everything" +
			" for Parallelcoin!",
		Copyright: "Legacy portions derived from btcsuite/btcd under" +
			" ISC licence. The remainder is already in your" +
			" possession. Use it wisely.",
		Action: func(c *cli.Context) error {
			fmt.Println("no subcommand requested")
			cli.ShowAppHelpAndExit(c, 1)
			return nil
		},
		Before: beforeFunc(cx),
		After: func(c *cli.Context) error {
			log.TRACE("subcommand completed")
			return nil
		},
		Commands: []cli.Command{
			util.NewCommand("version",
				"print version and exit",
				func(c *cli.Context) error {
					fmt.Println(c.App.Name, c.App.Version)
					return nil
				},
				util.SubCommands(),
				"v"),
			util.NewCommand("ctl",
				"send RPC commands to a node or wallet and print the result",
				ctlHandle(cx),
				util.SubCommands(
					util.NewCommand(
						"listcommands",
						"list commands available at endpoint",
						ctlHandleList,
						nil,
						"list", "l",
					),
				),
				"c"),
			util.NewCommand("node",
				"start parallelcoin full node",
				nodeHandle(cx),
				util.SubCommands(
					util.NewCommand("dropaddrindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropAddrIndex = true
							return nodeHandle(cx)(c)
						},
						util.SubCommands(),
					),
					util.NewCommand("droptxindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropTxIndex = true
							return nodeHandle(cx)(c)
						},
						util.SubCommands(),
					),
					util.NewCommand("dropcfindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
						},
						util.SubCommands(),
					),
				),
				"n",
			),
			util.NewCommand("wallet",
				"start parallelcoin wallet server",
				walletHandle(cx),
				util.SubCommands(),
				"w",
			),
			util.NewCommand("shell",
				"start combined wallet/node shell",
				shellHandle(cx),
				util.SubCommands(),
				"s",
			),
			util.NewCommand(
				"gui",
				"start GUI",
				guiHandle(cx),
				util.SubCommands(),
			),
			util.NewCommand("kopach",
				"standalone miner for clusters",
				kopachHandle(cx),
				util.SubCommands(
					// apputil.NewCommand("bench",
					// 	"generate a set of benchmarks of each algorithm",
					// 	func(c *cli.Context) error {
					// 		return bench.Benchmark(cx)(c)
					// 	},
					// 	apputil.SubCommands(),
					// ),
				),
			),
		},
		Flags: []cli.Flag{
			altsrc.NewStringFlag(cli.StringFlag{
				Name:        "datadir, D",
				Value:       *cx.Config.DataDir,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_DATADIR",
				Destination: cx.Config.DataDir,
			}),
			util.BoolTrue("save, i",
				"save settings as effective from invocation",
				&cx.StateCfg.Save,
			),
			altsrc.NewStringFlag(cli.StringFlag{
				Name:        "loglevel, l",
				Value:       *cx.Config.LogLevel,
				Usage:       "sets the base for all subsystem logging",
				EnvVar:      "POD_LOGLEVEL",
				Destination: cx.Config.LogLevel,
			}),
			util.String(
				"network, n",
				"connect to mainnet/testnet/regtest/simnet",
				"mainnet",
				cx.Config.Network),
			util.String(
				"username",
				"sets the username for services",
				"server",
				cx.Config.Username),
			util.String(
				"password",
				"sets the password for services",
				genPassword(),
				cx.Config.Password),
			util.String(
				"serveruser",
				"sets the username for clients of services",
				"client",
				cx.Config.ServerUser),
			util.String(
				"serverpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.ServerPass),
			util.String(
				"limituser",
				"sets the limited rpc username",
				"limit",
				cx.Config.LimitUser),
			util.String(
				"limitpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.LimitPass),
			util.String(
				"rpccert",
				"File containing the certificate file",
				util.Join(*cx.Config.DataDir, "rpc.cert"),
				cx.Config.RPCCert),
			util.String(
				"rpckey",
				"File containing the certificate key",
				util.Join(*cx.Config.DataDir, "rpc.key"),
				cx.Config.RPCKey),
			util.String(
				"cafile",
				"File containing root certificates to authenticate a TLS"+
					" connections with pod",
				util.Join(*cx.Config.DataDir, "cafile"),
				cx.Config.CAFile),
			util.Bool(
				"clienttls",
				"Enable TLS for client connections",
				cx.Config.TLS),
			util.Bool(
				"servertls",
				"Enable TLS for server connections",
				cx.Config.ServerTLS),
			util.String(
				"proxy",
				"Connect via SOCKS5 proxy",
				"",
				cx.Config.Proxy),
			util.String(
				"proxyuser",
				"Username for proxy server",
				"user",
				cx.Config.ProxyUser),
			util.String(
				"proxypass",
				"Password for proxy server",
				"pa55word",
				cx.Config.ProxyPass),
			util.Bool(
				"onion",
				"Enable connecting to tor hidden services",
				cx.Config.Onion),
			util.String(
				"onionproxy",
				"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0."+
					"0.1:9050)",
				"127.0.0.1:9050",
				cx.Config.OnionProxy),
			util.String(
				"onionuser",
				"Username for onion proxy server",
				"user",
				cx.Config.OnionProxyUser),
			util.String(
				"onionpass",
				"Password for onion proxy server",
				genPassword(),
				cx.Config.OnionProxyPass),
			util.Bool(
				"torisolation",
				"Enable Tor stream isolation by randomizing user credentials"+
					" for each connection.",
				cx.Config.TorIsolation),
			util.String(
				"group",
				"zeroconf testnet group identifier (whitelist connections)",
				"",
				cx.Config.Group),
			util.Bool(
				"nodiscovery",
				"disable zeroconf peer discovery",
				cx.Config.NoDiscovery),
			util.StringSlice(
				"addpeer",
				"Add a peer to connect with at startup",
				cx.Config.AddPeers),
			util.StringSlice(
				"connect",
				"Connect only to the specified peers at startup",
				cx.Config.ConnectPeers),
			util.Bool(
				"nolisten",
				"Disable listening for incoming connections -- NOTE:"+
					" Listening is automatically disabled if the --connect or"+
					" --proxy options are used without also specifying listen"+
					" interfaces via --listen",
				cx.Config.DisableListen),
			util.StringSlice(
				"listen",
				"Add an interface/port to listen for connections",
				cx.Config.Listeners),
			util.Int(
				"maxpeers",
				"Max number of inbound and outbound peers",
				node.DefaultMaxPeers,
				cx.Config.MaxPeers),
			util.Bool(
				"nobanning",
				"Disable banning of misbehaving peers",
				cx.Config.DisableBanning),
			util.Duration(
				"banduration",
				"How long to ban misbehaving peers",
				time.Hour*24,
				cx.Config.BanDuration),
			util.Int(
				"banthreshold",
				"Maximum allowed ban score before disconnecting and"+
					" banning misbehaving peers.",
				node.DefaultBanThreshold,
				cx.Config.BanThreshold),
			util.StringSlice(
				"whitelist",
				"Add an IP network or IP that will not be banned. (eg. 192."+
					"168.1.0/24 or ::1)",
				cx.Config.Whitelists),
			util.String(
				"rpcconnect",
				"Hostname/IP and port of pod RPC server to connect to",
				"127.0.0.1:11048",
				cx.Config.RPCConnect),
			util.StringSlice(
				"rpclisten",
				"Add an interface/port to listen for RPC connections",
				cx.Config.RPCListeners),
			util.Int(
				"rpcmaxclients",
				"Max number of RPC clients for standard connections",
				node.DefaultMaxRPCClients,
				cx.Config.RPCMaxClients),
			util.Int(
				"rpcmaxwebsockets",
				"Max number of RPC websocket connections",
				node.DefaultMaxRPCWebsockets,
				cx.Config.RPCMaxWebsockets),
			util.Int(
				"rpcmaxconcurrentreqs",
				"Max number of RPC requests that may be"+
					" processed concurrently",
				node.DefaultMaxRPCConcurrentReqs,
				cx.Config.RPCMaxConcurrentReqs),
			util.Bool(
				"rpcquirks",
				"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE:"+
					" Discouraged unless interoperability issues need to be worked"+
					" around",
				cx.Config.RPCQuirks),
			util.Bool(
				"norpc",
				"Disable built-in RPC server -- NOTE: The RPC server"+
					" is disabled by default if no rpcuser/rpcpass or"+
					" rpclimituser/rpclimitpass is specified",
				cx.Config.DisableRPC),
			util.Bool(
				"nodnsseed",
				"Disable DNS seeding for peers",
				cx.Config.DisableDNSSeed),
			util.StringSlice(
				"externalip",
				"Add an ip to the list of local addresses we claim to"+
					" listen on to peers",
				cx.Config.ExternalIPs),
			util.StringSlice(
				"addcheckpoint",
				"Add a custom checkpoint.  Format: '<height>:<hash>'",
				cx.Config.AddCheckpoints),
			util.Bool(
				"nocheckpoints",
				"Disable built-in checkpoints.  Don't do this unless"+
					" you know what you're doing.",
				cx.Config.DisableCheckpoints),
			util.String(
				"dbtype",
				"Database backend to use for the Block Chain",
				node.DefaultDbType,
				cx.Config.DbType),
			util.String(
				"profile",
				"Enable HTTP profiling on given port -- NOTE port"+
					" must be between 1024 and 65536",
				"",
				cx.Config.Profile),
			util.String(
				"cpuprofile",
				"Write CPU profile to the specified file",
				"",
				cx.Config.CPUProfile),
			util.Bool(
				"upnp",
				"Use UPnP to map our listening port outside of NAT",
				cx.Config.Upnp),
			util.Float64(
				"minrelaytxfee",
				"The minimum transaction fee in DUO/kB to be"+
					" considered a non-zero fee.",
				mempool.DefaultMinRelayTxFee.ToDUO(),
				cx.Config.MinRelayTxFee),
			util.Float64(
				"limitfreerelay",
				"Limit relay of transactions with no transaction"+
					" fee to the given amount in thousands of bytes per minute",
				node.DefaultFreeTxRelayLimit,
				cx.Config.FreeTxRelayLimit),
			util.Bool(
				"norelaypriority",
				"Do not require free or low-fee transactions to have"+
					" high priority for relaying",
				cx.Config.NoRelayPriority),
			util.Duration(
				"trickleinterval",
				"Minimum time between attempts to send new"+
					" inventory to a connected peer",
				node.DefaultTrickleInterval,
				cx.Config.TrickleInterval),
			util.Int(
				"maxorphantx",
				"Max number of orphan transactions to keep in memory",
				node.DefaultMaxOrphanTransactions,
				cx.Config.MaxOrphanTxs),
			util.String(
				"algo",
				"Sets the algorithm for the CPU miner ( blake14lr,"+
					" cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog,"+
					" skein, x11 default is 'random')",
				"random",
				cx.Config.Algo),
			util.Bool(
				"generate",
				"Generate (mine) DUO using the CPU",
				cx.Config.Generate),
			util.Int(
				"genthreads",
				"Number of CPU threads to use with CPU miner"+
					" -1 = all cores",
				-1,
				cx.Config.GenThreads),
			util.String(
				"controller",
				"address to bind miner controller listener",
				genPassword(),
				cx.Config.Controller),
			util.Bool(
				"nocontroller",
				"disable zeroconf kcp miner controller",
				cx.Config.NoController),
			util.StringSlice(
				"miningaddr",
				"Add the specified payment address to the list of"+
					" addresses to use for generated blocks, at least one is "+
					"required if generate or minerlistener are set",
				cx.Config.MiningAddrs),
			util.String(
				"minerpass",
				"password to authorise sending work to a miner",
				genPassword(),
				cx.Config.MinerPass),
			util.Int(
				"blockminsize",
				"Minimum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMin,
				cx.Config.BlockMinSize),
			util.Int(
				"blockmaxsize",
				"Maximum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMax,
				cx.Config.BlockMaxSize),
			util.Int(
				"blockminweight",
				"Minimum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMin,
				cx.Config.BlockMinWeight),
			util.Int(
				"blockmaxweight",
				"Maximum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMax,
				cx.Config.BlockMaxWeight),
			util.Int(
				"blockprioritysize",
				"Size in bytes for high-priority/low-fee"+
					" transactions when creating a block",
				mempool.DefaultBlockPrioritySize,
				cx.Config.BlockPrioritySize),
			util.StringSlice(
				"uacomment",
				"Comment to add to the user agent -- See BIP 14 for"+
					" more information.",
				cx.Config.UserAgentComments),
			util.Bool(
				"nopeerbloomfilters",
				"Disable bloom filtering support",
				cx.Config.NoPeerBloomFilters),
			util.Bool(
				"nocfilters",
				"Disable committed filtering (CF) support",
				cx.Config.NoCFilters),
			util.Int(
				"sigcachemaxsize",
				"The maximum number of entries in the"+
					" signature verification cache",
				node.DefaultSigCacheMaxSize,
				cx.Config.SigCacheMaxSize),
			util.Bool(
				"blocksonly",
				"Do not accept transactions from remote peers.",
				cx.Config.BlocksOnly),
			util.BoolTrue(
				"notxindex",
				"Disable the transaction index which makes all transactions"+
					" available via the getrawtransaction RPC",
				cx.Config.TxIndex),
			util.BoolTrue(
				"noaddrindex",
				"Disable address-based transaction index which"+
					" makes the searchrawtransactions RPC available",
				cx.Config.AddrIndex,
			),
			util.Bool(
				"relaynonstd",
				"Relay non-standard transactions regardless of the default"+
					" settings for the active network.",
				cx.Config.RelayNonStd), util.Bool("rejectnonstd",
				"Reject non-standard transactions regardless of"+
					" the default settings for the active network.",
				cx.Config.RejectNonStd),
			util.Bool(
				"noinitialload",
				"Defer wallet creation/opening on startup and"+
					" enable loading wallets over RPC",
				cx.Config.NoInitialLoad),
			util.Bool(
				"walletconnect, wc",
				"connect to wallet instead of full node",
				cx.Config.Wallet),
			util.String(
				"walletserver, ws",
				"set wallet server to connect to",
				"127.0.0.1:11046",
				cx.Config.WalletServer),
			util.String(
				"walletpass",
				"The public wallet password -- Only required if"+
					" the wallet was created with one",
				"",
				cx.Config.WalletPass),
			util.Bool(
				"onetimetlskey",
				"Generate a new TLS certpair at startup, but"+
					" only write the certificate to disk",
				cx.Config.OneTimeTLSKey),
			util.Bool(
				"tlsskipverify",
				"skip verifying tls certificates",
				cx.Config.TLSSkipVerify),
			util.StringSlice(
				"walletrpclisten",
				"Listen for wallet RPC connections on this"+
					" interface/port (default port: 11046, testnet: 21046,"+
					" simnet: 41046)",
				cx.Config.WalletRPCListeners),
			util.Int(
				"walletrpcmaxclients",
				"Max number of legacy RPC clients for"+
					" standard connections",
				8,
				cx.Config.WalletRPCMaxClients),
			util.Int(
				"walletrpcmaxwebsockets",
				"Max number of legacy RPC websocket connections",
				8,
				cx.Config.WalletRPCMaxWebsockets,
			),
			util.StringSlice(
				"experimentalrpclisten",
				"Listen for RPC connections on this interface/port",
				cx.Config.ExperimentalRPCListeners),
			util.Bool(
				"nodeoff",
				"Starts GUI with node turned off",
				cx.Config.NodeOff),
			util.Bool( // TODO remove this
				"testnodeoff",
				"Starts GUI with testnode turned off",
				cx.Config.TestNodeOff),
			util.Bool(
				"walletoff",
				"Starts GUI with wallet turned off",
				cx.Config.WalletOff,
			),
		},
	}
}

func genPassword() string {
	s, err := hdkeychain.GenerateSeed(16)
	if err != nil {
		panic("can't do nothing without entropy! " + err.Error())
	}
	return base58.Encode(s)
}
