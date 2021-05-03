package cpuminer

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"go.uber.org/atomic"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/util"
)

var tn = time.Now()

// CPUMiner provides facilities for solving blocks (mining) using the CPU in a
// concurrency-safe manner.  It consists of two main goroutines -- a speed
// monitor and a controller for worker goroutines which generate and solve
// blocks.  The number of goroutines can be set via the SetMaxGoRoutines
// function, but the default is based on the number of processor cores in the
// system which is typically sufficient.
type CPUMiner struct {
	sync.Mutex
	b                 *blockchain.BlockChain
	g                 *mining.BlkTmplGenerator
	cfg               Config
	numWorkers        uint32
	started           bool
	discreteMining    bool
	submitBlockLock   sync.Mutex
	wg                sync.WaitGroup
	workerWg          sync.WaitGroup
	updateNumWorkers  qu.C
	queryHashesPerSec chan float64
	updateHashes      chan uint64
	speedMonitorQuit  qu.C
	quit              qu.C
	rotator           atomic.Uint64
}

// Config is a descriptor containing the cpu miner configuration.
type Config struct {
	// Blockchain gives access for the miner to information about the chain
	//
	Blockchain *blockchain.BlockChain
	// ChainParams identifies which chain parameters the cpu miner is associated
	// with.
	ChainParams *netparams.Params
	// BlockTemplateGenerator identifies the instance to use in order to
	// generate block templates that the miner will attempt to solve.
	BlockTemplateGenerator *mining.BlkTmplGenerator
	// MiningAddrs is a list of payment addresses to use for the generated
	// blocks.  Each generated block will randomly choose one of them.
	MiningAddrs []util.Address
	// ProcessBlock defines the function to call with any solved blocks. It
	// typically must run the provided block through the same set of rules and
	// handling as any other block coming from the network.
	ProcessBlock func(*util.Block, blockchain.BehaviorFlags) (bool, error)
	// ConnectedCount defines the function to use to obtain how many other peers
	// the server is connected to.  This is used by the automatic persistent
	// mining routine to determine whether or it should attempt mining.  This is
	// useful because there is no point in mining when not connected to any
	// peers since there would no be anyone to send any found blocks to.
	ConnectedCount func() int32
	// IsCurrent defines the function to use to obtain whether or not the block
	// chain is current.  This is used by the automatic persistent mining
	// routine to determine whether or it should attempt mining. This is useful
	// because there is no point in mining if the chain is not current since any
	// solved blocks would be on a side chain and and up orphaned anyways.
	IsCurrent func() bool
	// Algo is the name of the type of PoW used for the block header.
	Algo string
	// NumThreads is the number of threads set in the configuration for the
	// CPUMiner
	NumThreads uint32
	// Solo sets whether the miner will run when not connected
	Solo bool
}

const (
	// // maxNonce is the maximum value a nonce can be in a block header.
	// maxNonce = 2 ^ 32 - 1
	// maxExtraNonce is the maximum value an extra nonce used in a coinbase
	// transaction can be.
	maxExtraNonce = 2 ^ 64 - 1
	// hpsUpdateSecs is the number of seconds to wait in between each update to the hashes per second monitor.
	hpsUpdateSecs = 1
	// hashUpdateSec is the number of seconds each worker waits in between
	// notifying the speed monitor with how many hashes have been completed
	// while they are actively searching for a solution.  This is done to reduce
	// the amount of syncs between the workers that must be done to keep track
	// of the hashes per second.
	hashUpdateSecs = 9
)

var (
	// defaultNumWorkers is the default number of workers to use for mining and
	// is based on the number of processor cores.  This helps ensure the system
	// stays reasonably responsive under heavy load.
	defaultNumWorkers = uint32(runtime.NumCPU())
)

// GenerateNBlocks generates the requested number of blocks. It is self
// contained in that it creates block templates and attempts to solve them
// while detecting when it is performing stale work and reacting accordingly by
// generating a new block template.  When a block is solved, it is submitted.
// The function returns a list of the hashes of generated blocks.
func (m *CPUMiner) GenerateNBlocks(workerNumber uint32, n uint32,
	algo string) ([]*chainhash.Hash, error) {
	m.Lock()
	L.Warnf("generating %s blocks...", m.cfg.Algo)
	// Respond with an error if server is already mining.
	if m.started || m.discreteMining {
		m.Unlock()
		return nil, errors.New("server is already CPU mining; call " +
			"`setgenerate 0` before calling discrete `generate` commands")
	}
	m.started = true
	m.discreteMining = true
	m.speedMonitorQuit = qu.Ter()
	m.wg.Add(1)
	go m.speedMonitor()
	m.Unlock()
	L.Warnf("generating %d blocks", n)
	i := uint32(0)
	blockHashes := make([]*chainhash.Hash, n)
	// Start a ticker which is used to signal checks for stale work and updates to the speed monitor.
	ticker := time.NewTicker(time.Second * hashUpdateSecs)
	defer ticker.Stop()
	for {
		// Read updateNumWorkers in case someone tries a `setgenerate` while
		// we're generating. We can ignore it as the `generate` RPC call only
		// uses 1 worker.
		select {
		case <-m.updateNumWorkers:
		default:
		}
		// Grab the lock used for block submission, since the current block will
		// be changing and this would otherwise end up building a new block
		// template on a block that is in the process of becoming stale.
		m.submitBlockLock.Lock()
		curHeight := m.g.BestSnapshot().Height
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		payToAddr := m.cfg.MiningAddrs[rand.Intn(len(m.cfg.MiningAddrs))]
		// Create a new block template using the available transactions in the
		// memory pool as a source of transactions to potentially include in the
		// block.
		template, e := m.g.NewBlockTemplate(workerNumber, payToAddr, algo)
		m.submitBlockLock.Unlock()
		if e != nil  {
			L.			L.Warnf("failed to create new block template:", err)
			continue
		}
		// Attempt to solve the block.  The function will exit early with false
		// when conditions that trigger a stale block, so a new block template
		// can be generated.  When the return is true a solution was found, so
		// submit the solved block.
		if m.solveBlock(workerNumber, template.Block, curHeight+1,
			m.cfg.ChainParams.Name == "testnet", ticker, nil) {
			block := util.NewBlock(template.Block)
			m.submitBlock(block)
			blockHashes[i] = block.Hash()
			i++
			if i == n {
				L.Warnf("generated %d blocks", i)
				m.Lock()
				close(m.speedMonitorQuit)
				m.wg.Wait()
				m.started = false
				m.discreteMining = false
				m.Unlock()
				return blockHashes, nil
			}
		}
	}
}

// GetAlgo returns the algorithm currently configured for the miner
func (m *CPUMiner) GetAlgo() (name string) {
	return m.cfg.Algo
}

// HashesPerSecond returns the number of hashes per second the mining process
// is performing.  0 is returned if the miner is not currently running. This
// function is safe for concurrent access.
func (m *CPUMiner) HashesPerSecond() float64 {
	m.Lock()
	defer m.Unlock()
	// Nothing to do if the miner is not currently running.
	if !m.started {
		return 0
	}
	return <-m.queryHashesPerSec
}

// IsMining returns whether or not the CPU miner has been started and is
// therefore currently mining. This function is safe for concurrent access.
func (m *CPUMiner) IsMining() bool {
	m.Lock()
	defer m.Unlock()
	return m.started
}

// NumWorkers returns the number of workers which are running to solve blocks.
// This function is safe for concurrent access.
func (m *CPUMiner) NumWorkers() int32 {
	m.Lock()
	defer m.Unlock()
	return int32(m.numWorkers)
}

// SetAlgo sets the algorithm for the CPU miner
func (m *CPUMiner) SetAlgo(
	name string) {
	m.cfg.Algo = name
}

// SetNumWorkers sets the number of workers to create which solve blocks.  Any
// negative values will cause a default number of workers to be used which is
// based on the number of processor cores in the system.  A value of 0 will
// cause all CPU mining to be stopped. This function is safe for concurrent
// access.
func (m *CPUMiner) SetNumWorkers(numWorkers int32) {
	if numWorkers == 0 {
		m.Stop()
	}
	// Don't lock until after the first check since Stop does its own locking.
	m.Lock()
	defer m.Unlock()
	// Use default if provided value is negative.
	if numWorkers < 0 {
		m.numWorkers = defaultNumWorkers
	} else {
		m.numWorkers = uint32(numWorkers)
	}
	// When the miner is already running, notify the controller about the the change.
	if m.started {
		m.updateNumWorkers <- struct{}{}
	}
}

// Start begins the CPU mining process as well as the speed monitor used to
// track hashing metrics.  Calling this function when the CPU miner has already
// been started will have no effect.
// This function is safe for concurrent access.
func (m *CPUMiner) Start() {
	if len(m.cfg.MiningAddrs) < 1 {
		return
	}
	m.Lock()
	defer m.Unlock()
	L.inf.Ln("starting cpu miner")
	// Nothing to do if the miner is already running or if running in discrete mode (using GenerateNBlocks).
	if m.started || m.discreteMining {
		return
	}
	// randomize the starting point so all network is mining different
	m.rotator.Store(uint64(rand.Intn(
		len(fork.List[fork.GetCurrent(m.b.BestSnapshot().Height)].Algos))))
	m.quit = qu.Ter()
	m.speedMonitorQuit = qu.Ter()
	m.wg.Add(2)
	go m.speedMonitor()
	go m.miningWorkerController()
	m.started = true
	L.trc.Ln("CPU miner started mining", m.cfg.Algo, m.cfg.NumThreads)
}

// Stop gracefully stops the mining process by signalling all workers, and the
// speed monitor to quit.  Calling this function when the CPU miner has not
// already been started will have no effect. This function is safe for
// concurrent access.
func (m *CPUMiner) Stop() {
	m.Lock()
	defer m.Unlock()
	// Nothing to do if the miner is not currently running or if running in discrete mode (using GenerateNBlocks).
	if !m.started || m.discreteMining {
		return
	}
	close(m.quit)
	m.wg.Wait()
	m.started = false
	L.wrn.Ln("CPU miner stopped")
}

// generateBlocks is a worker that is controlled by the miningWorkerController.
// It is self contained in that it creates block templates and attempts to
// solve them while detecting when it is performing stale work and reacting
// accordingly by generating a new block template.  When a block is solved, it
// is submitted. It must be run as a goroutine.
func (m *CPUMiner) generateBlocks(workerNumber uint32, quit qu.C) {
	// Start a ticker which is used to signal checks for stale work and updates
	// to the speed monitor.
	ticker := time.NewTicker(time.Second / 3) // * hashUpdateSecs)
	defer ticker.Stop()
out:
	for i := 0; ; i++ {
		//L.trc.Ln(workerNumber, "generateBlocksLoop start")
		// Quit when the miner is stopped.
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}
		// Wait until there is a connection to at least one other peer since
		// there is no way to relay a found block or receive transactions to work
		// on when there are no connected peers.
		if (m.cfg.ConnectedCount() == 0 || !m.cfg.IsCurrent()) &&
			(m.cfg.ChainParams.Net == wire.MainNet && !m.cfg.Solo) {
			log.Print(log.Composite("server has no peers, waiting...",
				"STATUS", true), "\r")
			time.Sleep(time.Second)
			continue
		}
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}
		// No point in searching for a solution before the chain is synced.  Also,
		// grab the same lock as used for block submission, since the current
		// block will be changing and this would otherwise end up building a new
		// block template on a block that is in the process of becoming stale.
		curHeight := m.g.BestSnapshot().Height
		if curHeight != 0 && !m.cfg.IsCurrent() && !m.cfg.Solo {
			L.Warnf("server is not current yet, waiting")
			time.Sleep(time.Second)
			continue
		}
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}
		// choose the algorithm on a rolling cycle
		counter := m.rotator.Load()
		//counter /= uint64(len(fork.List[fork.GetCurrent(curHeight+1)].Algos))*2
		m.rotator.Add(1)
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
			// L.wrn.Ln("algo", algo)
		}
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		payToAddr := m.cfg.MiningAddrs[rand.Intn(len(m.cfg.MiningAddrs))]
		// Create a new block template using the available transactions in the
		// memory pool as a source of transactions to potentially include in the
		// block.
		template, e := m.g.NewBlockTemplate(workerNumber, payToAddr, algo)
		if e != nil  {
			L.wrn.Ln("failed to create new block template:", err)
			continue
		}
		// Attempt to solve the block.  The function will exit early with false
		// when conditions that trigger a stale block, so a new block template
		// can be generated.  When the return is true a solution was found, so
		// submit the solved block.
		//L.trc.Ln("attempting to solve block")
		select {
		case <-quit:
			break out
		default:
			// Non-blocking select to fall through
		}
		if m.solveBlock(workerNumber, template.Block, curHeight+1,
			m.cfg.ChainParams.Name == "testnet", ticker, quit) {
			block := util.NewBlock(template.Block)
			m.submitBlock(block)
		}
	}
	L.trc.Ln("cpu miner worker finished")
	m.workerWg.Done()
}

// miningWorkerController launches the worker goroutines that are used to
// generate block templates and solve them.  It also provides the ability to
// dynamically adjust the number of running worker goroutines. It must be run
// as a goroutine.
func (m *CPUMiner) miningWorkerController() {
	L.trc.Ln("starting mining worker controller")
	// launchWorkers groups common code to launch a specified number of workers
	// for generating blocks.
	var runningWorkers qu.Ters
	launchWorkers := func(numWorkers uint32) {
		for i := uint32(0); i < numWorkers; i++ {
			quit := qu.Ter()
			runningWorkers = append(runningWorkers, quit)
			m.workerWg.Add(1)
			go m.generateBlocks(i, quit)
		}
	}
	L.Tracef("spawning %d worker(s)", m.numWorkers)
	// Launch the current number of workers by default.
	runningWorkers = make(qu.Ters, 0, m.numWorkers)
	launchWorkers(m.numWorkers)
out:
	for {
		select {
		// Update the number of running workers.
		case <-m.updateNumWorkers:
			// No change.
			numRunning := uint32(len(runningWorkers))
			if m.numWorkers == numRunning {
				continue
			}
			// Add new workers.
			if m.numWorkers > numRunning {
				launchWorkers(m.numWorkers - numRunning)
				continue
			}
			// Signal the most recently created goroutines to exit.
			for i := numRunning - 1; i >= m.numWorkers; i-- {
				close(runningWorkers[i])
				runningWorkers[i] = nil
				runningWorkers = runningWorkers[:i]
			}
		case <-m.quit:
			for _, quit := range runningWorkers {
				close(quit)
			}
			break out
		}
	}
	// Wait until all workers shut down to stop the speed monitor since they rely on being able to send updates to it.
	m.workerWg.Wait()
	close(m.speedMonitorQuit)
	m.wg.Done()
}

// solveBlock attempts to find some combination of a nonce, extra nonce, and
// current timestamp which makes the passed block hash to a value less than the
// target difficulty.  The timestamp is updated periodically and the passed
// block is modified with all tweaks during this process.  This means that when
// the function returns true, the block is ready for submission. This function
// will return early with false when conditions that trigger a stale block such
// as a new block showing up or periodically when there are new transactions
// and enough time has elapsed without finding a solution.
func (m *CPUMiner) solveBlock(workerNumber uint32, msgBlock *wire.MsgBlock,
	blockHeight int32, testnet bool, ticker *time.Ticker,
	quit qu.C) bool {
	//L.trc.Ln("running solveBlock")
	// algoName := fork.GetAlgoName(
	// 	msgBlock.Header.Version, m.b.BestSnapshot().Height)
	// Choose a random extra nonce offset for this block template and worker.
	enOffset, e := wire.RandomUint64()
	if e != nil  {
		L.		L.Warnf("unexpected error while generating random extra nonce"+
			" offset:",
			err)
		enOffset = 0
	}
	// Create some convenience variables.
	header := &msgBlock.Header
	targetDifficulty := fork.CompactToBig(header.Bits)
	// Initial state.
	lastGenerated := time.Now()
	lastTxUpdate := m.g.GetTxSource().LastUpdated()
	hashesCompleted := uint64(0)
	// Note that the entire extra nonce range is iterated and the offset is
	// added relying on the fact that overflow will wrap around 0 as provided by
	// the Go spec.
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
		//L.trc.Ln("updating extraNonce")
		e := m.g.UpdateExtraNonce(msgBlock, blockHeight, extraNonce+enOffset)
		if e != nil  {
			L.		}
		// Search through the entire nonce range for a solution while
		// periodically checking for early quit and stale block conditions along
		// with updates to the speed monitor.
		var shifter uint64 = 16
		if testnet {
			shifter = 16
		}
		rn, _ := wire.RandomUint64()
		if rn > 1<<63-1<<shifter {
			rn -= 1 << shifter
		}
		rn += 1 << shifter
		rNonce := uint32(rn)
		mn := uint32(1 << 8)
		switch {
		case m.cfg.NumThreads < 2:
			mn = uint32(1 << 7)
		case m.cfg.NumThreads < 4:
			mn = uint32(1 << 8)
		case m.cfg.NumThreads < 6:
			mn = uint32(1 << 9)
		case m.cfg.NumThreads < 8:
			mn = uint32(1 << 10)
		}
		// if testnet {
		// 	mn = 1 << shifter
		// }
		if fork.GetCurrent(blockHeight) == 0 {
			mn = 1 << 16 * m.cfg.NumThreads
		}
		var i uint32
		algo := fork.GetAlgoName(msgBlock.Header.Version,
			blockHeight)
		defer func() {
			L.Tracef("wrkr %d finished %d rounds of %s", workerNumber,
				i-rNonce, algo)
		}()
		L.trc.Ln("starting round from ", rNonce, algo, mn)
		for i = rNonce; i <= rNonce+mn; i++ {
			// if time.Now().Sub(now) > time.Second*3 {
			// 	return false
			// }
			select {
			case <-quit:
				return false
			case <-ticker.C:
				m.updateHashes <- hashesCompleted
				hashesCompleted = 0
				// The current block is stale if the best block has changed.
				best := m.g.BestSnapshot()
				if !header.PrevBlock.IsEqual(&best.Hash) {
					return false
				}
				// The current block is stale if the memory pool has been updated
				// since the block template was generated and it has been at least
				// one minute.
				if lastTxUpdate != m.g.GetTxSource().LastUpdated() &&
					time.Now().After(lastGenerated.Add(time.Minute)) {
					return false
				}
				e := m.g.UpdateBlockTime(workerNumber, msgBlock)
				if e != nil  {
					L.				}
			default:
			}
			var incr uint64 = 1
			header.Nonce = i
			hash := header.BlockHashWithAlgos(blockHeight)
			hashesCompleted += incr
			// The block is solved when the new block hash is less than the target
			// difficulty.  Yay!
			bigHash := blockchain.HashToBig(&hash)
			if bigHash.Cmp(targetDifficulty) <= 0 {
				m.updateHashes <- hashesCompleted
				return true
			}
		}
		return false
	}
	return false
}

// speedMonitor handles tracking the number of hashes per second the mining
// process is performing.  It must be run as a goroutine.
func (m *CPUMiner) speedMonitor() {
	var hashesPerSec float64
	var totalHashes uint64 = 1
	ticker := time.NewTicker(time.Second * hpsUpdateSecs)
	defer ticker.Stop()
out:
	for i := 0; ; i++ {
		select {
		// Periodic updates from the workers with how many hashes they have
		// performed.
		case numHashes := <-m.updateHashes:
			totalHashes += numHashes
		// Time to update the hashes per second.
		case <-ticker.C:
			//curHashesPerSec := float64(totalHashes) / hpsUpdateSecs
			//if hashesPerSec == 0 {
			//	hashesPerSec = curHashesPerSec
			//}
			//hashesPerSec = (hashesPerSec + curHashesPerSec) / 2
			////totalHashes = 0
			//if hashesPerSec != 0 {
			//since := fmt.Sprint(time.Now().Sub(log.StartupTime) / time.
			//	Second * time.Second)
			since := uint64(time.Now().Sub(tn)/time.Second) + 1
			log.Print(log.Composite(fmt.Sprintf(
				"--> Hash speed: %d hash/s av since start",
				totalHashes/since),
				"STATUS", true),
				"\r")
			//}
		// Request for the number of hashes per second.
		case m.queryHashesPerSec <- hashesPerSec:
			// Nothing to do.
		case <-m.speedMonitorQuit:
			break out
		}
	}
	m.wg.Done()
}

// submitBlock submits the passed block to network after ensuring it passes all
// of the consensus validation rules.
func (m *CPUMiner) submitBlock(block *util.Block) bool {
	L.trc.Ln("submitting block")
	m.submitBlockLock.Lock()
	defer m.submitBlockLock.Unlock()
	// TODO: This nonsense and the workgroup are the entirely wrong way to
	//  write a cryptocurrency miner.
	//  It is not critical work so it should not use a workgroup but rather,
	//  a semaphore, which can yank all the threads to stop as soon as they
	//  select on the semaphore,
	//  whereas this stops 'when all the jobs are finished' for what purpose?
	//  This miner will be eliminated once the replacement is complete.
	//  End result of this is node waits for miners to stop,
	//  which sometimes takes 5 seconds and almost every other height it has
	//  two submissions processing on one mutex and second and others are
	//  always stale. So also, the submitlock needs to be revised,
	//  I think submitlock and waitgroup together are far better replaced by
	//  a semaphore. Miner should STOP DEAD when a solution is found and wait
	//  for more work. The kopach controller will stop all miners in the
	//  network when it receives submissions prejudicially because it is
	//  better to save power than catch one block in a thousand from an
	//  economics poinnt of view.
	//  Every cycle degrades the value and brings closer the hardware failure
	//  so don't work unless there is a very good reason to.
	// Ensure the block is not stale since a new block could have shown up while
	// the solution was being found.  Typically that condition is detected and
	// all work on the stale block is halted to start work on a new block, but
	// the check only happens periodically, so it is possible a block was found
	// and submitted in between.
	msgBlock := block.MsgBlock()
	if !msgBlock.Header.PrevBlock.IsEqual(&m.g.BestSnapshot().Hash) {
		L.trc.Ln(
			"Block submitted via CPU miner with previous block", msgBlock.Header.PrevBlock,
			"is stale", msgBlock.Header.Version,
			msgBlock.BlockHashWithAlgos(block.Height()))
		return false
	}
	//L.trc.Ln("found block is fresh ", m.cfg.ProcessBlock)
	// Process this block using the same rules as blocks coming from other
	// nodes.  This will in turn relay it to the network like normal.
	isOrphan, e := m.cfg.ProcessBlock(block, blockchain.BFNone)
	if e != nil  {
		L.		// Anything other than a rule violation is an unexpected error, so log
		// that error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			L.Warnf(
				"Unexpected error while processing block submitted via CPU miner:", err,
			)
			return false
		}
		L.Warnf("block submitted via CPU miner rejected:", err)
		return false
	}
	if isOrphan {
		L.wrn.Ln("block is an orphan")
		return false
	}
	//L.trc.Ln("the block was accepted")
	coinbaseTx := block.MsgBlock().Transactions[0].TxOut[0]
	prevHeight := block.Height() - 1
	prevBlock, _ := m.b.BlockByHeight(prevHeight)
	prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
	since := block.MsgBlock().Header.Timestamp.Unix() - prevTime
	bHash := block.MsgBlock().BlockHashWithAlgos(block.Height())
	L.Warnf("new block height %d %08x %s%10d %08x %v %s %ds since prev",
		block.Height(),
		prevBlock.MsgBlock().Header.Bits,
		bHash,
		block.MsgBlock().Header.Timestamp.Unix(),
		block.MsgBlock().Header.Bits,
		util.Amount(coinbaseTx.Value),
		fork.GetAlgoName(block.MsgBlock().Header.Version, block.Height()),
		since)
	return true
}

// New returns a new instance of a CPU miner for the provided configuration.
// Use Start to begin the mining process.  See the documentation for CPUMiner
// type for more details.
func New(cfg *Config) *CPUMiner {
	return &CPUMiner{
		b:                 cfg.Blockchain,
		g:                 cfg.BlockTemplateGenerator,
		cfg:               *cfg,
		numWorkers:        cfg.NumThreads,
		updateNumWorkers:  qu.Ter(),
		queryHashesPerSec: make(chan float64),
		updateHashes:      make(chan uint64),
	}
}
