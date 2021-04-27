package blockchain

import (
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/bits"
	"github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/log"

	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
)

// BehaviorFlags is a bitmask defining tweaks to the normal behavior when
// performing chain processing and consensus rules checks.
type BehaviorFlags uint32

const (
	// BFFastAdd may be set to indicate that several checks can be avoided for the
	// block since it is already known to fit into the chain due to already proving
	// it correct links into the chain up to a known checkpoint. This is primarily
	// used for headers-first mode.
	BFFastAdd BehaviorFlags = 1 << iota
	// BFNoPoWCheck may be set to indicate the proof of work check which ensures a
	// block hashes to a value less than the required target will not be performed.
	BFNoPoWCheck
	// BFNone is a convenience value to specifically indicate no flags.
	BFNone BehaviorFlags = 0
)

// ProcessBlock is the main workhorse for handling insertion of new blocks into
// the block chain. It includes functionality such as rejecting duplicate
// blocks, ensuring blocks follow all rules, orphan handling, and insertion into
// the block chain along with best chain selection and reorganization.
//
// When no errors occurred during processing, the first return value indicates
// whether or not the block is on the main chain and the second indicates
// whether or not the block is an orphan.
//
// This function is safe for concurrent access.
func (b *BlockChain) ProcessBlock(
	workerNumber uint32, candidateBlock *block.Block,
	flags BehaviorFlags, blockHeight int32,
) (bool, bool, error,) {
	T.Ln("blockchain.ProcessBlock", blockHeight, log.Caller("\nfrom", 1))
	var prevBlock *block.Block
	var e error
	prevBlock, e = b.BlockByHash(&candidateBlock.WireBlock().Header.PrevBlock)
	if prevBlock != nil {
		blockHeight = prevBlock.Height() + 1
	} else {
		return false, false, e
	}
	// trc.S(prevBlock)
	b.ChainLock.Lock()
	defer b.ChainLock.Unlock()
	fastAdd := flags&BFFastAdd == BFFastAdd
	blockHash := candidateBlock.Hash()
	hf := fork.GetCurrent(blockHeight)
	bhwa := candidateBlock.WireBlock().BlockHashWithAlgos
	var algo int32
	switch hf {
	case 0:
		if candidateBlock.WireBlock().Header.Version != 514 {
			algo = 2
		} else {
			algo = 514
		}
	case 1:
		algo = candidateBlock.WireBlock().Header.Version
	}
	// The candidateBlock must not already exist in the main chain or side chains.
	var exists bool
	if exists, e = b.blockExists(blockHash); E.Chk(e) {
		return false, false, e
	}
	if exists {
		str := ruleError(ErrDuplicateBlock, fmt.Sprintf("already have candidateBlock %v", bhwa(blockHeight).String()))
		E.Ln(str)
		return false, false, str
	}
	// The candidateBlock must not already exist as an orphan.
	if _, exists := b.orphans[*blockHash]; exists {
		str := ruleError(ErrDuplicateBlock, fmt.Sprintf("already have candidateBlock (orphan)"))
		E.Ln(str)
		return false, false, str
	}
	// Perform preliminary sanity checks on the candidateBlock and its transactions.
	var DoNotCheckPow bool
	pl := fork.GetMinDiff(fork.GetAlgoName(algo, blockHeight), blockHeight)
	T.F("powLimit %d %s %d %064x", algo, fork.GetAlgoName(algo, blockHeight), blockHeight, pl)
	ph := &candidateBlock.WireBlock().Header.PrevBlock
	pn := b.Index.LookupNode(ph)
	if pn == nil {
		return false, false, errors.New("could not find parent block of candidate block")
	}
	var pb *BlockNode
	pb = pn.GetLastWithAlgo(algo)
	if pb == nil {
		DoNotCheckPow = true
	}
	T.F("checkBlockSanity powLimit %d %s %d %064x ts %v", algo, fork.GetAlgoName(algo, blockHeight), blockHeight, pl,
		pn.Header().Timestamp,
	)
	if e = checkBlockSanity(
		candidateBlock,
		pl,
		b.timeSource,
		flags,
		DoNotCheckPow,
		blockHeight,
		pn.Header().Timestamp,
	); E.Chk(e) {
		return false, false, e
	}
	T.Ln("searching back to checkpoints")
	// Find the previous checkpoint and perform some additional checks based on the
	// checkpoint. This provides a few nice properties such as preventing old side
	// chain blocks before the last checkpoint, rejecting easy to mine, but
	// otherwise bogus, blocks that could be used to eat memory, and ensuring
	// expected (versus claimed) proof of work requirements since the previous
	// checkpoint are met.
	blockHeader := &candidateBlock.WireBlock().Header
	var checkpointNode *BlockNode
	if checkpointNode, e = b.findPreviousCheckpoint(); E.Chk(e) {
		return false, false, e
	}
	if checkpointNode != nil {
		// Ensure the candidateBlock timestamp is after the checkpoint timestamp.
		checkpointTime := time.Unix(checkpointNode.timestamp, 0)
		if blockHeader.Timestamp.Before(checkpointTime) {
			str := fmt.Sprintf(
				"candidateBlock %v has timestamp %v before last checkpoint timestamp %v",
				bhwa(blockHeight).String(), blockHeader.Timestamp, checkpointTime,
			)
			T.Ln(str)
			return false, false, ruleError(ErrCheckpointTimeTooOld, str)
		}
		if !fastAdd {
			// Even though the checks prior to now have already ensured the proof of work
			// exceeds the claimed amount, the claimed amount is a field in the candidateBlock header
			// which could be forged. This check ensures the proof of work is at least the
			// minimum expected based on elapsed time since the last checkpoint and maximum
			// adjustment allowed by the retarget rules.
			duration := blockHeader.Timestamp.Sub(checkpointTime)
			requiredTarget := bits.CompactToBig(
				b.calcEasiestDifficulty(
					checkpointNode.bits, duration,
				),
			)
			currentTarget := bits.CompactToBig(blockHeader.Bits)
			if currentTarget.Cmp(requiredTarget) > 0 {
				str := fmt.Sprintf(
					"processing: candidateBlock target difficulty of %064x is too low when compared to the"+
						" previous checkpoint", currentTarget,
				)
				E.Ln(str)
				return false, false, ruleError(ErrDifficultyTooLow, str)
			}
		}
	}
	T.Ln("handling orphans")
	// Handle orphan blocks.
	prevHash := &blockHeader.PrevBlock
	var prevHashExists bool
	if prevHashExists, e = b.blockExists(prevHash); E.Chk(e) {
		return false, false, e
	}
	if !prevHashExists {
		D.C(
			func() string {
				return fmt.Sprintf(
					"adding orphan candidateBlock %v with parent %v",
					bhwa(blockHeight).String(),
					prevHash,
				)
			},
		)
		b.addOrphanBlock(candidateBlock)
		return false, true, nil
	}
	// The candidateBlock has passed all context independent checks and appears sane enough
	// to potentially accept it into the candidateBlock chain.
	T.Ln("maybe accept candidateBlock")
	var isMainChain bool
	if isMainChain, e = b.maybeAcceptBlock(workerNumber, candidateBlock, flags); E.Chk(e) {
		return false, false, e
	}
	// Accept any orphan blocks that depend on this candidateBlock (they are no longer
	// orphans) and repeat for those accepted blocks until there are no more.
	if isMainChain {
		T.Ln("new candidateBlock on main chain")
		// Traces(candidateBlock)
	}
	if e = b.processOrphans(workerNumber, blockHash, flags); E.Chk(e) {
		return false, false, e
	}
	T.F(
		"accepted candidateBlock %d %v %s",
		blockHeight, bhwa(blockHeight).String(), fork.GetAlgoName(
			candidateBlock.WireBlock().
				Header.Version, blockHeight,
		),
	)
	T.Ln("finished blockchain.ProcessBlock")
	return isMainChain, false, nil
}

// blockExists determines whether a block with the given hash exists either in
// the main chain or any side chains.
//
// This function is safe for concurrent access.
func (b *BlockChain) blockExists(hash *chainhash.Hash) (exists bool, e error) {
	// Chk block index first (could be main chain or side chain blocks).
	if b.Index.HaveBlock(hash) {
		return true, nil
	}
	// Check in the database.
	e = b.db.View(
		func(dbTx database.Tx) (e error) {
			exists, e = dbTx.HasBlock(hash)
			if e != nil || !exists {
				return e
			}
			// Ignore side chain blocks in the database. This is necessary because there is
			// not currently any record of the associated block index data such as its block
			// height, so it's not yet possible to efficiently load the block and do
			// anything useful with it. Ultimately the entire block index should be
			// serialized instead of only the current main chain so it can be consulted
			// directly.
			if _, e = dbFetchHeightByHash(dbTx, hash); E.Chk(e) {
			}
			if isNotInMainChainErr(e) {
				exists = false
				return nil
			}
			return e
		},
	)
	return exists, e
}

// processOrphans determines if there are any orphans which depend on the passed
// block hash (they are no longer orphans if true) and potentially accepts them.
// It repeats the process for the newly accepted blocks ( to detect further
// orphans which may no longer be orphans) until there are no more. The flags do
// not modify the behavior of this function directly, however they are needed to
// pass along to maybeAcceptBlock.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) processOrphans(
	workerNumber uint32, hash *chainhash.Hash,
	flags BehaviorFlags,
) (e error) {
	// Start with processing at least the passed hash. Leave a little room for
	// additional orphan blocks that need to be processed without needing to grow
	// the array in the common case.
	processHashes := make([]*chainhash.Hash, 0, 10)
	processHashes = append(processHashes, hash)
	for len(processHashes) > 0 {
		// Pop the first hash to process from the slice.
		processHash := processHashes[0]
		processHashes[0] = nil // Prevent GC leak.
		processHashes = processHashes[1:]
		// Look up all orphans that are parented by the block we just accepted. This
		// will typically only be one, but it could be multiple if multiple blocks are
		// mined and broadcast around the same time. The one with the most proof of work
		// will eventually win out. An indexing for loop is intentionally used over a
		// range here as range does not reevaluate the slice on each iteration nor does
		// it adjust the index for the modified slice.
		for i := 0; i < len(b.prevOrphans[*processHash]); i++ {
			orphan := b.prevOrphans[*processHash][i]
			if orphan == nil {
				D.F(
					"found a nil entry at index %d in the orphan dependency list for block %v",
					i, processHash,
				)
				continue
			}
			// Remove the orphan from the orphan pool.
			orphanHash := orphan.block.Hash()
			b.removeOrphanBlock(orphan)
			i--
			// Potentially accept the block into the block chain.
			var e error
			if _, e = b.maybeAcceptBlock(workerNumber, orphan.block, flags); E.Chk(e) {
				return e
			}
			// Add this block to the list of blocks to process so any orphan blocks that
			// depend on this block are handled too.
			processHashes = append(processHashes, orphanHash)
		}
	}
	return nil
}
