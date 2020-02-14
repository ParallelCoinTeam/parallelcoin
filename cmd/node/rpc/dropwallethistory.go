package rpc

import (
	"encoding/binary"
	"fmt"
	"path/filepath"
	
	"github.com/urfave/cli"
	
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

func DropWalletHistory(cfg *pod.Config) func(c *cli.Context) error {
	return func(c *cli.Context) error {
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
		return err
	}
}
