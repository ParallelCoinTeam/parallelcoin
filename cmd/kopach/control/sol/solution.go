package sol

import (
	"bytes"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/pod/pkg/chain/wire"
)

// Magic is the marker for packets containing a solution
var Magic = []byte{'s', 'o', 'l', 1}

type Solution struct {
	Port int32
	// *wire.MsgBlock
	Bytes []byte
}

func Get(port int32, mb *wire.MsgBlock) []byte {
	var buf []byte
	wr := bytes.NewBuffer(buf)
	var err error
	if err = mb.Serialize(wr); Check(err) {
	}
	s := Solution{Port: port, Bytes: wr.Bytes()} // MsgBlock: mb}
	return gotiny.Marshal(&s)
}

//
// type Container struct {
// 	simplebuffer.Container
// }
//
// func GetSolContainer(port uint32, b *wire.MsgBlock) *Container {
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
// func (sC *Container) GetMsgBlock() *wire.MsgBlock {
// 	// Traces(sC.Data)
// 	buff := sC.Get(1)
// 	// Traces(buff)
// 	decoded := Block.New().DecodeOne(buff)
// 	// Traces(decoded)
// 	got := decoded.Get()
// 	// Traces(got)
// 	return got
// }
//
// func (sC *Container) GetSenderPort() int32 {
// 	buff := sC.Get(0)
// 	decoded := Int32.New().DecodeOne(buff)
// 	got := decoded.Get()
// 	return got
// }
