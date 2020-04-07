package ring

import (
	"context"
	"github.com/marusama/semaphore"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/logi"
)

type Entry struct {
	Sem     semaphore.Semaphore
	Buf     []*logi.Entry
	Cursor  int
	Full    bool
	Clicked int
	Buttons []gel.Button
	Hiders  []gel.Button
}

func NewEntry(size int) *Entry {
	return &Entry{
		Sem:     semaphore.New(1),
		Buf:     make([]*logi.Entry, size),
		Cursor:  0,
		Clicked: -1,
		Buttons: make([]gel.Button, size),
		Hiders:  make([]gel.Button, size),
	}
}

// Clear sets the buffer back to initial state
func (b *Entry) Clear() {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		b.Cursor = 0
		b.Full = false
	}
}

// Len returns the length of the buffer
func (b *Entry) Len() (out int) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		if b.Full {
			out = len(b.Buf)
		} else {
			out = b.Cursor
		}
	}
	return
}

// Get returns the value at the given index or nil if nothing
func (b *Entry) Get(i int) (out *logi.Entry) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		bl := len(b.Buf)
		cursor := i
		if i < bl {
			if b.Full {
				cursor = i + b.Cursor
				if cursor >= bl {
					cursor -= bl
				}
			}
			//Debug("get entry", i, "len", bl, "cursor", b.Cursor, "position",
			//	cursor)
			out = b.Buf[cursor]
		}
	}
	return
}

// GetButton returns the gel.Button of the entry
func (b *Entry) GetButton(i int) (out *gel.Button) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		bl := len(b.Buf)
		cursor := i
		if i < bl {
			if b.Full {
				cursor = i + b.Cursor
				if cursor >= bl {
					cursor -= bl
				}
			}
			//Debug("get entry", i, "len", bl, "cursor", b.Cursor, "position",
			//	cursor)
			out = &b.Buttons[cursor]
		}
	}
	return
}

// GetHider returns the gel.Button of the entry
func (b *Entry) GetHider(i int) (out *gel.Button) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		bl := len(b.Buf)
		cursor := i
		if i < bl {
			if b.Full {
				cursor = i + b.Cursor
				if cursor >= bl {
					cursor -= bl
				}
			}
			//Debug("get entry", i, "len", bl, "cursor", b.Cursor, "position",
			//	cursor)
			out = &b.Hiders[cursor]
		}
	}
	return
}

func (b *Entry) Add(value *logi.Entry) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
		defer b.Sem.Release(1)
		if b.Cursor == len(b.Buf) {
			b.Cursor = 0
			if !b.Full {
				b.Full = true
			}
		}
		b.Buf[b.Cursor] = value
		b.Cursor++
	}
}

func (b *Entry) ForEach(fn func(v *logi.Entry) error) (err error) {
	if err := b.Sem.Acquire(context.Background(), 1); !Check(err) {
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
		b.Sem.Release(1)
	}
	return
}
