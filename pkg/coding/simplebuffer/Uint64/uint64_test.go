package Uint64

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestInt32(t *testing.T) {
	by, err := hex.DecodeString("deadbeefcafeb00b")
	if err != nil {
		panic(err)
	}
	bits := binary.BigEndian.Uint64(by)
	bt := New()
	bt.Put(bits)
	bt2 := New()
	bt2.Decode(bt.Encode())
	if bt.Get() != bt2.Get() {
		t.Fail()
	}
}
