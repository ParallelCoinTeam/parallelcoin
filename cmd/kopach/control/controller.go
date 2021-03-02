package control

import (
	"container/ring"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"net"
	"sync"
	"time"
	
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util/quit"
	
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/comm/transport"
	rav "github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const (
	MaxDatagramSize      = 8192
	UDP4MulticastAddress = "224.0.0.1:11049"
	BufferSize           = 4096
)

type Controller struct {
	multiConn              *transport.Channel
	active                 atomic.Bool
	quit                   qu.C
	cx                     *conte.Xt
	isMining               atomic.Bool
	height                 atomic.Int32
	blockTemplateGenerator *mining.BlkTmplGenerator
	msgBlockTemplate       *templates.Message
	coinbases              atomic.Value
	transactions           atomic.Value
	txMx                   sync.Mutex
	oldBlocks              atomic.Value
	prevHash               atomic.Value
	lastTxUpdate           atomic.Value
	lastGenerated          atomic.Value
	pauseShards            [][]byte
	sendAddresses          []*net.UDPAddr
	buffer                 *ring.Ring
	began                  time.Time
	otherNodes             map[uint64]*nodeSpec
	uuid                   uint64
	hashCount              atomic.Uint64
	hashSampleBuf          *rav.BufferUint64
	lastNonce              int32
	walletClient           *rpcclient.Client
}

type nodeSpec struct {
	time.Time
	addr string
}

// Run starts up a controller
func Run(cx *conte.Xt) (quit qu.C) {
	if *cx.Config.DisableController {
		Info("controller is disabled")
		return
	}
	cx.Controller.Store(true)
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		Warn("not running controller without RPC enabled")
		return
	}
	if len(*cx.Config.P2PListeners) < 1 || *cx.Config.DisableListen {
		Warn("not running controller without p2p listener enabled", *cx.Config.P2PListeners)
		return
	}
	nS := make(map[uint64]*nodeSpec)
	c := &Controller{
		quit:                   qu.T(),
		cx:                     cx,
		sendAddresses:          []*net.UDPAddr{},
		blockTemplateGenerator: getBlkTemplateGenerator(cx),
		buffer:                 ring.New(BufferSize),
		began:                  time.Now(),
		otherNodes:             nS,
		uuid:                   cx.UUID,
		hashSampleBuf:          rav.NewBufferUint64(100),
	}
	c.isMining.Store(true)
	// maintain connection to wallet if it is available
	var err error
	go c.walletRPCWatcher()
	c.prevHash.Store(&chainhash.Hash{})
	quit = c.quit
	c.lastTxUpdate.Store(time.Now().UnixNano())
	c.lastGenerated.Store(time.Now().UnixNano())
	c.height.Store(0)
	c.active.Store(false)
	if c.multiConn, err = transport.NewBroadcastChannel(
		"controller", c, *cx.Config.MinerPass, transport.DefaultPort, MaxDatagramSize, handlersMulticast,
		quit,
	); Check(err) {
		c.quit.Q()
		return
	}
	if c.pauseShards = transport.GetShards(p2padvt.Get(cx)); Check(err) {
	} else {
		c.active.Store(true)
	}
	interrupt.AddHandler(
		func() {
			Debug("miner controller shutting down")
			c.active.Store(false)
			if err = c.multiConn.SendMany(pause.Magic, c.pauseShards); Check(err) {
			}
			if err = c.multiConn.Close(); Check(err) {
			}
			c.quit.Q()
		},
	)
	Debug("sending broadcasts to:", UDP4MulticastAddress)
	
	go c.advertiserAndRebroadcaster()
	return
}
