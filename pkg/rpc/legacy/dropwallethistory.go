package legacy

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	"time"
	
	"github.com/urfave/cli"
	
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

func DropWalletHistory(w *wallet.Wallet) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg := w.PodConfig
		fmt.Println("\n", cfg)
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
		db, err := walletdb.Open("bdb",
			filepath.Join(*cfg.DataDir,
				*cfg.Network, "wallet.db"))
		if err != nil {
			log.Println("failed to open database:", err)
			return err
		}
		defer db.Close()
		log.Println("dropping wtxmgr namespace")
		err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) error {
			err := tx.DeleteTopLevelBucket(wtxmgrNamespace)
			if err != nil && err != walletdb.ErrBucketNotFound {
				return err
			}
			ns, err := tx.CreateTopLevelBucket(wtxmgrNamespace)
			if err != nil {
				log.ERROR(err)
				return err
			}
			err = wtxmgr.Create(ns)
			if err != nil {
				log.ERROR(err)
				return err
			}
			ns = tx.ReadWriteBucket(waddrmgrNamespace).NestedReadWriteBucket(syncBucketName)
			startBlock := ns.Get(startBlockName)
			err = ns.Put(syncedToName, startBlock)
			if err != nil {
				log.ERROR(err)
				return err
			}
			recentBlocks := make([]byte, 40)
			copy(recentBlocks[0:4], startBlock[0:4])
			copy(recentBlocks[8:], startBlock[4:])
			binary.LittleEndian.PutUint32(recentBlocks[4:8], uint32(1))
			return ns.Put(recentBlocksName, recentBlocks)
		})
		if err != nil {
			log.ERROR(err)
			return err
		}
		if w != nil {
			// Rescan chain to ensure balance is correctly regenerated
			job := &wallet.RescanJob{
				InitialSync: true,
			}
			// Submit rescan job and log when the import has completed.
			// Do not block on finishing the rescan.  The rescan success
			// or failure is logged elsewhere, and the channel is not
			// required to be read, so discard the return value.
			errC := w.SubmitRescan(job)
			select {
			case err := <-errC:
				log.ERROR(err)
			case <-time.After(time.Second * 5):
				break
			}
		}
		return err
	}
}
