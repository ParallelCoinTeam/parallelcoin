package kopach

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

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
	firstSender   string
	lastSent      time.Time
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
		lastSent:      time.Now(),
	}
	w.active.Store(false)
	for i := range w.workers {
		log.DEBUG("sending pass to worker", i)
		err := w.workers[i].SendPass(*cx.Config.MinerPass)
		if err != nil {
			log.ERROR(err)
		}
	}
	err = w.conn.Listen(handlers, w, &w.lastSent, &w.firstSender)
	if err != nil {
		log.ERROR(err)
		cancel()
		return
	}
	// controller watcher thread
	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				//log.DEBUG("tick", w.lastSent, w.firstSender)
				// if the last message sent was 3 seconds ago the server is
				// almost certainly disconnected or crashed so clear firstSender
				w.mx.Lock()
				since := time.Now().Sub(w.lastSent)
				wasSending := since > time.Second*3 && w.firstSender != ""
				w.mx.Unlock()
				if wasSending {
					log.DEBUG("previous current controller has stopped" +
						" broadcasting")
					// when this string is clear other broadcasts will be
					// listened to
					w.mx.Lock()
					w.firstSender = ""
					w.mx.Unlock()
					// pause the workers
					for i := range w.workers {
						log.DEBUG("sending pause to worker", i)
						err := w.workers[i].Pause()
						if err != nil {
							log.ERROR(err)
						}
					}
				}
			case <-quit:
			}
		}
	}()
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
			j := job.LoadMinerContainer(b)
			ips := j.GetIPs()
			cP := j.GetControllerListenerPort()
			addr := net.JoinHostPort(ips[0].String(), fmt.Sprint(cP))
			w.mx.Lock()
			otherSent := w.firstSender != addr && w.firstSender != ""
			w.mx.Unlock()
			if otherSent {
				// ignore other controllers while one is active and received
				// first
				log.DEBUG("ignoring other controller", addr)
				return
			} else {
				w.mx.Lock()
				w.firstSender = addr
				w.lastSent = time.Now()
				w.mx.Unlock()
			}
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
			for i := range w.workers {
				log.DEBUG("sending pause to worker", i)
				err := w.workers[i].Pause()
				if err != nil {
					log.ERROR(err)
				}
			}
			w.mx.Lock()
			// clear the firstSender
			w.firstSender = ""
			w.mx.Unlock()
			return
		}
	},
}
