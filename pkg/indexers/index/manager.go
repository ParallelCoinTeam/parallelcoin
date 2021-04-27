package index

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
)

var (
	// indexTipsBucketName is the name of the db bucket used to house the current tip of each index.
	indexTipsBucketName = []byte("idxtips")
)

// The index manager tracks the current tip of each index by using a parent bucket that contains an entry for index.
//
// The serialized format for an index tip is:
//
//   [<block hash><block height>],...
//   Field           Type             Size
//   block hash      chainhash.Hash   chainhash.HashSize
//   block height    uint32           4 bytes

// dbPutIndexerTip uses an existing database transaction to update or add the current tip for the given index to the
// provided values.
func dbPutIndexerTip(dbTx database.Tx, idxKey []byte, hash *chainhash.Hash, height int32) (e error) {
	serialized := make([]byte, chainhash.HashSize+4)
	copy(serialized, hash[:])
	byteOrder.PutUint32(serialized[chainhash.HashSize:], uint32(height))
	indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
	return indexesBucket.Put(idxKey, serialized)
}

// dbFetchIndexerTip uses an existing database transaction to retrieve the hash and height of the current tip for the
// provided index.
func dbFetchIndexerTip(dbTx database.Tx, idxKey []byte) (*chainhash.Hash, int32, error) {
	indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
	serialized := indexesBucket.Get(idxKey)
	if len(serialized) < chainhash.HashSize+4 {
		return nil, 0, database.DBError{
			ErrorCode: database.ErrCorruption,
			Description: fmt.Sprintf(
				"unexpected end of data for "+
					"index %q tip", string(idxKey),
			),
		}
	}
	var hash chainhash.Hash
	copy(hash[:], serialized[:chainhash.HashSize])
	height := int32(byteOrder.Uint32(serialized[chainhash.HashSize:]))
	return &hash, height, nil
}

// dbIndexConnectBlock adds all of the index entries associated with the given block using the provided indexer and
// updates the tip of the indexer accordingly. An error will be returned if the current tip for the indexer is not the
// previous block for the passed block.
func dbIndexConnectBlock(
	dbTx database.Tx, indexer Indexer, block *block.Block,
	stxo []blockchain.SpentTxOut,
) (e error) {
	// Assert that the block being connected properly connects to the current tip of the index.
	idxKey := indexer.Key()
	var curTipHash *chainhash.Hash
	if curTipHash, _, e = dbFetchIndexerTip(dbTx, idxKey); E.Chk(e) {
		return e
	}
	if !curTipHash.IsEqual(&block.WireBlock().Header.PrevBlock) {
		return AssertError(
			fmt.Sprintf(
				"dbIndexConnectBlock must be "+
					"called with a block that extends the current index "+
					"tip (%s, tip %s, block %s)", indexer.Name(),
				curTipHash, block.Hash(),
			),
		)
	}
	// Notify the indexer with the connected block so it can index it.
	if e := indexer.ConnectBlock(dbTx, block, stxo); E.Chk(e) {
		return e
	}
	// Update the current index tip.
	return dbPutIndexerTip(dbTx, idxKey, block.Hash(), block.Height())
}

// dbIndexDisconnectBlock removes all of the index entries associated with the given block using the provided indexer
// and updates the tip of the indexer accordingly. An error will be returned if the current tip for the indexer is not
// the passed block.
func dbIndexDisconnectBlock(
	dbTx database.Tx, indexer Indexer, block *block.Block,
	stxo []blockchain.SpentTxOut,
) (e error) {
	// Assert that the block being disconnected is the current tip of the index.
	idxKey := indexer.Key()
	var curTipHash *chainhash.Hash
	if curTipHash, _, e = dbFetchIndexerTip(dbTx, idxKey); E.Chk(e) {
		return e
	}
	if !curTipHash.IsEqual(block.Hash()) {
		return AssertError(
			fmt.Sprintf(
				"dbIndexDisconnectBlock must "+
					"be called with the block at the current index tip "+
					"(%s, tip %s, block %s)", indexer.Name(),
				curTipHash, block.Hash(),
			),
		)
	}
	// Notify the indexer with the disconnected block so it can remove all of
	// the appropriate entries.
	if e := indexer.DisconnectBlock(dbTx, block, stxo); E.Chk(e) {
		return e
	}
	// Update the current index tip.
	prevHash := &block.WireBlock().Header.PrevBlock
	return dbPutIndexerTip(dbTx, idxKey, prevHash, block.Height()-1)
}

// Manager defines an index manager that manages multiple optional indexes and implements the blockchain. IndexManager
// interface so it can be seamlessly plugged into normal chain processing.
type Manager struct {
	db             database.DB
	enabledIndexes []Indexer
}

// Ensure the Manager type implements the blockchain.IndexManager interface.
var _ blockchain.IndexManager = (*Manager)(nil)

// indexDropKey returns the key for an index which indicates it is in the process of being dropped.
func indexDropKey(idxKey []byte) []byte {
	dropKey := make([]byte, len(idxKey)+1)
	dropKey[0] = 'd'
	copy(dropKey[1:], idxKey)
	return dropKey
}

// maybeFinishDrops determines if each of the enabled indexes are in the middle of being dropped and finishes dropping
// them when the are. This is necessary because dropping and index has to be done in several atomic steps rather than
// one big atomic step due to the massive number of entries.
func (m *Manager) maybeFinishDrops(interrupt <-chan struct{}) (e error) {
	indexNeedsDrop := make([]bool, len(m.enabledIndexes))
	if e = m.db.View(
		func(dbTx database.Tx) (e error) {
			// None of the indexes needs to be dropped if the index tips bucket hasn't been created yet.
			indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
			if indexesBucket == nil {
				return nil
			}
			// Mark the indexer as requiring a drop if one is already in progress.
			for i, indexer := range m.enabledIndexes {
				dropKey := indexDropKey(indexer.Key())
				if indexesBucket.Get(dropKey) != nil {
					indexNeedsDrop[i] = true
				}
			}
			return nil
		},
	); E.Chk(e) {
		return e
	}
	if interruptRequested(interrupt) {
		return errInterruptRequested
	}
	// Finish dropping any of the enabled indexes that are already in the
	// middle of being dropped.
	for i, indexer := range m.enabledIndexes {
		if !indexNeedsDrop[i] {
			continue
		}
		I.C(
			func() string {
				return fmt.Sprintf("Resuming %s drop", indexer.Name())
			},
		)
		if e = dropIndex(m.db, indexer.Key(), indexer.Name(), interrupt); E.Chk(e) {
			return e
		}
	}
	return nil
}

// maybeCreateIndexes determines if each of the enabled indexes have already been created and creates them if not.
func (m *Manager) maybeCreateIndexes(dbTx database.Tx) (e error) {
	indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
	for _, indexer := range m.enabledIndexes {
		// Nothing to do if the index tip already exists.
		idxKey := indexer.Key()
		if indexesBucket.Get(idxKey) != nil {
			continue
		}
		// The tip for the index does not exist, so create it and invoke the create callback for the index so it can
		// perform any one-time initialization it requires.
		if e := indexer.Create(dbTx); E.Chk(e) {
			return e
		}
		// Set the tip for the index to values which represent an uninitialized index.
		e := dbPutIndexerTip(dbTx, idxKey, &chainhash.Hash{}, -1)
		if e != nil {
			return e
		}
	}
	return nil
}

// Init initializes the enabled indexes. This is called during chain initialization and primarily consists of catching
// up all indexes to the current best chain tip. This is necessary since each index can be disabled and re-enabled at
// any time and attempting to catch-up indexes at the same time new blocks are being downloaded would lead to an overall
// longer time to catch up due to the I/O contention. This is part of the blockchain.IndexManager interface.
func (m *Manager) Init(chain *blockchain.BlockChain, interrupt <-chan struct{}) (e error) {
	// Nothing to do when no indexes are enabled.
	if len(m.enabledIndexes) == 0 {
		return nil
	}
	if interruptRequested(interrupt) {
		return errInterruptRequested
	}
	// Finish and drops that were previously interrupted.
	if e = m.maybeFinishDrops(interrupt); E.Chk(e) {
		return e
	}
	// Create the initial state for the indexes as needed.
	e = m.db.Update(
		func(dbTx database.Tx) (e error) {
			// Create the bucket for the current tips as needed.
			meta := dbTx.Metadata()
			if _, e = meta.CreateBucketIfNotExists(indexTipsBucketName); E.Chk(e) {
				return e
			}
			return m.maybeCreateIndexes(dbTx)
		},
	)
	if E.Chk(e) {
		return e
	}
	// Initialize each of the enabled indexes.
	for _, indexer := range m.enabledIndexes {
		if e = indexer.Init(); E.Chk(e) {
			return e
		}
	}
	// Rollback indexes to the main chain if their tip is an orphaned fork. This is fairly unlikely, but it can happen
	// if the chain is reorganized while the index is disabled. This has to be done in reverse order because later
	// indexes can depend on earlier ones.
	var height int32
	var hash *chainhash.Hash
	for i := len(m.enabledIndexes); i > 0; i-- {
		indexer := m.enabledIndexes[i-1]
		// Fetch the current tip for the index.
		e = m.db.View(
			func(dbTx database.Tx) (e error) {
				idxKey := indexer.Key()
				hash, height, e = dbFetchIndexerTip(dbTx, idxKey)
				return e
			},
		)
		if e != nil {
			return e
		}
		// Nothing to do if the index does not have any entries yet.
		if height == -1 {
			continue
		}
		// Loop until the tip is a block that exists in the main chain.
		initialHeight := height
		for !chain.MainChainHasBlock(hash) {
			// At this point the index tip is orphaned, so load the orphaned block from the database directly and
			// disconnect it from the index. The block has to be loaded directly since it is no longer in the main chain
			// and thus the chain. BlockByHash function would error.
			var b *block.Block
			e = m.db.View(
				func(dbTx database.Tx) (e error) {
					blockBytes, e := dbTx.FetchBlock(hash)
					if e != nil {
						return e
					}
					b, e = block.NewFromBytes(blockBytes)
					if e != nil {
						return e
					}
					b.SetHeight(height)
					return e
				},
			)
			if e != nil {
				return e
			}
			// We'll also grab the set of outputs spent by this block so we can remove them from the index.
			var spentTxos []blockchain.SpentTxOut
			spentTxos, e = chain.FetchSpendJournal(b)
			if e != nil {
				return e
			}
			// With the block and stxo set for that block retrieved, we can now update the index itself.
			e = m.db.Update(
				func(dbTx database.Tx) (e error) {
					// Remove all of the index entries associated with the block and update the indexer tip.
					e = dbIndexDisconnectBlock(
						dbTx, indexer, b, spentTxos,
					)
					if e != nil {
						return e
					}
					// Update the tip to the previous block.
					hash = &b.WireBlock().Header.PrevBlock
					height--
					return nil
				},
			)
			if e != nil {
				return e
			}
			if interruptRequested(interrupt) {
				return errInterruptRequested
			}
		}
		if initialHeight != height {
			I.F(
				"removed %d orphaned blocks from %s (heights %d to %d)",
				initialHeight-height,
				indexer.Name(),
				height+1,
				initialHeight,
			)
		}
	}
	// Fetch the current tip heights for each index along with tracking the lowest one so the catchup code only needs to
	// start at the earliest block and is able to skip connecting the block for the indexes that don't need it.
	bestHeight := chain.BestSnapshot().Height
	lowestHeight := bestHeight
	indexerHeights := make([]int32, len(m.enabledIndexes))
	e = m.db.View(
		func(dbTx database.Tx) (e error) {
			for i, indexer := range m.enabledIndexes {
				idxKey := indexer.Key()
				_, height, e := dbFetchIndexerTip(dbTx, idxKey)
				if e != nil {
					return e
				}
				T.F(
					"current %s tip (height %d, hash %v)",
					indexer.Name(),
					height,
					hash,
				)
				indexerHeights[i] = height
				if height < lowestHeight {
					lowestHeight = height
				}
			}
			return nil
		},
	)
	if e != nil {
		return e
	}
	// Nothing to index if all of the indexes are caught up.
	if lowestHeight == bestHeight {
		return nil
	}
	// Create a progress logger for the indexing process below.
	progressLogger := newBlockProgressLogger(
		"Indexed",
	)
	// At this point, one or more indexes are behind the current best chain tip and need to be caught up, so log the
	// details and loop through each block that needs to be indexed.
	I.F(
		"catching up indexes from height %d to %d",
		lowestHeight,
		bestHeight,
	)
	for height := lowestHeight + 1; height <= bestHeight; height++ {
		// Load the block for the height since it is required to index it.
		var blk *block.Block
		blk, e = chain.BlockByHeight(height)
		if e != nil {
			return e
		}
		if interruptRequested(interrupt) {
			return errInterruptRequested
		}
		// Connect the block for all indexes that need it.
		var spentTxos []blockchain.SpentTxOut
		for i, indexer := range m.enabledIndexes {
			// Skip indexes that don't need to be updated with this block.
			if indexerHeights[i] >= height {
				continue
			}
			// When the index requires all of the referenced txouts and they haven't been loaded yet, they need to be
			// retrieved from the spend journal.
			if spentTxos == nil && indexNeedsInputs(indexer) {
				spentTxos, e = chain.FetchSpendJournal(blk)
				if e != nil {
					return e
				}
			}
			e := m.db.Update(
				func(dbTx database.Tx) (e error) {
					return dbIndexConnectBlock(
						dbTx, indexer, blk, spentTxos,
					)
				},
			)
			if e != nil {
				return e
			}
			indexerHeights[i] = height
		}
		// Log indexing progress.
		progressLogger.LogBlockHeight(blk)
		if interruptRequested(interrupt) {
			return errInterruptRequested
		}
	}
	I.Ln("indexes caught up to height", bestHeight)
	return nil
}

// indexNeedsInputs returns whether or not the index needs access to the txouts referenced by the transaction inputs
// being indexed.
func indexNeedsInputs(index Indexer) bool {
	if idx, ok := index.(NeedsInputser); ok {
		return idx.NeedsInputs()
	}
	return false
}

// // dbFetchTx looks up the passed transaction hash in the transaction index
// and loads it from the database.
// func dbFetchTx(// 	dbTx database.Tx, hash *chainhash.Hash) (*wire.MsgTx,
// error) {
// 	// Look up the location of the transaction.
// 	blockRegion, e := dbFetchTxIndexEntry(dbTx, hash)
// 	if e != nil  {
// DB// 		return nil, e
// 	}
// 	if blockRegion == nil {
// 		return nil, fmt.Errorf("transaction %v not found", hash)
// 	}
// 	// Load the raw transaction bytes from the database.
// 	txBytes, e := dbTx.FetchBlockRegion(blockRegion)
// 	if e != nil  {
//		DB// 		return nil, e
// 	}
// 	// Deserialize the transaction.
// 	var msgTx wire.MsgTx
// 	e = msgTx.Deserialize(bytes.NewReader(txBytes))
// 	if e != nil  {
//		DB// 		return nil, e
// 	}
// 	return &msgTx, nil
// }

// ConnectBlock must be invoked when a block is extending the main chain. It keeps track of the state of each index it
// is managing, performs some sanity checks, and invokes each indexer. This is part of the blockchain.IndexManager
// interface.
func (m *Manager) ConnectBlock(
	dbTx database.Tx, block *block.Block,
	stxos []blockchain.SpentTxOut,
) (e error) {
	// Call each of the currently active optional indexes with the block being connected so they can update accordingly.
	for _, index := range m.enabledIndexes {
		e := dbIndexConnectBlock(dbTx, index, block, stxos)
		if e != nil {
			return e
		}
	}
	return nil
}

// DisconnectBlock must be invoked when a block is being disconnected from the end of the main chain. It keeps track of
// the state of each index it is managing, performs some sanity checks, and invokes each indexer to remove the index
// entries associated with the block. This is part of the blockchain.IndexManager interface.
func (m *Manager) DisconnectBlock(
	dbTx database.Tx, block *block.Block,
	stxo []blockchain.SpentTxOut,
) (e error) {
	// Call each of the currently active optional indexes with the block being disconnected so they can update
	// accordingly.
	for _, index := range m.enabledIndexes {
		e := dbIndexDisconnectBlock(dbTx, index, block, stxo)
		if e != nil {
			return e
		}
	}
	return nil
}

// NewManager returns a new index manager with the provided indexes enabled. The manager returned satisfies the
// blockchain. IndexManager interface and thus cleanly plugs into the normal blockchain processing path.
func NewManager(db database.DB, enabledIndexes []Indexer) *Manager {
	return &Manager{
		db:             db,
		enabledIndexes: enabledIndexes,
	}
}

// dropIndex drops the passed index from the database. Since indexes can be massive, it deletes the index in multiple
// database transactions in order to keep memory usage to reasonable levels. It also marks the drop in progress so the
// drop can be resumed if it is stopped before it is done before the index can be used again.
func dropIndex(db database.DB, idxKey []byte, idxName string, interrupt <-chan struct{}) (e error) {
	// Nothing to do if the index doesn't already exist.
	var needsDelete bool
	e = db.View(
		func(dbTx database.Tx) (e error) {
			indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
			if indexesBucket != nil && indexesBucket.Get(idxKey) != nil {
				needsDelete = true
			}
			return nil
		},
	)
	if E.Chk(e) {
		return
	}
	if !needsDelete {
		W.F("not dropping %s because it does not exist", idxName)
		return
	}
	// Mark that the index is in the process of being dropped so that it can be resumed on the next start if interrupted
	// before the process is complete.
	I.F("dropping all %s entries.  This might take a while...", idxName)
	e = db.Update(
		func(dbTx database.Tx) (e error) {
			indexesBucket := dbTx.Metadata().Bucket(indexTipsBucketName)
			return indexesBucket.Put(indexDropKey(idxKey), idxKey)
		},
	)
	if e != nil {
		return e
	}
	// Since the indexes can be so large, attempting to simply delete the bucket in a single database transaction would
	// result in massive memory usage and likely crash many systems due to ulimits. In order to avoid this, use a cursor
	// to delete a maximum number of entries out of the bucket at a time. Recurse buckets depth-first to delete any
	// sub-buckets.
	const maxDeletions = 2000000
	var totalDeleted uint64
	// Recurse through all buckets in the index, cataloging each for later deletion.
	var subBuckets [][][]byte
	var subBucketClosure func(database.Tx, []byte, [][]byte) error
	subBucketClosure = func(
		dbTx database.Tx,
		subBucket []byte, tlBucket [][]byte,
	) (e error) {
		// Get full bucket name and append to subBuckets for later deletion.
		var bucketName [][]byte
		if (tlBucket == nil) || (len(tlBucket) == 0) {
			bucketName = append(bucketName, subBucket)
		} else {
			bucketName = append(tlBucket, subBucket)
		}
		subBuckets = append(subBuckets, bucketName)
		// Recurse sub-buckets to append to subBuckets slice.
		bucket := dbTx.Metadata()
		for _, subBucketName := range bucketName {
			bucket = bucket.Bucket(subBucketName)
		}
		return bucket.ForEachBucket(
			func(k []byte) (e error) {
				return subBucketClosure(dbTx, k, bucketName)
			},
		)
	}
	// Call subBucketClosure with top-level bucket.
	e = db.View(
		func(dbTx database.Tx) (e error) {
			return subBucketClosure(dbTx, idxKey, nil)
		},
	)
	if e != nil {
		return nil
	}
	// Iterate through each sub-bucket in reverse, deepest-first, deleting all keys inside them and then dropping the
	// buckets themselves.
	for i := range subBuckets {
		bucketName := subBuckets[len(subBuckets)-1-i]
		// Delete maxDeletions key/value pairs at a time.
		for numDeleted := maxDeletions; numDeleted == maxDeletions; {
			numDeleted = 0
			e = db.Update(
				func(dbTx database.Tx) (e error) {
					subBucket := dbTx.Metadata()
					for _, subBucketName := range bucketName {
						subBucket = subBucket.Bucket(subBucketName)
					}
					cursor := subBucket.Cursor()
					for ok := cursor.First(); ok; ok = cursor.Next() &&
						numDeleted < maxDeletions {
						if e := cursor.Delete(); E.Chk(e) {
							return e
						}
						numDeleted++
					}
					return nil
				},
			)
			if e != nil {
				return e
			}
			if numDeleted > 0 {
				totalDeleted += uint64(numDeleted)
				I.F(
					"deleted %d keys (%d total) from %s", numDeleted,
					totalDeleted, idxName,
				)
			}
		}
		if interruptRequested(interrupt) {
			return errInterruptRequested
		}
		// Drop the bucket itself.
		e = db.Update(
			func(dbTx database.Tx) (e error) {
				bucket := dbTx.Metadata()
				for j := 0; j < len(bucketName)-1; j++ {
					bucket = bucket.Bucket(bucketName[j])
				}
				return bucket.DeleteBucket(bucketName[len(bucketName)-1])
			},
		)
		if e != nil {
		}
	}
	// Call extra index specific deinitialization for the transaction index.
	if idxName == txIndexName {
		if e = dropBlockIDIndex(db); E.Chk(e) {
			return e
		}
	}
	// Remove the index tip, index bucket, and in-progress drop flag now that all index entries have been removed.
	e = db.Update(
		func(dbTx database.Tx) (e error) {
			meta := dbTx.Metadata()
			indexesBucket := meta.Bucket(indexTipsBucketName)
			if e := indexesBucket.Delete(idxKey); E.Chk(e) {
				return e
			}
			return indexesBucket.Delete(indexDropKey(idxKey))
		},
	)
	if e != nil {
		return e
	}
	I.Ln("dropped", idxName)
	return nil
}
