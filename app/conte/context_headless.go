// +build headless

package conte

import (
	"sync"

	"github.com/urfave/cli"
	"go.uber.org/atomic"

	"github.com/stalker-loki/pod/app/appdata"
	"github.com/stalker-loki/pod/cmd/node/state"
	"github.com/stalker-loki/pod/pkg/chain/config/netparams"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/stalker-loki/pod/pkg/rpc/chainrpc"
	"github.com/stalker-loki/pod/pkg/util/lang"
	"github.com/stalker-loki/pod/pkg/wallet"
)

type _dtype int

var _d _dtype

// Xt as in conte.Xt stores all the common state data used in pod
type Xt struct {
	sync.Mutex
	// App is the heart of the application system,
	// this creates and initialises it.
	App *cli.App
	// Config is the pod all-in-one server config
	Config *pod.Config
	// StateCfg is a reference to the main node state configuration struct
	StateCfg *state.Config
	// ActiveNet is the active net parameters
	ActiveNet *netparams.Params
	// Language libraries
	Language *lang.Lexicon
	// DataDir is the default data dir
	DataDir string
	// Node is the run state of the node
	Node *atomic.Value
	// NodeKill is the killswitch for the Node
	NodeKill chan struct{}
	// Wallet is the run state of the wallet
	Wallet *atomic.Value
	// WalletKill is the killswitch for the Wallet
	WalletKill chan struct{}
	// // Window is the fyne window when running GUI
	// Window *fyne.Window
	// RPCServer is needed to directly query data
	RPCServer *chainrpc.Server
	// WalletServer is needed to query the wallet
	WalletServer *wallet.Wallet
	// RealNode is the main node
	RealNode *chainrpc.Node
}

// GetNewContext returns a fresh new context
func GetNewContext(appName, appLang, subtext string) *Xt {
	return &Xt{
		App:      cli.NewApp(),
		Config:   pod.EmptyConfig(),
		StateCfg: new(state.Config),
		Language: lang.ExportLanguage(appLang),
		DataDir:  appdata.Dir(appName, false),
	}
}
