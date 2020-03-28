package ring

import (
	"context"
	"github.com/p9c/pod/pkg/logi"
)
import "github.com/marusama/semaphore"

type Entry struct {
	Sem    semaphore.Semaphore
	Buf    []*logi.Entry
	Cursor int
	Full   bool
}

func NewEntry(size int) *Entry {
	return &Entry{
		Sem:    semaphore.New(1),
		Buf:    make([]*logi.Entry, size),
		Cursor: -1,
	}
}

// Get returns the value at the given index or nil if nothing
func (b *Entry) Get(index int) (out *logi.Entry) {
	if err := b.Sem.Acquire(context.Background(), 1); !L.Check(err) {
		bl := len(b.Buf)
		if index < bl {
			cursor := b.Cursor + index
			if cursor > bl {
				cursor = cursor - bl
			}
			return b.Buf[cursor]
		}
		b.Sem.Release(1)
	}
	return
}

func (b *Entry) Add(value *logi.Entry) {
	if err := b.Sem.Acquire(context.Background(), 1); !L.Check(err) {
		b.Cursor++
		if b.Cursor == len(b.Buf) {
			b.Cursor = 0
			if !b.Full {
				b.Full = true
			}
		}
		b.Buf[b.Cursor] = value
		b.Sem.Release(1)
	}
}

func (b *Entry) ForEach(fn func(v *logi.Entry) error) (err error) {
	if err := b.Sem.Acquire(context.Background(), 1); !L.Check(err) {
		c := b.Cursor
		i := c + 1
		if i == len(b.Buf) {
			// L.Debug("hit the end")
			i = 0
		}
		if !b.Full {
			// L.Debug("buffer not yet full")
			i = 0
		}
		// L.Debug(b.Buf)
		for ; ; i++ {
			if i == len(b.Buf) {
				// L.Debug("passed the end")
				i = 0
			}
			if i == c {
				// L.Debug("reached cursor again")
				break
			}
			// L.Debug(i, b.Cursor)
			if err = fn(b.Buf[i]); err != nil {
				break
			}
		}
		b.Sem.Release(1)
	}
	return
}
