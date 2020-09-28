// Copyright (c) 2015-2016 The btcsuite developers
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/p9c/pkg/app/slog"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"github.com/p9c/pod/app/appdata"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	"github.com/p9c/pod/pkg/db/walletdb"
	_ "github.com/p9c/pod/pkg/db/walletdb/bdb"
)

const defaultNet = "mainnet"

var datadir = appdata.Dir("mod", false)

// Flags.
var opts = struct {
	Force  bool   `short:"f" description:"Force removal without prompt"`
	DbPath string `long:"db" description:"Path to wallet database"`
}{
	Force:  false,
	DbPath: filepath.Join(datadir, defaultNet, "wallet.db"),
}

func init() {
	_, err := flags.Parse(&opts)
	if err != nil {
		slog.Error(err)
		os.Exit(1)
	}
}

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

func yes(
	s string) bool {
	switch s {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}
func no(
	s string) bool {
	switch s {
	case "n", "N", "no", "No":
		return true
	default:
		return false
	}
}
func main() {
	os.Exit(mainInt())
}
func mainInt() int {
	fmt.Println("Database path:", opts.DbPath)
	_, err := os.Stat(opts.DbPath)
	if os.IsNotExist(err) {
		fmt.Println("Database file does not exist")
		return 1
	}
	for !opts.Force {
		fmt.Print("Drop all mod transaction history? [y/N] ")
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		if !scanner.Scan() {
			// Exit on EOF.
			return 0
		}
		err := scanner.Err()
		if err != nil {
			slog.Error(err)
			return 1
		}
		resp := scanner.Text()
		if yes(resp) {
			break
		}
		if no(resp) || resp == "" {
			return 0
		}
		fmt.Println("Enter yes or no.")
	}
	db, err := walletdb.Open("bdb", opts.DbPath)
	if err != nil {
		fmt.Println("failed to open database:", err)
		return 1
	}
	defer db.Close()
	fmt.Println("dropping wtxmgr namespace")
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (err error) {
		err := tx.DeleteTopLevelBucket(wtxmgrNamespace)
		if err != nil && err != walletdb.ErrBucketNotFound {
			return err
		}
		ns, err := tx.CreateTopLevelBucket(wtxmgrNamespace)
		if err != nil {
			slog.Error(err)
			return err
		}
		err = wtxmgr.Create(ns)
		if err != nil {
			slog.Error(err)
			return err
		}
		ns = tx.ReadWriteBucket(waddrmgrNamespace).NestedReadWriteBucket(syncBucketName)
		startBlock := ns.Get(startBlockName)
		err = ns.Put(syncedToName, startBlock)
		if err != nil {
			slog.Error(err)
			return err
		}
		recentBlocks := make([]byte, 40)
		copy(recentBlocks[0:4], startBlock[0:4])
		copy(recentBlocks[8:], startBlock[4:])
		binary.LittleEndian.PutUint32(recentBlocks[4:8], uint32(1))
		return ns.Put(recentBlocksName, recentBlocks)
	})
	if err != nil {
		slog.Error(err)
		fmt.Println("Failed to drop and re-create namespace:", err)
		return 1
	}
	return 0
}
