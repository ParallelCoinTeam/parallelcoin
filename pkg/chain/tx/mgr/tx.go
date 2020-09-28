package wtxmgr

import (
	"bytes"
	"github.com/p9c/pkg/app/slog"
	"time"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/db/walletdb"
	"github.com/p9c/pod/pkg/util"
)

type (
	// Block contains the minimum amount of data to uniquely identify any
	// block on
	// either the best or side chain.
	Block struct {
		Hash   chainhash.Hash
		Height int32
	}
	// BlockMeta contains the unique identification for a block and any metadata
	// pertaining to the block.  At the moment, this additional metadata only
	// includes the block time from the block header.
	BlockMeta struct {
		Block
		Time time.Time
	}
	// blockRecord is an in-memory representation of the block record saved
	// in the
	// database.
	blockRecord struct {
		Block
		Time         time.Time
		transactions []chainhash.Hash
	}
	// incidence records the block hash and blockchain height of a mined
	// transaction.
	// Since a transaction hash alone is not enough to uniquely identify a mined
	// transaction (duplicate transaction hashes are allowed),
	// the incidence is used
	// instead.
	incidence struct {
		txHash chainhash.Hash
		block  Block
	}
	// indexedIncidence records the transaction incidence and an input or output
	// index.
	indexedIncidence struct {
		incidence
		index uint32
	}
	// // debit records the debits a transaction record makes from previous
	// wallet
	// // transaction credits.
	// debit struct {
	// 	// txHash chainhash.Hash
	// 	// index  uint32
	// 	// amount util.Amount
	// 	// spends indexedIncidence
	// }
	// credit describes a transaction output which was or is spendable by
	// wallet.
	credit struct {
		outPoint wire.OutPoint
		block    Block
		amount   util.Amount
		change   bool
		spentBy  indexedIncidence // Index == ^uint32(0) if unspent
	}
	// TxRecord represents a transaction managed by the Store.
	TxRecord struct {
		MsgTx        wire.MsgTx
		Hash         chainhash.Hash
		Received     time.Time
		SerializedTx []byte // Optional: may be nil
	}
	// Credit is the type representing a transaction output which was spent or
	// is still spendable by wallet.  A UTXO is an unspent Credit, but not all
	// Credits are UTXOs.
	Credit struct {
		wire.OutPoint
		BlockMeta
		Amount       util.Amount
		PkScript     []byte
		Received     time.Time
		FromCoinBase bool
	}
	// Store implements a transaction store for storing and managing wallet
	// transactions.
	Store struct {
		chainParams *netparams.Params
		// Event callbacks.  These execute in the same goroutine as the wtxmgr
		// caller.
		NotifyUnspent func(hash *chainhash.Hash, index uint32)
	}
)

// NewTxRecord creates a new transaction record that may be inserted into the
// store.  It uses memoization to save the transaction hash and the serialized
// transaction.
func NewTxRecord(serializedTx []byte, received time.Time) (rec *TxRecord, err error) {
	rec = &TxRecord{
		Received:     received,
		SerializedTx: serializedTx,
	}
	if err = rec.MsgTx.Deserialize(bytes.NewReader(serializedTx)); slog.Check(err) {
		str := "failed to deserialize transaction"
		err = storeError(ErrInput, str, err)
		slog.Debug(err)
		return
	}
	copy(rec.Hash[:], chainhash.DoubleHashB(serializedTx))
	return
}

// NewTxRecordFromMsgTx creates a new transaction record that may be inserted
// into the store.
func NewTxRecordFromMsgTx(msgTx *wire.MsgTx, received time.Time) (rec *TxRecord, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	if err = msgTx.Serialize(buf); slog.Check(err) {
		str := "failed to serialize transaction"
		err = storeError(ErrInput, str, err)
		slog.Debug(err)
		return
	}
	rec = &TxRecord{
		MsgTx:        *msgTx,
		Received:     received,
		SerializedTx: buf.Bytes(),
		Hash:         msgTx.TxHash(),
	}
	return
}

// DoUpgrades performs any necessary upgrades to the transaction history contained in the wallet database, namespaced by
// the top level bucket key namespaceKey.
func DoUpgrades(db walletdb.DB, namespaceKey []byte) (err error) {
	// No upgrades
	return
}

// Open opens the wallet transaction store from a walletdb namespace.
// If the store does not exist, ErrNoExist is returned.
func Open(ns walletdb.ReadBucket, chainParams *netparams.Params) (s *Store, err error) {
	// Open the store.
	if err = openStore(ns); slog.Check(err) {
		return
	}
	s = &Store{chainParams, nil} // TODO: set callbacks
	return
}

// Create creates a new persistent transaction store in the walletdb
// namespace.
// Creating the store when one already exists in this namespace will error with
// ErrAlreadyExists.
func Create(ns walletdb.ReadWriteBucket) (err error) {
	return createStore(ns)
}

// updateMinedBalance updates the mined balance within the store,
// if changed, after processing the given transaction record.
func (s *Store) updateMinedBalance(ns walletdb.ReadWriteBucket, rec *TxRecord, block *BlockMeta) (err error) {
	// Fetch the mined balance in case we need to update it.
	var minedBalance util.Amount
	if minedBalance, err = fetchMinedBalance(ns); slog.Check(err) {
		return
	}
	// Add a debit record for each unspent credit spent by this transaction.
	// The index is set in each iteration below.
	spender := indexedIncidence{
		incidence: incidence{
			txHash: rec.Hash,
			block:  block.Block,
		},
	}
	newMinedBalance := minedBalance
	for i, input := range rec.MsgTx.TxIn {
		unspentKey, credKey := existsUnspent(ns, &input.PreviousOutPoint)
		if credKey == nil {
			// Debits for unmined transactions are not explicitly
			// tracked.  Instead, all previous outputs spent by any
			// unmined transaction are added to a map for quick
			// lookups when it must be checked whether a mined
			// output is unspent or not.
			//
			// Tracking individual debits for unmined transactions
			// could be added later to simplify (and increase
			// performance of) determining some details that need
			// the previous outputs (e.g. determining a fee), but at
			// the moment that is not done (and a db lookup is used
			// for those cases instead).  There is also a good
			// chance that all unmined transaction handling will
			// move entirely to the db rather than being handled in
			// memory for atomicity reasons, so the simplist
			// implementation is currently used.
			continue
		}
		// If this output is relevant to us, we'll mark the it as spent
		// and remove its amount from the store.
		spender.index = uint32(i)
		var amt util.Amount
		if amt, err = spendCredit(ns, credKey, &spender); slog.Check(err) {
			return
		}
		if err = putDebit(
			ns, &rec.Hash, uint32(i), amt, &block.Block, credKey,
		); slog.Check(err) {
			return
		}
		if err = deleteRawUnspent(ns, unspentKey); slog.Check(err) {
			return
		}
		newMinedBalance -= amt
	}
	// For each output of the record that is marked as a credit, if the
	// output is marked as a credit by the unconfirmed store, remove the
	// marker and mark the output as a credit in the db.
	//
	// Moved credits are added as unspents, even if there is another
	// unconfirmed transaction which spends them.
	cred := credit{
		outPoint: wire.OutPoint{Hash: rec.Hash},
		block:    block.Block,
		spentBy:  indexedIncidence{index: ^uint32(0)},
	}
	it := makeUnminedCreditIterator(ns, &rec.Hash)
	var amount util.Amount
	var change bool
	var index uint32
	for it.next() {
		// TODO: This should use the raw apis.  The credit value (it.cv)
		// can be moved from unmined directly to the credits bucket.
		// The key needs a modification to include the block
		// height/hash.
		if index, err = fetchRawUnminedCreditIndex(it.ck); slog.Check(err) {
			return
		}
		if amount, change, err = fetchRawUnminedCreditAmountChange(it.cv); slog.Check(err) {
			return
		}
		cred.outPoint.Index = index
		cred.amount = amount
		cred.change = change
		if err = putUnspentCredit(ns, &cred); slog.Check(err) {
			return
		}
		if err = putUnspent(ns, &cred.outPoint, &block.Block); slog.Check(err) {
			return
		}
		newMinedBalance += amount
	}
	if slog.Check(it.err) {
		return it.err
	}
	// Update the balance if it has changed.
	if newMinedBalance != minedBalance {
		return putMinedBalance(ns, newMinedBalance)
	}
	return
}

// deleteUnminedTx deletes an unmined transaction from the store.
// NOTE: This should only be used once the transaction has been mined.
func (s *Store) deleteUnminedTx(ns walletdb.ReadWriteBucket, rec *TxRecord) (err error) {
	for i := range rec.MsgTx.TxOut {
		k := canonicalOutPoint(&rec.Hash, uint32(i))
		if err = deleteRawUnminedCredit(ns, k); slog.Check(err) {
			return
		}
	}
	return deleteRawUnmined(ns, rec.Hash[:])
}

// InsertTx records a transaction as belonging to a wallet's transaction
// history.  If block is nil, the transaction is considered unspent, and the
// transaction's index must be unset.
func (s *Store) InsertTx(ns walletdb.ReadWriteBucket, rec *TxRecord, block *BlockMeta) (err error) {
	if block == nil {
		return s.insertMemPoolTx(ns, rec)
	}
	return s.insertMinedTx(ns, rec, block)
}

// RemoveUnminedTx attempts to remove an unmined transaction from the
// transaction store. This is to be used in the scenario that a transaction
// that we attempt to rebroadcast, turns out to double spend one of our
// existing inputs. This function we remove the conflicting transaction
// identified by the tx record, and also recursively remove all transactions
// that depend on it.
func (s *Store) RemoveUnminedTx(ns walletdb.ReadWriteBucket, rec *TxRecord) (err error) {
	// As we already have a tx record, we can directly call the
	// RemoveConflict method. This will do the job of recursively removing
	// this unmined transaction, and any transactions that depend on it.
	return RemoveConflict(ns, rec)
}

// insertMinedTx inserts a new transaction record for a mined
// transaction into the database under the confirmed bucket.
// It guarantees that, if the tranasction was previously unconfirmed,
// then it will take care of cleaning up the unconfirmed state.
// All other unconfirmed double spend attempts will be removed as well.
func (s *Store) insertMinedTx(ns walletdb.ReadWriteBucket, rec *TxRecord,
	block *BlockMeta) (err error) {
	// If a transaction record for this hash and block already exists, we
	// can exit early.
	if _, v := existsTxRecord(ns, &rec.Hash, &block.Block); v != nil {
		return nil
	}
	// If a block record does not yet exist for any transactions from this
	// block, insert a block record first. Otherwise, update it by adding
	// the transaction hash to the set of transactions from this block.
	blockKey, blockValue := existsBlockRecord(ns, block.Height)
	if blockValue == nil {
		err = putBlockRecord(ns, block, &rec.Hash)
	} else {
		if blockValue, err = appendRawBlockRecord(blockValue, &rec.Hash); slog.Check(err) {
			return
		}
		err = putRawBlockRecord(ns, blockKey, blockValue)
	}
	if slog.Check(err) {
		return
	}
	if err = putTxRecord(ns, rec, &block.Block); slog.Check(err) {
		return
	}
	// Determine if this transaction has affected our balance, and if so,
	// update it.
	if err = s.updateMinedBalance(ns, rec, block); slog.Check(err) {
		return
	}
	// If this transaction previously existed within the store as unmined,
	// we'll need to remove it from the unmined bucket.
	if v := existsRawUnmined(ns, rec.Hash[:]); v != nil {
		slog.Infof("marking unconfirmed transaction %v mined in block %d", &rec.Hash, block.Height)
		if err = s.deleteUnminedTx(ns, rec); slog.Check(err) {
			return
		}
	}
	// As there may be unconfirmed transactions that are invalidated by this
	// transaction (either being duplicates, or double spends), remove them
	// from the unconfirmed set.  This also handles removing unconfirmed
	// transaction spend chains if any other unconfirmed transactions spend
	// outputs of the removed double spend.
	return s.removeDoubleSpends(ns, rec)
}

// AddCredit marks a transaction record as containing a transaction output
// spendable by wallet.  The output is added unspent, and is marked spent
// when a new transaction spending the output is inserted into the store.
//
// TODO(jrick): This should not be necessary.  Instead, pass the indexes
// that are known to contain credits when a transaction or merkleblock is
// inserted into the store.
func (s *Store) AddCredit(ns walletdb.ReadWriteBucket, rec *TxRecord, block *BlockMeta, index uint32, change bool) (err error) {
	if int(index) >= len(rec.MsgTx.TxOut) {
		str := "transaction output does not exist"
		err = storeError(ErrInput, str, nil)
		slog.Debug(err)
		return
	}
	var isNew bool
	if isNew, err = s.addCredit(ns, rec, block, index, change); slog.Check(err) && isNew && s.NotifyUnspent != nil {
		s.NotifyUnspent(&rec.Hash, index)
	}
	return err
}

func // addCredit is an AddCredit helper that runs in an update transaction.
// The bool return specifies whether the unspent output is newly added (
// true) or a duplicate (false).
(s *Store) addCredit(ns walletdb.ReadWriteBucket, rec *TxRecord, block *BlockMeta, index uint32, change bool) (b bool, err error) {
	if block == nil {
		// If the outpoint that we should mark as credit already exists
		// within the store, either as unconfirmed or confirmed, then we
		// have nothing left to do and can exit.
		k := canonicalOutPoint(&rec.Hash, index)
		if existsRawUnminedCredit(ns, k) != nil {
			return
		}
		if existsRawUnspent(ns, k) != nil {
			return
		}
		v := valueUnminedCredit(util.Amount(rec.MsgTx.TxOut[index].Value), change)
		return true, putRawUnminedCredit(ns, k, v)
	}
	var k, v []byte
	if k, v = existsCredit(ns, &rec.Hash, index, &block.Block); v != nil {
		return
	}
	txOutAmt := util.Amount(rec.MsgTx.TxOut[index].Value)
	slog.Tracef("marking transaction %v output %d (%v) spendable", rec.Hash, index, txOutAmt)
	cred := credit{
		outPoint: wire.OutPoint{
			Hash:  rec.Hash,
			Index: index,
		},
		block:   block.Block,
		amount:  txOutAmt,
		change:  change,
		spentBy: indexedIncidence{index: ^uint32(0)},
	}
	v = valueUnspentCredit(&cred)
	if err = putRawCredit(ns, k, v); slog.Check(err) {
		return
	}
	var minedBalance util.Amount
	if minedBalance, err = fetchMinedBalance(ns); slog.Check(err) {
		return
	}
	if err = putMinedBalance(ns, minedBalance+txOutAmt); slog.Check(err) {
		return
	}
	return true, putUnspent(ns, &cred.outPoint, &block.Block)
}

// Rollback removes all blocks at height onwards,
// moving any transactions within each block to the unconfirmed pool.
func (s *Store) Rollback(ns walletdb.ReadWriteBucket, height int32) (err error) {
	return s.rollback(ns, height)
}

func (s *Store) rollback(ns walletdb.ReadWriteBucket, height int32) (err error) {
	minedBalance, err := fetchMinedBalance(ns)
	if err != nil {
		slog.Error(err)
		return err
	}
	// Keep track of all credits that were removed from coinbase
	// transactions.  After detaching all blocks, if any transaction record
	// exists in unmined that spends these outputs, remove them and their
	// spend chains.
	//
	// It is necessary to keep these in memory and fix the unmined
	// transactions later since blocks are removed in increasing order.
	var coinBaseCredits []wire.OutPoint
	var heightsToRemove []int32
	it := makeReverseBlockIterator(ns)
	for it.prev() {
		b := &it.elem
		if it.elem.Height < height {
			break
		}
		heightsToRemove = append(heightsToRemove, it.elem.Height)
		slog.Tracef(
			"rolling back %d transactions from block %v height %d",
			len(b.transactions), b.Hash, b.Height,
		)
		for i := range b.transactions {
			txHash := &b.transactions[i]
			recKey := keyTxRecord(txHash, &b.Block)
			recVal := existsRawTxRecord(ns, recKey)
			var rec TxRecord
			if err = readRawTxRecord(txHash, recVal, &rec); slog.Check(err) {
				return
			}
			if err = deleteTxRecord(ns, txHash, &b.Block); slog.Check(err) {
				return
			}
			// Handle coinbase transactions specially since they are
			// not moved to the unconfirmed store.  A coinbase cannot
			// contain any debits, but all credits should be removed
			// and the mined balance decremented.
			if blockchain.IsCoinBaseTx(&rec.MsgTx) {
				op := wire.OutPoint{Hash: rec.Hash}
				for i, output := range rec.MsgTx.TxOut {
					k, v := existsCredit(ns, &rec.Hash,
						uint32(i), &b.Block)
					if v == nil {
						continue
					}
					op.Index = uint32(i)
					coinBaseCredits = append(coinBaseCredits, op)
					unspentKey, credKey := existsUnspent(ns, &op)
					if credKey != nil {
						minedBalance -= util.Amount(output.Value)
						if err = deleteRawUnspent(ns, unspentKey); slog.Check(err) {
							return
						}
					}
					if err = deleteRawCredit(ns, k); slog.Check(err) {
						return
					}
				}
				continue
			}
			if err = putRawUnmined(ns, txHash[:], recVal); slog.Check(err) {
				return
			}
			// For each debit recorded for this transaction, mark
			// the credit it spends as unspent (as long as it still
			// exists) and delete the debit.  The previous output is
			// recorded in the unconfirmed store for every previous
			// output, not just debits.
			for i, input := range rec.MsgTx.TxIn {
				prevOut := &input.PreviousOutPoint
				prevOutKey := canonicalOutPoint(&prevOut.Hash,
					prevOut.Index)
				if err = putRawUnminedInput(ns, prevOutKey, rec.Hash[:]); slog.Check(err) {
					return
				}
				// If this input is a debit, remove the debit record and mark the credit that it spent as
				// unspent, incrementing the mined balance.
				var debKey, credKey []byte
				if debKey, credKey, err = existsDebit(ns, &rec.Hash, uint32(i), &b.Block); slog.Check(err) {
					return
				}
				if debKey == nil {
					continue
				}
				// unspendRawCredit does not error in case the no credit exists for this key, but this
				// behavior is correct.  Since blocks are removed in increasing order, this credit
				// may have already been removed from a previously removed transaction record in
				// this rollback.
				var amt util.Amount
				if amt, err = unspendRawCredit(ns, credKey); slog.Check(err) {
					return
				}
				if err = deleteRawDebit(ns, debKey); slog.Check(err) {
					return
				}
				// If the credit was previously removed in the rollback, the credit amount is zero.  Only
				// mark the previously spent credit as unspent if it still exists.
				if amt == 0 {
					continue
				}
				var unspentVal []byte
				if unspentVal, err = fetchRawCreditUnspentValue(credKey); slog.Check(err) {
					return
				}
				minedBalance += amt
				if err = putRawUnspent(ns, prevOutKey, unspentVal); slog.Check(err) {
					return
				}
			}
			// For each detached non-coinbase credit, move the
			// credit output to unmined.  If the credit is marked
			// unspent, it is removed from the utxo set and the
			// mined balance is decremented.
			//
			// TODO: use a credit iterator
			var amt util.Amount
			var change bool
			for i, output := range rec.MsgTx.TxOut {
				k, v := existsCredit(ns, &rec.Hash, uint32(i), &b.Block)
				if v == nil {
					continue
				}
				if amt, change, err = fetchRawCreditAmountChange(v); slog.Check(err) {
					return
				}
				outPointKey := canonicalOutPoint(&rec.Hash, uint32(i))
				unminedCredVal := valueUnminedCredit(amt, change)
				if err = putRawUnminedCredit(ns, outPointKey, unminedCredVal); slog.Check(err) {
					return
				}
				if err = deleteRawCredit(ns, k); slog.Check(err) {
					return
				}
				credKey := existsRawUnspent(ns, outPointKey)
				if credKey != nil {
					minedBalance -= util.Amount(output.Value)
					if err = deleteRawUnspent(ns, outPointKey); slog.Check(err) {
						return
					}
				}
			}
		}
		// reposition cursor before deleting this k/v pair and advancing to the
		// previous.
		it.reposition(it.elem.Height)
		// Avoid cursor deletion until bolt issue #620 is resolved.
		// err = it.delete()
		// if err != nil {
		// 	return err
		// }
	}
	if slog.Check(it.err) {
		return it.err
	}
	// Delete the block records outside of the iteration since cursor deletion
	// is broken.
	for _, h := range heightsToRemove {
		if err = deleteBlockRecord(ns, h); slog.Check(err) {
			return
		}
	}
	for _, op := range coinBaseCredits {
		opKey := canonicalOutPoint(&op.Hash, op.Index)
		unminedSpendTxHashKeys := fetchUnminedInputSpendTxHashes(ns, opKey)
		for _, unminedSpendTxHashKey := range unminedSpendTxHashKeys {
			unminedVal := existsRawUnmined(ns, unminedSpendTxHashKey[:])
			// If the spending transaction spends multiple outputs
			// from the same transaction, we'll find duplicate
			// entries within the store, so it's possible we're
			// unable to find it if the conflicts have already been
			// removed in a previous iteration.
			if unminedVal == nil {
				continue
			}
			var unminedRec TxRecord
			unminedRec.Hash = unminedSpendTxHashKey
			if err = readRawTxRecord(&unminedRec.Hash, unminedVal, &unminedRec); slog.Check(err) {
				return
			}
			slog.Debugf("transaction %v spends a removed coinbase output -- removing as well %s", unminedRec.Hash)
			if err = RemoveConflict(ns, &unminedRec); slog.Check(err) {
				return
			}
		}
	}
	return putMinedBalance(ns, minedBalance)
}

// UnspentOutputs returns all unspent received transaction outputs.
// The order is undefined.
func (s *Store) UnspentOutputs(ns walletdb.ReadBucket) (unspent []Credit, err error) {
	var op wire.OutPoint
	var block Block
	if err = ns.NestedReadBucket(bucketUnspent).ForEach(func(k, v []byte) (err error) {
		if err = readCanonicalOutPoint(k, &op); slog.Check(err) {
			return
		}
		if existsRawUnminedInput(ns, k) != nil {
			// Output is spent by an unmined transaction.
			// Skip this k/v pair.
			return
		}
		if err = readUnspentBlock(v, &block); slog.Check(err) {
			return
		}
		var blockTime time.Time
		if blockTime, err = fetchBlockTime(ns, block.Height); slog.Check(err) {
			return
		}
		// TODO(jrick): reading the entire transaction should be avoidable. Creating the credit only requires the
		//  output amount and pkScript.
		var rec *TxRecord
		if rec, err = fetchTxRecord(ns, &op.Hash, &block); slog.Check(err) {
			return
		}
		txOut := rec.MsgTx.TxOut[op.Index]
		cred := Credit{
			OutPoint: op,
			BlockMeta: BlockMeta{
				Block: block,
				Time:  blockTime,
			},
			Amount:       util.Amount(txOut.Value),
			PkScript:     txOut.PkScript,
			Received:     rec.Received,
			FromCoinBase: blockchain.IsCoinBaseTx(&rec.MsgTx),
		}
		unspent = append(unspent, cred)
		return
	}); slog.Check(err) {
		if _, ok := err.(TxMgrError); ok {
			return
		}
		str := "failed iterating unspent bucket"
		err = storeError(ErrDatabase, str, err)
		slog.Debug(err)
		return
	}
	if err = ns.NestedReadBucket(bucketUnminedCredits).ForEach(func(k, v []byte) (err error) {
		if existsRawUnminedInput(ns, k) != nil {
			// Output is spent by an unmined transaction.
			// Skip to next unmined credit.
			return
		}
		if err = readCanonicalOutPoint(k, &op); slog.Check(err) {
			return
		}
		// TODO(jrick): Reading/parsing the entire transaction record
		// just for the output amount and script can be avoided.
		recVal := existsRawUnmined(ns, op.Hash[:])
		var rec TxRecord
		if err = readRawTxRecord(&op.Hash, recVal, &rec); slog.Check(err) {
			return
		}
		txOut := rec.MsgTx.TxOut[op.Index]
		cred := Credit{
			OutPoint: op,
			BlockMeta: BlockMeta{
				Block: Block{Height: -1},
			},
			Amount:       util.Amount(txOut.Value),
			PkScript:     txOut.PkScript,
			Received:     rec.Received,
			FromCoinBase: blockchain.IsCoinBaseTx(&rec.MsgTx),
		}
		unspent = append(unspent, cred)
		return
	}); slog.Check(err) {
		if _, ok := err.(TxMgrError); ok {
			return
		}
		str := "failed iterating unmined credits bucket"
		err = storeError(ErrDatabase, str, err)
		slog.Debug(err)
		return
	}
	return
}

// Balance returns the spendable wallet balance (total value of all unspent
// transaction outputs) given a minimum of minConf confirmations, calculated
// at a current chain height of curHeight.  Coinbase outputs are only included
// in the balance if maturity has been reached.
//
// Balance may return unexpected results if syncHeight is lower than the block
// height of the most recent mined transaction in the store.
func (s *Store) Balance(ns walletdb.ReadBucket, minConf int32, syncHeight int32) (bal util.Amount, err error) {
	if bal, err = fetchMinedBalance(ns); slog.Check(err) {
		return
	}
	// Subtract the balance for each credit that is spent by an unmined
	// transaction.
	var op wire.OutPoint
	var block Block
	if err = ns.NestedReadBucket(bucketUnspent).ForEach(func(k, v []byte) (err error) {
		if err = readCanonicalOutPoint(k, &op); slog.Check(err) {
			return
		}
		if err = readUnspentBlock(v, &block); slog.Check(err) {
			return
		}
		if existsRawUnminedInput(ns, k) != nil {
			_, v := existsCredit(ns, &op.Hash, op.Index, &block)
			var amt util.Amount
			if amt, err = fetchRawCreditAmount(v); slog.Check(err) {
				return
			}
			bal -= amt
		}
		return
	}); slog.Check(err) {
		if _, ok := err.(TxMgrError); ok {
			return
		}
		str := "failed iterating unspent outputs"
		err = storeError(ErrDatabase, str, err)
		slog.Debug(err)
		return
	}
	// Decrement the balance for any unspent credit with less than
	// minConf confirmations and any (unspent) immature coinbase credit.
	coinbaseMaturity := int32(s.chainParams.CoinbaseMaturity)
	stopConf := minConf
	if coinbaseMaturity > stopConf {
		stopConf = coinbaseMaturity
	}
	lastHeight := syncHeight - stopConf
	blockIt := makeReadReverseBlockIterator(ns)
	for blockIt.prev() {
		block := &blockIt.elem
		if block.Height < lastHeight {
			break
		}
		for i := range block.transactions {
			txHash := &block.transactions[i]
			var rec *TxRecord
			if rec, err = fetchTxRecord(ns, txHash, &block.Block); slog.Check(err) {
				return
			}
			numOuts := uint32(len(rec.MsgTx.TxOut))
			for i := uint32(0); i < numOuts; i++ {
				// Avoid double decrementing the credit amount
				// if it was already removed for being spent by
				// an unmined tx.
				opKey := canonicalOutPoint(txHash, i)
				if existsRawUnminedInput(ns, opKey) != nil {
					continue
				}
				_, v := existsCredit(ns, txHash, i, &block.Block)
				if v == nil {
					continue
				}
				var amt util.Amount
				var spent bool
				if amt, spent, err = fetchRawCreditAmountSpent(v); slog.Check(err) {
					return
				}
				if spent {
					continue
				}
				confs := syncHeight - block.Height + 1
				if confs < minConf || (blockchain.IsCoinBaseTx(&rec.MsgTx) &&
					confs < coinbaseMaturity) {
					bal -= amt
				}
			}
		}
	}
	if slog.Check(blockIt.err) {
		return 0, blockIt.err
	}
	// If unmined outputs are included, increment the balance for each
	// output that is unspent.
	if minConf == 0 {
		if err = ns.NestedReadBucket(bucketUnminedCredits).ForEach(func(k, v []byte) (err error) {
			if existsRawUnminedInput(ns, k) != nil {
				// Output is spent by an unmined transaction.
				// Skip to next unmined credit.
				return
			}
			var amount util.Amount
			if amount, err = fetchRawUnminedCreditAmount(v); slog.Check(err) {
				return
			}
			bal += amount
			return
		}); slog.Check(err) {
			if _, ok := err.(TxMgrError); ok {
				return
			}
			str := "failed to iterate over unmined credits bucket"
			err = storeError(ErrDatabase, str, err)
			slog.Debug(err)
			return
		}
	}
	return
}
