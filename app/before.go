package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/urfave/cli"

	"github.com/parallelcointeam/parallelcoin/app/apputil"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/config/netparams"
	"github.com/parallelcointeam/parallelcoin/pkg/chain/fork"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
)

func beforeFunc(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log <- cl.Trace{"running beforeFunc", cl.Ine()}
		// if user set datadir this is first thing to configure
		if c.IsSet("datadir") {
			*cx.Config.DataDir = c.String("datadir")
		}
		*cx.Config.ConfigFile =
			*cx.Config.DataDir + string(
				os.PathSeparator) +
				podConfigFilename
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(*cx.Config.ConfigFile) {
			log <- cl.Trace{"loading config", cl.Ine()}
			b, err := ioutil.ReadFile(*cx.Config.ConfigFile)
			if err == nil {
				err = toml.Unmarshal(b, cx.Config)
				if err != nil {
					fmt.Println("error unmarshalling config", err)
					os.Exit(1)
				}
			} else {
				fmt.Println("unexpected error reading configuration file:", err)
				os.Exit(1)
			}
		} // if file didn't exist we save it in Configure after parsing CLI args
		if c.String("loglevel") != "" {
			log <- cl.Trace{"set loglevel", c.String("loglevel"), cl.Ine()}
			*cx.Config.LogLevel = c.String("loglevel")
			cl.Register.SetAllLevels(*cx.Config.LogLevel)
		}
		if c.IsSet("network") {
			log <- cl.Trace{"set network", c.String("network"), cl.Ine()}
			*cx.Config.Network = c.String("network")
			log <- cl.Trace{"network is set to", *cx.Config.Network, cl.Ine()}
			switch *cx.Config.Network {
			case "testnet", "testnet3", "t":
				cx.Log <- cl.Trace{"on testnet", cl.Ine()}
				*cx.Config.TestNet3 = true
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.TestNet3Params
				fork.IsTestnet = true
			case "regtestnet", "regressiontest", "r":
				cx.Log <- cl.Trace{"on regression testnet"}
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = true
				cx.ActiveNet = &netparams.RegressionTestParams
			case "simnet", "s":
				cx.Log <- cl.Trace{"on simnet"}
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = true
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.SimNetParams
			default:
				if *cx.Config.Network != "mainnet" &&
					*cx.Config.Network != "m" {
					cx.Log <- cl.Warn{"using mainnet for node"}
				}
				cx.Log <- cl.Trace{"on mainnet"}
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.MainNetParams
			}
		}
		if c.IsSet("username") {
			log <- cl.Trace{"set username", c.String("username"), cl.Ine()}
			*cx.Config.Username = c.String("username")
		}
		if c.IsSet("password") {
			log <- cl.Trace{"set password", c.String("password"), cl.Ine()}
			*cx.Config.Password = c.String("password")
		}
		if c.IsSet("serveruser") {
			log <- cl.Trace{"set serveruser", c.String("serveruser"), cl.Ine()}
			*cx.Config.ServerUser = c.String("serveruser")
		}
		if c.IsSet("serverpass") {
			log <- cl.Trace{"set serverpass", c.String("serverpass"), cl.Ine()}
			*cx.Config.ServerPass = c.String("serverpass")
		}
		if c.IsSet("limituser") {
			log <- cl.Trace{"set limituser", c.String("limituser"), cl.Ine()}
			*cx.Config.LimitUser = c.String("limituser")
		}
		if c.IsSet("limitpass") {
			log <- cl.Trace{"set limitpass", c.String("limitpass"), cl.Ine()}
			*cx.Config.LimitPass = c.String("limitpass")
		}
		if c.IsSet("rpccert") {
			log <- cl.Trace{"set rpccert", c.String("rpccert"), cl.Ine()}
			*cx.Config.RPCCert = c.String("rpccert")
		}
		if c.IsSet("rpckey") {
			log <- cl.Trace{"set rpckey", c.String("rpckey"), cl.Ine()}
			*cx.Config.RPCKey = c.String("rpckey")
		}
		if c.IsSet("cafile") {
			log <- cl.Trace{"set cafile", c.String("cafile"), cl.Ine()}
			*cx.Config.CAFile = c.String("cafile")
		}
		if c.IsSet("clienttls") {
			log <- cl.Trace{"set clienttls", c.Bool("clienttls"), cl.Ine()}
			*cx.Config.TLS = c.Bool("clienttls")
		}
		if c.IsSet("servertls") {
			log <- cl.Trace{"set servertls", c.Bool("servertls"), cl.Ine()}
			*cx.Config.ServerTLS = c.Bool("servertls")
		}
		if c.IsSet("tlsskipverify") {
			log <- cl.Trace{"set tlsskipverify ", c.Bool("tlsskipverify"),
				cl.Ine()}
			*cx.Config.TLSSkipVerify = c.Bool("tlsskipverify")
		}
		if c.IsSet("proxy") {
			log <- cl.Trace{"set proxy", c.String("proxy"), cl.Ine()}
			*cx.Config.Proxy = c.String("proxy")
		}
		if c.IsSet("proxyuser") {
			log <- cl.Trace{"set proxyuser", c.String("proxyuser"), cl.Ine()}
			*cx.Config.ProxyUser = c.String("proxyuser")
		}
		if c.IsSet("proxypass") {
			log <- cl.Trace{"set proxypass", c.String("proxypass"), cl.Ine()}
			*cx.Config.ProxyPass = c.String("proxypass")
		}
		if c.IsSet("onion") {
			log <- cl.Trace{"set onion", c.Bool("onion"), cl.Ine()}
			*cx.Config.Onion = c.Bool("onion")
		}
		if c.IsSet("onionproxy") {
			log <- cl.Trace{"set onionproxy", c.String("onionproxy"), cl.Ine()}
			*cx.Config.OnionProxy = c.String("onionproxy")
		}
		if c.IsSet("onionuser") {
			log <- cl.Trace{"set onionuser", c.String("onionuser"), cl.Ine()}
			*cx.Config.OnionProxyUser = c.String("onionuser")
		}
		if c.IsSet("onionpass") {
			log <- cl.Trace{"set onionpass", c.String("onionpass"), cl.Ine()}
			*cx.Config.OnionProxyPass = c.String("onionpass")
		}
		if c.IsSet("torisolation") {
			log <- cl.Trace{"set torisolation", c.Bool("torisolation"), cl.Ine()}
			*cx.Config.TorIsolation = c.Bool("torisolation")
		}
		if c.IsSet("addpeer") {
			log <- cl.Trace{"set addpeer", c.StringSlice("addpeer"), cl.Ine()}
			*cx.Config.AddPeers = c.StringSlice("addpeer")
		}
		if c.IsSet("connect") {
			log <- cl.Trace{"set connect", c.StringSlice("connect"), cl.Ine()}
			*cx.Config.ConnectPeers = c.StringSlice("connect")
		}
		if c.IsSet("nolisten") {
			log <- cl.Trace{"set nolisten", c.Bool("nolisten"), cl.Ine()}
			*cx.Config.DisableListen = c.Bool("nolisten")
		}
		if c.IsSet("listen") {
			log <- cl.Trace{"set listen", c.StringSlice("listen"), cl.Ine()}
			*cx.Config.Listeners = c.StringSlice("listen")
		}
		if c.IsSet("maxpeers") {
			log <- cl.Trace{"set maxpeers", c.Int("maxpeers"), cl.Ine()}
			*cx.Config.MaxPeers = c.Int("maxpeers")
		}
		if c.IsSet("nobanning") {
			log <- cl.Trace{"set nobanning", c.Bool("nobanning"), cl.Ine()}
			*cx.Config.DisableBanning = c.Bool("nobanning")
		}
		if c.IsSet("banduration") {
			log <- cl.Trace{"set banduration", c.Duration("banduration"), cl.Ine()}
			*cx.Config.BanDuration = c.Duration("banduration")
		}
		if c.IsSet("banthreshold") {
			log <- cl.Trace{"set banthreshold", c.Int("banthreshold"), cl.Ine()}
			*cx.Config.BanThreshold = c.Int("banthreshold")
		}
		if c.IsSet("whitelist") {
			log <- cl.Trace{"set whitelist", c.StringSlice("whitelist"), cl.Ine()}
			*cx.Config.Whitelists = c.StringSlice("whitelist")
		}
		if c.IsSet("rpcconnect") {
			log <- cl.Trace{"set rpcconnect", c.String("rpcconnect"), cl.Ine()}
			*cx.Config.RPCConnect = c.String("rpcconnect")
		}
		if c.IsSet("rpclisten") {
			log <- cl.Trace{"set rpclisten", c.StringSlice("rpclisten"), cl.Ine()}
			*cx.Config.RPCListeners = c.StringSlice("rpclisten")
		}
		if c.IsSet("rpcmaxclients") {
			log <- cl.Trace{"set rpcmaxclients", c.Int("rpcmaxclients"), cl.Ine()}
			*cx.Config.RPCMaxClients = c.Int("rpcmaxclients")
		}
		if c.IsSet("rpcmaxwebsockets") {
			log <- cl.Trace{"set rpcmaxwebsockets", c.Int("rpcmaxwebsockets"), cl.Ine()}
			*cx.Config.RPCMaxWebsockets = c.Int("rpcmaxwebsockets")
		}
		if c.IsSet("rpcmaxconcurrentreqs") {
			log <- cl.Trace{"set rpcmaxconcurrentreqs", c.Int("rpcmaxconcurrentreqs"), cl.Ine()}
			*cx.Config.RPCMaxConcurrentReqs = c.Int("rpcmaxconcurrentreqs")
		}
		if c.IsSet("rpcquirks") {
			log <- cl.Trace{"set rpcquirks", c.Bool("rpcquirks"), cl.Ine()}
			*cx.Config.RPCQuirks = c.Bool("rpcquirks")
		}
		if c.IsSet("norpc") {
			log <- cl.Trace{"set norpc", c.Bool("norpc"), cl.Ine()}
			*cx.Config.DisableRPC = c.Bool("norpc")
		}
		if c.IsSet("nodnsseed") {
			log <- cl.Trace{"set nodnsseed", c.Bool("nodnsseed"), cl.Ine()}
			*cx.Config.DisableDNSSeed = c.Bool("nodnsseed")
		}
		if c.IsSet("externalip") {
			log <- cl.Trace{"set externalip", c.StringSlice("externalip"), cl.Ine()}
			*cx.Config.ExternalIPs = c.StringSlice("externalip")
		}
		if c.IsSet("addcheckpoint") {
			log <- cl.Trace{"set addcheckpoint", c.StringSlice("addcheckpoint"), cl.Ine()}
			*cx.Config.AddCheckpoints = c.StringSlice("addcheckpoint")
		}
		if c.IsSet("nocheckpoints") {
			log <- cl.Trace{"set nocheckpoints", c.Bool("nocheckpoints"), cl.Ine()}
			*cx.Config.DisableCheckpoints = c.Bool("nocheckpoints")
		}
		if c.IsSet("dbtype") {
			log <- cl.Trace{"set dbtype", c.String("dbtype"), cl.Ine()}
			*cx.Config.DbType = c.String("dbtype")
		}
		if c.IsSet("profile") {
			log <- cl.Trace{"set profile", c.String("profile"), cl.Ine()}
			*cx.Config.Profile = c.String("profile")
		}
		if c.IsSet("cpuprofile") {
			log <- cl.Trace{"set cpuprofile", c.String("cpuprofile"), cl.Ine()}
			*cx.Config.CPUProfile = c.String("cpuprofile")
		}
		if c.IsSet("upnp") {
			log <- cl.Trace{"set upnp", c.Bool("upnp"), cl.Ine()}
			*cx.Config.Upnp = c.Bool("upnp")
		}
		if c.IsSet("minrelaytxfee") {
			log <- cl.Trace{"set minrelaytxfee", c.Float64("minrelaytxfee"), cl.Ine()}
			*cx.Config.MinRelayTxFee = c.Float64("minrelaytxfee")
		}
		if c.IsSet("limitfreerelay") {
			log <- cl.Trace{"set limitfreerelay", c.Float64("limitfreerelay"), cl.Ine()}
			*cx.Config.FreeTxRelayLimit = c.Float64("limitfreerelay")
		}
		if c.IsSet("norelaypriority") {
			log <- cl.Trace{"set norelaypriority", c.Bool("norelaypriority"), cl.Ine()}
			*cx.Config.NoRelayPriority = c.Bool("norelaypriority")
		}
		if c.IsSet("trickleinterval") {
			log <- cl.Trace{"set trickleinterval", c.Duration("trickleinterval"), cl.Ine()}
			*cx.Config.TrickleInterval = c.Duration("trickleinterval")
		}
		if c.IsSet("maxorphantx") {
			log <- cl.Trace{"set maxorphantx", c.Int("maxorphantx"), cl.Ine()}
			*cx.Config.MaxOrphanTxs = c.Int("maxorphantx")
		}
		if c.IsSet("algo") {
			log <- cl.Trace{"set algo", c.String("algo"), cl.Ine()}
			*cx.Config.Algo = c.String("algo")
		}
		if c.IsSet("generate") {
			log <- cl.Trace{"set generate", c.Bool("generate"), cl.Ine()}
			*cx.Config.Generate = c.Bool("generate")
		}
		if c.IsSet("genthreads") {
			log <- cl.Trace{"set genthreads", c.Int("genthreads"), cl.Ine()}
			*cx.Config.GenThreads = c.Int("genthreads")
		}
		if c.IsSet("nocontroller") {
			log <- cl.Trace{"set nocontroller",
				c.String("nocontroller"), cl.Ine()}
			*cx.Config.NoController = c.Bool("nocontroller")
		}
		if c.IsSet("miningaddr") {
			log <- cl.Trace{"set miningaddr", c.StringSlice("miningaddr"), cl.Ine()}
			*cx.Config.MiningAddrs = c.StringSlice("miningaddr")
		}
		if c.IsSet("minerpass") {
			log <- cl.Trace{"set minerpass", c.String("minerpass"), cl.Ine()}
			*cx.Config.MinerPass = c.String("minerpass")
		}
		if c.IsSet("group") {
			log <- cl.Trace{"set group", c.String("group"), cl.Ine()}
			*cx.Config.Group = c.String("group")
		}
		if c.IsSet("nodiscovery") {
			log <- cl.Trace{"set nodiscovery",
				c.String("nodiscovery"), cl.Ine()}
			*cx.Config.NoDiscovery = c.Bool("nodiscovery")
		}
		if c.IsSet("blockminsize") {
			log <- cl.Trace{"set blockminsize", c.Int("blockminsize"), cl.Ine()}
			*cx.Config.BlockMinSize = c.Int("blockminsize")
		}
		if c.IsSet("blockmaxsize") {
			log <- cl.Trace{"set blockmaxsize", c.Int("blockmaxsize"), cl.Ine()}
			*cx.Config.BlockMaxSize = c.Int("blockmaxsize")
		}
		if c.IsSet("blockminweight") {
			log <- cl.Trace{"set blockminweight", c.Int("blockminweight"), cl.Ine()}
			*cx.Config.BlockMinWeight = c.Int("blockminweight")
		}
		if c.IsSet("blockmaxweight") {
			log <- cl.Trace{"set blockmaxweight", c.Int("blockmaxweight"), cl.Ine()}
			*cx.Config.BlockMaxWeight = c.Int("blockmaxweight")
		}
		if c.IsSet("blockprioritysize") {
			log <- cl.Trace{"set blockprioritysize", c.Int("blockprioritysize"), cl.Ine()}
			*cx.Config.BlockPrioritySize = c.Int("blockprioritysize")
		}
		if c.IsSet("uacomment") {
			log <- cl.Trace{"set uacomment", c.StringSlice("uacomment"), cl.Ine()}
			*cx.Config.UserAgentComments = c.StringSlice("uacomment")
		}
		if c.IsSet("nopeerbloomfilters") {
			log <- cl.Trace{"set nopeerbloomfilters", c.Bool("nopeerbloomfilters"), cl.Ine()}
			*cx.Config.NoPeerBloomFilters = c.Bool("nopeerbloomfilters")
		}
		if c.IsSet("nocfilters") {
			log <- cl.Trace{"set nocfilters", c.Bool("nocfilters"), cl.Ine()}
			*cx.Config.NoCFilters = c.Bool("nocfilters")
		}
		if c.IsSet("sigcachemaxsize") {
			log <- cl.Trace{"set sigcachemaxsize", c.Int("sigcachemaxsize"), cl.Ine()}
			*cx.Config.SigCacheMaxSize = c.Int("sigcachemaxsize")
		}
		if c.IsSet("blocksonly") {
			log <- cl.Trace{"set blocksonly", c.Bool("blocksonly"), cl.Ine()}
			*cx.Config.BlocksOnly = c.Bool("blocksonly")
		}
		if c.IsSet("notxindex") {
			log <- cl.Trace{"set notxindex", c.Bool("notxindex"), cl.Ine()}
			*cx.Config.TxIndex = c.Bool("notxindex")
		}
		if c.IsSet("noaddrindex") {
			log <- cl.Trace{"set noaddrindex", c.Bool("noaddrindex"), cl.Ine()}
			*cx.Config.AddrIndex = c.Bool("noaddrindex")
		}
		if c.IsSet("relaynonstd") {
			log <- cl.Trace{"set relaynonstd", c.Bool("relaynonstd"), cl.Ine()}
			*cx.Config.RelayNonStd = c.Bool("relaynonstd")
		}
		if c.IsSet("rejectnonstd") {
			log <- cl.Trace{"set rejectnonstd", c.Bool("rejectnonstd"), cl.Ine()}
			*cx.Config.RejectNonStd = c.Bool("rejectnonstd")
		}
		if c.IsSet("noinitialload") {
			log <- cl.Trace{"set noinitialload", c.Bool("noinitialload"), cl.Ine()}
			*cx.Config.NoInitialLoad = c.Bool("noinitialload")
		}
		if c.IsSet("walletconnect") {
			log <- cl.Trace{"set walletconnect", c.Bool("walletconnect"), cl.Ine()}
			*cx.Config.Wallet = c.Bool("walletconnect")
		}
		if c.IsSet("walletserver") {
			log <- cl.Trace{"set walletserver", c.String("walletserver"), cl.Ine()}
			*cx.Config.WalletServer = c.String("walletserver")
		}
		if c.IsSet("walletpass") {
			log <- cl.Trace{"set walletpass", c.String("walletpass"), cl.Ine()}
			*cx.Config.WalletPass = c.String("walletpass")
		}
		if c.IsSet("onetimetlskey") {
			log <- cl.Trace{"set onetimetlskey", c.Bool("onetimetlskey"), cl.Ine()}
			*cx.Config.OneTimeTLSKey = c.Bool("onetimetlskey")
		}
		if c.IsSet("walletrpclisten") {
			log <- cl.Trace{"set walletrpclisten", c.StringSlice("walletrpclisten"), cl.Ine()}
			*cx.Config.WalletRPCListeners = c.StringSlice("walletrpclisten")
		}
		if c.IsSet("walletrpcmaxclients") {
			log <- cl.Trace{"set walletrpcmaxclients", c.Int("walletrpcmaxclients"), cl.Ine()}
			*cx.Config.WalletRPCMaxClients = c.Int("walletrpcmaxclients")
		}
		if c.IsSet("walletrpcmaxwebsockets") {
			log <- cl.Trace{"set walletrpcmaxwebsockets", c.Int("walletrpcmaxwebsockets"), cl.Ine()}
			*cx.Config.WalletRPCMaxWebsockets = c.Int("walletrpcmaxwebsockets")
		}
		if c.IsSet("experimentalrpclisten") {
			log <- cl.Trace{"set experimentalrpclisten", c.StringSlice("experimentalrpclisten"), cl.Ine()}
			*cx.Config.ExperimentalRPCListeners = c.StringSlice("experimentalrpclisten")
		}
		if c.IsSet("nodeoff") {
			log <- cl.Trace{"set nodeoff", c.Bool("nodeoff"), cl.Ine()}
			*cx.Config.NodeOff = c.Bool("nodeoff")
		}
		if c.IsSet("testnodeoff") {
			log <- cl.Trace{"set testnodeoff", c.Bool("testnodeoff"), cl.Ine()}
			*cx.Config.TestNodeOff = c.Bool("testnodeoff")
		}
		if c.IsSet("walletoff") {
			log <- cl.Trace{"set walletoff", c.Bool("walletoff"), cl.Ine()}
			*cx.Config.WalletOff = c.Bool("walletoff")
		}
		if c.IsSet("save") {
			log <- cl.Trace{"set save", c.Bool("save"), cl.Ine()}
			cx.StateCfg.Save = true
		}
		return nil
	}
}
