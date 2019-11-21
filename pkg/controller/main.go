//+build ignore

package controller

import (
	"context"
	"crypto/cipher"
	"errors"
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

type MsgBuffer struct {
	Buffers    [][]byte
	First      time.Time
	Decoded    bool
	Superseded bool
}

type Controller struct {
	active        *atomic.Bool
	buffers       map[string]*MsgBuffer
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
	cancel = func() {}
	log.DEBUG("miner controller starting")
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &Controller{
		active:        &atomic.Bool{},
		buffers:       make(map[string]*MsgBuffer),
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
	uc, err := net.DialUDP("udp", nil, MCAddress)
	if err != nil {
		log.ERROR(err)
	}
	ctrl.conns = []*net.UDPConn{uc}
	ctrl.sendAddresses = []*net.UDPAddr{MCAddress}
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
	// send out the First broadcast
	bTG := getBlkTemplateGenerator(cx)
	msgBase := pause.GetPauseContainer(cx)
	lisP := msgBase.GetControllerListenerPort()
	addy := fmt.Sprintf("%s:%d", msgBase.GetIPs()[0], lisP)
	listenAddress, err := net.ResolveUDPAddr("udp", addy)
	submitListenConn, err := net.ListenUDP("udp", listenAddress)
	if err != nil {
		log.ERROR(err)
		return
	}
	adv := advertisment.Get(cx)
	pauseShards, err = sendNewBlockTemplate(cx, bTG, adv,
		ctrl.sendAddresses, ctrl.conns, ctrl.oldBlocks, ctrl.ciph)
	if err != nil {
		log.ERROR(err)
		ctrl.active.Store(false)
	} else {
		ctrl.active.Store(true)
	}
	cx.RealNode.Chain.Subscribe(getNotifier(ctrl.active, bTG, ctrl.ciph,
		ctrl.conns, cx, adv, ctrl.oldBlocks, ctrl.sendAddresses,
		ctrl.subMx))
	log.DEBUG("miner controller submit port listening on",
		submitListenConn.LocalAddr())
	cancel, err = Listen(submitListenConn, getMsgHandler(ctrl))
	if err != nil {
		log.DEBUG(err)
		return
	}
	go rebroadcaster(ctrl)
	go submitter(ctrl)
	select {
	case <-ctx.Done():
		ctrl.active.Store(false)
	case <-interrupt.HandlersDone:
	}
	log.TRACE("controller exiting")
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
			//
		case <-ctrl.ctx.Done():
			break out
		}
	}
}

func getMsgHandler(ctrl *Controller) func(a *net.UDPAddr, n int, b []byte) {
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

				if !bn.Decoded {
					payload := b[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					bn.Buffers = append(bn.Buffers, newP)
					if len(bn.Buffers) >= 3 {
						// try to decode it
						var cipherText []byte
						//log.SPEW(bn.Buffers)
						cipherText, err = fec.Decode(bn.Buffers)
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
						bn.Decoded = true
						ctrl.submitChan <- msg
					}
				} else {
					for i := range ctrl.buffers {
						if i != nonce {
							// Superseded blocks can be deleted from the
							// Buffers,
							// we don't add more data for the already
							// Decoded
							ctrl.buffers[i].Superseded = true
						}
					}
				}
			} else {
				ctrl.buffers[nonce] = &MsgBuffer{[][]byte{}, time.Now(),
					false, false}
				payload := b[16:n]
				newP := make([]byte, len(payload))
				copy(newP, payload)
				ctrl.buffers[nonce].Buffers = append(ctrl.buffers[nonce].Buffers,
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
	if template == nil {
		return nil, errors.New("could not get template")
	}
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
	if len(*cx.Config.MiningAddrs) < 1 {
		return
	}
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
	rebroadcastTicker := time.NewTicker(time.Second * 2)
out:
	for {
		select {
		case <-rebroadcastTicker.C:
			for i := range ctrl.conns {
				oB := ctrl.oldBlocks.Load().([][]byte)
				if len(oB) == 0 {
					log.DEBUG("template is empty")
					break
				}
				err := SendShards(
					ctrl.sendAddresses[i],
					oB,
					ctrl.conns[i])
				if err != nil {
					log.ERROR(err)
				}
			}
		case <-ctrl.ctx.Done():
			break out
			//default:
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

func getNotifier(active *atomic.Bool, bTG *mining.BlkTmplGenerator,
	ciph cipher.AEAD, conns []*net.UDPConn, cx *conte.Xt,
	msgBase simplebuffer.Serializers, oldBlocks *atomic.Value, sendAddresses []*net.UDPAddr,
	subMx *sync.Mutex,
) func(n *chain.Notification) {
	return func(n *chain.Notification) {
		if active.Load() {
			// First to arrive locks out any others while processing
			switch n.Type {
			case chain.NTBlockAccepted:
				subMx.Lock()
				defer subMx.Unlock()
				log.TRACE("received new chain notification")
				// construct work message
				//log.SPEW(n)
				_, ok := n.Data.(*util.Block)
				if !ok {
					log.WARN("chain accepted notification is not a block")
					break
				}
				template := getNewBlockTemplate(cx, bTG)
				if template != nil {
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
}
