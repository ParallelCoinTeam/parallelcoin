package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

func // Configure loads and sanitises the configuration from urfave/cli
Configure(cx *conte.Xt) {
	log.TRACE("configuring pod")
	// theoretically, the configuration should be accessed only when locked
	cfg := cx.Config
	st := cx.StateCfg
	initDictionary(cfg)
	initDataDir(cfg)
	initTLSStuffs(cfg, st)
	initConfigFile(cfg)
	initLogDir(cfg)
	initParams(cx)
	initListeners(cx)
	initLogLevel(cfg)
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
	if st.Save || !apputil.FileExists(*cx.Config.ConfigFile) {
		log.TRACE("saving configuration")
		save.Pod(cx.Config)
	}
}
