package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"os"
	
	ec "github.com/p9c/pod/pkg/coding/elliptic"
)

func main() {
	
	fi, err := os.Create("secp256k1.go")
	
	if err != nil {
		Error(err)
		Fatal(err)
	}
	defer func() {
		if err := fi.Close(); Check(err) {
		}
	}()
	
	// Compress the serialized byte points.
	serialized := ec.S256().SerializedBytePoints()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	
	if _, err := w.Write(serialized); err != nil {
		Error(err)
		os.Exit(1)
	}
	w.Close()
	
	// Encode the compressed byte points with base64.
	encoded := make([]byte, base64.StdEncoding.EncodedLen(compressed.Len()))
	base64.StdEncoding.Encode(encoded, compressed.Bytes())
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintln(fi, "package ec")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintln(fi, "// Auto-generated file (see genprecomps.go)")
	_, _ = fmt.Fprintln(fi, "// DO NOT EDIT")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintf(fi, "var secp256k1BytePoints = %q\n", string(encoded))
	a1, b1, a2, b2 := ec.S256().EndomorphismVectors()
	_, _ = fmt.Fprintln(fi,
		"// The following values are the computed linearly "+
			"independent vectors needed to make use of the secp256k1 "+
			"endomorphism:")
	_, _ = fmt.Fprintf(fi, "// a1: %x\n", a1)
	_, _ = fmt.Fprintf(fi, "// b1: %x\n", b1)
	_, _ = fmt.Fprintf(fi, "// a2: %x\n", a2)
	_, _ = fmt.Fprintf(fi, "// b2: %x\n", b2)
}
