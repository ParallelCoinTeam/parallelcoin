package Int16

import "encoding/binary"

// Int16 is a 16 bit signed integer
type Int16 struct {
	Bytes [2]byte
}

func New() *Int16 {
	return &Int16{}
}

func (b *Int16) DecodeOne(by []byte) *Int16 {
	b.Decode(by)
	return b
}

func (b *Int16) Decode(by []byte) (out []byte) {
	if len(by) >= 2 {
		b.Bytes = [2]byte{by[0], by[1]}
		if len(by) > 2 {
			out = by[2:]
		}
	}
	return
}

func (b *Int16) Encode() []byte {
	return b.Bytes[:]
}

func (b *Int16) Get() int16 {
	return int16(binary.BigEndian.Uint16(b.Bytes[:]))
}

func (b *Int16) Put(bits int16) *Int16 {
	binary.BigEndian.PutUint16(b.Bytes[:], uint16(bits))
	return b
}
