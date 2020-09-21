package app

import (
	"encoding/json"
	"fmt"
	"github.com/stalker-loki/pod/app/save"
	"github.com/stalker-loki/pod/pkg/util/logi"
	"github.com/stalker-loki/pod/pkg/util/logi/serve"
	"io/ioutil"
	pseudo "math/rand"
	"os"
	"runtime"
	"time"

	"github.com/stalker-loki/app/slog"
	"github.com/stalker-loki/pod/app/apputil"
	"github.com/stalker-loki/pod/app/conte"
	config "github.com/stalker-loki/pod/pkg/chain/config"
	"github.com/stalker-loki/pod/pkg/chain/config/netparams"
	"github.com/stalker-loki/pod/pkg/chain/fork"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/urfave/cli"
)

func beforeFunc(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		slog.Debug("running beforeFunc")
		cx.AppContext = c
		// if user set datadir this is first thing to configure
		if c.IsSet("datadir") {
			*cx.Config.DataDir = c.String("datadir")
			cx.DataDir = c.String("datadir")
		}
		if c.IsSet("pipelog") {
			slog.Warn("pipe logger enabled")
			*cx.Config.PipeLog = c.Bool("pipelog")
			serve.Log(cx.KillAll, save.Filters(*cx.Config.DataDir))
		}
		if c.IsSet("walletfile") {
			*cx.Config.WalletFile = c.String("walletfile")
		}
		*cx.Config.ConfigFile =
			*cx.Config.DataDir + string(os.PathSeparator) + podConfigFilename
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(*cx.Config.ConfigFile) {
			b, err := ioutil.ReadFile(*cx.Config.ConfigFile)
			if err == nil {
				cx.Config, cx.ConfigMap = pod.EmptyConfig()
				err = json.Unmarshal(b, cx.Config)
				if err != nil {
					slog.Error("error unmarshalling config", err)
					os.Exit(1)
				}
			} else {
				slog.Fatal("unexpected error reading configuration file:", err)
				os.Exit(1)
			}
		} else {
			*cx.Config.ConfigFile = ""
			slog.Debug("will save config after configuration")
			cx.StateCfg.Save = true
		}
		if c.IsSet("loglevel") {
			slog.Trace("set loglevel", c.String("loglevel"))
			*cx.Config.LogLevel = c.String("loglevel")
			color := true
			if runtime.GOOS == "windows" {
				color = false
			}
			logi.L.SetLevel(*cx.Config.LogLevel, color, "pod")
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
					slog.Warn("using mainnet for node")
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
			*cx.Config.Listeners = c.StringSlice("listen")
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
				slog.Warn("set lan", c.Bool("lan"))
				*cx.Config.LAN = c.Bool("lan")
				cx.ActiveNet.DNSSeeds = []config.DNSSeed{}
			} else {
				*cx.Config.LAN = false
			}
		}
		if c.IsSet("controller") {
			*cx.Config.Controller = c.String("controller")
		}
		if c.IsSet("miningaddrs") {
			*cx.Config.MiningAddrs = c.StringSlice("miningaddrs")
		}
		if c.IsSet("minerpass") {
			*cx.Config.MinerPass = c.String("minerpass")
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
		pseudo.Seed(time.Now().UnixNano())
		nonce := fmt.Sprintf("nonce%0x", pseudo.Uint32())
		if cx.Config.UserAgentComments == nil {
			cx.Config.UserAgentComments = &cli.StringSlice{nonce}
		} else {
			*cx.Config.UserAgentComments = append(cli.StringSlice{nonce}, *cx.Config.UserAgentComments...)
		}
		if c.IsSet("uacomment") {
			*cx.Config.UserAgentComments = append(*cx.Config.UserAgentComments,
				c.StringSlice("uacomment")...)
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
		if c.IsSet("save") {
			slog.Info("saving configuration")
			cx.StateCfg.Save = true
		}
		return nil
	}
}
