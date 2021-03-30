package podcfg

import (
	"encoding/json"
	"fmt"
	uberatomic "go.uber.org/atomic"
	"strconv"
	"strings"
)

type Int struct {
	Metadata
	hook  []func(i int64)
	value *uberatomic.Int64
	def   int64
}

// NewInt creates a new Int with a given default value
func NewInt(m Metadata, def int64) *Int {
	return &Int{value: uberatomic.NewInt64(def), Metadata: m, def: def}
}

// SetName sets the name for the generator
func (x *Int) SetName(name string) {
	x.Metadata.Option = strings.ToLower(name)
	x.Metadata.Name = name
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Int) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Int) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Int) ReadInput(input string) (o Option, e error) {
	if input == "" {
		e = fmt.Errorf("integer number option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing characters
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var v int64
	if v, e = strconv.ParseInt(input, 10, 64); E.Chk(e) {
		return
	}
	x.value.Store(v)
	return x, e
}

// LoadInput sets the value from a string (this is the same as the above but differs for Strings)
func (x *Int) LoadInput(input string) (o Option, e error) {
	return x.ReadInput(input)
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
