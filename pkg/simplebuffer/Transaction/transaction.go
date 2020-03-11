package Transaction

import (
	"bytes"
	"encoding/binary"
	"github.com/p9c/pod/pkg/chain/wire"
	log "github.com/p9c/logi"
)

type Transaction struct {
	Length uint32
	Bytes  []byte
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) DecodeOne(b []byte) *Transaction {
	t.Decode(b)
	return t
}

func (t *Transaction) Decode(b []byte) (out []byte) {
	if len(b) >= 4 {
		t.Length = binary.BigEndian.Uint32(b[:4])
		if len(b) >= 4+int(t.Length) {
			t.Bytes = b[4 : 4+t.Length]
			if len(b) > 4+int(t.Length) {
				out = b[4+t.Length:]
			}
		}
	}
	return
}

func (t *Transaction) Encode() (out []byte) {
	out = make([]byte, 4+len(t.Bytes))
	binary.BigEndian.PutUint32(out[:4], t.Length)
	copy(out[4:], t.Bytes)
	return
}

func (t *Transaction) Get() (txs *wire.MsgTx) {
	txs = new(wire.MsgTx)
	buffer := bytes.NewBuffer(t.Bytes)
	err := txs.Deserialize(buffer)
	if err != nil {
		log.L.Error(err)
	}
	return
}

func (t *Transaction) Put(txs *wire.MsgTx) *Transaction {
	var buffer bytes.Buffer
	err := txs.Serialize(&buffer)
	if err != nil {
		log.L.Error(err)
		return t
	}
	t.Bytes = buffer.Bytes()
	t.Length = uint32(len(t.Bytes))
	return t
}
