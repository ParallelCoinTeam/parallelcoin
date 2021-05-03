package transport

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// DecryptMessage attempts to decode the received message
func DecryptMessage(creator string, ciph cipher.AEAD, data []byte) (msg []byte, e error) {
	nonceSize := ciph.NonceSize()
	msg, e = ciph.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if e != nil {
		e = errors.New(fmt.Sprintf("%s %s", creator, e.Error()))
	} else {
		D.Ln("decrypted message", hex.EncodeToString(data[:nonceSize]))
	}
	return
}

// EncryptMessage encrypts a message, if the nonce is given it uses that
// otherwise it generates a new one. If there is no cipher this just returns a
// message with the given magic prepended.
func EncryptMessage(creator string, ciph cipher.AEAD, magic []byte, nonce, data []byte) (msg []byte, e error) {
	if ciph != nil {
		if nonce == nil {
			nonce, e = GetNonce(ciph)
		}
		msg = append(append(magic, nonce...), ciph.Seal(nil, nonce, data, nil)...)
	} else {
		msg = append(magic, data...)
	}
	return
}

// GetNonce reads from a cryptographicallly secure random number source
func GetNonce(ciph cipher.AEAD) (nonce []byte, e error) {
	// get a nonce for the packet, it is both message ID and salt
	nonce = make([]byte, ciph.NonceSize())
	if _, e = io.ReadFull(rand.Reader, nonce); E.Chk(e) {
	}
	return
}
