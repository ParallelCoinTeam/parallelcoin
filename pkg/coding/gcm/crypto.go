package gcm

import (
	"crypto/aes"
	"crypto/cipher"

	"golang.org/x/crypto/argon2"
)

// GetCipher returns a GCM cipher given a password string. Note that this cipher must be renewed every 4gb of encrypted
// data
func GetCipher(password string) (gcm cipher.AEAD, e error) {
	bytes := []byte(password)
	var c cipher.Block
	if c, e = aes.NewCipher(argon2.IDKey(reverse(bytes), bytes, 1, 64*1024, 4, 32)); err.Chk(e) {
	}
	if gcm, e = cipher.NewGCM(c); err.Chk(e) {
	}
	return
}

func reverse(b []byte) []byte {
	for i := range b {
		b[i] = b[len(b)-1]
	}
	return b
}
