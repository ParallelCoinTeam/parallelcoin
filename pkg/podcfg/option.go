// Package podcfg implements a concurrent/parallel active application configuration system for multi-process
// applications to share a configuration as well as keep in sync with each other.
//
// This file contains all of the data types stored in a podcfg.Config and the various accessors and methods relevant to
// them. There is a basic byte slice-as-string type which is intended to eventually cover proper security practices for
// storing password information.
package podcfg

type (
	// Option is an interface to simplify concurrent-safe access to a variety of types of configuration item
	Option interface {
		ReadInput(string) (o Option, e error)
		GetMetadata() *Metadata
		Name() string
		String() string
		MarshalJSON() (b []byte, e error)
		UnmarshalJSON(data []byte) (e error)
		GetAllOptionStrings() []string
		Type() interface{}
	}
	// Metadata is the information about the option to be used by interface code and other presentations of the data
	Metadata struct {
		Option      string
		Aliases     []string
		Group       string
		Label       string
		Description string
		Type        string
		Widget      string
		Options     []string
		OmitEmpty   bool
	}
)

func (m Metadata) GetAllOptionStrings() (opts []string) {
	opts = append([]string{m.Option}, m.Aliases...)
	return opts
}
