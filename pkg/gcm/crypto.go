package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	
	"golang.org/x/crypto/argon2"
	
	"github.com/p9c/pod/pkg/log"
)

// GetCipher returns a GCM cipher given a password string. Note that this cipher
// must be renewed every 4gb of encrypted data
func GetCipher(password string) (gcm cipher.AEAD, err error) {
	bytes := []byte(password)
	if c, err := aes.NewCipher(argon2.IDKey(reverse(bytes), bytes, 1, 64*1024, 4, 32)); log.Check(err) {
	} else if gcm, err = cipher.NewGCM(c); log.Check(err) {
	}
	return
}

func reverse(b []byte) []byte {
	for i := range b {
		b[i] = b[len(b)-1]
	}
	return b
}
