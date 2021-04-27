// +build headless

package state

import (
	"fmt"
	"math/rand"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/p9c/matrjoska/cmd/ctrl"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainclient"
	"github.com/p9c/matrjoska/pod/config"
	"github.com/p9c/matrjoska/pod/podcfgs"

	"github.com/p9c/qu"

	"go.uber.org/atomic"

	"github.com/p9c/matrjoska/cmd/node/active"
	"github.com/p9c/matrjoska/pkg/chainrpc"
	"github.com/p9c/matrjoska/pkg/util/lang"
)

type _dtype int

var _d _dtype

// State stores all the common state data used in pod
type State struct {
	sync.Mutex
	WaitGroup sync.WaitGroup
	KillAll   qu.C
	Config *config.Config
	// ConfigMap
	ConfigMap map[string]interface{}
	// StateCfg is a reference to the main node state configuration struct
	StateCfg *active.Config
	// ActiveNet is the active net parameters
	ActiveNet *chaincfg.Params
	// Language libraries
	Language *lang.Lexicon
	// // DataDir is the default data dir
	// DataDir string
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
	// // WalletServer is needed to query the wallet
	// WalletServer *wallet.Wallet
	// ChainClientReady signals when the chain client is ready
	ChainClientReady qu.C
	// ChainClient is the wallet's chain RPC client
	ChainClient *chainclient.RPCClient
	// RealNode is the main node
	RealNode *chainrpc.Node
	// Hashrate is the current total hashrate from kopach workers taking work from this node
	Hashrate atomic.Uint64
	// Controller is the state of the controller
	Controller *ctrl.State
	// OtherNodesCounter is the count of nodes connected automatically on the LAN
	OtherNodesCounter atomic.Int32
	// IsGUI indicates if we have the possibility of terminal input
	IsGUI        bool
	waitChangers []string
	waitCounter  int
	Syncing      *atomic.Bool
}

func (cx *State) WaitAdd() {
	cx.WaitGroup.Add(1)
	_, file, line, _ := runtime.Caller(1)
	record := fmt.Sprintf("+ %s:%d", file, line)
	cx.waitChangers = append(cx.waitChangers, record)
	cx.waitCounter++
	D.Ln("added to waitgroup", record, cx.waitCounter)
	D.Ln(cx.PrintWaitChangers())
}

func (cx *State) WaitDone() {
	_, file, line, _ := runtime.Caller(1)
	record := fmt.Sprintf("- %s:%d", file, line)
	cx.waitChangers = append(cx.waitChangers, record)
	cx.waitCounter--
	D.Ln("removed from waitgroup", record, cx.waitCounter)
	D.Ln(cx.PrintWaitChangers())
	qu.PrintChanState()
	cx.WaitGroup.Done()
}

func (cx *State) WaitWait() {
	D.Ln(cx.PrintWaitChangers())
	cx.WaitGroup.Wait()
}

func (cx *State) PrintWaitChangers() string {
	o := "Calls that change context waitgroup values:\n"
	for i := range cx.waitChangers {
		o += strings.Repeat(" ", 48)
		o += cx.waitChangers[i] + "\n"
	}
	o += strings.Repeat(" ", 48)
	o += "current total:"
	o += fmt.Sprint(cx.waitCounter)
	return o
}

// GetNewContext returns a fresh new context
func GetNewContext() (s *State, e error) {
	config := podcfgs.GetDefaultConfig()
	if e = config.Initialize(); E.Chk(e) {
		return
	}
	chainClientReady := qu.T()
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	s = &State{
		ChainClientReady: chainClientReady,
		KillAll:          qu.T(),
		App:              cli.NewApp(),
		Config:           config,
		ConfigMap:        config.Map,
		StateCfg:         new(active.Config),
		// Language:         lang.ExportLanguage(appLang),
		// DataDir:          appdata.Dir(appName, false),
		NodeChan: make(chan *chainrpc.Server),
		Syncing:  atomic.NewBool(true),
	}
	return
}

func GetContext(cx *State) *chainrpc.Context {
	return &chainrpc.Context{
		Config: cx.Config, StateCfg: cx.StateCfg, ActiveNet: cx.ActiveNet,
		Hashrate: cx.Hashrate,
	}
}

func (cx *State) IsCurrent() (is bool) {
	rn := cx.RealNode
	cc := rn.ConnectedCount()
	othernodes := cx.OtherNodesCounter.Load()
	if !*cx.Config.LAN {
		cc -= othernodes
	}
	D.Ln(cc, "nodes connected")
	connected := cc > 0
	is = rn.Chain.IsCurrent() &&
		rn.SyncManager.IsCurrent() &&
		connected &&
		rn.Chain.BestChain.Height() >= rn.HighestKnown.Load() || *cx.Config.Solo
	D.Ln(
		"is current:", is, "-", rn.Chain.IsCurrent(), rn.SyncManager.IsCurrent(),
		*cx.Config.Solo, "connected", rn.HighestKnown.Load(), rn.Chain.BestChain.Height(),
		othernodes,
	)
	return is
}
