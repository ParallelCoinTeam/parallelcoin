package json_test

import (
	"testing"

	"github.com/p9c/pod/pkg/rpc/json"
)

// TestErrorCodeStringer tests the stringized output for the ErrorCode type.
func TestErrorCodeStringer(	t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   json.ErrorCode
		want string
	}{
		{json.ErrDuplicateMethod, "ErrDuplicateMethod"},
		{json.ErrInvalidUsageFlags, "ErrInvalidUsageFlags"},
		{json.ErrInvalidType, "ErrInvalidType"},
		{json.ErrEmbeddedType, "ErrEmbeddedType"},
		{json.ErrUnexportedField, "ErrUnexportedField"},
		{json.ErrUnsupportedFieldType, "ErrUnsupportedFieldType"},
		{json.ErrNonOptionalField, "ErrNonOptionalField"},
		{json.ErrNonOptionalDefault, "ErrNonOptionalDefault"},
		{json.ErrMismatchedDefault, "ErrMismatchedDefault"},
		{json.ErrUnregisteredMethod, "ErrUnregisteredMethod"},
		{json.ErrNumParams, "ErrNumParams"},
		{json.ErrMissingDescription, "ErrMissingDescription"},
		{0xffff, "Unknown ErrorCode (65535)"},
	}
	// Detect additional error codes that don't have the stringer added.
	if len(tests)-1 != int(json.TstNumErrorCodes) {
		t.Errorf("It appears an error code was added without adding an " +
			"associated stringer test")
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.String()
		if result != test.want {
			t.Errorf("String #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}

// TestError tests the error output for the Error type.
func TestError(	t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   json.Error
		want string
	}{
		{
			json.Error{Description: "some error"},
			"some error",
		},
		{
			json.Error{Description: "human-readable error"},
			"human-readable error",
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		result := test.in.Error()
		if result != test.want {
			t.Errorf("Error #%d\n got: %s want: %s", i, result,
				test.want)
			continue
		}
	}
}
