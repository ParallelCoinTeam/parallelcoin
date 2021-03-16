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
	"time"
)

const (
	MaxDatagramSize      = 8192
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

// State stores the state of the controller
type State struct {
	cfg               *pod.Config
	node              *chainrpc.Node
	rpcServer         *chainrpc.Server
	stateCfg          *state.Config
	mempoolUpdateChan qu.C
	uuid              uint64
	start, stop, quit qu.C
	blockUpdate       chan *util.Block
	generator         *mining.BlkTmplGenerator
	nextAddress       util.Address
	walletClient      *rpcclient.Client
	msgBlockTemplate  *templates.Message
	templateShards    [][]byte
	multiConn         *transport.Channel
	otherNodes        map[uint64]*nodeSpec
	otherNodeCount    *atomic.Int32
	hashSampleBuf     *rav.BufferUint64
	hashCount         atomic.Uint64
	lastNonce         int32
	mining            *atomic.Bool
}

type nodeSpec struct {
	time.Time
	addr string
}

// New creates a new controller
func New(
	cfg *pod.Config,
	stateCfg *state.Config,
	node *chainrpc.Node,
	rpcServer *chainrpc.Server,
	otherNodeCount *atomic.Int32,
	mempoolUpdateChan qu.C,
	uuid uint64,
	killall qu.C,
) (s *State) {
	var e error
	if *cfg.DisableController {
		wrn.Ln("controller is disabled")
		return
	}
	quit := qu.T()
	s = &State{
		cfg:               cfg,
		node:              node,
		rpcServer:         rpcServer,
		stateCfg:          stateCfg,
		mempoolUpdateChan: mempoolUpdateChan,
		otherNodes:        make(map[uint64]*nodeSpec),
		otherNodeCount:    otherNodeCount,
		quit:              quit,
		uuid:              uuid,
		start:             qu.Ts(1),
		stop:              qu.Ts(1),
		blockUpdate:       make(chan *util.Block, 1),
		hashSampleBuf:     rav.NewBufferUint64(100),
		mining:            atomic.NewBool(false),
	}
	s.generator = s.getBlkTemplateGenerator()
	var mc *transport.Channel
	if mc, e = transport.NewBroadcastChannel(
		"controller",
		s,
		*cfg.MinerPass,
		transport.DefaultPort,
		MaxDatagramSize,
		handlersMulticast,
		quit,
	); err.Chk(e) {
		return
	}
	s.multiConn = mc
	go func() {
		dbg.Ln("starting shutdown signal watcher")
		select {
		case <-killall:
			dbg.Ln("received killall signal, signalling to quit controller")
			s.Shutdown()
		case <-s.quit:
			dbg.Ln("received quit signal, breaking out of shutdown signal watcher")
		}
	}()
	node.Chain.Subscribe(
		func(n *blockchain.Notification) {
			switch n.Type {
			case blockchain.NTBlockConnected:
				dbg.Ln("received block connected notification")
				if b, ok := n.Data.(*util.Block); !ok {
					wrn.Ln("block notification is not a block")
					break
				} else {
					s.blockUpdate <- b
				}
			}
		},
	)
	return
}

// Start up the controller
func (s *State) Start() {
	dbg.Ln("calling start controller")
	s.start.Signal()
}

// Stop the controller
func (s *State) Stop() {
	dbg.Ln("calling stop controller")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	dbg.Ln("sending shutdown signal to controller")
	s.quit.Q()
}

func (s *State) startWallet() (e error) {
	dbg.Ln("getting configured TLS certificates")
	certs := pod.ReadCAFile(s.cfg)
	dbg.Ln("establishing wallet connection")
	if s.walletClient, e = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         *s.cfg.WalletServer,
			Endpoint:     "ws",
			User:         *s.cfg.Username,
			Pass:         *s.cfg.Password,
			TLS:          *s.cfg.TLS,
			Certificates: certs,
		}, nil, s.quit,
	); err.Chk(e) {
	}
	return
}

func (s *State) updateBlockTemplate() {
	dbg.Ln("getting current chain tip")
	var e error
	s.node.Chain.ChainLock.Lock() // previously this was done before the above, it might be jumping the gun on a new block
	h := s.node.Chain.BestChain.Tip().Header().BlockHash()
	var blk *util.Block
	if blk, e = s.node.Chain.BlockByHash(&h); err.Chk(e) {
		s.node.Chain.ChainLock.Unlock()
		return
	}
	s.node.Chain.ChainLock.Unlock()
	dbg.Ln("updating block from chain tip")
	dbg.S(blk)
	if e = s.doBlockUpdate(blk); err.Chk(e) {
	}
}

// Run must be started as a goroutine, central routing for the business of the
// controller
//
// For increased simplicity, every type of work runs in one thread, only signalling
// from background goroutines to trigger state changes.
func (s *State) Run() {
	dbg.Ln("starting controller server")
	var e error
	if *s.cfg.DisableController {
		wrn.Ln("controller is disabled")
		return
	}
	ticker := time.NewTicker(time.Second)
out:
	for {
		dbg.Ln("controller now pausing")
		s.mining.Store(false)
		if s.walletClient.Disconnected() {
			dbg.Ln("wallet client is disconnected, retrying")
			if e = s.startWallet(); !err.Chk(e) {
				dbg.Ln("wallet client is connected, switching to running")
				
			}
		}
		s.updateBlockTemplate()
	pausing:
		for {
			select {
			case <-s.mempoolUpdateChan:
				s.updateBlockTemplate()
			case <-s.blockUpdate:
				dbg.Ln("received new block update while paused")
				// if e = s.doBlockUpdate(bu); err.Chk(e) {
				// }
				s.updateBlockTemplate()
			case <-ticker.C:
				dbg.Ln("controller ticker running")
				s.Advertise()
				s.checkConnectivity()
			case <-s.start.Wait():
				dbg.Ln("received start signal while paused")
				if s.walletClient.Disconnected() {
					dbg.Ln("wallet client is disconnected, retrying")
					if e = s.startWallet(); !err.Chk(e) {
						dbg.Ln("wallet client is connected, switching to running")
						s.updateBlockTemplate()
					}
				}
				break pausing
			case <-s.stop.Wait():
				dbg.Ln("received stop signal while paused")
			case <-s.quit.Wait():
				dbg.Ln("received quit signal while paused")
				break out
			}
		}
		dbg.Ln("controller now running")
		// if s.templateShards == nil || len(s.templateShards) < 1 {
		s.updateBlockTemplate()
		// }
		s.mining.Store(true)
	running:
		for {
			select {
			case <-s.mempoolUpdateChan:
				s.updateBlockTemplate()
				dbg.Ln("sending out templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); err.Chk(e) {
				}
			case <-s.blockUpdate:
				dbg.Ln("received new block update while running")
				s.updateBlockTemplate()
				dbg.Ln("sending out templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); err.Chk(e) {
				}
			case <-ticker.C:
				// dbg.Ln("controller ticker running")
				// dbg.Ln("checking if wallet is connected")
				s.Advertise()
				s.checkConnectivity()
				// dbg.Ln("resending current templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); err.Chk(e) {
				}
				if s.walletClient.Disconnected() {
					dbg.Ln("wallet client has disconnected, switching to pausing")
					break running
				}
			case <-s.start.Wait():
				dbg.Ln("received start signal while running")
			case <-s.stop.Wait():
				dbg.Ln("received stop signal while running")
				break running
			case <-s.quit.Wait():
				dbg.Ln("received quit signal while running")
				break out
			}
		}
	}
}

func (s *State) checkConnectivity() {
	// if !*s.cfg.Generate || *s.cfg.GenThreads == 0 {
	// 	dbg.Ln("no need to check connectivity if we aren't mining")
	// 	return
	// }
	if *s.cfg.Solo {
		dbg.Ln("in solo mode, mining anyway")
		s.Start()
		return
	}
	trc.Ln("checking connectivity state")
	ps := make(chan chainrpc.PeerSummaries, 1)
	s.node.PeerState <- ps
	trc.Ln("sent peer list query")
	var lanPeers int
	var totalPeers int
	select {
	case connState := <-ps:
		trc.Ln("received peer list query response")
		totalPeers = len(connState)
		for i := range connState {
			if routeable.IPNet.Contains(connState[i].IP) {
				lanPeers++
			}
		}
		if *s.cfg.LAN {
			// if there is no peers on lan and solo was not set, stop mining
			if lanPeers == 0 {
				trc.Ln("no lan peers while in lan mode, stopping mining")
				s.Stop()
			} else {
				s.Start()
			}
		} else {
			if totalPeers-lanPeers == 0 {
				// we have no peers on the internet, stop mining
				trc.Ln("no internet peers, stopping mining")
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
	trc.Ln(totalPeers, "total peers", lanPeers, "lan peers solo:", *s.cfg.Solo, "lan:", *s.cfg.LAN)
}

func (s *State) Advertise() {
	trc.Ln("sending out advertisment")
	var e error
	if e = s.multiConn.SendMany(
		p2padvt.Magic,
		transport.GetShards(p2padvt.Get(s.uuid, s.cfg, s.node)),
	); err.Chk(e) {
	}
}

func (s *State) doBlockUpdate(prev *util.Block) (e error) {
	if s.nextAddress == nil {
		dbg.Ln("getting new address for templates")
		if s.nextAddress, e = s.GetNewAddressFromMiningAddrs(); err.Chk(e) {
			if s.nextAddress, e = s.GetNewAddressFromWallet(); err.Chk(e) {
				return
			}
		}
	}
	dbg.Ln("getting templates...", prev.MsgBlock().Header.Timestamp)
	if s.msgBlockTemplate, e = s.GetMsgBlockTemplate(prev, s.nextAddress); err.Chk(e) {
		// return
	}
	dbg.Ln(s.msgBlockTemplate.Timestamp)
	dbg.Ln("caching error corrected message shards...")
	s.templateShards = transport.GetShards(s.msgBlockTemplate.Serialize())
	return
}

func (s *State) getBlkTemplateGenerator() *mining.BlkTmplGenerator {
	dbg.Ln("getting a block template generator")
	return mining.NewBlkTmplGenerator(
		&mining.Policy{
			BlockMinWeight:    uint32(*s.cfg.BlockMinWeight),
			BlockMaxWeight:    uint32(*s.cfg.BlockMaxWeight),
			BlockMinSize:      uint32(*s.cfg.BlockMinSize),
			BlockMaxSize:      uint32(*s.cfg.BlockMaxSize),
			BlockPrioritySize: uint32(*s.cfg.BlockPrioritySize),
			TxMinFreeFee:      s.stateCfg.ActiveMinRelayTxFee,
		},
		s.node.ChainParams,
		s.node.TxMemPool,
		s.node.Chain,
		s.node.TimeSource,
		s.node.SigCache,
		s.node.HashCache,
	)
}

// GetMsgBlockTemplate gets a Message building on given block paying to a given
// address
func (s *State) GetMsgBlockTemplate(prev *util.Block, addr util.Address) (mbt *templates.Message, e error) {
	trc.Ln("GetMsgBlockTemplate")
	rand.Seed(time.Now().Unix())
	mbt = &templates.Message{
		Nonce:     rand.Uint64(),
		UUID:      s.uuid,
		PrevBlock: prev.MsgBlock().BlockHash(),
		Height:    prev.Height() + 1,
		Bits:      make(templates.Diffs),
		Merkles:   make(templates.Merkles),
	}
	//
	// mbt.Timestamp = prev.MsgBlock().Header.Timestamp.Truncate(time.Second).Add(time.Second)
	// dbg.Ln("initial timestamp", mbt.Timestamp)
	// tn := time.Now().Truncate(time.Second)
	// if fork.GetCurrent(mbt.Height) < 1 {
	// 	dbg.Ln("on legacy consensus")
	// 	mbt.Timestamp = s.generator.BestSnapshot().MedianTime.Add(time.Second)
	// } else {
	// 	if tn.After(mbt.Timestamp.Add(time.Second)) {
	// 		dbg.Ln("adjusted timestamp", tn)
	// 		mbt.Timestamp = tn
	// 	}
	// }
	for next, curr, more := fork.AlgoVerIterator(mbt.Height); more(); next() {
		var templateX *mining.BlockTemplate
		if templateX, e = s.generator.NewBlockTemplate(addr, fork.GetAlgoName(curr(), mbt.Height)); err.Chk(e) {
		} else {
			newB := templateX.Block
			newH := newB.Header
			mbt.Timestamp = newH.Timestamp
			mbt.Bits[curr()] = newH.Bits
			mbt.Merkles[curr()] = newH.MerkleRoot
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
			dbg.Ln("have access to a wallet, generating address")
			if addr, e = s.walletClient.GetNewAddress("default"); err.Chk(e) {
			} else {
				dbg.Ln("-------- found address", addr)
			}
		}
	} else {
		e = errors.New("no wallet available for new address")
		dbg.Ln(e)
	}
	return
}

// GetNewAddressFromMiningAddrs tries to get an address from the mining
// addresses list in the configuration file
func (s *State) GetNewAddressFromMiningAddrs() (addr util.Address, e error) {
	if s.cfg.MiningAddrs == nil {
		e = errors.New("mining addresses is nil")
		dbg.Ln(e)
		return
	}
	if len(*s.cfg.MiningAddrs) < 1 {
		e = errors.New("no mining addresses")
		dbg.Ln(e)
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
	string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (e error) {
	dbg.Ln("processing advertisment message", src, dst)
	s := ctx.(*State)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	uuid := j.UUID
	if uuid == s.uuid {
		// dbg.Ln("ignoring own advertisment message")
		return
	}
	if _, ok := s.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		inf.Ln("connecting to lan peer with same PSK", j.IPs, j.UUID)
		// try all IPs
		if *s.cfg.AutoListen {
			s.cfg.P2PConnect = &cli.StringSlice{}
		}
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			if e = s.rpcServer.Cfg.ConnMgr.Connect(
				peerIP,
				true,
			); err.Chk(e) {
				continue
			}
			dbg.Ln("connected to peer via address", peerIP)
			s.otherNodes[uuid] = &nodeSpec{}
			s.otherNodes[uuid].addr = addr
			break
		}
	}
	// update last seen time for uuid for garbage collection of stale disconnected
	// nodes
	s.otherNodes[uuid].Time = time.Now()
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range s.otherNodes {
		if time.Now().Sub(s.otherNodes[i].Time) > time.Second*9 {
			// also remove from connection manager
			if e = s.rpcServer.Cfg.ConnMgr.RemoveByAddr(s.otherNodes[i].addr); err.Chk(e) {
			}
			dbg.Ln("deleting", s.otherNodes[i])
			delete(s.otherNodes, i)
		}
	}
	on := int32(len(s.otherNodes))
	s.otherNodeCount.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (e error) {
	dbg.Ln("received solution", src, dst)
	s := ctx.(*State)
	var so sol.Solution
	gotiny.Unmarshal(b, &so)
	// dbg.S(so)
	if s.msgBlockTemplate == nil {
		dbg.Ln("template is nil, solution is not for this controller")
		return
	}
	if s.uuid != so.UUID {
		dbg.Ln("solution not from current controller", s.uuid, s.msgBlockTemplate.UUID, so.UUID)
		return
	}
	// todo: s.msgBlockTemplate needs to be changed into an array of 3 to keep the
	//  last 3, then search them here
	if so.Nonce != s.msgBlockTemplate.Nonce {
		dbg.Ln("sollution nonce is not known by this controller")
	}
	var newHeader *wire.BlockHeader
	if newHeader, e = so.Decode(); err.Chk(e) {
		return
	}
	if newHeader.PrevBlock != s.msgBlockTemplate.PrevBlock {
		dbg.Ln("block submitted by kopach miner worker is stale")
		return
	}
	var msgBlock *wire.MsgBlock
	if msgBlock, e = s.msgBlockTemplate.Reconstruct(newHeader); err.Chk(e) {
		return
	}
	dbg.Ln("sending pause to workers")
	if e = s.multiConn.SendMany(pause.Magic, transport.GetShards(p2padvt.Get(s.uuid, s.cfg, s.node))); err.Chk(e) {
		return
	}
	dbg.Ln("clearing current block template")
	s.msgBlockTemplate = nil
	dbg.Ln("signalling controller to enter pause mode")
	s.Stop()
	block := util.NewBlock(msgBlock)
	var isOrphan bool
	dbg.Ln("submitting block for processing")
	if isOrphan, e = s.node.SyncManager.ProcessBlock(block, blockchain.BFNone); err.Chk(e) {
		// Anything other than a rule violation is an unexpected error, so log that
		// error as an internal error.
		if _, ok := e.(blockchain.RuleError); !ok {
			wrn.F(
				"Unexpected error while processing block submitted via kopach miner:", err,
			)
			return
		} else {
			wrn.Ln("block submitted via kopach miner rejected:", err)
			if isOrphan {
				dbg.Ln("block is an orphan")
				return
			}
			return
		}
	}
	dbg.Ln("clearing address used for block")
	s.nextAddress = nil
	dbg.Ln("the block was accepted, new height", block.Height())
	trc.C(
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
	); err.Chk(e) {
	}
	return av.Value()
}
