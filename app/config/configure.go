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
	if cx.StateCfg.Save {
		Debug("save was set")
	}
	initLogLevel(cx.Config)
	Debug("set log level")
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initDictionary(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initParams(cx)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initDataDir(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initTLSStuffs(cx.Config, cx.StateCfg)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initConfigFile(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initLogDir(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initWalletFile(cx)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	initListeners(cx, commandName, initial)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	// Don't add peers from the config file when in regression test mode.
	if ((*cx.Config.Network)[0] == 'r') && len(*cx.Config.AddPeers) > 0 {
		*cx.Config.AddPeers = nil
	}
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	normalizeAddresses(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	setRelayReject(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateDBtype(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateProfilePort(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateBanDuration(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateWhitelists(cx.Config, cx.StateCfg)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validatePeerLists(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	configListener(cx.Config, cx.ActiveNet)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateUsers(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	configRPC(cx.Config, cx.ActiveNet)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validatePolicies(cx.Config, cx.StateCfg)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateOnions(cx.Config)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	validateMiningStuff(cx.Config, cx.StateCfg, cx.ActiveNet)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	setDiallers(cx.Config, cx.StateCfg)
	if cx.StateCfg.Save {
		Debug("save was set", commandName, initial)
	}
	// if the user set the save flag, or file doesn't exist save the file now
	if cx.StateCfg.Save || !apputil.FileExists(*cx.Config.ConfigFile) {
		Trace("saving configuration")
		save.Pod(cx.Config)
		cx.StateCfg.Save = false
	}
	if cx.ActiveNet.Name == netparams.TestNet3Params.Name {
		fork.IsTestnet = true
	}
}
