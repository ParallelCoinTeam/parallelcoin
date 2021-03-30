package duration

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/opts/meta"
	"github.com/p9c/pod/pkg/opts/opt"
	uberatomic "go.uber.org/atomic"
	"strings"
	"time"
)

// Opt stores an time.Duration configuration value
type Opt struct {
	meta.Data
	hook  []func(d time.Duration)
	Value *uberatomic.Duration
	Def   time.Duration
}

// New creates a new Opt with a given default value set
func New(m meta.Data, def time.Duration) *Opt {
	return &Opt{Value: uberatomic.NewDuration(def), Data: m, Def: def}
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
	var v time.Duration
	if v, e = time.ParseDuration(input); !E.Chk(e) {
		x.Value.Store(v)
	}
	return
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
func (x *Opt) AddHooks(hook ...func(d time.Duration)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...func(d time.Duration)) {
	x.hook = hook
}

// V returns the value stored
func (x *Opt) V() time.Duration {
	return x.Value.Load()
}

// Set the value stored
func (x *Opt) Set(d time.Duration) *Opt {
	x.Value.Store(d)
	return x
}

// String returns a string representation of the value
func (x *Opt) String() string {
	return fmt.Sprintf("%s: %v", x.Data.Option, x.V())
}

// MarshalJSON returns the json representation
func (x *Opt) MarshalJSON() (b []byte, e error) {
	v := x.Value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	v := x.Value.Load()
	e = json.Unmarshal(data, &v)
	x.Value.Store(v)
	return
}
