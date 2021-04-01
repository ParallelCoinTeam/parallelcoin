package app

import (
	"fmt"
	"github.com/p9c/pod/pkg/constant"
	"github.com/p9c/pod/pkg/opts"
	"github.com/p9c/pod/pkg/pod"
	walletrpc2 "github.com/p9c/pod/pkg/walletrpc"
	"github.com/p9c/pod/version"
	"os"
	"path/filepath"
	"time"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/kopach_worker"
	"github.com/p9c/pod/cmd/node/mempool"
	"github.com/p9c/pod/cmd/walletmain"
	au "github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/base58"
	"github.com/p9c/pod/pkg/database/blockdb"
	"github.com/p9c/pod/pkg/podconfig"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/util/interrupt"
)

// getApp defines the pod app
func getApp(cx *pod.State) (a *cli.App) {
	return &cli.App{
		Name:        "pod",
		Version:     version.Get(),
		Description: cx.Language.RenderText("goApp_DESCRIPTION"),
		Copyright:   cx.Language.RenderText("goApp_COPYRIGHT"),
		Action:      GUIHandle(cx),
		Before:      beforeFunc(cx),
		After: func(c *cli.Context) (e error) {
			opts.D.Ln("subcommand completed", os.Args)
			if interrupt.Restart {
			}
			return nil
		},
		Commands: []cli.Command{
			au.Command(
				"version", "print version and exit",
				func(c *cli.Context) (e error) {
					fmt.Println(c.App.Name, c.App.Version)
					return nil
				}, au.SubCommands(), nil, "v",
			),
			au.Command(
				"gui", "start wallet GUI", GUIHandle(cx),
				au.SubCommands(), nil,
			),
			au.Command(
				"ctl",
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
				), nil, "c",
			),
			au.Command(
				"rpc", "start parallelcoin full node for vps/rpc services usage",
				rpcNodeHandle(cx),
				au.SubCommands(
					au.Command(
						"dropaddrindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropAddrIndex = true
							return nodeHandle(cx)(c)
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"droptxindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropTxIndex = true
							return nodeHandle(cx)(c)
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"dropindexes",
						"drop all of the indexes",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropAddrIndex = true
							cx.StateCfg.DropTxIndex = true
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"dropcfindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"resetchain",
						"reset the chain",
						func(c *cli.Context) (e error) {
							podconfig.Configure(cx, true)
							dbName := blockdb.NamePrefix + "_" + cx.Config.DbType.V()
							if cx.Config.DbType.V() == "sqlite" {
								dbName += ".db"
							}
							dbPath := filepath.Join(
								filepath.Join(
									cx.Config.DataDir.V(),
									cx.ActiveNet.Name,
								), dbName,
							)
							if e = os.RemoveAll(dbPath); opts.E.Chk(e) {
							}
							return nodeHandle(cx)(c)
						},
						au.SubCommands(),
						nil,
					),
				), nil, "n",
			),
			au.Command(
				"node", "start parallelcoin full node",
				nodeHandle(cx), au.SubCommands(
					au.Command(
						"dropaddrindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropAddrIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"droptxindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropTxIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"dropindexes",
						"drop all of the indexes",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropAddrIndex = true
							cx.StateCfg.DropTxIndex = true
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"dropcfindex",
						"drop the address search index",
						func(c *cli.Context) (e error) {
							cx.StateCfg.DropCfIndex = true
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
					au.Command(
						"resetchain",
						"reset the chain",
						func(c *cli.Context) (e error) {
							podconfig.Configure(cx, true)
							dbName := blockdb.NamePrefix + "_" + cx.Config.DbType.V()
							if cx.Config.DbType.V() == "sqlite" {
								dbName += ".db"
							}
							dbPath := filepath.Join(
								filepath.Join(
									cx.Config.DataDir.V(),
									cx.ActiveNet.Name,
								), dbName,
							)
							if e = os.RemoveAll(dbPath); opts.E.Chk(e) {
							}
							return nodeHandle(cx)(c)
							// return nil
						},
						au.SubCommands(),
						nil,
					),
				), nil, "n",
			),
			au.Command(
				"wallet", "start parallelcoin wallet server",
				WalletHandle(cx), au.SubCommands(
					au.Command(
						"drophistory", "drop the transaction history in the wallet (for "+
							"development and testing as well as clearing up transaction mess)",
						func(c *cli.Context) (e error) {
							podconfig.Configure(cx, true)
							opts.I.Ln("dropping wallet history")
							go func() {
								opts.D.Ln("starting wallet")
								if e = walletmain.Main(cx); opts.E.Chk(e) {
									// os.Exit(1)
								} else {
									opts.D.Ln("wallet started")
								}
							}()
							// D.Ln("waiting for walletChan")
							// cx.WalletServer = <-cx.WalletChan
							// D.Ln("walletChan sent")
							e = walletrpc2.DropWalletHistory(cx.WalletServer, cx.Config)(c)
							return
						}, au.SubCommands(), nil,
					),
				), nil, "w",
			),
			au.Command(
				"shell", "start combined wallet/node shell",
				ShellHandle(cx), au.SubCommands(), nil, "s",
			),
			au.Command(
				"kopach", "standalone miner for clusters",
				KopachHandle(cx), au.SubCommands(), nil, "k",
			),
			au.Command(
				"worker",
				"single thread parallelcoin miner controlled with binary IPC interface on stdin/stdout; "+
					"internal use, must have network name string as second arg after worker and nothing before;"+
					" communicates via net/rpc encoding/gob as default over stdio",
				kopach_worker.KopachWorkerHandle(cx),
				au.SubCommands(),
				nil,
			),
			au.Command(
				"init",
				"steps through creation of new wallet and initialization for a network with these specified "+
					"in the main",
				initHandle(cx),
				au.SubCommands(),
				nil,
				"I",
			),
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "datadir, D",
				Value:       cx.Config.DataDir.V(),
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_DATADIR",
				Destination: cx.Config.DataDir.Ptr(),
			},
			cli.BoolFlag{
				Name: "pipelog, P",
				Usage: "enables pipe logger (" +
					"setting only activates on use of cli flag or environment" +
					" variable as it alters stdin/out behaviour)",
				EnvVar:      "POD_PIPELOG",
				Destination: cx.Config.PipeLog.Ptr(),
			},
			cli.StringFlag{
				Name:        "lang, L",
				Value:       cx.Config.Language.V(),
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_LANGUAGE",
				Destination: cx.Config.Language.Ptr(),
			},
			cli.StringFlag{
				Name:        "walletfile, WF",
				Value:       cx.Config.WalletFile.V(),
				Usage:       "sets the data directory base for a pod instance",
				EnvVar:      "POD_WALLETFILE",
				Destination: cx.Config.WalletFile.Ptr(),
			},
			au.BoolTrue(
				"save, i",
				"save settings as effective from invocation",
				&cx.StateCfg.Save,
			),
			cli.StringFlag{
				Name:        "loglevel, l",
				Value:       cx.Config.LogLevel.V(),
				Usage:       "sets the base for all subsystem logging",
				EnvVar:      "POD_LOGLEVEL",
				Destination: cx.Config.LogLevel.Ptr(),
			},
			au.StringSlice(
				"highlight",
				"define the set of packages whose logs will have attention-grabbing keywords to aid scanning logs",
				cx.Config.Hilite.Ptr(),
			),
			au.StringSlice(
				"logfilter",
				"define the set of packages whose logs will not print",
				cx.Config.LogFilter.Ptr(),
			),
			au.String(
				"network, n",
				"connect to mainnet/testnet/regtest/simnet",
				"mainnet",
				cx.Config.Network.Ptr(),
			),
			au.String(
				"username",
				"sets the username for services",
				"server",
				cx.Config.Username.Ptr(),
			),
			au.String(
				"password",
				"sets the password for services",
				genPassword(),
				cx.Config.Password.Ptr(),
			),
			au.String(
				"serveruser",
				"sets the username for clients of services",
				"client",
				cx.Config.ServerUser.Ptr(),
			),
			au.String(
				"serverpass",
				"sets the password for clients of services",
				genPassword(),
				cx.Config.ServerPass.Ptr(),
			),
			au.String(
				"limituser",
				"sets the limited rpc username",
				"limit",
				cx.Config.LimitUser.Ptr(),
			),
			au.String(
				"limitpass",
				"sets the limited rpc password",
				genPassword(),
				cx.Config.LimitPass.Ptr(),
			),
			au.String(
				"rpccert",
				"File containing the certificate file",
				"",
				cx.Config.RPCCert.Ptr(),
			),
			au.String(
				"rpckey",
				"File containing the certificate key",
				"",
				cx.Config.RPCKey.Ptr(),
			),
			au.String(
				"cafile",
				"File containing root certificates to authenticate a TLS"+
					" connections with pod",
				"",
				cx.Config.CAFile.Ptr(),
			),
			au.BoolTrue(
				"clienttls",
				"Enable TLS for client connections",
				cx.Config.TLS.Ptr(),
			),
			au.BoolTrue(
				"servertls",
				"Enable TLS for server connections",
				cx.Config.ServerTLS.Ptr(),
			),
			au.String(
				"proxy",
				"Connect via SOCKS5 proxy",
				"",
				cx.Config.Proxy.Ptr(),
			),
			au.String(
				"proxyuser",
				"Username for proxy server",
				"user",
				cx.Config.ProxyUser.Ptr(),
			),
			au.String(
				"proxypass",
				"Password for proxy server",
				"pa55word",
				cx.Config.ProxyPass.Ptr(),
			),
			au.Bool(
				"onion",
				"Enable connecting to tor hidden services",
				cx.Config.Onion.Ptr(),
			),
			au.String(
				"onionproxy",
				"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0."+
					"0.1:9050)",
				"127.0.0.1:9050",
				cx.Config.OnionProxy.Ptr(),
			),
			au.String(
				"onionuser",
				"Username for onion proxy server",
				"user",
				cx.Config.OnionProxyUser.Ptr(),
			),
			au.String(
				"onionpass",
				"Password for onion proxy server",
				genPassword(),
				cx.Config.OnionProxyPass.Ptr(),
			),
			au.Bool(
				"torisolation",
				"Enable Tor stream isolation by randomizing user credentials"+
					" for each connection.",
				cx.Config.TorIsolation.Ptr(),
			),
			au.StringSlice(
				"addpeer",
				"Add a peer to connect with at startup",
				cx.Config.AddPeers.Ptr(),
			),
			au.StringSlice(
				"connect",
				"Connect only to the specified peers at startup",
				cx.Config.ConnectPeers.Ptr(),
			),
			au.Bool(
				"nolisten",
				"Disable listening for incoming connections -- NOTE:"+
					" Listening is automatically disabled if the --connect or"+
					" --proxy options are used without also specifying listen"+
					" interfaces via --listen",
				cx.Config.DisableListen.Ptr(),
			),
			au.BoolTrue(
				"autolisten",
				"enable automatically populating p2p and controller reachable addresses",
				cx.Config.AutoListen.Ptr(),
			),
			au.StringSlice(
				"p2pconnect",
				"Addresses that are configured to receive inbound connections",
				cx.Config.P2PConnect.Ptr(),
			),
			au.StringSlice(
				"listen",
				"Add an interface/port to listen for connections",
				cx.Config.P2PListeners.Ptr(),
			),
			au.Int(
				"maxpeers",
				"Max number of inbound and outbound peers",
				constant.DefaultMaxPeers,
				cx.Config.MaxPeers.Ptr(),
			),
			au.Bool(
				"nobanning",
				"Disable banning of misbehaving peers",
				cx.Config.DisableBanning.Ptr(),
			),
			au.Duration(
				"banduration",
				"How long to ban misbehaving peers",
				time.Hour*24,
				cx.Config.BanDuration.Ptr(),
			),
			au.Int(
				"banthreshold",
				"Maximum allowed ban score before disconnecting and"+
					" banning misbehaving peers.",
				constant.DefaultBanThreshold,
				cx.Config.BanThreshold.Ptr(),
			),
			au.StringSlice(
				"whitelist",
				"Add an IP network or IP that will not be banned. (eg. 192."+
					"168.1.0/24 or ::1)",
				cx.Config.Whitelists.Ptr(),
			),
			au.String(
				"rpcconnect",
				"Hostname/IP and port of pod RPC server to connect to",
				"",
				cx.Config.RPCConnect.Ptr(),
			),
			au.StringSlice(
				"rpclisten",
				"Add an interface/port to listen for RPC connections",
				cx.Config.RPCListeners.Ptr(),
			),
			au.Int(
				"rpcmaxclients",
				"Max number of RPC clients for standard connections",
				constant.DefaultMaxRPCClients,
				cx.Config.RPCMaxClients.Ptr(),
			),
			au.Int(
				"rpcmaxwebsockets",
				"Max number of RPC websocket connections",
				constant.DefaultMaxRPCWebsockets,
				cx.Config.RPCMaxWebsockets.Ptr(),
			),
			au.Int(
				"rpcmaxconcurrentreqs",
				"Max number of RPC requests that may be"+
					" processed concurrently",
				constant.DefaultMaxRPCConcurrentReqs,
				cx.Config.RPCMaxConcurrentReqs.Ptr(),
			),
			au.Bool(
				"rpcquirks",
				"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE:"+
					" Discouraged unless interoperability issues need to be worked"+
					" around",
				cx.Config.RPCQuirks.Ptr(),
			),
			au.Bool(
				"norpc",
				"Disable built-in RPC server -- NOTE: The RPC server"+
					" is disabled by default if no rpcuser/rpcpass or"+
					" rpclimituser/rpclimitpass is specified",
				cx.Config.DisableRPC.Ptr(),
			),
			au.Bool(
				"nodnsseed",
				"Disable DNS seeding for peers",
				cx.Config.DisableDNSSeed.Ptr(),
			),
			au.StringSlice(
				"externalip",
				"Add an ip to the list of local addresses we claim to"+
					" listen on to peers",
				cx.Config.ExternalIPs.Ptr(),
			),
			au.StringSlice(
				"addcheckpoint",
				"Add a custom checkpoint.  Format: '<height>:<hash>'",
				cx.Config.AddCheckpoints.Ptr(),
			),
			au.Bool(
				"nocheckpoints",
				"Disable built-in checkpoints.  Don't do this unless"+
					" you know what you're doing.",
				cx.Config.DisableCheckpoints.Ptr(),
			),
			au.String(
				"dbtype",
				"Database backend to use for the Block Chain",
				constant.DefaultDbType,
				cx.Config.DbType.Ptr(),
			),
			au.String(
				"profile",
				"Enable HTTP profiling on given port -- NOTE port"+
					" must be between 1024 and 65536",
				"",
				cx.Config.Profile.Ptr(),
			),
			au.String(
				"cpuprofile",
				"Write CPU profile to the specified file",
				"",
				cx.Config.CPUProfile.Ptr(),
			),
			au.Bool(
				"upnp",
				"Use UPnP to map our listening port outside of NAT",
				cx.Config.UPNP.Ptr(),
			),
			au.Float64(
				"minrelaytxfee",
				"The minimum transaction fee in DUO/kB to be"+
					" considered a non-zero fee.",
				mempool.DefaultMinRelayTxFee.ToDUO(),
				cx.Config.MinRelayTxFee.Ptr(),
			),
			au.Float64(
				"limitfreerelay",
				"Limit relay of transactions with no transaction"+
					" fee to the given amount in thousands of bytes per minute",
				constant.DefaultFreeTxRelayLimit,
				cx.Config.FreeTxRelayLimit.Ptr(),
			),
			au.Bool(
				"norelaypriority",
				"Do not require free or low-fee transactions to have"+
					" high priority for relaying",
				cx.Config.NoRelayPriority.Ptr(),
			),
			au.Duration(
				"trickleinterval",
				"Minimum time between attempts to send new"+
					" inventory to a connected peer",
				constant.DefaultTrickleInterval,
				cx.Config.TrickleInterval.Ptr(),
			),
			au.Int(
				"maxorphantx",
				"Max number of orphan transactions to keep in memory",
				constant.DefaultMaxOrphanTransactions,
				cx.Config.MaxOrphanTxs.Ptr(),
			),
			au.Bool(
				"generate, g",
				"Generate (mine) DUO using the CPU",
				cx.Config.Generate.Ptr(),
			),
			au.Int(
				"genthreads, G",
				"Number of CPU threads to use with CPU miner"+
					" -1 = all cores",
				1,
				cx.Config.GenThreads.Ptr(),
			),
			au.Bool(
				"solo",
				"mine DUO even if not connected to the network",
				cx.Config.Solo.Ptr(),
			),
			au.Bool(
				"lan",
				"mine duo if not connected to nodes on internet",
				cx.Config.LAN.Ptr(),
			),
			au.Bool(
				"controller",
				"enables multicast",
				cx.Config.Controller.Ptr(),
			),
			au.Bool(
				"autoports",
				"uses random automatic ports for p2p & rpc",
				cx.Config.AutoPorts.Ptr(),
			),
			au.StringSlice(
				"miningaddr",
				"Add the specified payment address to the list of"+
					" addresses to use for generated blocks, at least one is "+
					"required if generate or minerlistener are set",
				cx.Config.MiningAddrs.Ptr(),
			),
			au.String(
				"minerpass",
				"password to authorise sending work to a miner",
				genPassword(),
				cx.Config.MulticastPass.Ptr(),
			),
			au.Int(
				"blockminsize",
				"Minimum block size in bytes to be used when"+
					" creating a block",
				constant.BlockMaxSizeMin,
				cx.Config.BlockMinSize.Ptr(),
			),
			au.Int(
				"blockmaxsize",
				"Maximum block size in bytes to be used when"+
					" creating a block",
				constant.BlockMaxSizeMax,
				cx.Config.BlockMaxSize.Ptr(),
			),
			au.Int(
				"blockminweight",
				"Minimum block weight to be used when creating"+
					" a block",
				constant.BlockMaxWeightMin,
				cx.Config.BlockMinWeight.Ptr(),
			),
			au.Int(
				"blockmaxweight",
				"Maximum block weight to be used when creating"+
					" a block",
				constant.BlockMaxWeightMax,
				cx.Config.BlockMaxWeight.Ptr(),
			),
			au.Int(
				"blockprioritysize",
				"Size in bytes for high-priority/low-fee"+
					" transactions when creating a block",
				mempool.DefaultBlockPrioritySize,
				cx.Config.BlockPrioritySize.Ptr(),
			),
			au.StringSlice(
				"uacomment",
				"Comment to add to the user agent -- See BIP 14 for"+
					" more information.",
				cx.Config.UserAgentComments.Ptr(),
			),
			au.Bool(
				"nopeerbloomfilters",
				"Disable bloom filtering support",
				cx.Config.NoPeerBloomFilters.Ptr(),
			),
			au.Bool(
				"nocfilters",
				"Disable committed filtering (CF) support",
				cx.Config.NoCFilters.Ptr(),
			),
			au.Int(
				"sigcachemaxsize",
				"The maximum number of entries in the"+
					" signature verification cache",
				constant.DefaultSigCacheMaxSize,
				cx.Config.SigCacheMaxSize.Ptr(),
			),
			au.Bool(
				"blocksonly",
				"Do not accept transactions from remote peers.",
				cx.Config.BlocksOnly.Ptr(),
			),
			au.BoolTrue(
				"txindex",
				"Disable the transaction index which makes all transactions available via the getrawtransaction RPC",
				cx.Config.TxIndex.Ptr(),
			),
			au.BoolTrue(
				"addrindex",
				"Disable address-based transaction index which makes the searchrawtransactions RPC available",
				cx.Config.AddrIndex.Ptr(),
			),
			au.Bool(
				"relaynonstd",
				"Relay non-standard transactions regardless of the default settings for the active network.",
				cx.Config.RelayNonStd.Ptr(),
			), au.Bool(
				"rejectnonstd",
				"Reject non-standard transactions regardless of the default settings for the active network.",
				cx.Config.RejectNonStd.Ptr(),
			),
			au.Bool(
				"noinitialload",
				"Defer wallet creation/opening on startup and enable loading wallets over RPC (loading not yet implemented)",
				cx.Config.NoInitialLoad.Ptr(),
			),
			au.Bool(
				"walletconnect, wc",
				"connect to wallet instead of full node",
				cx.Config.Wallet.Ptr(),
			),
			au.String(
				"walletserver, ws",
				"set wallet server to connect to",
				"127.0.0.1:11046",
				cx.Config.WalletServer.Ptr(),
			),
			cli.StringFlag{
				Name:        "walletpass",
				Value:       cx.Config.WalletPass.V(),
				Usage:       "The public wallet password -- Only required if the wallet was created with one",
				EnvVar:      "POD_WALLETPASS",
				Destination: cx.Config.WalletPass.Ptr(),
			},
			au.Bool(
				"onetimetlskey",
				"Generate a new TLS certificate pair at startup, but only write the certificate to disk",
				cx.Config.OneTimeTLSKey.Ptr(),
			),
			au.Bool(
				"tlsskipverify",
				"skip verifying tls certificates",
				cx.Config.TLSSkipVerify.Ptr(),
			),
			au.StringSlice(
				"walletrpclisten",
				"Listen for wallet RPC connections on this"+
					" interface/port (default port: 11046, testnet: 21046,"+
					" simnet: 41046)",
				cx.Config.WalletRPCListeners.Ptr(),
			),
			au.Int(
				"walletrpcmaxclients",
				"Max number of legacy RPC clients for"+
					" standard connections",
				8,
				cx.Config.WalletRPCMaxClients.Ptr(),
			),
			au.Int(
				"walletrpcmaxwebsockets",
				"Max number of legacy RPC websocket connections",
				8,
				cx.Config.WalletRPCMaxWebsockets.Ptr(),
			),
			au.Bool(
				"nodeoff",
				"Starts with node turned off",
				cx.Config.NodeOff.Ptr(),
			),
			au.Bool(
				"walletoff",
				"Starts with wallet turned off",
				cx.Config.WalletOff.Ptr(),
			),
			au.Bool(
				"discover",
				"enable LAN multicast peer discovery in GUI wallet",
				cx.Config.WalletOff.Ptr(),
			),
			au.Bool(
				"delaystart",
				"pauses for 3 seconds before starting, for internal use with restart function",
				nil,
			),
			au.Bool(
				"darktheme",
				"sets the dark theme on the gui interface",
				cx.Config.DarkTheme.Ptr(),
			),
			au.Bool(
				"notty",
				"tells pod there is no keyboard input available",
				nil,
			),
			au.Bool(
				"runasservice",
				"tells wallet to shut down when the wallet locks",
				cx.Config.RunAsService.Ptr(),
			),
		},
	}
}

func genPassword() string {
	s, e := hdkeychain.GenerateSeed(16)
	if e != nil {
		panic("can't do nothing without entropy! " + e.Error())
	}
	return base58.Encode(s)
}
