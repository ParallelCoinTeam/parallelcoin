package cfgutil

// ExplicitString is a string value implementing the flags.Marshaler and flags.Unmarshaler interfaces so it may be used
// as a config struct field. It records whether the value was explicitly set by the flags package.
//
// This is useful when behavior must be modified depending on whether a flag was set by the user or left as a default.
//
// Without recording this, it would be impossible to determine whether flag with a default value was unmodified or
// explicitly set to the default.
type ExplicitString struct {
	Value         string
	explicitlySet bool
}

// NewExplicitString creates a string flag with the provided default value.
func NewExplicitString(defaultValue string) *ExplicitString {
	return &ExplicitString{Value: defaultValue, explicitlySet: false}
}

// ExplicitlySet returns whether the flag was explicitly set through the flags.Unmarshaler interface.
func (es *ExplicitString) ExplicitlySet() bool { return es.explicitlySet }

// MarshalFlag implements the flags.Marshaler interface.
func (es *ExplicitString) MarshalFlag() (string, error) {
	return es.Value, nil
}

// UnmarshalFlag implements the flags.Unmarshaler interface.
func (es *ExplicitString) UnmarshalFlag(value string) (e error) {
	es.Value = value
	es.explicitlySet = true
	return nil
}
