package indexers

import (
	"errors"
	"github.com/stalker-loki/app/slog"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/coding/gcs"
	"github.com/p9c/pod/pkg/coding/gcs/builder"
	database "github.com/p9c/pod/pkg/db"
	"github.com/p9c/pod/pkg/util"
)

const (
	// cfIndexName is the human-readable name for the index.
	cfIndexName = "committed filter index"
)

// Committed filters come in one flavor currently: basic. They are generated and dropped in pairs, and both are indexed by a block's hash.  Besides holding different content, they also live in different buckets.
var (
	// cfIndexParentBucketKey is the name of the parent bucket used to house the index. The rest of the buckets live below this bucket.
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

// dbFetchFilterIdxEntry retrieves a data blob from the filter index database. An entry's absence is not considered an error.
func dbFetchFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash) (entry []byte, err error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Get(h[:]), nil
}

// dbStoreFilterIdxEntry stores a data blob in the filter index database.
func dbStoreFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash, f []byte) (err error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Put(h[:], f)
}

// dbDeleteFilterIdxEntry deletes a data blob from the filter index database.
func dbDeleteFilterIdxEntry(dbTx database.Tx, key []byte, h *chainhash.Hash) (err error) {
	idx := dbTx.Metadata().Bucket(cfIndexParentBucketKey).Bucket(key)
	return idx.Delete(h[:])
}

// CFIndex implements a committed filter (cf) by hash index.
type CFIndex struct {
	db          database.DB
	chainParams *netparams.Params
}

// Ensure the CfIndex type implements the Indexer interface.
var _ Indexer = (*CFIndex)(nil)

// Ensure the CfIndex type implements the NeedsInputser interface.
var _ NeedsInputser = (*CFIndex)(nil)

// NeedsInputs signals that the index requires the referenced inputs in order to properly create the index. This implements the NeedsInputser interface.
func (idx *CFIndex) NeedsInputs() bool {
	return true
}

// Init initializes the hash-based cf index. This is part of the Indexer interface.
func (idx *CFIndex) Init() (err error) {
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

// Create is invoked when the indexer manager determines the index needs to be created for the first time. It creates buckets for the two hash-based cf indexes (regular only currently).
func (idx *CFIndex) Create(dbTx database.Tx) (err error) {
	meta := dbTx.Metadata()
	cfIndexParentBucket, err := meta.CreateBucket(cfIndexParentBucketKey)
	if err != nil {
		slog.Error(err)
		return err
	}
	for _, bucketName := range cfIndexKeys {
		_, err = cfIndexParentBucket.CreateBucket(bucketName)
		if err != nil {
			slog.Error(err)
			return err
		}
	}
	for _, bucketName := range cfHeaderKeys {
		_, err = cfIndexParentBucket.CreateBucket(bucketName)
		if err != nil {
			slog.Error(err)
			return err
		}
	}
	for _, bucketName := range cfHashKeys {
		_, err = cfIndexParentBucket.CreateBucket(bucketName)
		if err != nil {
			slog.Error(err)
			return err
		}
	}
	return nil
}

// storeFilter stores a given filter, and performs the steps needed to generate the filter's header.
func storeFilter(dbTx database.Tx, block *util.Block, f *gcs.Filter,
	filterType wire.FilterType) (err error) {
	if uint8(filterType) > maxFilterType {
		return errors.New("unsupported filter type")
	}
	// Figure out which buckets to use.
	fKey := cfIndexKeys[filterType]
	hKey := cfHeaderKeys[filterType]
	hashKey := cfHashKeys[filterType]
	// Start by storing the filter.
	h := block.Hash()
	var filterBytes []byte
	if filterBytes, err = f.NBytes(); slog.Check(err) {
		return
	}
	if err = dbStoreFilterIdxEntry(dbTx, fKey, h, filterBytes); slog.Check(err) {
		return
	}
	// Next store the filter hash.
	var filterHash chainhash.Hash
	if filterHash, err = builder.GetFilterHash(f); slog.Check(err) {
		return
	}
	if err = dbStoreFilterIdxEntry(dbTx, hashKey, h, filterHash[:]); slog.Check(err) {
		return
	}
	// Then fetch the previous block's filter header.
	var prevHeader *chainhash.Hash
	ph := &block.MsgBlock().Header.PrevBlock
	if ph.IsEqual(&zeroHash) {
		prevHeader = &zeroHash
	} else {
		var pfh []byte
		if pfh, err = dbFetchFilterIdxEntry(dbTx, hKey, ph); slog.Check(err) {
			return
		}
		// Construct the new block's filter header, and store it.
		if prevHeader, err = chainhash.NewHash(pfh); slog.Check(err) {
			return
		}
	}
	var fh chainhash.Hash
	if fh, err = builder.MakeHeaderForFilter(f, *prevHeader); slog.Check(err) {
		return
	}
	return dbStoreFilterIdxEntry(dbTx, hKey, h, fh[:])
}

// ConnectBlock is invoked by the index manager when a new block has been connected to the main chain. This indexer adds a hash-to-cf mapping for every passed block. This is part of the Indexer interface.
func (idx *CFIndex) ConnectBlock(dbTx database.Tx, block *util.Block,
	stxos []blockchain.SpentTxOut) (err error) {
	prevScripts := make([][]byte, len(stxos))
	for i, stxo := range stxos {
		prevScripts[i] = stxo.PkScript
	}
	var f *gcs.Filter
	if f, err = builder.BuildBasicFilter(block.MsgBlock(), prevScripts); slog.Check(err) {
		return
	}
	return storeFilter(dbTx, block, f, wire.GCSFilterRegular)
}

// DisconnectBlock is invoked by the index manager when a block has been disconnected from the main chain.  This indexer removes the hash-to-cf mapping for every passed block. This is part of the Indexer interface.
func (idx *CFIndex) DisconnectBlock(dbTx database.Tx, block *util.Block,
	_ []blockchain.SpentTxOut) (err error) {
	for _, key := range cfIndexKeys {
		if err = dbDeleteFilterIdxEntry(dbTx, key, block.Hash()); slog.Check(err) {
			return
		}
	}
	for _, key := range cfHeaderKeys {
		if err = dbDeleteFilterIdxEntry(dbTx, key, block.Hash()); slog.Check(err) {
			return
		}
	}
	for _, key := range cfHashKeys {
		if err = dbDeleteFilterIdxEntry(dbTx, key, block.Hash()); slog.Check(err) {
			return
		}
	}
	return
}

// entryByBlockHash fetches a filter index entry of a particular type (eg. filter, filter header, etc) for a filter type
// and block hash.
func (idx *CFIndex) entryByBlockHash(filterTypeKeys [][]byte,
	filterType wire.FilterType, h *chainhash.Hash) (entry []byte, err error) {
	if uint8(filterType) > maxFilterType {
		err = errors.New("unsupported filter type")
		slog.Debug(err)
		return
	}
	key := filterTypeKeys[filterType]
	err = idx.db.View(func(dbTx database.Tx) (err error) {
		entry, err = dbFetchFilterIdxEntry(dbTx, key, h)
		return
	})
	return
}

// entriesByBlockHashes batch fetches a filter index entry of a particular type (eg. filter, filter header, etc) for a
// filter type and slice of block hashes.
func (idx *CFIndex) entriesByBlockHashes(filterTypeKeys [][]byte,
	filterType wire.FilterType, blockHashes []*chainhash.Hash) (entries [][]byte, err error) {
	if uint8(filterType) > maxFilterType {
		err = errors.New("unsupported filter type")
		slog.Debug(err)
		return
	}
	key := filterTypeKeys[filterType]
	entries = make([][]byte, 0, len(blockHashes))
	var entry []byte
	if err = idx.db.View(func(dbTx database.Tx) (err error) {
		for _, blockHash := range blockHashes {
			if entry, err = dbFetchFilterIdxEntry(dbTx, key, blockHash); slog.Check(err) {
				return
			}
			entries = append(entries, entry)
		}
		return
	}); slog.Check(err) {
	}
	return
}

// FilterByBlockHash returns the serialized contents of a block's basic or committed filter.
func (idx *CFIndex) FilterByBlockHash(h *chainhash.Hash,
	filterType wire.FilterType) (entry []byte, err error) {
	return idx.entryByBlockHash(cfIndexKeys, filterType, h)
}

// FiltersByBlockHashes returns the serialized contents of a block's basic or committed filter for a set of blocks by hash.
func (idx *CFIndex) FiltersByBlockHashes(blockHashes []*chainhash.Hash,
	filterType wire.FilterType) (entries [][]byte, err error) {
	return idx.entriesByBlockHashes(cfIndexKeys, filterType, blockHashes)
}

// FilterHeaderByBlockHash returns the serialized contents of a block's basic committed filter header.
func (idx *CFIndex) FilterHeaderByBlockHash(h *chainhash.Hash,
	filterType wire.FilterType) (entry []byte, err error) {
	return idx.entryByBlockHash(cfHeaderKeys, filterType, h)
}

// FilterHeadersByBlockHashes returns the serialized contents of a block's basic committed filter header for a set of blocks by hash.
func (idx *CFIndex) FilterHeadersByBlockHashes(blockHashes []*chainhash.Hash,
	filterType wire.FilterType) (entry [][]byte, err error) {
	return idx.entriesByBlockHashes(cfHeaderKeys, filterType, blockHashes)
}

// FilterHashByBlockHash returns the serialized contents of a block's basic committed filter hash.
func (idx *CFIndex) FilterHashByBlockHash(h *chainhash.Hash,
	filterType wire.FilterType) (entry []byte, err error) {
	return idx.entryByBlockHash(cfHashKeys, filterType, h)
}

// FilterHashesByBlockHashes returns the serialized contents of a block's basic committed filter hash for a set of blocks by hash.
func (idx *CFIndex) FilterHashesByBlockHashes(blockHashes []*chainhash.Hash,
	filterType wire.FilterType) (entries [][]byte, err error) {
	return idx.entriesByBlockHashes(cfHashKeys, filterType, blockHashes)
}

// NewCfIndex returns a new instance of an indexer that is used to create a mapping of the hashes of all blocks in the blockchain to their respective committed filters.
// It implements the Indexer interface which plugs into the IndexManager that in turn is used by the blockchain package. This allows the index to be seamlessly maintained along with the chain.
func NewCfIndex(db database.DB, chainParams *netparams.Params) *CFIndex {
	return &CFIndex{db: db, chainParams: chainParams}
}

// DropCfIndex drops the CF index from the provided database if exists.
func DropCfIndex(db database.DB, interrupt <-chan struct{}) (err error) {
	return dropIndex(db, cfIndexParentBucketKey, cfIndexName, interrupt)
}
