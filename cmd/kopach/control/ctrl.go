package control

import (
	"errors"
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/cmd/kopach/control/peersummary"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/blockchain"
	"github.com/p9c/pod/pkg/blockchain/fork"
	"github.com/p9c/pod/pkg/blockchain/mining"
	"github.com/p9c/pod/pkg/blockchain/wire"
	"github.com/p9c/pod/pkg/comm/transport"
	rav "github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/rpc/rpcclient"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/qu"
	"github.com/p9c/pod/pkg/util/routeable"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"sync"
	"time"
)

const (
	MaxDatagramSize      = 8192 * 3
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

// State stores the state of the controller
type State struct {
	sync.Mutex
	Syncing           *atomic.Bool
	cfg               *pod.Config
	node              *chainrpc.Node
	connMgr           chainrpc.ServerConnManager
	stateCfg          *state.Config
	mempoolUpdateChan qu.C
	uuid              uint64
	start, stop, quit qu.C
	blockUpdate       chan *util.Block
	generator         *mining.BlkTmplGenerator
	nextAddress       util.Address
	walletClient      *rpcclient.Client
	msgBlockTemplates *templates.RecentMessages
	templateShards    [][]byte
	multiConn         *transport.Channel
	otherNodes        map[uint64]*nodeSpec
	hashSampleBuf     *rav.BufferUint64
	hashCount         atomic.Uint64
	lastNonce         int32
	lastBlockUpdate   atomic.Int64
}

type nodeSpec struct {
	time.Time
	addr string
}

// New creates a new controller
func New(
	syncing *atomic.Bool,
	cfg *pod.Config,
	stateCfg *state.Config,
	node *chainrpc.Node,
	connMgr chainrpc.ServerConnManager,
	mempoolUpdateChan qu.C,
	uuid uint64,
	killall qu.C,
) (s *State) {
	var e error
	quit := qu.T()
	D.Ln("creating othernodes map")
	s = &State{
		Syncing:           syncing,
		cfg:               cfg,
		node:              node,
		connMgr:           connMgr,
		stateCfg:          stateCfg,
		mempoolUpdateChan: mempoolUpdateChan,
		otherNodes:        make(map[uint64]*nodeSpec),
		quit:              quit,
		uuid:              uuid,
		start:             qu.Ts(2),
		stop:              qu.Ts(2),
		blockUpdate:       make(chan *util.Block, 1),
		hashSampleBuf:     rav.NewBufferUint64(100),
		msgBlockTemplates: templates.NewRecentMessages(),
	}
	s.lastBlockUpdate.Store(time.Now().Add(-time.Second * 3).Unix())
	s.generator = chainrpc.GetBlkTemplateGenerator(node, cfg, stateCfg)
	var mc *transport.Channel
	if mc, e = transport.NewBroadcastChannel(
		"controller",
		s,
		*cfg.MinerPass,
		transport.DefaultPort,
		MaxDatagramSize,
		handlersMulticast,
		quit,
	); E.Chk(e) {
		return
	}
	s.multiConn = mc
	go func() {
		D.Ln("starting shutdown signal watcher")
		select {
		case <-killall:
			D.Ln("received killall signal, signalling to quit controller")
			s.Shutdown()
		case <-s.quit:
			D.Ln("received quit signal, breaking out of shutdown signal watcher")
		}
	}()
	node.Chain.Subscribe(
		func(n *blockchain.Notification) {
			switch n.Type {
			case blockchain.NTBlockConnected:
				if s.Syncing.Load() {
					return
				}
				D.Ln("received block connected notification")
				if b, ok := n.Data.(*util.Block); !ok {
					W.Ln("block notification is not a block")
					break
				} else {
					s.blockUpdate <- b
				}
			}
		},
	)
	return
}

// todo: the stop

// Start up the controller
func (s *State) Start() {
	D.Ln("calling start controller")
	s.start.Signal()
}

// Stop the controller
func (s *State) Stop() {
	D.Ln("calling stop controller")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	D.Ln("sending shutdown signal to controller")
	s.quit.Q()
}

func (s *State) startWallet() (e error) {
	D.Ln("getting configured TLS certificates")
	certs := pod.ReadCAFile(s.cfg)
	D.Ln("establishing wallet connection")
	if s.walletClient, e = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         *s.cfg.WalletServer,
			Endpoint:     "ws",
			User:         *s.cfg.Username,
			Pass:         *s.cfg.Password,
			TLS:          *s.cfg.TLS,
			Certificates: certs,
		}, nil, s.quit,
	); E.Chk(e) {
	}
	return
}

func (s *State) updateBlockTemplate() (e error) {
	D.Ln("getting current chain tip")
	// s.node.Chain.ChainLock.Lock() // previously this was done before the above, it might be jumping the gun on a new block
	h := s.node.Chain.BestSnapshot().Hash
	var blk *util.Block
	if blk, e = s.node.Chain.BlockByHash(&h); E.Chk(e) {
		return
	}
	// s.node.Chain.ChainLock.Unlock()
	D.Ln("updating block from chain tip")
	// D.S(blk)
	if e = s.doBlockUpdate(blk); E.Chk(e) {
	}
	return
}

// Run must be started as a goroutine, central routing for the business of the
// controller
//
// For increased simplicity, every type of work runs in one thread, only signalling
// from background goroutines to trigger state changes.
func (s *State) Run() {
	D.Ln("starting controller server")
	var e error
	ticker := time.NewTicker(time.Second)
out:
	for {
		// if !s.Syncing.Load() {
		// 	if s.walletClient.Disconnected() {
		// 		D.Ln("wallet client is disconnected, retrying")
		// 		if e = s.startWallet(); !E.Chk(e) {
		// 			continue
		// 		}
		// 		select {
		// 		case <-time.After(time.Second):
		// 			continue
		// 		case <-s.quit:
		// 			break out
		// 		}
		// 	}
		// } else {
		// 	select {
		// 	case <-time.After(time.Second):
		// 		continue
		// 	case <-s.quit:
		// 		break out
		// 	}
		// }
		// // D.Ln("wallet client is connected, switching to running")
		// // if e = s.updateBlockTemplate(); E.Chk(e) {
		// // }
		D.Ln("controller now pausing")
		*s.cfg.Controller = false
	pausing:
		for {
			select {
			case <-s.mempoolUpdateChan:
				// D.Ln("mempool update chan signal")
				// if e = s.updateBlockTemplate(); E.Chk(e) {
				// }
			case /* bu :=*/ <-s.blockUpdate:
				// D.Ln("received new block update while paused")
				// if e = s.doBlockUpdate(bu); E.Chk(e) {
				// }
				// // s.updateBlockTemplate()
			case <-ticker.C:
				D.Ln("controller ticker running")
				// s.Advertise()
				// s.checkConnectivity()
			case <-s.start.Wait():
				D.Ln("received start signal while paused")
				if s.walletClient.Disconnected() {
					D.Ln("wallet client is disconnected, retrying")
					if e = s.startWallet(); E.Chk(e) {
						// s.updateBlockTemplate()
						break
					}
				}
				D.Ln("wallet client is connected, switching to running")
				break pausing
			case <-s.stop.Wait():
				D.Ln("received stop signal while paused")
			case <-s.quit.Wait():
				D.Ln("received quit signal while paused")
				break out
			}
		}
		// if s.templateShards == nil || len(s.templateShards) < 1 {
		// }
		D.Ln("controller now running")
		if e = s.updateBlockTemplate(); E.Chk(e) {
		}
		*s.cfg.Controller = true
	running:
		for {
			select {
			case <-s.mempoolUpdateChan:
				D.Ln("mempoolUpdateChan updating block templates")
				if e = s.updateBlockTemplate(); E.Chk(e) {
					break
				}
				D.Ln("sending out templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
				}
			case bu := <-s.blockUpdate:
				// _ = bu
				D.Ln("received new block update while running")
				if e = s.doBlockUpdate(bu); E.Chk(e) {
					break
				}
				D.Ln("sending out templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
					break
				}
			case <-ticker.C:
				D.Ln("checking if wallet is connected")
				s.checkConnectivity()
				D.Ln("resending current templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
					break
				}
				if s.walletClient.Disconnected() {
					D.Ln("wallet client has disconnected, switching to pausing")
					break running
				}
			case <-s.start.Wait():
				D.Ln("received start signal while running")
			case <-s.stop.Wait():
				D.Ln("received stop signal while running")
				break running
			case <-s.quit.Wait():
				D.Ln("received quit signal while running")
				break out
			}
		}
	}
}

func (s *State) checkConnectivity() {
	// if !*s.cfg.Generate || *s.cfg.GenThreads == 0 {
	// 	D.Ln("no need to check connectivity if we aren't mining")
	// 	return
	// }
	if *s.cfg.Solo {
		D.Ln("in solo mode, mining anyway")
		s.Start()
		return
	}
	T.Ln("checking connectivity state")
	ps := make(chan peersummary.PeerSummaries, 1)
	s.node.PeerState <- ps
	T.Ln("sent peer list query")
	var lanPeers int
	var totalPeers int
	select {
	case connState := <-ps:
		T.Ln("received peer list query response")
		totalPeers = len(connState)
		for i := range connState {
			if routeable.IPNet.Contains(connState[i].IP) {
				lanPeers++
			}
		}
		if *s.cfg.LAN {
			// if there is no peers on lan and solo was not set, stop mining
			if lanPeers == 0 {
				T.Ln("no lan peers while in lan mode, stopping mining")
				s.Stop()
			} else {
				s.Start()
			}
		} else {
			if totalPeers-lanPeers == 0 {
				// we have no peers on the internet, stop mining
				T.Ln("no internet peers, stopping mining")
				s.Stop()
			} else {
				s.Start()
			}
		}
		break
		// quit waiting if we are shutting down
	case <-s.quit:
		break
	}
	T.Ln(totalPeers, "total peers", lanPeers, "lan peers solo:", *s.cfg.Solo, "lan:", *s.cfg.LAN)
}

//
// func (s *State) Advertise() {
// 	if !*s.cfg.Discovery {
// 		return
// 	}
// 	T.Ln("sending out advertisment")
// 	var e error
// 	if e = s.multiConn.SendMany(
// 		p2padvt.Magic,
// 		transport.GetShards(p2padvt.Get(s.uuid, s.cfg)),
// 	); E.Chk(e) {
// 	}
// }

func (s *State) doBlockUpdate(prev *util.Block) (e error) {
	if s.Syncing.Load() {
		return nil
	}
	D.Ln("do block update")
	if s.nextAddress == nil {
		D.Ln("getting new address for templates")
		// if s.nextAddress, e = s.GetNewAddressFromMiningAddrs(); T.Chk(e) {
		if s.nextAddress, e = s.GetNewAddressFromWallet(); T.Chk(e) {
			s.Stop()
			return
		}
		// }
	}
	D.Ln("getting templates...", prev.MsgBlock().Header.Timestamp)
	var tpl *templates.Message
	if tpl, e = s.GetMsgBlockTemplate(prev, s.nextAddress); E.Chk(e) {
		s.Stop()
		return
	}
	s.msgBlockTemplates.Add(tpl)
	D.Ln(tpl.Timestamp)
	D.Ln("caching error corrected message shards...")
	s.templateShards = transport.GetShards(tpl.Serialize())
	return
}

// GetMsgBlockTemplate gets a Message building on given block paying to a given
// address
func (s *State) GetMsgBlockTemplate(prev *util.Block, addr util.Address) (mbt *templates.Message, e error) {
	T.Ln("GetMsgBlockTemplate")
	rand.Seed(time.Now().Unix())
	mbt = &templates.Message{
		Nonce:     rand.Uint64(),
		UUID:      s.uuid,
		PrevBlock: prev.MsgBlock().BlockHash(),
		Height:    prev.Height() + 1,
		Bits:      make(templates.Diffs),
		Merkles:   make(templates.Merkles),
	}
	for next, curr, more := fork.AlgoVerIterator(mbt.Height); more(); next() {
		D.Ln("creating template for", curr())
		var templateX *mining.BlockTemplate
		if templateX, e = s.generator.NewBlockTemplate(
			addr,
			fork.GetAlgoName(curr(), mbt.Height),
		); D.Chk(e) || templateX == nil {
		} else {
			// I.S(templateX)
			newB := templateX.Block
			newH := newB.Header
			mbt.Timestamp = newH.Timestamp
			mbt.Bits[curr()] = newH.Bits
			mbt.Merkles[curr()] = newH.MerkleRoot
			D.Ln("merkle for", curr(), mbt.Merkles[curr()])
			mbt.SetTxs(curr(), newB.Transactions)
		}
	}
	return
}

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (s *State) GetNewAddressFromWallet() (addr util.Address, e error) {
	if s.walletClient != nil {
		if !s.walletClient.Disconnected() {
			D.Ln("have access to a wallet, generating address")
			if addr, e = s.walletClient.GetNewAddress("default"); E.Chk(e) {
			} else {
				D.Ln("-------- found address", addr)
			}
		}
	} else {
		e = errors.New("no wallet available for new address")
		D.Ln(e)
	}
	return
}

// GetNewAddressFromMiningAddrs tries to get an address from the mining
// addresses list in the configuration file
func (s *State) GetNewAddressFromMiningAddrs() (addr util.Address, e error) {
	if s.cfg.MiningAddrs == nil {
		e = errors.New("mining addresses is nil")
		D.Ln(e)
		return
	}
	if len(*s.cfg.MiningAddrs) < 1 {
		e = errors.New("no mining addresses")
		D.Ln(e)
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*s.cfg.MiningAddrs))
	addr = s.stateCfg.ActiveMiningAddrs[p2a]
	// remove the address from the state
	if p2a == 0 {
		s.stateCfg.ActiveMiningAddrs = s.stateCfg.ActiveMiningAddrs[1:]
	} else {
		s.stateCfg.ActiveMiningAddrs = append(
			s.stateCfg.ActiveMiningAddrs[:p2a],
			s.stateCfg.ActiveMiningAddrs[p2a+1:]...,
		)
	}
	// update the config
	var ma cli.StringSlice
	for i := range s.stateCfg.ActiveMiningAddrs {
		ma = append(ma, s.stateCfg.ActiveMiningAddrs[i].String())
	}
	*s.cfg.MiningAddrs = ma
	save.Pod(s.cfg)
	return
}

var handlersMulticast = transport.Handlers{
	string(sol.Magic): processSolMsg,
	// string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (e error) {
	D.Ln("processing advertisment message", src, dst)
	s := ctx.(*State)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	var uuid uint64
	uuid = j.UUID
	// I.Ln("uuid of advertisment", uuid, s.otherNodes)
	if uuid == s.uuid {
		D.Ln("ignoring own advertisment message")
		return
	}
	if _, ok := s.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		I.Ln("connecting to lan peer with same PSK", j.IPs, uuid)
		s.otherNodes[uuid] = &nodeSpec{}
		s.otherNodes[uuid].Time = time.Now()
		// try all IPs
		if *s.cfg.AutoListen {
			s.cfg.P2PConnect = &cli.StringSlice{}
		}
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			if e = s.connMgr.Connect(
				peerIP,
				true,
			); E.Chk(e) {
				continue
			}
			D.Ln("connected to peer via address", peerIP)
			s.otherNodes[uuid].addr = peerIP
			break
		}
		I.Ln("otherNodes", s.otherNodes)
	} else {
		// update last seen time for uuid for garbage collection of stale disconnected
		// nodes
		s.otherNodes[uuid].Time = time.Now()
	}
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range s.otherNodes {
		if time.Now().Sub(s.otherNodes[i].Time) > time.Second*6 {
			// also remove from connection manager
			if e = s.connMgr.RemoveByAddr(s.otherNodes[i].addr); E.Chk(e) {
			}
			D.Ln("deleting", s.otherNodes[i])
			delete(s.otherNodes, i)
		}
	}
	// on := int32(len(s.otherNodes))
	// s.otherNodeCount.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (e error) {
	I.Ln("received solution", src, dst)
	s := ctx.(*State)
	var so sol.Solution
	gotiny.Unmarshal(b, &so)
	tpl := s.msgBlockTemplates.Find(so.Nonce)
	if tpl == nil {
		I.Ln("solution nonce", so.Nonce, "is not known by this controller")
		return
	}
	if so.UUID != s.uuid {
		I.Ln("solution is for another controller")
		return
	}
	var newHeader *wire.BlockHeader
	if newHeader, e = so.Decode(); E.Chk(e) {
		return
	}
	if newHeader.PrevBlock != tpl.PrevBlock {
		I.Ln("block submitted by kopach miner worker is stale")
		return
	}
	var msgBlock *wire.MsgBlock
	if msgBlock, e = tpl.Reconstruct(newHeader); E.Chk(e) {
		I.Ln("failed to construct new header")
		return
	}
	
	I.Ln("sending pause to workers")
	if e = s.multiConn.SendMany(pause.Magic, transport.GetShards(p2padvt.Get(s.uuid, s.cfg))); E.Chk(e) {
		return
	}
	I.Ln("signalling controller to enter pause mode")
	s.Stop()
	defer s.Start()
	block := util.NewBlock(msgBlock)
	block.SetHeight(tpl.Height)
	var isOrphan bool
	I.Ln("submitting block for processing")
	if isOrphan, e = s.node.SyncManager.ProcessBlock(block, blockchain.BFNone); E.Chk(e) {
		// Anything other than a rule violation is an unexpected error, so log that
		// error as an internal error.
		if _, ok := e.(blockchain.RuleError); !ok {
			W.F(
				"Unexpected error while processing block submitted via kopach miner:", e,
			)
			return
		} else {
			W.Ln("block submitted via kopach miner rejected:", e)
			if isOrphan {
				W.Ln("block is an orphan")
				return
			}
			return
		}
	}
	I.Ln("the block was accepted, new height", block.Height())
	I.C(
		func() string {
			bmb := block.MsgBlock()
			coinbaseTx := bmb.Transactions[0].TxOut[0]
			prevHeight := block.Height() - 1
			prevBlock, _ := s.node.Chain.BlockByHeight(prevHeight)
			prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
			since := bmb.Header.Timestamp.Unix() - prevTime
			bHash := bmb.BlockHashWithAlgos(block.Height())
			return fmt.Sprintf(
				"new block height %d %08x %s%10d %08x %v %s %ds since prev",
				block.Height(),
				prevBlock.MsgBlock().Header.Bits,
				bHash,
				bmb.Header.Timestamp.Unix(),
				bmb.Header.Bits,
				util.Amount(coinbaseTx.Value),
				fork.GetAlgoName(
					bmb.Header.Version,
					block.Height(),
				), since,
			)
		},
	)
	I.Ln("clearing address used for block")
	s.nextAddress = nil
	return
}

// hashrate reports from workers
func processHashrateMsg(ctx interface{}, src net.Addr, dst string, b []byte) (e error) {
	s := ctx.(*State)
	var hr hashrate.Hashrate
	gotiny.Unmarshal(b, &hr)
	// only count each one once
	if s.lastNonce == hr.Nonce {
		return
	}
	s.lastNonce = hr.Nonce
	// add to total hash counts
	s.hashCount.Add(uint64(hr.Count))
	return
}

func (s *State) hashReport() float64 {
	s.hashSampleBuf.Add(s.hashCount.Load())
	av := ewma.NewMovingAverage()
	var i int
	var prev uint64
	if e := s.hashSampleBuf.ForEach(
		func(v uint64) (e error) {
			if i < 1 {
				prev = v
			} else {
				interval := v - prev
				av.Add(float64(interval))
				prev = v
			}
			i++
			return nil
		},
	); E.Chk(e) {
	}
	return av.Value()
}
