// Package zero
// Copyright (c) 2015 The btcsuite developers
package zero

// Bytea32 clears the 32-byte array by filling it with the zero value. This is used to explicitly clear private key
// material from memory.
func Bytea32(b *[32]byte) {
	*b = [32]byte{}
}

// Bytea64 clears the 64-byte array by filling it with the zero value. This is used to explicitly clear sensitive
// material from memory.
func Bytea64(b *[64]byte) {
	*b = [64]byte{}
}
