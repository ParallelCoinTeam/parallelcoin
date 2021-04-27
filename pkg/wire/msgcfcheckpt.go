package wire

import (
	"errors"
	"fmt"
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

const (
	// CFCheckptInterval is the gap (in number of blocks) between each filter header checkpoint.
	CFCheckptInterval = 1000
	// maxCFHeadersLen is the max number of filter headers we will attempt to decode.
	maxCFHeadersLen = 100000
)

// ErrInsaneCFHeaderCount signals that we were asked to decode an unreasonable number of cfilter headers.
var ErrInsaneCFHeaderCount = errors.New(
	"refusing to decode unreasonable number of filter headers",
)

// MsgCFCheckpt implements the Message interface and represents a bitcoin cfcheckpt message. It is used to deliver
// committed filter header information in response to a getcfcheckpt message (MsgGetCFCheckpt). See MsgGetCFCheckpt for
// details on requesting the headers.
type MsgCFCheckpt struct {
	FilterType    FilterType
	StopHash      chainhash.Hash
	FilterHeaders []*chainhash.Hash
}

// AddCFHeader adds a new committed filter header to the message.
func (msg *MsgCFCheckpt) AddCFHeader(header *chainhash.Hash) (e error) {
	if len(msg.FilterHeaders) == cap(msg.FilterHeaders) {
		str := fmt.Sprintf(
			"FilterHeaders has insufficient capacity for "+
				"additional header: len = %d", len(msg.FilterHeaders),
		)
		return messageError("MsgCFCheckpt.AddCFHeader", str)
	}
	msg.FilterHeaders = append(msg.FilterHeaders, header)
	return nil
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface
// implementation.
func (msg *MsgCFCheckpt) BtcDecode(r io.Reader, pver uint32, _ MessageEncoding) (e error) {
	// Read filter type
	if e = readElement(r, &msg.FilterType); E.Chk(e) {
		return
	}
	// Read stop hash
	if e = readElement(r, &msg.StopHash); E.Chk(e) {
		return
	}
	// Read number of filter headers
	var count uint64
	if count, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	// Refuse to decode an insane number of cfheaders.
	if count > maxCFHeadersLen {
		return ErrInsaneCFHeaderCount
	}
	// Create a contiguous slice of hashes to deserialize into in order to reduce
	// the number of allocations.
	msg.FilterHeaders = make([]*chainhash.Hash, count)
	for i := uint64(0); i < count; i++ {
		var cfh chainhash.Hash
		if e = readElement(r, &cfh); E.Chk(e) {
			return
		}
		msg.FilterHeaders[i] = &cfh
	}
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This
// is part of the Message interface implementation.
func (msg *MsgCFCheckpt) BtcEncode(w io.Writer, pver uint32, _ MessageEncoding) (e error) {
	// Write filter type
	if e = writeElement(w, msg.FilterType); E.Chk(e) {
		return
	}
	// Write stop hash
	if e = writeElement(w, msg.StopHash); E.Chk(e) {
		return
	}
	// Write length of FilterHeaders slice
	count := len(msg.FilterHeaders)
	e = WriteVarInt(w, pver, uint64(count))
	if E.Chk(e) {
		return
	}
	for _, cfh := range msg.FilterHeaders {
		if e = writeElement(w, cfh); E.Chk(e) {
			return
		}
	}
	return
}

// Deserialize decodes a filter header from r into the receiver using a format
// that is suitable for long-term storage such as a database. This function
// differs from BtcDecode in that BtcDecode decodes from the bitcoin wire
// protocol as it was sent across the network. The wire encoding can technically
// differ depending on the protocol version and doesn't even really need to
// match the format of a stored filter header at all. As of the time this
// comment was written, the encoded filter header is the same in both instances,
// but there is a distinct difference and separating the two allows the API to
// be flexible enough to deal with changes.
func (msg *MsgCFCheckpt) Deserialize(r io.Reader) (e error) {
	// At the current time, there is no difference between the wire encoding and the
	// stable long-term storage format. As a result, make use of BtcDecode.
	return msg.BtcDecode(r, 0, BaseEncoding)
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgCFCheckpt) Command() string {
	return CmdCFCheckpt
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver. This is part of the Message interface implementation.
func (msg *MsgCFCheckpt) MaxPayloadLength(pver uint32) uint32 {
	// Message size depends on the blockchain height, so return general limit for
	// all messages.
	return MaxMessagePayload
}

// NewMsgCFCheckpt returns a new bitcoin cfheaders message that conforms to the
// Message interface. See MsgCFCheckpt for details.
func NewMsgCFCheckpt(filterType FilterType, stopHash *chainhash.Hash, headersCount int,) *MsgCFCheckpt {
	return &MsgCFCheckpt{
		FilterType:    filterType,
		StopHash:      *stopHash,
		FilterHeaders: make([]*chainhash.Hash, 0, headersCount),
	}
}
