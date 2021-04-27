package blockchain

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/hardfork"
)

// maybeAcceptBlock potentially accepts a block into the block chain
// and, if accepted, returns whether or not it is on the main chain.
// It performs several validation checks which depend on its position within
// the block chain before adding it.
// The block is expected to have already gone through ProcessBlock before
// calling this function with it.
// The flags are also passed to checkBlockContext and connectBestChain.
// See their documentation for how the flags modify their behavior.
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) maybeAcceptBlock(workerNumber uint32, block *block.Block, flags BehaviorFlags) (bool, error) {
	T.Ln("maybeAcceptBlock starting")
	// The height of this block is one more than the referenced previous block.
	prevHash := &block.WireBlock().Header.PrevBlock
	prevNode := b.Index.LookupNode(prevHash)
	if prevNode == nil {
		str := fmt.Sprintf("previous block %s is unknown", prevHash)
		E.Ln(str)
		return false, ruleError(ErrPreviousBlockUnknown, str)
	} else if b.Index.NodeStatus(prevNode).KnownInvalid() {
		str := fmt.Sprintf("previous block %s is known to be invalid", prevHash)
		E.Ln(str)
		return false, ruleError(ErrInvalidAncestorBlock, str)
	}
	blockHeight := prevNode.height + 1
	T.Ln("block not found, good, setting height", blockHeight)
	block.SetHeight(blockHeight)
	// // To deal with multiple mining algorithms, we must check first the block header version. Rather than pass the
	// // direct previous by height, we look for the previous of the same algorithm and pass that.
	// if blockHeight < b.params.BIP0034Height {
	//
	// }
	T.Ln("sanitizing header versions for legacy")
	var DoNotCheckPow bool
	var pn *BlockNode
	var a int32 = 2
	if block.WireBlock().Header.Version == 514 {
		a = 514
	}
	var aa int32 = 2
	if prevNode.version == 514 {
		aa = 514
	}
	if a != aa {
		var i int64
		pn = prevNode
		for ; i < b.params.AveragingInterval-1; i++ {
			pn = pn.GetLastWithAlgo(a)
			if pn == nil {
				break
			}
		}
	}
	T.Ln("check for blacklisted addresses")
	txs := block.Transactions()
	for i := range txs {
		if ContainsBlacklisted(b, txs[i], hardfork.Blacklist) {
			return false, ruleError(ErrBlacklisted, "block contains a blacklisted address ")
		}
	}
	T.Ln("found no blacklisted addresses")
	var e error
	if pn != nil {
		// The block must pass all of the validation rules which depend on the position
		// of the block within the block chain.
		if e = b.checkBlockContext(block, prevNode, flags, DoNotCheckPow); E.Chk(e) {
			return false, e
		}
	}
	// Insert the block into the database if it's not already there. Even though it
	// is possible the block will ultimately fail to connect, it has already passed
	// all proof-of-work and validity tests which means it would be prohibitively
	// expensive for an attacker to fill up the disk with a bunch of blocks that
	// fail to connect. This is necessary since it allows block download to be
	// decoupled from the much more expensive connection logic. It also has some
	// other nice properties such as making blocks that never become part of the
	// main chain or blocks that fail to connect available for further analysis.
	T.Ln("inserting block into database")
	if e = b.db.Update(func(dbTx database.Tx) (e error) {
		return dbStoreBlock(dbTx, block)
	},
	); E.Chk(e) {
		return false, e
	}
	// Create a new block node for the block and add it to the node index. Even if the block ultimately gets connected
	// to the main chain, it starts out on a side chain.
	blockHeader := &block.WireBlock().Header
	newNode := NewBlockNode(blockHeader, prevNode)
	newNode.status = statusDataStored
	b.Index.AddNode(newNode)
	T.Ln("flushing db")
	if e = b.Index.flushToDB(); E.Chk(e) {
		return false, e
	}
	
	// Connect the passed block to the chain while respecting proper chain selection
	// according to the chain with the most proof of work. This also handles
	// validation of the transaction scripts.
	T.Ln("connecting to best chain")
	var isMainChain bool
	if isMainChain, e = b.connectBestChain(newNode, block, flags); E.Chk(e) {
		return false, e
	}
	// Notify the caller that the new block was accepted into the block chain. The caller would typically want to react
	// by relaying the inventory to other peers.
	T.Ln("sending out block notifications for block accepted")
	b.ChainLock.Unlock()
	b.sendNotification(NTBlockAccepted, block)
	b.ChainLock.Lock()
	return isMainChain, nil
}
