//+build ignore



//go:generate go run ../tools/genmsghandle/main.go kopach controller.Blocks broadcast.TplBlock github.com/p9c/pod/pkg/controller msghandle.go
package kopach

import (
	"fmt"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"github.com/ugorji/go/codec"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

type workerConfig struct {
	bytes          []byte
	templates      controllerold.Templates
	enc            *codec.Encoder
	rotator        *atomic.Uint64
	outAddr        *net.UDPAddr
	blockSemaphore *chan struct{}
	submitLock     *sync.Mutex
	msgHandle      *msgHandle
	quit           chan struct{}
	workerNumber   int
	threads        int
	curHeight      int32
}

// Main is the entry point for the kopach miner
func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	log.WARN("starting kopach standalone miner worker")
	returnChan := make(chan *controllerold.Blocks)
	m := newMsgHandle(*cx.Config.MinerPass, returnChan)
	blockSemaphore := make(chan struct{})
	outAddr, err := broadcast.New(*cx.Config.BroadcastAddress)
	if err != nil {
		log.ERROR(err)
		return
	}
	// mining work dispatch goroutine
	var started atomic.Bool
	var rotator atomic.Uint64
	var submitLock sync.Mutex
	var oldBlocks controllerold.Blocks
	go func() {
	workOut:
		for {
			select {
			case bt := <-m.returnChan:
				// to handle receiving multiple new identical blocks from 2 or more servers in a LAN, if the hash of
				// the raw template is the same then we have already seen it and drop the message
				if len(oldBlocks.Templates) > 0 {
					var newPrev chainhash.Hash
					oldPrev := oldBlocks.Templates[0].Block.BlockHash()
					for tp := range bt.Templates {
						// maps are not linear so we have to search for the match in the second
						if bt.Templates[tp].Block.Header.Version == oldBlocks.Templates[0].Block.Header.Version {
							newPrev = bt.Templates[tp].Block.BlockHash()
						}
					}
					if oldPrev.IsEqual(&newPrev) {
						// ignore the message
						continue
					}
				}
				switch {
				// if the channel is returning nil it has been closed
				case bt == nil:
					break workOut
				// received a normal block template
				default:
					var mh codec.MsgpackHandle
					if len(bt.Templates) < 1 {
						close(blockSemaphore)
						submitLock.Lock()
						log.WARN("received empty templates, halting work")
						blockSemaphore = make(chan struct{})
						submitLock.Unlock()
						break
					}
					// If a worker is running and the block templates are not marked new, ignore
					if started.Load() {
						if !bt.New && blockSemaphore != nil {
							//log.TRACE("already started, block is not new, ignoring")
							break
						}
					} else {
						log.WARN("starting mining")
						started.Store(true)
					}
					// TODO: handle multiple servers later sending templates
					// if workers are working, stop them
					if blockSemaphore != nil {
						log.WARN("stopping currently running miners")
						submitLock.Lock()
						close(blockSemaphore)
						blockSemaphore = make(chan struct{})
						submitLock.Unlock()
					}
					curHeight := bt.Templates[0].Height
					// create a copy of the templates for each worker thread
					numWorkers := *cx.Config.GenThreads
					templates := bt.Templates.Copy(numWorkers)
					for i := 0; i < numWorkers; i++ {
						bytes := make([]byte, 0, broadcast.MaxDatagramSize)
						enc := codec.NewEncoderBytes(&bytes, &mh)
						go mine(workerConfig{
							templates:      templates[i],
							enc:            enc,
							rotator:        &rotator,
							outAddr:        outAddr,
							blockSemaphore: &blockSemaphore,
							submitLock:     &submitLock,
							msgHandle:      m,
							quit:           quit,
							workerNumber:   i,
							threads:        *cx.Config.GenThreads,
							curHeight:      curHeight,
						})
					}
				}
			case <-quit:
				close(m.returnChan)
				started.Store(false)
				break workOut
			}
		}
	}()
	// quit goroutine that ensures context is cancelled
	go func() {
		cancel := broadcast.Listen(broadcast.DefaultAddress, m.msgHandler)
	out:
		for {
			select {
			case <-quit:
				log.DEBUG("quitting on quit channel close")
				cancel()
				break out
			}
		}
		wg.Done()
	}()
}

func mine(cfg workerConfig) {
	// start up worker
	tn := time.Now()
	log.DEBUG("starting worker", cfg.workerNumber, tn)
threadOut:
	for {
		// choose the algorithm on a rolling cycle
		counter := cfg.rotator.Load()
		cfg.rotator.Add(1)
		algo := "sha256d"
		switch fork.GetCurrent(cfg.curHeight) {
		case 0:
			if counter&1 == 1 {
				algo = "sha256d"
			} else {
				algo = "scrypt"
			}
		case 1:
			l9 := uint64(len(fork.P9AlgoVers))
			mod := counter % l9
			algo = fork.P9AlgoVers[int32(mod+5)]
		}
		log.WARN("worker", cfg.workerNumber, "algo", algo)
		algoVer := fork.GetAlgoVer(algo, cfg.curHeight)
		var msgBlock *wire.MsgBlock
		found := false
		for j := range cfg.templates {
			if cfg.templates[j].Block.Header.Version == algoVer {
				msgBlock = cfg.templates[j].Block
				found = true
			}
		}
		if !found { // this really shouldn't happen
			break threadOut
		}
		// start attempting to solve block
		enOffset, err := wire.RandomUint64()
		if err != nil {
			log.WARNF("unexpected error while generating random extra nonce offset:", err)
			enOffset = 0
		}
		// Create some convenience variables.
		header := &msgBlock.Header
		targetDifficulty := fork.CompactToBig(header.Bits)
		// Initial state.
		hashesCompleted := uint64(0)
		eN, _ := wire.RandomUint64()
		extraNonce := eN
		// use a random extra nonce to ensure no duplicated work
		err2 := UpdateExtraNonce(msgBlock, cfg.curHeight, extraNonce+enOffset)
		if err2 != nil {
			log.WARN(err2)
		}
		var shifter uint64 = 16
		rn, _ := wire.RandomUint64()
		if rn > 1<<63-1<<shifter {
			rn -= 1 << shifter
		}
		rn += 1 << shifter
		rNonce := uint32(rn)
		mn := uint32(27)
		mn = 1 << 16 // * uint32(cfg.threads)
		var nonce uint32
		log.DEBUG("starting round from ", rNonce, algo)
		for nonce = rNonce; nonce <= rNonce+mn; nonce++ {
			select {
			case <-cfg.quit:
				log.DEBUG("quitting after", nonce-rNonce, "rounds of", algo)
				break
			default:
			}
			var incr uint64 = 1
			header.Nonce = nonce
			hash := header.BlockHashWithAlgos(cfg.curHeight)
			hashesCompleted += incr
			// The block is solved when the new
			// block hash is less than the target
			// difficulty.  Yay!
			bigHash := blockchain.HashToBig(&hash)
			if bigHash.Cmp(targetDifficulty) <= 0 {
				// broadcast solved block:
				// first stop all competing later submissions
				cfg.submitLock.Lock()
				// all other threads will terminate when the semaphore is
				// closed and they finish a work cycle
				if *cfg.blockSemaphore == nil {
					close(*cfg.blockSemaphore)
					*cfg.blockSemaphore = nil
				}
				log.WARN("found block", hash)
				// serialize the block
				cfg.bytes = cfg.bytes[:0]
				cfg.enc.ResetBytes(&cfg.bytes)
				err := cfg.enc.Encode(msgBlock)
				if err != nil {
					cfg.submitLock.Unlock()
					log.ERROR(err)
					break
				}
				log.SPEW(header)
				err = broadcast.Send(cfg.outAddr, cfg.bytes, *cfg.msgHandle.ciph,
					broadcast.Solution)
				if err != nil {
					log.ERROR(err)
				} else {
					log.INFO("submitted block")
				}
				cfg.submitLock.Unlock()
				break threadOut
			}
		}
		log.DEBUG("finished", nonce-rNonce, "rounds of", algo)
		select {
		// Job times out after 6 minutes (10 blocks)
		case <-time.After(6 * time.Minute):
			log.DEBUG("stopping due to no new blocks in 6 minutes")
			break threadOut
		case <-cfg.quit:
			log.DEBUG("stopping on quit signal")
			break threadOut
		case <-*cfg.blockSemaphore:
			log.DEBUG("stopping on close of semaphore")
			break threadOut
		default:
		}
	}
	log.DEBUG("worker", cfg.workerNumber, tn, "stopped")
}

// UpdateExtraNonce updates the extra nonce in the coinbase script of the
// passed block by regenerating the coinbase script with the passed value and
// block height.  It also recalculates and updates the new merkle root that
// results from changing the coinbase script.
func UpdateExtraNonce(msgBlock *wire.MsgBlock,
	blockHeight int32, extraNonce uint64) error {
	coinbaseScript, err := standardCoinbaseScript(blockHeight, extraNonce)
	if err != nil {
		return err
	}
	if len(coinbaseScript) > blockchain.MaxCoinbaseScriptLen {
		return fmt.Errorf(
			"coinbase transaction script length of %d is out of range (min: %d, max: %d)",
			len(coinbaseScript), blockchain.MinCoinbaseScriptLen,
			blockchain.MaxCoinbaseScriptLen)
	}
	msgBlock.Transactions[0].TxIn[0].SignatureScript = coinbaseScript
	// TODO(davec): A util.Solution should use saved in the state to avoid
	//  recalculating all of the other transaction hashes.
	//  block.Transaction[0].InvalidateCache() Recalculate the merkle root with
	//  the updated extra nonce.
	block := util.NewBlock(msgBlock)
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
