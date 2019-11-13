package controller

import (
	"context"
	chain "github.com/p9c/pod/pkg/chain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
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
	WorkMagic     = [4]byte{'w', 'o', 'r', 'k'}
	PauseMagic    = [4]byte{'p', 'a', 'u', 's'}
	SolutionMagic = [4]byte{'s', 'o', 'l', 'v'}
)

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	var active atomic.Bool
	var oldBlocks atomic.Value
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
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	var sendAddresses []*net.UDPAddr
	for i := range MCAddresses {
		_, err := net.ListenUDP("udp", MCAddresses[i])
		if err == nil {
			sendAddresses = append(sendAddresses, MCAddresses[i])
		}
	}
	var conns []*net.UDPConn
	for i := range sendAddresses {
		conn, err := net.ListenUDP("udp", sendAddresses[i])
		if err != nil {
			log.ERROR(err)
		} else {
			conns = append(conns, conn)
		}
	}
	var pauseShards [][]byte
	pM := GetMessageBase(cx).CreateContainer(PauseMagic)
	shards, err := Shards(pM.Data, PauseMagic, *ciph)
	if err != nil {
		log.TRACE(err)
	}
	pauseShards = shards
	defer func() {
		log.DEBUG("miner controller shutting down")
		for i := range sendAddresses {
			err := SendShards(sendAddresses[i], pauseShards, conns[i])
			if err != nil {
				log.ERROR(err)
			}
		}
		for i := range conns {
			log.DEBUG("stopping listener on", conns[i].LocalAddr())
			err := conns[i].Close()
			if err != nil {
				log.ERROR(err)
			}
		}
	}()
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
	active.Store(true)
	for i := range sendAddresses {
		shards, err := Send(sendAddresses[i], fMC.Data, WorkMagic, *ciph,
			conns[i])
		if err != nil {
			log.ERROR(err)
		} else {
			oldBlocks.Store(shards)
		}
	}
	//log.SPEW(messageBase.CreateContainer(WorkMagic))
	// There is no unsubscribe but we can use an atomic to disable the
	// function instead - this also ensures that new work doesn't start
	// once the context is cancelled below
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
				for i := range sendAddresses {
					shards, err := Send(sendAddresses[i], mC.Data,
						WorkMagic, *ciph, conns[i])
					if err != nil {
						log.TRACE(err)
					}
					oldBlocks.Store(shards)
				}
			}
		}
	})
	go func() {
		rebroadcastTicker := time.NewTicker(time.Second)
	out:
		for {
			select {
			case <-rebroadcastTicker.C:
				for i := range conns {
					oB := oldBlocks.Load().([][]byte)
					err = SendShards(sendAddresses[i], oB, conns[i])
					if err != nil {
						log.TRACE(err)
					}
				}
			case <-ctx.Done():
				active.Store(false)
				break out
			default:
			}
		}
	}()
	select {
	case <-ctx.Done():
	case <-interrupt.HandlersDone:
	}
	log.DEBUG("controller exiting")
	active.Store(false)
	return
}
