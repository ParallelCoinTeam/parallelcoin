package podcfg

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"strings"
	"sync/atomic"
)

type Strings struct {
	Metadata
	hook  []func(s []string)
	value *atomic.Value
	def   []string
}

// NewStrings  creates a new podcfg.Strings with default values set
func NewStrings(m Metadata, def []string, hook ...func(s []string)) *Strings {
	as := &atomic.Value{}
	v := cli.StringSlice(def)
	as.Store(&v)
	return &Strings{value: as, Metadata: m, def: def, hook: hook}
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Strings) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *Strings) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string. For this option this means appending to the list
func (x *Strings) ReadInput(input string) (o Option, e error) {
	if input == "" {
		e = fmt.Errorf("string option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
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

// LoadInput sets the value from a string. For this option this means appending to the list
func (x *Strings) LoadInput(input string) (o Option, e error) {
	if input == "" {
		e = fmt.Errorf("string option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
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

// Name returns the name of the option
func (x *Strings) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Strings) AddHooks(hook ...func(b []string)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Strings) SetHooks(hook ...func(b []string)) {
	x.hook = hook
}

// V returns the stored value
func (x *Strings) V() *cli.StringSlice {
	return x.value.Load().(*cli.StringSlice)
}

// Len returns the length of the slice of strings
func (x *Strings) Len() int {
	return len(x.S())
}

// Set the slice of strings stored
func (x *Strings) Set(ss []string) *Strings {
	sss := cli.StringSlice(ss)
	x.value.Store(&sss)
	return x
}

// S returns the value as a slice of string
func (x *Strings) S() []string {
	return *x.value.Load().(*cli.StringSlice)
}

// String returns a string representation of the value
func (x *Strings) String() string {
	return fmt.Sprint(x.Metadata.Option, ": ", x.S())
}

// MarshalJSON returns the json representation of
func (x *Strings) MarshalJSON() (b []byte, e error) {
	xs := x.value.Load().(*cli.StringSlice)
	return json.Marshal(xs)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Strings) UnmarshalJSON(data []byte) (e error) {
	v := &cli.StringSlice{}
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}
