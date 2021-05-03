package integer

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	uberatomic "go.uber.org/atomic"

	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/sanitizers"
)

// Opt stores an int configuration value
type Opt struct {
	meta.Data
	hook     []Hook
	Min, Max int
	clamp    func(input int) (result int)
	Value    *uberatomic.Int64
	Def      int64
}

type Hook func(i int) error

// New creates a new Opt with a given default value
func New(m meta.Data, def int64, min, max int, hook ...Hook) *Opt {
	return &Opt{
		Value: uberatomic.NewInt64(def),
		Data:  m,
		Def:   def,
		Min:   min,
		Max:   max,
		hook:  hook,
		clamp: sanitizers.ClampInt(min, max),
	}
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
		// the following removes leading and trailing '='
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var v int64
	if v, e = strconv.ParseInt(input, 10, 64); E.Chk(e) {
		return
	}
	if e = x.Set(int(v)); E.Chk(e) {
	}
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
func (x *Opt) AddHooks(hook ...Hook) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...Hook) {
	x.hook = hook
}

// V returns the stored int
func (x *Opt) V() int {
	return int(x.Value.Load())
}

func (x *Opt) runHooks(ii int) (e error) {
	for i := range x.hook {
		if e = x.hook[i](ii); E.Chk(e) {
			break
		}
	}
	return
}

// Set the value stored
func (x *Opt) Set(i int) (e error) {
	i = x.clamp(i)
	if e = x.runHooks(i); !E.Chk(e) {
		x.Value.Store(int64(i))
	}
	return
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
	e = x.Set(int(v))
	return
}
