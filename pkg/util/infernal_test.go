package util

import (
	"golang.org/x/crypto/ripemd160"

	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/pkg/util/base58"
	"github.com/p9c/pod/pkg/util/bech32"
	ec "github.com/p9c/pod/pkg/util/elliptic"
)

// SetBlockBytes sets the internal serialized block byte buffer to the passed buffer.  It is used to inject errors and is only available to the test package.
func (b *Block) SetBlockBytes(buf []byte) {
	b.serializedBlock = buf
}

// TstAppDataDir makes the internal appDataDir function available to the test package.
// nolint
func TstAppDataDir(goos, appName string, roaming bool) string {
	return appdata.AppDataDir(goos, appName, roaming)
}

// TstAddressPubKeyHash makes an AddressPubKeyHash, setting the unexported fields with the parameters hash and netID.
// nolint
func TstAddressPubKeyHash(hash [ripemd160.Size]byte,
	netID byte) *AddressPubKeyHash {
	return &AddressPubKeyHash{
		hash:  hash,
		netID: netID,
	}
}

// TstAddressScriptHash makes an AddressScriptHash, setting the unexported fields with the parameters hash and netID.
// nolint
func TstAddressScriptHash(hash [ripemd160.Size]byte,
	netID byte) *AddressScriptHash {
	return &AddressScriptHash{
		hash:  hash,
		netID: netID,
	}
}

// TstAddressWitnessPubKeyHash creates an AddressWitnessPubKeyHash, initiating the fields as given.
// nolint
func TstAddressWitnessPubKeyHash(version byte, program [20]byte,
	hrp string) *AddressWitnessPubKeyHash {
	return &AddressWitnessPubKeyHash{
		hrp:            hrp,
		witnessVersion: version,
		witnessProgram: program,
	}
}

// TstAddressWitnessScriptHash creates an AddressWitnessScriptHash, initiating the fields as given.
// nolint
func TstAddressWitnessScriptHash(version byte, program [32]byte,
	hrp string) *AddressWitnessScriptHash {
	return &AddressWitnessScriptHash{
		hrp:            hrp,
		witnessVersion: version,
		witnessProgram: program,
	}
}

// TstAddressPubKey makes an AddressPubKey, setting the unexported fields with the parameters.
// nolint
func TstAddressPubKey(serializedPubKey []byte, pubKeyFormat PubKeyFormat,
	netID byte) *AddressPubKey {
	pubKey, _ := ec.ParsePubKey(serializedPubKey, ec.S256())
	return &AddressPubKey{
		pubKeyFormat: pubKeyFormat,
		pubKey:       (*ec.PublicKey)(pubKey),
		pubKeyHashID: netID,
	}
}

// TstAddressSAddr returns the expected script address bytes for P2PKH and P2SH bitcoin addresses.
// nolint
func TstAddressSAddr(addr string) []byte {
	decoded := base58.Decode(addr)
	return decoded[1 : 1+ripemd160.Size]
}

// TstAddressSegwitSAddr returns the expected witness program bytes for bech32 encoded P2WPKH and P2WSH bitcoin addresses.
// nolint
func TstAddressSegwitSAddr(addr string) []byte {
	_, data, err := bech32.Decode(addr)
	if err != nil {
		return []byte{}
	}
	// First byte is version, rest is base 32 encoded data.
	data, err = bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return []byte{}
	}
	return data
}
