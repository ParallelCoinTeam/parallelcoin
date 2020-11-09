package config

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/wallet"
)

// Configure loads and sanitises the configuration from urfave/cli
func Configure(cx *conte.Xt, commandName string, initial bool) {
	Debug("running Configure", commandName)
	cx.WalletChan = make(chan *wallet.Wallet)
	cx.NodeChan = make(chan *chainrpc.Server)
	// theoretically, the configuration should be accessed only when locked
	// cfg := cx.Config
	Debug("DATADIR", *cx.Config.DataDir)
	initLogLevel(cx.Config)
	Debug("set log level")
	initDictionary(cx.Config)
	initParams(cx)
	initDataDir(cx.Config)
	initTLSStuffs(cx.Config, cx.StateCfg)
	initConfigFile(cx.Config)
	initLogDir(cx.Config)
	initWalletFile(cx)
	initListeners(cx, commandName, initial)
	// Don't add peers from the config file when in regression test mode.
	if ((*cx.Config.Network)[0] == 'r') && len(*cx.Config.AddPeers) > 0 {
		*cx.Config.AddPeers = nil
	}
	normalizeAddresses(cx.Config)
	setRelayReject(cx.Config)
	validateDBtype(cx.Config)
	validateProfilePort(cx.Config)
	validateBanDuration(cx.Config)
	validateWhitelists(cx.Config, cx.StateCfg)
	validatePeerLists(cx.Config)
	configListener(cx.Config, cx.ActiveNet)
	validateUsers(cx.Config)
	configRPC(cx.Config, cx.ActiveNet)
	validatePolicies(cx.Config, cx.StateCfg)
	validateOnions(cx.Config)
	validateMiningStuff(cx.Config, cx.StateCfg, cx.ActiveNet)
	setDiallers(cx.Config, cx.StateCfg)
	// if the user set the save flag, or file doesn't exist save the file now
	if cx.StateCfg.Save || !apputil.FileExists(*cx.Config.ConfigFile) {
		Trace("saving configuration")
		save.Pod(cx.Config)
	}
	if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
		fork.IsTestnet = true
	}
}
