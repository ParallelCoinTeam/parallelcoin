package Uint64

import "encoding/binary"

// Uint64 is a 32 bit value that stores an uint64
type Uint64 struct {
	Bytes [8]byte
}

func New() *Uint64 {
	return &Uint64{}
}

func (b *Uint64) DecodeOne(by []byte) *Uint64 {
	b.Decode(by)
	return b
}

func (b *Uint64) Decode(by []byte) (out []byte) {
	if len(by) >= 8 {
		b.Bytes = [8]byte{
			by[0], by[1], by[2], by[3],
			by[4], by[5], by[6], by[7],
		}
		if len(by) > 8 {
			out = by[8:]
		}
	}
	return
}

func (b *Uint64) Encode() []byte {
	return b.Bytes[:]
}

func (b *Uint64) Get() uint64 {
	return binary.BigEndian.Uint64(b.Bytes[:])
}

func (b *Uint64) Put(bits uint64) *Uint64 {
	binary.BigEndian.PutUint64(b.Bytes[:], bits)
	return b
}
