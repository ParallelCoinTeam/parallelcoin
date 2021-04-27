package block

import (
	"bytes"
	"fmt"
	"github.com/p9c/matrjoska/pkg/util"
	"io"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/wire"
)

// OutOfRangeError describes an error due to accessing an element that is out of range.
type OutOfRangeError string

// BlockHeightUnknown is the value returned for a block height that is unknown. This is typically because the block has
// not been inserted into the main chain yet.
const BlockHeightUnknown = int32(-1)

// Error satisfies the error interface and prints human-readable errors.
func (e OutOfRangeError) Error() string {
	return string(e)
}

// SetBlockBytes sets the internal serialized block byte buffer to the passed buffer. It is used to inject errors and is
// only available to the test package.
func (b *Block) SetBlockBytes(buf []byte) {
	b.serializedBlock = buf
}

// Block defines a bitcoin block that provides easier and more efficient manipulation of raw blocks. It also memorizes
// hashes for the block and its transactions on their first access so subsequent accesses don't have to repeat the
// relatively expensive hashing operations.
type Block struct {
	msgBlock                 *wire.Block     // Underlying WireBlock
	serializedBlock          []byte          // Serialized bytes for the block
	serializedBlockNoWitness []byte          // Serialized bytes for block w/o witness data
	blockHash                *chainhash.Hash // Cached block hash
	blockHeight              int32           // Height in the main block chain
	transactions             []*util.Tx      // Transactions
	txnsGenerated            bool            // ALL wrapped transactions generated
}

// WireBlock returns the underlying wire.Block for the Block.
func (b *Block) WireBlock() *wire.Block {
	// Return the cached block.
	return b.msgBlock
}

// Bytes returns the serialized bytes for the Block.
// This is equivalent to calling Serialize on the underlying wire.Block,
// however it caches the result so subsequent calls are more efficient.
func (b *Block) Bytes() ([]byte, error) {
	// Return the cached serialized bytes if it has already been generated.
	if len(b.serializedBlock) != 0 {
		return b.serializedBlock, nil
	}
	// Serialize the Block.
	w := bytes.NewBuffer(make([]byte, 0, b.msgBlock.SerializeSize()))
	e := b.msgBlock.Serialize(w)
	if e != nil {
		return nil, e
	}
	serializedBlock := w.Bytes()
	// Cache the serialized bytes and return them.
	b.serializedBlock = serializedBlock
	return serializedBlock, nil
}

// BytesNoWitness returns the serialized bytes for the block with transactions
// encoded without any witness data.
func (b *Block) BytesNoWitness() ([]byte, error) {
	// Return the cached serialized bytes if it has already been generated.
	if len(b.serializedBlockNoWitness) != 0 {
		return b.serializedBlockNoWitness, nil
	}
	// Serialize the Block.
	var w bytes.Buffer
	e := b.msgBlock.SerializeNoWitness(&w)
	if e != nil {
		return nil, e
	}
	serializedBlock := w.Bytes()
	// Cache the serialized bytes and return them.
	b.serializedBlockNoWitness = serializedBlock
	return serializedBlock, nil
}

// Hash returns the block identifier hash for the Block. This is equivalent to calling BlockHash on the underlying
// wire.Block, however it caches the result so subsequent calls are more efficient.
func (b *Block) Hash() *chainhash.Hash {
	// Return the cached block hash if it has already been generated.
	if b.blockHash != nil {
		return b.blockHash
	}
	// Cache the block hash and return it.
	hash := b.msgBlock.BlockHash()
	b.blockHash = &hash
	return &hash
}

// Tx returns a wrapped transaction (util.Tx) for the transaction at the specified index in the Block. The supplied
// index is 0 based. That is to say, the first transaction in the block is txNum 0. This is nearly equivalent to
// accessing the raw transaction (wire.MsgTx) from the underlying wire.Block, however the wrapped transaction has
// some helpful properties such as caching the hash so subsequent calls are more efficient.
func (b *Block) Tx(txNum int) (*util.Tx, error) {
	// Ensure the requested transaction is in range.
	numTx := uint64(len(b.msgBlock.Transactions))
	if txNum < 0 || uint64(txNum) > numTx {
		str := fmt.Sprintf(
			"transaction index %d is out of range - max %d",
			txNum, numTx-1,
		)
		return nil, OutOfRangeError(str)
	}
	// Generate slice to hold all of the wrapped transactions if needed.
	if len(b.transactions) == 0 {
		b.transactions = make([]*util.Tx, numTx)
	}
	// Return the wrapped transaction if it has already been generated.
	if b.transactions[txNum] != nil {
		return b.transactions[txNum], nil
	}
	// Generate and cache the wrapped transaction and return it.
	newTx := util.NewTx(b.msgBlock.Transactions[txNum])
	newTx.SetIndex(txNum)
	b.transactions[txNum] = newTx
	return newTx, nil
}

// Transactions returns a slice of wrapped transactions (util.Tx) for all transactions in the Block. This is nearly
// equivalent to accessing the raw transactions (wire.MsgTx) in the underlying wire.Block, however it instead
// provides easy access to wrapped versions (util.Tx) of them.
func (b *Block) Transactions() []*util.Tx {
	// Return transactions if they have ALL already been generated. This flag is necessary because the wrapped
	// transactions are lazily generated in a sparse fashion.
	if b.txnsGenerated {
		return b.transactions
	}
	// Generate slice to hold all of the wrapped transactions if needed.
	if len(b.transactions) == 0 {
		b.transactions = make([]*util.Tx, len(b.msgBlock.Transactions))
	}
	// Generate and cache the wrapped transactions for all that haven't already been done.
	for i, tx := range b.transactions {
		if tx == nil {
			newTx := util.NewTx(b.msgBlock.Transactions[i])
			newTx.SetIndex(i)
			b.transactions[i] = newTx
		}
	}
	b.txnsGenerated = true
	return b.transactions
}

// TxHash returns the hash for the requested transaction number in the Block. The supplied index is 0 based. That is to
// say, the first transaction in the block is txNum 0. This is equivalent to calling TxHash on the underlying
// wire.MsgTx, however it caches the result so subsequent calls are more efficient.
func (b *Block) TxHash(txNum int) (*chainhash.Hash, error) {
	// Attempt to get a wrapped transaction for the specified index. It will be created lazily if needed or simply
	// return the cached version if it has already been generated.
	tx, e := b.Tx(txNum)
	if e != nil {
		return nil, e
	}
	// Defer to the wrapped transaction which will return the cached hash if it has already been generated.
	return tx.Hash(), nil
}

// TxLoc returns the offsets and lengths of each transaction in a raw block. It is used to allow fast indexing into
// transactions within the raw byte stream.
func (b *Block) TxLoc() ([]wire.TxLoc, error) {
	rawMsg, e := b.Bytes()
	if e != nil {
		return nil, e
	}
	rbuf := bytes.NewBuffer(rawMsg)
	var mblock wire.Block
	txLocs, e := mblock.DeserializeTxLoc(rbuf)
	if e != nil {
		return nil, e
	}
	return txLocs, e
}

// Height returns the saved height of the block in the block chain. This value will be BlockHeightUnknown if it hasn't
// already explicitly been set.
func (b *Block) Height() int32 {
	return b.blockHeight
}

// SetHeight sets the height of the block in the block chain.
func (b *Block) SetHeight(height int32) {
	b.blockHeight = height
}

// NewBlock returns a new instance of a bitcoin block given an underlying wire.Block.  See Block.
func NewBlock(msgBlock *wire.Block) *Block {
	return &Block{
		msgBlock:    msgBlock,
		blockHeight: BlockHeightUnknown,
	}
}

// NewFromBytes returns a new instance of a bitcoin block given the serialized bytes.  See Block.
func NewFromBytes(serializedBlock []byte) (*Block, error) {
	br := bytes.NewReader(serializedBlock)
	b, e := NewFromReader(br)
	if e != nil {
		return nil, e
	}
	b.serializedBlock = serializedBlock
	return b, nil
}

// NewFromReader returns a new instance of a bitcoin block given a Reader to deserialize the block.  See Block.
func NewFromReader(r io.Reader) (*Block, error) {
	// Deserialize the bytes into a Block.
	var msgBlock wire.Block
	e := msgBlock.Deserialize(r)
	if e != nil {
		return nil, e
	}
	b := Block{
		msgBlock:    &msgBlock,
		blockHeight: BlockHeightUnknown,
	}
	return &b, nil
}

// NewFromBlockAndBytes returns a new instance of a bitcoin block given an underlying wire.Block and the
// serialized bytes for it. See Block.
func NewFromBlockAndBytes(msgBlock *wire.Block, serializedBlock []byte) *Block {
	return &Block{
		msgBlock:        msgBlock,
		serializedBlock: serializedBlock,
		blockHeight:     BlockHeightUnknown,
	}
}
