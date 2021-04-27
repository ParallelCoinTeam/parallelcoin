package blockchain

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/wire"
)

// blockHdrOffset defines the offsets into a v1 block index row for the block header.
//
// The serialized block index row format is:
//
//   <blocklocation><blockheader>
const blockHdrOffset = 12

// errInterruptRequested indicates that an operation was cancelled due to a user-requested interrupt.
var errInterruptRequested = errors.New("interrupt requested")

// interruptRequested returns true when the provided channel has been closed.
//
// This simplifies early shutdown slightly since the caller can just use an if statement instead of a select.
func interruptRequested(interrupted <-chan struct{}) bool {
	select {
	case <-interrupted:
		return true
	default:
	}
	return false
}

// blockChainContext represents a particular block's placement in the block chain. This is used by the block index
// migration to track block metadata that will be written to disk.
type blockChainContext struct {
	parent    *chainhash.Hash
	children  []*chainhash.Hash
	height    int32
	mainChain bool
}

// migrateBlockIndex migrates all block entries from the v1 block index bucket to the v2 bucket. The v1 bucket stores
// all block entries keyed by block hash, whereas the v2 bucket stores the exact same values, but keyed instead by block
// height + hash.
func migrateBlockIndex(db database.DB) (e error) {
	// Hardcoded bucket names so updates to the global values do not affect old upgrades.
	v1BucketName := []byte("ffldb-blockidx")
	v2BucketName := []byte("blockheaderidx")
	e = db.Update(
		func(dbTx database.Tx) (e error) {
			v1BlockIdxBucket := dbTx.Metadata().Bucket(v1BucketName)
			if v1BlockIdxBucket == nil {
				return fmt.Errorf("bucket %s does not exist", v1BucketName)
			}
			I.Ln(
				"Re-indexing block information in the database. " +
					"This might take a while",
			)
			var v2BlockIdxBucket database.Bucket
			v2BlockIdxBucket, e = dbTx.Metadata().CreateBucketIfNotExists(v2BucketName)
			if e != nil {
				return e
			}
			// Get tip of the main chain.
			serializedData := dbTx.Metadata().Get(chainStateKeyName)
			state, e := deserializeBestChainState(serializedData)
			if e != nil {
				return e
			}
			tip := &state.hash
			// Scan the old block index bucket and construct a mapping of each block to parent block and all child blocks.
			blocksMap, e := readBlockTree(v1BlockIdxBucket)
			if e != nil {
				return e
			}
			// Use the block graph to calculate the height of each block.
			e = determineBlockHeights(blocksMap)
			if e != nil {
				return e
			}
			// Find blocks on the main chain with the block graph and current tip.
			determineMainChainBlocks(blocksMap, tip)
			// Now that we have heights for all blocks, scan the old block index
			// bucket and insert all rows into the new one.
			return v1BlockIdxBucket.ForEach(
				func(hashBytes, blockRow []byte) (e error) {
					endOffset := blockHdrOffset + blockHdrSize
					headerBytes := blockRow[blockHdrOffset:endOffset:endOffset]
					var hash chainhash.Hash
					copy(hash[:], hashBytes[0:chainhash.HashSize])
					chainContext := blocksMap[hash]
					if chainContext.height == -1 {
						return fmt.Errorf(
							"unable to calculate chain height for "+
								"stored block %s", hash,
						)
					}
					// Mark blocks as valid if they are part of the main chain.
					status := statusDataStored
					if chainContext.mainChain {
						status |= statusValid
					}
					// Write header to v2 bucket
					value := make([]byte, blockHdrSize+1)
					copy(value[0:blockHdrSize], headerBytes)
					value[blockHdrSize] = byte(status)
					key := blockIndexKey(&hash, uint32(chainContext.height))
					e = v2BlockIdxBucket.Put(key, value)
					if e != nil {
						return e
					}
					// Delete header from v1 bucket
					truncatedRow := blockRow[0:blockHdrOffset:blockHdrOffset]
					return v1BlockIdxBucket.Put(hashBytes, truncatedRow)
				},
			)
		},
	)
	if e != nil {
		return e
	}
	I.Ln("Block database migration complete")
	return nil
}

// readBlockTree reads the old block index bucket and constructs a mapping of each block to its parent block and all
// child blocks. This mapping represents the full tree of blocks. This function does not populate the height or
// mainChain fields of the returned blockChainContext values.
func readBlockTree(v1BlockIdxBucket database.Bucket) (blocksMap map[chainhash.Hash]*blockChainContext, e error) {
	blocksMap = make(map[chainhash.Hash]*blockChainContext)
	e = v1BlockIdxBucket.ForEach(
		func(_, blockRow []byte) (e error) {
			var header wire.BlockHeader
			endOffset := blockHdrOffset + blockHdrSize
			headerBytes := blockRow[blockHdrOffset:endOffset:endOffset]
			e = header.Deserialize(bytes.NewReader(headerBytes))
			if e != nil {
				return e
			}
			blockHash := header.BlockHash()
			prevHash := header.PrevBlock
			if blocksMap[blockHash] == nil {
				blocksMap[blockHash] = &blockChainContext{height: -1}
			}
			if blocksMap[prevHash] == nil {
				blocksMap[prevHash] = &blockChainContext{height: -1}
			}
			blocksMap[blockHash].parent = &prevHash
			blocksMap[prevHash].children =
				append(blocksMap[prevHash].children, &blockHash)
			return nil
		},
	)
	return blocksMap, e
}

// determineBlockHeights takes a map of block hashes to a slice of child hashes and uses it to compute the height for
// each block.
//
// The function assigns a height of 0 to the genesis hash and explores the tree of blocks breadth-first, assigning a
// height to every block with a path back to the genesis block.
//
// This function modifies the height field on the blocksMap entries.
func determineBlockHeights(blocksMap map[chainhash.Hash]*blockChainContext) (e error) {
	queue := list.New()
	// The genesis block is included in blocksMap as a child of the zero hash because that is the value of the PrevBlock
	// field in the genesis header.
	preGenesisContext, exists := blocksMap[zeroHash]
	if !exists || len(preGenesisContext.children) == 0 {
		return fmt.Errorf("unable to find genesis block")
	}
	for _, genesisHash := range preGenesisContext.children {
		blocksMap[*genesisHash].height = 0
		queue.PushBack(genesisHash)
	}
	for e := queue.Front(); e != nil; e = queue.Front() {
		queue.Remove(e)
		hash := e.Value.(*chainhash.Hash)
		height := blocksMap[*hash].height
		// For each block with this one as a parent, assign it a height and
		// push to queue for future processing.
		for _, childHash := range blocksMap[*hash].children {
			blocksMap[*childHash].height = height + 1
			queue.PushBack(childHash)
		}
	}
	return nil
}

// determineMainChainBlocks traverses the block graph down from the tip to determine which block hashes that are part of
// the main chain. This function modifies the mainChain field on the blocksMap entries.
func determineMainChainBlocks(blocksMap map[chainhash.Hash]*blockChainContext, tip *chainhash.Hash) {
	for nextHash := tip; *nextHash != zeroHash; nextHash = blocksMap[*nextHash].parent {
		blocksMap[*nextHash].mainChain = true
	}
}

// deserializeUtxoEntryV0 decodes a utxo entry from the passed
// serialized byte slice according to the legacy version 0 format into a map
// of utxos keyed by the output index within the transaction.
// The map is necessary because the previous format encoded all unspent
// outputs for a transaction using a single entry,
// whereas the new format encodes each unspent output individually.
//
// The legacy format is as follows:
//
//   <version><height><header code><unspentness bitmap>[<compressed txouts>,...]
//
//   Field                Type     Size
//
//   version              VLQ      variable
//
//   block height         VLQ      variable
//
//   header code          VLQ      variable
//
//   unspentness bitmap   []byte   variable
//
//   compressed txouts
//
//     compressed amount  VLQ      variable
//     compressed script  []byte   variable
//
// The serialized header code format is:
//
//   bit 0 - containing transaction is a coinbase
//
//   bit 1 - output zero is unspent
//
//   bit 2 - output one is unspent
//
//   bits 3-x - number of bytes in unspentness bitmap.  When both bits 1 and 2
//     are unset, it encodes N-1 since there must be at least one unspent
//     output.
//
// The rationale for the header code scheme is as follows:
//
//   - Transactions which only pay to a single output and a change output are
//   extremely common, thus an extra byte for the unspentness bitmap can be
//   avoided for them by encoding those two outputs in the low order bits.
//
//   - Given it is encoded as a VLQ which can encode values up to 127 with a
//   single byte, that leaves 4 bits to represent the number of bytes in the
//   unspentness bitmap while still only consuming a single byte for the
//   header code.  In other words, an unspentness bitmap with up to 120
//   transaction outputs can be encoded with a single-byte header code.
//   This covers the vast majority of transactions.
//
//   - Encoding N-1 bytes when both bits 1 and 2 are unset allows an additional
//   8 outpoints to be encoded before causing the header code to require an
//   additional byte.
//
// Example 1:
//
// From tx in main blockchain:
//
// Blk 1, 0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098
//    010103320496b538e853519c726a2c91e61ec11600ae1390813a627c66fb8be7947be63c52
//    <><><><------------------------------------------------------------------>
//     | | \--------\                               |
//     | height     |                      compressed txout 0
//  version    header code
//
//  - version: 1
//
//  - height: 1
//
//  - header code: 0x03 (coinbase, output zero unspent, 0 bytes of unspentness)
//
//  - unspentness: Nothing since it is zero bytes
//
//  - compressed txout 0:
//
//    - 0x32: VLQ-encoded compressed amount for 5000000000 (50 DUO)
//
//    - 0x04: special script type pay-to-pubkey
//
//    - 0x96...52: x-coordinate of the pubkey
//
// Example 2:
//
// From tx in main blockchain:
//
// Blk 113931, 4a16969aa4764dd7507fc1de7f0baa4850a246de90c45e59a3207f9a26b5036f
//    0185f90b0a011200e2ccd6ec7c6e2e581349c77e067385fa8236bf8a800900b8025be1b3efc63b0ad48e7f9f10e87544528d58
//    <><----><><><------------------------------------------><-------------------------------------------->
//     |    |  | \-------------------\            |                            |
//  version |  \--------\       unspentness
//  |                    compressed txout 2
//        height     header code          compressed txout 0
//
//  - version: 1
//
//  - height: 113931
//
//  - header code: 0x0a (output zero unspent, 1 byte in unspentness bitmap)
//
//  - unspentness: [0x01] (bit 0 is set, so output 0+2 = 2 is unspent)
//    NOTE: It's +2 since the first two outputs are encoded in the header code
//
//  - compressed txout 0:
//
//    - 0x12: VLQ-encoded compressed amount for 20000000 (0.2 DUO)
//
//    - 0x00: special script type pay-to-pubkey-hash
//
//    - 0xe2...8a: pubkey hash
//
//  - compressed txout 2:
//
//    - 0x8009: VLQ-encoded compressed amount for 15000000 (0.15 DUO)
//
//    - 0x00: special script type pay-to-pubkey-hash
//
//    - 0xb8...58: pubkey hash
//
// Example 3:
//
// From tx in main blockchain:
//
// Blk 338156, 1b02d1c8cfef60a189017b9a420c682cf4a0028175f2f563209e4ff61c8c3620
//    0193d06c100000108ba5b9e763011dd46a006572d820e448e12d2bbb38640bc718e6
//    <><----><><----><-------------------------------------------------->
//     |    |  |   \-----------------\            |
//  version |  \--------\       unspentness       |
//        height     header code          compressed txout 22
//
//  - version: 1
//
//  - height: 338156
//
//  - header code: 0x10 (2+1 = 3 bytes in unspentness bitmap) NOTE: It's +1 since neither bit 1 nor 2 are set, so N-1 is
//  encoded.
//
//  - unspentness: [0x00 0x00 0x10] (bit 20 is set, so output 20+2 = 22 is unspent)
//    NOTE: It's +2 since the first two outputs are encoded in the header code
//
//  - compressed txout 22:
//
//    - 0x8ba5b9e763: VLQ-encoded compressed amount for 366875659 (3.66875659 DUO)
//
//    - 0x01: special script type pay-to-script-hash
//
//    - 0x1d...e6: script hash
func deserializeUtxoEntryV0(serialized []byte) (map[uint32]*UtxoEntry, error) {
	// Deserialize the version.
	// NOTE: Ignore version since it is no longer used in the new format.
	_, bytesRead := deserializeVLQ(serialized)
	offset := bytesRead
	if offset >= len(serialized) {
		return nil, errDeserialize("unexpected end of data after version")
	}
	// Deserialize the block height.
	blockHeight, bytesRead := deserializeVLQ(serialized[offset:])
	offset += bytesRead
	if offset >= len(serialized) {
		return nil, errDeserialize("unexpected end of data after height")
	}
	// Deserialize the header code.
	code, bytesRead := deserializeVLQ(serialized[offset:])
	offset += bytesRead
	if offset >= len(serialized) {
		return nil, errDeserialize("unexpected end of data after header")
	}
	// Decode the header code.
	//
	// Bit 0 indicates whether the containing transaction is a coinbase.
	//
	// Bit 1 indicates output 0 is unspent.
	//
	// Bit 2 indicates output 1 is unspent.
	//
	// Bits 3-x encodes the number of non-zero unspentness bitmap bytes that follow. When both output 0 and 1 are spent,
	// it encodes N-1.
	isCoinBase := code&0x01 != 0
	output0Unspent := code&0x02 != 0
	output1Unspent := code&0x04 != 0
	numBitmapBytes := code >> 3
	if !output0Unspent && !output1Unspent {
		numBitmapBytes++
	}
	// Ensure there are enough bytes left to deserialize the unspentness bitmap.
	if uint64(len(serialized[offset:])) < numBitmapBytes {
		return nil, errDeserialize(
			"unexpected end of data for " +
				"unspentness bitmap",
		)
	}
	// Add sparse output for unspent outputs 0 and 1 as needed based on the details provided by the header code.
	var outputIndexes []uint32
	if output0Unspent {
		outputIndexes = append(outputIndexes, 0)
	}
	if output1Unspent {
		outputIndexes = append(outputIndexes, 1)
	}
	// Decode the unspentness bitmap adding a sparse output for each unspent output.
	for i := uint32(0); i < uint32(numBitmapBytes); i++ {
		unspentBits := serialized[offset]
		for j := uint32(0); j < 8; j++ {
			if unspentBits&0x01 != 0 {
				// The first 2 outputs are encoded via the header code, so adjust the output number accordingly.
				outputNum := 2 + i*8 + j
				outputIndexes = append(outputIndexes, outputNum)
			}
			unspentBits >>= 1
		}
		offset++
	}
	// Map to hold all of the converted outputs.
	entries := make(map[uint32]*UtxoEntry)
	// All entries will need to potentially be marked as a coinbase.
	var packedFlags txoFlags
	if isCoinBase {
		packedFlags |= tfCoinBase
	}
	// Decode and add all of the utxos.
	for i, outputIndex := range outputIndexes {
		// Decode the next utxo.
		amount, pkScript, bytesRead, e := decodeCompressedTxOut(
			serialized[offset:],
		)
		if e != nil {
			return nil, errDeserialize(
				fmt.Sprintf(
					"unable to "+
						"decode utxo at index %d: %v", i, e,
				),
			)
		}
		offset += bytesRead
		// Create a new utxo entry with the details deserialized above.
		entries[outputIndex] = &UtxoEntry{
			amount:      int64(amount),
			pkScript:    pkScript,
			blockHeight: int32(blockHeight),
			packedFlags: packedFlags,
		}
	}
	return entries, nil
}

// upgradeUtxoSetToV2 migrates the utxo set entries from version 1 to 2 in batches. It is guaranteed to updated if this
// returns without failure.
func upgradeUtxoSetToV2(db database.DB, interrupt <-chan struct{}) (e error) {
	// Hardcoded bucket names so updates to the global values do not affect old upgrades.
	var (
		v1BucketName = []byte("utxoset")
		v2BucketName = []byte("utxosetv2")
	)
	I.Ln("Upgrading utxo set to v2.  This will take a while")
	start := time.Now()
	// Create the new utxo set bucket as needed.
	e = db.Update(
		func(dbTx database.Tx) (e error) {
			_, e = dbTx.Metadata().CreateBucketIfNotExists(v2BucketName)
			return e
		},
	)
	if e != nil {
		return e
	}
	// doBatch contains the primary logic for upgrading the utxo set from version 1 to 2 in batches. This is done
	// because the utxo set can be huge and thus attempting to migrate in a single database transaction would result in
	// massive memory usage and could potentially crash on many systems due to ulimits. It returns the number of utxos
	// processed.
	const maxUtxos = 200000
	doBatch := func(dbTx database.Tx) (uint32, error) {
		v1Bucket := dbTx.Metadata().Bucket(v1BucketName)
		v2Bucket := dbTx.Metadata().Bucket(v2BucketName)
		v1Cursor := v1Bucket.Cursor()
		// Migrate utxos so long as the max number of utxos for this batch has not been exceeded.
		var numUtxos uint32
		for ok := v1Cursor.First(); ok && numUtxos < maxUtxos; ok =
			v1Cursor.Next() {
			// Old key was the transaction hash.
			oldKey := v1Cursor.Key()
			var txHash chainhash.Hash
			copy(txHash[:], oldKey)
			// Deserialize the old entry which included all utxos for the given transaction.
			var utxos map[uint32]*UtxoEntry
			utxos, e = deserializeUtxoEntryV0(v1Cursor.Value())
			if e != nil {
				return 0, e
			}
			// Add an entry for each utxo into the new bucket using the new format.
			for txOutIdx, utxo := range utxos {
				var reserialized []byte
				reserialized, e = serializeUtxoEntry(utxo)
				if e != nil {
					return 0, e
				}
				key := outpointKey(
					wire.OutPoint{
						Hash:  txHash,
						Index: txOutIdx,
					},
				)
				e = v2Bucket.Put(*key, reserialized)
				// NOTE: The key is intentionally not recycled here since the database interface contract prohibits
				// modifications. It will be garbage collected normally when the database is done with it.
				if e != nil {
					return 0, e
				}
			}
			// Remove old entry.
			e = v1Bucket.Delete(oldKey)
			if e != nil {
				return 0, e
			}
			numUtxos += uint32(len(utxos))
			if interruptRequested(interrupt) {
				// No error here so the database transaction is not cancelled and therefore outstanding work is written
				// to disk.
				break
			}
		}
		return numUtxos, nil
	}
	// Migrate all entries in batches for the reasons mentioned above.
	var totalUtxos uint64
	for {
		var numUtxos uint32
		e = db.Update(
			func(dbTx database.Tx) (e error) {
				numUtxos, e = doBatch(dbTx)
				return e
			},
		)
		if e != nil {
			return e
		}
		if interruptRequested(interrupt) {
			return errInterruptRequested
		}
		if numUtxos == 0 {
			break
		}
		totalUtxos += uint64(numUtxos)
		I.F("migrated %d utxos (%d total)", numUtxos, totalUtxos)
	}
	// Remove the old bucket and update the utxo set version once it has been fully migrated.
	e = db.Update(
		func(dbTx database.Tx) (e error) {
			e = dbTx.Metadata().DeleteBucket(v1BucketName)
			if e != nil {
				return e
			}
			return dbPutVersion(dbTx, utxoSetVersionKeyName, 2)
		},
	)
	if e != nil {
		return e
	}
	seconds := int64(time.Since(start) / time.Second)
	I.F(
		"Done upgrading utxo set.  Total utxos: %d in %d seconds",
		totalUtxos, seconds,
	)
	return nil
}

// maybeUpgradeDbBuckets checks the database version of the buckets used by this package and performs any needed
// upgrades to bring them to the latest version. All buckets used by this package are guaranteed to be the latest
// version if this function returns without error.
func (b *BlockChain) maybeUpgradeDbBuckets(interrupt <-chan struct{}) (e error) {
	// Load or create bucket versions as needed.
	var utxoSetVersion uint32
	e = b.db.Update(
		func(dbTx database.Tx) (e error) {
			// Load the utxo set version from the database or create it and initialize it to version 1 if it doesn't exist.
			utxoSetVersion, e = dbFetchOrCreateVersion(
				dbTx,
				utxoSetVersionKeyName, 1,
			)
			return e
		},
	)
	if e != nil {
		return e
	}
	// Update the utxo set to v2 if needed.
	if utxoSetVersion < 2 {
		if e := upgradeUtxoSetToV2(b.db, interrupt); E.Chk(e) {
			return e
		}
	}
	return nil
}
