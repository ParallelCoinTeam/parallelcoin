package podcfg

import (
	"encoding/json"
	"fmt"
	uberatomic "go.uber.org/atomic"
	"strings"
)

type Bool struct {
	Metadata
	hook  []func(b bool)
	value *uberatomic.Bool
	def   bool
}

// NewBool creates a new podcfg.Bool with default values set
func NewBool(m Metadata, def bool, hook ...func(b bool)) *Bool {
	return &Bool{value: uberatomic.NewBool(def), Metadata: m, def: def, hook: hook}
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Bool) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Bool) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *Bool) ReadInput(s string) (opt Option, e error) {
	// if the input is empty, the user intends the default
	if s == "" {
		x.value.Store(x.def)
		return
	}
	s = strings.ToLower(s)
	switch s {
	case "t", "true":
		x.value.Store(true)
	case "f", "false":
		x.value.Store(false)
	default:
		e = fmt.Errorf("input on option %s: '%s' is not valid for a boolean flag", x.Name(), s)
	}
	return
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
