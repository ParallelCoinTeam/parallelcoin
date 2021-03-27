package headerfs

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"
	
	"github.com/p9c/pod/pkg/walletdb"
	_ "github.com/p9c/pod/pkg/walletdb/bdb"
)

func createTestIndex() (func(), *headerIndex, error) {
	tempDir, e := ioutil.TempDir("", "neutrino")
	if e != nil {
		return nil, nil, e
	}
	db, e := walletdb.Create("bdb", tempDir+"/test.db")
	if e != nil {
		return nil, nil, e
	}
	cleanUp := func() {
		if e = os.RemoveAll(tempDir); E.Chk(e) {
		}
		if e = db.Close(); E.Chk(e) {
		}
	}
	filterDB, e := newHeaderIndex(db, Block)
	if e != nil {
		return nil, nil, e
	}
	return cleanUp, filterDB, nil
}

func TestAddHeadersIndexRetrieve(t *testing.T) {
	var e error
	var hIndex *headerIndex
	var cleanUp func()
	if cleanUp, hIndex, e = createTestIndex(); !E.Chk(e) {
		defer cleanUp()
	} else {
		t.Fatalf("unable to create test db: %v", e)
	}
	// First, we'll create a a series of random headers that we'll use to write into the database.
	const numHeaders = 100
	headerEntries := make(headerBatch, numHeaders)
	headerIndex := make(map[uint32]headerEntry)
	for i := uint32(0); i < numHeaders; i++ {
		var header headerEntry
		if _, e = rand.Read(header.hash[:]); E.Chk(e) {
			t.Fatalf("unable to read header: %v", e)
		}
		header.height = i
		headerEntries[i] = header
		headerIndex[i] = header
	}
	// With the headers constructed, we'll write them to disk in a single batch.
	if e = hIndex.addHeaders(headerEntries); E.Chk(e) {
		t.Fatalf("unable to add headers: %v", e)
	}
	// Next, verify that the database tip matches the _final_ header inserted.
	dbTip, dbHeight, e := hIndex.chainTip()
	if e != nil {
		t.Fatalf("unable to obtain chain tip: %v", e)
	}
	lastEntry := headerIndex[numHeaders-1]
	if dbHeight != lastEntry.height {
		t.Fatalf(
			"height doesn't match: expected %v, got %v",
			lastEntry.height, dbHeight,
		)
	}
	if !bytes.Equal(dbTip[:], lastEntry.hash[:]) {
		t.Fatalf(
			"tip doesn't match: expected %x, got %x",
			lastEntry.hash[:], dbTip[:],
		)
	}
	// For each header written, check that we're able to retrieve the entry both by hash and height.
	for i, he := range headerEntries {
		var height uint32
		height, e = hIndex.heightFromHash(&he.hash)
		if e != nil {
			t.Fatalf("unable to retreive height(%v): %v", i, e)
		}
		if height != he.height {
			t.Fatalf(
				"height doesn't match: expected %v, got %v",
				he.height, height,
			)
		}
	}
	// Next if we truncate the index by one, then we should end up at the second to last entry for the tip.
	newTip := headerIndex[numHeaders-2]
	if e = hIndex.truncateIndex(&newTip.hash, true); E.Chk(e) {
		t.Fatalf("unable to truncate index: %v", e)
	}
	// This time the database tip should be the _second_ to last entry inserted.
	dbTip, dbHeight, e = hIndex.chainTip()
	if e != nil {
		t.Fatalf("unable to obtain chain tip: %v", e)
	}
	lastEntry = headerIndex[numHeaders-2]
	if dbHeight != lastEntry.height {
		t.Fatalf(
			"height doesn't match: expected %v, got %v",
			lastEntry.height, dbHeight,
		)
	}
	if !bytes.Equal(dbTip[:], lastEntry.hash[:]) {
		t.Fatalf(
			"tip doesn't match: expected %x, got %x",
			lastEntry.hash[:], dbTip[:],
		)
	}
}
