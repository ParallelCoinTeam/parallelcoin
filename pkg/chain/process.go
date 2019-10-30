package blockchain

import (
	"fmt"
	"time"

	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
)

type // BehaviorFlags is a bitmask defining tweaks to the normal behavior when
	// performing chain processing and consensus rules checks.
	BehaviorFlags uint32

const (
	// BFFastAdd may be set to indicate that several checks can be avoided
	// for the block since it is already known to fit into the chain due to
	// already proving it correct links into the chain up to a known
	// checkpoint.  This is primarily used for headers-first mode.
	BFFastAdd BehaviorFlags = 1 << iota
	// BFNoPoWCheck may be set to indicate the proof of work check which
	// ensures a block hashes to a value less than the required target will
	// not be performed.
	BFNoPoWCheck
	// BFNone is a convenience value to specifically indicate no flags.
	BFNone BehaviorFlags = 0
)

func // ProcessBlock is the main workhorse for handling insertion of new blocks
// into the block chain.  It includes functionality such as rejecting
// duplicate blocks, ensuring blocks follow all rules, orphan handling,
// and insertion into the block chain along with best chain selection and
// reorganization. When no errors occurred during processing,
// the first return value indicates whether or not the block is on the main
// chain and the second indicates whether or not the block is an orphan.
// This function is safe for concurrent access.
(b *BlockChain) ProcessBlock(workerNumber uint32, block *util.Block,
	flags BehaviorFlags, height int32) (bool, bool, error) {
	// log.WARN("blockchain.ProcessBlock NEW MAYBE BLOCK", height)
	blockHeight := height
	bb, _ := b.BlockByHash(&block.MsgBlock().Header.PrevBlock)
	if bb != nil {
		blockHeight = bb.Height() + 1
	}
	b.chainLock.Lock()
	defer b.chainLock.Unlock()
	fastAdd := flags&BFFastAdd == BFFastAdd
	blockHash := block.Hash()
	hf := fork.GetCurrent(blockHeight)
	blockHashWithAlgo := block.MsgBlock().BlockHashWithAlgos(blockHeight).String()
	// log.WARNC(func() string {
	// 		return "processing block " + blockHashWithAlgo
	// 	})
	var algo int32
	switch hf {
	case 0:
		if block.MsgBlock().Header.Version != 514 {
			algo = 2
		} else {
			algo = 514
		}
	case 1:
		algo = block.MsgBlock().Header.Version
	}
	// The block must not already exist in the main chain or side chains.
	exists, err := b.blockExists(blockHash)
	if err != nil {
		log.ERROR(err)
		log.ERROR(err)
		return false, false, err
	}
	if exists {
		str := fmt.Sprintf("already have block %v", blockHashWithAlgo)
		return false, false, ruleError(ErrDuplicateBlock, str)
	}
	// The block must not already exist as an orphan.
	if _, exists := b.orphans[*blockHash]; exists {
		str := fmt.Sprintf(
			"already have block (orphan) %v", blockHashWithAlgo)
		return false, false, ruleError(ErrDuplicateBlock, str)
	}
	// Perform preliminary sanity checks on the block and its transactions.
	var DoNotCheckPow bool
	pl := fork.GetMinDiff(fork.GetAlgoName(algo, blockHeight), blockHeight)
	// log.WARNF("powLimit %d %s %d %064x", algo, fork.GetAlgoName(algo,
	// 	blockHeight), blockHeight, pl)
	ph := &block.MsgBlock().Header.PrevBlock
	pn := b.Index.LookupNode(ph)
	if pn == nil {
		// log.WARN("found no previous node")
		DoNotCheckPow = true
	}
	pb := pn.GetLastWithAlgo(algo)
	if pb == nil {
		// pl = &chaincfg.AllOnes !!!!!!!!!!!!!!!!!!
		DoNotCheckPow = true
	}
	// log.WARNF("checkBlockSanity powLimit %d %s %d %064x", algo,
	// 	fork.GetAlgoName(algo, blockHeight), blockHeight, pl)
	err = checkBlockSanity(block, pl, b.timeSource, flags, DoNotCheckPow, blockHeight)
	if err != nil {
		log.ERROR("block processing error: ", err)
		return false, false, err
	}
	// log.WARN("searching back to checkpoints")
	// Find the previous checkpoint and perform some additional checks based
	// on the checkpoint.  This provides a few
	// nice properties such as preventing old side chain blocks before the
	// last checkpoint, rejecting easy to mine,
	// but otherwise bogus, blocks that could be used to eat memory,
	// and ensuring expected (versus claimed) proof of
	// work requirements since the previous checkpoint are met.
	blockHeader := &block.MsgBlock().Header
	checkpointNode, err := b.findPreviousCheckpoint()
	if err != nil {
		log.ERROR(err)
		log.ERROR(err)
		return false, false, err
	}
	if checkpointNode != nil {
		// Ensure the block timestamp is after the checkpoint timestamp.
		checkpointTime := time.Unix(checkpointNode.timestamp, 0)
		if blockHeader.Timestamp.Before(checkpointTime) {
			str := fmt.Sprintf("block %v has timestamp %v before "+
				"last checkpoint timestamp %v", blockHashWithAlgo,
				blockHeader.Timestamp, checkpointTime)
			return false, false, ruleError(ErrCheckpointTimeTooOld, str)
		}
		if !fastAdd {
			// Even though the checks prior to now have already ensured the
			// proof of work exceeds the claimed amount,
			// the claimed amount is a field in the block header which could
			// be forged.  This check ensures the proof of work is at least
			// the minimum expected based on elapsed time since the last
			// checkpoint and maximum adjustment allowed by the retarget rules.
			duration := blockHeader.Timestamp.Sub(checkpointTime)
			requiredTarget := fork.CompactToBig(b.calcEasiestDifficulty(
				checkpointNode.bits, duration))
			currentTarget := fork.CompactToBig(blockHeader.Bits)
			if currentTarget.Cmp(requiredTarget) > 0 {
				str := fmt.Sprintf("processing: block target difficulty of"+
					" %064x is too low when compared to the previous checkpoint", currentTarget)
				return false, false, ruleError(ErrDifficultyTooLow, str)
			}
		}
	}
	// log.WARN("handling orphans")
	// Handle orphan blocks.
	prevHash := &blockHeader.PrevBlock
	prevHashExists, err := b.blockExists(prevHash)
	if err != nil {
		log.ERROR(err)
		log.ERROR(err)
		return false, false, err
	}
	if !prevHashExists {
		log.WARNC(func() string {
			return fmt.Sprintf(
				"adding orphan block %v with parent %v",
				blockHashWithAlgo,
				prevHash,
			)
		})
		b.addOrphanBlock(block)
		return false, true, nil
	}
	// The block has passed all context independent checks and appears sane
	// enough to potentially accept it into the block chain.
	// log.WARN("maybe accept block")
	isMainChain, err := b.maybeAcceptBlock(workerNumber, block, flags)
	if err != nil {
		log.ERROR(err)
		log.ERROR(err)
		return false, false, err
	}
	// Accept any orphan blocks that depend on this block (they are no longer
	// orphans) and repeat for those accepted blocks until there are no more.
	if isMainChain {
		log.WARN("new block on main chain")
	}
	err = b.processOrphans(workerNumber, blockHash, flags)
	if err != nil {
		log.ERROR(err)
		log.ERROR(err)
		return false, false, err
	}
	// log.DEBUGF("accepted block %d %v",
	// 	blockHeight, blockHashWithAlgo, fork.GetAlgoName(block.MsgBlock().
	// 		Header.Version, blockHeight))

	// log.WARN("finished blockchain.ProcessBlock")
	return isMainChain, false, nil
}

func // blockExists determines whether a block with the given hash exists
// either in the main chain or any side chains.
// This function is safe for concurrent access.
(b *BlockChain) blockExists(hash *chainhash.Hash) (bool, error) {
	// Check block index first (could be main chain or side chain blocks).
	if b.Index.HaveBlock(hash) {
		return true, nil
	}
	// Check in the database.
	var exists bool
	err := b.db.View(func(dbTx database.Tx) error {
		var err error
		exists, err = dbTx.HasBlock(hash)
		if err != nil || !exists {
			return err
		}
		// Ignore side chain blocks in the database.
		// This is necessary because there is not currently any record of the
		// associated block index data such as its block height,
		// so it's not yet possible to efficiently load the block and do
		// anything useful with it.
		// Ultimately the entire block index should be serialized instead of
		// only the current main chain so it can be consulted directly.
		_, err = dbFetchHeightByHash(dbTx, hash)
		if isNotInMainChainErr(err) {
			exists = false
			return nil
		}
		return err
	})
	return exists, err
}

func // processOrphans determines if there are any orphans which depend on the
// passed block hash (they are no longer orphans if true) and potentially
// accepts them. It repeats the process for the newly accepted blocks (
// to detect further orphans which may no longer be orphans) until there are
// no more. The flags do not modify the behavior of this function directly,
// however they are needed to pass along to maybeAcceptBlock.
// This function MUST be called with the chain state lock held (for writes).
(b *BlockChain) processOrphans(workerNumber uint32, hash *chainhash.Hash,
	flags BehaviorFlags) error {
	// Start with processing at least the passed hash.
	// Leave a little room for additional orphan blocks that need to be
	// processed without needing to grow the array in the common case.
	processHashes := make([]*chainhash.Hash, 0, 10)
	processHashes = append(processHashes, hash)
	for len(processHashes) > 0 {
		// Pop the first hash to process from the slice.
		processHash := processHashes[0]
		processHashes[0] = nil // Prevent GC leak.
		processHashes = processHashes[1:]
		// Look up all orphans that are parented by the block we just
		// accepted.  This will typically only be one,
		// but it could be multiple if multiple blocks are mined and
		// broadcast around the same time.
		// The one with the most proof of work will eventually win out.
		// An indexing for loop is intentionally used over a range here as
		// range does not reevaluate the slice on each iteration nor does it
		// adjust the index for the modified slice.
		for i := 0; i < len(b.prevOrphans[*processHash]); i++ {
			orphan := b.prevOrphans[*processHash][i]
			if orphan == nil {
				log.TRACEF("found a nil entry at index %d in the orphan"+
					" dependency list for block %v", i, processHash)
				continue
			}
			// Remove the orphan from the orphan pool.
			orphanHash := orphan.block.Hash()
			b.removeOrphanBlock(orphan)
			i--
			// Potentially accept the block into the block chain.
			_, err := b.maybeAcceptBlock(workerNumber, orphan.block, flags)
			if err != nil {
				log.ERROR(err)
				log.ERROR(err)
				return err
			}
			// Add this block to the list of blocks to process so any orphan
			// blocks that depend on this block are handled too.
			processHashes = append(processHashes, orphanHash)
		}
	}
	return nil
}
