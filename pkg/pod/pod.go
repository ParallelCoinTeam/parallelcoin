package pod

import (
	"reflect"
	"sort"
	"sync"
	"time"

	"github.com/p9c/pod/app/appdata"
	log "github.com/p9c/pod/pkg/logi"

	"github.com/davecgh/go-spew/spew"

	"github.com/urfave/cli"

	"github.com/p9c/pod/pkg/chain/fork"
)

const AppName = "pod"

type Schema struct {
	Groups Groups `json:"groups"`
}
type Groups []Group

type Group struct {
	Legend string `json:"legend"`
	Fields `json:"fields"`
}

type Fields []Field

type Field struct {
	Group       string   `json:"group"`
	Type        string   `json:"type"`
	Label       string   `json:"label"`
	Slug        string   `json:"slug"`
	Description string   `json:"help"`
	InputType   string   `json:"inputType"`
	Featured    string   `json:"featured"`
	Model       string   `json:"model"`
	Datatype    string   `json:"datatype"`
	Options     []string `json:"options"`
}

func GetConfigSchema(cfg *Config, cfgMap map[string]interface{}) Schema {
	t := reflect.TypeOf(*cfg)
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
	rawFields := make(map[string]Fields)
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
			Label:       field.Tag.Get("label"),
			Description: field.Tag.Get("description"),
			InputType:   field.Tag.Get("inputType"),
			Featured:    field.Tag.Get("featured"),
			Options:     options,
			Datatype:    field.Type.String(),
			Model:       field.Tag.Get("json"),
			// Value:       cfgMap[field.Tag.Get("model")],
		}
		if f.Group != "" {
			rawFields[f.Group] = append(rawFields[f.Group], f)
		}
		// groups = append(groups, f.Group)
	}
	spew.Dump(groups)
	var outGroups Groups
	var rf []string
	for i := range rawFields {
		rf = append(rf, i)
	}
	sort.Strings(rf)
	for i := range rf {
		rf[i], rf[len(rf)-1-i] = rf[len(rf)-1-i], rf[i]
	}
	for i := range rf {
		group := Group{
			Legend: rf[i],
			Fields: rawFields[rf[i]],
		}
		outGroups = append(Groups{group}, outGroups...)
	}
	return Schema{
		Groups: outGroups,
	}
}

// Config is
type Config struct {
	sync.Mutex
	AddCheckpoints         *cli.StringSlice `group:"debug" label:"AddCheckpoints" description:"add custom checkpoints" type:"stringSlice" inputType:"text" json:"AddCheckpoints"`
	AddPeers               *cli.StringSlice `group:"node" label:"Add Peers" description:"Manually adds addresses to try to connect to" type:"stringSlice" inputType:"text" json:"AddPeers"`
	AddrIndex              *bool            `group:"node" label:"Addr Index" description:"maintain a full address-based transaction index which makes the searchrawtransactions RPC available" type:"switch" json:"AddrIndex"`
	Algo                   *string          `group:"mining" label:"Algo" description:"algorithm to mine, random is best" type:"input" inputType:"text" json:"Algo"`
	AutoPorts              *bool            `group:"node" label:"Automatic	Ports" description:"with controller enabled p2p, rpc and controller ports are randomized" type:"switch" json:"AutoPorts"`
	BanDuration            *time.Duration   `group:"debug" label:"Ban Duration" description:"how long a ban of a misbehaving peer lasts" type:"input" inputType:"time" json:"BanDuration"`
	BanThreshold           *int             `group:"debug" label:"Ban Threshold" description:"ban score that triggers a ban (default 100)" type:"input" inputType:"number" json:"BanThreshold"`
	BlockMaxSize           *int             `group:"mining" label:"Block Max Size" description:"maximum block size in bytes to be used when creating a block" type:"input" inputType:"number" json:"BlockMaxSize"`
	BlockMaxWeight         *int             `group:"mining" label:"Block Max Weight" description:"maximum block weight to be used when creating a block" type:"input" inputType:"number" json:"BlockMaxWeight"`
	BlockMinSize           *int             `group:"mining" label:"Block Min Size" description:"minimum block size in bytes to be used when creating a block" type:"input" inputType:"number" json:"BlockMinSize"`
	BlockMinWeight         *int             `group:"mining" label:"Block Min Weight" description:"minimum block weight to be used when creating a block" type:"input" inputType:"number" json:"BlockMinWeight"`
	BlockPrioritySize      *int             `group:"mining" label:"Block Priority Size" description:"size in bytes for high-priority/low-fee transactions when creating a block" type:"input" inputType:"number" json:"BlockPrioritySize"`
	BlocksOnly             *bool            `group:"node" label:"Blocks Only" description:"do not accept transactions from remote peers" type:"switch" json:"BlocksOnly"`
	CAFile                 *string          `group:"tls" label:"CA File" description:"certificate authority file for TLS certificate validation" type:"input" inputType:"text" json:"CAFile"`
	ConfigFile             *string          `json:"ConfigFile"`
	ConnectPeers           *cli.StringSlice `group:"node" label:"Connect Peers" description:"Connect ONLY to these addresses (disables inbound connections)" type:"stringSlice" inputType:"text" json:"ConnectPeers"`
	Controller             *string          `group:"mining" label:"Controller Receiver" description:"address to bind miner controller to" json:"Controller"`
	CPUProfile             *string          `group:"debug" label:"CPU Profile" description:"write cpu profile to this file" type:"input" inputType:"text" json:"CPUProfile"`
	DataDir                *string          `group:"config" label:"Data Dir" description:"Root folder where application data is stored" type:"input" inputType:"text" json:"DataDir"`
	DbType                 *string          `group:"debug" label:"Db Type" description:"type of database storage engine to use (only one right now)" type:"input" inputType:"text" json:"DbType"`
	DisableBanning         *bool            `group:"debug" label:"Disable Banning" description:"Disables banning of misbehaving peers" type:"switch" json:"DisableBanning"`
	DisableCheckpoints     *bool            `group:"debug" label:"Disable Checkpoints" description:"disables all checkpoints" type:"switch" json:"DisableCheckpoints"`
	DisableDNSSeed         *bool            `group:"node" label:"Disable DNS Seed" description:"disable seeding of addresses to peers" type:"switch" json:"DisableDNSSeed"`
	DisableListen          *bool            `group:"node" label:"Disable Listen" description:"Disables inbound connections for the peer to peer network" type:"switch" json:"DisableListen"`
	DisableRPC             *bool            `group:"rpc" label:"Disable RPC" description:"disable rpc servers" type:"switch" json:"DisableRPC"`
	ExternalIPs            *cli.StringSlice `group:"node" label:"External IPs" description:"extra addresses to tell peers they can connect to" type:"stringSlice" inputType:"text" json:"ExternalIPs"`
	FreeTxRelayLimit       *float64         `group:"policy" label:"Free Tx Relay Limit" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute" type:"input" inputType:"decimal" json:"FreeTxRelayLimit"`
	Generate               *bool            `group:"mining" label:"Generate" description:"turn on built in CPU miner" type:"switch" json:"Generate"`
	GenThreads             *int             `group:"mining" label:"Gen Threads" description:"number of CPU threads to mine using" type:"input" inputType:"number" json:"GenThreads"`
	Language               *string          `group:"config" label:"Language" description:"User interface language i18 localization" type:"input" inputType:"text" json:"Language"`
	LimitPass              *string          `group:"rpc" label:"Limit Pass" description:"limited user password" type:"input" inputType:"password" json:"LimitPass"`
	LimitUser              *string          `group:"rpc" label:"Limit User" description:"limited user name" type:"input" inputType:"text" json:"LimitUser"`
	Listeners              *cli.StringSlice `group:"node" label:"Listeners" description:"List of addresses to bind the node listener to" type:"stringSlice" inputType:"text" json:"Listeners"`
	LogDir                 *string          `group:"config" label:"Log Dir" description:"Folder where log files are written" type:"input" inputType:"text" json:"LogDir"`
	LogLevel               *string          `group:"config" label:"Log Level" description:"Verbosity of log printouts" type:"input" inputType:"text" json:"LogLevel"`
	MaxOrphanTxs           *int             `group:"policy" label:"Max Orphan Txs" description:"max number of orphan transactions to keep in memory" type:"input" inputType:"number" json:"MaxOrphanTxs"`
	MaxPeers               *int             `group:"node" label:"Max Peers" description:"Maximum number of peers to hold connections with" type:"input" inputType:"number" json:"MaxPeers"`
	MinerPass              *string          `group:"mining" label:"Miner Pass" description:"password that encrypts the connection to the mining controller" type:"input" inputType:"text" json:"MinerPass"`
	MiningAddrs            *cli.StringSlice `group:"mining" label:"Mining Addrs" description:"addresses to pay block rewards to (TODO, make this auto)" type:"stringSlice" inputType:"text" json:"MiningAddrs"`
	MinRelayTxFee          *float64         `group:"policy" label:"Min Relay Tx Fee" description:"the minimum transaction fee in DUO/kB to be considered a non-zero fee" type:"input" inputType:"decimal" json:"MinRelayTxFee"`
	Network                *string          `group:"node" label:"Network" description:"Which network are you connected to (eg.: mainnet, testnet)" type:"input" inputType:"text" json:"Network"`
	NoCFilters             *bool            `group:"node" label:"No CFilters" description:"disable committed filtering (CF) support" type:"switch" json:"NoCFilters"`
	NodeOff                *bool            `group:"debug" label:"Node Off" description:"turn off the node backend" type:"switch" json:"NodeOff"`
	NoInitialLoad          *bool            `json:"NoInitialLoad"`
	NoPeerBloomFilters     *bool            `group:"node" label:"No Peer Bloom Filters" description:"disable bloom filtering support" type:"switch" json:"NoPeerBloomFilters"`
	NoRelayPriority        *bool            `group:"policy" label:"No Relay Priority" description:"do not require free or low-fee transactions to have high priority for relaying" type:"switch" json:"NoRelayPriority"`
	OneTimeTLSKey          *bool            `group:"wallet" label:"One Time TLS Key" description:"generate a new TLS certpair at startup, but only write the certificate to disk" type:"switch" json:"OneTimeTLSKey"`
	Onion                  *bool            `group:"proxy" label:"Onion" description:"enable tor proxy" type:"switch" json:"Onion"`
	OnionProxy             *string          `group:"proxy" label:"Onion Proxy" description:"address of tor proxy you want to connect to" type:"input" inputType:"text" json:"OnionProxy"`
	OnionProxyPass         *string          `group:"proxy" label:"Onion Proxy Pass" description:"password for tor proxy" type:"input" inputType:"password" json:"OnionProxyPass"`
	OnionProxyUser         *string          `group:"proxy" label:"Onion Proxy User" description:"tor proxy username" type:"input" inputType:"text" json:"OnionProxyUser"`
	Password               *string          `group:"rpc" label:"Password" description:"password for client RPC connections" type:"input" inputType:"text" json:"Password"`
	Profile                *string          `group:"debug" label:"Profile" description:"http profiling on given port (1024-40000)" type:"input" inputType:"text" json:"Profile"`
	Proxy                  *string          `group:"proxy" label:"Proxy" description:"address of proxy to connect to for outbound connections" type:"input" inputType:"text" json:"Proxy"`
	ProxyPass              *string          `group:"proxy" label:"Proxy Pass" description:"proxy password, if required" type:"input" inputType:"password" json:"ProxyPass"`
	ProxyUser              *string          `group:"proxy" label:"ProxyUser" description:"proxy username, if required" type:"input" inputType:"text" json:"ProxyUser"`
	RejectNonStd           *bool            `group:"node" label:"Reject Non Std" description:"reject non-standard transactions regardless of the default settings for the active network" type:"switch" json:"RejectNonStd"`
	RelayNonStd            *bool            `group:"node" label:"Relay Non Std" description:"relay non-standard transactions regardless of the default settings for the active network" type:"switch" json:"RelayNonStd"`
	RPCCert                *string          `group:"rpc" label:"RPC Cert" description:"location of rpc TLS certificate" type:"input" inputType:"text" json:"RPCCert"`
	RPCConnect             *string          `group:"wallet" label:"RPC Connect" description:"full node RPC for wallet" type:"input" inputType:"text" json:"RPCConnect"`
	RPCKey                 *string          `group:"rpc" label:"RPC Key" description:"location of rpc TLS key" type:"input" inputType:"text" json:"RPCKey"`
	RPCListeners           *cli.StringSlice `group:"rpc" label:"RPC Listeners" description:"addresses to listen for RPC connections" type:"stringSlice" inputType:"text" json:"RPCListeners"`
	RPCMaxClients          *int             `group:"rpc" label:"RPC Max Clients" description:"maximum number of clients for regular RPC" type:"input" inputType:"number" json:"RPCMaxClients"`
	RPCMaxConcurrentReqs   *int             `group:"rpc" label:"RPC Max Concurrent Reqs" description:"maximum number of requests to process concurrently" type:"input" inputType:"number" json:"RPCMaxConcurrentReqs"`
	RPCMaxWebsockets       *int             `group:"rpc" label:"RPC Max Websockets" description:"maximum number of websocket clients to allow" type:"input" inputType:"number" json:"RPCMaxWebsockets"`
	RPCQuirks              *bool            `group:"rpc" label:"RPC Quirks" description:"enable bugs that replicate bitcoin core RPC's JSON" type:"switch" json:"RPCQuirks"`
	ServerPass             *string          `group:"rpc" label:"Server Pass" description:"password for server connections" type:"input" inputType:"password" json:"ServerPass"`
	ServerTLS              *bool            `group:"wallet" label:"Server TLS" description:"Enable TLS for the wallet connection to node RPC server" type:"switch" json:"ServerTLS"`
	ServerUser             *string          `group:"rpc" label:"Server User" description:"username for server connections" type:"input" inputType:"text" json:"ServerUser"`
	SigCacheMaxSize        *int             `group:"node" label:"Sig Cache Max Size" description:"the maximum number of entries in the signature verification cache" type:"input" inputType:"number" json:"SigCacheMaxSize"`
	Solo                   *bool            `group:"mining" label:"Solo Generate" description:"mine even if not connected to a network" type:"switch" json:"Solo"`
	TLS                    *bool            `group:"tls" label:"TLS" description:"enable TLS for RPC connections" type:"switch" json:"TLS"`
	TLSSkipVerify          *bool            `group:"tls" label:"TLS Skip Verify" description:"skip TLS certificate verification (ignore CA errors)" type:"switch" json:"TLSSkipVerify"`
	TorIsolation           *bool            `group:"proxy" label:"Tor Isolation" description:"makes a separate proxy connection for each connection" type:"switch" json:"TorIsolation"`
	TrickleInterval        *time.Duration   `group:"policy" label:"Trickle Interval" description:"minimum time between attempts to send new inventory to a connected peer" type:"input" inputType:"time" json:"TrickleInterval"`
	TxIndex                *bool            `group:"node" label:"Tx Index" description:"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC" type:"switch" json:"TxIndex"`
	UPNP                   *bool            `group:"node" label:"UPNP" description:"enable UPNP for NAT traversal" type:"switch" json:"UPNP"`
	UserAgentComments      *cli.StringSlice `group:"node" label:"User Agent Comments" description:"Comment to add to the user agent -- See BIP 14 for more information" type:"stringSlice" inputType:"text" json:"UserAgentComments"`
	Username               *string          `group:"rpc" label:"Username" description:"password for client RPC connections" type:"input" inputType:"text" json:"Username"`
	Wallet                 *bool            `json:"Wallet"`
	WalletFile             *string          `group:"config" label:"Wallet File" description:"Wallet database file" type:"input" inputType:"text" featured:"true" json:"WalletFile"`
	WalletOff              *bool            `group:"debug" label:"Wallet Off" description:"turn off the wallet backend" type:"switch" json:"WalletOff"`
	WalletPass             *string          `group:"wallet" label:"Wallet Pass" description:"password encrypting public data in wallet" type:"input" inputType:"text" json:"WalletPass"`
	WalletRPCListeners     *cli.StringSlice `group:"wallet" label:"Legacy RPC Listeners" description:"addresses for wallet RPC server to listen on" type:"stringSlice" inputType:"text" json:"WalletRPCListeners"`
	WalletRPCMaxClients    *int             `group:"wallet" label:"Legacy RPC Max Clients" description:"maximum number of RPC clients allowed for wallet RPC" type:"input" inputType:"number" json:"WalletRPCMaxClients"`
	WalletRPCMaxWebsockets *int             `group:"wallet" label:"Legacy RPC Max Websockets" description:"maximum number of websocket clients allowed for wallet RPC" type:"input" inputType:"number" json:"WalletRPCMaxWebsockets"`
	WalletServer           *string          `group:"wallet" label:"Wallet Server" description:"node address to connect wallet server to" type:"input" inputType:"text" json:"WalletServer"`
	Whitelists             *cli.StringSlice `group:"debug" label:"Whitelists" description:"peers that you don't want to ever ban" type:"stringSlice" inputType:"text" json:"Whitelists"`
	LAN                    *bool            `json:"LAN"`
}

func EmptyConfig() (c *Config, conf map[string]interface{}) {
	datadir := appdata.Dir(AppName, false)
	c = &Config{
		AddCheckpoints:         newStringSlice(),
		AddPeers:               newStringSlice(),
		AddrIndex:              newbool(),
		Algo:                   newstring(),
		AutoPorts:              newbool(),
		BanDuration:            newDuration(),
		BanThreshold:           newint(),
		BlockMaxSize:           newint(),
		BlockMaxWeight:         newint(),
		BlockMinSize:           newint(),
		BlockMinWeight:         newint(),
		BlockPrioritySize:      newint(),
		BlocksOnly:             newbool(),
		CAFile:                 newstring(),
		ConfigFile:             newstring(),
		ConnectPeers:           newStringSlice(),
		Controller:             newstring(),
		CPUProfile:             newstring(),
		DataDir:                &datadir,
		DbType:                 newstring(),
		DisableBanning:         newbool(),
		DisableCheckpoints:     newbool(),
		DisableDNSSeed:         newbool(),
		DisableListen:          newbool(),
		DisableRPC:             newbool(),
		ExternalIPs:            newStringSlice(),
		FreeTxRelayLimit:       new(float64),
		Generate:               newbool(),
		GenThreads:             newint(),
		LAN:                    newbool(),
		Language:               newstring(),
		LimitPass:              newstring(),
		LimitUser:              newstring(),
		Listeners:              newStringSlice(),
		LogDir:                 newstring(),
		LogLevel:               newstring(),
		MaxOrphanTxs:           newint(),
		MaxPeers:               newint(),
		MinerPass:              newstring(),
		MiningAddrs:            newStringSlice(),
		MinRelayTxFee:          new(float64),
		Network:                newstring(),
		NoCFilters:             newbool(),
		NodeOff:                newbool(),
		NoInitialLoad:          newbool(),
		NoPeerBloomFilters:     newbool(),
		NoRelayPriority:        newbool(),
		OneTimeTLSKey:          newbool(),
		Onion:                  newbool(),
		OnionProxy:             newstring(),
		OnionProxyPass:         newstring(),
		OnionProxyUser:         newstring(),
		Password:               newstring(),
		Profile:                newstring(),
		Proxy:                  newstring(),
		ProxyPass:              newstring(),
		ProxyUser:              newstring(),
		RejectNonStd:           newbool(),
		RelayNonStd:            newbool(),
		RPCCert:                newstring(),
		RPCConnect:             newstring(),
		RPCKey:                 newstring(),
		RPCListeners:           newStringSlice(),
		RPCMaxClients:          newint(),
		RPCMaxConcurrentReqs:   newint(),
		RPCMaxWebsockets:       newint(),
		RPCQuirks:              newbool(),
		ServerPass:             newstring(),
		ServerTLS:              newbool(),
		ServerUser:             newstring(),
		SigCacheMaxSize:        newint(),
		Solo:                   newbool(),
		TLS:                    newbool(),
		TLSSkipVerify:          newbool(),
		TorIsolation:           newbool(),
		TrickleInterval:        newDuration(),
		TxIndex:                newbool(),
		UPNP:                   newbool(),
		UserAgentComments:      newStringSlice(),
		Username:               newstring(),
		Wallet:                 newbool(),
		WalletFile:             newstring(),
		WalletOff:              newbool(),
		WalletPass:             newstring(),
		WalletRPCListeners:     newStringSlice(),
		WalletRPCMaxClients:    newint(),
		WalletRPCMaxWebsockets: newint(),
		WalletServer:           newstring(),
		Whitelists:             newStringSlice(),
	}
	conf = map[string]interface{}{
		"AddCheckpoints":         c.AddCheckpoints,
		"AddPeers":               c.AddPeers,
		"AddrIndex":              c.AddrIndex,
		"Algo":                   c.Algo,
		"AutoPorts":              c.AutoPorts,
		"BanDuration":            c.BanDuration,
		"BanThreshold":           c.BanThreshold,
		"BlockMaxSize":           c.BlockMaxSize,
		"BlockMaxWeight":         c.BlockMaxWeight,
		"BlockMinSize":           c.BlockMinSize,
		"BlockMinWeight":         c.BlockMinWeight,
		"BlockPrioritySize":      c.BlockPrioritySize,
		"BlocksOnly":             c.BlocksOnly,
		"CAFile":                 c.CAFile,
		"ConfigFile":             c.ConfigFile,
		"ConnectPeers":           c.ConnectPeers,
		"Controller":             c.Controller,
		"CPUProfile":             c.CPUProfile,
		"DataDir":                c.DataDir,
		"DbType":                 c.DbType,
		"DisableBanning":         c.DisableBanning,
		"DisableCheckpoints":     c.DisableCheckpoints,
		"DisableDNSSeed":         c.DisableDNSSeed,
		"DisableListen":          c.DisableListen,
		"DisableRPC":             c.DisableRPC,
		"ExternalIPs":            c.ExternalIPs,
		"FreeTxRelayLimit":       c.FreeTxRelayLimit,
		"Generate":               c.Generate,
		"GenThreads":             c.GenThreads,
		"LAN":                    c.LAN,
		"Language":               c.Language,
		"LimitPass":              c.LimitPass,
		"LimitUser":              c.LimitUser,
		"Listeners":              c.Listeners,
		"LogDir":                 c.LogDir,
		"LogLevel":               c.LogLevel,
		"MaxOrphanTxs":           c.MaxOrphanTxs,
		"MaxPeers":               c.MaxPeers,
		"MinerPass":              c.MinerPass,
		"MiningAddrs":            c.MiningAddrs,
		"MinRelayTxFee":          c.MinRelayTxFee,
		"Network":                c.Network,
		"NoCFilters":             c.NoCFilters,
		"NodeOff":                c.NodeOff,
		"NoInitialLoad":          c.NoInitialLoad,
		"NoPeerBloomFilters":     c.NoPeerBloomFilters,
		"NoRelayPriority":        c.NoRelayPriority,
		"OneTimeTLSKey":          c.OneTimeTLSKey,
		"Onion":                  c.Onion,
		"OnionProxy":             c.OnionProxy,
		"OnionProxyPass":         c.OnionProxyPass,
		"OnionProxyUser":         c.OnionProxyUser,
		"Password":               c.Password,
		"Profile":                c.Profile,
		"Proxy":                  c.Proxy,
		"ProxyPass":              c.ProxyPass,
		"ProxyUser":              c.ProxyUser,
		"RejectNonStd":           c.RejectNonStd,
		"RelayNonStd":            c.RelayNonStd,
		"RPCCert":                c.RPCCert,
		"RPCConnect":             c.RPCConnect,
		"RPCKey":                 c.RPCKey,
		"RPCListeners":           c.RPCListeners,
		"RPCMaxClients":          c.RPCMaxClients,
		"RPCMaxConcurrentReqs":   c.RPCMaxConcurrentReqs,
		"RPCMaxWebsockets":       c.RPCMaxWebsockets,
		"RPCQuirks":              c.RPCQuirks,
		"ServerPass":             c.ServerPass,
		"ServerTLS":              c.ServerTLS,
		"ServerUser":             c.ServerUser,
		"SigCacheMaxSize":        c.SigCacheMaxSize,
		"Solo":                   c.Solo,
		"TLS":                    c.TLS,
		"TLSSkipVerify":          c.TLSSkipVerify,
		"TorIsolation":           c.TorIsolation,
		"TrickleInterval":        c.TrickleInterval,
		"TxIndex":                c.TxIndex,
		"UPNP":                   c.UPNP,
		"UserAgentComments":      c.UserAgentComments,
		"Username":               c.Username,
		"Wallet":                 c.Wallet,
		"WalletFile":             c.WalletFile,
		"WalletOff":              c.WalletOff,
		"WalletPass":             c.WalletPass,
		"WalletRPCListeners":     c.WalletRPCListeners,
		"WalletRPCMaxClients":    c.WalletRPCMaxClients,
		"WalletRPCMaxWebsockets": c.WalletRPCMaxWebsockets,
		"WalletServer":           c.WalletServer,
		"Whitelists":             c.Whitelists,
	}
	return
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
