package wire

import (
	"bytes"
	"fmt"
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// defaultTransactionAlloc is the default size used for the backing array for transactions. The transaction array will
// dynamically grow as needed, but this figure is intended to provide enough space for the number of transactions in the
// vast majority of blocks without needing to grow the backing array multiple times.
const defaultTransactionAlloc = 2048

// MaxBlocksPerMsg is the maximum number of blocks allowed per message.
const MaxBlocksPerMsg = 500

// MaxBlockPayload is the maximum bytes a block message can be in bytes. After
// Segregated Witness, the max block payload has been raised to 4MB.
const MaxBlockPayload = 4000000

// maxTxPerBlock is the maximum number of transactions that could possibly fit into a block.
const maxTxPerBlock = (MaxBlockPayload / minTxPayload) + 1

// TxLoc holds locator data for the offset and length of where a transaction is located within a Block data buffer.
type TxLoc struct {
	TxStart int
	TxLen   int
}

// Block implements the Message interface and represents a bitcoin block message. It is used to deliver block and
// transaction information in response to a getdata message (MsgGetData) for a given block hash.
type Block struct {
	Header       BlockHeader
	Transactions []*MsgTx
}

// AddTransaction adds a transaction to the message.
func (msg *Block) AddTransaction(tx *MsgTx) (e error) {
	msg.Transactions = append(msg.Transactions, tx)
	return nil
}

// ClearTransactions removes all transactions from the message.
func (msg *Block) ClearTransactions() {
	msg.Transactions = make([]*MsgTx, 0, defaultTransactionAlloc)
}

// BtcDecode decodes r using the bitcoin protocol encoding into the receiver. This is part of the Message interface
// implementation. See Deserialize for decoding blocks stored to disk, such as in a database, as opposed to decoding
// blocks from the wire.
func (msg *Block) BtcDecode(r io.Reader, pver uint32, enc MessageEncoding) (e error) {
	if e = readBlockHeader(r, pver, &msg.Header); E.Chk(e) {
		return
	}
	var txCount uint64
	if txCount, e = ReadVarInt(r, pver); E.Chk(e) {
		return
	}
	// Prevent more transactions than could possibly fit into a block. It would be possible to cause memory exhaustion
	// and panics without a sane upper bound on this count.
	if txCount > maxTxPerBlock {
		str := fmt.Sprintf(
			"too many transactions to fit into a block [count %d, max %d]",
			txCount, maxTxPerBlock,
		)
		return messageError("Block.BtcDecode", str)
	}
	msg.Transactions = make([]*MsgTx, 0, txCount)
	for i := uint64(0); i < txCount; i++ {
		tx := MsgTx{}
		if e = tx.BtcDecode(r, pver, enc); E.Chk(e) {
			return
		}
		msg.Transactions = append(msg.Transactions, &tx)
	}
	return
}

// Deserialize decodes a block from r into the receiver using a format that is
// suitable for long-term storage such as a database while respecting the
// Version field in the block.
//
// This function differs from BtcDecode in that BtcDecode decodes from the
// bitcoin wire protocol as it was sent across the network. The wire encoding
// can technically differ depending on the protocol version and doesn't even
// really need to match the format of a stored block at all.
//
// As of the time this comment was written, the encoded block is the same in
// both instances, but there is a distinct difference and separating the two
// allows the API to be flexible enough to deal with changes.
func (msg *Block) Deserialize(r io.Reader) (e error) {
	// At the current time, there is no difference between the wire encoding at
	// protocol version 0 and the stable long-term storage format. As a result, make
	// use of BtcDecode. Passing an encoding type of WitnessEncoding to BtcEncode
	// for the MessageEncoding parameter indicates that the transactions within the
	// block are expected to be serialized according to the new serialization
	// structure defined in BIP0141.
	return msg.BtcDecode(r, 0, BaseEncoding)
}

// DeserializeNoWitness decodes a block from r into the receiver similar to
// Deserialize, however DeserializeWitness strips all (if any) witness data from
// the transactions within the block before encoding them.
func (msg *Block) DeserializeNoWitness(r io.Reader) (e error) {
	return msg.BtcDecode(r, 0, BaseEncoding)
}

// DeserializeTxLoc decodes r in the same manner Deserialize does, but it takes
// a byte buffer instead of a generic reader and returns a slice containing the
// start and length of each transaction within the raw data that is being
// deserialized.
func (msg *Block) DeserializeTxLoc(r *bytes.Buffer) (txLocs []TxLoc, e error) {
	fullLen := r.Len()
	// At the current time, there is no difference between the wire encoding at protocol version 0 and the stable
	// long-term storage format. As a result, make use of existing wire protocol functions.
	if e = readBlockHeader(r, 0, &msg.Header); E.Chk(e) {
		return
	}
	var txCount uint64
	if txCount, e = ReadVarInt(r, 0); E.Chk(e) {
		return
	}
	// Prevent more transactions than could possibly fit into a block. It would be possible to cause memory exhaustion
	// and panics without a sane upper bound on this count.
	if txCount > maxTxPerBlock {
		str := fmt.Sprintf(
			"too many transactions to fit into a block [count %d, max %d]", txCount, maxTxPerBlock,
		)
		return nil, messageError("Block.DeserializeTxLoc", str)
	}
	// Deserialize each transaction while keeping track of its location within the byte stream.
	msg.Transactions = make([]*MsgTx, 0, txCount)
	txLocs = make([]TxLoc, txCount)
	for i := uint64(0); i < txCount; i++ {
		txLocs[i].TxStart = fullLen - r.Len()
		tx := MsgTx{}
		if e = tx.Deserialize(r); E.Chk(e) {
			return
		}
		msg.Transactions = append(msg.Transactions, &tx)
		txLocs[i].TxLen = (fullLen - r.Len()) - txLocs[i].TxStart
	}
	return
}

// BtcEncode encodes the receiver to w using the bitcoin protocol encoding. This is part of the Message interface
// implementation. See Serialize for encoding blocks to be stored to disk, such as in a database, as opposed to encoding
// blocks for the wire.
func (msg *Block) BtcEncode(w io.Writer, pver uint32, enc MessageEncoding) (e error) {
	if e = writeBlockHeader(w, pver, &msg.Header); E.Chk(e) {
		return
	}
	if e = WriteVarInt(w, pver, uint64(len(msg.Transactions))); E.Chk(e) {
		return
	}
	for _, tx := range msg.Transactions {
		if e = tx.BtcEncode(w, pver, enc); E.Chk(e) {
			return
		}
	}
	return
}

// Serialize encodes the block to w using a format that suitable for long-term storage such as a database while
// respecting the Version field in the block. This function differs from BtcEncode in that BtcEncode encodes the block
// to the bitcoin wire protocol in order to be sent across the network. The wire encoding can technically differ
// depending on the protocol version and doesn't even really need to match the format of a stored block at all. As of
// the time this comment was written, the encoded block is the same in both instances, but there is a distinct
// difference and separating the two allows the API to be flexible enough to deal with changes.
func (msg *Block) Serialize(w io.Writer) (e error) {
	// At the current time, there is no difference between the wire encoding at
	// protocol version 0 and the stable long-term storage format. As a result, make
	// use of BtcEncode. Passing WitnessEncoding as the encoding type here indicates
	// that each of the transactions should be serialized using the witness
	// serialization structure defined in BIP0141.
	return msg.BtcEncode(w, 0, BaseEncoding)
}

// SerializeNoWitness encodes a block to w using an identical format to
// Serialize, with all (if any) witness data stripped from all transactions.
// This method is provided in additon to the regular Serialize, in order to
// allow one to selectively encode transaction witness data to non-upgraded
// peers which are unaware of the new encoding.
func (msg *Block) SerializeNoWitness(w io.Writer) (e error) {
	return msg.BtcEncode(w, 0, BaseEncoding)
}

// SerializeSize returns the number of bytes it would take to serialize the
// block, factoring in any witness data within transaction.
func (msg *Block) SerializeSize() int {
	// Block header bytes + Serialized varint size for the number of transactions.
	n := blockHeaderLen + VarIntSerializeSize(uint64(len(msg.Transactions)))
	for _, tx := range msg.Transactions {
		n += tx.SerializeSize()
	}
	return n
}

// SerializeSizeStripped returns the number of bytes it would take to serialize
// the block, excluding any witness data (if any).
func (msg *Block) SerializeSizeStripped() int {
	// Block header bytes + Serialized varint size for the number of transactions.
	n := blockHeaderLen + VarIntSerializeSize(uint64(len(msg.Transactions)))
	for _, tx := range msg.Transactions {
		n += tx.SerializeSizeStripped()
	}
	return n
}

// Command returns the protocol command string for the message. This is part of the Message interface implementation.
func (msg *Block) Command() string {
	return CmdBlock
}

// MaxPayloadLength returns the maximum length the payload can be for the receiver. This is part of the Message
// interface implementation.
func (msg *Block) MaxPayloadLength(pver uint32) uint32 {
	// Block header at 80 bytes + transaction count + max transactions which can vary up to the MaxBlockPayload
	// (including the block header and transaction count).
	return MaxBlockPayload
}

// BlockHash computes the block identifier hash for this block.
func (msg *Block) BlockHash() chainhash.Hash {
	return msg.Header.BlockHash()
}

// BlockHashWithAlgos computes the block identifier hash for this block.
func (msg *Block) BlockHashWithAlgos(h int32) chainhash.Hash {
	return msg.Header.BlockHashWithAlgos(h)
}

// TxHashes returns a slice of hashes of all of transactions in this block.
func (msg *Block) TxHashes() ([]chainhash.Hash, error) {
	hashList := make([]chainhash.Hash, 0, len(msg.Transactions))
	for _, tx := range msg.Transactions {
		hashList = append(hashList, tx.TxHash())
	}
	return hashList, nil
}

// NewMsgBlock returns a new bitcoin block message that conforms to the Message interface.  See Block for details.
func NewMsgBlock(blockHeader *BlockHeader) *Block {
	return &Block{
		Header:       *blockHeader,
		Transactions: make([]*MsgTx, 0, defaultTransactionAlloc),
	}
}
