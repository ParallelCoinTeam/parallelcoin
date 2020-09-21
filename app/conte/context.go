// +build !headless

package conte

import (
	"github.com/stalker-loki/app/slog"
	"sync"

	"go.uber.org/atomic"

	"github.com/urfave/cli"

	"github.com/stalker-loki/pod/app/appdata"
	"github.com/stalker-loki/pod/cmd/node/state"
	"github.com/stalker-loki/pod/pkg/chain/config/netparams"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/stalker-loki/pod/pkg/rpc/chainrpc"
	"github.com/stalker-loki/pod/pkg/util/lang"
	"github.com/stalker-loki/pod/pkg/wallet"
	"github.com/stalker-loki/pod/pkg/wallet/chain"
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
	// AppContext is the urfave/cli app context
	AppContext *cli.Context
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
	Node atomic.Bool
	// NodeReady is closed when it is ready then always returns
	NodeReady chan struct{}
	// NodeKill is the killswitch for the Node
	NodeKill chan struct{}
	// Wallet is the run state of the wallet
	Wallet atomic.Bool
	// WalletKill is the killswitch for the Wallet
	WalletKill chan struct{}
	// RPCServer is needed to directly query data
	RPCServer *chainrpc.Server
	NodeChan  chan *chainrpc.Server
	// WalletServer is needed to query the wallet
	WalletServer *wallet.Wallet
	// WalletChan is a channel used to return the wallet server pointer when it starts
	WalletChan chan *wallet.Wallet
	// ChainClientChan returns the chainclient
	ChainClientReady chan struct{}
	// ChainClient is the wallet's chain RPC client
	ChainClient *chain.RPCClient
	// RealNode is the main node
	RealNode *chainrpc.Node
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
	chainClientReady := make(chan struct{})
	return &Xt{
		ChainClientReady: chainClientReady,
		KillAll:          make(chan struct{}),
		App:              cli.NewApp(),
		Config:           config,
		ConfigMap:        configMap,
		StateCfg:         new(state.Config),
		Language:         lang.ExportLanguage(appLang),
		DataDir:          appdata.Dir(appName, false),
	}
}

func GetContext(cx *Xt) *chainrpc.Context {
	return &chainrpc.Context{
		Config: cx.Config, StateCfg: cx.StateCfg, ActiveNet: cx.ActiveNet,
		Hashrate: cx.Hashrate,
	}
}

func (cx *Xt) IsCurrent() (is bool) {
	rn := cx.RealNode
	cc := rn.ConnectedCount()
	othernodes := cx.OtherNodes.Load()
	if !*cx.Config.LAN {
		cc -= othernodes
		// Debug("LAN disabled, non-lan node count:", cc)
	}
	// Debug("LAN enabled", *cx.Config.LAN, "othernodes", othernodes, "node's connect count", cc)
	connected := cc > 0
	if *cx.Config.Solo {
		connected = true
	}
	is = rn.Chain.IsCurrent() &&
		rn.SyncManager.IsCurrent() &&
		connected &&
		rn.Chain.BestChain.Height() >= rn.HighestKnown.Load()
	slog.Trace("is current:", is, "-", rn.Chain.IsCurrent(),
		rn.SyncManager.IsCurrent(), !*cx.Config.Solo,
		"connected", rn.HighestKnown.Load(),
		rn.Chain.BestChain.Height(),
	)
	return is
}
