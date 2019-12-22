package app

import (
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

func // Configure loads and sanitises the configuration from urfave/cli
Configure(cx *conte.Xt, ctx *cli.Context) {
	log.TRACE("configuring pod")
	cx.StateCfg.Save = false
	// theoretically, the configuration should be accessed only when locked
	cfg := cx.Config
	initLogLevel(cfg)
	initDictionary(cfg)
	initDataDir(cfg)
	initTLSStuffs(cfg, cx.StateCfg)
	initConfigFile(cfg)
	initLogDir(cfg)
	initParams(cx)
	initWalletFile(cx)
	initListeners(cx, ctx)
	// Don't add peers from the config file when in regression test mode.
	if ((*cfg.Network)[0] == 'r') && len(*cfg.AddPeers) > 0 {
		*cfg.AddPeers = nil
	}
	normalizeAddresses(cfg)
	setAlgo(cfg)
	setRelayReject(cfg)
	validateDBtype(cfg)
	validateProfilePort(cfg)
	validateBanDuration(cfg)
	validateWhitelists(cfg, cx.StateCfg)
	validatePeerLists(cfg)
	configListener(cfg, cx.ActiveNet)
	validateUsers(cfg)
	configRPC(cfg, cx.ActiveNet)
	validatePolicies(cfg, cx.StateCfg)
	validateOnions(cfg)
	validateMiningStuff(cfg, cx.StateCfg, cx.ActiveNet)
	setDiallers(cfg, cx.StateCfg)
	// if the user set the save flag, or file doesn't exist save the file now
	if cx.StateCfg.Save || !apputil.FileExists(*cx.Config.ConfigFile) {
		log.TRACE("saving configuration")
		save.Pod(cx.Config)
	}
}
