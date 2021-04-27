package wtxmgr

import (
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pkg/wire"
)

// insertMemPoolTx inserts the unmined transaction record. It also marks previous outputs referenced by the inputs as
// spent.
func (s *Store) insertMemPoolTx(ns walletdb.ReadWriteBucket, rec *TxRecord) (e error) {
	// Chk whether the transaction has already been added to the unconfirmed bucket.
	if existsRawUnmined(ns, rec.Hash[:]) != nil {
		// TODO: compare serialized txs to ensure this isn't a hash
		//  collision?
		return nil
	}
	// Since transaction records within the store are keyed by their transaction _and_ block confirmation, we'll iterate
	// through the transaction's outputs to determine if we've already seen them to prevent from adding this transaction
	// to the unconfirmed bucket.
	for i := range rec.MsgTx.TxOut {
		k := canonicalOutPoint(&rec.Hash, uint32(i))
		if existsRawUnspent(ns, k) != nil {
			return nil
		}
	}
	I.Ln("inserting unconfirmed transaction", rec.Hash)
	v, e := valueTxRecord(rec)
	if e != nil {
		return e
	}
	e = putRawUnmined(ns, rec.Hash[:], v)
	if e != nil {
		return e
	}
	for _, input := range rec.MsgTx.TxIn {
		prevOut := &input.PreviousOutPoint
		k := canonicalOutPoint(&prevOut.Hash, prevOut.Index)
		e = putRawUnminedInput(ns, k, rec.Hash[:])
		if e != nil {
			return e
		}
	}
	// TODO: increment credit amount for each credit (but those are unknown here currently).
	return nil
}

// removeDoubleSpends checks for any unmined transactions which would introduce a double spend if tx was added to the
// store (either as a confirmed or unmined transaction). Each conflicting transaction and all transactions which spend
// it are recursively removed.
func (s *Store) removeDoubleSpends(ns walletdb.ReadWriteBucket, rec *TxRecord) (e error) {
	for _, input := range rec.MsgTx.TxIn {
		prevOut := &input.PreviousOutPoint
		prevOutKey := canonicalOutPoint(&prevOut.Hash, prevOut.Index)
		doubleSpendHashes := fetchUnminedInputSpendTxHashes(ns, prevOutKey)
		for _, doubleSpendHash := range doubleSpendHashes {
			doubleSpendVal := existsRawUnmined(ns, doubleSpendHash[:])
			// If the spending transaction spends multiple outputs from the same transaction, we'll find duplicate
			// entries within the store, so it's possible we're unable to find it if the conflicts have already been
			// removed in a previous iteration.
			if doubleSpendVal == nil {
				continue
			}
			var doubleSpend TxRecord
			doubleSpend.Hash = doubleSpendHash
			e := readRawTxRecord(
				&doubleSpend.Hash, doubleSpendVal, &doubleSpend,
			)
			if e != nil {
				return e
			}
			D.Ln(
				"removing double spending transaction", doubleSpend.Hash,
			)
			if e := RemoveConflict(ns, &doubleSpend); E.Chk(e) {
				return e
			}
		}
	}
	return nil
}

// RemoveConflict removes an unmined transaction record and all spend chains deriving from it from the store.
//
// This is designed to remove transactions that would otherwise result in double spend conflicts if left in the store,
// and to remove transactions that spend coinbase transactions on reorgs.
func RemoveConflict(ns walletdb.ReadWriteBucket, rec *TxRecord) (e error) {
	// For each potential credit for this record, each spender (if any) must be recursively removed as well. Once the
	// spenders are removed, the credit is deleted.
	for i := range rec.MsgTx.TxOut {
		k := canonicalOutPoint(&rec.Hash, uint32(i))
		spenderHashes := fetchUnminedInputSpendTxHashes(ns, k)
		for _, spenderHash := range spenderHashes {
			spenderVal := existsRawUnmined(ns, spenderHash[:])
			// If the spending transaction spends multiple outputs from the same transaction, we'll find duplicate
			// entries within the store, so it's possible we're unable to find it if the conflicts have already been
			// removed in a previous iteration.
			if spenderVal == nil {
				continue
			}
			var spender TxRecord
			spender.Hash = spenderHash
			e := readRawTxRecord(&spender.Hash, spenderVal, &spender)
			if e != nil {
				return e
			}
			D.F(
				"transaction %v is part of a removed conflict chain -- removing as well %s",
				spender.Hash,
			)
			if e := RemoveConflict(ns, &spender); E.Chk(e) {
				return e
			}
		}
		if e := deleteRawUnminedCredit(ns, k); E.Chk(e) {
			return e
		}
	}
	// If this tx spends any previous credits (either mined or unmined), set each unspent. Mined transactions are only
	// marked spent by having the output in the unmined inputs bucket.
	for _, input := range rec.MsgTx.TxIn {
		prevOut := &input.PreviousOutPoint
		k := canonicalOutPoint(&prevOut.Hash, prevOut.Index)
		if e := deleteRawUnminedInput(ns, k); E.Chk(e) {
			return e
		}
	}
	return deleteRawUnmined(ns, rec.Hash[:])
}

// UnminedTxs returns the underlying transactions for all unmined transactions which are not known to have been mined in
// a block. Transactions are guaranteed to be sorted by their dependency order.
func (s *Store) UnminedTxs(ns walletdb.ReadBucket) ([]*wire.MsgTx, error) {
	recSet, e := s.unminedTxRecords(ns)
	if e != nil {
		return nil, e
	}
	recs := dependencySort(recSet)
	txs := make([]*wire.MsgTx, 0, len(recs))
	for _, rec := range recs {
		txs = append(txs, &rec.MsgTx)
	}
	return txs, nil
}

func (s *Store) unminedTxRecords(ns walletdb.ReadBucket) (map[chainhash.Hash]*TxRecord, error) {
	unmined := make(map[chainhash.Hash]*TxRecord)
	e := ns.NestedReadBucket(bucketUnmined).ForEach(
		func(k, v []byte) (e error) {
			var txHash chainhash.Hash
			e = readRawUnminedHash(k, &txHash)
			if e != nil {
				return e
			}
			rec := new(TxRecord)
			e = readRawTxRecord(&txHash, v, rec)
			if e != nil {
				return e
			}
			unmined[rec.Hash] = rec
			return nil
		},
	)
	return unmined, e
}

// UnminedTxHashes returns the hashes of all transactions not known to have been mined in a block.
func (s *Store) UnminedTxHashes(ns walletdb.ReadBucket) ([]*chainhash.Hash, error) {
	return s.unminedTxHashes(ns)
}

func (s *Store) unminedTxHashes(ns walletdb.ReadBucket) (hashes []*chainhash.Hash, e error) {
	e = ns.NestedReadBucket(bucketUnmined).ForEach(
		func(k, v []byte) (e error) {
			hash := new(chainhash.Hash)
			e = readRawUnminedHash(k, hash)
			if e == nil {
				hashes = append(hashes, hash)
			}
			return e
		},
	)
	return hashes, e
}
