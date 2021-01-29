package sol

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/pkg/chain/wire"
)

// SolutionMagic is the marker for packets containing a solution
var SolutionMagic = []byte{'s', 'o', 'l', 'v'}

type Solution struct {
	Port int32
	*wire.MsgBlock
}

func Get(port int32, mb *wire.MsgBlock) []byte {
	s := Solution{Port: port, MsgBlock: mb}
	return gotiny.Marshal(&s)
}
//
// type Container struct {
// 	simplebuffer.Container
// }
//
// func GetSolContainer(port uint32, b *wire.MsgBlock) *Container {
// 	mB := Block.New().Put(b)
// 	srs := simplebuffer.Serializers{Int32.New().Put(int32(port)), mB}.CreateContainer(SolutionMagic)
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
