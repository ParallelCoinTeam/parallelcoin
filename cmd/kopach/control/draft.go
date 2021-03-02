package control

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/chain/mining"
	qu "github.com/p9c/pod/pkg/util/quit"
)

// State stores the state of the controller
type State struct {
	cx                *conte.Xt
	start, stop, quit qu.C
	generator         *mining.BlkTmplGenerator
}

// New creates a new controller
func New(cx *conte.Xt) (s *State) {
	s = &State{
		cx:        cx,
		quit:      qu.T(),
		start:     qu.T(),
		stop:      qu.T(),
		generator: getBlkTemplateGenerator(cx),
	}
	go func() {
		// this ensures that if the main app is shutting down that all goroutines are
		// stopped
		select {
		case <-cx.KillAll:
			s.quit.Q()
		}
	}()
	return
}

// Start up the controller
func (s *State) Start() {
	Debug("calling start controller")
	s.start.Signal()
}

// Stop the controller
func (s *State) Stop() {
	Debug("stop")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	Debug("quit")
	s.quit.Q()
}

// Run must be start as a goroutine, central routing for the business of the
// controller
func (s *State) Run() {
	if *s.cx.Config.DisableController {
		Warn("controller is disabled")
		return
	}
	if len(*s.cx.Config.RPCListeners) < 1 || *s.cx.Config.DisableRPC {
		Warn("not running controller without RPC enabled")
		return
	}
	if len(*s.cx.Config.P2PListeners) < 1 || *s.cx.Config.DisableListen {
		Warn("not running controller without p2p listener enabled", *s.cx.Config.P2PListeners)
		return
	}
out:
	for {
	pausing:
		for {
			select {
			case <-s.start.Wait():
				Debug("received start signal while paused")
				break pausing
			case <-s.stop.Wait():
				Debug("received stop signal while paused")
				break
			case <-s.quit.Wait():
				Debug("received quit signal whiel paused")
				break out
			}
			Debug("establishing wallet connection")
			
			Debug("wallet is connected, switching to running")
			
			Debug("wallet did not connect, pausing before retry")
		}
	running:
		for {
			select {
			case <-s.start.Wait():
				Debug("received start signal while running")
				break
			case <-s.stop.Wait():
				Debug("received stop signal while running")
				break running
			case <-s.quit.Wait():
				Debug("received quit signal while running")
				break out
			}
			Debug("checking if wallet is connected")
			
			Debug("wallet not connected, switching to pausing")
		}
	}
}

func getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
	policy := mining.Policy{
		BlockMinWeight:    uint32(*cx.Config.BlockMinWeight),
		BlockMaxWeight:    uint32(*cx.Config.BlockMaxWeight),
		BlockMinSize:      uint32(*cx.Config.BlockMinSize),
		BlockMaxSize:      uint32(*cx.Config.BlockMaxSize),
		BlockPrioritySize: uint32(*cx.Config.BlockPrioritySize),
		TxMinFreeFee:      cx.StateCfg.ActiveMinRelayTxFee,
	}
	s := cx.RealNode
	return mining.NewBlkTmplGenerator(
		&policy,
		s.ChainParams,
		s.TxMemPool,
		s.Chain,
		s.TimeSource,
		s.SigCache,
		s.HashCache,
	)
}
