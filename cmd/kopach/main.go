package kopach

import (
	"context"
	"net"
	"sync"

	"go.uber.org/atomic"

	"github.com/p9c/pod/cmd/kopach/client"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"github.com/p9c/pod/pkg/transport"
)

type Worker struct {
	active        *atomic.Bool
	conn          *transport.Connection
	ctx           context.Context
	cx            *conte.Xt
	mx            *sync.Mutex
	sendAddresses []*net.UDPAddr
	workers       []*client.Client
}

type handleFunc map[string]func(c *Worker) func(b []byte) (err error)

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
	var workers []*client.Client
	// start up the workers
	for i := 0; i < *cx.Config.GenThreads; i++ {
		// TODO: this needs to be made into a subcommand
		cmd := worker.Spawn("go", "run", "cmd/kopach/kopach_worker/main.go")
		workers = append(workers, client.New(cmd.StdConn))
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
				err = hnd(w)(nB)
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
var handlers = transport.HandleFunc{
	string(job.WorkMagic): func(ctx interface{}) func(b []byte) (err error) {
		return func(b []byte) (err error) {
			w := ctx.(*Worker)
			_ = w
			log.DEBUG("received job")
			//log.SPEW(b)
			j := job.LoadMinerContainer(b)
			log.DEBUG(j.String())
			return
		}
	},
	string(pause.PauseMagic): func(ctx interface{}) func(b []byte) (err error) {
		return func(b []byte) (err error) {
			w := ctx.(*Worker)
			_ = w
			log.DEBUG("received pause")
			//log.SPEW(b)
			j := pause.LoadPauseContainer(b)
			log.DEBUG(j.Count())
			log.DEBUG(j.GetIPs())
			log.DEBUG(j.GetP2PListenersPort())
			log.DEBUG(j.GetRPCListenersPort())
			log.DEBUG(j.GetControllerListenerPort())
			return
		}
	},
}
