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
	i := b.Cursor + 1
	if i > len(b.Buf)-1 {
		i = 0
	}
	if !b.Full {
		i = 0
	}
	// log.DEBUG(b.Buf)
	for ; ; i++ {
		if i > len(b.Buf)-1 {
			i = 0
		}
		if i == b.Cursor {
			break
		}
		// log.DEBUG(i, b.Cursor)
		if err = fn(b.Buf[i]); err != nil {
			break
		}
	}
	return
}
