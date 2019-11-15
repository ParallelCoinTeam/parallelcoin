package Hash

import (
	"encoding/hex"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"testing"
)

func TestHash(t *testing.T) {
	by, err := hex.DecodeString(
		"00c44981699c4b621fe89b32148a64fc11fb680fa484ab1abe0e6fba4fcca462")
	var bhash chainhash.Hash
	err = bhash.SetBytes(by)
	if err != nil {
		panic(err)
	}
	h := New()
	h.Put(bhash)
	h2 := New()
	h2.Decode(h.Encode())
	if !h.Get().IsEqual(h2.Get()) {
		t.Fail()
	}
}
