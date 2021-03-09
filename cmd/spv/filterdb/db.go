package filterdb

import (
	"fmt"
	
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	chainhash "github.com/p9c/pod/pkg/blockchain/chainhash"
	"github.com/p9c/pod/pkg/coding/gcs"
	"github.com/p9c/pod/pkg/coding/gcs/builder"
	"github.com/p9c/pod/pkg/database/walletdb"
)

var (
	// filterBucket is the name of the root bucket for this package. Within this bucket, sub-buckets are stored which
	// themselves store the actual filters.
	filterBucket = []byte("filter-store")
	// regBucket is the bucket that stores the regular filters.
	regBucket = []byte("regular")
)

// FilterType is a enum-like type that represents the various filter types currently defined.
type FilterType uint8

const (
	// RegularFilter is the filter type of regular filters which contain outputs and pkScript data pushes.
	RegularFilter FilterType = iota
)

var (
	// ErrFilterNotFound is returned when a filter for a target block hash is unable to be located.
	ErrFilterNotFound = fmt.Errorf("unable to find filter")
)

// FilterDatabase is an interface which represents an object that is capable of storing and retrieving filters according
// to their corresponding block hash and also their filter type.
// TODO(roasbeef): similar interface for headerfs?
type FilterDatabase interface {
	// PutFilter stores a filter with the given hash and type to persistent storage.
	PutFilter(*chainhash.Hash, *gcs.Filter, FilterType) error
	// FetchFilter attempts to fetch a filter with the given hash and type from persistent storage. In the case that a
	// filter matching the target block hash cannot be found, then ErrFilterNotFound is to be returned.
	FetchFilter(*chainhash.Hash, FilterType) (*gcs.Filter, error)
}

// FilterStore is an implementation of the FilterDatabase interface which is backed by boltdb.
type FilterStore struct {
	db walletdb.DB
	// chainParams netparams.Params
}

// A compile-time check to ensure the FilterStore adheres to the FilterDatabase interface.
var _ FilterDatabase = (*FilterStore)(nil)

// New creates a new instance of the FilterStore given an already open database, and the target chain parameters.
func New(db walletdb.DB, params netparams.Params) (*FilterStore, error) {
	e := walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			// As part of our initial setup, we'll try to create the top level filter bucket. If this already exists, then
			// we can exit early.
			filters, e := tx.CreateTopLevelBucket(filterBucket)
			if e != nil {
				return e
			}
			// If the main bucket doesn't already exist, then we'll need to create the sub-buckets, and also initialize them
			// with the genesis filters.
			genesisBlock := params.GenesisBlock
			genesisHash := params.GenesisHash
			// First we'll create the bucket for the regular filters.
			regFilters, e := filters.CreateBucketIfNotExists(regBucket)
			if e != nil {
				return e
			}
			// With the bucket created, we'll now construct the initial basic genesis filter and store it within the
			// database.
			basicFilter, e := builder.BuildBasicFilter(genesisBlock, nil)
			if e != nil {
				return e
			}
			return putFilter(regFilters, genesisHash, basicFilter)
		},
	)
	if e != nil && e != walletdb.ErrBucketExists {
		return nil, e
	}
	return &FilterStore{
			db: db,
		},
		nil
}

// putFilter stores a filter in the database according to the corresponding block hash. The passed bucket is expected to
// be the proper bucket for the passed filter type.
func putFilter(
	bucket walletdb.ReadWriteBucket, hash *chainhash.Hash,
	filter *gcs.Filter,
) (e error) {
	if filter == nil {
		return bucket.Put(hash[:], nil)
	}
	bytes, e := filter.NBytes()
	if e != nil {
		return e
	}
	return bucket.Put(hash[:], bytes)
}

// PutFilter stores a filter with the given hash and type to persistent storage.
//
// NOTE: This method is a part of the FilterDatabase interface.
func (f *FilterStore) PutFilter(
	hash *chainhash.Hash,
	filter *gcs.Filter, fType FilterType,
) (e error) {
	return walletdb.Update(
		f.db, func(tx walletdb.ReadWriteTx) (e error) {
			filters := tx.ReadWriteBucket(filterBucket)
			var targetBucket walletdb.ReadWriteBucket
			switch fType {
			case RegularFilter:
				targetBucket = filters.NestedReadWriteBucket(regBucket)
			default:
				return fmt.Errorf("unknown filter type: %v", fType)
			}
			if filter == nil {
				return targetBucket.Put(hash[:], nil)
			}
			bytes, e := filter.NBytes()
			if e != nil {
				return e
			}
			return targetBucket.Put(hash[:], bytes)
		},
	)
}

// FetchFilter attempts to fetch a filter with the given hash and type from persistent storage.
//
// NOTE: This method is a part of the FilterDatabase interface.
func (f *FilterStore) FetchFilter(
	blockHash *chainhash.Hash,
	filterType FilterType,
) (*gcs.Filter, error) {
	var filter *gcs.Filter
	e := walletdb.View(
		f.db, func(tx walletdb.ReadTx) (e error) {
			filters := tx.ReadBucket(filterBucket)
			var targetBucket walletdb.ReadBucket
			switch filterType {
			case RegularFilter:
				targetBucket = filters.NestedReadBucket(regBucket)
			default:
				return fmt.Errorf("unknown filter type")
			}
			filterBytes := targetBucket.Get(blockHash[:])
			if filterBytes == nil {
				return ErrFilterNotFound
			}
			if len(filterBytes) == 0 {
				return nil
			}
			dbFilter, e := gcs.FromNBytes(
				builder.DefaultP, builder.DefaultM, filterBytes,
			)
			if e != nil {
				return e
			}
			filter = dbFilter
			return nil
		},
	)
	if e != nil {
		return nil, e
	}
	return filter, nil
}
