package sol

import (
	"bytes"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/matrjoska/pkg/wire"
)

// Magic is the marker for packets containing a solution
var Magic = []byte{'s', 'o', 'l', 1}

type Solution struct {
	Nonce uint64
	UUID  uint64
	// *wire.Block
	Bytes []byte
}

// Encode a message for a solution
func Encode(nonce uint64, uuid uint64, mb *wire.BlockHeader) []byte {
	var buf []byte
	wr := bytes.NewBuffer(buf)
	var e error
	if e = mb.Serialize(wr); E.Chk(e) {
	}
	s := Solution{Nonce: nonce, UUID: uuid, Bytes: wr.Bytes()} // Block: mb}
	return gotiny.Marshal(&s)
}

// Decode an encoded solution message to a wire.BlockHeader
func (s *Solution) Decode() (mb *wire.BlockHeader, e error) {
	buf := bytes.NewBuffer(s.Bytes)
	mb = &wire.BlockHeader{}
	if e = mb.Deserialize(buf); E.Chk(e) {
	}
	return
}

//
// type Container struct {
// 	simplebuffer.Container
// }
//
// func GetSolContainer(port uint32, b *wire.Block) *Container {
// 	mB := Block.New().Put(b)
// 	srs := simplebuffer.Serializers{Int32.New().Put(int32(port)), mB}.CreateContainer(Magic)
// 	return &Container{*srs}
// }
//
// func LoadSolContainer(b []byte) (out *Container) {
// 	out = &Container{}
// 	out.Data = b
// 	return
// }
//
//
// func (sC *Container) GetMsgBlock() *wire.Block {
// 	// Traces(sC.Data)
// 	buff := sC.Encode(1)
// 	// Traces(buff)
// 	decoded := Block.New().DecodeOne(buff)
// 	// Traces(decoded)
// 	got := decoded.Encode()
// 	// Traces(got)
// 	return got
// }

//
// func (sC *Container) GetSenderPort() int32 {
// 	buff := sC.Encode(0)
// 	decoded := Int32.New().DecodeOne(buff)
// 	got := decoded.Encode()
// 	return got
// }
