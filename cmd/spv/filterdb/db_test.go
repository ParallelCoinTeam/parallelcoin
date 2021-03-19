package filterdb

import (
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"testing"
	
	"github.com/p9c/pod/pkg/blockchain/chaincfg"
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	"github.com/p9c/pod/pkg/blockchain/chainhash"
	"github.com/p9c/pod/pkg/coding/gcs"
	"github.com/p9c/pod/pkg/coding/gcs/builder"
	"github.com/p9c/pod/pkg/database/walletdb"
	_ "github.com/p9c/pod/pkg/database/walletdb/bdb"
)

func createTestDatabase() (func(), FilterDatabase, error) {
	tempDir, e := ioutil.TempDir("", "neutrino")
	if e != nil {
		return nil, nil, e
	}
	db, e := walletdb.Create("bdb", tempDir+"/test.db")
	if e != nil {
		return nil, nil, e
	}
	cleanUp := func() {
		if e := os.RemoveAll(tempDir); E.Chk(e) {
		}
		if e := db.Close(); E.Chk(e) {
		}
	}
	filterDB, e := New(db, netparams.SimNetParams)
	if e != nil {
		return nil, nil, e
	}
	return cleanUp, filterDB, nil
}

func TestGenesisFilterCreation(t *testing.T) {
	var e error
	var cleanUp func()
	var dB FilterDatabase
	if cleanUp, dB, e = createTestDatabase(); !E.Chk(e) {
		defer cleanUp()
	} else {
		t.Fatalf("unable to create test db: %v", err)
	}
	genesisHash := chaincfg.SimNetParams.GenesisHash
	// With the database initialized, we should be able to fetch the
	// regular filter for the genesis block.
	regGenesisFilter, e := dB.FetchFilter(genesisHash, RegularFilter)
	if e != nil {
		t.Fatalf("unable to fetch regular genesis filter: %v", err)
	}
	// The regular filter should be non-nil as the gensis block's output and the coinbase txid should be indexed.
	if regGenesisFilter == nil {
		t.Fatalf("regular genesis filter is nil")
	}
	
}
func genRandFilter(numElements uint32) (*gcs.Filter, error) {
	elements := make([][]byte, numElements)
	for i := uint32(0); i < numElements; i++ {
		var elem [20]byte
		if _, e = rand.Read(elem[:]); E.Chk(e) {
			return nil, e
		}
		elements[i] = elem[:]
	}
	var key [16]byte
	if _, e = rand.Read(key[:]); E.Chk(e) {
		return nil, e
	}
	filter, e := gcs.BuildGCSFilter(
		builder.DefaultP, builder.DefaultM, key, elements,
	)
	if e != nil {
		return nil, e
	}
	return filter, nil
}

func TestFilterStorage(t *testing.T) {
	// TODO(roasbeef): use testing.Quick
	var cleanUp func()
	var dB FilterDatabase
	var e error
	if cleanUp, dB, e = createTestDatabase(); !E.Chk(e) {
		defer cleanUp()
	} else {
		t.Fatalf("unable to create test db: %v", err)
	}
	// We'll generate a random block hash to create our test filters against.
	var randHash chainhash.Hash
	if _, e = rand.Read(randHash[:]); E.Chk(e) {
		t.Fatalf("unable to generate random hash: %v", err)
	}
	// First, we'll create and store a random fitler for the regular filter type for the block hash generate above.
	regFilter, e := genRandFilter(100)
	if e != nil {
		t.Fatalf("unable to create random filter: %v", err)
	}
	e = dB.PutFilter(&randHash, regFilter, RegularFilter)
	if e != nil {
		t.Fatalf("unable to store regular filter: %v", err)
	}
	// With the filter stored, we should be able to retrieve the filter without any issue, and it should match the
	// stored filter exactly.
	regFilterDB, e := dB.FetchFilter(&randHash, RegularFilter)
	if e != nil {
		t.Fatalf("unable to retrieve reg filter: %v", err)
	}
	if !reflect.DeepEqual(regFilter, regFilterDB) {
		t.Fatalf("regular filter doesn't match!")
	}
}
