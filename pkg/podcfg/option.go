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
		LoadInput(input string) (o Option, e error)
		ReadInput(input string) (o Option, e error)
		GetMetadata() *Metadata
		Name() string
		String() string
		MarshalJSON() (b []byte, e error)
		UnmarshalJSON(data []byte) (e error)
		GetAllOptionStrings() []string
		Type() interface{}
		SetName(string)
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
	Name        string
}
)

func (m Metadata) GetAllOptionStrings() (opts []string) {
	opts = append([]string{m.Option}, m.Aliases...)
	return opts
}

// Configs is the source location for the Config items, which is used to generate the Config struct
type Configs map[string]Option
type ConfigSliceElement struct {
	Opt Option
	Name string
}
type ConfigSlice []ConfigSliceElement

func (c ConfigSlice) Len() int           { return len(c) }
func (c ConfigSlice) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c ConfigSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

