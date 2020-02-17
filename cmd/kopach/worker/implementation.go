package worker

import (
	"crypto/cipher"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
	
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
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

const RoundsPerAlgo = 23

type Worker struct {
	sem          sem.T
	conn         net.Conn
	dispatchConn *transport.Connection
	ciph         cipher.AEAD
	Quit         chan struct{}
	run          sem.T
	block        *util.Block
	msgBlock     *wire.MsgBlock
	bitses       map[int32]uint32
	hashes       map[int32]*chainhash.Hash
	roller       *Counter
	startNonce   uint32
	startChan    chan struct{}
	stopChan     chan struct{}
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
	conn net.Conn,
	s sem.T,
	quit chan struct{},
) *Worker {
	log.DEBUG("creating new worker")
	msgBlock := &wire.MsgBlock{Header: wire.BlockHeader{}}
	w := &Worker{
		sem:       s,
		conn:      conn,
		Quit:      quit,
		run:       sem.New(1),
		block:     util.NewBlock(msgBlock),
		msgBlock:  msgBlock,
		roller:    NewCounter(RoundsPerAlgo),
		startChan: make(chan struct{}),
		stopChan:  make(chan struct{}),
	}
	// with this we can report cumulative hash counts as well as using it to
	// distribute algorithms evenly
	w.startNonce = uint32(w.roller.C)
	go func() {
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
			log.DEBUG("worker running")
			// Run state
		running:
			for {
				select {
				case <-w.startChan:
					// drain start channel in run mode
					continue
				case <-w.stopChan:
					w.block = nil
					w.bitses = nil
					w.hashes = nil
					break running
				case <-w.Quit:
					log.DEBUG("worker stopping on pausing message")
					break pausing
				default:
					if w.block == nil || w.bitses == nil || w.hashes == nil {
						// log.INFO("stop was called before we started working")
					} else {
						// work
						nH := w.block.Height()
						w.msgBlock.Header.Version = w.roller.GetAlgoVer()
						w.msgBlock.Header.MerkleRoot = *w.hashes[w.msgBlock.Header.Version]
						w.msgBlock.Header.Bits = w.bitses[w.msgBlock.Header.Version]
						select {
						case <-w.stopChan:
							w.block = nil
							w.bitses = nil
							w.hashes = nil
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
							err := w.dispatchConn.Send(srs.Data, sol.SolutionMagic)
							if err != nil {
								log.ERROR(err)
							}
							log.DEBUG("sent solution")
							break running
						}
						nextAlgo := w.roller.GetAlgoVer()
						w.msgBlock.Header.Version = nextAlgo
						w.msgBlock.Header.Bits = w.bitses[w.msgBlock.Header.Version]
						w.msgBlock.Header.Nonce++
						// if we have completed a cycle report the hashrate on starting new algo
						if w.roller.C%w.roller.RoundsPerAlgo == 0 {
							// since := int(time.Now().Sub(tn)/time.Second) + 1
							// total := w.roller.C - int(w.startNonce)
							// _, _ = fmt.Fprintf(os.Stderr,
							// 	"\r %9d hash/s %s       \r", total/since, fork.GetAlgoName(w.msgBlock.Header.Version, nH))
							// send out broadcast containing worker nonce and algorithm and count of blocks
							hashReport := hashrate.Get(w.roller.RoundsPerAlgo, nextAlgo, nH)
							err := w.dispatchConn.Send(hashReport.Data, hashrate.HashrateMagic)
							if err != nil {
								log.ERROR(err)
							}
						}
					}
				}
			}
			log.DEBUG("worker pausing")
		}
		log.DEBUG("worker finished")
	}()
	return w
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(s sem.T) (w *Worker, conn net.Conn) {
	log.L.SetLevel("trace", true)
	quit := make(chan struct{})
	conn = stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(
		conn,
		s,
		quit), conn
}

// NewJob is a delivery of a new job for the worker,
// this makes the miner start mining from pause or pause,
// prepare the work and restart
func (w *Worker) NewJob(job *job.Container, reply *bool) (err error) {
	// log.DEBUG("running NewJob RPC method")
	// if w.dispatchConn.SendConn == nil || len(w.dispatchConn.SendConn) < 1 {
	log.DEBUG("loading dispatch connection from job message")
	log.TRACE(job.String())
	// if there is no dispatch connection, make one.
	// If there is one but the server died or was disconnected the
	// connection the existing dispatch connection is nilled and this
	// will run. If there is no controllers on the network,
	// the worker pauses
	ips := job.GetIPs()
	var addresses []string
	for i := range ips {
		// generally there is only one but if a server had two interfaces
		// to different lans it would send both
		addresses = append(addresses, ips[i].String()+":"+
			fmt.Sprint(job.GetControllerListenerPort()))
	}
	err = w.dispatchConn.SetSendConn(addresses...)
	if err != nil {
		log.ERROR(err)
	}
	// }
	// log.SPEW(w.dispatchConn)
	*reply = true
	// halting current work
	w.stopChan <- struct{}{}
	w.bitses = job.GetBitses()
	w.hashes = job.GetHashes()
	newHeight := job.GetNewHeight()
	w.roller.Algos = []int32{}
	for i := range w.bitses {
		// we don't need to know net params if version numbers come with jobs
		w.roller.Algos = append(w.roller.Algos, i)
	}
	w.msgBlock.Header.PrevBlock = *job.GetPrevBlockHash()
	// TODO: ensure worker time sync - ntp? time wrapper with skew adjustment
	w.msgBlock.Header.Version = w.roller.GetAlgoVer()
	w.msgBlock.Header.Bits = w.bitses[w.msgBlock.Header.Version]
	rand.Seed(time.Now().UnixNano())
	w.msgBlock.Header.Nonce = rand.Uint32()
	// log.TRACE(w.hashes)
	if w.hashes != nil {
		w.msgBlock.Header.MerkleRoot = *w.hashes[w.msgBlock.Header.Version]
	} else {
		return errors.New("failed to decode merkle roots")
	}
	w.msgBlock.Header.Timestamp = time.Now()
	// halting current work
	w.stopChan <- struct{}{}
	// create the unique extra nonce for this worker,
	// which creates a different merkel root
	// extraNonce, err := wire.RandomUint64()
	// if err != nil {
	// 	log.ERROR(err)
	// 	return
	// }
	// log.TRACE("updating extra nonce")
	// err = UpdateExtraNonce(w.msgBlock, newHeight, extraNonce)
	// if err != nil {
	// 	log.ERROR(err)
	// 	return
	// }
	// log.SPEW(w.msgBlock)
	// make the work select block start running
	w.block = util.NewBlock(w.msgBlock)
	w.block.SetHeight(newHeight)
	// halting current work
	w.stopChan <- struct{}{}
	// log.INFO("height", newHeight)
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
	log.DEBUG("receiving dispatch password")
	conn, err := transport.NewConnection("", "", pass, controller.MaxDatagramSize, nil, false)
	if err != nil {
		log.ERROR(err)
	}
	w.dispatchConn = conn
	*reply = true
	return
}

// UpdateExtraNonce updates the extra nonce in the coinbase script of the
// passed block by regenerating the coinbase script with the passed value and
// block height.  It also recalculates and updates the new merkle root that
// results from changing the coinbase script.
func UpdateExtraNonce(msgBlock *wire.MsgBlock, blockHeight int32,
	extraNonce uint64) error {
	if msgBlock == nil {
		log.ERROR("cannot update a nil MsgBlock")
	}
	log.DEBUG("UpdateExtraNonce")
	coinbaseScript, err := standardCoinbaseScript(blockHeight, extraNonce)
	if err != nil {
		return err
	}
	if len(coinbaseScript) > blockchain.MaxCoinbaseScriptLen {
		return fmt.Errorf(
			"coinbase transaction script length of %d is out of range ("+
				"min: %d, max: %d)",
			len(coinbaseScript),
			blockchain.MinCoinbaseScriptLen,
			blockchain.MaxCoinbaseScriptLen)
	}
	// log.SPEW(msgBlock.Transactions)
	msgBlock.Transactions[0].TxIn[0].SignatureScript = coinbaseScript
	// TODO(davec): A util.Solution should use saved in the state to avoid
	//  recalculating all of the other transaction hashes.
	//  block.Transaction[0].InvalidateCache() Recalculate the merkle root with
	//  the updated extra nonce.
	block := util.NewBlock(msgBlock)
	log.DEBUG("recalculating merkle root")
	merkles := blockchain.BuildMerkleTreeStore(block.Transactions(), false)
	msgBlock.Header.MerkleRoot = *merkles[len(merkles)-1]
	return nil
}

// standardCoinbaseScript returns a standard script suitable for use as the
// signature script of the coinbase transaction of a new block.  In particular,
// it starts with the block height that is required by version 2 blocks and
// adds the extra nonce as well as additional coinbase flags.
func standardCoinbaseScript(nextBlockHeight int32, extraNonce uint64) ([]byte, error) {
	return txscript.NewScriptBuilder().AddInt64(int64(nextBlockHeight)).
		AddInt64(int64(extraNonce)).AddData([]byte(mining.CoinbaseFlags)).
		Script()
}

// calcMerkleRoot creates a merkle tree from the slice of transactions and returns the root of the tree.
func calcMerkleRoot(txns []*wire.MsgTx) chainhash.Hash {
	if len(txns) == 0 {
		return chainhash.Hash{}
	}
	utilTxns := make([]*util.Tx, 0, len(txns))
	for _, tx := range txns {
		utilTxns = append(utilTxns, util.NewTx(tx))
	}
	merkles := blockchain.BuildMerkleTreeStore(utilTxns, false)
	return *merkles[len(merkles)-1]
}
