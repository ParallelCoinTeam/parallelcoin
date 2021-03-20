package votingpool_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	wtxmgr "github.com/p9c/pod/pkg/blockchain/tx/wtxmgr"
	txscript "github.com/p9c/pod/pkg/blockchain/tx/txscript"
	"github.com/p9c/pod/pkg/database/walletdb"
	_ "github.com/p9c/pod/pkg/database/walletdb/bdb"
	"github.com/p9c/pod/pkg/util"
	waddrmgr "github.com/p9c/pod/pkg/wallet/waddrmgr"
	"github.com/p9c/pod/pkg/wallet/votingpool"
)

var (
	pubPassphrase  = []byte("pubPassphrase")
	privPassphrase = []byte("privPassphrase")
	seed           = bytes.Repeat([]byte{0x2a, 0x64, 0xdf, 0x08}, 8)
	fastScrypt     = &waddrmgr.ScryptOptions{N: 16, R: 8, P: 1}
)

func createWaddrmgr(ns walletdb.ReadWriteBucket, params *netparams.Params) (*waddrmgr.Manager, error) {
	e := waddrmgr.Create(ns, seed, pubPassphrase, privPassphrase, params,
		fastScrypt, time.Now())
	if e != nil  {
		return nil, e
	}
	return waddrmgr.Open(ns, pubPassphrase, params)
}
func ExampleCreate() {
	// Create a new walletdb.DB. See the walletdb docs for instructions on how
	// to do that.
	db, dbTearDown, e := createWalletDB()
	if e != nil  {
		votingpool.		return
	}
	defer dbTearDown()
	dbtx, e := db.BeginReadWriteTx()
	if e != nil  {
		votingpool.		return
	}
	defer func() {
		e := dbtx.Commit()
		if e != nil  {
			votingpool.		}
	}()
	// Create a new walletdb namespace for the address manager.
	mgrNamespace, e := dbtx.CreateTopLevelBucket([]byte("waddrmgr"))
	if e != nil  {
		votingpool.		return
	}
	// Create the address manager.
	mgr, e := createWaddrmgr(mgrNamespace, &netparams.MainNetParams)
	if e != nil  {
		votingpool.		return
	}
	// Create a walletdb namespace for votingpools.
	vpNamespace, e := dbtx.CreateTopLevelBucket([]byte("votingpool"))
	if e != nil  {
		votingpool.		return
	}
	// Create a voting pool.
	_, e = votingpool.Create(vpNamespace, mgr, []byte{0x00})
	if e != nil  {
		votingpool.		return
	}
	// Output:
	//
}

// This example demonstrates how to create a voting pool with one
// series and get a deposit address for that series.
func Example_depositAddress() {
	// Create the address manager and votingpool DB namespace. See the example
	// for the Create() function for more info on how this is done.
	teardown, db, mgr := exampleCreateDBAndMgr()
	defer teardown()
	e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		ns := votingpoolNamespace(tx)
		// Create the voting pool.
		pool, e := votingpool.Create(ns, mgr, []byte{0x00})
		if e != nil  {
			return e
		}
		// Create a 2-of-3 series.
		seriesID := uint32(1)
		requiredSignatures := uint32(2)
		pubKeys := []string{
			"xpub661MyMwAqRbcFDDrR5jY7LqsRioFDwg3cLjc7tML3RRcfYyhXqqgCH5SqMSQdpQ1Xh8EtVwcfm8psD8zXKPcRaCVSY4GCqbb3aMEs27GitE",
			"xpub661MyMwAqRbcGsxyD8hTmJFtpmwoZhy4NBBVxzvFU8tDXD2ME49A6JjQCYgbpSUpHGP1q4S2S1Pxv2EqTjwfERS5pc9Q2yeLkPFzSgRpjs9",
			"xpub661MyMwAqRbcEbc4uYVXvQQpH9L3YuZLZ1gxCmj59yAhNy33vXxbXadmRpx5YZEupNSqWRrR7PqU6duS2FiVCGEiugBEa5zuEAjsyLJjKCh",
		}
		e = pool.CreateSeries(ns, votingpool.CurrentVersion, seriesID, requiredSignatures, pubKeys)
		if e != nil  {
			return e
		}
		// Create a deposit address.
		addr, e := pool.DepositScriptAddress(seriesID, votingpool.Branch(0), votingpool.Index(1))
		if e != nil  {
			return e
		}
		fmt.Println("Generated deposit address:", addr.EncodeAddress())
		return nil
	})
	if e != nil  {
		votingpool.		return
	}
	// Output:
	// Generated deposit address: 51pQm3LmtcK6e4rgGoJDpdCw2N4uWZB9wr
}

// This example demonstrates how to empower a series by loading the private
// key for one of the series' public keys.
func Example_empowerSeries() {
	// Create the address manager and votingpool DB namespace. See the example
	// for the Create() function for more info on how this is done.
	teardown, db, mgr := exampleCreateDBAndMgr()
	defer teardown()
	// Create a pool and a series. See the DepositAddress example for more info
	// on how this is done.
	pool, seriesID := exampleCreatePoolAndSeries(db, mgr)
	e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		ns := votingpoolNamespace(tx)
		addrmgrNs := addrmgrNamespace(tx)
		// Now empower the series with one of its private keys. Notice that in order
		// to do that we need to unlock the address manager.
		e := mgr.Unlock(addrmgrNs, privPassphrase)
		if e != nil  {
			return e
		}
		defer func() {
			e := mgr.Lock()
			if e != nil  {
				votingpool.			}
		}()
		privKey := "xprv9s21ZrQH143K2j9PK4CXkCu8sgxkpUxCF7p1KVwiV5tdnkeYzJXReUkxz5iB2FUzTXC1L15abCDG4RMxSYT5zhm67uvsnLYxuDhZfoFcB6a"
		return pool.EmpowerSeries(ns, seriesID, privKey)
	})
	if e != nil  {
		votingpool.		return
	}
	// Output:
	//
}

// This example demonstrates how to use the Pool.StartWithdrawal method.
func Example_startWithdrawal() {
	// Create the address manager and votingpool DB namespace. See the example
	// for the Create() function for more info on how this is done.
	teardown, db, mgr := exampleCreateDBAndMgr()
	defer teardown()
	// Create a pool and a series. See the DepositAddress example for more info
	// on how this is done.
	pool, seriesID := exampleCreatePoolAndSeries(db, mgr)
	e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		ns := votingpoolNamespace(tx)
		addrmgrNs := addrmgrNamespace(tx)
		txmgrNs := txmgrNamespace(tx)
		// Create the transaction store for later use.
		txstore := exampleCreateTxStore(txmgrNs)
		// Unlock the manager
		e := mgr.Unlock(addrmgrNs, privPassphrase)
		if e != nil  {
			return e
		}
		defer func() {
			e := mgr.Lock()
			if e != nil  {
				votingpool.			}
		}()
		addr, _ := util.DecodeAddress("1MirQ9bwyQcGVJPwKUgapu5ouK2E2Ey4gX", mgr.ChainParams())
		pkScript, _ := txscript.PayToAddrScript(addr)
		requests := []votingpool.OutputRequest{
			{
				PkScript:    pkScript,
				Address:     addr,
				Amount:      1e6,
				Server:      "server-id",
				Transaction: 123,
			},
		}
		changeStart, e := pool.ChangeAddress(seriesID, votingpool.Index(0))
		if e != nil  {
			return e
		}
		// This is only needed because we have not used any deposit addresses from
		// the series, and we cannot create a WithdrawalAddress for an unused
		// branch/idx pair.
		e = pool.EnsureUsedAddr(ns, addrmgrNs, seriesID, votingpool.Branch(1), votingpool.Index(0))
		if e != nil  {
			return e
		}
		startAddr, e := pool.WithdrawalAddress(ns, addrmgrNs, seriesID, votingpool.Branch(1), votingpool.Index(0))
		if e != nil  {
			return e
		}
		lastSeriesID := seriesID
		dustThreshold := util.Amount(1e4)
		currentBlock := int32(19432)
		roundID := uint32(0)
		_, e = pool.StartWithdrawal(ns, addrmgrNs,
			roundID, requests, *startAddr, lastSeriesID, *changeStart, txstore, txmgrNs, currentBlock,
			dustThreshold)
		return e
	})
	if e != nil  {
		votingpool.		return
	}
	// Output:
	//
}
func createWalletDB() (walletdb.DB, func(), error) {
	dir, e := ioutil.TempDir("", "votingpool_example")
	if e != nil  {
		return nil, nil, e
	}
	db, e := walletdb.Create("bdb", filepath.Join(dir, "wallet.db"))
	if e != nil  {
		return nil, nil, e
	}
	dbTearDown := func() {
		if e := db.Close(); votingpool.E.Chk(e) {
		}
		if e := os.RemoveAll(dir); votingpool.E.Chk(e) {
		}
	}
	return db, dbTearDown, nil
}

var (
	addrmgrNamespaceKey    = []byte("addrmgr")
	txmgrNamespaceKey      = []byte("txmgr")
	votingpoolNamespaceKey = []byte("votingpool")
)

func addrmgrNamespace(dbtx walletdb.ReadWriteTx) walletdb.ReadWriteBucket {
	return dbtx.ReadWriteBucket(addrmgrNamespaceKey)
}
func txmgrNamespace(dbtx walletdb.ReadWriteTx) walletdb.ReadWriteBucket {
	return dbtx.ReadWriteBucket(txmgrNamespaceKey)
}
func votingpoolNamespace(dbtx walletdb.ReadWriteTx) walletdb.ReadWriteBucket {
	return dbtx.ReadWriteBucket(votingpoolNamespaceKey)
}
func exampleCreateDBAndMgr() (teardown func(), db walletdb.DB, mgr *waddrmgr.Manager) {
	var dbTearDown func()
	var e error
	if db, dbTearDown, e = createWalletDB(); votingpool.E.Chk(e){
		dbTearDown()
		panic(e)
	}
	// Create a new walletdb namespace for the address manager.
	e = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		addrmgrNs, e := tx.CreateTopLevelBucket(addrmgrNamespaceKey)
		if e != nil  {
			return e
		}
		_, e = tx.CreateTopLevelBucket(votingpoolNamespaceKey)
		if e != nil  {
			return e
		}
		_, e = tx.CreateTopLevelBucket(txmgrNamespaceKey)
		if e != nil  {
			return e
		}
		// Create the address manager
		mgr, e = createWaddrmgr(addrmgrNs, &netparams.MainNetParams)
		return e
	})
	if e != nil  {
		dbTearDown()
		panic(e)
	}
	teardown = func() {
		mgr.Close()
		dbTearDown()
	}
	return teardown, db, mgr
}
func exampleCreatePoolAndSeries(db walletdb.DB, mgr *waddrmgr.Manager) (pool *votingpool.Pool, seriesID uint32) {
	e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		ns := votingpoolNamespace(tx)
		var e error
		pool, e = votingpool.Create(ns, mgr, []byte{0x00})
		if e != nil  {
			return e
		}
		// Create a 2-of-3 series.
		seriesID = uint32(1)
		requiredSignatures := uint32(2)
		pubKeys := []string{
			"xpub661MyMwAqRbcFDDrR5jY7LqsRioFDwg3cLjc7tML3RRcfYyhXqqgCH5SqMSQdpQ1Xh8EtVwcfm8psD8zXKPcRaCVSY4GCqbb3aMEs27GitE",
			"xpub661MyMwAqRbcGsxyD8hTmJFtpmwoZhy4NBBVxzvFU8tDXD2ME49A6JjQCYgbpSUpHGP1q4S2S1Pxv2EqTjwfERS5pc9Q2yeLkPFzSgRpjs9",
			"xpub661MyMwAqRbcEbc4uYVXvQQpH9L3YuZLZ1gxCmj59yAhNy33vXxbXadmRpx5YZEupNSqWRrR7PqU6duS2FiVCGEiugBEa5zuEAjsyLJjKCh",
		}
		e = pool.CreateSeries(ns, votingpool.CurrentVersion, seriesID, requiredSignatures, pubKeys)
		if e != nil  {
			return e
		}
		return pool.ActivateSeries(ns, seriesID)
	})
	if e != nil  {
		panic(e)
	}
	return pool, seriesID
}
func exampleCreateTxStore(ns walletdb.ReadWriteBucket) *wtxmgr.Store {
	e := wtxmgr.Create(ns)
	if e != nil  {
		panic(e)
	}
	s, e := wtxmgr.Open(ns, &netparams.MainNetParams)
	if e != nil  {
		panic(e)
	}
	return s
}
