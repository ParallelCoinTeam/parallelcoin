package worker

import (
	"crypto/cipher"
	"errors"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/VividCortex/ewma"
	"go.uber.org/atomic"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/kopachctrl"
	"github.com/p9c/pod/pkg/kopachctrl/hashrate"
	"github.com/p9c/pod/pkg/kopachctrl/job"
	"github.com/p9c/pod/pkg/kopachctrl/sol"
	"github.com/p9c/pod/pkg/ring"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const RoundsPerAlgo = 69

type Worker struct {
	mx            sync.Mutex
	pipeConn      *stdconn.StdConn
	multicastConn net.Conn
	unicastConn   net.Conn
	dispatchConn  *transport.Channel
	dispatchReady atomic.Bool
	ciph          cipher.AEAD
	Quit          chan struct{}
	run           sem.T
	block         atomic.Value
	senderPort    atomic.Uint32
	msgBlock      atomic.Value // *wire.MsgBlock
	bitses        atomic.Value
	hashes        atomic.Value
	lastMerkle    *chainhash.Hash
	roller        *Counter
	startNonce    uint32
	startChan     chan struct{}
	stopChan      chan struct{}
	running       atomic.Bool
	hashCount     atomic.Uint64
	hashSampleBuf *ring.BufferUint64
}

type Counter struct {
	rpa           int32
	C             atomic.Int32
	Algos         atomic.Value // []int32
	RoundsPerAlgo atomic.Int32
}

// NewCounter returns an initialized algorithm rolling counter that ensures
// each miner does equal amounts of every algorithm
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
func (c *Counter) GetAlgoVer() (ver int32) {
	// the formula below rolls through versions with blocks roundsPerAlgo
	// long for each algorithm by its index
	algs := c.Algos.Load().([]int32)
	// L.Debug(algs)
	if c.RoundsPerAlgo.Load() < 1 {
		L.Debug("RoundsPerAlgo is", c.RoundsPerAlgo.Load(), len(algs))
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

func (w *Worker) hashReport() {
	w.hashSampleBuf.Add(w.hashCount.Load())
	av := ewma.NewMovingAverage(15)
	var i int
	var prev uint64
	if err := w.hashSampleBuf.ForEach(func(v uint64) error {
		if i < 1 {
			prev = v
		} else {
			interval := v - prev
			av.Add(float64(interval))
			prev = v
		}
		i++
		return nil
	}); L.Check(err) {
	}
	// L.Info("kopach",w.hashSampleBuf.Cursor, w.hashSampleBuf.Buf)
	L.Tracef("average hashrate %.2f", av.Value())
}

// NewWithConnAndSemaphore is exposed to enable use an actual network
// connection while retaining the same RPC API to allow a worker to be
// configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(conn *stdconn.StdConn, quit chan struct{}) *Worker {
	L.Debug("creating new worker")
	msgBlock := wire.MsgBlock{Header: wire.BlockHeader{}}
	w := &Worker{
		pipeConn:      conn,
		Quit:          quit,
		roller:        NewCounter(RoundsPerAlgo),
		startChan:     make(chan struct{}),
		stopChan:      make(chan struct{}),
		hashSampleBuf: ring.NewBufferUint64(1000),
	}
	w.msgBlock.Store(msgBlock)
	w.block.Store(util.NewBlock(&msgBlock))
	w.dispatchReady.Store(false)
	// with this we can report cumulative hash counts as well as using it to
	// distribute algorithms evenly
	// tn := time.Now()
	w.startNonce = uint32(w.roller.C.Load())
	interrupt.AddHandler(func() {
		L.Debug("worker quitting")
		close(w.Quit)
		// w.pipeConn.Close()
		w.dispatchReady.Store(false)
	})
	go func(w *Worker) {
		L.Debug("main work loop starting")
		sampleTicker := time.NewTicker(time.Second)
	out:
		for {
			// Pause state
			L.Debug("worker pausing")
		pausing:
			for {
				select {
				case <-sampleTicker.C:
					w.hashReport()
					break
				case <-w.stopChan:
					L.Trace("received pause signal while paused")
					// drain stop channel in pause
					break
				case <-w.startChan:
					L.Trace("received start signal")
					break pausing
				case <-w.Quit:
					L.Trace("quitting")
					break out
				}
			}
			// Run state
			L.Debug("worker running")
		running:
			for {
				select {
				case <-sampleTicker.C:
					w.hashReport()
					break
				case <-w.startChan:
					L.Trace("received start signal while running")
					// drain start channel in run mode
					break
				case <-w.stopChan:
					L.Trace("received pause signal while running")
					// w.block.Store(&util.Block{})
					// w.bitses.Store((blockchain.TargetBits)(nil))
					// w.hashes.Store((map[int32]*chainhash.Hash)(nil))
					break running
				case <-w.Quit:
					L.Trace("worker stopping while running")
					break out
				default:
					if w.block.Load() == nil || w.bitses.Load() == nil ||
						w.hashes.Load() == nil || !w.dispatchReady.Load() {
						// L.Info("stop was called before we started working")
					} else {
						// work
						nH := w.block.Load().(*util.Block).Height()
						hv := w.roller.GetAlgoVer()
						// L.Debug(hv)
						h := w.hashes.Load().(map[int32]*chainhash.Hash)
						// L.Debug("hashes", hv, h)
						mmb := w.msgBlock.Load().(wire.MsgBlock)
						mb := &mmb
						mb.Header.Version = hv
						if h != nil {
							mr, ok := h[hv]
							if !ok {
								continue
							}
							mb.Header.MerkleRoot = *mr
						} else {
							continue
						}
						b := w.bitses.Load().(blockchain.TargetBits)
						if bb, ok := b[mb.Header.Version]; ok {
							mb.Header.Bits = bb
						} else {
							continue
						}
						mb.Header.Timestamp = time.Now()
						var nextAlgo int32
						if w.roller.C.Load()%w.roller.RoundsPerAlgo.Load() == 0 {
							select {
							case <-w.Quit:
								L.Trace("worker stopping on pausing message")
								break out
							default:
							}
							// send out broadcast containing worker nonce and algorithm and count of blocks
							w.hashCount.Store(w.hashCount.Load() + uint64(w.roller.RoundsPerAlgo.Load()))
							nextAlgo = w.roller.C.Load() + 1
							hashReport := hashrate.Get(w.roller.RoundsPerAlgo.Load(), nextAlgo, nH)
							err := w.dispatchConn.SendMany(hashrate.HashrateMagic,
								transport.GetShards(hashReport.Data))
							if err != nil {
								L.Error(err)
							}
						}
						hash := mb.Header.BlockHashWithAlgos(nH)
						bigHash := blockchain.HashToBig(&hash)
						if bigHash.Cmp(fork.CompactToBig(mb.Header.Bits)) <= 0 {
							srs := sol.GetSolContainer(w.senderPort.Load(), mb)
							err := w.dispatchConn.SendMany(sol.SolutionMagic,
								transport.GetShards(srs.Data))
							if err != nil {
								L.Error(err)
							}
							L.Trace("sent solution")

							break running
						}
						mb.Header.Version = nextAlgo
						mb.Header.Bits = w.bitses.Load().(blockchain.TargetBits)[mb.Header.Version]
						mb.Header.Nonce++
						w.msgBlock.Store(*mb)
					}
				}
			}
		}
		L.Debug("worker finished")
	}(w)
	return w
}

// New initialises the state for a worker, loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(quit chan struct{}) (w *Worker, conn net.Conn) {
	// L.SetLevel("trace", true)
	sc := stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(&sc, quit), &sc
}

// NewJob is a delivery of a new job for the worker, this makes the miner start
// mining from pause or pause, prepare the work and restart
func (w *Worker) NewJob(job *job.Container, reply *bool) (err error) {
	L.Trace("starting new job")
	if !w.dispatchReady.Load() { // || !w.running.Load() {
		*reply = true
		return
	}
	j := job.Struct()
	w.bitses.Store(j.Bitses)
	w.hashes.Store(j.Hashes)
	if j.Hashes[5].IsEqual(w.lastMerkle) {
		// L.Debug("not a new job")
		*reply = true
		return
	}
	var algos []int32
	for i := range j.Bitses {
		// we don't need to know net params if version numbers come with jobs
		algos = append(algos, i)
	}
	w.lastMerkle = j.Hashes[5]
	*reply = true
	// halting current work
	w.stopChan <- struct{}{}
	newHeight := job.GetNewHeight()

	if len(algos) > 0 {
		// if we didn't get them in the job don't update the old
		w.roller.Algos.Store(algos)
	}
	mbb := w.msgBlock.Load().(wire.MsgBlock)
	mb := &mbb
	mb.Header.PrevBlock = *job.GetPrevBlockHash()
	// TODO: ensure worker time sync - ntp? time wrapper with skew adjustment
	hv := w.roller.GetAlgoVer()
	mb.Header.Version = hv
	var ok bool
	mb.Header.Bits, ok = j.Bitses[mb.Header.Version]
	if !ok {
		return errors.New("bits are empty")
	}
	rand.Seed(time.Now().UnixNano())
	mb.Header.Nonce = rand.Uint32()
	if j.Hashes == nil {
		return errors.New("failed to decode merkle roots")
	} else {
		hh, ok := j.Hashes[hv]
		if !ok {
			return errors.New("could not get merkle root from job")
		}
		mb.Header.MerkleRoot = *hh
	}
	mb.Header.Timestamp = time.Now()
	// make the work select block start running
	bb := util.NewBlock(mb)
	bb.SetHeight(newHeight)
	w.block.Store(bb)
	w.msgBlock.Store(*mb)
	w.senderPort.Store(uint32(job.GetControllerListenerPort()))
	// halting current work
	// w.stopChan <- struct{}{}
	w.startChan <- struct{}{}
	return
}

// Pause signals the worker to stop working,
// releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	L.Trace("pausing from IPC")
	w.running.Store(false)
	w.stopChan <- struct{}{}
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	L.Debug("stopping from IPC")
	w.stopChan <- struct{}{}
	defer close(w.Quit)
	*reply = true
	return
}

// SendPass gives the encryption key configured in the kopach controller (
// pod) configuration to allow workers to dispatch their solutions
func (w *Worker) SendPass(pass string, reply *bool) (err error) {
	L.Debug("receiving dispatch password", pass)
	rand.Seed(time.Now().UnixNano())
	// sp := fmt.Sprint(rand.Intn(32767) + 1025)
	// rp := fmt.Sprint(rand.Intn(32767) + 1025)
	var conn *transport.Channel
	conn, err = transport.NewBroadcastChannel("kopachworker", w, pass,
		transport.DefaultPort, kopachctrl.MaxDatagramSize, transport.Handlers{}, w.Quit)
	if err != nil {
		L.Error(err)
	}
	w.dispatchConn = conn
	w.dispatchReady.Store(true)
	*reply = true
	return
}
