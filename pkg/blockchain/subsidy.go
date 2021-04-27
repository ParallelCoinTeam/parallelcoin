package blockchain

import (
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/hardfork"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

// CreateHardForkSubsidyTx creates the transaction that must be on the hard fork activation block in place of a standard
// coinbase transaction.
//
// The main difference is the value set on this coinbase and that it pays out to multiple addresses, several being to
// the developers and to a 3 of 4 multisig to the development team for marketing and ongoing development costs multisig
// tx: NUM_SIGS PUBKEY PUBKEY PUBKEY... NUM_PUBKEYS OP_CHECKMULTISIG
func CreateHardForkSubsidyTx(
	params *chaincfg.Params,
	coinbaseScript []byte,
	nextBlockHeight int32,
	addr btcaddr.Address,
	version int32,
) (*util.Tx, error) {
	payees := hardfork.Payees
	if params.Net == wire.TestNet3 {
		payees = hardfork.TestnetPayees
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	tx.AddTxIn(
		&wire.TxIn{
			// Coinbase transactions have no inputs, so previous outpoint is zero hash and max index.
			PreviousOutPoint: *wire.NewOutPoint(
				&chainhash.Hash{},
				wire.MaxPrevOutIndex,
			),
			SignatureScript: coinbaseScript,
			Sequence:        wire.MaxTxInSequenceNum,
		},
	)
	for i := range payees {
		script, _ := txscript.PayToAddrScript(payees[i].Address)
		tx.AddTxOut(
			&wire.TxOut{
				Value:    int64(payees[i].Amount),
				PkScript: script,
			},
		)
	}
	// Add Core multisig payment
	builder := txscript.NewScriptBuilder()
	keylist := hardfork.CorePubkeyBytes
	if params.Net == wire.TestNet3 {
		keylist = hardfork.TestnetCorePubkeyBytes
	}
	builder.AddOp(txscript.OP_3).
		AddData(keylist[0]).
		AddData(keylist[1]).
		AddData(keylist[2]).
		AddData(keylist[3]).
		AddOp(txscript.OP_4).
		AddOp(txscript.OP_CHECKMULTISIG)
	script, _ := builder.Script()
	amt := hardfork.CoreAmount
	if params.Net == wire.TestNet3 {
		amt = hardfork.TestnetCoreAmount
	}
	tx.AddTxOut(
		&wire.TxOut{
			Value:    int64(amt),
			PkScript: script,
		},
	)
	// add miner's reward based on last non-hf reward
	script, _ = txscript.PayToAddrScript(addr)
	tx.AddTxOut(
		&wire.TxOut{
			Value:    CalcBlockSubsidy(nextBlockHeight+1, params, version),
			PkScript: script,
		},
	)
	return util.NewTx(tx), nil
}
