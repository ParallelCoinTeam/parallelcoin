package ring

type BufferUint64 struct {
	Buf    []uint64
	Cursor int
	Full   bool
}

func NewBufferUint64(size int) *BufferUint64 {
	return &BufferUint64{
		Buf:    make([]uint64, size),
		Cursor: 0,
	}
}

func (b *BufferUint64) Add(value uint64) {
	b.Cursor++
	if b.Cursor > len(b.Buf)-1 {
		b.Cursor = 0
		if !b.Full {
			b.Full = true
		}
	}
	b.Buf[b.Cursor] = value
}

func (b *BufferUint64) ForEach(fn func(v uint64) error) (err error) {
	newest := b.Cursor + 1
	if newest >= len(b.Buf) {
		newest = 0
	}
	if !b.Full {
		newest = 0
	}
	for i := newest; i != b.Cursor; i++ {
		if i >= len(b.Buf) {
			i = 0
		}
		if err = fn(b.Buf[i]); err != nil {
			break
		}
	}
	return
}
