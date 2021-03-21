package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	"os"
	
	"github.com/p9c/pod/pkg/coding/ecc"
)

func main() {
	
	fi, e := os.Create("secp256k1.go")
	
	if e != nil {
		F.Ln(e)
	}
	defer func() {
		if e := fi.Close(); E.Chk(e) {
		}
	}()
	
	// todo this needs fixing lol
	
	// Compress the serialized byte points.
	serialized := ecc.S256().SerializedBytePoints()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	
	if _, e = w.Write(serialized); E.Chk(e) {
		os.Exit(1)
	}
	if e := w.Close(); E.Chk(e) {
	}
	
	// Encode the compressed byte points with base64.
	encoded := make([]byte, base64.StdEncoding.EncodedLen(compressed.Len()))
	base64.StdEncoding.Encode(encoded, compressed.Bytes())
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi, "")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintln(fi, "package ecc")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintln(fi, "// Auto-generated file (see genprecomps.go)")
	_, _ = fmt.Fprintln(fi, "// DO NOT EDIT")
	_, _ = fmt.Fprintln(fi)
	_, _ = fmt.Fprintf(fi, "var secp256k1BytePoints = %q\n", string(encoded))
	a1, b1, a2, b2 := ecc.S256().EndomorphismVectors()
	_, _ = fmt.Fprintln(
		fi,
		"// The following values are the computed linearly "+
			"independent vectors needed to make use of the secp256k1 "+
			"endomorphism:",
	)
	_, _ = fmt.Fprintf(fi, "// a1: %x\n", a1)
	_, _ = fmt.Fprintf(fi, "// b1: %x\n", b1)
	_, _ = fmt.Fprintf(fi, "// a2: %x\n", a2)
	_, _ = fmt.Fprintf(fi, "// b2: %x\n", b2)
}

var subsystem = logg.AddLoggerSubsystem()
var F, E, W, I, D, T logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// // var _ = logg.AddFilteredSubsystem(subsystem)
	// // var _ = logg.AddHighlightedSubsystem(subsystem)
	// F.Ln("F.Ln")
	// E.Ln("E.Ln")
	// W.Ln("W.Ln")
	// I.Ln("inf.Ln")
	// D.Ln("D.Ln")
	// F.Ln("T.Ln")
	// F.F("%s", "F.F")
	// E.F("%s", "E.F")
	// W.F("%s", "W.F")
	// I.F("%s", "I.F")
	// D.F("%s", "D.F")
	// T.F("%s", "T.F")
	// ftl.C(func() string { return "ftl.C" })
	// err.C(func() string { return "err.C" })
	// W.C(func() string { return "W.C" })
	// I.C(func() string { return "inf.C" })
	// D.C(func() string { return "D.C" })
	// T.C(func() string { return "T.C" })
	// ftl.C(func() string { return "ftl.C" })
	// E.Chk(errors.New("E.Chk"))
	// W.Chk(errors.New("W.Chk"))
	// I.Chk(errors.New("inf.Chk"))
	// D.Chk(errors.New("D.Chk"))
	// T.Chk(errors.New("T.Chk"))
}
