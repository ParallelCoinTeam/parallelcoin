package config

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/stalker-loki/app/slog"
)

// Configure loads and sanitises the configuration from urfave/cli
func Configure(cx *conte.Xt, commandName string, initial bool) {
	slog.Debug("running Configure", commandName)
	cx.WalletChan = make(chan *wallet.Wallet)
	cx.NodeChan = make(chan *chainrpc.Server)
	// theoretically, the configuration should be accessed only when locked
	cfg := cx.Config
	initLogLevel(cfg)
	initDictionary(cfg)
	initParams(cx)
	initDataDir(cfg)
	initTLSStuffs(cfg, cx.StateCfg)
	initConfigFile(cfg)
	initLogDir(cfg)
	initWalletFile(cx)
	initListeners(cx, commandName, initial)
	// Don't add peers from the config file when in regression test mode.
	if ((*cfg.Network)[0] == 'r') && len(*cfg.AddPeers) > 0 {
		*cfg.AddPeers = nil
	}
	normalizeAddresses(cfg)
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
		slog.Trace("saving configuration")
		save.Pod(cx.Config)
	}
}
