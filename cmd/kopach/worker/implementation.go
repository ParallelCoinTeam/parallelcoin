package worker

import (
	"crypto/cipher"
	"errors"
	"fmt"
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
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/controller/hashrate"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/controller/sol"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/ring"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
)

const RoundsPerAlgo = 25

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
	hashCount     atomic.Uint64
	hashSampleBuf *ring.BufferUint64
}

type Counter struct {
	C             atomic.Int32
	Algos         atomic.Value // []int32
	RoundsPerAlgo atomic.Int32
}

// func (w *Worker) Close() {
// 	if err := w.multicastConn.Close(); log.Check(err) {
// 	}
// 	if err := w.unicastConn.Close(); log.Check(err) {
// 	}
// 	if err := w.dispatchConn.Close(); log.Check(err) {
// 	}
//
// }

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
	return
}

// GetAlgoVer returns the next algo version based on the current configuration
func (c *Counter) GetAlgoVer() (ver int32) {
	// the formula below rolls through versions with blocks roundsPerAlgo
	// long for each algorithm by its index
	if c.RoundsPerAlgo.Load() < 1 || len(c.Algos.Load().([]int32)) < 1 {
		log.DEBUG("RoundsPerAlgo is", c.RoundsPerAlgo.Load(), len(c.Algos.Load().([]int32)))
	}
	ver = c.Algos.Load().([]int32)[(c.C.Load()/
		c.RoundsPerAlgo.Load())%
		int32(len(c.Algos.Load().([]int32)))]
	c.C.Add(1)
	return
}

func (w *Worker) hashReport() {
	w.hashSampleBuf.Add(w.hashCount.Load())
	av := ewma.NewMovingAverage()
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
	}); log.Check(err) {
	}
	// log.INFO("kopach",w.hashSampleBuf.Cursor, w.hashSampleBuf.Buf)
	log.INFOF("average hashrate %.2f", av.Value())
}

// NewWithConnAndSemaphore is exposed to enable use an actual network
// connection while retaining the same RPC API to allow a worker to be
// configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(conn *stdconn.StdConn, quit chan struct{}, ) *Worker {
	log.DEBUG("creating new worker")
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
		log.DEBUG("worker quitting")
		close(w.Quit)
		// w.pipeConn.Close()
		w.dispatchReady.Store(false)
	})
	go func(w *Worker) {
		log.DEBUG("main work loop starting")
		sampleTicker := time.NewTicker(time.Second)
	pausing:
		for {
			// Pause state
			select {
			case <-sampleTicker.C:
				w.hashReport()
			case <-w.stopChan:
				// drain stop channel in pause
				continue
			case <-w.startChan:
				break
			case <-w.Quit:
				log.DEBUG("worker stopping on pausing message")
				break pausing
			}
			log.TRACE("worker running")
			// Run state
		running:
			for {
				select {
				case <-sampleTicker.C:
					w.hashReport()
				case <-w.startChan:
					// drain start channel in run mode
					continue
				case <-w.stopChan:
					w.block.Store(&util.Block{})
					w.bitses.Store((map[int32]uint32)(nil))
					w.hashes.Store((map[int32]*chainhash.Hash)(nil))
					break running
				case <-w.Quit:
					log.DEBUG("worker stopping while running")
					break pausing
				default:
					if w.block.Load() == nil || w.bitses.Load() == nil || w.hashes.Load() == nil ||
						!w.dispatchReady.Load() {
						// log.INFO("stop was called before we started working")
					} else {
						// work
						nH := w.block.Load().(*util.Block).Height()
						hv := w.roller.GetAlgoVer()
						mmb := w.msgBlock.Load().(wire.MsgBlock)
						mb := &mmb
						mb.Header.Version = hv
						h := w.hashes.Load().(map[int32]*chainhash.Hash)
						if h != nil {
							mb.Header.MerkleRoot = *h[mb.Header.Version]
						} else {
							continue
						}
						b := w.bitses.Load().(map[int32]uint32)
						if bb, ok := b[mb.Header.Version]; ok {
							mb.Header.Bits = bb
						} else {
							continue
						}
						select {
						case <-w.stopChan:
							w.block.Store(&util.Block{})
							w.bitses.Store((map[int32]uint32)(nil))
							w.hashes.Store((map[int32]*chainhash.Hash)(nil))
							break running
						case <-w.Quit:
							log.DEBUG("worker stopping in the middle of it")
							break pausing
						default:
						}
						var nextAlgo int32
						if w.roller.C.Load()%w.roller.RoundsPerAlgo.Load() == 0 {
							select {
							case <-w.Quit:
								log.DEBUG("worker stopping on pausing message")
								break pausing
							default:
							}
							// log.DEBUG("sending hashcount")
							// send out broadcast containing worker nonce and algorithm and count of blocks
							w.hashCount.Store(w.hashCount.Load() + uint64(w.roller.RoundsPerAlgo.Load()))
							nextAlgo = w.roller.C.Load() + 1
							hashReport := hashrate.Get(w.roller.RoundsPerAlgo.Load(), nextAlgo, nH)
							err := w.dispatchConn.SendMany(hashrate.HashrateMagic,
								transport.GetShards(hashReport.Data))
							if err != nil {
								log.ERROR(err)
							}
						}
						hash := mb.Header.BlockHashWithAlgos(nH)
						bigHash := blockchain.HashToBig(&hash)
						if bigHash.Cmp(fork.CompactToBig(mb.Header.Bits)) <= 0 {
							log.DEBUGC(func() string {
								return fmt.Sprintln(
									"solution found h:", nH,
									hash.String(),
									fork.List[fork.GetCurrent(nH)].
										AlgoVers[mb.Header.Version],
									"total hashes since startup",
									w.roller.C.Load()-int32(w.startNonce),
									fork.IsTestnet,
									mb.Header.Version,
									mb.Header.Bits,
									mb.Header.MerkleRoot.String(),
									hash,
								)
							})
							log.SPEW(mb)
							srs := sol.GetSolContainer(w.senderPort.Load(), mb)
							select {
							case <-w.Quit:
								log.DEBUG("worker stopping in the middle of it")
								break pausing
							default:
							}
							err := w.dispatchConn.SendMany(sol.SolutionMagic,
								transport.GetShards(srs.Data))
							if err != nil {
								log.ERROR(err)
							}
							log.DEBUG("sent solution")
							break running
						}
						mb.Header.Version = nextAlgo
						mb.Header.Bits = w.bitses.Load().(map[int32]uint32)[mb.Header.Version]
						mb.Header.Nonce++
						w.msgBlock.Store(*mb)
						// if we have completed a cycle report the hashrate on starting new algo
						// log.DEBUG(w.hashCount.Load(), uint64(w.roller.RoundsPerAlgo), w.roller.C)
						select {
						case <-w.Quit:
							log.DEBUG("worker stopping in the middle of it")
							break pausing
						default:
						}
						
					}
				}
			}
			log.TRACE("worker pausing")
		}
		log.TRACE("worker finished")
		// w.Close()
	}(w)
	return w
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(quit chan struct{}) (w *Worker, conn net.Conn) {
	// log.L.SetLevel("trace", true)
	sc := stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(&sc, quit), &sc
}

// NewJob is a delivery of a new job for the worker,
// this makes the miner start mining from pause or pause,
// prepare the work and restart
func (w *Worker) NewJob(job *job.Container, reply *bool) (err error) {
	if !w.dispatchReady.Load() {
		*reply = true
		return
	}
	// log.DEBUG("running NewJob RPC method")
	// if w.dispatchConn.SendConn == nil || len(w.dispatchConn.SendConn) < 1 {
	// log.DEBUG("loading dispatch connection from job message")
	// log.TRACE(job.String())
	// if there is no dispatch connection, make one.
	// If there is one but the server died or was disconnected the
	// connection the existing dispatch connection is nilled and this
	// will run. If there is no controllers on the network,
	// the worker pauses
	// ips := job.GetIPs()
	hashes := job.GetHashes()
	if hashes[5].IsEqual(w.lastMerkle) {
		// log.DEBUG("not a new job")
		*reply = true
		return
	}
	w.lastMerkle = hashes[5]
	// var addresses []string
	// for i := range ips {
	// 	generally there is only one but if a server had two interfaces
	// 	to different LANs it would send both
	// 	addresses = append(addresses, ips[i].String()+":"+
	// 		fmt.Sprint(job.GetControllerListenerPort()))
	// }
	// address := ips[0].String() + ":" + fmt.Sprint(job.GetControllerListenerPort())
	// remoteAddress := address
	// if w.dispatchConn != nil {
	// 	ra := w.dispatchConn.Sender
	// 	if ra != nil {
	// 		remoteAddress = ra.RemoteAddr().String()
	// 	}
	// }
	// if address != remoteAddress {
	// 	log.DEBUG("setting destination", address)
	// 	err = w.dispatchConn.SetDestination(address)
	// 	if err != nil {
	// 		log.ERROR(err)
	// 	}
	// }
	// }
	// log.SPEW(w.dispatchConn)
	*reply = true
	// halting current work
	w.stopChan <- struct{}{}
	bitses := job.GetBitses()
	w.bitses.Store(bitses)
	w.hashes.Store(hashes)
	newHeight := job.GetNewHeight()
	var algos []int32
	for i := range bitses {
		// we don't need to know net params if version numbers come with jobs
		algos = append(algos, i)
	}
	w.roller.Algos.Store(algos)
	mbb := w.msgBlock.Load().(wire.MsgBlock)
	mb := &mbb
	mb.Header.PrevBlock = *job.GetPrevBlockHash()
	// TODO: ensure worker time sync - ntp? time wrapper with skew adjustment
	hv := w.roller.GetAlgoVer()
	mb.Header.Version = hv
	b := w.bitses.Load().(map[int32]uint32)
	var ok bool
	mb.Header.Bits, ok = b[mb.Header.Version]
	if !ok {
		return errors.New("bits are empty")
	}
	rand.Seed(time.Now().UnixNano())
	mb.Header.Nonce = rand.Uint32()
	if w.hashes.Load() == nil {
		return errors.New("failed to decode merkle roots")
	} else {
		h := w.hashes.Load().(map[int32]*chainhash.Hash)
		hh, ok := h[hv]
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
	w.stopChan <- struct{}{}
	w.startChan <- struct{}{}
	return
}

// Pause signals the worker to stop working,
// releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	log.DEBUG("pausing from IPC")
	w.stopChan <- struct{}{}
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	log.DEBUG("stopping from IPC")
	w.stopChan <- struct{}{}
	defer close(w.Quit)
	*reply = true
	return
}

// SendPass gives the encryption key configured in the kopach controller (
// pod) configuration to allow workers to dispatch their solutions
func (w *Worker) SendPass(pass string, reply *bool) (err error) {
	log.DEBUG("receiving dispatch password", pass)
	rand.Seed(time.Now().UnixNano())
	// sp := fmt.Sprint(rand.Intn(32767) + 1025)
	// rp := fmt.Sprint(rand.Intn(32767) + 1025)
	var conn *transport.Channel
	conn, err = transport.NewBroadcastChannel("kopachworker", w, pass,
		transport.DefaultPort, controller.MaxDatagramSize, transport.Handlers{}, w.Quit)
	if err != nil {
		log.ERROR(err)
	}
	w.dispatchConn = conn
	w.dispatchReady.Store(true)
	*reply = true
	return
}
