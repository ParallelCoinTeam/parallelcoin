package worker

import (
	"crypto/cipher"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
	
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/cmd/kopach/control"
	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/transport"
	"github.com/p9c/pod/pkg/data/ring"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const RoundsPerAlgo = 100

type Worker struct {
	mx               sync.Mutex
	id               string
	pipeConn         *stdconn.StdConn
	dispatchConn     *transport.Channel
	dispatchReady    atomic.Bool
	ciph             cipher.AEAD
	quit             qu.C
	templatesMessage *templates.Message
	uuid             atomic.Uint64
	roller           *Counter
	startNonce       uint32
	startChan        qu.C
	stopChan         qu.C
	running          atomic.Bool
	hashCount        atomic.Uint64
	hashSampleBuf    *ring.BufferUint64
}

type Counter struct {
	rpa           int32
	C             atomic.Int32
	Algos         atomic.Value // []int32
	RoundsPerAlgo atomic.Int32
}

// NewCounter returns an initialized algorithm rolling counter that ensures each
// miner does equal amounts of every algorithm
func NewCounter(roundsPerAlgo int32) (c *Counter) {
	// these will be populated when work arrives
	var algos []int32
	// Start the counter at a random position
	rand.Seed(time.Now().UnixNano())
	c = &Counter{}
	c.C.Store(int32(rand.Intn(int(roundsPerAlgo)+1) + 1))
	c.Algos.Store(algos)
	c.RoundsPerAlgo.Store(roundsPerAlgo)
	c.rpa = roundsPerAlgo
	return
}

// GetAlgoVer returns the next algo version based on the current configuration
func (c *Counter) GetAlgoVer(height int32) (ver int32) {
	// the formula below rolls through versions with blocks roundsPerAlgo long for each algorithm by its index
	algs := fork.GetAlgoVerSlice(height)
	// Debug(algs)
	if c.RoundsPerAlgo.Load() < 1 {
		Debug("RoundsPerAlgo is", c.RoundsPerAlgo.Load(), len(algs))
		return 0
	}
	if len(algs) > 0 {
		ver = algs[(c.C.Load()/
			c.RoundsPerAlgo.Load())%
			int32(len(algs))]
		c.C.Add(1)
	}
	return
}
//
// func (w *Worker) hashReport() {
// 	w.hashSampleBuf.Add(w.hashCount.Load())
// 	av := ewma.NewMovingAverage(15)
// 	var i int
// 	var prev uint64
// 	if err := w.hashSampleBuf.ForEach(
// 		func(v uint64) error {
// 			if i < 1 {
// 				prev = v
// 			} else {
// 				interval := v - prev
// 				av.Add(float64(interval))
// 				prev = v
// 			}
// 			i++
// 			return nil
// 		},
// 	); Check(err) {
// 	}
// 	// Info("kopach",w.hashSampleBuf.Cursor, w.hashSampleBuf.Buf)
// 	Tracef("average hashrate %.2f", av.Value())
// }

// NewWithConnAndSemaphore is exposed to enable use an actual network connection while retaining the same RPC API to
// allow a worker to be configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(id string, conn *stdconn.StdConn, quit qu.C, uuid uint64) *Worker {
	Debug("creating new worker")
	// msgBlock := wire.MsgBlock{Header: wire.BlockHeader{}}
	w := &Worker{
		id:            id,
		pipeConn:      conn,
		quit:          quit,
		roller:        NewCounter(RoundsPerAlgo),
		startChan:     qu.T(),
		stopChan:      qu.T(),
		hashSampleBuf: ring.NewBufferUint64(1000),
	}
	w.uuid.Store(uuid)
	w.dispatchReady.Store(false)
	// with this we can report cumulative hash counts as well as using it to distribute algorithms evenly
	w.startNonce = uint32(w.roller.C.Load())
	interrupt.AddHandler(
		func() {
			Debug("worker quitting")
			w.stopChan <- struct{}{}
			// _ = w.pipeConn.Close()
			w.dispatchReady.Store(false)
		},
	)
	go worker(w)
	return w
}

func worker(w *Worker) {
	Debug("main work loop starting")
	// sampleTicker := time.NewTicker(time.Second)
	var nonce uint32
out:
	for {
		// Pause state
		Trace("worker pausing")
	pausing:
		for {
			select {
			// case <-sampleTicker.C:
			// 	// w.hashReport()
			// 	break
			case <-w.stopChan.Wait():
				Debug("received pause signal while paused")
				// drain stop channel in pause
				break
			case <-w.startChan.Wait():
				Debug("received start signal")
				break pausing
			case <-w.quit.Wait():
				Debug("quitting")
				break out
			}
		}
		// Run state
		Trace("worker running")
	running:
		for {
			select {
			// case <-sampleTicker.C:
			// 	// w.hashReport()
			// 	break
			case <-w.startChan.Wait():
				Debug("received start signal while running")
				// drain start channel in run mode
				break
			case <-w.stopChan.Wait():
				Debug("received pause signal while running")
				break running
			case <-w.quit.Wait():
				Debug("worker stopping while running")
				break out
			default:
				if w.templatesMessage == nil || !w.dispatchReady.Load() {
					Debug("not ready to work")
				} else {
					// Debug("starting mining round")
					newHeight := w.templatesMessage.Height
					vers := w.roller.GetAlgoVer(newHeight)
					nonce++
					tn := time.Now().Round(time.Second)
					if tn.After(w.templatesMessage.Timestamp.Round(time.Second)) {
						w.templatesMessage.Timestamp = tn
					}
					if w.roller.C.Load()%w.roller.RoundsPerAlgo.Load() == 0 {
						// Debug("switching algorithms", w.roller.C.Load())
						// send out broadcast containing worker nonce and algorithm and count of blocks
						w.hashCount.Store(w.hashCount.Load() + uint64(w.roller.RoundsPerAlgo.Load()))
						hashReport := hashrate.Get(w.roller.RoundsPerAlgo.Load(), vers, newHeight, w.id)
						err := w.dispatchConn.SendMany(
							hashrate.Magic,
							transport.GetShards(hashReport),
						)
						if err != nil {
							Error(err)
						}
						// reseed the nonce
						rand.Seed(time.Now().UnixNano())
						nonce = rand.Uint32()
						select {
						case <-w.quit.Wait():
							Debug("breaking out of work loop")
							break out
						case <-w.stopChan.Wait():
							Debug("received pause signal while running")
							break running
						default:
						}
					}
					blockHeader := w.templatesMessage.GenBlockHeader(vers)
					blockHeader.Nonce = nonce
					// Debugs(w.templatesMessage)
					// Debugs(blockHeader)
					hash := blockHeader.BlockHashWithAlgos(newHeight)
					bigHash := blockchain.HashToBig(&hash)
					if bigHash.Cmp(fork.CompactToBig(blockHeader.Bits)) <= 0 {
						Debug("found solution", newHeight)
						srs := sol.Encode(w.uuid.Load(), blockHeader)
						err := w.dispatchConn.SendMany(
							sol.Magic,
							transport.GetShards(srs),
						)
						if err != nil {
							Error(err)
						}
						Debug("sent solution")
						w.templatesMessage = nil
						select {
						case <-w.quit.Wait():
							Debug("breaking out of work loop")
							break out
						default:
						}
						break running
					}
					// Debug("completed mining round")
				}
			}
		}
	}
	Debug("worker finished")
	interrupt.Request()
}

// New initialises the state for a worker, loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(id string, quit qu.C, uuid uint64) (w *Worker, conn net.Conn) {
	// log.L.SetLevel("trace", true)
	sc := stdconn.New(os.Stdin, os.Stdout, quit)
	
	return NewWithConnAndSemaphore(id, sc, quit, uuid), sc
}

// NewJob is a delivery of a new job for the worker, this makes the miner start
// mining from pause or pause, prepare the work and restart
func (w *Worker) NewJob(j *templates.Message, reply *bool) (err error) {
	// Trace("received new job")
	if !w.dispatchReady.Load() {
		Debug("dispatch not ready")
		*reply = true
		return
	}
	if w.templatesMessage != nil {
		if j.PrevBlock == w.templatesMessage.PrevBlock {
			// Trace("not a new job")
			*reply = true
			return
		}
	}
	// Debugs(j)
	*reply = true
	Debug("halting current work")
	w.stopChan <- struct{}{}
	// load the job into the template
	if w.templatesMessage == nil {
		w.templatesMessage = j
	} else {
		*w.templatesMessage = *j
	}
	Debug("switching to new job")
	w.startChan <- struct{}{}
	return
}

// Pause signals the worker to stop working, releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	Trace("pausing from IPC")
	w.running.Store(false)
	w.stopChan <- struct{}{}
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	Debug("stopping from IPC")
	w.stopChan <- struct{}{}
	defer w.quit.Q()
	*reply = true
	// time.Sleep(time.Second * 3)
	// os.Exit(0)
	return
}

// SendPass gives the encryption key configured in the kopach controller ( pod) configuration to allow workers to
// dispatch their solutions
func (w *Worker) SendPass(pass string, reply *bool) (err error) {
	Debug("receiving dispatch password", pass)
	rand.Seed(time.Now().UnixNano())
	// sp := fmt.Sprint(rand.Intn(32767) + 1025)
	// rp := fmt.Sprint(rand.Intn(32767) + 1025)
	var conn *transport.Channel
	conn, err = transport.NewBroadcastChannel(
		"kopachworker",
		w,
		pass,
		transport.DefaultPort,
		control.MaxDatagramSize,
		transport.Handlers{},
		w.quit,
	)
	if err != nil {
		Error(err)
	}
	w.dispatchConn = conn
	w.dispatchReady.Store(true)
	*reply = true
	return
}
