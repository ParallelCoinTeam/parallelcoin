//+build !headless

package conte

import (
	"sync"
	
	"go.uber.org/atomic"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/lang"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/wallet"
)

type _dtype int

var _d _dtype

// Xt as in conte.Xt stores all the common state data used in pod
type Xt struct {
	sync.Mutex
	WaitGroup sync.WaitGroup
	KillAll   chan struct{}
	// App is the heart of the application system,
	// this creates and initialises it.
	App *cli.App
	// Config is the pod all-in-one server config
	Config *pod.Config
	// ConfigMap
	ConfigMap map[string]interface{}
	// StateCfg is a reference to the main node state configuration struct
	StateCfg *state.Config
	// ActiveNet is the active net parameters
	ActiveNet *netparams.Params
	// Language libraries
	Language *lang.Lexicon
	// DataDir is the default data dir
	DataDir string
	// Node is the run state of the node
	Node atomic.Value
	// NodeReady is closed when it is ready then always returns
	NodeReady chan struct{}
	// NodeKill is the killswitch for the Node
	NodeKill chan struct{}
	// Wallet is the run state of the wallet
	Wallet atomic.Value
	// WalletKill is the killswitch for the Wallet
	WalletKill chan struct{}
	// RPCServer is needed to directly query data
	RPCServer *rpc.Server
	NodeChan  chan *rpc.Server
	// WalletServer is needed to query the wallet
	WalletServer *wallet.Wallet
	// WalletChan is a channel used to return the wallet server pointer when it starts
	WalletChan chan *wallet.Wallet
	// RealNode is the main node
	RealNode *rpc.Node
	// Hashrate is the current total hashrate from kopach workers taking work from this node
	Hashrate atomic.Uint64
	// Controller is the run state indicator of the controller
	Controller atomic.Bool
	// OtherNodes is the count of nodes connected automatically on the LAN
	OtherNodes atomic.Int32
}

// GetNewContext returns a fresh new context
func GetNewContext(appName, appLang, subtext string) *Xt {
	hr := &atomic.Value{}
	hr.Store(int(0))
	config, configMap := pod.EmptyConfig()
	return &Xt{
		KillAll:   make(chan struct{}),
		App:       cli.NewApp(),
		Config:    config,
		ConfigMap: configMap,
		StateCfg:  new(state.Config),
		Language:  lang.ExportLanguage(appLang),
		DataDir:   appdata.Dir(appName, false),
	}
}

func GetContext(cx *Xt) *rpc.Context {
	return &rpc.Context{
		Config: cx.Config, StateCfg: cx.StateCfg, ActiveNet: cx.ActiveNet,
		Hashrate: cx.Hashrate,
	}
}

func (cx *Xt) IsCurrent() (is bool) {
	cc := cx.RealNode.ConnectedCount()
	othernodes := cx.OtherNodes.Load()
	if !*cx.Config.LAN {
		cc -= othernodes
		// log.DEBUG("LAN disabled, non-lan node count:", cc)
	}
	// log.DEBUG("LAN enabled", *cx.Config.LAN, "othernodes", othernodes, "node's connect count", cc)
	connected := cc > 0
	if *cx.Config.Solo {
		connected = true
	}
	is = cx.RealNode.Chain.IsCurrent() && cx.RealNode.SyncManager.IsCurrent() &&
		connected
	
	log.TRACE("is current:",is, "-", cx.
		RealNode.Chain.IsCurrent(), cx.
		RealNode.SyncManager.IsCurrent(), !*cx.
		Config.Solo,
		"connected", cx.RealNode.ConnectedCount(), cx.RealNode.ConnectedCount() > 0)
	return is
}
