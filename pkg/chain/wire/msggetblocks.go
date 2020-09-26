package wire

import (
	"fmt"
	"github.com/stalker-loki/app/slog"
	"io"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

// MaxBlockLocatorsPerMsg is the maximum number of block locator hashes allowed per message.
const MaxBlockLocatorsPerMsg = 500

// MsgGetBlocks implements the Message interface and represents a bitcoin getblocks message.  It is used to request a list of blocks starting after the last known hash in the slice of block locator hashes.  The list is returned via an inv message (MsgInv) and is limited by a specific hash to stop at or the maximum number of blocks per message, which is currently 500. the HashStop field to the hash at which to stop and use AddBlockLocatorHash to build up the list of block locator hashes. The algorithm for building the block locator hashes should be to add the hashes in reverse order until you reach the genesis block.  In order to keep the list of locator hashes to a reasonable number of entries, first add the most recent 10 block hashes, then double the step each loop iteration to decrease the number of hashes the further away from head and closer to the genesis block you get.
type MsgGetBlocks struct {
	ProtocolVersion    uint32
	BlockLocatorHashes []*chainhash.Hash
	HashStop           chainhash.Hash
}

// AddBlockLocatorHash adds a new block locator hash to the message.
func (msg *MsgGetBlocks) AddBlockLocatorHash(hash *chainhash.Hash) (err error) {
	if len(msg.BlockLocatorHashes)+1 > MaxBlockLocatorsPerMsg {
		str := fmt.Sprintf("too many block locator hashes for message [max %v]",
			MaxBlockLocatorsPerMsg)
		return messageError("MsgGetBlocks.AddBlockLocatorHash", str)
	}
	msg.BlockLocatorHashes = append(msg.BlockLocatorHashes, hash)
	return nil
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface implementation.
func (msg *MsgGetBlocks) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) (err error) {
	if err = readElement(r, &msg.ProtocolVersion); slog.Check(err) {
		return
	}
	// Read num block locator hashes and limit to max.
	var count uint64
	if count, err = ReadVarInt(r, pver); slog.Check(err) {
		return
	}
	if count > MaxBlockLocatorsPerMsg {
		if err = messageError("MsgGetBlocks.BtcDecode", fmt.Sprintf("too many block locator hashes for message "+
			"[count %v, max %v]", count, MaxBlockLocatorsPerMsg)); slog.Check(err) {
		}
		return
	}
	// Create a contiguous slice of hashes to deserialize into in order to reduce the number of allocations.
	locatorHashes := make([]chainhash.Hash, count)
	msg.BlockLocatorHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := &locatorHashes[i]
		if err = readElement(r, hash); slog.Check(err) {
			return
		}
		if err = msg.AddBlockLocatorHash(hash); slog.Check(err) {
		}
	}
	return readElement(r, &msg.HashStop)
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface implementation.
func (msg *MsgGetBlocks) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) (err error) {
	count := len(msg.BlockLocatorHashes)
	if count > MaxBlockLocatorsPerMsg {
		str := fmt.Sprintf("too many block locator hashes for message "+
			"[count %v, max %v]", count, MaxBlockLocatorsPerMsg)
		return messageError("MsgGetBlocks.BtcEncode", str)
	}
	if err = writeElement(w, msg.ProtocolVersion); slog.Check(err) {
		return
	}
	if err = WriteVarInt(w, pver, uint64(count)); slog.Check(err) {
		return
	}
	for _, hash := range msg.BlockLocatorHashes {
		if err = writeElement(w, hash); slog.Check(err) {
			return
		}
	}
	return writeElement(w, &msg.HashStop)
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgGetBlocks) Command() string {
	return CmdGetBlocks
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver.  This is part of the Message interface implementation.
func (msg *MsgGetBlocks) MaxPayloadLength(pver uint32) uint32 {
	// Protocol version 4 bytes + num hashes (varInt) + max block locator hashes + hash stop.
	return 4 + MaxVarIntPayload + (MaxBlockLocatorsPerMsg * chainhash.HashSize) + chainhash.HashSize
}

// NewMsgGetBlocks returns a new bitcoin getblocks message that conforms to the Message interface using the passed parameters and defaults for the remaining fields.
func NewMsgGetBlocks(hashStop *chainhash.Hash) *MsgGetBlocks {
	return &MsgGetBlocks{
		ProtocolVersion:    ProtocolVersion,
		BlockLocatorHashes: make([]*chainhash.Hash, 0, MaxBlockLocatorsPerMsg),
		HashStop:           *hashStop,
	}
}
