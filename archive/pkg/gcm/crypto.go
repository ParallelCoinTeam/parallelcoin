package gcm

import (
	"crypto/aes"
	"crypto/cipher"

	"golang.org/x/crypto/argon2"
)

// GetCipher returns a GCM cipher given a password string. Note that this cipher must be renewed every 4gb of encrypted
// data
func GetCipher(password []byte) (gcm cipher.AEAD, e error) {
	bytes := make([]byte, len(password))
	rb := make([]byte, len(password))
	copy(bytes, password)
	copy(rb, password)
	var c cipher.Block
	rb = reverse(bytes)
	ark := argon2.IDKey(rb, bytes, 1, 64*1024, 4, 32)
	if c, e = aes.NewCipher(ark); E.Chk(e) {
		return
	}
	if gcm, e = cipher.NewGCM(c); E.Chk(e) {
	}
	for i := range bytes {
		bytes[i] = 0
		rb[i] = 0
	}
	return
}

func reverse(b []byte) []byte {
	for i := range b {
		b[i], b[len(b)-1] = b[len(b)-1], b[i]
	}
	return b
}
