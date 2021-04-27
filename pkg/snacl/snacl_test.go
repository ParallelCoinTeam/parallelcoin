package snacl

import (
	"bytes"
	"testing"
)

var (
	password = []byte("sikrit")
	message  = []byte("this is a secret message of sorts")
	key      *SecretKey
	params   []byte
	blob     []byte
)

func TestNewSecretKey(t *testing.T) {
	var e error
	key, e = NewSecretKey(&password, DefaultN, DefaultR, DefaultP)
	if e != nil {
		return
	}
}
func TestMarshalSecretKey(t *testing.T) {
	params = key.Marshal()
}
func TestUnmarshalSecretKey(t *testing.T) {
	var sk SecretKey
	if e := sk.Unmarshal(params); E.Chk(e) {
		t.Errorf("unexpected unmarshal error: %v", e)
		return
	}
	if e := sk.DeriveKey(&password); E.Chk(e) {
		t.Errorf("unexpected DeriveKey error: %v", e)
		return
	}
	if !bytes.Equal(sk.Key[:], key.Key[:]) {
		t.Errorf("keys not equal")
	}
}
func TestUnmarshalSecretKeyInvalid(t *testing.T) {
	var sk SecretKey
	if e := sk.Unmarshal(params); E.Chk(e) {
		t.Errorf("unexpected unmarshal error: %v", e)
		return
	}
	p := []byte("wrong password")
	if e := sk.DeriveKey(&p); e != ErrInvalidPassword {
		t.Errorf("wrong password didn't fail")
		return
	}
}
func TestEncrypt(t *testing.T) {
	var e error
	blob, e = key.Encrypt(message)
	if e != nil {
		return
	}
}
func TestDecrypt(t *testing.T) {
	decryptedMessage, e := key.Decrypt(blob)
	if e != nil {
		return
	}
	if !bytes.Equal(decryptedMessage, message) {
		t.Errorf("decryption failed")
		return
	}
}
func TestDecryptCorrupt(t *testing.T) {
	blob[len(blob)-15] = blob[len(blob)-15] + 1
	_, e := key.Decrypt(blob)
	if e == nil {
		t.Errorf("corrupt message decrypted")
		return
	}
}
func TestZero(t *testing.T) {
	var zeroKey [32]byte
	key.Zero()
	if !bytes.Equal(key.Key[:], zeroKey[:]) {
		t.Errorf("zero key failed")
	}
}
func TestDeriveKey(t *testing.T) {
	if e := key.DeriveKey(&password); E.Chk(e) {
		t.Errorf("unexpected DeriveKey key failure: %v", e)
	}
}
func TestDeriveKeyInvalid(t *testing.T) {
	bogusPass := []byte("bogus")
	if e := key.DeriveKey(&bogusPass); e != ErrInvalidPassword {
		t.Errorf("unexpected DeriveKey key failure: %v", e)
	}
}
