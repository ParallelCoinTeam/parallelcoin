package indexers

import (
	"errors"
	"github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/gcs"
	"github.com/p9c/matrjoska/pkg/gcs/builder"
	"github.com/p9c/matrjoska/pkg/wire"
)

const (
	// cfIndexName is the human-readable name for the index.
	cfIndexName = "committed filter index"
)

// Committed filters come in one flavor currently: basic. They are generated and dropped in pairs, and both are indexed
// by a block's hash. Besides holding different content, they also live in different buckets.
var (
	// cfIndexParentBucketKey is the name of the parent bucket used to house the index. The rest of the buckets live
	// below this bucket.
	cfIndexParentBucketKey = []byte("cfindexparentbucket")
	// cfIndexKeys is an array of db bucket names used to house indexes of block hashes to cfilters.
	cfIndexKeys = [][]byte{
		[]byte("cf0byhashidx"),
	}
	// cfHeaderKeys is an array of db bucket names used to house indexes of block hashes to cf headers.
	cfHeaderKeys = [][]byte{
		[]byte("cf0headerbyhashidx"),
	}
	// cfHashKeys is an array of db bucket names used to house indexes of block hashes to cf hashes.
	cfHashKeys = [][]byte{
		[]byte("cf0hashbyhashidx"),
	}
	maxFilterType = uint8(len(cfHeaderKeys) - 1)
	// zeroHash is the chainhash.Hash value of all zero bytes, defined here for convenience.
	zeroHash chainhash.Hash
)

// dbFetchFilterIdxEntry retrieves a data blob from the filter index database. An entry's absence is not considered an
// error.
func dbFetchFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash) ([]byte, error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Get(h[:]), nil
}

// dbStoreFilterIdxEntry stores a data blob in the filter index database.
func dbStoreFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash, f []byte) (e error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Put(h[:], f)
}

// dbDeleteFilterIdxEntry deletes a data blob from the filter index database.
func dbDeleteFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash) (e error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Delete(h[:])
}

// CFIndex implements a committed filter (cf) by hash index.
type CFIndex struct {
	db          database.DB
	chainParams *chaincfg.Params
}

// Ensure the CfIndex type implements the Indexer interface.
var _ Indexer = (*CFIndex)(nil)

// Ensure the CfIndex type implements the NeedsInputser interface.
var _ NeedsInputser = (*CFIndex)(nil)

// NeedsInputs signals that the index requires the referenced inputs in order to properly create the index. This
// implements the NeedsInputser interface.
func (idx *CFIndex) NeedsInputs() bool {
	return true
}

// Init initializes the hash-based cf index. This is part of the Indexer interface.
func (idx *CFIndex) Init() (e error) {
	return nil // Nothing to do.
}

// Key returns the database key to use for the index as a byte slice. This is part of the Indexer interface.
func (idx *CFIndex) Key() []byte {
	return cfIndexParentBucketKey
}

// Name returns the human-readable name of the index. This is part of the Indexer interface.
func (idx *CFIndex) Name() string {
	return cfIndexName
}

// Create is invoked when the indexer manager determines the index needs to be created for the first time. It creates
// buckets for the two hash-based cf indexes (regular only currently).
func (idx *CFIndex) Create(dbTx database.Tx) (e error) {
	meta := dbTx.Metadata()
	cfIndexParentBucket, e := meta.CreateBucket(cfIndexParentBucketKey)
	if e != nil {
		return e
	}
	for _, bucketName := range cfIndexKeys {
		_, e = cfIndexParentBucket.CreateBucket(bucketName)
		if e != nil {
			return e
		}
	}
	for _, bucketName := range cfHeaderKeys {
		_, e = cfIndexParentBucket.CreateBucket(bucketName)
		if e != nil {
			return e
		}
	}
	for _, bucketName := range cfHashKeys {
		_, e = cfIndexParentBucket.CreateBucket(bucketName)
		if e != nil {
			return e
		}
	}
	return nil
}

// storeFilter stores a given filter, and performs the steps needed to generate the filter's header.
func storeFilter(
	dbTx database.Tx, block *block.Block, f *gcs.Filter,
	filterType wire.FilterType,
) (e error) {
	if uint8(filterType) > maxFilterType {
		return errors.New("unsupported filter type")
	}
	// Figure out which buckets to use.
	fkey := cfIndexKeys[filterType]
	hkey := cfHeaderKeys[filterType]
	hashkey := cfHashKeys[filterType]
	// Start by storing the filter.
	h := block.Hash()
	filterBytes, e := f.NBytes()
	if e != nil {
		return e
	}
	e = dbStoreFilterIdxEntry(dbTx, fkey, h, filterBytes)
	if e != nil {
		return e
	}
	// Next store the filter hash.
	filterHash, e := builder.GetFilterHash(f)
	if e != nil {
		return e
	}
	e = dbStoreFilterIdxEntry(dbTx, hashkey, h, filterHash[:])
	if e != nil {
		return e
	}
	// Then fetch the previous block's filter header.
	var prevHeader *chainhash.Hash
	ph := &block.WireBlock().Header.PrevBlock
	if ph.IsEqual(&zeroHash) {
		prevHeader = &zeroHash
	} else {
		var pfh []byte
		pfh, e = dbFetchFilterIdxEntry(dbTx, hkey, ph)
		if e != nil {
			return e
		}
		// Construct the new block's filter header, and store it.
		prevHeader, e = chainhash.NewHash(pfh)
		if e != nil {
			return e
		}
	}
	fh, e := builder.MakeHeaderForFilter(f, *prevHeader)
	if e != nil {
		return e
	}
	return dbStoreFilterIdxEntry(dbTx, hkey, h, fh[:])
}

// ConnectBlock is invoked by the index manager when a new block has been connected to the main chain. This indexer adds
// a hash-to-cf mapping for every passed block. This is part of the Indexer interface.
func (idx *CFIndex) ConnectBlock(
	dbTx database.Tx, block *block.Block,
	stxos []blockchain.SpentTxOut,
) (e error) {
	prevScripts := make([][]byte, len(stxos))
	for i, stxo := range stxos {
		prevScripts[i] = stxo.PkScript
	}
	f, e := builder.BuildBasicFilter(block.WireBlock(), prevScripts)
	if e != nil {
		return e
	}
	return storeFilter(dbTx, block, f, wire.GCSFilterRegular)
}

// DisconnectBlock is invoked by the index manager when a block has been disconnected from the main chain. This indexer
// removes the hash-to-cf mapping for every passed block. This is part of the Indexer interface.
func (idx *CFIndex) DisconnectBlock(
	dbTx database.Tx, block *block.Block,
	_ []blockchain.SpentTxOut,
) (e error) {
	for _, key := range cfIndexKeys {
		e := dbDeleteFilterIdxEntry(dbTx, key, block.Hash())
		if e != nil {
			return e
		}
	}
	for _, key := range cfHeaderKeys {
		e := dbDeleteFilterIdxEntry(dbTx, key, block.Hash())
		if e != nil {
			return e
		}
	}
	for _, key := range cfHashKeys {
		e := dbDeleteFilterIdxEntry(dbTx, key, block.Hash())
		if e != nil {
			return e
		}
	}
	return nil
}

// entryByBlockHash fetches a filter index entry of a particular type (eg. filter, filter header, etc) for a filter type
// and block hash.
func (idx *CFIndex) entryByBlockHash(
	filterTypeKeys [][]byte,
	filterType wire.FilterType, h *chainhash.Hash,
) (entry []byte, e error) {
	if uint8(filterType) > maxFilterType {
		return nil, errors.New("unsupported filter type")
	}
	key := filterTypeKeys[filterType]
	e = idx.db.View(
		func(dbTx database.Tx) (e error) {
			entry, e = dbFetchFilterIdxEntry(dbTx, key, h)
			return e
		},
	)
	return entry, e
}

// entriesByBlockHashes batch fetches a filter index entry of a particular type (eg. filter, filter header, etc) for a
// filter type and slice of block hashes.
func (idx *CFIndex) entriesByBlockHashes(
	filterTypeKeys [][]byte,
	filterType wire.FilterType, blockHashes []*chainhash.Hash,
) (entries [][]byte, e error) {
	if uint8(filterType) > maxFilterType {
		return nil, errors.New("unsupported filter type")
	}
	key := filterTypeKeys[filterType]
	entries = make([][]byte, 0, len(blockHashes))
	e = idx.db.View(
		func(dbTx database.Tx) (e error) {
			for _, blockHash := range blockHashes {
				entry, e := dbFetchFilterIdxEntry(dbTx, key, blockHash)
				if e != nil {
					return e
				}
				entries = append(entries, entry)
			}
			return nil
		},
	)
	return entries, e
}

// FilterByBlockHash returns the serialized contents of a block's basic or committed filter.
func (idx *CFIndex) FilterByBlockHash(
	h *chainhash.Hash,
	filterType wire.FilterType,
) ([]byte, error) {
	return idx.entryByBlockHash(cfIndexKeys, filterType, h)
}

// FiltersByBlockHashes returns the serialized contents of a block's basic or committed filter for a set of blocks by
// hash.
func (idx *CFIndex) FiltersByBlockHashes(
	blockHashes []*chainhash.Hash,
	filterType wire.FilterType,
) ([][]byte, error) {
	return idx.entriesByBlockHashes(cfIndexKeys, filterType, blockHashes)
}

// FilterHeaderByBlockHash returns the serialized contents of a block's basic committed filter header.
func (idx *CFIndex) FilterHeaderByBlockHash(
	h *chainhash.Hash,
	filterType wire.FilterType,
) ([]byte, error) {
	return idx.entryByBlockHash(cfHeaderKeys, filterType, h)
}

// FilterHeadersByBlockHashes returns the serialized contents of a block's basic committed filter header for a set of
// blocks by hash.
func (idx *CFIndex) FilterHeadersByBlockHashes(
	blockHashes []*chainhash.Hash,
	filterType wire.FilterType,
) ([][]byte, error) {
	return idx.entriesByBlockHashes(cfHeaderKeys, filterType, blockHashes)
}

// FilterHashByBlockHash returns the serialized contents of a block's basic committed filter hash.
func (idx *CFIndex) FilterHashByBlockHash(
	h *chainhash.Hash,
	filterType wire.FilterType,
) ([]byte, error) {
	return idx.entryByBlockHash(cfHashKeys, filterType, h)
}

// FilterHashesByBlockHashes returns the serialized contents of a block's basic committed filter hash for a set of
// blocks by hash.
func (idx *CFIndex) FilterHashesByBlockHashes(
	blockHashes []*chainhash.Hash,
	filterType wire.FilterType,
) ([][]byte, error) {
	return idx.entriesByBlockHashes(cfHashKeys, filterType, blockHashes)
}

// NewCfIndex returns a new instance of an indexer that is used to create a mapping of the hashes of all blocks in the
// blockchain to their respective committed filters. It implements the Indexer interface which plugs into the
// IndexManager that in turn is used by the blockchain package. This allows the index to be seamlessly maintained along
// with the chain.
func NewCfIndex(db database.DB, chainParams *chaincfg.Params) *CFIndex {
	return &CFIndex{db: db, chainParams: chainParams}
}

// DropCfIndex drops the CF index from the provided database if exists.
func DropCfIndex(db database.DB, interrupt qu.C) (e error) {
	return dropIndex(db, cfIndexParentBucketKey, cfIndexName, interrupt)
}
