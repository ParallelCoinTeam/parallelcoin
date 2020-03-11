package util

import (
	"crypto/sha256"
	log "github.com/p9c/logi"
	"hash"

	"golang.org/x/crypto/ripemd160"
)

// Calculate the hash of hasher over buf.
func calcHash(buf []byte, hasher hash.Hash) []byte {
	_, err := hasher.Write(buf)
	if err != nil {
		log.L.Error(err)
	}
	return hasher.Sum(nil)
}

// Hash160 calculates the hash ripemd160(sha256(b)).
func Hash160(buf []byte) []byte {
	return calcHash(calcHash(buf, sha256.New()), ripemd160.New())
}
