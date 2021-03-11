package templates

import (
	"errors"
	"github.com/niubaoshu/gotiny"
	chainhash "github.com/p9c/pod/pkg/blockchain/chainhash"
	"github.com/p9c/pod/pkg/blockchain/wire"
	"time"
)

type (
	// Diffs is a bundle of difficulty bits
	Diffs map[int32]uint32
	// Merkles is a bundle of merkle roots
	Merkles map[int32]chainhash.Hash
	// Txs is a set of transactions
	Txs map[int32][]*wire.MsgTx
)

// Message describes a message that a mining worker can use to
// construct a block to mine on.
type Message struct {
	UUID      uint64
	Height    int32
	PrevBlock chainhash.Hash
	Bits      Diffs
	Merkles   Merkles
	txs       Txs
	Timestamp time.Time
}

// SetTxs writes to the private, non-serialized transactions field
func (m *Message) SetTxs(ver int32, txs []*wire.MsgTx) {
	if m.txs == nil {
		m.txs = make(Txs)
	}
	m.txs[ver] = txs
}

// GetTxs returns the transactions
func (m *Message) GetTxs() (txs map[int32][]*wire.MsgTx) {
	return m.txs
}

// Serialize the Message for the wire
func (m *Message) Serialize() []byte {
	return gotiny.Marshal(m)
}

// DeserializeMsgBlockTemplate takes a message expected to be a Message
// and reconstitutes it
func DeserializeMsgBlockTemplate(b []byte) (m *Message) {
	m = &Message{}
	gotiny.Unmarshal(b, m)
	return
}

// GenBlockHeader generate a block given a version number to use for mining
// The nonce is empty, date can be updated, version changes merkle and target bits.
// All the data required for this is in the exported fields that are serialized
// for over the wire
func (m *Message) GenBlockHeader(vers int32) *wire.BlockHeader {
	return &wire.BlockHeader{
		Version:    vers,
		PrevBlock:  m.PrevBlock,
		MerkleRoot: m.Merkles[vers],
		Timestamp:  m.Timestamp,
		Bits:       m.Bits[vers],
	}
}

// Reconstruct takes a received block from the wire and reattaches the transactions
func (m *Message) Reconstruct(hdr *wire.BlockHeader) (mb *wire.MsgBlock, e error) {
	if hdr.PrevBlock != m.PrevBlock {
		e = errors.New("block is not for same parent block")
		return
	}
	mb = &wire.MsgBlock{Header: *hdr, Transactions: m.txs[hdr.Version]}
	return
}
