package worker

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/controller/job"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/util"
)

const RoundsPerAlgo = 1 << 8

type Worker struct {
	sem        sem.T
	conn       net.Conn
	Quit       chan struct{}
	run        sem.T
	block      *util.Block
	msgBlock   *wire.MsgBlock
	bitses     map[int32]uint32
	roller     *Counter
	startNonce int
}

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
		sem:      s,
		conn:     conn,
		Quit:     quit,
		run:      sem.New(1),
		block:    util.NewBlock(msgBlock),
		msgBlock: msgBlock,
		roller:   NewCounter(RoundsPerAlgo),
	}
	// with this we can report cumulative hash counts as well as using it to
	// distribute algorithms evenly
	w.startNonce = w.roller.C
	ticker := time.NewTicker(time.Second)
	//w.sem.Acquire()
	go func() {
		log.DEBUG("main work loop starting")
	out:
		for {
			select {
			case <-w.sem.Release():
				log.DEBUG("pausing work")
				// release runner semaphore so worker won't run until it's
				// reacquired by NewJob.
				// This select will block here until the run semaphore is
				// acquired, ie, no cpu cycles in this goroutine (
				// which is main really)
				select {
				case <-w.sem.Release():
					log.DEBUG("semaphore released from pause")
				case <-w.run.Release():
					log.DEBUG("run semaphore released from pause")
					w.run.Acquire()
					break
				case <-w.Quit:
					log.DEBUG("quitting from pause")
					break out
				}
				log.DEBUG("pause acquiring run semaphore")
				w.run.Acquire()
			case <-w.Quit:
				// quit when w.Stop() is called
				log.DEBUG("worker stopping on quit message")
				break out
			case <-w.run.Release():
				log.DEBUG("run semaphore released")
				// do a round
				hash := w.msgBlock.BlockHashWithAlgos(w.block.Height())
				bigHash := blockchain.HashToBig(&hash)
				if bigHash.Cmp(
					fork.CompactToBig(w.msgBlock.Header.Bits)) <= 0 {
					// yay we win!
					log.DEBUG("found a block")
					log.SPEW(w.block.MsgBlock())
					w.sem.Acquire()
					log.DEBUG("found block acquired semaphore")
					break
				}
				log.DEBUGF("%065x",bigHash.Bytes())
				// prepare for next round
				w.msgBlock.Header.Nonce++
				// trigger another to follow, this can be stopped by releasing
				w.run.Acquire()
				log.DEBUG("runner reacquired semaphore")
			case <-ticker.C:
				log.DEBUG("timestamp update ticker")
				w.msgBlock.Header.Timestamp = time.Now()
			}
		}
		log.DEBUG("worker finished")
	}()
	return w
}

// New initialises the state for a worker,
// loading the work function handler that runs a round of processing between
// checking quit signal and work semaphore
func New(s sem.T) (w *Worker, conn net.Conn) {
	quit := make(chan struct{})
	conn = stdconn.New(os.Stdin, os.Stdout, quit)
	return NewWithConnAndSemaphore(
		conn,
		s,
		quit), conn
}

// NewJob is a delivery of a new job for the worker, this starts a miner thread
func (w *Worker) NewJob(job *job.Container, reply *bool) (err error) {
	*reply = true
	// previous thread loses its semaphore when a new job arrives,
	// this acts as a mutex because this worker is single thread with RPC
	// handler calls changing main thread's activity mode,
	// this semaphore ensures the worker is not accessing what we are about
	// to update
	// This is concurrent code but only single thread,
	// so a semaphore acts here as a mutex between listener and miner
	log.DEBUG("new job acquiring pause semaphore")
	w.sem.Acquire()
	log.DEBUG("new job acquired pause semaphore")
	// load the new value in the bitses and MsgBlock
	w.bitses = job.GetBitses()
	newHeight := job.GetNewHeight()
	w.roller.Algos = []int32{}
	for i := range w.bitses {
		// we don't need to know net params if version numbers come with jobs
		w.roller.Algos = append(w.roller.Algos, i)
	}
	w.block.SetHeight(newHeight)
	w.msgBlock.Header.PrevBlock = *job.GetPrevBlockHash()
	// TODO: ensure worker time sync - ntp? time wrapper with skew adjustment
	w.msgBlock.Header.Version = w.roller.GetAlgoVer()
	w.msgBlock.Header.Bits = w.bitses[w.msgBlock.Header.Version]
	rand.Seed(time.Now().UnixNano())
	w.msgBlock.Header.Nonce = rand.Uint32()
	w.msgBlock.Transactions = job.GetTxs()
	// create the unique extra nonce for this worker,
	// which creates a different merkel root
	extraNonce, err := wire.RandomUint64()
	if err != nil {
		log.ERROR(err)
		return
	}
	log.DEBUG("updating extra nonce")
	err = UpdateExtraNonce(w.msgBlock, newHeight, extraNonce)
	if err != nil {
		log.ERROR(err)
		return
	}
	// make the work select block start running
	//log.DEBUG("releasing pause semaphore")
	//<-w.sem.Release()
	log.DEBUG("acquiring run semaphore to start new work")
	w.run.Acquire()
	return
}

// Pause signals the worker to stop working,
// releases its semaphore and the worker is then idle
func (w *Worker) Pause(_ int, reply *bool) (err error) {
	log.DEBUG("pausing from IPC")
	w.sem.Acquire()
	log.DEBUG("pause from IPC acquired semaphore")
	//w.run.Acquire()
	//<-w.run.Release()
	*reply = true
	return
}

// Stop signals the worker to quit
func (w *Worker) Stop(_ int, reply *bool) (err error) {
	log.DEBUG("stopping from IPC")
	w.sem.Acquire()
	defer close(w.Quit)
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
	log.SPEW(msgBlock.Transactions)
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
