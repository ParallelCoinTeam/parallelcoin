package base58_test

import (
	"bytes"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/base58"
)

func BenchmarkBase58Encode(b *testing.B) {
	b.StopTimer()
	data := bytes.Repeat([]byte{0xff}, 5000)
	b.SetBytes(int64(len(data)))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		base58.Encode(data)
	}
}
func BenchmarkBase58Decode(b *testing.B) {
	b.StopTimer()
	data := bytes.Repeat([]byte{0xff}, 5000)
	encoded := base58.Encode(data)
	b.SetBytes(int64(len(encoded)))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		base58.Decode(encoded)
	}
}
