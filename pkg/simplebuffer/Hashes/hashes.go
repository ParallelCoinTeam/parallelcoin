package Hashes

import (
	"encoding/binary"
	"sync"
	
	chainhash "github.com/p9c/pod/pkg/chain/hash"
)

type Hashes struct {
	sync.Mutex
	Length  byte
	Byteses map[int32][]byte
}

func NewHashes() *Hashes {
	return &Hashes{Byteses: make(map[int32][]byte)}
}

func (b *Hashes) DecodeOne(by []byte) *Hashes {
	b.Decode(by)
	return b
}

func (b *Hashes) Decode(by []byte) (out []byte) {
	b.Lock()
	defer b.Unlock()
	// log.L.Traces(by)
	if len(by) >= 7 {
		nB := by[0]
		if len(by) >= int(nB)*8 {
			for i := 0; i < int(nB); i++ {
				algoVer := int32(binary.BigEndian.Uint32(by[1+i*36 : 1+i*36+4]))
				// log.L.Debug("algoVer", algoVer, by[1+i*8+4:1+i*8+8], b.Byteses)
				b.Byteses[algoVer] = by[1+i*36+4 : 1+i*36+36]
			}
		}
		bL := int(nB)*36 + 1
		if len(by) > bL {
			out = by[bL:]
		}
	}
	// log.L.Traces(b.Byteses)
	return
}

func (b *Hashes) Encode() (out []byte) {
	b.Lock()
	defer b.Unlock()
	out = []byte{b.Length}
	for algoVer := range b.Byteses {
		by := make([]byte, 4)
		binary.BigEndian.PutUint32(by, uint32(algoVer))
		out = append(out, append(by, b.Byteses[algoVer]...)...)
	}
	// log.L.Traces(out)
	return
}

func (b *Hashes) Get() (out map[int32]*chainhash.Hash) {
	b.Lock()
	defer b.Unlock()
	out = make(map[int32]*chainhash.Hash)
	for algoVer := range b.Byteses {
		oB := b.Byteses[algoVer][:32]
		cH := chainhash.Hash{}
		err := cH.SetBytes(oB)
		// if one fails all others after it will too, this should be prevented by the HMAC
		if err == nil {
			out[algoVer] = &cH
		}
	}
	return
}

func (b *Hashes) Put(in map[int32]*chainhash.Hash) *Hashes {
	b.Lock()
	defer b.Unlock()
	// log.L.Traces(in)
	b.Length = byte(len(in))
	b.Byteses = make(map[int32][]byte, b.Length)
	for algoVer := range in {
		b.Byteses[algoVer] = in[algoVer].CloneBytes()
	}
	// log.L.Traces(b.Byteses)
	return b
}
