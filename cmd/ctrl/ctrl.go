package ctrl

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/VividCortex/ewma"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/qu"
	"go.uber.org/atomic"

	"github.com/p9c/matrjoska/cmd/node/active"
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chainrpc"
	"github.com/p9c/matrjoska/pkg/chainrpc/hashrate"
	"github.com/p9c/matrjoska/pkg/chainrpc/job"
	"github.com/p9c/matrjoska/pkg/chainrpc/p2padvt"
	"github.com/p9c/matrjoska/pkg/chainrpc/pause"
	"github.com/p9c/matrjoska/pkg/chainrpc/sol"
	"github.com/p9c/matrjoska/pkg/chainrpc/templates"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pkg/mining"
	rav "github.com/p9c/matrjoska/pkg/ring"
	"github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/transport"
	"github.com/p9c/matrjoska/pkg/wire"
	"github.com/p9c/matrjoska/pod/config"
)

const (
	BufferSize = 4096
)

// State stores the state of the controller
type State struct {
	sync.Mutex
	cfg               *config.Config
	node              *chainrpc.Node
	connMgr           chainrpc.ServerConnManager
	stateCfg          *active.Config
	mempoolUpdateChan qu.C
	uuid              uint64
	start, stop, quit qu.C
	blockUpdate       chan *block.Block
	generator         *mining.BlkTmplGenerator
	nextAddress       btcaddr.Address
	walletClient      *rpcclient.Client
	msgBlockTemplates *templates.RecentMessages
	templateShards    [][]byte
	multiConn         *transport.Channel
	otherNodes        map[uint64]*nodeSpec
	hashSampleBuf     *rav.BufferUint64
	hashCount         atomic.Uint64
	lastNonce         int32
	lastBlockUpdate   atomic.Int64
	certs             []byte
}

type nodeSpec struct {
	time.Time
	addr string
}

// New creates a new controller
func New(
	syncing *atomic.Bool,
	cfg *config.Config,
	stateCfg *active.Config,
	node *chainrpc.Node,
	connMgr chainrpc.ServerConnManager,
	mempoolUpdateChan qu.C,
	uuid uint64,
	killall qu.C,
	start, stop qu.C,
) (s *State, e error) {
	quit := qu.T()
	I.Ln("creating othernodes map")
	I.Ln("getting configured TLS certificates")
	certs := cfg.ReadCAFile()
	s = &State{
		cfg:               cfg,
		node:              node,
		connMgr:           connMgr,
		stateCfg:          stateCfg,
		mempoolUpdateChan: mempoolUpdateChan,
		otherNodes:        make(map[uint64]*nodeSpec),
		quit:              quit,
		uuid:              uuid,
		start:             start,
		stop:              stop,
		blockUpdate:       make(chan *block.Block, 1),
		hashSampleBuf:     rav.NewBufferUint64(100),
		msgBlockTemplates: templates.NewRecentMessages(),
		certs:             certs,
	}
	s.lastBlockUpdate.Store(time.Now().Add(-time.Second * 3).Unix())
	s.generator = chainrpc.GetBlkTemplateGenerator(node, cfg, stateCfg)
	var mc *transport.Channel
	I.S(cfg.MulticastPass.V(), cfg.MulticastPass.Bytes())
	if mc, e = transport.NewBroadcastChannel(
		"controller",
		s,
		cfg.MulticastPass.Bytes(),
		transport.DefaultPort,
		constant.MaxDatagramSize,
		handlersMulticast,
		quit,
	); E.Chk(e) {
		return
	}
	s.multiConn = mc
	go func() {
		I.Ln("starting shutdown signal watcher")
		select {
		case <-killall:
			I.Ln("received killall signal, signalling to quit controller")
			s.Shutdown()
		case <-s.quit:
			I.Ln("received quit signal, breaking out of shutdown signal watcher")
		}
	}()
	node.Chain.Subscribe(
		func(n *blockchain.Notification) {
			switch n.Type {
			case blockchain.NTBlockConnected:
				T.Ln("received block connected notification")
				if b, ok := n.Data.(*block.Block); !ok {
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
	I.Ln("calling start controller")
	s.start.Signal()
}

// Stop the controller
func (s *State) Stop() {
	I.Ln("calling stop controller")
	s.stop.Signal()
}

// Shutdown the controller
func (s *State) Shutdown() {
	I.Ln("sending shutdown signal to controller")
	s.quit.Q()
}

func (s *State) startWallet() (e error) {
	if s.walletClient, e = rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         s.cfg.WalletServer.V(),
			Endpoint:     "ws",
			User:         s.cfg.Username.V(),
			Pass:         s.cfg.Password.V(),
			TLS:          s.cfg.ServerTLS.True(),
			Certificates: s.certs,
		}, nil, s.quit,
	); T.Chk(e) {
		return
	}
	I.Ln("established wallet connection")
	return
}

func (s *State) updateBlockTemplate() (e error) {
	I.Ln("getting current chain tip")
	// s.node.Chain.ChainLock.Lock() // previously this was done before the above, it might be jumping the gun on a new block
	h := s.node.Chain.BestSnapshot().Hash
	var blk *block.Block
	if blk, e = s.node.Chain.BlockByHash(&h); E.Chk(e) {
		return
	}
	// s.node.Chain.ChainLock.Unlock()
	I.Ln("updating block from chain tip")
	// D.S(blk)
	s.doBlockUpdate(blk)
	I.Ln("sending out templates...")
	if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
		return
	}
	return
}

// Run must be started as a goroutine, central routing for the business of the
// controller
//
// For increased simplicity, every type of work runs in one thread, only signalling
// from background goroutines to trigger state changes.
func (s *State) Run() {
	I.Ln("starting controller")
	var e error
	ticker := time.NewTicker(time.Second)
out:
	for {
		// if !s.Syncing.Load() {
		// 	if s.walletClient.Disconnected() {
		// 		I.Ln("wallet client is disconnected, retrying")
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
		// // I.Ln("wallet client is connected, switching to running")
		// // if e = s.updateBlockTemplate(); E.Chk(e) {
		// // }
		I.Ln("controller now pausing")
	pausing:
		for {
			select {
			case <-s.mempoolUpdateChan:
				// I.Ln("mempool update chan signal")
				// if e = s.updateBlockTemplate(); E.Chk(e) {
				// }
			case /* bu :=*/ <-s.blockUpdate:
				// I.Ln("received new block update while paused")
				// if e = s.doBlockUpdate(bu); E.Chk(e) {
				// }
				// // s.updateBlockTemplate()
			case <-ticker.C:
				// D.Ln("controller ticker running")
				// s.Advertise()
				// s.checkConnected()
				if s.walletClient.Disconnected() {
					I.Ln("wallet client is disconnected, retrying")
					if e = s.startWallet(); e != nil { // T.Chk(e) {
						// s.updateBlockTemplate()
						break
					}
				}
				I.Ln("wallet client is connected, switching to running")
				break pausing
			case <-s.start.Wait():
				I.Ln("received start signal while paused")
				if !s.checkConnected() {
					I.Ln("not connected")
					break
				}
				break pausing
			case <-s.stop.Wait():
				I.Ln("received stop signal while paused")
			case <-s.quit.Wait():
				I.Ln("received quit signal while paused")
				break out
			}
		}
		// if s.templateShards == nil || len(s.templateShards) < 1 {
		// }
		I.Ln("controller now running")
		if e = s.updateBlockTemplate(); E.Chk(e) {
		}
	running:
		for {
			select {
			case <-s.mempoolUpdateChan:
				I.Ln("mempoolUpdateChan updating block templates")
				if e = s.updateBlockTemplate(); E.Chk(e) {
					break
				}
				I.Ln("sending out templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
				}
			case bu := <-s.blockUpdate:
				// _ = bu
				go func(){

					I.Ln("received new block update while running")
					s.doBlockUpdate(bu)
					I.Ln("sending out templates...")
					if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
						return
					}
				}()
			case <-ticker.C:
				// T.Ln("checking if wallet is connected")
				if !s.checkConnected() {
					break running
				}
				// I.Ln("resending current templates...")
				if e = s.multiConn.SendMany(job.Magic, s.templateShards); E.Chk(e) {
					break
				}
				if s.walletClient.Disconnected() {
					I.Ln("wallet client has disconnected, switching to pausing")
					break running
				}
			case <-s.start.Wait():
				I.Ln("received start signal while running")
			case <-s.stop.Wait():
				I.Ln("received stop signal while running")
				break running
			case <-s.quit.Wait():
				I.Ln("received quit signal while running")
				break out
			}
		}
	}
}

func (s *State) checkConnected() (connected bool) {
	return true
	// if !*s.cfg.Generate || *s.cfg.GenThreads == 0 {
	// 	I.Ln("no need to check connectivity if we aren't mining")
	// 	return
	// }
	// if s.cfg.Solo.True() {
	// 	// I.Ln("in solo mode, mining anyway")
	// 	// s.Start()
	// 	return true
	// }
	// T.Ln("checking connectivity state")
	// ps := make(chan peersummary.PeerSummaries, 1)
	// s.node.PeerState <- ps
	// T.Ln("sent peer list query")
	// var lanPeers int
	// var totalPeers int
	// select {
	// case connState := <-ps:
	// 	T.Ln("received peer list query response")
	// 	totalPeers = len(connState)
	// 	for i := range connState {
	// 		if routeable.IPNet.Contains(connState[i].IP) {
	// 			lanPeers++
	// 		}
	// 	}
	// 	if s.cfg.LAN.True() {
	// 		// if there is no peers on lan and solo was not set, stop mining
	// 		if lanPeers == 0 {
	// 			T.Ln("no lan peers while in lan mode, stopping mining")
	// 			s.Stop()
	// 		} else {
	// 			s.Start()
	// 			connected = true
	// 		}
	// 	} else {
	// 		if totalPeers-lanPeers == 0 {
	// 			// we have no peers on the internet, stop mining
	// 			T.Ln("no internet peers, stopping mining")
	// 			s.Stop()
	// 		} else {
	// 			s.Start()
	// 			connected = true
	// 		}
	// 	}
	// 	break
	// 	// quit waiting if we are shutting down
	// case <-s.quit:
	// 	break
	// }
	// T.Ln(totalPeers, "total peers", lanPeers, "lan peers solo:",
	// 	s.cfg.Solo.True(), "lan:", s.cfg.LAN.True())
	// return
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

func (s *State) doBlockUpdate(prev *block.Block) (e error) {
	I.Ln("do block update")
	if s.nextAddress == nil {
		I.Ln("getting new address for templates")
		// if s.nextAddress, e = s.GetNewAddressFromMiningAddrs(); T.Chk(e) {
		if s.nextAddress, e = s.GetNewAddressFromWallet(); T.Chk(e) {
			s.Stop()
			return
		}
		// }
	}
	I.Ln("getting templates...", prev.WireBlock().Header.Timestamp)
	var tpl *templates.Message
	if tpl, e = s.GetMsgBlockTemplate(prev, s.nextAddress); E.Chk(e) {
		s.Stop()
		return
	}
	// I.S(tpl)
	s.msgBlockTemplates.Add(tpl)
	// I.Ln(tpl.Timestamp)
	I.Ln("caching error corrected message shards...")
	srl := tpl.Serialize()
	// I.S(srl)
	s.templateShards = transport.GetShards(srl)
	// var dt []byte
	// if _, e = fec.Decode(s.templateShards); E.Chk(e) {
	// }
	// I.S(dt)
	return
}

// GetMsgBlockTemplate gets a Message building on given block paying to a given
// address
func (s *State) GetMsgBlockTemplate(
	prev *block.Block, addr btcaddr.Address,
) (mbt *templates.Message, e error) {
	T.Ln("GetMsgBlockTemplate")
	rand.Seed(time.Now().Unix())
	mbt = &templates.Message{
		Nonce:     rand.Uint64(),
		UUID:      s.uuid,
		PrevBlock: prev.WireBlock().BlockHash(),
		Height:    prev.Height() + 1,
		Bits:      make(templates.Diffs),
		Merkles:   make(templates.Merkles),
	}
	for next, curr, more := fork.AlgoVerIterator(mbt.Height); more(); next() {
		// I.Ln("creating template for", curr())
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
			// I.Ln("merkle for", curr(), mbt.Merkles[curr()])
			mbt.SetTxs(curr(), newB.Transactions)
		}
	}
	return
}

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (s *State) GetNewAddressFromWallet() (addr btcaddr.Address, e error) {
	if s.walletClient != nil {
		if !s.walletClient.Disconnected() {
			I.Ln("have access to a wallet, generating address")
			if addr, e = s.walletClient.GetNewAddress("default"); E.Chk(e) {
			} else {
				I.Ln("-------- found address", addr)
			}
		}
	} else {
		e = errors.New("no wallet available for new address")
		I.Ln(e)
	}
	return
}

//
// // GetNewAddressFromMiningAddrs tries to get an address from the mining
// // addresses list in the configuration file
// func (s *State) GetNewAddressFromMiningAddrs() (addr btcaddr.Address, e error) {
// 	if s.cfg.MiningAddrs == nil {
// 		e = errors.New("mining addresses is nil")
// 		I.Ln(e)
// 		return
// 	}
// 	if len(s.cfg.MiningAddrs.S()) < 1 {
// 		e = errors.New("no mining addresses")
// 		I.Ln(e)
// 		return
// 	}
// 	// Choose a payment address at random.
// 	rand.Seed(time.Now().UnixNano())
// 	p2a := rand.Intn(len(s.cfg.MiningAddrs.S()))
// 	addr = s.stateCfg.ActiveMiningAddrs[p2a]
// 	// remove the address from the state
// 	if p2a == 0 {
// 		s.stateCfg.ActiveMiningAddrs = s.stateCfg.ActiveMiningAddrs[1:]
// 	} else {
// 		s.stateCfg.ActiveMiningAddrs = append(
// 			s.stateCfg.ActiveMiningAddrs[:p2a],
// 			s.stateCfg.ActiveMiningAddrs[p2a+1:]...,
// 		)
// 	}
// 	// update the config
// 	var ma cli.StringSlice
// 	for i := range s.stateCfg.ActiveMiningAddrs {
// 		ma = append(ma, s.stateCfg.ActiveMiningAddrs[i].String())
// 	}
// 	s.cfg.MiningAddrs.Set(ma)
// 	podcfg.Save(s.cfg)
// 	return
// }

var handlersMulticast = transport.Handlers{
	string(sol.Magic): processSolMsg,
	// string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(
	ctx interface{}, src net.Addr, dst string, b []byte,
) (e error) {
	I.Ln("processing advertisment message", src, dst)
	s := ctx.(*State)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	var uuid uint64
	uuid = j.UUID
	// I.Ln("uuid of advertisment", uuid, s.otherNodes)
	if uuid == s.uuid {
		I.Ln("ignoring own advertisment message")
		return
	}
	if _, ok := s.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		I.Ln("connecting to lan peer with same PSK", j.IPs, uuid)
		s.otherNodes[uuid] = &nodeSpec{}
		s.otherNodes[uuid].Time = time.Now()
		// try all IPs
		// if s.cfg.AutoListen.True() {
		// 	s.cfg.P2PConnect.Set(cli.StringSlice{})
		// }
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			if e = s.connMgr.Connect(
				peerIP,
				true,
			); E.Chk(e) {
				continue
			}
			I.Ln("connected to peer via address", peerIP)
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
			I.Ln("deleting", s.otherNodes[i])
			delete(s.otherNodes, i)
		}
	}
	// on := int32(len(s.otherNodes))
	// s.otherNodeCount.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(
	ctx interface{}, src net.Addr, dst string, b []byte,
) (e error) {
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
		I.Ln("blk submitted by kopach miner worker is stale")
		return
	}
	var msgBlock *wire.Block
	if msgBlock, e = tpl.Reconstruct(newHeader); E.Chk(e) {
		I.Ln("failed to construct new header")
		return
	}

	I.Ln("sending pause to workers")
	if e = s.multiConn.SendMany(pause.Magic, transport.GetShards(p2padvt.Get(s.uuid, (s.cfg.P2PListeners.S())[0])),
	); E.Chk(e) {
		return
	}
	I.Ln("signalling controller to enter pause mode")
	s.Stop()
	defer s.Start()
	blk := block.NewBlock(msgBlock)
	blk.SetHeight(tpl.Height)
	var isOrphan bool
	I.Ln("submitting blk for processing")
	if isOrphan, e = s.node.SyncManager.ProcessBlock(blk, blockchain.BFNone); E.Chk(e) {
		// Anything other than a rule violation is an unexpected error, so log that
		// error as an internal error.
		if _, ok := e.(blockchain.RuleError); !ok {
			W.F(
				"Unexpected error while processing blk submitted via kopach miner:", e,
			)
			return
		} else {
			W.Ln("blk submitted via kopach miner rejected:", e)
			if isOrphan {
				W.Ln("blk is an orphan")
				return
			}
			return
		}
	}
	I.Ln("the blk was accepted, new height", blk.Height())
	I.C(
		func() string {
			bmb := blk.WireBlock()
			coinbaseTx := bmb.Transactions[0].TxOut[0]
			prevHeight := blk.Height() - 1
			var prevBlock *block.Block
			if prevBlock, e = s.node.Chain.BlockByHeight(prevHeight); E.Chk(e) {
			}
			if prevBlock == nil {
				return "prevblock nil while generating log"
			}
			prevTime := prevBlock.WireBlock().Header.Timestamp.Unix()
			since := bmb.Header.Timestamp.Unix() - prevTime
			bHash := bmb.BlockHashWithAlgos(blk.Height())
			return fmt.Sprintf(
				"new blk height %d %08x %s%10d %08x %v %s %ds since prev",
				blk.Height(),
				prevBlock.WireBlock().Header.Bits,
				bHash,
				bmb.Header.Timestamp.Unix(),
				bmb.Header.Bits,
				amt.Amount(coinbaseTx.Value),
				fork.GetAlgoName(
					bmb.Header.Version,
					blk.Height(),
				), since,
			)
		},
	)
	I.Ln("clearing address used for blk")
	s.nextAddress = nil
	return
}

// hashrate reports from workers
func processHashrateMsg(
	ctx interface{}, src net.Addr, dst string, b []byte,
) (e error) {
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
