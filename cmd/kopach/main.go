package kopach

import "C"
import (
	"bytes"
	"context"
	"crypto/cipher"
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
	"math/big"
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

type Miner struct {
	buffers         map[string]*msgBuffer
	ciph            cipher.AEAD
	currFork        int
	cx              *conte.Xt
	hashesPerAlgo   map[int32]*atomic.Uint64
	hf1Height       int32
	jobChan         chan controller.MinerContainer
	latestHeight    *atomic.Uint32
	loopCounter     int
	mx              *sync.Mutex
	numPerAlgo      uint32
	quit            chan struct{}
	rotator         *atomic.Uint64
	serverCounter   map[string]struct{}
	serverCounterMx *sync.Mutex
	ticker          *time.Ticker
	working         *atomic.Bool
}

// Main the main thread of the kopach miner
func Main(cx *conte.Xt, quit chan struct{}, wg *sync.WaitGroup) {
	log.DEBUG("kopach miner starting")
	wg.Add(1)
	var cancel context.CancelFunc
	var err error
	m := &Miner{
		buffers:         make(map[string]*msgBuffer),
		ciph:            gcm.GetCipher(*cx.Config.MinerPass),
		currFork:        0,
		cx:              cx,
		hashesPerAlgo:   make(map[int32]*atomic.Uint64),
		hf1Height:       0,
		jobChan:         make(chan controller.MinerContainer, *cx.Config.GenThreads),
		latestHeight:    &atomic.Uint32{},
		loopCounter:     0,
		mx:              &sync.Mutex{},
		numPerAlgo:      uint32(1 << 8),
		quit:            nil,
		rotator:         &atomic.Uint64{},
		serverCounter:   make(map[string]struct{}),
		serverCounterMx: &sync.Mutex{},
		// 100 times per second we check whether to stop or start new work
		ticker:  time.NewTicker(time.Millisecond * 10),
		working: &atomic.Bool{},
	}
	//var lastHash atomic.Uint64
	for _, j := range controller.MCAddresses {
		i := j
		cancel, err = controller.Listen(i, getListener(m))
		if err != nil {
			continue
		}
		m.working.Store(false)
		// this atomic stores a list of ip addresss unique to each server
		// on the lan in order to make a threshold to listen to a work
		// pause message
		log.DEBUG("listener started", i.IP, i.Port, i.Zone, i.String(),
			i.Network())
		m.hf1Height = fork.List[1].ActivationHeight
		if fork.IsTestnet {
			m.hf1Height = fork.List[1].TestnetStart
		}
		m.currFork = fork.GetCurrent(int32(m.latestHeight.Load()))
		rand.Seed(time.Now().UnixNano())
		m.rotator.Store(uint64(rand.Intn(len(fork.List[fork.GetCurrent(
			int32(m.latestHeight.Load()))].AlgoVers))))
		for i := 0; i < *cx.Config.GenThreads; i++ {
			m.loopCounter = 1
			// start the rolling algorithm cycle on a random starting point
			go mine(m)(i, time.Now())
		}
	}
	select {
	case <-quit:
		log.DEBUG("kopach miner shutting down")
		cancel()
	}
	wg.Done()
}

func getListener(m *Miner) func(a *net.UDPAddr, n int, b []byte) {
	return func(a *net.UDPAddr, n int, b []byte) {
		var err error
		m.mx.Lock()
		defer m.mx.Unlock()
		if n < 16 {
			log.ERROR("received short broadcast message")
			return
		}
		magic := string(b[12:16])
		if magic == string(controller.WorkMagic[:]) ||
			magic == string(controller.PauseMagic[:]) {
			nonce := string(b[:12])
			if bn, ok := m.buffers[nonce]; ok {
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
						msg, err := (m.ciph).Open(nil, []byte(nonce),
							cipherText, nil)
						if err != nil {
							log.ERROR(err)
							return
						}
						bn.decoded = true
						//log.DEBUG(magic)
						switch magic {
						case string(controller.WorkMagic[:]):
							mC := controller.LoadMinerContainer(msg)
							for i := range m.buffers {
								if i != nonce {
									if m.buffers[i].superseded {
										log.DEBUGF("deleting buffer %x", i)
										delete(m.buffers, i)
									}
								}
							}

							m.latestHeight.Store(uint32(mC.GetNewHeight()))
							for i := 0; i < *m.cx.Config.GenThreads; i++ {
								m.jobChan <- mC
							}
							// channels loaded,
							// enable mining if not already enabled
							log.DEBUG("signalling to start work")
							m.working.Store(true)
						case string(controller.PauseMagic[:]):
							pC := controller.LoadPauseContainer(msg)
							srvr := pC.GetIPs()
							var signature string
							for i := range srvr {
								signature += srvr[i].String()
							}
							m.serverCounterMx.Lock()
							delete(m.serverCounter, signature)
							if len(m.serverCounter) < 1 {
								// no currently active servers, pause work
								m.working.Store(false)
								log.WARN("pause message received and no" +
									" active servers on network, pausing")
							}
							m.serverCounterMx.Unlock()

						}
						//log.SPEW(msg)

					}
				} else {
					for i := range m.buffers {
						if i != nonce {
							// superseded blocks can be deleted from the
							// buffers,
							// we don't add more data for the already
							// decoded
							m.buffers[i].superseded = true
						}
					}
					fmt.Printf("received rebroadcast of %x %v\r", nonce, time.Now())
				}
			} else {
				m.buffers[nonce] = &msgBuffer{[][]byte{}, time.Now(),
					false, false}
				payload := b[16:n]
				newP := make([]byte, len(payload))
				copy(newP, payload)
				m.buffers[nonce].buffers = append(m.buffers[nonce].buffers,
					newP)
				//log.DEBUGF("%x", payload)
			}
			//log.DEBUGF("%v %v %012x %s", i, a, nonce, magic)

		}
	}
}

func mine(m *Miner) func(wrkr int, startup time.Time) {
	return func(wrkr int, startup time.Time) {
		// each worker has its own copy so there is no races when
		// it is updated
		var mC controller.MinerContainer
		var blk *util.Block
		var header *wire.BlockHeader
		var rNonce, algoCount uint32
		var targetDifficulty *big.Int
		rn, _ := wire.RandomUint64()
		rNonce = uint32(rn)
		aVs := fork.List[m.currFork].AlgoVers
		aVL := len(aVs)
		log.DEBUG("starting worker", wrkr)
	out:
		for {
			if m.working.Load() && header != nil {
				targetDifficulty = &fork.SecondPowLimit
				if int32(m.latestHeight.Load()) >= m.hf1Height {
					m.currFork = 1
					aVs = fork.List[m.currFork].AlgoVers
					aVL = len(aVs)
				}
				if m.loopCounter%50 == 0 {
					log.DEBUG("worker", wrkr, m.loopCounter,
						"rounds since startup, ",
						time.Now().Sub(startup)/time.
							Duration(m.loopCounter), "/ round")
				}
				m.loopCounter++
				algoCount++
				if algoCount > m.numPerAlgo {
					algoCount = 0
					m.rotator.Inc()
					var ver int32
					choice := int(m.rotator.Load()) % aVL
					if m.currFork == 0 {
						if choice == 0 {
							ver = 2
						} else {
							ver = 514
						}
					} else if m.currFork == 1 {
						ver = int32(choice + 5)
					}
					header.Version = ver
					header.Bits = mC.GetBitses()[ver]
					targetDifficulty = fork.CompactToBig(header.Bits)
				}
				// run a round of hashing
				header.Nonce = rNonce
				// we just keep incrementing this number and let
				// it roll over
				rNonce++
				hash := header.BlockHashWithAlgos(int32(m.latestHeight.Load()))
				if _, ok := m.hashesPerAlgo[header.Version]; !ok {
					if header.Version == 0 {
						panic("wtf")
					}
					m.hashesPerAlgo[header.Version] = &atomic.Uint64{}
					m.hashesPerAlgo[header.Version].Store(0)
				}
				m.hashesPerAlgo[header.Version].Inc()
				bigHash := blockchain.HashToBig(&hash)
				if bigHash.Cmp(targetDifficulty) <= 0 {
					log.WARN("worker", wrkr, "solution found",
						m.latestHeight.Load(), hash, header.Version)
					// now we should stop mining
					m.working.Store(false)
					// construct a message to submit the solution
					var buffer bytes.Buffer
					err := blk.MsgBlock().Serialize(&buffer)
					if err != nil {
						log.ERROR(err)
					}
					ips := mC.GetIPs()
					shards, err := controller.Shards(buffer.
						Bytes(), controller.SolutionMagic, m.ciph)
					if err != nil {
						log.ERROR(err)
						break out
					}
					log.DEBUG(ips)
					port := mC.GetControllerListenerPort()
					var ipA *net.UDPAddr
					for k := range ips {
						host := ips[k].String()
						ipA = &net.UDPAddr{IP: net.ParseIP(host),
							Port: int(port)}
						log.SPEW(ipA)
						//conn, err = net.ListenUDP("udp",
						//	ipA)
						//if err != nil {
						//	log.ERROR(err)
						//	continue
						//}
						var n, cumulative int
						for i := range shards {
							conn, err := net.DialUDP("udp", nil, ipA)
							n, err = conn.Write(shards[i])
							if err != nil {
								log.ERROR(err)
								return
							}
							cumulative += n
							err = conn.Close()
							if err != nil {
								log.ERROR(err)
								continue
							}
						}
						log.DEBUGF("sent %v bytes to %v port %v %v",
							cumulative, ipA.IP, ipA.Port, time.Now())
						//err = controller.SendShards(ipA, shards, conn)
						//if err != nil {
						//	log.ERROR(err)
						//	continue
						//}
					}
				}
			}
			select {
			case mC = <-m.jobChan:
				log.WARN("new job")
				// get the signature of the sender of the job so
				// when it sends a pause job it can be removed
				// from the serverCounter map
				srvr := mC.GetIPs()
				var signature string
				for i := range srvr {
					signature += srvr[i].String()
				}
				m.serverCounterMx.Lock()
				if _, ok := m.serverCounter[signature]; ok {
					break
				}
				m.serverCounterMx.Unlock()
				// generate the msgblock for hashing
				txs := mC.GetTxs()
				rn, _ := wire.RandomUint64()
				var ver, cnt int32
				choice := int(m.rotator.Load()) % aVL
				for ii := range aVs {
					if int(cnt) == choice {
						ver = ii
						break
					}
					cnt++
				}
				targetBits := mC.GetBitses()[ver]
				targetDifficulty = fork.CompactToBig(targetBits)
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
				if _, ok := m.hashesPerAlgo[header.Version]; !ok {
					m.hashesPerAlgo[header.Version] = &atomic.Uint64{}
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
					int32(m.latestHeight.Load()),
					extraNonce+enOffset)
				if err != nil {
					log.WARN(err)
				}
				m.working.Store(true)
			case <-m.ticker.C:
				// only check whether to quit on the ticker
				select {
				case <-m.quit:
					log.DEBUG("worker", wrkr, "shutting down")
					break out
				default:
				}
			default:
			}
		}
	}
}

// UpdateExtraNonce updates the extra nonce in the coinbase script of the
// passed block by regenerating the coinbase script with the passed value and
// block height.  It also recalculates and updates the new merkle root that
// results from changing the coinbase script.
func UpdateExtraNonce(msgBlock *wire.MsgBlock, blockHeight int32,
	extraNonce uint64) error {
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
