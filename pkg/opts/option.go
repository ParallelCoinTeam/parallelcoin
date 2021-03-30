package opts

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
	Opt  Option
	Name string
}
type ConfigSlice []ConfigSliceElement

func (c ConfigSlice) Len() int           { return len(c) }
func (c ConfigSlice) Less(i, j int) bool { return c[i].Name < c[j].Name }
func (c ConfigSlice) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

