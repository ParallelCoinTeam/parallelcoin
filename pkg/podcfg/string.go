package podcfg

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
)

// String stores a string configuration value as bytes
type String struct {
	Metadata
	hook  []func(s Strice)
	value *atomic.Value
	def   string
}

// NewString creates a new String with a given default value set
func NewString(m Metadata, def string) *String {
	v := &atomic.Value{}
	v.Store([]byte(def))
	return &String{value: v, Metadata: m, def: def}
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *String) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the option type
func (x *String) GetMetadata() *Metadata {
	return &x.Metadata
}

// ReadInput sets the value from a string
func (x *String) ReadInput(s string) (o Option, e error) {
	if s == "" {
		e = fmt.Errorf("string option %s %v may not be empty", x.Name(), x.Metadata.Aliases)
		return
	}
	if strings.HasPrefix(s, "=") {
		// the following removes leading and trailing characters
		s = strings.Join(strings.Split(s, "=")[1:], "=")
	}
	x.Set(s)
	return x, e
}

// Name returns the name of the option
func (x *String) Name() string {
	return x.Metadata.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *String) AddHooks(hook ...func(f Strice)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *String) SetHooks(hook ...func(f Strice)) {
	x.hook = hook
}

// V returns the stored string
func (x *String) V() string {
	return string(x.value.Load().([]byte))
}

// Empty returns true if the string is empty
func (x *String) Empty() bool {
	return len(x.value.Load().(Strice)) == 0
}

// Bytes returns the raw bytes in the underlying storage
func (x *String) Bytes() []byte {
	return x.value.Load().([]byte)
}

// Set the value stored
func (x *String) Set(s string) *String {
	x.value.Store([]byte(s))
	return x
}

// SetBytes sets the string from bytes
func (x *String) SetBytes(s []byte) *String {
	x.value.Store(s)
	return x
}

// String returns a string representation of the value
func (x *String) String() string {
	return fmt.Sprintf("%s: '%s'", x.Metadata.Option, x.V())
}

// MarshalJSON returns the json representation
func (x *String) MarshalJSON() (b []byte, e error) {
	v := string(x.value.Load().([]byte))
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *String) UnmarshalJSON(data []byte) (e error) {
	v := x.value.Load().([]byte)
	e = json.Unmarshal(data, &v)
	x.value.Store(v)
	return
}
