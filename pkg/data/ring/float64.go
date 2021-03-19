package ring

type BufferFloat64 struct {
	Buf    []float64
	Cursor int
	Full   bool
}

// NewBufferFloat64 creates a new ring buffer of float64 values
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

// Len returns the length of the buffer, which grows until it fills, after which
// this will always return the size of the buffer
func (b *BufferFloat64) Len() (length int) {
	if b.Full {
		return len(b.Buf)
	}
	return b.Cursor
}

// Add a new value to the cursor position of the ring buffer
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

// ForEach is an iterator that can be used to process every element in the
// buffer
func (b *BufferFloat64) ForEach(fn func(v float64) error) (e error) {
	c := b.Cursor
	i := c + 1
	if i == len(b.Buf) {
		// D.Ln("hit the end")
		i = 0
	}
	if !b.Full {
		// D.Ln("buffer not yet full")
		i = 0
	}
	// D.Ln(b.Buf)
	for ; ; i++ {
		if i == len(b.Buf) {
			// D.Ln("passed the end")
			i = 0
		}
		if i == c {
			// D.Ln("reached cursor again")
			break
		}
		// D.Ln(i, b.Cursor)
		if e = fn(b.Buf[i]); e != nil {
			break
		}
	}
	return
}
