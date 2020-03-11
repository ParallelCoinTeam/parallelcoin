package Hash

import (
	log "github.com/p9c/logi"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

type Hash struct {
	Hash *chainhash.Hash
}

func New() *Hash {
	return &Hash{Hash: new(chainhash.Hash)}
}

func (h *Hash) DecodeOne(b []byte) *Hash {
	h.Decode(b)
	return h
}

func (h *Hash) Decode(b []byte) (out []byte) {
	if len(b) >= 32 {
		err := h.Hash.SetBytes(b[:32])
		if err != nil {
			log.L.Error(err)
			return
		}
		if len(b) > 32 {
			out = b[32:]
		}
	}
	return
}

func (h *Hash) Encode() []byte {
	return h.Hash.CloneBytes()
}

func (h *Hash) Get() *chainhash.Hash {
	return h.Hash
}

func (h *Hash) Put(pH chainhash.Hash) *Hash {
	// this should avoid a copy
	h.Hash = &pH
	return h
}
