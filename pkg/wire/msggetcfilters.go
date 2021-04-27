package wire

import (
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// MaxGetCFiltersReqRange the maximum number of filters that may be requested in a getcfheaders message.
const MaxGetCFiltersReqRange = 1000

// MsgGetCFilters implements the Message interface and represents a bitcoin getcfilters message. It is used to request
// committed filters for a range of blocks.
type MsgGetCFilters struct {
	FilterType  FilterType
	StartHeight uint32
	StopHash    chainhash.Hash
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface
// implementation.
func (msg *MsgGetCFilters) BtcDecode(r io.Reader, pver uint32, _ MessageEncoding) (e error) {
	if e = readElement(r, &msg.FilterType); E.Chk(e) {
		return
	}
	if e = readElement(r, &msg.StartHeight); E.Chk(e) {
		return
	}
	return readElement(r, &msg.StopHash)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface
// implementation.
func (msg *MsgGetCFilters) BtcEncode(w io.Writer, pver uint32, _ MessageEncoding) (e error) {
	if e = writeElement(w, msg.FilterType); E.Chk(e) {
		return
	}
	if e = writeElement(w, &msg.StartHeight); E.Chk(e) {
		return
	}
	return writeElement(w, &msg.StopHash)
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgGetCFilters) Command() string {
	return CmdGetCFilters
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver. This is part of the Message
// interface implementation.
func (msg *MsgGetCFilters) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + 4 + chainhash.HashSize
}

// NewMsgGetCFilters returns a new bitcoin getcfilters message that conforms to the Message interface using the passed
// parameters and defaults for the remaining fields.
func NewMsgGetCFilters(
	filterType FilterType, startHeight uint32,
	stopHash *chainhash.Hash,
) *MsgGetCFilters {
	return &MsgGetCFilters{
		FilterType:  filterType,
		StartHeight: startHeight,
		StopHash:    *stopHash,
	}
}
