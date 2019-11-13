package controller

import (
	"context"
	"crypto/cipher"
	"fmt"
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
	var sendAddresses []*net.UDPAddr
	var conns []*net.UDPConn
	var pauseShards [][]byte
	var subMx sync.Mutex
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	msgBase := GetMessageBase(cx)
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
	// test the addresses and collate the ones that work
	for i := range MCAddresses {
		_, err := net.ListenUDP("udp", MCAddresses[i])
		if err == nil {
			sendAddresses = append(sendAddresses, MCAddresses[i])
			conn, err := net.ListenUDP("udp", MCAddresses[i])
			if err != nil {
				log.ERROR(err)
			} else {
				conns = append(conns, conn)
			}
		}
	}
	// create pause message ready for shutdown handler next
	pM := GetMessageBase(cx).CreateContainer(PauseMagic)

	pauseShards, err := Shards(pM.Data, PauseMagic, ciph)
	if err != nil {
		log.TRACE(err)
	}
	oldBlocks.Store(pauseShards)
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
	log.DEBUG("sending broadcasts from:", sendAddresses)
	// send out the first broadcast
	bTG := getBlkTemplateGenerator(cx)
	tpl, err := bTG.NewBlockTemplate(0, cx.StateCfg.ActiveMiningAddrs[0],
		"sha256d")
	if err != nil {
		log.ERROR(err)
		return
	}
	mC := GetMinerContainer(cx, util.NewBlock(tpl.Block), msgBase)
	lisP := mC.GetControllerListenerPort()
	listenAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", lisP))
	if err != nil {
		log.ERROR(err)
		return
	}
	pauseShards, err = sendNewBlockTemplate(cx, bTG, msgBase, sendAddresses,
		conns, &oldBlocks, ciph)
	if err != nil {
		log.ERROR(err)
	}
	active.Store(true)
	log.DEBUG("miner controller starting")
	cx.RealNode.Chain.Subscribe(getNotifier(&active, bTG, ciph, conns, cx,
		msgBase, &oldBlocks, sendAddresses, &subMx))
	go rebroadcaster(conns, &oldBlocks, sendAddresses, ctx)
	cancel, err = Listen(listenAddress, getListener())
	if err != nil {
		log.DEBUG(err)
		return
	}
	select {
	case <-ctx.Done():
		active.Store(false)
	case <-interrupt.HandlersDone:
	}
	log.DEBUG("controller exiting")
	cancel()
	active.Store(false)
	return
}

func getListener() func(a *net.UDPAddr, n int, b []byte) {
	return func(a *net.UDPAddr, n int, b []byte) {
		log.DEBUG(a)
		received := b[:n]
		_ = received

	}
}

func sendNewBlockTemplate(
	cx *conte.Xt,
	bTG *mining.BlkTmplGenerator,
	msgBase Serializers,
	sendAddresses []*net.UDPAddr,
	conns []*net.UDPConn,
	oldBlocks *atomic.Value,
	ciph cipher.AEAD,
) (shards [][]byte, err error) {
	template := getNewBlockTemplate(cx, bTG)
	msgB := template.Block
	fMC := GetMinerContainer(cx, util.NewBlock(msgB), msgBase)
	for i := range sendAddresses {
		shards, err = Send(sendAddresses[i], fMC.Data, WorkMagic, ciph,
			conns[i])
		if err != nil {
			log.ERROR(err)
		}
		oldBlocks.Store(shards)
	}
	return
}

func getNewBlockTemplate(
	cx *conte.Xt,
	bTG *mining.BlkTmplGenerator,
) (
	template *mining.BlockTemplate) {
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(*cx.Config.
		MiningAddrs))]
	template, err := bTG.NewBlockTemplate(0, payToAddr,
		"sha256d")
	if err != nil {
		log.ERROR(err)
	}
	return
}

func rebroadcaster(
	conns []*net.UDPConn,
	oldBlocks *atomic.Value,
	sendAddresses []*net.UDPAddr,
	ctx context.Context,
) {
	rebroadcastTicker := time.NewTicker(time.Second)
out:
	for {
		select {
		case <-rebroadcastTicker.C:
			for i := range conns {
				oB := oldBlocks.Load().([][]byte)
				err := SendShards(
					sendAddresses[i],
					oB,
					conns[i])
				if err != nil {
					log.TRACE(err)
				}
			}
		case <-ctx.Done():
			break out
		default:
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
	return mining.NewBlkTmplGenerator(&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache, s.Algo)
}

func getNotifier(
	active *atomic.Bool,
	bTG *mining.BlkTmplGenerator,
	ciph cipher.AEAD,
	conns []*net.UDPConn,
	cx *conte.Xt,
	msgBase Serializers,
	oldBlocks *atomic.Value,
	sendAddresses []*net.UDPAddr,
	subMx *sync.Mutex,
) func(n *chain.Notification) {
	return func(n *chain.Notification) {
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
				template := getNewBlockTemplate(cx, bTG)
				msgB := template.Block
				mC := GetMinerContainer(cx, util.NewBlock(msgB), msgBase)
				for i := range sendAddresses {
					shards, err := Send(sendAddresses[i], mC.Data,
						WorkMagic, ciph, conns[i])
					if err != nil {
						log.TRACE(err)
					}
					oldBlocks.Store(shards)
				}
			}
		}
	}
}
