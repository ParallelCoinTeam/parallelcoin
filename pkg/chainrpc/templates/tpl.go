package templates

import (
	"errors"
	"time"

	"github.com/niubaoshu/gotiny"

	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/wire"
)

type (
	// Diffs is a bundle of difficulty bits
	Diffs map[int32]uint32
	// Merkles is a bundle of merkle roots
	Merkles map[int32]chainhash.Hash
	// Txs is a set of transactions
	Txs map[int32][]*wire.MsgTx
)

// Message describes a message that a mining worker can use to construct a block
// to mine on, as well the extra data required to reconstruct any version block
// mined by the miner
type Message struct {
	Nonce     uint64
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
func (m *Message) Reconstruct(hdr *wire.BlockHeader) (mb *wire.Block, e error) {
	if hdr.PrevBlock != m.PrevBlock {
		e = errors.New("block is not for same parent block")
		return
	}
	mb = &wire.Block{Header: *hdr, Transactions: m.txs[hdr.Version]}
	return
}

// RecentMessages keeps a buffer of four previously created messages so that
// solutions found after a new template is created can be submitted
type RecentMessages struct {
	msgs   [4]*Message
	cursor int
}

func NewRecentMessages() *RecentMessages {
	return &RecentMessages{
		msgs:   [4]*Message{},
		cursor: 0,
	}
}

// Add a message to the RecentMessages, we just write to the current cursor
// position, and then advance it, back to zero if it exceeds the buffer length,
// overwriting the first, and so on
func (rm *RecentMessages) Add(msg *Message) {
	D.Ln("adding template with cursor", rm.cursor)
	rm.msgs[rm.cursor] = msg
	rm.cursor++
	if rm.cursor >= 4 {
		rm.cursor = 0
	}
}

// Len returns the number of elements in the buffer
func (rm *RecentMessages) Len() (o int) {
	for i := range rm.msgs {
		if rm.msgs[i] != nil {
			o++
		}
	}
	return
}

// Find checks whether the given nonce matches any of the cached Message's and remove it from the list
func (rm *RecentMessages) Find(nonce uint64) *Message {
	for i := range rm.msgs {
		if rm.msgs[i] != nil {
			I.Ln("recent message", i, rm.msgs[i].Nonce, nonce)
			if rm.msgs[i].Nonce == nonce {
				I.Ln("found message", nonce)
				msg := rm.msgs[i]
				rm.msgs[i] = nil
				rm.cursor = i
				return msg
			}
		}
	}
	return nil
}
