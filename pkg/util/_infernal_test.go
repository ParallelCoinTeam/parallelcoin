package util

import (
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"golang.org/x/crypto/ripemd160"
	
	"github.com/p9c/matrjoska/pkg/appdata"
	"github.com/p9c/matrjoska/pkg/base58"
	"github.com/p9c/matrjoska/pkg/bech32"
	ec "github.com/p9c/matrjoska/pkg/ecc"
)

// TstAppDataDir makes the internal appDataDir function available to the test package.
func TstAppDataDir(goos, appName string, roaming bool) string {
	return appdata.GetDataDir(goos, appName, roaming)
}

// TstAddressPubKeyHash makes an PubKeyHash, setting the unexported fields with the parameters hash and netID.
func TstAddressPubKeyHash(hash [ripemd160.Size]byte,
	netID byte,
) *btcaddr.PubKeyHash {
	return &btcaddr.PubKeyHash{
		Hash:  hash,
		NetID: netID,
	}
}

// TstAddressScriptHash makes an ScriptHash, setting the unexported fields with the parameters hash and netID.
func TstAddressScriptHash(hash [ripemd160.Size]byte,
	netID byte,
) *btcaddr.ScriptHash {
	return &btcaddr.ScriptHash{
		Hash:  hash,
		NetID: netID,
	}
}

// // TstAddressWitnessPubKeyHash creates an AddressWitnessPubKeyHash, initiating
// // the fields as given.
// func TstAddressWitnessPubKeyHash(version byte, program [20]byte,
// 	hrp string) *AddressWitnessPubKeyHash {
// 	return &AddressWitnessPubKeyHash{
// 		hrp:            hrp,
// 		witnessVersion: version,
// 		witnessProgram: program,
// 	}
// }
//
// // TstAddressWitnessScriptHash creates an AddressWitnessScriptHash, initiating
// // the fields as given.
// func TstAddressWitnessScriptHash(version byte, program [32]byte,
// 	hrp string) *AddressWitnessScriptHash {
// 	return &AddressWitnessScriptHash{
// 		hrp:            hrp,
// 		witnessVersion: version,
// 		witnessProgram: program,
// 	}
// }

// TstAddressPubKey makes an PubKey, setting the unexported fields with the parameters.
func TstAddressPubKey(serializedPubKey []byte, pubKeyFormat btcaddr.PubKeyFormat,
	netID byte,
) *btcaddr.PubKey {
	pubKey, _ := ec.ParsePubKey(serializedPubKey, ec.S256())
	return &btcaddr.PubKey{
		PubKeyFormat: pubKeyFormat,
		pubKey:       pubKey,
		pubKeyHashID: netID,
	}
}

// TstAddressSAddr returns the expected script address bytes for P2PKH and P2SH bitcoin addresses.
func TstAddressSAddr(addr string) []byte {
	decoded := base58.Decode(addr)
	return decoded[1 : 1+ripemd160.Size]
}

// TstAddressSegwitSAddr returns the expected witness program bytes for bech32
// encoded P2WPKH and P2WSH bitcoin addresses.
func TstAddressSegwitSAddr(addr string) []byte {
	_, data, e := bech32.Decode(addr)
	if e != nil {
		return []byte{}
	}
	// First byte is version, rest is base 32 encoded data.
	data, e = bech32.ConvertBits(data[1:], 5, 8, false)
	if e != nil {
		return []byte{}
	}
	return data
}
