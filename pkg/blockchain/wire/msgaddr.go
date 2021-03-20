package wire

import (
	"fmt"
	"io"
)

// MaxAddrPerMsg is the maximum number of addresses that can be in a single bitcoin addr message (MsgAddr).
const MaxAddrPerMsg = 1000

// MsgAddr implements the Message interface and represents a bitcoin addr message. It is used to provide a list of known
// active peers on the network. An active peer is considered one that has transmitted a message within the last 3 hours.
// Nodes which have not transmitted in that time frame should be forgotten. Each message is limited to a maximum number
// of addresses, which is currently 1000. As a result, multiple messages must be used to relay the full list. Use the
// AddAddress function to build up the list of known addresses when sending an addr message to another peer.
type MsgAddr struct {
	AddrList []*NetAddress
}

// AddAddress adds a known active peer to the message.
func (msg *MsgAddr) AddAddress(na *NetAddress) (e error) {
	if len(msg.AddrList)+1 > MaxAddrPerMsg {
		str := fmt.Sprintf(
			"too many addresses in message [max %v]",
			MaxAddrPerMsg,
		)
		return messageError("MsgAddr.AddAddress", str)
	}
	msg.AddrList = append(msg.AddrList, na)
	return nil
}

// AddAddresses adds multiple known active peers to the message.
func (msg *MsgAddr) AddAddresses(netAddrs ...*NetAddress) (e error) {
	for _, na := range netAddrs {
		e = msg.AddAddress(na)
		if E.Chk(e) {
			return
		}
	}
	return
}

// ClearAddresses removes all addresses from the message.
func (msg *MsgAddr) ClearAddresses() {
	msg.AddrList = []*NetAddress{}
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface
// implementation.
func (msg *MsgAddr) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) (e error) {
	var count uint64
	if count, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	// Limit to max addresses per message.
	if count > MaxAddrPerMsg {
		str := fmt.Sprintf(
			"too many addresses for message "+
				"[count %v, max %v]", count, MaxAddrPerMsg,
		)
		return messageError("MsgAddr.BtcDecode", str)
	}
	addrList := make([]NetAddress, count)
	msg.AddrList = make([]*NetAddress, 0, count)
	for i := uint64(0); i < count; i++ {
		na := &addrList[i]
		if e = readNetAddress(r, pver, na, true); E.Chk(e) {
			return
		}
		if e = msg.AddAddress(na); E.Chk(e) {
		}
	}
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface
// implementation.
func (msg *MsgAddr) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) (e error) {
	// Protocol versions before MultipleAddressVersion only allowed 1 address per message.
	count := len(msg.AddrList)
	if pver < MultipleAddressVersion && count > 1 {
		str := fmt.Sprintf(
			"too many addresses for message of "+
				"protocol version %v [count %v, max 1]", pver, count,
		)
		return messageError("MsgAddr.BtcEncode", str)
	}
	if count > MaxAddrPerMsg {
		str := fmt.Sprintf(
			"too many addresses for message "+
				"[count %v, max %v]", count, MaxAddrPerMsg,
		)
		return messageError("MsgAddr.BtcEncode", str)
	}
	if e = WriteVarInt(w, pver, uint64(count)); E.Chk(e) {
		return
	}
	for _, na := range msg.AddrList {
		if e = writeNetAddress(w, pver, na, true); E.Chk(e) {
			return
		}
	}
	return
}

// Command returns the protocol command string for the message. This is part of the Message interface implementation.
func (msg *MsgAddr) Command() string {
	return CmdAddr
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver. This is part of the Message
// interface implementation.
func (msg *MsgAddr) MaxPayloadLength(pver uint32) uint32 {
	if pver < MultipleAddressVersion {
		// Num addresses (varInt) + a single net addresses.
		return MaxVarIntPayload + maxNetAddressPayload(pver)
	}
	// Num addresses (varInt) + max allowed addresses.
	return MaxVarIntPayload + (MaxAddrPerMsg * maxNetAddressPayload(pver))
}

// NewMsgAddr returns a new bitcoin addr message that conforms to the Message interface. See MsgAddr for details.
func NewMsgAddr() *MsgAddr {
	return &MsgAddr{
		AddrList: make([]*NetAddress, 0, MaxAddrPerMsg),
	}
}
