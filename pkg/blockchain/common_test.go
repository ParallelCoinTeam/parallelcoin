package blockchain

import (
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/txscript"
	
	"github.com/p9c/matrjoska/pkg/database"
	_ "github.com/p9c/matrjoska/pkg/database/ffldb"
	"github.com/p9c/matrjoska/pkg/wire"
)

const (
	// testDbType is the database backend type to use for the tests.
	testDbType = "ffldb"
	// testDbRoot is the root directory used to create all test databases.
	testDbRoot = "testdbs"
	// blockDataNet is the expected network in the test block data.
	blockDataNet = wire.MainNet
)

// filesExists returns whether or not the named file or directory exists.
func fileExists(name string) bool {
	if _, e := os.Stat(name); E.Chk(e) {
		if os.IsNotExist(e) {
			return false
		}
	}
	return true
}

// isSupportedDbType returns whether or not the passed database type is currently supported.
func isSupportedDbType(dbType string) bool {
	supportedDrivers := database.SupportedDrivers()
	for _, driver := range supportedDrivers {
		if dbType == driver {
			return true
		}
	}
	return false
}

// loadBlocks reads files containing bitcoin block data (gzipped but otherwise in the format bitcoind writes) from disk
// and returns them as an array of util.Block. This is largely borrowed from the test code in pod.
func loadBlocks(filename string) (blocks []*block.Block, e error) {
	filename = filepath.Join("tstdata/", filename)
	var network = wire.MainNet
	var dr io.Reader
	var fi io.ReadCloser
	fi, e = os.Open(filename)
	if e != nil {
		return
	}
	if strings.HasSuffix(filename, ".bz2") {
		dr = bzip2.NewReader(fi)
	} else {
		dr = fi
	}
	defer func() {
		if e = fi.Close(); E.Chk(e) {
		}
	}()
	var blk *block.Block
	height := int64(1)
	for ; ; height++ {
		var rintbuf uint32
		e = binary.Read(dr, binary.LittleEndian, &rintbuf)
		if e == io.EOF {
			// hit end of file at expected offset: no warning
			// height--
			e = nil
		}
		if rintbuf != uint32(network) {
			break
		}
		e = binary.Read(dr, binary.LittleEndian, &rintbuf)
		blocklen := rintbuf
		rbytes := make([]byte, blocklen)
		// read block
		_, e = dr.Read(rbytes)
		if e != nil {
			fmt.Println(e)
		}
		blk, e = block.NewFromBytes(rbytes)
		if e != nil {
			return
		}
		blocks = append(blocks, blk)
	}
	return
}

// chainSetup is used to create a new db and chain instance with the genesis block already inserted. In addition to the
// new chain instance, it returns a teardown function the caller should invoke when done testing to clean up.
func chainSetup(dbName string, netparams *chaincfg.Params) (chain *BlockChain, teardown func(), e error) {
	if !isSupportedDbType(testDbType) {
		return nil, nil, fmt.Errorf("unsupported db type %v", testDbType)
	}
	// Handle memory database specially since it doesn't need the disk specific handling.
	var db database.DB
	if testDbType == "memdb" {
		var ndb database.DB
		ndb, e = database.Create(testDbType)
		if e != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", e)
		}
		db = ndb
		// Setup a teardown function for cleaning up. This function is returned to the caller to be invoked when it is
		// done testing.
		teardown = func() {
			if e = db.Close(); E.Chk(e) {
			}
		}
	} else {
		// Create the root directory for test databases.
		if !fileExists(testDbRoot) {
			if e = os.MkdirAll(testDbRoot, 0700); E.Chk(e) {
				e = fmt.Errorf(
					"unable to create test db "+
						"root: %v", e,
				)
				return nil, nil, e
			}
		}
		// Create a new database to store the accepted blocks into.
		dbPath := filepath.Join(testDbRoot, dbName)
		_ = os.RemoveAll(dbPath)
		var ndb database.DB
		ndb, e = database.Create(testDbType, dbPath, blockDataNet)
		if e != nil {
			return nil, nil, fmt.Errorf("error creating db: %v", e)
		}
		db = ndb
		// Setup a teardown function for cleaning up. This function is returned to the caller to be invoked when it is
		// done testing.
		teardown = func() {
			if e = db.Close(); E.Chk(e) {
			}
			if e = os.RemoveAll(dbPath); E.Chk(e) {
			}
			if e = os.RemoveAll(testDbRoot); E.Chk(e) {
			}
		}
	}
	// Copy the chain netparams to ensure any modifications the tests do to the chain parameters do not affect the
	// global instance.
	paramsCopy := *netparams
	// Create the main chain instance.
	chain, e = New(
		&Config{
			DB:          db,
			ChainParams: &paramsCopy,
			Checkpoints: nil,
			TimeSource:  NewMedianTime(),
			SigCache:    txscript.NewSigCache(1000),
		},
	)
	if e != nil {
		teardown()
		e = fmt.Errorf("failed to create chain instance: %v", e)
		return nil, nil, e
	}
	return chain, teardown, nil
}

// loadUtxoView returns a utxo view loaded from a file.
func loadUtxoView(filename string) (*UtxoViewpoint, error) {
	// The utxostore file format is:
	//
	// <tx hash><output index><serialized utxo len><serialized utxo>
	//
	// The output index and serialized utxo len are little endian uint32s and the serialized utxo uses the format
	// described in chainio.go.
	filename = filepath.Join("tstdata", filename)
	fi, e := os.Open(filename)
	if e != nil {
		return nil, e
	}
	// Choose read based on whether the file is compressed or not.
	var r io.Reader
	if strings.HasSuffix(filename, ".bz2") {
		r = bzip2.NewReader(fi)
	} else {
		r = fi
	}
	defer func() {
		if e := fi.Close(); E.Chk(e) {
		}
	}()
	view := NewUtxoViewpoint()
	for {
		// Hash of the utxo entry.
		var hash chainhash.Hash
		_, e := io.ReadAtLeast(r, hash[:], len(hash[:]))
		if e != nil {
			// Expected EOF at the right offset.
			if e == io.EOF {
				break
			}
			return nil, e
		}
		// Output index of the utxo entry.
		var index uint32
		e = binary.Read(r, binary.LittleEndian, &index)
		if e != nil {
			return nil, e
		}
		// Num of serialized utxo entry bytes.
		var numBytes uint32
		e = binary.Read(r, binary.LittleEndian, &numBytes)
		if e != nil {
			return nil, e
		}
		// Serialized utxo entry.
		serialized := make([]byte, numBytes)
		_, e = io.ReadAtLeast(r, serialized, int(numBytes))
		if e != nil {
			return nil, e
		}
		// Deserialize it and add it to the view.
		entry, e := deserializeUtxoEntry(serialized)
		if e != nil {
			return nil, e
		}
		view.Entries()[wire.OutPoint{Hash: hash, Index: index}] = entry
	}
	return view, nil
}

// convertUtxoStore reads a utxostore from the legacy format and writes it back out using the latest format. It is only
// useful for converting utxostore data used in the tests, which has already been done. However, the code is left
// available for future reference.
func convertUtxoStore(r io.Reader, w io.Writer) (e error) {
	// The old utxostore file format was:
	//
	// <tx hash><serialized utxo len><serialized utxo>
	//
	// The serialized utxo len was a little endian uint32 and the serialized utxo uses the format described in
	// upgrade.go.
	littleEndian := binary.LittleEndian
	for {
		// Hash of the utxo entry.
		var hash chainhash.Hash
		_, e := io.ReadAtLeast(r, hash[:], len(hash[:]))
		if e != nil {
			// Expected EOF at the right offset.
			if e == io.EOF {
				break
			}
			return e
		}
		// Num of serialized utxo entry bytes.
		var numBytes uint32
		e = binary.Read(r, littleEndian, &numBytes)
		if e != nil {
			return e
		}
		// Serialized utxo entry.
		serialized := make([]byte, numBytes)
		_, e = io.ReadAtLeast(r, serialized, int(numBytes))
		if e != nil {
			return e
		}
		// Deserialize the entry.
		entries, e := deserializeUtxoEntryV0(serialized)
		if e != nil {
			return e
		}
		// Loop through all of the utxos and write them out in the new format.
		for outputIdx, entry := range entries {
			// Reserialize the entries using the new format.
			serialized, e := serializeUtxoEntry(entry)
			if e != nil {
				return e
			}
			// Write the hash of the utxo entry.
			_, e = w.Write(hash[:])
			if e != nil {
				return e
			}
			// Write the output index of the utxo entry.
			e = binary.Write(w, littleEndian, outputIdx)
			if e != nil {
				return e
			}
			// Write num of serialized utxo entry bytes.
			e = binary.Write(w, littleEndian, uint32(len(serialized)))
			if e != nil {
				return e
			}
			// Write the serialized utxo.
			_, e = w.Write(serialized)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

// TstSetCoinbaseMaturity makes the ability to set the coinbase maturity available when running tests.
func (b *BlockChain) TstSetCoinbaseMaturity(maturity uint16) {
	b.params.CoinbaseMaturity = maturity
}

// newFakeChain returns a chain that is usable for syntetic tests. It is important to note that this chain has no
// database associated with it, so it is not usable with all functions and the tests must take care when making use of
// it.
func newFakeChain(params *chaincfg.Params) *BlockChain {
	// Create a genesis block node and block index index populated with it for use when creating the fake chain below.
	node := NewBlockNode(&params.GenesisBlock.Header, nil)
	index := newBlockIndex(nil, params)
	index.AddNode(node)
	targetTimespan := params.TargetTimespan
	targetTimePerBlock := params.TargetTimePerBlock
	adjustmentFactor := params.RetargetAdjustmentFactor
	return &BlockChain{
		params:              params,
		timeSource:          NewMedianTime(),
		minRetargetTimespan: targetTimespan / adjustmentFactor,
		maxRetargetTimespan: targetTimespan * adjustmentFactor,
		blocksPerRetarget:   int32(targetTimespan / targetTimePerBlock),
		Index:               index,
		BestChain:           newChainView(node),
	}
}

// newFakeNode creates a block node connected to the passed parent with the provided fields populated and fake values
// for the other fields.
func newFakeNode(parent *BlockNode, blockVersion int32, bits uint32, timestamp time.Time) *BlockNode {
	// Make up a header and create a block node from it.
	header := &wire.BlockHeader{
		Version:   blockVersion,
		PrevBlock: parent.hash,
		Bits:      bits,
		Timestamp: timestamp,
	}
	return NewBlockNode(header, parent)
}
