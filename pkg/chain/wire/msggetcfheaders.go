package wire

import (
	"github.com/stalker-loki/app/slog"
	"io"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

// MsgGetCFHeaders is a message similar to MsgGetHeaders, but for committed filter headers. It allows to set the FilterType field to get headers in the chain of basic (0x00) or extended (0x01) headers.
type MsgGetCFHeaders struct {
	FilterType  FilterType
	StartHeight uint32
	StopHash    chainhash.Hash
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) BtcDecode(r io.Reader, pver uint32, _ MessageEncoding) (err error) {
	if err = readElement(r, &msg.FilterType); slog.Check(err) {
		return
	}
	if err = readElement(r, &msg.StartHeight); slog.Check(err) {
		return
	}
	err = readElement(r, &msg.StopHash)
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) BtcEncode(w io.Writer, pver uint32, _ MessageEncoding) (err error) {
	if err = writeElement(w, msg.FilterType); slog.Check(err) {
		return
	}
	if err = writeElement(w, &msg.StartHeight); slog.Check(err) {
		return
	}
	err = writeElement(w, &msg.StopHash)
	return
}

// Command returns the protocol command string for the message.
// This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) Command() string {
	return CmdGetCFHeaders
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + 4 + chainhash.HashSize
}

// NewMsgGetCFHeaders returns a new bitcoin getcfheader message that conforms to the Message interface using the passed
// parameters and defaults for the remaining fields.
func NewMsgGetCFHeaders(filterType FilterType, startHeight uint32, stopHash *chainhash.Hash) *MsgGetCFHeaders {
	return &MsgGetCFHeaders{
		FilterType:  filterType,
		StartHeight: startHeight,
		StopHash:    *stopHash,
	}
}
