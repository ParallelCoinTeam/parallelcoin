package app

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/fork"
	"github.com/p9c/pod/pkg/opts"
	"io/ioutil"
	prand "math/rand"
	"os"
	"path/filepath"
	"time"
	
	"github.com/p9c/pod/cmd/spv"
	"github.com/p9c/pod/pkg/pipe/serve"
	"github.com/p9c/pod/version"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/pod"
)

func beforeFunc(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		opts.D.Ln("running beforeFunc")
		cx.AppContext = c
		// if user set datadir this is first thing to configure
		if c.IsSet("datadir") {
			cx.Config.DataDir.Set(c.String("datadir"))
			opts.D.Ln("datadir", *cx.Config.DataDir)
		}
		opts.D.Ln(c.IsSet("D"), c.IsSet("datadir"))
		// // propagate datadir path to interrupt for restart handling
		// interrupt.DataDir = cx.DataDir
		// if there is a delaystart requested, pause for 3 seconds
		if c.IsSet("delaystart") {
			time.Sleep(time.Second * 3)
		}
		if c.IsSet("pipelog") {
			opts.D.Ln("pipe logger enabled")
			cx.Config.PipeLog.Set(c.Bool("pipelog"))
			serve.Log(cx.KillAll, fmt.Sprint(os.Args))
		}
		if c.IsSet("walletfile") {
			cx.Config.WalletFile.Set(c.String("walletfile"))
		}
		cx.Config.ConfigFile.Set(filepath.Join(cx.Config.DataDir.V(), opts.PodConfigFilename))
		// we are going to assume the config is not manually misedited
		if apputil.FileExists(cx.Config.ConfigFile.V()) {
			b, e := ioutil.ReadFile(cx.Config.ConfigFile.V())
			if e == nil {
				cx.Config, cx.ConfigMap = podcfg.New()
				e = json.Unmarshal(b, cx.Config)
				if e != nil {
					opts.E.Ln("error unmarshalling config", e)
					// os.Exit(1)
					return e
				}
			} else {
				opts.F.Ln("unexpected error reading configuration file:", e)
				// os.Exit(1)
				return e
			}
		} else {
			cx.Config.ConfigFile.Set("")
			opts.D.Ln("will save config after configuration")
			cx.StateCfg.Save = true
		}
		if c.IsSet("loglevel") {
			opts.T.Ln("set loglevel", c.String("loglevel"))
			cx.Config.LogLevel.Set(c.String("loglevel"))
		}
		log.SetLogLevel(cx.Config.LogLevel.V())
		if cx.Config.PipeLog.False() {
			// if/when running further instances of the same version no reason
			// to print the version message again
			opts.D.Ln("\nrunning", os.Args, version.Get())
		}
		if c.IsSet("network") {
			cx.Config.Network.Set(c.String("network"))
			switch cx.Config.Network.V() {
			case "testnet", "testnet3", "t":
				cx.ActiveNet = &chaincfg.TestNet3Params
				fork.IsTestnet = true
				// fork.HashReps = 3
			case "regtestnet", "regressiontest", "r":
				fork.IsTestnet = true
				cx.ActiveNet = &chaincfg.RegressionTestParams
			case "simnet", "s":
				fork.IsTestnet = true
				cx.ActiveNet = &chaincfg.SimNetParams
			default:
				if cx.Config.Network.V() != "mainnet" &&
					cx.Config.Network.V() != "m" {
					opts.D.Ln("using mainnet for node")
				}
				cx.ActiveNet = &chaincfg.MainNetParams
			}
		}
		if c.IsSet("username") {
			cx.Config.Username.Set(c.String("username"))
		}
		if c.IsSet("password") {
			cx.Config.Password.Set(c.String("password"))
		}
		if c.IsSet("serveruser") {
			cx.Config.ServerUser.Set(c.String("serveruser"))
		}
		if c.IsSet("serverpass") {
			cx.Config.ServerPass.Set(c.String("serverpass"))
		}
		if c.IsSet("limituser") {
			cx.Config.LimitUser.Set(c.String("limituser"))
		}
		if c.IsSet("limitpass") {
			cx.Config.LimitPass.Set(c.String("limitpass"))
		}
		if c.IsSet("rpccert") {
			cx.Config.RPCCert.Set(c.String("rpccert"))
		}
		if c.IsSet("rpckey") {
			cx.Config.RPCKey.Set(c.String("rpckey"))
		}
		if c.IsSet("cafile") {
			cx.Config.CAFile.Set(c.String("cafile"))
		}
		if c.IsSet("clienttls") {
			cx.Config.TLS.Set(c.Bool("clienttls"))
		}
		if c.IsSet("servertls") {
			cx.Config.ServerTLS.Set(c.Bool("servertls"))
		}
		if c.IsSet("tlsskipverify") {
			cx.Config.TLSSkipVerify.Set(c.Bool("tlsskipverify"))
		}
		if c.IsSet("proxy") {
			cx.Config.Proxy.Set(c.String("proxy"))
		}
		if c.IsSet("proxyuser") {
			cx.Config.ProxyUser.Set(c.String("proxyuser"))
		}
		if c.IsSet("proxypass") {
			cx.Config.ProxyPass.Set(c.String("proxypass"))
		}
		if c.IsSet("onion") {
			cx.Config.Onion.Set(c.Bool("onion"))
		}
		if c.IsSet("onionproxy") {
			cx.Config.OnionProxy.Set(c.String("onionproxy"))
		}
		if c.IsSet("onionuser") {
			cx.Config.OnionProxyUser.Set(c.String("onionuser"))
		}
		if c.IsSet("onionpass") {
			cx.Config.OnionProxyPass.Set(c.String("onionpass"))
		}
		if c.IsSet("torisolation") {
			cx.Config.TorIsolation.Set(c.Bool("torisolation"))
		}
		if c.IsSet("addpeer") {
			cx.Config.AddPeers.Set(c.StringSlice("addpeer"))
		}
		if c.IsSet("connect") {
			cx.Config.ConnectPeers.Set(c.StringSlice("connect"))
		}
		if c.IsSet("nolisten") {
			cx.Config.DisableListen.Set(c.Bool("nolisten"))
		}
		if c.IsSet("listen") {
			cx.Config.P2PListeners.Set(c.StringSlice("listen"))
		}
		if c.IsSet("maxpeers") {
			cx.Config.MaxPeers.Set(c.Int("maxpeers"))
		}
		if c.IsSet("nobanning") {
			cx.Config.DisableBanning.Set(c.Bool("nobanning"))
		}
		if c.IsSet("banduration") {
			cx.Config.BanDuration.Set(c.Duration("banduration"))
		}
		if c.IsSet("banthreshold") {
			cx.Config.BanThreshold.Set(c.Int("banthreshold"))
		}
		if c.IsSet("whitelist") {
			cx.Config.Whitelists.Set(c.StringSlice("whitelist"))
		}
		if c.IsSet("rpcconnect") {
			cx.Config.RPCConnect.Set(c.String("rpcconnect"))
		}
		if c.IsSet("rpclisten") {
			cx.Config.RPCListeners.Set(c.StringSlice("rpclisten"))
		}
		if c.IsSet("rpcmaxclients") {
			cx.Config.RPCMaxClients.Set(c.Int("rpcmaxclients"))
		}
		if c.IsSet("rpcmaxwebsockets") {
			cx.Config.RPCMaxWebsockets.Set(c.Int("rpcmaxwebsockets"))
		}
		if c.IsSet("rpcmaxconcurrentreqs") {
			cx.Config.RPCMaxConcurrentReqs.Set(c.Int("rpcmaxconcurrentreqs"))
		}
		if c.IsSet("rpcquirks") {
			cx.Config.RPCQuirks.Set(c.Bool("rpcquirks"))
		}
		if c.IsSet("norpc") {
			cx.Config.DisableRPC.Set(c.Bool("norpc"))
		}
		if c.IsSet("nodnsseed") {
			cx.Config.DisableDNSSeed.Set(c.Bool("nodnsseed"))
			spv.DisableDNSSeed = c.Bool("nodnsseed")
		}
		if c.IsSet("externalip") {
			cx.Config.ExternalIPs.Set(c.StringSlice("externalip"))
		}
		if c.IsSet("addcheckpoint") {
			cx.Config.AddCheckpoints.Set(c.StringSlice("addcheckpoint"))
		}
		if c.IsSet("nocheckpoints") {
			cx.Config.DisableCheckpoints.Set(c.Bool("nocheckpoints"))
		}
		if c.IsSet("dbtype") {
			cx.Config.DbType.Set(c.String("dbtype"))
		}
		if c.IsSet("profile") {
			cx.Config.Profile.Set(c.String("profile"))
		}
		if c.IsSet("cpuprofile") {
			cx.Config.CPUProfile.Set(c.String("cpuprofile"))
		}
		if c.IsSet("upnp") {
			cx.Config.UPNP.Set(c.Bool("upnp"))
		}
		if c.IsSet("minrelaytxfee") {
			cx.Config.MinRelayTxFee.Set(c.Float64("minrelaytxfee"))
		}
		if c.IsSet("limitfreerelay") {
			cx.Config.FreeTxRelayLimit.Set(c.Float64("limitfreerelay"))
		}
		if c.IsSet("norelaypriority") {
			cx.Config.NoRelayPriority.Set(c.Bool("norelaypriority"))
		}
		if c.IsSet("trickleinterval") {
			cx.Config.TrickleInterval.Set(c.Duration("trickleinterval"))
		}
		if c.IsSet("maxorphantx") {
			cx.Config.MaxOrphanTxs.Set(c.Int("maxorphantx"))
		}
		if c.IsSet("generate") {
			cx.Config.Generate.Set(c.Bool("generate"))
		}
		if c.IsSet("genthreads") {
			cx.Config.GenThreads.Set(c.Int("genthreads"))
		}
		if c.IsSet("solo") {
			cx.Config.Solo.Set(c.Bool("solo"))
		}
		if c.IsSet("autoports") {
			cx.Config.AutoPorts.Set(c.Bool("autoports"))
		}
		if c.IsSet("lan") {
			// if LAN is turned on we need to remove the seeds from netparams not on mainnet
			// mainnet is never in lan mode
			// if LAN is turned on it means by default we are on testnet
			cx.ActiveNet = &chaincfg.TestNet3Params
			if cx.ActiveNet.Name != "mainnet" {
				opts.D.Ln("set lan", c.Bool("lan"))
				cx.Config.LAN.Set(c.Bool("lan"))
				cx.ActiveNet.DNSSeeds = []chaincfg.DNSSeed{}
			} else {
				cx.Config.LAN.F()
			}
		}
		// if c.IsSet("controller") {
		//    *cx.Config.Controller.Set(.String("controller"))
		// }
		// if c.IsSet("controllerconnect") {
		//    *cx.Config.ControllerConnect.Set(c.StringSlice("controllerconnect"))
		// }
		if c.IsSet("miningaddrs") {
			cx.Config.MiningAddrs.Set(c.StringSlice("miningaddrs"))
		}
		if c.IsSet("minerpass") {
			cx.Config.MulticastPass.Set(c.String("minerpass"))
			opts.D.Ln("--------- set minerpass", *cx.Config.MulticastPass)
			cx.StateCfg.Save = true
		}
		if c.IsSet("blockminsize") {
			cx.Config.BlockMinSize.Set(c.Int("blockminsize"))
		}
		if c.IsSet("blockmaxsize") {
			cx.Config.BlockMaxSize.Set(c.Int("blockmaxsize"))
		}
		if c.IsSet("blockminweight") {
			cx.Config.BlockMinWeight.Set(c.Int("blockminweight"))
		}
		if c.IsSet("blockmaxweight") {
			cx.Config.BlockMaxWeight.Set(c.Int("blockmaxweight"))
		}
		if c.IsSet("blockprioritysize") {
			cx.Config.BlockPrioritySize.Set(c.Int("blockprioritysize"))
		}
		prand.Seed(time.Now().UnixNano())
		nonce := fmt.Sprintf("nonce%0x", prand.Uint32())
		if cx.Config.UserAgentComments == nil {
			cx.Config.UserAgentComments.Set(cli.StringSlice{nonce})
		} else {
			cx.Config.UserAgentComments.Set(append(cli.StringSlice{nonce}, cx.Config.UserAgentComments.S()...))
		}
		if c.IsSet("uacomment") {
			cx.Config.UserAgentComments.Set(
				append(
					cx.Config.UserAgentComments.S(),
					c.StringSlice("uacomment")...,
				),
			)
		}
		if c.IsSet("nopeerbloomfilters") {
			cx.Config.NoPeerBloomFilters.Set(c.Bool("nopeerbloomfilters"))
		}
		if c.IsSet("nocfilters") {
			cx.Config.NoCFilters.Set(c.Bool("nocfilters"))
		}
		if c.IsSet("sigcachemaxsize") {
			cx.Config.SigCacheMaxSize.Set(c.Int("sigcachemaxsize"))
		}
		if c.IsSet("blocksonly") {
			cx.Config.BlocksOnly.Set(c.Bool("blocksonly"))
		}
		if c.IsSet("notxindex") {
			cx.Config.TxIndex.Set(c.Bool("notxindex"))
		}
		if c.IsSet("noaddrindex") {
			cx.Config.AddrIndex.Set(c.Bool("noaddrindex"))
		}
		if c.IsSet("relaynonstd") {
			cx.Config.RelayNonStd.Set(c.Bool("relaynonstd"))
		}
		if c.IsSet("rejectnonstd") {
			cx.Config.RejectNonStd.Set(c.Bool("rejectnonstd"))
		}
		if c.IsSet("noinitialload") {
			cx.Config.NoInitialLoad.Set(c.Bool("noinitialload"))
		}
		if c.IsSet("walletconnect") {
			cx.Config.Wallet.Set(c.Bool("walletconnect"))
		}
		if c.IsSet("walletserver") {
			cx.Config.WalletServer.Set(c.String("walletserver"))
		}
		if c.IsSet("walletpass") {
			cx.Config.WalletPass.Set(c.String("walletpass"))
		} else {
			// if this is not set, the config will be storing the hash and hashes on save, so we set explicitly to empty
			// as otherwise it would have the hex of the hash of the password here
			cx.Config.WalletPass.Set("")
		}
		if c.IsSet("onetimetlskey") {
			cx.Config.OneTimeTLSKey.Set(c.Bool("onetimetlskey"))
		}
		if c.IsSet("walletrpclisten") {
			cx.Config.WalletRPCListeners.Set(c.StringSlice("walletrpclisten"))
		}
		if c.IsSet("walletrpcmaxclients") {
			cx.Config.WalletRPCMaxClients.Set(c.Int("walletrpcmaxclients"))
		}
		if c.IsSet("walletrpcmaxwebsockets") {
			cx.Config.WalletRPCMaxWebsockets.Set(c.Int("walletrpcmaxwebsockets"))
		}
		if c.IsSet("nodeoff") {
			cx.Config.NodeOff.Set(c.Bool("nodeoff"))
		}
		if c.IsSet("walletoff") {
			cx.Config.WalletOff.Set(c.Bool("walletoff"))
		}
		if c.IsSet("darktheme") {
			cx.Config.DarkTheme.Set(c.Bool("darktheme"))
		}
		if c.IsSet("notty") {
			cx.IsGUI = true
		}
		if c.IsSet("controller") {
			cx.Config.Controller.Set(c.Bool("controller"))
		}
		if c.IsSet("save") {
			opts.I.Ln("saving configuration")
			cx.StateCfg.Save = true
		}
		// // if e = routeable.Discover(); E.Chk(e) {
		// // 	// TODO: this should trigger the display of this lack of internet
		// // }
		// go func() {
		// out:
		// 	for {
		// 		select {
		// 		case <-time.After(time.Second * 10):
		// 			if e = routeable.Discover(); E.Chk(e) {
		// 				// TODO: this should trigger the display of this lack of internet
		// 			}
		// 		case <-cx.KillAll:
		// 			break out
		// 		}
		// 	}
		// }()
		return nil
	}
}
