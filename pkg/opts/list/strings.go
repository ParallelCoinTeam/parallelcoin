package list

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/opts/meta"
	"github.com/p9c/pod/pkg/opts/opt"
	"strings"
	"sync/atomic"
)

// Opt stores a string slice configuration value
type Opt struct {
	meta.Data
	hook  []func(s []string)
	value *atomic.Value
	Def   []string
}

// New  creates a new Opt with default values set
func New(m meta.Data, def []string, hook ...func(s []string)) *Opt {
	as := &atomic.Value{}
	as.Store(def)
	return &Opt{value: as, Data: m, Def: def, hook: hook}
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

// ReadInput sets the value from a string. For this opt this means appending to the list
func (x *Opt) ReadInput(input string) (o opt.Option, e error) {
	if input == "" {
		e = fmt.Errorf("string opt %s %v may not be empty", x.Name(), x.Data.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	// if value has a comma in it, it's a list of items, so split them and join them
	slice := x.S()
	if strings.Contains(input, ",") {
		x.Set(append(slice, strings.Split(input, ",")...))
	} else {
		x.Set(append(slice, input))
	}
	return x, e
}

// LoadInput sets the value from a string. For this opt this means appending to the list
func (x *Opt) LoadInput(input string) (o opt.Option, e error) {
	if input == "" {
		e = fmt.Errorf("string opt %s %v may not be empty", x.Name(), x.Data.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	var slice []string
	// if value has a comma in it, it's a list of items, so split them and join them
	if strings.Contains(input, ",") {
		x.Set(append(slice, strings.Split(input, ",")...))
	} else {
		x.Set(append(slice, input))
	}
	return x, e
}

// Name returns the name of the opt
func (x *Opt) Name() string {
	return x.Data.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Opt) AddHooks(hook ...func(b []string)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...func(b []string)) {
	x.hook = hook
}

// V returns the stored value
func (x *Opt) V() []string {
	return x.value.Load().([]string)
}

// Len returns the length of the slice of strings
func (x *Opt) Len() int {
	return len(x.S())
}

// Set the slice of strings stored
func (x *Opt) Set(ss []string) *Opt {
	x.value.Store(ss)
	return x
}

// S returns the value as a slice of string
func (x *Opt) S() []string {
	return x.value.Load().([]string)
}

// String returns a string representation of the value
func (x *Opt) String() string {
	return fmt.Sprint(x.Data.Option, ": ", x.S())
}

// MarshalJSON returns the json representation of
func (x *Opt) MarshalJSON() (b []byte, e error) {
	xs := x.value.Load().([]string)
	return json.Marshal(xs)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	var v []string
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}
