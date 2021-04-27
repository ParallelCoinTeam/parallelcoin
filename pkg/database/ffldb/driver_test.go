package ffldb_test

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/block"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/database/ffldb"
)

// dbType is the database type name for this driver.
const dbType = "ffldb"

// TestCreateOpenFail ensures that errors related to creating and opening a database are handled properly.
func TestCreateOpenFail(t *testing.T) {
	t.Parallel()
	// Ensure that attempting to open a database that doesn't exist returns the expected error.
	wantErrCode := database.ErrDbDoesNotExist
	_, e := database.Open(dbType, "noexist", blockDataNet)
	if !checkDbError(t, "Open", e, wantErrCode) {
		return
	}
	// Ensure that attempting to open a database with the wrong number of parameters returns the expected error.
	wantErr := fmt.Errorf(
		"invalid arguments to %s.Open -- expected "+
			"database path and block network", dbType,
	)
	_, e = database.Open(dbType, 1, 2, 3)
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Open: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure that attempting to open a database with an invalid type for the first parameter returns the expected
	// error.
	wantErr = fmt.Errorf(
		"first argument to %s.Open is invalid -- "+
			"expected database path string", dbType,
	)
	_, e = database.Open(dbType, 1, blockDataNet)
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Open: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure that attempting to open a database with an invalid type for the second parameter returns the expected
	// error.
	wantErr = fmt.Errorf(
		"second argument to %s.Open is invalid -- "+
			"expected block network", dbType,
	)
	_, e = database.Open(dbType, "noexist", "invalid")
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Open: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure that attempting to create a database with the wrong number of parameters returns the expected error.
	wantErr = fmt.Errorf(
		"invalid arguments to %s.Create -- expected "+
			"database path and block network", dbType,
	)
	_, e = database.Create(dbType, 1, 2, 3)
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Create: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure that attempting to create a database with an invalid type for the first parameter returns the expected
	// error.
	wantErr = fmt.Errorf(
		"first argument to %s.Create is invalid -- "+
			"expected database path string", dbType,
	)
	_, e = database.Create(dbType, 1, blockDataNet)
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Create: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure that attempting to create a database with an invalid type for the second parameter returns the expected
	// error.
	wantErr = fmt.Errorf(
		"second argument to %s.Create is invalid -- "+
			"expected block network", dbType,
	)
	_, e = database.Create(dbType, "noexist", "invalid")
	if e != nil && e.Error() != wantErr.Error() {
		t.Errorf(
			"Create: did not receive expected error - got %v, "+
				"want %v", e, wantErr,
		)
		return
	}
	// Ensure operations against a closed database return the expected error.
	dbPath := filepath.Join(os.TempDir(), "ffldb-createfail")
	_ = os.RemoveAll(dbPath)
	db, e := database.Create(dbType, dbPath, blockDataNet)
	if e != nil {
		t.Errorf("Create: unexpected error: %v", e)
		return
	}
	defer func() {
		if e = os.RemoveAll(dbPath); ffldb.E.Chk(e) {
		}
	}()
	func() {
		if e = db.Close(); ffldb.E.Chk(e) {
		}
	}()
	wantErrCode = database.ErrDbNotOpen
	e = db.View(
		func(tx database.Tx) (e error) {
			return nil
		},
	)
	if !checkDbError(t, "View", e, wantErrCode) {
		return
	}
	wantErrCode = database.ErrDbNotOpen
	e = db.Update(
		func(tx database.Tx) (e error) {
			return nil
		},
	)
	if !checkDbError(t, "Update", e, wantErrCode) {
		return
	}
	wantErrCode = database.ErrDbNotOpen
	_, e = db.Begin(false)
	if !checkDbError(t, "Begin(false)", e, wantErrCode) {
		return
	}
	wantErrCode = database.ErrDbNotOpen
	_, e = db.Begin(true)
	if !checkDbError(t, "Begin(true)", e, wantErrCode) {
		return
	}
	wantErrCode = database.ErrDbNotOpen
	e = db.Close()
	if !checkDbError(t, "Close", e, wantErrCode) {
		return
	}
}

// TestPersistence ensures that values stored are still valid after closing and reopening the database.
func TestPersistence(t *testing.T) {
	t.Parallel()
	// Create a new database to run tests against.
	dbPath := filepath.Join(os.TempDir(), "ffldb-persistencetest")
	_ = os.RemoveAll(dbPath)
	db, e := database.Create(dbType, dbPath, blockDataNet)
	if e != nil {
		t.Errorf("Failed to create test database (%s) %v", dbType, e)
		return
	}
	defer func() {
		if e = os.RemoveAll(dbPath); ffldb.E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); ffldb.E.Chk(e) {
		}
	}()
	// Create a bucket, put some values into it, and store a block so they can be tested for existence on re-open.
	bucket1Key := []byte("bucket1")
	storeValues := map[string]string{
		"b1key1": "foo1",
		"b1key2": "foo2",
		"b1key3": "foo3",
	}
	genesisBlock := block.NewBlock(chaincfg.MainNetParams.GenesisBlock)
	genesisHash := chaincfg.MainNetParams.GenesisHash
	e = db.Update(
		func(tx database.Tx) (e error) {
			metadataBucket := tx.Metadata()
			if metadataBucket == nil {
				return fmt.Errorf("metadata: unexpected nil bucket")
			}
			bucket1, e := metadataBucket.CreateBucket(bucket1Key)
			if e != nil {
				return fmt.Errorf(
					"createBucket: unexpected error: %v",
					e,
				)
			}
			for k, v := range storeValues {
				e := bucket1.Put([]byte(k), []byte(v))
				if e != nil {
					return fmt.Errorf(
						"put: unexpected error: %v",
						e,
					)
				}
			}
			if e := tx.StoreBlock(genesisBlock); E.Chk(e) {
				return fmt.Errorf(
					"StoreBlock: unexpected error: %v",
					e,
				)
			}
			return nil
		},
	)
	if e != nil {
		t.Errorf("Update: unexpected error: %v", e)
		return
	}
	// Close and reopen the database to ensure the values persist.
	if e = db.Close(); ffldb.E.Chk(e) {
	}
	db, e = database.Open(dbType, dbPath, blockDataNet)
	if e != nil {
		t.Errorf("Failed to open test database (%s) %v", dbType, e)
		return
	}
	defer func() {
		if e = db.Close(); ffldb.E.Chk(e) {
		}
	}()
	// Ensure the values previously stored in the 3rd namespace still exist and are correct.
	e = db.View(
		func(tx database.Tx) (e error) {
			metadataBucket := tx.Metadata()
			if metadataBucket == nil {
				return fmt.Errorf("metadata: unexpected nil bucket")
			}
			bucket1 := metadataBucket.Bucket(bucket1Key)
			if bucket1 == nil {
				return fmt.Errorf("bucket1: unexpected nil bucket")
			}
			for k, v := range storeValues {
				gotVal := bucket1.Get([]byte(k))
				if !reflect.DeepEqual(gotVal, []byte(v)) {
					return fmt.Errorf(
						"get: key '%s' does not match expected value - got %s, want %s",
						k, gotVal, v,
					)
				}
			}
			genesisBlockBytes, _ := genesisBlock.Bytes()
			gotBytes, e := tx.FetchBlock(genesisHash)
			if e != nil {
				return fmt.Errorf(
					"fetchBlock: unexpected error: %v",
					e,
				)
			}
			if !reflect.DeepEqual(gotBytes, genesisBlockBytes) {
				return fmt.Errorf("fetchBlock: stored block mismatch")
			}
			return nil
		},
	)
	if e != nil {
		t.Errorf("view: unexpected error: %v", e)
		return
	}
}

// // TestInterface performs all interfaces tests for this database driver.
// func TestInterface(// 	t *testing.T) {
// 	t.Parallel()
// 	// Create a new database to run tests against.
// 	dbPath := filepath.Join(os.TempDir(), "ffldb-interfacetest")
// 	_ = os.RemoveAll(dbPath)
// 	db, e := database.Create(dbType, dbPath, blockDataNet)
// 	if e != nil  {
// 		t.Errorf("Failed to create test database (%s) %v", dbType, err)
// 		return
// 	}
// 	defer os.RemoveAll(dbPath)
// 	defer db.Close()
// 	// Ensure the driver type is the expected value.
// 	gotDbType := db.Type()
// 	if gotDbType != dbType {
// 		t.Errorf("Type: unepxected driver type - got %v, want %v",
// 			gotDbType, dbType)
// 		return
// 	}
// 	// Run all of the interface tests against the database.
// 	runtime.GOMAXPROCS(runtime.NumCPU())
// 	// Change the maximum file size to a small value to force multiple flat files with the test data set.
// 	ffldb.TstRunWithMaxBlockFileSize(db, 2048, func() {
// 		testInterface(t, db)
// 	})
// }
