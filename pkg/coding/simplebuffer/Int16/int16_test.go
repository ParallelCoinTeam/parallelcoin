package Int16

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestInt16(t *testing.T) {
	by, err := hex.DecodeString("beef")
	if err != nil {
		panic(err)
	}
	bits := binary.BigEndian.Uint16(by)
	bt := New()
	bt.Put(int16(bits))
	bt2 := New()
	bt2.Decode(bt.Encode())
	if bt.Get() != bt2.Get() {
		t.Fail()
	}
}
