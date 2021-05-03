package opt

import (
	"github.com/p9c/opts/meta"
)

type (
	// Option is an interface to simplify concurrent-safe access to a variety of types of configuration item
	Option interface {
		LoadInput(input string) (o Option, e error)
		ReadInput(input string) (o Option, e error)
		GetMetadata() *meta.Data
		Name() string
		String() string
		MarshalJSON() (b []byte, e error)
		UnmarshalJSON(data []byte) (e error)
		GetAllOptionStrings() []string
		Type() interface{}
		SetName(string)
	}
)
