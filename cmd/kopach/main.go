package kopach

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/cmd/kopach/client"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/pause"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type HashCount struct {
	uint64
	Time time.Time
}

type Worker struct {
	active        *atomic.Bool
	conn          *transport.Channel
	ctx           context.Context
	cx            *conte.Xt
	mx            *sync.Mutex
	sendAddresses []*net.UDPAddr
	workers       []*client.Client
	FirstSender   string
	lastSent      time.Time
	Status        atomic.String
	HashTick      chan HashCount
	LastHash      string
}

func KopachHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		// log.L.SetLevel("trace", true)
		log.DEBUG("miner controller starting")
		quit := make(chan struct{})
		ctx, cancel := context.WithCancel(context.Background())
		w := &Worker{
			active:        &atomic.Bool{},
			ctx:           ctx,
			cx:            cx,
			mx:            &sync.Mutex{},
			sendAddresses: []*net.UDPAddr{},
			lastSent:      time.Now(),
		}
		w.conn, err = transport.
			NewBroadcastChannel("kopachmain", w, *cx.Config.MinerPass, 11049,
				controller.MaxDatagramSize, handlers)
		if err != nil {
			log.ERROR(err)
			cancel()
			return
		}
		var wks []*worker.Worker
		// start up the workers
		for i := 0; i < *cx.Config.GenThreads; i++ {
			log.DEBUG("starting worker", i)
			cmd := worker.Spawn(os.Args[0], "worker",
				cx.ActiveNet.Name)
			wks = append(wks, cmd)
			w.workers = append(w.workers, client.New(cmd.StdConn))
		}
		interrupt.AddHandler(func() {
			log.DEBUG("KopachHandle interrupt")
			for i := range w.workers {
				if err := wks[i].Kill(); log.Check(err) {
				}
			}
		})
		w.active.Store(false)
		for i := range w.workers {
			log.DEBUG("sending pass to worker", i)
			err := w.workers[i].SendPass(*cx.Config.MinerPass)
			if err != nil {
				log.ERROR(err)
			}
		}
		// controller watcher thread
		go func() {
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					// log.DEBUG("tick", w.lastSent, w.FirstSender)
					// if the last message sent was 3 seconds ago the server is
					// almost certainly disconnected or crashed so clear FirstSender
					w.mx.Lock()
					since := time.Now().Sub(w.lastSent)
					wasSending := since > time.Second*3 && w.FirstSender != ""
					w.mx.Unlock()
					if wasSending {
						log.DEBUG("previous current controller has stopped" +
							" broadcasting")
						// when this string is clear other broadcasts will be
						// listened to
						w.mx.Lock()
						w.FirstSender = ""
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
		return
	}
}

// these are the handlers for specific message types.
var handlers = transport.Handlers{
	string(job.WorkMagic): func(ctx interface{}, src *net.UDPAddr, dst string,
		b []byte) (err error) {
		w := ctx.(*Worker)
		j := job.LoadContainer(b)
		h := j.GetHashes()
		w.mx.Lock()
		ips := j.GetIPs()
		cP := j.GetControllerListenerPort()
		addr := net.JoinHostPort(ips[0].String(), fmt.Sprint(cP))
		otherSent := w.FirstSender != addr && w.FirstSender != ""
		w.FirstSender = addr
		w.lastSent = time.Now()
		w.mx.Unlock()
		if len(h) > 0 {
			// log.DEBUG(h)
			hS := h[5].String()
			if w.LastHash == hS {
				log.TRACE("not responding to same job")
				return
			} else {
				w.LastHash = hS
			}
		}
		log.TRACE("received job")
		// if newHash == w.LastHash {
		// 	return
		// } else {
		// 	w.LastHash = newHash
		// }
		if otherSent {
			// ignore other controllers while one is active and received
			// first
			log.DEBUG("ignoring other controller", addr)
			return
		} else {
			w.mx.Lock()
			w.FirstSender = addr
			w.lastSent = time.Now()
			w.mx.Unlock()
		}
		for i := range w.workers {
			log.TRACE("sending job to worker", i)
			err := w.workers[i].NewJob(&j)
			if err != nil {
				log.ERROR(err)
			}
			// log.SPEW(j)
		}
		return
	},
	string(pause.PauseMagic): func(ctx interface{}, src *net.UDPAddr, dst string,
		b []byte) (err error) {
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
		// clear the FirstSender
		w.FirstSender = ""
		w.mx.Unlock()
		return
	},
}
