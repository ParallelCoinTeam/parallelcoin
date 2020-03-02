package sol

import (
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/Block"
	"github.com/p9c/pod/pkg/simplebuffer/Int32"
)


// SolutionMagic is the marker for packets containing a solution
var SolutionMagic = []byte{'s', 'o', 'l', 'v'}

type SolContainer struct {
	simplebuffer.Container
}

func GetSolContainer(port uint32, b *wire.MsgBlock) *SolContainer {
	mB := Block.New().Put(b)
	srs := simplebuffer.Serializers{Int32.New().Put(int32(port)), mB}.CreateContainer(SolutionMagic)
	return &SolContainer{*srs}
}

func LoadSolContainer(b []byte) (out *SolContainer) {
	out = &SolContainer{}
	out.Data = b
	return
}

func (sC *SolContainer) GetMsgBlock() *wire.MsgBlock {
	// log.SPEW(sC.Data)
	buff := sC.Get(1)
	// log.SPEW(buff)
	decoded := Block.New().DecodeOne(buff)
	// log.SPEW(decoded)
	got := decoded.Get()
	// log.SPEW(got)
	return got
}

func (sC *SolContainer) GetSenderPort() int32 {
	buff := sC.Get(0)
	decoded := Int32.New().DecodeOne(buff)
	got := decoded.Get()
	return got
}