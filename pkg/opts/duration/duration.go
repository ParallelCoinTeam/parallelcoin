package duration

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	uberatomic "go.uber.org/atomic"

	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/sanitizers"
)

// Opt stores an time.Duration configuration value
type Opt struct {
	meta.Data
	hook     []Hook
	clamp    func(input time.Duration) (result time.Duration)
	Min, Max time.Duration
	Value    *uberatomic.Duration
	Def      time.Duration
}

type Hook func(d time.Duration) error

// New creates a new Opt with a given default value set
func New(m meta.Data, def time.Duration, min, max time.Duration, hook ...Hook) *Opt {
	return &Opt{
		Value: uberatomic.NewDuration(def),
		Data:  m,
		Def:   def,
		Min:   min,
		Max:   max,
		hook:  hook,
		clamp: sanitizers.ClampDuration(min, max),
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
		e = fmt.Errorf("duration opt %s %v may not be empty", x.Name(), x.Data.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing '='
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var v time.Duration
	if v, e = time.ParseDuration(input); E.Chk(e) {
		return
	}
	if e = x.Set(v); E.Chk(e) {
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
func (x *Opt) AddHooks(hook ...Hook) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...Hook) {
	x.hook = hook
}

// V returns the value stored
func (x *Opt) V() time.Duration {
	return x.Value.Load()
}

func (x *Opt) runHooks(d time.Duration) (e error) {
	for i := range x.hook {
		if e = x.hook[i](d); E.Chk(e) {
			break
		}
	}
	return
}

// Set the value stored
func (x *Opt) Set(d time.Duration) (e error) {
	d = x.clamp(d)
	if e = x.runHooks(d); !E.Chk(e) {
		x.Value.Store(d)
	}
	return
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
	e = x.Set(v)
	return
}
