package wire

import (
	"bytes"
	"testing"
)

func TestMemPool(t *testing.T) {
	pver := ProtocolVersion
	enc := BaseEncoding
	// Ensure the command is expected value.
	wantCmd := "mempool"
	msg := NewMsgMemPool()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgMemPool: wrong command - got %v want %v",
			cmd, wantCmd,
		)
	}
	// Ensure max payload is expected value.
	wantPayload := uint32(0)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload,
		)
	}
	// Test encode with latest protocol version.
	var buf bytes.Buffer
	e := msg.BtcEncode(&buf, pver, enc)
	if e != nil {
		t.Errorf("encode of MsgMemPool failed %v err <%v>", msg, e)
	}
	// Older protocol versions should fail encode since message didn't exist yet.
	oldPver := BIP0035Version - 1
	e = msg.BtcEncode(&buf, oldPver, enc)
	if e == nil {
		s := "encode of MsgMemPool passed for old protocol version %v err <%v>"
		t.Errorf(s, msg, e)
	}
	// Test decode with latest protocol version.
	readmsg := NewMsgMemPool()
	e = readmsg.BtcDecode(&buf, pver, enc)
	if e != nil {
		t.Errorf("decode of MsgMemPool failed [%v] err <%v>", buf, e)
	}
	// Older protocol versions should fail decode since message didn't exist yet.
	e = readmsg.BtcDecode(&buf, oldPver, enc)
	if e == nil {
		s := "decode of MsgMemPool passed for old protocol version %v err <%v>"
		t.Errorf(s, msg, e)
	}
}
