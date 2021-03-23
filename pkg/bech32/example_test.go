package bech32_test

import (
	"encoding/hex"
	"fmt"
	
	"github.com/p9c/pod/pkg/bech32"
)

// This example demonstrates how to decode a bech32 encoded string.
func ExampleDecode() {
	encoded := "bc1pw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7k7grplx"
	hrp, decoded, e := bech32.Decode(encoded)
	if e != nil  {
		E.Ln("Error:", e)
	}
	// Show the decoded data.
	fmt.Println("Decoded human-readable part:", hrp)
	fmt.Println("Decoded Data:", hex.EncodeToString(decoded))
	// Output:
	// Decoded human-readable part: bc
	// Decoded Data: 010e140f070d1a001912060b0d081504140311021d030c1d03040f1814060e1e160e140f070d1a001912060b0d081504140311021d030c1d03040f1814060e1e16
}

// This example demonstrates how to encode data into a bech32 string.
func ExampleEncode() {
	data := []byte("Test data")
	// Convert test data to base32:
	conv, e := bech32.ConvertBits(data, 8, 5, true)
	if e != nil  {
		E.Ln("Error:", e)
	}
	encoded, e := bech32.Encode("customHrp!11111q", conv)
	if e != nil  {
		E.Ln("Error:", e)
	}
	// Show the encoded data.
	fmt.Println("Encoded Data:", encoded)
	// Output:
	// Encoded Data: customHrp!11111q123jhxapqv3shgcgumastr
}
