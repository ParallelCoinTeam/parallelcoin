package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	"os"
	
	"github.com/p9c/pod/pkg/coding/ecc"
)

func main() {
	
	fi, e := os.Create("secp256k1.go")
	
	if e != nil  {
				ftl.Ln(e)
	}
	defer func() {
		if e := fi.Close(); err.Chk(e) {
		}
	}()
	
	// Compress the serialized byte points.
	serialized := ecc.S256().SerializedBytePoints()
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	
	if _, e = w.Write(serialized); err.Chk(e) {
				os.Exit(1)
	}
	if e := w.Close(); err.Chk(e) {
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
	_, _ = fmt.Fprintln(fi,
		"// The following values are the computed linearly "+
			"independent vectors needed to make use of the secp256k1 "+
			"endomorphism:")
	_, _ = fmt.Fprintf(fi, "// a1: %x\n", a1)
	_, _ = fmt.Fprintf(fi, "// b1: %x\n", b1)
	_, _ = fmt.Fprintf(fi, "// a2: %x\n", a2)
	_, _ = fmt.Fprintf(fi, "// b2: %x\n", b2)
}

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// var _ = logg.AddFilteredSubsystem(subsystem)
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	ftl.Ln("ftl.Ln")
	err.Ln("err.Ln")
	wrn.Ln("wrn.Ln")
	inf.Ln("inf.Ln")
	dbg.Ln("dbg.Ln")
	trc.Ln("trc.Ln")
	ftl.F("%s", "ftl.F")
	err.F("%s", "err.F")
	wrn.F("%s", "wrn.F")
	inf.F("%s", "inf.F")
	dbg.F("%s", "dbg.F")
	trc.F("%s", "trc.F")
	ftl.C(func() string { return "ftl.C" })
	err.C(func() string { return "err.C" })
	wrn.C(func() string { return "wrn.C" })
	inf.C(func() string { return "inf.C" })
	dbg.C(func() string { return "dbg.C" })
	trc.C(func() string { return "trc.C" })
	ftl.C(func() string { return "ftl.C" })
	err.Chk(errors.New("err.Chk"))
	wrn.Chk(errors.New("wrn.Chk"))
	inf.Chk(errors.New("inf.Chk"))
	dbg.Chk(errors.New("dbg.Chk"))
	trc.Chk(errors.New("trc.Chk"))
}
