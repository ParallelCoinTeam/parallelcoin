package txscript

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// tstCheckScriptError ensures the type of the two passed errors are of the same type (either both nil or both of type
// ScriptError) and their error codes match when not nil.
func tstCheckScriptError(gotErr, wantErr error) (e error) {
	// Ensure the error code is of the expected type and the error code matches the value specified in the test
	// instance.
	if reflect.TypeOf(gotErr) != reflect.TypeOf(wantErr) {
		return fmt.Errorf("wrong error - got %T (%[1]v), want %T",
			gotErr, wantErr,
		)
	}
	if gotErr == nil {
		return nil
	}
	// Ensure the want error type is a script error.
	werr, ok := wantErr.(ScriptError)
	if !ok {
		return fmt.Errorf("unexpected test error type %T", wantErr)
	}
	// Ensure the error codes match. It's safe to use a raw type assert here since the code above already proved they
	// are the same type and the want error is a script error.
	gotErrorCode := gotErr.(ScriptError).ErrorCode
	if gotErrorCode != werr.ErrorCode {
		return fmt.Errorf("mismatched error code - got %v (%v), want %v",
			gotErrorCode, gotErr, werr.ErrorCode,
		)
	}
	return nil
}

// TestStack tests that all of the stack operations work as expected.
func TestStack(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		before    [][]byte
		operation func(*stack) error
		err       error
		after     [][]byte
	}{
		{
			"noop",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				return nil
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {5}},
		},
		{
			"peek underflow (byte)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				_, e = s.PeekByteArray(5)
				return e
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"peek underflow (int)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				_, e = s.PeekInt(5)
				return e
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"peek underflow (bool)",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				_, e = s.PeekBool(5)
				return e
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"pop",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				val, e := s.PopByteArray()
				if e != nil {
					return e
				}
				if !bytes.Equal(val, []byte{5}) {
					return errors.New("not equal")
				}
				return e
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}},
		},
		{
			"pop everything",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				for i := 0; i < 5; i++ {
					_, e := s.PopByteArray()
					if e != nil {
						return e
					}
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop underflow",
			[][]byte{{1}, {2}, {3}, {4}, {5}},
			func(s *stack) (e error) {
				for i := 0; i < 6; i++ {
					_, e := s.PopByteArray()
					if e != nil {
						return e
					}
				}
				return nil
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"pop bool",
			[][]byte{nil},
			func(s *stack) (e error) {
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop bool",
			[][]byte{{1}},
			func(s *stack) (e error) {
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if !val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop bool",
			nil,
			func(s *stack) (e error) {
				_, e = s.PopBool()
				return e
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"popInt 0",
			[][]byte{{0x0}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != 0 {
					return errors.New("0 != 0 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt -0",
			[][]byte{{0x80}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != 0 {
					return errors.New("-0 != 0 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt 1",
			[][]byte{{0x01}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != 1 {
					return errors.New("1 != 1 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt 1 leading 0",
			[][]byte{{0x01, 0x00, 0x00, 0x00}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != 1 {
					fmt.Printf("%v != %v\n", v, 1)
					return errors.New("1 != 1 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt -1",
			[][]byte{{0x81}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != -1 {
					return errors.New("-1 != -1 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"popInt -1 leading 0",
			[][]byte{{0x01, 0x00, 0x00, 0x80}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != -1 {
					fmt.Printf("%v != %v\n", v, -1)
					return errors.New("-1 != -1 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		// Triggers the multibyte case in asInt
		{
			"popInt -513",
			[][]byte{{0x1, 0x82}},
			func(s *stack) (e error) {
				v, e := s.PopInt()
				if e != nil {
					return e
				}
				if v != -513 {
					fmt.Printf("%v != %v\n", v, -513)
					return errors.New("1 != 1 on popInt")
				}
				return nil
			},
			nil,
			nil,
		},
		// Confirm that the asInt code doesn't modify the base data.
		{
			"peekint nomodify -1",
			[][]byte{{0x01, 0x00, 0x00, 0x80}},
			func(s *stack) (e error) {
				v, e := s.PeekInt(0)
				if e != nil {
					return e
				}
				if v != -1 {
					fmt.Printf("%v != %v\n", v, -1)
					return errors.New("-1 != -1 on popInt")
				}
				return nil
			},
			nil,
			[][]byte{{0x01, 0x00, 0x00, 0x80}},
		},
		{
			"PushInt 0",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(0))
				return nil
			},
			nil,
			[][]byte{{}},
		},
		{
			"PushInt 1",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(1))
				return nil
			},
			nil,
			[][]byte{{0x1}},
		},
		{
			"PushInt -1",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(-1))
				return nil
			},
			nil,
			[][]byte{{0x81}},
		},
		{
			"PushInt two bytes",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(256))
				return nil
			},
			nil,
			// little endian.. *sigh*
			[][]byte{{0x00, 0x01}},
		},
		{
			"PushInt leading zeros",
			nil,
			func(s *stack) (e error) {
				// this will have the highbit set
				s.PushInt(scriptNum(128))
				return nil
			},
			nil,
			[][]byte{{0x80, 0x00}},
		},
		{
			"dup",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.DupN(1)
			},
			nil,
			[][]byte{{1}, {1}},
		},
		{
			"dup2",
			[][]byte{{1}, {2}},
			func(s *stack) (e error) {
				return s.DupN(2)
			},
			nil,
			[][]byte{{1}, {2}, {1}, {2}},
		},
		{
			"dup3",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.DupN(3)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {1}, {2}, {3}},
		},
		{
			"dup0",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.DupN(0)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"dup-1",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.DupN(-1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"dup too much",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.DupN(2)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"PushBool true",
			nil,
			func(s *stack) (e error) {
				s.PushBool(true)
				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"PushBool false",
			nil,
			func(s *stack) (e error) {
				s.PushBool(false)
				return nil
			},
			nil,
			[][]byte{nil},
		},
		{
			"PushBool PopBool",
			nil,
			func(s *stack) (e error) {
				s.PushBool(true)
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if !val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"PushBool PopBool 2",
			nil,
			func(s *stack) (e error) {
				s.PushBool(false)
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"PushInt PopBool",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(1))
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if !val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"PushInt PopBool 2",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(0))
				val, e := s.PopBool()
				if e != nil {
					return e
				}
				if val {
					return errors.New("unexpected value")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"Nip top",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.NipN(0)
			},
			nil,
			[][]byte{{1}, {2}},
		},
		{
			"Nip middle",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.NipN(1)
			},
			nil,
			[][]byte{{1}, {3}},
		},
		{
			"Nip low",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.NipN(2)
			},
			nil,
			[][]byte{{2}, {3}},
		},
		{
			"Nip too much",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				// bite off more than we can chew
				return s.NipN(3)
			},
			scriptError(ErrInvalidStackOperation, ""),
			[][]byte{{2}, {3}},
		},
		{
			"keep on tucking",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.Tuck()
			},
			nil,
			[][]byte{{1}, {3}, {2}, {3}},
		},
		{
			"a little tucked up",
			[][]byte{{1}}, // too few arguments for tuck
			func(s *stack) (e error) {
				return s.Tuck()
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"all tucked up",
			nil, // too few arguments  for tuck
			func(s *stack) (e error) {
				return s.Tuck()
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"drop 1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}},
		},
		{
			"drop 2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(2)
			},
			nil,
			[][]byte{{1}, {2}},
		},
		{
			"drop 3",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(3)
			},
			nil,
			[][]byte{{1}},
		},
		{
			"drop 4",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(4)
			},
			nil,
			nil,
		},
		{
			"drop 4/5",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(5)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"drop invalid",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.DropN(0)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Rot1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.RotN(1)
			},
			nil,
			[][]byte{{1}, {3}, {4}, {2}},
		},
		{
			"Rot2",
			[][]byte{{1}, {2}, {3}, {4}, {5}, {6}},
			func(s *stack) (e error) {
				return s.RotN(2)
			},
			nil,
			[][]byte{{3}, {4}, {5}, {6}, {1}, {2}},
		},
		{
			"Rot too little",
			[][]byte{{1}, {2}},
			func(s *stack) (e error) {
				return s.RotN(1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Rot0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.RotN(0)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Swap1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.SwapN(1)
			},
			nil,
			[][]byte{{1}, {2}, {4}, {3}},
		},
		{
			"Swap2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.SwapN(2)
			},
			nil,
			[][]byte{{3}, {4}, {1}, {2}},
		},
		{
			"Swap too little",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.SwapN(1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Swap0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.SwapN(0)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Over1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.OverN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {3}},
		},
		{
			"Over2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.OverN(2)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {1}, {2}},
		},
		{
			"Over too little",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.OverN(1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Over0",
			[][]byte{{1}, {2}, {3}},
			func(s *stack) (e error) {
				return s.OverN(0)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Pick1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.PickN(1)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {3}},
		},
		{
			"Pick2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.PickN(2)
			},
			nil,
			[][]byte{{1}, {2}, {3}, {4}, {2}},
		},
		{
			"Pick too little",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.PickN(1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Roll1",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.RollN(1)
			},
			nil,
			[][]byte{{1}, {2}, {4}, {3}},
		},
		{
			"Roll2",
			[][]byte{{1}, {2}, {3}, {4}},
			func(s *stack) (e error) {
				return s.RollN(2)
			},
			nil,
			[][]byte{{1}, {3}, {4}, {2}},
		},
		{
			"Roll too little",
			[][]byte{{1}},
			func(s *stack) (e error) {
				return s.RollN(1)
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
		{
			"Peek bool",
			[][]byte{{1}},
			func(s *stack) (e error) {
				// Peek bool is otherwise pretty well tested, just check it works.
				val, e := s.PeekBool(0)
				if e != nil {
					return e
				}
				if !val {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"Peek bool 2",
			[][]byte{nil},
			func(s *stack) (e error) {
				// Peek bool is otherwise pretty well tested, just check it works.
				val, e := s.PeekBool(0)
				if e != nil {
					return e
				}
				if val {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{nil},
		},
		{
			"Peek int",
			[][]byte{{1}},
			func(s *stack) (e error) {
				// Peek int is otherwise pretty well tested, just check it works.
				val, e := s.PeekInt(0)
				if e != nil {
					return e
				}
				if val != 1 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{{1}},
		},
		{
			"Peek int 2",
			[][]byte{{0}},
			func(s *stack) (e error) {
				// Peek int is otherwise pretty well tested, just check it works.
				val, e := s.PeekInt(0)
				if e != nil {
					return e
				}
				if val != 0 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			[][]byte{{0}},
		},
		{
			"pop int",
			nil,
			func(s *stack) (e error) {
				s.PushInt(scriptNum(1))
				// Peek int is otherwise pretty well tested, just check it works.
				val, e := s.PopInt()
				if e != nil {
					return e
				}
				if val != 1 {
					return errors.New("invalid result")
				}
				return nil
			},
			nil,
			nil,
		},
		{
			"pop empty",
			nil,
			func(s *stack) (e error) {
				// Peek int is otherwise pretty well tested, just check it works.
				_, e = s.PopInt()
				return e
			},
			scriptError(ErrInvalidStackOperation, ""),
			nil,
		},
	}
	var e error
	for _, test := range tests {
		// Setup the initial stack state and perform the test operation.
		s := stack{}
		for i := range test.before {
			s.PushByteArray(test.before[i])
		}
		e = test.operation(&s)
		// Ensure the error code is of the expected type and the error code matches the value specified in the test instance.
		if e = tstCheckScriptError(e, test.err); e != nil {
			t.Errorf("%s: %v", test.name, e)
			continue
		}
		// Ensure the resulting stack is the expected length.
		if int32(len(test.after)) != s.Depth() {
			t.Errorf("%s: stack depth doesn't match expected: %v "+
				"vs %v", test.name, len(test.after),
				s.Depth(),
			)
			continue
		}
		// Ensure all items of the resulting stack are the expected values.
		for i := range test.after {
			var val []byte
			val, e = s.PeekByteArray(s.Depth() - int32(i) - 1)
			if e != nil {
				t.Errorf("%s: can't peek %dth stack entry: %v",
					test.name, i, e,
				)
				break
			}
			if !bytes.Equal(val, test.after[i]) {
				t.Errorf("%s: %dth stack entry doesn't match "+
					"expected: %v vs %v", test.name, i, val,
					test.after[i],
				)
				break
			}
		}
	}
}
