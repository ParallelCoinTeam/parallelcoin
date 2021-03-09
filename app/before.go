package app

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	"github.com/p9c/pod/pkg/util/routeable"
	"io/ioutil"
	prand "math/rand"
	"os"
	"time"
	
	"github.com/p9c/pod/cmd/spv"
	"github.com/p9c/pod/pkg/pipe/serve"
	"github.com/p9c/pod/version"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/blockchain/chaincfg"
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	"github.com/p9c/pod/pkg/blockchain/fork"
	"github.com/p9c/pod/pkg/pod"
)

func beforeFunc(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		dbg.Ln("running beforeFunc")
		cx.AppContext = c
		// if user set datadir this is first thing to configure
		if c.IsSet("datadir") {
			*cx.Config.DataDir = c.String("datadir")
			cx.DataDir = c.String("datadir")
			dbg.Ln("datadir", *cx.Config.DataDir)
		}
		dbg.Ln(c.IsSet("D"), c.IsSet("datadir"))
		// // propagate datadir path to interrupt for restart handling
		// interrupt.DataDir = cx.DataDir
		// if there is a delaystart requested, pause for 3 seconds
		if c.IsSet("delaystart") {
			time.Sleep(time.Second * 3)
		}
		if c.IsSet("pipelog") {
			dbg.Ln("pipe logger enabled")
			*cx.Config.PipeLog = c.Bool("pipelog")
			serve.Log(cx.KillAll, fmt.Sprint(os.Args))
		}
		if c.IsSet("walletfile") {
			*cx.Config.WalletFile = c.String("walletfile")
		}
		*cx.Config.ConfigFile =
			*cx.Config.DataDir + string(os.PathSeparator) + podConfigFilename
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(*cx.Config.ConfigFile) {
			b, e := ioutil.ReadFile(*cx.Config.ConfigFile)
			if e == nil {
				cx.Config, cx.ConfigMap = pod.EmptyConfig()
				e = json.Unmarshal(b, cx.Config)
				if e != nil {
					err.Ln("error unmarshalling config", e)
					// os.Exit(1)
					return e
				}
			} else {
				ftl.Ln("unexpected error reading configuration file:", e)
				// os.Exit(1)
				return e
			}
		} else {
			*cx.Config.ConfigFile = ""
			dbg.Ln("will save config after configuration")
			cx.StateCfg.Save = true
		}
		if c.IsSet("loglevel") {
			trc.Ln("set loglevel", c.String("loglevel"))
			*cx.Config.LogLevel = c.String("loglevel")
		}
		logg.SetLogLevel(*cx.Config.LogLevel)
		if !*cx.Config.PipeLog {
			// if/when running further instances of the same version no reason
			// to print the version message again
			dbg.Ln("\nrunning", os.Args, version.Get())
		}
		if c.IsSet("network") {
			*cx.Config.Network = c.String("network")
			switch *cx.Config.Network {
			case "testnet", "testnet3", "t":
				cx.ActiveNet = &netparams.TestNet3Params
				fork.IsTestnet = true
				// fork.HashReps = 3
			case "regtestnet", "regressiontest", "r":
				fork.IsTestnet = true
				cx.ActiveNet = &netparams.RegressionTestParams
			case "simnet", "s":
				fork.IsTestnet = true
				cx.ActiveNet = &netparams.SimNetParams
			default:
				if *cx.Config.Network != "mainnet" &&
					*cx.Config.Network != "m" {
					dbg.Ln("using mainnet for node")
				}
				cx.ActiveNet = &netparams.MainNetParams
			}
		}
		if c.IsSet("username") {
			*cx.Config.Username = c.String("username")
		}
		if c.IsSet("password") {
			*cx.Config.Password = c.String("password")
		}
		if c.IsSet("serveruser") {
			*cx.Config.ServerUser = c.String("serveruser")
		}
		if c.IsSet("serverpass") {
			*cx.Config.ServerPass = c.String("serverpass")
		}
		if c.IsSet("limituser") {
			*cx.Config.LimitUser = c.String("limituser")
		}
		if c.IsSet("limitpass") {
			*cx.Config.LimitPass = c.String("limitpass")
		}
		if c.IsSet("rpccert") {
			*cx.Config.RPCCert = c.String("rpccert")
		}
		if c.IsSet("rpckey") {
			*cx.Config.RPCKey = c.String("rpckey")
		}
		if c.IsSet("cafile") {
			*cx.Config.CAFile = c.String("cafile")
		}
		if c.IsSet("clienttls") {
			*cx.Config.TLS = c.Bool("clienttls")
		}
		if c.IsSet("servertls") {
			*cx.Config.ServerTLS = c.Bool("servertls")
		}
		if c.IsSet("tlsskipverify") {
			*cx.Config.TLSSkipVerify = c.Bool("tlsskipverify")
		}
		if c.IsSet("proxy") {
			*cx.Config.Proxy = c.String("proxy")
		}
		if c.IsSet("proxyuser") {
			*cx.Config.ProxyUser = c.String("proxyuser")
		}
		if c.IsSet("proxypass") {
			*cx.Config.ProxyPass = c.String("proxypass")
		}
		if c.IsSet("onion") {
			*cx.Config.Onion = c.Bool("onion")
		}
		if c.IsSet("onionproxy") {
			*cx.Config.OnionProxy = c.String("onionproxy")
		}
		if c.IsSet("onionuser") {
			*cx.Config.OnionProxyUser = c.String("onionuser")
		}
		if c.IsSet("onionpass") {
			*cx.Config.OnionProxyPass = c.String("onionpass")
		}
		if c.IsSet("torisolation") {
			*cx.Config.TorIsolation = c.Bool("torisolation")
		}
		if c.IsSet("addpeer") {
			*cx.Config.AddPeers = c.StringSlice("addpeer")
		}
		if c.IsSet("connect") {
			*cx.Config.ConnectPeers = c.StringSlice("connect")
		}
		if c.IsSet("nolisten") {
			*cx.Config.DisableListen = c.Bool("nolisten")
		}
		if c.IsSet("listen") {
			*cx.Config.P2PListeners = c.StringSlice("listen")
		}
		if c.IsSet("maxpeers") {
			*cx.Config.MaxPeers = c.Int("maxpeers")
		}
		if c.IsSet("nobanning") {
			*cx.Config.DisableBanning = c.Bool("nobanning")
		}
		if c.IsSet("banduration") {
			*cx.Config.BanDuration = c.Duration("banduration")
		}
		if c.IsSet("banthreshold") {
			*cx.Config.BanThreshold = c.Int("banthreshold")
		}
		if c.IsSet("whitelist") {
			*cx.Config.Whitelists = c.StringSlice("whitelist")
		}
		if c.IsSet("rpcconnect") {
			*cx.Config.RPCConnect = c.String("rpcconnect")
		}
		if c.IsSet("rpclisten") {
			*cx.Config.RPCListeners = c.StringSlice("rpclisten")
		}
		if c.IsSet("rpcmaxclients") {
			*cx.Config.RPCMaxClients = c.Int("rpcmaxclients")
		}
		if c.IsSet("rpcmaxwebsockets") {
			*cx.Config.RPCMaxWebsockets = c.Int("rpcmaxwebsockets")
		}
		if c.IsSet("rpcmaxconcurrentreqs") {
			*cx.Config.RPCMaxConcurrentReqs = c.Int("rpcmaxconcurrentreqs")
		}
		if c.IsSet("rpcquirks") {
			*cx.Config.RPCQuirks = c.Bool("rpcquirks")
		}
		if c.IsSet("norpc") {
			*cx.Config.DisableRPC = c.Bool("norpc")
		}
		if c.IsSet("nodnsseed") {
			*cx.Config.DisableDNSSeed = c.Bool("nodnsseed")
			spv.DisableDNSSeed = c.Bool("nodnsseed")
		}
		if c.IsSet("externalip") {
			*cx.Config.ExternalIPs = c.StringSlice("externalip")
		}
		if c.IsSet("addcheckpoint") {
			*cx.Config.AddCheckpoints = c.StringSlice("addcheckpoint")
		}
		if c.IsSet("nocheckpoints") {
			*cx.Config.DisableCheckpoints = c.Bool("nocheckpoints")
		}
		if c.IsSet("dbtype") {
			*cx.Config.DbType = c.String("dbtype")
		}
		if c.IsSet("profile") {
			*cx.Config.Profile = c.String("profile")
		}
		if c.IsSet("cpuprofile") {
			*cx.Config.CPUProfile = c.String("cpuprofile")
		}
		if c.IsSet("upnp") {
			*cx.Config.UPNP = c.Bool("upnp")
		}
		if c.IsSet("minrelaytxfee") {
			*cx.Config.MinRelayTxFee = c.Float64("minrelaytxfee")
		}
		if c.IsSet("limitfreerelay") {
			*cx.Config.FreeTxRelayLimit = c.Float64("limitfreerelay")
		}
		if c.IsSet("norelaypriority") {
			*cx.Config.NoRelayPriority = c.Bool("norelaypriority")
		}
		if c.IsSet("trickleinterval") {
			*cx.Config.TrickleInterval = c.Duration("trickleinterval")
		}
		if c.IsSet("maxorphantx") {
			*cx.Config.MaxOrphanTxs = c.Int("maxorphantx")
		}
		if c.IsSet("generate") {
			*cx.Config.Generate = c.Bool("generate")
		}
		if c.IsSet("genthreads") {
			*cx.Config.GenThreads = c.Int("genthreads")
		}
		if c.IsSet("solo") {
			*cx.Config.Solo = c.Bool("solo")
		}
		if c.IsSet("autoports") {
			*cx.Config.AutoPorts = c.Bool("autoports")
		}
		if c.IsSet("lan") {
			// if LAN is turned on we need to remove the seeds from netparams not on mainnet
			// mainnet is never in lan mode
			// if LAN is turned on it means by default we are on testnet
			cx.ActiveNet = &netparams.TestNet3Params
			if cx.ActiveNet.Name != "mainnet" {
				dbg.Ln("set lan", c.Bool("lan"))
				*cx.Config.LAN = c.Bool("lan")
				cx.ActiveNet.DNSSeeds = []chaincfg.DNSSeed{}
			} else {
				*cx.Config.LAN = false
			}
		}
		// if c.IsSet("controller") {
		// 	*cx.Config.Controller = c.String("controller")
		// }
		// if c.IsSet("controllerconnect") {
		// 	*cx.Config.ControllerConnect = c.StringSlice("controllerconnect")
		// }
		if c.IsSet("miningaddrs") {
			*cx.Config.MiningAddrs = c.StringSlice("miningaddrs")
		}
		if c.IsSet("minerpass") {
			*cx.Config.MinerPass = c.String("minerpass")
			dbg.Ln("--------- set minerpass", *cx.Config.MinerPass)
			cx.StateCfg.Save = true
		}
		if c.IsSet("blockminsize") {
			*cx.Config.BlockMinSize = c.Int("blockminsize")
		}
		if c.IsSet("blockmaxsize") {
			*cx.Config.BlockMaxSize = c.Int("blockmaxsize")
		}
		if c.IsSet("blockminweight") {
			*cx.Config.BlockMinWeight = c.Int("blockminweight")
		}
		if c.IsSet("blockmaxweight") {
			*cx.Config.BlockMaxWeight = c.Int("blockmaxweight")
		}
		if c.IsSet("blockprioritysize") {
			*cx.Config.BlockPrioritySize = c.Int("blockprioritysize")
		}
		prand.Seed(time.Now().UnixNano())
		nonce := fmt.Sprintf("nonce%0x", prand.Uint32())
		if cx.Config.UserAgentComments == nil {
			cx.Config.UserAgentComments = &cli.StringSlice{nonce}
		} else {
			*cx.Config.UserAgentComments = append(cli.StringSlice{nonce}, *cx.Config.UserAgentComments...)
		}
		if c.IsSet("uacomment") {
			*cx.Config.UserAgentComments = append(
				*cx.Config.UserAgentComments,
				c.StringSlice("uacomment")...,
			)
		}
		if c.IsSet("nopeerbloomfilters") {
			*cx.Config.NoPeerBloomFilters = c.Bool("nopeerbloomfilters")
		}
		if c.IsSet("nocfilters") {
			*cx.Config.NoCFilters = c.Bool("nocfilters")
		}
		if c.IsSet("sigcachemaxsize") {
			*cx.Config.SigCacheMaxSize = c.Int("sigcachemaxsize")
		}
		if c.IsSet("blocksonly") {
			*cx.Config.BlocksOnly = c.Bool("blocksonly")
		}
		if c.IsSet("notxindex") {
			*cx.Config.TxIndex = c.Bool("notxindex")
		}
		if c.IsSet("noaddrindex") {
			*cx.Config.AddrIndex = c.Bool("noaddrindex")
		}
		if c.IsSet("relaynonstd") {
			*cx.Config.RelayNonStd = c.Bool("relaynonstd")
		}
		if c.IsSet("rejectnonstd") {
			*cx.Config.RejectNonStd = c.Bool("rejectnonstd")
		}
		if c.IsSet("noinitialload") {
			*cx.Config.NoInitialLoad = c.Bool("noinitialload")
		}
		if c.IsSet("walletconnect") {
			*cx.Config.Wallet = c.Bool("walletconnect")
		}
		if c.IsSet("walletserver") {
			*cx.Config.WalletServer = c.String("walletserver")
		}
		if c.IsSet("walletpass") {
			*cx.Config.WalletPass = c.String("walletpass")
		} else {
			// if this is not set, the config will be storing the hash and hashes on save, so we set explicitly to empty
			// as otherwise it would have the hex of the hash of the password here
			*cx.Config.WalletPass = ""
		}
		if c.IsSet("onetimetlskey") {
			*cx.Config.OneTimeTLSKey = c.Bool("onetimetlskey")
		}
		if c.IsSet("walletrpclisten") {
			*cx.Config.WalletRPCListeners = c.StringSlice("walletrpclisten")
		}
		if c.IsSet("walletrpcmaxclients") {
			*cx.Config.WalletRPCMaxClients = c.Int("walletrpcmaxclients")
		}
		if c.IsSet("walletrpcmaxwebsockets") {
			*cx.Config.WalletRPCMaxWebsockets = c.Int("walletrpcmaxwebsockets")
		}
		if c.IsSet("nodeoff") {
			*cx.Config.NodeOff = c.Bool("nodeoff")
		}
		if c.IsSet("walletoff") {
			*cx.Config.WalletOff = c.Bool("walletoff")
		}
		if c.IsSet("darktheme") {
			*cx.Config.DarkTheme = c.Bool("darktheme")
		}
		if c.IsSet("notty") {
			cx.IsGUI = true
		}
		if c.IsSet("disablecontroller") {
			*cx.Config.DisableController = c.Bool("disablecontroller")
		}
		if c.IsSet("save") {
			inf.Ln("saving configuration")
			cx.StateCfg.Save = true
		}
		// if e = routeable.Discover(); err.Chk(e) {
		// 	// TODO: this should trigger the display of this lack of internet
		// }
		go func() {
		out:
			for {
				select {
				case <-time.After(time.Second * 10):
					if e = routeable.Discover(); err.Chk(e) {
						// TODO: this should trigger the display of this lack of internet
					}
				case <-cx.KillAll:
					break out
				}
			}
		}()
		return nil
	}
}
