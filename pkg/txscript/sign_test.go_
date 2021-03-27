package txscript

import (
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/btcaddr"
	"github.com/p9c/pod/pkg/chaincfg"
	"testing"
	
	"github.com/p9c/pod/pkg/chainhash"
	ec "github.com/p9c/pod/pkg/ecc"
	"github.com/p9c/pod/pkg/wire"
)

type addressToKey struct {
	key        *ec.PrivateKey
	compressed bool
}

func mkGetKey(keys map[string]addressToKey) KeyDB {
	if keys == nil {
		return KeyClosure(
			func(addr btcaddr.Address) (
				*ec.PrivateKey,
				bool, error,
			) {
				return nil, false, errors.New("nope")
			},
		)
	}
	return KeyClosure(
		func(addr btcaddr.Address) (
			*ec.PrivateKey,
			bool, error,
		) {
			a2k, ok := keys[addr.EncodeAddress()]
			if !ok {
				return nil, false, errors.New("nope")
			}
			return a2k.key, a2k.compressed, nil
		},
	)
}
func mkGetScript(scripts map[string][]byte) ScriptDB {
	if scripts == nil {
		return ScriptClosure(
			func(addr btcaddr.Address) ([]byte, error) {
				return nil, errors.New("nope")
			},
		)
	}
	return ScriptClosure(
		func(addr btcaddr.Address) ([]byte, error) {
			script, ok := scripts[addr.EncodeAddress()]
			if !ok {
				return nil, errors.New("nope")
			}
			return script, nil
		},
	)
}
func checkScripts(msg string, tx *wire.MsgTx, idx int, inputAmt int64, sigScript, pkScript []byte) (e error) {
	tx.TxIn[idx].SignatureScript = sigScript
	vm, e := NewEngine(
		pkScript, tx, idx,
		ScriptBip16|ScriptVerifyDERSignatures, nil, nil, inputAmt,
	)
	if e != nil {
		return fmt.Errorf(
			"failed to make script engine for %s: %v",
			msg, e,
		)
	}
	e = vm.Execute()
	if e != nil {
		return fmt.Errorf(
			"invalid script signature for %s: %v", msg,
			e,
		)
	}
	return nil
}
func signAndCheck(
	msg string, tx *wire.MsgTx, idx int, inputAmt int64, pkScript []byte,
	hashType SigHashType, kdb KeyDB, sdb ScriptDB,
	previousScript []byte,
) (e error) {
	sigScript, e := SignTxOutput(
		&chaincfg.TestNet3Params, tx, idx,
		pkScript, hashType, kdb, sdb, nil,
	)
	if e != nil {
		return fmt.Errorf("failed to sign output %s: %v", msg, e)
	}
	return checkScripts(msg, tx, idx, inputAmt, sigScript, pkScript)
}
func TestSignTxOutput(t *testing.T) {
	t.Parallel()
	// make key
	// make script based on key.
	// sign with magic pixie dust.
	hashTypes := []SigHashType{
		SigHashOld, // no longer used but should act like all
		SigHashAll,
		SigHashNone,
		SigHashSingle,
		SigHashAll | SigHashAnyOneCanPay,
		SigHashNone | SigHashAnyOneCanPay,
		SigHashSingle | SigHashAnyOneCanPay,
	}
	inputAmounts := []int64{5, 10, 15}
	tx := &wire.MsgTx{
		Version: 1,
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 0,
				},
				Sequence: 4294967295,
			},
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 1,
				},
				Sequence: 4294967295,
			},
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 2,
				},
				Sequence: 4294967295,
			},
		},
		TxOut: []*wire.TxOut{
			{
				Value: 1,
			},
			{
				Value: 2,
			},
			{
				Value: 3,
			},
		},
		LockTime: 0,
	}
	// Pay to Pubkey Hash (uncompressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i], pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to Pubkey Hash (uncompressed) (merging with correct)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), sigScript,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(msg, tx, i, inputAmounts[i], sigScript, pkScript)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to Pubkey Hash (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to Pubkey Hash (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), sigScript,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, pkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(nil), sigScript,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(msg, tx, i, inputAmounts[i], sigScript, pkScript)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to PubKey (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, pkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(nil), sigScript,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, pkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// As before, but with p2sh now. Pay to Pubkey Hash (uncompressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			var ad *btcaddr.PubKeyHash
			ad, e = btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(ad)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
				break
			}
			scriptAddr, e := NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(
				scriptAddr,
			)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to Pubkey Hash (uncompressed) with duplicate merge
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
				break
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(
				scriptAddr,
			)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to Pubkey Hash (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(
				scriptAddr,
			)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to Pubkey Hash (compressed) with duplicate merge
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKeyHash(
				btcaddr.Hash160(pk), &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(
				scriptAddr,
			)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to PubKey (uncompressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(
				scriptAddr,
			)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to PubKey (uncompressed) with duplicate merge
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeUncompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, false},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Pay to PubKey (compressed)
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk := (*ec.PublicKey)(&key.PublicKey).
				SerializeCompressed()
			address, e := btcaddr.NewPubKey(
				pk,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := PayToAddrScript(address)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := address.NewAddressScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// by the above loop, this should be valid, now sign again and merge.
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address.EncodeAddress(): {key, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s a "+
						"second time: %v", msg, e,
				)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"twice signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Basic Multisig
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key1, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk1 := (*ec.PublicKey)(&key1.PublicKey).
				SerializeCompressed()
			address1, e := btcaddr.NewPubKey(
				pk1,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			key2, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey 2 for %s: %v",
					msg, e,
				)
				break
			}
			pk2 := (*ec.PublicKey)(&key2.PublicKey).
				SerializeCompressed()
			address2, e := btcaddr.NewPubKey(
				pk2,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address 2 for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := MultiSigScript(
				[]*btcaddr.PubKey{address1, address2},
				2,
			)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := btcaddr.NewScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			if e := signAndCheck(
				msg, tx, i, inputAmounts[i],
				scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address1.EncodeAddress(): {key1, true},
						address2.EncodeAddress(): {key2, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			); E.Chk(e) {
				break
			}
		}
	}
	// Two part multisig, sign with one key then the other.
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key1, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk1 := (*ec.PublicKey)(&key1.PublicKey).
				SerializeCompressed()
			address1, e := btcaddr.NewPubKey(
				pk1,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			key2, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey 2 for %s: %v",
					msg, e,
				)
				break
			}
			pk2 := (*ec.PublicKey)(&key2.PublicKey).
				SerializeCompressed()
			address2, e := btcaddr.NewPubKey(
				pk2,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address 2 for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := MultiSigScript(
				[]*btcaddr.PubKey{address1, address2},
				2,
			)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := btcaddr.NewScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address1.EncodeAddress(): {key1, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// Only 1 out of 2 signed, this *should* fail.
			if checkScripts(
				msg, tx, i, inputAmounts[i], sigScript,
				scriptPkScript,
			) == nil {
				t.Errorf("part signed script valid for %s", msg)
				break
			}
			// Sign with the other key and merge
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address2.EncodeAddress(): {key2, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), sigScript,
			)
			if e != nil {
				t.Errorf("failed to sign output %s: %v", msg, e)
				break
			}
			e = checkScripts(
				msg, tx, i, inputAmounts[i], sigScript,
				scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"fully signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
	// Two part multisig, sign with one key then both, check key dedup correctly.
	for _, hashType := range hashTypes {
		for i := range tx.TxIn {
			msg := fmt.Sprintf("%d:%d", hashType, i)
			key1, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey for %s: %v",
					msg, e,
				)
				break
			}
			pk1 := (*ec.PublicKey)(&key1.PublicKey).
				SerializeCompressed()
			address1, e := btcaddr.NewPubKey(
				pk1,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address for %s: %v",
					msg, e,
				)
				break
			}
			key2, e := ec.NewPrivateKey(ec.S256())
			if e != nil {
				t.Errorf(
					"failed to make privKey 2 for %s: %v",
					msg, e,
				)
				break
			}
			pk2 := (*ec.PublicKey)(&key2.PublicKey).
				SerializeCompressed()
			address2, e := btcaddr.NewPubKey(
				pk2,
				&chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make address 2 for %s: %v",
					msg, e,
				)
				break
			}
			pkScript, e := MultiSigScript(
				[]*btcaddr.PubKey{address1, address2},
				2,
			)
			if e != nil {
				t.Errorf(
					"failed to make pkscript "+
						"for %s: %v", msg, e,
				)
			}
			scriptAddr, e := btcaddr.NewScriptHash(
				pkScript, &chaincfg.TestNet3Params,
			)
			if e != nil {
				t.Errorf(
					"failed to make p2sh addr for %s: %v",
					msg, e,
				)
				break
			}
			scriptPkScript, e := PayToAddrScript(scriptAddr)
			if e != nil {
				t.Errorf(
					"failed to make script pkscript for "+
						"%s: %v", msg, e,
				)
				break
			}
			sigScript, e := SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address1.EncodeAddress(): {key1, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), nil,
			)
			if e != nil {
				t.Errorf(
					"failed to sign output %s: %v", msg,
					e,
				)
				break
			}
			// Only 1 out of 2 signed, this *should* fail.
			if checkScripts(
				msg, tx, i, inputAmounts[i], sigScript,
				scriptPkScript,
			) == nil {
				t.Errorf("part signed script valid for %s", msg)
				break
			}
			// Sign with the other key and merge
			sigScript, e = SignTxOutput(
				&chaincfg.TestNet3Params,
				tx, i, scriptPkScript, hashType,
				mkGetKey(
					map[string]addressToKey{
						address1.EncodeAddress(): {key1, true},
						address2.EncodeAddress(): {key2, true},
					},
				), mkGetScript(
					map[string][]byte{
						scriptAddr.EncodeAddress(): pkScript,
					},
				), sigScript,
			)
			if e != nil {
				t.Errorf("failed to sign output %s: %v", msg, e)
				break
			}
			// Now we should pass.
			e = checkScripts(
				msg, tx, i, inputAmounts[i],
				sigScript, scriptPkScript,
			)
			if e != nil {
				t.Errorf(
					"fully signed script invalid for "+
						"%s: %v", msg, e,
				)
				break
			}
		}
	}
}

type tstInput struct {
	txout              *wire.TxOut
	sigscriptGenerates bool
	inputValidates     bool
	indexOutOfRange    bool
}
type tstSigScript struct {
	name               string
	inputs             []tstInput
	hashType           SigHashType
	compress           bool
	scriptAtWrongIndex bool
}

var coinbaseOutPoint = &wire.OutPoint{
	Index: (1 << 32) - 1,
}

// Pregenerated private key, with associated public key and pkScripts for the uncompressed and compressed hash160.
var (
	privKeyD = []byte{
		0x6b, 0x0f, 0xd8, 0xda, 0x54, 0x22, 0xd0, 0xb7,
		0xb4, 0xfc, 0x4e, 0x55, 0xd4, 0x88, 0x42, 0xb3, 0xa1, 0x65,
		0xac, 0x70, 0x7f, 0x3d, 0xa4, 0x39, 0x5e, 0xcb, 0x3b, 0xb0,
		0xd6, 0x0e, 0x06, 0x92,
	}
	// pubkeyX = []byte{0xb2, 0x52, 0xf0, 0x49, 0x85, 0x78, 0x03, 0x03, 0xc8,
	// 	0x7d, 0xce, 0x51, 0x7f, 0xa8, 0x69, 0x0b, 0x91, 0x95, 0xf4,
	// 	0xf3, 0x5c, 0x26, 0x73, 0x05, 0x05, 0xa2, 0xee, 0xbc, 0x09,
	// 	0x38, 0x34, 0x3a}
	// pubkeyY = []byte{0xb7, 0xc6, 0x7d, 0xb2, 0xe1, 0xff, 0xc8, 0x43, 0x1f,
	// 	0x63, 0x32, 0x62, 0xaa, 0x60, 0xc6, 0x83, 0x30, 0xbd, 0x24,
	// 	0x7e, 0xef, 0xdb, 0x6f, 0x2e, 0x8d, 0x56, 0xf0, 0x3c, 0x9f,
	// 	0x6d, 0xb6, 0xf8}
	uncompressedPkScript = []byte{
		0x76, 0xa9, 0x14, 0xd1, 0x7c, 0xb5,
		0xeb, 0xa4, 0x02, 0xcb, 0x68, 0xe0, 0x69, 0x56, 0xbf, 0x32,
		0x53, 0x90, 0x0e, 0x0a, 0x86, 0xc9, 0xfa, 0x88, 0xac,
	}
	compressedPkScript = []byte{
		0x76, 0xa9, 0x14, 0x27, 0x4d, 0x9f, 0x7f,
		0x61, 0x7e, 0x7c, 0x7a, 0x1c, 0x1f, 0xb2, 0x75, 0x79, 0x10,
		0x43, 0x65, 0x68, 0x27, 0x9d, 0x86, 0x88, 0xac,
	}
	shortPkScript = []byte{
		0x76, 0xa9, 0x14, 0xd1, 0x7c, 0xb5,
		0xeb, 0xa4, 0x02, 0xcb, 0x68, 0xe0, 0x69, 0x56, 0xbf, 0x32,
		0x53, 0x90, 0x0e, 0x0a, 0x88, 0xac,
	}
	// uncompressedAddrStr = "1L6fd93zGmtzkK6CsZFVVoCwzZV3MUtJ4F"
	// compressedAddrStr   = "14apLppt9zTq6cNw8SDfiJhk9PhkZrQtYZ"
)

// Pretend output amounts.
const coinbaseVal = 2500000000
const fee = 5000000

var sigScriptTests = []tstSigScript{
	{
		name: "one input uncompressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "two inputs uncompressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "one input compressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "two inputs compressed",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, compressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashNone",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashNone,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashSingle",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashSingle,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType SigHashAnyoneCanPay",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAnyOneCanPay,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "hashType non-standard",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           0x04,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "invalid compression",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     false,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           true,
		scriptAtWrongIndex: false,
	},
	{
		name: "short PkScript",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, shortPkScript),
				sigscriptGenerates: false,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           false,
		scriptAtWrongIndex: false,
	},
	{
		name: "valid script at wrong index",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           false,
		scriptAtWrongIndex: true,
	},
	{
		name: "index out of range",
		inputs: []tstInput{
			{
				txout:              wire.NewTxOut(coinbaseVal, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
			{
				txout:              wire.NewTxOut(coinbaseVal+fee, uncompressedPkScript),
				sigscriptGenerates: true,
				inputValidates:     true,
				indexOutOfRange:    false,
			},
		},
		hashType:           SigHashAll,
		compress:           false,
		scriptAtWrongIndex: true,
	},
}

// Test the sigscript generation for valid and invalid inputs, all hashTypes, and with and without compression.  This test creates sigscripts to spend fake coinbase inputs, as sigscripts cannot be created for the MsgTxs in txTests, since they come from the blockchain and we don't have the private keys.
func TestSignatureScript(t *testing.T) {
	t.Parallel()
	privKey, _ := ec.PrivKeyFromBytes(ec.S256(), privKeyD)
nexttest:
	for i := range sigScriptTests {
		tx := wire.NewMsgTx(wire.TxVersion)
		output := wire.NewTxOut(500, []byte{OP_RETURN})
		tx.AddTxOut(output)
		for range sigScriptTests[i].inputs {
			txin := wire.NewTxIn(coinbaseOutPoint, nil, nil)
			tx.AddTxIn(txin)
		}
		var script []byte
		var e error
		for j := range tx.TxIn {
			var idx int
			if sigScriptTests[i].inputs[j].indexOutOfRange {
				t.Errorf("at test %v", sigScriptTests[i].name)
				idx = len(sigScriptTests[i].inputs)
			} else {
				idx = j
			}
			script, e = SignatureScript(
				tx, idx,
				sigScriptTests[i].inputs[j].txout.PkScript,
				sigScriptTests[i].hashType, privKey,
				sigScriptTests[i].compress,
			)
			if (e == nil) != sigScriptTests[i].inputs[j].sigscriptGenerates {
				if e == nil {
					t.Errorf(
						"passed test '%v' incorrectly",
						sigScriptTests[i].name,
					)
				} else {
					t.Errorf(
						"failed test '%v': %v",
						sigScriptTests[i].name, e,
					)
				}
				continue nexttest
			}
			if !sigScriptTests[i].inputs[j].sigscriptGenerates {
				// done with this test
				continue nexttest
			}
			tx.TxIn[j].SignatureScript = script
		}
		// If testing using a correct sigscript but for an incorrect index, use last input script for first input.  Requires > 0 inputs for test.
		if sigScriptTests[i].scriptAtWrongIndex {
			tx.TxIn[0].SignatureScript = script
			sigScriptTests[i].inputs[0].inputValidates = false
		}
		// Validate tx input scripts
		scriptFlags := ScriptBip16 | ScriptVerifyDERSignatures
		for j := range tx.TxIn {
			vm, e := NewEngine(
				sigScriptTests[i].
					inputs[j].txout.PkScript, tx, j, scriptFlags, nil, nil, 0,
			)
			if e != nil {
				t.Errorf(
					"cannot create script vm for test %v: %v",
					sigScriptTests[i].name, e,
				)
				continue nexttest
			}
			e = vm.Execute()
			if (e == nil) != sigScriptTests[i].inputs[j].inputValidates {
				if e == nil {
					t.Errorf(
						"passed test '%v' validation incorrectly: %v",
						sigScriptTests[i].name, e,
					)
				} else {
					t.Errorf(
						"failed test '%v' validation: %v",
						sigScriptTests[i].name, e,
					)
				}
				continue nexttest
			}
		}
	}
}
