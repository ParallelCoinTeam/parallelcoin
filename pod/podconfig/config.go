package podconfig

import (
	
	// This ensures the database drivers get registered
	// _ "github.com/p9c/matrjoska/pkg/database/ffldb"
)

//
// // Config defines the configuration options for pod. See loadConfig for details on the configuration load process.
// type Config struct {
// 	ShowVersion          *bool            `short:"True" long:"version" description:"Display version information and exit"`
// 	ConfigFile           *string          `short:"C" long:"configfile" description:"Path to configuration file"`
// 	DataDir              *string          `short:"b" long:"datadir" description:"Directory to store data"`
// 	LogDir               *string          `long:"logdir" description:"Directory to log output."`
// 	DebugLevel           *string          `long:"debuglevel" description:"baseline debug level for all subsystems unless specified"`
// 	AddPeers             *cli.StringSlice `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
// 	ConnectPeers         *cli.StringSlice `long:"connect" description:"Connect only to the specified peers at startup"`
// 	DisableListen        *bool            `long:"nolisten" description:"Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen"`
// 	P2PListeners            *cli.StringSlice `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 11047, testnet: 21047)"`
// 	MaxPeers             *int             `long:"maxpeers" description:"Max number of inbound and outbound peers"`
// 	DisableBanning       *bool            `long:"nobanning" description:"Disable banning of misbehaving peers"`
// 	BanDuration          *time.Duration   `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h, d}.  Minimum 1 second"`
// 	BanThreshold         *int             `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`
// 	Whitelists           *cli.StringSlice `long:"whitelist" description:"Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)"`
// 	RPCUser              *string          `short:"u" long:"rpcuser" description:"Username for RPC connections"`
// 	RPCPass              *string          `short:"P" long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
// 	RPCLimitUser         *string          `long:"rpclimituser" description:"Username for limited RPC connections"`
// 	RPCLimitPass         *string          `long:"rpclimitpass" default-mask:"-" description:"Password for limited RPC connections"`
// 	RPCListeners         *cli.StringSlice `long:"rpclisten" description:"Add an interface/port to listen for RPC connections (default port: 11048, testnet: 21048) gives sha256d block templates"`
// 	RPCCert              *string          `long:"rpccert" description:"File containing the certificate file"`
// 	RPCKey               *string          `long:"rpckey" description:"File containing the certificate key"`
// 	RPCMaxClients        *int             `long:"rpcmaxclients" description:"Max number of RPC clients for standard connections"`
// 	RPCMaxWebsockets     *int             `long:"rpcmaxwebsockets" description:"Max number of RPC websocket connections"`
// 	RPCMaxConcurrentReqs *int             `long:"rpcmaxconcurrentreqs" description:"Max number of concurrent RPC requests that may be processed concurrently"`
// 	RPCQuirks            *bool            `long:"rpcquirks" description:"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
// 	DisableRPC           *bool            `long:"norpc" description:"Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
// 	TLS                  *bool            `long:"tls" description:"Enable TLS for the RPC server"`
// 	DisableDNSSeed       *bool            `long:"nodnsseed" description:"Disable DNS seeding for peers"`
// 	ExternalIPs          *cli.StringSlice `long:"externalip" description:"Add an ip to the list of local addresses we claim to listen on to peers"`
// 	Proxy                *string          `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
// 	ProxyUser            *string          `long:"proxyuser" description:"Username for proxy server"`
// 	ProxyPass            *string          `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
// 	OnionProxy           *string          `long:"onion" description:"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
// 	OnionProxyUser       *string          `long:"onionuser" description:"Username for onion proxy server"`
// 	OnionProxyPass       *string          `long:"onionpass" default-mask:"-" description:"Password for onion proxy server"`
// 	Onion                *bool            `long:"noonion" description:"Disable connecting to tor hidden services"`
// 	TorIsolation         *bool            `long:"torisolation" description:"Enable Tor stream isolation by randomizing user credentials for each connection."`
// 	TestNet3             *bool            `long:"testnet" description:"Use the test network"`
// 	RegressionTest       *bool            `long:"regtest" description:"Use the regression test network"`
// 	SimNet               *bool            `long:"simnet" description:"Use the simulation test network"`
// 	AddCheckpoints       *cli.StringSlice `long:"addcheckpoint" description:"Add a custom checkpoint.  Format: '<height>:<hash>'"`
// 	DisableCheckpoints   *bool            `long:"nocheckpoints" description:"Disable built-in checkpoints.  Don't do this unless you know what you're doing."`
// 	DbType               *string          `long:"dbtype" description:"Database backend to use for the Block Chain"`
// 	Profile              *string          `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
// 	CPUProfile           *string          `long:"cpuprofile" description:"Write CPU profile to the specified file"`
// 	Upnp                 *bool            `long:"upnp" description:"Use UPnP to map our listening port outside of NAT"`
// 	MinRelayTxFee        *float64         `long:"minrelaytxfee" description:"The minimum transaction fee in DUO/kB to be considered a non-zero fee."`
// 	FreeTxRelayLimit     *float64         `long:"limitfreerelay" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
// 	NoRelayPriority      *bool            `long:"norelaypriority" description:"Do not require free or low-fee transactions to have high priority for relaying"`
// 	TrickleInterval      *time.Duration   `long:"trickleinterval" description:"Minimum time between attempts to send new inventory to a connected peer"`
// 	MaxOrphanTxs         *int             `long:"maxorphantx" description:"Max number of orphan transactions to keep in memory"`
// 	Algo                 *string          `long:"algo" description:"Sets the algorithm for the CPU miner ( blake14lr, cryptonight7v2, keccak, lyra2rev2, scrypt, sha256d, stribog, skein, x11 default is 'random')"`
// 	Generate             *bool            `long:"generate" description:"Generate (mine) bitcoins using the CPU"`
// 	GenThreads           *int             `long:"genthreads" description:"Number of CPU threads to use with CPU miner -1 = all cores"`
// 	MiningAddrs          *cli.StringSlice `long:"miningaddrs" description:"Add the specified payment address to the list of addresses to use for generated blocks, at least one is required if generate or minerport are set"`
// 	MinerListener        *string          `long:"minerlistener" description:"listen address for miner controller"`
// 	MinerPass            *string          `long:"minerpass" description:"Encryption password required for miner clients to subscribe to work updates, for use over insecure connections"`
// 	BlockMinSize         *int             `long:"blockminsize" description:"Mininum block size in bytes to be used when creating a block"`
// 	BlockMaxSize         *int             `long:"blockmaxsize" description:"Maximum block size in bytes to be used when creating a block"`
// 	BlockMinWeight       *int             `long:"blockminweight" description:"Mininum block weight to be used when creating a block"`
// 	BlockMaxWeight       *int             `long:"blockmaxweight" description:"Maximum block weight to be used when creating a block"`
// 	BlockPrioritySize    *int             `long:"blockprioritysize" description:"Size in bytes for high-priority/low-fee transactions when creating a block"`
// 	UserAgentComments    *cli.StringSlice `long:"uacomment" description:"Comment to add to the user agent -- See BIP 14 for more information."`
// 	NoPeerBloomFilters   *bool            `long:"nopeerbloomfilters" description:"Disable bloom filtering support"`
// 	NoCFilters           *bool            `long:"nocfilters" description:"Disable committed filtering (CF) support"`
// 	DropCfIndex          *bool            `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
// 	SigCacheMaxSize      *int             `long:"sigcachemaxsize" description:"The maximum number of entries in the signature verification cache"`
// 	BlocksOnly           *bool            `long:"blocksonly" description:"Do not accept transactions from remote peers."`
// 	TxIndex              *bool            `long:"txindex" description:"Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
// 	DropTxIndex          *bool            `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
// 	AddrIndex            *bool            `long:"addrindex" description:"Maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
// 	DropAddrIndex        *bool            `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
// 	RelayNonStd          *bool            `long:"relaynonstd" description:"Relay non-standard transactions regardless of the default settings for the active network."`
// 	RejectNonStd         *bool            `long:"rejectnonstd" description:"Reject non-standard transactions regardless of the default settings for the active network."`
// }

// serviceOptions defines the configuration options for the daemon as a service on Windows.
type serviceOptions struct {
	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
}

var (
	// defaultHomeDir is the default home directory location (
	// this should be centralised)
	
	// KnownDbTypes stores the currently supported database drivers
	// KnownDbTypes = database.SupportedDrivers()
	// runServiceCommand is only set to a real function on Windows.
	// It is used to parse and execute service commands specified via the -s flag.
	runServiceCommand func(string) error
)

// // newCheckpointFromStr parses checkpoints in the '<height>:<hash>' format.
// func newCheckpointFromStr(checkpoint string) (chaincfg.Checkpoint, error) {
// 	parts := strings.Split(checkpoint, ":")
// 	if len(parts) != 2 {
// 		return chaincfg.Checkpoint{}, fmt.Errorf(
// 			"unable to parse "+
// 				"checkpoint %q -- use the syntax <height>:<hash>",
// 			checkpoint,
// 		)
// 	}
// 	height, e := strconv.ParseInt(parts[0], 10, 32)
// 	if e != nil {
// 		return chaincfg.Checkpoint{}, fmt.Errorf(
// 			"unable to parse "+
// 				"checkpoint %q due to malformed height", checkpoint,
// 		)
// 	}
// 	if len(parts[1]) == 0 {
// 		return chaincfg.Checkpoint{}, fmt.Errorf(
// 			"unable to parse "+
// 				"checkpoint %q due to missing hash", checkpoint,
// 		)
// 	}
// 	hash, e := chainhash.NewHashFromStr(parts[1])
// 	if e != nil {
// 		return chaincfg.Checkpoint{}, fmt.Errorf(
// 			"unable to parse "+
// 				"checkpoint %q due to malformed hash", checkpoint,
// 		)
// 	}
// 	return chaincfg.Checkpoint{
// 			Height: int32(height),
// 			Hash:   hash,
// 		},
// 		nil
// }

//
// // NewConfigParser returns a new command line flags parser.
// func NewConfigParser(cfg *Config, so *serviceOptions, options flags.Options) *flags.Parser {
// 	parser := flags.NewParser(cfg, options)
// 	if runtime.GOOS == "windows" {
// 		_, e := parser.AddGroup("Service Options", "Service Options", so)
// 		if e != nil {
// 			panic(e)
// 		}
// 	}
// 	return parser
// }

// // ParseCheckpoints checks the checkpoint strings for valid syntax ( '<height>:<hash>') and parses them to
// // chaincfg.Checkpoint instances.
// func ParseCheckpoints(checkpointStrings []string) ([]chaincfg.Checkpoint, error) {
// 	if len(checkpointStrings) == 0 {
// 		return nil, nil
// 	}
// 	checkpoints := make([]chaincfg.Checkpoint, len(checkpointStrings))
// 	for i, cpString := range checkpointStrings {
// 		checkpoint, e := newCheckpointFromStr(cpString)
// 		if e != nil {
// 			return nil, e
// 		}
// 		checkpoints[i] = checkpoint
// 	}
// 	return checkpoints, nil
// }
//
// // ValidDbType returns whether or not dbType is a supported database type.
// func ValidDbType(dbType string) bool {
// 	for _, knownType := range KnownDbTypes {
// 		if dbType == knownType {
// 			return true
// 		}
// 	}
// 	return false
// }
//
// // validLogLevel returns whether or not logLevel is a valid debug log level.
// func validLogLevel(logLevel string) bool {
// 	switch logLevel {
// 	case "trace":
// 		fallthrough
// 	case "debug":
// 		fallthrough
// 	case "info":
// 		fallthrough
// 	case "warn":
// 		fallthrough
// 	case "error":
// 		fallthrough
// 	case "critical":
// 		return true
// 	}
// 	return false
// }

// // createDefaultConfig copies the file sample-pod.
// conf to the given destination path, and populates it with some randomly generated RPC username and password.
// func createDefaultConfigFile(destinationPath string) (e error) {
// 	// Create the destination directory if it does not exists
// 	e := os.MkdirAll(filepath.Dir(destinationPath), 0700)
// 	if e != nil  {
//		DB// 		return e
// 	}
// 	// We generate a random user and password
// 	randomBytes := make([]byte, 20)
// 	_, e = rand.Read(randomBytes)
// 	if e != nil  {
//		DB// 		return e
// 	}
// 	generatedRPCUser := base64.StdEncoding.EncodeToString(randomBytes)
// 	_, e = rand.Read(randomBytes)
// 	if e != nil  {
//		DB// 		return e
// 	}
// 	generatedRPCPass := base64.StdEncoding.EncodeToString(randomBytes)
// 	var bb bytes.Buffer
// 	bb.Write(samplePodConf)
// 	dest, e := os.OpenFile(destinationPath,
// 		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
// 	if e != nil  {
//		DB// 		return e
// 	}
// 	defer dest.Close()
// 	reader := bufio.NewReader(&bb)
// 	for e != io.EOF {
// 		var line string
// 		line, e = reader.ReadString('\n')
// 		if e != nil  && e != io.EOF {
// 			return e
// 		}
// 		if strings.Contains(line, "rpcuser=") {
// 			line = "rpcuser=" + generatedRPCUser + "\n"
// 		} else if strings.Contains(line, "rpcpass=") {
// 			line = "rpcpass=" + generatedRPCPass + "\n"
// 		}
// 		if _, e = dest.WriteString(line); E.Chk(e) {
// 			return e
// 		}
// 	}
// 	return nil
// }

/*
// loadConfig initializes and parses the config using a config file and command line options.
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
// The above results in pod functioning properly without any config settings while still allowing the user to override settings with config files and command line options.  Command line options always take precedence.
func loadConfig() (
	*Config, []string,
	error,
) {
	// Default config.
	cfg := Config{
		ConfigFile:           DefaultConfigFile,
		MaxPeers:             DefaultMaxPeers,
		BanDuration:          DefaultBanDuration,
		BanThreshold:         DefaultBanThreshold,
		RPCMaxClients:        DefaultMaxRPCClients,
		RPCMaxWebsockets:     DefaultMaxRPCWebsockets,
		RPCMaxConcurrentReqs: DefaultMaxRPCConcurrentReqs,
		DataDir:              DefaultDataDir,
		LogDir:               DefaultLogDir,
		DbType:               DefaultDbType,
		RPCKey:               DefaultRPCKeyFile,
		RPCCert:              DefaultRPCCertFile,
		MinRelayTxFee:        mempool.DefaultMinRelayTxFee.ToDUO(),
		FreeTxRelayLimit:     DefaultFreeTxRelayLimit,
		TrickleInterval:      DefaultTrickleInterval,
		BlockMinSize:         DefaultBlockMinSize,
		BlockMaxSize:         DefaultBlockMaxSize,
		BlockMinWeight:       DefaultBlockMinWeight,
		BlockMaxWeight:       DefaultBlockMaxWeight,
		BlockPrioritySize:    mempool.DefaultBlockPrioritySize,
		MaxOrphanTxs:         DefaultMaxOrphanTransactions,
		SigCacheMaxSize:      DefaultSigCacheMaxSize,
		Generate:             DefaultGenerate,
		GenThreads:           1,
		TxIndex:              DefaultTxIndex,
		AddrIndex:            DefaultAddrIndex,
		Algo:                 DefaultAlgo,
	}
	// Service options which are only added on Windows.
	serviceOpts := serviceOptions{}
	// Pre-parse the command line options to see if an alternative config file or the version flag was specified.  Any errors aside from the help message error can be ignored here since they will be caught by the final parse below.
	preCfg := cfg
	preParser := NewConfigParser(&preCfg, &serviceOpts, flags.HelpFlag)
	_, e := preParser.Parse()
	if e != nil  {
		if e, ok := err.(*flags.DBError); ok && e.Type == flags.ErrHelp {
			DB			return nil, nil, e
		}
	}
	// Show the version and exit if the version flag was specified.
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Println(appName, "version", Version())
		os.Exit(0)
	}
	// Perform service command and exit if specified.  Invalid service commands show an appropriate error.  Only runs on Windows since the runServiceCommand function will be nil when not on Windows.
	if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
		e := runServiceCommand(serviceOpts.ServiceCommand)
		if e != nil  {
			fmt.Fprintln(os.Stderr, e)
		}
		os.Exit(0)
	}
	// Load additional config from file.
	var configFileError error
	parser := NewConfigParser(&cfg, &serviceOpts, flags.Default)
	if !(preCfg.RegressionTest || preCfg.SimNet) || preCfg.ConfigFile !=
		DefaultConfigFile {
		if _, e = os.Stat(preCfg.ConfigFile); os.IsNotExist(e) {
			e := createDefaultConfigFile(preCfg.ConfigFile)
			if e != nil  {
				fmt.Fprintf(os.Stderr, "DBError creating a "+
					"default config file: %v\n", e)
			}
		}
		e := flags.NewIniParser(parser).ParseFile(preCfg.ConfigFile)
		if e != nil  {
			if _, ok := err.(*os.PathError); !ok {
				fmt.Fprintf(os.Stderr, "DBError parsing config "+
					"file: %v\n", e)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, e
			}
			configFileError = err
		}
	}
	// Don't add peers from the config file when in regression test mode.
	if preCfg.RegressionTest && len(cfg.AddPeers) > 0 {
		cfg.AddPeers = nil
	}
	// Parse command line options again to ensure they take precedence.
	remainingArgs, e := parser.Parse()
	if e != nil  {
		if e, ok := err.(*flags.DBError); !ok || e.Type != flags.ErrHelp {
			fmt.Fprintln(os.Stderr, usageMessage)
		}
		return nil, nil, e
	}
	// Create the home directory if it doesn't already exist.
	funcName := "loadConfig"
	e = os.MkdirAll(DefaultHomeDir, 0700)
	if e != nil  {
		// Show a nicer error message if it's because a symlink is linked to a directory that does not exist (probably because it's not mounted).
		if e, ok := err.(*os.PathError); ok && os.IsExist(e) {
			if link, lerr := os.Readlink(e.Path); lerr == nil {
				str := "is symlink %s -> %s mounted?"
				e = fmt.Errorf(str, e.Path, link)
			}
		}
		str := "%s: Failed to create home directory: %v"
		e := fmt.Errorf(str, funcName, e)
		fmt.Fprintln(os.Stderr, e)
		return nil, nil, e
	}
	// Multiple networks can't be selected simultaneously.
	numNets := 0
	// Count number of network flags passed; assign active network netparams while we're at it
	if cfg.TestNet3 {
		numNets++
		ActiveNetParams = &TestNet3Params
		fork.IsTestnet = true
	}
	if cfg.RegressionTest {
		numNets++
		ActiveNetParams = &RegressionTestParams
		fork.IsTestnet = true
	}
	if cfg.SimNet {
		numNets++
		// Also disable dns seeding on the simulation test network.
		ActiveNetParams = &SimNetParams
		cfg.DisableDNSSeed = true
		fork.IsTestnet = true
	}
	if numNets > 1 {
		str := "%s: The testnet, regtest, segnet, and simnet netparams " +
			"can't be used together -- choose one of the four"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Set the mining algorithm correctly, default to sha256d if unrecognised
	switch cfg.Algo {
	case "blake14lr", "cryptonight7v2", "keccak", "lyra2rev2", "scrypt", "skein", "x11", "stribog", "random", "easy":
	default:
		cfg.Algo = fork.SHA256d
	}
	// Set the default policy for relaying non-standard transactions according to the default of the active network. The set configuration value takes precedence over the default value for the selected network.
	relayNonStd := ActiveNetParams.RelayNonStdTxs
	switch {
	case cfg.RelayNonStd && cfg.RejectNonStd:
		str := "%s: rejectnonstd and relaynonstd cannot be used " +
			"together -- choose only one"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	case cfg.RejectNonStd:
		relayNonStd = false
	case cfg.RelayNonStd:
		relayNonStd = true
	}
	cfg.RelayNonStd = relayNonStd
	// Append the network type to the data directory so it is "namespaced" per network.  In addition to the block database, there are other pieces of data that are saved to disk such as address manager state. All data is specific to a network, so namespacing the data directory means each individual piece of serialized data does not have to worry about changing names per network and such.
	cfg.DataDir = CleanAndExpandPath(cfg.DataDir)
	cfg.DataDir = filepath.Join(cfg.DataDir, NetName(ActiveNetParams))
	// Append the network type to the log directory so it is "namespaced" per network in the same fashion as the data directory.
	cfg.LogDir = CleanAndExpandPath(cfg.LogDir)
	cfg.LogDir = filepath.Join(cfg.LogDir, NetName(ActiveNetParams))
	// Initialize log rotation.  After log rotation has been initialized, the logger variables may be used.
	// initLogRotator(filepath.Join(cfg.LogDir, DefaultLogFilename))
	// Validate database type.
	if !ValidDbType(cfg.DbType) {
		str := "%s: The specified database type [%v] is invalid -- supported types %v"
		e := fmt.Errorf(str, funcName, cfg.DbType, KnownDbTypes)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Validate profile port number
	if cfg.Profile != "" {
		log <- cl.trace{"profiling to", cfg.Profile}
		profilePort, e := strconv.Atoi(cfg.Profile)
		if e != nil  || profilePort < 1024 || profilePort > 65535 {
			str := "%s: The profile port must be between 1024 and 65535"
			e := fmt.Errorf(str, funcName)
			fmt.Fprintln(os.Stderr, e)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
	}
	// Don't allow ban durations that are too short.
	if cfg.BanDuration < time.Second {
		str := "%s: The banduration opt may not be less than 1s -- parsed [%v]"
		e := fmt.Errorf(str, funcName, cfg.BanDuration)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Validate any given whitelisted IP addresses and networks.
	if len(StateCfg.ActiveWhitelists) > 0 {
		var ip net.IP
		StateCfg.ActiveWhitelists = make([]*net.IPNet, 0, len(StateCfg.ActiveWhitelists))
		for _, addr := range cfg.Whitelists {
			_, ipnet, e := net.ParseCIDR(addr)
			if e != nil  {
				ip = net.ParseIP(addr)
				if ip == nil {
					str := "%s: The whitelist value of '%s' is invalid"
					e = fmt.Errorf(str, funcName, addr)
					fmt.Fprintln(os.Stderr, e)
					fmt.Fprintln(os.Stderr, usageMessage)
					return nil, nil, e
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
			StateCfg.ActiveWhitelists = append(StateCfg.ActiveWhitelists, ipnet)
		}
	}
	// --addPeer and --connect do not mix.
	if len(cfg.AddPeers) > 0 && len(cfg.ConnectPeers) > 0 {
		str := "%s: the --addpeer and --connect options can not be " +
			"mixed"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// --proxy or --connect without --listen disables listening.
	if (cfg.Proxy != "" || len(cfg.ConnectPeers) > 0) &&
		len(cfg.P2PListeners) == 0 {
		cfg.DisableListen = true
	}
	// Connect means no DNS seeding.
	if len(cfg.ConnectPeers) > 0 {
		cfg.DisableDNSSeed = true
	}
	// Add the default listener if none were specified. The default listener is all addresses on the listen port for the network we are to connect to.
	if len(cfg.P2PListeners) == 0 {
		cfg.P2PListeners = []string{
			net.JoinHostPort("", ActiveNetParams.DefaultPort),
		}
	}
	// Chk to make sure limited and admin users don't have the same username
	if cfg.RPCUser == cfg.RPCLimitUser && cfg.RPCUser != "" {
		str := "%s: --rpcuser and --rpclimituser must not specify the " +
			"same username"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Chk to make sure limited and admin users don't have the same password
	if cfg.RPCPass == cfg.RPCLimitPass && cfg.RPCPass != "" {
		str := "%s: --rpcpass and --rpclimitpass must not specify the " +
			"same password"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// The RPC server is disabled if no username or password is provided.
	if (cfg.RPCUser == "" || cfg.RPCPass == "") &&
		(cfg.RPCLimitUser == "" || cfg.RPCLimitPass == "") {
		cfg.DisableRPC = true
	}
	if cfg.DisableRPC {
		log <- cl.Inf("RPC service is disabled")
		// Default RPC to listen on localhost only.
		if !cfg.DisableRPC && len(cfg.RPCListeners) == 0 {
			addrs, e := net.LookupHost("127.0.0.1:11048")
			if e != nil  {
				return nil, nil, e
			}
			cfg.RPCListeners = make([]string, 0, len(addrs))
			for _, addr := range addrs {
				addr = net.JoinHostPort(addr, ActiveNetParams.RPCPort)
				cfg.RPCListeners = append(cfg.RPCListeners, addr)
			}
		}
		if cfg.RPCMaxConcurrentReqs < 0 {
			str := "%s: The rpcmaxwebsocketconcurrentrequests opt may not be less than 0 -- parsed [%d]"
			e := fmt.Errorf(str, funcName, cfg.RPCMaxConcurrentReqs)
			fmt.Fprintln(os.Stderr, e)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
	}
	// Validate the the minrelaytxfee.
	StateCfg.ActiveMinRelayTxFee, e = util.NewAmount(cfg.MinRelayTxFee)
	if e != nil  {
		str := "%s: invalid minrelaytxfee: %v"
		e := fmt.Errorf(str, funcName, e)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Limit the max block size to a sane value.
	if cfg.BlockMaxSize < BlockMaxSizeMin || cfg.BlockMaxSize >
		BlockMaxSizeMax {
		str := "%s: The blockmaxsize opt must be in between %d " +
			"and %d -- parsed [%d]"
		e := fmt.Errorf(str, funcName, BlockMaxSizeMin,
			BlockMaxSizeMax, cfg.BlockMaxSize)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Limit the max block weight to a sane value.
	if cfg.BlockMaxWeight < BlockMaxWeightMin ||
		cfg.BlockMaxWeight > BlockMaxWeightMax {
		str := "%s: The blockmaxweight opt must be in between %d " +
			"and %d -- parsed [%d]"
		e := fmt.Errorf(str, funcName, BlockMaxWeightMin,
			BlockMaxWeightMax, cfg.BlockMaxWeight)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Limit the max orphan count to a sane vlue.
	if cfg.MaxOrphanTxs < 0 {
		str := "%s: The maxorphantx opt may not be less than 0 " +
			"-- parsed [%d]"
		e := fmt.Errorf(str, funcName, cfg.MaxOrphanTxs)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Limit the block priority and minimum block txsizes to max block size.
	cfg.BlockPrioritySize = minUint32(cfg.BlockPrioritySize, cfg.BlockMaxSize)
	cfg.BlockMinSize = minUint32(cfg.BlockMinSize, cfg.BlockMaxSize)
	cfg.BlockMinWeight = minUint32(cfg.BlockMinWeight, cfg.BlockMaxWeight)
	switch {
	// If the max block size isn't set, but the max weight is, then we'll set the limit for the max block size to a safe limit so weight takes precedence.
	case cfg.BlockMaxSize == DefaultBlockMaxSize &&
		cfg.BlockMaxWeight != DefaultBlockMaxWeight:
		cfg.BlockMaxSize = blockchain.MaxBlockBaseSize - 1000
	// If the max block weight isn't set, but the block size is, then we'll scale the set weight accordingly based on the max block size value.
	case cfg.BlockMaxSize != DefaultBlockMaxSize &&
		cfg.BlockMaxWeight == DefaultBlockMaxWeight:
		cfg.BlockMaxWeight = cfg.BlockMaxSize * blockchain.WitnessScaleFactor
	}
	// Look for illegal characters in the user agent comments.
	for _, uaComment := range cfg.UserAgentComments {
		if strings.ContainsAny(uaComment, "/:()") {
			e := fmt.Errorf("%s: The following characters must not "+
				"appear in user agent comments: '/', ':', '(', ')'",
				funcName)
			fmt.Fprintln(os.Stderr, e)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
	}
	// --txindex and --droptxindex do not mix.
	if cfg.TxIndex && cfg.DropTxIndex {
		e := fmt.Errorf("%s: the --txindex and --droptxindex "+
			"options may  not be activated at the same time",
			funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// --addrindex and --dropaddrindex do not mix.
	if cfg.AddrIndex && cfg.DropAddrIndex {
		e := fmt.Errorf("%s: the --addrindex and --dropaddrindex "+
			"options may not be activated at the same time",
			funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// --addrindex and --droptxindex do not mix.
	if cfg.AddrIndex && cfg.DropTxIndex {
		e := fmt.Errorf("%s: the --addrindex and --droptxindex options may not be activated at the same time because the address index relies on the transaction index", funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Chk mining addresses are valid and saved parsed versions.
	StateCfg.ActiveMiningAddrs = make([]util.Address, 0, len(cfg.MiningAddrs))
	for _, strAddr := range cfg.MiningAddrs {
		addr, e := util.DecodeAddress(strAddr, Activechaincfg.Params)
		if e != nil  {
			str := "%s: mining address '%s' failed to decode: %v"
			e := fmt.Errorf(str, funcName, strAddr, e)
			fmt.Fprintln(os.Stderr, e)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
		if !addr.IsForNet(Activechaincfg.Params) {
			str := "%s: mining address '%s' is on the wrong network"
			e := fmt.Errorf(str, funcName, strAddr)
			fmt.Fprintln(os.Stderr, e)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
		StateCfg.ActiveMiningAddrs = append(StateCfg.ActiveMiningAddrs, addr)
	}
	// Ensure there is at least one mining address when the generate flag is set.
	if (cfg.Generate || cfg.MinerListener != "") && len(cfg.MiningAddrs) == 0 {
		str := "%s: the generate flag is set, but there are no mining addresses specified "
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	if cfg.MinerPass != "" {
		StateCfg.ActiveMinerKey = fork.Argon2i([]byte(cfg.MinerPass))
	}
	// Add default port to all listener addresses if needed and remove duplicate addresses.
	cfg.P2PListeners = NormalizeAddresses(cfg.P2PListeners,
		ActiveNetParams.DefaultPort)
	// Add default port to all rpc listener addresses if needed and remove duplicate addresses.
	cfg.RPCListeners = NormalizeAddresses(cfg.RPCListeners,
		ActiveNetParams.RPCPort)
	if !cfg.DisableRPC && !cfg.TLS {
		for _, addr := range cfg.RPCListeners {
			if e != nil  {
				str := "%s: RPC listen interface '%s' is invalid: %v"
				e := fmt.Errorf(str, funcName, addr, e)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return nil, nil, e
			}
		}
	}
	// Add default port to all added peer addresses if needed and remove duplicate addresses.
	cfg.AddPeers = NormalizeAddresses(cfg.AddPeers,
		ActiveNetParams.DefaultPort)
	cfg.ConnectPeers = NormalizeAddresses(cfg.ConnectPeers,
		ActiveNetParams.DefaultPort)
	// --noonion and --onion do not mix.
	if cfg.NoOnion && cfg.OnionProxy != "" {
		e := fmt.Errorf("%s: the --noonion and --onion options may not be activated at the same time", funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Chk the checkpoints for syntax errors.
	StateCfg.AddedCheckpoints, e = ParseCheckpoints(cfg.AddCheckpoints)
	if e != nil  {
		str := "%s: DBError parsing checkpoints: %v"
		e := fmt.Errorf(str, funcName, err)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Tor stream isolation requires either proxy or onion proxy to be set.
	if cfg.TorIsolation && cfg.Proxy == "" && cfg.OnionProxy == "" {
		str := "%s: Tor stream isolation requires either proxy or onionproxy to be set"
		e := fmt.Errorf(str, funcName)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return nil, nil, e
	}
	// Setup dial and DNS resolution (lookup) functions depending on the specified options.  The default is to use the standard net.DialTimeout function as well as the system DNS resolver.  When a proxy is specified, the dial function is set to the proxy specific dial function and the lookup is set to use tor (unless --noonion is specified in which case the system DNS resolver is used).
	StateCfg.Dial = net.DialTimeout
	StateCfg.Lookup = net.LookupIP
	if cfg.Proxy != "" {
		_, _, e = net.SplitHostPort(cfg.Proxy)
		if e != nil  {
			str := "%s: Proxy address '%s' is invalid: %v"
			e := fmt.Errorf(str, funcName, cfg.Proxy, err)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
		// Tor isolation flag means proxy credentials will be overridden unless there is also an onion proxy configured in which case that one will be overridden.
		torIsolation := false
		if cfg.TorIsolation && cfg.OnionProxy == "" &&
			(cfg.ProxyUser != "" || cfg.ProxyPass != "") {
			torIsolation = true
			fmt.Fprintln(os.Stderr, "Tor isolation set -- overriding specified proxy user credentials")
		}
		proxy := &socks.Proxy{
			Addr:         cfg.Proxy,
			Username:     cfg.ProxyUser,
			Password:     cfg.ProxyPass,
			TorIsolation: torIsolation,
		}
		StateCfg.Dial = proxy.DialTimeout
		// Treat the proxy as tor and perform DNS resolution through it unless the --noonion flag is set or there is an onion-specific proxy configured.
		if !cfg.NoOnion && cfg.OnionProxy == "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.Proxy)
			}
		}
	}
	// Setup onion address dial function depending on the specified options. The default is to use the same dial function selected above.  However, when an onion-specific proxy is specified, the onion address dial function is set to use the onion-specific proxy while leaving the normal dial function as selected above.  This allows .onion address traffic to be routed through a different proxy than normal traffic.
	if cfg.OnionProxy != "" {
		_, _, e = net.SplitHostPort(cfg.OnionProxy)
		if e != nil  {
			str := "%s: Onion proxy address '%s' is invalid: %v"
			e := fmt.Errorf(str, funcName, cfg.OnionProxy, err)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return nil, nil, e
		}
		// Tor isolation flag means onion proxy credentials will be overridden.
		if cfg.TorIsolation &&
			(cfg.OnionProxyUser != "" || cfg.OnionProxyPass != "") {
			fmt.Fprintln(os.Stderr, "Tor isolation set -- "+
				"overriding specified onionproxy user "+
				"credentials ")
		}
		StateCfg.Oniondial = func(network, addr string, timeout time.Duration) (net.Conn, error) {
			proxy := &socks.Proxy{
				Addr:         cfg.OnionProxy,
				Username:     cfg.OnionProxyUser,
				Password:     cfg.OnionProxyPass,
				TorIsolation: cfg.TorIsolation,
			}
			return proxy.DialTimeout(network, addr, timeout)
		}
		// When configured in bridge mode (both --onion and --proxy are configured), it means that the proxy configured by --proxy is not a tor proxy, so override the DNS resolution to use the onion-specific proxy.
		if cfg.Proxy != "" {
			StateCfg.Lookup = func(host string) ([]net.IP, error) {
				return connmgr.TorLookupIP(host, cfg.OnionProxy)
			}
		}
	} else {
		StateCfg.Oniondial = StateCfg.Dial
	}
	// Specifying --noonion means the onion address dial function results in an error.
	if cfg.NoOnion {
		StateCfg.Oniondial = func(a, b string, t time.Duration) (net.Conn, error) {
			return nil, errors.New("tor has been disabled")
		}
	}
	// wrn about missing config file only after all other configuration is done.  This prevents the warning on help messages and invalid options.  Note this should go directly before the return.
	if configFileError != nil {
		log <- cl.wrn{configFileError}
	}
	return &cfg, remainingArgs, nil
}
// minUint32 is a helper function to return the minimum of two uint32s. This avoids a math import and the need to cast to floats.
func minUint32(	a, b uint32,
	) uint32 {
		if a < b {
			return a
		}
		return b
	}
*/
