package walletdb_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/p9c/matrjoska/pkg/walletdb"
	_ "github.com/p9c/matrjoska/pkg/walletdb/bdb"
)

// This example demonstrates creating a new database.
func ExampleCreate() {
	// This example assumes the bdb (bolt db) driver is imported.
	//
	// import (
	// 	"github.com/p9c/matrjoska/cmd/wallet/db"
	// 	_ "github.com/p9c/matrjoska/cmd/wallet/db/bdb"
	// )
	//
	// Create a database and schedule it to be closed and removed on exit. Typically you wouldn't want to remove the
	// database right away like this, but it's done here in the example to ensure the example cleans up after itself.
	dbPath := filepath.Join(os.TempDir(), "examplecreate.db")
	db, e := walletdb.Create("bdb", dbPath)
	if E.Chk(e) {
		return
	}
	defer func() {
		if e := os.Remove(dbPath); walletdb.E.Chk(e) {
		}
	}()
	defer func() {
		if e := db.Close(); walletdb.E.Chk(e) {
		}
	}()
	// Output:
}

// exampleNum is used as a counter in the exampleLoadDB function to provided a unique database name for each example.
var exampleNum = 0

// exampleLoadDB is used in the examples to elide the setup code.
func exampleLoadDB() (db walletdb.DB, teardownFunc func(), e error) {
	dbName := fmt.Sprintf("exampleload%d.db", exampleNum)
	dbPath := filepath.Join(os.TempDir(), dbName)
	db, e = walletdb.Create("bdb", dbPath)
	if e != nil {
		return nil, nil, e
	}
	teardownFunc = func() {
		if e = db.Close(); walletdb.E.Chk(e) {
		}
		if e = os.Remove(dbPath); walletdb.E.Chk(e) {
		}
	}
	exampleNum++
	return db, teardownFunc, e
}

// This example demonstrates creating a new top level bucket.
func ExampleDB_createTopLevelBucket() {
	// Load a database for the purposes of this example and schedule it to be closed and removed on exit. See the Create
	// example for more details on what this step is doing.
	db, teardownFunc, e := exampleLoadDB()
	if E.Chk(e) {
		return
	}
	defer teardownFunc()
	dbtx, e := db.BeginReadWriteTx()
	if E.Chk(e) {
		return
	}
	defer func() {
		e = dbtx.Commit()
		if E.Chk(e) {
			return
		}
	}()
	// Get or create a bucket in the database as needed. This bucket is what is typically passed to specific
	// sub-packages so they have their own area to work in without worrying about conflicting keys.
	bucketKey := []byte("walletsubpackage")
	bucket, e := dbtx.CreateTopLevelBucket(bucketKey)
	if E.Chk(e) {
		return
	}
	// Prevent unused error.
	_ = bucket
	// Output:
}

// This example demonstrates creating a new database, getting a namespace from it, and using a managed read-write
// transaction against the namespace to store and retrieve data.
func Example_basicUsage() {
	// This example assumes the bdb (bolt db) driver is imported.
	//
	// import (
	// 	"github.com/p9c/matrjoska/cmd/wallet/db"
	// 	_ "github.com/p9c/matrjoska/cmd/wallet/db/bdb"
	// )
	//
	// Create a database and schedule it to be closed and removed on exit. Typically you wouldn't want to remove the
	// database right away like this, but it's done here in the example to ensure the example cleans up after itself.
	dbPath := filepath.Join(os.TempDir(), "exampleusage.db")
	db, e := walletdb.Create("bdb", dbPath)
	if E.Chk(e) {
		return
	}
	defer func() {
		if e = os.Remove(dbPath); walletdb.E.Chk(e) {
		}
	}()
	defer func() {
		if e = db.Close(); walletdb.E.Chk(e) {
		}
	}()
	// Get or create a bucket in the database as needed. This bucket is what is typically passed to specific
	// sub-packages so they have their own area to work in without worrying about conflicting keys.
	bucketKey := []byte("walletsubpackage")
	e = walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			bucket := tx.ReadWriteBucket(bucketKey)
			if bucket == nil {
				_, e = tx.CreateTopLevelBucket(bucketKey)
				if e != nil {
					return e
				}
			}
			return nil
		},
	)
	if E.Chk(e) {
		return
	}
	// Use the Update function of the namespace to perform a managed read-write transaction. The transaction will
	// automatically be rolled back if the supplied inner function returns a non-nil error.
	e = walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			// All data is stored against the root bucket of the namespace, or nested buckets of the root bucket. It's not
			// really necessary to store it in a separate variable like this, but it has been done here for the purposes of
			// the example to illustrate.
			rootBucket := tx.ReadWriteBucket(bucketKey)
			// Store a key/value pair directly in the root bucket.
			key := []byte("mykey")
			value := []byte("myvalue")
			if e = rootBucket.Put(key, value); E.Chk(e) {
				return e
			}
			// Read the key back and ensure it matches.
			if !bytes.Equal(rootBucket.Get(key), value) {
				return fmt.Errorf("unexpected value for key '%s'", key)
			}
			// Create a new nested bucket under the root bucket.
			nestedBucketKey := []byte("mybucket")
			nestedBucket, e := rootBucket.CreateBucket(nestedBucketKey)
			if e != nil {
				return e
			}
			// The key from above that was set in the root bucket does not exist in this new nested bucket.
			if nestedBucket.Get(key) != nil {
				return fmt.Errorf("key '%s' is not expected nil", key)
			}
			return nil
		},
	)
	if E.Chk(e) {
		return
	}
	// Output:
}
