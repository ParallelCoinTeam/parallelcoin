package Uint32

import "encoding/binary"

// Uint32 is a 32 bit value that stores an int32 (used for block height).
// I don't think the sign is preserved but block heights are never negative
// except with special semantics
type Uint32 struct {
	Bytes [4]byte
}

func New() *Uint32 {
	return &Uint32{}
}

func (b *Uint32) DecodeOne(by []byte) *Uint32 {
	b.Decode(by)
	return b
}

func (b *Uint32) Decode(by []byte) (out []byte) {
	if len(by) >= 4 {
		b.Bytes = [4]byte{by[0], by[1], by[2], by[3]}
		if len(by) > 4 {
			out = by[4:]
		}
	}
	return
}

func (b *Uint32) Encode() []byte {
	return b.Bytes[:]
}

func (b *Uint32) Get() uint32 {
	return binary.BigEndian.Uint32(b.Bytes[:])
}

func (b *Uint32) Put(bits uint32) *Uint32 {
	binary.BigEndian.PutUint32(b.Bytes[:], bits)
	return b
}
