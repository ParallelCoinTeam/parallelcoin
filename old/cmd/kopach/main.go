package kopach

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/p9c/pod/pkg/control/sol"
	"github.com/p9c/pod/pkg/control/templates"
	"github.com/p9c/pod/pkg/logg"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/podcfg"
	"net"
	"os"
	"runtime"
	"time"
	
	"github.com/p9c/pod/pkg/control/p2padvt"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/VividCortex/ewma"
	"github.com/urfave/cli"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/cmd/kopach/client"
	"github.com/p9c/pod/pkg/chainhash"
	"github.com/p9c/pod/pkg/control"
	"github.com/p9c/pod/pkg/control/hashrate"
	"github.com/p9c/pod/pkg/control/job"
	"github.com/p9c/pod/pkg/control/pause"
	"github.com/p9c/pod/pkg/pipe/stdconn/worker"
	rav "github.com/p9c/pod/pkg/ring"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var maxThreads = float32(runtime.NumCPU())

type HashCount struct {
	uint64
	Time time.Time
}

type SolutionData struct {
	time       time.Time
	height     int
	algo       string
	hash       string
	indexHash  string
	version    int32
	prevBlock  string
	merkleRoot string
	timestamp  time.Time
	bits       uint32
	nonce      uint32
}

type Worker struct {
	id                  string
	cx                  *pod.State
	height              int32
	active              atomic.Bool
	conn                *transport.Channel
	ctx                 context.Context
	quit                qu.C
	sendAddresses       []*net.UDPAddr
	clients             []*client.Client
	workers             []*worker.Worker
	FirstSender         atomic.Uint64
	lastSent            atomic.Int64
	Status              atomic.String
	HashTick            chan HashCount
	LastHash            *chainhash.Hash
	StartChan, StopChan qu.C
	SetThreads          chan int
	PassChan            chan string
	solutions           []SolutionData
	solutionCount       int
	Update              qu.C
	hashCount           atomic.Uint64
	hashSampleBuf       *rav.BufferUint64
	hashrate            float64
	lastNonce           uint64
}

func (w *Worker) Start() {
	// if !*cx.Config.Generate {
	// 	D.Ln("called start but not running generate")
	// 	return
	// }
	D.Ln("starting up kopach workers")
	w.workers = []*worker.Worker{}
	w.clients = []*client.Client{}
	for i := 0; i < *w.cx.Config.GenThreads; i++ {
		D.Ln("starting worker", i)
		cmd, _ := worker.Spawn(w.quit, os.Args[0], "worker", w.id, w.cx.ActiveNet.Name, *w.cx.Config.LogLevel)
		w.workers = append(w.workers, cmd)
		w.clients = append(w.clients, client.New(cmd.StdConn))
	}
	for i := range w.clients {
		T.Ln("sending pass to worker", i)
		e := w.clients[i].SendPass(*w.cx.Config.MinerPass)
		if e != nil {
		}
	}
	D.Ln("setting workers to active")
	w.active.Store(true)
	
}

func (w *Worker) Stop() {
	var e error
	for i := range w.clients {
		if e = w.clients[i].Pause(); E.Chk(e) {
		}
		if e = w.clients[i].Stop(); E.Chk(e) {
		}
		if e = w.clients[i].Close(); E.Chk(e) {
		}
	}
	for i := range w.workers {
		// if e = w.workers[i].Interrupt(); !E.Chk(e) {
		// }
		if e = w.workers[i].Kill(); !E.Chk(e) {
		}
		D.Ln("stopped worker", i)
	}
	w.active.Store(false)
	w.quit.Q()
}

func Handle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		D.Ln("miner controller starting")
		randomBytes := make([]byte, 4)
		if _, e = rand.Read(randomBytes); E.Chk(e) {
		}
		w := &Worker{
			id:            fmt.Sprintf("%x", randomBytes),
			cx:            cx,
			quit:          cx.KillAll,
			sendAddresses: []*net.UDPAddr{},
			StartChan:     qu.T(),
			StopChan:      qu.T(),
			SetThreads:    make(chan int),
			solutions:     make([]SolutionData, 0, 2048),
			Update:        qu.T(),
			hashSampleBuf: rav.NewBufferUint64(1000),
		}
		w.lastSent.Store(time.Now().UnixNano())
		w.active.Store(false)
		D.Ln("opening broadcast channel listener")
		w.conn, e = transport.NewBroadcastChannel(
			"kopachmain", w, *cx.Config.MinerPass,
			transport.DefaultPort, control.MaxDatagramSize, handlers,
			w.quit,
		)
		if e != nil {
			return
		}
		// start up the workers
		if *cx.Config.Generate {
			w.Start()
			interrupt.AddHandler(
				func() {
					w.Stop()
				},
			)
		}
		// controller watcher thread
		go func() {
			D.Ln("starting controller watcher")
			ticker := time.NewTicker(time.Second)
			logger := time.NewTicker(time.Second)
		out:
			for {
				select {
				case <-ticker.C:
					// if the last message sent was 3 seconds ago the server is almost certainly disconnected or crashed
					// so clear FirstSender
					since := time.Now().Sub(time.Unix(0, w.lastSent.Load()))
					wasSending := since > time.Second*6 && w.FirstSender.Load() != 0
					if wasSending {
						D.Ln("previous current controller has stopped broadcasting", since, w.FirstSender.Load())
						// when this string is clear other broadcasts will be listened to
						w.FirstSender.Store(0)
						// pause the workers
						for i := range w.clients {
							D.Ln("sending pause to worker", i)
							e := w.clients[i].Pause()
							if e != nil {
							}
						}
					}
					if interrupt.Requested() {
						w.StopChan <- struct{}{}
						w.quit.Q()
					}
				case <-logger.C:
					w.hashrate = w.HashReport()
					if interrupt.Requested() {
						w.StopChan <- struct{}{}
						w.quit.Q()
					}
				case <-w.StartChan.Wait():
					D.Ln("received signal on StartChan")
					*cx.Config.Generate = true
					podcfg.Save(cx.Config)
					w.Start()
				case <-w.StopChan.Wait():
					D.Ln("received signal on StopChan")
					*cx.Config.Generate = false
					podcfg.Save(cx.Config)
					w.Stop()
				case s := <-w.PassChan:
					D.Ln("received signal on PassChan", s)
					*cx.Config.MinerPass = s
					podcfg.Save(cx.Config)
					w.Stop()
					w.Start()
				case n := <-w.SetThreads:
					D.Ln("received signal on SetThreads", n)
					*cx.Config.GenThreads = n
					podcfg.Save(cx.Config)
					if *cx.Config.Generate {
						// always sanitise
						if n < 0 {
							n = int(maxThreads)
						}
						if n > int(maxThreads) {
							n = int(maxThreads)
						}
						w.Stop()
						w.Start()
					}
				case <-w.quit.Wait():
					D.Ln("stopping from quit")
					interrupt.Request()
					break out
				}
			}
			D.Ln("finished kopach miner work loop")
			logg.LogChanDisabled.Store(true)
		}()
		D.Ln("listening on", control.UDP4MulticastAddress)
		<-w.quit
		I.Ln("kopach shutting down") // , interrupt.GoroutineDump())
		// <-interrupt.HandlersDone
		I.Ln("kopach finished shutdown")
		return
	}
}

// these are the handlers for specific message types.
var handlers = transport.Handlers{
	string(hashrate.Magic): func(ctx interface{}, src net.Addr, dst string, b []byte) (e error) {
		c := ctx.(*Worker)
		if !c.active.Load() {
			D.Ln("not active")
			return
		}
		var hr hashrate.Hashrate
		gotiny.Unmarshal(b, &hr)
		// if this is not one of our workers reports ignore it
		if hr.ID != c.id {
			return
		}
		count := hr.Count
		hc := c.hashCount.Load() + uint64(count)
		c.hashCount.Store(hc)
		return
	},
	string(job.Magic): func(
		ctx interface{}, src net.Addr, dst string,
		b []byte,
	) (e error) {
		w := ctx.(*Worker)
		if !w.active.Load() {
			T.Ln("not active")
			return
		}
		jr := templates.Message{}
		gotiny.Unmarshal(b, &jr)
		w.height = jr.Height
		cN := jr.UUID
		firstSender := w.FirstSender.Load()
		otherSent := firstSender != cN && firstSender != 0
		if otherSent {
			T.Ln("ignoring other controller job", jr.Nonce, jr.UUID)
			// ignore other controllers while one is active and received first
			return
		}
		// if jr.Nonce == w.lastNonce {
		// 	I.Ln("same job again, ignoring (NOT)")
		// 	// return
		// }
		// w.lastNonce = jr.Nonce
		// w.FirstSender.Store(cN)
		T.Ln("received job, starting workers on it", jr.Nonce, jr.UUID)
		w.lastSent.Store(time.Now().UnixNano())
		for i := range w.clients {
			if e = w.clients[i].NewJob(&jr); E.Chk(e) {
			}
		}
		return
	},
	string(pause.Magic): func(ctx interface{}, src net.Addr, dst string, b []byte) (e error) {
		w := ctx.(*Worker)
		var advt p2padvt.Advertisment
		gotiny.Unmarshal(b, &advt)
		// p := pause.LoadPauseContainer(b)
		fs := w.FirstSender.Load()
		ni := advt.IPs
		// ni := p.GetIPs()[0].String()
		np := advt.UUID
		// np := p.GetControllerListenerPort()
		// ns := net.JoinHostPort(strings.Split(ni.String(), ":")[0], fmt.Sprint(np))
		D.Ln("received pause from server at", ni, np, "stopping", len(w.clients), "workers stopping")
		if fs == np {
			for i := range w.clients {
				// D.Ln("sending pause to worker", i, fs, np)
				e := w.clients[i].Pause()
				if e != nil {
				}
			}
		}
		w.FirstSender.Store(0)
		return
	},
	string(sol.Magic): func(
		ctx interface{}, src net.Addr, dst string,
		b []byte,
	) (e error) {
		// w := ctx.(*Worker)
		// I.Ln("shuffling work due to solution on network")
		// w.FirstSender.Store(0)
		// 	D.Ln("solution detected from miner at", src)
		// 	portSlice := strings.Split(w.FirstSender.Load(), ":")
		// 	if len(portSlice) < 2 {
		// 		D.Ln("error with solution", w.FirstSender.Load(), portSlice)
		// 		return
		// 	}
		// 	// port := portSlice[1]
		// 	// j := sol.LoadSolContainer(b)
		// 	// senderPort := j.GetSenderPort()
		// 	// if fmt.Sprint(senderPort) == port {
		// 	// // W.Ln("we found a solution")
		// 	// // prepend to list of solutions for GUI display if enabled
		// 	// if *w.cx.Config.KopachGUI {
		// 	// 	// D.Ln("length solutions", len(w.solutions))
		// 	// 	blok := j.GetMsgBlock()
		// 	// 	w.solutions = append(
		// 	// 		w.solutions, []SolutionData{
		// 	// 			{
		// 	// 				time:   time.Now(),
		// 	// 				height: int(w.height),
		// 	// 				algo: fmt.Sprint(
		// 	// 					fork.GetAlgoName(blok.Header.Version, w.height),
		// 	// 				),
		// 	// 				hash:       blok.Header.BlockHashWithAlgos(w.height).String(),
		// 	// 				indexHash:  blok.Header.BlockHash().String(),
		// 	// 				version:    blok.Header.Version,
		// 	// 				prevBlock:  blok.Header.PrevBlock.String(),
		// 	// 				merkleRoot: blok.Header.MerkleRoot.String(),
		// 	// 				timestamp:  blok.Header.Timestamp,
		// 	// 				bits:       blok.Header.Bits,
		// 	// 				nonce:      blok.Header.Nonce,
		// 	// 			},
		// 	// 		}...,
		// 	// 	)
		// 	// 	if len(w.solutions) > 2047 {
		// 	// 		w.solutions = w.solutions[len(w.solutions)-2047:]
		// 	// 	}
		// 	// 	w.solutionCount = len(w.solutions)
		// 	// 	w.Update <- struct{}{}
		// 	// }
		// 	// }
		// 	// D.Ln("no longer listening to", w.FirstSender.Load())
		// 	// w.FirstSender.Store("")
		return
	},
}

func (w *Worker) HashReport() float64 {
	// T.Ln("generating hash report")
	w.hashSampleBuf.Add(w.hashCount.Load())
	av := ewma.NewMovingAverage()
	var i int
	var prev uint64
	if e := w.hashSampleBuf.ForEach(
		func(v uint64) (e error) {
			if i < 1 {
				prev = v
			} else {
				interval := v - prev
				av.Add(float64(interval))
				prev = v
			}
			i++
			return nil
		},
	); E.Chk(e) {
	}
	average := av.Value()
	D.Ln("hashrate average", average)
	return average
}
