package transport

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func DecryptMessage(creator string, ciph cipher.AEAD, data []byte) (msg []byte, err error) {
	nonceSize := ciph.NonceSize()
	msg, err = ciph.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s %s", creator, err.Error()))
	} else {
		Debug("decrypted message", hex.EncodeToString(data[:nonceSize]))
	}
	return
}

// EncryptMessage encrypts a message, if the nonce is given it uses that otherwise it generates a new one. If there is
// no cipher this just returns a message with the given magic prepended.
func EncryptMessage(creator string, ciph cipher.AEAD, magic []byte, nonce, data []byte) (msg []byte, err error) {
	if ciph != nil {
		if nonce == nil {
			nonce, err = GetNonce(ciph)
		}
		msg = append(append(magic, nonce...), ciph.Seal(nil, nonce, data, nil)...)
	} else {
		msg = append(magic, data...)
	}

	return
}

func GetNonce(ciph cipher.AEAD) (nonce []byte, err error) {
	// get a nonce for the packet, it is both message ID and salt
	nonce = make([]byte, ciph.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); Check(err) {
	}
	return
}
