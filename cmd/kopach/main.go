//go:generate go run ../tools/genmsghandle/main.go kopach controller.Blocks broadcast.TplBlock github.com/p9c/pod/pkg/controller msghandle.go
package kopach

import (
	"fmt"
	"github.com/p9c/pod/pkg/broadcast"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"github.com/ugorji/go/codec"
	"go.uber.org/atomic"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/conte"
)

// Main is the entry point for the kopach miner
func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	log.WARN("starting kopach standalone miner worker")
	returnChan := make(chan *controller.Blocks)
	m := newMsgHandle(*cx.Config.MinerPass, returnChan)
	blockSemaphore := make(chan struct{})
	outAddr, err := broadcast.New(*cx.Config.BroadcastAddress)
	if err != nil {
		log.ERROR(err)
		return
	}
	// create buffer and load into msgpack codec
	var mh codec.MsgpackHandle
	bytes := make([]byte, 0, broadcast.MaxDatagramSize)
	enc := codec.NewEncoderBytes(&bytes, &mh)
	var rotator atomic.Uint64
	var started atomic.Bool
	// mining work dispatch goroutine
	go func() {
	workOut:
		for {
			select {
			case bt := <-m.returnChan:
				switch {
				// if the channel is returning nil it has been closed
				case bt == nil:
					break workOut
				// received a normal block template
				default:
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
					// if workers are working, stop them
					if blockSemaphore != nil {
						close(blockSemaphore)
						blockSemaphore = make(chan struct{})
					}
					curHeight := bt.Templates[0].Height
					for i := 0; i < *cx.Config.GenThreads; i++ {
						curr :=i
						// start up worker
						go func() {
							tn := time.Now()
							log.DEBUG("starting worker", curr, tn)
							j := curr
						threadOut:
							for {
								// choose the algorithm on a rolling cycle
								counter := rotator.Load()
								algo := "sha256d"
								switch fork.GetCurrent(curHeight + 1) {
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
								rotator.Add(1)
								log.WARN("worker", j, "algo", algo)
								algoVer := fork.GetAlgoVer(algo, curHeight+1)
								var msgBlock *wire.MsgBlock
								found := false
								for j := range bt.Templates {
									if bt.Templates[j].Block.Header.Version == algoVer {
										msgBlock = bt.Templates[j].Block
										found = true
									}
								}
								if !found { // this really shouldn't happen
									break threadOut
								}
								// start attempting to solve block
								enOffset, err := wire.RandomUint64()
								if err != nil {
									log.WARNF(
										"unexpected error while generating"+
											" random extra nonce offset:", err)
									enOffset = 0
								}
								// Create some convenience variables.
								header := &msgBlock.Header
								targetDifficulty := fork.CompactToBig(header.Bits)
								// Initial state.
								hashesCompleted := uint64(0)
								eN, _ := wire.RandomUint64()
								did := false
								extraNonce := eN
								// we only do this once
								for !did {
									did = true
									// use a random extra nonce to ensure no
									// duplicated work
									err := UpdateExtraNonce(msgBlock,
										curHeight+1, extraNonce+enOffset)
									if err != nil {
										log.WARN(err)
									}
									var shifter uint64 = 16
									rn, _ := wire.RandomUint64()
									if rn > 1<<63-1<<shifter {
										rn -= 1 << shifter
									}
									rn += 1 << shifter
									rNonce := uint32(rn)
									mn := uint32(27)
									mn = 1 << 8 * uint32(*cx.Config.GenThreads)
									var i uint32
									//log.TRACE("starting round from ", rNonce)
									for i = rNonce; i <= rNonce+mn; i++ {
										select {
										case <-quit:
											break
										default:
										}
										var incr uint64 = 1
										header.Nonce = i
										hash := header.BlockHashWithAlgos(
											curHeight + 1)
										hashesCompleted += incr
										// The block is solved when the new
										// block hash is less than the target
										// difficulty.  Yay!
										bigHash := blockchain.HashToBig(&hash)
										if bigHash.Cmp(targetDifficulty) <= 0 {
											log.WARN("found block")
											// broadcast solved block:
											// first stop all work
											if blockSemaphore == nil {
												close(blockSemaphore)
												blockSemaphore = nil
											}
											// serialize the block
											bytes = bytes[:0]
											enc.ResetBytes(&bytes)
											err := enc.Encode(msgBlock)
											if err != nil {
												log.ERROR(err)
												break
											}
											err = broadcast.Send(outAddr,
												bytes, *m.ciph, broadcast.Solution)
											break threadOut
										}
									}
								}
								select {
								case <-quit:
									break threadOut
								case <-blockSemaphore:
									break threadOut
								default:
								}
							}
							log.DEBUG("worker", j, tn, "stopped")
							started.Store(false)
						}()
					}
				}
			case <-quit:
				close(m.returnChan)
				break workOut
			}
		}
	}()
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
	//  block.Transactions[0].InvalidateCache() Recalculate the merkle root with
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
