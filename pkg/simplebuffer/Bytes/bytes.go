package Bytes

import "encoding/binary"

// Bytes plain old bytes. Maximum length from 32 bits int
type Bytes struct {
	Bytes []byte
}

func New() *Bytes {
	return &Bytes{}
}

func (b *Bytes) DecodeOne(by []byte) *Bytes {
	b.Decode(by)
	return b
}

func (b *Bytes) Decode(by []byte) (out []byte) {
	if len(by) >= 4 {
		length := binary.BigEndian.Uint32(by[:4])
		if len(by) >= 4+int(length) {
			out = by[4 : 4+length]
			b.Bytes = out
		}
	}
	return
}

func (b *Bytes) Encode() []byte {
	by := make([]byte, 4)
	binary.BigEndian.PutUint32(by, uint32(len(b.Bytes)))
	return append(by, b.Bytes...)
}

func (b *Bytes) Get() []byte {
	return b.Bytes
}

func (b *Bytes) Put(by []byte) *Bytes {
	b.Bytes = by
	return b
}
