package transport

//
// type (
// 	Buf     [][]byte
// 	BufIter struct {
// 		Buf
// 		cursor int
// 	}
// )
//
// // NewBufIter returns a new BufIter loaded with a given slice of Buf
// func NewBufIter(buf Buf) *BufIter {
// 	return &BufIter{Buf: buf}
// }
//
// // Len returns the length of the buffer
// func (b *BufIter) Len() int {
// 	return len(b.Buf)
// }
//
// // Get returns the currently selected Buf
// func (b *BufIter) Get() []byte {
// 	return b.Buf[b.cursor]
// }
//
// // More returns true if the cursor is not at the end
// func (b *BufIter) More() bool {
// 	return b.cursor < len(b.Buf)
// }
//
// // Reset sets a buf iter to zero, not necessary to use first time iterating
// func (b *BufIter) Reset() {
// 	b.cursor = 0
// }
//
// // Next returns the next item in a buffer
// func (b *BufIter) Next() (buf []byte) {
// 	if b.cursor > len(b.Buf) {
// 	} else {
// 		buf = b.Buf[b.cursor]
// 		b.cursor++
// 	}
// 	return
// }
