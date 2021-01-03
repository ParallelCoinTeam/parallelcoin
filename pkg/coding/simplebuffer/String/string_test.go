package String

import (
	"testing"
)

func TestString(t *testing.T) {
	by := "this is a test"
	bt := New()
	bt.Put(by)
	bt2 := New()
	bt2.Decode(bt.Encode())
	if bt.Get() != bt2.Get() {
		t.Fail()
	}
}
