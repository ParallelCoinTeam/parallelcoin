package builder_test

import (
	"encoding/hex"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/gcs"
	"github.com/p9c/matrjoska/pkg/gcs/builder"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/wire"
)

var (
	// // No need to allocate an err variable in every test
	// e error
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
	testKey = [16]byte{
		0x4c, 0xb1, 0xab, 0x12, 0x57, 0x62, 0x1e, 0x41,
		0x3b, 0x8b, 0x0e, 0x26, 0x64, 0x8d, 0x4a, 0x15,
	}
	testHash = "000000000000000000496d7ff9bd2c96154a8d64260e8b3b411e625712abb14c"
	testAddr = "3Nxwenay9Z8Lc9JBiywExpnEFiLp6Afp8v"
	witness  = [][]byte{
		{
			0x4c, 0xb1, 0xab, 0x12, 0x57, 0x62, 0x1e, 0x41,
			0x3b, 0x8b, 0x0e, 0x26, 0x64, 0x8d, 0x4a, 0x15,
			0x3b, 0x8b, 0x0e, 0x26, 0x64, 0x8d, 0x4a, 0x15,
			0x3b, 0x8b, 0x0e, 0x26, 0x64, 0x8d, 0x4a, 0x15,
		},
		{
			0xdd, 0xa3, 0x5a, 0x14, 0x88, 0xfb, 0x97, 0xb6,
			0xeb, 0x3f, 0xe6, 0xe9, 0xef, 0x2a, 0x25, 0x81,
			0x4e, 0x39, 0x6f, 0xb5, 0xdc, 0x29, 0x5f, 0xe9,
			0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98,
			0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98,
			0x94, 0xb9, 0x67, 0x89, 0xb2, 0x1a, 0x03, 0x98,
		},
	}
)

// TestUseBlockHash tests using a block hash as a filter key.
func TestUseBlockHash(t *testing.T) {
	// Block hash #448710, pretty high difficulty.
	hash, e := chainhash.NewHashFromStr(testHash)
	if e != nil {
		t.Fatalf("Hash from string failed: %s", e.Error())
	}
	// wire.OutPoint
	outPoint := wire.OutPoint{
		Hash:  *hash,
		Index: 4321,
	}
	// util.Address
	addr, e := btcaddr.Decode(testAddr, &chaincfg.MainNetParams)
	if e != nil {
		t.Fatalf("Address decode failed: %s", e.Error())
	}
	addrBytes, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("Address script podbuild failed: %s", e.Error())
	}
	// Create a GCS with a key hash and check that the key is derived correctly, then test it.
	b := builder.WithKeyHash(hash)
	key, e := b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with key hash failed: %s",
			e.Error(),
		)
	}
	if key != testKey {
		t.Fatalf(
			"Key not derived correctly from key hash:\n%s\n%s",
			hex.EncodeToString(key[:]),
			hex.EncodeToString(testKey[:]),
		)
	}
	BuilderTest(b, hash, builder.DefaultP, outPoint, addrBytes, witness, t)
	// Create a GCS with a key hash and non-default P and test it.
	b = builder.WithKeyHashPM(hash, 30, 90)
	BuilderTest(b, hash, 30, outPoint, addrBytes, witness, t)
	// Create a GCS with a random key, set the key from a hash manually, check that the key is correct, and test it.
	b = builder.WithRandomKey()
	b.SetKeyFromHash(hash)
	key, e = b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with known key failed: %s",
			e.Error(),
		)
	}
	if key != testKey {
		t.Fatalf(
			"Key not copied correctly from known key:\n%s\n%s",
			hex.EncodeToString(key[:]),
			hex.EncodeToString(testKey[:]),
		)
	}
	BuilderTest(b, hash, builder.DefaultP, outPoint, addrBytes, witness, t)
	// Create a GCS with a random key and test it.
	b = builder.WithRandomKey()
	key1, e := b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with random key failed: %s",
			e.Error(),
		)
	}
	t.Logf("Random Key 1: %s", hex.EncodeToString(key1[:]))
	BuilderTest(b, hash, builder.DefaultP, outPoint, addrBytes, witness, t)
	// Create a GCS with a random key and non-default P and test it.
	b = builder.WithRandomKeyPM(30, 90)
	key2, e := b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with random key failed: %s",
			e.Error(),
		)
	}
	t.Logf("Random Key 2: %s", hex.EncodeToString(key2[:]))
	if key2 == key1 {
		t.Fatalf("Random keys are the same!")
	}
	BuilderTest(b, hash, 30, outPoint, addrBytes, witness, t)
	// Create a GCS with a known key and test it.
	b = builder.WithKey(testKey)
	key, e = b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with known key failed: %s",
			e.Error(),
		)
	}
	if key != testKey {
		t.Fatalf(
			"Key not copied correctly from known key:\n%s\n%s",
			hex.EncodeToString(key[:]),
			hex.EncodeToString(testKey[:]),
		)
	}
	BuilderTest(b, hash, builder.DefaultP, outPoint, addrBytes, witness, t)
	// Create a GCS with a known key and non-default P and test it.
	b = builder.WithKeyPM(testKey, 30, 90)
	key, e = b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with known key failed: %s",
			e.Error(),
		)
	}
	if key != testKey {
		t.Fatalf(
			"Key not copied correctly from known key:\n%s\n%s",
			hex.EncodeToString(key[:]),
			hex.EncodeToString(testKey[:]),
		)
	}
	BuilderTest(b, hash, 30, outPoint, addrBytes, witness, t)
	// Create a GCS with a known key and too-high P and ensure error works throughout all functions that use it.
	b = builder.WithRandomKeyPM(33, 99).SetKeyFromHash(hash).SetKey(testKey)
	b.SetP(30).AddEntry(hash.CloneBytes()).AddEntries(contents).
		AddHash(hash).AddEntry(addrBytes)
	_, e = b.Key()
	if e != gcs.ErrPTooBig {
		t.Fatalf("No error on P too big!")
	}
	_, e = b.Build()
	if e != gcs.ErrPTooBig {
		t.Fatalf("No error on P too big!")
	}
}
func BuilderTest(
	b *builder.GCS, hash *chainhash.Hash, p uint8,
	outPoint wire.OutPoint, addrBytes []byte, witness wire.TxWitness,
	t *testing.T,
) {
	key, e := b.Key()
	if e != nil {
		t.Fatalf(
			"Builder instantiation with key hash failed: %s",
			e.Error(),
		)
	}
	// Build a filter and test matches.
	b.AddEntries(contents)
	f, e := b.Build()
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
	if f.P() != p {
		t.Fatalf("Filter built with wrong probability")
	}
	match, e := f.Match(key, []byte("Nate"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e)
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	match, e = f.Match(key, []byte("weks"))
	if e != nil {
		t.Fatalf("Filter match failed: %s", e)
	}
	if match {
		t.Logf(
			"False positive match, should be 1 in 2**%d!",
			builder.DefaultP,
		)
	}
	// Add a hash, podbuild a filter, and test matches
	b.AddHash(hash)
	f, e = b.Build()
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
	match, e = f.Match(key, hash.CloneBytes())
	if e != nil {
		t.Fatalf("Filter match failed: %s", e)
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	// Add a script, podbuild a filter, and test matches
	b.AddEntry(addrBytes)
	f, e = b.Build()
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
	match, e = f.MatchAny(key, [][]byte{addrBytes})
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e)
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	// Add a routine witness stack, podbuild a filter, and test that it matches.
	b.AddWitness(witness)
	f, e = b.Build()
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
	match, e = f.MatchAny(key, witness)
	if e != nil {
		t.Fatalf("Filter match any failed: %s", e)
	}
	if !match {
		t.Fatal("Filter didn't match when it should have!")
	}
	// Chk that adding duplicate items does not increase filter size.
	originalSize := f.N()
	b.AddEntry(addrBytes)
	b.AddWitness(witness)
	f, e = b.Build()
	if e != nil {
		t.Fatalf("Filter podbuild failed: %s", e.Error())
	}
	if f.N() != originalSize {
		t.Fatal("Filter size increased with duplicate items")
	}
}
