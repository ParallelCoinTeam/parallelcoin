package wallet

import (
	"bytes"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"strings"
	
	"github.com/p9c/matrjoska/pkg/chainclient"
	"github.com/p9c/matrjoska/pkg/txscript"
	wm "github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/walletdb"
	tm "github.com/p9c/matrjoska/pkg/wtxmgr"
)

func (w *Wallet) handleChainNotifications() {
	defer w.wg.Done()
	if w == nil {
		panic("w should not be nil")
	}
	chainClient, e := w.requireChainClient()
	if e != nil {
		E.Ln("handleChainNotifications called without RPC client", e)
		return
	}
	sync := func(w *Wallet) {
		if w.db != nil {
			// At the moment there is no recourse if the rescan fails for some reason, however, the wallet will not be
			// marked synced and many methods will error early since the wallet is known to be out of date.
			e := w.syncWithChain()
			if e != nil && !w.ShuttingDown() {
				W.Ln("unable to synchronize wallet to chain:", e)
			}
		}
	}
	catchUpHashes := func(
		w *Wallet, client chainclient.Interface,
		height int32,
	) (e error) {
		// TODO(aakselrod): There's a race condition here, which happens when a reorg occurs between the rescanProgress
		//  notification and the last GetBlockHash call. The solution when using pod is to make pod send blockconnected
		//  notifications with each block the way Neutrino does, and get rid of the loop. The other alternative is to
		//  check the final hash and, if it doesn't match the original hash returned by the notification, to roll back
		//  and restart the rescan.
		I.F(
			"handleChainNotifications: catching up block hashes to height %d, this might take a while", height,
		)
		e = walletdb.Update(
			w.db, func(tx walletdb.ReadWriteTx) (e error) {
				ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				startBlock := w.Manager.SyncedTo()
				for i := startBlock.Height + 1; i <= height; i++ {
					hash, e := client.GetBlockHash(int64(i))
					if e != nil {
						return e
					}
					header, e := chainClient.GetBlockHeader(hash)
					if e != nil {
						return e
					}
					bs := wm.BlockStamp{
						Height:    i,
						Hash:      *hash,
						Timestamp: header.Timestamp,
					}
					e = w.Manager.SetSyncedTo(ns, &bs)
					if e != nil {
						return e
					}
				}
				return nil
			},
		)
		if e != nil {
			E.F(
				"failed to update address manager sync state for height %d: %v",
				height, e,
			)
		}
		I.Ln("done catching up block hashes")
		return e
	}
	for {
		select {
		case n, ok := <-chainClient.Notifications():
			if !ok {
				return
			}
			var notificationName string
			var e error
			switch n := n.(type) {
			case chainclient.ClientConnected:
				if w != nil {
					go sync(w)
				}
			case chainclient.BlockConnected:
				e = walletdb.Update(
					w.db, func(tx walletdb.ReadWriteTx) (e error) {
						return w.connectBlock(tx, tm.BlockMeta(n))
					},
				)
				notificationName = "blockconnected"
			case chainclient.BlockDisconnected:
				e = walletdb.Update(
					w.db, func(tx walletdb.ReadWriteTx) (e error) {
						return w.disconnectBlock(tx, tm.BlockMeta(n))
					},
				)
				notificationName = "blockdisconnected"
			case chainclient.RelevantTx:
				e = walletdb.Update(
					w.db, func(tx walletdb.ReadWriteTx) (e error) {
						return w.addRelevantTx(tx, n.TxRecord, n.Block)
					},
				)
				notificationName = "recvtx/redeemingtx"
			case chainclient.FilteredBlockConnected:
				// Atomically update for the whole block.
				if len(n.RelevantTxs) > 0 {
					e = walletdb.Update(
						w.db, func(
							tx walletdb.ReadWriteTx,
						) (e error) {
							for _, rec := range n.RelevantTxs {
								e = w.addRelevantTx(
									tx, rec,
									n.Block,
								)
								if e != nil {
									return e
								}
							}
							return nil
						},
					)
				}
				notificationName = "filteredblockconnected"
			// The following require some database maintenance, but also need to be reported to the wallet's rescan
			// goroutine.
			case *chainclient.RescanProgress:
				e = catchUpHashes(w, chainClient, n.Height)
				notificationName = "rescanprogress"
				select {
				case w.rescanNotifications <- n:
				case <-w.quitChan().Wait():
					return
				}
			case *chainclient.RescanFinished:
				e = catchUpHashes(w, chainClient, n.Height)
				notificationName = "rescanprogress"
				w.SetChainSynced(true)
				select {
				case w.rescanNotifications <- n:
				case <-w.quitChan().Wait():
					return
				}
			}
			if e != nil {
				// On out-of-sync blockconnected notifications, only send a debug message.
				errStr := "failed to process consensus server " +
					"notification (name: `%s`, detail: `%v`)"
				if notificationName == "blockconnected" &&
					strings.Contains(
						e.Error(),
						"couldn't get hash from database",
					) {
					D.F(errStr, notificationName, e)
				} else {
					E.F(errStr, notificationName, e)
				}
			}
		case <-w.quit.Wait():
			return
		}
	}
}

// connectBlock handles a chain server notification by marking a wallet that's currently in-sync with the chain server
// as being synced up to the passed block.
func (w *Wallet) connectBlock(dbtx walletdb.ReadWriteTx, b tm.BlockMeta) (e error) {
	addrmgrNs := dbtx.ReadWriteBucket(waddrmgrNamespaceKey)
	bs := wm.BlockStamp{
		Height:    b.Height,
		Hash:      b.Hash,
		Timestamp: b.Time,
	}
	e = w.Manager.SetSyncedTo(addrmgrNs, &bs)
	if e != nil {
		return e
	}
	// Notify interested clients of the connected block.
	//
	// TODO: move all notifications outside of the database transaction.
	w.NtfnServer.notifyAttachedBlock(dbtx, &b)
	return nil
}

// disconnectBlock handles a chain server reorganize by rolling back all block history from the reorged block for a
// wallet in-sync with the chain server.
func (w *Wallet) disconnectBlock(dbtx walletdb.ReadWriteTx, b tm.BlockMeta) (e error) {
	addrmgrNs := dbtx.ReadWriteBucket(waddrmgrNamespaceKey)
	txmgrNs := dbtx.ReadWriteBucket(wtxmgrNamespaceKey)
	if !w.ChainSynced() {
		return nil
	}
	// Disconnect the removed block and all blocks after it if we know about the disconnected block. Otherwise, the
	// block is in the future.
	if b.Height <= w.Manager.SyncedTo().Height {
		hash, e := w.Manager.BlockHash(addrmgrNs, b.Height)
		if e != nil {
			return e
		}
		if bytes.Equal(hash[:], b.Hash[:]) {
			bs := wm.BlockStamp{
				Height: b.Height - 1,
			}
			hash, e = w.Manager.BlockHash(addrmgrNs, bs.Height)
			if e != nil {
				return e
			}
			b.Hash = *hash
			client := w.ChainClient()
			header, e := client.GetBlockHeader(hash)
			if e != nil {
				return e
			}
			bs.Timestamp = header.Timestamp
			e = w.Manager.SetSyncedTo(addrmgrNs, &bs)
			if e != nil {
				return e
			}
			e = w.TxStore.Rollback(txmgrNs, b.Height)
			if e != nil {
				return e
			}
		}
	}
	// Notify interested clients of the disconnected block.
	w.NtfnServer.notifyDetachedBlock(&b.Hash)
	return nil
}
func (w *Wallet) addRelevantTx(dbtx walletdb.ReadWriteTx, rec *tm.TxRecord, block *tm.BlockMeta) (e error) {
	addrmgrNs := dbtx.ReadWriteBucket(waddrmgrNamespaceKey)
	txmgrNs := dbtx.ReadWriteBucket(wtxmgrNamespaceKey)
	// At the moment all notified transactions are assumed to actually be relevant. This assumption will not hold true
	// when SPV support is added, but until then, simply insert the transaction because there should either be one or
	// more relevant inputs or outputs.
	e = w.TxStore.InsertTx(txmgrNs, rec, block)
	if e != nil {
		return e
	}
	// Chk every output to determine whether it is controlled by a wallet key. If so, mark the output as a credit.
	for i, output := range rec.MsgTx.TxOut {
		var addrs []btcaddr.Address
		_, addrs, _, e = txscript.ExtractPkScriptAddrs(
			output.PkScript,
			w.chainParams,
		)
		if e != nil {
			// Non-standard outputs are skipped.
			continue
		}
		for _, addr := range addrs {
			ma, e := w.Manager.Address(addrmgrNs, addr)
			if e == nil {
				// TODO: Credits should be added with the account they belong to, so tm is able to track per-account
				//  balances.
				e = w.TxStore.AddCredit(
					txmgrNs, rec, block, uint32(i),
					ma.Internal(),
				)
				if e != nil {
					return e
				}
				e = w.Manager.MarkUsed(addrmgrNs, addr)
				if e != nil {
					return e
				}
				T.Ln("marked address used:", addr)
				continue
			}
			// Missing addresses are skipped. Other errors should be propagated.
			if !wm.IsError(e, wm.ErrAddressNotFound) {
				return e
			}
		}
	}
	// Send notification of mined or unmined transaction to any interested clients.
	//
	// TODO: Avoid the extra db hits.
	if block == nil {
		details, e := w.TxStore.UniqueTxDetails(txmgrNs, &rec.Hash, nil)
		if e != nil {
			E.Ln("cannot query transaction details for notification:", e)
		}
		// It's possible that the transaction was not found within the wallet's set of unconfirmed transactions due to
		// it already being confirmed, so we'll avoid notifying it.
		//
		// TODO(wilmer): ideally we should find the culprit to why we're receiving an additional unconfirmed
		//  chain.RelevantTx notification from the chain backend.
		if details != nil {
			w.NtfnServer.notifyUnminedTransaction(dbtx, details)
		}
	} else {
		details, e := w.TxStore.UniqueTxDetails(txmgrNs, &rec.Hash, &block.Block)
		if e != nil {
			E.Ln("cannot query transaction details for notification:", e)
		}
		// We'll only notify the transaction if it was found within the wallet's set of confirmed transactions.
		if details != nil {
			w.NtfnServer.notifyMinedTransaction(dbtx, details, block)
		}
	}
	return nil
}
