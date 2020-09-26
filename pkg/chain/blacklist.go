package blockchain

import (
	"github.com/p9c/pod/pkg/chain/fork"
	script "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/util"
)

// ContainsBlacklisted returns true if one of the given addresses
// is found in the transaction
func ContainsBlacklisted(b *BlockChain, tx *util.Tx, blacklist []util.Address) (hasBlacklisted bool) {
	// in tests this function is not relevant
	if b == nil {
		return false
	}
	// blacklist only applies from hard fork
	if fork.GetCurrent(b.BestSnapshot().Height) < 1 {
		return false
	}
	var addrs []util.Address
	// first decode transaction and collect all addresses in the transaction outputs
	txo := tx.MsgTx().TxOut
	for i := range txo {
		scr := txo[i].PkScript
		_, a, _, _ := script.ExtractPkScriptAddrs(scr, b.params)
		addrs = append(addrs, a...)
	}
	// next get the addresses from the input transactions outpoints
	txi := tx.MsgTx().TxIn
	for i := range txi {
		bb, err := b.BlockByHash(&txi[i].PreviousOutPoint.Hash)
		if err == nil {
			txs := bb.MsgBlock().Transactions
			for j := range txs {
				txitxo := txs[j].TxOut
				for k := range txitxo {
					scr := txitxo[k].PkScript
					_, a, _, _ := script.ExtractPkScriptAddrs(scr, b.params)
					addrs = append(addrs, a...)
				}
			}
		}
	}
	// check if the set of addresses intersects with the blacklist
	return Intersects(addrs, blacklist)
}

// Intersects returns whether one slice of byte slices contains a match in
// another
func Intersects(a, b []util.Address) bool {
	for x := range a {
		for y := range b {
			if a[x].EncodeAddress() == b[y].EncodeAddress() {
				// If we find one match the two arrays intersect
				return true
			}
		}
	}
	return false
}
