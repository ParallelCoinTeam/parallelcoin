package atom

import (
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"time"
	
	"go.uber.org/atomic"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// import all the atomics from uber atomic
type (
	Int32    struct{ *atomic.Int32 }
	Int64    struct{ *atomic.Int64 }
	Uint32   struct{ *atomic.Uint32 }
	Uint64   struct{ *atomic.Uint64 }
	Bool     struct{ *atomic.Bool }
	Float64  struct{ *atomic.Float64 }
	Duration struct{ *atomic.Duration }
	Value    struct{ *atomic.Value }
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
	t := &Int64{atomic.NewInt64(tt.UnixNano())}
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
// comparisons as it is fine to do these non-atomically between Load / Store and
// they are (slightly) long operations)
type Hash struct {
	*Value
}

// NewHash creates a Hash.
func NewHash(tt chainhash.Hash) *Hash {
	t := &Value{
		Value: &atomic.Value{},
	}
	t.Store(tt)
	return &Hash{Value: t}
}

// Load atomically loads the wrapped value.
// The returned value copied so as to prevent mutation by concurrent users
// of the atomic, as arrays, slices and maps are pass-by-reference variables
func (at *Hash) Load() chainhash.Hash {
	o := at.Value.Load().(chainhash.Hash)
	var v chainhash.Hash
	copy(v[:], o[:])
	return v
}

// Store atomically stores the passed value.
// The passed value is copied so further mutations are not propagated.
func (at *Hash) Store(h chainhash.Hash) {
	var v chainhash.Hash
	copy(v[:], h[:])
	at.Value.Store(v)
}

// Swap atomically swaps the wrapped chainhash.Hash and returns the old value.
func (at *Hash) Swap(n chainhash.Hash) chainhash.Hash {
	o := at.Value.Load().(chainhash.Hash)
	at.Value.Store(n)
	return o
}

// Address is an atomic wrapper around util.Address
type Address struct {
	*atomic.String
	ForNet *chaincfg.Params
}

// NewAddress creates a Hash.
func NewAddress(tt btcaddr.Address, forNet *chaincfg.Params) *Address {
	t := atomic.NewString(tt.EncodeAddress())
	return &Address{String: t, ForNet: forNet}
}

// Load atomically loads the wrapped value.
func (at *Address) Load() btcaddr.Address {
	addr, e := btcaddr.Decode(at.String.Load(), at.ForNet)
	if e != nil {
		return nil
	}
	return addr
}

// Store atomically stores the passed value.
// The passed value is copied so further mutations are not propagated.
func (at *Address) Store(h btcaddr.Address) {
	at.String.Store(h.EncodeAddress())
}

// Swap atomically swaps the wrapped util.Address and returns the old value.
func (at *Address) Swap(n btcaddr.Address) btcaddr.Address {
	o := at.Load()
	at.Store(n)
	return o
}

// ListTransactionsResult is an atomic wrapper around
// []btcjson.ListTransactionsResult
type ListTransactionsResult struct {
	v *Value
}

// NewListTransactionsResult creates a btcjson.ListTransactionsResult.
func NewListTransactionsResult(ltr []btcjson.ListTransactionsResult) *ListTransactionsResult {
	t := &Value{
		Value: &atomic.Value{},
	}
	v := make([]btcjson.ListTransactionsResult, len(ltr))
	copy(v, ltr)
	t.Store(v)
	return &ListTransactionsResult{v: t}
}

// Load atomically loads the wrapped value.
// Note that it is copied and the stored value remains as it is
func (at *ListTransactionsResult) Load() []btcjson.ListTransactionsResult {
	ltr := at.v.Load().([]btcjson.ListTransactionsResult)
	v := make([]btcjson.ListTransactionsResult, len(ltr))
	copy(v, ltr)
	return v
}

// Store atomically stores the passed value.
// Note that it is copied and the passed value remains as it is
func (at *ListTransactionsResult) Store(ltr []btcjson.ListTransactionsResult) {
	v := make([]btcjson.ListTransactionsResult, len(ltr))
	copy(v, ltr)
	at.v.Store(v)
}

// Swap atomically swaps the wrapped chainhash.ListTransactionsResult and
// returns the old value.
func (at *ListTransactionsResult) Swap(
	n []btcjson.ListTransactionsResult,
) []btcjson.ListTransactionsResult {
	o := at.v.Load().([]btcjson.ListTransactionsResult)
	at.v.Store(n)
	return o
}

// Len returns the length of the []btcjson.ListTransactionsResult
func (at *ListTransactionsResult) Len() int {
	return len(at.v.Load().([]btcjson.ListTransactionsResult))
}
