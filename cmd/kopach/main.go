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
		log.DEBUG("starting worker", i)
		cmd := worker.Spawn("go", "run", "cmd/kopach/kopach_worker/main.go",
			cx.ActiveNet.Name)
		workers = append(workers, client.New(cmd.StdConn))
	}
	w := &Worker{
		conn:          conn,
		active:        &atomic.Bool{},
		ctx:           ctx,
		cx:            cx,
		mx:            &sync.Mutex{},
		sendAddresses: []*net.UDPAddr{},
		workers:       workers,
	}
	w.active.Store(false)
	for i := range w.workers {
		log.DEBUG("sending pass to worker", i)
		err := w.workers[i].SendPass(*cx.Config.MinerPass)
		if err != nil {
			log.ERROR(err)
		}
	}
	err = w.conn.Listen(handlers, w)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	log.DEBUG("listening on", controller.UDP4MulticastAddress)
	<-quit
	log.INFO("kopach shutting down")
}

// these are the handlers for specific message types.
var handlers = transport.HandleFunc{
	string(job.WorkMagic): func(ctx interface{}) func(b []byte) (err error) {
		return func(b []byte) (err error) {
			log.DEBUG("received job")
			w := ctx.(*Worker)
			_ = w
			//log.SPEW(b)
			j := job.LoadMinerContainer(b)
			log.DEBUG(j.String())
			//log.DEBUG("workers", len(w.workers))
			for i := range w.workers {
				log.DEBUG("sending job to worker", i)
				err := w.workers[i].NewJob(&j)
				if err != nil {
					log.ERROR(err)
				}
			}
			return
		}
	},
	string(pause.PauseMagic): func(ctx interface{}) func(b []byte) (err error) {
		return func(b []byte) (err error) {
			log.DEBUG("received pause")
			w := ctx.(*Worker)
			_ = w
			//log.SPEW(b)
			j := pause.LoadPauseContainer(b)
			log.DEBUG(j.Count())
			log.DEBUG(j.GetIPs())
			log.DEBUG(j.GetP2PListenersPort())
			log.DEBUG(j.GetRPCListenersPort())
			log.DEBUG(j.GetControllerListenerPort())
			for i := range w.workers {
				log.DEBUG("sending pause to worker", i)
				err := w.workers[i].Pause()
				if err != nil {
					log.ERROR(err)
				}
			}

			return
		}
	},
}
