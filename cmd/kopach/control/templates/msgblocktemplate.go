package templates

import (
	"errors"
	"github.com/niubaoshu/gotiny"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"time"
)

type (
	// Diffs is a bundle of difficulty bits
	Diffs map[int32]uint32
	// Merkles is a bundle of merkle roots
	Merkles map[int32]chainhash.Hash
	// Coinbases is a bundle of coinbase transactions
	Coinbases map[int32]*wire.MsgTx
	// Txs is a set of transactions
	Txs []*wire.MsgTx
)

// Message describes a message that a mining worker can use to
// construct a block to mine on.
type Message struct {
	Nonce     uint64
	Height    int32
	PrevBlock chainhash.Hash
	Bits      Diffs
	Merkles   Merkles
	coinbases Coinbases
	txs       Txs
	Timestamp time.Time
}

// ResetCoinbases clears to the private, non-serialized coinbases field
func (m *Message) ResetCoinbases() {
	m.coinbases = make(Coinbases)
}

// GetCoinbases clears to the private, non-serialized coinbases field
func (m *Message) GetCoinbases() Coinbases {
	return m.coinbases
}

func (m *Message) SetCoinbase(ver int32, tx *wire.MsgTx) {
	m.coinbases[ver] = tx
}

func (m *Message) GetCoinbase(ver int32) (tx *wire.MsgTx) {
	return m.coinbases[ver]
}

// SetTxs writes to the private, non-serialized transactions field
func (m *Message) SetTxs(txs []*wire.MsgTx) {
	m.txs = txs
}

// GetTxs returns the transactions
func (m *Message) GetTxs() (txs []*wire.MsgTx) {
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
func (m *Message) Reconstruct(hdr *wire.BlockHeader) (mb *wire.MsgBlock, err error) {
	if hdr.PrevBlock != m.PrevBlock {
		err = errors.New("block is not for same parent block")
		Error(err)
		return
	}
	mb = &wire.MsgBlock{Header: *hdr}
	// the coinbase is the last transaction
	txs := append(m.txs, m.coinbases[mb.Header.Version])
	for _, tx := range txs {
		if err = mb.AddTransaction(tx); Check(err) {
			return
		}
	}
	return
}
