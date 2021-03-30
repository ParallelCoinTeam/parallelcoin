package podcfg

import (
	"encoding/json"
	"fmt"
	uberatomic "go.uber.org/atomic"
	"strconv"
	"strings"
)

type Float struct {
	Metadata
	hook  []func(f float64)
	value *uberatomic.Float64
	def   float64
}

// NewFloat returns a new Float value set to a default value
func NewFloat(m Metadata, def float64) *Float {
	return &Float{value: uberatomic.NewFloat64(def), Metadata: m, def: def}
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Float) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Float) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Float) ReadInput(input string) (o Option, e error) {
	if input == "" {
		e = fmt.Errorf("floating point number option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing characters
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var v float64
	if v, e = strconv.ParseFloat(input, 64); E.Chk(e) {
		return
	}
	x.value.Store(v)
	return x, e
}

// LoadInput sets the value from a string (this is the same as the above but differs for Strings)
func (x *Float) LoadInput(input string) (o Option, e error) {
	return x.ReadInput(input)
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
