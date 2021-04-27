package wtxmgr

import (
	"bytes"
	"encoding/hex"
	"github.com/p9c/matrjoska/pkg/amt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/walletdb"
	_ "github.com/p9c/matrjoska/pkg/walletdb/bdb"
	"github.com/p9c/matrjoska/pkg/wire"
)

// Received transaction output for mainnet outpoint
// 61d3696de4c888730cbe06b0ad8ecb6d72d6108e893895aa9bc067bd7eba3fad:0
var (
	TstRecvSerializedTx, _          = hex.DecodeString("010000000114d9ff358894c486b4ae11c2a8cf7851b1df64c53d2e511278eff17c22fb7373000000008c493046022100995447baec31ee9f6d4ec0e05cb2a44f6b817a99d5f6de167d1c75354a946410022100c9ffc23b64d770b0e01e7ff4d25fbc2f1ca8091053078a247905c39fce3760b601410458b8e267add3c1e374cf40f1de02b59213a82e1d84c2b94096e22e2f09387009c96debe1d0bcb2356ffdcf65d2a83d4b34e72c62eccd8490dbf2110167783b2bffffffff0280969800000000001976a914479ed307831d0ac19ebc5f63de7d5f1a430ddb9d88ac38bfaa00000000001976a914dadf9e3484f28b385ddeaa6c575c0c0d18e9788a88ac00000000")
	TstRecvTx, _                    = util.NewTxFromBytes(TstRecvSerializedTx)
	TstRecvTxSpendingTxBlockHash, _ = chainhash.NewHashFromStr("00000000000000017188b968a371bab95aa43522665353b646e41865abae02a4")
	// TstRecvAmt                      = int64(10000000)
	TstRecvTxBlockDetails = &BlockMeta{
		Block: Block{Hash: *TstRecvTxSpendingTxBlockHash, Height: 276425},
		Time:  time.Unix(1387737310, 0),
	}
	TstRecvCurrentHeight = int32(284498) // mainnet blockchain height at time of writing
	// TstRecvTxOutConfirms       = 8074          // hardcoded number of confirmations given the above block height
	TstSpendingSerializedTx, _ = hex.DecodeString("0100000003ad3fba7ebd67c09baa9538898e10d6726dcb8eadb006be0c7388c8e46d69d361000000006b4830450220702c4fbde5532575fed44f8d6e8c3432a2a9bd8cff2f966c3a79b2245a7c88db02210095d6505a57e350720cb52b89a9b56243c15ddfcea0596aedc1ba55d9fb7d5aa0012103cccb5c48a699d3efcca6dae277fee6b82e0229ed754b742659c3acdfed2651f9ffffffffdbd36173f5610e34de5c00ed092174603761595d90190f790e79cda3e5b45bc2010000006b483045022000fa20735e5875e64d05bed43d81b867f3bd8745008d3ff4331ef1617eac7c44022100ad82261fc57faac67fc482a37b6bf18158da0971e300abf5fe2f9fd39e107f58012102d4e1caf3e022757512c204bf09ff56a9981df483aba3c74bb60d3612077c9206ffffffff65536c9d964b6f89b8ef17e83c6666641bc495cb27bab60052f76cd4556ccd0d040000006a473044022068e3886e0299ffa69a1c3ee40f8b6700f5f6d463a9cf9dbf22c055a131fc4abc02202b58957fe19ff1be7a84c458d08016c53fbddec7184ac5e633f2b282ae3420ae012103b4e411b81d32a69fb81178a8ea1abaa12f613336923ee920ffbb1b313af1f4d2ffffffff02ab233200000000001976a91418808b2fbd8d2c6d022aed5cd61f0ce6c0a4cbb688ac4741f011000000001976a914f081088a300c80ce36b717a9914ab5ec8a7d283988ac00000000")
	TstSpendingTx, _           = util.NewTxFromBytes(TstSpendingSerializedTx)
	TstSpendingTxBlockHeight   = int32(279143)
	TstSignedTxBlockHash, _    = chainhash.NewHashFromStr("00000000000000017188b968a371bab95aa43522665353b646e41865abae02a4")
	TstSignedTxBlockDetails    = &BlockMeta{
		Block: Block{Hash: *TstSignedTxBlockHash, Height: TstSpendingTxBlockHeight},
		Time:  time.Unix(1389114091, 0),
	}
)

func testDB() (walletdb.DB, func(), error) {
	tmpDir, e := ioutil.TempDir("", "wtxmgr_test")
	if e != nil {
		return nil, func() {
		}, e
	}
	db, e := walletdb.Create("bdb", filepath.Join(tmpDir, "db"))
	return db, func() {
		if e = os.RemoveAll(tmpDir); E.Chk(e) {
		}
	}, e
}

var namespaceKey = []byte("txstore")

func testStore() (*Store, walletdb.DB, func(), error) {
	var e error
	var tmpDir string
	tmpDir, e = ioutil.TempDir("", "wtxmgr_test")
	if e != nil {
		return nil, nil, func() {
		}, e
	}
	var db walletdb.DB
	db, e = walletdb.Create("bdb", filepath.Join(tmpDir, "db"))
	if e != nil {
		if e = os.RemoveAll(tmpDir); E.Chk(e) {
		}
		return nil, nil, nil, e
	}
	teardown := func() {
		if e = db.Close(); E.Chk(e) {
		}
		if e = os.RemoveAll(tmpDir); E.Chk(e) {
		}
	}
	var s *Store
	e = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
		ns, e := tx.CreateTopLevelBucket(namespaceKey)
		if e != nil {
			return e
		}
		e = Create(ns)
		if e != nil {
			return e
		}
		s, e = Open(ns, &chaincfg.TestNet3Params)
		return e
	},
	)
	return s, db, teardown, e
}
func serializeTx(tx *util.Tx) []byte {
	var buf bytes.Buffer
	e := tx.MsgTx().Serialize(&buf)
	if e != nil {
		panic(e)
	}
	return buf.Bytes()
}
func TestInsertsCreditsDebitsRollbacks(t *testing.T) {
	t.Parallel()
	// Create a double spend of the received blockchain transaction.
	dupRecvTx, _ := util.NewTxFromBytes(TstRecvSerializedTx)
	// Switch txout amount to 1 DUO. Transaction store doesn't validate txs, so this is fine for testing a double spend
	// removal.
	TstDupRecvAmount := int64(1e8)
	newDupMsgTx := dupRecvTx.MsgTx()
	newDupMsgTx.TxOut[0].Value = TstDupRecvAmount
	TstDoubleSpendTx := util.NewTx(newDupMsgTx)
	TstDoubleSpendSerializedTx := serializeTx(TstDoubleSpendTx)
	// Create a "signed" (with invalid sigs) tx that spends output 0 of the double spend.
	spendingTx := wire.NewMsgTx(wire.TxVersion)
	spendingTxIn := wire.NewTxIn(wire.NewOutPoint(TstDoubleSpendTx.Hash(), 0), []byte{0, 1, 2, 3, 4}, nil)
	spendingTx.AddTxIn(spendingTxIn)
	spendingTxOut1 := wire.NewTxOut(1e7, []byte{5, 6, 7, 8, 9})
	spendingTxOut2 := wire.NewTxOut(9e7, []byte{10, 11, 12, 13, 14})
	spendingTx.AddTxOut(spendingTxOut1)
	spendingTx.AddTxOut(spendingTxOut2)
	TstSpendingTx = util.NewTx(spendingTx)
	TstSpendingSerializedTx = serializeTx(TstSpendingTx)
	var _ = TstSpendingTx
	tests := []struct {
		name     string
		f        func(*Store, walletdb.ReadWriteBucket) (*Store, error)
		bal, unc amt.Amount
		unspents map[wire.OutPoint]struct{}
		unmined  map[chainhash.Hash]struct{}
	}{
		{
			name: "new store",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				return s, nil
			},
			bal:      0,
			unc:      0,
			unspents: map[wire.OutPoint]struct{}{},
			unmined:  map[chainhash.Hash]struct{}{},
		},
		{
			name: "txout insert",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, nil)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, nil, 0, false)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstRecvTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstRecvTx.Hash(): {},
			},
		},
		{
			name: "insert duplicate unconfirmed",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, nil)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, nil, 0, false)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstRecvTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstRecvTx.Hash(): {},
			},
		},
		{
			name: "confirmed txout insert",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstRecvTxBlockDetails)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, TstRecvTxBlockDetails, 0, false)
				return s, e
			},
			bal: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstRecvTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
		{
			name: "insert duplicate confirmed",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstRecvTxBlockDetails)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, TstRecvTxBlockDetails, 0, false)
				return s, e
			},
			bal: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstRecvTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
		{
			name: "rollback confirmed credit",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				e := s.Rollback(ns, TstRecvTxBlockDetails.Height)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstRecvTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstRecvTx.Hash(): {},
			},
		},
		{
			name: "insert confirmed double spend",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstDoubleSpendSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstRecvTxBlockDetails)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, TstRecvTxBlockDetails, 0, false)
				return s, e
			},
			bal: amt.Amount(TstDoubleSpendTx.MsgTx().TxOut[0].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstDoubleSpendTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
		{
			name: "insert unconfirmed debit",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstSpendingSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, nil)
				return s, e
			},
			bal:      0,
			unc:      0,
			unspents: map[wire.OutPoint]struct{}{},
			unmined: map[chainhash.Hash]struct{}{
				*TstSpendingTx.Hash(): {},
			},
		},
		{
			name: "insert unconfirmed debit again",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstDoubleSpendSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstRecvTxBlockDetails)
				return s, e
			},
			bal:      0,
			unc:      0,
			unspents: map[wire.OutPoint]struct{}{},
			unmined: map[chainhash.Hash]struct{}{
				*TstSpendingTx.Hash(): {},
			},
		},
		{
			name: "insert change (index 0)",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstSpendingSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, nil)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, nil, 0, true)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 0,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstSpendingTx.Hash(): {},
			},
		},
		{
			name: "insert output back to this own wallet (index 1)",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstSpendingSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, nil)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, nil, 1, true)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value + TstSpendingTx.MsgTx().TxOut[1].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 0,
				}: {},
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 1,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstSpendingTx.Hash(): {},
			},
		},
		{
			name: "confirm signed tx",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstSpendingSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstSignedTxBlockDetails)
				return s, e
			},
			bal: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value + TstSpendingTx.MsgTx().TxOut[1].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 0,
				}: {},
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 1,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
		{
			name: "rollback after spending tx",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				e := s.Rollback(ns, TstSignedTxBlockDetails.Height+1)
				return s, e
			},
			bal: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value + TstSpendingTx.MsgTx().TxOut[1].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 0,
				}: {},
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 1,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
		{
			name: "rollback spending tx block",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				e := s.Rollback(ns, TstSignedTxBlockDetails.Height)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value + TstSpendingTx.MsgTx().TxOut[1].Value),
			unspents: map[wire.OutPoint]struct{}{
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 0,
				}: {},
				{
					Hash:  *TstSpendingTx.Hash(),
					Index: 1,
				}: {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstSpendingTx.Hash(): {},
			},
		},
		{
			name: "rollback double spend tx block",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				e := s.Rollback(ns, TstRecvTxBlockDetails.Height)
				return s, e
			},
			bal: 0,
			unc: amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value + TstSpendingTx.MsgTx().TxOut[1].Value),
			unspents: map[wire.OutPoint]struct{}{
				*wire.NewOutPoint(TstSpendingTx.Hash(), 0): {},
				*wire.NewOutPoint(TstSpendingTx.Hash(), 1): {},
			},
			unmined: map[chainhash.Hash]struct{}{
				*TstDoubleSpendTx.Hash(): {},
				*TstSpendingTx.Hash():    {},
			},
		},
		{
			name: "insert original recv txout",
			f: func(s *Store, ns walletdb.ReadWriteBucket) (*Store, error) {
				rec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
				if e != nil {
					return nil, e
				}
				e = s.InsertTx(ns, rec, TstRecvTxBlockDetails)
				if e != nil {
					return nil, e
				}
				e = s.AddCredit(ns, rec, TstRecvTxBlockDetails, 0, false)
				return s, e
			},
			bal: amt.Amount(TstRecvTx.MsgTx().TxOut[0].Value),
			unc: 0,
			unspents: map[wire.OutPoint]struct{}{
				*wire.NewOutPoint(TstRecvTx.Hash(), 0): {},
			},
			unmined: map[chainhash.Hash]struct{}{},
		},
	}
	s, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	for _, test := range tests {
		e := walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(namespaceKey)
			tmpStore, e := test.f(s, ns)
			if e != nil {
				t.Fatalf("%s: got error: %v", test.name, e)
			}
			s = tmpStore
			bal, e := s.Balance(ns, 1, TstRecvCurrentHeight)
			if e != nil {
				t.Fatalf("%s: Confirmed Balance failed: %v", test.name, e)
			}
			if bal != test.bal {
				t.Fatalf("%s: balance mismatch: expected: %d, got: %d", test.name, test.bal, bal)
			}
			unc, e := s.Balance(ns, 0, TstRecvCurrentHeight)
			if e != nil {
				t.Fatalf("%s: Unconfirmed Balance failed: %v", test.name, e)
			}
			unc -= bal
			if unc != test.unc {
				t.Fatalf("%s: unconfirmed balance mismatch: expected %d, got %d", test.name, test.unc, unc)
			}
			// Chk that unspent outputs match expected.
			unspent, e := s.UnspentOutputs(ns)
			if e != nil {
				t.Fatalf("%s: failed to fetch unspent outputs: %v", test.name, e)
			}
			for _, cred := range unspent {
				if _, ok := test.unspents[cred.OutPoint]; !ok {
					t.Errorf("%s: unexpected unspent output: %v", test.name, cred.OutPoint)
				}
				delete(test.unspents, cred.OutPoint)
			}
			if len(test.unspents) != 0 {
				t.Fatalf("%s: missing expected unspent output(s)", test.name)
			}
			// Chk that unmined txs match expected.
			unmined, e := s.UnminedTxs(ns)
			if e != nil {
				t.Fatalf("%s: cannot load unmined transactions: %v", test.name, e)
			}
			for _, tx := range unmined {
				txHash := tx.TxHash()
				if _, ok := test.unmined[txHash]; !ok {
					t.Fatalf("%s: unexpected unmined tx: %v", test.name, txHash)
				}
				delete(test.unmined, txHash)
			}
			if len(test.unmined) != 0 {
				t.Fatalf("%s: missing expected unmined tx(s)", test.name)
			}
			return nil
		},
		)
		if e != nil {
			t.Fatal(e)
		}
	}
}
func TestFindingSpentCredits(t *testing.T) {
	t.Parallel()
	s, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	dbtx, e := db.BeginReadWriteTx()
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		e = dbtx.Commit()
		if e != nil {
			t.Log(e)
		}
	}()
	ns := dbtx.ReadWriteBucket(namespaceKey)
	// Insert transaction and credit which will be spent.
	recvRec, e := NewTxRecord(TstRecvSerializedTx, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, recvRec, TstRecvTxBlockDetails)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, recvRec, TstRecvTxBlockDetails, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	// Insert confirmed transaction which spends the above credit.
	spendingRec, e := NewTxRecord(TstSpendingSerializedTx, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spendingRec, TstSignedTxBlockDetails)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spendingRec, TstSignedTxBlockDetails, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	bal, e := s.Balance(ns, 1, TstSignedTxBlockDetails.Height)
	if e != nil {
		t.Fatal(e)
	}
	expectedBal := amt.Amount(TstSpendingTx.MsgTx().TxOut[0].Value)
	if bal != expectedBal {
		t.Fatalf("bad balance: %v != %v", bal, expectedBal)
	}
	unspents, e := s.UnspentOutputs(ns)
	if e != nil {
		t.Fatal(e)
	}
	op := wire.NewOutPoint(TstSpendingTx.Hash(), 0)
	if unspents[0].OutPoint != *op {
		t.Fatal("unspent outpoint doesn't match expected")
	}
	if len(unspents) > 1 {
		t.Fatal("has more than one unspent credit")
	}
}
func newCoinBase(outputValues ...int64) *wire.MsgTx {
	tx := wire.MsgTx{
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{Index: ^uint32(0)},
			},
		},
	}
	for _, val := range outputValues {
		tx.TxOut = append(tx.TxOut, &wire.TxOut{Value: val})
	}
	return &tx
}
func spendOutput(txHash *chainhash.Hash, index uint32, outputValues ...int64) *wire.MsgTx {
	tx := wire.MsgTx{
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{Hash: *txHash, Index: index},
			},
		},
	}
	for _, val := range outputValues {
		tx.TxOut = append(tx.TxOut, &wire.TxOut{Value: val})
	}
	return &tx
}
func spendOutputs(outputs []wire.OutPoint, outputValues ...int64) *wire.MsgTx {
	tx := &wire.MsgTx{}
	for _, output := range outputs {
		tx.TxIn = append(tx.TxIn, &wire.TxIn{PreviousOutPoint: output})
	}
	for _, value := range outputValues {
		tx.TxOut = append(tx.TxOut, &wire.TxOut{Value: value})
	}
	return tx
}
func TestCoinbases(t *testing.T) {
	t.Parallel()
	s, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	dbtx, e := db.BeginReadWriteTx()
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		e = dbtx.Commit()
		if e != nil {
			t.Log(e)
		}
	}()
	ns := dbtx.ReadWriteBucket(namespaceKey)
	b100 := BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	cb := newCoinBase(20e8, 10e8, 30e8)
	cbRec, e := NewTxRecordFromMsgTx(cb, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	// Insert coinbase and mark outputs 0 and 2 as credits.
	e = s.InsertTx(ns, cbRec, &b100)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, cbRec, &b100, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, cbRec, &b100, 2, false)
	if e != nil {
		t.Fatal(e)
	}
	coinbaseMaturity := int32(chaincfg.TestNet3Params.CoinbaseMaturity)
	// Balance should be 0 if the coinbase is immature, 50 DUO at and beyond
	// maturity.
	//
	// Outputs when depth is below maturity are never included, no matter
	// the required number of confirmations.  Matured outputs which have
	// greater depth than minConf are still excluded.
	type balTest struct {
		height  int32
		minConf int32
		bal     amt.Amount
	}
	balTests := []balTest{
		// Next block it is still immature
		{
			height:  b100.Height + coinbaseMaturity - 2,
			minConf: 0,
			bal:     0,
		},
		{
			height:  b100.Height + coinbaseMaturity - 2,
			minConf: coinbaseMaturity,
			bal:     0,
		},
		// Next block it matures
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: 0,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: 1,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: coinbaseMaturity - 1,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: coinbaseMaturity,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: coinbaseMaturity + 1,
			bal:     0,
		},
		// Matures at this block
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: 0,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: 1,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity + 1,
			bal:     50e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity + 2,
			bal:     0,
		},
	}
	for i, tst := range balTests {
		var bal amt.Amount
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks after inserting coinbase")
	}
	// Spend an output from the coinbase tx in an unmined transaction when
	// the next block will mature the coinbase.
	spenderATime := time.Now()
	spenderA := spendOutput(&cbRec.Hash, 0, 5e8, 15e8)
	spenderARec, e := NewTxRecordFromMsgTx(spenderA, spenderATime)
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spenderARec, nil)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderARec, nil, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	balTests = []balTest{
		// Next block it matures
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: 0,
			bal:     35e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: 1,
			bal:     30e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: coinbaseMaturity,
			bal:     30e8,
		},
		{
			height:  b100.Height + coinbaseMaturity - 1,
			minConf: coinbaseMaturity + 1,
			bal:     0,
		},
		// Matures at this block
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: 0,
			bal:     35e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: 1,
			bal:     30e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity,
			bal:     30e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity + 1,
			bal:     30e8,
		},
		{
			height:  b100.Height + coinbaseMaturity,
			minConf: coinbaseMaturity + 2,
			bal:     0,
		},
	}
	balTestsBeforeMaturity := balTests
	for i, tst := range balTests {
		var bal amt.Amount
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks after spending coinbase with unmined transaction")
	}
	// Mine the spending transaction in the block the coinbase matures.
	bMaturity := BlockMeta{
		Block: Block{Height: b100.Height + coinbaseMaturity},
		Time:  time.Now(),
	}
	e = s.InsertTx(ns, spenderARec, &bMaturity)
	if e != nil {
		t.Fatal(e)
	}
	balTests = []balTest{
		// Maturity height
		{
			height:  bMaturity.Height,
			minConf: 0,
			bal:     35e8,
		},
		{
			height:  bMaturity.Height,
			minConf: 1,
			bal:     35e8,
		},
		{
			height:  bMaturity.Height,
			minConf: 2,
			bal:     30e8,
		},
		{
			height:  bMaturity.Height,
			minConf: coinbaseMaturity,
			bal:     30e8,
		},
		{
			height:  bMaturity.Height,
			minConf: coinbaseMaturity + 1,
			bal:     30e8,
		},
		{
			height:  bMaturity.Height,
			minConf: coinbaseMaturity + 2,
			bal:     0,
		},
		// Next block after maturity height
		{
			height:  bMaturity.Height + 1,
			minConf: 0,
			bal:     35e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: 2,
			bal:     35e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: 3,
			bal:     30e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: coinbaseMaturity + 2,
			bal:     30e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: coinbaseMaturity + 3,
			bal:     0,
		},
	}
	var bal amt.Amount
	for i, tst := range balTests {
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks mining coinbase spending transaction")
	}
	// Create another spending transaction which spends the credit from the
	// first spender.  This will be used to test removing the entire
	// conflict chain when the coinbase is later reorged out.
	//
	// Use the same output amount as spender A and mark it as a credit.
	// This will mean the balance tests should report identical results.
	spenderBTime := time.Now()
	spenderB := spendOutput(&spenderARec.Hash, 0, 5e8)
	spenderBRec, e := NewTxRecordFromMsgTx(spenderB, spenderBTime)
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spenderBRec, &bMaturity)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderBRec, &bMaturity, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	for i, tst := range balTests {
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks mining second spending transaction")
	}
	// Reorg out the block that matured the coinbase and check balances
	// again.
	e = s.Rollback(ns, bMaturity.Height)
	if e != nil {
		t.Fatal(e)
	}
	balTests = balTestsBeforeMaturity
	for i, tst := range balTests {
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks after reorging maturity block")
	}
	// Reorg out the block which contained the coinbase.  There should be no
	// more transactions in the store (since the previous outputs referenced
	// by the spending tx no longer exist), and the balance will always be
	// zero.
	e = s.Rollback(ns, b100.Height)
	if e != nil {
		t.Fatal(e)
	}
	balTests = []balTest{
		// Current height
		{
			height:  b100.Height - 1,
			minConf: 0,
			bal:     0,
		},
		{
			height:  b100.Height - 1,
			minConf: 1,
			bal:     0,
		},
		// Next height
		{
			height:  b100.Height,
			minConf: 0,
			bal:     0,
		},
		{
			height:  b100.Height,
			minConf: 1,
			bal:     0,
		},
	}
	for i, tst := range balTests {
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks after reorging coinbase block")
	}
	unminedTxs, e := s.UnminedTxs(ns)
	if e != nil {
		t.Fatal(e)
	}
	if len(unminedTxs) != 0 {
		t.Fatalf("Should have no unmined transactions after coinbase reorg, found %d", len(unminedTxs))
	}
}

// Test moving multiple transactions from unmined buckets to the same block.
func TestMoveMultipleToSameBlock(t *testing.T) {
	t.Parallel()
	s, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	dbtx, e := db.BeginReadWriteTx()
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		e = dbtx.Commit()
		if e != nil {
			t.Log(e)
		}
	}()
	ns := dbtx.ReadWriteBucket(namespaceKey)
	b100 := BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	cb := newCoinBase(20e8, 30e8)
	cbRec, e := NewTxRecordFromMsgTx(cb, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	// Insert coinbase and mark both outputs as credits.
	e = s.InsertTx(ns, cbRec, &b100)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, cbRec, &b100, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, cbRec, &b100, 1, false)
	if e != nil {
		t.Fatal(e)
	}
	// Create and insert two unmined transactions which spend both coinbase
	// outputs.
	spenderATime := time.Now()
	spenderA := spendOutput(&cbRec.Hash, 0, 1e8, 2e8, 18e8)
	spenderARec, e := NewTxRecordFromMsgTx(spenderA, spenderATime)
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spenderARec, nil)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderARec, nil, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderARec, nil, 1, false)
	if e != nil {
		t.Fatal(e)
	}
	spenderBTime := time.Now()
	spenderB := spendOutput(&cbRec.Hash, 1, 4e8, 8e8, 18e8)
	spenderBRec, e := NewTxRecordFromMsgTx(spenderB, spenderBTime)
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spenderBRec, nil)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderBRec, nil, 0, false)
	if e != nil {
		t.Fatal(e)
	}
	e = s.AddCredit(ns, spenderBRec, nil, 1, false)
	if e != nil {
		t.Fatal(e)
	}
	coinbaseMaturity := int32(chaincfg.TestNet3Params.CoinbaseMaturity)
	// Mine both transactions in the block that matures the coinbase.
	bMaturity := BlockMeta{
		Block: Block{Height: b100.Height + coinbaseMaturity},
		Time:  time.Now(),
	}
	e = s.InsertTx(ns, spenderARec, &bMaturity)
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, spenderBRec, &bMaturity)
	if e != nil {
		t.Fatal(e)
	}
	// Chk that both transactions can be queried at the maturity block.
	detailsA, e := s.UniqueTxDetails(ns, &spenderARec.Hash, &bMaturity.Block)
	if e != nil {
		t.Fatal(e)
	}
	if detailsA == nil {
		t.Fatal("No details found for first spender")
	}
	detailsB, e := s.UniqueTxDetails(ns, &spenderBRec.Hash, &bMaturity.Block)
	if e != nil {
		t.Fatal(e)
	}
	if detailsB == nil {
		t.Fatal("No details found for second spender")
	}
	// Verify that the balance was correctly updated on the block record
	// append and that no unmined transactions remain.
	balTests := []struct {
		height  int32
		minConf int32
		bal     amt.Amount
	}{
		// Maturity height
		{
			height:  bMaturity.Height,
			minConf: 0,
			bal:     15e8,
		},
		{
			height:  bMaturity.Height,
			minConf: 1,
			bal:     15e8,
		},
		{
			height:  bMaturity.Height,
			minConf: 2,
			bal:     0,
		},
		// Next block after maturity height
		{
			height:  bMaturity.Height + 1,
			minConf: 0,
			bal:     15e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: 2,
			bal:     15e8,
		},
		{
			height:  bMaturity.Height + 1,
			minConf: 3,
			bal:     0,
		},
	}
	for i, tst := range balTests {
		var bal amt.Amount
		bal, e = s.Balance(ns, tst.minConf, tst.height)
		if e != nil {
			t.Fatalf("Balance test %d: Store.Balance failed: %v", i, e)
		}
		if bal != tst.bal {
			t.Errorf("Balance test %d: Got %v Expected %v", i, bal, tst.bal)
		}
	}
	if t.Failed() {
		t.Fatal("Failed balance checks after moving both coinbase spenders")
	}
	unminedTxs, e := s.UnminedTxs(ns)
	if e != nil {
		t.Fatal(e)
	}
	if len(unminedTxs) != 0 {
		t.Fatalf("Should have no unmined transactions mining both, found %d", len(unminedTxs))
	}
}

// Test the optional-ness of the serialized transaction in a TxRecord.
// NewTxRecord and NewTxRecordFromMsgTx both save the serialized transaction, so
// manually strip it out to test this code path.
func TestInsertUnserializedTx(t *testing.T) {
	t.Parallel()
	s, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	dbtx, e := db.BeginReadWriteTx()
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		e = dbtx.Commit()
		if e != nil {
			t.Log(e)
		}
	}()
	ns := dbtx.ReadWriteBucket(namespaceKey)
	tx := newCoinBase(50e8)
	rec, e := NewTxRecordFromMsgTx(tx, timeNow())
	if e != nil {
		t.Fatal(e)
	}
	b100 := makeBlockMeta(100)
	e = s.InsertTx(ns, stripSerializedTx(rec), &b100)
	if e != nil {
		t.Fatalf("Insert for stripped TxRecord failed: %v", e)
	}
	// Ensure it can be retreived successfully.
	details, e := s.UniqueTxDetails(ns, &rec.Hash, &b100.Block)
	if e != nil {
		t.Fatal(e)
	}
	rec2, e := NewTxRecordFromMsgTx(&details.MsgTx, rec.Received)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(rec.SerializedTx, rec2.SerializedTx) {
		t.Fatal("Serialized txs for coinbase do not match")
	}
	// Now test that path with an unmined transaction.
	tx = spendOutput(&rec.Hash, 0, 50e8)
	rec, e = NewTxRecordFromMsgTx(tx, timeNow())
	if e != nil {
		t.Fatal(e)
	}
	e = s.InsertTx(ns, rec, nil)
	if e != nil {
		t.Fatal(e)
	}
	details, e = s.UniqueTxDetails(ns, &rec.Hash, nil)
	if e != nil {
		t.Fatal(e)
	}
	rec2, e = NewTxRecordFromMsgTx(&details.MsgTx, rec.Received)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(rec.SerializedTx, rec2.SerializedTx) {
		t.Fatal("Serialized txs for coinbase spender do not match")
	}
}

// TestRemoveUnminedTx tests that if we add an umined transaction, then we're
// able to remove that unmined transaction later along with any of its
// descendants. Any balance modifications due to the unmined transaction should
// be revered.
func TestRemoveUnminedTx(t *testing.T) {
	t.Parallel()
	store, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	// In order to reproduce real-world scenarios, we'll use a new database transaction for each interaction with the
	// wallet.
	//
	// We'll start off the test by creating a new coinbase output at height 100 and inserting it into the store.
	b100 := &BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	initialBalance := int64(1e8)
	cb := newCoinBase(initialBalance)
	cbRec, e := NewTxRecordFromMsgTx(cb, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, cbRec, b100); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, cbRec, b100, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Determine the maturity height for the coinbase output created.
	coinbaseMaturity := int32(chaincfg.TestNet3Params.CoinbaseMaturity)
	maturityHeight := b100.Block.Height + coinbaseMaturity
	// checkBalance is a helper function that compares the balance of the store with the expected value. The
	// includeUnconfirmed boolean can be used to include the unconfirmed balance as a part of the total balance.
	checkBalance := func(expectedBalance amt.Amount,
		includeUnconfirmed bool,
	) {
		t.Helper()
		minConfs := int32(1)
		if includeUnconfirmed {
			minConfs = 0
		}
		commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
			t.Helper()
			var b amt.Amount
			b, e = store.Balance(ns, minConfs, maturityHeight)
			if e != nil {
				t.Fatalf("unable to retrieve balance: %v", e)
			}
			if b != expectedBalance {
				t.Fatalf("expected balance of %d, got %d",
					expectedBalance, b,
				)
			}
		},
		)
	}
	// Since we don't have any unconfirmed transactions within the store, the total balance reflecting confirmed and
	// unconfirmed outputs should match the initial balance.
	checkBalance(amt.Amount(initialBalance), false)
	checkBalance(amt.Amount(initialBalance), true)
	// Then, we'll create an unconfirmed spend for the coinbase output and insert it into the store.
	b101 := &BlockMeta{
		Block: Block{Height: 201},
		Time:  time.Now(),
	}
	changeAmount := int64(4e7)
	spendTx := spendOutput(&cbRec.Hash, 0, 5e7, changeAmount)
	spendTxRec, e := NewTxRecordFromMsgTx(spendTx, b101.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, spendTxRec, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, spendTxRec, nil, 1, true)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// With the unconfirmed spend inserted into the store, we'll query it for its unconfirmed tranasctions to ensure it
	// was properly added.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		unminedTxs, e := store.UnminedTxs(ns)
		if e != nil {
			t.Fatalf("unable to query for unmined txs: %v", e)
		}
		if len(unminedTxs) != 1 {
			t.Fatalf("expected 1 mined tx, instead got %v",
				len(unminedTxs),
			)
		}
		unminedTxHash := unminedTxs[0].TxHash()
		spendTxHash := spendTx.TxHash()
		if !unminedTxHash.IsEqual(&spendTxHash) {
			t.Fatalf("mismatch tx hashes: expected %v, got %v",
				spendTxHash, unminedTxHash,
			)
		}
	},
	)
	// Now that an unconfirmed spend exists, there should no longer be any confirmed balance. The total balance should
	// now all be unconfirmed and it should match the change amount of the unconfirmed spend tranasction.
	checkBalance(0, false)
	checkBalance(amt.Amount(changeAmount), true)
	// Now, we'll remove the unconfirmed spend tranaction from the store.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.RemoveUnminedTx(ns, spendTxRec); E.Chk(e) {
			t.Fatal(e)
		}
	},
	)
	// We'll query the store one last time for its unconfirmed transactions to ensure the unconfirmed spend was properly
	// removed above.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		unminedTxs, e := store.UnminedTxs(ns)
		if e != nil {
			t.Fatalf("unable to query for unmined txs: %v", e)
		}
		if len(unminedTxs) != 0 {
			t.Fatalf("expected 0 mined txs, instead got %v",
				len(unminedTxs),
			)
		}
	},
	)
	// Finally, the total balance (including confirmed and unconfirmed) should once again match the initial balance, as
	// the uncofirmed spend has already been removed.
	checkBalance(amt.Amount(initialBalance), false)
	checkBalance(amt.Amount(initialBalance), true)
}

// commitDBTx is a helper function that allows us to perform multiple operations on a specific database's bucket as a
// single atomic operation.
func commitDBTx(t *testing.T, store *Store, db walletdb.DB,
	f func(walletdb.ReadWriteBucket),
) {
	t.Helper()
	dbTx, e := db.BeginReadWriteTx()
	if e != nil {
		t.Fatal(e)
	}
	defer func() {
		e := dbTx.Commit()
		if e != nil {
			t.Log(e)
		}
	}()
	ns := dbTx.ReadWriteBucket(namespaceKey)
	f(ns)
}

// testInsertDoubleSpendTx is a helper test which double spends an output. The boolean parameter indicates whether the
// first spending transaction should be the one confirmed. This test ensures that when a double spend occurs and both
// spending transactions are present in the mempool, if one of them confirms, then the remaining conflicting
// transactions within the mempool should be removed from the wallet's store.
func testInsertMempoolDoubleSpendTx(t *testing.T, first bool) {
	store, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	// In order to reproduce real-world scenarios, we'll use a new database transaction for each interaction with the
	// wallet.
	//
	// We'll start off the test by creating a new coinbase output at height 100 and inserting it into the store.
	b100 := BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	cb := newCoinBase(1e8)
	cbRec, e := NewTxRecordFromMsgTx(cb, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, cbRec, &b100); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, cbRec, &b100, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Then, we'll create two spends from the same coinbase output, in order to replicate a double spend scenario.
	firstSpend := spendOutput(&cbRec.Hash, 0, 5e7, 5e7)
	firstSpendRec, e := NewTxRecordFromMsgTx(firstSpend, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	secondSpend := spendOutput(&cbRec.Hash, 0, 4e7, 6e7)
	secondSpendRec, e := NewTxRecordFromMsgTx(secondSpend, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	// We'll insert both of them into the store without confirming them.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, firstSpendRec, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, firstSpendRec, nil, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, secondSpendRec, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, secondSpendRec, nil, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Ensure that both spends are found within the unconfirmed transactions in the wallet's store.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		unminedTxs, e := store.UnminedTxs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(unminedTxs) != 2 {
			t.Fatalf("expected 2 unmined txs, got %v",
				len(unminedTxs),
			)
		}
	},
	)
	// Then, we'll confirm either the first or second spend, depending on the boolean passed, with a height deep enough
	// that allows us to succesfully spend the coinbase output.
	coinbaseMaturity := int32(chaincfg.TestNet3Params.CoinbaseMaturity)
	bMaturity := BlockMeta{
		Block: Block{Height: b100.Height + coinbaseMaturity},
		Time:  time.Now(),
	}
	var confirmedSpendRec *TxRecord
	if first {
		confirmedSpendRec = firstSpendRec
	} else {
		confirmedSpendRec = secondSpendRec
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		e := store.InsertTx(ns, confirmedSpendRec, &bMaturity)
		if e != nil {
			t.Fatal(e)
		}
		e = store.AddCredit(
			ns, confirmedSpendRec, &bMaturity, 0, false,
		)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// This should now trigger the store to remove any other pending double spends for this coinbase output, as we've
	// already seen the confirmed one. Therefore, we shouldn't see any other unconfirmed transactions within it. We also
	// ensure that the transaction that confirmed and is now listed as a UTXO within the wallet is the correct one.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		unminedTxs, e := store.UnminedTxs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(unminedTxs) != 0 {
			t.Fatalf("expected 0 unmined txs, got %v",
				len(unminedTxs),
			)
		}
		minedTxs, e := store.UnspentOutputs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(minedTxs) != 1 {
			t.Fatalf("expected 1 mined tx, got %v", len(minedTxs))
		}
		if !minedTxs[0].Hash.IsEqual(&confirmedSpendRec.Hash) {
			t.Fatalf("expected confirmed tx hash %v, got %v",
				confirmedSpendRec.Hash, minedTxs[0].Hash,
			)
		}
	},
	)
}

// TestInsertMempoolDoubleSpendConfirmedFirstTx ensures that when a double spend occurs and both spending transactions
// are present in the mempool, if the first spend seen is confirmed, then the second spend transaction within the
// mempool should be removed from the wallet's store.
func TestInsertMempoolDoubleSpendConfirmedFirstTx(t *testing.T) {
	t.Parallel()
	testInsertMempoolDoubleSpendTx(t, true)
}

// TestInsertMempoolDoubleSpendConfirmedFirstTx ensures that when a double spend occurs and both spending transactions
// are present in the mempool, if the second spend seen is confirmed, then the first spend transaction within the
// mempool should be removed from the wallet's store.
func TestInsertMempoolDoubleSpendConfirmSecondTx(t *testing.T) {
	t.Parallel()
	testInsertMempoolDoubleSpendTx(t, false)
}

// TestInsertConfirmedDoubleSpendTx tests that when one or more double spends occur and a spending transaction confirms
// that was not known to the wallet, then the unconfirmed double spends within the mempool should be removed from the
// wallet's store.
func TestInsertConfirmedDoubleSpendTx(t *testing.T) {
	t.Parallel()
	store, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	// In order to reproduce real-world scenarios, we'll use a new database transaction for each interaction with the
	// wallet.
	//
	// We'll start off the test by creating a new coinbase output at height 100 and inserting it into the store.
	b100 := BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	cb1 := newCoinBase(1e8)
	cbRec1, e := NewTxRecordFromMsgTx(cb1, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, cbRec1, &b100); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, cbRec1, &b100, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Then, we'll create three spends from the same coinbase output. The first two will remain unconfirmed, while the
	// last should confirm and remove the remaining unconfirmed from the wallet's store.
	firstSpend1 := spendOutput(&cbRec1.Hash, 0, 5e7)
	firstSpendRec1, e := NewTxRecordFromMsgTx(firstSpend1, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, firstSpendRec1, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, firstSpendRec1, nil, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	secondSpend1 := spendOutput(&cbRec1.Hash, 0, 4e7)
	secondSpendRec1, e := NewTxRecordFromMsgTx(secondSpend1, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, secondSpendRec1, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, secondSpendRec1, nil, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// We'll also create another output and have one unconfirmed and one confirmed spending transaction also spend it.
	cb2 := newCoinBase(2e8)
	cbRec2, e := NewTxRecordFromMsgTx(cb2, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, cbRec2, &b100); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, cbRec2, &b100, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	firstSpend2 := spendOutput(&cbRec2.Hash, 0, 5e7)
	firstSpendRec2, e := NewTxRecordFromMsgTx(firstSpend2, time.Now())
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, firstSpendRec2, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, firstSpendRec2, nil, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// At this point, we should see all unconfirmed transactions within the store.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		var unminedTxs []*wire.MsgTx
		unminedTxs, e = store.UnminedTxs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(unminedTxs) != 3 {
			t.Fatalf("expected 3 unmined txs, got %d",
				len(unminedTxs),
			)
		}
	},
	)
	// Then, we'll insert the confirmed spend at a height deep enough that allows us to successfully spend the coinbase
	// outputs.
	coinbaseMaturity := int32(chaincfg.TestNet3Params.CoinbaseMaturity)
	bMaturity := BlockMeta{
		Block: Block{Height: b100.Height + coinbaseMaturity},
		Time:  time.Now(),
	}
	outputsToSpend := []wire.OutPoint{
		{Hash: cbRec1.Hash, Index: 0},
		{Hash: cbRec2.Hash, Index: 0},
	}
	confirmedSpend := spendOutputs(outputsToSpend, 3e7)
	confirmedSpendRec, e := NewTxRecordFromMsgTx(
		confirmedSpend, bMaturity.Time,
	)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		e := store.InsertTx(ns, confirmedSpendRec, &bMaturity)
		if e != nil {
			t.Fatal(e)
		}
		e = store.AddCredit(
			ns, confirmedSpendRec, &bMaturity, 0, false,
		)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Now that the confirmed spend exists within the store, we should no longer see the unconfirmed spends within it.
	// We also ensure that the transaction that confirmed and is now listed as a UTXO within the wallet is the correct
	// one.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		unminedTxs, e := store.UnminedTxs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(unminedTxs) != 0 {
			t.Fatalf("expected 0 unmined txs, got %v",
				len(unminedTxs),
			)
		}
		minedTxs, e := store.UnspentOutputs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(minedTxs) != 1 {
			t.Fatalf("expected 1 mined tx, got %v", len(minedTxs))
		}
		if !minedTxs[0].Hash.IsEqual(&confirmedSpendRec.Hash) {
			t.Fatalf("expected confirmed tx hash %v, got %v",
				confirmedSpend, minedTxs[0].Hash,
			)
		}
	},
	)
}

// TestAddDuplicateCreditAfterConfirm aims to test the case where a duplicate unconfirmed credit is added to the store
// after the intial credit has already confirmed. This can lead to outputs being duplicated in the store, which can lead
// to creating double spends when querying the wallet's UTXO set.
func TestAddDuplicateCreditAfterConfirm(t *testing.T) {
	t.Parallel()
	store, db, teardown, e := testStore()
	if e != nil {
		t.Fatal(e)
	}
	defer teardown()
	// In order to reproduce real-world scenarios, we'll use a new database transaction for each interaction with the
	// wallet.
	//
	// We'll start off the test by creating a new coinbase output at height 100 and inserting it into the store.
	b100 := &BlockMeta{
		Block: Block{Height: 100},
		Time:  time.Now(),
	}
	cb := newCoinBase(1e8)
	cbRec, e := NewTxRecordFromMsgTx(cb, b100.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e = store.InsertTx(ns, cbRec, b100); E.Chk(e) {
			t.Fatal(e)
		}
		e = store.AddCredit(ns, cbRec, b100, 0, false)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// We'll confirm that there is one unspent output in the store, which should be the coinbase output created above.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		var minedTxs []Credit
		minedTxs, e = store.UnspentOutputs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(minedTxs) != 1 {
			t.Fatalf("expected 1 mined tx, got %v", len(minedTxs))
		}
		if !minedTxs[0].Hash.IsEqual(&cbRec.Hash) {
			t.Fatalf("expected tx hash %v, got %v", cbRec.Hash,
				minedTxs[0].Hash,
			)
		}
	},
	)
	// Then, we'll create an unconfirmed spend for the coinbase output.
	b101 := &BlockMeta{
		Block: Block{Height: 101},
		Time:  time.Now(),
	}
	spendTx := spendOutput(&cbRec.Hash, 0, 5e7, 4e7)
	spendTxRec, e := NewTxRecordFromMsgTx(spendTx, b101.Time)
	if e != nil {
		t.Fatal(e)
	}
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, spendTxRec, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, spendTxRec, nil, 1, true)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Confirm the spending transaction at the next height.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, spendTxRec, b101); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, spendTxRec, b101, 1, true)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// We should see one unspent output within the store once again, this time being the change output of the spending
	// transaction.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		minedTxs, e := store.UnspentOutputs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(minedTxs) != 1 {
			t.Fatalf("expected 1 mined txs, got %v", len(minedTxs))
		}
		if !minedTxs[0].Hash.IsEqual(&spendTxRec.Hash) {
			t.Fatalf("expected tx hash %v, got %v", spendTxRec.Hash,
				minedTxs[0].Hash,
			)
		}
	},
	)
	// Now, we'll insert the spending transaction once again, this time as unconfirmed. This can happen if the backend
	// happens to forward an unconfirmed chain.RelevantTx notification to the client even after it has confirmed, which
	// results in us adding it to the store once again.
	//
	// TODO(wilmer): ideally this shouldn't happen, so we should identify the real reason for this.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		if e := store.InsertTx(ns, spendTxRec, nil); E.Chk(e) {
			t.Fatal(e)
		}
		e := store.AddCredit(ns, spendTxRec, nil, 1, true)
		if e != nil {
			t.Fatal(e)
		}
	},
	)
	// Finally, we'll ensure the change output is still the only unspent output within the store.
	commitDBTx(t, store, db, func(ns walletdb.ReadWriteBucket) {
		minedTxs, e := store.UnspentOutputs(ns)
		if e != nil {
			t.Fatal(e)
		}
		if len(minedTxs) != 1 {
			t.Fatalf("expected 1 mined txs, got %v", len(minedTxs))
		}
		if !minedTxs[0].Hash.IsEqual(&spendTxRec.Hash) {
			t.Fatalf("expected tx hash %v, got %v", spendTxRec.Hash,
				minedTxs[0].Hash,
			)
		}
	},
	)
}
