package database_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/p9c/matrjoska/pkg/database"
	_ "github.com/p9c/matrjoska/pkg/database/ffldb"
	"github.com/p9c/matrjoska/pkg/wire"
)

// This example demonstrates creating a new database.
func ExampleCreate() {
	// This example assumes the ffldb driver is imported.
	//
	// import (
	// 	"github.com/p9c/matrjoska/pkg/db"
	// 	_ "github.com/p9c/matrjoska/pkg/db/ffldb"
	// )
	//
	// Create a database and schedule it to be closed and removed on exit. Typically you wouldn't want to remove the
	// database right away like this, nor put it in the temp directory, but it's done here to ensure the example cleans
	// up after itself.
	dbPath := filepath.Join(os.TempDir(), "examplecreate")
	db, e := database.Create("ffldb", dbPath, wire.MainNet)
	if e != nil {
		return
	}
	defer func() {
		if e := os.RemoveAll(dbPath); database.E.Chk(e) {
		}
	}()
	defer func() {
		if e := db.Close(); database.E.Chk(e) {
		}
	}()
	// Output:
}

// This example demonstrates creating a new database and using a managed read-write transaction to store and retrieve
// metadata.
func Example_basicUsage() {
	// This example assumes the ffldb driver is imported.
	//
	// import (
	// 	"github.com/p9c/matrjoska/pkg/db"
	// 	_ "github.com/p9c/matrjoska/pkg/db/ffldb"
	// )
	//
	// Create a database and schedule it to be closed and removed on exit. Typically you wouldn't want to remove the
	// database right away like this, nor put it in the temp directory, but it's done here to ensure the example cleans
	// up after itself.
	dbPath := filepath.Join(os.TempDir(), "exampleusage")
	var e error
	var db database.DB
	db, e = database.Create("ffldb", dbPath, wire.MainNet)
	if e != nil {
		return
	}
	defer func() {
		if e = os.RemoveAll(dbPath); database.E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); database.E.Chk(e) {
		}
	}()
	// Use the Update function of the database to perform a managed read-write transaction. The transaction will
	// automatically be rolled back if the supplied inner function returns a non-nil error.
	e = db.Update(
		func(tx database.Tx) (e error) {
			// Store a key/value pair directly in the metadata bucket. Typically a nested bucket would be used for a given
			// feature, but this example is using the metadata bucket directly for simplicity.
			key := []byte("mykey")
			value := []byte("myvalue")
			if e = tx.Metadata().Put(key, value); E.Chk(e) {
				return e
			}
			// Read the key back and ensure it matches.
			if !bytes.Equal(tx.Metadata().Get(key), value) {
				return fmt.Errorf("unexpected value for key '%s'", key)
			}
			// Create a new nested bucket under the metadata bucket.
			nestedBucketKey := []byte("mybucket")
			var nestedBucket database.Bucket
			nestedBucket, e = tx.Metadata().CreateBucket(nestedBucketKey)
			if e != nil {
				return e
			}
			// The key from above that was set in the metadata bucket does not exist in this new nested bucket.
			if nestedBucket.Get(key) != nil {
				return fmt.Errorf("key '%s' is not expected nil", key)
			}
			return nil
		},
	)
	if e != nil {
		return
	}
	// Output:
}

// // This example demonstrates creating a new database, using a managed read-write
// // transaction to store a block, and using a managed read-only transaction to
// // fetch the block.
// func Example_blockStorageAndRetrieval() {
// 	// This example assumes the ffldb driver is imported.
// 	//
// 	// import (
// 	// 	"github.com/p9c/matrjoska/pkg/db"
// 	// 	_ "github.com/p9c/matrjoska/pkg/db/ffldb"
// 	// )
// 	// Create a database and schedule it to be closed and removed on exit.
// 	// Typically you wouldn't want to remove the database right away like
// 	// this, nor put it in the temp directory, but it's done here to ensure
// 	// the example cleans up after itself.
// 	dbPath := filepath.Join(os.TempDir(), "exampleblkstorage")
// 	db, e := database.Create("ffldb", dbPath, wire.MainNet)
// 	if e != nil  {
// 		DB// 		return
// 	}
// 	defer os.RemoveAll(dbPath)
// 	defer db.Close()
// 	// Use the Update function of the database to perform a managed
// 	// read-write transaction and store a genesis block in the database as
// 	// and example.
// 	e = db.Update(func(tx database.Tx) (e error) {
// 		genesisBlock := chaincfg.MainNetParams.GenesisBlock
// 		return tx.StoreBlock(util.NewBlock(genesisBlock))
// 	})
// 	if e != nil  {
// 		DB// 		return
// 	}
// 	// Use the View function of the database to perform a managed read-only
// 	// transaction and fetch the block stored above.
// 	var loadedBlockBytes []byte
// 	e = db.Update(func(tx database.Tx) (e error) {
// 		genesisHash := chaincfg.MainNetParams.GenesisHash
// 		blockBytes, e := tx.FetchBlock(genesisHash)
// 		if e != nil  {
// 			return e
// 		}
// 		// As documented, all data fetched from the database is only
// 		// valid during a database transaction in order to support
// 		// zero-copy backends.  Thus, make a copy of the data so it
// 		// can be used outside of the transaction.
// 		loadedBlockBytes = make([]byte, len(blockBytes))
// 		copy(loadedBlockBytes, blockBytes)
// 		return nil
// 	})
// 	if e != nil  {
// 		DB// 		return
// 	}
// 	// Typically at this point, the block could be deserialized via the
// 	// wire.Block.Deserialize function or used in its serialized form
// 	// depending on need.  However, for this example, just display the
// 	// number of serialized bytes to show it was loaded as expected.
// 	fmt.Printf("Serialized block size: %d bytes\n", len(loadedBlockBytes))
// 	// Output:
// 	// Serialized block size: 285 bytes
// }
