package votingpool

// Helpers to create parameterized objects to use in tests.
import (
	"bytes"
	amount2 "github.com/p9c/pod/pkg/amt"
	"github.com/p9c/pod/pkg/walletrpc"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
	
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/chainhash"
	"github.com/p9c/pod/pkg/txscript"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/waddrmgr"
	"github.com/p9c/pod/pkg/walletdb"
	"github.com/p9c/pod/pkg/wire"
	"github.com/p9c/pod/pkg/wtxmgr"
)

var (
	// seed is the master seed used to create extended keys.
	seed           = bytes.Repeat([]byte{0x2a, 0x64, 0xdf, 0x08}, 8)
	pubPassphrase  = []byte("_DJr{fL4H0O}*-0\n:V1izc)(6BomK")
	privPassphrase = []byte("81lUHXnOMZ@?XXd7O9xyDIWIbXX-lj")
	uniqueCounter  = uint32(0)
	// The block height where all our test inputs are created.
	TstInputsBlock = int32(10)
)

func getUniqueID() uint32 {
	return atomic.AddUint32(&uniqueCounter, 1)
}

// createWithdrawalTx creates a withdrawalTx with the given input and output amounts.
func createWithdrawalTx(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, inputAmounts []int64,
	outputAmounts []int64,
) *withdrawalTx {
	net := pool.Manager().ChainParams()
	tx := newWithdrawalTx(defaultTxOptions)
	_, credits := TstCreateCreditsOnNewSeries(t, dbtx, pool, inputAmounts)
	for _, c := range credits {
		tx.addInput(c)
	}
	for i, amount := range outputAmounts {
		request := TstNewOutputRequest(
			t, uint32(i), "34eVkREKgvvGASZW7hkgE2uNc1yycntMK6", amount2.Amount(amount), net,
		)
		tx.addOutput(request)
	}
	return tx
}
func createMsgTx(pkScript []byte, amts []int64) *wire.MsgTx {
	msgtx := &wire.MsgTx{
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 0xffffffff,
				},
				SignatureScript: []byte{txscript.OP_NOP},
				Sequence:        0xffffffff,
			},
		},
		LockTime: 0,
	}
	for _, amt := range amts {
		msgtx.AddTxOut(wire.NewTxOut(amt, pkScript))
	}
	return msgtx
}
func TstNewDepositScript(t *testing.T, p *Pool, seriesID uint32, branch Branch, idx Index) []byte {
	script, e := p.DepositScript(seriesID, branch, idx)
	if e != nil {
		t.Fatalf("Failed to create deposit script for series %d, branch %d, index %d: %v",
			seriesID, branch, idx, e,
		)
	}
	return script
}
func TstRNamespaces(tx walletdb.ReadTx) (votingpoolNs, addrmgrNs walletdb.ReadBucket) {
	return tx.ReadBucket(votingpoolNamespaceKey), tx.ReadBucket(addrmgrNamespaceKey)
}
func TstRWNamespaces(tx walletdb.ReadWriteTx) (votingpoolNs, addrmgrNs walletdb.ReadWriteBucket) {
	return tx.ReadWriteBucket(votingpoolNamespaceKey), tx.ReadWriteBucket(addrmgrNamespaceKey)
}

// func TstTxStoreRWNamespace(// 	tx walletdb.ReadWriteTx) walletdb.ReadWriteBucket {
// 	return tx.ReadWriteBucket(txmgrNamespaceKey)
// }

// TstEnsureUsedAddr ensures the addresses defined by the given series/branch and index==0..idx are present in the set
// of used addresses for the given Pool.
func TstEnsureUsedAddr(t *testing.T, dbtx walletdb.ReadWriteTx, p *Pool, seriesID uint32, branch Branch, idx Index,
) []byte {
	ns, addrmgrNs := TstRWNamespaces(dbtx)
	addr, e := p.getUsedAddr(ns, addrmgrNs, seriesID, branch, idx)
	if e != nil {
		t.Fatal(e)
	} else if addr != nil {
		var script []byte
		TstRunWithManagerUnlocked(t, p.Manager(), addrmgrNs, func() {
			script, e = addr.Script()
		},
		)
		if e != nil {
			t.Fatal(e)
		}
		return script
	}
	TstRunWithManagerUnlocked(t, p.Manager(), addrmgrNs, func() {
		e = p.EnsureUsedAddr(ns, addrmgrNs, seriesID, branch, idx)
	},
	)
	if e != nil {
		t.Fatal(e)
	}
	return TstNewDepositScript(t, p, seriesID, branch, idx)
}
func TstCreatePkScript(t *testing.T, dbtx walletdb.ReadWriteTx, p *Pool, seriesID uint32, branch Branch, idx Index,
) []byte {
	script := TstEnsureUsedAddr(t, dbtx, p, seriesID, branch, idx)
	addr, e := p.addressFor(script)
	if e != nil {
		t.Fatal(e)
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatal(e)
	}
	return pkScript
}

type TstSeriesDef struct {
	ReqSigs  uint32
	PubKeys  []string
	PrivKeys []string
	SeriesID uint32
	Inactive bool
}

// TstCreateSeries creates a new Series for every definition in the given slice of TstSeriesDef. If the definition
// includes any private keys, the Series is empowered with them.
func TstCreateSeries(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, definitions []TstSeriesDef) {
	ns, addrmgrNs := TstRWNamespaces(dbtx)
	for _, def := range definitions {
		e := pool.CreateSeries(ns, CurrentVersion, def.SeriesID, def.ReqSigs, def.PubKeys)
		if e != nil {
			t.Fatalf("Cannot creates series %d: %v", def.SeriesID, e)
		}
		TstRunWithManagerUnlocked(t, pool.Manager(), addrmgrNs, func() {
			for _, key := range def.PrivKeys {
				if e := pool.EmpowerSeries(ns, def.SeriesID, key); E.Chk(e) {
					t.Fatal(e)
				}
			}
		},
		)
		pool.Series(def.SeriesID).active = !def.Inactive
	}
}
func TstCreateMasterKey(t *testing.T, seed []byte) *hdkeychain.ExtendedKey {
	key, e := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if e != nil {
		t.Fatal(e)
	}
	return key
}

// // createMasterKeys creates count master ExtendedKeys with unique seeds.
// func createMasterKeys(	t *testing.T, count int) []*hdkeychain.ExtendedKey {
// 	keys := make([]*hdkeychain.ExtendedKey, count)
// 	for i := range keys {
// 		keys[i] = TstCreateMasterKey(t, bytes.Repeat(uint32ToBytes(getUniqueID()), 4))
// 	}
// 	return keys
// }

// TstCreateSeriesDef creates a TstSeriesDef with a unique SeriesID, the given reqSigs and the raw public/private keys
// extracted from the list of private keys. The new series will be empowered with all private keys.
func TstCreateSeriesDef(t *testing.T, pool *Pool, reqSigs uint32, keys []*hdkeychain.ExtendedKey) TstSeriesDef {
	pubKeys := make([]string, len(keys))
	privKeys := make([]string, len(keys))
	for i, key := range keys {
		privKeys[i] = key.String()
		pubkey, _ := key.Neuter()
		pubKeys[i] = pubkey.String()
	}
	seriesID := uint32(len(pool.seriesLookup)) + 1
	return TstSeriesDef{
		ReqSigs: reqSigs, SeriesID: seriesID, PubKeys: pubKeys, PrivKeys: privKeys,
	}
}
func TstCreatePoolAndTxStore(t *testing.T) (tearDown func(), db walletdb.DB, pool *Pool, store *wtxmgr.Store) {
	teardown, db, pool := TstCreatePool(t)
	store = TstCreateTxStore(t, db)
	return teardown, db, pool, store
}

// TstCreateCreditsOnNewSeries creates a new Series (with a unique ID) and a slice of credits locked to the series'
// address with branch==1 and index==0. The new Series will use a 2-of-3 configuration and will be empowered with all of
// its private keys.
func TstCreateCreditsOnNewSeries(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, amounts []int64) (uint32,
	[]Credit,
) {
	masters := []*hdkeychain.ExtendedKey{
		TstCreateMasterKey(t, bytes.Repeat(uint32ToBytes(getUniqueID()), 4)),
		TstCreateMasterKey(t, bytes.Repeat(uint32ToBytes(getUniqueID()), 4)),
		TstCreateMasterKey(t, bytes.Repeat(uint32ToBytes(getUniqueID()), 4)),
	}
	def := TstCreateSeriesDef(t, pool, 2, masters)
	TstCreateSeries(t, dbtx, pool, []TstSeriesDef{def})
	return def.SeriesID, TstCreateSeriesCredits(t, dbtx, pool, def.SeriesID, amounts)
}

// TstCreateSeriesCredits creates a new credit for every item in the amounts slice, locked to the given series' address
// with branch==1 and index==0.
func TstCreateSeriesCredits(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, seriesID uint32, amounts []int64,
) []Credit {
	addr := TstNewWithdrawalAddress(t, dbtx, pool, seriesID, Branch(1), Index(0))
	pkScript, e := txscript.PayToAddrScript(addr.addr)
	if e != nil {
		t.Fatal(e)
	}
	msgTx := createMsgTx(pkScript, amounts)
	txHash := msgTx.TxHash()
	credits := make([]Credit, len(amounts))
	for i := range msgTx.TxOut {
		c := wtxmgr.Credit{
			OutPoint: wire.OutPoint{
				Hash:  txHash,
				Index: uint32(i),
			},
			BlockMeta: wtxmgr.BlockMeta{
				Block: wtxmgr.Block{Height: TstInputsBlock},
			},
			Amount:   amount2.Amount(msgTx.TxOut[i].Value),
			PkScript: msgTx.TxOut[i].PkScript,
		}
		credits[i] = NewCredit(c, *addr)
	}
	return credits
}

// TstCreateSeriesCreditsOnStore inserts a new credit in the given store for every item in the amounts slice. These
// credits are locked to the votingpool address composed of the given seriesID, branch==1 and index==0.
func TstCreateSeriesCreditsOnStore(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, seriesID uint32,
	amounts []int64,
	store *wtxmgr.Store,
) []Credit {
	branch := Branch(1)
	idx := Index(0)
	pkScript := TstCreatePkScript(t, dbtx, pool, seriesID, branch, idx)
	eligible := make([]Credit, len(amounts))
	for i, credit := range TstCreateCreditsOnStore(t, dbtx, store, pkScript, amounts) {
		eligible[i] = NewCredit(credit, *TstNewWithdrawalAddress(t, dbtx, pool, seriesID, branch, idx))
	}
	return eligible
}

// TstCreateCreditsOnStore inserts a new credit in the given store for every item in the amounts slice.
func TstCreateCreditsOnStore(t *testing.T, dbtx walletdb.ReadWriteTx, s *wtxmgr.Store, pkScript []byte, amounts []int64,
) []wtxmgr.Credit {
	msgTx := createMsgTx(pkScript, amounts)
	meta := &wtxmgr.BlockMeta{
		Block: wtxmgr.Block{Height: TstInputsBlock},
	}
	rec, e := wtxmgr.NewTxRecordFromMsgTx(msgTx, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	txmgrNs := dbtx.ReadWriteBucket(txmgrNamespaceKey)
	if e := s.InsertTx(txmgrNs, rec, meta); E.Chk(e) {
		t.Fatal("Failed to create inputs: ", e)
	}
	credits := make([]wtxmgr.Credit, len(msgTx.TxOut))
	for i := range msgTx.TxOut {
		if e := s.AddCredit(txmgrNs, rec, meta, uint32(i), false); E.Chk(e) {
			t.Fatal("Failed to create inputs: ", e)
		}
		credits[i] = wtxmgr.Credit{
			OutPoint: wire.OutPoint{
				Hash:  rec.Hash,
				Index: uint32(i),
			},
			BlockMeta: *meta,
			Amount:    amount2.Amount(msgTx.TxOut[i].Value),
			PkScript:  msgTx.TxOut[i].PkScript,
		}
	}
	return credits
}

var (
	addrmgrNamespaceKey    = []byte("waddrmgr")
	votingpoolNamespaceKey = []byte("votingpool")
	txmgrNamespaceKey      = []byte("testtxstore")
)

// TstCreatePool creates a Pool on a fresh walletdb and returns it. It also returns a teardown function that closes the
// Manager and removes the directory used to store the database.
func TstCreatePool(t *testing.T) (tearDownFunc func(), db walletdb.DB, pool *Pool) {
	// This should be moved somewhere else eventually as not all of our tests call this function, but right now the only
	// opt would be to have the t.Parallel() call in each of our tests.
	t.Parallel()
	// Create a new wallet DB and addr manager.
	dir, e := ioutil.TempDir("", "pool_test")
	if e != nil {
		t.Fatalf("Failed to create db dir: %v", e)
	}
	db, e = walletdb.Create("bdb", filepath.Join(dir, "wallet.db"))
	if e != nil {
		t.Fatalf("Failed to create wallet DB: %v", e)
	}
	var addrMgr *waddrmgr.Manager
	e = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		addrmgrNs, e := tx.CreateTopLevelBucket(addrmgrNamespaceKey)
		if e != nil {
			return e
		}
		votingpoolNs, e := tx.CreateTopLevelBucket(votingpoolNamespaceKey)
		if e != nil {
			return e
		}
		fastScrypt := &waddrmgr.ScryptOptions{N: 16, R: 8, P: 1}
		e = waddrmgr.Create(addrmgrNs, seed, pubPassphrase, privPassphrase,
			&chaincfg.MainNetParams, fastScrypt, time.Now(),
		)
		if e != nil {
			return e
		}
		addrMgr, e = waddrmgr.Open(addrmgrNs, pubPassphrase, &chaincfg.MainNetParams)
		if e != nil {
			return e
		}
		pool, e = Create(votingpoolNs, addrMgr, []byte{0x00})
		return e
	},
	)
	if e != nil {
		t.Fatalf("Could not set up DB: %v", e)
	}
	tearDownFunc = func() {
		addrMgr.Close()
		if e := db.Close(); E.Chk(e) {
		}
		if e := os.RemoveAll(dir); E.Chk(e) {
		}
	}
	return tearDownFunc, db, pool
}
func TstCreateTxStore(t *testing.T, db walletdb.DB) *wtxmgr.Store {
	var store *wtxmgr.Store
	e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		txmgrNs, e := tx.CreateTopLevelBucket(txmgrNamespaceKey)
		if e != nil {
			return e
		}
		e = wtxmgr.Create(txmgrNs)
		if e != nil {
			return e
		}
		store, e = wtxmgr.Open(txmgrNs, &chaincfg.MainNetParams)
		return e
	},
	)
	if e != nil {
		t.Fatalf("Failed to create txmgr: %v", e)
	}
	return store
}
func TstNewOutputRequest(t *testing.T, transaction uint32, address string, amount amount2.Amount,
	net *chaincfg.Params,
) OutputRequest {
	addr, e := walletrpc.DecodeAddress(address, net)
	if e != nil {
		t.Fatalf("Unable to decode address %s", address)
	}
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		t.Fatalf("Unable to generate pkScript for %v", addr)
	}
	return OutputRequest{
		PkScript:    pkScript,
		Address:     addr,
		Amount:      amount,
		Server:      "server",
		Transaction: transaction,
	}
}

// func TstNewWithdrawalOutput(// 	r OutputRequest, status outputStatus,
// 	outpoints []OutBailmentOutpoint) *WithdrawalOutput {
// 	output := &WithdrawalOutput{
// 		request:   r,
// 		status:    status,
// 		outpoints: outpoints,
// 	}
// 	return output
// }
func TstNewWithdrawalAddress(t *testing.T, dbtx walletdb.ReadWriteTx, p *Pool, seriesID uint32, branch Branch,
	index Index,
) (addr *WithdrawalAddress) {
	TstEnsureUsedAddr(t, dbtx, p, seriesID, branch, index)
	ns, addrmgrNs := TstRNamespaces(dbtx)
	var e error
	TstRunWithManagerUnlocked(t, p.Manager(), addrmgrNs, func() {
		addr, e = p.WithdrawalAddress(ns, addrmgrNs, seriesID, branch, index)
	},
	)
	if e != nil {
		t.Fatalf("Failed to get WithdrawalAddress: %v", e)
	}
	return addr
}
func TstNewChangeAddress(t *testing.T, p *Pool, seriesID uint32, idx Index) (addr *ChangeAddress) {
	addr, e := p.ChangeAddress(seriesID, idx)
	if e != nil {
		t.Fatalf("Failed to get ChangeAddress: %v", e)
	}
	return addr
}
func TstConstantFee(fee amount2.Amount) func() amount2.Amount {
	return func() amount2.Amount { return fee }
}
func createAndFulfillWithdrawalRequests(t *testing.T, dbtx walletdb.ReadWriteTx, pool *Pool, roundID uint32,
) withdrawalInfo {
	params := pool.Manager().ChainParams()
	seriesID, eligible := TstCreateCreditsOnNewSeries(t, dbtx, pool, []int64{2e6, 4e6})
	requests := []OutputRequest{
		TstNewOutputRequest(t, 1, "34eVkREKgvvGASZW7hkgE2uNc1yycntMK6", 3e6, params),
		TstNewOutputRequest(t, 2, "3PbExiaztsSYgh6zeMswC49hLUwhTQ86XG", 2e6, params),
	}
	changeStart := TstNewChangeAddress(t, pool, seriesID, 0)
	dustThreshold := amount2.Amount(1e4)
	startAddr := TstNewWithdrawalAddress(t, dbtx, pool, seriesID, 1, 0)
	lastSeriesID := seriesID
	w := newWithdrawal(roundID, requests, eligible, *changeStart)
	if e := w.fulfillRequests(); E.Chk(e) {
		t.Fatal(e)
	}
	return withdrawalInfo{
		requests:      requests,
		startAddress:  *startAddr,
		changeStart:   *changeStart,
		lastSeriesID:  lastSeriesID,
		dustThreshold: dustThreshold,
		status:        *w.status,
	}
}
