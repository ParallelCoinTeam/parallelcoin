package index

import (
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/wire"
)

const (
	// txIndexName is the human-readable name for the index.
	txIndexName = "transaction index"
)

var (
	// txIndexKey is the key of the transaction index and the db bucket used
	// to house it.
	txIndexKey = []byte("txbyhashidx")
	// idByHashIndexBucketName is the name of the db bucket used to house the
	// block id -> block hash index.
	idByHashIndexBucketName = []byte("idbyhashidx")
	// hashByIDIndexBucketName is the name of the db bucket used to house the
	// block hash -> block id index.
	hashByIDIndexBucketName = []byte("hashbyididx")
	// errNoBlockIDEntry is an error that indicates a requested entry does
	// not exist in the block ID index.
	errNoBlockIDEntry = errors.New("no entry in the block ID index")
)

// The transaction index consists of an entry for every transaction in the main chain. In order to significantly
// optimize the space requirements a separate index which provides an internal mapping between each block that has been
// indexed and a unique ID for use within the hash to location mappings. The ID is simply a sequentially incremented
// uint32.
//
// This is useful because it is only 4 bytes versus 32 bytes hashes and thus saves a ton of space in the index. There
// are three buckets used in total.
//
// The first bucket maps the hash of each transaction to the specific block location. The second bucket maps the hash of
// each block to the unique ID and the third maps that ID back to the block hash.
//
// NOTE: Although it is technically possible for multiple transactions to have the same hash as long as the previous
// transaction with the same hash is fully spent, this code only stores the most recent one because doing otherwise
// would add a non-trivial amount of space and overhead for something that will realistically never happen per the
// probability and even if it did, the old one must be fully spent and so the most likely transaction a caller would
// want for a given hash is the most recent one anyways.
//
// The serialized format for keys and values in the block hash to ID bucket is:
//
//   <hash> = <ID>
//   Field           Type              Size
//   hash            chainhash.Hash    32 bytes
//   ID              uint32            4 bytes
//   -----
//   Total: 36 bytes
//
// The serialized format for keys and values in the ID to block hash bucket is:
//
//   <ID> = <hash>
//   Field           Type              Size
//   ID              uint32            4 bytes
//   hash            chainhash.Hash    32 bytes
//   -----
//   Total: 36 bytes
// The serialized format for the keys and values in the tx index bucket is:
//   <txhash> = <block id><start offset><tx length>
//   Field           Type              Size
//   txhash          chainhash.Hash    32 bytes
//   block id        uint32            4 bytes
//   start offset    uint32          4 bytes
//   tx length       uint32          4 bytes
//   -----
//   Total: 44 bytes
//
// dbPutBlockIDIndexEntry uses an existing database transaction to update or add the index entries for the hash to id
// and id to hash mappings for the provided values.
func dbPutBlockIDIndexEntry(dbTx database.Tx, hash *chainhash.Hash, id uint32) (e error) {
	// Serialize the height for use in the index entries.
	var serializedID [4]byte
	byteOrder.PutUint32(serializedID[:], id)
	// Add the block hash to ID mapping to the index.
	meta := dbTx.Metadata()
	hashIndex := meta.Bucket(idByHashIndexBucketName)
	if e := hashIndex.Put(hash[:], serializedID[:]); E.Chk(e) {
		return e
	}
	// Add the block ID to hash mapping to the index.
	idIndex := meta.Bucket(hashByIDIndexBucketName)
	return idIndex.Put(serializedID[:], hash[:])
}

// dbRemoveBlockIDIndexEntry uses an existing database transaction remove index entries from the hash to id and id to
// hash mappings for the provided hash.
func dbRemoveBlockIDIndexEntry(dbTx database.Tx, hash *chainhash.Hash) (e error) {
	// Remove the block hash to ID mapping.
	meta := dbTx.Metadata()
	hashIndex := meta.Bucket(idByHashIndexBucketName)
	serializedID := hashIndex.Get(hash[:])
	if serializedID == nil {
		return nil
	}
	if e := hashIndex.Delete(hash[:]); E.Chk(e) {
		return e
	}
	// Remove the block ID to hash mapping.
	idIndex := meta.Bucket(hashByIDIndexBucketName)
	return idIndex.Delete(serializedID)
}

// dbFetchBlockIDByHash uses an existing database transaction to retrieve the block id for the provided hash from the
// index.
func dbFetchBlockIDByHash(dbTx database.Tx, hash *chainhash.Hash) (uint32, error) {
	hashIndex := dbTx.Metadata().Bucket(idByHashIndexBucketName)
	serializedID := hashIndex.Get(hash[:])
	if serializedID == nil {
		return 0, errNoBlockIDEntry
	}
	return byteOrder.Uint32(serializedID), nil
}

// dbFetchBlockHashBySerializedID uses an existing database transaction to retrieve the hash for the provided serialized
// block id from the index.
func dbFetchBlockHashBySerializedID(dbTx database.Tx, serializedID []byte) (*chainhash.Hash, error) {
	idIndex := dbTx.Metadata().Bucket(hashByIDIndexBucketName)
	hashBytes := idIndex.Get(serializedID)
	if hashBytes == nil {
		return nil, errNoBlockIDEntry
	}
	var hash chainhash.Hash
	copy(hash[:], hashBytes)
	return &hash, nil
}

// dbFetchBlockHashByID uses an existing database transaction to retrieve the hash for the provided block id from the
// index.
func dbFetchBlockHashByID(dbTx database.Tx, id uint32) (*chainhash.Hash, error) {
	var serializedID [4]byte
	byteOrder.PutUint32(serializedID[:], id)
	return dbFetchBlockHashBySerializedID(dbTx, serializedID[:])
}

// putTxIndexEntry serializes the provided values according to the format described about for a transaction index entry.
// The target byte slice must be at least large enough to handle the number of bytes defined by the txEntrySize constant
// or it will panic.
func putTxIndexEntry(target []byte, blockID uint32, txLoc wire.TxLoc) {
	byteOrder.PutUint32(target, blockID)
	byteOrder.PutUint32(target[4:], uint32(txLoc.TxStart))
	byteOrder.PutUint32(target[8:], uint32(txLoc.TxLen))
}

// dbPutTxIndexEntry uses an existing database transaction to update the transaction index given the provided serialized
// data that is expected to have been serialized putTxIndexEntry.
func dbPutTxIndexEntry(dbTx database.Tx, txHash *chainhash.Hash, serializedData []byte) (e error) {
	txIndex := dbTx.Metadata().Bucket(txIndexKey)
	return txIndex.Put(txHash[:], serializedData)
}

// dbFetchTxIndexEntry uses an existing database transaction to fetch the block region for the provided transaction hash
// from the transaction index. When there is no entry for the provided hash, nil will be returned for the both the
// region and the error.
func dbFetchTxIndexEntry(dbTx database.Tx, txHash *chainhash.Hash) (*database.BlockRegion, error) {
	// Load the record from the database and return now if it doesn't exist.
	txIndex := dbTx.Metadata().Bucket(txIndexKey)
	serializedData := txIndex.Get(txHash[:])
	if len(serializedData) == 0 {
		return nil, nil
	}
	// Ensure the serialized data has enough bytes to properly deserialize.
	if len(serializedData) < 12 {
		return nil, database.DBError{
			ErrorCode: database.ErrCorruption,
			Description: fmt.Sprintf(
				"corrupt transaction index "+
					"entry for %s", txHash,
			),
		}
	}
	// Load the block hash associated with the block ID.
	hash, e := dbFetchBlockHashBySerializedID(dbTx, serializedData[0:4])
	if e != nil {
		return nil, database.DBError{
			ErrorCode: database.ErrCorruption,
			Description: fmt.Sprintf(
				"corrupt transaction index "+
					"entry for %s: %v", txHash, e,
			),
		}
	}
	// Deserialize the final entry.
	region := database.BlockRegion{Hash: &chainhash.Hash{}}
	copy(region.Hash[:], hash[:])
	region.Offset = byteOrder.Uint32(serializedData[4:8])
	region.Len = byteOrder.Uint32(serializedData[8:12])
	return &region, nil
}

// dbAddTxIndexEntries uses an existing database transaction to add a transaction index entry for every transaction in
// the passed block.
func dbAddTxIndexEntries(dbTx database.Tx, block *block.Block, blockID uint32) (e error) {
	// The offset and length of the transactions within the serialized block.
	txLocs, e := block.TxLoc()
	if e != nil {
		return e
	}
	// As an optimization, allocate a single slice big enough to hold all of the serialized transaction index entries
	// for the block and serialize them directly into the slice. Then, pass the appropriate subslice to the database to
	// be written. This approach significantly cuts down on the number of required allocations.
	offset := 0
	serializedValues := make([]byte, len(block.Transactions())*txEntrySize)
	for i, tx := range block.Transactions() {
		putTxIndexEntry(serializedValues[offset:], blockID, txLocs[i])
		endOffset := offset + txEntrySize
		e := dbPutTxIndexEntry(
			dbTx, tx.Hash(),
			serializedValues[offset:endOffset:endOffset],
		)
		if e != nil {
			return e
		}
		offset += txEntrySize
	}
	return nil
}

// dbRemoveTxIndexEntry uses an existing database transaction to remove the most recent transaction index entry for the
// given hash.
func dbRemoveTxIndexEntry(dbTx database.Tx, txHash *chainhash.Hash) (e error) {
	txIndex := dbTx.Metadata().Bucket(txIndexKey)
	serializedData := txIndex.Get(txHash[:])
	if len(serializedData) == 0 {
		return fmt.Errorf(
			"can't remove non-existent transaction %s "+
				"from the transaction index", txHash,
		)
	}
	return txIndex.Delete(txHash[:])
}

// dbRemoveTxIndexEntries uses an existing database transaction to remove the latest transaction entry for every
// transaction in the passed block.
func dbRemoveTxIndexEntries(dbTx database.Tx, block *block.Block) (e error) {
	for _, tx := range block.Transactions() {
		e := dbRemoveTxIndexEntry(dbTx, tx.Hash())
		if e != nil {
			return e
		}
	}
	return nil
}

// TxIndex implements a transaction by hash index. That is to say, it supports querying all transactions by their hash.
type TxIndex struct {
	db         database.DB
	curBlockID uint32
}

// Ensure the TxIndex type implements the Indexer interface.
var _ Indexer = (*TxIndex)(nil)

// Init initializes the hash-based transaction index. In particular, it finds the highest used block ID and stores it
// for later use when connecting or disconnecting blocks. This is part of the Indexer interface.
func (idx *TxIndex) Init() (e error) {
	// Find the latest known block id field for the internal block id index and initialize it. This is done because it's
	// a lot more efficient to do a single search at initialize time than it is to write another value to the database
	// on every update.
	e = idx.db.View(
		func(dbTx database.Tx) (e error) {
			// Scan forward in large gaps to find a block id that doesn't exist yet to serve as an upper bound for the
			// binary search below.
			var highestKnown, nextUnknown uint32
			testBlockID := uint32(1)
			increment := uint32(100000)
			for {
				_, e := dbFetchBlockHashByID(dbTx, testBlockID)
				if e != nil {
					// F.Ln(err)
					nextUnknown = testBlockID
					break
				}
				highestKnown = testBlockID
				testBlockID += increment
			}
			T.F("forward scan (highest known %d, next unknown %d)", highestKnown, nextUnknown)
			// No used block IDs due to new database.
			if nextUnknown == 1 {
				return nil
			}
			// Use a binary search to find the final highest used block id. This will take at most ceil(log_2(increment))
			// attempts.
			for {
				testBlockID = (highestKnown + nextUnknown) / 2
				_, e := dbFetchBlockHashByID(dbTx, testBlockID)
				if e != nil {
					// F.Ln(err)
					nextUnknown = testBlockID
				} else {
					highestKnown = testBlockID
				}
				T.F("binary scan (highest known %d, next unknown %d)", highestKnown, nextUnknown)
				if highestKnown+1 == nextUnknown {
					break
				}
			}
			idx.curBlockID = highestKnown
			return nil
		},
	)
	if e != nil {
		return e
	}
	F.Ln("current internal block ID:", idx.curBlockID)
	return nil
}

// Key returns the database key to use for the index as a byte slice. This is part of the Indexer interface.
func (idx *TxIndex) Key() []byte {
	return txIndexKey
}

// Name returns the human-readable name of the index. This is part of the Indexer interface.
func (idx *TxIndex) Name() string {
	return txIndexName
}

// Create is invoked when the indexer manager determines the index needs to be created for the first time. It creates
// the buckets for the hash-based transaction index and the internal block ID indexes. This is part of the Indexer
// interface.
func (idx *TxIndex) Create(dbTx database.Tx) (e error) {
	meta := dbTx.Metadata()
	if _, e = meta.CreateBucket(idByHashIndexBucketName); E.Chk(e) {
		return e
	}
	if _, e = meta.CreateBucket(hashByIDIndexBucketName); E.Chk(e) {
		return e
	}
	_, e = meta.CreateBucket(txIndexKey)
	return e
}

// ConnectBlock is invoked by the index manager when a new block has been connected to the main chain. This indexer adds
// a hash-to-transaction mapping for every transaction in the passed block. This is part of the Indexer interface.
func (idx *TxIndex) ConnectBlock(dbTx database.Tx, block *block.Block, stxos []blockchain.SpentTxOut) (e error) {
	// Increment the internal block ID to use for the block being connected and add all of the transactions in the block
	// to the index.
	newBlockID := idx.curBlockID + 1
	if e = dbAddTxIndexEntries(dbTx, block, newBlockID); E.Chk(e) {
		return e
	}
	// Add the new block ID index entry for the block being connected and update the current internal block ID
	// accordingly.
	e = dbPutBlockIDIndexEntry(dbTx, block.Hash(), newBlockID)
	if e != nil {
		return e
	}
	idx.curBlockID = newBlockID
	return nil
}

// DisconnectBlock is invoked by the index manager when a block has been disconnected from the main chain.
//
// This indexer removes the hash-to-transaction mapping for every transaction in the block.
//
// This is part of the Indexer interface.
func (idx *TxIndex) DisconnectBlock(
	dbTx database.Tx, block *block.Block,
	stxos []blockchain.SpentTxOut,
) (e error) {
	// Remove all of the transactions in the block from the index.
	if e = dbRemoveTxIndexEntries(dbTx, block); E.Chk(e) {
		return e
	}
	// Remove the block ID index entry for the block being disconnected and decrement the current internal block ID to
	// account for it.
	if e = dbRemoveBlockIDIndexEntry(dbTx, block.Hash()); E.Chk(e) {
		return e
	}
	idx.curBlockID--
	return nil
}

// TxBlockRegion returns the block region for the provided transaction hash from the transaction index.
//
// The block region can in turn be used to load the raw transaction bytes.
//
// When there is no entry for the provided hash, nil will be returned for the both the entry and the error.
//
// This function is safe for concurrent access.
func (idx *TxIndex) TxBlockRegion(hash *chainhash.Hash) (region *database.BlockRegion, e error) {
	e = idx.db.View(
		func(dbTx database.Tx) (e error) {
			region, e = dbFetchTxIndexEntry(dbTx, hash)
			return e
		},
	)
	return region, e
}

// NewTxIndex returns a new instance of an indexer that is used to create a mapping of the hashes of all transactions in
// the blockchain to the respective block, location within the block, and size of the transaction.
//
// It implements the Indexer interface which plugs into the IndexManager that in turn is used by the blockchain package.
//
// This allows the index to be seamlessly maintained along with the chain.
func NewTxIndex(db database.DB) *TxIndex {
	return &TxIndex{db: db}
}

// dropBlockIDIndex drops the internal block id index.
func dropBlockIDIndex(db database.DB) (e error) {
	return db.Update(
		func(dbTx database.Tx) (e error) {
			meta := dbTx.Metadata()
			e = meta.DeleteBucket(idByHashIndexBucketName)
			if e != nil {
				return e
			}
			return meta.DeleteBucket(hashByIDIndexBucketName)
		},
	)
}

// DropTxIndex drops the transaction index from the provided database if it exists. Since the address index relies on
// it, the address index will also be dropped when it exists.
func DropTxIndex(db database.DB, interrupt qu.C) (e error) {
	e = dropIndex(db, addrIndexKey, addrIndexName, interrupt)
	if e != nil {
		return e
	}
	return dropIndex(db, txIndexKey, txIndexName, interrupt)
}
