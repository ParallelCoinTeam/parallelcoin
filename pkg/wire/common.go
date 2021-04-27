package wire

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

const (
	// MaxVarIntPayload is the maximum payload size for a variable length integer.
	MaxVarIntPayload = 9
	// binaryFreeListMaxItems is the number of buffers to keep in the free list to
	// use for binary serialization and deserialization.
	binaryFreeListMaxItems = 16384
)

var (
	// littleEndian is a convenience variable since binary.LittleEndian is quite
	// long.
	littleEndian = binary.LittleEndian
	// bigEndian is a convenience variable since binary.BigEndian is quite long.
	bigEndian = binary.BigEndian
)

// binaryFreeList defines a concurrent safe free list of byte slices (up to the
// maximum number defined by the binaryFreeListMaxItems constant) that have a
// cap of 8 (thus it supports up to a uint64). It is used to provide temporary
// buffers for deserializing primitive numbers to and from their binary encoding
// in order to greatly reduce the number of allocations required. For
// convenience, functions are provided for each of the primitive unsigned
// integers that automatically obtain a buffer from the free list, perform the
// necessary binary conversion, read from or write to the given io.Reader or
// io.Writer, and return the buffer to the free list.
type binaryFreeList chan []byte

// Borrow returns a byte slice from the free list with a length of 8.  A new buffer is allocated if there are not any available on the free list.
func (l binaryFreeList) Borrow() []byte {
	var buf []byte
	select {
	case buf = <-l:
	default:
		buf = make([]byte, 8)
	}
	return buf[:8]
}

// Return puts the provided byte slice back on the free list. The buffer MUST have been obtained via the Borrow function
// and therefore have a cap of 8.
func (l binaryFreeList) Return(buf []byte) {
	select {
	case l <- buf:
	default:
		// Let it go to the garbage collector.
	}
}

// Uint8 reads a single byte from the provided reader using a buffer from the free list and returns it as a uint8.
func (l binaryFreeList) Uint8(r io.Reader) (rv uint8, e error) {
	buf := l.Borrow()[:1]
	if _, e = io.ReadFull(r, buf); E.Chk(e) {
		l.Return(buf)
		return 0, e
	}
	rv = buf[0]
	l.Return(buf)
	return rv, nil
}

// Uint16 reads two bytes from the provided reader using a buffer from the free
// list, converts it to a number using the provided byte order, and returns the
// resulting uint16.
func (l binaryFreeList) Uint16(r io.Reader, byteOrder binary.ByteOrder) (rv uint16, e error) {
	buf := l.Borrow()[:2]
	if _, e = io.ReadFull(r, buf); E.Chk(e) {
		l.Return(buf)
		return
	}
	rv = byteOrder.Uint16(buf)
	l.Return(buf)
	return rv, nil
}

// Uint32 reads four bytes from the provided reader using a buffer from the free
// list, converts it to a number using the provided byte order, and returns the
// resulting uint32.
func (l binaryFreeList) Uint32(r io.Reader, byteOrder binary.ByteOrder) (rv uint32, e error) {
	buf := l.Borrow()[:4]
	if _, e = io.ReadFull(r, buf); E.Chk(e) {
		l.Return(buf)
		return 0, e
	}
	rv = byteOrder.Uint32(buf)
	l.Return(buf)
	return rv, nil
}

// Uint64 reads eight bytes from the provided reader using a buffer from the
// free list, converts it to a number using the provided byte order, and returns
// the resulting uint64.
func (l binaryFreeList) Uint64(r io.Reader, byteOrder binary.ByteOrder) (rv uint64, e error) {
	buf := l.Borrow()[:8]
	if _, e = io.ReadFull(r, buf); E.Chk(e) {
		l.Return(buf)
		return
	}
	rv = byteOrder.Uint64(buf)
	l.Return(buf)
	return rv, nil
}

// PutUint8 copies the provided uint8 into a buffer from the free list and writes the resulting byte to the given
// writer.
func (l binaryFreeList) PutUint8(w io.Writer, val uint8) (e error) {
	buf := l.Borrow()[:1]
	buf[0] = val
	_, e = w.Write(buf)
	l.Return(buf)
	return e
}

// PutUint16 serializes the provided uint16 using the given byte order into a buffer from the free list and writes the
// resulting two bytes to the given writer.
func (l binaryFreeList) PutUint16(w io.Writer, byteOrder binary.ByteOrder, val uint16) (e error) {
	buf := l.Borrow()[:2]
	byteOrder.PutUint16(buf, val)
	_, e = w.Write(buf)
	l.Return(buf)
	return
}

// PutUint32 serializes the provided uint32 using the given byte order into a buffer from the free list and writes the
// resulting four bytes to the given writer.
func (l binaryFreeList) PutUint32(w io.Writer, byteOrder binary.ByteOrder, val uint32) (e error) {
	buf := l.Borrow()[:4]
	byteOrder.PutUint32(buf, val)
	_, e = w.Write(buf)
	l.Return(buf)
	return
}

// PutUint64 serializes the provided uint64 using the given byte order into a buffer from the free list and writes the
// resulting eight bytes to the given writer.
func (l binaryFreeList) PutUint64(w io.Writer, byteOrder binary.ByteOrder, val uint64) (e error) {
	buf := l.Borrow()[:8]
	byteOrder.PutUint64(buf, val)
	_, e = w.Write(buf)
	l.Return(buf)
	return
}

// binarySerializer provides a free list of buffers to use for serializing and deserializing primitive integer values to
// and from io.Readers and io.Writers.
var binarySerializer binaryFreeList = make(chan []byte, binaryFreeListMaxItems)

// errNonCanonicalVarInt is the common format string used for non-canonically encoded variable length integer errors.
var errNonCanonicalVarInt = "non-canonical varint %x - discriminant %x must " +
	"encode a value greater than %x"

// uint32Time represents a unix timestamp encoded with a uint32. It is used as a way to signal the readElement function
// how to decode a timestamp into a Go time.Time since it is otherwise ambiguous.
type uint32Time time.Time

// int64Time represents a unix timestamp encoded with an int64. It is used as a way to signal the readElement function
// how to decode a timestamp into a Go time.Time since it is otherwise ambiguous.
type int64Time time.Time

// readElement reads the next sequence of bytes from r using little endian depending on the concrete type of element
// pointed to.
func readElement(r io.Reader, element interface{}) (e error) {
	// Attempt to read the element based on the concrete type via fast type assertions first.
	switch l := element.(type) {
	case *int32:
		var rv uint32
		if rv, e = binarySerializer.Uint32(r, littleEndian); E.Chk(e) {
			return
		}
		*l = int32(rv)
		return
	case *uint32:
		var rv uint32
		if rv, e = binarySerializer.Uint32(r, littleEndian); E.Chk(e) {
			return
		}
		*l = rv
		return
	case *int64:
		var rv uint64
		if rv, e = binarySerializer.Uint64(r, littleEndian); E.Chk(e) {
			return
		}
		*l = int64(rv)
		return
	case *uint64:
		var rv uint64
		if rv, e = binarySerializer.Uint64(r, littleEndian); E.Chk(e) {
			return
		}
		*l = rv
		return
	case *bool:
		var rv uint8
		if rv, e = binarySerializer.Uint8(r); E.Chk(e) {
			return
		}
		if rv == 0x00 {
			*l = false
		} else {
			*l = true
		}
		return nil
	// Unix timestamp encoded as a uint32.
	case *uint32Time:
		var rv uint32
		if rv, e = binarySerializer.Uint32(r, binary.LittleEndian); E.Chk(e) {
			return
		}
		*l = uint32Time(time.Unix(int64(rv), 0))
		return
	// Unix timestamp encoded as an int64.
	case *int64Time:
		var rv uint64
		if rv, e = binarySerializer.Uint64(r, binary.LittleEndian); E.Chk(e) {
			return
		}
		*l = int64Time(time.Unix(int64(rv), 0))
		return
	// Message header checksum.
	case *[4]byte:
		if _, e = io.ReadFull(r, l[:]); E.Chk(e) {
		}
		return
	// Message header command.
	case *[CommandSize]uint8:
		if _, e = io.ReadFull(r, l[:]); E.Chk(e) {
		}
		return
	// IP address.
	case *[16]byte:
		if _, e = io.ReadFull(r, l[:]); E.Chk(e) {
		}
		return
	case *chainhash.Hash:
		if _, e = io.ReadFull(r, l[:]); E.Chk(e) {
			return
		}
		return nil
	case *ServiceFlag:
		var rv uint64
		if rv, e = binarySerializer.Uint64(r, littleEndian); E.Chk(e) {
			return
		}
		*l = ServiceFlag(rv)
		return
	case *InvType:
		var rv uint32
		if rv, e = binarySerializer.Uint32(r, littleEndian); E.Chk(e) {
			return
		}
		*l = InvType(rv)
		return
	case *BitcoinNet:
		var rv uint32
		if rv, e = binarySerializer.Uint32(r, littleEndian); E.Chk(e) {
			return
		}
		*l = BitcoinNet(rv)
		return
	case *BloomUpdateType:
		var rv uint8
		if rv, e = binarySerializer.Uint8(r); E.Chk(e) {
			return
		}
		*l = BloomUpdateType(rv)
		return
	case *RejectCode:
		var rv uint8
		if rv, e = binarySerializer.Uint8(r); E.Chk(e) {
			return
		}
		*l = RejectCode(rv)
		return
	}
	// Fall back to the slower binary.Read if a fast path was not available above.
	return binary.Read(r, littleEndian, element)
}

// readElements reads multiple items from r.  It is equivalent to multiple calls to readElement.
func readElements(r io.Reader, elements ...interface{}) (e error) {
	for _, element := range elements {
		if e = readElement(r, element); E.Chk(e) {
			return
		}
	}
	return
}

// writeElement writes the little endian representation of element to w.
func writeElement(w io.Writer, element interface{}) (e error) {
	// Attempt to write the element based on the concrete type via fast type assertions first.
	switch el := element.(type) {
	case int32:
		if e = binarySerializer.PutUint32(w, littleEndian, uint32(el)); E.Chk(e) {
			return
		}
		return
	case uint32:
		if e = binarySerializer.PutUint32(w, littleEndian, el); E.Chk(e) {
		}
		return
	case int64:
		if e = binarySerializer.PutUint64(w, littleEndian, uint64(el)); E.Chk(e) {
		}
		return
	case uint64:
		if e = binarySerializer.PutUint64(w, littleEndian, el); E.Chk(e) {
		}
		return
	case bool:
		if el {
			e = binarySerializer.PutUint8(w, 0x01)
		} else {
			e = binarySerializer.PutUint8(w, 0x00)
		}
		if E.Chk(e) {
			return
		}
		return nil
	// Message header checksum.
	case [4]byte:
		if _, e = w.Write(el[:]); E.Chk(e) {
			return
		}
		return
	// Message header command.
	case [CommandSize]uint8:
		if _, e = w.Write(el[:]); E.Chk(e) {
		}
		return
	// IP address.
	case [16]byte:
		if _, e = w.Write(el[:]); E.Chk(e) {
		}
		return
	case *chainhash.Hash:
		if _, e = w.Write(el[:]); E.Chk(e) {
		}
		return
	case ServiceFlag:
		if e = binarySerializer.PutUint64(w, littleEndian, uint64(el)); E.Chk(e) {
		}
		return
	case InvType:
		if e = binarySerializer.PutUint32(w, littleEndian, uint32(el)); E.Chk(e) {
		}
		return
	case BitcoinNet:
		if e = binarySerializer.PutUint32(w, littleEndian, uint32(el)); E.Chk(e) {
		}
		return
	case BloomUpdateType:
		if e = binarySerializer.PutUint8(w, uint8(el)); E.Chk(e) {
		}
		return
	case RejectCode:
		if e = binarySerializer.PutUint8(w, uint8(el)); E.Chk(e) {
			return
		}
		return
	}
	// Fall back to the slower binary.Write if a fast path was not available above.
	return binary.Write(w, littleEndian, element)
}

// writeElements writes multiple items to w.  It is equivalent to multiple calls to writeElement.
func writeElements(w io.Writer, elements ...interface{}) (e error) {
	for _, element := range elements {
		if e = writeElement(w, element); E.Chk(e) {
			return
		}
	}
	return
}

// ReadVarInt reads a variable length integer from r and returns it as a uint64.
func ReadVarInt(r io.Reader, pver uint32) (rv uint64, e error) {
	var discriminant uint8
	if discriminant, e = binarySerializer.Uint8(r); E.Chk(e) {
		return
	}
	switch discriminant {
	case 0xff:
		var sv uint64
		if sv, e = binarySerializer.Uint64(r, littleEndian); E.Chk(e) {
			return
		}
		rv = sv
		// The encoding is not canonical if the value could have been encoded using fewer bytes.
		min := uint64(0x100000000)
		if rv < min {
			return 0, messageError(
				"ReadVarInt", fmt.Sprintf(
					errNonCanonicalVarInt, rv, discriminant, min,
				),
			)
		}
	case 0xfe:
		var sv uint32
		if sv, e = binarySerializer.Uint32(r, littleEndian); E.Chk(e) {
			return
		}
		rv = uint64(sv)
		// The encoding is not canonical if the value could have been encoded using fewer bytes.
		min := uint64(0x10000)
		if rv < min {
			return 0, messageError(
				"ReadVarInt", fmt.Sprintf(
					errNonCanonicalVarInt, rv, discriminant, min,
				),
			)
		}
	case 0xfd:
		var sv uint16
		if sv, e = binarySerializer.Uint16(r, littleEndian); E.Chk(e) {
			return
		}
		rv = uint64(sv)
		// The encoding is not canonical if the value could have been encoded using fewer bytes.
		min := uint64(0xfd)
		if rv < min {
			return 0, messageError(
				"ReadVarInt", fmt.Sprintf(
					errNonCanonicalVarInt, rv, discriminant, min,
				),
			)
		}
	default:
		rv = uint64(discriminant)
	}
	return
}

// WriteVarInt serializes val to w using a variable number of bytes depending on its value.
func WriteVarInt(w io.Writer, pver uint32, val uint64) (e error) {
	if val < 0xfd {
		return binarySerializer.PutUint8(w, uint8(val))
	}
	if val <= math.MaxUint16 {
		if e = binarySerializer.PutUint8(w, 0xfd); E.Chk(e) {
			return
		}
		return binarySerializer.PutUint16(w, littleEndian, uint16(val))
	}
	if val <= math.MaxUint32 {
		if e = binarySerializer.PutUint8(w, 0xfe); E.Chk(e) {
			return
		}
		return binarySerializer.PutUint32(w, littleEndian, uint32(val))
	}
	if e = binarySerializer.PutUint8(w, 0xff); E.Chk(e) {
		return
	}
	return binarySerializer.PutUint64(w, littleEndian, val)
}

// VarIntSerializeSize returns the number of bytes it would take to serialize val as a variable length integer.
func VarIntSerializeSize(val uint64) int {
	// The value is small enough to be represented by itself, so it's just 1 byte.
	if val < 0xfd {
		return 1
	}
	// Discriminant 1 byte plus 2 bytes for the uint16.
	if val <= math.MaxUint16 {
		return 3
	}
	// Discriminant 1 byte plus 4 bytes for the uint32.
	if val <= math.MaxUint32 {
		return 5
	}
	// Discriminant 1 byte plus 8 bytes for the uint64.
	return 9
}

// ReadVarString reads a variable length string from r and returns it as a Go string. A variable length string is
// encoded as a variable length integer containing the length of the string followed by the bytes that represent the
// string itself. An error is returned if the length is greater than the maximum block payload size since it helps
// protect against memory exhaustion attacks and forced panics through malformed messages.
func ReadVarString(r io.Reader, pver uint32) (s string, e error) {
	var count uint64
	if count, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	// Prevent variable length strings that are larger than the maximum message size. It would be possible to cause
	// memory exhaustion and panics without a sane upper bound on this count.
	if count > MaxMessagePayload {
		str := fmt.Sprintf(
			"variable length string is too long "+
				"[count %d, max %d]", count, MaxMessagePayload,
		)
		return "", messageError("ReadVarString", str)
	}
	buf := make([]byte, count)
	if _, e = io.ReadFull(r, buf); E.Chk(e) {
		return
	}
	return string(buf), e
}

// WriteVarString serializes str to w as a variable length integer containing the length of the string followed by the
// bytes that represent the string itself.
func WriteVarString(w io.Writer, pver uint32, str string) (e error) {
	if e = WriteVarInt(w, pver, uint64(len(str))); E.Chk(e) {
		return
	}
	_, e = w.Write([]byte(str))
	return
}

// ReadVarBytes reads a variable length byte array. A byte array is encoded as a varInt containing the length of the
// array followed by the bytes themselves. An error is returned if the length is greater than the passed maxAllowed
// parameter which helps protect against memory exhaustion attacks and forced panics through malformed messages. The
// fieldName parameter is only used for the error message so it provides more context in the error.
func ReadVarBytes(r io.Reader, pver uint32, maxAllowed uint32, fieldName string) (b []byte, e error) {
	var count uint64
	if count, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	// Prevent byte array larger than the max message size. It would be possible to cause memory exhaustion and panics
	// without a sane upper bound on this count.
	if count > uint64(maxAllowed) {
		str := fmt.Sprintf(
			"%s is larger than the max allowed size "+
				"[count %d, max %d]", fieldName, count, maxAllowed,
		)
		e = messageError("ReadVarBytes", str)
		return
	}
	b = make([]byte, count)
	if _, e = io.ReadFull(r, b); E.Chk(e) {
	}
	return
}

// WriteVarBytes serializes a variable length byte array to w as a varInt containing the number of bytes, followed by
// the bytes themselves.
func WriteVarBytes(w io.Writer, pver uint32, bytes []byte) (e error) {
	slen := uint64(len(bytes))
	if e = WriteVarInt(w, pver, slen); E.Chk(e) {
		return
	}
	_, e = w.Write(bytes)
	return
}

// randomUint64 returns a cryptographically random uint64 value. This unexported version takes a reader primarily to
// ensure the error paths can be properly tested by passing a fake reader in the tests.
func randomUint64(r io.Reader) (rv uint64, e error) {
	if rv, e = binarySerializer.Uint64(r, bigEndian); E.Chk(e) {
	}
	return
}

// RandomUint64 returns a cryptographically random uint64 value.
func RandomUint64() (rv uint64, e error) {
	return randomUint64(rand.Reader)
}
