// Package wallet Copyright (c) 2015-2016 The btcsuite developers
package wallet

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chainclient"
	"sort"

	ec "github.com/p9c/matrjoska/pkg/ecc"
	"github.com/p9c/matrjoska/pkg/txauthor"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pkg/wire"
	"github.com/p9c/matrjoska/pkg/wtxmgr"
)

// byAmount defines the methods needed to satisify sort.Interface to txsort credits by their output amount.
type byAmount []wtxmgr.Credit

func (s byAmount) Len() int           { return len(s) }
func (s byAmount) Less(i, j int) bool { return s[i].Amount < s[j].Amount }
func (s byAmount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func makeInputSource(eligible []wtxmgr.Credit) txauthor.InputSource {
	// Pick largest outputs first. This is only done for compatibility with previous
	// tx creation code, not because it's a good idea.
	sort.Sort(sort.Reverse(byAmount(eligible)))
	// Current inputs and their total value. These are closed over by the returned
	// input source and reused across multiple calls.
	currentTotal := amt.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, len(eligible))
	currentScripts := make([][]byte, 0, len(eligible))
	currentInputValues := make([]amt.Amount, 0, len(eligible))
	return func(target amt.Amount) (
		amt.Amount, []*wire.TxIn,
		[]amt.Amount, [][]byte, error,
	) {
		for currentTotal < target && len(eligible) != 0 {
			nextCredit := &eligible[0]
			eligible = eligible[1:]
			nextInput := wire.NewTxIn(&nextCredit.OutPoint, nil, nil)
			currentTotal += nextCredit.Amount
			currentInputs = append(currentInputs, nextInput)
			currentScripts = append(currentScripts, nextCredit.PkScript)
			currentInputValues = append(currentInputValues, nextCredit.Amount)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
}

// secretSource is an implementation of txauthor.SecretSource for the wallet's address manager.
type secretSource struct {
	*waddrmgr.Manager
	addrmgrNs walletdb.ReadBucket
}

// GetKey gets the private key for an address if it is available
func (s secretSource) GetKey(addr btcaddr.Address) (privKey *ec.PrivateKey, cmpr bool, e error) {
	var ma waddrmgr.ManagedAddress
	ma, e = s.Address(s.addrmgrNs, addr)
	if e != nil {
		return nil, false, e
	}
	var ok bool
	var mpka waddrmgr.ManagedPubKeyAddress
	mpka, ok = ma.(waddrmgr.ManagedPubKeyAddress)
	if !ok {
		e = fmt.Errorf(
			"managed address type for %v is `%T` but "+
				"want waddrmgr.ManagedPubKeyAddress", addr, ma,
		)
		return nil, false, e
	}
	if privKey, e = mpka.PrivKey(); E.Chk(e) {
		return nil, false, e
	}
	return privKey, ma.Compressed(), nil
}

// GetScript returns pay to script transaction
func (s secretSource) GetScript(addr btcaddr.Address) ([]byte, error) {
	ma, e := s.Address(s.addrmgrNs, addr)
	if e != nil {
		return nil, e
	}
	msa, ok := ma.(waddrmgr.ManagedScriptAddress)
	if !ok {
		e := fmt.Errorf(
			"managed address type for %v is `%T` but "+
				"want waddrmgr.ManagedScriptAddress", addr, ma,
		)
		return nil, e
	}
	return msa.Script()
}

// txToOutputs creates a signed transaction which includes each output from
// outputs. Previous outputs to reedeem are chosen from the passed account's
// UTXO set and minconf policy. An additional output may be added to return
// change to the wallet. An appropriate fee is included based on the wallet's
// current relay fee. The wallet must be unlocked to create the transaction.
func (w *Wallet) txToOutputs(
	outputs []*wire.TxOut, account uint32,
	minconf int32, feeSatPerKb amt.Amount,
) (tx *txauthor.AuthoredTx, e error) {
	var chainClient chainclient.Interface
	if chainClient, e = w.requireChainClient(); E.Chk(e) {
		return nil, e
	}
	e = walletdb.Update(
		w.db, func(dbtx walletdb.ReadWriteTx) (e error) {
			addrmgrNs := dbtx.ReadWriteBucket(waddrmgrNamespaceKey)
			// Get current block's height and hash.
			var bs *waddrmgr.BlockStamp
			if bs, e = chainClient.BlockStamp(); E.Chk(e) {
				return
			}
			var eligible []wtxmgr.Credit
			if eligible, e = w.findEligibleOutputs(dbtx, account, minconf, bs); E.Chk(e) {
				return
			}
			inputSource := makeInputSource(eligible)
			changeSource := func() (b []byte, e error) {
				// Derive the change output script. As a hack to allow spending from the
				// imported account, change addresses are created from account 0.
				var changeAddr btcaddr.Address
				if account == waddrmgr.ImportedAddrAccount {
					changeAddr, e = w.newChangeAddress(addrmgrNs, 0)
				} else {
					changeAddr, e = w.newChangeAddress(addrmgrNs, account)
				}
				if E.Chk(e) {
					return
				}
				return txscript.PayToAddrScript(changeAddr)
			}
			if tx, e = txauthor.NewUnsignedTransaction(outputs, feeSatPerKb, inputSource, changeSource); E.Chk(e) {
				return
			}
			// Randomize change position, if change exists, before signing. This doesn't
			// affect the serialize size, so the change amount will still be valid.
			if tx.ChangeIndex >= 0 {
				tx.RandomizeChangePosition()
			}
			return tx.AddAllInputScripts(secretSource{w.Manager, addrmgrNs})
		},
	)
	if E.Chk(e) {
		return
	}
	if e = validateMsgTx(tx.Tx, tx.PrevScripts, tx.PrevInputValues); E.Chk(e) {
		return
	}
	if tx.ChangeIndex >= 0 && account == waddrmgr.ImportedAddrAccount {
		changeAmount := amt.Amount(tx.Tx.TxOut[tx.ChangeIndex].Value)
		W.F(
			"spend from imported account produced change: moving %v from imported account into default account.",
			changeAmount,
		)
	}
	return
}
func (w *Wallet) findEligibleOutputs(
	dbtx walletdb.ReadTx,
	account uint32,
	minconf int32,
	bs *waddrmgr.BlockStamp,
) ([]wtxmgr.Credit, error) {
	addrmgrNs := dbtx.ReadBucket(waddrmgrNamespaceKey)
	txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)
	unspent, e := w.TxStore.UnspentOutputs(txmgrNs)
	if e != nil {
		return nil, e
	}
	// TODO: Eventually all of these filters (except perhaps output locking) should
	//  be handled by the call to UnspentOutputs (or similar). Because one of these
	//  filters requires matching the output script to the desired
	//  account, this change depends on making wtxmgr a waddrmgr dependancy and
	//  requesting unspent outputs for a single account.
	eligible := make([]wtxmgr.Credit, 0, len(unspent))
	for i := range unspent {
		output := &unspent[i]
		// Only include this output if it meets the required number of confirmations.
		// Coinbase transactions must have have reached maturity before their outputs
		// may be spent.
		if !confirmed(minconf, output.Height, bs.Height) {
			continue
		}
		if output.FromCoinBase {
			target := int32(w.chainParams.CoinbaseMaturity)
			if !confirmed(target, output.Height, bs.Height) {
				continue
			}
		}
		// Locked unspent outputs are skipped.
		if w.LockedOutpoint(output.OutPoint) {
			continue
		}
		// Only include the output if it is associated with the passed account.
		//
		// TODO: Handle multisig outputs by determining if enough of the addresses are
		//  controlled.
		var addrs []btcaddr.Address
		if _, addrs, _, e = txscript.ExtractPkScriptAddrs(
			output.PkScript, w.chainParams,
		); E.Chk(e) || len(addrs) != 1 {
			continue
		}
		var addrAcct uint32
		if _, addrAcct, e = w.Manager.AddrAccount(addrmgrNs, addrs[0]); E.Chk(e) ||
			addrAcct != account {
			continue
		}
		eligible = append(eligible, *output)
	}
	return eligible, nil
}

// validateMsgTx verifies transaction input scripts for tx. All previous output
// scripts from outputs redeemed by the transaction, in the same order they are
// spent, must be passed in the prevScripts slice.
func validateMsgTx(tx *wire.MsgTx, prevScripts [][]byte, inputValues []amt.Amount) (e error) {
	hashCache := txscript.NewTxSigHashes(tx)
	for i, prevScript := range prevScripts {
		var vm *txscript.Engine
		vm, e = txscript.NewEngine(
			prevScript, tx, i,
			txscript.StandardVerifyFlags, nil, hashCache, int64(inputValues[i]),
		)
		if E.Chk(e) {
			return fmt.Errorf("cannot create script engine: %s", e)
		}
		e = vm.Execute()
		if E.Chk(e) {
			return fmt.Errorf("cannot validate transaction: %s", e)
		}
	}
	return nil
}
