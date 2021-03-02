package control

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/walletmain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	qu "github.com/p9c/pod/pkg/util/quit"
	"time"
)

type BlockUpdate struct {
	hash   chainhash.Hash
	height int32
	t      time.Time
}

// State stores the state of the controller
type State struct {
	cx                *conte.Xt
	start, stop, quit qu.C
	running           bool
	blockUpdate       chan *BlockUpdate
	generator         *mining.BlkTmplGenerator
	walletClient      *rpcclient.Client
}

// New creates a new controller
func New(cx *conte.Xt) (s *State) {
	s = &State{
		cx:          cx,
		quit:        qu.T(),
		start:       qu.Ts(1),
		stop:        qu.Ts(1),
		blockUpdate: make(chan *BlockUpdate, 1),
		generator:   getBlkTemplateGenerator(cx),
	}
	go func() {
		Debug("starting shutdown signal watcher")
		select {
		case <-cx.KillAll:
			Debug("received killall signal, signalling to quit controller")
			s.Shutdown()
		case <-cx.NodeKill:
			Debug("received nodekill signal, signalling to quit controller")
			s.Shutdown()
		case <-s.quit:
			Debug("received quit signal, breaking out of shutdown signal watcher")
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
	Debug("calling stop controller")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	Debug("sending shutdown signal to controller")
	s.quit.Q()
}

func (s *State) startWallet() (err error) {
	Debug("getting configured TLS certificates")
	certs := walletmain.ReadCAFile(s.cx.Config)
	Debug("establishing wallet connection")
	if s.walletClient, err = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         *s.cx.Config.WalletServer,
			Endpoint:     "ws",
			User:         *s.cx.Config.Username,
			Pass:         *s.cx.Config.Password,
			TLS:          *s.cx.Config.TLS,
			Certificates: certs,
		}, s.chainNotifier(), s.quit,
	); Check(err) {
	}
	return
}

// Run must be start as a goroutine, central routing for the business of the
// controller
//
// For increased simplicity, every type of work runs in one thread, only signalling
// from background goroutines to trigger state changes.
func (s *State) Run() {
	Debug("starting controller server")
	var err error
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
	ticker := time.NewTicker(time.Second)
	if err = s.startWallet(); !Check(err) {
		s.running = true
		s.start.Signal()
		Debug("getting templates...")
	}
out:
	for {
		Debug("controller now pausing")
	pausing:
		for {
			select {
			case bu := <-s.blockUpdate:
				Debug("received new block update while paused")
				doBlockUpdate(bu)
			case <-ticker.C:
				Debug("controller ticker running")
				if s.running {
					s.start.Signal()
				}
				Debug("do things that run anyway like the p2padvt")
			case <-s.start.Wait():
				Debug("received start signal while paused")
				if s.walletClient.Disconnected() {
					Debug("wallet client is disconnected, retrying")
					if err = s.startWallet(); !Check(err) {
						Debug("wallet client is connected, switching to running")
						s.running = true
						break pausing
					}
				}
			case <-s.stop.Wait():
				Debug("received stop signal while paused")
				s.running = false
			case <-s.quit.Wait():
				Debug("received quit signal while paused")
				break out
			}
		}
		Debug("controller now running")
	running:
		for {
			select {
			case bu := <-s.blockUpdate:
				Debug("received new block update while running")
				doBlockUpdate(bu)
			case <-ticker.C:
				Debug("do tickery things like p2padvt")
				Debug("resending current templates...")
				if s.running {
					Debug("controller ticker running")
					Debug("checking if wallet is connected")
					if s.walletClient.Disconnected() {
						Debug("wallet client has disconnected, switching to pausing")
						s.stop.Signal()
						break
					}
				} else {
					break running
				}
			case <-s.start.Wait():
				Debug("received start signal while running")
			case <-s.stop.Wait():
				Debug("received stop signal while running")
				s.running = false
				break running
			case <-s.quit.Wait():
				Debug("received quit signal while running")
				break out
			}
		}
		Debug("disconnecting wallet client if it was connected")
		if !s.walletClient.Disconnected() {
			s.walletClient.Disconnect()
		}
	}
}

func doBlockUpdate(bu *BlockUpdate) {
	Debug("getting templates...")
	Debugs(bu)
	Debugs("caching templates")
	Debug("sending out templates...")
}

func getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
	Debug("getting a block template generator")
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

func (s *State) chainNotifier() *rpcclient.NotificationHandlers {
	return &rpcclient.NotificationHandlers{
		OnBlockConnected: func(hash *chainhash.Hash, height int32, t time.Time) {
			Debug("updating for new block")
			s.blockUpdate <- &BlockUpdate{
				hash:   *hash,
				height: height,
				t:      t,
			}
		},
	}
}
