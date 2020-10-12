package kopach

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/atomic"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/client"
	"github.com/p9c/pod/cmd/kopach/control"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/comm/transport"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type HashCount struct {
	uint64
	Time time.Time
}

type SolutionData struct {
	time   time.Time
	height int
	block  *wire.MsgBlock
}

type Worker struct {
	cx                  *conte.Xt
	height              int32
	active              atomic.Bool
	conn                *transport.Channel
	ctx                 context.Context
	quit                chan struct{}
	sendAddresses       []*net.UDPAddr
	clients             []*client.Client
	workers             []*worker.Worker
	FirstSender         atomic.String
	lastSent            atomic.Int64
	Status              atomic.String
	HashTick            chan HashCount
	LastHash            *chainhash.Hash
	StartChan, StopChan chan struct{}
	SetThreads          chan int
	solutions           []SolutionData
	solutionCount       int
}

func (w *Worker) Start(cx *conte.Xt) {
	// if !*cx.Config.Generate {
	// 	Debug("called start but not running generate")
	// 	return
	// }
	Debug("starting up kopach workers")
	w.workers = []*worker.Worker{}
	w.clients = []*client.Client{}
	for i := 0; i < *cx.Config.GenThreads; i++ {
		Debug("starting worker", i)
		cmd, _ := worker.Spawn(os.Args[0], "worker", fmt.Sprint("worker", i), cx.ActiveNet.Name, *cx.Config.LogLevel)
		w.workers = append(w.workers, cmd)
		w.clients = append(w.clients, client.New(cmd.StdConn))
	}
	for i := range w.clients {
		Debug("sending pass to worker", i)
		err := w.clients[i].SendPass(*cx.Config.MinerPass)
		if err != nil {
			Error(err)
		}
	}
	w.active.Store(true)
	interrupt.AddHandler(func() {
		w.Stop()
	})
}

func (w *Worker) Stop() {
	var err error
	for i := range w.clients {
		if err = w.clients[i].Stop(); Check(err) {
		}
		if err = w.clients[i].Close(); Check(err) {
		}
	}
	for i := range w.workers {
		if err = w.workers[i].Interrupt(); !Check(err) {
		}
		if err = w.workers[i].Kill(); !Check(err) {
		}
		Debug("stopped worker", i)
	}
	w.active.Store(false)
}

func Handle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		Debug("miner controller starting")
		ctx, cancel := context.WithCancel(context.Background())
		w := &Worker{
			cx:            cx,
			ctx:           ctx,
			quit:          cx.KillAll,
			sendAddresses: []*net.UDPAddr{},
			StartChan:     make(chan struct{}),
			StopChan:      make(chan struct{}),
			SetThreads:    make(chan int),
			solutions:     make([]SolutionData, 0, 201),
		}
		w.lastSent.Store(time.Now().UnixNano())
		w.active.Store(false)
		Debug("opening broadcast channel listener")
		w.conn, err = transport.NewBroadcastChannel("kopachmain", w, *cx.Config.MinerPass,
			transport.DefaultPort, control.MaxDatagramSize, handlers,
			cx.KillAll)
		if err != nil {
			Error(err)
			cancel()
			return
		}
		// start up the workers
		if *cx.Config.Generate {
			w.Start(cx)
		}
		// controller watcher thread
		go func() {
			Debug("starting controller watcher")
			ticker := time.NewTicker(time.Second)
		out:
			for {
				select {
				case <-ticker.C:
					// if the last message sent was 3 seconds ago the server is almost certainly disconnected or crashed
					// so clear FirstSender
					since := time.Now().Sub(time.Unix(0, w.lastSent.Load()))
					wasSending := since > time.Second*3 && w.FirstSender.Load() != ""
					if wasSending {
						Debug("previous current controller has stopped broadcasting", since, w.FirstSender.Load())
						// when this string is clear other broadcasts will be listened to
						w.FirstSender.Store("")
						// pause the workers
						for i := range w.clients {
							Debug("sending pause to worker", i)
							err := w.clients[i].Pause()
							if err != nil {
								Error(err)
							}
						}
					}
				case <-w.StartChan:
					*cx.Config.Generate = true
					save.Pod(cx.Config)
					w.Start(cx)
				case <-w.StopChan:
					*cx.Config.Generate = false
					save.Pod(cx.Config)
					w.Stop()
				case n := <-w.SetThreads:
					*cx.Config.GenThreads = n
					save.Pod(cx.Config)
					if *cx.Config.Generate {
						// always sanitise
						if n < 0 {
							n = 0
						}
						if n > int(maxThreads) {
							n = int(maxThreads)
						}
						w.Stop()
						w.Start(cx)
					}
				case <-cx.KillAll:
					Debug("stopping from killall")
					// close(w.quit)
					break out
				case <-w.quit:
					Debug("stopping from quit")
					break out
				}
			}
		}()
		Debug("listening on", control.UDP4MulticastAddress)
		if *cx.Config.KopachGUI {
			Info("opening miner controller GUI")
			go Run(w, cx)
		}
		<-w.quit
		Info("kopach shutting down")
		return
	}
}

// these are the handlers for specific message types.
var handlers = transport.Handlers{
	string(job.Magic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		w := ctx.(*Worker)
		if !w.active.Load() {
			Debug("not active")
			return
		}
		j := job.LoadContainer(b)
		ips := j.GetIPs()
		w.height = j.GetNewHeight()
		cP := j.GetControllerListenerPort()
		addr := net.JoinHostPort(ips[0].String(), fmt.Sprint(cP))
		firstSender := w.FirstSender.Load()
		otherSent := firstSender != addr && firstSender != ""
		if otherSent {
			Trace("ignoring other controller job")
			// ignore other controllers while one is active and received first
			return
		}
		w.FirstSender.Store(addr)
		w.lastSent.Store(time.Now().UnixNano())
		for i := range w.clients {
			err := w.clients[i].NewJob(&j)
			if err != nil {
				Error(err)
			}
		}
		return
	},
	string(pause.PauseMagic): func(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
		w := ctx.(*Worker)
		p := pause.LoadPauseContainer(b)
		fs := w.FirstSender.Load()
		ni := p.GetIPs()[0].String()
		np := p.GetControllerListenerPort()
		ns := net.JoinHostPort(ni, fmt.Sprint(np))
		if fs == ns {
			for i := range w.clients {
				Debug("sending pause to worker", i, fs, ns)
				err := w.clients[i].Pause()
				if err != nil {
					Error(err)
				}
			}
		}
		return
	},
	string(sol.SolutionMagic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		w := ctx.(*Worker)
		portSlice := strings.Split(w.FirstSender.Load(), ":")
		if len(portSlice) < 2 {
			Debug("error with solution", w.FirstSender.Load(), portSlice)
			return
		}
		port := portSlice[1]
		j := sol.LoadSolContainer(b)
		senderPort := j.GetSenderPort()
		if fmt.Sprint(senderPort) == port {
			Warn("we found a solution")
			// prepend to list of solutions for GUI display if enabled
			if *w.cx.Config.KopachGUI {
				Debug("length solutions", len(w.solutions))
				w.solutions = append([]SolutionData{{time: time.Now(), height: int(w.height), block: j.GetMsgBlock()},
				}, w.solutions...)
				if len(w.solutions) > 200 {
					w.solutions = w.solutions[:200]
				}
				w.solutionCount = len(w.solutions)
			}
		}
		w.FirstSender.Store("")
		return
	},
}
