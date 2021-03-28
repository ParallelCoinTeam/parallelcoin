// Package podcfg implements a concurrent/parallel active application configuration system for multi-process
// applications to share a configuration as well as keep in sync with each other.
//
// This file contains all of the data types stored in a podcfg.Config and the various accessors and methods relevant to
// them. There is a basic byte slice-as-string type which is intended to eventually cover proper security practices for
// storing password information.
package podcfg

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	uberatomic "go.uber.org/atomic"
	"sync/atomic"
	"time"
)

type (
	Option interface {
		ReadInput(string) Option
		GetMetadata() *Metadata
		Name() string
	}
	Metadata struct {
		Option      string
		Aliases     []string
		Group       string
		Label       string
		Description string
		Type        string
		Widget      string
		Options     []string
		OmitEmpty   bool
	}
	Bool struct {
		Metadata
		hook  []func(b bool)
		value *uberatomic.Bool
		def   bool
	}
	Strings struct {
		Metadata
		hook  []func(s []string)
		value *atomic.Value
		def   []string
	}
	Float struct {
		Metadata
		hook  []func(f float64)
		value *uberatomic.Float64
		def   float64
	}
	Int struct {
		Metadata
		hook  []func(i int64)
		value *uberatomic.Int64
		def   int64
	}
	String struct {
		Metadata
		hook  []func(s Strice)
		value *atomic.Value
		def   string
	}
	Duration struct {
		Metadata
		hook  []func(d time.Duration)
		value *uberatomic.Duration
		def   time.Duration
	}
	// Strice is a wrapper around byte slices to enable optional security features and possibly better performance for
	// bulk comparison and editing. There isn't any extensive editing primitives for this purpose,
	Strice []byte
)

// S returns the underlying bytes converted into string
func (s *Strice) S() string {
	return string(*s)
}

// E returns the byte at the requested index in the string
func (s *Strice) E(elem int) byte {
	if s.Len() > elem {
		return (*s)[elem]
	}
	return 0
}

// Len returns the length of the string in bytes
func (s *Strice) Len() int {
	return len(*s)
}

// Equal returns true if two Strices are equal in both length and content
func (s *Strice) Equal(sb *Strice) bool {
	if s.Len() == sb.Len() {
		for i := range *s {
			if s.E(i) != sb.E(i) {
				return false
			}
		}
		return true
	}
	return false
}

// Cat two Strices together
func (s *Strice) Cat(sb *Strice) *Strice {
	*s = append(*s, *sb...)
	return s
}

// Find returns true if a match of a substring is found and if found, the position in the first string that the second
// string starts, the number of matching characters from the start of the search Strice, or -1 if not found.
//
// You specify a minimum length match and it will trawl through it systematically until it finds the first match of the
// minimum length.
func (s *Strice) Find(sb *Strice, minLengthMatch int) (found bool, extent, pos int) {
	// can't be a substring if it's longer
	if sb.Len() > s.Len() {
		return
	}
	for pos = range *s {
		// if we find a match, grab onto it
		if s.E(pos) == sb.E(pos) {
			extent++
			// this exhaustively searches for a match between the two strings, but we do not restrict the match to the
			// minimum, maximising the ways this function can be used for simple position tests and editing
			for srchPos := 1; srchPos < sb.Len() || srchPos+pos < s.Len(); srchPos++ {
				// the first element is skipped
				if s.E(srchPos+pos) != sb.E(srchPos) {
					break
				}
				extent++
			}
			// the above loop ends when the bytes stop matching, then if it is under the minimum length requested, it
			// continues. Note that we are not mutating `i` so it iterates for a match comprehensively.
			if extent < minLengthMatch {
				// reset the extent
				extent = 0
			} else {
				break
			}
		}
	}
	return
}

// HasPrefix returns true if the given string forms the beginning of the current string
func (s *Strice) HasPrefix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == 0 {
			return true
		}
	}
	return false
}

// HasSuffix returns true if the given string forms the ending of the current string
func (s *Strice) HasSuffix(sb *Strice) bool {
	found, _, pos := s.Find(sb, sb.Len())
	if found {
		if pos == s.Len()-sb.Len()-1 {
			return true
		}
	}
	return false
}

// Dup copies a string and returns it
func (s *Strice) Dup() *Strice {
	ns := make(Strice, s.Len())
	copy(ns, *s)
	return &ns
}

// Wipe zeroes the bytes of a string
func (s *Strice) Wipe() {
	for i := range *s {
		(*s)[i] = 0
	}
}

// Split the string by a given cutset
func (s *Strice) Split(cutset string) (out []*Strice) {
	// convert immutable string type to Strice bytes
	c := Strice(cutset)
	// need the pointer to call the methods
	cs := &c
	// copy the bytes so we can guarantee the original is unmodified
	cp := s.Dup()
	for {
		// locate the next instance of the cutset
		found, _, pos := s.Find(cp, cp.Len())
		if found {
			// add the found section to the return slice
			before := (*s)[:pos+cp.Len()]
			out = append(out, &before)
			// trim off the prefix and cutslice from the working copy
			*cs = (*cs)[pos+cp.Len():]
			// continue to search for more instances of the cutset
			continue
		} else {
			// once we get not found, the searching is over and whatever we have, we return
			break
		}
	}
	return
}

// NewBool creates a new podcfg.Bool with default values set
func NewBool(m Metadata, def bool, hook ...func(b bool)) *Bool {
	return &Bool{value: uberatomic.NewBool(def), Metadata: m, def: def, hook: hook}
}

// GetMetadata returns the metadata of the option type
func (x *Bool) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Bool) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *Bool) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Bool) AddHooks(hook ...func(b bool)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Bool) SetHooks(hook ...func(b bool)) {
	x.hook = hook
}

// True returns whether the value is set to true (it returns the value)
func (x *Bool) True() bool {
	return x.value.Load()
}

// False returns whether the value is false (it returns the inverse of the value)
func (x *Bool) False() bool {
	return !x.value.Load()
}

// Flip changes the value to its opposite
func (x *Bool) Flip() {
	x.value.Toggle()
}

// Set changes the value currently stored
func (x *Bool) Set(b bool) *Bool {
	x.value.Store(b)
	return x
}

// String returns a string form of the value
func (x *Bool) String() string {
	return fmt.Sprint(x.Metadata.Option, ": ", x.True())
}

// T sets the value to true
func (x *Bool) T() *Bool {
	x.value.Store(true)
	return x
}

// F sets the value to false
func (x *Bool) F() *Bool {
	x.value.Store(false)
	return x
}

// MarshalJSON returns the json representation of a Bool
func (x *Bool) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of a Bool
func (x *Bool) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewStrings  creates a new podcfg.Strings with default values set
func NewStrings(m Metadata, def []string, hook ...func(s []string)) *Strings {
	as := &atomic.Value{}
	v := cli.StringSlice(def)
	as.Store(&v)
	return &Strings{value: as, Metadata: m, def: def, hook: hook}
}

// GetMetadata returns the metadata of the option type
func (x *Strings) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Strings) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *Strings) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Strings) AddHooks(hook ...func(b []string)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Strings) SetHooks(hook ...func(b []string)) {
	x.hook = hook
}

// V returns the stored value
func (x *Strings) V() *cli.StringSlice {
	return x.value.Load().(*cli.StringSlice)
}

// Len returns the length of the slice of strings
func (x *Strings) Len() int {
	return len(x.S())
}

// Set the slice of strings stored
func (x *Strings) Set(ss []string) *Strings {
	sss := cli.StringSlice(ss)
	x.value.Store(&sss)
	return x
}

// S returns the value as a slice of string
func (x *Strings) S() []string {
	return *x.value.Load().(*cli.StringSlice)
}

// String returns a string representation of the value
func (x *Strings) String() string {
	return fmt.Sprint(x.Metadata.Option, ": ", x.S())
}

// MarshalJSON returns the json representation of
func (x *Strings) MarshalJSON() (b []byte, e error) {
	xs := x.value.Load().(*cli.StringSlice)
	return json.Marshal(xs)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Strings) UnmarshalJSON(data []byte) (e error) {
	v := &cli.StringSlice{}
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewFloat returns a new Float value set to a default value
func NewFloat(m Metadata, def float64) *Float {
	return &Float{value: uberatomic.NewFloat64(def), Metadata: m, def: def}
}

// GetMetadata returns the metadata of the option type
func (x *Float) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Float) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *Float) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Float) AddHooks(hook ...func(f float64)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Float) SetHooks(hook ...func(f float64)) {
	x.hook = hook
}

// V returns the value stored
func (x *Float) V() float64 {
	return x.value.Load()
}

// Set the value stored
func (x *Float) Set(f float64) *Float {
	x.value.Store(f)
	return x
}

// String returns a string representation of the value
func (x *Float) String() string {
	return fmt.Sprintf("%s: %0.8f", x.Metadata.Option, x.V())
}

// MarshalJSON returns the json representation of
func (x *Float) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Float) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewInt creates a new Int with a given default value
func NewInt(m Metadata, def int64) *Int {
	return &Int{value: uberatomic.NewInt64(def), Metadata: m, def: def}
}

// GetMetadata returns the metadata of the option type
func (x *Int) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Int) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *Int) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Int) AddHooks(hook ...func(f int64)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Int) SetHooks(hook ...func(f int64)) {
	x.hook = hook
}

// V returns the stored int
func (x *Int) V() int {
	return int(x.value.Load())
}

// Set the value stored
func (x *Int) Set(i int) *Int {
	x.value.Store(int64(i))
	return x
}

// String returns the string stored
func (x *Int) String() string {
	return fmt.Sprintf("%s: %d", x.Metadata.Option, x.V())
}

// MarshalJSON returns the json representation of
func (x *Int) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Int) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewString creates a new String with a given default value set
func NewString(m Metadata, def string) *String {
	v := &atomic.Value{}
	v.Store([]byte(def))
	return &String{value: v, Metadata: m, def: def}
}

// GetMetadata returns the metadata of the option type
func (x *String) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *String) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *String) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *String) AddHooks(hook ...func(f Strice)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *String) SetHooks(hook ...func(f Strice)) {
	x.hook = hook
}

// V returns the stored string
func (x *String) V() string {
	return string(x.value.Load().([]byte))
}

// Empty returns true if the string is empty
func (x *String) Empty() bool {
	return len(x.value.Load().([]byte)) == 0
}

// Bytes returns the raw bytes in the underlying storage
func (x *String) Bytes() []byte {
	return x.value.Load().([]byte)
}

// Set the value stored
func (x *String) Set(s string) *String {
	x.value.Store([]byte(s))
	return x
}

func (x *String) SetBytes(s []byte) *String {
	x.value.Store(s)
	return x
}

// String returns a string representation of the value
func (x *String) String() string {
	return fmt.Sprintf("%s: '%s'", x.Metadata.Option, x.V())
}

// MarshalJSON returns the json representation
func (x *String) MarshalJSON() (b []byte, e error) {
	v := string(x.value.Load().([]byte))
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *String) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load().([]byte)
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}

// NewDuration creates a new Duration with a given default value set
func NewDuration(m Metadata, def time.Duration) *Duration {
	return &Duration{value: uberatomic.NewDuration(def), Metadata: m, def: def}
}

// GetMetadata returns the metadata of the option type
func (x *Duration) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Duration) ReadInput(s string) Option {
	return x
}

// Name returns the name of the option
func (x *Duration) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Duration) AddHooks(hook ...func(d time.Duration)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Duration) SetHooks(hook ...func(d time.Duration)) {
	x.hook = hook
}

// V returns the value stored
func (x *Duration) V() time.Duration {
	return x.value.Load()
}

// Set the value stored
func (x *Duration) Set(d time.Duration) *Duration {
	x.value.Store(d)
	return x
}

// String returns a string representation of the value
func (x *Duration) String() string {
	return fmt.Sprintf("%s: %v", x.Metadata.Option, x.V())
}

// MarshalJSON returns the json representation
func (x *Duration) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *Duration) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}
