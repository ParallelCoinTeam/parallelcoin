package wire

import (
	"github.com/stalker-loki/app/slog"
	"io"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

// MsgGetCFCheckpt is a request for filter headers at evenly spaced intervals throughout the blockchain history. It allows to set the FilterType field to get headers in the chain of basic (0x00) or extended (0x01) headers.
type MsgGetCFCheckpt struct {
	FilterType FilterType
	StopHash   chainhash.Hash
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) BtcDecode(r io.Reader, pver uint32, _ MessageEncoding) (err error) {
	if err = readElement(r, &msg.FilterType); slog.Check(err) {
		return
	}
	return readElement(r, &msg.StopHash)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) BtcEncode(w io.Writer, pver uint32, _ MessageEncoding) (err error) {
	if err = writeElement(w, msg.FilterType); slog.Check(err) {
		return
	}
	return writeElement(w, &msg.StopHash)
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) Command() string {
	return CmdGetCFCheckpt
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFCheckpt) MaxPayloadLength(pver uint32) uint32 {
	// Filter type + uint32 + block hash
	return 1 + chainhash.HashSize
}

// NewMsgGetCFCheckpt returns a new bitcoin getcfcheckpt message that conforms to the Message interface using the passed parameters and defaults for the remaining fields.
func NewMsgGetCFCheckpt(filterType FilterType, stopHash *chainhash.Hash) *MsgGetCFCheckpt {
	return &MsgGetCFCheckpt{
		FilterType: filterType,
		StopHash:   *stopHash,
	}
}
