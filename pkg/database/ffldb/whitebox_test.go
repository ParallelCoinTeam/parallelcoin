package ffldb

// This file is part of the ffldb package rather than the ffldb_test package as it provides whitebox testing.
import (
	"compress/bzip2"
	"encoding/binary"
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"testing"
	
	"github.com/btcsuite/goleveldb/leveldb"
	ldberrors "github.com/btcsuite/goleveldb/leveldb/errors"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/wire"
)

var (
	// blockDataNet is the expected network in the test block data.
	blockDataNet = wire.MainNet
	// blockDataFile is the path to a file containing the first 256 blocks of the block chain.
	blockDataFile = filepath.Join("..", "tstdata", "blocks1-256.bz2")
	// errSubTestFail is used to signal that a sub test returned false.
	errSubTestFail = fmt.Errorf("sub test failure")
)

// loadBlocks loads the blocks contained in the tstdata directory and returns a slice of them.
func loadBlocks(t *testing.T, dataFile string, network wire.BitcoinNet) ([]*block.Block, error) {
	// Open the file that contains the blocks for reading.
	fi, e := os.Open(dataFile)
	if e != nil {
		t.Errorf("failed to open file %v, err %v", dataFile, e)
		return nil, e
	}
	defer func() {
		if e := fi.Close(); E.Chk(e) {
			t.Errorf(
				"failed to close file %v %v", dataFile,
				e,
			)
		}
	}()
	dr := bzip2.NewReader(fi)
	// Set the first block as the genesis block.
	blocks := make([]*block.Block, 0, 256)
	genesis := block.NewBlock(chaincfg.MainNetParams.GenesisBlock)
	blocks = append(blocks, genesis)
	// Load the remaining blocks.
	for height := 1; ; height++ {
		var net uint32
		e := binary.Read(dr, binary.LittleEndian, &net)
		if e == io.EOF {
			// Hit end of file at the expected offset.  No error.
			break
		}
		if e != nil {
			t.Errorf(
				"Failed to load network type for block %d: %v",
				height, e,
			)
			return nil, e
		}
		if net != uint32(network) {
			t.Errorf(
				"Block doesn't match network: %v expects %v",
				net, network,
			)
			return nil, e
		}
		var blockLen uint32
		e = binary.Read(dr, binary.LittleEndian, &blockLen)
		if e != nil {
			t.Errorf(
				"Failed to load block size for block %d: %v",
				height, e,
			)
			return nil, e
		}
		// Read the block.
		blockBytes := make([]byte, blockLen)
		_, e = io.ReadFull(dr, blockBytes)
		if e != nil {
			t.Errorf("Failed to load block %d: %v", height, e)
			return nil, e
		}
		// Deserialize and store the block.
		block, e := block.NewFromBytes(blockBytes)
		if e != nil {
			t.Errorf("Failed to parse block %v: %v", height, e)
			return nil, e
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// checkDbError ensures the passed error is a database.DBError with an error code that matches the passed  error code.
func checkDbError(t *testing.T, testName string, gotErr error, wantErrCode database.ErrorCode) bool {
	dbErr, ok := gotErr.(database.DBError)
	if !ok {
		t.Errorf(
			"%s: unexpected error type - got %T, want %T",
			testName, gotErr, database.DBError{},
		)
		return false
	}
	if dbErr.ErrorCode != wantErrCode {
		t.Errorf(
			"%s: unexpected error code - got %s (%s), want %s",
			testName, dbErr.ErrorCode, dbErr.Description,
			wantErrCode,
		)
		return false
	}
	return true
}

// testContext is used to store context information about a running test which is passed into helper functions.
type testContext struct {
	t            *testing.T
	db           database.DB
	files        map[uint32]*lockableFile
	maxFileSizes map[uint32]int64
	blocks       []*block.Block
}

// TestConvertErr ensures the leveldb error to database error conversion works as expected.
func TestConvertErr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err         error
		wantErrCode database.ErrorCode
	}{
		{&ldberrors.ErrCorrupted{}, database.ErrCorruption},
		{leveldb.ErrClosed, database.ErrDbNotOpen},
		{leveldb.ErrSnapshotReleased, database.ErrTxClosed},
		{leveldb.ErrIterReleased, database.ErrTxClosed},
	}
	for i, test := range tests {
		gotErr := convertErr("test", test.err)
		if gotErr.ErrorCode != test.wantErrCode {
			t.Errorf("convertErr #%d unexpected error - got %v, want %v", i, gotErr.ErrorCode, test.wantErrCode)
			continue
		}
	}
}

// TestCornerCases ensures several corner cases which can happen when opening a database and/or block files work as expected.
func TestCornerCases(t *testing.T) {
	t.Parallel()
	// Create a file at the datapase path to force the open below to fail.
	dbPath := filepath.Join(os.TempDir(), "ffldb-errors")
	_ = os.RemoveAll(dbPath)
	fi, e := os.Create(dbPath)
	if e != nil {
		t.Errorf("os.Create: unexpected error: %v", e)
		return
	}
	if e = fi.Close(); E.Chk(e) {
	}
	// Ensure creating a new database fails when a file exists where a directory is needed.
	testName := "openDB: fail due to file at target location"
	wantErrCode := database.ErrDriverSpecific
	var idb database.DB
	if idb, e = openDB(dbPath, blockDataNet, true); E.Chk(e) {
	}
	if !checkDbError(t, testName, e, wantErrCode) {
		if e = idb.Close(); E.Chk(e) {
		}
		if e = os.RemoveAll(dbPath); E.Chk(e) {
		}
		return
	}
	// Remove the file and create the database to run tests against.  It should be successful this time.
	_ = os.RemoveAll(dbPath)
	idb, e = openDB(dbPath, blockDataNet, true)
	if e != nil {
		t.Errorf("openDB: unexpected error: %v", e)
		return
	}
	defer func() {
		if e = os.RemoveAll(dbPath); E.Chk(e) {
		}
		if e = idb.Close(); E.Chk(e) {
		}
	}()
	// Ensure attempting to write to a file that can't be created returns the expected error.
	testName = "writeBlock: open file failure"
	filePath := blockFilePath(dbPath, 0)
	if e = os.Mkdir(filePath, 0755); E.Chk(e) {
		t.Errorf("os.Mkdir: unexpected error: %v", e)
		return
	}
	store := idb.(*db).store
	_, e = store.writeBlock([]byte{0x00})
	if !checkDbError(t, testName, e, database.ErrDriverSpecific) {
		return
	}
	_ = os.RemoveAll(filePath)
	// Close the underlying leveldb database out from under the database.
	ldb := idb.(*db).cache.ldb
	if e = ldb.Close(); E.Chk(e) {
	}
	// Ensure initilization errors in the underlying database work as expected.
	testName = "initDB: reinitialization"
	wantErrCode = database.ErrDbNotOpen
	e = initDB(ldb)
	if !checkDbError(t, testName, e, wantErrCode) {
		return
	}
	// Ensure the View handles errors in the underlying leveldb database properly.
	testName = "View: underlying leveldb error"
	wantErrCode = database.ErrDbNotOpen
	e = idb.View(
		func(tx database.Tx) (e error) {
			return nil
		},
	)
	if !checkDbError(t, testName, e, wantErrCode) {
		return
	}
	// Ensure the Update handles errors in the underlying leveldb database properly.
	testName = "Update: underlying leveldb error"
	e = idb.Update(
		func(tx database.Tx) (e error) {
			return nil
		},
	)
	if !checkDbError(t, testName, e, wantErrCode) {
		return
	}
}

// resetDatabase removes everything from the opened database associated with the test context including all metadata and the mock files.
func resetDatabase(tc *testContext) bool {
	// Reset the metadata.
	e := tc.db.Update(
		func(tx database.Tx) (e error) {
			// Remove all the keys using a cursor while also generating a list of buckets.  It's not safe to remove keys during ForEach iteration nor is it safe to remove buckets during cursor iteration, so this dual approach is needed.
			var bucketNames [][]byte
			cursor := tx.Metadata().Cursor()
			for ok := cursor.First(); ok; ok = cursor.Next() {
				if cursor.Value() != nil {
					if e = cursor.Delete(); E.Chk(e) {
						return e
					}
				} else {
					bucketNames = append(bucketNames, cursor.Key())
				}
			}
			// Remove the buckets.
			for _, k := range bucketNames {
				if e = tx.Metadata().DeleteBucket(k); E.Chk(e) {
					return e
				}
			}
			_, e = tx.Metadata().CreateBucket(blockIdxBucketName)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("Update: unexpected error: %v", e)
		return false
	}
	// Reset the mock files.
	store := tc.db.(*db).store
	wc := store.writeCursor
	wc.curFile.Lock()
	if wc.curFile.file != nil {
		if e := wc.curFile.file.Close(); E.Chk(e) {
		}
		wc.curFile.file = nil
	}
	wc.curFile.Unlock()
	wc.Lock()
	wc.curFileNum = 0
	wc.curOffset = 0
	wc.Unlock()
	tc.files = make(map[uint32]*lockableFile)
	tc.maxFileSizes = make(map[uint32]int64)
	return true
}

// testWriteFailures tests various failures paths when writing to the block files.
func testWriteFailures(tc *testContext) bool {
	if !resetDatabase(tc) {
		return false
	}
	// Ensure file sync errors during flush return the expected error.
	store := tc.db.(*db).store
	testName := "flush: file sync failure"
	store.writeCursor.Lock()
	oldFile := store.writeCursor.curFile
	store.writeCursor.curFile = &lockableFile{
		file: &mockFile{forceSyncErr: true, maxSize: -1},
	}
	store.writeCursor.Unlock()
	var e error
	e = tc.db.(*db).cache.flush()
	if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
		return false
	}
	store.writeCursor.Lock()
	store.writeCursor.curFile = oldFile
	store.writeCursor.Unlock()
	// Force errors in the various error paths when writing data by using mock files with a limited max size.
	block0Bytes, _ := tc.blocks[0].Bytes()
	tests := []struct {
		fileNum uint32
		maxSize int64
	}{
		// Force an error when writing the network bytes.
		{fileNum: 0, maxSize: 2},
		// Force an error when writing the block size.
		{fileNum: 0, maxSize: 6},
		// Force an error when writing the block.
		{fileNum: 0, maxSize: 17},
		// Force an error when writing the checksum.
		{fileNum: 0, maxSize: int64(len(block0Bytes)) + 10},
		// Force an error after writing enough blocks for force multiple
		// files.
		{fileNum: 15, maxSize: 1},
	}
	for i, test := range tests {
		if !resetDatabase(tc) {
			return false
		}
		// Ensure storing the specified number of blocks using a mock file that fails the write fails when the transaction is committed, not when the block is stored.
		tc.maxFileSizes = map[uint32]int64{test.fileNum: test.maxSize}
		e = tc.db.Update(
			func(tx database.Tx) (e error) {
				for i, block := range tc.blocks {
					e := tx.StoreBlock(block)
					if e != nil {
						tc.t.Errorf(
							"StoreBlock (%d): unexpected "+
								"error: %v", i, e,
						)
						return errSubTestFail
					}
				}
				return nil
			},
		)
		testName := fmt.Sprintf(
			"Force update commit failure - test "+
				"%d, fileNum %d, maxsize %d", i, test.fileNum,
			test.maxSize,
		)
		if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
			tc.t.Errorf("%v", e)
			return false
		}
		// Ensure the commit rollback removed all extra files and data.
		if len(tc.files) != 1 {
			tc.t.Errorf(
				"Update rollback: new not removed - want "+
					"1 file, got %d", len(tc.files),
			)
			return false
		}
		if _, ok := tc.files[0]; !ok {
			tc.t.Error("Update rollback: file 0 does not exist")
			return false
		}
		file := tc.files[0].file.(*mockFile)
		if len(file.data) != 0 {
			tc.t.Errorf(
				"Update rollback: file did not truncate - "+
					"want len 0, got len %d", len(file.data),
			)
			return false
		}
	}
	return true
}

// testBlockFileErrors ensures the database returns expected errors with various file-related issues such as closed and missing files.
func testBlockFileErrors(tc *testContext) bool {
	if !resetDatabase(tc) {
		return false
	}
	// Ensure errors in blockFile and openFile when requesting invalid file numbers.
	store := tc.db.(*db).store
	testName := "blockFile invalid file open"
	_, e := store.blockFile(^uint32(0))
	if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
		return false
	}
	testName = "openFile invalid file open"
	_, e = store.openFile(^uint32(0))
	if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
		return false
	}
	// Insert the first block into the mock file.
	e = tc.db.Update(
		func(tx database.Tx) (e error) {
			e = tx.StoreBlock(tc.blocks[0])
			if e != nil {
				tc.t.Errorf("StoreBlock: unexpected error: %v", e)
				return errSubTestFail
			}
			return nil
		},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("Update: unexpected error: %v", e)
		}
		return false
	}
	// Ensure errors in readBlock and readBlockRegion when requesting a file number that doesn't exist.
	block0Hash := tc.blocks[0].Hash()
	testName = "readBlock invalid file number"
	invalidLoc := blockLocation{
		blockFileNum: ^uint32(0),
		blockLen:     80,
	}
	_, e = store.readBlock(block0Hash, invalidLoc)
	if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
		return false
	}
	testName = "readBlockRegion invalid file number"
	_, e = store.readBlockRegion(invalidLoc, 0, 80)
	if !checkDbError(tc.t, testName, e, database.ErrDriverSpecific) {
		return false
	}
	// Close the block file out from under the database.
	store.writeCursor.curFile.Lock()
	if e = store.writeCursor.curFile.file.Close(); E.Chk(e) {
	}
	store.writeCursor.curFile.Unlock()
	// Ensure failures in FetchBlock and FetchBlockRegion(s) since the underlying file they need to read from has been closed.
	e = tc.db.View(
		func(tx database.Tx) (e error) {
			testName = "FetchBlock closed file"
			wantErrCode := database.ErrDriverSpecific
			_, e = tx.FetchBlock(block0Hash)
			if !checkDbError(tc.t, testName, e, wantErrCode) {
				return errSubTestFail
			}
			testName = "FetchBlockRegion closed file"
			regions := []database.BlockRegion{
				{
					Hash:   block0Hash,
					Len:    80,
					Offset: 0,
				},
			}
			_, e = tx.FetchBlockRegion(&regions[0])
			if !checkDbError(tc.t, testName, e, wantErrCode) {
				return errSubTestFail
			}
			testName = "FetchBlockRegions closed file"
			_, e = tx.FetchBlockRegions(regions)
			if !checkDbError(tc.t, testName, e, wantErrCode) {
				return errSubTestFail
			}
			return nil
		},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("View: unexpected error: %v", e)
		}
		return false
	}
	return true
}

// testCorruption ensures the database returns expected errors under various corruption scenarios.
func testCorruption(tc *testContext) bool {
	if !resetDatabase(tc) {
		return false
	}
	// Insert the first block into the mock file.
	e := tc.db.Update(
		func(tx database.Tx) (e error) {
			e = tx.StoreBlock(tc.blocks[0])
			if e != nil {
				tc.t.Errorf("StoreBlock: unexpected error: %v", e)
				return errSubTestFail
			}
			return nil
		},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("Update: unexpected error: %v", e)
		}
		return false
	}
	// Ensure corruption is detected by intentionally modifying the bytes stored to the mock file and reading the block.
	block0Bytes, _ := tc.blocks[0].Bytes()
	block0Hash := tc.blocks[0].Hash()
	tests := []struct {
		offset      uint32
		fixChecksum bool
		wantErrCode database.ErrorCode
	}{
		// One of the network bytes.  The checksum needs to be fixed so the invalid network is detected.
		{2, true, database.ErrDriverSpecific},
		// The same network byte, but this time don't fix the checksum to ensure the corruption is detected.
		{2, false, database.ErrCorruption},
		// One of the block length bytes.
		{6, false, database.ErrCorruption},
		// Random header byte.
		{17, false, database.ErrCorruption},
		// Random transaction byte.
		{90, false, database.ErrCorruption},
		// Random checksum byte.
		{uint32(len(block0Bytes)) + 10, false, database.ErrCorruption},
	}
	e = tc.db.View(
		func(tx database.Tx) (e error) {
			data := tc.files[0].file.(*mockFile).data
			for i, test := range tests {
				// Corrupt the byte at the offset by a single bit.
				data[test.offset] ^= 0x10
				// Fix the checksum if requested to force other errors.
				fileLen := len(data)
				var oldChecksumBytes [4]byte
				copy(oldChecksumBytes[:], data[fileLen-4:])
				if test.fixChecksum {
					toSum := data[:fileLen-4]
					cksum := crc32.Checksum(toSum, castagnoli)
					binary.BigEndian.PutUint32(data[fileLen-4:], cksum)
				}
				testName := fmt.Sprintf(
					"FetchBlock (test #%d): "+
						"corruption", i,
				)
				_, e = tx.FetchBlock(block0Hash)
				if !checkDbError(tc.t, testName, e, test.wantErrCode) {
					return errSubTestFail
				}
				// Reset the corrupted data back to the original.
				data[test.offset] ^= 0x10
				if test.fixChecksum {
					copy(data[fileLen-4:], oldChecksumBytes[:])
				}
			}
			return nil
		},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("View: unexpected error: %v", e)
		}
		return false
	}
	return true
}

// // TestFailureScenarios ensures several failure scenarios such as database corruption, block file write failures, and rollback failures are handled correctly.
// func TestFailureScenarios(// 	t *testing.T) {
// 	// Create a new database to run tests against.
// 	dbPath := filepath.Join(os.TempDir(), "ffldb-failurescenarios")
// 	_ = os.RemoveAll(dbPath)
// 	idb, e := database.Create(dbType, dbPath, blockDataNet)
// 	if e != nil  {
// 		t.Errorf("Failed to create test database (%s) %v", dbType, err)
// 		return
// 	}
// 	defer os.RemoveAll(dbPath)
// 	defer idb.Close()
// 	// Create a test context to pass around.
// 	tc := &testContext{
// 		t:            t,
// 		db:           idb,
// 		files:        make(map[uint32]*lockableFile),
// 		maxFileSizes: make(map[uint32]int64),
// 	}
// 	// Change the maximum file size to a small value to force multiple flat files with the test data set and replace the file-related functions to make use of mock files in memory.  This allows injection of various file-related errors.
// 	store := idb.(*db).store
// 	store.maxBlockFileSize = 1024 // 1KiB
// 	store.openWriteFileFunc = func(fileNum uint32) (filer, error) {
// 		if file, ok := tc.files[fileNum]; ok {
// 			// "Reopen" the file.
// 			file.Lock()
// 			mock := file.file.(*mockFile)
// 			mock.Lock()
// 			mock.closed = false
// 			mock.Unlock()
// 			file.Unlock()
// 			return mock, nil
// 		}
// 		// Limit the max size of the mock file as specified in the test context.
// 		maxSize := int64(-1)
// 		if maxFileSize, ok := tc.maxFileSizes[fileNum]; ok {
// 			maxSize = int64(maxFileSize)
// 		}
// 		file := &mockFile{maxSize: int64(maxSize)}
// 		tc.files[fileNum] = &lockableFile{file: file}
// 		return file, nil
// 	}
// 	store.openFileFunc = func(fileNum uint32) (*lockableFile, error) {
// 		// Force error when trying to open max file num.
// 		if fileNum == ^uint32(0) {
// 			return nil, makeDbErr(database.ErrDriverSpecific,
// 				"test", nil)
// 		}
// 		if file, ok := tc.files[fileNum]; ok {
// 			// "Reopen" the file.
// 			file.Lock()
// 			mock := file.file.(*mockFile)
// 			mock.Lock()
// 			mock.closed = false
// 			mock.Unlock()
// 			file.Unlock()
// 			return file, nil
// 		}
// 		file := &lockableFile{file: &mockFile{}}
// 		tc.files[fileNum] = file
// 		return file, nil
// 	}
// 	store.deleteFileFunc = func(fileNum uint32) (e error) {
// 		if file, ok := tc.files[fileNum]; ok {
// 			file.Lock()
// 			file.file.Close()
// 			file.Unlock()
// 			delete(tc.files, fileNum)
// 			return nil
// 		}
// 		str := fmt.Sprintf("file %d does not exist", fileNum)
// 		return makeDbErr(database.ErrDriverSpecific, str, nil)
// 	}
// 	// Load the test blocks and save in the test context for use throughout the tests.
// 	blocks, e := loadBlocks(t, blockDataFile, blockDataNet)
// 	if e != nil  {
// 		t.Errorf("loadBlocks: Unexpected error: %v", err)
// 		return
// 	}
// 	tc.blocks = blocks
// 	// Test various failures paths when writing to the block files.
// 	if !testWriteFailures(tc) {
// 		return
// 	}
// 	// Test various file-related issues such as closed and missing files.
// 	if !testBlockFileErrors(tc) {
// 		return
// 	}
// 	// Test various corruption scenarios.
// 	testCorruption(tc)
// }
