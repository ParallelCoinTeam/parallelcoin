package bdb

import (
	"github.com/stalker-loki/app/slog"
	"io"
	"os"

	bolt "github.com/coreos/bbolt"

	"github.com/p9c/pod/pkg/db/walletdb"
)

// convertErr converts some bolt errors to the equivalent walletdb error.
func convertErr(er error) (err error) {
	switch er {
	// Database open/create errors.
	case bolt.ErrDatabaseNotOpen:
		err = walletdb.ErrDbNotOpen
	case bolt.ErrInvalid:
		err = walletdb.ErrInvalid
	// Transaction errors.
	case bolt.ErrTxNotWritable:
		err = walletdb.ErrTxNotWritable
	case bolt.ErrTxClosed:
		err = walletdb.ErrTxClosed
	// Value/bucket errors.
	case bolt.ErrBucketNotFound:
		err = walletdb.ErrBucketNotFound
	case bolt.ErrBucketExists:
		err = walletdb.ErrBucketExists
	case bolt.ErrBucketNameRequired:
		err = walletdb.ErrBucketNameRequired
	case bolt.ErrKeyRequired:
		err = walletdb.ErrKeyRequired
	case bolt.ErrKeyTooLarge:
		err = walletdb.ErrKeyTooLarge
	case bolt.ErrValueTooLarge:
		err = walletdb.ErrValueTooLarge
	case bolt.ErrIncompatibleValue:
		err = walletdb.ErrIncompatibleValue
	}
	// Return the original error if none of the above applies.
	return
}

// transaction represents a database transaction.  It can either by read-only or
// read-write and implements the walletdb Tx interfaces.  The transaction
// provides a root bucket against which all read and writes occur.
type transaction struct {
	boltTx *bolt.Tx
}

func (tx *transaction) ReadBucket(key []byte) walletdb.ReadBucket {
	return tx.ReadWriteBucket(key)
}
func (tx *transaction) ReadWriteBucket(key []byte) walletdb.ReadWriteBucket {
	boltBucket := tx.boltTx.Bucket(key)
	if boltBucket == nil {
		return nil
	}
	return (*bucket)(boltBucket)
}
func (tx *transaction) CreateTopLevelBucket(key []byte) (bkt walletdb.ReadWriteBucket, err error) {
	var boltBucket *bolt.Bucket
	if boltBucket, err = tx.boltTx.CreateBucket(key); slog.Check(err) {
		err = convertErr(err)
		slog.Debug(err)
		return
	}
	bkt = (*bucket)(boltBucket)
	return
}
func (tx *transaction) DeleteTopLevelBucket(key []byte) (err error) {
	if err = tx.boltTx.DeleteBucket(key); slog.Check(err) {
		err = convertErr(err)
	}
	return
}

// Commit commits all changes that have been made through the root bucket and
// all of its sub-buckets to persistent storage.
//
// This function is part of the walletdb.Tx interface implementation.
func (tx *transaction) Commit() (err error) {
	return convertErr(tx.boltTx.Commit())
}

// Rollback undoes all changes that have been made to the root bucket and all of
// its sub-buckets.
//
// This function is part of the walletdb.Tx interface implementation.
func (tx *transaction) Rollback() (err error) {
	return convertErr(tx.boltTx.Rollback())
}

// bucket is an internal type used to represent a collection of key/value pairs
// and implements the walletdb Bucket interfaces.
type bucket bolt.Bucket

// Enforce bucket implements the walletdb Bucket interfaces.
var _ walletdb.ReadWriteBucket = (*bucket)(nil)

// NestedReadWriteBucket retrieves a nested bucket with the given key.  Returns
// nil if the bucket does not exist.
//
// This function is part of the walletdb.ReadWriteBucket interface implementation.
func (b *bucket) NestedReadWriteBucket(key []byte) walletdb.ReadWriteBucket {
	boltBucket := (*bolt.Bucket)(b).Bucket(key)
	// Don't return a non-nil interface to a nil pointer.
	if boltBucket == nil {
		return nil
	}
	return (*bucket)(boltBucket)
}
func (b *bucket) NestedReadBucket(key []byte) walletdb.ReadBucket {
	return b.NestedReadWriteBucket(key)
}

// CreateBucket creates and returns a new nested bucket with the given key.
// Returns ErrBucketExists if the bucket already exists, ErrBucketNameRequired
// if the key is empty, or ErrIncompatibleValue if the key value is otherwise
// invalid.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) CreateBucket(key []byte) (bkt walletdb.ReadWriteBucket, err error) {
	var boltBucket *bolt.Bucket
	if boltBucket, err = (*bolt.Bucket)(b).CreateBucket(key); slog.Check(err) {
		err = convertErr(err)
		slog.Debug(err)
		return
	}
	bkt = (*bucket)(boltBucket)
	return
}

// CreateBucketIfNotExists creates and returns a new nested bucket with the
// given key if it does not already exist.  Returns ErrBucketNameRequired if the
// key is empty or ErrIncompatibleValue if the key value is otherwise invalid.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) CreateBucketIfNotExists(key []byte) (bkt walletdb.ReadWriteBucket, err error) {
	var boltBucket *bolt.Bucket
	if boltBucket, err = (*bolt.Bucket)(b).CreateBucketIfNotExists(key); slog.Check(err) {
		err = convertErr(err)
		slog.Debug(err)
		return
	}
	bkt = (*bucket)(boltBucket)
	return
}

// DeleteNestedBucket removes a nested bucket with the given key.  Returns
// ErrTxNotWritable if attempted against a read-only transaction and
// ErrBucketNotFound if the specified bucket does not exist.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) DeleteNestedBucket(key []byte) (err error) {
	return convertErr((*bolt.Bucket)(b).DeleteBucket(key))
}

// ForEach invokes the passed function with every key/value pair in the bucket.
// This includes nested buckets, in which case the value is nil, but it does not
// include the key/value pairs within those nested buckets.
//
// NOTE: The values returned by this function are only valid during a
// transaction.  Attempting to access them after a transaction has ended will
// likely result in an access violation.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) ForEach(fn func(k, v []byte) error) (err error) {
	return convertErr((*bolt.Bucket)(b).ForEach(fn))
}

// Put saves the specified key/value pair to the bucket.  Keys that do not
// already exist are added and keys that already exist are overwritten.  Returns
// ErrTxNotWritable if attempted against a read-only transaction.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) Put(key, value []byte) (err error) {
	return convertErr((*bolt.Bucket)(b).Put(key, value))
}

// Get returns the value for the given key.  Returns nil if the key does
// not exist in this bucket (or nested buckets).
//
// NOTE: The value returned by this function is only valid during a
// transaction.  Attempting to access it after a transaction has ended
// will likely result in an access violation.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) Get(key []byte) []byte {
	return (*bolt.Bucket)(b).Get(key)
}

// Delete removes the specified key from the bucket.  Deleting a key that does
// not exist does not return an error.  Returns ErrTxNotWritable if attempted
// against a read-only transaction.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) Delete(key []byte) (err error) {
	return convertErr((*bolt.Bucket)(b).Delete(key))
}
func (b *bucket) ReadCursor() walletdb.ReadCursor {
	return b.ReadWriteCursor()
}

// ReadWriteCursor returns a new cursor, allowing for iteration over the bucket's
// key/value pairs and nested buckets in forward or backward order.
//
// This function is part of the walletdb.Bucket interface implementation.
func (b *bucket) ReadWriteCursor() walletdb.ReadWriteCursor {
	return (*cursor)((*bolt.Bucket)(b).Cursor())
}

// cursor represents a cursor over key/value pairs and nested buckets of a
// bucket.
//
// Note that open cursors are not tracked on bucket changes and any
// modifications to the bucket, with the exception of cursor.Delete, invalidate
// the cursor. After invalidation, the cursor must be repositioned, or the keys
// and values returned may be unpredictable.
type cursor bolt.Cursor

// Delete removes the current key/value pair the cursor is at without
// invalidating the cursor. Returns ErrTxNotWritable if attempted on a read-only
// transaction, or ErrIncompatibleValue if attempted when the cursor points to a
// nested bucket.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) Delete() (err error) {
	return convertErr((*bolt.Cursor)(c).Delete())
}

// First positions the cursor at the first key/value pair and returns the pair.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) First() (key, value []byte) {
	return (*bolt.Cursor)(c).First()
}

// Last positions the cursor at the last key/value pair and returns the pair.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) Last() (key, value []byte) {
	return (*bolt.Cursor)(c).Last()
}

// Next moves the cursor one key/value pair forward and returns the new pair.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) Next() (key, value []byte) {
	return (*bolt.Cursor)(c).Next()
}

// Prev moves the cursor one key/value pair backward and returns the new pair.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) Prev() (key, value []byte) {
	return (*bolt.Cursor)(c).Prev()
}

// Seek positions the cursor at the passed seek key. If the key does not exist,
// the cursor is moved to the next key after seek. Returns the new pair.
//
// This function is part of the walletdb.Cursor interface implementation.
func (c *cursor) Seek(seek []byte) (key, value []byte) {
	return (*bolt.Cursor)(c).Seek(seek)
}

// db represents a collection of namespaces which are persisted and implements
// the walletdb.Db interface.  All database access is performed through
// transactions which are obtained through the specific Namespace.
type db bolt.DB

// Enforce db implements the walletdb.Db interface.
var _ walletdb.DB = (*db)(nil)

func (db *db) beginTx(writable bool) (t *transaction, err error) {
	boltTx, err := (*bolt.DB)(db).Begin(writable)
	if err != nil {
		slog.Error(err)
		return nil, convertErr(err)
	}
	return &transaction{boltTx: boltTx}, nil
}
func (db *db) BeginReadTx() (t walletdb.ReadTx, err error) {
	return db.beginTx(false)
}
func (db *db) BeginReadWriteTx() (t walletdb.ReadWriteTx, err error) {
	return db.beginTx(true)
}

// Copy writes a copy of the database to the provided writer.  This call will
// start a read-only transaction to perform all operations.
//
// This function is part of the walletdb.Db interface implementation.
func (db *db) Copy(w io.Writer) (err error) {
	return convertErr((*bolt.DB)(db).View(func(tx *bolt.Tx) (err error) {
		return tx.Copy(w)
	}))
}

// Close cleanly shuts down the database and syncs all data.
//
// This function is part of the walletdb.Db interface implementation.
func (db *db) Close() (err error) {
	return convertErr((*bolt.DB)(db).Close())
}

// filesExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// openDB opens the database at the provided path.  walletdb.ErrDbDoesNotExist
// is returned if the database doesn't exist and the create flag is not set.
func openDB(dbPath string, create bool) (t walletdb.DB, err error) {
	if !create && !fileExists(dbPath) {
		return nil, walletdb.ErrDbDoesNotExist
	}
	var boltDB *bolt.DB
	boltDB, err = bolt.Open(dbPath, 0600, nil)
	return (*db)(boltDB), convertErr(err)
}
