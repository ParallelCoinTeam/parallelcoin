package blockchain

import (
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/fork"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/wire"
)

// blockStatus is a bit field representing the validation state of the block.
type blockStatus byte

const (
	// statusDataStored indicates that the block's payload is stored on disk.
	statusDataStored blockStatus = 1 << iota
	// statusValid indicates that the block has been fully validated.
	statusValid
	// statusValidateFailed indicates that the block has failed validation.
	statusValidateFailed
	// statusInvalidAncestor indicates that one of the block's ancestors has has failed validation, thus the block is
	// also invalid.
	statusInvalidAncestor
	// statusNone indicates that the block has no validation state flags set.
	//
	// NOTE: This must be defined last in order to avoid influencing iota.
	statusNone blockStatus = 0
)

// HaveData returns whether the full block data is stored in the database. This will return false for a block node where
// only the header is downloaded or kept.
func (status blockStatus) HaveData() bool {
	return status&statusDataStored != 0
}

// KnownValid returns whether the block is known to be valid. This will return false for a valid block that has not been
// fully validated yet.
func (status blockStatus) KnownValid() bool {
	return status&statusValid != 0
}

// KnownInvalid returns whether the block is known to be invalid. This may be because the block itself failed validation
// or any of its ancestors is invalid. This will return false for invalid blocks that have not been proven invalid yet.
func (status blockStatus) KnownInvalid() bool {
	return status&(statusValidateFailed|statusInvalidAncestor) != 0
}

// BlockNode represents a block within the block chain and is primarily used to aid in selecting the best chain to be
// the main chain. The main chain is stored into the block database.
type BlockNode struct {
	// NOTE: Additions, deletions, or modifications to the order of the definitions in this struct should not be changed
	// without considering how it affects alignment on 64-bit platforms. The current order is specifically crafted to
	// result in minimal padding. There will be hundreds of thousands of these in memory, so a few extra bytes of
	// padding adds up. parent is the parent block for this node.
	parent *BlockNode
	// hash is the double sha 256 of the block.
	hash chainhash.Hash
	// workSum is the total amount of work in the chain up to and including this node.
	workSum *big.Int
	// height is the position in the block chain.
	height int32
	// Some fields from block headers to aid in best chain selection and reconstructing headers from memory. These must
	// be treated as immutable and are intentionally ordered to avoid padding on 64-bit platforms.
	version    int32
	bits       uint32
	nonce      uint32
	timestamp  int64
	merkleRoot chainhash.Hash
	// status is a bitfield representing the validation state of the block. The status field, unlike the other fields,
	// may be written to and so should only be accessed using the concurrent -safe NodeStatus method on blockIndex once
	// the node has been added to the global index.
	status blockStatus
	// Diffs is the computed difficulty targets for a block to be connected to this one
	Diffs atomic.Value
}

// initBlockNode initializes a block node from the given header and parent node, calculating the height and workSum from
// the respective fields on the parent. This function is NOT safe for concurrent access. It must only be called when
// initially creating a node.
func initBlockNode(node *BlockNode, blockHeader *wire.BlockHeader, parent *BlockNode) {
	*node = BlockNode{
		hash:       blockHeader.BlockHash(),
		version:    blockHeader.Version,
		bits:       blockHeader.Bits,
		nonce:      blockHeader.Nonce,
		timestamp:  blockHeader.Timestamp.Unix(),
		merkleRoot: blockHeader.MerkleRoot,
	}
	if parent != nil {
		node.parent = parent
		node.height = parent.height + 1
		node.workSum = CalcWork(blockHeader.Bits, node.height, node.version)
		parent.workSum = CalcWork(parent.bits, parent.height, parent.version)
		node.workSum = node.workSum.Add(parent.workSum, node.workSum)
	}
}

// NewBlockNode returns a new block node for the given block header and parent node, calculating the height and workSum
// from the respective fields on the parent. This function is NOT safe for concurrent access.
func NewBlockNode(blockHeader *wire.BlockHeader, parent *BlockNode) *BlockNode {
	var node BlockNode
	initBlockNode(&node, blockHeader, parent)
	return &node
}

// Header constructs a block header from the node and returns it. This function is safe for concurrent access.
func (node *BlockNode) Header() wire.BlockHeader {
	// No lock is needed because all accessed fields are immutable.
	prevHash := &zeroHash
	if node.parent != nil {
		prevHash = &node.parent.hash
	}
	return wire.BlockHeader{
		Version:    node.version,
		PrevBlock:  *prevHash,
		MerkleRoot: node.merkleRoot,
		Timestamp:  time.Unix(node.timestamp, 0),
		Bits:       node.bits,
		Nonce:      node.nonce,
	}
}

// Ancestor returns the ancestor block node at the provided height by following the chain backwards from this node. The
// returned block will be nil when a height is requested that is after the height of the passed node or is less than
// zero. This function is safe for concurrent access.
func (node *BlockNode) Ancestor(height int32) *BlockNode {
	if height < 0 || height > node.height {
		return nil
	}
	n := node
	for ; n != nil && n.height != height; n = n.parent {
		// Intentionally left blank
	}
	return n
}

// RelativeAncestor returns the ancestor block node a relative 'distance' blocks before this node. This is equivalent to
// calling Ancestor with the node's height minus provided distance. This function is safe for concurrent access.
func (node *BlockNode) RelativeAncestor(distance int32) *BlockNode {
	return node.Ancestor(node.height - distance)
}

// CalcPastMedianTime calculates the median time of the previous few blocks prior to, and including, the block node.
// This function is safe for concurrent access.
func (node *BlockNode) CalcPastMedianTime() time.Time {
	// Create a slice of the previous few block timestamps used to calculate the median per the number defined by the
	// constant medianTimeBlocks.
	timestamps := make([]int64, medianTimeBlocks)
	numNodes := 0
	iterNode := node
	for i := 0; i < medianTimeBlocks && iterNode != nil; i++ {
		timestamps[i] = iterNode.timestamp
		numNodes++
		iterNode = iterNode.parent
	}
	// Prune the slice to the actual number of available timestamps which will be fewer than desired near the beginning
	// of the block chain and txsort them.
	timestamps = timestamps[:numNodes]
	sort.Sort(timeSorter(timestamps))
	// NOTE: The consensus rules incorrectly calculate the median for even numbers of blocks. A true median averages the
	// middle two elements for a set with an even number of elements in it. Since the constant for the previous number
	// of blocks to be used is odd, this is only an issue for a few blocks near the beginning of the chain. I suspect
	// this is an optimization even though the result is slightly wrong for a few of the first blocks since after the
	// first few blocks, there will always be an odd number of blocks in the set per the constant. This code follows
	// suit to ensure the same rules are used, however, be aware that should the medianTimeBlocks constant ever be
	// changed to an even number, this code will be wrong.
	medianTimestamp := timestamps[numNodes/2]
	return time.Unix(medianTimestamp, 0)
}

// blockIndex provides facilities for keeping track of an in-memory index of the block chain. Although the name block
// chain suggests a single chain of blocks, it is actually a tree-shaped structure where any node can have multiple
// children. However, there can only be one active branch which does indeed form a chain from the tip all the way back
// to the genesis block.
type blockIndex struct {
	// The following fields are set when the instance is created and can't be changed afterwards, so there is no need to
	// protect them with a separate mutex.
	db          database.DB
	chainParams *chaincfg.Params
	sync.RWMutex
	index map[chainhash.Hash]*BlockNode
	dirty map[*BlockNode]struct{}
}

// newBlockIndex returns a new empty instance of a block index. The index will be dynamically populated as block nodes
// are loaded from the database and manually added.
func newBlockIndex(db database.DB, chainParams *chaincfg.Params) *blockIndex {
	return &blockIndex{
		db:          db,
		chainParams: chainParams,
		index:       make(map[chainhash.Hash]*BlockNode),
		dirty:       make(map[*BlockNode]struct{}),
	}
}

// HaveBlock returns whether or not the block index contains the provided hash. This function is safe for concurrent
// access.
func (bi *blockIndex) HaveBlock(hash *chainhash.Hash) bool {
	bi.RLock()
	_, hasBlock := bi.index[*hash]
	bi.RUnlock()
	return hasBlock
}

// LookupNode returns the block node identified by the provided hash. It will return nil if there is no entry for the
// hash. This function is safe for concurrent access.
func (bi *blockIndex) LookupNode(hash *chainhash.Hash) *BlockNode {
	bi.RLock()
	node := bi.index[*hash]
	bi.RUnlock()
	return node
}

// AddNode adds the provided node to the block index and marks it as dirty. Duplicate entries are not checked so it is
// up to caller to avoid adding them. This function is safe for concurrent access.
func (bi *blockIndex) AddNode(node *BlockNode) {
	bi.Lock()
	bi.addNode(node)
	bi.dirty[node] = struct{}{}
	bi.Unlock()
}

// addNode adds the provided node to the block index, but does not mark it as dirty. This can be used while initializing
// the block index. This function is NOT safe for concurrent access.
func (bi *blockIndex) addNode(node *BlockNode) {
	bi.index[node.hash] = node
}

// NodeStatus provides concurrent-safe access to the status field of a node. This function is safe for concurrent
// access.
func (bi *blockIndex) NodeStatus(node *BlockNode) blockStatus {
	bi.RLock()
	status := node.status
	bi.RUnlock()
	return status
}

// SetStatusFlags flips the provided status flags on the block node to on, regardless of whether they were on or off
// previously. This does not unset any flags currently on. This function is safe for concurrent access.
func (bi *blockIndex) SetStatusFlags(node *BlockNode, flags blockStatus) {
	bi.Lock()
	node.status |= flags
	bi.dirty[node] = struct{}{}
	bi.Unlock()
}

// UnsetStatusFlags flips the provided status flags on the block node to off, regardless of whether they were on or off
// previously. This function is safe for concurrent access.
func (bi *blockIndex) UnsetStatusFlags(node *BlockNode, flags blockStatus) {
	bi.Lock()
	node.status &^= flags
	bi.dirty[node] = struct{}{}
	bi.Unlock()
}

// flushToDB writes all dirty block nodes to the database. If all writes succeed, this clears the dirty set.
func (bi *blockIndex) flushToDB() (e error) {
	bi.Lock()
	if len(bi.dirty) == 0 {
		bi.Unlock()
		return nil
	}
	e = bi.db.Update(
		func(dbTx database.Tx) (e error) {
			for node := range bi.dirty {
				e := dbStoreBlockNode(dbTx, node)
				if e != nil {
					E.Ln(e)
					return e
				}
			}
			return nil
		},
	)
	// If write was successful, clear the dirty set.
	if e == nil {
		bi.dirty = make(map[*BlockNode]struct{})
	}
	bi.Unlock()
	return e
}

// GetAlgo returns the algorithm of a block node
func (node *BlockNode) GetAlgo() int32 {
	return node.version
}

// GetLastWithAlgo returns the newest block from node with specified algo
func (node *BlockNode) GetLastWithAlgo(algo int32) (prev *BlockNode) {
	if node == nil {
		return
	}
	if fork.GetCurrent(node.height+1) == 0 {
		// F.Ln("checking pre-hardfork algo versions")
		if algo != 514 &&
			algo != 2 {
			D.Ln("irregular version", algo, "block, assuming 2 (sha256d)")
			algo = 2
		}
	}
	prev = node
	for {
		if prev == nil {
			return nil
		}
		// Tracef("node %d %d %8x", prev.height, prev.version, prev.bits)
		prevversion := prev.version
		if fork.GetCurrent(prev.height) == 0 {
			// F.Ln("checking pre-hardfork algo versions")
			if prev.version != 514 &&
				prev.version != 2 {
				D.Ln("irregular version block", prev.version, ", assuming 2 (sha256d)")
				prevversion = 2
			}
		}
		if prevversion == algo {
			// Tracef(
			//	"found height %d version %d prev version %d prev bits %8x",
			//	prev.height, prev.version, prevversion, prev.bits)
			return
		}
		prev = prev.RelativeAncestor(1)
	}
}

// if node == nil {
// 	F.Ln("this node is nil")
// 	return nil
// }
// prev = node.RelativeAncestor(1)
// if prev == nil {
// 	F.Ln("the previous node was nil")
// 	return nil
// }
// prevFork := fork.GetCurrent(prev.height)
// if prevFork == 0 {
// 	if algo != 514 &&
// 		algo != 2 {
// 		F.Ln("bogus version halcyon", algo)
// 		algo = 2
// 	}
// }
// if prev.version == algo {
// 	Tracef("found previous %d %d %08x", prev.height, prev.version,
// 	prev.bits)
// 	return prev
// }
// prev = prev.RelativeAncestor(1)
// for {
// 	if prev == nil {
// 		F.Ln("passed through genesis")
// 		return nil
// 	}
// 	F.Ln(prev.height)
// 	prevVersion := prev.version
// 	if fork.GetCurrent(prev.height) == 0 {
// 		if prevVersion != 514 &&
// 			prevVersion != 2 {
// 			F.Ln("bogus version", prevVersion)
// 			prevVersion = 2
// 		}
// 	}
// 	if prevVersion == algo {
// 		Tracef("found previous %d %d %08x", prev.height, prev.version,
// 			prev.bits)
// 		return prev
// 	} else {
// 		F.Ln(prev.height)
// 		prev = prev.RelativeAncestor(1)
// 	}
// }
// }
