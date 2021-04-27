package worker

import (
	"crypto/cipher"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/chainrpc/templates"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/matrjoska/pkg/pipe"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chainrpc/hashrate"
	"github.com/p9c/matrjoska/pkg/chainrpc/sol"

	"go.uber.org/atomic"

	"github.com/p9c/interrupt"

	"github.com/p9c/matrjoska/pkg/ring"
	"github.com/p9c/matrjoska/pkg/transport"
)

const CountPerRound = 81

type Worker struct {
	mx               sync.Mutex
	id               string
	pipeConn         *pipe.StdConn
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
func NewCounter(countPerRound int32) (c *Counter) {
	// these will be populated when work arrives
	var algos []int32
	// Start the counter at a random position
	rand.Seed(time.Now().UnixNano())
	c = &Counter{}
	c.C.Store(int32(rand.Intn(int(countPerRound)+1) + 1))
	c.Algos.Store(algos)
	c.RoundsPerAlgo.Store(countPerRound)
	c.rpa = countPerRound
	return
}

// GetAlgoVer returns the next algo version based on the current configuration
func (c *Counter) GetAlgoVer(height int32) (ver int32) {
	// the formula below rolls through versions with blocks roundsPerAlgo long for each algorithm by its index
	algs := fork.GetAlgoVerSlice(height)
	// D.Ln(algs)
	if c.RoundsPerAlgo.Load() < 1 {
		D.Ln("CountPerRound is", c.RoundsPerAlgo.Load(), len(algs))
		return 0
	}
	if len(algs) > 0 {
		ver = algs[c.C.Load()%int32(len(algs))]
		// ver = algs[(c.C.Load()/
		// 	c.CountPerRound.Load())%
		// 	int32(len(algs))]
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
// 	if e := w.hashSampleBuf.ForEach(
// 		func(v uint64) (e error) {
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
// 	); E.Chk(e) {
// 	}
// 	// I.Ln("kopach",w.hashSampleBuf.Cursor, w.hashSampleBuf.Buf)
// 	Tracef("average hashrate %.2f", av.Value())
// }

// NewWithConnAndSemaphore is exposed to enable use an actual network connection while retaining the same RPC API to
// allow a worker to be configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(id string, conn *pipe.StdConn, quit qu.C, uuid uint64) *Worker {
	T.Ln("creating new worker")
	// msgBlock := wire.WireBlock{Header: wire.BlockHeader{}}
	w := &Worker{
		id:            id,
		pipeConn:      conn,
		quit:          quit,
		roller:        NewCounter(CountPerRound),
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
			D.Ln("worker", id, "quitting")
			w.stopChan <- struct{}{}
			// _ = w.pipeConn.Close()
			w.dispatchReady.Store(false)
		},
	)
	go worker(w)
	return w
}

func worker(w *Worker) {
	D.Ln("main work loop starting")
	// sampleTicker := time.NewTicker(time.Second)
	var nonce uint32
out:
	for {
		// Pause state
		T.Ln("worker pausing")
	pausing:
		for {
			select {
			// case <-sampleTicker.C:
			// 	// w.hashReport()
			// 	break
			case <-w.stopChan.Wait():
				D.Ln("received pause signal while paused")
				// drain stop channel in pause
				break
			case <-w.startChan.Wait():
				D.Ln("received start signal")
				break pausing
			case <-w.quit.Wait():
				D.Ln("quitting")
				break out
			}
		}
		// Run state
		T.Ln("worker running")
	running:
		for {
			select {
			// case <-sampleTicker.C:
			// 	// w.hashReport()
			// 	break
			case <-w.startChan.Wait():
				D.Ln("received start signal while running")
				// drain start channel in run mode
				break
			case <-w.stopChan.Wait():
				D.Ln("received pause signal while running")
				break running
			case <-w.quit.Wait():
				D.Ln("worker stopping while running")
				break out
			default:
				if w.templatesMessage == nil || !w.dispatchReady.Load() {
					D.Ln("not ready to work")
				} else {
					// I.Ln("starting mining round")
					newHeight := w.templatesMessage.Height
					vers := w.roller.GetAlgoVer(newHeight)
					nonce++
					tn := time.Now().Round(time.Second)
					if tn.After(w.templatesMessage.Timestamp.Round(time.Second)) {
						w.templatesMessage.Timestamp = tn
					}
					if w.roller.C.Load()%w.roller.RoundsPerAlgo.Load() == 0 {
						D.Ln("switching algorithms", w.roller.C.Load())
						// send out broadcast containing worker nonce and algorithm and count of blocks
						w.hashCount.Store(w.hashCount.Load() + uint64(w.roller.RoundsPerAlgo.Load()))
						hashReport := hashrate.Get(w.roller.RoundsPerAlgo.Load(), vers, newHeight, w.id)
						e := w.dispatchConn.SendMany(
							hashrate.Magic,
							transport.GetShards(hashReport),
						)
						if e != nil {
						}
						// reseed the nonce
						rand.Seed(time.Now().UnixNano())
						nonce = rand.Uint32()
						select {
						case <-w.quit.Wait():
							D.Ln("breaking out of work loop")
							break out
						case <-w.stopChan.Wait():
							D.Ln("received pause signal while running")
							break running
						default:
						}
					}
					blockHeader := w.templatesMessage.GenBlockHeader(vers)
					blockHeader.Nonce = nonce
					// D.S(w.templatesMessage)
					// D.S(blockHeader)
					hash := blockHeader.BlockHashWithAlgos(newHeight)
					bigHash := blockchain.HashToBig(&hash)
					if bigHash.Cmp(bits.CompactToBig(blockHeader.Bits)) <= 0 {
						D.Ln("found solution", newHeight, w.templatesMessage.Nonce, w.templatesMessage.UUID)
						srs := sol.Encode(w.templatesMessage.Nonce, w.templatesMessage.UUID, blockHeader)
						e := w.dispatchConn.SendMany(
							sol.Magic,
							transport.GetShards(srs),
						)
						if e != nil {
						}
						D.Ln("sent solution")
						w.templatesMessage = nil
						select {
						case <-w.quit.Wait():
							D.Ln("breaking out of work loop")
							break out
						default:
						}
						break running
					}
					// D.Ln("completed mining round")
				}
			}
		}
	}
	D.Ln("worker finished")
	interrupt.Request()
}

// New initialises the state for a worker, loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(id string, quit qu.C, uuid uint64) (w *Worker, conn net.Conn) {
	// log.L.SetLevel("trace", true)
	sc := pipe.New(os.Stdin, os.Stdout, quit)
	
	return NewWithConnAndSemaphore(id, sc, quit, uuid), sc
}

// NewJob is a delivery of a new job for the worker, this makes the miner start
// mining from pause or pause, prepare the work and restart
func (w *Worker) NewJob(j *templates.Message, reply *bool) (e error) {
	// T.Ln("received new job")
	if !w.dispatchReady.Load() {
		D.Ln("dispatch not ready")
		*reply = true
		return
	}
	if w.templatesMessage != nil {
		if j.PrevBlock == w.templatesMessage.PrevBlock {
			// T.Ln("not a new job")
			*reply = true
			return
		}
	}
	// D.S(j)
	*reply = true
	D.Ln("halting current work")
	w.stopChan <- struct{}{}
	D.Ln("halt signal sent")
	// load the job into the template
	if w.templatesMessage == nil {
		w.templatesMessage = j
	} else {
		*w.templatesMessage = *j
	}
	D.Ln("switching to new job")
	w.startChan <- struct{}{}
	D.Ln("start signal sent")
	return
}

// Pause signals the worker to stop working, releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (e error) {
	T.Ln("pausing from IPC")
	w.running.Store(false)
	w.stopChan <- struct{}{}
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (e error) {
	D.Ln("stopping from IPC")
	w.stopChan <- struct{}{}
	defer w.quit.Q()
	*reply = true
	// time.Sleep(time.Second * 3)
	// os.Exit(0)
	return
}

// SendPass gives the encryption key configured in the kopach controller ( pod) configuration to allow workers to
// dispatch their solutions
func (w *Worker) SendPass(pass []byte, reply *bool) (e error) {
	D.Ln("receiving dispatch password", pass)
	rand.Seed(time.Now().UnixNano())
	// sp := fmt.Sprint(rand.Intn(32767) + 1025)
	// rp := fmt.Sprint(rand.Intn(32767) + 1025)
	var conn *transport.Channel
	conn, e = transport.NewBroadcastChannel(
		"kopachworker",
		w,
		pass,
		transport.DefaultPort,
		constant.MaxDatagramSize,
		transport.Handlers{},
		w.quit,
	)
	if e != nil {
	}
	w.dispatchConn = conn
	w.dispatchReady.Store(true)
	*reply = true
	return
}
