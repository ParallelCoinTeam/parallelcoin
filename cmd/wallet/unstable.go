package wallet

import (
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pkg/wtxmgr"
)

// UnstableAPI exposes unstable api in the wallet
type UnstableAPI struct {
	w *Wallet
}

// ExposeUnstableAPI exposes additional unstable public APIs for a Wallet. These APIs may be changed or removed at any
// time. Currently this type exists to ease the transation (particularly for the legacy JSON-RPC server) from using
// exported manager packages to a unified wallet package that exposes all functionality by itself. New code should not
// be written using this API.
func ExposeUnstableAPI(w *Wallet) UnstableAPI {
	return UnstableAPI{w}
}

// TxDetails calls wtxmgr.Store.TxDetails under a single database view transaction.
func (u UnstableAPI) TxDetails(txHash *chainhash.Hash) (details *wtxmgr.TxDetails, e error) {
	e = walletdb.View(
		u.w.db, func(dbtx walletdb.ReadTx) (e error) {
			txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
			details, e = u.w.TxStore.TxDetails(txmgrNs, txHash)
			return e
		},
	)
	return
}

// RangeTransactions calls wtxmgr.Store.RangeTransactions under a single database view transaction.
func (u UnstableAPI) RangeTransactions(begin, end int32, f func([]wtxmgr.TxDetails) (bool, error)) error {
	return walletdb.View(
		u.w.db, func(dbtx walletdb.ReadTx) (e error) {
			txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
			return u.w.TxStore.RangeTransactions(txmgrNs, begin, end, f)
		},
	)
}
