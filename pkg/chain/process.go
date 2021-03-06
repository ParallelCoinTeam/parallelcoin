package blockchain

import (
	"fmt"
	"time"
	
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/util"
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
func (b *BlockChain) ProcessBlock(workerNumber uint32, candidateBlock *util.Block, flags BehaviorFlags, height int32,) (
	bool, bool, error,
) {
	Trace("blockchain.ProcessBlock", height)
	blockHeight := height
	var prevBlock *util.Block
	var err error
	prevBlock, err = b.BlockByHash(&candidateBlock.MsgBlock().Header.PrevBlock)
	if prevBlock != nil {
		blockHeight = prevBlock.Height() + 1
	} else {
		return false, false, err
	}
	Debug(">> ", fork.GetCurrent(blockHeight) > 0)
	if fork.GetCurrent(blockHeight) > 0 {
		Debug("checking for plan 9 hard fork invariant of timestamp always progressing")
		// prevTS := candidateBlock.MsgBlock().Header.Timestamp
		Debug(candidateBlock.MsgBlock().Header.Timestamp.Round(time.Second).Sub(prevBlock.MsgBlock().Header.Timestamp.Round(time.Second)))
		if candidateBlock.MsgBlock().Header.Timestamp.Round(time.Second).Sub(prevBlock.MsgBlock().Header.Timestamp.Round(time.Second)) <= time.Second {
			return false, false, ruleError(
				ErrTimeTooOld,
				fmt.Sprint("new blocks cannot be less than one second ahead of the chain tip"),
			)
		}
	}
	b.chainLock.Lock()
	defer b.chainLock.Unlock()
	fastAdd := flags&BFFastAdd == BFFastAdd
	blockHash := candidateBlock.Hash()
	hf := fork.GetCurrent(blockHeight)
	bhwa := candidateBlock.MsgBlock().BlockHashWithAlgos
	var algo int32
	switch hf {
	case 0:
		if candidateBlock.MsgBlock().Header.Version != 514 {
			algo = 2
		} else {
			algo = 514
		}
	case 1:
		algo = candidateBlock.MsgBlock().Header.Version
	}
	// The candidateBlock must not already exist in the main chain or side chains.
	var exists bool
	if exists, err = b.blockExists(blockHash); Check(err) {
		return false, false, err
	}
	if exists {
		str := ruleError(ErrDuplicateBlock, fmt.Sprintf("already have candidateBlock %v", bhwa(blockHeight).String()))
		Error(str)
		return false, false, str
	}
	// The candidateBlock must not already exist as an orphan.
	if _, exists := b.orphans[*blockHash]; exists {
		str := ruleError(ErrDuplicateBlock, fmt.Sprintf("already have candidateBlock (orphan)"))
		Error(str)
		return false, false, str
	}
	// Perform preliminary sanity checks on the candidateBlock and its transactions.
	var DoNotCheckPow bool
	pl := fork.GetMinDiff(fork.GetAlgoName(algo, blockHeight), blockHeight)
	Tracef("powLimit %d %s %d %064x", algo, fork.GetAlgoName(algo, blockHeight), blockHeight, pl)
	ph := &candidateBlock.MsgBlock().Header.PrevBlock
	pn := b.Index.LookupNode(ph)
	var pb *BlockNode
	if pn == nil {
		DoNotCheckPow = true
	} else {
		pb = pn.GetLastWithAlgo(algo)
		if pb == nil {
			DoNotCheckPow = true
		}
	}
	Trace("checkBlockSanity powLimit %d %s %d %064x", algo, fork.GetAlgoName(algo, blockHeight), blockHeight, pl)
	if err = checkBlockSanity(candidateBlock, pl, b.timeSource, flags, DoNotCheckPow, blockHeight); Check(err) {
		return false, false, err
	}
	Trace("searching back to checkpoints")
	// Find the previous checkpoint and perform some additional checks based on the
	// checkpoint. This provides a few nice properties such as preventing old side
	// chain blocks before the last checkpoint, rejecting easy to mine, but
	// otherwise bogus, blocks that could be used to eat memory, and ensuring
	// expected (versus claimed) proof of work requirements since the previous
	// checkpoint are met.
	blockHeader := &candidateBlock.MsgBlock().Header
	var checkpointNode *BlockNode
	if checkpointNode, err = b.findPreviousCheckpoint(); Check(err) {
		return false, false, err
	}
	if checkpointNode != nil {
		// Ensure the candidateBlock timestamp is after the checkpoint timestamp.
		checkpointTime := time.Unix(checkpointNode.timestamp, 0)
		if blockHeader.Timestamp.Before(checkpointTime) {
			str := fmt.Sprintf(
				"candidateBlock %v has timestamp %v before last checkpoint timestamp %v",
				bhwa(blockHeight).String(), blockHeader.Timestamp, checkpointTime,
			)
			Trace(str)
			return false, false, ruleError(ErrCheckpointTimeTooOld, str)
		}
		if !fastAdd {
			// Even though the checks prior to now have already ensured the proof of work
			// exceeds the claimed amount, the claimed amount is a field in the candidateBlock header
			// which could be forged. This check ensures the proof of work is at least the
			// minimum expected based on elapsed time since the last checkpoint and maximum
			// adjustment allowed by the retarget rules.
			duration := blockHeader.Timestamp.Sub(checkpointTime)
			requiredTarget := fork.CompactToBig(
				b.calcEasiestDifficulty(
					checkpointNode.bits, duration,
				),
			)
			currentTarget := fork.CompactToBig(blockHeader.Bits)
			if currentTarget.Cmp(requiredTarget) > 0 {
				str := fmt.Sprintf(
					"processing: candidateBlock target difficulty of %064x is too low when compared to the"+
						" previous checkpoint", currentTarget,
				)
				Error(str)
				return false, false, ruleError(ErrDifficultyTooLow, str)
			}
		}
	}
	Trace("handling orphans")
	// Handle orphan blocks.
	prevHash := &blockHeader.PrevBlock
	var prevHashExists bool
	if prevHashExists, err = b.blockExists(prevHash); Check(err) {
		return false, false, err
	}
	if !prevHashExists {
		Debugc(
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
	Trace("maybe accept candidateBlock")
	var isMainChain bool
	if isMainChain, err = b.maybeAcceptBlock(workerNumber, candidateBlock, flags); Check(err) {
		return false, false, err
	}
	// Accept any orphan blocks that depend on this candidateBlock (they are no longer
	// orphans) and repeat for those accepted blocks until there are no more.
	if isMainChain {
		Trace("new candidateBlock on main chain")
		// Traces(candidateBlock)
	}
	if err = b.processOrphans(workerNumber, blockHash, flags); Check(err) {
		return false, false, err
	}
	Tracef(
		"accepted candidateBlock %d %v %s",
		blockHeight, bhwa(blockHeight).String(), fork.GetAlgoName(
			candidateBlock.MsgBlock().
				Header.Version, blockHeight,
		),
	)
	Trace("finished blockchain.ProcessBlock")
	return isMainChain, false, nil
}

// blockExists determines whether a block with the given hash exists either in
// the main chain or any side chains.
//
// This function is safe for concurrent access.
func (b *BlockChain) blockExists(hash *chainhash.Hash) (bool, error) {
	// Check block index first (could be main chain or side chain blocks).
	if b.Index.HaveBlock(hash) {
		return true, nil
	}
	// Check in the database.
	var exists bool
	err := b.db.View(
		func(dbTx database.Tx) error {
			var err error
			exists, err = dbTx.HasBlock(hash)
			if err != nil || !exists {
				return err
			}
			// Ignore side chain blocks in the database. This is necessary because there is
			// not currently any record of the associated block index data such as its block
			// height, so it's not yet possible to efficiently load the block and do
			// anything useful with it. Ultimately the entire block index should be
			// serialized instead of only the current main chain so it can be consulted
			// directly.
			if _, err = dbFetchHeightByHash(dbTx, hash); Check(err) {
			}
			if isNotInMainChainErr(err) {
				exists = false
				return nil
			}
			return err
		},
	)
	return exists, err
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
) error {
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
				Debugf(
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
			var err error
			if _, err = b.maybeAcceptBlock(workerNumber, orphan.block, flags); Check(err) {
				return err
			}
			// Add this block to the list of blocks to process so any orphan blocks that
			// depend on this block are handled too.
			processHashes = append(processHashes, orphanHash)
		}
	}
	return nil
}
