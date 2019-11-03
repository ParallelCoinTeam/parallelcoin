package wallet

import (
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

// UnstableAPI exposes unstable api in the wallet
type UnstableAPI struct {
	w *Wallet
}

// ExposeUnstableAPI exposes additional unstable public APIs for a Wallet.  These APIs
// may be changed or removed at any time.  Currently this type exists to ease
// the transation (particularly for the legacy JSON-RPC server) from using
// exported manager packages to a unified wallet package that exposes all
// functionality by itself.  New code should not be written using this API.
func ExposeUnstableAPI(w *Wallet) UnstableAPI {
	return UnstableAPI{w}
}

// TxDetails calls wtxmgr.Store.TxDetails under a single database view transaction.
func (u UnstableAPI) TxDetails(txHash *chainhash.Hash) (*wtxmgr.TxDetails, error) {
	var details *wtxmgr.TxDetails
	err := walletdb.View(u.w.db, func(dbtx walletdb.ReadTx) error {
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
		var err error
		details, err = u.w.TxStore.TxDetails(txmgrNs, txHash)
		return err
	})
	return details, err
}

// RangeTransactions calls wtxmgr.Store.RangeTransactions under a single
// database view tranasction.
func (u UnstableAPI) RangeTransactions(begin, end int32, f func([]wtxmgr.TxDetails) (bool, error)) error {
	return walletdb.View(u.w.db, func(dbtx walletdb.ReadTx) error {
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
		return u.w.TxStore.RangeTransactions(txmgrNs, begin, end, f)
	})
}
