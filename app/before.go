package app

import (
	"encoding/json"
	"fmt"
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
			log.INFO("setting datadir", *cx.Config.DataDir)
		}
		*cx.Config.ConfigFile =
			*cx.Config.DataDir + string(
				os.PathSeparator) +
				podConfigFilename
		log.INFO("config file set to", *cx.Config.ConfigFile)
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(*cx.Config.ConfigFile) {
			log.TRACE("loading config")
			b, err := ioutil.ReadFile(*cx.Config.ConfigFile)
			log.INFO("loaded config")
			if err == nil {
				*cx.Config = *pod.EmptyConfig()
				err = json.Unmarshal(b, cx.Config)
				if err != nil {
					fmt.Println("error unmarshalling config", err)
					os.Exit(1)
				}
				log.INFO("unmarshalled config")
			} else {
				fmt.Println("unexpected error reading configuration file:", err)
				os.Exit(1)
			}
		} else {
			log.INFO("will save config after configuration")
			cx.StateCfg.Save = true
		}
		log.TRACE("checking log level")
		log.TRACE("set loglevel", c.String("loglevel"))
		if c.String("loglevel") != "" {
			*cx.Config.LogLevel = c.String("loglevel")
			log.L.SetLevel(*cx.Config.LogLevel, true)
		}
		log.TRACE("checking network")
		log.TRACE("set network", c.String("network"))
		if c.IsSet("network") {
			log.TRACE("network is set to", *cx.Config.Network)
			*cx.Config.Network = c.String("network")
			switch *cx.Config.Network {
			case "testnet", "testnet3", "t":
				log.TRACE("on testnet")
				cx.ActiveNet = &netparams.TestNet3Params
				fork.IsTestnet = true
			case "regtestnet", "regressiontest", "r":
				log.TRACE("on regression testnet")
				cx.ActiveNet = &netparams.RegressionTestParams
			case "simnet", "s":
				log.TRACE("on simnet")
				cx.ActiveNet = &netparams.SimNetParams
			default:
				if *cx.Config.Network != "mainnet" &&
					*cx.Config.Network != "m" {
					log.WARN("using mainnet for node")
					log.TRACE("on mainnet")
				}
				cx.ActiveNet = &netparams.MainNetParams
			}
		}
		log.TRACE("set username", c.String("username"))
		if c.IsSet("username") {
			*cx.Config.Username = c.String("username")
		}
		log.TRACE("set password", c.String("password"))
		if c.IsSet("password") {
			*cx.Config.Password = c.String("password")
		}
		log.TRACE("set serveruser", c.String("serveruser"))
		if c.IsSet("serveruser") {
			*cx.Config.ServerUser = c.String("serveruser")
		}
		log.TRACE("set serverpass", c.String("serverpass"))
		if c.IsSet("serverpass") {
			*cx.Config.ServerPass = c.String("serverpass")
		}
		log.TRACE("set limituser", c.String("limituser"))
		if c.IsSet("limituser") {
			*cx.Config.LimitUser = c.String("limituser")
		}
		log.TRACE("set limitpass", c.String("limitpass"))
		if c.IsSet("limitpass") {
			*cx.Config.LimitPass = c.String("limitpass")
		}
		log.TRACE("set rpccert", c.String("rpccert"))
		if c.IsSet("rpccert") {
			*cx.Config.RPCCert = c.String("rpccert")
		}
		log.TRACE("set rpckey", c.String("rpckey"))
		if c.IsSet("rpckey") {
			*cx.Config.RPCKey = c.String("rpckey")
		}
		log.TRACE("set cafile", c.String("cafile"))
		if c.IsSet("cafile") {
			*cx.Config.CAFile = c.String("cafile")
		}
		log.TRACE("set clienttls", c.Bool("clienttls"))
		if c.IsSet("clienttls") {
			*cx.Config.TLS = c.Bool("clienttls")
		}
		log.TRACE("set servertls", c.Bool("servertls"))
		if c.IsSet("servertls") {
			*cx.Config.ServerTLS = c.Bool("servertls")
		}
		log.TRACE("set tlsskipverify ", c.Bool("tlsskipverify"))
		if c.IsSet("tlsskipverify") {
			*cx.Config.TLSSkipVerify = c.Bool("tlsskipverify")
		}
		log.TRACE("set proxy", c.String("proxy"))
		if c.IsSet("proxy") {
			*cx.Config.Proxy = c.String("proxy")
		}
		log.TRACE("set proxyuser", c.String("proxyuser"))
		if c.IsSet("proxyuser") {
			*cx.Config.ProxyUser = c.String("proxyuser")
		}
		log.TRACE("set proxypass", c.String("proxypass"))
		if c.IsSet("proxypass") {
			*cx.Config.ProxyPass = c.String("proxypass")
		}
		log.TRACE("set onion", c.Bool("onion"))
		if c.IsSet("onion") {
			*cx.Config.Onion = c.Bool("onion")
		}
		log.TRACE("set onionproxy", c.String("onionproxy"))
		if c.IsSet("onionproxy") {
			*cx.Config.OnionProxy = c.String("onionproxy")
		}
		log.TRACE("set onionuser", c.String("onionuser"))
		if c.IsSet("onionuser") {
			*cx.Config.OnionProxyUser = c.String("onionuser")
		}
		log.TRACE("set onionpass", c.String("onionpass"))
		if c.IsSet("onionpass") {
			*cx.Config.OnionProxyPass = c.String("onionpass")
		}
		log.TRACE("set torisolation", c.Bool("torisolation"))
		if c.IsSet("torisolation") {
			*cx.Config.TorIsolation = c.Bool("torisolation")
		}
		log.TRACE("set addpeer", c.StringSlice("addpeer"))
		if c.IsSet("addpeer") {
			*cx.Config.AddPeers = c.StringSlice("addpeer")
		}
		log.TRACE("set connect", c.StringSlice("connect"))
		if c.IsSet("connect") {
			*cx.Config.ConnectPeers = c.StringSlice("connect")
		}
		log.TRACE("set nolisten", c.Bool("nolisten"))
		if c.IsSet("nolisten") {
			*cx.Config.DisableListen = c.Bool("nolisten")
		}
		log.TRACE("set listen", c.StringSlice("listen"))
		if c.IsSet("listen") {
			*cx.Config.Listeners = c.StringSlice("listen")
		}
		log.TRACE("set maxpeers", c.Int("maxpeers"))
		if c.IsSet("maxpeers") {
			*cx.Config.MaxPeers = c.Int("maxpeers")
		}
		log.TRACE("set nobanning", c.Bool("nobanning"))
		if c.IsSet("nobanning") {
			*cx.Config.DisableBanning = c.Bool("nobanning")
		}
		log.TRACE("set banduration", c.Duration("banduration"))
		if c.IsSet("banduration") {
			*cx.Config.BanDuration = c.Duration("banduration")
		}
		log.TRACE("set banthreshold", c.Int("banthreshold"))
		if c.IsSet("banthreshold") {
			*cx.Config.BanThreshold = c.Int("banthreshold")
		}
		log.TRACE("set whitelist", c.StringSlice("whitelist"))
		if c.IsSet("whitelist") {
			*cx.Config.Whitelists = c.StringSlice("whitelist")
		}
		log.TRACE("set rpcconnect", c.String("rpcconnect"))
		if c.IsSet("rpcconnect") {
			*cx.Config.RPCConnect = c.String("rpcconnect")
		}
		log.TRACE("set rpclisten", c.StringSlice("rpclisten"))
		if c.IsSet("rpclisten") {
			*cx.Config.RPCListeners = c.StringSlice("rpclisten")
		}
		log.TRACE("set rpcmaxclients", c.Int("rpcmaxclients"))
		if c.IsSet("rpcmaxclients") {
			*cx.Config.RPCMaxClients = c.Int("rpcmaxclients")
		}
		log.TRACE("set rpcmaxwebsockets", c.Int("rpcmaxwebsockets"))
		if c.IsSet("rpcmaxwebsockets") {
			*cx.Config.RPCMaxWebsockets = c.Int("rpcmaxwebsockets")
		}
		log.TRACE("set rpcmaxconcurrentreqs", c.Int("rpcmaxconcurrentreqs"))
		if c.IsSet("rpcmaxconcurrentreqs") {
			*cx.Config.RPCMaxConcurrentReqs = c.Int("rpcmaxconcurrentreqs")
		}
		log.TRACE("set rpcquirks", c.Bool("rpcquirks"))
		if c.IsSet("rpcquirks") {
			*cx.Config.RPCQuirks = c.Bool("rpcquirks")
		}
		log.TRACE("set norpc", c.Bool("norpc"))
		if c.IsSet("norpc") {
			*cx.Config.DisableRPC = c.Bool("norpc")
		}
		log.TRACE("set nodnsseed", c.Bool("nodnsseed"))
		if c.IsSet("nodnsseed") {
			*cx.Config.DisableDNSSeed = c.Bool("nodnsseed")
		}
		log.TRACE("set externalip", c.StringSlice("externalip"))
		if c.IsSet("externalip") {
			*cx.Config.ExternalIPs = c.StringSlice("externalip")
		}
		log.TRACE("set addcheckpoint", c.StringSlice("addcheckpoint"))
		if c.IsSet("addcheckpoint") {
			*cx.Config.AddCheckpoints = c.StringSlice("addcheckpoint")
		}
		log.TRACE("set nocheckpoints", c.Bool("nocheckpoints"))
		if c.IsSet("nocheckpoints") {
			*cx.Config.DisableCheckpoints = c.Bool("nocheckpoints")
		}
		log.TRACE("set dbtype", c.String("dbtype"))
		if c.IsSet("dbtype") {
			*cx.Config.DbType = c.String("dbtype")
		}
		log.TRACE("set profile", c.String("profile"))
		if c.IsSet("profile") {
			*cx.Config.Profile = c.String("profile")
		}
		log.TRACE("set cpuprofile", c.String("cpuprofile"))
		if c.IsSet("cpuprofile") {
			*cx.Config.CPUProfile = c.String("cpuprofile")
		}
		log.TRACE("set upnp", c.Bool("upnp"))
		if c.IsSet("upnp") {
			*cx.Config.Upnp = c.Bool("upnp")
		}
		log.TRACE("set minrelaytxfee", c.Float64("minrelaytxfee"))
		if c.IsSet("minrelaytxfee") {
			*cx.Config.MinRelayTxFee = c.Float64("minrelaytxfee")
		}
		log.TRACE("set limitfreerelay", c.Float64("limitfreerelay"))
		if c.IsSet("limitfreerelay") {
			*cx.Config.FreeTxRelayLimit = c.Float64("limitfreerelay")
		}
		log.TRACE("set norelaypriority", c.Bool("norelaypriority"))
		if c.IsSet("norelaypriority") {
			*cx.Config.NoRelayPriority = c.Bool("norelaypriority")
		}
		log.TRACE("set trickleinterval", c.Duration("trickleinterval"))
		if c.IsSet("trickleinterval") {
			*cx.Config.TrickleInterval = c.Duration("trickleinterval")
		}
		log.TRACE("set maxorphantx", c.Int("maxorphantx"))
		if c.IsSet("maxorphantx") {
			*cx.Config.MaxOrphanTxs = c.Int("maxorphantx")
		}
		log.TRACE("set algo", c.String("algo"))
		if c.IsSet("algo") {
			*cx.Config.Algo = c.String("algo")
		}
		log.TRACE("set generate", c.Bool("generate"))
		if c.IsSet("generate") {
			*cx.Config.Generate = c.Bool("generate")
		}
		log.TRACE("set genthreads", c.Int("genthreads"))
		if c.IsSet("genthreads") {
			*cx.Config.GenThreads = c.Int("genthreads")
		}
		log.TRACE("set nocontroller", c.String("nocontroller"))
		if c.IsSet("nocontroller") {
			*cx.Config.NoController = c.Bool("nocontroller")
		}
		log.TRACE("set broadcast", c.Bool("broadcast"))
		if c.IsSet("broadcast") {
			*cx.Config.Broadcast = c.Bool("broadcast")
		}
		log.TRACE("set workers", c.StringSlice("workers"))
		if c.IsSet("workers") {
			*cx.Config.Workers = c.StringSlice("workers")
		}
		log.TRACE("set miningaddr", c.StringSlice("miningaddr"))
		if c.IsSet("miningaddr") {
			*cx.Config.MiningAddrs = c.StringSlice("miningaddr")
		}
		log.TRACE("set minerpass", c.String("minerpass"))
		if c.IsSet("minerpass") {
			*cx.Config.MinerPass = c.String("minerpass")
		}
		log.TRACE("set blockminsize", c.Int("blockminsize"))
		if c.IsSet("blockminsize") {
			*cx.Config.BlockMinSize = c.Int("blockminsize")
		}
		log.TRACE("set blockmaxsize", c.Int("blockmaxsize"))
		if c.IsSet("blockmaxsize") {
			*cx.Config.BlockMaxSize = c.Int("blockmaxsize")
		}
		log.TRACE("set blockminweight", c.Int("blockminweight"))
		if c.IsSet("blockminweight") {
			*cx.Config.BlockMinWeight = c.Int("blockminweight")
		}
		log.TRACE("set blockmaxweight", c.Int("blockmaxweight"))
		if c.IsSet("blockmaxweight") {
			*cx.Config.BlockMaxWeight = c.Int("blockmaxweight")
		}
		log.TRACE("set blockprioritysize", c.Int("blockprioritysize"))
		if c.IsSet("blockprioritysize") {
			*cx.Config.BlockPrioritySize = c.Int("blockprioritysize")
		}
		log.TRACE("set uacomment", c.StringSlice("uacomment"))
		if c.IsSet("uacomment") {
			*cx.Config.UserAgentComments = c.StringSlice("uacomment")
		}
		log.TRACE("set nopeerbloomfilters", c.Bool("nopeerbloomfilters"))
		if c.IsSet("nopeerbloomfilters") {
			*cx.Config.NoPeerBloomFilters = c.Bool("nopeerbloomfilters")
		}
		log.TRACE("set nocfilters", c.Bool("nocfilters"))
		if c.IsSet("nocfilters") {
			*cx.Config.NoCFilters = c.Bool("nocfilters")
		}
		log.TRACE("set sigcachemaxsize", c.Int("sigcachemaxsize"))
		if c.IsSet("sigcachemaxsize") {
			*cx.Config.SigCacheMaxSize = c.Int("sigcachemaxsize")
		}
		log.TRACE("set blocksonly", c.Bool("blocksonly"))
		if c.IsSet("blocksonly") {
			*cx.Config.BlocksOnly = c.Bool("blocksonly")
		}
		log.TRACE("set notxindex", c.Bool("notxindex"))
		if c.IsSet("notxindex") {
			*cx.Config.TxIndex = c.Bool("notxindex")
		}
		log.TRACE("set noaddrindex", c.Bool("noaddrindex"))
		if c.IsSet("noaddrindex") {
			*cx.Config.AddrIndex = c.Bool("noaddrindex")
		}
		log.TRACE("set relaynonstd", c.Bool("relaynonstd"))
		if c.IsSet("relaynonstd") {
			*cx.Config.RelayNonStd = c.Bool("relaynonstd")
		}
		log.TRACE("set rejectnonstd", c.Bool("rejectnonstd"))
		if c.IsSet("rejectnonstd") {
			*cx.Config.RejectNonStd = c.Bool("rejectnonstd")
		}
		log.TRACE("set noinitialload", c.Bool("noinitialload"))
		if c.IsSet("noinitialload") {
			*cx.Config.NoInitialLoad = c.Bool("noinitialload")
		}
		log.TRACE("set walletconnect", c.Bool("walletconnect"))
		if c.IsSet("walletconnect") {
			*cx.Config.Wallet = c.Bool("walletconnect")
		}
		log.TRACE("set walletserver", c.String("walletserver"))
		if c.IsSet("walletserver") {
			*cx.Config.WalletServer = c.String("walletserver")
		}
		log.TRACE("set walletpass", c.String("walletpass"))
		if c.IsSet("walletpass") {
			*cx.Config.WalletPass = c.String("walletpass")
		}
		log.TRACE("set onetimetlskey", c.Bool("onetimetlskey"))
		if c.IsSet("onetimetlskey") {
			*cx.Config.OneTimeTLSKey = c.Bool("onetimetlskey")
		}
		log.TRACE("set walletrpclisten", c.StringSlice("walletrpclisten"))
		if c.IsSet("walletrpclisten") {
			*cx.Config.WalletRPCListeners = c.StringSlice("walletrpclisten")
		}
		log.TRACE("set walletrpcmaxclients", c.Int("walletrpcmaxclients"))
		if c.IsSet("walletrpcmaxclients") {
			*cx.Config.WalletRPCMaxClients = c.Int("walletrpcmaxclients")
		}
		log.TRACE("set walletrpcmaxwebsockets", c.Int("walletrpcmaxwebsockets"))
		if c.IsSet("walletrpcmaxwebsockets") {
			*cx.Config.WalletRPCMaxWebsockets = c.Int("walletrpcmaxwebsockets")
		}
		log.TRACE("set experimentalrpclisten", c.StringSlice("experimentalrpclisten"))
		if c.IsSet("experimentalrpclisten") {
			*cx.Config.ExperimentalRPCListeners = c.StringSlice("experimentalrpclisten")
		}
		log.TRACE("set nodeoff", c.Bool("nodeoff"))
		if c.IsSet("nodeoff") {
			*cx.Config.NodeOff = c.Bool("nodeoff")
		}
		log.TRACE("set testnodeoff", c.Bool("testnodeoff"))
		if c.IsSet("testnodeoff") {
			*cx.Config.TestNodeOff = c.Bool("testnodeoff")
		}
		log.TRACE("set walletoff", c.Bool("walletoff"))
		if c.IsSet("walletoff") {
			*cx.Config.WalletOff = c.Bool("walletoff")
		}
		log.TRACE("set save", c.Bool("save"))
		if c.IsSet("save") {
			// cx.StateCfg.Save = true
			log.INFO("saving configuration")
			save.Pod(cx.Config)
		}
		return nil
	}
}
