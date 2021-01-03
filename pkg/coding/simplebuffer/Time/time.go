package Time

import (
	"encoding/binary"
	"time"
)

// Time stores a 64 bit time stamp
type Time struct {
	Bytes [8]byte
}

// New creates a new Time.Time
func New() *Time {
	return &Time{}
}

// DecodeOne decodes just the first element it finds in the slice
func (b *Time) DecodeOne(by []byte) *Time {
	b.Decode(by)
	return b
}

// Decode decodes the next element and returns the remainder
func (b *Time) Decode(by []byte) (out []byte) {
	if len(by) >= 4 {
		b.Bytes = [8]byte{by[0], by[1], by[2], by[3], by[4], by[5], by[6], by[7]}
		if len(by) > 8 {
			out = by[:8]
		}
	}
	return
}

// Encode the Time.Time to bytes
func (b *Time) Encode() []byte {
	return b.Bytes[:]
}

// Get returns the decoded form of the Time.Time
func (b *Time) Get() time.Time {
	t := binary.BigEndian.Uint64(b.Bytes[:8])
	return time.Unix(0, int64(t))
}

// Put stores the provided time.Time into a Time.Time
func (b *Time) Put(t time.Time) *Time {
	binary.BigEndian.PutUint64(b.Bytes[:], uint64(t.UnixNano()))
	return b
}
