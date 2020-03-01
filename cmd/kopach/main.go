package kopach

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
	
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/cmd/kopach/client"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
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
	active        atomic.Bool
	conn          *transport.Channel
	ctx           context.Context
	cx            *conte.Xt
	sendAddresses []*net.UDPAddr
	workers       []*client.Client
	FirstSender   atomic.String
	lastSent      atomic.Int64
	Status        atomic.String
	HashTick      chan HashCount
	LastHash      *chainhash.Hash
}

func KopachHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		// log.L.SetLevel("trace", true)
		log.DEBUG("miner controller starting")
		ctx, cancel := context.WithCancel(context.Background())
		w := &Worker{
			ctx:           ctx,
			cx:            cx,
			sendAddresses: []*net.UDPAddr{},
		}
		w.lastSent.Store(time.Now().UnixNano())
		w.active.Store(false)
		log.DEBUG("opening broadcast channel listener")
		w.conn, err = transport.
			NewBroadcastChannel("kopachmain", w, *cx.Config.MinerPass,
				transport.DefaultPort, controller.MaxDatagramSize, handlers)
		if err != nil {
			log.ERROR(err)
			cancel()
			return
		}
		var wks []*worker.Worker
		// start up the workers
		log.DEBUG("starting up kopach workers")
		for i := 0; i < *cx.Config.GenThreads; i++ {
			log.DEBUG("starting worker", i)
			cmd := worker.Spawn(os.Args[0], "worker",
				cx.ActiveNet.Name, *cx.Config.LogLevel)
			wks = append(wks, cmd)
			w.workers = append(w.workers, client.New(cmd.StdConn))
		}
		interrupt.AddHandler(func() {
			log.DEBUG("KopachHandle interrupt")
			for i := range w.workers {
				if err := wks[i].Kill(); log.Check(err) {
				}
				log.DEBUG("stopped worker", i)
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
		w.active.Store(true)
		// controller watcher thread
		go func() {
			log.DEBUG("starting controller watcher")
			ticker := time.NewTicker(time.Second)
		out:
			for {
				select {
				case <-ticker.C:
					// log.DEBUG("tick", w.lastSent, w.FirstSender)
					// if the last message sent was 3 seconds ago the server is
					// almost certainly disconnected or crashed so clear FirstSender
					since := time.Now().Sub(time.Unix(0 , int64(w.lastSent.Load())))
					wasSending := since > time.Second*3 && w.FirstSender.Load() != ""
					if wasSending {
						log.DEBUG("previous current controller has stopped" +
							" broadcasting")
						// when this string is clear other broadcasts will be
						// listened to
						w.FirstSender.Store("")
						// pause the workers
						for i := range w.workers {
							log.DEBUG("sending pause to worker", i)
							err := w.workers[i].Pause()
							if err != nil {
								log.ERROR(err)
							}
						}
					}
				case <-cx.KillAll:
					break out
				}
			}
		}()
		log.DEBUG("listening on", controller.UDP4MulticastAddress)
		<-cx.KillAll
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
		// h := j.GetHashes()
		ips := j.GetIPs()
		cP := j.GetControllerListenerPort()
		addr := net.JoinHostPort(ips[0].String(), fmt.Sprint(cP))
		firstSender := w.FirstSender.Load()
		otherSent := firstSender != addr && firstSender != ""
		if otherSent {
			// ignore other controllers while one is active and received first
			return
		} else {
			w.FirstSender.Store(addr)
			w.lastSent.Store(time.Now().UnixNano())
		}
		// log.TRACE("received job")
		for i := range w.workers {
			// log.TRACE("sending job to worker", i)
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
		// clear the FirstSender
		w.FirstSender.Store("")
		return
	},
}
