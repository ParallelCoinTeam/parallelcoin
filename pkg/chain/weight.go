package blockchain

import (
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util"
)

// TODO: this might be wrong stuff for non-segwit, being changed to remove
//  segwit as its complexity weighs the whole thing down when not in use

const (
	// MaxBlockWeight defines the maximum block weight, where "block weight" is interpreted as defined in BIP0141. A
	// block's weight is calculated as the sum of the of bytes in the existing transactions and header, plus the weight
	// of each byte within a transaction. The weight of a "base" byte is 4, while the weight of a witness byte is 1. As
	// a result, for a block to be valid, the BlockWeight MUST be less than, or equal to MaxBlockWeight.
	MaxBlockWeight = 4000000
	// MaxBlockBaseSize is the maximum number of bytes within a block which can be allocated to non-witness data.
	MaxBlockBaseSize = 1000000
	// MaxBlockSigOpsCost is the maximum number of signature operations allowed for a block. It is calculated via a
	// weighted algorithm which weights segregated witness sig ops lower than regular sig ops.
	MaxBlockSigOpsCost = 80000
	// MinTxOutputWeight is the minimum possible weight for a transaction output.
	MinTxOutputWeight = wire.MinTxOutPayload
	// MaxOutputsPerBlock is the maximum number of transaction outputs there can be in a block of max weight size.
	MaxOutputsPerBlock = MaxBlockWeight / MinTxOutputWeight
)

// GetBlockWeight computes the value of the weight metric for a given block. Currently the weight metric is simply the
// sum of the block's serialized size without any witness data scaled proportionally by the WitnessScaleFactor, and the
// block's serialized size including any witness data.
func GetBlockWeight(blk *util.Block) int64 {
	msgBlock := blk.MsgBlock()
	baseSize := msgBlock.SerializeSizeStripped()
	totalSize := msgBlock.SerializeSize()
	// (baseSize * 3) + totalSize
	return int64(baseSize + totalSize)
}

// GetTransactionWeight computes the value of the weight metric for a given transaction. Currently the weight metric is
// simply the sum of the transactions's serialized size without any witness data scaled proportionally by the
// WitnessScaleFactor, and the transaction's serialized size including any witness data.
func GetTransactionWeight(tx *util.Tx) int64 {
	msgTx := tx.MsgTx()
	baseSize := msgTx.SerializeSizeStripped()
	totalSize := msgTx.SerializeSize()
	// (baseSize * 3) + totalSize
	return int64(baseSize + totalSize)
}

// GetSigOpCost returns the unified sig op cost for the passed transaction
func GetSigOpCost(tx *util.Tx, isCoinBaseTx bool, utxoView *UtxoViewpoint, bip16 bool) (int, error) {
	numSigOps := CountSigOps(tx)
	if bip16 {
		numP2SHSigOps, err := CountP2SHSigOps(tx, isCoinBaseTx, utxoView)
		if err != nil {
			Error(err)
			return 0, nil
		}
		numSigOps += numP2SHSigOps
	}
	return numSigOps, nil
}
