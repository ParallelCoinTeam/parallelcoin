package conte

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
	
	"github.com/p9c/pod/pkg/util/quit"
	
	"go.uber.org/atomic"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/util/lang"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/p9c/pod/pkg/wallet/chain"
)

type _dtype int

var _d _dtype

// Xt as in conte.Xt stores all the common state data used in pod
type Xt struct {
	sync.Mutex
	WaitGroup sync.WaitGroup
	KillAll   qu.C
	// App is the heart of the application system, this creates and initialises it.
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
	NodeReady qu.C
	// NodeKill is the killswitch for the Node
	NodeKill qu.C
	// Wallet is the run state of the wallet
	Wallet atomic.Bool
	// WalletKill is the killswitch for the Wallet
	WalletKill qu.C
	// RPCServer is needed to directly query data
	RPCServer *chainrpc.Server
	// NodeChan relays the chain RPC server to the main
	NodeChan chan *chainrpc.Server
	// WalletServer is needed to query the wallet
	WalletServer *wallet.Wallet
	// // WalletChan is a channel used to return the wallet server pointer when it starts
	// WalletChan chan *wallet.Wallet
	// ChainClientChan returns the chainclient
	ChainClientReady qu.C
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
	// IsGUI indicates if we have the possibility of terminal input
	IsGUI        bool
	UUID         uint64
	waitChangers []string
	waitCounter  int
}

func (cx *Xt) WaitAdd() {
	cx.WaitGroup.Add(1)
	_, file, line, _ := runtime.Caller(1)
	record := fmt.Sprintf("+ %s:%d", file, line)
	cx.waitChangers = append(cx.waitChangers, record)
	cx.waitCounter++
	Debug("added to waitgroup", record, cx.waitCounter)
	Debug(cx.PrintWaitChangers())
}

func (cx *Xt) WaitDone() {
	_, file, line, _ := runtime.Caller(1)
	record := fmt.Sprintf("- %s:%d", file, line)
	cx.waitChangers = append(cx.waitChangers, record)
	cx.waitCounter--
	Debug("removed from waitgroup", record, cx.waitCounter)
	Debug(cx.PrintWaitChangers())
	qu.PrintChanState()
	cx.WaitGroup.Done()
}

func (cx *Xt) WaitWait() {
	Debug(cx.PrintWaitChangers())
	cx.WaitGroup.Wait()
}

func (cx *Xt) PrintWaitChangers() string {
	o := "Calls that change context waitgroup values:\n"
	for i := range cx.waitChangers {
		o += cx.waitChangers[i] + "\n"
	}
	o += "current total:"
	o += fmt.Sprint(cx.waitCounter)
	return o
}

// GetNewContext returns a fresh new context
func GetNewContext(appName, appLang, subtext string) *Xt {
	config, configMap := pod.EmptyConfig()
	chainClientReady := qu.T()
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	cx := &Xt{
		ChainClientReady: chainClientReady,
		KillAll:          qu.T(),
		App:              cli.NewApp(),
		Config:           config,
		ConfigMap:        configMap,
		StateCfg:         new(state.Config),
		Language:         lang.ExportLanguage(appLang),
		DataDir:          appdata.Dir(appName, false),
		NodeChan:         make(chan *chainrpc.Server),
		UUID:             uint64(rand.Int63()),
	}
	return cx
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
	}
	Trace(cc, "nodes connected")
	connected := cc > 0
	is = rn.Chain.IsCurrent() &&
		rn.SyncManager.IsCurrent() &&
		connected &&
		rn.Chain.BestChain.Height() >= rn.HighestKnown.Load() || *cx.Config.Solo
	Trace(
		"is current:", is, "-", rn.Chain.IsCurrent(), rn.SyncManager.IsCurrent(),
		*cx.Config.Solo, "connected", rn.HighestKnown.Load(), rn.Chain.BestChain.Height(),
		othernodes,
	)
	return is
}
