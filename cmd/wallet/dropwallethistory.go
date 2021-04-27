package wallet

import (
	"encoding/binary"
	"path/filepath"

	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pkg/wtxmgr"
	"github.com/p9c/matrjoska/pod/config"
)

func DropWalletHistory(w *Wallet, cfg *config.Config) (e error) {
	var (
		// Namespace keys.
		syncBucketName    = []byte("sync")
		waddrmgrNamespace = []byte("waddrmgr")
		wtxmgrNamespace   = []byte("wtxmgr")
		// Sync related key names (sync bucket).
		syncedToName     = []byte("syncedto")
		startBlockName   = []byte("startblock")
		recentBlocksName = []byte("recentblocks")
	)
	dbPath := filepath.Join(
		cfg.DataDir.V(),
		cfg.Network.V(), "wallet.db",
	)
	// I.Ln("dbPath", dbPath)
	var db walletdb.DB
	db, e = walletdb.Open("bdb", dbPath)
	if E.Chk(e) {
		// DBError("failed to open database:", err)
		return e
	}
	defer db.Close()
	D.Ln("dropping wtxmgr namespace")
	e = walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			D.Ln("deleting top level bucket")
			if e = tx.DeleteTopLevelBucket(wtxmgrNamespace); E.Chk(e) {
			}
			if e != nil && e != walletdb.ErrBucketNotFound {
				return e
			}
			var ns walletdb.ReadWriteBucket
			D.Ln("creating new top level bucket")
			if ns, e = tx.CreateTopLevelBucket(wtxmgrNamespace); E.Chk(e) {
				return e
			}
			if e = wtxmgr.Create(ns); E.Chk(e) {
				return e
			}
			ns = tx.ReadWriteBucket(waddrmgrNamespace).NestedReadWriteBucket(syncBucketName)
			startBlock := ns.Get(startBlockName)
			D.Ln("putting start block", startBlock)
			if e = ns.Put(syncedToName, startBlock); E.Chk(e) {
				return e
			}
			recentBlocks := make([]byte, 40)
			copy(recentBlocks[0:4], startBlock[0:4])
			copy(recentBlocks[8:], startBlock[4:])
			binary.LittleEndian.PutUint32(recentBlocks[4:8], uint32(1))
			defer D.Ln("put recent blocks")
			return ns.Put(recentBlocksName, recentBlocks)
		},
	)
	if E.Chk(e) {
		return e
	}
	D.Ln("updated wallet")
	// if w != nil {
	// 	// Rescan chain to ensure balance is correctly regenerated
	// 	job := &wallet.RescanJob{
	// 		InitialSync: true,
	// 	}
	// 	// Submit rescan job and log when the import has completed.
	// 	// Do not block on finishing the rescan.  The rescan success
	// 	// or failure is logged elsewhere, and the channel is not
	// 	// required to be read, so discard the return value.
	// 	errC := w.SubmitRescan(job)
	// 	select {
	// 	case e := <-errC:
	// 		DB		// 		// case <-time.After(time.Second * 5):
	// 		// 	break
	// 	}
	// }
	return e
}
