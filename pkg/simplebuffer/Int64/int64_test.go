package Int64

import (
	"encoding/binary"
	"encoding/hex"
	"testing"
)

func TestInt64(t *testing.T) {
	by, err := hex.DecodeString("deadbeefcafe8008")
	if err != nil {
		panic(err)
	}
	bits := binary.BigEndian.Uint64(by)
	bt := New()
	bt.Put(int64(bits))
	bt2 := New()
	bt2.Decode(bt.Encode())
	if bt.Get() != bt2.Get() {
		t.Fail()
	}
}
