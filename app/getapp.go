package app

import (
	"fmt"
	"time"

	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/base58"
	"github.com/p9c/pod/pkg/util/hdkeychain"
)

func // getApp defines the pod app
getApp(cx *conte.Xt) (a *cli.App) {
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
			apputil.NewCommand("version",
				"print version and exit",
				func(c *cli.Context) error {
					fmt.Println(c.App.Name, c.App.Version)
					return nil
				},
				apputil.SubCommands(),
				"v"),
			apputil.NewCommand("ctl",
				"send RPC commands to a node or wallet and print the result",
				ctlHandle(cx),
				apputil.SubCommands(
					apputil.NewCommand(
						"listcommands",
						"list commands available at endpoint",
						ctlHandleList,
						nil,
						"list", "l",
					),
				),
				"c"),
			apputil.NewCommand("node",
				"start parallelcoin full node",
				nodeHandle(cx),
				apputil.SubCommands(
					apputil.NewCommand("dropaddrindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropAddrIndex = true
							return nodeHandle(cx)(c)
						},
						apputil.SubCommands(),
					),
					apputil.NewCommand("droptxindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropTxIndex = true
							return nodeHandle(cx)(c)
						},
						apputil.SubCommands(),
					),
					apputil.NewCommand("dropcfindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
						},
						apputil.SubCommands(),
					),
				),
				"n",
			),
			apputil.NewCommand("wallet",
				"start parallelcoin wallet server",
				walletHandle(cx),
				apputil.SubCommands(),
				"w",
			),
			apputil.NewCommand("shell",
				"start combined wallet/node shell",
				shellHandle(cx),
				apputil.SubCommands(),
				"s",
			),
			apputil.NewCommand(
				"gui",
				"start GUI",
				guiHandle(cx),
				apputil.SubCommands(),
			),
			apputil.NewCommand("kopach",
				"standalone miner for clusters",
				kopachHandle(cx),
				apputil.SubCommands(
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
			apputil.BoolTrue("save, i",
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
			apputil.String(
				"network, n",
				"connect to mainnet/testnet/regtest/simnet",
				"mainnet",
				cx.Config.Network),
			apputil.String(
				"username",
				"sets the username for services",
				"server",
				cx.Config.Username),
			apputil.String(
				"password",
				"sets the password for services",
				genPassword(),
				cx.Config.Password),
			apputil.String(
				"serveruser",
				"sets the username for clients of services",
				"client",
				cx.Config.ServerUser),
			apputil.String(
				"serverpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.ServerPass),
			apputil.String(
				"limituser",
				"sets the limited rpc username",
				"limit",
				cx.Config.LimitUser),
			apputil.String(
				"limitpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.LimitPass),
			apputil.String(
				"rpccert",
				"File containing the certificate file",
				apputil.Join(*cx.Config.DataDir, "rpc.cert"),
				cx.Config.RPCCert),
			apputil.String(
				"rpckey",
				"File containing the certificate key",
				apputil.Join(*cx.Config.DataDir, "rpc.key"),
				cx.Config.RPCKey),
			apputil.String(
				"cafile",
				"File containing root certificates to authenticate a TLS"+
					" connections with pod",
				apputil.Join(*cx.Config.DataDir, "cafile"),
				cx.Config.CAFile),
			apputil.Bool(
				"clienttls",
				"Enable TLS for client connections",
				cx.Config.TLS),
			apputil.Bool(
				"servertls",
				"Enable TLS for server connections",
				cx.Config.ServerTLS),
			apputil.String(
				"proxy",
				"Connect via SOCKS5 proxy",
				"",
				cx.Config.Proxy),
			apputil.String(
				"proxyuser",
				"Username for proxy server",
				"user",
				cx.Config.ProxyUser),
			apputil.String(
				"proxypass",
				"Password for proxy server",
				"pa55word",
				cx.Config.ProxyPass),
			apputil.Bool(
				"onion",
				"Enable connecting to tor hidden services",
				cx.Config.Onion),
			apputil.String(
				"onionproxy",
				"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0."+
					"0.1:9050)",
				"127.0.0.1:9050",
				cx.Config.OnionProxy),
			apputil.String(
				"onionuser",
				"Username for onion proxy server",
				"user",
				cx.Config.OnionProxyUser),
			apputil.String(
				"onionpass",
				"Password for onion proxy server",
				genPassword(),
				cx.Config.OnionProxyPass),
			apputil.Bool(
				"torisolation",
				"Enable Tor stream isolation by randomizing user credentials"+
					" for each connection.",
				cx.Config.TorIsolation),
			apputil.StringSlice(
				"addpeer",
				"Add a peer to connect with at startup",
				cx.Config.AddPeers),
			apputil.StringSlice(
				"connect",
				"Connect only to the specified peers at startup",
				cx.Config.ConnectPeers),
			apputil.Bool(
				"nolisten",
				"Disable listening for incoming connections -- NOTE:"+
					" Listening is automatically disabled if the --connect or"+
					" --proxy options are used without also specifying listen"+
					" interfaces via --listen",
				cx.Config.DisableListen),
			apputil.StringSlice(
				"listen",
				"Add an interface/port to listen for connections",
				cx.Config.Listeners),
			apputil.Int(
				"maxpeers",
				"Max number of inbound and outbound peers",
				node.DefaultMaxPeers,
				cx.Config.MaxPeers),
			apputil.Bool(
				"nobanning",
				"Disable banning of misbehaving peers",
				cx.Config.DisableBanning),
			apputil.Duration(
				"banduration",
				"How long to ban misbehaving peers",
				time.Hour*24,
				cx.Config.BanDuration),
			apputil.Int(
				"banthreshold",
				"Maximum allowed ban score before disconnecting and"+
					" banning misbehaving peers.",
				node.DefaultBanThreshold,
				cx.Config.BanThreshold),
			apputil.StringSlice(
				"whitelist",
				"Add an IP network or IP that will not be banned. (eg. 192."+
					"168.1.0/24 or ::1)",
				cx.Config.Whitelists),
			apputil.String(
				"rpcconnect",
				"Hostname/IP and port of pod RPC server to connect to",
				"127.0.0.1:11048",
				cx.Config.RPCConnect),
			apputil.StringSlice(
				"rpclisten",
				"Add an interface/port to listen for RPC connections",
				cx.Config.RPCListeners),
			apputil.Int(
				"rpcmaxclients",
				"Max number of RPC clients for standard connections",
				node.DefaultMaxRPCClients,
				cx.Config.RPCMaxClients),
			apputil.Int(
				"rpcmaxwebsockets",
				"Max number of RPC websocket connections",
				node.DefaultMaxRPCWebsockets,
				cx.Config.RPCMaxWebsockets),
			apputil.Int(
				"rpcmaxconcurrentreqs",
				"Max number of RPC requests that may be"+
					" processed concurrently",
				node.DefaultMaxRPCConcurrentReqs,
				cx.Config.RPCMaxConcurrentReqs),
			apputil.Bool(
				"rpcquirks",
				"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE:"+
					" Discouraged unless interoperability issues need to be worked"+
					" around",
				cx.Config.RPCQuirks),
			apputil.Bool(
				"norpc",
				"Disable built-in RPC server -- NOTE: The RPC server"+
					" is disabled by default if no rpcuser/rpcpass or"+
					" rpclimituser/rpclimitpass is specified",
				cx.Config.DisableRPC),
			apputil.Bool(
				"nodnsseed",
				"Disable DNS seeding for peers",
				cx.Config.DisableDNSSeed),
			apputil.StringSlice(
				"externalip",
				"Add an ip to the list of local addresses we claim to"+
					" listen on to peers",
				cx.Config.ExternalIPs),
			apputil.StringSlice(
				"addcheckpoint",
				"Add a custom checkpoint.  Format: '<height>:<hash>'",
				cx.Config.AddCheckpoints),
			apputil.Bool(
				"nocheckpoints",
				"Disable built-in checkpoints.  Don't do this unless"+
					" you know what you're doing.",
				cx.Config.DisableCheckpoints),
			apputil.String(
				"dbtype",
				"Database backend to use for the Block Chain",
				node.DefaultDbType,
				cx.Config.DbType),
			apputil.String(
				"profile",
				"Enable HTTP profiling on given port -- NOTE port"+
					" must be between 1024 and 65536",
				"",
				cx.Config.Profile),
			apputil.String(
				"cpuprofile",
				"Write CPU profile to the specified file",
				"",
				cx.Config.CPUProfile),
			apputil.Bool(
				"upnp",
				"Use UPnP to map our listening port outside of NAT",
				cx.Config.Upnp),
			apputil.Float64(
				"minrelaytxfee",
				"The minimum transaction fee in DUO/kB to be"+
					" considered a non-zero fee.",
				mempool.DefaultMinRelayTxFee.ToDUO(),
				cx.Config.MinRelayTxFee),
			apputil.Float64(
				"limitfreerelay",
				"Limit relay of transactions with no transaction"+
					" fee to the given amount in thousands of bytes per minute",
				node.DefaultFreeTxRelayLimit,
				cx.Config.FreeTxRelayLimit),
			apputil.Bool(
				"norelaypriority",
				"Do not require free or low-fee transactions to have"+
					" high priority for relaying",
				cx.Config.NoRelayPriority),
			apputil.Duration(
				"trickleinterval",
				"Minimum time between attempts to send new"+
					" inventory to a connected peer",
				node.DefaultTrickleInterval,
				cx.Config.TrickleInterval),
			apputil.Int(
				"maxorphantx",
				"Max number of orphan transactions to keep in memory",
				node.DefaultMaxOrphanTransactions,
				cx.Config.MaxOrphanTxs),
			apputil.String(
				"algo",
				"Sets the algorithm for the CPU miner ( blake14lr,"+
					" cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog,"+
					" skein, x11 default is 'random')",
				"random",
				cx.Config.Algo),
			apputil.Bool(
				"generate, g",
				"Generate (mine) DUO using the CPU",
				cx.Config.Generate),
			apputil.Int(
				"genthreads, G",
				"Number of CPU threads to use with CPU miner"+
					" -1 = all cores",
				-1,
				cx.Config.GenThreads),
			apputil.Bool(
				"broadcast",
				"enable broadcasting blocks for workers to mine on",
				cx.Config.Broadcast),
			apputil.StringSlice(
				"workers",
				"addresses to send out blocks to when broadcast is not enabled",
				cx.Config.Workers),
			apputil.Bool(
				"nocontroller",
				"disable zeroconf kcp miner controller",
				cx.Config.NoController),
			apputil.StringSlice(
				"miningaddr",
				"Add the specified payment address to the list of"+
					" addresses to use for generated blocks, at least one is "+
					"required if generate or minerlistener are set",
				cx.Config.MiningAddrs),
			apputil.String(
				"minerpass",
				"password to authorise sending work to a miner",
				genPassword(),
				cx.Config.MinerPass),
			apputil.Int(
				"blockminsize",
				"Minimum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMin,
				cx.Config.BlockMinSize),
			apputil.Int(
				"blockmaxsize",
				"Maximum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMax,
				cx.Config.BlockMaxSize),
			apputil.Int(
				"blockminweight",
				"Minimum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMin,
				cx.Config.BlockMinWeight),
			apputil.Int(
				"blockmaxweight",
				"Maximum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMax,
				cx.Config.BlockMaxWeight),
			apputil.Int(
				"blockprioritysize",
				"Size in bytes for high-priority/low-fee"+
					" transactions when creating a block",
				mempool.DefaultBlockPrioritySize,
				cx.Config.BlockPrioritySize),
			apputil.StringSlice(
				"uacomment",
				"Comment to add to the user agent -- See BIP 14 for"+
					" more information.",
				cx.Config.UserAgentComments),
			apputil.Bool(
				"nopeerbloomfilters",
				"Disable bloom filtering support",
				cx.Config.NoPeerBloomFilters),
			apputil.Bool(
				"nocfilters",
				"Disable committed filtering (CF) support",
				cx.Config.NoCFilters),
			apputil.Int(
				"sigcachemaxsize",
				"The maximum number of entries in the"+
					" signature verification cache",
				node.DefaultSigCacheMaxSize,
				cx.Config.SigCacheMaxSize),
			apputil.Bool(
				"blocksonly",
				"Do not accept transactions from remote peers.",
				cx.Config.BlocksOnly),
			apputil.BoolTrue(
				"notxindex",
				"Disable the transaction index which makes all transactions"+
					" available via the getrawtransaction RPC",
				cx.Config.TxIndex),
			apputil.BoolTrue(
				"noaddrindex",
				"Disable address-based transaction index which"+
					" makes the searchrawtransactions RPC available",
				cx.Config.AddrIndex,
			),
			apputil.Bool(
				"relaynonstd",
				"Relay non-standard transactions regardless of the default"+
					" settings for the active network.",
				cx.Config.RelayNonStd), apputil.Bool("rejectnonstd",
				"Reject non-standard transactions regardless of"+
					" the default settings for the active network.",
				cx.Config.RejectNonStd),
			apputil.Bool(
				"noinitialload",
				"Defer wallet creation/opening on startup and"+
					" enable loading wallets over RPC",
				cx.Config.NoInitialLoad),
			apputil.Bool(
				"walletconnect, wc",
				"connect to wallet instead of full node",
				cx.Config.Wallet),
			apputil.String(
				"walletserver, ws",
				"set wallet server to connect to",
				"127.0.0.1:11046",
				cx.Config.WalletServer),
			apputil.String(
				"walletpass",
				"The public wallet password -- Only required if"+
					" the wallet was created with one",
				"",
				cx.Config.WalletPass),
			apputil.Bool(
				"onetimetlskey",
				"Generate a new TLS certpair at startup, but"+
					" only write the certificate to disk",
				cx.Config.OneTimeTLSKey),
			apputil.Bool(
				"tlsskipverify",
				"skip verifying tls certificates",
				cx.Config.TLSSkipVerify),
			apputil.StringSlice(
				"walletrpclisten",
				"Listen for wallet RPC connections on this"+
					" interface/port (default port: 11046, testnet: 21046,"+
					" simnet: 41046)",
				cx.Config.WalletRPCListeners),
			apputil.Int(
				"walletrpcmaxclients",
				"Max number of legacy RPC clients for"+
					" standard connections",
				8,
				cx.Config.WalletRPCMaxClients),
			apputil.Int(
				"walletrpcmaxwebsockets",
				"Max number of legacy RPC websocket connections",
				8,
				cx.Config.WalletRPCMaxWebsockets,
			),
			apputil.StringSlice(
				"experimentalrpclisten",
				"Listen for RPC connections on this interface/port",
				cx.Config.ExperimentalRPCListeners),
			apputil.Bool(
				"nodeoff",
				"Starts GUI with node turned off",
				cx.Config.NodeOff),
			apputil.Bool( // TODO remove this
				"testnodeoff",
				"Starts GUI with testnode turned off",
				cx.Config.TestNodeOff),
			apputil.Bool(
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
