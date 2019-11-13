package kopach

import "C"
import (
	"context"
	"fmt"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller"
	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"go.uber.org/atomic"
	"math/rand"
	"net"
	"sync"
	"time"
)

type msgBuffer struct {
	buffers    [][]byte
	first      time.Time
	decoded    bool
	superseded bool
}

var nilKey = make([]byte, 32)

func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	log.DEBUG("kopach miner starting")
	wg.Add(1)
	ciph := gcm.GetCipher(*cx.Config.MinerPass)
	var cancel context.CancelFunc
	var err error
	buffers := make(map[string]*msgBuffer)
	jobChan := make(chan controller.MinerContainer, *cx.Config.GenThreads)
	serverCounter := make(map[string]struct{})
	var serverCounterMx sync.Mutex
	var working atomic.Bool
	var latestHeight atomic.Uint32
	latestHeight.Store(0)
	var starter chan struct{}
	starter = make(chan struct{})
	//var lastHash atomic.Uint64
	for _, j := range controller.MCAddresses {
		i := j
		var mx sync.Mutex
		cancel, err = controller.Listen(i, func(a *net.UDPAddr, n int,
			b []byte) {
			mx.Lock()
			defer mx.Unlock()
			if n < 16 {
				log.ERROR("received short broadcast message")
				return
			}
			magic := string(b[12:16])
			if magic == string(controller.WorkMagic[:]) ||
				magic == string(controller.PauseMagic[:]) {
				nonce := string(b[:12])
				if bn, ok := buffers[nonce]; ok {
					if !bn.decoded {
						payload := b[16:n]
						newP := make([]byte, len(payload))
						copy(newP, payload)
						bn.buffers = append(bn.buffers, newP)
						if len(bn.buffers) >= 3 {
							// try to decode it
							var cipherText []byte
							//log.SPEW(bn.buffers)
							cipherText, err = fec.Decode(bn.buffers)
							if err != nil {
								log.ERROR(err)
								return
							}
							//log.SPEW(cipherText)
							msg, err := ciph.Open(nil, []byte(nonce), cipherText,
								nil)
							if err != nil {
								log.ERROR(err)
								return
							}
							bn.decoded = true
							//log.DEBUG(magic)
							switch magic {
							case string(controller.WorkMagic[:]):
								mC := controller.LoadMinerContainer(msg)
								for i := range buffers {
									if i != nonce {
										if buffers[i].superseded {
											log.DEBUGF("deleting buffer %x", i)
											delete(buffers, i)
										}
									}
								}
								if latestHeight.Load() == 0 {
									close(starter)
								}
								latestHeight.Store(uint32(mC.GetNewHeight()))
								for i := 0; i < *cx.Config.GenThreads; i++ {
									jobChan <- mC
								}
								// channels loaded,
								// enable mining if not already enabled
								log.DEBUG("signalling to start work")
								working.Store(true)
							case string(controller.PauseMagic[:]):
								pC := controller.LoadPauseContainer(msg)
								srvr := pC.GetIPs()
								var signature string
								for i := range srvr {
									signature += srvr[i].String()
								}
								serverCounterMx.Lock()
								delete(serverCounter, signature)
								if len(serverCounter) < 1 {
									// no currently active servers, pause work
									working.Store(false)
									log.WARN("pause message received and no" +
										" active servers on network, pausing")
								}
								serverCounterMx.Unlock()

							}
							//log.SPEW(msg)

						}
					} else {
						for i := range buffers {
							if i != nonce {
								// superseded blocks can be deleted from the
								// buffers,
								// we don't add more data for the already
								// decoded
								buffers[i].superseded = true
							}
						}
						fmt.Printf("received rebroadcast of %x %v\r", nonce, time.Now())
					}
				} else {
					buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(),
						false, false}
					payload := b[16:n]
					newP := make([]byte, len(payload))
					copy(newP, payload)
					buffers[nonce].buffers = append(buffers[nonce].buffers,
						newP)
					//log.DEBUGF("%x", payload)
				}
				//log.DEBUGF("%v %v %012x %s", i, a, nonce, magic)

			}
		})
		// end of handler
		if err != nil {
			continue
		}
		// we only need to start one of them,
		// ipv6 is preferred. if the same OS or any modern vanilla configured
		// system will have ipv6 unless the router doesn't support it,
		// but the enabling or disabling of it should be the same on workers
		// as the nodes,
		// unlikely most users will encounter this with the multicast protocol
		if cancel != nil {
			<-starter
			// 100 times per second we check whether to stop or start new work
			ticker := time.NewTicker(time.Millisecond * 10)
			working.Store(false)
			hashesPerAlgo := make(map[int32]*atomic.Uint64)
			// this atomic stores a list of ip addresss unique to each server
			// on the lan in order to make a threshold to listen to a work
			// pause message
			log.DEBUG("listener started", i.IP, i.Port, i.Zone, i.String(),
				i.Network())
			var rotator atomic.Uint64
			hf1Height := fork.List[1].ActivationHeight
			if fork.IsTestnet {
				hf1Height = fork.List[1].TestnetStart
			}
			currFork := fork.GetCurrent(int32(latestHeight.Load()))
			numPerAlgo := uint32(1 << 8)
			for i := 0; i < *cx.Config.GenThreads; i++ {
				loopCounter := 1
				// start the rolling algorithm cycle on a random starting point
				rand.Seed(time.Now().UnixNano())
				rotator.Store(uint64(rand.Intn(len(fork.List[fork.GetCurrent(
					int32(latestHeight.Load()))].AlgoVers))))
				go func(wrkr int, startup time.Time) {
					log.DEBUG("starting worker", wrkr)
					// each worker has its own copy so there is no races when
					// it is updated
					var mC controller.MinerContainer
					var blk *util.Block
					var header *wire.BlockHeader
					var rNonce, algoCount uint32
					rn, _ := wire.RandomUint64()
					rNonce = uint32(rn)
					aVs := fork.List[currFork].AlgoVers
					aVL := len(aVs)
				out:
					for {
						var targetDifficulty = &fork.SecondPowLimit
						if working.Load() {
							if int32(latestHeight.Load()) >= hf1Height {
								currFork = 1
								aVs = fork.List[currFork].AlgoVers
								aVL = len(aVs)
							}
							if loopCounter%50 == 0 {
								log.DEBUG("worker", wrkr, loopCounter,
									"rounds since startup, ",
									time.Now().Sub(startup)/time.
										Duration(loopCounter), "/ round")
							}
							loopCounter++
							algoCount++
							if algoCount > numPerAlgo {
								algoCount = 0
								rotator.Inc()
								var ver int32
								choice := int(rotator.Load()) % aVL
								if currFork == 0 {
									if choice == 0 {
										ver = 2
									} else {
										ver = 514
									}
								} else if currFork == 1 {
									ver = int32(choice + 5)
								}
								header.Version = ver
								header.Bits = mC.GetBitses()[ver]
								targetDifficulty = fork.CompactToBig(header.Bits)
							}
							//header = &blk.MsgBlock().Header
							// run a round of hashing
							//log.DEBUG("working")
							header.Nonce = rNonce
							// we just keep incrementing this number and let
							// it roll over
							rNonce++
							//log.SPEW(header)
							hash := header.BlockHashWithAlgos(int32(latestHeight.Load()))
							//log.DEBUG(hash[:32])
							if _, ok := hashesPerAlgo[header.Version]; !ok {
								if header.Version == 0 {
									panic("wtf")
								}
								hashesPerAlgo[header.Version] = &atomic.Uint64{}
								hashesPerAlgo[header.Version].Store(0)
							}
							hashesPerAlgo[header.Version].Inc()
							bigHash := blockchain.HashToBig(&hash)
							//log.DEBUGF("%064x", targetDifficulty)
							if bigHash.Cmp(targetDifficulty) <= 0 {
								log.WARN("worker", wrkr, "solution found",
									latestHeight.Load(), hash, header.Version)
								for i := range hashesPerAlgo {
									log.DEBUG(i, fork.List[fork.GetCurrent(
										int32(latestHeight.Load()))].
										AlgoVers[i], hashesPerAlgo[i].Load())
								}
								log.SPEW(header)
							}
						}
						select {
						case mC = <-jobChan:
							//close(starter)
							log.WARN("new job")
							// get the signature of the sender of the job so
							// when it sends a pause job it can be removed
							// from the serverCounter map
							srvr := mC.GetIPs()
							var signature string
							for i := range srvr {
								signature += srvr[i].String()
							}
							serverCounterMx.Lock()
							if _, ok := serverCounter[signature]; ok {
								break
							}
							serverCounterMx.Unlock()
							if wrkr == 0 {
								//log.DEBUG("P2PListenersPort",
								//	mC.GetP2PListenersPort())
								//log.DEBUG("RPCListenersPort",
								//	mC.GetRPCListenersPort())
								//log.DEBUG("ControllerListenerPort",
								//	mC.GetControllerListenerPort())
								//log.DEBUGF("h %d %v",
								//	mC.GetNewHeight(), mC.GetPrevBlockHash())
								//log.DEBUG(mC.GetBitses())
								//log.SPEW(mC.GetTxs())
							}
							// generate the msgblock for hashing
							txs := mC.GetTxs()
							rn, _ := wire.RandomUint64()
							var ver, cnt int32
							choice := int(rotator.Load()) % aVL
							for ii := range aVs {
								if int(cnt) == choice {
									ver = ii
									break
								}
								cnt++
							}
							log.DEBUG("choice", choice, ver)
							targetBits := mC.GetBitses()[ver]
							targetDifficulty = fork.CompactToBig(targetBits)
							//rotator.Inc()
							//log.SPEW(txs)
							blk = util.NewBlock(&wire.MsgBlock{
								Header: wire.BlockHeader{
									Version:   ver,
									PrevBlock: *mC.GetPrevBlockHash(),
									Timestamp: time.Now(),
									Bits:      targetBits,
									Nonce:     uint32(rn),
								},
								Transactions: txs,
							})
							header = &blk.MsgBlock().Header

							header.Version = ver
							header.Bits = targetBits
							targetDifficulty = fork.CompactToBig(targetBits)
							if _, ok := hashesPerAlgo[header.Version]; !ok {
								hashesPerAlgo[header.Version] = &atomic.Uint64{}
							}
							// use a random extra nonce to ensure no
							// duplicated work,
							// as well as for this case putting a merkle root
							// in there as that solely depends on the list of
							// transactions,
							// which come newly generated in every update (
							// new block accepted and transaction list change)
							enOffset, err := wire.RandomUint64()
							if err != nil {
								log.WARNF("unexpected error while generating"+
									" random extra nonce offset:", err)
								enOffset = 0
							}
							eN, _ := wire.RandomUint64()
							extraNonce := eN
							err = UpdateExtraNonce(blk.MsgBlock(),
								int32(latestHeight.Load()),
								extraNonce+enOffset)
							if err != nil {
								log.WARN(err)
							}
							//if wrkr == 0 {
							//log.SPEW(blk)
							//}
						case <-ticker.C:
							// only check whether to quit on the ticker
							select {
							case <-quit:
								break out
							default:
							}
						default:
							//log.DEBUG("spinning")
						}
					}
					log.DEBUG("worker", wrkr, "shutting down")
				}(i, time.Now())
			}
		} else {
			log.ERROR("failed to start listener on", i)
		}
		select {
		case <-quit:
			log.DEBUG("kopach miner shutting down")
			cancel()
		}
	}
	wg.Done()
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
	msgBlock.
		Transactions[0].
		TxIn[0].
		SignatureScript =
		coinbaseScript
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
