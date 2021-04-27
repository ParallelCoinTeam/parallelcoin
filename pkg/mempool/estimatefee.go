package mempool

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/block"
	"io"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/mining"
	"github.com/p9c/matrjoska/pkg/util"
)

// DUOPerKilobyte is number with units of parallelcoins per kilobyte.
type DUOPerKilobyte float64

// FeeEstimator manages the data necessary to create fee estimations. It is safe for concurrent access.
type FeeEstimator struct {
	maxRollback uint32
	binSize     int32
	// The maximum number of replacements that can be made in a single bin per block. Default is
	// estimateFeeMaxReplacements
	maxReplacements int32
	// The minimum number of blocks that can be registered with the fee estimator before it will provide answers.
	minRegisteredBlocks uint32
	// The last known height.
	lastKnownHeight int32
	// The number of blocks that have been registered.
	numBlocksRegistered uint32
	mtx                 sync.RWMutex
	observed            map[chainhash.Hash]*observedTransaction
	bin                 [estimateFeeDepth][]*observedTransaction
	// The cached estimates.
	cached []SatoshiPerByte
	// Transactions that have been removed from the bins. This allows us to revert in case of an orphaned block.
	dropped []*registeredBlock
}

// FeeEstimatorState represents a saved FeeEstimator that can be restored with data from an earlier session of the
// program.
type FeeEstimatorState []byte

// SatoshiPerByte is number with units of satoshis per byte.
type SatoshiPerByte float64

// estimateFeeSet is a set of txs that can that is sorted by the fee per kb rate.
type estimateFeeSet struct {
	feeRate []SatoshiPerByte
	bin     [estimateFeeDepth]uint32
}

// observedTransaction represents an observed transaction and some additional data required for the fee estimation
// algorithm.
type observedTransaction struct {
	// A transaction hash.
	hash chainhash.Hash
	// The fee per byte of the transaction in satoshis.
	feeRate SatoshiPerByte
	// The block height when it was observed.
	observed int32
	// The height of the block in which it was mined. If the transaction has not yet been mined, it is zero.
	mined int32
}

// observedTxSet is a set of txs that can that is sorted by hash. It exists for serialization purposes so that a
// serialized state always comes out the same.
type observedTxSet []*observedTransaction

// registeredBlock has the hash of a block and the list of transactions it mined which had been previously observed by
// the FeeEstimator. It is used if Rollback is called to reverse the effect of registering a block.
type registeredBlock struct {
	hash         chainhash.Hash
	transactions []*observedTransaction
}

// TODO incorporate Alex Morcos' modifications to Gavin's initial model
//  https://lists.linuxfoundation.org/pipermail/bitcoin-dev/2014-October/006824.html
const (
	// estimateFeeDepth is the maximum number of blocks before a transaction is confirmed that we want to track.
	estimateFeeDepth = 25
	// estimateFeeBinSize is the number of txs stored in each bin.
	estimateFeeBinSize = 100
	// estimateFeeMaxReplacements is the max number of replacements that can be made by the txs found in a given block.
	estimateFeeMaxReplacements = 10
	// DefaultEstimateFeeMaxRollback is the default number of rollbacks allowed by the fee estimator for orphaned
	// blocks.
	DefaultEstimateFeeMaxRollback = 2
	// DefaultEstimateFeeMinRegisteredBlocks is the default minimum number of blocks which must be observed by the fee
	// estimator before will provide fee estimations.
	DefaultEstimateFeeMinRegisteredBlocks = 3
	bytePerKb                             = 1000
	duoPerSatoshi                         = 1e-8
)

// In case the format for the serialized version of the FeeEstimator changes, we use a version number. If the version
// number changes, it does not make sense to try to upgrade a previous version to a new version. Instead, just start fee
// estimation over.
const estimateFeeSaveVersion = 1

var (
	// EstimateFeeDatabaseKey is the key that we use to store the fee estimator in the database.
	EstimateFeeDatabaseKey = []byte("estimatefee")
)

// EstimateFee estimates the fee per byte to have a tx confirmed a given number of blocks from now.
func (ef *FeeEstimator) EstimateFee(numBlocks uint32) (DUOPerKilobyte, error) {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	// If the number of registered blocks is below the minimum, return an error.
	if ef.numBlocksRegistered < ef.minRegisteredBlocks {
		return -1, errors.New("not enough blocks have been observed")
	}
	if numBlocks == 0 {
		return -1, errors.New("cannot confirm transaction in zero blocks")
	}
	if numBlocks > estimateFeeDepth {
		return -1, fmt.Errorf(
			"can only estimate fees for up to %d blocks from now",
			estimateFeeBinSize,
		)
	}
	// If there are no cached results, generate them.
	if ef.cached == nil {
		ef.cached = ef.estimates()
	}
	return ef.cached[int(numBlocks)-1].ToBtcPerKb(), nil
}

func // LastKnownHeight returns the height of the last block which was
// registered.
(ef *FeeEstimator) LastKnownHeight() int32 {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	return ef.lastKnownHeight
}

// ObserveTransaction is called when a new transaction is observed in the mempool.
func (ef *FeeEstimator) ObserveTransaction(
	t *TxDesc,
) {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	// If we haven't seen a block yet we don't know when this one arrived, so we ignore it.
	if ef.lastKnownHeight == mining.UnminedHeight {
		return
	}
	hash := *t.Tx.Hash()
	if _, ok := ef.observed[hash]; !ok {
		size := uint32(GetTxVirtualSize(t.Tx))
		ef.observed[hash] = &observedTransaction{
			hash:     hash,
			feeRate:  NewSatoshiPerByte(amt.Amount(t.Fee), size),
			observed: t.Height,
			mined:    mining.UnminedHeight,
		}
	}
}

// RegisterBlock informs the fee estimator of a new block to take into account.
func (ef *FeeEstimator) RegisterBlock(
	block *block.Block,
) (e error) {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	// The previous sorted list is invalid, so delete it.
	ef.cached = nil
	height := block.Height()
	if height != ef.lastKnownHeight+1 && ef.lastKnownHeight != mining.UnminedHeight {
		return fmt.Errorf(
			"intermediate block not recorded; current height is %d; new height is %d",
			ef.lastKnownHeight, height,
		)
	}
	// Update the last known height.
	ef.lastKnownHeight = height
	ef.numBlocksRegistered++
	// Randomly order txs in block.
	transactions := make(map[*util.Tx]struct{})
	for _, t := range block.Transactions() {
		transactions[t] = struct{}{}
	}
	// Count the number of replacements we make per bin so that we don't replace too many.
	var replacementCounts [estimateFeeDepth]int
	// Keep track of which txs were dropped in case of an orphan block.
	dropped := &registeredBlock{
		hash:         *block.Hash(),
		transactions: make([]*observedTransaction, 0, 100),
	}
	// Go through the txs in the block.
	for t := range transactions {
		hash := *t.Hash()
		// Have we observed this tx in the mempool?
		o, ok := ef.observed[hash]
		if !ok {
			continue
		}
		// Put the observed tx in the appropriate bin.
		blocksToConfirm := height - o.observed - 1
		// This shouldn't happen if the fee estimator works correctly, but return an error if it does.
		if o.mined != mining.UnminedHeight {
			E.Ln("Estimate fee: transaction ", hash, " has already been mined")
			return errors.New("transaction has already been mined")
		}
		// This shouldn't happen but check just in case to avoid an out-of -bounds array index later.
		if blocksToConfirm >= estimateFeeDepth {
			continue
		}
		// Make sure we do not replace too many transactions per min.
		if replacementCounts[blocksToConfirm] == int(ef.maxReplacements) {
			continue
		}
		o.mined = height
		replacementCounts[blocksToConfirm]++
		bin := ef.bin[blocksToConfirm]
		// Remove a random element and replace it with this new tx.
		if len(bin) == int(ef.binSize) {
			// Don't drop transactions we have just added from this same block.
			l := int(ef.binSize) - replacementCounts[blocksToConfirm]
			drop := rand.Intn(l)
			dropped.transactions = append(dropped.transactions, bin[drop])
			bin[drop] = bin[l-1]
			bin[l-1] = o
		} else {
			bin = append(bin, o)
		}
		ef.bin[blocksToConfirm] = bin
	}
	// Go through the mempool for txs that have been in too long.
	for hash, o := range ef.observed {
		if o.mined == mining.UnminedHeight && height-o.observed >= estimateFeeDepth {
			delete(ef.observed, hash)
		}
	}
	// Add dropped list to history.
	if ef.maxRollback == 0 {
		return nil
	}
	if uint32(len(ef.dropped)) == ef.maxRollback {
		ef.dropped = append(ef.dropped[1:], dropped)
	} else {
		ef.dropped = append(ef.dropped, dropped)
	}
	return nil
}

// Rollback unregisters a recently registered block from the FeeEstimator. This can be used to reverse the effect of an
// orphaned block on the fee estimator. The maximum number of rollbacks allowed is given by maxRollbacks. Note: not
// everything can be rolled back because some transactions are deleted if they have been observed too long ago. That
// means the result of Rollback won't always be exactly the same as if the last block had not happened, but it should be
// close enough.
func (ef *FeeEstimator) Rollback(hash *chainhash.Hash) (e error) {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	// Find this block in the stack of recent registered blocks.
	var n int
	for n = 1; n <= len(ef.dropped); n++ {
		if ef.dropped[len(ef.dropped)-n].hash.IsEqual(hash) {
			break
		}
	}
	if n > len(ef.dropped) {
		return errors.New("no such block was recently registered")
	}
	for i := 0; i < n; i++ {
		ef.rollback()
	}
	return nil
}

// Save records the current state of the FeeEstimator to a []byte that can be restored later.
func (ef *FeeEstimator) Save() FeeEstimatorState {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	// TODO figure out what the capacity should be.
	w := bytes.NewBuffer(make([]byte, 0))
	e := binary.Write(
		w, binary.BigEndian, uint32(estimateFeeSaveVersion),
	)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	// Insert basic parameters.
	e = binary.Write(w, binary.BigEndian, &ef.maxRollback)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	e = binary.Write(w, binary.BigEndian, &ef.binSize)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	e = binary.Write(w, binary.BigEndian, &ef.maxReplacements)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	e = binary.Write(w, binary.BigEndian, &ef.minRegisteredBlocks)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	e = binary.Write(w, binary.BigEndian, &ef.lastKnownHeight)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	e = binary.Write(w, binary.BigEndian, &ef.numBlocksRegistered)
	if e != nil {
		F.Ln("failed to write fee estimates", e)
	}
	// Put all the observed transactions in a sorted list.
	var txCount uint32
	ots := make([]*observedTransaction, len(ef.observed))
	for hash := range ef.observed {
		ots[txCount] = ef.observed[hash]
		txCount++
	}
	sort.Sort(observedTxSet(ots))
	txCount = 0
	observed := make(map[*observedTransaction]uint32)
	e = binary.Write(w, binary.BigEndian, uint32(len(ef.observed)))
	if e != nil {
		F.Ln("failed to write:", e)
	}
	for _, ot := range ots {
		ot.Serialize(w)
		observed[ot] = txCount
		txCount++
	}
	// Save all the right bins.
	for _, list := range ef.bin {
		e = binary.Write(w, binary.BigEndian, uint32(len(list)))
		if e != nil {
			F.Ln("failed to write:", e)
		}
		for _, o := range list {
			e = binary.Write(w, binary.BigEndian, observed[o])
			if e != nil {
				F.Ln("failed to write:", e)
			}
		}
	}
	// Dropped transactions.
	e = binary.Write(w, binary.BigEndian, uint32(len(ef.dropped)))
	if e != nil {
		F.Ln("failed to write:", e)
	}
	for _, registered := range ef.dropped {
		registered.serialize(w, observed)
	}
	// Commit the tx and return.
	return w.Bytes()
}

// estimates returns the set of all fee estimates from 1 to estimateFeeDepth confirmations from now.
func (ef *FeeEstimator) estimates() []SatoshiPerByte {
	set := ef.newEstimateFeeSet()
	estimates := make([]SatoshiPerByte, estimateFeeDepth)
	for i := 0; i < estimateFeeDepth; i++ {
		estimates[i] = set.estimateFee(i + 1)
	}
	return estimates
}

// newEstimateFeeSet creates a temporary data structure that can be used to find all fee estimates.
func (ef *FeeEstimator) newEstimateFeeSet() *estimateFeeSet {
	set := &estimateFeeSet{}
	capacity := 0
	for i, b := range ef.bin {
		l := len(b)
		set.bin[i] = uint32(l)
		capacity += l
	}
	set.feeRate = make([]SatoshiPerByte, capacity)
	i := 0
	for _, b := range ef.bin {
		for _, o := range b {
			set.feeRate[i] = o.feeRate
			i++
		}
	}
	sort.Sort(set)
	return set
}

// rollback rolls back the effect of the last block in the stack of registered blocks.
func (ef *FeeEstimator) rollback() {
	// The previous sorted list is invalid, so delete it.
	ef.cached = nil
	// pop the last list of dropped txs from the stack.
	last := len(ef.dropped) - 1
	if last == -1 {
		// Cannot really happen because the exported calling function only rolls back a block already known to be in the
		// list of dropped transactions.
		return
	}
	dropped := ef.dropped[last]
	// where we are in each bin as we replace txs?
	var replacementCounters [estimateFeeDepth]int
	// Go through the txs in the dropped block.
	for _, o := range dropped.transactions {
		// Which bin was this tx in?
		blocksToConfirm := o.mined - o.observed - 1
		bin := ef.bin[blocksToConfirm]
		var counter = replacementCounters[blocksToConfirm]
		// Continue to go through that bin where we left off.
		for {
			if counter >= len(bin) {
				// Panic, as we have entered an unrecoverable invalid state.
				panic(
					errors.New(
						"illegal state: cannot rollback dropped transaction",
					),
				)
			}
			prev := bin[counter]
			if prev.mined == ef.lastKnownHeight {
				prev.mined = mining.UnminedHeight
				bin[counter] = o
				counter++
				break
			}
			counter++
		}
		replacementCounters[blocksToConfirm] = counter
	}
	// Continue going through bins to find other txs to remove which did not replace any other when they were entered.
	for i, j := range replacementCounters {
		for {
			l := len(ef.bin[i])
			if j >= l {
				break
			}
			prev := ef.bin[i][j]
			if prev.mined == ef.lastKnownHeight {
				prev.mined = mining.UnminedHeight
				newBin := append(ef.bin[i][0:j], ef.bin[i][j+1:l]...)
				// TODO This line should prevent an unintentional memory leak
				//  but it causes a panic when it is uncommented.
				// ef.bin[i][j] = nil
				ef.bin[i] = newBin
				continue
			}
			j++
		}
	}
	ef.dropped = ef.dropped[0:last]
	// The number of blocks the fee estimator has seen is decremented.
	ef.numBlocksRegistered--
	ef.lastKnownHeight--
}
func (b *estimateFeeSet) Len() int           { return len(b.feeRate) }
func (b *estimateFeeSet) Less(i, j int) bool { return b.feeRate[i] > b.feeRate[j] }
func (b *estimateFeeSet) Swap(i, j int) {
	b.feeRate[i], b.feeRate[j] = b.feeRate[j], b.feeRate[i]
}

// estimateFee returns the estimated fee for a transaction to confirm in confirmations blocks from now, given the data
// set we have collected.
func (b *estimateFeeSet) estimateFee(confirmations int) SatoshiPerByte {
	if confirmations <= 0 {
		return SatoshiPerByte(math.Inf(1))
	}
	if confirmations > estimateFeeDepth {
		return 0
	}
	// We don't have any transactions!
	if len(b.feeRate) == 0 {
		return 0
	}
	var min, max int
	for i := 0; i < confirmations-1; i++ {
		min += int(b.bin[i])
	}
	max = min + int(b.bin[confirmations-1]) - 1
	if max < min {
		max = min
	}
	feeIndex := (min + max) / 2
	if feeIndex >= len(b.feeRate) {
		feeIndex = len(b.feeRate) - 1
	}
	return b.feeRate[feeIndex]
}
func (o *observedTransaction) Serialize(w io.Writer) {
	e := binary.Write(w, binary.BigEndian, o.hash)
	if e != nil {
		F.Ln("failed to serialize observed transaction:", e)
	}
	e = binary.Write(w, binary.BigEndian, o.feeRate)
	if e != nil {
		F.Ln("failed to serialize observed transaction:", e)
	}
	e = binary.Write(w, binary.BigEndian, o.observed)
	if e != nil {
		F.Ln("failed to serialize observed transaction:", e)
	}
	e = binary.Write(w, binary.BigEndian, o.mined)
	if e != nil {
		F.Ln("failed to serialize observed transaction:", e)
	}
}

func (rb *registeredBlock) serialize(
	w io.Writer,
	txs map[*observedTransaction]uint32,
) {
	e := binary.Write(w, binary.BigEndian, rb.hash)
	if e != nil {
		F.Ln("failed to write:", e)
	}
	e = binary.Write(w, binary.BigEndian, uint32(len(rb.transactions)))
	if e != nil {
		F.Ln("failed to write:", e)
	}
	for _, o := range rb.transactions {
		e = binary.Write(w, binary.BigEndian, txs[o])
		if e != nil {
			F.Ln("failed to write:", e)
		}
	}
}

// Fee returns the fee for a transaction of a given size for the given fee rate.
func (rate SatoshiPerByte) Fee(size uint32) amt.Amount {
	// If our rate is the error value, return that.
	if rate == SatoshiPerByte(-1) {
		return amt.Amount(-1)
	}
	return amt.Amount(float64(rate) * float64(size))
}

// ToBtcPerKb returns a float value that represents the given SatoshiPerByte converted to satoshis per kb.
func (rate SatoshiPerByte) ToBtcPerKb() DUOPerKilobyte {
	// If our rate is the error value, return that.
	if rate == SatoshiPerByte(-1.0) {
		return -1.0
	}
	return DUOPerKilobyte(float64(rate) * bytePerKb * duoPerSatoshi)
}

func (q observedTxSet) Len() int { return len(q) }
func (q observedTxSet) Less(i, j int) bool {
	return strings.Compare(q[i].hash.String(), q[j].hash.String()) < 0
}
func (q observedTxSet) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

// NewFeeEstimator creates a FeeEstimator for which at most maxRollback blocks can be unregistered and which returns an
// error unless minRegisteredBlocks have been registered with it.
func NewFeeEstimator(maxRollback, minRegisteredBlocks uint32) *FeeEstimator {
	return &FeeEstimator{
		maxRollback:         maxRollback,
		minRegisteredBlocks: minRegisteredBlocks,
		lastKnownHeight:     mining.UnminedHeight,
		binSize:             estimateFeeBinSize,
		maxReplacements:     estimateFeeMaxReplacements,
		observed:            make(map[chainhash.Hash]*observedTransaction),
		dropped:             make([]*registeredBlock, 0, maxRollback),
	}
}

// NewSatoshiPerByte creates a SatoshiPerByte from an Amount and a size in bytes.
func NewSatoshiPerByte(fee amt.Amount, size uint32) SatoshiPerByte {
	return SatoshiPerByte(float64(fee) / float64(size))
}

// RestoreFeeEstimator takes a FeeEstimatorState that was previously returned by Save and restores it to a FeeEstimator
func RestoreFeeEstimator(data FeeEstimatorState) (*FeeEstimator, error) {
	r := bytes.NewReader(data)
	// Chk version
	var version uint32
	e := binary.Read(r, binary.BigEndian, &version)
	if e != nil {
		return nil, e
	}
	if version != estimateFeeSaveVersion {
		return nil, fmt.Errorf(
			"incorrect version: expected %d found %d",
			estimateFeeSaveVersion, version,
		)
	}
	ef := &FeeEstimator{
		observed: make(map[chainhash.Hash]*observedTransaction),
	}
	// Read basic parameters.
	e = binary.Read(r, binary.BigEndian, &ef.maxRollback)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ef.binSize)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ef.maxReplacements)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ef.minRegisteredBlocks)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ef.lastKnownHeight)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ef.numBlocksRegistered)
	if e != nil {
		F.Ln("failed to read", e)
	}
	// Read transactions.
	var numObserved uint32
	observed := make(map[uint32]*observedTransaction)
	e = binary.Read(r, binary.BigEndian, &numObserved)
	if e != nil {
		F.Ln("failed to read", e)
	}
	for i := uint32(0); i < numObserved; i++ {
		var ot *observedTransaction
		ot, e = deserializeObservedTransaction(r)
		if e != nil {
			return nil, e
		}
		observed[i] = ot
		ef.observed[ot.hash] = ot
	}
	// Read bins.
	for i := 0; i < estimateFeeDepth; i++ {
		var numTransactions uint32
		e = binary.Read(r, binary.BigEndian, &numTransactions)
		if e != nil {
			F.Ln("failed to read", e)
		}
		bin := make([]*observedTransaction, numTransactions)
		for j := uint32(0); j < numTransactions; j++ {
			var index uint32
			e = binary.Read(r, binary.BigEndian, &index)
			if e != nil {
				F.Ln("failed to read", e)
			}
			var exists bool
			bin[j], exists = observed[index]
			if !exists {
				return nil, fmt.Errorf(
					"invalid transaction reference %d",
					index,
				)
			}
		}
		ef.bin[i] = bin
	}
	// Read dropped transactions.
	var numDropped uint32
	e = binary.Read(r, binary.BigEndian, &numDropped)
	if e != nil {
		F.Ln("failed to read", e)
	}
	ef.dropped = make([]*registeredBlock, numDropped)
	for i := uint32(0); i < numDropped; i++ {
		var e error
		ef.dropped[int(i)], e = deserializeRegisteredBlock(r, observed)
		if e != nil {
			return nil, e
		}
	}
	return ef, nil
}
func deserializeObservedTransaction(r io.Reader) (*observedTransaction, error) {
	ot := observedTransaction{}
	// The first 32 bytes should be a hash.
	e := binary.Read(r, binary.BigEndian, &ot.hash)
	if e != nil {
		F.Ln("failed to read", e)
	}
	// The next 8 are SatoshiPerByte
	e = binary.Read(r, binary.BigEndian, &ot.feeRate)
	if e != nil {
		F.Ln("failed to read", e)
	}
	// And next there are two uint32's.
	e = binary.Read(r, binary.BigEndian, &ot.observed)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &ot.mined)
	if e != nil {
		F.Ln("failed to read", e)
	}
	return &ot, nil
}
func deserializeRegisteredBlock(
	r io.Reader,
	txs map[uint32]*observedTransaction,
) (*registeredBlock, error) {
	var lenTransactions uint32
	rb := &registeredBlock{}
	e := binary.Read(r, binary.BigEndian, &rb.hash)
	if e != nil {
		F.Ln("failed to read", e)
	}
	e = binary.Read(r, binary.BigEndian, &lenTransactions)
	if e != nil {
		F.Ln("failed to read", e)
	}
	rb.transactions = make([]*observedTransaction, lenTransactions)
	for i := uint32(0); i < lenTransactions; i++ {
		var index uint32
		e = binary.Read(r, binary.BigEndian, &index)
		if e != nil {
			F.Ln("failed to read", e)
		}
		rb.transactions[i] = txs[index]
	}
	return rb, nil
}
