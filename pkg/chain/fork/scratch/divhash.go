package main

import (
	"encoding/hex"
	"fmt"
	"math/big"

	log "github.com/p9c/logi"
)

func main() {
	log.L.SetLevel("trace", true)

	blockbytes, _ := hex.DecodeString(
		"05000000b415eee3e60bec16f601fa272f938344df32f230f24de0220818cb009c" +
			"f5901120a7bd6af8214ba26e3e15c9b98162fa358085fca88eb3a49fc138d0522" +
			"0ef5018a3d95df632032064e7041c")
	for i := 0; i < 1; i++ {
		blocklen := len(blockbytes)
		rfb := reverse(blockbytes[:blocklen/2])
		firsthalf := append(blockbytes, rfb...)
		fhc := make([]byte, len(firsthalf))
		copy(fhc, firsthalf)
		secondhalf := append(blockbytes, reverse(blockbytes[blocklen/2:])...)
		shc := make([]byte, len(secondhalf))
		copy(shc, secondhalf)
		bbb := make([]byte, len(blockbytes))
		copy(bbb, blockbytes)
		bl := big.NewInt(0).SetBytes(bbb)
		log.L.Debug("bl", len(fmt.Sprint(bl)), bl)
		fh := big.NewInt(0).SetBytes(fhc)
		log.L.Traces(fhc)
		log.L.Debug("fh", len(fmt.Sprint(fh)), fh)
		sh := big.NewInt(0).SetBytes(shc)
		sqfh := fh.Mul(fh, fh)
		sqsh := sh.Mul(sh, sh)
		sqsq := fh.Mul(sqfh, sqsh)
		divd := sqsq.Div(sqsq, bl)
		divdb := divd.Bytes()
		dlen := len(divdb)
		ddd := make([]byte, dlen)
		copy(ddd, reverse(divdb))
		dddB := big.NewInt(0).SetBytes(ddd)
		log.L.Traces(ddd)
		log.L.Debug(dddB)
	}
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1-i]
	}
	return out
}

func trimLeadingZeroes(b []byte) []byte {
	var i int
	for i = len(b) - 1; i >= 0; i++ {
		if b[0] != 0 {
			break
		}
	}
	return b[:i]
}
