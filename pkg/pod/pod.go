package pod

import (
	"reflect"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/urfave/cli"

	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
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
	var leveloptions, network, algos []string
	for i := range cl.Levels {
		leveloptions = append(leveloptions, i)
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
			options = leveloptions
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
// nolint
type Config struct {
	sync.Mutex
	ConfigFile               *string
	DataDir                  *string          `category:"config" name:"DataDir" description:"Root folder where application data is stored" type:"input" inputType:"text" model:"DataDir" featured:"false"`
	LogDir                   *string          `category:"config" name:"LogDir" description:"Folder where log files are written" type:"input" inputType:"text" model:"LogDir" featured:"false"`
	Network                  *string          `category:"node" name:"Network" description:"Which network are you connected to (eg.: mainnet, testnet)" type:"input" inputType:"text" model:"Network" featured:"false"`
	LogLevel                 *string          `category:"config" name:"LogLevel" description:"Verbosity of log printouts" type:"input" inputType:"text" model:"LogLevel" featured:"false"`
	Subsystems               *cli.StringSlice `category:"config" name:"Subsystems" description:"Specific systems' verbosity levels" type:"input" inputType:"text" model:"Subsystems" featured:"false"`
	Group                    *string          `category:"debug" name:"Zeroconf Group" description:"if set to non-empty zeroconf will set all found peers to connect peers (dns seeding can be on, this is mainly for testnets)"`
	NoDiscovery              *bool            `category:"node" name:"NoDiscovery" description:"disables zeroconf peer autodiscovery"`
	AddPeers                 *cli.StringSlice `category:"node" name:"Add Peers" description:"Manually adds addresses to try to connect to" type:"input" inputType:"text" model:"AddPeers" featured:"false"`
	ConnectPeers             *cli.StringSlice `category:"node" name:"Connect Peers" description:"Connect ONLY to these addresses (disables inbound connections)" type:"input" inputType:"text" model:"ConnectPeers" featured:"false"`
	MaxPeers                 *int             `category:"node" name:"MaxPeers" description:"Maximum number of peers to hold connections with" type:"input" inputType:"number" model:"MaxPeers" featured:"false"`
	Listeners                *cli.StringSlice `category:"node" name:"Listeners" description:"List of addresses to bind the node listener to" type:"input" inputType:"text" model:"Listeners" featured:"false"`
	DisableListen            *bool            `category:"node" name:"DisableListen" description:"Disables inbound connections for the peer to peer network" type:"switch" model:"DisableListen" featured:"false"`
	DisableBanning           *bool            `category:"debug" name:"Disable  Banning" description:"Disables banning of misbehaving peers" type:"switch" model:"DisableBanning" featured:"false"`
	BanDuration              *time.Duration   `category:"debug" name:"Ban Duration" description:"how long a ban of a misbehaving peer lasts" type:"input" inputType:"text" model:"BanDuration" featured:"false"`
	BanThreshold             *int             `category:"debug" name:"BanThreshold" description:"ban score that triggers a ban (default 100)" type:"input" inputType:"number" model:"BanThreshold" featured:"false"`
	Whitelists               *cli.StringSlice `category:"debug" name:"Whitelists" description:"peers that you don't want to ever ban" type:"input" inputType:"text" model:"Whitelists" featured:"false"`
	Username                 *string          `category:"rpc" name:"Username" description:"password for client RPC connections" type:"input" inputType:"text" model:"Username" featured:"false"`
	Password                 *string          `category:"rpc" name:"Password" type:"password" description:"password for client RPC connections" type:"input" inputType:"text" model:"Password" featured:"false"`
	ServerUser               *string          `category:"rpc" name:"ServerUser" description:"username for server connections" type:"input" inputType:"text" model:"ServerUser" featured:"false"`
	ServerPass               *string          `category:"rpc" name:"ServerPass" type:"password" description:"password for server connections" type:"input" inputType:"text" model:"ServerPass" featured:"false"`
	LimitUser                *string          `category:"rpc" name:"LimitUser" description:"limited user name" type:"input" inputType:"text" model:"LimitUser" featured:"false"`
	LimitPass                *string          `category:"rpc" name:"LimitPass" type:"password" description:"limited user password" type:"input" inputType:"text" model:"LimitPass" featured:"false"`
	RPCConnect               *string          `category:"wallet" name:"RPCConnect" description:"full node RPC for wallet" type:"input" inputType:"text" model:"RPCConnect" featured:"false"`
	RPCListeners             *cli.StringSlice `category:"rpc" name:"RPC Listeners" description:"addresses to listen for RPC connections" type:"input" inputType:"text" model:"RPCListeners" featured:"false"`
	RPCCert                  *string          `category:"rpc" name:"RPC Cert" description:"location of rpc TLS certificate" type:"input" inputType:"text" model:"RPCCert" featured:"false"`
	RPCKey                   *string          `category:"rpc" name:"RPC Key" description:"location of rpc TLS key" type:"input" inputType:"text" model:"RPCKey" featured:"false"`
	RPCMaxClients            *int             `category:"rpc" name:"RPC MaxClients" description:"maximum number of clients for regular RPC" type:"input" inputType:"number" model:"RPCMaxClients" featured:"false"`
	RPCMaxWebsockets         *int             `category:"rpc" name:"RPC MaxWebsockets" description:"maximum number of websocket clients to allow" type:"input" inputType:"number" model:"RPCMaxWebsockets" featured:"false"`
	RPCMaxConcurrentReqs     *int             `category:"rpc" name:"RPC MaxConcurrent Reqs" description:"maximum number of requests to process concurrently" type:"input" inputType:"number" model:"RPCMaxConcurrentReqs" featured:"false"`
	RPCQuirks                *bool            `category:"rpc" name:"RPCQuirks" description:"enable bugs that replicate bitcoin core RPC's JSON" type:"switch" model:"RPCQuirks" featured:"false"`
	DisableRPC               *bool            `category:"rpc" name:"DisableRPC" description:"disable rpc servers" type:"switch" model:"DisableRPC" featured:"false"`
	TLS                      *bool            `category:"tls" name:"TLS" description:"enable TLS for RPC connections" type:"switch" model:"TLS" featured:"false"`
	DisableDNSSeed           *bool            `category:"node" name:"Disable DNS Seed" description:"disable seeding of addresses to peers" type:"switch" model:"DisableDNSSeed" featured:"false"`
	ExternalIPs              *cli.StringSlice `category:"node" name:"External IPs" description:"extra addresses to tell peers they can connect to" type:"input" inputType:"text" model:"ExternalIPs" featured:"false"`
	Proxy                    *string          `category:"proxy" name:"Proxy" description:"address of proxy to connect to for outbound connections" type:"input" inputType:"text" model:"Proxy" featured:"false"`
	ProxyUser                *string          `category:"proxy" name:"ProxyUser" description:"proxy username, if required" type:"input" inputType:"text" model:"ProxyUser" featured:"false"`
	ProxyPass                *string          `category:"proxy" name:"ProxyPass" type:"password" description:"proxy password, if required" type:"input" inputType:"text" model:"ProxyPass" featured:"false"`
	OnionProxy               *string          `category:"proxy" name:"OnionProxy" description:"address of tor proxy you want to connect to" type:"input" inputType:"text" model:"OnionProxy" featured:"false"`
	OnionProxyUser           *string          `category:"proxy" name:"OnionProxy User" description:"tor proxy username" type:"input" inputType:"text" model:"OnionProxyUser" featured:"false"`
	OnionProxyPass           *string          `category:"proxy" name:"OnionProxy Pass" type:"password" description:"password for tor proxy" type:"input" inputType:"text" model:"OnionProxyPass" featured:"false"`
	Onion                    *bool            `category:"proxy" name:"Onion" description:"enable tor proxy" type:"switch" model:"Onion" featured:"false"`
	TorIsolation             *bool            `category:"proxy" name:"TorIsolation" description:"makes a separate proxy connection for each connection" type:"switch" model:"TorIsolation" featured:"false"`
	TestNet3                 *bool
	RegressionTest           *bool
	SimNet                   *bool
	AddCheckpoints           *cli.StringSlice `category:"debug" name:"AddCheckpoints" description:"add custom checkpoints"`
	DisableCheckpoints       *bool            `category:"debug" name:"Disable Checkpoints" description:"disables all checkpoints"`
	DbType                   *string          `category:"debug" name:"Db Type" description:"type of database storage engine to use (only one right now)"`
	Profile                  *string          `category:"debug" name:"Profile" description:"http profiling on given port (1024-40000)"`
	CPUProfile               *string          `category:"debug" name:"CPU Profile" description:"write cpu profile to this file"`
	Upnp                     *bool            `category:"node" name:"Upnp" description:"enable UPNP for NAT traversal"`
	MinRelayTxFee            *float64         `category:"policy" name:"Min Relay Tx Fee" description:"the minimum transaction fee in DUO/kB to be considered a non-zero fee"`
	FreeTxRelayLimit         *float64         `category:"policy" name:"Free Tx Relay Limit" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
	NoRelayPriority          *bool            `category:"policy" name:"No Relay Priority" description:"do not require free or low-fee transactions to have high priority for relaying"`
	TrickleInterval          *time.Duration   `category:"policy" name:"Trickle Interval" description:"minimum time between attempts to send new inventory to a connected peer"`
	MaxOrphanTxs             *int             `category:"policy" name:"Max Orphan Txs" description:"max number of orphan transactions to keep in memory"`
	Algo                     *string          `category:"mining" name:"Algo" description:"algorithm to mine, random is best"`
	Generate                 *bool            `category:"mining" name:"Generate" description:"turn on built in CPU miner"`
	GenThreads               *int             `category:"mining" name:"Gen Threads" description:"number of CPU threads to mine using"`
	Controller               *string          `category:"mining" name:"Controller Listener" description:"address to bind miner controller to"`
	NoController             *bool            `category:"node" name:"Disable Controller" description:"disables the zeroconf peer discovery/miner controller system"`
	MiningAddrs              *cli.StringSlice `category:"mining" name:"Mining Addrs" description:"addresses to pay block rewards to (TODO, make this auto)"`
	MinerPass                *string          `category:"mining" name:"Miner Pass" description:"password that encrypts the connection to the mining controller"`
	BlockMinSize             *int             `category:"mining" name:"Block Min Size" description:"mininum block size in bytes to be used when creating a block"`
	BlockMaxSize             *int             `category:"mining" name:"Block Max Size" description:"maximum block size in bytes to be used when creating a block"`
	BlockMinWeight           *int             `category:"mining" name:"Block Min Weight" description:"mininum block weight to be used when creating a block"`
	BlockMaxWeight           *int             `category:"mining" name:"Block Max Weight" description:"maximum block weight to be used when creating a block"`
	BlockPrioritySize        *int             `category:"mining" name:"Block Priority Size" description:"size in bytes for high-priority/low-fee transactions when creating a block"`
	UserAgentComments        *cli.StringSlice `category:"node" name:"User Agent Comments" description:"Comment to add to the user agent -- See BIP 14 for more information"`
	NoPeerBloomFilters       *bool            `category:"node" name:"No Peer Bloom Filters" description:"disable bloom filtering support"`
	NoCFilters               *bool            `category:"node" name:"No CFilters" description:"disable committed filtering (CF) support"`
	SigCacheMaxSize          *int             `category:"node" name:"Sig Cache Max Size" description:"the maximum number of entries in the signature verification cache"`
	BlocksOnly               *bool            `category:"node" name:"BlockC Only" description:"do not accept transactions from remote peers"`
	TxIndex                  *bool            `category:"node" name:"Tx Index" description:"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	AddrIndex                *bool            `category:"node" name:"Addr Index" description:"maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
	RelayNonStd              *bool            `category:"node" name:"Relay Non Std" description:"relay non-standard transactions regardless of the default settings for the active network"`
	RejectNonStd             *bool            `category:"node" name:"Reject Non Std" description:"reject non-standard transactions regardless of the default settings for the active network"`
	TLSSkipVerify            *bool            `category:"tls" name:"TLS Skip Verify" description:"skip TLS certificate verification (ignore CA errors)"`
	Wallet                   *bool
	NoInitialLoad            *bool
	WalletPass               *string          `category:"wallet" name:"Wallet Pass" description:"password encrypting public data in wallet" type:"input" inputType:"text" model:"WalletPass" featured:"false"`
	WalletServer             *string          `category:"wallet" name:"Node Address to connect wallet server to" type:"input" inputType:"text" model:"WalletServer" featured:"false"`
	CAFile                   *string          `category:"tls" name:"CA File" description:"certificate authority file for TLS certificate validation" type:"input" inputType:"text" model:"CAFile" featured:"false"`
	OneTimeTLSKey            *bool            `category:"wallet" name:"OneTime TLS Key" description:"generate a new TLS certpair at startup,  but only write the certificate to disk" type:"switch" model:"OneTimeTLSKey" featured:"false"`
	ServerTLS                *bool            `category:"wallet" name:"Server TLS" description:"Enable TLS for the wallet connection to node RPC server" type:"switch" model:"ServerTLS" featured:"false"`
	WalletRPCListeners       *cli.StringSlice `category:"wallet" name:"Legacy RPC Listeners" description:"addresses for wallet RPC server to listen on" type:"input" inputType:"text" model:"WalletRPCListeners" featured:"false"`
	WalletRPCMaxClients      *int             `category:"wallet" name:"Legacy RPC Max Clients" description:"maximum number of RPC clients allowed for wallet RPC" type:"input" inputType:"number" model:"WalletRPCMaxClients" featured:"false"`
	WalletRPCMaxWebsockets   *int             `category:"wallet" name:"Legacy RPC Max Websockets" description:"maximum number of websocket clients allowed for wallet RPC" type:"input" inputType:"number" model:"WalletRPCMaxWebsockets" featured:"false"`
	ExperimentalRPCListeners *cli.StringSlice `category:"wallet" name:"Experimental RPC Listeners" description:"addresses for experimental RPC listeners to listen on" type:"input" inputType:"text" model:"ExperimentalRPCListeners" featured:"false"`
	NodeOff                  *bool            `category:"debug" name:"NodeOff" description:"turn off the node backend" type:"switch" model:"NodeOff" featured:"false"`
	TestNodeOff              *bool            `category:"debug" name:"TestNodeOff" description:"turn off the testnode (testnet only)" type:"switch" model:"TestNodeOff" featured:"false"`
	WalletOff                *bool            `category:"debug" name:"WalletOff" description:"turn off the wallet backend" type:"switch" model:"WalletOff" featured:"false"`
}

func PodDefConfig() *Config {
	return &Config{
		ConfigFile:               new(string),
		DataDir:                  new(string),
		LogDir:                   new(string),
		Network:                  new(string),
		LogLevel:                 new(string),
		Subsystems:               new(cli.StringSlice),
		NoDiscovery:              new(bool),
		AddPeers:                 new(cli.StringSlice),
		ConnectPeers:             new(cli.StringSlice),
		MaxPeers:                 new(int),
		Listeners:                new(cli.StringSlice),
		DisableListen:            new(bool),
		DisableBanning:           new(bool),
		BanDuration:              new(time.Duration),
		BanThreshold:             new(int),
		Whitelists:               new(cli.StringSlice),
		Username:                 new(string),
		Password:                 new(string),
		ServerUser:               new(string),
		ServerPass:               new(string),
		LimitUser:                new(string),
		LimitPass:                new(string),
		RPCConnect:               new(string),
		RPCListeners:             new(cli.StringSlice),
		RPCCert:                  new(string),
		RPCKey:                   new(string),
		RPCMaxClients:            new(int),
		RPCMaxWebsockets:         new(int),
		RPCMaxConcurrentReqs:     new(int),
		RPCQuirks:                new(bool),
		DisableRPC:               new(bool),
		TLS:                      new(bool),
		DisableDNSSeed:           new(bool),
		ExternalIPs:              new(cli.StringSlice),
		Proxy:                    new(string),
		ProxyUser:                new(string),
		ProxyPass:                new(string),
		OnionProxy:               new(string),
		OnionProxyUser:           new(string),
		OnionProxyPass:           new(string),
		Onion:                    new(bool),
		TorIsolation:             new(bool),
		TestNet3:                 new(bool),
		RegressionTest:           new(bool),
		SimNet:                   new(bool),
		AddCheckpoints:           new(cli.StringSlice),
		DisableCheckpoints:       new(bool),
		DbType:                   new(string),
		Profile:                  new(string),
		CPUProfile:               new(string),
		Upnp:                     new(bool),
		MinRelayTxFee:            new(float64),
		FreeTxRelayLimit:         new(float64),
		NoRelayPriority:          new(bool),
		TrickleInterval:          new(time.Duration),
		MaxOrphanTxs:             new(int),
		Algo:                     new(string),
		Generate:                 new(bool),
		GenThreads:               new(int),
		Controller:               new(string),
		NoController:             new(bool),
		MiningAddrs:              new(cli.StringSlice),
		MinerPass:                new(string),
		Group:                    new(string),
		BlockMinSize:             new(int),
		BlockMaxSize:             new(int),
		BlockMinWeight:           new(int),
		BlockMaxWeight:           new(int),
		BlockPrioritySize:        new(int),
		UserAgentComments:        new(cli.StringSlice),
		NoPeerBloomFilters:       new(bool),
		NoCFilters:               new(bool),
		SigCacheMaxSize:          new(int),
		BlocksOnly:               new(bool),
		TxIndex:                  new(bool),
		AddrIndex:                new(bool),
		RelayNonStd:              new(bool),
		RejectNonStd:             new(bool),
		TLSSkipVerify:            new(bool),
		Wallet:                   new(bool),
		NoInitialLoad:            new(bool),
		WalletPass:               new(string),
		WalletServer:             new(string),
		CAFile:                   new(string),
		OneTimeTLSKey:            new(bool),
		ServerTLS:                new(bool),
		WalletRPCListeners:       new(cli.StringSlice),
		WalletRPCMaxClients:      new(int),
		WalletRPCMaxWebsockets:   new(int),
		ExperimentalRPCListeners: new(cli.StringSlice),
		NodeOff:                  new(bool),
		TestNodeOff:              new(bool),
		WalletOff:                new(bool),
	}
}
