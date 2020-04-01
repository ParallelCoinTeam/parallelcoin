package Int64

import "encoding/binary"

// Int64 is a 32 bit value that stores an int32 (used for block height).
// I don't think the sign is preserved but block heights are never negative
// except with special semantics
type Int64 struct {
	Bytes [8]byte
}

func New() *Int64 {
	return &Int64{}
}

func (b *Int64) DecodeOne(by []byte) *Int64 {
	b.Decode(by)
	return b
}

func (b *Int64) Decode(by []byte) (out []byte) {
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

func (b *Int64) Encode() []byte {
	return b.Bytes[:]
}

func (b *Int64) Get() int64 {
	return int64(binary.BigEndian.Uint64(b.Bytes[:]))
}

func (b *Int64) Put(bits int64) *Int64 {
	binary.BigEndian.PutUint64(b.Bytes[:], uint64(bits))
	return b
}
