package ecc_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	
	"github.com/p9c/pod/pkg/chainhash"
	ec "github.com/p9c/pod/pkg/ecc"
)

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// var _ = logg.AddFilteredSubsystem(subsystem)
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	F.Ln("F.Ln")
	E.Ln("E.Ln")
	W.Ln("W.Ln")
	I.Ln("I.Ln")
	D.Ln("D.Ln")
	F.Ln("T.Ln")
	F.F("%s", "F.F")
	E.F("%s", "E.F")
	W.F("%s", "W.F")
	I.F("%s", "I.F")
	D.F("%s", "D.F")
	T.F("%s", "T.F")
	ftl.C(func() string { return "ftl.C" })
	err.C(func() string { return "err.C" })
	W.C(func() string { return "W.C" })
	I.C(func() string { return "inf.C" })
	D.C(func() string { return "D.C" })
	T.C(func() string { return "T.C" })
	ftl.C(func() string { return "ftl.C" })
	E.Chk(errors.New("E.Chk"))
	W.Chk(errors.New("W.Chk"))
	I.Chk(errors.New("inf.Chk"))
	D.Chk(errors.New("D.Chk"))
	T.Chk(errors.New("T.Chk"))
}

// This example demonstrates decrypting a message using a private key that is first parsed from raw bytes.
func Example_decryptMessage() {
	// Decode the hex-encoded private key.
	pkBytes, e := hex.DecodeString(
		"a11b0a4e1a132305652ee7a8eb7848f6ad" +
			"5ea381e3ce20a2c086a2e388230811",
	)
	if e != nil {
		E.Ln(e)
		return
	}
	privKey, _ := ec.PrivKeyFromBytes(ec.S256(), pkBytes)
	ciphertext, e := hex.DecodeString(
		"35f644fbfb208bc71e57684c3c8b437402ca" +
			"002047a2f1b38aa1a8f1d5121778378414f708fe13ebf7b4a7bb74407288c1958969" +
			"00207cf4ac6057406e40f79961c973309a892732ae7a74ee96cd89823913b8b8d650" +
			"a44166dc61ea1c419d47077b748a9c06b8d57af72deb2819d98a9d503efc59fc8307" +
			"d14174f8b83354fac3ff56075162",
	)
	if e != nil {
		E.Ln(e)
	}
	// Try decrypting the message.
	plaintext, e := ec.Decrypt(privKey, ciphertext)
	if e != nil {
		E.Ln(e)
		return
	}
	fmt.Println(string(plaintext))
	// Output:
	// test message
}

// This example demonstrates encrypting a message for a public key that is first
// parsed from raw bytes, then decrypting it using the corresponding private key.
func Example_encryptMessage() {
	// Decode the hex-encoded pubkey of the recipient.
	pubKeyBytes, e := hex.DecodeString(
		"04115c42e757b2efb7671c578530ec191a1" +
			"359381e6a71127a9d37c486fd30dae57e76dc58f693bd7e7010358ce6b165e483a29" +
			"21010db67ac11b1b51b651953d2",
	) // uncompressed pubkey
	if e != nil {
		return
	}
	pubKey, e := ec.ParsePubKey(pubKeyBytes, ec.S256())
	if e != nil {
		E.Ln(e)
		return
	}
	// Encrypt a message decryptable by the private key corresponding to pubKey
	message := "test message"
	ciphertext, e := ec.Encrypt(pubKey, []byte(message))
	if e != nil {
		E.Ln(e)
		return
	}
	// Decode the hex-encoded private key.
	pkBytes, e := hex.DecodeString(
		"a11b0a4e1a132305652ee7a8eb7848f6ad" +
			"5ea381e3ce20a2c086a2e388230811",
	)
	if e != nil {
		E.Ln(e)
		return
	}
	// note that we already have corresponding pubKey
	privKey, _ := ec.PrivKeyFromBytes(ec.S256(), pkBytes)
	// Try decrypting and verify if it's the same message.
	plaintext, e := ec.Decrypt(privKey, ciphertext)
	if e != nil {
		E.Ln(e)
		return
	}
	fmt.Println(string(plaintext))
	// Output:
	// test message
}

// This example demonstrates signing a message with a secp256k1 private key that
// is first parsed form raw bytes and serializing the generated signature.
func Example_signMessage() {
	// Decode a hex-encoded private key.
	pkBytes, e := hex.DecodeString(
		"22a47fa09a223f2aa079edf85a7c2d4f87" +
			"20ee63e502ee2869afab7de234b80c",
	)
	if e != nil {
		E.Ln(e)
		return
	}
	privKey, pubKey := ec.PrivKeyFromBytes(ec.S256(), pkBytes)
	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, e := privKey.Sign(messageHash)
	if e != nil {
		E.Ln(e)
		return
	}
	// Serialize and display the signature.
	fmt.Printf("Serialized Signature: %x\n", signature.Serialize())
	// Verify the signature for the message using the public key.
	verified := signature.Verify(messageHash, pubKey)
	fmt.Printf("Signature Verified? %v\n", verified)
	// Output:
	// Serialized Signature: 304402201008e236fa8cd0f25df4482dddbb622e8a8b26ef0ba731719458de3ccd93805b022032f8ebe514ba5f672466eba334639282616bb3c2f0ab09998037513d1f9e3d6d
	// Signature Verified? true
}

// This example demonstrates verifying a secp256k1 signature against a public
// key that is first parsed from raw bytes.  The signature is also parsed from
// raw bytes.
func Example_verifySignature() {
	// Decode hex-encoded serialized public key.
	pubKeyBytes, e := hex.DecodeString(
		"02a673638cb9587cb68ea08dbef685c" +
			"6f2d2a751a8b3c6f2a7e9a4999e6e4bfaf5",
	)
	if e != nil {
		E.Ln(e)
		return
	}
	pubKey, e := ec.ParsePubKey(pubKeyBytes, ec.S256())
	if e != nil {
		E.Ln(e)
		return
	}
	// Decode hex-encoded serialized signature.
	sigBytes, e := hex.DecodeString(
		"30450220090ebfb3690a0ff115bb1b38b" +
			"8b323a667b7653454f1bccb06d4bbdca42c2079022100ec95778b51e707" +
			"1cb1205f8bde9af6592fc978b0452dafe599481c46d6b2e479",
	)
	if e != nil {
		E.Ln(e)
		return
	}
	signature, e := ec.ParseSignature(sigBytes, ec.S256())
	if e != nil {
		E.Ln(e)
		return
	}
	// Verify the signature for the message using the public key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	verified := signature.Verify(messageHash, pubKey)
	fmt.Println("Signature Verified?", verified)
	// Output:
	// Signature Verified? true
}
