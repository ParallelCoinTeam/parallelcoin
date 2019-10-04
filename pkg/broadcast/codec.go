package broadcast

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/p9c/pod/pkg/fec"
	"github.com/p9c/pod/pkg/log"
)

// Encode creates Reed Solomon shards and encrypts them using
// the provided GCM cipher function (from pkg/gcm).
// Message type is given in the first byte of each shard so nodes can quickly
// eliminate erroneous or irrelevant messages
func Encode(ciph cipher.AEAD, bytes []byte, typ []byte) (shards [][]byte,
	err error) {
	if len(bytes) > 1<<32 {
		log.WARN("GCM ciphers should only encode a maximum of 4gb per nonce" +
			" per key")
	}
	var clearText [][]byte
	clearText, err = fec.Encode(bytes)
	if err != nil {
		log.ERROR(err)
		return
	}
	// the nonce groups a broadcast's pieces,
	// the listener will gather them by this criteria.
	// The decoder does not enforce this but a message can be identified by
	// its' nonce due to using the same for each piece of the message
	nonce := make([]byte, ciph.NonceSize())
	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.ERROR(err)
		return
	}
	for i := range clearText {
		shards = append(shards, append(append(typ, nonce...), ciph.Seal(nil, nonce,
			clearText[i], nil)...))
	}
	log.SPEW(shards)
	return
}

func Decode(ciph cipher.AEAD, shards [][]byte) (bytes []byte, err error) {
	plainShards := make([][]byte, len(shards))
	nonceSize := ciph.NonceSize()
	for i := range shards {
		if len(shards[i]) < nonceSize {
			errMsg := []interface{}{"shard size too small, got",
				len(shards[i]), "expected minimum", nonceSize}
			log.ERROR(errMsg...)
			return nil, errors.New(fmt.Sprintln(errMsg...))
		}
		nonce, cipherText := shards[i][:nonceSize], shards[i][nonceSize:]
		var plaintext []byte
		plaintext, err = ciph.Open(nil, nonce, cipherText, nil)
		if err != nil {
			log.ERROR(err)
			return
		}
		plainShards[i] = plaintext
	}
	return fec.Decode(plainShards)
}
