package btcjson_test

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestAssignField tests the assignField function handles supported combinations properly.
func TestAssignField(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		dest     interface{}
		src      interface{}
		expected interface{}
	}{
		{
			name:     "same types",
			dest:     int8(0),
			src:      int8(100),
			expected: int8(100),
		},
		{
			name: "same types - more source pointers",
			dest: int8(0),
			src: func() interface{} {
				i := int8(100)
				return &i
			}(),
			expected: int8(100),
		},
		{
			name: "same types - more dest pointers",
			dest: func() interface{} {
				i := int8(0)
				return &i
			}(),
			src:      int8(100),
			expected: int8(100),
		},
		{
			name: "convertible types - more source pointers",
			dest: int16(0),
			src: func() interface{} {
				i := int8(100)
				return &i
			}(),
			expected: int16(100),
		},
		{
			name: "convertible types - both pointers",
			dest: func() interface{} {
				i := int8(0)
				return &i
			}(),
			src: func() interface{} {
				i := int16(100)
				return &i
			}(),
			expected: int8(100),
		},
		{
			name:     "convertible types - int16 -> int8",
			dest:     int8(0),
			src:      int16(100),
			expected: int8(100),
		},
		{
			name:     "convertible types - int16 -> uint8",
			dest:     uint8(0),
			src:      int16(100),
			expected: uint8(100),
		},
		{
			name:     "convertible types - uint16 -> int8",
			dest:     int8(0),
			src:      uint16(100),
			expected: int8(100),
		},
		{
			name:     "convertible types - uint16 -> uint8",
			dest:     uint8(0),
			src:      uint16(100),
			expected: uint8(100),
		},
		{
			name:     "convertible types - float32 -> float64",
			dest:     float64(0),
			src:      float32(1.5),
			expected: float64(1.5),
		},
		{
			name:     "convertible types - float64 -> float32",
			dest:     float32(0),
			src:      float64(1.5),
			expected: float32(1.5),
		},
		{
			name:     "convertible types - string -> bool",
			dest:     false,
			src:      "true",
			expected: true,
		},
		{
			name:     "convertible types - string -> int8",
			dest:     int8(0),
			src:      "100",
			expected: int8(100),
		},
		{
			name:     "convertible types - string -> uint8",
			dest:     uint8(0),
			src:      "100",
			expected: uint8(100),
		},
		{
			name:     "convertible types - string -> float32",
			dest:     float32(0),
			src:      "1.5",
			expected: float32(1.5),
		},
		{
			name: "convertible types - typecase string -> string",
			dest: "",
			src: func() interface{} {
				type foo string
				return foo("foo")
			}(),
			expected: "foo",
		},
		{
			name:     "convertible types - string -> array",
			dest:     [2]string{},
			src:      `["test","test2"]`,
			expected: [2]string{"test", "test2"},
		},
		{
			name:     "convertible types - string -> slice",
			dest:     []string{},
			src:      `["test","test2"]`,
			expected: []string{"test", "test2"},
		},
		{
			name:     "convertible types - string -> struct",
			dest:     struct{ A int }{},
			src:      `{"A":100}`,
			expected: struct{ A int }{100},
		},
		{
			name:     "convertible types - string -> map",
			dest:     map[string]float64{},
			src:      `{"1Address":1.5}`,
			expected: map[string]float64{"1Address": 1.5},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		dst := reflect.New(reflect.TypeOf(test.dest)).Elem()
		src := reflect.ValueOf(test.src)
		e := btcjson.TstAssignField(1, "testField", dst, src)
		if e != nil {
			t.Errorf("Test #%d (%s) unexpected error: %v", i,
				test.name, e,
			)
			continue
		}
		// Indirect through to the base types to ensure their values are the same.
		for dst.Kind() == reflect.Ptr {
			dst = dst.Elem()
		}
		if !reflect.DeepEqual(dst.Interface(), test.expected) {
			t.Errorf("Test #%d (%s) unexpected value - got %v, "+
				"want %v", i, test.name, dst.Interface(),
				test.expected,
			)
			continue
		}
	}
}

// TestAssignFieldErrors tests the assignField function error paths.
func TestAssignFieldErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dest interface{}
		src  interface{}
		e    btcjson.GeneralError
	}{
		{
			name: "general incompatible int -> string",
			dest: string(rune(0)),
			src:  0,
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow source int -> dest int",
			dest: int8(0),
			src:  128,
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow source int -> dest uint",
			dest: uint8(0),
			src:  256,
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "int -> float",
			dest: float32(0),
			src:  256,
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow source uint64 -> dest int64",
			dest: int64(0),
			src:  uint64(1 << 63),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow source uint -> dest int",
			dest: int8(0),
			src:  uint(128),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow source uint -> dest uint",
			dest: uint8(0),
			src:  uint(256),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "uint -> float",
			dest: float32(0),
			src:  uint(256),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "float -> int",
			dest: 0,
			src:  float32(1.0),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow float64 -> float32",
			dest: float32(0),
			src:  float64(math.MaxFloat64),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> bool",
			dest: true,
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> int",
			dest: int8(0),
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow string -> int",
			dest: int8(0),
			src:  "128",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> uint",
			dest: uint8(0),
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow string -> uint",
			dest: uint8(0),
			src:  "256",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> float",
			dest: float32(0),
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "overflow string -> float",
			dest: float32(0),
			src:  "1.7976931348623157e+308",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> array",
			dest: [3]int{},
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> slice",
			dest: []int{},
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> struct",
			dest: struct{ A int }{},
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid string -> map",
			dest: map[string]int{},
			src:  "foo",
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		dst := reflect.New(reflect.TypeOf(test.dest)).Elem()
		src := reflect.ValueOf(test.src)
		e := btcjson.TstAssignField(1, "testField", dst, src)
		if reflect.TypeOf(e) != reflect.TypeOf(test.e) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%[3]v), "+
				"want %T", i, test.name, e, test.e,
			)
			continue
		}
		gotErrorCode := e.(btcjson.GeneralError).ErrorCode
		if gotErrorCode != test.e.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v (%v), want %v", i, test.name, gotErrorCode,
				e, test.e.ErrorCode,
			)
			continue
		}
	}
}

// TestNewCmdErrors ensures the error paths of NewCmd behave as expected.
func TestNewCmdErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		method string
		args   []interface{}
		e      btcjson.GeneralError
	}{
		{
			name:   "unregistered command",
			method: "boguscommand",
			args:   []interface{}{},
			e:      btcjson.GeneralError{ErrorCode: btcjson.ErrUnregisteredMethod},
		},
		{
			name:   "too few parameters to command with required + optional",
			method: "getblock",
			args:   []interface{}{},
			e:      btcjson.GeneralError{ErrorCode: btcjson.ErrNumParams},
		},
		{
			name:   "too many parameters to command with no optional",
			method: "getblockcount",
			args:   []interface{}{"123"},
			e:      btcjson.GeneralError{ErrorCode: btcjson.ErrNumParams},
		},
		{
			name:   "incorrect parameter type",
			method: "getblock",
			args:   []interface{}{1},
			e:      btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		_, e := btcjson.NewCmd(test.method, test.args...)
		if reflect.TypeOf(e) != reflect.TypeOf(test.e) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, e, e, test.e,
			)
			continue
		}
		gotErrorCode := e.(btcjson.GeneralError).ErrorCode
		if gotErrorCode != test.e.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v (%v), want %v", i, test.name, gotErrorCode,
				e, test.e.ErrorCode,
			)
			continue
		}
	}
}

// TestMarshalCmdErrors  tests the error paths of the MarshalCmd function.
func TestMarshalCmdErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		id   interface{}
		cmd  interface{}
		e    btcjson.GeneralError
	}{
		{
			name: "unregistered type",
			id:   1,
			cmd:  (*int)(nil),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrUnregisteredMethod},
		},
		{
			name: "nil instance of registered type",
			id:   1,
			cmd:  (*btcjson.GetBlockCmd)(nil),
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "nil instance of registered type",
			id:   []int{0, 1},
			cmd:  &btcjson.GetBlockCountCmd{},
			e:    btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		_, e := btcjson.MarshalCmd(test.id, test.cmd)
		if reflect.TypeOf(e) != reflect.TypeOf(test.e) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, e, e, test.e,
			)
			continue
		}
		gotErrorCode := e.(btcjson.GeneralError).ErrorCode
		if gotErrorCode != test.e.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v (%v), want %v", i, test.name, gotErrorCode,
				e, test.e.ErrorCode,
			)
			continue
		}
	}
}

// TestUnmarshalCmdErrors  tests the error paths of the UnmarshalCmd function.
func TestUnmarshalCmdErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		request btcjson.Request
		e       btcjson.GeneralError
	}{
		{
			name: "unregistered type",
			request: btcjson.Request{
				Jsonrpc: "1.0",
				Method:  "bogusmethod",
				Params:  nil,
				ID:      nil,
			},
			e: btcjson.GeneralError{ErrorCode: btcjson.ErrUnregisteredMethod},
		},
		{
			name: "incorrect number of netparams",
			request: btcjson.Request{
				Jsonrpc: "1.0",
				Method:  "getblockcount",
				Params:  []json.RawMessage{[]byte(`"bogusparam"`)},
				ID:      nil,
			},
			e: btcjson.GeneralError{ErrorCode: btcjson.ErrNumParams},
		},
		{
			name: "invalid type for a parameter",
			request: btcjson.Request{
				Jsonrpc: "1.0",
				Method:  "getblock",
				Params:  []json.RawMessage{[]byte("1")},
				ID:      nil,
			},
			e: btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
		{
			name: "invalid JSON for a parameter",
			request: btcjson.Request{
				Jsonrpc: "1.0",
				Method:  "getblock",
				Params:  []json.RawMessage{[]byte(`"1`)},
				ID:      nil,
			},
			e: btcjson.GeneralError{ErrorCode: btcjson.ErrInvalidType},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		_, e := btcjson.UnmarshalCmd(&test.request)
		if reflect.TypeOf(e) != reflect.TypeOf(test.e) {
			t.Errorf("Test #%d (%s) wrong error - got %T (%v), "+
				"want %T", i, test.name, e, e, test.e,
			)
			continue
		}
		gotErrorCode := e.(btcjson.GeneralError).ErrorCode
		if gotErrorCode != test.e.ErrorCode {
			t.Errorf("Test #%d (%s) mismatched error code - got "+
				"%v (%v), want %v", i, test.name, gotErrorCode,
				e, test.e.ErrorCode,
			)
			continue
		}
	}
}
