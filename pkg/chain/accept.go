package blockchain

import (
	"fmt"

	"github.com/p9c/pod/pkg/chain/hardfork"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/log"

	"github.com/p9c/pod/pkg/util"
)

func // maybeAcceptBlock potentially accepts a block into the block chain
// and, if accepted, returns whether or not it is on the main chain.
// It performs several validation checks which depend on its position within
// the block chain before adding it.
// The block is expected to have already gone through ProcessBlock before
// calling this function with it.
// The flags are also passed to checkBlockContext and connectBestChain.
// See their documentation for how the flags modify their behavior.
// This function MUST be called with the chain state lock held (for writes).
(b *BlockChain) maybeAcceptBlock(workerNumber uint32, block *util.Block,
	flags BehaviorFlags) (bool, error) {
	// log.WARN("maybeAcceptBlock")
	// The height of this block is one more than the referenced previous block.
	prevHash := &block.MsgBlock().Header.PrevBlock
	prevNode := b.Index.LookupNode(prevHash)
	if prevNode == nil {
		str := fmt.Sprintf("previous block %s is unknown", prevHash)
		log.WARN(str)
		return false, ruleError(ErrPreviousBlockUnknown, str)
	} else if b.Index.NodeStatus(prevNode).KnownInvalid() {
		str := fmt.Sprintf("previous block %s is known to be invalid", prevHash)
		log.WARN(str)
		return false, ruleError(ErrInvalidAncestorBlock, str)
	}
	blockHeight := prevNode.height + 1
	block.SetHeight(blockHeight)
	// To deal with multiple mining algorithms,
	// we must check first the block header version.
	// Rather than pass the direct previous by height,
	// we look for the previous of the same algorithm and pass that.
	var DoNotCheckPow bool
	var pn *blockNode
	var a int32 = 2
	if block.MsgBlock().Header.Version == 514 {
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
	// log.WARN("check for blacklisted addresses")
	txs := block.Transactions()
	for i := range txs {
		if ContainsBlacklisted(b, txs[i], hardfork.Blacklist) {
			return false, ruleError(ErrBlacklisted, "block contains a blacklisted address ")
		}
	}
	// log.WARN("found no blacklisted addresses")
	var err error
	if pn != nil {
		// The block must pass all of the validation rules which depend on
		// the  position of the block within the block chain.
		err = b.checkBlockContext(workerNumber, block, prevNode, flags,
			DoNotCheckPow)
		if err != nil {
			log.ERROR(err)
			return false, err
		}
	}
	// Insert the block into the database if it's not already there.
	// Even though it is possible the block will ultimately fail to connect,
	// it has already passed all proof-of-work and validity tests which means
	// it would be prohibitively expensive for an attacker to fill up the
	// disk with a bunch of blocks that fail to connect.
	// This is necessary since it allows block download to be decoupled from
	// the much more expensive connection logic.
	// It also has some other nice properties such as making blocks that
	// never become part of the main chain or blocks that fail to connect
	// available for further analysis.
	err = b.db.Update(func(dbTx database.Tx) error {
		return dbStoreBlock(dbTx, block)
	})
	if err != nil {
		log.ERROR(err)
		return false, err
	}
	// log.WARN("creating new block node for new block")
	// Create a new block node for the block and add it to the node index.
	// Even if the block ultimately gets connected to the main chain,
	// it starts out on a side chain.
	blockHeader := &block.MsgBlock().Header
	newNode := newBlockNode(blockHeader, prevNode)
	newNode.status = statusDataStored
	b.Index.AddNode(newNode)
	err = b.Index.flushToDB()
	if err != nil {
		log.ERROR(err)
		return false, err
	}
	// Connect the passed block to the chain while respecting proper chain
	// selection according to the chain with the most proof of work.
	// This also handles validation of the transaction scripts.
	isMainChain, err := b.connectBestChain(newNode, block, flags)
	if err != nil {
		log.ERROR(err)
		return false, err
	}
	// Notify the caller that the new block was accepted into the block
	// chain.  The caller would typically want to react by relaying the
	// inventory to other peers.
	b.chainLock.Unlock()
	b.sendNotification(NTBlockAccepted, block)
	b.chainLock.Lock()
	return isMainChain, nil
}
