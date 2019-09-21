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
		L.Trace("running beforeFunc")
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
			L.Trace("loading config")
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
			L.Trace("set loglevel", c.String("loglevel"))
			*cx.Config.LogLevel = c.String("loglevel")
			cl.Register.SetAllLevels(*cx.Config.LogLevel)
		}
		if c.IsSet("network") {
			L.Trace("set network", c.String("network"))
			*cx.Config.Network = c.String("network")
			L.Trace("network is set to", *cx.Config.Network)
			switch *cx.Config.Network {
			case "testnet", "testnet3", "t":
				L.Trace("on testnet")
				*cx.Config.TestNet3 = true
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.TestNet3Params
				fork.IsTestnet = true
			case "regtestnet", "regressiontest", "r":
				L.Trace("on regression testnet")
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = true
				cx.ActiveNet = &netparams.RegressionTestParams
			case "simnet", "s":
				L.Trace("on simnet")
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = true
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.SimNetParams
			default:
				if *cx.Config.Network != "mainnet" &&
					*cx.Config.Network != "m" {
					cx.Log <- cl.Warn{"using mainnet for node"}
				}
				L.Trace("on mainnet")
				*cx.Config.TestNet3 = false
				*cx.Config.SimNet = false
				*cx.Config.RegressionTest = false
				cx.ActiveNet = &netparams.MainNetParams
			}
		}
		if c.IsSet("username") {
			L.Trace("set username", c.String("username"))
			*cx.Config.Username = c.String("username")
		}
		if c.IsSet("password") {
			L.Trace("set password", c.String("password"))
			*cx.Config.Password = c.String("password")
		}
		if c.IsSet("serveruser") {
			L.Trace("set serveruser", c.String("serveruser"))
			*cx.Config.ServerUser = c.String("serveruser")
		}
		if c.IsSet("serverpass") {
			L.Trace("set serverpass", c.String("serverpass"))
			*cx.Config.ServerPass = c.String("serverpass")
		}
		if c.IsSet("limituser") {
			L.Trace("set limituser", c.String("limituser"))
			*cx.Config.LimitUser = c.String("limituser")
		}
		if c.IsSet("limitpass") {
			L.Trace("set limitpass", c.String("limitpass"))
			*cx.Config.LimitPass = c.String("limitpass")
		}
		if c.IsSet("rpccert") {
			L.Trace("set rpccert", c.String("rpccert"))
			*cx.Config.RPCCert = c.String("rpccert")
		}
		if c.IsSet("rpckey") {
			L.Trace("set rpckey", c.String("rpckey"))
			*cx.Config.RPCKey = c.String("rpckey")
		}
		if c.IsSet("cafile") {
			L.Trace("set cafile", c.String("cafile"))
			*cx.Config.CAFile = c.String("cafile")
		}
		if c.IsSet("clienttls") {
			L.Trace("set clienttls", c.Bool("clienttls"))
			*cx.Config.TLS = c.Bool("clienttls")
		}
		if c.IsSet("servertls") {
			L.Trace("set servertls", c.Bool("servertls"))
			*cx.Config.ServerTLS = c.Bool("servertls")
		}
		if c.IsSet("tlsskipverify") {
			L.Trace("set tlsskipverify ", c.Bool("tlsskipverify"))
			*cx.Config.TLSSkipVerify = c.Bool("tlsskipverify")
		}
		if c.IsSet("proxy") {
			L.Trace("set proxy", c.String("proxy"))
			*cx.Config.Proxy = c.String("proxy")
		}
		if c.IsSet("proxyuser") {
			L.Trace("set proxyuser", c.String("proxyuser"))
			*cx.Config.ProxyUser = c.String("proxyuser")
		}
		if c.IsSet("proxypass") {
			L.Trace("set proxypass", c.String("proxypass"))
			*cx.Config.ProxyPass = c.String("proxypass")
		}
		if c.IsSet("onion") {
			L.Trace("set onion", c.Bool("onion"))
			*cx.Config.Onion = c.Bool("onion")
		}
		if c.IsSet("onionproxy") {
			L.Trace("set onionproxy", c.String("onionproxy"))
			*cx.Config.OnionProxy = c.String("onionproxy")
		}
		if c.IsSet("onionuser") {
			L.Trace("set onionuser", c.String("onionuser"))
			*cx.Config.OnionProxyUser = c.String("onionuser")
		}
		if c.IsSet("onionpass") {
			L.Trace("set onionpass", c.String("onionpass"))
			*cx.Config.OnionProxyPass = c.String("onionpass")
		}
		if c.IsSet("torisolation") {
			L.Trace("set torisolation", c.Bool("torisolation"))
			*cx.Config.TorIsolation = c.Bool("torisolation")
		}
		if c.IsSet("addpeer") {
			L.Trace("set addpeer", c.StringSlice("addpeer"))
			*cx.Config.AddPeers = c.StringSlice("addpeer")
		}
		if c.IsSet("connect") {
			L.Trace("set connect", c.StringSlice("connect"))
			*cx.Config.ConnectPeers = c.StringSlice("connect")
		}
		if c.IsSet("nolisten") {
			L.Trace("set nolisten", c.Bool("nolisten"))
			*cx.Config.DisableListen = c.Bool("nolisten")
		}
		if c.IsSet("listen") {
			L.Trace("set listen", c.StringSlice("listen"))
			*cx.Config.Listeners = c.StringSlice("listen")
		}
		if c.IsSet("maxpeers") {
			L.Trace("set maxpeers", c.Int("maxpeers"))
			*cx.Config.MaxPeers = c.Int("maxpeers")
		}
		if c.IsSet("nobanning") {
			L.Trace("set nobanning", c.Bool("nobanning"))
			*cx.Config.DisableBanning = c.Bool("nobanning")
		}
		if c.IsSet("banduration") {
			L.Trace("set banduration", c.Duration("banduration"))
			*cx.Config.BanDuration = c.Duration("banduration")
		}
		if c.IsSet("banthreshold") {
			L.Trace("set banthreshold", c.Int("banthreshold"))
			*cx.Config.BanThreshold = c.Int("banthreshold")
		}
		if c.IsSet("whitelist") {
			L.Trace("set whitelist", c.StringSlice("whitelist"))
			*cx.Config.Whitelists = c.StringSlice("whitelist")
		}
		if c.IsSet("rpcconnect") {
			L.Trace("set rpcconnect", c.String("rpcconnect"))
			*cx.Config.RPCConnect = c.String("rpcconnect")
		}
		if c.IsSet("rpclisten") {
			L.Trace("set rpclisten", c.StringSlice("rpclisten"))
			*cx.Config.RPCListeners = c.StringSlice("rpclisten")
		}
		if c.IsSet("rpcmaxclients") {
			L.Trace("set rpcmaxclients", c.Int("rpcmaxclients"))
			*cx.Config.RPCMaxClients = c.Int("rpcmaxclients")
		}
		if c.IsSet("rpcmaxwebsockets") {
			L.Trace("set rpcmaxwebsockets", c.Int("rpcmaxwebsockets"))
			*cx.Config.RPCMaxWebsockets = c.Int("rpcmaxwebsockets")
		}
		if c.IsSet("rpcmaxconcurrentreqs") {
			L.Trace("set rpcmaxconcurrentreqs",
				c.Int("rpcmaxconcurrentreqs"))
			*cx.Config.RPCMaxConcurrentReqs = c.Int("rpcmaxconcurrentreqs")
		}
		if c.IsSet("rpcquirks") {
			L.Trace("set rpcquirks", c.Bool("rpcquirks"))
			*cx.Config.RPCQuirks = c.Bool("rpcquirks")
		}
		if c.IsSet("norpc") {
			L.Trace("set norpc", c.Bool("norpc"))
			*cx.Config.DisableRPC = c.Bool("norpc")
		}
		if c.IsSet("nodnsseed") {
			L.Trace("set nodnsseed", c.Bool("nodnsseed"))
			*cx.Config.DisableDNSSeed = c.Bool("nodnsseed")
		}
		if c.IsSet("externalip") {
			L.Trace("set externalip", c.StringSlice("externalip"))
			*cx.Config.ExternalIPs = c.StringSlice("externalip")
		}
		if c.IsSet("addcheckpoint") {
			L.Trace("set addcheckpoint", c.StringSlice("addcheckpoint"))
			*cx.Config.AddCheckpoints = c.StringSlice("addcheckpoint")
		}
		if c.IsSet("nocheckpoints") {
			L.Trace("set nocheckpoints", c.Bool("nocheckpoints"))
			*cx.Config.DisableCheckpoints = c.Bool("nocheckpoints")
		}
		if c.IsSet("dbtype") {
			L.Trace("set dbtype", c.String("dbtype"))
			*cx.Config.DbType = c.String("dbtype")
		}
		if c.IsSet("profile") {
			L.Trace("set profile", c.String("profile"))
			*cx.Config.Profile = c.String("profile")
		}
		if c.IsSet("cpuprofile") {
			L.Trace("set cpuprofile", c.String("cpuprofile"))
			*cx.Config.CPUProfile = c.String("cpuprofile")
		}
		if c.IsSet("upnp") {
			L.Trace("set upnp", c.Bool("upnp"))
			*cx.Config.Upnp = c.Bool("upnp")
		}
		if c.IsSet("minrelaytxfee") {
			L.Trace("set minrelaytxfee", c.Float64("minrelaytxfee"))
			*cx.Config.MinRelayTxFee = c.Float64("minrelaytxfee")
		}
		if c.IsSet("limitfreerelay") {
			L.Trace("set limitfreerelay", c.Float64("limitfreerelay"))
			*cx.Config.FreeTxRelayLimit = c.Float64("limitfreerelay")
		}
		if c.IsSet("norelaypriority") {
			L.Trace("set norelaypriority", c.Bool("norelaypriority"))
			*cx.Config.NoRelayPriority = c.Bool("norelaypriority")
		}
		if c.IsSet("trickleinterval") {
			L.Trace("set trickleinterval", c.Duration("trickleinterval"))
			*cx.Config.TrickleInterval = c.Duration("trickleinterval")
		}
		if c.IsSet("maxorphantx") {
			L.Trace("set maxorphantx", c.Int("maxorphantx"))
			*cx.Config.MaxOrphanTxs = c.Int("maxorphantx")
		}
		if c.IsSet("algo") {
			L.Trace("set algo", c.String("algo"))
			*cx.Config.Algo = c.String("algo")
		}
		if c.IsSet("generate") {
			L.Trace("set generate", c.Bool("generate"))
			*cx.Config.Generate = c.Bool("generate")
		}
		if c.IsSet("genthreads") {
			L.Trace("set genthreads", c.Int("genthreads"))
			*cx.Config.GenThreads = c.Int("genthreads")
		}
		if c.IsSet("nocontroller") {
			L.Trace("set nocontroller",
				c.String("nocontroller"))
			*cx.Config.NoController = c.Bool("nocontroller")
		}
		if c.IsSet("miningaddr") {
			L.Trace("set miningaddr", c.StringSlice("miningaddr"))
			*cx.Config.MiningAddrs = c.StringSlice("miningaddr")
		}
		if c.IsSet("minerpass") {
			L.Trace("set minerpass", c.String("minerpass"))
			*cx.Config.MinerPass = c.String("minerpass")
		}
		if c.IsSet("group") {
			L.Trace("set group", c.String("group"))
			*cx.Config.Group = c.String("group")
		}
		if c.IsSet("nodiscovery") {
			L.Trace("set nodiscovery",
				c.String("nodiscovery"))
			*cx.Config.NoDiscovery = c.Bool("nodiscovery")
		}
		if c.IsSet("blockminsize") {
			L.Trace("set blockminsize", c.Int("blockminsize"))
			*cx.Config.BlockMinSize = c.Int("blockminsize")
		}
		if c.IsSet("blockmaxsize") {
			L.Trace("set blockmaxsize", c.Int("blockmaxsize"))
			*cx.Config.BlockMaxSize = c.Int("blockmaxsize")
		}
		if c.IsSet("blockminweight") {
			L.Trace("set blockminweight", c.Int("blockminweight"))
			*cx.Config.BlockMinWeight = c.Int("blockminweight")
		}
		if c.IsSet("blockmaxweight") {
			L.Trace("set blockmaxweight", c.Int("blockmaxweight"))
			*cx.Config.BlockMaxWeight = c.Int("blockmaxweight")
		}
		if c.IsSet("blockprioritysize") {
			L.Trace("set blockprioritysize", c.Int("blockprioritysize"))
			*cx.Config.BlockPrioritySize = c.Int("blockprioritysize")
		}
		if c.IsSet("uacomment") {
			L.Trace("set uacomment", c.StringSlice("uacomment"))
			*cx.Config.UserAgentComments = c.StringSlice("uacomment")
		}
		if c.IsSet("nopeerbloomfilters") {
			L.Trace("set nopeerbloomfilters", c.Bool("nopeerbloomfilters"))
			*cx.Config.NoPeerBloomFilters = c.Bool("nopeerbloomfilters")
		}
		if c.IsSet("nocfilters") {
			L.Trace("set nocfilters", c.Bool("nocfilters"))
			*cx.Config.NoCFilters = c.Bool("nocfilters")
		}
		if c.IsSet("sigcachemaxsize") {
			L.Trace("set sigcachemaxsize", c.Int("sigcachemaxsize"))
			*cx.Config.SigCacheMaxSize = c.Int("sigcachemaxsize")
		}
		if c.IsSet("blocksonly") {
			L.Trace("set blocksonly", c.Bool("blocksonly"))
			*cx.Config.BlocksOnly = c.Bool("blocksonly")
		}
		if c.IsSet("notxindex") {
			L.Trace("set notxindex", c.Bool("notxindex"))
			*cx.Config.TxIndex = c.Bool("notxindex")
		}
		if c.IsSet("noaddrindex") {
			L.Trace("set noaddrindex", c.Bool("noaddrindex"))
			*cx.Config.AddrIndex = c.Bool("noaddrindex")
		}
		if c.IsSet("relaynonstd") {
			L.Trace("set relaynonstd", c.Bool("relaynonstd"))
			*cx.Config.RelayNonStd = c.Bool("relaynonstd")
		}
		if c.IsSet("rejectnonstd") {
			L.Trace("set rejectnonstd", c.Bool("rejectnonstd"))
			*cx.Config.RejectNonStd = c.Bool("rejectnonstd")
		}
		if c.IsSet("noinitialload") {
			L.Trace("set noinitialload", c.Bool("noinitialload"))
			*cx.Config.NoInitialLoad = c.Bool("noinitialload")
		}
		if c.IsSet("walletconnect") {
			L.Trace("set walletconnect", c.Bool("walletconnect"))
			*cx.Config.Wallet = c.Bool("walletconnect")
		}
		if c.IsSet("walletserver") {
			L.Trace("set walletserver", c.String("walletserver"))
			*cx.Config.WalletServer = c.String("walletserver")
		}
		if c.IsSet("walletpass") {
			L.Trace("set walletpass", c.String("walletpass"))
			*cx.Config.WalletPass = c.String("walletpass")
		}
		if c.IsSet("onetimetlskey") {
			L.Trace("set onetimetlskey", c.Bool("onetimetlskey"))
			*cx.Config.OneTimeTLSKey = c.Bool("onetimetlskey")
		}
		if c.IsSet("walletrpclisten") {
			L.Trace("set walletrpclisten", c.StringSlice("walletrpclisten"))
			*cx.Config.WalletRPCListeners = c.StringSlice("walletrpclisten")
		}
		if c.IsSet("walletrpcmaxclients") {
			L.Trace("set walletrpcmaxclients", c.Int("walletrpcmaxclients"))
			*cx.Config.WalletRPCMaxClients = c.Int("walletrpcmaxclients")
		}
		if c.IsSet("walletrpcmaxwebsockets") {
			L.Trace("set walletrpcmaxwebsockets",
				c.Int("walletrpcmaxwebsockets"))
			*cx.Config.WalletRPCMaxWebsockets = c.Int("walletrpcmaxwebsockets")
		}
		if c.IsSet("experimentalrpclisten") {
			L.Trace("set experimentalrpclisten",
				c.StringSlice("experimentalrpclisten"))
			*cx.Config.ExperimentalRPCListeners = c.StringSlice("experimentalrpclisten")
		}
		if c.IsSet("nodeoff") {
			L.Trace("set nodeoff", c.Bool("nodeoff"))
			*cx.Config.NodeOff = c.Bool("nodeoff")
		}
		if c.IsSet("testnodeoff") {
			L.Trace("set testnodeoff", c.Bool("testnodeoff"))
			*cx.Config.TestNodeOff = c.Bool("testnodeoff")
		}
		if c.IsSet("walletoff") {
			L.Trace("set walletoff", c.Bool("walletoff"))
			*cx.Config.WalletOff = c.Bool("walletoff")
		}
		if c.IsSet("save") {
			L.Trace("set save", c.Bool("save"))
			cx.StateCfg.Save = true
		}
		return nil
	}
}
