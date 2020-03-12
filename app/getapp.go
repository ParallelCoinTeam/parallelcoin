package app

import (
	"os"
	"path/filepath"
	"time"

	"github.com/p9c/util/interrupt"

	"github.com/p9c/rpc/legacy"
	"github.com/p9c/wallet/walletmain"

	"github.com/p9c/kopach/kopach_worker"
	"github.com/p9c/pod/cmd/node/blockdb"

	"github.com/urfave/cli"

	log "github.com/p9c/logi"
	"github.com/p9c/util/base58"
	"github.com/p9c/util/hdkeychain"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/pkg/conte"
)

func // GetApp defines the pod app
GetApp(cx *conte.Xt) (a *cli.App) {
	return &cli.App{
		Name:        "pod",
		Version:     "v0.0.1",
		Description: cx.Language.RenderText("goApp_DESCRIPTION"),
		Copyright:   cx.Language.RenderText("goApp_COPYRIGHT"),
		Action:      guiHandle(cx),
		Before:      beforeFunc(cx),
		After: func(c *cli.Context) error {
			log.L.Trace("subcommand completed")
			if interrupt.Restart {
			}
			return nil
		},
		Commands: []cli.Command{
			apputil.NewCommand("version", "print version and exit", func(c *cli.Context) error {
				log.Println(c.App.Name, c.App.Version)
				return nil
			}, apputil.SubCommands(), nil, "v"),
			apputil.NewCommand("ctl", "send RPC commands to a node or wallet and print the result", ctlHandle(cx), apputil.SubCommands(
				apputil.NewCommand(
					"listcommands",
					"list commands available at endpoint",
					ctlHandleList,
					apputil.SubCommands(),
					nil,
					"list",
					"l",
				),
			), nil, "c"),
			apputil.NewCommand("ctlgui", "GUI interface send RPC commands to a node or wallet and print the result", ctlGUIHandle(cx), apputil.SubCommands(), nil, "C"),
			apputil.NewCommand("node", "start parallelcoin full node", nodeHandle(cx), apputil.SubCommands(
				apputil.NewCommand("dropaddrindex",
					"drop the address search index",
					func(c *cli.Context) error {
						cx.StateCfg.DropAddrIndex = true
						// return nodeHandle(cx)(c)
						return nil
					},
					apputil.SubCommands(),
					nil,
				),
				apputil.NewCommand("droptxindex",
					"drop the address search index",
					func(c *cli.Context) error {
						cx.StateCfg.DropTxIndex = true
						// return nodeHandle(cx)(c)
						return nil
					},
					apputil.SubCommands(),
					nil,
				),
				apputil.NewCommand("dropindexes",
					"drop all of the indexes",
					func(c *cli.Context) error {
						cx.StateCfg.DropAddrIndex = true
						cx.StateCfg.DropTxIndex = true
						cx.StateCfg.DropCfIndex = true
						// return nodeHandle(cx)(c)
						return nil
					},
					apputil.SubCommands(),
					nil,
				),
				apputil.NewCommand("dropcfindex",
					"drop the address search index",
					func(c *cli.Context) error {
						cx.StateCfg.DropCfIndex = true
						// return nodeHandle(cx)(c)
						return nil
					},
					apputil.SubCommands(),
					nil,
				),
				apputil.NewCommand("resetchain",
					"reset the chain",
					func(c *cli.Context) (err error) {
						dbName := blockdb.NamePrefix + "_" + *cx.Config.DbType
						if *cx.Config.DbType == "sqlite" {
							dbName += ".db"
						}
						dbPath := filepath.Join(filepath.Join(*cx.Config.DataDir,
							cx.ActiveNet.Name), dbName)
						if err = os.RemoveAll(dbPath); log.L.Check(err) {
						}
						// return nodeHandle(cx)(c)
						return nil
					},
					apputil.SubCommands(),
					nil,
				),
			), nil, "n"),
			apputil.NewCommand("wallet", "start parallelcoin wallet server", WalletHandle(cx), apputil.SubCommands(
				apputil.NewCommand("drophistory",
					"drop the transaction history in the wallet ("+
						"for development and testing as well as clearing up"+
						" transaction mess)",
					func(c *cli.Context) (err error) {
						Configure(cx, c)
						log.L.Info("dropping wallet history")
						go func() {
							log.L.Warn("starting wallet")
							if err = walletmain.Main(cx); log.L.Check(err) {
								os.Exit(1)
							} else {
								log.L.Debug("wallet started")
							}
						}()
						log.L.Debug("waiting for walletChan")
						cx.WalletServer = <-cx.WalletChan
						log.L.Debug("walletChan sent")
						err = legacy.DropWalletHistory(cx.WalletServer)(c)
						return
					},
					apputil.SubCommands(),
					nil,
				),
			), nil, "w"),
			apputil.NewCommand("shell", "start combined wallet/node shell", shellHandle(cx), apputil.SubCommands(), nil, "s"),
			apputil.NewCommand("gui", "start GUI", guiHandle(cx), apputil.SubCommands(), nil),
			apputil.NewCommand("kopach", "standalone miner for clusters", KopachHandle(cx), apputil.SubCommands(), nil, "k"),
			apputil.NewCommand("worker", "single thread parallelcoin miner controlled with binary IPC"+
				" interface on stdin/stdout; internal use, must have network name string as second arg after worker and"+
				"nothing before; communicates via net/rpc encoding/gob as default over stdio", kopach_worker.KopachWorkerHandle(cx), apputil.SubCommands(), nil),
			apputil.NewCommand("init",
				"steps through creation of new wallet and initialization for a network with these specified in the main",
				initHandle(cx),
				apputil.SubCommands(),
				nil,
				"I"),
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "lang, L",
				Value:       *cx.Config.Language,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_LANGUAGE",
				Destination: cx.Config.Language,
			},
			cli.StringFlag{
				Name:        "datadir, D",
				Value:       *cx.Config.DataDir,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_DATADIR",
				Destination: cx.Config.DataDir,
			},
			cli.StringFlag{
				Name:        "walletfile, WF",
				Value:       *cx.Config.WalletFile,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_WALLETFILE",
				Destination: cx.Config.WalletFile,
			},
			apputil.BoolTrue("save, i",
				"save settings as effective from invocation",
				&cx.StateCfg.Save,
			),
			cli.StringFlag{
				Name:        "loglevel, l",
				Value:       *cx.Config.LogLevel,
				Usage:       "sets the base for all subsystem logging",
				EnvVar:      "POD_LOGLEVEL",
				Destination: cx.Config.LogLevel,
			},
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
				"sets the limited rpc password",
				genPassword(),
				cx.Config.LimitPass),
			apputil.String(
				"rpccert",
				"File containing the certificate file",
				"",
				cx.Config.RPCCert),
			apputil.String(
				"rpckey",
				"File containing the certificate key",
				"",
				cx.Config.RPCKey),
			apputil.String(
				"cafile",
				"File containing root certificates to authenticate a TLS"+
					" connections with pod",
				"",
				cx.Config.CAFile),
			apputil.BoolTrue(
				"clienttls",
				"Enable TLS for client connections",
				cx.Config.TLS),
			apputil.BoolTrue(
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
				"",
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
				cx.Config.UPNP),
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
				// TODO: remove this as mining only one algo is
				//  not advisable
				"algo",
				"Sets the algorithm for the CPU miner ( blake14lr,"+
					" cn7v2, keccak, lyra2rev2, scrypt, sha256d, stribog,"+
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
				"solo",
				"mine DUO even if not connected to the network",
				cx.Config.Solo),
			apputil.Bool(
				"lan",
				"mine duo if not connected to nodes on internet",
				cx.Config.LAN),
			apputil.String(
				"controller",
				"port controller listens on for solutions from workers"+
					" and other node peers",
				":0",
				cx.Config.Controller),
			apputil.Bool(
				"autoports",
				"uses random automatic ports for p2p, rpc and controller",
				cx.Config.AutoPorts),
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
			apputil.Bool(
				"notxindex",
				"Disable the transaction index which makes all transactions"+
					" available via the getrawtransaction RPC",
				cx.Config.TxIndex),
			apputil.Bool(
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
			// apputil.StringSlice(
			// 	"experimentalrpclisten",
			// 	"Listen for RPC connections on this interface/port",
			// 	cx.Config.ExperimentalRPCListeners),
			apputil.Bool(
				"nodeoff",
				"Starts GUI with node turned off",
				cx.Config.NodeOff),
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
