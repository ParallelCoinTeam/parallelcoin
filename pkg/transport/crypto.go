package transport

import (
	"crypto/cipher"
	"crypto/rand"
	"io"
	
	"github.com/p9c/pod/pkg/log"
)

func decryptMessage(ciph cipher.AEAD, data []byte) (msg []byte, err error) {
	nonceSize := ciph.NonceSize()
	msg, err = ciph.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	return
}

func encryptMessage(magic []byte, ciph cipher.AEAD, data []byte) (msg []byte, err error) {
	if ciph != nil {
		var nonce []byte
		nonce, err = getNonce(ciph)
		msg = append(append(magic, nonce...), ciph.Seal(nil, nonce, data, nil)...)
		return
	} else {
		return append(magic, data...), err
	}
}

func getNonce(ciph cipher.AEAD) (nonce []byte, err error) {
	// get a nonce for the packet, it is both message ID and salt
	nonce = make([]byte, ciph.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); log.Check(err) {
	}
	return
}
