// +build ignore

package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	ec "github.com/parallelcointeam/parallelcoin/pkg/util/elliptic"
)

func main() {

	fi, err := os.Create("secp256k1.go")

	if err != nil {

		log.Fatal(err)
	}
	defer fi.Close()

	// Compress the serialized byte points.
	serialized := ec.S256().SerializedBytePoints()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)

	if _, err := w.Write(serialized); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
	w.Close()

	// Encode the compressed byte points with base64.
	encoded := make([]byte, base64.StdEncoding.EncodedLen(compressed.Len()))
	base64.StdEncoding.Encode(encoded, compressed.Bytes())
	fmt.Fprintln(fi, "")
	fmt.Fprintln(fi, "")
	fmt.Fprintln(fi, "")
	fmt.Fprintln(fi)
	fmt.Fprintln(fi, "package ec")
	fmt.Fprintln(fi)
	fmt.Fprintln(fi, "// Auto-generated file (see genprecomps.go)")
	fmt.Fprintln(fi, "// DO NOT EDIT")
	fmt.Fprintln(fi)
	fmt.Fprintf(fi, "var secp256k1BytePoints = %q\n", string(encoded))
	a1, b1, a2, b2 := ec.S256().EndomorphismVectors()
	fmt.Println("The following values are the computed linearly " +
		"independent vectors needed to make use of the secp256k1 " +
		"endomorphism:")
	fmt.Printf("a1: %x\n", a1)
	fmt.Printf("b1: %x\n", b1)
	fmt.Printf("a2: %x\n", a2)
	fmt.Printf("b2: %x\n", b2)
}
