package gcs_test

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/gcs"
)

var (
	// No need to allocate an e variable in every test
	e error
	// Collision probability for the tests (1/2**19)
	P = uint8(19)
	// Modulus value for the tests.
	M uint64 = 784931
	// Filters are conserved between tests but we must define with an interface which functions we're testing because
	// the gcsFilter type isn't exported
	filter, filter2, filter3 /*, filter4, filter5*/ *gcs.Filter
	// We need to use the same key for building and querying the filters
	key [gcs.KeySize]byte
	// List of values for building a filter
	contents = [][]byte{
		[]byte("Alex"),
		[]byte("Bob"),
		[]byte("Charlie"),
		[]byte("Dick"),
		[]byte("Ed"),
		[]byte("Frank"),
		[]byte("George"),
		[]byte("Harry"),
		[]byte("Ilya"),
		[]byte("John"),
		[]byte("Kevin"),
		[]byte("Larry"),
		[]byte("Michael"),
		[]byte("Nate"),
		[]byte("Owen"),
		[]byte("Paul"),
		[]byte("Quentin"),
	}
	// List of values for querying a filter using MatchAny()
	contents2 = [][]byte{
		[]byte("Alice"),
		[]byte("Betty"),
		[]byte("Charmaine"),
		[]byte("Donna"),
		[]byte("Edith"),
		[]byte("Faina"),
		[]byte("Georgia"),
		[]byte("Hannah"),
		[]byte("Ilsbeth"),
		[]byte("Jennifer"),
		[]byte("Kayla"),
		[]byte("Lena"),
		[]byte("Michelle"),
		[]byte("Natalie"),
		[]byte("Ophelia"),
		[]byte("Peggy"),
		[]byte("Queenie"),
	}
)

// TestGCSFilterBuild builds a test filter with a randomized key. For Bitcoin use, deterministic filter generation is
// desired. Therefore, a key that's derived deterministically would be required.
func TestGCSFilterBuild(t *testing.T) {
	for i := 0; i < gcs.KeySize; i += 4 {
		binary.BigEndian.PutUint32(key[i:], rand.Uint32())
	}
	filter, e = gcs.BuildGCSFilter(P, M, key, contents)
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
}

// TestGCSFilterCopy deserializes and serializes a filter to create a copy.
func TestGCSFilterCopy(t *testing.T) {
	var serialized2 []byte
	serialized2, e = filter.Bytes()
	if e != nil {
		t.Fatalf("Filter Hash() failed: %v", e)
	}
	filter2, e = gcs.FromBytes(filter.N(), P, M, serialized2)
	if e != nil {
		t.Fatalf("Filter copy failed: %s", e.Error())
	}
	var serialized3 []byte
	serialized3, e = filter.NBytes()
	if e != nil {
		t.Fatalf("Filter NBytes() failed: %v", e)
	}
	filter3, e = gcs.FromNBytes(filter.P(), M, serialized3)
	if e != nil {
		t.Fatalf("Filter copy failed: %s", e.Error())
	}
}

// TestGCSFilterMetadata checks that the filter metadata is built and copied correctly.
func TestGCSFilterMetadata(t *testing.T) {
	if filter.P() != P {
		t.Fatal("P not correctly stored in filter metadata")
	}
	if filter.N() != uint32(len(contents)) {
		t.Fatal("N not correctly stored in filter metadata")
	}
	if filter.P() != filter2.P() {
		t.Fatal("P doesn't match between copied filters")
	}
	if filter.P() != filter3.P() {
		t.Fatal("P doesn't match between copied filters")
	}
	if filter.N() != filter2.N() {
		t.Fatal("N doesn't match between copied filters")
	}
	if filter.N() != filter3.N() {
		t.Fatal("N doesn't match between copied filters")
	}
	var serialized []byte
	serialized, e = filter.Bytes()
	if e != nil {
		t.Fatalf("Filter Hash() failed: %v", e)
	}
	var serialized2 []byte
	serialized2, e = filter2.Bytes()
	if e != nil {
		t.Fatalf("Filter Hash() failed: %v", e)
	}
	if !bytes.Equal(serialized, serialized2) {
		t.Fatal("Hash don't match between copied filters")
	}
	var serialized3 []byte
	serialized3, e = filter3.Bytes()
	if e != nil {
		t.Fatalf("Filter Hash() failed: %v", e)
	}
	if !bytes.Equal(serialized, serialized3) {
		t.Fatal("Hash don't match between copied filters")
	}
	var serialized4 []byte
	serialized4, e = filter3.Bytes()
	if e != nil {
		t.Fatalf("Filter Hash() failed: %v", e)
	}
	if !bytes.Equal(serialized, serialized4) {
		t.Fatal("Hash don't match between copied filters")
	}
}

// TestGCSFilterMatch checks that both the built and copied filters match correctly, logging any false positives without
// failing on them.
func TestGCSFilterMatch(t *testing.T) {
	match, e = filter.Match(key, []byte("Nate"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	match, e = filter2.Match(key, []byte("Nate"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	match, e = filter.Match(key, []byte("Quentin"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	match, e = filter2.Match(key, []byte("Quentin"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	match, e = filter.Match(key, []byte("Nates"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
	match, e = filter2.Match(key, []byte("Nates"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
	match, e = filter.Match(key, []byte("Quentins"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
	match, e = filter2.Match(key, []byte("Quentins"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
}

// TestGCSFilterMatchAny checks that both the built and copied filters match a list correctly, logging any false
// positives without failing on them.
func TestGCSFilterMatchAny(t *testing.T) {
	match, e = filter.MatchAny(key, contents2)
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
	match, e = filter2.MatchAny(key, contents2)
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e.Error())
	}
	if match {
		t.Logf("False positive match, should be 1 in 2**%d!", P)
	}
	contents2 = append(contents2, []byte("Nate"))
	match, e = filter.MatchAny(key, contents2)
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match any when it should have!")
	}
	match, e = filter2.MatchAny(key, contents2)
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e.Error())
	}
	if !match {
		t.Fatal("Filter didn't match any when it should have!")
	}
}
