package podconfig

const (
	appName           = "pod"
	confExt           = ".json"
	appLanguage       = "en"
	podConfigFilename = appName + confExt
	PARSER            = "json"
)

var funcName = "loadConfig"

// func initDictionary(cfg *opts.Config) {
// 	if cfg.Locale == nil || cfg.Locale.V() == "" {
// 		cfg.Locale.Set(Lang("en"))
// 	}
// 	T.Ln("lang set to", *cfg.Locale)
// }

// func validatePort(port string) bool {
// 	var e error
// 	var p int64
// 	if p, e = strconv.ParseInt(port, 10, 32); E.Chk(e) {
// 		return false
// 	}
// 	if p < 1024 || p > 65535 {
// 		return false
// 	}
// 	return true
// }

// func initListeners(cx *pod.State, commandName string, initial bool) {
// 	cfg := cx.Config
// 	var e error
// 	var fP int
// 	if fP, e = GetFreePort(); E.Chk(e) {
// 	}
// 	// if cfg.AutoListen.True() {
// 	// 	_, allAddresses := routeable.GetAddressesAndInterfaces()
// 	// 	p2pAddresses := cli.StringSlice{}
// 	// 	for addr := range allAddresses {
// 	// 		p2pAddresses = append(p2pAddresses, net.JoinHostPort(addr, cx.ActiveNet.DefaultPort))
// 	// 	}
// 	// 	cfg.P2PConnect.Set(p2pAddresses)
// 	// }
// 	// if cfg.P2PListeners.Len() < 1 && !cfg.DisableListen.True() && cfg.ConnectPeers.Len() < 1 {
// 	// 	cfg.P2PListeners.Set([]string{fmt.Sprintf("0.0.0.0:" + cx.ActiveNet.DefaultPort)})
// 	// 	cx.StateCfg.Save = true
// 	// 	D.Ln("P2PListeners")
// 	// }
// 	// if cfg.RPCListeners.Len() < 1 {
// 	// 	address := fmt.Sprintf("127.0.0.1:%s", cx.ActiveNet.RPCClientPort)
// 	// 	cfg.RPCListeners.Set([]string{address})
// 	// 	cfg.RPCConnect.Set(fmt.Sprintf("127.0.0.1:%s", cx.ActiveNet.RPCClientPort))
// 	// 	D.Ln("setting save flag because rpc listeners is empty and rpc is not disabled")
// 	// 	cx.StateCfg.Save = true
// 	// 	D.Ln("RPCListeners")
// 	// }
// 	// if cfg.WalletRPCListeners.Len() < 1 && !cfg.DisableRPC.True() {
// 	// 	address := fmt.Sprintf("127.0.0.1:" + cx.ActiveNet.WalletRPCServerPort)
// 	// 	cfg.WalletRPCListeners.Set([]string{address})
// 	// 	cfg.WalletServer.Set(address)
// 	// 	D.Ln(
// 	// 		"setting save flag because wallet rpc listeners is empty and" +
// 	// 			" rpc is not disabled",
// 	// 	)
// 	// 	cx.StateCfg.Save = true
// 	// 	D.Ln("WalletRPCListeners")
// 	// }
// 	// if cx.Config.AutoPorts.True() || !initial {
// 	// 	if fP, e = GetFreePort(); E.Chk(e) {
// 	// 	}
// 	// 	cfg.P2PListeners.Set([]string{"0.0.0.0:" + fmt.Sprint(fP)})
// 	// 	if fP, e = GetFreePort(); E.Chk(e) {
// 	// 	}
// 	// 	cfg.RPCListeners.Set([]string{"127.0.0.1:" + fmt.Sprint(fP)})
// 	// 	if fP, e = GetFreePort(); E.Chk(e) {
// 	// 	}
// 	// 	cfg.WalletRPCListeners.Set([]string{"127.0.0.1:" + fmt.Sprint(fP)})
// 	// 	cx.StateCfg.Save = true
// 	// 	D.Ln("autoports")
// 	// } else {
// 		// sanitize user input and set auto on any that fail
// 		// l := cfg.P2PListeners.S()
// 		// r := cfg.RPCListeners.S()
// 		// w := cfg.WalletRPCListeners.S()
// 		// for i := range l {
// 		// 	if _, p, e := net.SplitHostPort((l)[i]); !E.Chk(e) {
// 		// 		if !validatePort(p) {
// 		// 			if fP, e = GetFreePort(); E.Chk(e) {
// 		// 			}
// 		// 			l[i] = "0.0.0.0:" + fmt.Sprint(fP)
// 		// 			cx.StateCfg.Save = true
// 		// 			D.Ln("port not validate P2PListeners")
// 		// 		}
// 		// 	}
// 		// }
// 		// for i := range r {
// 		// 	if _, p, e := net.SplitHostPort((r)[i]); !E.Chk(e) {
// 		// 		if !validatePort(p) {
// 		// 			if fP, e = GetFreePort(); E.Chk(e) {
// 		// 			}
// 		// 			r[i] = "127.0.0.1:" + fmt.Sprint(fP)
// 		// 			cx.StateCfg.Save = true
// 		// 			D.Ln("port not validate RPCListeners")
// 		// 		}
// 		// 	}
// 		// }
// 		// for i := range w {
// 		// 	if _, p, e := net.SplitHostPort((w)[i]); !E.Chk(e) {
// 		// 		if !validatePort(p) {
// 		// 			if fP, e = GetFreePort(); E.Chk(e) {
// 		// 			}
// 		// 			w[i] = "127.0.0.1:" + fmt.Sprint(fP)
// 		// 			cx.StateCfg.Save = true
// 		// 			D.Ln("port not validate WalletRPCListeners")
// 		// 		}
// 		// 	}
// 		// }
// 	}
// 	// if cfg.LAN.True() && cx.ActiveNet.Name != "mainnet" {
// 	// 	cfg.DisableDNSSeed.T()
// 	// }
// 	// if cfg.WalletRPCListeners.Len() > 0 {
// 	// 	cfg.WalletServer.Set(cfg.WalletRPCListeners.S()[0])
// 	// }
// }

// // GetFreePort asks the kernel for free open ports that are ready to use.
// func GetFreePort() (int, error) {
// 	var port int
// 	addr, e := net.ResolveTCPAddr("tcp", "localhost:0")
// 	if e != nil {
// 		return 0, e
// 	}
// 	var l *net.TCPListener
// 	l, e = net.ListenTCP("tcp", addr)
// 	if e != nil {
// 		return 0, e
// 	}
// 	defer func() {
// 		if e := l.Close(); E.Chk(e) {
// 		}
// 	}()
// 	port = l.Addr().(*net.TCPAddr).Port
// 	return port, nil
// }

// func normalizeAddresses(cfg *opts.Config) {
// 	T.Ln("normalising addresses")
// 	port := fmt.Sprint(constant.DefaultP2PPort)
// 	nrm := normalize.StringSliceAddresses
// 	nrm(cfg.AddPeers.V(), port)
// 	nrm(cfg.ConnectPeers.V(), port)
// 	nrm(cfg.P2PListeners.V(), port)
// 	nrm(cfg.Whitelists.V(), port)
// 	// nrm(cfg.RPCListeners, port)
// }

// func configListener(cfg *opts.Config, params *chaincfg.Params) {
// 	// --proxy or --connect without --listen disables listening.
// 	T.Ln("checking proxy/connect for disabling listening")
// 	if (cfg.ProxyAddress.V() != "" ||
// 		cfg.ConnectPeers.Len() > 0) &&
// 		cfg.P2PListeners.Len() == 0 {
// 		cfg.DisableListen.T()
// 	}
// 	// Add the default listener if none were specified. The default listener is all
// 	// addresses on the listen port for the network we are to connect to.
// 	T.Ln("checking if listener was set")
// 	if cfg.P2PListeners.Len() == 0 {
// 		cfg.P2PListeners.Set([]string{"0.0.0.0:" + params.DefaultPort})
// 	}
// }

// func configRPC(cfg *opts.Config, params *chaincfg.Params) {
// 	// // The RPC server is disabled if no username or password is provided.
// 	// // T.Ln("checking rpc server has a login enabled")
// 	// // if (cfg.Username.Empty() || cfg.Password.Empty()) &&
// 	// // 	(cfg.LimitUser.Empty() || cfg.LimitPass.Empty()) {
// 	// // 	cfg.DisableRPC.T()
// 	// // }
// 	// // if cfg.DisableRPC.True() {
// 	// // 	T.Ln("RPC service is disabled")
// 	// // }
// 	// T.Ln("checking rpc server has listeners set")
// 	// if cfg.DisableRPC.False() && cfg.RPCListeners.Len() == 0 {
// 	// 	D.Ln("looking up default listener")
// 	// 	addrs, e := net.LookupHost(constant.DefaultRPCListener)
// 	// 	if e != nil {
// 	// 		E.Ln(e)
// 	// 		// os.Exit(1)
// 	// 	}
// 	// 	tmp := make([]string, 0, len(addrs))
// 	// 	D.Ln("setting listeners")
// 	// 	for _, addr := range addrs {
// 	// 		tmp = append(tmp, addr)
// 	// 		addr = net.JoinHostPort(addr, params.RPCClientPort)
// 	// 	}
// 	// 	cfg.RPCListeners.Set(tmp)
// 	// }
// 	// // T.Ln("checking rpc max concurrent requests")
// 	// // if cfg.RPCMaxConcurrentReqs.V() < 0 {
// 	// // 	str := "%s: The rpcmaxwebsocketconcurrentrequests opt may not be" +
// 	// // 		" less than 0 -- parsed [%d]"
// 	// // 	e := fmt.Errorf(str, funcName, cfg.RPCMaxConcurrentReqs.V())
// 	// // 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// // 	// os.Exit(1)
// 	// // }
// 	// T.Ln("checking rpc listener addresses")
// 	// nrms := normalize.Addresses
// 	// // Add default port to all added peer addresses if needed and remove duplicate addresses.
// 	// cfg.AddPeers.Set(nrms(cfg.AddPeers.S(), params.DefaultPort))
// 	// cfg.ConnectPeers.Set(nrms(cfg.ConnectPeers.S(), params.DefaultPort))
// }

// func validatePolicies(cfg *opts.Config, stateConfig *state.Config) {
// 	// var e error
// 	// // Validate the the minrelaytxfee.
// 	// T.Ln("checking min relay tx fee")
// 	// stateConfig.ActiveMinRelayTxFee, e = amt.NewAmount(cfg.MinRelayTxFee.V())
// 	// if e != nil {
// 	// 	E.Ln(e)
// 	// 	str := "%s: invalid minrelaytxfee: %v"
// 	// 	e = fmt.Errorf(str, funcName, e)
// 	// 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// }
// 	// // Limit the max block size to a sane value.
// 	// T.Ln("checking max block size")
// 	// if cfg.BlockMaxSize.V() < constant.BlockMaxSizeMin ||
// 	// 	cfg.BlockMaxSize.V() > constant.BlockMaxSizeMax {
// 	// 	str := "%s: The blockmaxsize opt must be in between %d and %d -- parsed [%d]"
// 	// 	e = fmt.Errorf(
// 	// 		str, funcName, constant.BlockMaxSizeMin,
// 	// 		constant.BlockMaxSizeMax, cfg.BlockMaxSize.V(),
// 	// 	)
// 	// 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// }
// 	// Limit the max block weight to a sane value.
// 	// T.Ln("checking max block weight")
// 	// if cfg.BlockMaxWeight.V() < constant.BlockMaxWeightMin ||
// 	// 	cfg.BlockMaxWeight.V() > constant.BlockMaxWeightMax {
// 	// 	str := "%s: The blockmaxweight opt must be in between %d and %d -- parsed [%d]"
// 	// 	e = fmt.Errorf(
// 	// 		str, funcName, constant.BlockMaxWeightMin,
// 	// 		constant.BlockMaxWeightMax, cfg.BlockMaxWeight.V(),
// 	// 	)
// 	// 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// }
// 	// Limit the max orphan count to a sane vlue.
// 	// T.Ln("checking max orphan limit")
// 	// if cfg.MaxOrphanTxs.V() < 0 {
// 	// 	str := "%s: The maxorphantx opt may not be less than 0 -- parsed [%d]"
// 	// 	e = fmt.Errorf(str, funcName, cfg.MaxOrphanTxs.V())
// 	// 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// }
// 	// // Limit the block priority and minimum block txsizes to max block size.
// 	// T.Ln("validating block priority and minimum size/weight")
// 	// cfg.BlockPrioritySize.Set(
// 	// 	int(
// 	// 		apputil.MinUint32(
// 	// 			uint32(cfg.BlockPrioritySize.V()),
// 	// 			uint32(cfg.BlockMaxSize.V()),
// 	// 		),
// 	// 	),
// 	// )
// 	// cfg.BlockMinSize.Set(
// 	// 	int(
// 	// 		apputil.MinUint32(
// 	// 			uint32(cfg.BlockMinSize.V()),
// 	// 			uint32(cfg.BlockMaxSize.V()),
// 	// 		),
// 	// 	),
// 	// )
// 	// cfg.BlockMinWeight.Set(
// 	// 	int(
// 	// 		apputil.MinUint32(
// 	// 			uint32(cfg.BlockMinWeight.V()),
// 	// 			uint32(cfg.BlockMaxWeight.V()),
// 	// 		),
// 	// 	),
// 	// )
// 	// switch {
// 	// // If the max block size isn't set, but the max weight is, then we'll set the
// 	// // limit for the max block size to a safe limit so weight takes precedence.
// 	// case cfg.BlockMaxSize.V() == constant.DefaultBlockMaxSize &&
// 	// 	cfg.BlockMaxWeight.V() != constant.DefaultBlockMaxWeight:
// 	// 	cfg.BlockMaxSize.Set(blockchain.MaxBlockBaseSize - 1000)
// 	// 	// If the max block weight isn't set, but the block size is, then we'll scale
// 	// 	// the set weight accordingly based on the max block size value.
// 	// case cfg.BlockMaxSize.V() != constant.DefaultBlockMaxSize &&
// 	// 	cfg.BlockMaxWeight.V() == constant.DefaultBlockMaxWeight:
// 	// 	cfg.BlockMaxWeight.Set(cfg.BlockMaxSize.V() * blockchain.WitnessScaleFactor)
// 	// }
// 	// Look for illegal characters in the user agent comments.
// 	// T.Ln("checking user agent comments", cfg.UserAgentComments)
// 	// for _, uaComment := range cfg.UserAgentComments.S() {
// 	// 	if strings.ContainsAny(uaComment, "/:()") {
// 	// 		e = fmt.Errorf(
// 	// 			"%s: The following characters must not "+
// 	// 				"appear in user agent comments: '/', ':', '(', ')'",
// 	// 			funcName,
// 	// 		)
// 	// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 	// 	}
// 	// }
// 	// // Chk the checkpoints for syntax errors.
// 	// T.Ln("checking the checkpoints")
// 	// stateConfig.AddedCheckpoints, e = ParseCheckpoints(
// 	// 	cfg.AddCheckpoints.S(),
// 	// )
// 	// if e != nil {
// 	// 	E.Ln(e)
// 	// 	str := "%s: err parsing checkpoints: %v"
// 	// 	e = fmt.Errorf(str, funcName, e)
// 	// 	_, _ = fmt.Fprintln(os.Stderr, e)
// 	// }
// }

// func validateOnions(cfg *opts.Config) {
// 	// --onionproxy and not --onion are contradictory
// 	// TODO: this is kinda stupid hm? switch *and* toggle by presence of flag value, one should be enough
// 	if cfg.OnionEnabled.True() && !cfg.OnionProxyAddress.Empty() {
// 		E.Ln("onion enabled but no onionproxy has been configured")
// 		T.Ln("halting to avoid exposing IP address")
// 	}
// 	// Tor stream isolation requires either proxy or onion proxy to be set.
// 	if cfg.TorIsolation.True() &&
// 		cfg.ProxyAddress.Empty() &&
// 		cfg.OnionProxyAddress.Empty() {
// 		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
// 		e := fmt.Errorf(str, funcName)
// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 		// os.Exit(1)
// 	}
// 	if cfg.OnionEnabled.False() {
// 		cfg.OnionProxyAddress.Set("")
// 	}
//
// }

// func setDiallers(cfg *opts.Config, stateConfig *state.Config) {
// 	// Setup dial and DNS resolution (lookup) functions depending on the specified
// 	// options. The default is to use the standard net.DialTimeout function as well
// 	// as the system DNS resolver. When a proxy is specified, the dial function is
// 	// set to the proxy specific dial function and the lookup is set to use tor
// 	// (unless --noonion is specified in which case the system DNS resolver is
// 	// used).
// 	T.Ln("setting network dialer and lookup")
// 	stateConfig.Dial = net.DialTimeout
// 	stateConfig.Lookup = net.LookupIP
// 	var e error
// 	if !cfg.ProxyAddress.Empty() {
// 		T.Ln("we are loading a proxy!")
// 		_, _, e = net.SplitHostPort(cfg.ProxyAddress.V())
// 		if e != nil {
// 			E.Ln(e)
// 			str := "%s: Proxy address '%s' is invalid: %v"
// 			e = fmt.Errorf(str, funcName, cfg.ProxyAddress.V(), e)
// 			fmt.Fprintln(os.Stderr, e)
// 			// os.Exit(1)
// 		}
// 		// Tor isolation flag means proxy credentials will be overridden unless there is
// 		// also an onion proxy configured in which case that one will be overridden.
// 		torIsolation := false
// 		if cfg.TorIsolation.True() && cfg.OnionProxyAddress.Empty() &&
// 			(!cfg.ProxyUser.Empty() || !cfg.ProxyPass.Empty()) {
// 			torIsolation = true
// 			W.Ln(
// 				"Tor isolation set -- overriding specified" +
// 					" proxy user credentials",
// 			)
// 		}
// 		proxy := &socks.Proxy{
// 			Addr:         cfg.ProxyAddress.V(),
// 			Username:     cfg.ProxyUser.V(),
// 			Password:     cfg.ProxyPass.V(),
// 			TorIsolation: torIsolation,
// 		}
// 		stateConfig.Dial = proxy.DialTimeout
// 		// Treat the proxy as tor and perform DNS resolution through it unless the
// 		// --noonion flag is set or there is an onion-specific proxy configured.
// 		if cfg.OnionEnabled.True() &&
// 			cfg.OnionProxyAddress.Empty() {
// 			stateConfig.Lookup = func(host string) ([]net.IP, error) {
// 				return connmgr.TorLookupIP(host, cfg.ProxyAddress.V())
// 			}
// 		}
// 	}
// 	// Setup onion address dial function depending on the specified options. The
// 	// default is to use the same dial function selected above. However, when an
// 	// onion-specific proxy is specified, the onion address dial function is set to
// 	// use the onion-specific proxy while leaving the normal dial function as
// 	// selected above. This allows .onion address traffic to be routed through a
// 	// different proxy than normal traffic.
// 	T.Ln("setting up tor proxy if enabled")
// 	if !cfg.OnionProxyAddress.Empty() {
// 		if _, _, e = net.SplitHostPort(cfg.OnionProxyAddress.V()); E.Chk(e) {
// 			e = fmt.Errorf("%s: Onion proxy address '%s' is invalid: %v",
// 				funcName, cfg.OnionProxyAddress.V(), e,
// 			)
// 			// _, _ = fmt.Fprintln(os.Stderr, e)
// 		}
// 		// Tor isolation flag means onion proxy credentials will be overridden.
// 		if cfg.TorIsolation.True() &&
// 			(!cfg.OnionProxyUser.Empty() || !cfg.OnionProxyPass.Empty()) {
// 			W.Ln(
// 				"Tor isolation set - overriding specified onionproxy user" +
// 					" credentials",
// 			)
// 		}
// 	}
// 	T.Ln("setting onion dialer")
// 	stateConfig.Oniondial =
// 		func(network, addr string, timeout time.Duration) (net.Conn, error) {
// 			proxy := &socks.Proxy{
// 				Addr:         cfg.OnionProxyAddress.V(),
// 				Username:     cfg.OnionProxyUser.V(),
// 				Password:     cfg.OnionProxyPass.V(),
// 				TorIsolation: cfg.TorIsolation.True(),
// 			}
// 			return proxy.DialTimeout(network, addr, timeout)
// 		}
//
// 	// When configured in bridge mode (both --onion and --proxy are configured), it
// 	// means that the proxy configured by --proxy is not a tor proxy, so override
// 	// the DNS resolution to use the onion-specific proxy.
// 	T.Ln("setting proxy lookup")
// 	if !cfg.ProxyAddress.Empty() {
// 		stateConfig.Lookup = func(host string) ([]net.IP, error) {
// 			return connmgr.TorLookupIP(host, cfg.OnionProxyAddress.V())
// 		}
// 	} else {
// 		stateConfig.Oniondial = stateConfig.Dial
// 	}
// 	// Specifying --noonion means the onion address dial function results in an error.
// 	if cfg.OnionEnabled.False() {
// 		stateConfig.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
// 			return nil, errors.New("tor has been disabled")
// 		}
// 	}
// }

// func validateMiningStuff(
// 	cfg *opts.Config, state *state.Config,
// 	params *chaincfg.Params,
// ) {
// 	if state == nil {
// 		panic("state is nil")
// 	}
// 	// // Chk mining addresses are valid and saved parsed versions.
// 	// T.Ln("checking mining addresses")
// 	// aml := 99
// 	// if cfg.MiningAddrs != nil {
// 	// 	aml = cfg.MiningAddrs.Len()
// 	// } else {
// 	// 	D.Ln("MiningAddrs is nil")
// 	// 	return
// 	// }
// 	// state.ActiveMiningAddrs = make([]btcaddr.Address, 0, aml)
// 	// for _, strAddr := range cfg.MiningAddrs.S() {
// 	// 	addr, e := btcaddr.Decode(strAddr, params)
// 	// 	if e != nil {
// 	// 		E.Ln(e)
// 	// 		str := "%s: mining address '%s' failed to decode: %v"
// 	// 		e = fmt.Errorf(str, funcName, strAddr, e)
// 	// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 	// 		// os.Exit(1)
// 	// 		continue
// 	// 	}
// 	// 	if !addr.IsForNet(params) {
// 	// 		str := "%s: mining address '%s' is on the wrong network"
// 	// 		e := fmt.Errorf(str, funcName, strAddr)
// 	// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 	// 		// os.Exit(1)
// 	// 		continue
// 	// 	}
// 	// 	state.ActiveMiningAddrs = append(state.ActiveMiningAddrs, addr)
// 	// }
// 	if cfg.MulticastPass.Empty() {
// 		D.Ln("--------------- generating new miner key")
// 		cfg.MulticastPass.Set(hex.EncodeToString(forkhash.Argon2i(cfg.MulticastPass.Bytes())))
// 		state.ActiveMinerKey = cfg.MulticastPass.Bytes()
// 	}
// }

// func validateUsers(cfg *opts.Config) {
// 	// Chk to make sure limited and admin users don't have the same username
// 	T.Ln("checking admin and limited username is different")
// 	if !cfg.Username.Empty() &&
// 		cfg.Username.V() == cfg.LimitUser.V() {
// 		str := "%s: --username and --limituser must not specify the same username"
// 		e := fmt.Errorf(str, funcName)
// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 	}
// 	// Chk to make sure limited and admin users don't have the same password
// 	T.Ln("checking limited and admin passwords are not the same")
// 	if !cfg.Password.Empty() &&
// 		cfg.Password.V() == cfg.LimitPass.V() {
// 		str := "%s: --password and --limitpass must not specify the same password"
// 		e := fmt.Errorf(str, funcName)
// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 		// os.Exit(1)
// 	}
// }

// func initDataDir(cfg *opts.Config) {
// 	if cfg.DataDir == nil || cfg.DataDir.V() == "" {
// 		D.Ln("setting default data dir")
// 		cfg.DataDir.Set(appdata.Dir("pod", false))
// 	}
// 	T.Ln("datadir set to", *cfg.DataDir)
// }

// func initWalletFile(cx *pod.State) {
// 	if cx.Config.WalletFile == nil || cx.Config.WalletFile.V() == "" {
// 		cx.Config.WalletFile.Set(filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, constant.DbName))
// 	}
// 	T.Ln("wallet file set to", *cx.Config.WalletFile, *cx.Config.Network)
// }

// func initConfigFile(cfg *opts.Config) {
// 	if cfg.ConfigFile.V() == "" {
// 		cfg.ConfigFile.Set(filepath.Join(cfg.DataDir.V(), podConfigFilename))
// 	}
// 	T.Ln("using config file:", *cfg.ConfigFile)
// }

// func initLogDir(cfg *opts.Config) {
// 	if cfg.LogDir.V() != "" {
// 		// logi.L.SetLogPaths(*cfg.LogDir, "pod")
// 		interrupt.AddHandler(
// 			func() {
// 				D.Ln("initLogDir interrupt")
// 				// _ = logi.L.LogFileHandle.Close()
// 			},
// 		)
// 	}
// }

// func initParams(cx *pod.State) {
// 	network := "mainnet"
// 	if cx.Config.Network != nil {
// 		network = cx.Config.Network.V()
// 	}
// 	switch network {
// 	case "testnet", "testnet3", "t":
// 		T.Ln("on testnet")
// 		cx.ActiveNet = &chaincfg.TestNet3Params
// 		fork.IsTestnet = true
// 	case "regtestnet", "regressiontest", "r":
// 		T.Ln("on regression testnet")
// 		cx.ActiveNet = &chaincfg.RegressionTestParams
// 	case "simnet", "s":
// 		T.Ln("on simnet")
// 		cx.ActiveNet = &chaincfg.SimNetParams
// 	default:
// 		if network != "mainnet" && network != "m" {
// 			D.Ln("using mainnet for node")
// 		}
// 		T.Ln("on mainnet")
// 		cx.ActiveNet = &chaincfg.MainNetParams
// 	}
// }

// func initTLSStuffs(cfg *opts.Config, st *state.Config) {
// 	isNew := false
// 	if cfg.RPCCert.V() == "" {
// 		cfg.RPCCert.Set(filepath.Join(cfg.DataDir.V(), "rpc.cert"))
// 		D.Ln("setting save flag because rpc cert path was not set")
// 		st.Save = true
// 		isNew = true
// 	}
// 	if cfg.RPCKey.V() == "" {
// 		cfg.RPCKey.Set(filepath.Join(cfg.DataDir.V(), "rpc.key"))
// 		D.Ln("setting save flag because rpc key path was not set")
// 		st.Save = true
// 		isNew = true
// 	}
// 	if cfg.CAFile.V() == "" {
// 		cfg.CAFile.Set(filepath.Join(cfg.DataDir.V(), "ca.cert"))
// 		D.Ln("setting save flag because CA cert path was not set")
// 		st.Save = true
// 		isNew = true
// 	}
// 	if isNew {
// 		// Now is the best time to make the certs
// 		I.Ln("generating TLS certificates")
// 		// Create directories for cert and key files if they do not yet exist.
// 		D.Ln("rpc tls ", *cfg.RPCCert, " ", *cfg.RPCKey)
// 		certDir, _ := filepath.Split(cfg.RPCCert.V())
// 		keyDir, _ := filepath.Split(cfg.RPCKey.V())
// 		var e error
// 		e = os.MkdirAll(certDir, 0700)
// 		if e != nil {
// 			E.Ln(e)
// 			return
// 		}
// 		e = os.MkdirAll(keyDir, 0700)
// 		if e != nil {
// 			E.Ln(e)
// 			return
// 		}
// 		// Generate cert pair.
// 		org := "pod/wallet autogenerated cert"
// 		validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
// 		cert, key, e := util.NewTLSCertPair(org, validUntil, nil)
// 		if e != nil {
// 			E.Ln(e)
// 			return
// 		}
// 		_, e = tls.X509KeyPair(cert, key)
// 		if e != nil {
// 			E.Ln(e)
// 			return
// 		}
// 		// Write cert and (potentially) the key files.
// 		e = ioutil.WriteFile(cfg.RPCCert.V(), cert, 0600)
// 		if e != nil {
// 			rmErr := os.Remove(cfg.RPCCert.V())
// 			if rmErr != nil {
// 				E.Ln("cannot remove written certificates:", rmErr)
// 			}
// 			return
// 		}
// 		e = ioutil.WriteFile(cfg.CAFile.V(), cert, 0600)
// 		if e != nil {
// 			rmErr := os.Remove(cfg.RPCCert.V())
// 			if rmErr != nil {
// 				E.Ln("cannot remove written certificates:", rmErr)
// 			}
// 			return
// 		}
// 		e = ioutil.WriteFile(cfg.RPCKey.V(), key, 0600)
// 		if e != nil {
// 			E.Ln(e)
// 			rmErr := os.Remove(cfg.RPCCert.V())
// 			if rmErr != nil {
// 				E.Ln("cannot remove written certificates:", rmErr)
// 			}
// 			rmErr = os.Remove(cfg.CAFile.V())
// 			if rmErr != nil {
// 				E.Ln("cannot remove written certificates:", rmErr)
// 			}
// 			return
// 		}
// 		I.Ln("done generating TLS certificates")
// 		return
// 	}
// }

// func initLogLevel(cfg *opts.Config) {
// 	loglevel := cfg.LogLevel.V()
// 	switch loglevel {
// 	case "trace", "debug", "info", "warn", "error", "fatal", "off":
// 		D.Ln("log level", loglevel)
// 	default:
// 		E.Ln("unrecognised loglevel", loglevel, "setting default info")
// 		cfg.LogLevel.Set("info")
// 	}
// 	log.SetLogLevel(cfg.LogLevel.V())
// }

// func setRelayReject(cfg *opts.Config) {
// 	relayNonStd := *cfg.RelayNonStd
// 	switch {
// 	case cfg.RelayNonStd.True() && cfg.RejectNonStd.True():
// 		errf := "%s: rejectnonstd and relaynonstd cannot be used together" +
// 			" -- choose only one, leaving neither activated"
// 		E.Ln(errf, funcName)
// 		// just leave both false
// 		cfg.RelayNonStd.F()
// 		cfg.RejectNonStd.F()
// 	case cfg.RejectNonStd.True():
// 		relayNonStd.F()
// 	case cfg.RelayNonStd.True():
// 		relayNonStd.F()
// 	}
// 	*cfg.RelayNonStd = relayNonStd
// }
//
// func validateDBtype(cfg *opts.Config) {
// 	// Validate database type.
// 	T.Ln("validating database type")
// 	if !ValidDbType(cfg.DbType.V()) {
// 		str := "%s: The specified database type [%v] is invalid -- " +
// 			"supported types %v"
// 		e := fmt.Errorf(str, funcName, *cfg.DbType, KnownDbTypes)
// 		E.Ln(funcName, e)
// 		// set to default
// 		cfg.DbType.Set(KnownDbTypes[0])
// 	}
// }

// func validateProfilePort(cfg *opts.Config) {
// 	// Validate profile port number
// 	T.Ln("validating profile port number")
// 	if cfg.Profile.V() != "" {
// 		profilePort, e := strconv.Atoi(cfg.Profile.V())
// 		if e != nil || profilePort < 1024 || profilePort > 65535 {
// 			str := "%s: The profile port must be between 1024 and 65535"
// 			e = fmt.Errorf(str, funcName)
// 			E.Ln(funcName, e)
// 			cfg.Profile.Set("")
// 		}
// 	}
// }

// func validateBanDuration(cfg *opts.Config) {
// 	// Don't allow ban durations that are too short.
// 	T.Ln("validating ban duration")
// 	if cfg.BanDuration.V() < time.Second {
// 		e := fmt.Errorf(
// 			"%s: The banduration opt may not be less than 1s -- parsed [%v]",
// 			funcName, *cfg.BanDuration,
// 		)
// 		I.Ln(funcName, e)
// 		cfg.BanDuration.Set(constant.DefaultBanDuration)
// 	}
// }
//
// func validateWhitelists(cfg *opts.Config, st *state.Config) {
// 	// Validate any given whitelisted IP addresses and networks.
// 	T.Ln("validating whitelists")
// 	if cfg.Whitelists.Len() > 0 {
// 		var ip net.IP
// 		st.ActiveWhitelists = make([]*net.IPNet, 0, cfg.Whitelists.Len())
// 		for _, addr := range cfg.Whitelists.S() {
// 			_, ipnet, e := net.ParseCIDR(addr)
// 			if e != nil {
// 				E.Ln(e)
// 				ip = net.ParseIP(addr)
// 				if ip == nil {
// 					str := e.Error() + " %s: The whitelist value of '%s' is invalid"
// 					e = fmt.Errorf(str, funcName, addr)
// 					E.Ln(e)
// 					_, _ = fmt.Fprintln(os.Stderr, e)
// 					interrupt.Request()
// 					// os.Exit(1)
// 				} else {
// 					var bits int
// 					if ip.To4() == nil {
// 						// IPv6
// 						bits = 128
// 					} else {
// 						bits = 32
// 					}
// 					ipnet = &net.IPNet{
// 						IP:   ip,
// 						Mask: net.CIDRMask(bits, bits),
// 					}
// 				}
// 			}
// 			st.ActiveWhitelists = append(st.ActiveWhitelists, ipnet)
// 		}
// 	}
// }

// func validatePeerLists(cfg *opts.Config) {
// 	T.Ln("checking addpeer and connectpeer lists")
// 	if cfg.AddPeers.Len() > 0 && cfg.ConnectPeers.Len() > 0 {
// 		e := fmt.Errorf(
// 			"%s: the --addpeer and --connect options can not be mixed",
// 			funcName,
// 		)
// 		_, _ = fmt.Fprintln(os.Stderr, e)
// 		// os.Exit(1)
// 	}
// }
