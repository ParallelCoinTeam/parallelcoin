package String

import "encoding/binary"

// String plain old bytes. Maximum length from 32 bits int
type String struct {
	Bytes []byte
}

func New() *String {
	return &String{}
}

func (b *String) DecodeOne(by []byte) *String {
	b.Decode(by)
	return b
}

func (b *String) Decode(by []byte) (out []byte) {
	if len(by) >= 4 {
		length := binary.BigEndian.Uint32(by[:4])
		if len(by) >= 4+int(length) {
			out = by[4 : 4+length]
			b.Bytes = out
		}
	}
	return
}

func (b *String) Encode() []byte {
	by := make([]byte, 4)
	binary.BigEndian.PutUint32(by, uint32(len(b.Bytes)))
	return append(by, b.Bytes...)
}

func (b *String) Get() string {
	return string(b.Bytes)
}

func (b *String) Put(s string) *String {
	b.Bytes = []byte(s)
	return b
}
