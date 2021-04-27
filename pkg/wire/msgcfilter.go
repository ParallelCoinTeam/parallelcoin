package wire

import (
	"fmt"
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// FilterType is used to represent a filter type.
type FilterType uint8

const (
	// GCSFilterRegular is the regular filter type.
	GCSFilterRegular FilterType = iota
)
const (
	// MaxCFilterDataSize is the maximum byte size of a committed filter. The
	// maximum size is currently defined as 256KiB.
	MaxCFilterDataSize = 256 * 1024
)

// MsgCFilter implements the Message interface and represents a bitcoin cfilter
// message. It is used to deliver a committed filter in response to a
// getcfilters (MsgGetCFilters) message.
type MsgCFilter struct {
	FilterType FilterType
	BlockHash  chainhash.Hash
	Data       []byte
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCFilter) BtcDecode(r io.Reader, pver uint32, _ MessageEncoding) (e error) {
	// Read filter type
	if e = readElement(r, &msg.FilterType); E.Chk(e) {
		return
	}
	// Read the hash of the filter's block
	if e = readElement(r, &msg.BlockHash); E.Chk(e) {
		return e
	}
	// Read filter data
	msg.Data, e = ReadVarBytes(r, pver, MaxCFilterDataSize, "cfilter data")
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This
// is part of the Message interface implementation.
func (msg *MsgCFilter) BtcEncode(w io.Writer, pver uint32, _ MessageEncoding) (e error) {
	size := len(msg.Data)
	if size > MaxCFilterDataSize {
		str := fmt.Sprintf(
			"cfilter size too large for message "+
				"[size %v, max %v]", size, MaxCFilterDataSize,
		)
		return messageError("MsgCFilter.BtcEncode", str)
	}
	if e = writeElement(w, msg.FilterType); E.Chk(e) {
		return
	}
	if e = writeElement(w, msg.BlockHash); E.Chk(e) {
		return
	}
	return WriteVarBytes(w, pver, msg.Data)
}

// Deserialize decodes a filter from r into the receiver using a format that is
// suitable for long-term storage such as a database. This function differs from
// BtcDecode in that BtcDecode decodes from the bitcoin wire protocol as it was
// sent across the network. The wire encoding can technically differ depending
// on the protocol version and doesn't even really need to match the format of a
// stored filter at all. As of the time this comment was written, the encoded
// filter is the same in both instances, but there is a distinct difference and
// separating the two allows the API to be flexible enough to with changes.
func (msg *MsgCFilter) Deserialize(r io.Reader) (e error) {
	// At the current time, there is no difference between the wire encoding and the
	// stable long-term storage format. As a result, make use of BtcDecode.
	return msg.BtcDecode(r, 0, BaseEncoding)
}

// Command returns the protocol command string for the message. This is part of
// the Message interface implementation.
func (msg *MsgCFilter) Command() string {
	return CmdCFilter
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver. This is part of the Message interface implementation.
func (msg *MsgCFilter) MaxPayloadLength(pver uint32) uint32 {
	return uint32(VarIntSerializeSize(MaxCFilterDataSize)) + MaxCFilterDataSize + chainhash.HashSize + 1
}

// NewMsgCFilter returns a new bitcoin cfilter message that conforms to the
// Message interface. See MsgCFilter for details.
func NewMsgCFilter(
	filterType FilterType, blockHash *chainhash.Hash,
	data []byte,
) *MsgCFilter {
	return &MsgCFilter{
		FilterType: filterType,
		BlockHash:  *blockHash,
		Data:       data,
	}
}
