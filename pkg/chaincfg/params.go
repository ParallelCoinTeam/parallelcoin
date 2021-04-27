package chaincfg

import (
	"strings"
	
	"github.com/p9c/matrjoska/pkg/chainhash"
)

// String returns the hostname of the DNS seed in human-readable form.
func (d DNSSeed) String() string {
	return d.Host
}

// Register registers the network parameters for a Bitcoin network. This may error with ErrDuplicateNet if the network
// is already registered (either due to a previous Register call, or the network being one of the default networks).
// Network parameters should be registered into this package by a main package as early as possible. Then, library
// packages may lookup networks or network parameters based on inputs and work regardless of the network being standard
// or not.
func Register(params *Params) (e error) {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	pubKeyHashAddrIDs[params.PubKeyHashAddrID] = struct{}{}
	scriptHashAddrIDs[params.ScriptHashAddrID] = struct{}{}
	hdPrivToPubKeyIDs[params.HDPrivateKeyID] = params.HDPublicKeyID[:]
	// // A valid Bech32 encoded segwit address always has as prefix the human-readable part for the given net followed by
	// // '1'.
	// bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}
	return nil
}

// mustRegister performs the same function as Register except it panics if there is an error. This should only be called
// from package init functions.
func mustRegister(params *Params) {
	if e := Register(params); E.Chk(e) {
		panic("failed to register network: " + e.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a pay-to-pubkey-hash address on any
// default or registered network. This is used when decoding an address string into a specific address type. It is up to
// the caller to check both this and IsScriptHashAddrID and decide whether an address is a pubkey hash address, script
// hash address, neither, or undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, ok := pubKeyHashAddrIDs[id]
	return ok
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a pay-to-script-hash address on any
// default or registered network. This is used when decoding an address string into a specific address type. It is up to
// the caller to check both this and IsPubKeyHashAddrID and decide whether an address is a pubkey hash address, script
// hash address, neither, or undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, ok := scriptHashAddrIDs[id]
	return ok
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit addresses on any default or registered
// network. This is used when decoding an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string) bool {
	prefix = strings.ToLower(prefix)
	_, ok := bech32SegwitPrefixes[prefix]
	return ok
}

// HDPrivateKeyToPublicKeyID accepts a private hierarchical deterministic extended key id and returns the associated
// public key id. When the provided id is not registered, the ErrUnknownHDKeyID error will be returned.
func HDPrivateKeyToPublicKeyID(id []byte) ([]byte, error) {
	if len(id) != 4 {
		return nil, ErrUnknownHDKeyID
	}
	var key [4]byte
	copy(key[:], id)
	pubBytes, ok := hdPrivToPubKeyIDs[key]
	if !ok {
		return nil, ErrUnknownHDKeyID
	}
	return pubBytes, nil
}

// newHashFromStr converts the passed big-endian hex string into a chainhash.Hash. It only differs from the one
// available in chainhash in that it panics on an error since it will only (and must only) be called with hard-coded,
// and therefore known good, hashes.
func newHashFromStr(hexStr string) *chainhash.Hash {
	hash, e := chainhash.NewHashFromStr(hexStr)
	if e != nil {
		// Ordinarily I don't like panics in library code since it can take applications down without them having a
		// chance to recover which is extremely annoying, however an exception is being made in this case because the
		// only way this can panic is if there is an error in the hard-coded hashes. Thus it will only ever potentially
		// panic on init and therefore is 100% predictable. loki: Panics are good when the condition should not happen!
		panic(e)
	}
	return hash
}
func init() {
	// Register all default networks when the package is initialized.
	mustRegister(&MainNetParams)
	mustRegister(&TestNet3Params)
	mustRegister(&RegressionTestParams)
	mustRegister(&SimNetParams)
}
