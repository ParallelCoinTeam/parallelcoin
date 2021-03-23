package config

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/spv"
	"github.com/p9c/pod/pkg/blockchain/chaincfg"
	"github.com/p9c/pod/pkg/blockchain/fork"
)

// Configure loads and sanitises the configuration from urfave/cli
func Configure(cx *conte.Xt, commandName string, initial bool) {
	initLogLevel(cx.Config)
	D.Ln("running Configure", commandName, *cx.Config.WalletPass)
	D.Ln("DATADIR", *cx.Config.DataDir)
	D.Ln("set log level")
	spv.DisableDNSSeed = *cx.Config.DisableDNSSeed
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
		cx.StateCfg.Save = false
		if commandName == "kopach" {
			return
		}
		D.Ln("saving configuration")
		save.Pod(cx.Config)
	}
	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
		fork.IsTestnet = true
	}
}
