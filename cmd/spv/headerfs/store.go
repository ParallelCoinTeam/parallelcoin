package headerfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/p9c/pod/pkg/coding/gcs"
	"github.com/stalker-loki/app/slog"
	"os"
	"path/filepath"
	"sync"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/coding/gcs/builder"
	"github.com/p9c/pod/pkg/db/walletdb"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
)

// BlockHeaderStore is an interface that provides an abstraction for a generic
// store for block headers.
type BlockHeaderStore interface {
	// ChainTip returns the best known block header and height for the
	// BlockHeaderStore.
	ChainTip() (*wire.BlockHeader, uint32, error)
	// LatestBlockLocator returns the latest block locator object based on
	// the tip of the current main chain from the PoV of the
	// BlockHeaderStore.
	LatestBlockLocator() (blockchain.BlockLocator, error)
	// FetchHeaderByHeight attempts to retrieve a target block header based
	// on a block height.
	FetchHeaderByHeight(height uint32) (*wire.BlockHeader, error)
	// FetchHeaderAncestors fetches the numHeaders block headers that are
	// the ancestors of the target stop hash. A total of numHeaders+1
	// headers will be returned, as we'll walk back numHeaders distance to
	// collect each header, then return the final header specified by the
	// stop hash. We'll also return the starting height of the header range
	// as well so callers can compute the height of each header without
	// knowing the height of the stop hash.
	FetchHeaderAncestors(uint32, *chainhash.Hash) ([]wire.BlockHeader,
		uint32, error)
	// HeightFromHash returns the height of a particular block header given
	// its hash.
	HeightFromHash(*chainhash.Hash) (uint32, error)
	// FetchHeader attempts to retrieve a block header determined by the
	// passed block height.
	FetchHeader(*chainhash.Hash) (*wire.BlockHeader, uint32, error)
	// WriteHeaders adds a set of headers to the BlockHeaderStore in a
	// single atomic transaction.
	WriteHeaders(...BlockHeader) error
	// RollbackLastBlock rolls back the BlockHeaderStore by a _single_
	// header. This method is meant to be used in the case of re-org which
	// disconnects the latest block header from the end of the main chain.
	// The information about the new header tip after truncation is
	// returned.
	RollbackLastBlock() (*waddrmgr.BlockStamp, error)
}

// headerBufPool is a pool of bytes.Buffer that will be re-used by the various
// headerStore implementations to batch their header writes to disk. By
// utilizing this variable we can minimize the total number of allocations when
// writing headers to disk.
var headerBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

// headerStore combines a on-disk set of headers within a flat file in addition
// to a database which indexes that flat file. Together, these two abstractions
// can be used in order to build an indexed header store for any type of
// "header" as it deals only with raw bytes, and leaves it to a higher layer to
// interpret those raw bytes accordingly.
//
// TODO(roasbeef): quickcheck coverage
type headerStore struct {
	// nolint
	mtx      sync.RWMutex
	filePath string
	file     *os.File
	*headerIndex
}

// newHeaderStore creates a new headerStore given an already open database, a
// target file path for the flat-file and a particular header type. The target
// file will be created as necessary.
func newHeaderStore(db walletdb.DB, filePath string,
	hType HeaderType) (hs *headerStore, err error) {
	var flatFileName string
	switch hType {
	case Block:
		flatFileName = "block_headers.bin"
	case RegularFilter:
		flatFileName = "reg_filter_headers.bin"
	default:
		return nil, fmt.Errorf("unrecognized filter type: %v", hType)
	}
	flatFileName = filepath.Join(filePath, flatFileName)
	// We'll open the file, creating it if necessary and ensuring that all
	// writes are actually appends to the end of the file.
	fileFlags := os.O_RDWR | os.O_APPEND | os.O_CREATE
	var headerFile *os.File
	if headerFile, err = os.OpenFile(flatFileName, fileFlags, 0644); slog.Check(err) {
		return
	}
	// With the file open, we'll then create the header index so we can
	// have random access into the flat files.
	var index *headerIndex
	if index, err = newHeaderIndex(db, hType); slog.Check(err) {
		return
	}
	hs = &headerStore{
		filePath:    filePath,
		file:        headerFile,
		headerIndex: index,
	}
	return
}

// blockHeaderStore is an implementation of the BlockHeaderStore interface, a
// fully fledged database for Bitcoin block headers. The blockHeaderStore
// combines a flat file to store the block headers with a database instance for
// managing the index into the set of flat files.
type blockHeaderStore struct {
	*headerStore
}

// A compile-time check to ensure the blockHeaderStore adheres to the
// BlockHeaderStore interface.
var _ BlockHeaderStore = (*blockHeaderStore)(nil)

// NewBlockHeaderStore creates a new instance of the blockHeaderStore based on
// a target file path, an open database instance, and finally a set of
// parameters for the target chain. These parameters are required as if this is
// the initial start up of the blockHeaderStore, then the initial genesis
// header will need to be inserted.
func NewBlockHeaderStore(filePath string, db walletdb.DB, netParams *netparams.Params) (b BlockHeaderStore, err error) {
	var hStore *headerStore
	if hStore, err = newHeaderStore(db, filePath, Block); slog.Check(err) {
		return
	}
	// With the header store created, we'll fetch the file size to see if
	// we need to initialize it with the first header or not.
	var fileInfo os.FileInfo
	if fileInfo, err = hStore.file.Stat(); slog.Check(err) {
		return
	}
	bhs := &blockHeaderStore{
		headerStore: hStore,
	}
	// If the size of the file is zero, then this means that we haven't yet
	// written the initial genesis header to disk, so we'll do so now.
	if fileInfo.Size() == 0 {
		genesisHeader := BlockHeader{
			BlockHeader: &netParams.GenesisBlock.Header,
			Height:      0,
		}
		if err = bhs.WriteHeaders(genesisHeader); slog.Check(err) {
			return
		}
	}
	// As a final initialization step (if this isn't the first time), we'll ensure that the header tip within the flat
	// files, is in sync with out database index.
	var tipHash *chainhash.Hash
	var tipHeight uint32
	if tipHash, tipHeight, err = hStore.chainTip(); slog.Check(err) {
		return
	}
	// First, we'll compute the size of the current file so we can
	// calculate the latest header written to disk.
	fileHeight := uint32(fileInfo.Size()/80) - 1
	// Using the file's current height, fetch the latest on-disk header.
	var latestFileHeader wire.BlockHeader
	if latestFileHeader, err = bhs.readHeader(fileHeight); slog.Check(err) {
		return
	}
	// If the index's tip hash, and the file on-disk match, then we're done here.
	var latestBlockHash chainhash.Hash
	latestBlockHash = latestFileHeader.BlockHash()
	if tipHash.IsEqual(&latestBlockHash) {
		return
	}
	// TODO(roasbeef): below assumes index can never get ahead?
	//  * we always update files _then_ indexes
	//  * need to dual pointer walk back for max safety
	// Otherwise, we'll need to truncate the file until it matches the
	// current index tip.
	for fileHeight > tipHeight {
		if err = bhs.singleTruncate(); slog.Check(err) {
			return
		}
		fileHeight--
	}
	b = bhs
	return
}

// FetchHeader attempts to retrieve a block header determined by the passed block height.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) FetchHeader(hash *chainhash.Hash) (header *wire.BlockHeader, height uint32, err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	// First, we'll query the index to obtain the block height of the passed block hash.
	if height, err = h.heightFromHash(hash); slog.Check(err) {
		return
	}
	// With the height known, we can now read the header from disk.
	var hdr wire.BlockHeader
	if hdr, err = h.readHeader(height); slog.Check(err) {
		return
	}
	header = &hdr
	return
}

// FetchHeaderByHeight attempts to retrieve a target block header based on a block height.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) FetchHeaderByHeight(height uint32) (header *wire.BlockHeader, err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	// For this query, we don't need to consult the index, and can instead
	// just seek into the flat file based on the target height and return
	// the full header.
	var hdr wire.BlockHeader
	if hdr, err = h.readHeader(height); slog.Check(err) {
		return
	}
	header = &hdr
	return
}

// FetchHeaderAncestors fetches the numHeaders block headers that are the
// ancestors of the target stop hash. A total of numHeaders+1 headers will be
// returned, as we'll walk back numHeaders distance to collect each header,
// then return the final header specified by the stop hash. We'll also return
// the starting height of the header range as well so callers can compute the
// height of each header without knowing the height of the stop hash.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) FetchHeaderAncestors(numHeaders uint32,
	stopHash *chainhash.Hash) (headers []wire.BlockHeader, startHeight uint32, err error) {
	// First, we'll find the final header in the range, this will be the
	// ending height of our scan.
	var endHeight uint32
	if endHeight, err = h.heightFromHash(stopHash); slog.Check(err) {
		return
	}
	startHeight = endHeight - numHeaders
	if headers, err = h.readHeaderRange(startHeight, endHeight); slog.Check(err) {
		return
	}
	return
}

// HeightFromHash returns the height of a particular block header given its
// hash.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) HeightFromHash(hash *chainhash.Hash) (u uint32, err error) {
	return h.heightFromHash(hash)
}

// RollbackLastBlock rolls back both the index, and on-disk header file by a
// _single_ header. This method is meant to be used in the case of re-org which
// disconnects the latest block header from the end of the main chain. The
// information about the new header tip after truncation is returned.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) RollbackLastBlock() (bs *waddrmgr.BlockStamp, err error) {
	// Lock store for write.
	h.mtx.Lock()
	defer h.mtx.Unlock()
	// First, we'll obtain the latest height that the index knows of.
	var chainTipHeight uint32
	if _, chainTipHeight, err = h.chainTip(); slog.Check(err) {
		return
	}
	// With this height obtained, we'll use it to read the latest header
	// from disk, so we can populate our return value which requires the
	// prev header hash.
	var bestHeader wire.BlockHeader
	if bestHeader, err = h.readHeader(chainTipHeight); slog.Check(err) {
		return
	}
	prevHeaderHash := bestHeader.PrevBlock
	// Now that we have the information we need to return from this
	// function, we can now truncate the header file, and then use the hash
	// of the prevHeader to set the proper index chain tip.
	if err = h.singleTruncate(); slog.Check(err) {
		return
	}
	if err = h.truncateIndex(&prevHeaderHash, true); slog.Check(err) {
		return
	}
	bs = &waddrmgr.BlockStamp{
		Height: int32(chainTipHeight) - 1,
		Hash:   prevHeaderHash,
	}
	return
}

// BlockHeader is a Bitcoin block header that also has its height included.
type BlockHeader struct {
	*wire.BlockHeader
	// Height is the height of this block header within the current main
	// chain.
	Height uint32
}

// toIndexEntry converts the BlockHeader into a matching headerEntry. This
// method is used when a header is to be written to disk.
func (b *BlockHeader) toIndexEntry() headerEntry {
	return headerEntry{
		hash:   b.BlockHash(),
		height: b.Height,
	}
}

// WriteHeaders writes a set of headers to disk and updates the index in a
// single atomic transaction.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) WriteHeaders(headers ...BlockHeader) (err error) {
	// Lock store for write.
	h.mtx.Lock()
	defer h.mtx.Unlock()
	// First, we'll grab a buffer from the write buffer pool so we an
	// reduce our total number of allocations, and also write the headers
	// in a single swoop.
	headerBuf := headerBufPool.Get().(*bytes.Buffer)
	headerBuf.Reset()
	defer headerBufPool.Put(headerBuf)
	// Next, we'll write out all the passed headers in series into the
	// buffer we just extracted from the pool.
	for _, header := range headers {
		if err := header.Serialize(headerBuf); err != nil {
			return err
		}
	}
	// With all the headers written to the buffer, we'll now write out the
	// entire batch in a single write call.
	if err := h.appendRaw(headerBuf.Bytes()); err != nil {
		return err
	}
	// Once those are written, we'll then collate all the headers into
	// headerEntry instances so we can write them all into the index in a
	// single atomic batch.
	headerLocations := make([]headerEntry, len(headers))
	for i, header := range headers {
		headerLocations[i] = header.toIndexEntry()
	}
	return h.addHeaders(headerLocations)
}

// blockLocatorFromHash takes a given block hash and then creates a block
// locator using it as the root of the locator. We'll start by taking a single
// step backwards, then keep doubling the distance until genesis after we get
// 10 locators.
//
// TODO(roasbeef): make into single transaction.
func (h *blockHeaderStore) blockLocatorFromHash(hash *chainhash.Hash) (locator blockchain.BlockLocator, err error) {
	// Append the initial hash
	locator = append(locator, hash)
	// If hash isn't found in DB or this is the genesis block, return the
	// locator as is
	var height uint32
	if height, err = h.heightFromHash(hash); height == 0 || slog.Check(err) {
		return
	}
	decrement := uint32(1)
	var blockHeader *wire.BlockHeader
	for height > 0 && len(locator) < wire.MaxBlockLocatorsPerMsg {
		// Decrement by 1 for the first 10 blocks, then double the jump
		// until we get to the genesis hash
		if len(locator) > 10 {
			decrement *= 2
		}
		if decrement > height {
			height = 0
		} else {
			height -= decrement
		}
		if blockHeader, err = h.FetchHeaderByHeight(height); slog.Check(err) {
			return
		}
		headerHash := blockHeader.BlockHash()
		locator = append(locator, &headerHash)
	}
	return
}

// LatestBlockLocator returns the latest block locator object based on the tip
// of the current main chain from the PoV of the database and flat files.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) LatestBlockLocator() (locator blockchain.BlockLocator, err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	var chainTipHash *chainhash.Hash
	if chainTipHash, _, err = h.chainTip(); slog.Check(err) {
		return
	}
	return h.blockLocatorFromHash(chainTipHash)
}

// BlockLocatorFromHash computes a block locator given a particular hash. The
// standard Bitcoin algorithm to compute block locators are employed.
func (h *blockHeaderStore) BlockLocatorFromHash(hash *chainhash.Hash) (locator blockchain.BlockLocator, err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	return h.blockLocatorFromHash(hash)
}

// CheckConnectivity cycles through all of the block headers on disk, from last
// to first, and makes sure they all connect to each other. Additionally, at
// each block header, we also ensure that the index entry for that height and
// hash also match up properly.
func (h *blockHeaderStore) CheckConnectivity() (err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	return walletdb.View(h.db, func(tx walletdb.ReadTx) (err error) {
		// First, we'll fetch the root bucket, in order to use that to
		// fetch the bucket that houses the header index.
		rootBucket := tx.ReadBucket(indexBucket)
		// With the header bucket retrieved, we'll now fetch the chain
		// tip so we can start our backwards scan.
		tipHash := rootBucket.Get(bitcoinTip)
		tipHeightBytes := rootBucket.Get(tipHash)
		// With the height extracted, we'll now read the _last_ block
		// header within the file before we kick off our connectivity
		// loop.
		tipHeight := binary.BigEndian.Uint32(tipHeightBytes)
		var header wire.BlockHeader
		if header, err = h.readHeader(tipHeight); slog.Check(err) {
			return
		}
		// We'll now cycle backwards, seeking backwards along the
		// header file to ensure each header connects properly and the
		// index entries are also accurate. To do this, we start from a
		// height of one before our current tip.
		var newHeader wire.BlockHeader
		for height := tipHeight - 1; height > 0; height-- {
			// First, read the block header for this block height,
			// and also compute the block hash for it.
			if newHeader, err = h.readHeader(height); slog.Check(err) {
				return fmt.Errorf("couldn't retrieve header %s: %s", header.PrevBlock, err)
			}
			newHeaderHash := newHeader.BlockHash()
			// With the header retrieved, we'll now fetch the
			// height for this current header hash to ensure the
			// on-disk state and the index matches up properly.
			indexHeightBytes := rootBucket.Get(newHeaderHash[:])
			if indexHeightBytes == nil {
				return fmt.Errorf("index and on-disk file out of sync at height: %v", height)
			}
			indexHeight := binary.BigEndian.Uint32(indexHeightBytes)
			// With the index entry retrieved, we'll now assert
			// that the height matches up with our current height
			// in this backwards walk.
			if indexHeight != height {
				return fmt.Errorf("index height isn't monotonically increasing")
			}
			// Finally, we'll assert that this new header is
			// actually the prev header of the target header from
			// the last loop. This ensures connectivity.
			if newHeader.BlockHash() != header.PrevBlock {
				return fmt.Errorf("block %s doesn't match block %s's PrevBlock (%s)",
					newHeader.BlockHash(), header.BlockHash(), header.PrevBlock)
			}
			// As all the checks have passed, we'll now reset our
			// header pointer to this current location, and
			// continue our backwards walk.
			header = newHeader
		}
		return
	})
}

// ChainTip returns the best known block header and height for the
// blockHeaderStore.
//
// NOTE: Part of the BlockHeaderStore interface.
func (h *blockHeaderStore) ChainTip() (latestHeader *wire.BlockHeader, tipHeight uint32, err error) {
	// Lock store for read.
	h.mtx.RLock()
	defer h.mtx.RUnlock()
	if _, tipHeight, err = h.chainTip(); slog.Check(err) {
		return
	}
	var hdr wire.BlockHeader
	if hdr, err = h.readHeader(tipHeight); slog.Check(err) {
		return
	}
	latestHeader = &hdr
	return
}

// FilterHeaderStore is an implementation of a fully fledged database for any
// variant of filter headers.  The FilterHeaderStore combines a flat file to
// store the block headers with a database instance for managing the index into
// the set of flat files.
type FilterHeaderStore struct {
	*headerStore
}

// NewFilterHeaderStore returns a new instance of the FilterHeaderStore based
// on a target file path, filter type, and target net parameters. These
// parameters are required as if this is the initial start up of the
// FilterHeaderStore, then the initial genesis filter header will need to be
// inserted.
func NewFilterHeaderStore(filePath string, db walletdb.DB,
	filterType HeaderType, netParams *netparams.Params) (fhs *FilterHeaderStore, err error) {
	var fStore *headerStore
	if fStore, err = newHeaderStore(db, filePath, filterType); slog.Check(err) {
		return
	}
	// With the header store created, we'll fetch the file size to see if
	// we need to initialize it with the first header or not.
	var fileInfo os.FileInfo
	if fileInfo, err = fStore.file.Stat(); slog.Check(err) {
		return
	}
	fhs = &FilterHeaderStore{
		fStore,
	}
	// TODO(roasbeef): also reconcile with block header state due to way
	//  roll back works atm
	// If the size of the file is zero, then this means that we haven't yet
	// written the initial genesis header to disk, so we'll do so now.
	if fileInfo.Size() == 0 {
		var genesisFilterHash chainhash.Hash
		switch filterType {
		case RegularFilter:
			var basicFilter *gcs.Filter
			if basicFilter, err = builder.BuildBasicFilter(
				netParams.GenesisBlock, nil,
			); slog.Check(err) {
				return
			}
			if genesisFilterHash, err = builder.MakeHeaderForFilter(
				basicFilter, netParams.GenesisBlock.Header.PrevBlock,
			); slog.Check(err) {
				return
			}
		default:
			err = fmt.Errorf("unknown filter type: %v", filterType)
			slog.Debug(err)
			return
		}
		genesisHeader := FilterHeader{
			HeaderHash: *netParams.GenesisHash,
			FilterHash: genesisFilterHash,
			Height:     0,
		}
		if err = fhs.WriteHeaders(genesisHeader); slog.Check(err) {
			return
		}
		return
	}
	// As a final initialization step, we'll ensure that the header tip
	// within the flat files, is in sync with out database index.
	var tipHash *chainhash.Hash
	var tipHeight uint32
	if tipHash, tipHeight, err = fhs.chainTip(); slog.Check(err) {
		return
	}
	// First, we'll compute the size of the current file so we can
	// calculate the latest header written to disk.
	fileHeight := uint32(fileInfo.Size()/32) - 1
	// Using the file's current height, fetch the latest on-disk header.
	var latestFileHeader *chainhash.Hash
	if latestFileHeader, err = fhs.readHeader(fileHeight); slog.Check(err) {
		return
	}
	// If the index's tip hash, and the file on-disk match, then we're
	// doing here.
	if tipHash.IsEqual(latestFileHeader) {
		return
	}
	// Otherwise, we'll need to truncate the file until it matches the
	// current index tip.
	for fileHeight > tipHeight {
		if err = fhs.singleTruncate(); slog.Check(err) {
			return
		}
		fileHeight--
	}
	// TODO(roasbeef): make above into func
	return fhs, nil
}

// FetchHeader returns the filter header that corresponds to the passed block
// height.
func (f *FilterHeaderStore) FetchHeader(hash *chainhash.Hash) (h *chainhash.Hash, err error) {
	// Lock store for read.
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	var height uint32
	if height, err = f.heightFromHash(hash); slog.Check(err) {
		return
	}
	return f.readHeader(height)
}

// FetchHeaderByHeight returns the filter header for a particular block height.
func (f *FilterHeaderStore) FetchHeaderByHeight(height uint32) (h *chainhash.Hash, err error) {
	// Lock store for read.
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return f.readHeader(height)
}

// FilterHeader represents a filter header (basic or extended). The filter
// header itself is coupled with the block height and hash of the filter's
// block.
type FilterHeader struct {
	// HeaderHash is the hash of the block header that this filter header
	// corresponds to.
	HeaderHash chainhash.Hash
	// FilterHash is the filter header itself.
	FilterHash chainhash.Hash
	// Height is the block height of the filter header in the main chain.
	Height uint32
}

// toIndexEntry converts the filter header into a index entry to be stored
// within the database.
func (f *FilterHeader) toIndexEntry() headerEntry {
	return headerEntry{
		hash:   f.HeaderHash,
		height: f.Height,
	}
}

// WriteHeaders writes a batch of filter headers to persistent storage. The
// headers themselves are appended to the flat file, and then the index updated
// to reflect the new entries.
func (f *FilterHeaderStore) WriteHeaders(headers ...FilterHeader) (err error) {
	// Lock store for write.
	f.mtx.Lock()
	defer f.mtx.Unlock()
	// If there are 0 headers to be written, return immediately. This
	// prevents the newTip assignment from panicking because of an index
	// of -1.
	if len(headers) == 0 {
		return nil
	}
	// First, we'll grab a buffer from the write buffer pool so we an
	// reduce our total number of allocations, and also write the headers
	// in a single swoop.
	headerBuf := headerBufPool.Get().(*bytes.Buffer)
	headerBuf.Reset()
	defer headerBufPool.Put(headerBuf)
	// Next, we'll write out all the passed headers in series into the
	// buffer we just extracted from the pool.
	for _, header := range headers {
		if _, err = headerBuf.Write(header.FilterHash[:]); slog.Check(err) {
			return
		}
	}
	// With all the headers written to the buffer, we'll now write out the
	// entire batch in a single write call.
	if err = f.appendRaw(headerBuf.Bytes()); slog.Check(err) {
		return
	}
	// As the block headers should already be written, we only need to
	// update the tip pointer for this particular header type.
	newTip := headers[len(headers)-1].toIndexEntry().hash
	return f.truncateIndex(&newTip, false)
}

// ChainTip returns the latest filter header and height known to the
// FilterHeaderStore.
func (f *FilterHeaderStore) ChainTip() (latestHeader *chainhash.Hash, tipHeight uint32, err error) {
	// Lock store for read.
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	if _, tipHeight, err = f.chainTip(); slog.Check(err) {
		return
	}
	if latestHeader, err = f.readHeader(tipHeight); slog.Check(err) {
		err = fmt.Errorf("unable to read header: %v", err)
		slog.Debug(err)
		return
	}
	return
}

// RollbackLastBlock rolls back both the index, and on-disk header file by a
// _single_ filter header. This method is meant to be used in the case of
// re-org which disconnects the latest filter header from the end of the main
// chain. The information about the latest header tip after truncation is
// returned.
func (f *FilterHeaderStore) RollbackLastBlock(newTip *chainhash.Hash) (bs *waddrmgr.BlockStamp, err error) {
	// Lock store for write.
	f.mtx.Lock()
	defer f.mtx.Unlock()
	// First, we'll obtain the latest height that the index knows of.
	var chainTipHeight uint32
	if _, chainTipHeight, err = f.chainTip(); slog.Check(err) {
		return
	}
	// With this height obtained, we'll use it to read what will be the new
	// chain tip from disk.
	newHeightTip := chainTipHeight - 1
	var newHeaderTip *chainhash.Hash
	if newHeaderTip, err = f.readHeader(newHeightTip); slog.Check(err) {
		return
	}
	// Now that we have the information we need to return from this
	// function, we can now truncate both the header file and the index.
	if err = f.singleTruncate(); slog.Check(err) {
		return
	}
	if err = f.truncateIndex(newTip, false); slog.Check(err) {
		return
	}
	// TODO(roasbeef): return chain hash also?
	bs = &waddrmgr.BlockStamp{
		Height: int32(newHeightTip),
		Hash:   *newHeaderTip,
	}
	return
}
