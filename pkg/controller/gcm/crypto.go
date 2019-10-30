package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"

	"github.com/btcsuite/golangcrypto/pbkdf2"

	"github.com/p9c/pod/pkg/log"
)

// GetCipher returns a GCM cipher given a password string. Note that this cipher
// must be renewed every 4gb of encrypted data
func GetCipher(password string) cipher.AEAD {
	key := pbkdf2.Key(reverse([]byte(password)), []byte(password),
		4096, 32,
		sha1.New)
	c, err := aes.NewCipher(key)
	if err != nil {
		log.ERROR(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.ERROR(err)
	}
	return gcm
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1]
	}
	return out
}
