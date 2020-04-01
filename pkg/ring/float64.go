package ring

type BufferFloat64 struct {
	Buf    []float64
	Cursor int
	Full   bool
}

func NewBufferFloat64(size int) *BufferFloat64 {
	return &BufferFloat64{
		Buf:    make([]float64, size),
		Cursor: -1,
	}
}

// Get returns the value at the given index or nil if nothing
func (b *BufferFloat64) Get(index int) (out *float64) {
	bl := len(b.Buf)
	if index < bl {
		cursor := b.Cursor + index
		if cursor > bl {
			cursor = cursor - bl
		}
		return &b.Buf[cursor]
	}
	return
}
func (b *BufferFloat64) Len() (length int) {
	if b.Full {
		return len(b.Buf)
	}
	return b.Cursor
}

func (b *BufferFloat64) Add(value float64) {
	b.Cursor++
	if b.Cursor == len(b.Buf) {
		b.Cursor = 0
		if !b.Full {
			b.Full = true
		}
	}
	b.Buf[b.Cursor] = value
}

func (b *BufferFloat64) ForEach(fn func(v float64) error) (err error) {
	c := b.Cursor
	i := c + 1
	if i == len(b.Buf) {
		// Debug("hit the end")
		i = 0
	}
	if !b.Full {
		// Debug("buffer not yet full")
		i = 0
	}
	// Debug(b.Buf)
	for ; ; i++ {
		if i == len(b.Buf) {
			// Debug("passed the end")
			i = 0
		}
		if i == c {
			// Debug("reached cursor again")
			break
		}
		// Debug(i, b.Cursor)
		if err = fn(b.Buf[i]); err != nil {
			break
		}
	}
	return
}
