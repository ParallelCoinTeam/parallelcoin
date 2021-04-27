package gcs_test

import (
	"encoding/binary"
	"math/rand"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/gcs"
)

func genRandFilterElements(numElements uint) ([][]byte, error) {
	testContents := make([][]byte, numElements)
	for i := range contents {
		randElem := make([]byte, 32)
		if _, e = rand.Read(randElem); E.Chk(e) {
			return nil, e
		}
		testContents[i] = randElem
	}
	return testContents, nil
}

var (
	generatedFilter *gcs.Filter
	// filterErr       error
)

// BenchmarkGCSFilterBuild benchmarks building a filter.
func BenchmarkGCSFilterBuild50000(b *testing.B) {
	b.StopTimer()
	var testKey [gcs.KeySize]byte
	for i := 0; i < gcs.KeySize; i += 4 {
		binary.BigEndian.PutUint32(testKey[i:], rand.Uint32())
	}
	randFilterElems, genErr := genRandFilterElements(50000)
	if e != nil {
		b.Fatalf("unable to generate random item: %v", genErr)
	}
	b.StartTimer()
	var localFilter *gcs.Filter
	for i := 0; i < b.N; i++ {
		localFilter, e = gcs.BuildGCSFilter(
			P, M, key, randFilterElems,
		)
		if e != nil {
			b.Fatalf("unable to generate filter: %v", e)
		}
	}
	generatedFilter = localFilter
}

// BenchmarkGCSFilterBuild benchmarks building a filter.
func BenchmarkGCSFilterBuild100000(b *testing.B) {
	b.StopTimer()
	var testKey [gcs.KeySize]byte
	for i := 0; i < gcs.KeySize; i += 4 {
		binary.BigEndian.PutUint32(testKey[i:], rand.Uint32())
	}
	randFilterElems, genErr := genRandFilterElements(100000)
	if e != nil {
		b.Fatalf("unable to generate random item: %v", genErr)
	}
	b.StartTimer()
	var localFilter *gcs.Filter
	for i := 0; i < b.N; i++ {
		localFilter, e = gcs.BuildGCSFilter(
			P, M, key, randFilterElems,
		)
		if e != nil {
			b.Fatalf("unable to generate filter: %v", e)
		}
	}
	generatedFilter = localFilter
}

var (
	match bool
)

// BenchmarkGCSFilterMatch benchmarks querying a filter for a single value.
func BenchmarkGCSFilterMatch(b *testing.B) {
	b.StopTimer()
	filter, e := gcs.BuildGCSFilter(P, M, key, contents)
	if e != nil {
		b.Fatalf("Failed to podbuild filter")
	}
	b.StartTimer()
	var (
		localMatch bool
	)
	for i := 0; i < b.N; i++ {
		_, e = filter.Match(key, []byte("Nate"))
		if e != nil {
			b.Fatalf("unable to match filter: %v", e)
		}
		localMatch, e = filter.Match(key, []byte("Nates"))
		if e != nil {
			b.Fatalf("unable to match filter: %v", e)
		}
	}
	match = localMatch
}

// BenchmarkGCSFilterMatchAny benchmarks querying a filter for a list of values.
func BenchmarkGCSFilterMatchAny(b *testing.B) {
	b.StopTimer()
	filter, e := gcs.BuildGCSFilter(P, M, key, contents)
	if e != nil {
		b.Fatalf("Failed to podbuild filter")
	}
	b.StartTimer()
	var (
		localMatch bool
	)
	for i := 0; i < b.N; i++ {
		localMatch, e = filter.MatchAny(key, contents2)
		if e != nil {
			b.Fatalf("unable to match filter: %v", e)
		}
	}
	match = localMatch
}
