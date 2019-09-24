package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/discovery"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/normalize"
	"github.com/p9c/pod/pkg/peer/connmgr"
	"github.com/p9c/pod/pkg/util"

	"github.com/btcsuite/go-socks/socks"

	"github.com/p9c/pod/cmd/node"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
)

// Configure loads and sanitises the configuration from urfave/cli
func Configure(cx *conte.Xt) {
	log.TRACE("configuring pod")
	var err error
	if cx.Config.DataDir == nil || *cx.Config.DataDir == "" {
		*cx.Config.DataDir = appdata.Dir("pod", false)
	}
	// theoretically, the configuration should be accessed only when locked
	cfg := cx.Config
	cfg.Lock()
	state := cx.StateCfg
	if *cfg.ConfigFile == "" {
		*cfg.ConfigFile =
			*cfg.DataDir + string(os.PathSeparator) + podConfigFilename
	}
	if *cfg.LogDir == "" {
		*cfg.LogDir = *cfg.DataDir
	}
	network := "mainnet"
	if cfg.Network != nil {
		network = *cfg.Network
	}
	switch network {
	case "testnet", "testnet3", "t":
		log.TRACE("on testnet")
		*cfg.TestNet3 = true
		*cfg.SimNet = false
		*cfg.RegressionTest = false
		cx.ActiveNet = &netparams.TestNet3Params
		fork.IsTestnet = true
	case "regtestnet", "regressiontest", "r":
		log.TRACE("on regression testnet")
		*cfg.TestNet3 = false
		*cfg.SimNet = false
		*cfg.RegressionTest = true
		cx.ActiveNet = &netparams.RegressionTestParams
	case "simnet", "s":
		log.TRACE("on simnet")
		*cfg.TestNet3 = false
		*cfg.SimNet = true
		*cfg.RegressionTest = false
		cx.ActiveNet = &netparams.SimNetParams
	default:
		if network != "mainnet" && network != "m" {
			log.WARN("using mainnet for node")
		}
		log.TRACE("on mainnet")
		*cfg.TestNet3 = false
		*cfg.SimNet = false
		*cfg.RegressionTest = false
		cx.ActiveNet = &netparams.MainNetParams
	}
	// setting up for zeroconf discovery advertisement
	routeable := discovery.GetRouteableInterface()
	cx.RouteableInterface = routeable
	addrs, _ := routeable.Addrs()
	routeableString := strings.Split(addrs[0].String(), "/")[0]
	cx.StopDiscovery, cx.RequestDiscoveryUpdate, err = discovery.
		Serve(cx.ActiveNet, cx.RouteableInterface, *cx.Config.Group)
	if err != nil {
		log.ERROR("error starting discovery server: ", err)
	}
	cx.StateCfg.DiscoveryUpdate = cx.RequestDiscoveryUpdate
	cx.StateCfg.RouteableAddress = routeableString

	if len(*cfg.Listeners) < 1 && !*cfg.DisableListen &&
		len(*cfg.ConnectPeers) < 1 {
		*cfg.Listeners = append(*cfg.Listeners, ":"+cx.ActiveNet.DefaultPort)
	}
	if len(*cfg.WalletRPCListeners) < 1 && !*cfg.DisableRPC {
		*cfg.WalletRPCListeners = append(*cfg.WalletRPCListeners,
			":"+cx.ActiveNet.RPCServerPort)
	}
	if len(*cfg.RPCListeners) < 1 {
		*cfg.RPCListeners = append(*cfg.RPCListeners,
			":"+cx.ActiveNet.RPCClientPort)
	}
	if *cfg.RPCCert == "" {
		*cfg.RPCCert =
			*cfg.DataDir + string(os.PathSeparator) + "rpc.cert"
	}
	if *cfg.RPCKey == "" {
		*cfg.RPCKey =
			*cfg.DataDir + string(os.PathSeparator) + "rpc.key"
	}
	if *cfg.CAFile == "" {
		*cfg.CAFile =
			*cfg.DataDir + string(os.PathSeparator) + "cafile"
	}
	loglevel := *cfg.LogLevel
	switch loglevel {
	case "trace", "debug", "info", "warn", "error", "fatal", "off":
		log.L.SetLevel(loglevel, true)
		log.TRACE("log level", loglevel)
	default:
		log.INFO("unrecognised loglevel", loglevel, "setting default info")
		*cfg.LogLevel = "info"
		log.L.SetLevel("info", true)
	}
	log.L.SetLevel(*cfg.LogLevel, true)
	if !*cfg.Onion {
		*cfg.OnionProxy = ""
	}

	log.TRACE("normalising addresses")
	port := node.DefaultPort
	nrm := normalize.StringSliceAddresses
	nrm(cfg.AddPeers, port)
	nrm(cfg.ConnectPeers, port)
	// nrm(cfg.Listeners, port)
	nrm(cfg.Whitelists, port)
	// nrm(cfg.RPCListeners, port)
	// Don't add peers from the config file when in regression test mode.
	if *cfg.RegressionTest && len(*cfg.AddPeers) > 0 {
		*cfg.AddPeers = nil
	}
	p9 := fork.P9AlgoVers
	// Set the mining algorithm correctly, default to random if unrecognised
	switch *cfg.Algo {
	case p9[0], p9[1], p9[2], p9[3], p9[4], p9[5], p9[6], p9[7], p9[8], "random", "easy":
	default:
		*cfg.Algo = "random"
	}
	log.TRACE("mining algorithm ", *cfg.Algo)
	relayNonStd := *cfg.RelayNonStd
	funcName := "loadConfig"
	switch {
	case *cfg.RelayNonStd && *cfg.RejectNonStd:
		errf := "%s: rejectnonstd and relaynonstd cannot be used together" +
			" -- choose only one"
		log.ERROR(errf, funcName)
		// just leave both false
		*cfg.RelayNonStd = false
		*cfg.RejectNonStd = false
	case *cfg.RejectNonStd:
		relayNonStd = false
	case *cfg.RelayNonStd:
		relayNonStd = true
	}
	*cfg.RelayNonStd = relayNonStd
	// Validate database type.
	log.TRACE("validating database type")
	if !node.ValidDbType(*cfg.DbType) {
		str := "%s: The specified database type [%v] is invalid -- " +
			"supported types %v"
		err := fmt.Errorf(str, funcName, *cfg.DbType, node.KnownDbTypes)
		log.ERROR(funcName, err)
		// set to default
		*cfg.DbType = node.KnownDbTypes[0]
	}
	// Validate profile port number
	log.TRACE("validating profile port number")
	if *cfg.Profile != "" {
		profilePort, err := strconv.Atoi(*cfg.Profile)
		if err != nil || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			err := fmt.Errorf(str, funcName)
			log.ERROR(funcName, err)
			*cfg.Profile = ""
		}
	}
	// Don't allow ban durations that are too short.
	log.TRACE("validating ban duration")
	if *cfg.BanDuration < time.Second {
		err := fmt.Errorf("%s: The banduration option may not be less than 1s -- parsed [%v]",
			funcName, *cfg.BanDuration)
		log.INFO(funcName, err)
		*cfg.BanDuration = node.DefaultBanDuration
	}
	// Validate any given whitelisted IP addresses and networks.
	log.TRACE("validating whitelists")
	if len(*cfg.Whitelists) > 0 {
		var ip net.IP
		state.ActiveWhitelists = make([]*net.IPNet, 0, len(*cfg.Whitelists))
		for _, addr := range *cfg.Whitelists {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				err = fmt.Errorf("%s '%s'", err.Error())
				ip = net.ParseIP(addr)
				if ip == nil {
					str := err.Error() + " %s: The whitelist value of '%s' is invalid"
					err = fmt.Errorf(str, funcName, addr)
					log.ERROR(err)
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
				var bits int
				if ip.To4() == nil {
					// IPv6
					bits = 128
				} else {
					bits = 32
				}
				ipnet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(bits, bits),
				}
			}
			state.ActiveWhitelists = append(state.ActiveWhitelists, ipnet)
		}
	}
	log.TRACE("checking addpeer and connectpeer lists")
	if len(*cfg.AddPeers) > 0 && len(*cfg.ConnectPeers) > 0 {
		err := fmt.Errorf(
			"%s: the --addpeer and --connect options can not be mixed",
			funcName)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// --proxy or --connect without --listen disables listening.
	log.TRACE("checking proxy/connect for disabling listening")
	if (*cfg.Proxy != "" || len(*cfg.ConnectPeers) > 0) &&
		len(*cfg.Listeners) == 0 {
		*cfg.DisableListen = true
	}
	// Add the default listener if none were specified. The default listener is
	// all addresses on the listen port for the network we are to connect to.
	log.TRACE("checking if listener was set")
	if len(*cfg.Listeners) == 0 {
		*cfg.Listeners = []string{":" + cx.ActiveNet.DefaultPort}
	}
	// Check to make sure limited and admin users don't have the same username
	log.TRACE("checking admin and limited username is different")
	if *cfg.Username != "" &&
		*cfg.Username == *cfg.LimitUser {
		str := "%s: --username and --limituser must not specify the same username"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Check to make sure limited and admin users don't have the same password
	log.TRACE("checking limited and admin passwords are not the same")
	if *cfg.Password != "" &&
		*cfg.Password == *cfg.LimitPass {
		str := "%s: --password and --limitpass must not specify the same password"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// The RPC server is disabled if no username or password is provided.
	log.TRACE("checking rpc server has a login enabled")
	if (*cfg.Username == "" || *cfg.Password == "") &&
		(*cfg.LimitUser == "" || *cfg.LimitPass == "") {
		*cfg.DisableRPC = true
	}
	if *cfg.DisableRPC {
		log.TRACE("RPC service is disabled")
	}
	log.TRACE("checking rpc server has listeners set")
	if !*cfg.DisableRPC && len(*cfg.RPCListeners) == 0 {
		log.DEBUG("looking up default listener")
		addrs, err := net.LookupHost(node.DefaultRPCListener)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		*cfg.RPCListeners = make([]string, 0, len(addrs))
		log.DEBUG("setting listeners")
		for _, addr := range addrs {
			addr = net.JoinHostPort(addr, cx.ActiveNet.RPCClientPort)
			*cfg.RPCListeners = append(*cfg.RPCListeners, addr)
		}
	}
	log.TRACE("checking rpc max concurrent requests")
	if *cfg.RPCMaxConcurrentReqs < 0 {
		str := "%s: The rpcmaxwebsocketconcurrentrequests option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *cfg.RPCMaxConcurrentReqs)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Validate the the minrelaytxfee.
	log.TRACE("checking min relay tx fee")
	state.ActiveMinRelayTxFee, err = util.NewAmount(*cfg.MinRelayTxFee)
	if err != nil {
		str := "%s: invalid minrelaytxfee: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Limit the max block size to a sane value.
	log.TRACE("checking max block size")
	if *cfg.BlockMaxSize < node.BlockMaxSizeMin ||
		*cfg.BlockMaxSize > node.BlockMaxSizeMax {
		str := "%s: The blockmaxsize option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxSizeMin,
			node.BlockMaxSizeMax, *cfg.BlockMaxSize)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Limit the max block weight to a sane value.
	log.TRACE("checking max block weight")
	if *cfg.BlockMaxWeight < node.BlockMaxWeightMin ||
		*cfg.BlockMaxWeight > node.BlockMaxWeightMax {
		str := "%s: The blockmaxweight option must be in between %d and %d -- parsed [%d]"
		err := fmt.Errorf(str, funcName, node.BlockMaxWeightMin,
			node.BlockMaxWeightMax, *cfg.BlockMaxWeight)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Limit the max orphan count to a sane vlue.
	log.TRACE("checking max orphan limit")
	if *cfg.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx option may not be less than 0 -- parsed [%d]"
		err := fmt.Errorf(str, funcName, *cfg.MaxOrphanTxs)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Limit the block priority and minimum block sizes to max block size.
	log.TRACE("validating block priority and minimum size/weight")
	*cfg.BlockPrioritySize = int(apputil.MinUint32(
		uint32(*cfg.BlockPrioritySize),
		uint32(*cfg.BlockMaxSize)))
	*cfg.BlockMinSize = int(apputil.MinUint32(
		uint32(*cfg.BlockMinSize),
		uint32(*cfg.BlockMaxSize)))
	*cfg.BlockMinWeight = int(apputil.MinUint32(
		uint32(*cfg.BlockMinWeight),
		uint32(*cfg.BlockMaxWeight)))
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set
	// the limit for the max block size to a safe limit so weight takes
	// precedence.
	case *cfg.BlockMaxSize == node.DefaultBlockMaxSize &&
		*cfg.BlockMaxWeight != node.DefaultBlockMaxWeight:
		*cfg.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll
	// scale the set weight accordingly based on the max block size value.
	case *cfg.BlockMaxSize != node.DefaultBlockMaxSize &&
		*cfg.BlockMaxWeight == node.DefaultBlockMaxWeight:
		*cfg.BlockMaxWeight = *cfg.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	log.TRACE("checking user agent comments", cfg.UserAgentComments)
	for _, uaComment := range *cfg.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			err := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	// Check mining addresses are valid and saved parsed versions.
	log.TRACE("checking mining addresses")
	state.ActiveMiningAddrs = make([]util.Address, 0, len(*cfg.MiningAddrs))
	for _, strAddr := range *cfg.MiningAddrs {
		addr, err := util.DecodeAddress(strAddr, cx.ActiveNet)
		if err != nil {
			str := "%s: mining address '%s' failed to decode: %v"
			err := fmt.Errorf(str, funcName, strAddr, err)
			fmt.Fprintln(os.Stderr, err)
			// os.Exit(1)
			continue
		}
		if !addr.IsForNet(cx.ActiveNet) {
			str := "%s: mining address '%s' is on the wrong network"
			err := fmt.Errorf(str, funcName, strAddr)
			fmt.Fprintln(os.Stderr, err)
			// os.Exit(1)
			continue
		}
		state.ActiveMiningAddrs = append(state.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (*cfg.Generate) && len(*cfg.MiningAddrs) == 0 { // || *cfg.MinerListener != ""
		// str := "%s: the generate flag is set, but there are no mining addresses specified "
		// err := fmt.Errorf(str, funcName)
		// fmt.Fprintln(os.Stderr, err)
		// os.Exit(1)
		*cfg.Generate = false
	}
	if *cfg.MinerPass != "" {
		state.ActiveMinerKey = fork.Argon2i([]byte(*cfg.MinerPass))
	}
	log.TRACE("checking rpc listener addresses")
	nrms := normalize.Addresses
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	// *cfg.RPCListeners = nrms(*cfg.RPCListeners, cx.ActiveNet.RPCClientPort)
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	// *cfg.Listeners = nrms(*cfg.Listeners, cx.ActiveNet.DefaultPort)
	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	*cfg.AddPeers = nrms(*cfg.AddPeers, cx.ActiveNet.DefaultPort)
	*cfg.ConnectPeers = nrms(*cfg.ConnectPeers,
		cx.ActiveNet.DefaultPort)
	// --onionproxy and not --onion are contradictory (TODO: this is kinda
	// stupid hm? switch *and* toggle by presence of flag value, one should be
	// enough)
	if !*cfg.Onion && *cfg.OnionProxy != "" {
		err := fmt.Errorf("%s: the --onionproxy and --onion options may not be activated at the same time", funcName)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Check the checkpoints for syntax errors.
	log.TRACE("checking the checkpoints")
	state.AddedCheckpoints, err = node.ParseCheckpoints(*cfg.AddCheckpoints)
	if err != nil {
		str := "%s: Error parsing checkpoints: %v"
		err := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if *cfg.TorIsolation &&
		*cfg.Proxy == "" &&
		*cfg.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		err := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// Setup dial and DNS resolution (lookup) functions depending on the
	// specified options.  The default is to use the standard net.DialTimeout
	// function as well as the system DNS resolver.  When a proxy is specified,
	// the dial function is set to the proxy specific dial function and the
	// lookup is set to use tor (unless --noonion is specified in which case the
	// system DNS resolver is used).
	log.TRACE("setting network dialer and lookup")
	state.Dial = net.DialTimeout
	state.Lookup = net.LookupIP
	if *cfg.Proxy != "" {
		log.TRACE("we are loading a proxy!")
		_, _, err := net.SplitHostPort(*cfg.Proxy)
		if err != nil {
			str := "%s: Proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *cfg.Proxy, err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Tor isolation flag means proxy credentials will be overridden unless
		// there is also an onion proxy configured in which case that one will be
		// overridden.
		torIsolation := false
		if *cfg.TorIsolation &&
			*cfg.OnionProxy == "" &&
			(*cfg.ProxyUser != "" ||
				*cfg.ProxyPass != "") {
			torIsolation = true
			log.WARN("Tor isolation set -- overriding specified" +
				" proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         *cfg.Proxy,
			Username:     *cfg.ProxyUser,
			Password:     *cfg.ProxyPass,
			TorIsolation: torIsolation,
		}
		state.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless
		// the --noonion flag is set or there is an onion-specific proxy
		// configured.
		if *cfg.Onion &&
			*cfg.OnionProxy == "" {
			state.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, *cfg.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The
	// default is to use the same dial function selected above.  However, when
	// an onion-specific proxy is specified, the onion address dial function is
	// set to use the onion-specific proxy while leaving the normal dial
	// function as selected above.  This allows .onion address traffic to be
	// routed through a different proxy than normal traffic.
	log.TRACE("setting up tor proxy if enabled")
	if *cfg.OnionProxy != "" {
		_, _, err := net.SplitHostPort(*cfg.OnionProxy)
		if err != nil {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			err := fmt.Errorf(str, funcName, *cfg.OnionProxy, err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if *cfg.TorIsolation &&
			(*cfg.OnionProxyUser != "" || *cfg.OnionProxyPass != "") {
			log.WARN("Tor isolation set - overriding specified onionproxy user" +
				" credentials")
		}
	}
	log.TRACE("setting onion dialer")
	state.Oniondial =
		func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         *cfg.OnionProxy,
				Username:     *cfg.OnionProxyUser,
				Password:     *cfg.OnionProxyPass,
				TorIsolation: *cfg.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
	// When configured in bridge mode (both --onion and --proxy are
	// configured), it means that the proxy configured by --proxy is not a
	// tor proxy, so override the DNS resolution to use the onion-specific
	// proxy.
	log.TRACE("setting proxy lookup")
	if *cfg.Proxy != "" {
		state.Lookup = func(host string) ([]net.IP, error) {
			return connmgr.TorLookupIP(host, *cfg.OnionProxy)
		}
	} else {
		state.Oniondial = state.Dial
	}
	// Specifying --noonion means the onion address dial function results in
	// an error.
	if !*cfg.Onion {
		state.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}
	// if the user set the save flag, or file doesn't exist save the file now
	if state.Save || !apputil.FileExists(*cx.Config.ConfigFile) {
		log.TRACE("saving configuration")
		save.Pod(cx.Config)
	}
	// if we are using discovery we override the listeners with ":0" and
	// the system takes care of interfaces and port allocation
	if !*cx.Config.NoDiscovery {
		*cx.Config.Listeners = []string{":0"}
		*cx.Config.RPCListeners = []string{":0"}
		*cx.Config.WalletRPCListeners = []string{":0"}
		*cx.Config.ExperimentalRPCListeners = []string{":0"}
	}
	cfg.Unlock()
}
