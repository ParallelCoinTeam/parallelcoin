package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"

	au "github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/kopach_worker"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/coding/base58"
	"github.com/p9c/pod/pkg/db/blockdb"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/util/interrupt"
)

// GetApp defines the pod app
func GetApp(cx *conte.Xt) (a *cli.App) {
	return &cli.App{
		Name:        "pod",
		Version:     "v0.0.1",
		Description: cx.Language.RenderText("goApp_DESCRIPTION"),
		Copyright:   cx.Language.RenderText("goApp_COPYRIGHT"),
		Action:      walletGUIHandle(cx),
		Before:      beforeFunc(cx),
		After: func(c *cli.Context) error {
			Trace("subcommand completed", os.Args)
			if interrupt.Restart {
			}
			return nil
		},
		Commands: []cli.Command{
			au.Command("version", "print version and exit",
				func(c *cli.Context) error {
					fmt.Println(c.App.Name, c.App.Version)
					return nil
				}, au.SubCommands(), nil, "v"),
			// apputil.NewCommand("gui", "run GUI",
			//	guiHandle(cx), apputil.SubCommands(), nil, "gui"),
			au.Command("gui", "start wallet GUI", walletGUIHandle(cx),
				au.SubCommands(), nil),
			au.Command("explorer", "start explorer GUI", explorerHandle(cx),
				au.SubCommands(), nil),
			au.Command("nodegui", "start node GUI", nodeGUIHandle(cx),
				au.SubCommands(), nil),
			au.Command("ctl",
				"send RPC commands to a node or wallet and print the result",
				ctlHandle(cx), au.SubCommands(
					au.Command(
						"listcommands",
						"list commands available at endpoint",
						ctlHandleList,
						au.SubCommands(),
						nil,
						"list",
						"l",
					),
				), nil, "c"),
			au.Command("node", "start parallelcoin full node",
				nodeHandle(cx), au.SubCommands(
					au.Command("dropaddrindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropAddrIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command("droptxindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropTxIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command("dropindexes",
						"drop all of the indexes",
						func(c *cli.Context) error {
							cx.StateCfg.DropAddrIndex = true
							cx.StateCfg.DropTxIndex = true
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command("dropcfindex",
						"drop the address search index",
						func(c *cli.Context) error {
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command("resetchain",
						"reset the chain",
						func(c *cli.Context) (err error) {
							config.Configure(cx, c.Command.Name, true)
							dbName := blockdb.NamePrefix + "_" + *cx.Config.DbType
							if *cx.Config.DbType == "sqlite" {
								dbName += ".db"
							}
							dbPath := filepath.Join(filepath.Join(*cx.Config.DataDir,
								cx.ActiveNet.Name), dbName)
							if err = os.RemoveAll(dbPath); Check(err) {
							}
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
				), nil, "n"),
			au.Command("wallet", "start parallelcoin wallet server",
				WalletHandle(cx), au.SubCommands(
					au.Command("drophistory", "drop the transaction history in the wallet (for "+
						"development and testing as well as clearing up transaction mess)",
						func(c *cli.Context) (err error) {
							config.Configure(cx, c.Command.Name, true)
							Info("dropping wallet history")
							go func() {
								Warn("starting wallet")
								if err = walletmain.Main(cx); Check(err) {
									// os.Exit(1)
								} else {
									Debug("wallet started")
								}
							}()
							Debug("waiting for walletChan")
							cx.WalletServer = <-cx.WalletChan
							Debug("walletChan sent")
							err = legacy.DropWalletHistory(cx.WalletServer)(c)
							return
						}, au.SubCommands(), nil),
				), nil, "w"),
			au.Command("shell", "start combined wallet/node shell",
				ShellHandle(cx), au.SubCommands(), nil, "s"),
			au.Command("kopach", "standalone miner for clusters",
				KopachHandle(cx), au.SubCommands(), nil, "k"),
			au.Command(
				"worker",
				"single thread parallelcoin miner controlled with binary IPC interface on stdin/stdout; "+
					"internal use, must have network name string as second arg after worker and nothing before;"+
					" communicates via net/rpc encoding/gob as default over stdio",
				kopach_worker.KopachWorkerHandle(cx),
				au.SubCommands(),
				nil,
			),
			au.Command("init",
				"steps through creation of new wallet and initialization for a network with these specified "+
					"in the main",
				initHandle(cx),
				au.SubCommands(),
				nil,
				"I"),
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "datadir, D",
				Value:       *cx.Config.DataDir,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_DATADIR",
				Destination: cx.Config.DataDir,
			},
			cli.BoolFlag{
				Name: "pipelog, P",
				Usage: "enables pipe logger (" +
					"setting only activates on use of cli flag or environment" +
					" variable as it alters stdin/out behaviour)",
				EnvVar:      "POD_PIPELOG",
				Destination: cx.Config.PipeLog,
			},
			cli.StringFlag{
				Name:        "lang, L",
				Value:       *cx.Config.Language,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_LANGUAGE",
				Destination: cx.Config.Language,
			},
			cli.StringFlag{
				Name:        "walletfile, WF",
				Value:       *cx.Config.WalletFile,
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_WALLETFILE",
				Destination: cx.Config.WalletFile,
			},
			au.BoolTrue("save, i",
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
			au.String(
				"network, n",
				"connect to mainnet/testnet/regtest/simnet",
				"mainnet",
				cx.Config.Network),
			au.String(
				"username",
				"sets the username for services",
				"server",
				cx.Config.Username),
			au.String(
				"password",
				"sets the password for services",
				genPassword(),
				cx.Config.Password),
			au.String(
				"serveruser",
				"sets the username for clients of services",
				"client",
				cx.Config.ServerUser),
			au.String(
				"serverpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.ServerPass),
			au.String(
				"limituser",
				"sets the limited rpc username",
				"limit",
				cx.Config.LimitUser),
			au.String(
				"limitpass",
				"sets the limited rpc password",
				genPassword(),
				cx.Config.LimitPass),
			au.String(
				"rpccert",
				"File containing the certificate file",
				"",
				cx.Config.RPCCert),
			au.String(
				"rpckey",
				"File containing the certificate key",
				"",
				cx.Config.RPCKey),
			au.String(
				"cafile",
				"File containing root certificates to authenticate a TLS"+
					" connections with pod",
				"",
				cx.Config.CAFile),
			au.BoolTrue(
				"clienttls",
				"Enable TLS for client connections",
				cx.Config.TLS),
			au.BoolTrue(
				"servertls",
				"Enable TLS for server connections",
				cx.Config.ServerTLS),
			au.String(
				"proxy",
				"Connect via SOCKS5 proxy",
				"",
				cx.Config.Proxy),
			au.String(
				"proxyuser",
				"Username for proxy server",
				"user",
				cx.Config.ProxyUser),
			au.String(
				"proxypass",
				"Password for proxy server",
				"pa55word",
				cx.Config.ProxyPass),
			au.Bool(
				"onion",
				"Enable connecting to tor hidden services",
				cx.Config.Onion),
			au.String(
				"onionproxy",
				"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0."+
					"0.1:9050)",
				"127.0.0.1:9050",
				cx.Config.OnionProxy),
			au.String(
				"onionuser",
				"Username for onion proxy server",
				"user",
				cx.Config.OnionProxyUser),
			au.String(
				"onionpass",
				"Password for onion proxy server",
				genPassword(),
				cx.Config.OnionProxyPass),
			au.Bool(
				"torisolation",
				"Enable Tor stream isolation by randomizing user credentials"+
					" for each connection.",
				cx.Config.TorIsolation),
			au.StringSlice(
				"addpeer",
				"Add a peer to connect with at startup",
				cx.Config.AddPeers),
			au.StringSlice(
				"connect",
				"Connect only to the specified peers at startup",
				cx.Config.ConnectPeers),
			au.Bool(
				"nolisten",
				"Disable listening for incoming connections -- NOTE:"+
					" Listening is automatically disabled if the --connect or"+
					" --proxy options are used without also specifying listen"+
					" interfaces via --listen",
				cx.Config.DisableListen),
			au.StringSlice(
				"listen",
				"Add an interface/port to listen for connections",
				cx.Config.Listeners),
			au.Int(
				"maxpeers",
				"Max number of inbound and outbound peers",
				node.DefaultMaxPeers,
				cx.Config.MaxPeers),
			au.Bool(
				"nobanning",
				"Disable banning of misbehaving peers",
				cx.Config.DisableBanning),
			au.Duration(
				"banduration",
				"How long to ban misbehaving peers",
				time.Hour*24,
				cx.Config.BanDuration),
			au.Int(
				"banthreshold",
				"Maximum allowed ban score before disconnecting and"+
					" banning misbehaving peers.",
				node.DefaultBanThreshold,
				cx.Config.BanThreshold),
			au.StringSlice(
				"whitelist",
				"Add an IP network or IP that will not be banned. (eg. 192."+
					"168.1.0/24 or ::1)",
				cx.Config.Whitelists),
			au.String(
				"rpcconnect",
				"Hostname/IP and port of pod RPC server to connect to",
				"",
				cx.Config.RPCConnect),
			au.StringSlice(
				"rpclisten",
				"Add an interface/port to listen for RPC connections",
				cx.Config.RPCListeners),
			au.Int(
				"rpcmaxclients",
				"Max number of RPC clients for standard connections",
				node.DefaultMaxRPCClients,
				cx.Config.RPCMaxClients),
			au.Int(
				"rpcmaxwebsockets",
				"Max number of RPC websocket connections",
				node.DefaultMaxRPCWebsockets,
				cx.Config.RPCMaxWebsockets),
			au.Int(
				"rpcmaxconcurrentreqs",
				"Max number of RPC requests that may be"+
					" processed concurrently",
				node.DefaultMaxRPCConcurrentReqs,
				cx.Config.RPCMaxConcurrentReqs),
			au.Bool(
				"rpcquirks",
				"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE:"+
					" Discouraged unless interoperability issues need to be worked"+
					" around",
				cx.Config.RPCQuirks),
			au.Bool(
				"norpc",
				"Disable built-in RPC server -- NOTE: The RPC server"+
					" is disabled by default if no rpcuser/rpcpass or"+
					" rpclimituser/rpclimitpass is specified",
				cx.Config.DisableRPC),
			au.Bool(
				"nodnsseed",
				"Disable DNS seeding for peers",
				cx.Config.DisableDNSSeed),
			au.StringSlice(
				"externalip",
				"Add an ip to the list of local addresses we claim to"+
					" listen on to peers",
				cx.Config.ExternalIPs),
			au.StringSlice(
				"addcheckpoint",
				"Add a custom checkpoint.  Format: '<height>:<hash>'",
				cx.Config.AddCheckpoints),
			au.Bool(
				"nocheckpoints",
				"Disable built-in checkpoints.  Don't do this unless"+
					" you know what you're doing.",
				cx.Config.DisableCheckpoints),
			au.String(
				"dbtype",
				"Database backend to use for the Block Chain",
				node.DefaultDbType,
				cx.Config.DbType),
			au.String(
				"profile",
				"Enable HTTP profiling on given port -- NOTE port"+
					" must be between 1024 and 65536",
				"",
				cx.Config.Profile),
			au.String(
				"cpuprofile",
				"Write CPU profile to the specified file",
				"",
				cx.Config.CPUProfile),
			au.Bool(
				"upnp",
				"Use UPnP to map our listening port outside of NAT",
				cx.Config.UPNP),
			au.Float64(
				"minrelaytxfee",
				"The minimum transaction fee in DUO/kB to be"+
					" considered a non-zero fee.",
				mempool.DefaultMinRelayTxFee.ToDUO(),
				cx.Config.MinRelayTxFee),
			au.Float64(
				"limitfreerelay",
				"Limit relay of transactions with no transaction"+
					" fee to the given amount in thousands of bytes per minute",
				node.DefaultFreeTxRelayLimit,
				cx.Config.FreeTxRelayLimit),
			au.Bool(
				"norelaypriority",
				"Do not require free or low-fee transactions to have"+
					" high priority for relaying",
				cx.Config.NoRelayPriority),
			au.Duration(
				"trickleinterval",
				"Minimum time between attempts to send new"+
					" inventory to a connected peer",
				node.DefaultTrickleInterval,
				cx.Config.TrickleInterval),
			au.Int(
				"maxorphantx",
				"Max number of orphan transactions to keep in memory",
				node.DefaultMaxOrphanTransactions,
				cx.Config.MaxOrphanTxs),
			au.Bool(
				"generate, g",
				"Generate (mine) DUO using the CPU",
				cx.Config.Generate),
			au.Int(
				"genthreads, G",
				"Number of CPU threads to use with CPU miner"+
					" -1 = all cores",
				1,
				cx.Config.GenThreads),
			au.Bool(
				"solo",
				"mine DUO even if not connected to the network",
				cx.Config.Solo),
			au.Bool(
				"lan",
				"mine duo if not connected to nodes on internet",
				cx.Config.LAN),
			au.String(
				"controller",
				"port controller listens on for solutions from workers"+
					" and other node peers",
				":0",
				cx.Config.Controller),
			au.Bool(
				"autoports",
				"uses random automatic ports for p2p, rpc and controller",
				cx.Config.AutoPorts),
			au.StringSlice(
				"miningaddr",
				"Add the specified payment address to the list of"+
					" addresses to use for generated blocks, at least one is "+
					"required if generate or minerlistener are set",
				cx.Config.MiningAddrs),
			au.String(
				"minerpass",
				"password to authorise sending work to a miner",
				genPassword(),
				cx.Config.MinerPass),
			au.Int(
				"blockminsize",
				"Minimum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMin,
				cx.Config.BlockMinSize),
			au.Int(
				"blockmaxsize",
				"Maximum block size in bytes to be used when"+
					" creating a block",
				node.BlockMaxSizeMax,
				cx.Config.BlockMaxSize),
			au.Int(
				"blockminweight",
				"Minimum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMin,
				cx.Config.BlockMinWeight),
			au.Int(
				"blockmaxweight",
				"Maximum block weight to be used when creating"+
					" a block",
				node.BlockMaxWeightMax,
				cx.Config.BlockMaxWeight),
			au.Int(
				"blockprioritysize",
				"Size in bytes for high-priority/low-fee"+
					" transactions when creating a block",
				mempool.DefaultBlockPrioritySize,
				cx.Config.BlockPrioritySize),
			au.StringSlice(
				"uacomment",
				"Comment to add to the user agent -- See BIP 14 for"+
					" more information.",
				cx.Config.UserAgentComments),
			au.Bool(
				"nopeerbloomfilters",
				"Disable bloom filtering support",
				cx.Config.NoPeerBloomFilters),
			au.Bool(
				"nocfilters",
				"Disable committed filtering (CF) support",
				cx.Config.NoCFilters),
			au.Int(
				"sigcachemaxsize",
				"The maximum number of entries in the"+
					" signature verification cache",
				node.DefaultSigCacheMaxSize,
				cx.Config.SigCacheMaxSize),
			au.Bool(
				"blocksonly",
				"Do not accept transactions from remote peers.",
				cx.Config.BlocksOnly),
			au.BoolTrue(
				"txindex",
				"Disable the transaction index which makes all transactions available via the getrawtransaction RPC",
				cx.Config.TxIndex),
			au.BoolTrue(
				"addrindex",
				"Disable address-based transaction index which makes the searchrawtransactions RPC available",
				cx.Config.AddrIndex,
			),
			au.Bool(
				"relaynonstd",
				"Relay non-standard transactions regardless of the default settings for the active network.",
				cx.Config.RelayNonStd), au.Bool("rejectnonstd",
				"Reject non-standard transactions regardless of the default settings for the active network.",
				cx.Config.RejectNonStd),
			au.Bool(
				"noinitialload",
				"Defer wallet creation/opening on startup and enable loading wallets over RPC (loading not yet implemented)",
				cx.Config.NoInitialLoad),
			au.Bool(
				"walletconnect, wc",
				"connect to wallet instead of full node",
				cx.Config.Wallet),
			au.String(
				"walletserver, ws",
				"set wallet server to connect to",
				"127.0.0.1:11046",
				cx.Config.WalletServer),
			au.String(
				"walletpass",
				"The public wallet password -- Only required if the wallet was created with one",
				"",
				cx.Config.WalletPass),
			au.Bool(
				"onetimetlskey",
				"Generate a new TLS certificate pair at startup, but only write the certificate to disk",
				cx.Config.OneTimeTLSKey),
			au.Bool(
				"tlsskipverify",
				"skip verifying tls certificates",
				cx.Config.TLSSkipVerify),
			au.StringSlice(
				"walletrpclisten",
				"Listen for wallet RPC connections on this"+
					" interface/port (default port: 11046, testnet: 21046,"+
					" simnet: 41046)",
				cx.Config.WalletRPCListeners),
			au.Int(
				"walletrpcmaxclients",
				"Max number of legacy RPC clients for"+
					" standard connections",
				8,
				cx.Config.WalletRPCMaxClients),
			au.Int(
				"walletrpcmaxwebsockets",
				"Max number of legacy RPC websocket connections",
				8,
				cx.Config.WalletRPCMaxWebsockets,
			),
			au.Bool(
				"nodeoff",
				"Starts with node turned off",
				cx.Config.NodeOff),
			au.Bool(
				"walletoff",
				"Starts with wallet turned off",
				cx.Config.WalletOff,
			),
			au.Bool(
				"delaystart",
				"pauses for 3 seconds before starting, for internal use with restart function",
				nil,
			),
			au.Bool(
				"kopachgui",
				"enables the GUI for the kopach miner",
				cx.Config.KopachGUI,
			),
			au.Bool(
				"gui",
				"enables the GUI",
				cx.Config.GUI,
			),
			au.Bool(
				"darktheme",
				"sets the dark theme on the gui interface",
				cx.Config.DarkTheme,
			),
			au.Bool(
				"notty",
				"tells pod there is no keyboard input available",
				nil,
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
