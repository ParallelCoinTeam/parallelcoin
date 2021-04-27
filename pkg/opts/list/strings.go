package list

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/p9c/opts/normalize"

	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/opt"
	"github.com/p9c/opts/sanitizers"
)

// Opt stores a string slice configuration value
type Opt struct {
	meta.Data
	hook  []Hook
	Value *atomic.Value
	Def   []string
}

type Hook func(s []string) error

// New  creates a new Opt with default values set
func New(m meta.Data, def []string, hook ...Hook) *Opt {
	as := &atomic.Value{}
	as.Store(def)
	return &Opt{Value: as, Data: m, Def: def, hook: hook}
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

// ReadInput adds the value from a string. For this opt this means appending to the list
func (x *Opt) ReadInput(input string) (o opt.Option, e error) {
	if input == "" {
		e = fmt.Errorf("string opt %s %v may not be empty", x.Name(), x.Data.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	// if value has a comma in it, it's a list of items, so split them and append them
	slice := x.S()
	if strings.Contains(input, ",") {
		split := strings.Split(input, ",")
		for i := range split {
			var cleaned string
			if cleaned, e = sanitizers.StringType(x.Data.Type, split[i], x.Data.DefaultPort); E.Chk(e) {
				return
			}
			if cleaned != "" {
				I.Ln("setting value for", x.Data.Name, cleaned)
				split[i] = cleaned
			}
		}
		e = x.Set(append(slice, split...))
	} else {
		var cleaned string
		if cleaned, e = sanitizers.StringType(x.Data.Type, input, x.Data.DefaultPort); E.Chk(e) {
			return
		}
		if cleaned != "" {
			I.Ln("setting value for", x.Data.Name, cleaned)
			input = cleaned
		}
		if e = x.Set(append(slice, input)); E.Chk(e) {
		}

	}
	// ensure there is no duplicates
	e = x.Set(normalize.RemoveDuplicateAddresses(x.V()))
	return x, e
}

// LoadInput sets the value from a string. For this opt this replacing the list
func (x *Opt) LoadInput(input string) (o opt.Option, e error) {
	old := x.V()
	_ = x.Set([]string{})
	if o, e = x.ReadInput(input); E.Chk(e) {
		// if input failed to parse, restore its prior state
		_ = x.Set(old)
	}
	return
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

// V returns the stored value
func (x *Opt) V() []string {
	return x.Value.Load().([]string)
}

// Len returns the length of the slice of strings
func (x *Opt) Len() int {
	return len(x.S())
}

func (x *Opt) runHooks(s []string) (e error) {
	for i := range x.hook {
		if e = x.hook[i](s); E.Chk(e) {
			break
		}
	}
	return
}

// Set the slice of strings stored
func (x *Opt) Set(ss []string) (e error) {
	if e = x.runHooks(ss); !E.Chk(e) {
		x.Value.Store(ss)
	}
	return
}

// S returns the value as a slice of string
func (x *Opt) S() []string {
	return x.Value.Load().([]string)
}

// String returns a string representation of the value
func (x *Opt) String() string {
	return fmt.Sprint(x.Data.Option, ": ", x.S())
}

// MarshalJSON returns the json representation of
func (x *Opt) MarshalJSON() (b []byte, e error) {
	xs := x.Value.Load().([]string)
	return json.Marshal(xs)
}

// UnmarshalJSON decodes a JSON representation of
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	var v []string
	e = json.Unmarshal(data, &v)
	x.Value.Store(v)
	return
}
