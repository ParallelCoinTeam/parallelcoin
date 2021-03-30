package podcfg

import (
	"encoding/json"
	"fmt"
	uberatomic "go.uber.org/atomic"
	"strings"
	"time"
)

type Duration struct {
	Metadata
	hook  []func(d time.Duration)
	value *uberatomic.Duration
	def   time.Duration
}

// NewDuration creates a new Duration with a given default value set
func NewDuration(m Metadata, def time.Duration) *Duration {
	return &Duration{value: uberatomic.NewDuration(def), Metadata: m, def: def}
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Duration) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Duration) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Duration) ReadInput(input string) (o Option, e error) {
	if input == "" {
		e = fmt.Errorf("integer number option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing characters
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var v time.Duration
	if v, e = time.ParseDuration(input); !E.Chk(e) {
		x.value.Store(v)
	}
	return
}

// LoadInput sets the value from a string (this is the same as the above but differs for Strings)
func (x *Duration) LoadInput(input string) (o Option, e error) {
	return x.ReadInput(input)
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
