package wire

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

// MaxUserAgentLen is the maximum allowed length for the user agent field in a version message (MsgVersion).
const MaxUserAgentLen = 256

// DefaultUserAgent for wire in the stack
const DefaultUserAgent = "/btcwire:0.5.0/"

// MsgVersion implements the Message interface and represents a bitcoin version message. It is used for a peer to
// advertise itself as soon as an outbound connection is made. The remote peer then uses this information along with its
// own to negotiate. The remote peer must then respond with a version message of its own containing the negotiated
// values followed by a verack message (MsgVerAck). This exchange must take place before any further communication is
// allowed to proceed.
type MsgVersion struct {
	// Version of the protocol the node is using.
	ProtocolVersion int32
	// Bitfield which identifies the enabled services.
	Services ServiceFlag
	// Time the message was generated.  This is encoded as an int64 on the wire.
	Timestamp time.Time
	// Address of the remote peer.
	AddrYou NetAddress
	// Address of the local peer.
	AddrMe NetAddress
	// Unique value associated with message that is used to detect self connections.
	Nonce uint64
	// The user agent that generated message. This is a encoded as a varString on the wire. This has a max length of
	// MaxUserAgentLen.
	UserAgent string
	// Last block seen by the generator of the version message.
	LastBlock int32
	// Don't announce transactions to peer.
	DisableRelayTx bool
}

// HasService returns whether the specified service is supported by the peer that generated the message.
func (msg *MsgVersion) HasService(service ServiceFlag) bool {
	return msg.Services&service == service
}

// AddService adds service as a supported service by the peer generating the message.
func (msg *MsgVersion) AddService(service ServiceFlag) {
	msg.Services |= service
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. The version message is special in that the
// protocol version hasn't been negotiated yet. As a result, the pver field is ignored and any fields which are added in
// new versions are optional. This also mean that r must be a *bytes.Buffer so the number of remaining bytes can be
// ascertained. This is part of the Message interface implementation.
func (msg *MsgVersion) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) (e error) {
	buf, ok := r.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("MsgVersion.BtcDecode reader is not a *bytes.Buffer")
	}
	if e = readElements(buf, &msg.ProtocolVersion, &msg.Services, (*int64Time)(&msg.Timestamp)); E.Chk(e) {
		return
	}
	e = readNetAddress(buf, pver, &msg.AddrYou, false)
	if e != nil {
		return
	}
	// Protocol versions >= 106 added a from address, nonce, and user agent field and they are only considered present
	// if there are bytes remaining in the message.
	if buf.Len() > 0 {
		e = readNetAddress(buf, pver, &msg.AddrMe, false)
		if e != nil {
			return
		}
	}
	if buf.Len() > 0 {
		e = readElement(buf, &msg.Nonce)
		if e != nil {
			return
		}
	}
	if buf.Len() > 0 {
		var userAgent string
		if userAgent, e = ReadVarString(buf, pver); E.Chk(e) {
			return
		}
		if e = validateUserAgent(userAgent); E.Chk(e) {
			return
		}
		msg.UserAgent = userAgent
	}
	// Protocol versions >= 209 added a last known block field. It is only
	// considered present if there are bytes remaining in the message.
	if buf.Len() > 0 {
		if e = readElement(buf, &msg.LastBlock);E.Chk(e){
			return
		}
	}
	// There was no relay transactions field before BIP0037Version, but the default
	// behavior prior to the addition of the field was to always relay transactions.
	if buf.Len() > 0 {
		// It's safe to ignore the error here since the buffer has at least one byte and
		// that byte will result in a boolean value regardless of its value. Also, the
		// wire encoding for the field is true when transactions should be relayed, so
		// reverse it for the DisableRelayTx field.
		var relayTx bool
		if e = readElement(r, &relayTx); E.Chk(e) {
		}
		msg.DisableRelayTx = !relayTx
	}
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface
// implementation.
func (msg *MsgVersion) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) (e error) {
	if e = validateUserAgent(msg.UserAgent); E.Chk(e) {
		return
	}
	if e = writeElements(
		w, msg.ProtocolVersion, msg.Services,
		msg.Timestamp.Unix(),
	); E.Chk(e) {
		return
	}
	if e = writeNetAddress(w, pver, &msg.AddrYou, false); E.Chk(e) {
		return
	}
	if e = writeNetAddress(w, pver, &msg.AddrMe, false); E.Chk(e) {
		return
	}
	if e = writeElement(w, msg.Nonce); E.Chk(e) {
		return
	}
	if e = WriteVarString(w, pver, msg.UserAgent); E.Chk(e) {
		return
	}
	if e = writeElement(w, msg.LastBlock); E.Chk(e) {
		return
	}
	// There was no relay transactions field before BIP0037Version.  Also, the wire encoding for the field is true when
	// transactions should be relayed, so reverse it from the DisableRelayTx field.
	if pver >= BIP0037Version {
		if e = writeElement(w, !msg.DisableRelayTx); E.Chk(e) {
			return
		}
	}
	return
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgVersion) Command() string {
	return CmdVersion
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver.  This is part of the Message
// interface implementation.
func (msg *MsgVersion) MaxPayloadLength(pver uint32) uint32 {
	// XXX: <= 106 different
	// Protocol version 4 bytes + services 8 bytes + timestamp 8 bytes +
	// remote and local net addresses + nonce 8 bytes + length of user
	// agent (varInt) + max allowed useragent length + last block 4 bytes +
	// relay transactions flag 1 byte.
	return 33 + (maxNetAddressPayload(pver) * 2) + MaxVarIntPayload +
		MaxUserAgentLen
}

// NewMsgVersion returns a new bitcoin version message that conforms to the Message interface using the passed
// parameters and defaults for the remaining fields.
func NewMsgVersion(
	me *NetAddress, you *NetAddress, nonce uint64,
	lastBlock int32,
) *MsgVersion {
	// Limit the timestamp to one second precision since the protocol doesn't support better.
	return &MsgVersion{
		ProtocolVersion: int32(ProtocolVersion),
		Services:        0,
		Timestamp:       time.Unix(time.Now().Unix(), 0),
		AddrYou:         *you,
		AddrMe:          *me,
		Nonce:           nonce,
		UserAgent:       DefaultUserAgent,
		LastBlock:       lastBlock,
		DisableRelayTx:  false,
	}
}

// validateUserAgent checks userAgent length against MaxUserAgentLen
func validateUserAgent(userAgent string) (e error) {
	if len(userAgent) > MaxUserAgentLen {
		str := fmt.Sprintf(
			"user agent too long [len %v, max %v]",
			len(userAgent), MaxUserAgentLen,
		)
		return messageError("MsgVersion", str)
	}
	return nil
}

// AddUserAgent adds a user agent to the user agent string for the version message.  The version string is not defined
// to any strict format, although it is recommended to use the form "major.minor.revision" e.g. "2.6.41".
func (msg *MsgVersion) AddUserAgent(
	name string, version string,
	comments ...string,
) (e error) {
	newUserAgent := fmt.Sprintf("%s:%s", name, version)
	if len(comments) != 0 {
		newUserAgent = fmt.Sprintf(
			"%s(%s)", newUserAgent,
			strings.Join(comments, "; "),
		)
	}
	newUserAgent = fmt.Sprintf("%s%s/", msg.UserAgent, newUserAgent)
	if e = validateUserAgent(newUserAgent); E.Chk(e) {
		return
	}
	msg.UserAgent = newUserAgent
	return
}
