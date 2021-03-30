package binary

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/opts"
	uberatomic "go.uber.org/atomic"
	"strings"
)

// Opt stores an boolean configuration value
type Opt struct {
	opts.Metadata
	hook  []func(b bool)
	value *uberatomic.Bool
	Def   bool
}

// New creates a new Opt with default values set
func New(m opts.Metadata, def bool, hook ...func(b bool)) *Opt {
	return &Opt{value: uberatomic.NewBool(def), Metadata: m, Def: def, hook: hook}
}

// SetName sets the name for the generator
func (x *Opt) SetName(name string) {
	x.Metadata.Option = strings.ToLower(name)
	x.Metadata.Name = name
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Opt) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Opt) GetMetadata() *opts.Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string.
// The value can be right up against the keyword or separated by a '='.
func (x *Opt) ReadInput(input string) (o opts.Option, e error) {
	// if the input is empty, the user intends the opposite of the default
	if input == "" {
		x.value.Store(!x.Def)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing characters
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	input = strings.ToLower(input)
	switch input {
	case "t", "true", "+":
		x.value.Store(true)
	case "f", "false", "-":
		x.value.Store(false)
	default:
		e = fmt.Errorf("input on option %s: '%s' is not valid for a boolean flag", x.Name(), input)
	}
	return
}

// LoadInput sets the value from a string (this is the same as the above but differs for Strings)
func (x *Opt) LoadInput(input string) (o opts.Option, e error) {
	return x.ReadInput(input)
}

// Name returns the name of the option
func (x *Opt) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Opt) AddHooks(hook ...func(b bool)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...func(b bool)) {
	x.hook = hook
}

// True returns whether the value is set to true (it returns the value)
func (x *Opt) True() bool {
	return x.value.Load()
}

// False returns whether the value is false (it returns the inverse of the value)
func (x *Opt) False() bool {
	return !x.value.Load()
}

// Flip changes the value to its opposite
func (x *Opt) Flip() {
	x.value.Toggle()
}

// Set changes the value currently stored
func (x *Opt) Set(b bool) *Opt {
	x.value.Store(b)
	return x
}

// String returns a string form of the value
func (x *Opt) String() string {
	return fmt.Sprint(x.Metadata.Option, ": ", x.True())
}

// T sets the value to true
func (x *Opt) T() *Opt {
	x.value.Store(true)
	return x
}

// F sets the value to false
func (x *Opt) F() *Opt {
	x.value.Store(false)
	return x
}

// MarshalJSON returns the json representation of a Opt
func (x *Opt) MarshalJSON() (b []byte, e error) {
	v := x.value.Load()
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation of a Opt
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load()
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}
