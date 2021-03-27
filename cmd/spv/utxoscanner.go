package spv

import (
	"container/heap"
	"github.com/p9c/pod/pkg/block"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/chainhash"
	am "github.com/p9c/pod/pkg/waddrmgr"
)

type (
	// GetUtxoRequest is a request to scan for InputWithScript from the height
	// BirthHeight.
	GetUtxoRequest struct {
		// Input is the target outpoint with script to watch for spentness.
		Input *InputWithScript
		// BirthHeight is the height at which we expect to find the original
		// unspent outpoint. This is also the height used when starting the search
		// for spends.
		BirthHeight uint32
		// resultChan either the spend report or error for this request.
		resultChan chan *getUtxoResult
		// result caches the first spend report or error returned for this
		// request.
		result *getUtxoResult
		// mu ensures the first response delivered via resultChan is in fact what
		// gets cached in result.
		mu   sync.Mutex
		quit qu.C
	}
	// A GetUtxoRequestPQ implements heap.Interface and holds GetUtxoRequests.
	// The queue maintains that heap. Pop() will always return the GetUtxo
	// request with the least starting height. This allows us to add new GetUtxo
	// requests to an already running batch.
	GetUtxoRequestPQ []*GetUtxoRequest
	// UtxoScanner batches calls to GetUtxo so that a single scan can search for
	// multiple outpoints. If a scan is in progress when a new element is added,
	// we check whether it can safely be added to the current batch, if not it
	// will be included in the next batch.
	UtxoScanner struct {
		started   uint32
		stopped   uint32
		cfg       *UtxoScannerConfig
		pq        GetUtxoRequestPQ
		nextBatch []*GetUtxoRequest
		mu        sync.Mutex
		cv        *sync.Cond
		wg        sync.WaitGroup
		quit      qu.C
		shutdown  qu.C
	}
	// UtxoScannerConfig exposes configurable methods for interacting with the
	// blockchain.
	UtxoScannerConfig struct {
		// BestSnapshot returns the block stamp of the current chain tip.
		BestSnapshot func() (*am.BlockStamp, error)
		// GetBlockHash returns the block hash at given height in main chain.
		GetBlockHash func(height int64) (*chainhash.Hash, error)
		// BlockFilterMatches checks the cfilter for the block hash for matches
		// against the rescan options.
		BlockFilterMatches func(ro *rescanOptions, blockHash *chainhash.Hash) (bool, error)
		// GetBlock fetches a block from the p2p network.
		GetBlock func(chainhash.Hash, ...QueryOption) (*block.Block, error)
	}
	// getUtxoResult is a simple pair type holding a spend report and error.
	getUtxoResult struct {
		report *SpendReport
		err    error
	}
)

// Result is callback returning either a spend report or an error.
func (r *GetUtxoRequest) Result(cancel qu.C) (*SpendReport, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case result := <-r.resultChan:
		// Cache the first result returned, in case we have multiple readers calling Result.
		if r.result == nil {
			r.result = result
		}
		return r.result.report, r.result.err
	case <-cancel.Wait():
		return nil, ErrGetUtxoCancelled
	case <-r.quit.Wait():
		return nil, ErrShuttingDown
	}
}

// deliver tries to deliver the report or error to any subscribers. If
// resultChan cannot accept a new update, this method will not block.
func (r *GetUtxoRequest) deliver(report *SpendReport, e error) {
	select {
	case r.resultChan <- &getUtxoResult{report, e}:
	default:
		W.F(
			"duplicate getutxo result delivered for outpoint=%v, spend=%v, err=%v",
			r.Input.OutPoint, report, e,
		)
	}
}

// IsEmpty returns true if the queue has no elements.
func (pq *GetUtxoRequestPQ) IsEmpty() bool {
	return pq.Len() == 0
}

// Peek returns the least height element in the queue without removing it.
func (pq *GetUtxoRequestPQ) Peek() *GetUtxoRequest {
	return (*pq)[0]
}

// Pop is called by the heap. Interface implementation to remove an element from
// the end of the backing store. The heap library will then maintain the heap
// invariant.
func (pq *GetUtxoRequestPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Push is called by the heap. Interface implementation to add an element to the
// end of the backing store. The heap library will then maintain the heap
// invariant.
func (pq *GetUtxoRequestPQ) Push(x interface{}) {
	item := x.(*GetUtxoRequest)
	*pq = append(*pq, item)
}

// Enqueue takes a GetUtxoRequest and adds it to the next applicable batch.
func (s *UtxoScanner) Enqueue(input *InputWithScript, birthHeight uint32) (*GetUtxoRequest, error) {
	D.F("enqueuing request for %s with birth height %d %s", input.OutPoint.String(), birthHeight)
	req := &GetUtxoRequest{
		Input:       input,
		BirthHeight: birthHeight,
		resultChan:  make(chan *getUtxoResult, 1),
		quit:        s.quit,
	}
	s.cv.L.Lock()
	select {
	case <-s.quit.Wait():
		s.cv.L.Unlock()
		return nil, ErrShuttingDown
	default:
	}
	// Insert the request into the queue and signal any threads that might be
	// waiting for new elements.
	heap.Push(&s.pq, req)
	s.cv.L.Unlock()
	s.cv.Signal()
	return req, nil
}

// Start begins running scan batches.
func (s *UtxoScanner) Start() (e error) {
	if !atomic.CompareAndSwapUint32(&s.started, 0, 1) {
		return nil
	}
	s.wg.Add(1)
	go s.batchManager()
	return nil
}

// Stop any in-progress scan.
func (s *UtxoScanner) Stop() (e error) {
	if !atomic.CompareAndSwapUint32(&s.stopped, 0, 1) {
		return nil
	}
	s.quit.Q()
batchShutdown:
	for {
		select {
		case <-s.shutdown.Wait():
			break batchShutdown
		case <-time.After(50 * time.Millisecond):
			s.cv.Signal()
		}
	}
	// Cancel all pending get utxo requests that were not pulled into the
	// batchManager's main goroutine.
	for !s.pq.IsEmpty() {
		pendingReq := heap.Pop(&s.pq).(*GetUtxoRequest)
		pendingReq.deliver(nil, ErrShuttingDown)
	}
	return nil
}

// batchManager is responsible for scheduling batches of UTXOs to scan. Any
// incoming requests whose start height has already been passed will be added to
// the next batch, which gets scheduled after the current batch finishes.
//
// NOTE: This method MUST be spawned as a goroutine.
func (s *UtxoScanner) batchManager() {
	defer s.shutdown.Q()
	for {
		s.cv.L.Lock()
		// Re-queue previously skipped requests for next batch.
		for _, request := range s.nextBatch {
			heap.Push(&s.pq, request)
		}
		s.nextBatch = nil
		// Wait for the queue to be non-empty.
		for s.pq.IsEmpty() {
			s.cv.Wait()
			select {
			case <-s.quit.Wait():
				s.cv.L.Unlock()
				return
			default:
			}
		}
		req := s.pq.Peek()
		s.cv.L.Unlock()
		// Break out now before starting a scan if a shutdown was requested.
		select {
		case <-s.quit.Wait():
			return
		default:
		}
		// Initiate a scan, starting from the birth height of the least-height
		// request currently in the queue.
		e := s.scanFromHeight(req.BirthHeight)
		if e != nil {
			E.F(
				"UXTO scan failed: %v", e,
			)
		}
	}
}

// dequeueAtHeight returns all GetUtxoRequests that have starting height of the
// given height.
func (s *UtxoScanner) dequeueAtHeight(height uint32) []*GetUtxoRequest {
	s.cv.L.Lock()
	defer s.cv.L.Unlock()
	// Take any requests that are too old to go in this batch and keep them for
	// the next batch.
	for !s.pq.IsEmpty() && s.pq.Peek().BirthHeight < height {
		item := heap.Pop(&s.pq).(*GetUtxoRequest)
		s.nextBatch = append(s.nextBatch, item)
	}
	var requests []*GetUtxoRequest
	for !s.pq.IsEmpty() && s.pq.Peek().BirthHeight == height {
		item := heap.Pop(&s.pq).(*GetUtxoRequest)
		requests = append(requests, item)
	}
	return requests
}

// scanFromHeight runs a single batch, pulling in any requests that get added
// above the batch's last processed height. If there was an error, then return
// the outstanding requests.
func (s *UtxoScanner) scanFromHeight(initHeight uint32) (e error) {
	// Before beginning the scan, grab the best block stamp we know of, which
	// will serve as an initial estimate for the end height of the scan.
	bestStamp, e := s.cfg.BestSnapshot()
	if e != nil {
		E.Ln(e)
		return e
	}
	var (
		// startHeight and endHeight bound the range of the current scan. If more
		// blocks are found while a scan is running, these values will be updated
		// afterwards to scan for the new blocks.
		startHeight = initHeight
		endHeight   = uint32(bestStamp.Height)
	)
	reporter := newBatchSpendReporter()
scanToEnd:
	// Scan forward through the blockchain and look for any transactions that
	// might spend the given UTXOs.
	for height := startHeight; height <= endHeight; height++ {
		// Before beginning to scan this height, check to see if the utxoscanner
		// has been signaled to exit.
		select {
		case <-s.quit.Wait():
			return reporter.FailRemaining(ErrShuttingDown)
		default:
		}
		var hash *chainhash.Hash
		hash, e = s.cfg.GetBlockHash(int64(height))
		if e != nil {
			E.Ln(e)
			return reporter.FailRemaining(e)
		}
		// If there are any new requests that can safely be added to this batch,
		// then try and fetch them.
		newReqs := s.dequeueAtHeight(height)
		// If an outpoint is created in this block, then fetch it regardless.
		// Otherwise check to see if the filter matches any of our watched
		// outpoints.
		fetch := len(newReqs) > 0
		if !fetch {
			options := rescanOptions{
				watchList: reporter.filterEntries,
			}
			var match bool
			match, e = s.cfg.BlockFilterMatches(&options, hash)
			if e != nil {
				E.Ln(e)
				return reporter.FailRemaining(e)
			}
			// If still no match is found, we have no reason to fetch this block,
			// and can continue to next height.
			if !match {
				continue
			}
		}
		// At this point, we've determined that we either (1) have new requests
		// which we need the block to scan for originating UTXOs, or (2) the
		// watchlist triggered a match against the neutrino filter. Before
		// fetching the block, check to see if the utxoscanner has been signaled
		// to exit so that we can exit the rescan before performing an expensive
		// operation.
		select {
		case <-s.quit.Wait():
			return reporter.FailRemaining(ErrShuttingDown)
		default:
		}
		T.F("fetching block height=%d hash=%s %s", height, hash)
		var blk *block.Block
		blk, e = s.cfg.GetBlock(*hash)
		if e != nil {
			E.Ln(e)
			return reporter.FailRemaining(e)
		}
		// Chk again to see if the utxoscanner has been signaled to exit.
		select {
		case <-s.quit.Wait():
			return reporter.FailRemaining(ErrShuttingDown)
		default:
		}
		D.F("processing block height=%d hash=%s %s", height, hash)
		reporter.ProcessBlock(blk.WireBlock(), newReqs, height)
	}
	// We've scanned up to the end height, now perform a check to see if we still
	// have any new blocks to process. If this is the first time through, we
	// might have a few blocks that were added since the scan started.
	currStamp, e := s.cfg.BestSnapshot()
	if e != nil {
		E.Ln(e)
		return reporter.FailRemaining(e)
	}
	// If the returned height is higher, we still have more blocks to go. Shift
	// the start and end heights and continue scanning.
	if uint32(currStamp.Height) > endHeight {
		startHeight = endHeight + 1
		endHeight = uint32(currStamp.Height)
		goto scanToEnd
	}
	reporter.NotifyUnspentAndUnfound()
	return nil
}

// We want Pop to give us the least BirthHeight.

func (pq GetUtxoRequestPQ) Len() int           { return len(pq) }
func (pq GetUtxoRequestPQ) Less(i, j int) bool { return pq[i].BirthHeight < pq[j].BirthHeight }
func (pq GetUtxoRequestPQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

// NewUtxoScanner creates a new instance of UtxoScanner using the given chain interface.
func NewUtxoScanner(cfg *UtxoScannerConfig) *UtxoScanner {
	scanner := &UtxoScanner{
		cfg:      cfg,
		quit:     qu.T(),
		shutdown: qu.T(),
	}
	scanner.cv = sync.NewCond(&scanner.mu)
	return scanner
}
