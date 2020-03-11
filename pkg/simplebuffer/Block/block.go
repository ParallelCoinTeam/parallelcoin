package Block

import (
	"bytes"
	"encoding/binary"

	"github.com/p9c/pod/pkg/chain/wire"
	log "github.com/p9c/logi"
)

type Block struct {
	Length uint32
	Bytes  []byte
}

func New() *Block {
	return &Block{}
}

func (B *Block) DecodeOne(b []byte) *Block {
	B.Decode(b)
	return B
}

func (B *Block) Decode(b []byte) (out []byte) {
	// log.L.Traces(b)
	if len(b) >= 4 {
		B.Length = binary.BigEndian.Uint32(b[:4])
		//log.L.Debug("length", B.Length)
		if len(b) >= 4+int(B.Length) {
			B.Bytes = b[4 : 4+B.Length]
			if len(b) > 4+int(B.Length) {
				out = b[4+B.Length:]
			}
		}
	}
	//log.L.Traces(out)
	return
}

func (B *Block) Encode() (out []byte) {
	out = make([]byte, 4+len(B.Bytes))
	binary.BigEndian.PutUint32(out[:4], B.Length)
	copy(out[4:], B.Bytes)
	return
}

func (B *Block) Get() (b *wire.MsgBlock) {
	b = &wire.MsgBlock{}
	buffer := bytes.NewBuffer(B.Bytes)
	err := b.Deserialize(buffer)
	if err != nil {
		log.L.Error(err)
	}
	return
}

func (B *Block) Put(b *wire.MsgBlock) *Block {
	var buffer bytes.Buffer
	err := b.Serialize(&buffer)
	if err != nil {
		log.L.Error(err)
		return B
	}
	B.Bytes = buffer.Bytes()
	B.Length = uint32(len(B.Bytes))
	return B
}
