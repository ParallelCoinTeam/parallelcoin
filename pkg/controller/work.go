package controller

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util"
)

// GetMinerWork returns a message broadcast by a node and each field is decoded
// where possible avoiding memory allocation (slicing the data). Yes,
// this is not concurrent safe, put a mutex in to share it.
// Using the same principles as used in FlatBuffers,
// we define a message type that instead of using a reflect based encoder,
// there is a creation function,
// and a set of methods that extracts the individual requested field without
// copying memory, or deserialize their contents which will be concurrent safe
// All of the fields are in the same order that they will be serialized to
func GetMinerWork(cx *conte.Xt, blk *util.Block) (out []Serializer) {
	h := &Hash{}
	h.PutHash(blk.MsgBlock().Header.PrevBlock)
	bits := Bits{}
	bits.PutBits(blk.MsgBlock().Header.Bits)
	out = []Serializer{
		GetRouteableIPs(),
		GetPort((*cx.Config.Listeners)[0]),
		GetPort((*cx.Config.RPCListeners)[0]),
		GetPort(*cx.Config.Controller),
		h,
		&bits,
	}
	txs := blk.MsgBlock().Transactions
	for i := range txs {
		t := &Transaction{}
		t.PutTx(txs[i])
		out = append(out, t)
	}
	return
}

func GetMessageBase(cx *conte.Xt) Serializers {
	return Serializers{
		GetRouteableIPs(),
		GetPort((*cx.Config.Listeners)[0]),
		GetPort((*cx.Config.RPCListeners)[0]),
		GetPort(*cx.Config.Controller),
	}
}