package integer

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/opts/meta"
	"github.com/p9c/pod/pkg/opts/opt"
	uberatomic "go.uber.org/atomic"
	"strconv"
	"strings"
)

// Opt stores an int configuration value
type Opt struct {
	meta.Data
	hook  []func(i int64)
	Value *uberatomic.Int64
	Def   int64
}

// New creates a new Opt with a given default value
func New(m meta.Data, def int64) *Opt {
	return &Opt{Value: uberatomic.NewInt64(def), Data: m, Def: def}
}

// SetName sets the name for the generator
func (x *Opt) SetName(name string) {
	x.Data.Option = strings.ToLower(name)
	x.Data.Name = name
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Opt) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the opt type
func (x *Opt) GetMetadata() *meta.Data {
	return &x.Data
}

// ReadInput sets the value from a string
func (x *Opt) ReadInput(input string) (o opt.Option, e error) {
	if input == "" {
		e = fmt.Errorf("integer number opt %s %v may not be empty", x.Name(), x.Data.Aliases)
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
	x.Value.Store(v)
	return x, e
}

// LoadInput sets the value from a string (this is the same as the above but differs for Strings)
func (x *Opt) LoadInput(input string) (o opt.Option, e error) {
	return x.ReadInput(input)
}

// Name returns the name of the opt
func (x *Opt) Name() string {
	return x.Data.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Opt) AddHooks(hook ...func(f int64)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...func(f int64)) {
	x.hook = hook
}

// V returns the stored int
func (x *Opt) V() int {
	return int(x.Value.Load())
}

// Set the value stored
func (x *Opt) Set(i int) *Opt {
	x.Value.Store(int64(i))
	return x
}

// String returns the string stored
func (x *Opt) String() string {
	return fmt.Sprintf("%s: %d", x.Data.Option, x.V())
}

// MarshalJSON returns the json representation of
func (x *Opt) MarshalJSON() (b []byte, e error) {
	v := x.Value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	v := x.Value.Load()
	e = json.Unmarshal(data, &v)
	x.Value.Store(v)
	return
}
