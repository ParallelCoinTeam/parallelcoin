package wire

import (
	"bytes"
	"reflect"
	"testing"
	
	"github.com/davecgh/go-spew/spew"
)

// TestFilterCLearLatest tests the MsgFilterClear API against the latest protocol version.
func TestFilterClearLatest(t *testing.T) {
	pver := ProtocolVersion
	msg := NewMsgFilterClear()
	// Ensure the command is expected value.
	wantCmd := "filterclear"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf(
			"NewMsgFilterClear: wrong command - got %v want %v",
			cmd, wantCmd,
		)
	}
	// Ensure max payload is expected value for latest protocol version.
	wantPayload := uint32(0)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf(
			"MaxPayloadLength: wrong max payload length for "+
				"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload,
		)
	}
}

// TestFilterClearCrossProtocol tests the MsgFilterClear API when encoding with the latest protocol version and decoding
// with BIP0031Version.
func TestFilterClearCrossProtocol(t *testing.T) {
	msg := NewMsgFilterClear()
	// Encode with latest protocol version.
	var buf bytes.Buffer
	e := msg.BtcEncode(&buf, ProtocolVersion, LatestEncoding)
	if e != nil {
		t.Errorf("encode of MsgFilterClear failed %v e <%v>", msg, e)
	}
	// Decode with old protocol version.
	var readmsg MsgFilterClear
	e = readmsg.BtcDecode(&buf, BIP0031Version, LatestEncoding)
	if e == nil {
		t.Errorf(
			"decode of MsgFilterClear succeeded when it "+
				"shouldn't have %v", msg,
		)
	}
}

// TestFilterClearWire tests the MsgFilterClear wire encode and decode for various protocol versions.
func TestFilterClearWire(t *testing.T) {
	msgFilterClear := NewMsgFilterClear()
	msgFilterClearEncoded := []byte{}
	tests := []struct {
		in   *MsgFilterClear // Message to encode
		out  *MsgFilterClear // Expected decoded message
		buf  []byte          // Wire encoding
		pver uint32          // Protocol version for wire encoding
		enc  MessageEncoding // Message encoding format
	}{
		// Latest protocol version.
		{
			msgFilterClear,
			msgFilterClear,
			msgFilterClearEncoded,
			ProtocolVersion,
			BaseEncoding,
		},
		// Protocol version BIP0037Version + 1.
		{
			msgFilterClear,
			msgFilterClear,
			msgFilterClearEncoded,
			BIP0037Version + 1,
			BaseEncoding,
		},
		// Protocol version BIP0037Version.
		{
			msgFilterClear,
			msgFilterClear,
			msgFilterClearEncoded,
			BIP0037Version,
			BaseEncoding,
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		e := test.in.BtcEncode(&buf, test.pver, test.enc)
		if e != nil {
			t.Errorf("BtcEncode #%d error %v", i, e)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf(
				"BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf),
			)
			continue
		}
		// Decode the message from wire format.
		var msg MsgFilterClear
		rbuf := bytes.NewReader(test.buf)
		e = msg.BtcDecode(rbuf, test.pver, test.enc)
		if e != nil {
			t.Errorf("BtcDecode #%d error %v", i, e)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf(
				"BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(msg), spew.Sdump(test.out),
			)
			continue
		}
	}
}

// TestFilterClearWireErrors performs negative tests against wire encode and decode of MsgFilterClear to confirm error
// paths work correctly.
func TestFilterClearWireErrors(t *testing.T) {
	pverNoFilterClear := BIP0037Version - 1
	wireErr := &MessageError{}
	baseFilterClear := NewMsgFilterClear()
	baseFilterClearEncoded := []byte{}
	tests := []struct {
		in       *MsgFilterClear // value to encode
		buf      []byte          // Wire encoding
		pver     uint32          // Protocol version for wire encoding
		enc      MessageEncoding // Message encoding format
		max      int             // Max size of fixed buffer to induce errors
		writeErr error           // Expected write error
		readErr  error           // Expected read error
	}{
		// Force error due to unsupported protocol version.
		{
			baseFilterClear, baseFilterClearEncoded,
			pverNoFilterClear, BaseEncoding, 4, wireErr, wireErr,
		},
	}
	t.Logf("Running %d tests", len(tests))
	var e error
	for i, test := range tests {
		// Encode to wire format.
		w := newFixedWriter(test.max)
		if e = test.in.BtcEncode(w, test.pver, test.enc); E.Chk(e) {
		}
		if reflect.TypeOf(e) != reflect.TypeOf(test.writeErr) {
			t.Errorf(
				"BtcEncode #%d wrong error got: %v, want: %v",
				i, e, test.writeErr,
			)
			continue
		}
		// For errors which are not of type MessageError, check them for equality.
		if _, ok := e.(*MessageError); !ok {
			if e != test.writeErr {
				t.Errorf(
					"BtcEncode #%d wrong error got: %v, want: %v", i, e, test.writeErr,
				)
				continue
			}
		}
		// Decode from wire format.
		var msg MsgFilterClear
		r := newFixedReader(test.max, test.buf)
		e = msg.BtcDecode(r, test.pver, test.enc)
		if reflect.TypeOf(e) != reflect.TypeOf(test.readErr) {
			t.Errorf(
				"BtcDecode #%d wrong error got: %v, want: %v",
				i, e, test.readErr,
			)
			continue
		}
		// For errors which are not of type MessageError, check them for equality.
		if _, ok := e.(*MessageError); !ok {
			if e != test.readErr {
				t.Errorf(
					"BtcDecode #%d wrong error got: %v, want: %v", i, e, test.readErr,
				)
				continue
			}
		}
	}
}
