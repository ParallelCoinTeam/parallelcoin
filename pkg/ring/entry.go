package ring

import (
	"context"
	"github.com/p9c/log"
	
	"github.com/marusama/semaphore"
)

type Entry struct {
	Sem     semaphore.Semaphore
	Buf     []*log.Entry
	Cursor  int
	Full    bool
	Clicked int
	// Buttons []gel.Button
	// Hiders  []gel.Button
}

// NewEntry creates a new entry ring buffer
func NewEntry(size int) *Entry {
	return &Entry{
		Sem:     semaphore.New(1),
		Buf:     make([]*log.Entry, size),
		Cursor:  0,
		Clicked: -1,
		// Buttons: make([]gel.Button, size),
		// Hiders:  make([]gel.Button, size),
	}
}

// Clear sets the buffer back to initial state
func (b *Entry) Clear() {
	var e error
	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
		defer b.Sem.Release(1)
		b.Cursor = 0
		b.Clicked = -1
		b.Full = false
	}
}

// Len returns the length of the buffer
func (b *Entry) Len() (out int) {
	var e error
	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
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
func (b *Entry) Get(i int) (out *log.Entry) {
	var e error
	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
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
			// D.Ln("get entry", i, "len", bl, "cursor", b.Cursor, "position",
			//	cursor)
			out = b.Buf[cursor]
		}
	}
	return
}

//
// // GetButton returns the gel.Button of the entry
// func (b *Entry) GetButton(i int) (out *gel.Button) {
// 	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
// 		defer b.Sem.Release(1)
// 		bl := len(b.Buf)
// 		cursor := i
// 		if i < bl {
// 			if b.Full {
// 				cursor = i + b.Cursor
// 				if cursor >= bl {
// 					cursor -= bl
// 				}
// 			}
// 			// D.Ln("get entry", i, "len", bl, "cursor", b.Cursor, "position",
// 			//	cursor)
// 			out = &b.Buttons[cursor]
// 		}
// 	}
// 	return
// }
//
// // GetHider returns the gel.Button of the entry
// func (b *Entry) GetHider(i int) (out *gel.Button) {
// 	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
// 		defer b.Sem.Release(1)
// 		bl := len(b.Buf)
// 		cursor := i
// 		if i < bl {
// 			if b.Full {
// 				cursor = i + b.Cursor
// 				if cursor >= bl {
// 					cursor -= bl
// 				}
// 			}
// 			// D.Ln("get entry", i, "len", bl, "cursor", b.Cursor, "position",
// 			//	cursor)
// 			out = &b.Hiders[cursor]
// 		}
// 	}
// 	return
// }

func (b *Entry) Add(value *log.Entry) {
	var e error
	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
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

func (b *Entry) ForEach(fn func(v *log.Entry) error) (e error) {
	if e = b.Sem.Acquire(context.Background(), 1); !E.Chk(e) {
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
			if e = fn(b.Buf[i]); E.Chk(e) {
				break
			}
		}
		b.Sem.Release(1)
	}
	return
}
