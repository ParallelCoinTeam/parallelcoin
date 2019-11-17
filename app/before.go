package app

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
)

func beforeFunc(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.INFO("running beforeFunc")
		// if user set datadir this is first thing to configure
		if c.IsSet("datadir") {
			*cx.Config.DataDir = c.String("datadir")
			cx.DataDir = c.String("datadir")
			log.TRACE("setting datadir", *cx.Config.DataDir)
		}
		*cx.Config.ConfigFile =
			*cx.Config.DataDir + string(os.PathSeparator) + podConfigFilename
		log.TRACE("config file set to", *cx.Config.ConfigFile)
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(*cx.Config.ConfigFile) {
			log.TRACE("loading config")
			b, err := ioutil.ReadFile(*cx.Config.ConfigFile)
			log.TRACE("loaded config")
			if err == nil {
				cx.Config = pod.EmptyConfig()
				err = json.Unmarshal(b, cx.Config)
				if err != nil {
					log.ERROR("error unmarshalling config", err)
					os.Exit(1)
				}
				log.TRACE("unmarshalled config")
			} else {
				log.FATAL("unexpected error reading configuration file:", err)
				os.Exit(1)
			}
		} else {
			*cx.Config.ConfigFile = ""
			log.TRACE("will save config after configuration")
			cx.StateCfg.Save = true
		}
		log.TRACE("checking log level")
		if c.String("loglevel") != "" {
			log.TRACE("set loglevel", c.String("loglevel"))
			*cx.Config.LogLevel = c.String("loglevel")
		}
		log.TRACE("checking network")
		if c.IsSet("network") {
			log.TRACE("set network", c.String("network"))
			*cx.Config.Network = c.String("network")
			switch *cx.Config.Network {
			case "testnet", "testnet3", "t":
				log.TRACE("on testnet")
				cx.ActiveNet = &netparams.TestNet3Params
				fork.IsTestnet = true
				fork.HashReps = 3
			case "regtestnet", "regressiontest", "r":
				log.TRACE("on regression testnet")
				fork.IsTestnet = true
				cx.ActiveNet = &netparams.RegressionTestParams
			case "simnet", "s":
				log.TRACE("on simnet")
				fork.IsTestnet = true
				cx.ActiveNet = &netparams.SimNetParams
			default:
				if *cx.Config.Network != "mainnet" &&
					*cx.Config.Network != "m" {
					log.WARN("using mainnet for node")
				}
				cx.ActiveNet = &netparams.MainNetParams
			}
		}
		if c.IsSet("username") {
			log.TRACE("set username", c.String("username"))
			*cx.Config.Username = c.String("username")
		}
		if c.IsSet("password") {
			log.TRACE("set password", c.String("password"))
			*cx.Config.Password = c.String("password")
		}
		if c.IsSet("serveruser") {
			log.TRACE("set serveruser", c.String("serveruser"))
			*cx.Config.ServerUser = c.String("serveruser")
		}
		if c.IsSet("serverpass") {
			log.TRACE("set serverpass", c.String("serverpass"))
			*cx.Config.ServerPass = c.String("serverpass")
		}
		if c.IsSet("limituser") {
			log.TRACE("set limituser", c.String("limituser"))
			*cx.Config.LimitUser = c.String("limituser")
		}
		if c.IsSet("limitpass") {
			log.TRACE("set limitpass", c.String("limitpass"))
			*cx.Config.LimitPass = c.String("limitpass")
		}
		if c.IsSet("rpccert") {
			log.TRACE("set rpccert", c.String("rpccert"))
			*cx.Config.RPCCert = c.String("rpccert")
		}
		if c.IsSet("rpckey") {
			log.TRACE("set rpckey", c.String("rpckey"))
			*cx.Config.RPCKey = c.String("rpckey")
		}
		if c.IsSet("cafile") {
			log.TRACE("set cafile", c.String("cafile"))
			*cx.Config.CAFile = c.String("cafile")
		}
		if c.IsSet("clienttls") {
			log.TRACE("set clienttls", c.Bool("clienttls"))
			*cx.Config.TLS = c.Bool("clienttls")
		}
		if c.IsSet("servertls") {
			log.TRACE("set servertls", c.Bool("servertls"))
			*cx.Config.ServerTLS = c.Bool("servertls")
		}
		if c.IsSet("tlsskipverify") {
			log.TRACE("set tlsskipverify ", c.Bool("tlsskipverify"))
			*cx.Config.TLSSkipVerify = c.Bool("tlsskipverify")
		}
		if c.IsSet("proxy") {
			log.TRACE("set proxy", c.String("proxy"))
			*cx.Config.Proxy = c.String("proxy")
		}
		if c.IsSet("proxyuser") {
			log.TRACE("set proxyuser", c.String("proxyuser"))
			*cx.Config.ProxyUser = c.String("proxyuser")
		}
		if c.IsSet("proxypass") {
			log.TRACE("set proxypass", c.String("proxypass"))
			*cx.Config.ProxyPass = c.String("proxypass")
		}
		if c.IsSet("onion") {
			log.TRACE("set onion", c.Bool("onion"))
			*cx.Config.Onion = c.Bool("onion")
		}
		if c.IsSet("onionproxy") {
			log.TRACE("set onionproxy", c.String("onionproxy"))
			*cx.Config.OnionProxy = c.String("onionproxy")
		}
		if c.IsSet("onionuser") {
			log.TRACE("set onionuser", c.String("onionuser"))
			*cx.Config.OnionProxyUser = c.String("onionuser")
		}
		if c.IsSet("onionpass") {
			log.TRACE("set onionpass", c.String("onionpass"))
			*cx.Config.OnionProxyPass = c.String("onionpass")
		}
		if c.IsSet("torisolation") {
			log.TRACE("set torisolation", c.Bool("torisolation"))
			*cx.Config.TorIsolation = c.Bool("torisolation")
		}
		if c.IsSet("addpeer") {
			log.TRACE("set addpeer", c.StringSlice("addpeer"))
			*cx.Config.AddPeers = c.StringSlice("addpeer")
		}
		if c.IsSet("connect") {
			log.TRACE("set connect", c.StringSlice("connect"))
			*cx.Config.ConnectPeers = c.StringSlice("connect")
		}
		if c.IsSet("nolisten") {
			log.TRACE("set nolisten", c.Bool("nolisten"))
			*cx.Config.DisableListen = c.Bool("nolisten")
		}
		if c.IsSet("listen") {
			log.TRACE("set listen", c.StringSlice("listen"))
			*cx.Config.Listeners = c.StringSlice("listen")
		}
		if c.IsSet("maxpeers") {
			log.TRACE("set maxpeers", c.Int("maxpeers"))
			*cx.Config.MaxPeers = c.Int("maxpeers")
		}
		if c.IsSet("nobanning") {
			log.TRACE("set nobanning", c.Bool("nobanning"))
			*cx.Config.DisableBanning = c.Bool("nobanning")
		}
		if c.IsSet("banduration") {
			log.TRACE("set banduration", c.Duration("banduration"))
			*cx.Config.BanDuration = c.Duration("banduration")
		}
		if c.IsSet("banthreshold") {
			log.TRACE("set banthreshold", c.Int("banthreshold"))
			*cx.Config.BanThreshold = c.Int("banthreshold")
		}
		if c.IsSet("whitelist") {
			log.TRACE("set whitelist", c.StringSlice("whitelist"))
			*cx.Config.Whitelists = c.StringSlice("whitelist")
		}
		if c.IsSet("rpcconnect") {
			log.TRACE("set rpcconnect", c.String("rpcconnect"))
			*cx.Config.RPCConnect = c.String("rpcconnect")
		}
		if c.IsSet("rpclisten") {
			log.TRACE("set rpclisten", c.StringSlice("rpclisten"))
			*cx.Config.RPCListeners = c.StringSlice("rpclisten")
		}
		if c.IsSet("rpcmaxclients") {
			log.TRACE("set rpcmaxclients", c.Int("rpcmaxclients"))
			*cx.Config.RPCMaxClients = c.Int("rpcmaxclients")
		}
		if c.IsSet("rpcmaxwebsockets") {
			log.TRACE("set rpcmaxwebsockets", c.Int("rpcmaxwebsockets"))
			*cx.Config.RPCMaxWebsockets = c.Int("rpcmaxwebsockets")
		}
		if c.IsSet("rpcmaxconcurrentreqs") {
			log.TRACE("set rpcmaxconcurrentreqs", c.Int("rpcmaxconcurrentreqs"))
			*cx.Config.RPCMaxConcurrentReqs = c.Int("rpcmaxconcurrentreqs")
		}
		if c.IsSet("rpcquirks") {
			log.TRACE("set rpcquirks", c.Bool("rpcquirks"))
			*cx.Config.RPCQuirks = c.Bool("rpcquirks")
		}
		if c.IsSet("norpc") {
			log.TRACE("set norpc", c.Bool("norpc"))
			*cx.Config.DisableRPC = c.Bool("norpc")
		}
		if c.IsSet("nodnsseed") {
			log.TRACE("set nodnsseed", c.Bool("nodnsseed"))
			*cx.Config.DisableDNSSeed = c.Bool("nodnsseed")
		}
		if c.IsSet("externalip") {
			log.TRACE("set externalip", c.StringSlice("externalip"))
			*cx.Config.ExternalIPs = c.StringSlice("externalip")
		}
		if c.IsSet("addcheckpoint") {
			log.TRACE("set addcheckpoint", c.StringSlice("addcheckpoint"))
			*cx.Config.AddCheckpoints = c.StringSlice("addcheckpoint")
		}
		if c.IsSet("nocheckpoints") {
			log.TRACE("set nocheckpoints", c.Bool("nocheckpoints"))
			*cx.Config.DisableCheckpoints = c.Bool("nocheckpoints")
		}
		if c.IsSet("dbtype") {
			log.TRACE("set dbtype", c.String("dbtype"))
			*cx.Config.DbType = c.String("dbtype")
		}
		if c.IsSet("profile") {
			log.TRACE("set profile", c.String("profile"))
			*cx.Config.Profile = c.String("profile")
		}
		if c.IsSet("cpuprofile") {
			log.TRACE("set cpuprofile", c.String("cpuprofile"))
			*cx.Config.CPUProfile = c.String("cpuprofile")
		}
		if c.IsSet("upnp") {
			log.TRACE("set upnp", c.Bool("upnp"))
			*cx.Config.UPNP = c.Bool("upnp")
		}
		if c.IsSet("minrelaytxfee") {
			log.TRACE("set minrelaytxfee", c.Float64("minrelaytxfee"))
			*cx.Config.MinRelayTxFee = c.Float64("minrelaytxfee")
		}
		if c.IsSet("limitfreerelay") {
			log.TRACE("set limitfreerelay", c.Float64("limitfreerelay"))
			*cx.Config.FreeTxRelayLimit = c.Float64("limitfreerelay")
		}
		if c.IsSet("norelaypriority") {
			log.TRACE("set norelaypriority", c.Bool("norelaypriority"))
			*cx.Config.NoRelayPriority = c.Bool("norelaypriority")
		}
		if c.IsSet("trickleinterval") {
			log.TRACE("set trickleinterval", c.Duration("trickleinterval"))
			*cx.Config.TrickleInterval = c.Duration("trickleinterval")
		}
		if c.IsSet("maxorphantx") {
			log.TRACE("set maxorphantx", c.Int("maxorphantx"))
			*cx.Config.MaxOrphanTxs = c.Int("maxorphantx")
		}
		if c.IsSet("algo") {
			log.TRACE("set algo", c.String("algo"))
			*cx.Config.Algo = c.String("algo")
		}
		if c.IsSet("generate") {
			log.TRACE("set generate", c.Bool("generate"))
			*cx.Config.Generate = c.Bool("generate")
		}
		if c.IsSet("genthreads") {
			log.TRACE("set genthreads", c.Int("genthreads"))
			*cx.Config.GenThreads = c.Int("genthreads")
		}
		if c.IsSet("solo") {
			log.WARN("set solo", c.Bool("solo"))
			*cx.Config.Solo = c.Bool("solo")
		}
		if c.IsSet("controller") {
			log.TRACE("set controller listener address", c.String("controller"))
			*cx.Config.Controller = c.String("controller")
		}
		if c.IsSet("autoports") {
			log.TRACE("set autoports", c.String("autoports"))
			*cx.Config.AutoPorts = c.Bool("autoports")
		}
		if c.IsSet("nocontroller") {
			log.TRACE("set nocontroller", c.String("nocontroller"))
			*cx.Config.NoController = c.Bool("nocontroller")
		}
		if c.IsSet("miningaddrs") {
			log.TRACE("set miningaddrs", c.StringSlice("miningaddrs"))
			*cx.Config.MiningAddrs = c.StringSlice("miningaddrs")
		}
		if c.IsSet("minerpass") {
			log.TRACE("set minerpass", c.String("minerpass"))
			*cx.Config.MinerPass = c.String("minerpass")
		}
		if c.IsSet("blockminsize") {
			log.TRACE("set blockminsize", c.Int("blockminsize"))
			*cx.Config.BlockMinSize = c.Int("blockminsize")
		}
		if c.IsSet("blockmaxsize") {
			log.TRACE("set blockmaxsize", c.Int("blockmaxsize"))
			*cx.Config.BlockMaxSize = c.Int("blockmaxsize")
		}
		if c.IsSet("blockminweight") {
			log.TRACE("set blockminweight", c.Int("blockminweight"))
			*cx.Config.BlockMinWeight = c.Int("blockminweight")
		}
		if c.IsSet("blockmaxweight") {
			log.TRACE("set blockmaxweight", c.Int("blockmaxweight"))
			*cx.Config.BlockMaxWeight = c.Int("blockmaxweight")
		}
		if c.IsSet("blockprioritysize") {
			log.TRACE("set blockprioritysize", c.Int("blockprioritysize"))
			*cx.Config.BlockPrioritySize = c.Int("blockprioritysize")
		}
		if c.IsSet("uacomment") {
			log.TRACE("set uacomment", c.StringSlice("uacomment"))
			*cx.Config.UserAgentComments = c.StringSlice("uacomment")
		}
		if c.IsSet("nopeerbloomfilters") {
			log.TRACE("set nopeerbloomfilters", c.Bool("nopeerbloomfilters"))
			*cx.Config.NoPeerBloomFilters = c.Bool("nopeerbloomfilters")
		}
		if c.IsSet("nocfilters") {
			log.TRACE("set nocfilters", c.Bool("nocfilters"))
			*cx.Config.NoCFilters = c.Bool("nocfilters")
		}
		if c.IsSet("sigcachemaxsize") {
			log.TRACE("set sigcachemaxsize", c.Int("sigcachemaxsize"))
			*cx.Config.SigCacheMaxSize = c.Int("sigcachemaxsize")
		}
		if c.IsSet("blocksonly") {
			log.TRACE("set blocksonly", c.Bool("blocksonly"))
			*cx.Config.BlocksOnly = c.Bool("blocksonly")
		}
		if c.IsSet("notxindex") {
			log.TRACE("set notxindex", c.Bool("notxindex"))
			*cx.Config.TxIndex = c.Bool("notxindex")
		}
		if c.IsSet("noaddrindex") {
			log.TRACE("set noaddrindex", c.Bool("noaddrindex"))
			*cx.Config.AddrIndex = c.Bool("noaddrindex")
		}
		if c.IsSet("relaynonstd") {
			log.TRACE("set relaynonstd", c.Bool("relaynonstd"))
			*cx.Config.RelayNonStd = c.Bool("relaynonstd")
		}
		if c.IsSet("rejectnonstd") {
			log.TRACE("set rejectnonstd", c.Bool("rejectnonstd"))
			*cx.Config.RejectNonStd = c.Bool("rejectnonstd")
		}
		if c.IsSet("noinitialload") {
			log.TRACE("set noinitialload", c.Bool("noinitialload"))
			*cx.Config.NoInitialLoad = c.Bool("noinitialload")
		}
		if c.IsSet("walletconnect") {
			log.TRACE("set walletconnect", c.Bool("walletconnect"))
			*cx.Config.Wallet = c.Bool("walletconnect")
		}
		if c.IsSet("walletserver") {
			log.TRACE("set walletserver", c.String("walletserver"))
			*cx.Config.WalletServer = c.String("walletserver")
		}
		if c.IsSet("walletpass") {
			log.TRACE("set walletpass", c.String("walletpass"))
			*cx.Config.WalletPass = c.String("walletpass")
		}
		if c.IsSet("onetimetlskey") {
			log.TRACE("set onetimetlskey", c.Bool("onetimetlskey"))
			*cx.Config.OneTimeTLSKey = c.Bool("onetimetlskey")
		}
		if c.IsSet("walletrpclisten") {
			log.TRACE("set walletrpclisten", c.StringSlice("walletrpclisten"))
			*cx.Config.WalletRPCListeners = c.StringSlice("walletrpclisten")
		}
		if c.IsSet("walletrpcmaxclients") {
			log.TRACE("set walletrpcmaxclients", c.Int("walletrpcmaxclients"))
			*cx.Config.WalletRPCMaxClients = c.Int("walletrpcmaxclients")
		}
		if c.IsSet("walletrpcmaxwebsockets") {
			log.TRACE("set walletrpcmaxwebsockets", c.Int("walletrpcmaxwebsockets"))
			*cx.Config.WalletRPCMaxWebsockets = c.Int("walletrpcmaxwebsockets")
		}
		if c.IsSet("experimentalrpclisten") {
			log.TRACE("set experimentalrpclisten", c.StringSlice("experimentalrpclisten"))
			*cx.Config.ExperimentalRPCListeners = c.StringSlice("experimentalrpclisten")
		}
		if c.IsSet("nodeoff") {
			log.TRACE("set nodeoff", c.Bool("nodeoff"))
			*cx.Config.NodeOff = c.Bool("nodeoff")
		}
		if c.IsSet("walletoff") {
			log.TRACE("set walletoff", c.Bool("walletoff"))
			*cx.Config.WalletOff = c.Bool("walletoff")
		}
		if c.IsSet("save") {
			log.TRACE("set save", c.Bool("save"))
			// cx.StateCfg.Save = true
			log.INFO("saving configuration")
			save.Pod(cx.Config)
		}
		return nil
	}
}
