// Copyright (c) 2015 The btcsuite developers
package wtxmgr_test

import (
	"fmt"
	"github.com/p9c/pkg/app/slog"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/db/walletdb"
)

var (
	// Spends: bogus
	// Outputs: 10 DUO
	exampleTxRecordA *wtxmgr.TxRecord
	// Spends: A:0
	// Outputs: 5 DUO, 5 DUO
	exampleTxRecordB *wtxmgr.TxRecord
)

func init() {
	tx := spendOutput(&chainhash.Hash{}, 0, 10e8)
	rec, err := wtxmgr.NewTxRecordFromMsgTx(tx, timeNow())
	if err != nil {
		panic(err)
	}
	exampleTxRecordA = rec
	tx = spendOutput(&exampleTxRecordA.Hash, 0, 5e8, 5e8)
	rec, err = wtxmgr.NewTxRecordFromMsgTx(tx, timeNow())
	if err != nil {
		panic(err)
	}
	exampleTxRecordB = rec
}

var exampleBlock100 = makeBlockMeta(100)

// This example demonstrates reporting the Store balance given an unmined and
// mined transaction given 0, 1, and 6 block confirmations.
func ExampleStore_Balance() {
	s, db, teardown, err := testStore()
	defer teardown()
	if err != nil {
		slog.Error(err)
		return
	}
	// Prints balances for 0 block confirmations, 1 confirmation, and 6
	// confirmations.
	printBalances := func(syncHeight int32) {
		dbtx, err := db.BeginReadTx()
		if err != nil {
			slog.Error(err)
			return
		}
		defer func() {
			err := dbtx.Rollback()
			if err != nil {
				slog.Error(err)
			}
		}()
		ns := dbtx.ReadBucket(namespaceKey)
		zeroConfBal, err := s.Balance(ns, 0, syncHeight)
		if err != nil {
			slog.Error(err)
			return
		}
		oneConfBal, err := s.Balance(ns, 1, syncHeight)
		if err != nil {
			slog.Error(err)
			return
		}
		sixConfBal, err := s.Balance(ns, 6, syncHeight)
		if err != nil {
			slog.Error(err)
			return
		}
		fmt.Printf("%v, %v, %v\n", zeroConfBal, oneConfBal, sixConfBal)
	}
	// Insert a transaction which outputs 10 DUO unmined and mark the output
	// as a credit.
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (err error) {
		ns := tx.ReadWriteBucket(namespaceKey)
		err := s.InsertTx(ns, exampleTxRecordA, nil)
		if err != nil {
			return err
		}
		return s.AddCredit(ns, exampleTxRecordA, nil, 0, false)
	})
	if err != nil {
		slog.Error(err)
		return
	}
	printBalances(100)
	// Mine the transaction in block 100 and print balances again with a
	// sync height of 100 and 105 blocks.
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (err error) {
		ns := tx.ReadWriteBucket(namespaceKey)
		return s.InsertTx(ns, exampleTxRecordA, &exampleBlock100)
	})
	if err != nil {
		slog.Error(err)
		return
	}
	printBalances(100)
	printBalances(105)
	// Output:
	// 10 DUO, 0 DUO, 0 DUO
	// 10 DUO, 10 DUO, 0 DUO
	// 10 DUO, 10 DUO, 10 DUO
}
func ExampleStore_Rollback() {
	s, db, teardown, err := testStore()
	defer teardown()
	if err != nil {
		slog.Error(err)
		return
	}
	err = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (err error) {
		ns := tx.ReadWriteBucket(namespaceKey)
		// Insert a transaction which outputs 10 DUO in a block at height 100.
		err := s.InsertTx(ns, exampleTxRecordA, &exampleBlock100)
		if err != nil {
			return err
		}
		// Rollback everything from block 100 onwards.
		err = s.Rollback(ns, 100)
		if err != nil {
			return err
		}
		// Assert that the transaction is now unmined.
		details, err := s.TxDetails(ns, &exampleTxRecordA.Hash)
		if err != nil {
			return err
		}
		if details == nil {
			return fmt.Errorf("no details found")
		}
		fmt.Println(details.Block.Height)
		return nil
	})
	if err != nil {
		slog.Error(err)
		return
	}
	// Output:
	// -1
}
func Example_basicUsage() {
	// Open the database.
	db, dbTeardown, err := testDB()
	defer dbTeardown()
	if err != nil {
		slog.Error(err)
		return
	}
	// Open a read-write transaction to operate on the database.
	dbtx, err := db.BeginReadWriteTx()
	if err != nil {
		slog.Error(err)
		return
	}
	defer func() {
		err := dbtx.Commit()
		if err != nil {
			slog.Error(err)
		}
	}()
	// Create a bucket for the transaction store.
	b, err := dbtx.CreateTopLevelBucket([]byte("txstore"))
	if err != nil {
		slog.Error(err)
		return
	}
	// Create and open the transaction store in the provided namespace.
	err = wtxmgr.Create(b)
	if err != nil {
		slog.Error(err)
		return
	}
	s, err := wtxmgr.Open(b, &netparams.TestNet3Params)
	if err != nil {
		slog.Error(err)
		return
	}
	// Insert an unmined transaction that outputs 10 DUO to a wallet address
	// at output 0.
	err = s.InsertTx(b, exampleTxRecordA, nil)
	if err != nil {
		slog.Error(err)
		return
	}
	err = s.AddCredit(b, exampleTxRecordA, nil, 0, false)
	if err != nil {
		slog.Error(err)
		return
	}
	// Insert a second transaction which spends the output, and creates two
	// outputs.  Mark the second one (5 DUO) as wallet change.
	err = s.InsertTx(b, exampleTxRecordB, nil)
	if err != nil {
		slog.Error(err)
		return
	}
	err = s.AddCredit(b, exampleTxRecordB, nil, 1, true)
	if err != nil {
		slog.Error(err)
		return
	}
	// Mine each transaction in a block at height 100.
	err = s.InsertTx(b, exampleTxRecordA, &exampleBlock100)
	if err != nil {
		slog.Error(err)
		return
	}
	err = s.InsertTx(b, exampleTxRecordB, &exampleBlock100)
	if err != nil {
		slog.Error(err)
		return
	}
	// Print the one confirmation balance.
	bal, err := s.Balance(b, 1, 100)
	if err != nil {
		slog.Error(err)
		return
	}
	fmt.Println(bal)
	// Fetch unspent outputs.
	utxos, err := s.UnspentOutputs(b)
	if err != nil {
		slog.Error(err)
	}
	expectedOutPoint := wire.OutPoint{Hash: exampleTxRecordB.Hash, Index: 1}
	for _, utxo := range utxos {
		fmt.Println(utxo.OutPoint == expectedOutPoint)
	}
	// Output:
	// 5 DUO
	// true
}
