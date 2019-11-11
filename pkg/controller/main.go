package controller

import (
	"context"
	"fmt"
	chain "github.com/p9c/pod/pkg/chain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"go.uber.org/atomic"
	"net"
	"sync"
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
					mB, ok := n.Data.(*util.Block)
					if !ok {
						log.WARN("chain accepted notification is not a block")
						break
					}
					mC := GetMinerContainer(cx, mB)
					for _, i := range []string{
						UDP4MulticastAddress,
						UDP6MulticastAddress,
					} {
						err := Send(net.JoinHostPort(i,
							fmt.Sprint(mC.GetControllerListenerPort())),
							mC.Data)
						if err != nil {
							log.ERROR(err)
						}
					}
					//mW := LoadMinerContainer(mC.Data)
					// send out srs.Data
					//log.SPEW(srs.Data)
					// the following decodes each element
					//mC := LoadMinerContainer(srs)
					//out := []interface{}{
					//	mW.GetIPs(),
					//	mW.GetP2PListenersPort(),
					//	mW.GetRPCListenersPort(),
					//	mW.GetControllerListenerPort(),
					//	mW.GetPrevBlockHash(),
					//	mW.GetBitses(),
					//	mW.GetTxs(),
					//}
					//log.SPEW(out)
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
