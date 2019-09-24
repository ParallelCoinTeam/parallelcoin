package pod

import (
	"reflect"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/urfave/cli"

	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/log"
)

type Schema struct {
	Groups []Group `json:"groups"`
}
type Group struct {
	Legend string  `json:"legend"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Group       string   `json:"group"`
	Type        string   `json:"type"`
	Name        string   `json:"label"`
	Description string   `json:"help"`
	InputType   string   `json:"inputType"`
	Featured    string   `json:"featured"`
	Model       string   `json:"model"`
	Datatype    string   `json:"datatype"`
	Options     []string `json:"options"`
}

func GetConfigSchema() Schema {
	t := reflect.TypeOf(Config{})
	var levelOptions, network, algos []string
	for _, i := range log.Levels {
		levelOptions = append(levelOptions, i)
	}
	algos = append(algos, "random")
	for _, x := range fork.P9AlgoVers {
		algos = append(algos, x)
	}
	network = []string{"mainnet", "testnet", "regtestnet", "simnet"}

	//  groups = []string{"config", "node", "debug", "rpc", "wallet", "proxy", "policy", "mining", "tls"}
	var groups []string
	rawFields := make(map[string][]Field)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		var options []string
		switch {
		case field.Name == "LogLevel":
			options = levelOptions
		case field.Name == "Network":
			options = network
		case field.Name == "Algo":
			options = algos
		}
		f := Field{
			Group:       field.Tag.Get("group"),
			Type:        field.Tag.Get("type"),
			Name:        field.Tag.Get("name"),
			Description: field.Tag.Get("description"),
			InputType:   field.Tag.Get("inputType"),
			Featured:    field.Tag.Get("featured"),
			Options:     options,
			Datatype:    field.Type.String(),
			Model:       field.Tag.Get("model"),
		}
		if f.Group != "" {
			rawFields[f.Group] = append(rawFields[f.Group], f)
		}
		// groups = append(groups, f.Group)
	}
	spew.Dump(groups)
	var outGroups []Group
	for fg, f := range rawFields {
		group := Group{
			Legend: fg,
			Fields: f,
		}
		outGroups = append(outGroups, group)
	}
	return Schema{
		Groups: outGroups,
	}
}

// Config is
type Config struct {
	sync.Mutex
	AddCheckpoints           *cli.StringSlice `group:"debug" name:"AddCheckpoints" description:"add custom checkpoints" type:"input" inputType:"text" model:"AddCheckpoints" featured:"false"`
	AddPeers                 *cli.StringSlice `group:"node" name:"Add Peers" description:"Manually adds addresses to try to connect to" type:"input" inputType:"text" model:"AddPeers" featured:"false"`
	AddrIndex                *bool            `group:"node" name:"Addr Index" description:"maintain a full address-based transaction index which makes the searchrawtransactions RPC available" type:"switch" model:"AddrIndex" featured:"false"`
	Algo                     *string          `group:"mining" name:"Algo" description:"algorithm to mine, random is best" type:"input" inputType:"text" model:"Algo" featured:"false"`
	BanDuration              *time.Duration   `group:"debug" name:"Ban Duration" description:"how long a ban of a misbehaving peer lasts" type:"input" inputType:"text" model:"BanDuration" featured:"false"`
	BanThreshold             *int             `group:"debug" name:"Ban Threshold" description:"ban score that triggers a ban (default 100)" type:"input" inputType:"number" model:"BanThreshold" featured:"false"`
	BlockMaxSize             *int             `group:"mining" name:"Block Max Size" description:"maximum block size in bytes to be used when creating a block" type:"input" inputType:"number" model:"BlockMaxSize" featured:"false"`
	BlockMaxWeight           *int             `group:"mining" name:"Block Max Weight" description:"maximum block weight to be used when creating a block" type:"input" inputType:"number" model:"BlockMaxWeight" featured:"false"`
	BlockMinSize             *int             `group:"mining" name:"Block Min Size" description:"minimum block size in bytes to be used when creating a block" type:"input" inputType:"number" model:"BlockMinSize" featured:"false"`
	BlockMinWeight           *int             `group:"mining" name:"Block Min Weight" description:"minimum block weight to be used when creating a block" type:"input" inputType:"number" model:"BlockMinWeight" featured:"false"`
	BlockPrioritySize        *int             `group:"mining" name:"Block Priority Size" description:"size in bytes for high-priority/low-fee transactions when creating a block" type:"input" inputType:"number" model:"BlockPrioritySize" featured:"false"`
	BlocksOnly               *bool            `group:"node" name:"Blocks Only" description:"do not accept transactions from remote peers" type:"switch" model:"BlocksOnly" featured:"false"`
	CAFile                   *string          `group:"tls" name:"CA File" description:"certificate authority file for TLS certificate validation" type:"input" inputType:"text" model:"CAFile" featured:"false"`
	ConfigFile               *string
	ConnectPeers             *cli.StringSlice `group:"node" name:"Connect Peers" description:"Connect ONLY to these addresses (disables inbound connections)" type:"input" inputType:"text" model:"ConnectPeers" featured:"false"`
	Controller               *string          `category:"mining" name:"Controller Listener" description:"address to bind miner controller to"`
	CPUProfile               *string          `group:"debug" name:"CPU Profile" description:"write cpu profile to this file" type:"input" inputType:"text" model:"CPUProfile" featured:"false"`
	DataDir                  *string          `group:"config" name:"Data Dir" description:"Root folder where application data is stored" type:"input" inputType:"text" model:"DataDir" featured:"false"`
	DbType                   *string          `group:"debug" name:"Db Type" description:"type of database storage engine to use (only one right now)" type:"input" inputType:"text" model:"DbType" featured:"false"`
	DisableBanning           *bool            `group:"debug" name:"Disable Banning" description:"Disables banning of misbehaving peers" type:"switch" model:"DisableBanning" featured:"false"`
	DisableCheckpoints       *bool            `group:"debug" name:"Disable Checkpoints" description:"disables all checkpoints" type:"switch" model:"DisableCheckpoints" featured:"false"`
	DisableDNSSeed           *bool            `group:"node" name:"Disable DNS Seed" description:"disable seeding of addresses to peers" type:"switch" model:"DisableDNSSeed" featured:"false"`
	DisableListen            *bool            `group:"node" name:"Disable Listen" description:"Disables inbound connections for the peer to peer network" type:"switch" model:"DisableListen" featured:"false"`
	DisableRPC               *bool            `group:"rpc" name:"Disable RPC" description:"disable rpc servers" type:"switch" model:"DisableRPC" featured:"false"`
	ExperimentalRPCListeners *cli.StringSlice `group:"wallet" name:"Experimental RPC Listeners" description:"addresses for experimental RPC listeners to listen on" type:"input" inputType:"text" model:"ExperimentalRPCListeners" featured:"false"`
	ExternalIPs              *cli.StringSlice `group:"node" name:"External IPs" description:"extra addresses to tell peers they can connect to" type:"input" inputType:"text" model:"ExternalIPs" featured:"false"`
	FreeTxRelayLimit         *float64         `group:"policy" name:"Free Tx Relay Limit" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute" type:"input" inputType:"text" model:"FreeTxRelayLimit" featured:"false"`
	Generate                 *bool            `group:"mining" name:"Generate" description:"turn on built in CPU miner" type:"switch" model:"Generate" featured:"false"`
	GenThreads               *int             `group:"mining" name:"Gen Threads" description:"number of CPU threads to mine using" type:"input" inputType:"number" model:"GenThreads" featured:"false"`
	Group                    *string          `category:"debug" name:"Zeroconf Group" description:"if set to non-empty zeroconf will set all found peers to connect peers (dns seeding can be on, this is mainly for testnets)"`
	LimitPass                *string          `group:"rpc" name:"Limit Pass" type:"password" description:"limited user password" type:"input" inputType:"text" model:"LimitPass" featured:"false"`
	LimitUser                *string          `group:"rpc" name:"Limit User" description:"limited user name" type:"input" inputType:"text" model:"LimitUser" featured:"false"`
	Listeners                *cli.StringSlice `group:"node" name:"Listeners" description:"List of addresses to bind the node listener to" type:"input" inputType:"text" model:"Listeners" featured:"false"`
	LogDir                   *string          `group:"config" name:"Log Dir" description:"Folder where log files are written" type:"input" inputType:"text" model:"LogDir" featured:"false"`
	LogLevel                 *string          `group:"config" name:"Log Level" description:"Verbosity of log printouts" type:"input" inputType:"text" model:"LogLevel" featured:"false"`
	MaxOrphanTxs             *int             `group:"policy" name:"Max Orphan Txs" description:"max number of orphan transactions to keep in memory" type:"input" inputType:"number" model:"MaxOrphanTxs" featured:"false"`
	MaxPeers                 *int             `group:"node" name:"Max Peers" description:"Maximum number of peers to hold connections with" type:"input" inputType:"number" model:"MaxPeers" featured:"false"`
	Broadcast                *bool            `group:"mining" name:"Broadcast" description:"enable broadcasting of blocks for workers to work on" type:"switch" model:"Broadcast" featured:"false"`
	Workers                  *cli.StringSlice `group:"mining" name:"MinerWorkers" description:"a list of addresses where workers are listening for blocks when not using lan broadcast" type:"input" inputType:"text" model:"Workers" featured:"false"`
	MinerPass                *string          `group:"mining" name:"Miner Pass" description:"password that encrypts the connection to the mining controller" type:"input" inputType:"text" model:"MinerPass" featured:"false"`
	MiningAddrs              *cli.StringSlice `group:"mining" name:"Mining Addrs" description:"addresses to pay block rewards to (TODO, make this auto)" type:"input" inputType:"text" model:"MiningAddrs" featured:"false"`
	MinRelayTxFee            *float64         `group:"policy" name:"Min Relay Tx Fee" description:"the minimum transaction fee in DUO/kB to be considered a non-zero fee" type:"input" inputType:"text" model:"MinRelayTxFee" featured:"false"`
	Network                  *string          `group:"node" name:"Network" description:"Which network are you connected to (eg.: mainnet, testnet)" type:"input" inputType:"text" model:"Network" featured:"false"`
	NoCFilters               *bool            `group:"node" name:"No CFilters" description:"disable committed filtering (CF) support" type:"switch" model:"NoCFilters" featured:"false"`
	NoController             *bool            `category:"node" name:"Disable Controller" description:"disables the zeroconf peer discovery/miner controller system"`
	NodeOff                  *bool            `group:"debug" name:"Node Off" description:"turn off the node backend" type:"switch" model:"NodeOff" featured:"false"`
	NoInitialLoad            *bool
	NoPeerBloomFilters       *bool   `group:"node" name:"No Peer Bloom Filters" description:"disable bloom filtering support" type:"switch" model:"NoPeerBloomFilters" featured:"false"`
	NoRelayPriority          *bool   `group:"policy" name:"No Relay Priority" description:"do not require free or low-fee transactions to have high priority for relaying" type:"switch" model:"NoRelayPriority" featured:"false"`
	OneTimeTLSKey            *bool   `group:"wallet" name:"One Time TLS Key" description:"generate a new TLS certpair at startup, but only write the certificate to disk" type:"switch" model:"OneTimeTLSKey" featured:"false"`
	Onion                    *bool   `group:"proxy" name:"Onion" description:"enable tor proxy" type:"switch" model:"Onion" featured:"false"`
	OnionProxy               *string `group:"proxy" name:"Onion Proxy" description:"address of tor proxy you want to connect to" type:"input" inputType:"text" model:"OnionProxy" featured:"false"`
	OnionProxyPass           *string `group:"proxy" name:"Onion Proxy Pass" type:"password" description:"password for tor proxy" type:"input" inputType:"text" model:"OnionProxyPass" featured:"false"`
	OnionProxyUser           *string `group:"proxy" name:"Onion Proxy User" description:"tor proxy username" type:"input" inputType:"text" model:"OnionProxyUser" featured:"false"`
	Password                 *string `group:"rpc" name:"Password" type:"password" description:"password for client RPC connections" type:"input" inputType:"text" model:"Password" featured:"false"`
	Profile                  *string `group:"debug" name:"Profile" description:"http profiling on given port (1024-40000)" type:"input" inputType:"text" model:"Profile" featured:"false"`
	Proxy                    *string `group:"proxy" name:"Proxy" description:"address of proxy to connect to for outbound connections" type:"input" inputType:"text" model:"Proxy" featured:"false"`
	ProxyPass                *string `group:"proxy" name:"Proxy Pass" type:"password" description:"proxy password, if required" type:"input" inputType:"text" model:"ProxyPass" featured:"false"`
	ProxyUser                *string `group:"proxy" name:"ProxyUser" description:"proxy username, if required" type:"input" inputType:"text" model:"ProxyUser" featured:"false"`
	RegressionTest           *bool
	RejectNonStd             *bool            `group:"node" name:"Reject Non Std" description:"reject non-standard transactions regardless of the default settings for the active network" type:"switch" model:"RejectNonStd" featured:"false"`
	RelayNonStd              *bool            `group:"node" name:"Relay Non Std" description:"relay non-standard transactions regardless of the default settings for the active network" type:"switch" model:"RelayNonStd" featured:"false"`
	RPCCert                  *string          `group:"rpc" name:"RPC Cert" description:"location of rpc TLS certificate" type:"input" inputType:"text" model:"RPCCert" featured:"false"`
	RPCConnect               *string          `group:"wallet" name:"RPC Connect" description:"full node RPC for wallet" type:"input" inputType:"text" model:"RPCConnect" featured:"false"`
	RPCKey                   *string          `group:"rpc" name:"RPC Key" description:"location of rpc TLS key" type:"input" inputType:"text" model:"RPCKey" featured:"false"`
	RPCListeners             *cli.StringSlice `group:"rpc" name:"RPC Listeners" description:"addresses to listen for RPC connections" type:"input" inputType:"text" model:"RPCListeners" featured:"false"`
	RPCMaxClients            *int             `group:"rpc" name:"RPC Max Clients" description:"maximum number of clients for regular RPC" type:"input" inputType:"number" model:"RPCMaxClients" featured:"false"`
	RPCMaxConcurrentReqs     *int             `group:"rpc" name:"RPC Max Concurrent Reqs" description:"maximum number of requests to process concurrently" type:"input" inputType:"number" model:"RPCMaxConcurrentReqs" featured:"false"`
	RPCMaxWebsockets         *int             `group:"rpc" name:"RPC Max Websockets" description:"maximum number of websocket clients to allow" type:"input" inputType:"number" model:"RPCMaxWebsockets" featured:"false"`
	RPCQuirks                *bool            `group:"rpc" name:"RPC Quirks" description:"enable bugs that replicate bitcoin core RPC's JSON" type:"switch" model:"RPCQuirks" featured:"false"`
	ServerPass               *string          `group:"rpc" name:"Server Pass" type:"password" description:"password for server connections" type:"input" inputType:"text" model:"ServerPass" featured:"false"`
	ServerTLS                *bool            `group:"wallet" name:"Server TLS" description:"Enable TLS for the wallet connection to node RPC server" type:"switch" model:"ServerTLS" featured:"false"`
	ServerUser               *string          `group:"rpc" name:"Server User" description:"username for server connections" type:"input" inputType:"text" model:"ServerUser" featured:"false"`
	SigCacheMaxSize          *int             `group:"node" name:"Sig Cache Max Size" description:"the maximum number of entries in the signature verification cache" type:"input" inputType:"number" model:"SigCacheMaxSize" featured:"false"`
	SimNet                   *bool
	Subsystems               *cli.StringSlice `group:"config" name:"Subsystems" description:"Specific systems' verbosity levels" type:"input" inputType:"text" model:"Subsystems" featured:"false"`
	TestNet3                 *bool
	TestNodeOff              *bool            `group:"debug" name:"Test Node Off" description:"turn off the testnode (testnet only)" type:"switch" model:"TestNodeOff" featured:"false"`
	TLS                      *bool            `group:"tls" name:"TLS" description:"enable TLS for RPC connections" type:"switch" model:"TLS" featured:"false"`
	TLSSkipVerify            *bool            `group:"tls" name:"TLS Skip Verify" description:"skip TLS certificate verification (ignore CA errors)" type:"switch" model:"TLSSkipVerify" featured:"false"`
	TorIsolation             *bool            `group:"proxy" name:"Tor Isolation" description:"makes a separate proxy connection for each connection" type:"switch" model:"TorIsolation" featured:"false"`
	TrickleInterval          *time.Duration   `group:"policy" name:"Trickle Interval" description:"minimum time between attempts to send new inventory to a connected peer" type:"input" inputType:"text" model:"TrickleInterval" featured:"false"`
	TxIndex                  *bool            `group:"node" name:"Tx Index" description:"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC" type:"switch" model:"TxIndex" featured:"false"`
	Upnp                     *bool            `group:"node" name:"Upnp" description:"enable UPNP for NAT traversal" type:"switch" model:"Upnp" featured:"false"`
	UserAgentComments        *cli.StringSlice `group:"node" name:"User Agent Comments" description:"Comment to add to the user agent -- See BIP 14 for more information" type:"input" inputType:"text" model:"UserAgentComments" featured:"false"`
	Username                 *string          `group:"rpc" name:"Username" description:"password for client RPC connections" type:"input" inputType:"text" model:"Username" featured:"false"`
	Wallet                   *bool
	WalletOff                *bool            `group:"debug" name:"Wallet Off" description:"turn off the wallet backend" type:"switch" model:"WalletOff" featured:"false"`
	WalletPass               *string          `group:"wallet" name:"Wallet Pass" description:"password encrypting public data in wallet" type:"input" inputType:"text" model:"WalletPass" featured:"false"`
	WalletRPCListeners       *cli.StringSlice `group:"wallet" name:"Legacy RPC Listeners" description:"addresses for wallet RPC server to listen on" type:"input" inputType:"text" model:"LegacyRPCListeners" featured:"false"`
	WalletRPCMaxClients      *int             `group:"wallet" name:"Legacy RPC Max Clients" description:"maximum number of RPC clients allowed for wallet RPC" type:"input" inputType:"number" model:"LegacyRPCMaxClients" featured:"false"`
	WalletRPCMaxWebsockets   *int             `group:"wallet" name:"Legacy RPC Max Websockets" description:"maximum number of websocket clients allowed for wallet RPC" type:"input" inputType:"number" model:"LegacyRPCMaxWebsockets" featured:"false"`
	WalletServer             *string          `group:"wallet" name:"node address to connect wallet server to" type:"input" inputType:"text" model:"WalletServer" featured:"false"`
	Whitelists               *cli.StringSlice `group:"debug" name:"Whitelists" description:"peers that you don't want to ever ban" type:"input" inputType:"text" model:"Whitelists" featured:"false"`
}

func EmptyConfig() *Config {
	return &Config{
		AddCheckpoints:           newStringSlice(),
		AddPeers:                 newStringSlice(),
		AddrIndex:                newbool(),
		Algo:                     newstring(),
		BanDuration:              newDuration(),
		BanThreshold:             newint(),
		BlockMaxSize:             newint(),
		BlockMaxWeight:           newint(),
		BlockMinSize:             newint(),
		BlockMinWeight:           newint(),
		BlockPrioritySize:        newint(),
		BlocksOnly:               newbool(),
		Broadcast:                newBool(),
		CAFile:                   newstring(),
		ConfigFile:               newstring(),
		ConnectPeers:             newStringSlice(),
		Controller:               newstring(),
		CPUProfile:               newstring(),
		DataDir:                  newstring(),
		DbType:                   newstring(),
		DisableBanning:           newbool(),
		DisableCheckpoints:       newbool(),
		DisableDNSSeed:           newbool(),
		DisableListen:            newbool(),
		DisableRPC:               newbool(),
		ExperimentalRPCListeners: newStringSlice(),
		ExternalIPs:              newStringSlice(),
		FreeTxRelayLimit:         new(float64),
		Generate:                 newbool(),
		GenThreads:               newint(),
		Group:                    newstring(),
		LimitPass:                newstring(),
		LimitUser:                newstring(),
		Listeners:                newStringSlice(),
		LogDir:                   newstring(),
		LogLevel:                 newstring(),
		MaxOrphanTxs:             newint(),
		MaxPeers:                 newint(),
		MinerPass:                newstring(),
		MiningAddrs:              newStringSlice(),
		MinRelayTxFee:            new(float64),
		Network:                  newstring(),
		NoCFilters:               newbool(),
		NoController:             newbool(),
		NodeOff:                  newbool(),
		NoDiscovery:              newbool(),
		NoInitialLoad:            newbool(),
		NoPeerBloomFilters:       newbool(),
		NoRelayPriority:          newbool(),
		OneTimeTLSKey:            newbool(),
		Onion:                    newbool(),
		OnionProxy:               newstring(),
		OnionProxyPass:           newstring(),
		OnionProxyUser:           newstring(),
		Password:                 newstring(),
		Profile:                  newstring(),
		Proxy:                    newstring(),
		ProxyPass:                newstring(),
		ProxyUser:                newstring(),
		RegressionTest:           newbool(),
		RejectNonStd:             newbool(),
		RelayNonStd:              newbool(),
		RPCCert:                  newstring(),
		RPCConnect:               newstring(),
		RPCKey:                   newstring(),
		RPCListeners:             newStringSlice(),
		RPCMaxClients:            newint(),
		RPCMaxConcurrentReqs:     newint(),
		RPCMaxWebsockets:         newint(),
		RPCQuirks:                newbool(),
		ServerPass:               newstring(),
		ServerTLS:                newbool(),
		ServerUser:               newstring(),
		SigCacheMaxSize:          newint(),
		SimNet:                   newbool(),
		Subsystems:               newStringSlice(),
		TestNet3:                 newbool(),
		TestNodeOff:              newbool(),
		TLS:                      newbool(),
		TLSSkipVerify:            newbool(),
		TorIsolation:             newbool(),
		TrickleInterval:          newDuration(),
		TxIndex:                  newbool(),
		Upnp:                     newbool(),
		UserAgentComments:        newStringSlice(),
		Username:                 newstring(),
		Wallet:                   newbool(),
		WalletOff:                newbool(),
		WalletPass:               newstring(),
		WalletRPCListeners:       newStringSlice(),
		WalletRPCMaxClients:      newint(),
		WalletRPCMaxWebsockets:   newint(),
		WalletServer:             newstring(),
		Whitelists:               newStringSlice(),
		Workers:                  newStringSlice(),
	}
}

func newbool() *bool {
	o := false
	return &o
}
func newStringSlice() *cli.StringSlice {
	o := cli.StringSlice{}
	return &o
}
func newfloat64() *float64 {
	o := 1.0
	return &o
}
func newint() *int {
	o := 1
	return &o
}
func newstring() *string {
	o := ""
	return &o
}
func newDuration() *time.Duration {
	o := time.Second - time.Second
	return &o
}
