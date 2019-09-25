package json_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/p9c/pod/pkg/rpc/json"
)

// TestUsageFlagStringer tests the stringized output for the UsageFlag type.
func TestUsageFlagStringer(	t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   json.UsageFlag
		want string
	}{
		{0, "0x0"},
		{json.UFWalletOnly, "UFWalletOnly"},
		{json.UFWebsocketOnly, "UFWebsocketOnly"},
		{json.UFNotification, "UFNotification"},
		{json.UFWalletOnly | json.UFWebsocketOnly,
			"UFWalletOnly|UFWebsocketOnly"},
		{json.UFWalletOnly | json.UFWebsocketOnly | (1 << 31),
			"UFWalletOnly|UFWebsocketOnly|0x80000000"},
	}
	// Detect additional usage flags that don't have the stringer added.
	numUsageFlags := 0
	highestUsageFlagBit := json.TstHighestUsageFlagBit
	for highestUsageFlagBit > 1 {
		numUsageFlags++
		highestUsageFlagBit >>= 1
	}
	if len(tests)-3 != numUsageFlags {
		t.Errorf("It appears a usage flag was added without adding " +
			"an associated stringer test")
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

// TestRegisterCmdErrors ensures the RegisterCmd function returns the expected error when provided with invalid types.
func TestRegisterCmdErrors(	t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		method  string
		cmdFunc func() interface{}
		flags   json.UsageFlag
		err     json.Error
	}{
		{
			name:   "duplicate method",
			method: "getblock",
			cmdFunc: func() interface{} {
				return struct{}{}
			},
			err: json.Error{ErrorCode: json.ErrDuplicateMethod},
		},
		{
			name:   "invalid usage flags",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			flags: json.TstHighestUsageFlagBit,
			err:   json.Error{ErrorCode: json.ErrInvalidUsageFlags},
		},
		{
			name:   "invalid type",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return 0
			},
			err: json.Error{ErrorCode: json.ErrInvalidType},
		},
		{
			name:   "invalid type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				return &[]string{}
			},
			err: json.Error{ErrorCode: json.ErrInvalidType},
		},
		{
			name:   "embedded field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				//nolint
				type test struct{ int }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrEmbeddedType},
		},
		{
			name:   "unexported field",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				//nolint
				type test struct{ a int }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnexportedField},
		},
		{
			name:   "unsupported field type 1",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A **int }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 2",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A chan int }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 3",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex64 }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 4",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A complex128 }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 5",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A func() }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "unsupported field type 6",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct{ A interface{} }
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrUnsupportedFieldType},
		},
		{
			name:   "required after optional",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int
					B int
				}
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrNonOptionalField},
		},
		{
			name:   "non-optional with default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A int `jsonrpcdefault:"1"`
				}
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrNonOptionalDefault},
		},
		{
			name:   "mismatched default",
			method: "registertestcmd",
			cmdFunc: func() interface{} {
				type test struct {
					A *int `jsonrpcdefault:"1.7"`
				}
				return (*test)(nil)
			},
			err: json.Error{ErrorCode: json.ErrMismatchedDefault},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		err := json.RegisterCmd(test.method, test.cmdFunc(),
			test.flags)
		if reflect.TypeOf(err) != reflect.TypeOf(test.err) {
			t.Errorf("Test #%d (%s) wrong error - got %T, "+
				"want %T", i, test.name, err, test.err)
			continue
		}
		gotErrorCode := err.(json.Error).ErrorCode
		if gotErrorCode != test.err.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v, want %v", i, test.name, gotErrorCode,
				test.err.ErrorCode)
			continue
		}
	}
}

// TestMustRegisterCmdPanic ensures the MustRegisterCmd function panics when used to register an invalid type.
func TestMustRegisterCmdPanic(	t *testing.T) {
	t.Parallel()
	// Setup a defer to catch the expected panic to ensure it actually
	// paniced.
	defer func() {
		if err := recover(); err == nil {
			t.Error("MustRegisterCmd did not panic as expected")
		}
	}()
	// Intentionally try to register an invalid type to force a panic.
	json.MustRegisterCmd("panicme", 0, 0)
}

// TestRegisteredCmdMethods tests the RegisteredCmdMethods function ensure it works as expected.
func TestRegisteredCmdMethods(	t *testing.T) {
	t.Parallel()
	// Ensure the registered methods are returned.
	methods := json.RegisteredCmdMethods()
	if len(methods) == 0 {
		t.Fatal("RegisteredCmdMethods: no methods")
	}
	// Ensure the returned methods are sorted.
	sortedMethods := make([]string, len(methods))
	copy(sortedMethods, methods)
	sort.Strings(sortedMethods)
	if !reflect.DeepEqual(sortedMethods, methods) {
		t.Fatal("RegisteredCmdMethods: methods are not sorted")
	}
}
