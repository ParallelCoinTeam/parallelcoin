package sol

import (
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/Block"
)


// SolutionMagic is the marker for packets containing a solution
var SolutionMagic = []byte{'s', 'o', 'l', 'v'}

type SolContainer struct {
	simplebuffer.Container
}

func GetSolContainer(b *wire.MsgBlock) *SolContainer {
	mB := Block.New().Put(b)
	srs := simplebuffer.Serializers{mB}.CreateContainer(SolutionMagic)
	return &SolContainer{*srs}
}

func LoadSolContainer(b []byte) (out *SolContainer) {
	out = &SolContainer{}
	out.Data = b
	return
}

func (sC *SolContainer) GetMsgBlock() *wire.MsgBlock {
	// log.SPEW(sC.Data)
	buff := sC.Get(0)
	// log.SPEW(buff)
	decoded := Block.New().DecodeOne(buff)
	// log.SPEW(decoded)
	got := decoded.Get()
	// log.SPEW(got)
	return got
}
