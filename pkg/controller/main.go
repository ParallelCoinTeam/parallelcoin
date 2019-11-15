package controller

import (
	"context"
	"crypto/cipher"
	"fmt"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/advertisment"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"sync"
	"time"
)

type msgBuffer struct {
	buffers    [][]byte
	first      time.Time
	decoded    bool
	superseded bool
}

type Controller struct {
	active        *atomic.Bool
	buffers       map[string]*msgBuffer
	ciph          cipher.AEAD
	conns         []*net.UDPConn
	ctx           context.Context
	cx            *conte.Xt
	mx            *sync.Mutex
	oldBlocks     *atomic.Value
	pauseShards   [][]byte
	sendAddresses []*net.UDPAddr
	subMx         *sync.Mutex
	submitChan    chan []byte
}

var (
	SolutionMagic = [4]byte{'s', 'o', 'l', 'v'}
)

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &Controller{
		active:        &atomic.Bool{},
		buffers:       make(map[string]*msgBuffer),
		ciph:          gcm.GetCipher(*cx.Config.MinerPass),
		conns:         []*net.UDPConn{},
		ctx:           ctx,
		cx:            cx,
		mx:            &sync.Mutex{},
		oldBlocks:     &atomic.Value{},
		pauseShards:   [][]byte{},
		sendAddresses: []*net.UDPAddr{},
		subMx:         &sync.Mutex{},
		submitChan:    make(chan []byte),
	}

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
			ctrl.sendAddresses = append(ctrl.sendAddresses, MCAddresses[i])
			conn, err := net.ListenUDP("udp", MCAddresses[i])
			if err != nil {
				log.ERROR(err)
			} else {
				ctrl.conns = append(ctrl.conns, conn)
			}
		}
	}
	// create pause message ready for shutdown handler next
	pM := pause.GetPauseContainer(cx)

	pauseShards, err := Shards(pM.Data, pause.PauseMagic, ctrl.ciph)
	if err != nil {
		log.TRACE(err)
	}
	ctrl.oldBlocks.Store(pauseShards)
	defer func() {
		log.DEBUG("miner controller shutting down")
		for i := range ctrl.sendAddresses {
			err := SendShards(ctrl.sendAddresses[i], pauseShards,
				ctrl.conns[i])
			if err != nil {
				log.ERROR(err)
			}
		}
		for i := range ctrl.conns {
			log.DEBUG("stopping listener on", ctrl.conns[i].LocalAddr())
			err := ctrl.conns[i].Close()
			if err != nil {
				log.ERROR(err)
			}
		}
	}()
	log.DEBUG("sending broadcasts from:", ctrl.sendAddresses)
	// send out the first broadcast
	bTG := getBlkTemplateGenerator(cx)
	tpl, err := bTG.NewBlockTemplate(0, cx.StateCfg.ActiveMiningAddrs[0],
		"sha256d")
	if err != nil {
		log.ERROR(err)
		return
	}
	msgBase := advertisment.Get(cx)
	mC := job.Get(cx, util.NewBlock(tpl.Block), msgBase)
	lisP := mC.GetControllerListenerPort()
	listenAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", lisP))
	if err != nil {
		log.ERROR(err)
		return
	}
	pauseShards, err = sendNewBlockTemplate(cx, bTG, msgBase,
		ctrl.sendAddresses, ctrl.conns, ctrl.oldBlocks, ctrl.ciph)
	if err != nil {
		log.ERROR(err)
	}
	ctrl.active.Store(true)
	log.DEBUG("miner controller starting")
	cx.RealNode.Chain.Subscribe(getNotifier(ctrl.active, bTG, ctrl.ciph,
		ctrl.conns, cx, msgBase, ctrl.oldBlocks, ctrl.sendAddresses,
		ctrl.subMx))
	go rebroadcaster(ctrl)
	go submitter(ctrl)
	cancel, err = Listen(listenAddress, getListener(ctrl))
	if err != nil {
		log.DEBUG(err)
		return
	}
	select {
	case <-ctx.Done():
		ctrl.active.Store(false)
	case <-interrupt.HandlersDone:
	}
	log.DEBUG("controller exiting")
	ctrl.active.Store(false)
	return
}

func submitter(ctrl *Controller) {
out:
	for {
		select {
		case msg := <-ctrl.submitChan:
			log.SPEW(msg)
			decodedB, err := util.NewBlockFromBytes(msg)
			if err != nil {
				log.ERROR(err)
				return
			}
			log.SPEW(decodedB)
		case <-ctrl.ctx.Done():
			break out
		}
	}
}

func getListener(ctrl *Controller) func(a *net.UDPAddr, n int, b []byte) {
	return func(a *net.UDPAddr, n int, b []byte) {
		var err error
		ctrl.mx.Lock()
		defer ctrl.mx.Unlock()
		if n < 16 {
			log.ERROR("received short broadcast message")
			return
		}
		magic := string(b[12:16])
		if magic == string(SolutionMagic[:]) {
			nonce := string(b[:12])
			if bn, ok := ctrl.buffers[nonce]; ok {

				if !bn.decoded {
					payload := b[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					bn.buffers = append(bn.buffers, newP)
					if len(bn.buffers) >= 3 {
						// try to decode it
						var cipherText []byte
						//log.SPEW(bn.buffers)
						cipherText, err = fec.Decode(bn.buffers)
						if err != nil {
							log.ERROR(err)
							return
						}
						//log.SPEW(cipherText)
						msg, err := ctrl.ciph.Open(nil, []byte(nonce),
							cipherText, nil)
						if err != nil {
							log.ERROR(err)
							return
						}
						bn.decoded = true
						ctrl.submitChan <- msg
					}
				} else {
					for i := range ctrl.buffers {
						if i != nonce {
							// superseded blocks can be deleted from the
							// buffers,
							// we don't add more data for the already
							// decoded
							ctrl.buffers[i].superseded = true
						}
					}
				}
			} else {
				ctrl.buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(),
					false, false}
				payload := b[16:n]
				newP := make([]byte, len(payload))
				copy(newP, payload)
				ctrl.buffers[nonce].buffers = append(ctrl.buffers[nonce].buffers,
					newP)
				//log.DEBUGF("%x", payload)
			}
			//log.DEBUGF("%v %v %012x %s", i, a, nonce, magic)
		}
	}
}

func sendNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator,
	msgBase simplebuffer.Serializers, sendAddresses []*net.UDPAddr, conns []*net.UDPConn,
	oldBlocks *atomic.Value, ciph cipher.AEAD, ) (shards [][]byte, err error) {
	template := getNewBlockTemplate(cx, bTG)
	msgB := template.Block
	fMC := job.Get(cx, util.NewBlock(msgB), msgBase)
	for i := range sendAddresses {
		shards, err = Send(sendAddresses[i], fMC.Data, job.WorkMagic, ciph,
			conns[i])
		if err != nil {
			log.ERROR(err)
		}
		oldBlocks.Store(shards)
	}
	return
}

func getNewBlockTemplate(cx *conte.Xt, bTG *mining.BlkTmplGenerator,
) (template *mining.BlockTemplate) {
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

func rebroadcaster(ctrl *Controller) {
	rebroadcastTicker := time.NewTicker(time.Second)
out:
	for {
		select {
		case <-rebroadcastTicker.C:
			for i := range ctrl.conns {
				oB := ctrl.oldBlocks.Load().([][]byte)
				err := SendShards(
					ctrl.sendAddresses[i],
					oB,
					ctrl.conns[i])
				if err != nil {
					log.TRACE(err)
				}
			}
		case <-ctrl.ctx.Done():
			break out
		//default:
		}
	}
}

func  getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
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

func getNotifier(active *atomic.Bool, bTG *mining.BlkTmplGenerator,
	ciph cipher.AEAD, conns []*net.UDPConn, cx *conte.Xt,
	msgBase simplebuffer.Serializers, oldBlocks *atomic.Value, sendAddresses []*net.UDPAddr,
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
				mC := job.Get(cx, util.NewBlock(msgB), msgBase)
				for i := range sendAddresses {
					shards, err := Send(sendAddresses[i], mC.Data,
						job.WorkMagic, ciph, conns[i])
					if err != nil {
						log.TRACE(err)
					}
					oldBlocks.Store(shards)
				}
			}
		}
	}
}
