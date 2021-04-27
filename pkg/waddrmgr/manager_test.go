// Copyright (c) 2014-2016 The btcsuite developers
package waddrmgr_test

import (
	"encoding/hex"
	"fmt"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"os"
	"reflect"
	"testing"
	"time"
	
	"github.com/davecgh/go-spew/spew"
	
	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/snacl"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/walletdb"
)

// // newHash converts the passed big-endian hex string into a chainhash.Hash.
// // It only differs from the one available in wire in that it panics on an
// // error since it will only (and must only) be called with hard-coded, and
// // therefore known good, hashes.
// func newHash(// 	hexStr string) *chainhash.Hash {
// 	hash, e := chainhash.NewHashFromStr(hexStr)
// 	if e != nil  {
// 		panic(e)
// 	}
// 	return hash
// }

// failingSecretKeyGen is a waddrmgr.SecretKeyGenerator that always returns
// snacl.ErrDecryptFailed.
func failingSecretKeyGen(
	passphrase *[]byte,
	config *waddrmgr.ScryptOptions,
) (*snacl.SecretKey, error) {
	return nil, snacl.ErrDecryptFailed
}

// testContext is used to store context information about a running test which
// is passed into helper functions. The useSpends field indicates whether or not
// the spend data should be empty or figure it out based on the specific test
// blocks provided.
//
// This is needed because the first loop where the blocks are inserted, the
// tests are running against the latest block and therefore none of the outputs
// can be spent yet.
//
// However, on subsequent runs, all blocks have been inserted and therefore some
// of the transaction outputs are spent.
type testContext struct {
	t            *testing.T
	db           walletdb.DB
	rootManager  *waddrmgr.Manager
	manager      *waddrmgr.ScopedKeyManager
	account      uint32
	create       bool
	unlocked     bool
	watchingOnly bool
}

// addrType is the type of address being tested
type addrType byte

const (
	addrPubKeyHash addrType = iota
	addrScriptHash
)

// expectedAddr is used to house the expected return values from a managed
// address. Not all fields for used for all managed address types.
type expectedAddr struct {
	address     string
	addressHash []byte
	internal    bool
	compressed  bool
	// used           bool
	imported       bool
	pubKey         []byte
	privKey        []byte
	privKeyWIF     string
	script         []byte
	derivationInfo waddrmgr.DerivationPath
}

// testNamePrefix is a helper to return a prefix to show for test errors based
// on the state of the test context.
func testNamePrefix(tc *testContext) string {
	prefix := "Open "
	if tc.create {
		prefix = "Create "
	}
	return prefix + fmt.Sprintf("account #%d", tc.account)
}

// testManagedPubKeyAddress ensures the data returned by all exported functions
// provided by the passed managed public key address matches the corresponding
// fields in the provided expected address.
//
// When the test context indicates the manager is unlocked, the private data
// will also be tested, otherwise, the functions which deal with private data
// are checked to ensure they return the correct error.
func testManagedPubKeyAddress(
	tc *testContext, prefix string,
	gotAddr waddrmgr.ManagedPubKeyAddress, wantAddr *expectedAddr,
) bool {
	// Ensure pubkey is the expected value for the managed address.
	var gpubBytes []byte
	if gotAddr.Compressed() {
		gpubBytes = gotAddr.PubKey().SerializeCompressed()
	} else {
		gpubBytes = gotAddr.PubKey().SerializeUncompressed()
	}
	if !reflect.DeepEqual(gpubBytes, wantAddr.pubKey) {
		tc.t.Errorf(
			"%s PubKey: unexpected public key - got %x, want "+
				"%x", prefix, gpubBytes, wantAddr.pubKey,
		)
		return false
	}
	// Ensure exported pubkey string is the expected value for the managed address.
	gpubHex := gotAddr.ExportPubKey()
	wantPubHex := hex.EncodeToString(wantAddr.pubKey)
	if gpubHex != wantPubHex {
		tc.t.Errorf(
			"%s ExportPubKey: unexpected public key - got %s, "+
				"want %s", prefix, gpubHex, wantPubHex,
		)
		return false
	}
	// Ensure that the derivation path has been properly re-set after the address
	// was read from disk.
	_, gotAddrPath, ok := gotAddr.DerivationInfo()
	if !ok && !gotAddr.Imported() {
		tc.t.Errorf(
			"%s PubKey: non-imported address has empty "+
				"derivation info", prefix,
		)
		return false
	}
	expectedDerivationInfo := wantAddr.derivationInfo
	if gotAddrPath != expectedDerivationInfo {
		tc.t.Errorf(
			"%s PubKey: wrong derivation info: expected %v, "+
				"got %v", prefix, spew.Sdump(gotAddrPath),
			spew.Sdump(expectedDerivationInfo),
		)
		return false
	}
	// Ensure private key is the expected value for the managed address. Since this
	// is only available when the manager is unlocked, also check for the expected
	// error when the manager is locked.
	gotPrivKey, e := gotAddr.PrivKey()
	switch {
	case tc.watchingOnly:
		// Confirm expected watching-only error.
		testName := fmt.Sprintf("%s PrivKey", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrWatchingOnly) {
			return false
		}
	case tc.unlocked:
		if e != nil {
			tc.t.Errorf(
				"%s PrivKey: unexpected error - got %v",
				prefix, e,
			)
			return false
		}
		gpriv := gotPrivKey.Serialize()
		if !reflect.DeepEqual(gpriv, wantAddr.privKey) {
			tc.t.Errorf(
				"%s PrivKey: unexpected private key - "+
					"got %x, want %x", prefix, gpriv, wantAddr.privKey,
			)
			return false
		}
	default:
		// Confirm expected locked error.
		testName := fmt.Sprintf("%s PrivKey", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrLocked) {
			return false
		}
	}
	// Ensure exported private key in Wallet Import Format (WIF) is the expected
	// value for the managed address. Since this is only available when the manager
	// is unlocked, also check for the expected error when the manager is locked.
	gotWIF, e := gotAddr.ExportPrivKey()
	switch {
	case tc.watchingOnly:
		// Confirm expected watching-only error.
		testName := fmt.Sprintf("%s ExportPrivKey", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrWatchingOnly) {
			return false
		}
	case tc.unlocked:
		if e != nil {
			tc.t.Errorf(
				"%s ExportPrivKey: unexpected error - "+
					"got %v", prefix, e,
			)
			return false
		}
		if gotWIF.String() != wantAddr.privKeyWIF {
			tc.t.Errorf(
				"%s ExportPrivKey: unexpected WIF - got "+
					"%v, want %v", prefix, gotWIF.String(),
				wantAddr.privKeyWIF,
			)
			return false
		}
	default:
		// Confirm expected locked error.
		testName := fmt.Sprintf("%s ExportPrivKey", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrLocked) {
			return false
		}
	}
	// Imported addresses should return a nil derivation info.
	if _, _, ok := gotAddr.DerivationInfo(); gotAddr.Imported() && ok {
		tc.t.Errorf("%s Imported: expected nil derivation info", prefix)
		return false
	}
	return true
}

// testManagedScriptAddress ensures the data returned by all exported functions
// provided by the passed managed script address matches the corresponding
// fields in the provided expected address.
//
// When the test context indicates the manager is unlocked, the private data
// will also be tested, otherwise, the functions which deal with private data
// are checked to ensure they return the correct error.
func testManagedScriptAddress(
	tc *testContext,
	prefix string,
	gotAddr waddrmgr.ManagedScriptAddress,
	wantAddr *expectedAddr,
) bool {
	// Ensure script is the expected value for the managed address. Since this is
	// only available when the manager is unlocked, also check for the expected
	// error when the manager is locked.
	gotScript, e := gotAddr.Script()
	switch {
	case tc.watchingOnly:
		// Confirm expected watching-only error.
		testName := fmt.Sprintf("%s Script", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrWatchingOnly) {
			return false
		}
	case tc.unlocked:
		if e != nil {
			tc.t.Errorf(
				"%s Script: unexpected error - got %v",
				prefix, e,
			)
			return false
		}
		if !reflect.DeepEqual(gotScript, wantAddr.script) {
			tc.t.Errorf(
				"%s Script: unexpected script - got %x, "+
					"want %x", prefix, gotScript, wantAddr.script,
			)
			return false
		}
	default:
		// Confirm expected locked error.
		testName := fmt.Sprintf("%s Script", prefix)
		if !checkManagerError(tc.t, testName, e, waddrmgr.ErrLocked) {
			return false
		}
	}
	return true
}

// testAddress ensures the data returned by all exported functions provided by
// the passed managed address matches the corresponding fields in the provided
// expected address. It also type asserts the managed address to determine its
// specific type and calls the corresponding testing functions accordingly.
//
// When the test context indicates the manager is unlocked, the private data
// will also be tested, otherwise, the functions which deal with private data
// are checked to ensure they return the correct error.
func testAddress(tc *testContext, prefix string, gotAddr waddrmgr.ManagedAddress, wantAddr *expectedAddr) bool {
	if gotAddr.Account() != tc.account {
		tc.t.Errorf(
			"ManagedAddress.Account: unexpected account - got "+
				"%d, want %d", gotAddr.Account(), tc.account,
		)
		return false
	}
	if gotAddr.Address().EncodeAddress() != wantAddr.address {
		tc.t.Errorf(
			"%s EncodeAddress: unexpected address - got %s, "+
				"want %s", prefix, gotAddr.Address().EncodeAddress(),
			wantAddr.address,
		)
		return false
	}
	if !reflect.DeepEqual(gotAddr.AddrHash(), wantAddr.addressHash) {
		tc.t.Errorf(
			"%s AddrHash: unexpected address hash - got %x, "+
				"want %x", prefix, gotAddr.AddrHash(),
			wantAddr.addressHash,
		)
		return false
	}
	if gotAddr.Internal() != wantAddr.internal {
		tc.t.Errorf(
			"%s Internal: unexpected internal flag - got %v, "+
				"want %v", prefix, gotAddr.Internal(), wantAddr.internal,
		)
		return false
	}
	if gotAddr.Compressed() != wantAddr.compressed {
		tc.t.Errorf(
			"%s Compressed: unexpected compressed flag - got "+
				"%v, want %v", prefix, gotAddr.Compressed(),
			wantAddr.compressed,
		)
		return false
	}
	if gotAddr.Imported() != wantAddr.imported {
		tc.t.Errorf(
			"%s Imported: unexpected imported flag - got %v, "+
				"want %v", prefix, gotAddr.Imported(), wantAddr.imported,
		)
		return false
	}
	switch addr := gotAddr.(type) {
	case waddrmgr.ManagedPubKeyAddress:
		if !testManagedPubKeyAddress(tc, prefix, addr, wantAddr) {
			return false
		}
	case waddrmgr.ManagedScriptAddress:
		if !testManagedScriptAddress(tc, prefix, addr, wantAddr) {
			return false
		}
	}
	return true
}

// testExternalAddresses tests several facets of external addresses such as
// generating multiple addresses via NextExternalAddresses, ensuring they can be
// retrieved by Address, and that they work properly when the manager is locked
// and unlocked.
func testExternalAddresses(tc *testContext) bool {
	prefix := testNamePrefix(tc) + " testExternalAddresses"
	var addrs []waddrmgr.ManagedAddress
	if tc.create {
		prefix = prefix + " NextExternalAddresses"
		e := walletdb.Update(
			tc.db, func(tx walletdb.ReadWriteTx) (e error) {
				ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				addrs, e = tc.manager.NextExternalAddresses(ns, tc.account, 5)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("%s: unexpected error: %v", prefix, e)
			return false
		}
		if len(addrs) != len(expectedExternalAddrs) {
			tc.t.Errorf(
				"%s: unexpected number of addresses - got "+
					"%d, want %d", prefix, len(addrs),
				len(expectedExternalAddrs),
			)
			return false
		}
	}
	// Setup a closure to test the results since the same tests need to be repeated
	// with the manager locked and unlocked.
	testResults := func() bool {
		// Ensure the returned addresses are the expected ones. When not in the create
		// phase, there will be no addresses in the addrs slice, so this really only
		// runs during the first phase of the tests.
		for i := 0; i < len(addrs); i++ {
			prefix = fmt.Sprintf("%s ExternalAddress #%d", prefix, i)
			if !testAddress(tc, prefix, addrs[i], &expectedExternalAddrs[i]) {
				return false
			}
		}
		// Ensure the last external address is the expected one.
		leaPrefix := prefix + " LastExternalAddress"
		var lastAddr waddrmgr.ManagedAddress
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				lastAddr, e = tc.manager.LastExternalAddress(ns, tc.account)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("%s: unexpected error: %v", leaPrefix, e)
			return false
		}
		if !testAddress(tc, leaPrefix, lastAddr, &expectedExternalAddrs[len(expectedExternalAddrs)-1]) {
			return false
		}
		// Now, use the Address API to retrieve each of the expected new addresses and
		// ensure they're accurate.
		chainParams := tc.manager.ChainParams()
		for i := 0; i < len(expectedExternalAddrs); i++ {
			pkHash := expectedExternalAddrs[i].addressHash
			utilAddr, e := btcaddr.NewPubKeyHash(
				pkHash, chainParams,
			)
			if e != nil {
				tc.t.Errorf(
					"%s NewPubKeyHash #%d: "+
						"unexpected error: %v", prefix, i, e,
				)
				return false
			}
			prefix = fmt.Sprintf("%s Address #%d", prefix, i)
			var addr waddrmgr.ManagedAddress
			e = walletdb.View(
				tc.db, func(tx walletdb.ReadTx) (e error) {
					ns := tx.ReadBucket(waddrmgrNamespaceKey)
					addr, e = tc.manager.Address(ns, utilAddr)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s: unexpected error: %v", prefix,
					e,
				)
				return false
			}
			if !testAddress(tc, prefix, addr, &expectedExternalAddrs[i]) {
				return false
			}
		}
		return true
	}
	// Since the manager is locked at this point, the public address information is
	// tested and the private functions are checked to ensure they return the
	// expected error.
	if !testResults() {
		return false
	}
	// Everything after this point involves retesting with an unlocked address
	// manager which is not possible for watching-only mode, so just exit now in
	// that case.
	if tc.watchingOnly {
		return true
	}
	// Unlock the manager and retest all of the addresses to ensure the private
	// information is valid as well.
	e := walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.rootManager.Unlock(ns, privPassphrase)
		},
	)
	if e != nil {
		tc.t.Errorf("Unlock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = true
	if !testResults() {
		return false
	}
	// Relock the manager for future tests.
	if e := tc.rootManager.Lock(); E.Chk(e) {
		tc.t.Errorf("Lock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = false
	return true
}

// testInternalAddresses tests several facets of internal addresses such as
// generating multiple addresses via NextInternalAddresses, ensuring they can be
// retrieved by Address, and that they work properly when the manager is locked
// and unlocked.
func testInternalAddresses(tc *testContext) bool {
	// When the address manager is not in watching-only mode, unlocked it first to
	// ensure that address generation works correctly when the address manager is
	// unlocked and then locked later. These tests reverse the order done in the
	// external tests which starts with a locked manager and unlock it afterwards.
	if !tc.watchingOnly {
		// Unlock the manager and retest all of the addresses to ensure the private
		// information is valid as well.
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				return tc.rootManager.Unlock(ns, privPassphrase)
			},
		)
		if e != nil {
			tc.t.Errorf("Unlock: unexpected error: %v", e)
			return false
		}
		tc.unlocked = true
	}
	prefix := testNamePrefix(tc) + " testInternalAddresses"
	var addrs []waddrmgr.ManagedAddress
	if tc.create {
		prefix = prefix + " NextInternalAddress"
		e := walletdb.Update(
			tc.db, func(tx walletdb.ReadWriteTx) (e error) {
				ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				addrs, e = tc.manager.NextInternalAddresses(ns, tc.account, 5)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("%s: unexpected error: %v", prefix, e)
			return false
		}
		if len(addrs) != len(expectedInternalAddrs) {
			tc.t.Errorf(
				"%s: unexpected number of addresses - got "+
					"%d, want %d", prefix, len(addrs),
				len(expectedInternalAddrs),
			)
			return false
		}
	}
	// Setup a closure to test the results since the same tests need to be repeated
	// with the manager locked and unlocked.
	testResults := func() bool {
		// Ensure the returned addresses are the expected ones. When not in the create
		// phase, there will be no addresses in the addrs slice, so this really only
		// runs during the first phase of the tests.
		for i := 0; i < len(addrs); i++ {
			prefix = fmt.Sprintf("%s InternalAddress #%d", prefix, i)
			if !testAddress(tc, prefix, addrs[i], &expectedInternalAddrs[i]) {
				return false
			}
		}
		// Ensure the last internal address is the expected one.
		liaPrefix := prefix + " LastInternalAddress"
		var lastAddr waddrmgr.ManagedAddress
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				lastAddr, e = tc.manager.LastInternalAddress(ns, tc.account)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("%s: unexpected error: %v", liaPrefix, e)
			return false
		}
		if !testAddress(tc, liaPrefix, lastAddr, &expectedInternalAddrs[len(expectedInternalAddrs)-1]) {
			return false
		}
		// Now, use the Address API to retrieve each of the expected new addresses and
		// ensure they're accurate.
		chainParams := tc.manager.ChainParams()
		for i := 0; i < len(expectedInternalAddrs); i++ {
			pkHash := expectedInternalAddrs[i].addressHash
			utilAddr, e := btcaddr.NewPubKeyHash(
				pkHash, chainParams,
			)
			if e != nil {
				tc.t.Errorf(
					"%s NewPubKeyHash #%d: "+
						"unexpected error: %v", prefix, i, e,
				)
				return false
			}
			prefix = fmt.Sprintf("%s Address #%d", prefix, i)
			var addr waddrmgr.ManagedAddress
			e = walletdb.View(
				tc.db, func(tx walletdb.ReadTx) (e error) {
					ns := tx.ReadBucket(waddrmgrNamespaceKey)
					addr, e = tc.manager.Address(ns, utilAddr)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s: unexpected error: %v", prefix,
					e,
				)
				return false
			}
			if !testAddress(tc, prefix, addr, &expectedInternalAddrs[i]) {
				return false
			}
		}
		return true
	}
	// The address manager could either be locked or unlocked here depending on
	// whether or not it's a watching-only manager. When it's unlocked, this will
	// test both the public and private address data are accurate. When it's locked,
	// it must be watching-only, so only the public address information is tested
	// and the private functions are checked to ensure they return the expected
	// ErrWatchingOnly error.
	if !testResults() {
		return false
	}
	// Everything after this point involves locking the address manager and
	// retesting the addresses with a locked manager. However, for watching-only
	// mode, this has already happened, so just exit now in that case.
	if tc.watchingOnly {
		return true
	}
	// Lock the manager and retest all of the addresses to ensure the public
	// information remains valid and the private functions return the expected
	// error.
	if e := tc.rootManager.Lock(); E.Chk(e) {
		tc.t.Errorf("Lock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = false
	if !testResults() {
		return false
	}
	return true
}

// testLocking tests the basic locking semantics of the address manager work as
// expected. Other tests ensure addresses behave as expected under locked and
// unlocked conditions.
func testLocking(tc *testContext) bool {
	if tc.unlocked {
		tc.t.Error("testLocking called with an unlocked manager")
		return false
	}
	if !tc.rootManager.IsLocked() {
		tc.t.Error("IsLocked: returned false on locked manager")
		return false
	}
	// Locking an already lock manager should return an error. The error should be
	// ErrLocked or ErrWatchingOnly depending on the type of the address manager.
	e := tc.rootManager.Lock()
	wantErrCode := waddrmgr.ErrLocked
	if tc.watchingOnly {
		wantErrCode = waddrmgr.ErrWatchingOnly
	}
	if !checkManagerError(tc.t, "Lock", e, wantErrCode) {
		return false
	}
	// Ensure unlocking with the correct passphrase doesn't return any unexpected
	// errors and the manager properly reports it is unlocked. Since watching-only
	// address managers can't be unlocked, also ensure the correct error for that
	// case.
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.rootManager.Unlock(ns, privPassphrase)
		},
	)
	if tc.watchingOnly {
		if !checkManagerError(tc.t, "Unlock", e, waddrmgr.ErrWatchingOnly) {
			return false
		}
	} else if e != nil {
		tc.t.Errorf("Unlock: unexpected error: %v", e)
		return false
	}
	if !tc.watchingOnly && tc.rootManager.IsLocked() {
		tc.t.Error("IsLocked: returned true on unlocked manager")
		return false
	}
	// Unlocking the manager again is allowed. Since watching-only address managers
	// can't be unlocked, also ensure the correct error for that case.
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.rootManager.Unlock(ns, privPassphrase)
		},
	)
	if tc.watchingOnly {
		if !checkManagerError(tc.t, "Unlock2", e, waddrmgr.ErrWatchingOnly) {
			return false
		}
	} else if e != nil {
		tc.t.Errorf("Unlock: unexpected error: %v", e)
		return false
	}
	if !tc.watchingOnly && tc.rootManager.IsLocked() {
		tc.t.Error("IsLocked: returned true on unlocked manager")
		return false
	}
	// Unlocking the manager with an invalid passphrase must result in an error and
	// a locked manager.
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.rootManager.Unlock(ns, []byte("invalidpassphrase"))
		},
	)
	wantErrCode = waddrmgr.ErrWrongPassphrase
	if tc.watchingOnly {
		wantErrCode = waddrmgr.ErrWatchingOnly
	}
	if !checkManagerError(tc.t, "Unlock", e, wantErrCode) {
		return false
	}
	if !tc.rootManager.IsLocked() {
		tc.t.Error(
			"IsLocked: manager is unlocked after failed unlock " +
				"attempt",
		)
		return false
	}
	return true
}

// testImportPrivateKey tests that importing private keys works properly. It
// ensures they can be retrieved by Address after they have been imported and
// the addresses give the expected values when the manager is locked and
// unlocked.
//
// This function expects the manager is already locked when called and returns
// with the manager locked.
func testImportPrivateKey(tc *testContext) bool {
	tests := []struct {
		name       string
		in         string
		blockstamp waddrmgr.BlockStamp
		expected   expectedAddr
	}{
		{
			name: "wif for uncompressed pubkey address",
			in:   "5HueCGU8rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
			expected: expectedAddr{
				address:     "1GAehh7TsJAHuUAeKZcXf5CnwuGuGgyX2S",
				addressHash: hexToBytes("a65d1a239d4ec666643d350c7bb8fc44d2881128"),
				internal:    false,
				imported:    true,
				compressed:  false,
				pubKey: hexToBytes(
					"04d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3" +
						"d66a2c5f10780d95b7df42645cd85228a6fb29940e858e7e558" +
						"42ae2bd115d1ed7cc0e82d934e929c97648cb0a",
				),
				privKey: hexToBytes("0c28fca386c7a227600b2fe50b7cae11ec86d3bf1fbe471be89827e19d72aa1d"),
				// privKeyWIF is set to the in field during tests
			},
		},
		{
			name: "wif for compressed pubkey address",
			in:   "KwdMAjGmerYanjeui5SHS7JkmpZvVipYvB2LJGU1ZxJwYvP98617",
			expected: expectedAddr{
				address:     "1LoVGDgRs9hTfTNJNuXKSpywcbdvwRXpmK",
				addressHash: hexToBytes("d9351dcbad5b8f3b8bfa2f2cdc85c28118ca9326"),
				internal:    false,
				imported:    true,
				compressed:  true,
				pubKey:      hexToBytes("02d0de0aaeaefad02b8bdc8a01a1b8b11c696bd3d66a2c5f10780d95b7df42645c"),
				privKey:     hexToBytes("0c28fca386c7a227600b2fe50b7cae11ec86d3bf1fbe471be89827e19d72aa1d"),
				// privKeyWIF is set to the in field during tests
			},
		},
	}
	// The manager must be unlocked to import a private key, however a watching-only
	// manager can't be unlocked.
	if !tc.watchingOnly {
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				return tc.rootManager.Unlock(ns, privPassphrase)
			},
		)
		if e != nil {
			tc.t.Errorf("Unlock: unexpected error: %v", e)
			return false
		}
		tc.unlocked = true
	}
	// Only import the private keys when in the create phase of testing.
	tc.account = waddrmgr.ImportedAddrAccount
	prefix := testNamePrefix(tc) + " testImportPrivateKey"
	if tc.create {
		for i, test := range tests {
			test.expected.privKeyWIF = test.in
			wif, e := util.DecodeWIF(test.in)
			if e != nil {
				tc.t.Errorf(
					"%s DecodeWIF #%d (%s): unexpected "+
						"error: %v", prefix, i, test.name, e,
				)
				continue
			}
			var addr waddrmgr.ManagedPubKeyAddress
			e = walletdb.Update(
				tc.db, func(tx walletdb.ReadWriteTx) (e error) {
					ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
					addr, e = tc.manager.ImportPrivateKey(ns, wif, &test.blockstamp)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s ImportPrivateKey #%d (%s): "+
						"unexpected error: %v", prefix, i,
					test.name, e,
				)
				continue
			}
			if !testAddress(
				tc, prefix+" ImportPrivateKey", addr,
				&test.expected,
			) {
				continue
			}
		}
	}
	// Setup a closure to test the results since the same tests need to be repeated
	// with the manager unlocked and locked.
	chainParams := tc.manager.ChainParams()
	testResults := func() bool {
		failed := false
		for i, test := range tests {
			test.expected.privKeyWIF = test.in
			// Use the Address API to retrieve each of the expected new addresses and ensure
			// they're accurate.
			utilAddr, e := btcaddr.NewPubKeyHash(
				test.expected.addressHash, chainParams,
			)
			if e != nil {
				tc.t.Errorf(
					"%s NewPubKeyHash #%d (%s): "+
						"unexpected error: %v", prefix, i,
					test.name, e,
				)
				failed = true
				continue
			}
			taPrefix := fmt.Sprintf(
				"%s Address #%d (%s)", prefix,
				i, test.name,
			)
			var ma waddrmgr.ManagedAddress
			e = walletdb.View(
				tc.db, func(tx walletdb.ReadTx) (e error) {
					ns := tx.ReadBucket(waddrmgrNamespaceKey)
					ma, e = tc.manager.Address(ns, utilAddr)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s: unexpected error: %v", taPrefix,
					e,
				)
				failed = true
				continue
			}
			if !testAddress(tc, taPrefix, ma, &test.expected) {
				failed = true
				continue
			}
		}
		return !failed
	}
	// The address manager could either be locked or unlocked here depending on
	// whether or not it's a watching-only manager. When it's unlocked, this will
	// test both the public and private address data are accurate. When it's locked,
	// it must be watching-only, so only the public address information is tested
	// and the private functions are checked to ensure they return the expected
	// ErrWatchingOnly error.
	if !testResults() {
		return false
	}
	// Everything after this point involves locking the address manager and
	// retesting the addresses with a locked manager. However, for watching-only
	// mode, this has already happened, so just exit now in that case.
	if tc.watchingOnly {
		return true
	}
	// Lock the manager and retest all of the addresses to ensure the private
	// information returns the expected error.
	if e := tc.rootManager.Lock(); E.Chk(e) {
		tc.t.Errorf("Lock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = false
	if !testResults() {
		return false
	}
	return true
}

// testImportScript tests that importing scripts works properly. It ensures they
// can be retrieved by Address after they have been imported and the addresses
// give the expected values when the manager is locked and unlocked.
//
// This function expects the manager is already locked when called and returns
// with the manager locked.
func testImportScript(tc *testContext) bool {
	tests := []struct {
		name       string
		in         []byte
		blockstamp waddrmgr.BlockStamp
		expected   expectedAddr
	}{
		{
			name: "p2sh uncompressed pubkey",
			in: hexToBytes(
				"41048b65a0e6bb200e6dac05e74281b1ab9a41e8" +
					"0006d6b12d8521e09981da97dd96ac72d24d1a7d" +
					"ed9493a9fc20fdb4a714808f0b680f1f1d935277" +
					"48b5e3f629ffac",
			),
			expected: expectedAddr{
				address:     "3MbyWAu9UaoBewR3cArF1nwf4aQgVwzrA5",
				addressHash: hexToBytes("da6e6a632d96dc5530d7b3c9f3017725d023093e"),
				internal:    false,
				imported:    true,
				compressed:  false,
				// script is set to the in field during tests.
			},
		},
		{
			name: "p2sh multisig",
			in: hexToBytes(
				"524104cb9c3c222c5f7a7d3b9bd152f363a0b6d5" +
					"4c9eb312c4d4f9af1e8551b6c421a6a4ab0e2910" +
					"5f24de20ff463c1c91fcf3bf662cdde4783d4799" +
					"f787cb7c08869b4104ccc588420deeebea22a7e9" +
					"00cc8b68620d2212c374604e3487ca08f1ff3ae1" +
					"2bdc639514d0ec8612a2d3c519f084d9a00cbbe3" +
					"b53d071e9b09e71e610b036aa24104ab47ad1939" +
					"edcb3db65f7fedea62bbf781c5410d3f22a7a3a5" +
					"6ffefb2238af8627363bdf2ed97c1f89784a1aec" +
					"db43384f11d2acc64443c7fc299cef0400421a53ae",
			),
			expected: expectedAddr{
				address:     "34CRZpt8j81rgh9QhzuBepqPi4cBQSjhjr",
				addressHash: hexToBytes("1b800cec1fe92222f36a502c139bed47c5959715"),
				internal:    false,
				imported:    true,
				compressed:  false,
				// script is set to the in field during tests.
			},
		},
	}
	// The manager must be unlocked to import a private key and also for testing
	// private data. However, a watching-only manager can't be unlocked.
	if !tc.watchingOnly {
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				return tc.rootManager.Unlock(ns, privPassphrase)
			},
		)
		if e != nil {
			tc.t.Errorf("Unlock: unexpected error: %v", e)
			return false
		}
		tc.unlocked = true
	}
	// Only import the scripts when in the create phase of testing.
	tc.account = waddrmgr.ImportedAddrAccount
	prefix := testNamePrefix(tc)
	if tc.create {
		for i, test := range tests {
			test.expected.script = test.in
			prefix = fmt.Sprintf(
				"%s ImportScript #%d (%s)", prefix,
				i, test.name,
			)
			var addr waddrmgr.ManagedScriptAddress
			e := walletdb.Update(
				tc.db, func(tx walletdb.ReadWriteTx) (e error) {
					ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
					addr, e = tc.manager.ImportScript(ns, test.in, &test.blockstamp)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s: unexpected error: %v", prefix,
					e,
				)
				continue
			}
			if !testAddress(tc, prefix, addr, &test.expected) {
				continue
			}
		}
	}
	// Setup a closure to test the results since the same tests need to be repeated
	// with the manager unlocked and locked.
	chainParams := tc.manager.ChainParams()
	testResults := func() bool {
		failed := false
		for i, test := range tests {
			test.expected.script = test.in
			// Use the Address API to retrieve each of the expected new addresses and ensure
			// they're accurate.
			utilAddr, e := btcaddr.NewScriptHash(
				test.in,
				chainParams,
			)
			if e != nil {
				tc.t.Errorf(
					"%s NewScriptHash #%d (%s): "+
						"unexpected error: %v", prefix, i,
					test.name, e,
				)
				failed = true
				continue
			}
			taPrefix := fmt.Sprintf(
				"%s Address #%d (%s)", prefix,
				i, test.name,
			)
			var ma waddrmgr.ManagedAddress
			e = walletdb.View(
				tc.db, func(tx walletdb.ReadTx) (e error) {
					ns := tx.ReadBucket(waddrmgrNamespaceKey)
					ma, e = tc.manager.Address(ns, utilAddr)
					return e
				},
			)
			if e != nil {
				tc.t.Errorf(
					"%s: unexpected error: %v", taPrefix,
					e,
				)
				failed = true
				continue
			}
			if !testAddress(tc, taPrefix, ma, &test.expected) {
				failed = true
				continue
			}
		}
		return !failed
	}
	// The address manager could either be locked or unlocked here depending on
	// whether or not it's a watching-only manager. When it's unlocked, this will
	// test both the public and private address data are accurate. When it's locked,
	// it must be watching-only, so only the public address information is tested
	// and the private functions are checked to ensure they return the expected
	// ErrWatchingOnly error.
	if !testResults() {
		return false
	}
	// Everything after this point involves locking the address manager and
	// retesting the addresses with a locked manager. However, for watching-only
	// mode, this has already happened, so just exit now in that case.
	if tc.watchingOnly {
		return true
	}
	// Lock the manager and retest all of the addresses to ensure the private
	// information returns the expected error.
	if e := tc.rootManager.Lock(); E.Chk(e) {
		tc.t.Errorf("Lock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = false
	if !testResults() {
		return false
	}
	return true
}

// testMarkUsed ensures used addresses are flagged as such.
func testMarkUsed(tc *testContext) bool {
	tests := []struct {
		name string
		typ  addrType
		in   []byte
	}{
		{
			name: "managed address",
			typ:  addrPubKeyHash,
			in:   hexToBytes("2ef94abb9ee8f785d087c3ec8d6ee467e92d0d0a"),
		},
		{
			name: "script address",
			typ:  addrScriptHash,
			in:   hexToBytes("da6e6a632d96dc5530d7b3c9f3017725d023093e"),
		},
	}
	prefix := "MarkUsed"
	chainParams := tc.manager.ChainParams()
	for i, test := range tests {
		addrHash := test.in
		var addr btcaddr.Address
		var e error
		switch test.typ {
		case addrPubKeyHash:
			addr, e = btcaddr.NewPubKeyHash(addrHash, chainParams)
		case addrScriptHash:
			addr, e = btcaddr.NewScriptHashFromHash(addrHash, chainParams)
		default:
			panic("unreachable")
		}
		if e != nil {
			tc.t.Errorf("%s #%d: NewAddress unexpected error: %v", prefix, i, e)
			continue
		}
		e = walletdb.Update(
			tc.db, func(tx walletdb.ReadWriteTx) (e error) {
				ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				maddr, e := tc.manager.Address(ns, addr)
				if e != nil {
					tc.t.Errorf("%s #%d: Address unexpected error: %v", prefix, i, e)
					return nil
				}
				if tc.create {
					// Test that initially the address is not flagged as used
					used := maddr.Used(ns)
					if used != false {
						tc.t.Errorf(
							"%s #%d: unexpected used flag -- got "+
								"%v, want %v", prefix, i, used, false,
						)
					}
				}
				e = tc.manager.MarkUsed(ns, addr)
				if e != nil {
					tc.t.Errorf("%s #%d: unexpected error: %v", prefix, i, e)
					return nil
				}
				used := maddr.Used(ns)
				if used != true {
					tc.t.Errorf(
						"%s #%d: unexpected used flag -- got "+
							"%v, want %v", prefix, i, used, true,
					)
				}
				return nil
			},
		)
		if e != nil {
			tc.t.Errorf("Unexpected error %v", e)
		}
	}
	return true
}

// testChangePassphrase ensures changes both the public and private passphrases
// works as intended.
func testChangePassphrase(tc *testContext) bool {
	// Force an error when changing the passphrase due to failure to generate a new
	// secret key by replacing the generation function one that intentionally
	// errors.
	testName := "ChangePassphrase (public) with invalid new secret key"
	oldKeyGen := waddrmgr.SetSecretKeyGen(failingSecretKeyGen)
	e := walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, pubPassphrase, pubPassphrase2, false, fastScrypt,
			)
		},
	)
	if !checkManagerError(tc.t, testName, e, waddrmgr.ErrCrypto) {
		return false
	}
	// Attempt to change public passphrase with invalid old passphrase.
	testName = "ChangePassphrase (public) with invalid old passphrase"
	waddrmgr.SetSecretKeyGen(oldKeyGen)
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, []byte("bogus"), pubPassphrase2, false, fastScrypt,
			)
		},
	)
	if !checkManagerError(tc.t, testName, e, waddrmgr.ErrWrongPassphrase) {
		return false
	}
	// Change the public passphrase.
	testName = "ChangePassphrase (public)"
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, pubPassphrase, pubPassphrase2, false, fastScrypt,
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", testName, e)
		return false
	}
	// Ensure the public passphrase was successfully changed.
	if !tc.rootManager.TstCheckPublicPassphrase(pubPassphrase2) {
		tc.t.Errorf("%s: passphrase does not match", testName)
		return false
	}
	// Change the private passphrase back to what it was.
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, pubPassphrase2, pubPassphrase, false, fastScrypt,
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", testName, e)
		return false
	}
	// Attempt to change private passphrase with invalid old passphrase. The error
	// should be ErrWrongPassphrase or ErrWatchingOnly depending on the type of the
	// address manager.
	testName = "ChangePassphrase (private) with invalid old passphrase"
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, []byte("bogus"), privPassphrase2, true, fastScrypt,
			)
		},
	)
	wantErrCode := waddrmgr.ErrWrongPassphrase
	if tc.watchingOnly {
		wantErrCode = waddrmgr.ErrWatchingOnly
	}
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	// Everything after this point involves testing that the private passphrase for
	// the address manager can be changed successfully. This is not possible for
	// watching-only mode, so just exit now in that case.
	if tc.watchingOnly {
		return true
	}
	// Change the private passphrase.
	testName = "ChangePassphrase (private)"
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, privPassphrase, privPassphrase2, true, fastScrypt,
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", testName, e)
		return false
	}
	// Unlock the manager with the new passphrase to ensure it changed as expected.
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.Unlock(ns, privPassphrase2)
		},
	)
	if e != nil {
		tc.t.Errorf(
			"%s: failed to unlock with new private "+
				"passphrase: %v", testName, e,
		)
		return false
	}
	tc.unlocked = true
	// Change the private passphrase back to what it was while the manager is
	// unlocked to ensure that path works properly as well.
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.ChangePassphrase(
				ns, privPassphrase2, privPassphrase, true, fastScrypt,
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", testName, e)
		return false
	}
	if tc.rootManager.IsLocked() {
		tc.t.Errorf("%s: manager is locked", testName)
		return false
	}
	// Relock the manager for future tests.
	if e := tc.rootManager.Lock(); E.Chk(e) {
		tc.t.Errorf("Lock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = false
	return true
}

// testNewAccount tests the new account creation func of the address manager
// works as expected.
func testNewAccount(tc *testContext) bool {
	if tc.watchingOnly {
		// Creating new accounts in watching-only mode should return ErrWatchingOnly
		e := walletdb.Update(
			tc.db, func(tx walletdb.ReadWriteTx) (e error) {
				ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
				_, e = tc.manager.NewAccount(ns, "test")
				return e
			},
		)
		if !checkManagerError(
			tc.t, "Create account in watching-only mode", e,
			waddrmgr.ErrWatchingOnly,
		) {
			tc.manager.Close()
			return false
		}
		return true
	}
	// Creating new accounts when wallet is locked should return ErrLocked
	e := walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.NewAccount(ns, "test")
			return e
		},
	)
	if !checkManagerError(
		tc.t, "Create account when wallet is locked", e,
		waddrmgr.ErrLocked,
	) {
		tc.manager.Close()
		return false
	}
	// Unlock the wallet to decrypt cointype keys required to derive account keys
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			e = tc.rootManager.Unlock(ns, privPassphrase)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("Unlock: unexpected error: %v", e)
		return false
	}
	tc.unlocked = true
	testName := "acct-create"
	expectedAccount := tc.account + 1
	if !tc.create {
		// Create a new account in open mode
		testName = "acct-open"
		expectedAccount++
	}
	var account uint32
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			account, e = tc.manager.NewAccount(ns, testName)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("NewAccount: unexpected error: %v", e)
		return false
	}
	if account != expectedAccount {
		tc.t.Errorf(
			"NewAccount "+
				"account mismatch -- got %d, "+
				"want %d", account, expectedAccount,
		)
		return false
	}
	// Test duplicate account name error
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.NewAccount(ns, testName)
			return e
		},
	)
	wantErrCode := waddrmgr.ErrDuplicateAccount
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	// Test account name validation
	testName = "" // Empty account names are not allowed
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.NewAccount(ns, testName)
			return e
		},
	)
	wantErrCode = waddrmgr.ErrInvalidAccount
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	testName = "imported" // A reserved account name
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.NewAccount(ns, testName)
			return e
		},
	)
	wantErrCode = waddrmgr.ErrInvalidAccount
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	return true
}

// testLookupAccount tests the basic account lookup func of the address manager
// works as expected.
func testLookupAccount(tc *testContext) bool {
	// Lookup accounts created earlier in testNewAccount
	expectedAccounts := map[string]uint32{
		waddrmgr.TstDefaultAccountName:   waddrmgr.DefaultAccountNum,
		waddrmgr.ImportedAddrAccountName: waddrmgr.ImportedAddrAccount,
	}
	for acctName, expectedAccount := range expectedAccounts {
		var account uint32
		e := walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				account, e = tc.manager.LookupAccount(ns, acctName)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("LookupAccount: unexpected error: %v", e)
			return false
		}
		if account != expectedAccount {
			tc.t.Errorf(
				"LookupAccount "+
					"account mismatch -- got %d, "+
					"want %d", account, expectedAccount,
			)
			return false
		}
	}
	// Test account not found error
	testName := "non existent account"
	e := walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.LookupAccount(ns, testName)
			return e
		},
	)
	wantErrCode := waddrmgr.ErrAccountNotFound
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	// Test last account
	var lastAccount uint32
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			lastAccount, e = tc.manager.LastAccount(ns)
			return e
		},
	)
	var expectedLastAccount uint32
	expectedLastAccount = 1
	if !tc.create {
		// Existing wallet manager will have 3 accounts
		expectedLastAccount = 2
	}
	if lastAccount != expectedLastAccount {
		tc.t.Errorf(
			"LookupAccount "+
				"account mismatch -- got %d, "+
				"want %d", lastAccount, expectedLastAccount,
		)
		return false
	}
	// Test account lookup for default account adddress
	var expectedAccount uint32
	for i, addr := range expectedAddrs {
		addr, e := btcaddr.NewPubKeyHash(
			addr.addressHash,
			tc.manager.ChainParams(),
		)
		if e != nil {
			tc.t.Errorf("AddrAccount #%d: unexpected error: %v", i, e)
			return false
		}
		var account uint32
		e = walletdb.View(
			tc.db, func(tx walletdb.ReadTx) (e error) {
				ns := tx.ReadBucket(waddrmgrNamespaceKey)
				account, e = tc.manager.AddrAccount(ns, addr)
				return e
			},
		)
		if e != nil {
			tc.t.Errorf("AddrAccount #%d: unexpected error: %v", i, e)
			return false
		}
		if account != expectedAccount {
			tc.t.Errorf(
				"AddrAccount "+
					"account mismatch -- got %d, "+
					"want %d", account, expectedAccount,
			)
			return false
		}
	}
	return true
}

// testRenameAccount tests the rename account func of the address manager works
// as expected.
func testRenameAccount(tc *testContext) bool {
	var acctName string
	e := walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			acctName, e = tc.manager.AccountName(ns, tc.account)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("AccountName: unexpected error: %v", e)
		return false
	}
	testName := acctName + "-renamed"
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.manager.RenameAccount(ns, tc.account, testName)
		},
	)
	if e != nil {
		tc.t.Errorf("RenameAccount: unexpected error: %v", e)
		return false
	}
	var newName string
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			newName, e = tc.manager.AccountName(ns, tc.account)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("AccountName: unexpected error: %v", e)
		return false
	}
	if newName != testName {
		tc.t.Errorf(
			"RenameAccount "+
				"account name mismatch -- got %s, "+
				"want %s", newName, testName,
		)
		return false
	}
	// Test duplicate account name error
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.manager.RenameAccount(ns, tc.account, testName)
		},
	)
	wantErrCode := waddrmgr.ErrDuplicateAccount
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	// Test old account name is no longer valid
	e = walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			_, e = tc.manager.LookupAccount(ns, acctName)
			return e
		},
	)
	wantErrCode = waddrmgr.ErrAccountNotFound
	if !checkManagerError(tc.t, testName, e, wantErrCode) {
		return false
	}
	return true
}

// testForEachAccount tests the retrieve all accounts func of the address
// manager works as expected.
func testForEachAccount(tc *testContext) bool {
	prefix := testNamePrefix(tc) + " testForEachAccount"
	expectedAccounts := []uint32{0, 1}
	if !tc.create {
		// Existing wallet manager will have 3 accounts
		expectedAccounts = append(expectedAccounts, 2)
	}
	// Imported account
	expectedAccounts = append(expectedAccounts, waddrmgr.ImportedAddrAccount)
	var accounts []uint32
	e := walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.manager.ForEachAccount(
				ns, func(account uint32) (e error) {
					accounts = append(accounts, account)
					return nil
				},
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", prefix, e)
		return false
	}
	if len(accounts) != len(expectedAccounts) {
		tc.t.Errorf(
			"%s: unexpected number of accounts - got "+
				"%d, want %d", prefix, len(accounts),
			len(expectedAccounts),
		)
		return false
	}
	for i, account := range accounts {
		if expectedAccounts[i] != account {
			tc.t.Errorf(
				"%s #%d: "+
					"account mismatch -- got %d, "+
					"want %d", prefix, i, account, expectedAccounts[i],
			)
		}
	}
	return true
}

// testForEachAccountAddress tests that iterating through the given account
// addresses using the manager API works as expected.
func testForEachAccountAddress(tc *testContext) bool {
	prefix := testNamePrefix(tc) + " testForEachAccountAddress"
	// Make a map of expected addresses
	expectedAddrMap := make(map[string]*expectedAddr, len(expectedAddrs))
	for i := 0; i < len(expectedAddrs); i++ {
		expectedAddrMap[expectedAddrs[i].address] = &expectedAddrs[i]
	}
	var addrs []waddrmgr.ManagedAddress
	e := walletdb.View(
		tc.db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return tc.manager.ForEachAccountAddress(
				ns, tc.account,
				func(maddr waddrmgr.ManagedAddress) (e error) {
					addrs = append(addrs, maddr)
					return nil
				},
			)
		},
	)
	if e != nil {
		tc.t.Errorf("%s: unexpected error: %v", prefix, e)
		return false
	}
	for i := 0; i < len(addrs); i++ {
		prefix = fmt.Sprintf("%s: #%d", prefix, i)
		gotAddr := addrs[i]
		wantAddr := expectedAddrMap[gotAddr.Address().String()]
		if !testAddress(tc, prefix, gotAddr, wantAddr) {
			return false
		}
		delete(expectedAddrMap, gotAddr.Address().String())
	}
	if len(expectedAddrMap) != 0 {
		tc.t.Errorf(
			"%s: unexpected addresses -- got %d, want %d", prefix,
			len(expectedAddrMap), 0,
		)
		return false
	}
	return true
}

// testManagerAPI tests the functions provided by the Manager API as well as the
// ManagedAddress, ManagedPubKeyAddress, and ManagedScriptAddress interfaces.
func testManagerAPI(tc *testContext) {
	testLocking(tc)
	testExternalAddresses(tc)
	testInternalAddresses(tc)
	testImportPrivateKey(tc)
	testImportScript(tc)
	testMarkUsed(tc)
	testChangePassphrase(tc)
	// Reset default account
	tc.account = 0
	testNewAccount(tc)
	testLookupAccount(tc)
	testForEachAccount(tc)
	testForEachAccountAddress(tc)
	// Rename account 1 "acct-create"
	tc.account = 1
	testRenameAccount(tc)
}

// testWatchingOnly tests various facets of a watching-only address manager such
// as running the full set of API tests against a newly converted copy as well
// as when it is opened from an existing namespace.
func testWatchingOnly(tc *testContext) bool {
	// Make a copy of the current database so the copy can be converted to watching
	// only.
	woMgrName := "mgrtestwo.bin"
	_ = os.Remove(woMgrName)
	fi, e := os.OpenFile(woMgrName, os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		tc.t.Errorf("%v", e)
		return false
	}
	if e = tc.db.Copy(fi); E.Chk(e) {
		if e = fi.Close(); waddrmgr.E.Chk(e) {
		}
		tc.t.Errorf("%v", e)
		return false
	}
	if e = fi.Close(); waddrmgr.E.Chk(e) {
	}
	defer func() {
		if e = os.Remove(woMgrName); waddrmgr.E.Chk(e) {
		}
	}()
	// Open the new database copy and get the address manager namespace.
	var db walletdb.DB
	if db, e = walletdb.Open("bdb", woMgrName); waddrmgr.E.Chk(e) {
		tc.t.Errorf("openDbNamespace: unexpected error: %v", e)
		return false
	}
	
	defer func() {
		if e = db.Close(); waddrmgr.E.Chk(e) {
		}
	}()
	// Open the manager using the namespace and convert it to watching-only.
	var mgr *waddrmgr.Manager
	e = walletdb.View(
		db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			mgr, e = waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("%v", e)
		return false
	}
	e = walletdb.Update(
		db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return mgr.ConvertToWatchingOnly(ns)
		},
	)
	if e != nil {
		tc.t.Errorf("%v", e)
		return false
	}
	// Run all of the manager API tests against the converted manager and close it.
	// We'll also retrieve the default scope (BIP0044) from the manager in order to
	// use.
	scopedMgr, e := mgr.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0044)
	if e != nil {
		tc.t.Errorf("unable to fetch bip 44 scope %v", e)
		return false
	}
	testManagerAPI(
		&testContext{
			t:            tc.t,
			db:           db,
			rootManager:  mgr,
			manager:      scopedMgr,
			account:      0,
			create:       false,
			watchingOnly: true,
		},
	)
	mgr.Close()
	// Open the watching-only manager and run all the tests again.
	e = walletdb.View(
		db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			mgr, e = waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
			return e
		},
	)
	if e != nil {
		tc.t.Errorf("Open Watching-Only: unexpected error: %v", e)
		return false
	}
	defer mgr.Close()
	scopedMgr, e = mgr.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0044)
	if e != nil {
		tc.t.Errorf("unable to fetch bip 44 scope %v", e)
		return false
	}
	testManagerAPI(
		&testContext{
			t:            tc.t,
			db:           db,
			rootManager:  mgr,
			manager:      scopedMgr,
			account:      0,
			create:       false,
			watchingOnly: true,
		},
	)
	return true
}

// testSync tests various facets of setting the manager sync state.
func testSync(tc *testContext) bool {
	// Ensure syncing the manager to nil results in the synced to state being the
	// earliest block (genesis block in this case).
	e := walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.SetSyncedTo(ns, nil)
		},
	)
	if e != nil {
		tc.t.Errorf("SetSyncedTo unexpected e on nil: %v", e)
		return false
	}
	blockStamp := waddrmgr.BlockStamp{
		Height: 0,
		Hash:   *chaincfg.MainNetParams.GenesisHash,
	}
	gotBlockStamp := tc.rootManager.SyncedTo()
	if gotBlockStamp != blockStamp {
		tc.t.Errorf(
			"SyncedTo unexpected block stamp on nil -- "+
				"got %v, want %v", gotBlockStamp, blockStamp,
		)
		return false
	}
	// If we update to a new more recent block time stamp, then upon retrieval it
	// should be returned as the best known state.
	latestHash, e := chainhash.NewHash(seed)
	if e != nil {
		tc.t.Errorf("%v", e)
		return false
	}
	blockStamp = waddrmgr.BlockStamp{
		Height:    1,
		Hash:      *latestHash,
		Timestamp: time.Unix(1234, 0),
	}
	e = walletdb.Update(
		tc.db, func(tx walletdb.ReadWriteTx) (e error) {
			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
			return tc.rootManager.SetSyncedTo(ns, &blockStamp)
		},
	)
	if e != nil {
		tc.t.Errorf("SetSyncedTo unexpected e on nil: %v", e)
		return false
	}
	gotBlockStamp = tc.rootManager.SyncedTo()
	if gotBlockStamp != blockStamp {
		tc.t.Errorf(
			"SyncedTo unexpected block stamp on nil -- "+
				"got %v, want %v", gotBlockStamp, blockStamp,
		)
		return false
	}
	return true
}

// // TestManager performs a full suite of tests against the address manager API.
// // It makes use of a test context because the address manager is persistent and
// // much of the testing involves having specific state.
// func TestManager(// 	t *testing.T) {
// 	t.Parallel()
// 	teardown, db := emptyDB(t)
// 	defer teardown()
// 	// Open manager that does not exist to ensure the expected error is
// 	// returned.
// 	e := walletdb.View(db, func(tx walletdb.ReadTx) (e error) {
// 		ns := tx.ReadBucket(waddrmgrNamespaceKey)
// 		_, e := waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
// 		return e
// 	})
// 	if !checkManagerError(t, "Open non-existent", e, waddrmgr.ErrNoExist) {
// 		return
// 	}
// 	// Create a new manager.
// 	var mgr *waddrmgr.Manager
// 	e = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
// 		ns, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
// 		if e != nil  {
// 			return e
// 		}
// 		e = waddrmgr.Create(
// 			ns, seed, pubPassphrase, privPassphrase,
// 			&chaincfg.MainNetParams, fastScrypt, time.Time{},
// 		)
// 		if e != nil  {
// 			return e
// 		}
// 		mgr, e = waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
// 		return e
// 	})
// 	if e != nil  {
// 		t.Errorf("Create/Open: unexpected error: %v", e)
// 		return
// 	}
// 	// NOTE: Not using deferred close here since part of the tests is
// 	// explicitly closing the manager and then opening the existing one.
// 	// Attempt to create the manager again to ensure the expected error is
// 	// returned.
// 	e = walletdb.Update(db, func(tx walletdb.ReadWriteTx) (e error) {
// 		ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 		return waddrmgr.Create(ns, seed, pubPassphrase, privPassphrase,
// 			&chaincfg.MainNetParams, fastScrypt, time.Time{})
// 	})
// 	if !checkManagerError(t, "Create existing", e, waddrmgr.ErrAlreadyExists) {
// 		mgr.Close()
// 		return
// 	}
// 	// Run all of the manager API tests in create mode and close the
// 	// manager after they've completed
// 	scopedMgr, e := mgr.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0044)
// 	if e != nil  {
// 		t.Fatalf("unable to fetch default scope: %v", e)
// 	}
// 	testManagerAPI(&testContext{
// 		t:            t,
// 		db:           db,
// 		manager:      scopedMgr,
// 		rootManager:  mgr,
// 		account:      0,
// 		create:       true,
// 		watchingOnly: false,
// 	})
// 	mgr.Close()
// 	// Ensure the expected error is returned if the latest manager version
// 	// constant is bumped without writing code to actually do the upgrade.
// 	*waddrmgr.TstLatestMgrVersion++
// 	e = walletdb.View(db, func(tx walletdb.ReadTx) (e error) {
// 		ns := tx.ReadBucket(waddrmgrNamespaceKey)
// 		_, e := waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
// 		return e
// 	})
// 	if !checkManagerError(t, "Upgrade needed", e, waddrmgr.ErrUpgrade) {
// 		return
// 	}
// 	*waddrmgr.TstLatestMgrVersion--
// 	// Open the manager and run all the tests again in open mode which
// 	// avoids reinserting new addresses like the create mode tests do.
// 	e = walletdb.View(db, func(tx walletdb.ReadTx) (e error) {
// 		ns := tx.ReadBucket(waddrmgrNamespaceKey)
// 		var e error
// 		mgr, e = waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
// 		return e
// 	})
// 	if e != nil  {
// 		t.Errorf("Open: unexpected error: %v", e)
// 		return
// 	}
// 	defer mgr.Close()
// 	scopedMgr, e = mgr.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0044)
// 	if e != nil  {
// 		t.Fatalf("unable to fetch default scope: %v", e)
// 	}
// 	tc := &testContext{
// 		t:            t,
// 		db:           db,
// 		manager:      scopedMgr,
// 		rootManager:  mgr,
// 		account:      0,
// 		create:       false,
// 		watchingOnly: false,
// 	}
// 	testManagerAPI(tc)
// 	// Now that the address manager has been tested in both the newly
// 	// created and opened modes, test a watching-only version.
// 	testWatchingOnly(tc)
// 	// Ensure that the manager sync state functionality works as expected.
// 	testSync(tc)
// 	// Unlock the manager so it can be closed with it unlocked to ensure
// 	// it works without issue.
// 	e = walletdb.View(db, func(tx walletdb.ReadTx) (e error) {
// 		ns := tx.ReadBucket(waddrmgrNamespaceKey)
// 		return mgr.Unlock(ns, privPassphrase)
// 	})
// 	if e != nil  {
// 		t.Errorf("Unlock: unexpected error: %v", e)
// 	}
// }

// TestEncryptDecryptErrors ensures that errors which occur while encrypting and
// decrypting data return the expected errors.
func TestEncryptDecryptErrors(t *testing.T) {
	t.Parallel()
	teardown, db, mgr := setupManager(t)
	defer teardown()
	invalidKeyType := waddrmgr.CryptoKeyType(0xff)
	var e error
	if _, e = mgr.Encrypt(invalidKeyType, []byte{}); e == nil {
		t.Fatalf("Encrypt accepted an invalid key type!")
	}
	if _, e = mgr.Decrypt(invalidKeyType, []byte{}); e == nil {
		t.Fatalf("Encrypt accepted an invalid key type!")
	}
	if !mgr.IsLocked() {
		t.Fatal("Manager should be locked at this point.")
	}
	// Now the mgr is locked and encrypting/decrypting with private keys should fail.
	_, e = mgr.Encrypt(waddrmgr.CKTPrivate, []byte{})
	checkManagerError(
		t, "encryption with private key fails when manager is locked",
		e, waddrmgr.ErrLocked,
	)
	_, e = mgr.Decrypt(waddrmgr.CKTPrivate, []byte{})
	checkManagerError(
		t, "decryption with private key fails when manager is locked",
		e, waddrmgr.ErrLocked,
	)
	// Unlock the manager for these tests
	e = walletdb.View(
		db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return mgr.Unlock(ns, privPassphrase)
		},
	)
	if e != nil {
		t.Fatal("Attempted to unlock the manager, but failed:", e)
	}
	// Make sure to cover the ErrCrypto error path in Encrypt.
	waddrmgr.TstRunWithFailingCryptoKeyPriv(
		mgr, func() {
			_, e = mgr.Encrypt(waddrmgr.CKTPrivate, []byte{})
		},
	)
	checkManagerError(t, "failed encryption", e, waddrmgr.ErrCrypto)
	// Make sure to cover the ErrCrypto error path in Decrypt.
	waddrmgr.TstRunWithFailingCryptoKeyPriv(
		mgr, func() {
			_, e = mgr.Decrypt(waddrmgr.CKTPrivate, []byte{})
		},
	)
	checkManagerError(t, "failed decryption", e, waddrmgr.ErrCrypto)
}

// TestEncryptDecrypt ensures that encrypting and decrypting data with the the
// various crypto key types works as expected.
func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()
	teardown, db, mgr := setupManager(t)
	defer teardown()
	plainText := []byte("this is a plaintext")
	// Make sure address manager is unlocked
	e := walletdb.View(
		db, func(tx walletdb.ReadTx) (e error) {
			ns := tx.ReadBucket(waddrmgrNamespaceKey)
			return mgr.Unlock(ns, privPassphrase)
		},
	)
	if e != nil {
		t.Fatal("Attempted to unlock the manager, but failed:", e)
	}
	keyTypes := []waddrmgr.CryptoKeyType{
		waddrmgr.CKTPublic,
		waddrmgr.CKTPrivate,
		waddrmgr.CKTScript,
	}
	for _, keyType := range keyTypes {
		cipherText, e := mgr.Encrypt(keyType, plainText)
		if e != nil {
			t.Fatalf("Failed to encrypt plaintext: %v", e)
		}
		decryptedCipherText, e := mgr.Decrypt(keyType, cipherText)
		if e != nil {
			t.Fatalf("Failed to decrypt plaintext: %v", e)
		}
		if !reflect.DeepEqual(decryptedCipherText, plainText) {
			t.Fatal("Got:", decryptedCipherText, ", want:", plainText)
		}
	}
}

// // TestScopedKeyManagerManagement tests that callers are able to properly
// // create, retrieve, and utilize new scoped managers outside the set of default
// // created scopes.
// func TestScopedKeyManagerManagement(t *testing.T) {
// 	t.Parallel()
// 	teardown, db := emptyDB(t)
// 	defer teardown()
// 	// We'll start the test by creating a new root manager that will be used for the
// 	// duration of the test.
// 	var mgr *waddrmgr.Manager
// 	e := walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
// 			if e != nil {
// 				return e
// 			}
// 			e = waddrmgr.Create(
// 				ns, seed, pubPassphrase, privPassphrase,
// 				&chaincfg.MainNetParams, fastScrypt, time.Time{},
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			mgr, e = waddrmgr.Open(
// 				ns, pubPassphrase, &chaincfg.MainNetParams,
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			return mgr.Unlock(ns, privPassphrase)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("create/open: unexpected error: %v", e)
// 	}
// 	// All the default scopes should have been created and loaded into memory upon
// 	// initial opening.
// 	for _, scope := range waddrmgr.DefaultKeyScopes {
// 		_, e := mgr.FetchScopedKeyManager(scope)
// 		if e != nil {
// 			t.Fatalf("unable to fetch scope %v: %v", scope, e)
// 		}
// 	}
// 	// Next, ensure that if we create an internal and external addrs for each of the
// 	// default scope types, then they're derived according to their schema.
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			for _, scope := range waddrmgr.DefaultKeyScopes {
// 				sMgr, e := mgr.FetchScopedKeyManager(scope)
// 				if e != nil {
// 					t.Fatalf("unable to fetch scope %v: %v", scope, e)
// 				}
// 				externalAddr, e := sMgr.NextExternalAddresses(
// 					ns, waddrmgr.DefaultAccountNum, 1,
// 				)
// 				if e != nil {
// 					t.Fatalf("unable to derive external addr: %v", e)
// 				}
// 				// The external address should match the prescribed addr schema for this scoped
// 				// key manager.
// 				if externalAddr[0].AddrType() != waddrmgr.ScopeAddrMap[scope].ExternalAddrType {
// 					t.Fatalf(
// 						"addr type mismatch: expected %v, got %v",
// 						externalAddr[0].AddrType(),
// 						waddrmgr.ScopeAddrMap[scope].ExternalAddrType,
// 					)
// 				}
// 				internalAddr, e := sMgr.NextInternalAddresses(
// 					ns, waddrmgr.DefaultAccountNum, 1,
// 				)
// 				if e != nil {
// 					t.Fatalf("unable to derive internal addr: %v", e)
// 				}
// 				// Similarly, the internal address should match the prescribed addr schema for
// 				// this scoped key manager.
// 				if internalAddr[0].AddrType() != waddrmgr.ScopeAddrMap[scope].InternalAddrType {
// 					t.Fatalf(
// 						"addr type mismatch: expected %v, got %v",
// 						internalAddr[0].AddrType(),
// 						waddrmgr.ScopeAddrMap[scope].InternalAddrType,
// 					)
// 				}
// 			}
// 			return e
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to read db: %v", e)
// 	}
// 	// Now that the manager is open, we'll create a "test" scope that we'll be
// 	// utilizing for the remainder of the test.
// 	testScope := waddrmgr.KeyScope{
// 		Purpose: 99,
// 		Coin:    0,
// 	}
// 	addrSchema := waddrmgr.ScopeAddrSchema{
// 		ExternalAddrType: waddrmgr.NestedWitnessPubKey,
// 		InternalAddrType: waddrmgr.WitnessPubKey,
// 	}
// 	var scopedMgr *waddrmgr.ScopedKeyManager
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			scopedMgr, e = mgr.NewScopedKeyManager(ns, testScope, addrSchema)
// 			if e != nil {
// 				return e
// 			}
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to read db: %v", e)
// 	}
// 	// The manager was just created, we should be able to look it up within the root
// 	// manager.
// 	if _, e = mgr.FetchScopedKeyManager(testScope); E.Chk(e) {
// 		t.Fatalf("attempt to read created mgr failed: %v", e)
// 	}
// 	var externalAddr, internalAddr []waddrmgr.ManagedAddress
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			// We'll now create a new external address to ensure we retrieve the proper
// 			// type.
// 			externalAddr, e = scopedMgr.NextExternalAddresses(
// 				ns, waddrmgr.DefaultAccountNum, 1,
// 			)
// 			if e != nil {
// 				t.Fatalf("unable to derive external addr: %v", e)
// 			}
// 			internalAddr, e = scopedMgr.NextInternalAddresses(
// 				ns, waddrmgr.DefaultAccountNum, 1,
// 			)
// 			if e != nil {
// 				t.Fatalf("unable to derive internal addr: %v", e)
// 			}
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("open: unexpected error: %v", e)
// 	}
// 	// Ensure that the type of the address matches as expected.
// 	if externalAddr[0].AddrType() != waddrmgr.NestedWitnessPubKey {
// 		t.Fatalf(
// 			"addr type mismatch: expected %v, got %v",
// 			waddrmgr.NestedWitnessPubKey, externalAddr[0].AddrType(),
// 		)
// 	}
// 	_, ok := externalAddr[0].Address().(*util.ScriptHash)
// 	if !ok {
// 		t.Fatalf("wrong type: %T", externalAddr[0].Address())
// 	}
// 	// We'll also create an internal address and ensure that the types match up
// 	// properly.
// 	if internalAddr[0].AddrType() != waddrmgr.WitnessPubKey {
// 		t.Fatalf(
// 			"addr type mismatch: expected %v, got %v",
// 			waddrmgr.WitnessPubKey, internalAddr[0].AddrType(),
// 		)
// 	}
// 	_, ok = internalAddr[0].Address().(*util.AddressWitnessPubKeyHash)
// 	if !ok {
// 		t.Fatalf("wrong type: %T", externalAddr[0].Address())
// 	}
// 	// We'll now simulate a restart by closing, then restarting the manager.
// 	mgr.Close()
// 	e = walletdb.View(
// 		db, func(tx walletdb.ReadTx) (e error) {
// 			ns := tx.ReadBucket(waddrmgrNamespaceKey)
// 			var e error
// 			mgr, e = waddrmgr.Open(ns, pubPassphrase, &chaincfg.MainNetParams)
// 			if e != nil {
// 				return e
// 			}
// 			return mgr.Unlock(ns, privPassphrase)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("open: unexpected error: %v", e)
// 	}
// 	defer mgr.Close()
// 	// We should be able to retrieve the new scoped manager that we just created.
// 	scopedMgr, e = mgr.FetchScopedKeyManager(testScope)
// 	if e != nil {
// 		t.Fatalf("attempt to read created mgr failed: %v", e)
// 	}
// 	// If we fetch the last generated external address, it should map exactly to the
// 	// address that we just generated.
// 	var lastAddr waddrmgr.ManagedAddress
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			lastAddr, e = scopedMgr.LastExternalAddress(
// 				ns, waddrmgr.DefaultAccountNum,
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("open: unexpected error: %v", e)
// 	}
// 	if !bytes.Equal(lastAddr.AddrHash(), externalAddr[0].AddrHash()) {
// 		t.Fatalf(
// 			"mismatch addr hashes: expected %x, got %x",
// 			externalAddr[0].AddrHash(), lastAddr.AddrHash(),
// 		)
// 	}
// 	// After the restart, all the default scopes should be been re-loaded.
// 	for _, scope := range waddrmgr.DefaultKeyScopes {
// 		_, e := mgr.FetchScopedKeyManager(scope)
// 		if e != nil {
// 			t.Fatalf("unable to fetch scope %v: %v", scope, e)
// 		}
// 	}
// 	// Finally, if we attempt to query the root manager for this last address, it
// 	// should be able to locate the private key, etc.
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			_, e := mgr.Address(ns, lastAddr.Address())
// 			if e != nil {
// 				return fmt.Errorf("unable to find addr: %v", e)
// 			}
// 			e = mgr.MarkUsed(ns, lastAddr.Address())
// 			if e != nil {
// 				return fmt.Errorf(
// 					"unable to mark addr as "+
// 						"used: %v", e,
// 				)
// 			}
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to find addr: %v", e)
// 	}
// }

// // TestRootHDKeyNeutering tests that callers are unable to create new scoped
// // managers once the root HD key has been deleted from the database.
// func TestRootHDKeyNeutering(t *testing.T) {
// 	t.Parallel()
// 	teardown, db := emptyDB(t)
// 	defer teardown()
// 	// We'll start the test by creating a new root manager that will be used for the
// 	// duration of the test.
// 	var mgr *waddrmgr.Manager
// 	e := walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
// 			if e != nil {
// 				return e
// 			}
// 			e = waddrmgr.Create(
// 				ns, seed, pubPassphrase, privPassphrase,
// 				&chaincfg.MainNetParams, fastScrypt, time.Time{},
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			mgr, e = waddrmgr.Open(
// 				ns, pubPassphrase, &chaincfg.MainNetParams,
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			return mgr.Unlock(ns, privPassphrase)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("create/open: unexpected error: %v", e)
// 	}
// 	defer mgr.Close()
// 	// With the root manager open, we'll now create a new scoped manager for usage
// 	// within this test.
// 	testScope := waddrmgr.KeyScope{
// 		Purpose: 99,
// 		Coin:    0,
// 	}
// 	addrSchema := waddrmgr.ScopeAddrSchema{
// 		ExternalAddrType: waddrmgr.NestedWitnessPubKey,
// 		InternalAddrType: waddrmgr.WitnessPubKey,
// 	}
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			_, e := mgr.NewScopedKeyManager(ns, testScope, addrSchema)
// 			if e != nil {
// 				return e
// 			}
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to read db: %v", e)
// 	}
// 	// With the manager created, we'll now neuter the root HD private key.
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			return mgr.NeuterRootKey(ns)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to read db: %v", e)
// 	}
// 	// If we try to create *another* scope, this should fail, as the root key is no
// 	// longer in the database.
// 	testScope = waddrmgr.KeyScope{
// 		Purpose: 100,
// 		Coin:    0,
// 	}
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			_, e := mgr.NewScopedKeyManager(ns, testScope, addrSchema)
// 			if e != nil {
// 				return e
// 			}
// 			return nil
// 		},
// 	)
// 	if e == nil {
// 		t.Fatalf("new scoped manager creation should have failed")
// 	}
// }

// // TestNewRawAccount tests that callers are able to properly create, and use raw
// // accounts created with only an account number, and not a string which is
// // eventually mapped to an account number.
// func TestNewRawAccount(t *testing.T) {
// 	t.Parallel()
// 	teardown, db := emptyDB(t)
// 	defer teardown()
// 	// We'll start the test by creating a new root manager that will be used for the
// 	// duration of the test.
// 	var mgr *waddrmgr.Manager
// 	e := walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns, e := tx.CreateTopLevelBucket(waddrmgrNamespaceKey)
// 			if e != nil {
// 				return e
// 			}
// 			e = waddrmgr.Create(
// 				ns, seed, pubPassphrase, privPassphrase,
// 				&chaincfg.MainNetParams, fastScrypt, time.Time{},
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			mgr, e = waddrmgr.Open(
// 				ns, pubPassphrase, &chaincfg.MainNetParams,
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			return mgr.Unlock(ns, privPassphrase)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("create/open: unexpected error: %v", e)
// 	}
// 	defer mgr.Close()
// 	// Now that we have the manager created, we'll fetch one of the default scopes
// 	// for usage within this test.
// 	scopedMgr, e := mgr.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0084)
// 	if e != nil {
// 		t.Fatalf("unable to fetch scope %v: %v", waddrmgr.KeyScopeBIP0084, e)
// 	}
// 	// With the scoped manager retrieved, we'll attempt to create a new raw account
// 	// by number.
// 	const accountNum = 1000
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			return scopedMgr.NewRawAccount(ns, accountNum)
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to create new account: %v", e)
// 	}
// 	// With the account created, we should be able to derive new addresses from the
// 	// account.
// 	var accountAddrNext waddrmgr.ManagedAddress
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			addrs, e := scopedMgr.NextExternalAddresses(
// 				ns, accountNum, 1,
// 			)
// 			if e != nil {
// 				return e
// 			}
// 			accountAddrNext = addrs[0]
// 			return nil
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to create addr: %v", e)
// 	}
// 	// Additionally, we should be able to manually derive specific target keys.
// 	var accountTargetAddr waddrmgr.ManagedAddress
// 	e = walletdb.Update(
// 		db, func(tx walletdb.ReadWriteTx) (e error) {
// 			ns := tx.ReadWriteBucket(waddrmgrNamespaceKey)
// 			keyPath := waddrmgr.DerivationPath{
// 				Account: accountNum,
// 				Branch:  0,
// 				Index:   0,
// 			}
// 			accountTargetAddr, e = scopedMgr.DeriveFromKeyPath(
// 				ns, keyPath,
// 			)
// 			return e
// 		},
// 	)
// 	if e != nil {
// 		t.Fatalf("unable to derive addr: %v", e)
// 	}
// 	// The two keys we just derived should match up perfectly.
// 	if accountAddrNext.AddrType() != accountTargetAddr.AddrType() {
// 		t.Fatalf(
// 			"wrong addr type: %v vs %v",
// 			accountAddrNext.AddrType(), accountTargetAddr.AddrType(),
// 		)
// 	}
// 	if !bytes.Equal(accountAddrNext.AddrHash(), accountTargetAddr.AddrHash()) {
// 		t.Fatalf(
// 			"wrong pubkey hash: %x vs %x", accountAddrNext.AddrHash(),
// 			accountTargetAddr.AddrHash(),
// 		)
// 	}
// }
