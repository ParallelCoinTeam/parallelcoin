package podconfig

// // Configure loads and sanitises the configuration from urfave/cli
// func Configure(cx *pod.State, initial bool) {
// 	commandName := cx.Config.RunningCommand.Name
// 	D.Ln("running Configure", commandName, "DATADIR", cx.Config.DataDir.V())
// 	// initListeners(cx, commandName, initial)
// 	if ((cx.Config.Network.V())[0] == 'r') && cx.Config.AddPeers.Len() > 0 {
// 		cx.Config.AddPeers.Set(nil)
// 	}
// 	// normalizeAddresses(cx.Config)
// 	// configListener(cx.Config, cx.ActiveNet)
// 	// configRPC(cx.Config, cx.ActiveNet)
// 	// validatePolicies(cx.Config, cx.StateCfg)
// 	// validateOnions(cx.Config)
// 	// setDiallers(cx.Config, cx.StateCfg)
// 	if cx.StateCfg.Save || !apputil.FileExists(cx.Config.ConfigFile.V()) {
// 		cx.StateCfg.Save = false
// 		if commandName == "kopach" {
// 			return
// 		}
// 		D.Ln("saving configuration")
// 		var e error
// 		if e = cx.Config.WriteToFile(cx.Config.ConfigFile.V()); E.Chk(e) {
// 		}
// 	}
// 	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
// 		fork.IsTestnet = true
// 	}
// 	// initLogLevel(cx.Config)
// 	// spv.DisableDNSSeed = cx.Config.DisableDNSSeed.True()
// 	// initDictionary(cx.Config)
// 	// initParams(cx)
// 	// initDataDir(cx.Config)
// 	// initTLSStuffs(cx.Config, cx.StateCfg)
// 	// initConfigFile(cx.Config)
// 	// initLogDir(cx.Config)
// 	// initWalletFile(cx)
// 	// Don't add peers from the config file when in regression test mode.
// 	// setRelayReject(cx.Config)
// 	// validateDBtype(cx.Config)
// 	// validateProfilePort(cx.Config)
// 	// validateBanDuration(cx.Config)
// 	// validateWhitelists(cx.Config, cx.StateCfg)
// 	// validatePeerLists(cx.Config)
// 	// validateUsers(cx.Config)
// 	// validateMiningStuff(cx.Config, cx.StateCfg, cx.ActiveNet)
// 	// if the user set the save flag, or file doesn't exist save the file now
// }
