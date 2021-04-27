package wire

import (
	"fmt"
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// maxFlagsPerMerkleBlock is the maximum number of flag bytes that could possibly fit into a merkle block. Since each
// transaction is represented by a single bit, this is the max number of transactions per block divided by 8 bits per
// byte. Then an extra one to cover partials.
const maxFlagsPerMerkleBlock = maxTxPerBlock / 8

// MsgMerkleBlock implements the Message interface and represents a bitcoin merkleblock message which is used to reset a
// Bloom filter. This message was not added until protocol version BIP0037Version.
type MsgMerkleBlock struct {
	Header       BlockHeader
	Transactions uint32
	Hashes       []*chainhash.Hash
	Flags        []byte
}

// AddTxHash adds a new transaction hash to the message.
func (msg *MsgMerkleBlock) AddTxHash(hash *chainhash.Hash) (e error) {
	if len(msg.Hashes)+1 > maxTxPerBlock {
		str := fmt.Sprintf(
			"too many tx hashes for message [max %v]",
			maxTxPerBlock,
		)
		return messageError("MsgMerkleBlock.AddTxHash", str)
	}
	msg.Hashes = append(msg.Hashes, hash)
	return nil
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface
// implementation.
func (msg *MsgMerkleBlock) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) (e error) {
	if pver < BIP0037Version {
		str := fmt.Sprintf(
			"merkleblock message invalid for protocol "+
				"version %d", pver,
		)
		return messageError("MsgMerkleBlock.BtcDecode", str)
	}
	if e = readBlockHeader(r, pver, &msg.Header); E.Chk(e) {
		return
	}
	if e = readElement(r, &msg.Transactions); E.Chk(e) {
		return
	}
	// Read num block locator hashes and limit to max.
	var count uint64
	if count, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	if count > maxTxPerBlock {
		str := fmt.Sprintf(
			"too many transaction hashes for message "+
				"[count %v, max %v]", count, maxTxPerBlock,
		)
		return messageError("MsgMerkleBlock.BtcDecode", str)
	}
	// Create a contiguous slice of hashes to deserialize into in order to reduce the number of allocations.
	hashes := make([]chainhash.Hash, count)
	msg.Hashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := &hashes[i]
		if e = readElement(r, hash); E.Chk(e) {
			return
		}
		if e = msg.AddTxHash(hash); E.Chk(e) {
		}
	}
	msg.Flags, e = ReadVarBytes(
		r, pver, maxFlagsPerMerkleBlock,
		"merkle block flags size",
	)
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface
// implementation.
func (msg *MsgMerkleBlock) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) (e error) {
	if pver < BIP0037Version {
		str := fmt.Sprintf(
			"merkleblock message invalid for protocol "+
				"version %d", pver,
		)
		return messageError("MsgMerkleBlock.BtcEncode", str)
	}
	// Read num transaction hashes and limit to max.
	numHashes := len(msg.Hashes)
	if numHashes > maxTxPerBlock {
		str := fmt.Sprintf(
			"too many transaction hashes for message "+
				"[count %v, max %v]", numHashes, maxTxPerBlock,
		)
		return messageError("MsgMerkleBlock.BtcDecode", str)
	}
	numFlagBytes := len(msg.Flags)
	if numFlagBytes > maxFlagsPerMerkleBlock {
		str := fmt.Sprintf(
			"too many flag bytes for message [count %v, "+
				"max %v]", numFlagBytes, maxFlagsPerMerkleBlock,
		)
		return messageError("MsgMerkleBlock.BtcDecode", str)
	}
	if e = writeBlockHeader(w, pver, &msg.Header); E.Chk(e) {
		return
	}
	e = writeElement(w, msg.Transactions)
	if e != nil {
		return
	}
	e = WriteVarInt(w, pver, uint64(numHashes))
	if e != nil {
		return
	}
	for _, hash := range msg.Hashes {
		e = writeElement(w, hash)
		if e != nil {
			return
		}
	}
	return WriteVarBytes(w, pver, msg.Flags)
}

// Command returns the protocol command string for the message.  This is part of the Message interface implementation.
func (msg *MsgMerkleBlock) Command() string {
	return CmdMerkleBlock
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver. This is part of the Message
// interface implementation.
func (msg *MsgMerkleBlock) MaxPayloadLength(pver uint32) uint32 {
	return MaxBlockPayload
}

// NewMsgMerkleBlock returns a new bitcoin merkleblock message that conforms to the Message interface. See
// MsgMerkleBlock for details.
func NewMsgMerkleBlock(bh *BlockHeader) *MsgMerkleBlock {
	return &MsgMerkleBlock{
		Header:       *bh,
		Transactions: 0,
		Hashes:       make([]*chainhash.Hash, 0),
		Flags:        make([]byte, 0),
	}
}
