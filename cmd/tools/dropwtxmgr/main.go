// Copyright (c) 2015-2016 The btcsuite developers
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/jessevdk/go-flags"
	
	"github.com/p9c/pod/pkg/appdata"
	"github.com/p9c/pod/pkg/walletdb"
	_ "github.com/p9c/pod/pkg/walletdb/bdb"
	"github.com/p9c/pod/pkg/wtxmgr"
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

	_, e := flags.Parse(&opts)
	if e != nil {
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
	s string,
) bool {
	switch s {
	case "y", "Y", "yes", "Yes":
		return true
	default:
		return false
	}
}
func no(
	s string,
) bool {
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
	_, e := os.Stat(opts.DbPath)
	if os.IsNotExist(e) {
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
		e := scanner.Err()
		if e != nil {
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
	db, e := walletdb.Open("bdb", opts.DbPath)
	if e != nil {
		fmt.Println("failed to open database:", e)
		return 1
	}
	defer func() {
		if e := db.Close(); E.Chk(e) {
			fmt.Println(e)
		}
	}()
	fmt.Println("dropping wtxmgr namespace")
	e = walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			e = tx.DeleteTopLevelBucket(wtxmgrNamespace)
			if e != nil && e != walletdb.ErrBucketNotFound {
				return e
			}
			ns, e := tx.CreateTopLevelBucket(wtxmgrNamespace)
			if e != nil {
				return e
			}
			e = wtxmgr.Create(ns)
			if e != nil {
				return e
			}
			ns = tx.ReadWriteBucket(waddrmgrNamespace).NestedReadWriteBucket(syncBucketName)
			startBlock := ns.Get(startBlockName)
			e = ns.Put(syncedToName, startBlock)
			if e != nil {
				return e
			}
			recentBlocks := make([]byte, 40)
			copy(recentBlocks[0:4], startBlock[0:4])
			copy(recentBlocks[8:], startBlock[4:])
			binary.LittleEndian.PutUint32(recentBlocks[4:8], uint32(1))
			return ns.Put(recentBlocksName, recentBlocks)
		},
	)
	if e != nil {
		fmt.Println("Failed to drop and re-create namespace:", e)
		return 1
	}
	return 0
}
