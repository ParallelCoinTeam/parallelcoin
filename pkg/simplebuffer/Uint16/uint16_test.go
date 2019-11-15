package Uint16

import "testing"

func TestUint16(t *testing.T) {
	var example uint16 = 11047
	u := New()
	u.Put(example)
	port2 := New()
	port2.Decode(u.Encode())
	if port2.Get() != u.Get() {
		t.Fail()
	}
}
