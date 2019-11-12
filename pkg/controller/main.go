package controller

import (
	"context"
	chain "github.com/p9c/pod/pkg/chain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Blocks struct {
	PrevBlock    *chainhash.Hash
	Bits         uint32
	Transactions []*wire.MsgTx
	Listeners    []string
}

var (
	WorkMagic = [4]byte{'w', 'o', 'r', 'k'}
)

// Blocks is a block broadcast message for miners to mine from
type Blocks struct {
	// New is a flag that distinguishes a newly accepted/connected block from a rebroadcast
	New bool
	// Payload is a map of bytes indexed by block version number
	Payload map[int32][]byte
}

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	var active atomic.Bool
	ctx, cancel = context.WithCancel(context.Background())
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		log.WARN("not running controller without RPC enabled")
		cancel()
		return
	}
	if len(*cx.Config.Listeners) < 1 || *cx.Config.DisableListen {
		log.WARN("not running controller without p2p listener enabled")
		cancel()
		return
	}

	var sendAddresses []*net.UDPAddr
	for i := range MCAddresses {
		_, err := net.ListenUDP("udp", MCAddresses[i])
		if err == nil {
			sendAddresses = append(sendAddresses, MCAddresses[i])
		}
	}
	//log.SPEW(sendAddresses)
	log.DEBUG("sending broadcasts from:", sendAddresses)
	policy := mining.Policy{
		BlockMinWeight:    uint32(*cx.Config.BlockMinWeight),
		BlockMaxWeight:    uint32(*cx.Config.BlockMaxWeight),
		BlockMinSize:      uint32(*cx.Config.BlockMinSize),
		BlockMaxSize:      uint32(*cx.Config.BlockMaxSize),
		BlockPrioritySize: uint32(*cx.Config.BlockPrioritySize),
		TxMinFreeFee:      cx.StateCfg.ActiveMinRelayTxFee,
	}
	s := cx.RealNode
	bTG := mining.NewBlkTmplGenerator(&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache, s.Algo)
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
		MiningAddrs))]
	algo := "sha256d"
	template, err := bTG.NewBlockTemplate(0, payToAddr,
		algo)
	if err != nil {
		log.ERROR(err)
	}
	msgB := template.Block
	msgBase := GetMessageBase(cx)
	fMC := GetMinerContainer(cx, util.NewBlock(msgB), msgBase)
	//log.DEBUG(fMC.GetIPs())
	//log.DEBUG(fMC.GetP2PListenersPort())
	//log.DEBUG(fMC.GetRPCListenersPort())
	//log.DEBUG(fMC.GetControllerListenerPort())
	//log.DEBUG(fMC.GetPrevBlockHash())
	//log.DEBUG(fMC.GetBitses())
	//log.SPEW(fMC.GetTxs())
	//log.SPEW(fMC.Data)
	for i := range sendAddresses {
		err := Send(sendAddresses[i], fMC.Data, WorkMagic)
		if err != nil {
			log.ERROR(err)
		}
	}
	//log.SPEW(messageBase.CreateContainer(WorkMagic))
	go func() {
		// There is no unsubscribe but we can use an atomic to disable the
		// function instead - this also ensures that new work doesn't start
		// once the context is cancelled below
		active.Store(true)
		var subMx sync.Mutex
		log.DEBUG("miner controller starting")
		cx.RealNode.Chain.Subscribe(func(n *chain.Notification) {
			if active.Load() {
				// first to arrive locks out any others while processing
				switch n.Type {
				case chain.NTBlockAccepted:
					subMx.Lock()
					defer subMx.Unlock()
					log.DEBUG("received new chain notification")
					// construct work message
					//log.SPEW(n)
					_, ok := n.Data.(*util.Block)
					if !ok {
						log.WARN("chain accepted notification is not a block")
						break
					}
					// Choose a payment address at random.
					rand.Seed(time.Now().UnixNano())
					payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
						MiningAddrs))]
					algo := "sha256d"
					template, err := bTG.NewBlockTemplate(0, payToAddr,
						algo)
					if err != nil {
						log.ERROR(err)
					}
					msgB := template.Block
					mC := GetMinerContainer(cx, util.NewBlock(msgB), msgBase)
					//log.DEBUG(mC.GetIPs())
					//log.DEBUG(mC.GetP2PListenersPort())
					//log.DEBUG(mC.GetRPCListenersPort())
					//log.DEBUG(mC.GetControllerListenerPort())
					//log.DEBUG(mC.GetPrevBlockHash())
					//log.DEBUG(mC.GetBitses())
					//log.SPEW(mC.GetTxs())
					//log.SPEW(mC.Data)
					for i := range sendAddresses {
						err := Send(sendAddresses[i], mC.Data, WorkMagic)
						if err != nil {
							log.TRACE(err)
						}
					}
				}
			}
		})
		select {
		case <-ctx.Done():
			log.DEBUG("miner controller shutting down")
			active.Store(false)
			break
		}
	}()
	return
}
