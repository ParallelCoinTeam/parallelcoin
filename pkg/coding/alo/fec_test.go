package fec_test

import (
	"crypto/rand"
	"testing"

	"github.com/stalker-loki/alo/pkg/coding/fec"
)

func MakeRandomBytes(size int) (p []byte) {
	p = make([]byte, size)
	_, _ = rand.Read(p)
	return
}

// func TestSegmentBytes(t *testing.T) {
// 	for dataLen := 256; dataLen < 65536; dataLen += 16 {
// 		b := MakeRandomBytes(dataLen)
// 		for size := 32; size < 65536; size *= 2 {
// 			s := fec.SegmentBytes(b, size)
// 			// slog.Debug(dataLen, size, spew.Sdump(s))
// 			if len(s) != fec.Pieces(dataLen, size) {
// 				t.Fatal(
// 					dataLen, size, len(s), "segments were not correctly split",
// 				)
// 			}
// 		}
// 	}
// }

func TestGetShards(t *testing.T) {
	// size := 1024
	for dataLen := 16388; dataLen < 65536; dataLen += 512 {
		// for red := 1; red <= 50; red += 1 {
		red := 300
		b := MakeRandomBytes(dataLen)
		// for size := 32; size < 65536; size *= 2 {
		// s := fec.SegmentBytes(b, size)
		// slog.Debug(dataLen, size, spew.Sdump(s))
		// if len(s) != fec.Pieces(dataLen, size) {
		// 	t.Fatal(
		// 		dataLen, size, len(s), "segments were not correctly split",
		// 	)
		// }
		// slog.Debug(spew.Sdump())
		_ = fec.GetShards(b, red)
		// t.Log(size, red, len(shards))
		// }
	}
	// }
}
