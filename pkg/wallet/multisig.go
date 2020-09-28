package wallet

import (
	"errors"

	"github.com/stalker-loki/app/slog"

	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/db/walletdb"
	"github.com/p9c/pod/pkg/util"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
)

// MakeMultiSigScript creates a multi-signature script that can be redeemed with
// nRequired signatures of the passed keys and addresses.  If the address is a
// P2PKH address, the associated pubkey is looked up by the wallet if possible,
// otherwise an error is returned for a missing pubkey.
//
// This function only works with pubkeys and P2PKH addresses derived from them.
func (w *Wallet) MakeMultiSigScript(addrs []util.Address, nRequired int) (scr []byte, err error) {
	pubKeys := make([]*util.AddressPubKey, len(addrs))
	var dbtx walletdb.ReadTx
	var addrmgrNs walletdb.ReadBucket
	defer func() {
		if dbtx != nil {
			if err = dbtx.Rollback(); slog.Check(err) {
			}
		}
	}()
	// The address list will made up either of addresses (pubkey hash), for which we need to look up the keys in wallet,
	// straight pubkeys, or a mixture of the two.
	for i, addr := range addrs {
		switch addr := addr.(type) {
		default:
			return nil, errors.New("cannot make multisig script for a non-secp256k1 public key or P2PKH address")
		case *util.AddressPubKey:
			pubKeys[i] = addr
		case *util.AddressPubKeyHash:
			if dbtx == nil {
				if dbtx, err = w.db.BeginReadTx(); slog.Check(err) {
					return
				}
				addrmgrNs = dbtx.ReadBucket(waddrmgrNamespaceKey)
			}
			var addrInfo waddrmgr.ManagedAddress
			if addrInfo, err = w.Manager.Address(addrmgrNs, addr); slog.Check(err) {
				return
			}
			serializedPubKey := addrInfo.(waddrmgr.ManagedPubKeyAddress).PubKey().SerializeCompressed()
			var pubKeyAddr *util.AddressPubKey
			if pubKeyAddr, err = util.NewAddressPubKey(serializedPubKey, w.chainParams); slog.Check(err) {
				return
			}
			pubKeys[i] = pubKeyAddr
		}
	}
	return txscript.MultiSigScript(pubKeys, nRequired)
}

// ImportP2SHRedeemScript adds a P2SH redeem script to the wallet.
func (w *Wallet) ImportP2SHRedeemScript(script []byte) (p2shAddr *util.AddressScriptHash, err error) {
	if err = walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) (err error) {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		// TODO(oga) blockstamp current block?
		bs := &waddrmgr.BlockStamp{
			Hash:   *w.ChainParams().GenesisHash,
			Height: 0,
		}
		// As this is a regular P2SH script, we'll import this into the BIP0044 scope.
		var bip44Mgr *waddrmgr.ScopedKeyManager
		if bip44Mgr, err = w.Manager.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0084); slog.Check(err) {
			return
		}
		var addrInfo waddrmgr.ManagedScriptAddress
		if addrInfo, err = bip44Mgr.ImportScript(addrmgrNs, script, bs); slog.Check(err) {
			// Don't care if it's already there, but still have to set the p2shAddr since the address manager didn't
			// return anything useful.
			if waddrmgr.IsError(err, waddrmgr.ErrDuplicateAddress) {
				// This function will never error as it always hashes the script to the correct length.
				p2shAddr, _ = util.NewAddressScriptHash(script, w.chainParams)
				return
			}
			return
		}
		p2shAddr = addrInfo.Address().(*util.AddressScriptHash)
		return
	}); slog.Check(err) {
	}
	return
}
