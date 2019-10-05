package kopach

import (
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"github.com/p9c/pod/pkg/broadcast"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"github.com/ugorji/go/codec"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/conte"
)

// Main is the entry point for the kopach miner
func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	log.WARN("starting kopach standalone miner worker")
	m := newMsgHandle(*cx.Config.MinerPass)
	m.blockChan = make(chan controller.Blocks)
	blockSemaphore := make(chan struct{})
	var rotator atomic.Uint64
	// mining work dispatch goroutine
	go func() {
	workOut:
		for {
			select {
			case bt := <-m.blockChan:
				switch {
				// if the channel is returning nil it has been closed
				case bt == nil:
					break workOut
				// empty block templates means stop work and don't start
				// new work
				case len(bt) < 1:
					close(blockSemaphore)
					blockSemaphore = make(chan struct{})
				// received a normal block template
				default:
					// if workers are working, stop them
					if blockSemaphore != nil {
						close(blockSemaphore)
						blockSemaphore = make(chan struct{})
					}
					curHeight := bt[0].Height
					for i := 0; i < *cx.Config.GenThreads; i++ {
						// start up worker
						go func() {
							tn := time.Now()
							log.WARN("starting worker", i, tn)
							j := i
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
								for j := range bt {
									if bt[j].Block.Header.Version == algoVer {
										msgBlock = bt[j].Block
										found = true
									}
								}
								if !found { // this really shouldn't happen
									break threadOut
								}
								// start attempting to solve block
								// Choose a random extra nonce offset for
								// this block template and worker.
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
								// now := time.Now()
								// for extraNonce := eN; extraNonce < eN+maxExtraNonce; extraNonce++ {
								did := false
								extraNonce := eN
								// we only do this once
								for !did {
									did = true
									// Update the extra nonce in the block template with the new value by
									// regenerating the coinbase script and setting the merkle root to the
									// new value.
									log.TRACE("updating extraNonce")
									err := UpdateExtraNonce(msgBlock,
										curHeight+1, extraNonce+enOffset)
									if err != nil {
										log.WARN(err)
									}
									// Search through the entire nonce range for a solution while
									// periodically checking for early quit and stale block conditions along
									// with updates to the speed monitor.
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
									defer func() {
										log.DEBUGF(
											"wrkr: %d finished %d rounds of"+
												" %s", j, i-rNonce-1,
											fork.GetAlgoName(
												msgBlock.Header.Version,
												curHeight+1))
									}()
									log.TRACE("starting round from ", rNonce)
									for i = rNonce; i <= rNonce+mn; i++ {
										// if time.Now().Sub(now) > time.Second*3 {
										// 	return false
										// }
										select {
										case <-quit:
											return
										default:
										}
										var incr uint64 = 1
										header.Nonce = i
										hash := header.BlockHashWithAlgos(
											curHeight + 1)
										hashesCompleted += incr
										// The block is solved when the new block hash is less than the target
										// difficulty.  Yay!
										bigHash := blockchain.HashToBig(&hash)
										if bigHash.Cmp(targetDifficulty) <= 0 {
											log.WARN("found block", )
											// broadcast solved block
											break
										}
									}
								}
								select {
								case <-blockSemaphore:
									break threadOut
								default:
								}
							}
							log.WARN("worker", j, tn, "stopped")
						}()
					}
				}
			case <-quit:
				close(m.blockChan)
				break workOut
			}
		}
	}()
	go func() {
	out:
		for {
			cancel := broadcast.Listen(broadcast.DefaultAddress, m.msgHandler)
			select {
			//case bt := <-blockChan:
			//	log.WARN("received block templates")
			//	if bt == nil || len(bt) < 1 {
			//		log.WARN("empty blocks, stopping work")
			//	}
			case <-quit:
				log.DEBUG("quitting on quit channel close")
				cancel()
				break out
			}
		}
		wg.Done()
	}()
}

type msgBuffer struct {
	buffers [][]byte
	first   time.Time
	decoded bool
}

type msgHandle struct {
	buffers   map[string]*msgBuffer
	ciph      *cipher.AEAD
	dec       *codec.Decoder
	decBuf    []byte
	blockChan chan controller.Blocks
}

func newMsgHandle(password string) (out *msgHandle) {
	out = &msgHandle{}
	out.buffers = make(map[string]*msgBuffer)
	ciph := gcm.GetCipher(password)
	out.ciph = &ciph
	var mh codec.MsgpackHandle
	out.decBuf = make([]byte, 0, broadcast.MaxDatagramSize)
	out.dec = codec.NewDecoderBytes(out.decBuf, &mh)
	return
}

func (m *msgHandle) msgHandler(src *net.UDPAddr, n int, b []byte) {
	// remove any expired message bundles in the cache
	var deleters []string
	for i := range m.buffers {
		if time.Now().Sub(m.buffers[i].first) > time.Millisecond*50 {
			deleters = append(deleters, i)
		}
	}
	for i := range deleters {
		log.TRACE("deleting old message buffer")
		delete(m.buffers, deleters[i])
	}
	b = b[:n]
	//log.SPEW(b)
	if n < 16 {
		log.ERROR("received short broadcast message")
		return
	}
	// snip off message magic bytes
	msgType := string(b[:8])
	b = b[8:]
	log.TRACE(n, " bytes read from ", src, "message type", msgType)
	if msgType == string(broadcast.Template) {
		log.TRACE("got block template shard")
		buffer := b
		nonce := string(b[:8])
		if x, ok := m.buffers[nonce]; ok {
			log.TRACE("additional shard with nonce", hex.EncodeToString([]byte(nonce)))
			if !x.decoded {
				log.TRACE("adding shard")
				x.buffers = append(x.buffers, buffer)
				lb := len(x.buffers)
				log.TRACE("have", lb, "buffers")
				if lb > 2 {
					// try to decode it
					bytes, err := broadcast.Decode(*m.ciph, x.buffers)
					if err != nil {
						log.ERROR(err)
						return
					}
					//log.SPEW(bytes)
					m.dec.ResetBytes(bytes)
					blocks := controller.Blocks{}
					err = m.dec.Decode(&blocks)
					if err != nil {
						log.ERROR(err)
					}
					log.INFO(len(blocks), "block templates received")
					x.decoded = true
					// mine on it
					m.blockChan <- blocks
				}
			} else if x.buffers != nil {
				log.TRACE("nilling buffers")
				x.buffers = nil
			} else {
				log.TRACE("ignoring already decoded message shard")
			}
		} else {
			log.TRACE("adding nonce", hex.EncodeToString([]byte(nonce)))
			m.buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(), false}
			m.buffers[nonce].buffers = append(m.buffers[nonce].buffers, b)
		}
	}
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
	// TODO(davec): A util.Block should use saved in the state to avoid
	// recalculating all of the other transaction hashes.
	// block.Transactions[0].InvalidateCache() Recalculate the merkle root with
	// the updated extra nonce.
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
