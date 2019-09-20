package sub

import (
	"encoding/binary"
	"fmt"
)

const (
	// _fatal = iota
	// _error
	// _warning
	// _info
	// _debug
	// _trace
	ipv4Format = "%d.%d.%d.%d:%d"
)

// EncodedAddrToString takes a string from bytes on the prefix containing an IPv4 address (it is passed around as a string for easy comparison) and returns the format used by net.Dial
func EncodedAddrToString(	encoded string) (out string) {
	if len(encoded) != 6 {
		return
	}
	e := []byte(encoded)
	out = fmt.Sprintf(ipv4Format,
		e[0], e[1], e[2], e[3],
		binary.LittleEndian.Uint16(e[4:6]),
	)
	return
}

// EncodeAddrToBytes takes a string in the format xxx.xxx.xxx.xxx:xxxxx and converts it to the encoded bytes format
func EncodeAddrToBytes(	addr string) (out []byte) {
	out = make([]byte, 4)
	o := make([]byte, 2)
	var ou16 uint16
	_, err := fmt.Sscanf(addr, ipv4Format, out[0], out[1], out[2], out[3], ou16)
	if err != nil {
		return []byte{}
	}
	binary.LittleEndian.PutUint16(o, ou16)
	out = append(out, o...)
	return
}
