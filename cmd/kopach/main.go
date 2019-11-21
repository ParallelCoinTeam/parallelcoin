package kopach

import (
	"context"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/transport"
	"go.uber.org/atomic"
	"net"
	"sync"
)

type Worker struct {
	active        *atomic.Bool
	conn          *transport.Connection
	ctx           context.Context
	cx            *conte.Xt
	mx            *sync.Mutex
	sendAddresses []*net.UDPAddr
}

func Main(cx *conte.Xt, quit chan struct{}) {
	log.DEBUG("miner controller starting")
	ctx, cancel := context.WithCancel(context.Background())
	conn, err := transport.NewConnection("", controller.UDP4MulticastAddress,
		*cx.Config.MinerPass, controller.MaxDatagramSize, ctx)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	w := &Worker{
		conn:          conn,
		active:        &atomic.Bool{},
		ctx:           ctx,
		cx:            cx,
		mx:            &sync.Mutex{},
		sendAddresses: []*net.UDPAddr{},
	}
	w.active.Store(false)
	err = w.conn.Listen(handlers)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	log.DEBUG("listening on", controller.UDP4MulticastAddress)
out:
	for {
		select {
		case nB := <-w.conn.ReceiveChan:
			//log.SPEW(nB)
			magicB := nB[:4]
			magic := string(magicB)
			if hnd, ok := handlers[magic]; ok {
				err = hnd(nB)
				if err != nil {
					log.ERROR(err)
				}
			}
			//switch magic {
			//case string(job.WorkMagic):
			//	log.DEBUG("work message")
			//case string(pause.PauseMagic):
			//	log.DEBUG("pause message")
			//}
		case <-quit:
			cancel()
			break out
		}
	}
}

// these are the handlers for specific message types.
// Controller only listens for submissions (currently)
var handlers = map[string]func(b []byte) (err error){
	string(job.WorkMagic): func(b []byte) (err error) {
		log.DEBUG("received job")
		log.SPEW(b)
		j := job.LoadMinerContainer(b)
		log.DEBUG(j.Count())
		log.DEBUG(j.GetIPs())
		log.DEBUG(j.GetP2PListenersPort())
		log.DEBUG(j.GetRPCListenersPort())
		log.DEBUG(j.GetControllerListenerPort())
		log.DEBUG(j.GetNewHeight())
		log.DEBUG(j.GetPrevBlockHash())
		log.DEBUG(j.GetBitses())
		log.SPEW(j.GetTxs())
		return
	},
	string(pause.PauseMagic): func(b []byte) (err error) {
		log.DEBUG("received pause")
		log.SPEW(b)
		j := pause.LoadPauseContainer(b)
		log.DEBUG(j.Count())
		log.DEBUG(j.GetIPs())
		log.DEBUG(j.GetP2PListenersPort())
		log.DEBUG(j.GetRPCListenersPort())
		log.DEBUG(j.GetControllerListenerPort())
		return
	},
}
