package worker

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	
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
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util"
)

const RoundsPerAlgo = 10

type Worker struct {
	sem           sem.T
	pipeConn      *stdconn.StdConn
	multicastConn net.Conn
	unicastConn   net.Conn
	dispatchConn  *transport.Channel
	dispatchReady atomic.Bool
	ciph          cipher.AEAD
	Quit          chan struct{}
	run           sem.T
	block         *util.Block
	msgBlock      *wire.MsgBlock
	bitses        atomic.Value
	hashes        atomic.Value
	lastMerkle    *chainhash.Hash
	roller        *Counter
	startNonce    uint32
	startChan     chan struct{}
	stopChan      chan struct{}
	// running    uint32
}

const (
	OFF uint32 = iota
	ON
)

type Counter struct {
	C             int
	Algos         []int32
	RoundsPerAlgo int
}

// NewCounter returns an initialized algorithm rolling counter that ensures
// each miner does equal amounts of every algorithm
func NewCounter(roundsPerAlgo int) (c *Counter) {
	// these will be populated when work arrives
	var algos []int32
	// Start the counter at a random position
	rand.Seed(time.Now().UnixNano())
	c = &Counter{
		C:             rand.Intn(roundsPerAlgo),
		Algos:         algos,
		RoundsPerAlgo: roundsPerAlgo,
	}
	return
}

// GetAlgoVer returns the next algo version based on the current configuration
func (c *Counter) GetAlgoVer() (ver int32) {
	// the formula below rolls through versions with blocks roundsPerAlgo
	// long for each algorithm by its index
	ver = c.Algos[(c.C/c.RoundsPerAlgo)%len(c.Algos)]
	c.C++
	return
}

// NewWithConnAndSemaphore is exposed to enable use an actual network
// connection while retaining the same RPC API to allow a worker to be
// configured to run on a bare metal system with a different launcher main
func NewWithConnAndSemaphore(
	conn *stdconn.StdConn,
	s sem.T,
	quit chan struct{},
) *Worker {
	log.DEBUG("creating new worker")
	msgBlock := &wire.MsgBlock{Header: wire.BlockHeader{}}
	w := &Worker{
		sem:       s,
		pipeConn:  conn,
		Quit:      quit,
		run:       sem.New(1),
		block:     util.NewBlock(msgBlock),
		msgBlock:  msgBlock,
		roller:    NewCounter(RoundsPerAlgo),
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
	}
	w.dispatchReady.Store(false)
	// with this we can report cumulative hash counts as well as using it to
	// distribute algorithms evenly
	w.startNonce = uint32(w.roller.C)
	go func(w *Worker) {
		log.DEBUG("main work loop starting")
	pausing:
		for {
			// Pause state
			select {
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
				case <-w.startChan:
					// drain start channel in run mode
					continue
				case <-w.stopChan:
					w.block = nil
					w.bitses.Store((map[int32]uint32)(nil))
					w.hashes.Store((map[int32]*chainhash.Hash)(nil))
					break running
				case <-w.Quit:
					log.DEBUG("worker stopping on pausing message")
					break pausing
				default:
					if w.block == nil || w.bitses.Load() == nil || w.hashes.Load() == nil ||
						!w.dispatchReady.Load() {
						// log.INFO("stop was called before we started working")
					} else {
						// work
						nH := w.block.Height()
						w.msgBlock.Header.Version = w.roller.GetAlgoVer()
						h := w.hashes.Load().(map[int32]*chainhash.Hash)
						if h != nil {
							w.msgBlock.Header.MerkleRoot = *h[w.msgBlock.Header.Version]
						}
						b := w.bitses.Load().(map[int32]uint32)
						if bb, ok := b[w.msgBlock.Header.Version]; ok {
							w.msgBlock.Header.Bits = bb
						}
						select {
						case <-w.stopChan:
							w.block = nil
							w.bitses.Store((map[int32]uint32)(nil))
							w.hashes.Store((map[int32]*chainhash.Hash)(nil))
							break running
						default:
						}
						hash := w.msgBlock.Header.BlockHashWithAlgos(nH)
						bigHash := blockchain.HashToBig(&hash)
						if bigHash.Cmp(fork.CompactToBig(w.msgBlock.Header.Bits)) <= 0 {
							log.DEBUGC(func() string {
								return fmt.Sprintln(
									"solution found h:", nH,
									hash.String(),
									fork.List[fork.GetCurrent(nH)].
										AlgoVers[w.msgBlock.Header.Version],
									"total hashes since startup",
									w.roller.C-int(w.startNonce),
									fork.IsTestnet,
									w.msgBlock.Header.Version,
									w.msgBlock.Header.Bits,
									w.msgBlock.Header.MerkleRoot.String(),
									hash,
								)
							})
							log.SPEW(w.msgBlock)
							srs := sol.GetSolContainer(w.msgBlock)
							err := w.dispatchConn.SendMany(sol.SolutionMagic,
								transport.GetShards(srs.Data))
							if err != nil {
								log.ERROR(err)
							}
							log.DEBUG("sent solution")
							break running
						}
						nextAlgo := w.roller.GetAlgoVer()
						w.msgBlock.Header.Version = nextAlgo
						w.msgBlock.Header.Bits = w.bitses.Load().(map[int32]uint32)[w.msgBlock.Header.Version]
						w.msgBlock.Header.Nonce++
						// if we have completed a cycle report the hashrate on starting new algo
						if w.roller.C%w.roller.RoundsPerAlgo == 0 {
							// since := int(time.Now().Sub(tn)/time.Second) + 1
							// total := w.roller.C - int(w.startNonce)
							// _, _ = fmt.Fprintf(os.Stderr,
							// 	"\r %9d hash/s %s       \r", total/since, fork.GetAlgoName(w.msgBlock.Header.Version, nH))
							// send out broadcast containing worker nonce and algorithm and count of blocks
							hashReport := hashrate.Get(w.roller.RoundsPerAlgo, nextAlgo, nH)
							err := w.dispatchConn.SendMany(hashrate.HashrateMagic,
								transport.GetShards(hashReport.Data))
							if err != nil {
								log.ERROR(err)
							}
						}
					}
				}
			}
			log.TRACE("worker pausing")
		}
		log.TRACE("worker finished")
	}(w)
	return w
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(s sem.T) (w *Worker, conn net.Conn) {
	// log.L.SetLevel("trace", true)
	quit := make(chan struct{})
	sc := stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(
		&sc,
		s,
		quit), &sc
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
	log.TRACE(job.String())
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
	// 	to different lans it would send both
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
	w.bitses.Store(job.GetBitses())
	w.hashes.Store(hashes)
	newHeight := job.GetNewHeight()
	w.roller.Algos = []int32{}
	for i := range w.bitses.Load().(map[int32]uint32) {
		// we don't need to know net params if version numbers come with jobs
		w.roller.Algos = append(w.roller.Algos, i)
	}
	w.msgBlock.Header.PrevBlock = *job.GetPrevBlockHash()
	// TODO: ensure worker time sync - ntp? time wrapper with skew adjustment
	w.msgBlock.Header.Version = w.roller.GetAlgoVer()
	b := w.bitses.Load().(map[int32]uint32)
	var ok bool
	w.msgBlock.Header.Bits, ok = b[w.msgBlock.Header.Version]
	if !ok {
		return errors.New("bits are empty")
	}
	rand.Seed(time.Now().UnixNano())
	w.msgBlock.Header.Nonce = rand.Uint32()
	if w.hashes.Load() == nil {
		return errors.New("failed to decode merkle roots")
	} else {
		h := w.hashes.Load().(map[int32]*chainhash.Hash)
		hh, ok := h[w.msgBlock.Header.Version]
		if !ok {
			return errors.New("could not get merkle root from job")
		}
		w.msgBlock.Header.MerkleRoot = *hh
	}
	w.msgBlock.Header.Timestamp = time.Now()
	// make the work select block start running
	w.block = util.NewBlock(w.msgBlock)
	w.block.SetHeight(newHeight)
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
		transport.DefaultPort, controller.MaxDatagramSize, nil)
	if err != nil {
		log.ERROR(err)
	}
	w.dispatchConn = conn
	w.dispatchReady.Store(true)
	*reply = true
	return
}
