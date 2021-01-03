package Byte

// Byte - plain old byte
type Byte struct {
	Byte byte
}

func New() *Byte {
	return &Byte{}
}

func (b *Byte) DecodeOne(by []byte) *Byte {
	b.Decode(by)
	return b
}

func (b *Byte) Decode(by []byte) (out []byte) {
	if len(by) >= 1 {
		out = by[0:]
		b.Byte = out[0]
	}
	return
}

func (b *Byte) Encode() []byte {
	return []byte{b.Byte}
}

func (b *Byte) Get() byte {
	return b.Byte
}

func (b *Byte) Put(by byte) *Byte {
	b.Byte = by
	return b
}
