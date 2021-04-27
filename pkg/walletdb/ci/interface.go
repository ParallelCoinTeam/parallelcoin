package ci

import (
	"fmt"
	"os"
	"reflect"
	
	"github.com/p9c/matrjoska/pkg/walletdb"
)

// errSubTestFail is used to signal that a sub test returned false.
var errSubTestFail = fmt.Errorf("sub test failure")

// testContext is used to store context information about a running test which is passed into helper functions.
type testContext struct {
	t           Tester
	db          walletdb.DB
	bucketDepth int
	isWritable  bool
}

// rollbackValues returns a copy of the provided map with all values set to an empty string. This is used to test that
// values are properly rolled back.
func rollbackValues(
	values map[string]string,
) map[string]string {
	retMap := make(map[string]string, len(values))
	for k := range values {
		retMap[k] = ""
	}
	return retMap
}

// testGetValues checks that all of the provided key/value pairs can be retrieved from the database and the retrieved
// values match the provided values.
func testGetValues(
	tc *testContext, bucket walletdb.ReadBucket, values map[string]string,
) bool {
	for k, v := range values {
		var vBytes []byte
		if v != "" {
			vBytes = []byte(v)
		}
		gotValue := bucket.Get([]byte(k))
		if !reflect.DeepEqual(gotValue, vBytes) {
			tc.t.Errorf("Get: unexpected value - got %s, want %s",
				gotValue, vBytes,
			)
			return false
		}
	}
	return true
}

// testPutValues stores all of the provided key/value pairs in the provided bucket while checking for errors.
func testPutValues(
	tc *testContext, bucket walletdb.ReadWriteBucket, values map[string]string,
) bool {
	for k, v := range values {
		var vBytes []byte
		if v != "" {
			vBytes = []byte(v)
		}
		if e := bucket.Put([]byte(k), vBytes); E.Chk(e) {
			tc.t.Errorf("Put: unexpected error: %v", e)
			return false
		}
	}
	return true
}

// testDeleteValues removes all of the provided key/value pairs from the provided bucket.
func testDeleteValues(
	tc *testContext, bucket walletdb.ReadWriteBucket, values map[string]string,
) bool {
	for k := range values {
		if e := bucket.Delete([]byte(k)); E.Chk(e) {
			tc.t.Errorf("Delete: unexpected error: %v", e)
			return false
		}
	}
	return true
}

// testNestedReadWriteBucket reruns the testBucketInterface against a nested bucket along with a counter to only test a
// couple of level deep.
func testNestedReadWriteBucket(
	tc *testContext, testBucket walletdb.ReadWriteBucket,
) bool {
	// Don't go more than 2 nested level deep.
	if tc.bucketDepth > 1 {
		return true
	}
	tc.bucketDepth++
	defer func() {
		tc.bucketDepth--
	}()
	return testReadWriteBucketInterface(tc, testBucket)
}

// testReadWriteBucketInterface ensures the bucket interface is working properly by exercising all of its functions.
func testReadWriteBucketInterface(
	tc *testContext, bucket walletdb.ReadWriteBucket,
) bool {
	// keyValues holds the keys and values to use when putting values into the bucket.
	var keyValues = map[string]string{
		"bucketkey1": "foo1",
		"bucketkey2": "foo2",
		"bucketkey3": "foo3",
	}
	if !testPutValues(tc, bucket, keyValues) {
		return false
	}
	if !testGetValues(tc, bucket, keyValues) {
		return false
	}
	// Iterate all of the keys using ForEach while making sure the stored values are the expected values.
	keysFound := make(map[string]struct{}, len(keyValues))
	e := bucket.ForEach(func(k, v []byte) (e error) {
		ks := string(k)
		wantV, ok := keyValues[ks]
		if !ok {
			return fmt.Errorf("ForEach: key '%s' should "+
				"exist", ks,
			)
		}
		if !reflect.DeepEqual(v, []byte(wantV)) {
			return fmt.Errorf("ForEach: value for key '%s' "+
				"does not match - got %s, want %s",
				ks, v, wantV,
			)
		}
		keysFound[ks] = struct{}{}
		return nil
	},
	)
	if e != nil {
		tc.t.Errorf("%v", e)
		return false
	}
	// Ensure all keys were iterated.
	for k := range keyValues {
		if _, ok := keysFound[k]; !ok {
			tc.t.Errorf("ForEach: key '%s' was not iterated "+
				"when it should have been", k,
			)
			return false
		}
	}
	// Delete the keys and ensure they were deleted.
	if !testDeleteValues(tc, bucket, keyValues) {
		return false
	}
	if !testGetValues(tc, bucket, rollbackValues(keyValues)) {
		return false
	}
	// Ensure creating a new bucket works as expected.
	testBucketName := []byte("testbucket")
	testBucket, e := bucket.CreateBucket(testBucketName)
	if e != nil {
		tc.t.Errorf("CreateBucket: unexpected error: %v", e)
		return false
	}
	if !testNestedReadWriteBucket(tc, testBucket) {
		return false
	}
	// Ensure creating a bucket that already exists fails with the expected error.
	wantErr := walletdb.ErrBucketExists
	if _, e = bucket.CreateBucket(testBucketName); e != wantErr {
		tc.t.Errorf("CreateBucket: unexpected error - got %v, "+
			"want %v", e, wantErr,
		)
		return false
	}
	// Ensure CreateBucketIfNotExists returns an existing bucket.
	testBucket, e = bucket.CreateBucketIfNotExists(testBucketName)
	if e != nil {
		tc.t.Errorf("CreateBucketIfNotExists: unexpected "+
			"error: %v", e,
		)
		return false
	}
	if !testNestedReadWriteBucket(tc, testBucket) {
		return false
	}
	// Ensure retrieving and existing bucket works as expected.
	testBucket = bucket.NestedReadWriteBucket(testBucketName)
	if !testNestedReadWriteBucket(tc, testBucket) {
		return false
	}
	// Ensure deleting a bucket works as intended.
	if e = bucket.DeleteNestedBucket(testBucketName); E.Chk(e) {
		tc.t.Errorf("DeleteNestedBucket: unexpected error: %v", e)
		return false
	}
	if b := bucket.NestedReadWriteBucket(testBucketName); b != nil {
		tc.t.Errorf("DeleteNestedBucket: bucket '%s' still exists",
			testBucketName,
		)
		return false
	}
	// Ensure deleting a bucket that doesn't exist returns the expected error.
	wantErr = walletdb.ErrBucketNotFound
	if e = bucket.DeleteNestedBucket(testBucketName); e != wantErr {
		tc.t.Errorf("DeleteNestedBucket: unexpected error - got %v, "+
			"want %v", e, wantErr,
		)
		return false
	}
	// Ensure CreateBucketIfNotExists creates a new bucket when it doesn't already exist.
	testBucket, e = bucket.CreateBucketIfNotExists(testBucketName)
	if e != nil {
		tc.t.Errorf("CreateBucketIfNotExists: unexpected "+
			"error: %v", e,
		)
		return false
	}
	if !testNestedReadWriteBucket(tc, testBucket) {
		return false
	}
	// Delete the test bucket to avoid leaving it around for future calls.
	if e := bucket.DeleteNestedBucket(testBucketName); E.Chk(e) {
		tc.t.Errorf("DeleteNestedBucket: unexpected error: %v", e)
		return false
	}
	if b := bucket.NestedReadWriteBucket(testBucketName); b != nil {
		tc.t.Errorf("DeleteNestedBucket: bucket '%s' still exists",
			testBucketName,
		)
		return false
	}
	return true
}

// testManualTxInterface ensures that manual transactions work as expected.
func testManualTxInterface(
	tc *testContext, bucketKey []byte,
) bool {
	db := tc.db
	// populateValues tests that populating values works as expected.
	//
	// When the writable flag is false, a read-only tranasction is created, standard bucket tests for read-only
	// transactions are performed, and the Commit function is checked to ensure it fails as expected.
	//
	// Otherwise, a read-write transaction is created, the values are written, standard bucket tests for read-write
	// transactions are performed, and then the transaction is either commited or rolled back depending on the flag.
	populateValues := func(writable, rollback bool, putValues map[string]string) bool {
		var dbtx walletdb.ReadTx
		var rootBucket walletdb.ReadBucket
		var e error
		if writable {
			dbtx, e = db.BeginReadWriteTx()
			if e != nil {
				tc.t.Errorf("BeginReadWriteTx: unexpected error %v", e)
				return false
			}
			rootBucket = dbtx.(walletdb.ReadWriteTx).ReadWriteBucket(bucketKey)
		} else {
			dbtx, e = db.BeginReadTx()
			if e != nil {
				tc.t.Errorf("BeginReadTx: unexpected error %v", e)
				return false
			}
			rootBucket = dbtx.ReadBucket(bucketKey)
		}
		if rootBucket == nil {
			tc.t.Errorf("ReadWriteBucket/ReadBucket: unexpected nil root bucket")
			_ = dbtx.Rollback()
			return false
		}
		if writable {
			tc.isWritable = writable
			if !testReadWriteBucketInterface(tc, rootBucket.(walletdb.ReadWriteBucket)) {
				_ = dbtx.Rollback()
				return false
			}
		}
		if !writable {
			// Rollback the transaction.
			if e := dbtx.Rollback(); E.Chk(e) {
				tc.t.Errorf("Commit: unexpected error %v", e)
				return false
			}
		} else {
			rootBucket := rootBucket.(walletdb.ReadWriteBucket)
			if !testPutValues(tc, rootBucket, putValues) {
				return false
			}
			if rollback {
				// Rollback the transaction.
				if e := dbtx.Rollback(); E.Chk(e) {
					tc.t.Errorf("Rollback: unexpected "+
						"error %v", e,
					)
					return false
				}
			} else {
				// The commit should succeed.
				if e := dbtx.(walletdb.ReadWriteTx).Commit(); E.Chk(e) {
					tc.t.Errorf("Commit: unexpected error "+
						"%v", e,
					)
					return false
				}
			}
		}
		return true
	}
	// checkValues starts a read-only transaction and checks that all of the key/value pairs specified in the
	// expectedValues parameter match what's in the database.
	checkValues := func(expectedValues map[string]string) bool {
		// Begin another read-only transaction to ensure...
		dbtx, e := db.BeginReadTx()
		if e != nil {
			tc.t.Errorf("BeginReadTx: unexpected error %v", e)
			return false
		}
		rootBucket := dbtx.ReadBucket(bucketKey)
		if rootBucket == nil {
			tc.t.Errorf("ReadBucket: unexpected nil root bucket")
			_ = dbtx.Rollback()
			return false
		}
		if !testGetValues(tc, rootBucket, expectedValues) {
			_ = dbtx.Rollback()
			return false
		}
		// Rollback the read-only transaction.
		if e := dbtx.Rollback(); E.Chk(e) {
			tc.t.Errorf("Commit: unexpected error %v", e)
			return false
		}
		return true
	}
	// deleteValues starts a read-write transaction and deletes the keys in the passed key/value pairs.
	deleteValues := func(values map[string]string) bool {
		dbtx, e := db.BeginReadWriteTx()
		if e != nil {
			tc.t.Errorf("BeginReadWriteTx: unexpected error %v", e)
			_ = dbtx.Rollback()
			return false
		}
		rootBucket := dbtx.ReadWriteBucket(bucketKey)
		if rootBucket == nil {
			tc.t.Errorf("RootBucket: unexpected nil root bucket")
			_ = dbtx.Rollback()
			return false
		}
		// Delete the keys and ensure they were deleted.
		if !testDeleteValues(tc, rootBucket, values) {
			_ = dbtx.Rollback()
			return false
		}
		if !testGetValues(tc, rootBucket, rollbackValues(values)) {
			_ = dbtx.Rollback()
			return false
		}
		// Commit the changes and ensure it was successful.
		if e := dbtx.Commit(); E.Chk(e) {
			tc.t.Errorf("Commit: unexpected error %v", e)
			return false
		}
		return true
	}
	// keyValues holds the keys and values to use when putting values into a bucket.
	var keyValues = map[string]string{
		"umtxkey1": "foo1",
		"umtxkey2": "foo2",
		"umtxkey3": "foo3",
	}
	// Ensure that attempting populating the values using a read-only transaction fails as expected.
	if !populateValues(false, true, keyValues) {
		return false
	}
	if !checkValues(rollbackValues(keyValues)) {
		return false
	}
	// Ensure that attempting populating the values using a read-write transaction and then rolling it back yields the
	// expected values.
	if !populateValues(true, true, keyValues) {
		return false
	}
	if !checkValues(rollbackValues(keyValues)) {
		return false
	}
	// Ensure that attempting populating the values using a read-write transaction and then committing it stores the
	// expected values.
	if !populateValues(true, false, keyValues) {
		return false
	}
	if !checkValues(keyValues) {
		return false
	}
	// Clean up the keys.
	if !deleteValues(keyValues) {
		return false
	}
	return true
}

// testNamespaceAndTxInterfaces creates a namespace using the provided key and tests all facets of it interface as well
// as transaction and bucket interfaces under it.
func testNamespaceAndTxInterfaces(
	tc *testContext, namespaceKey string,
) bool {
	namespaceKeyBytes := []byte(namespaceKey)
	e := walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
		_, e = tx.CreateTopLevelBucket(namespaceKeyBytes)
		return e
	},
	)
	if e != nil {
		tc.t.Errorf("CreateTopLevelBucket: unexpected error: %v", e)
		return false
	}
	defer func() {
		// Remove the namespace now that the tests are done for it.
		e = walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			return tx.DeleteTopLevelBucket(namespaceKeyBytes)
		},
		)
		if e != nil {
			tc.t.Errorf("DeleteTopLevelBucket: unexpected error: %v", e)
			return
		}
	}()
	if !testManualTxInterface(tc, namespaceKeyBytes) {
		return false
	}
	// keyValues holds the keys and values to use when putting values into a bucket.
	var keyValues = map[string]string{
		"mtxkey1": "foo1",
		"mtxkey2": "foo2",
		"mtxkey3": "foo3",
	}
	// Test the bucket interface via a managed read-only transaction.
	e = walletdb.View(tc.db, func(tx walletdb.ReadTx) (e error) {
		rootBucket := tx.ReadBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadBucket: unexpected nil root bucket")
		}
		return nil
	},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	// Test the bucket interface via a managed read-write transaction. Also, put a series of values and force a rollback
	// so the following code can ensure the values were not stored.
	forceRollbackError := fmt.Errorf("force rollback")
	e = walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
		rootBucket := tx.ReadWriteBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadWriteBucket: unexpected nil root bucket")
		}
		tc.isWritable = true
		if !testReadWriteBucketInterface(tc, rootBucket) {
			return errSubTestFail
		}
		if !testPutValues(tc, rootBucket, keyValues) {
			return errSubTestFail
		}
		// Return an error to force a rollback.
		return forceRollbackError
	},
	)
	if e != forceRollbackError {
		if e == errSubTestFail {
			return false
		}
		tc.t.Errorf("Update: inner function error not returned - got "+
			"%v, want %v", e, forceRollbackError,
		)
		return false
	}
	// Ensure the values that should have not been stored due to the forced rollback above were not actually stored.
	e = walletdb.View(tc.db, func(tx walletdb.ReadTx) (e error) {
		rootBucket := tx.ReadBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadBucket: unexpected nil root bucket")
		}
		if !testGetValues(tc, rootBucket, rollbackValues(keyValues)) {
			return errSubTestFail
		}
		return nil
	},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	// Store a series of values via a managed read-write transaction.
	e = walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
		rootBucket := tx.ReadWriteBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadWriteBucket: unexpected nil root bucket")
		}
		if !testPutValues(tc, rootBucket, keyValues) {
			return errSubTestFail
		}
		return nil
	},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	// Ensure the values stored above were committed as expected.
	e = walletdb.View(tc.db, func(tx walletdb.ReadTx) (e error) {
		rootBucket := tx.ReadBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadBucket: unexpected nil root bucket")
		}
		if !testGetValues(tc, rootBucket, keyValues) {
			return errSubTestFail
		}
		return nil
	},
	)
	if e != nil {
		E.Ln(e)
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	// Clean up the values stored above in a managed read-write transaction.
	e = walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
		rootBucket := tx.ReadWriteBucket(namespaceKeyBytes)
		if rootBucket == nil {
			return fmt.Errorf("ReadWriteBucket: unexpected nil root bucket")
		}
		if !testDeleteValues(tc, rootBucket, keyValues) {
			return errSubTestFail
		}
		return nil
	},
	)
	if e != nil {
		E.Ln(e)
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	return true
}

// testAdditionalErrors performs some tests for error cases not covered elsewhere in the tests and therefore improves
// negative test coverage.
func testAdditionalErrors(
	tc *testContext,
) bool {
	ns3Key := []byte("ns3")
	e := walletdb.Update(tc.db, func(tx walletdb.ReadWriteTx) (e error) {
		// Create a new namespace
		rootBucket, e := tx.CreateTopLevelBucket(ns3Key)
		if e != nil {
			return fmt.Errorf("CreateTopLevelBucket: unexpected error: %v", e)
		}
		// Ensure CreateBucket returns the expected error when no bucket key is specified.
		wantErr := walletdb.ErrBucketNameRequired
		if _, e = rootBucket.CreateBucket(nil); e != wantErr {
			return fmt.Errorf("CreateBucket: unexpected error - "+
				"got %v, want %v", e, wantErr,
			)
		}
		// Ensure DeleteNestedBucket returns the expected error when no bucket key is specified.
		wantErr = walletdb.ErrIncompatibleValue
		if e := rootBucket.DeleteNestedBucket(nil); e != wantErr {
			return fmt.Errorf("DeleteNestedBucket: unexpected error - "+
				"got %v, want %v", e, wantErr,
			)
		}
		// Ensure Put returns the expected error when no key is specified.
		wantErr = walletdb.ErrKeyRequired
		if e := rootBucket.Put(nil, nil); e != wantErr {
			return fmt.Errorf("put: unexpected error - got %v, want %v", e, wantErr)
		}
		return nil
	},
	)
	if e != nil {
		if e != errSubTestFail {
			tc.t.Errorf("%v", e)
		}
		return false
	}
	// Ensure that attempting to rollback or commit a transaction that is already closed returns the expected error.
	tx, e := tc.db.BeginReadWriteTx()
	if e != nil {
		tc.t.Errorf("Begin: unexpected error: %v", e)
		return false
	}
	if e := tx.Rollback(); E.Chk(e) {
		tc.t.Errorf("Rollback: unexpected error: %v", e)
		return false
	}
	wantErr := walletdb.ErrTxClosed
	if e := tx.Rollback(); e != wantErr {
		tc.t.Errorf("Rollback: unexpected error - got %v, want %v", e,
			wantErr,
		)
		return false
	}
	if e := tx.Commit(); e != wantErr {
		tc.t.Errorf("Commit: unexpected error - got %v, want %v", e,
			wantErr,
		)
		return false
	}
	return true
}

// TestInterface performs all interfaces tests for this database driver.
func TestInterface(
	t Tester, dbType, dbPath string,
) {
	db, e := walletdb.Create(dbType, dbPath)
	if e != nil {
		t.Errorf("Failed to create test database (%s) %v", dbType, e)
		return
	}
	defer func() {
		if e := os.Remove(dbPath); E.Chk(e) {
		}
	}()
	defer func() {
		if e := db.Close(); E.Chk(e) {
		}
	}()
	// Run all of the interface tests against the database. Create a test context to pass around.
	context := testContext{t: t, db: db}
	// Create a namespace and test the interface for it.
	if !testNamespaceAndTxInterfaces(&context, "ns1") {
		return
	}
	// Create a second namespace and test the interface for it.
	if !testNamespaceAndTxInterfaces(&context, "ns2") {
		return
	}
	// Chk a few more error conditions not covered elsewhere.
	if !testAdditionalErrors(&context) {
		return
	}
}
