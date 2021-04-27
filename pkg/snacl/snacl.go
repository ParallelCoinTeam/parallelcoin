package snacl

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"io"
	"runtime/debug"
	
	"github.com/btcsuite/golangcrypto/nacl/secretbox"
	"github.com/btcsuite/golangcrypto/scrypt"
	
	"github.com/p9c/matrjoska/pkg/util/zero"
)

var (
	prng = rand.Reader
	
	// ErrInvalidPassword ...
	ErrInvalidPassword = errors.New("invalid password")
	// ErrMalformed ...
	ErrMalformed = errors.New("malformed data")
	// ErrDecryptFailed ...
	ErrDecryptFailed = errors.New("unable to decrypt")
)

// Various constants needed for encryption scheme.
const (
	// Overhead const here expose secretbox's for convenience.
	Overhead  = secretbox.Overhead
	keySize   = 32
	NonceSize = 24
	DefaultN  = 16384 // 2^14
	DefaultR  = 8
	DefaultP  = 1
)

// CryptoKey represents a secret key which can be used to encrypt and decrypt data.
type CryptoKey [keySize]byte

// Encrypt encrypts the passed data.
func (ck *CryptoKey) Encrypt(in []byte) ([]byte, error) {
	var nonce [NonceSize]byte
	_, e := io.ReadFull(prng, nonce[:])
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	blob := secretbox.Seal(nil, in, &nonce, (*[keySize]byte)(ck))
	return append(nonce[:], blob...), nil
}

// Decrypt decrypts the passed data. The must be the output of the Encrypt function.
func (ck *CryptoKey) Decrypt(in []byte) ([]byte, error) {
	if len(in) < NonceSize {
		return nil, ErrMalformed
	}
	var nonce [NonceSize]byte
	copy(nonce[:], in[:NonceSize])
	blob := in[NonceSize:]
	opened, ok := secretbox.Open(nil, blob, &nonce, (*[keySize]byte)(ck))
	if !ok {
		return nil, ErrDecryptFailed
	}
	return opened, nil
}

// Zero clears the key by manually zeroing all memory. This is for security conscience application which wish to zero
// the memory after they've used it rather than waiting until it's reclaimed by the garbage collector. The key is no
// longer usable after this call.
func (ck *CryptoKey) Zero() {
	zero.Bytea32((*[keySize]byte)(ck))
}

// GenerateCryptoKey generates a new crypotgraphically random key.
func GenerateCryptoKey() (*CryptoKey, error) {
	var key CryptoKey
	_, e := io.ReadFull(prng, key[:])
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	return &key, nil
}

// Parameters are not secret and can be stored in plain text.
type Parameters struct {
	Salt   [keySize]byte
	Digest [sha256.Size]byte
	N      int
	R      int
	P      int
}

// SecretKey houses a crypto key and the parameters needed to derive it from a passphrase. It should only be used in
// memory.
type SecretKey struct {
	Key        *CryptoKey
	Parameters Parameters
}

// deriveKey fills out the Key field.
func (sk *SecretKey) deriveKey(password *[]byte) (e error) {
	key, e := scrypt.Key(*password, sk.Parameters.Salt[:],
		sk.Parameters.N,
		sk.Parameters.R,
		sk.Parameters.P,
		len(sk.Key),
	)
	if e != nil {
		E.Ln(e)
		return e
	}
	copy(sk.Key[:], key)
	zero.Bytes(key)
	// I'm not a fan of forced garbage collections, but scrypt allocates a ton of memory and calling it back to back
	// without a GC cycle in between means you end up needing twice the amount of memory.
	//
	// For example, if your scrypt parameters are such that you require 1GB and you call it twice in a row, without this
	// you end up allocating 2GB since the first GB probably hasn't been released yet.
	debug.FreeOSMemory()
	return nil
}

// Marshal returns the Parameters field marshalled into a format suitable for storage.
//
// This result of this can be stored in clear text.
func (sk *SecretKey) Marshal() []byte {
	params := &sk.Parameters
	// The marshalled format for the the netparams is as follows:
	//
	//   <salt><digest><N><R><P>
	//
	// KeySize + sha256.Size + N (8 bytes) + R (8 bytes) + P (8 bytes)
	marshalled := make([]byte, keySize+sha256.Size+24)
	b := marshalled
	copy(b[:keySize], params.Salt[:])
	b = b[keySize:]
	copy(b[:sha256.Size], params.Digest[:])
	b = b[sha256.Size:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.N))
	b = b[8:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.R))
	b = b[8:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.P))
	return marshalled
}

// Unmarshal unmarshalls the parameters needed to derive the secret key from a passphrase into sk.
func (sk *SecretKey) Unmarshal(marshalled []byte) (e error) {
	if sk.Key == nil {
		sk.Key = (*CryptoKey)(&[keySize]byte{})
	}
	// The marshalled format for the the netparams is as follows:
	//
	//   <salt><digest><N><R><P>
	//
	// KeySize + sha256.Size + N (8 bytes) + R (8 bytes) + P (8 bytes)
	if len(marshalled) != keySize+sha256.Size+24 {
		return ErrMalformed
	}
	params := &sk.Parameters
	copy(params.Salt[:], marshalled[:keySize])
	marshalled = marshalled[keySize:]
	copy(params.Digest[:], marshalled[:sha256.Size])
	marshalled = marshalled[sha256.Size:]
	params.N = int(binary.LittleEndian.Uint64(marshalled[:8]))
	marshalled = marshalled[8:]
	params.R = int(binary.LittleEndian.Uint64(marshalled[:8]))
	marshalled = marshalled[8:]
	params.P = int(binary.LittleEndian.Uint64(marshalled[:8]))
	return nil
}

// Zero zeroes the underlying secret key while leaving the parameters intact.
//
// This effectively makes the key unusable until it is derived again via the DeriveKey function.
func (sk *SecretKey) Zero() {
	sk.Key.Zero()
}

// DeriveKey derives the underlying secret key and ensures it matches the expected digest.
//
// This should only be called after previously calling the Zero function or on an initial Unmarshal.
func (sk *SecretKey) DeriveKey(password *[]byte) (e error) {
	if e = sk.deriveKey(password); E.Chk(e) {
		return
	}
	// verify password
	digest := sha256.Sum256(sk.Key[:])
	if subtle.ConstantTimeCompare(digest[:], sk.Parameters.Digest[:]) != 1 {
		return ErrInvalidPassword
	}
	return nil
}

// Encrypt encrypts in bytes and returns a JSON blob.
func (sk *SecretKey) Encrypt(in []byte) ([]byte, error) {
	return sk.Key.Encrypt(in)
}

// Decrypt takes in a JSON blob and returns it's decrypted form.
func (sk *SecretKey) Decrypt(in []byte) ([]byte, error) {
	return sk.Key.Decrypt(in)
}

// NewSecretKey returns a SecretKey structure based on the passed parameters.
func NewSecretKey(password *[]byte, N, r, p int) (sk *SecretKey, e error) {
	sk = &SecretKey{
		Key: (*CryptoKey)(&[keySize]byte{}),
	}
	// setup parameters
	sk.Parameters.N = N
	sk.Parameters.R = r
	sk.Parameters.P = p
	_, e = io.ReadFull(prng, sk.Parameters.Salt[:])
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	// derive key
	e = sk.deriveKey(password)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	// store digest
	sk.Parameters.Digest = sha256.Sum256(sk.Key[:])
	return sk, nil
}
