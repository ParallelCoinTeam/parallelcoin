package wallet

import (
	"errors"
	
	"github.com/p9c/pod/pkg/database/walletdb"
	"github.com/p9c/pod/pkg/txscript"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/wallet/waddrmgr"
)

// MakeMultiSigScript creates a multi-signature script that can be redeemed with nRequired signatures of the passed keys
// and addresses. If the address is a P2PKH address, the associated pubkey is looked up by the wallet if possible,
// otherwise an error is returned for a missing pubkey.
//
// This function only works with pubkeys and P2PKH addresses derived from them.
func (w *Wallet) MakeMultiSigScript(addrs []util.Address, nRequired int) ([]byte, error) {
	pubKeys := make([]*util.AddressPubKey, len(addrs))
	var dbtx walletdb.ReadTx
	var addrmgrNs walletdb.ReadBucket
	defer func() {
		if dbtx != nil {
			e := dbtx.Rollback()
			if e != nil {
			}
		}
	}()
	// The address list will made up either of addreseses (pubkey hash), for which we need to look up the keys in
	// wallet, straight pubkeys, or a mixture of the two.
	for i, addr := range addrs {
		switch addr := addr.(type) {
		default:
			return nil, errors.New(
				"cannot make multisig script for " +
					"a non-secp256k1 public key or P2PKH address",
			)
		case *util.AddressPubKey:
			pubKeys[i] = addr
		case *util.AddressPubKeyHash:
			if dbtx == nil {
				var e error
				dbtx, e = w.db.BeginReadTx()
				if e != nil {
					return nil, e
				}
				addrmgrNs = dbtx.ReadBucket(waddrmgrNamespaceKey)
			}
			addrInfo, e := w.Manager.Address(addrmgrNs, addr)
			if e != nil {
				return nil, e
			}
			serializedPubKey := addrInfo.(waddrmgr.ManagedPubKeyAddress).
				PubKey().SerializeCompressed()
			pubKeyAddr, e := util.NewAddressPubKey(
				serializedPubKey, w.chainParams,
			)
			if e != nil {
				return nil, e
			}
			pubKeys[i] = pubKeyAddr
		}
	}
	return txscript.MultiSigScript(pubKeys, nRequired)
}

// ImportP2SHRedeemScript adds a P2SH redeem script to the wallet.
func (w *Wallet) ImportP2SHRedeemScript(script []byte) (*util.AddressScriptHash, error) {
	var p2shAddr *util.AddressScriptHash
	e := walletdb.Update(
		w.db, func(tx walletdb.ReadWriteTx) (e error) {
			addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			// TODO(oga) blockstamp current block?
			bs := &waddrmgr.BlockStamp{
				Hash:   *w.ChainParams().GenesisHash,
				Height: 0,
			}
			// As this is a regular P2SH script, we'll import this into the BIP0044 scope.
			var bip44Mgr *waddrmgr.ScopedKeyManager
			bip44Mgr, e = w.Manager.FetchScopedKeyManager(
				waddrmgr.KeyScopeBIP0044,
			)
			if e != nil {
				return e
			}
			addrInfo, e := bip44Mgr.ImportScript(addrmgrNs, script, bs)
			if e != nil {
				// Don't care if it's already there, but still have to set the p2shAddr since the address manager didn't
				// return anything useful.
				if waddrmgr.IsError(e, waddrmgr.ErrDuplicateAddress) {
					// This function will never error as it always hashes the script to the correct length.
					p2shAddr, _ = util.NewAddressScriptHash(
						script,
						w.chainParams,
					)
					return nil
				}
				return e
			}
			p2shAddr = addrInfo.Address().(*util.AddressScriptHash)
			return nil
		},
	)
	return p2shAddr, e
}
