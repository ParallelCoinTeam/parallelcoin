package podconfig

import (
	"github.com/p9c/monorepo/pkg/apputil"
	"github.com/p9c/monorepo/pkg/chaincfg"
	"github.com/p9c/monorepo/pkg/fork"
	"github.com/p9c/monorepo/pkg/pod"
)

// Configure loads and sanitises the configuration from urfave/cli
func Configure(cx *pod.State, initial bool) {
	commandName := cx.Config.RunningCommand.Name
	initLogLevel(cx.Config)
	D.Ln("running Configure", commandName, "DATADIR", cx.Config.DataDir.V())
	// spv.DisableDNSSeed = cx.Config.DisableDNSSeed.True()
	initDictionary(cx.Config)
	initParams(cx)
	initDataDir(cx.Config)
	initTLSStuffs(cx.Config, cx.StateCfg)
	initConfigFile(cx.Config)
	initLogDir(cx.Config)
	initWalletFile(cx)
	initListeners(cx, commandName, initial)
	// Don't add peers from the config file when in regression test mode.
	if ((cx.Config.Network.V())[0] == 'r') && cx.Config.AddPeers.Len() > 0 {
		cx.Config.AddPeers.Set(nil)
	}
	normalizeAddresses(cx.Config)
	setRelayReject(cx.Config)
	validateDBtype(cx.Config)
	validateProfilePort(cx.Config)
	// validateBanDuration(cx.Config)
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
	if cx.StateCfg.Save || !apputil.FileExists(cx.Config.ConfigFile.V()) {
		cx.StateCfg.Save = false
		if commandName == "kopach" {
			return
		}
		D.Ln("saving configuration")
		var e error
		if e = cx.Config.WriteToFile(cx.Config.ConfigFile.V()); E.Chk(e) {
		}
	}
	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
		fork.IsTestnet = true
	}
}
