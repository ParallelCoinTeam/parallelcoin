package atom

import (
	"time"
	
	"go.uber.org/atomic"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

// import all the atomics from uber atomic
type (
	Int32    struct{ atomic.Int32 }
	Int64    struct{ atomic.Int64 }
	Uint32   struct{ atomic.Uint32 }
	Uint64   struct{ atomic.Uint64 }
	Bool     struct{ atomic.Bool }
	Float64  struct{ atomic.Float64 }
	Duration struct{ atomic.Duration }
	Value    struct{ atomic.Value }
)

// The following are types added for handling cryptocurrency data for
// ParallelCoin

// Time is an atomic wrapper around time.Time
// https://godoc.org/time#Time
type Time struct {
	v *Int64
}

// NewTime creates a Time.
func NewTime(tt time.Time) *Time {
	t := &Int64{}
	t.Store(tt.UnixNano())
	return &Time{v: t}
}

// Load atomically loads the wrapped value.
func (at *Time) Load() time.Time {
	return time.Unix(0, at.v.Load())
}

// Store atomically stores the passed value.
func (at *Time) Store(n time.Time) {
	at.v.Store(n.UnixNano())
}

// Add atomically adds to the wrapped time.Duration and returns the new value.
func (at *Time) Add(n time.Time) time.Time {
	return time.Unix(0, at.v.Add(n.UnixNano()))
}

// Sub atomically subtracts from the wrapped time.Duration and returns the new value.
func (at *Time) Sub(n time.Time) time.Time {
	return time.Unix(0, at.v.Sub(n.UnixNano()))
}

// Swap atomically swaps the wrapped time.Duration and returns the old value.
func (at *Time) Swap(n time.Time) time.Time {
	return time.Unix(0, at.v.Swap(n.UnixNano()))
}

// CAS is an atomic compare-and-swap.
func (at *Time) CAS(old, new time.Time) bool {
	return at.v.CAS(old.UnixNano(), new.UnixNano())
}

// Hash is an atomic wrapper around chainhash.Hash
// Note that there isn't really any reason to have CAS or arithmetic or
// comparisons as it is fine to do these non-atomically between Load/Store and
// they are (slightly) long operations)
type Hash struct {
	v *Value
}

// NewHash creates a Hash.
func NewHash(tt chainhash.Hash) *Hash {
	t := &Value{}
	t.Store(tt)
	return &Hash{v: t}
}

// Load atomically loads the wrapped value.
func (at *Hash) Load() chainhash.Hash {
	return at.v.Load().(chainhash.Hash)
}

// Store atomically stores the passed value.
func (at *Hash) Store(h chainhash.Hash) {
	at.v.Store(h)
}

// Swap atomically swaps the wrapped chainhash.Hash and returns the old value.
func (at *Hash) Swap(n chainhash.Hash) chainhash.Hash {
	o := at.v.Load().(chainhash.Hash)
	at.v.Store(n)
	return o
}
