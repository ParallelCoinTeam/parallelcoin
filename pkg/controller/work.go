package controller

import (
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"net"
)

type MinerContainer struct {
	Container
}

// GetMinerContainer returns a message broadcast by a node and each field is decoded
// where possible avoiding memory allocation (slicing the data). Yes,
// this is not concurrent safe, put a mutex in to share it.
// Using the same principles as used in FlatBuffers,
// we define a message type that instead of using a reflect based encoder,
// there is a creation function,
// and a set of methods that extracts the individual requested field without
// copying memory, or deserialize their contents which will be concurrent safe
// All of the fields are in the same order that they will be serialized to
func GetMinerContainer(cx *conte.Xt, mB *util.Block,
	msg Serializers) (out MinerContainer) {
	//msg := append(Serializers{}, GetMessageBase(cx)...)
	bH := cx.RealNode.Chain.BestSnapshot().Height + 1
	nBH := NewInt32().Put(bH)
	msg = append(msg, nBH)
	mH := NewHash().Put(*mB.Hash())
	msg = append(msg, mH)
	tip := cx.RealNode.Chain.BestChain.Tip()
	//// this should be the same as the block in the notification
	//tth := tip.Header()
	//tH := &tth
	//tbh := tH.BlockHash()
	//if tbh.IsEqual(mB.Hash()) {
	//	log.DEBUG("notification block is tip block")
	//} else {
	//	log.DEBUG("notification block is not tip block")
	//}
	bM := map[int32]uint32{}
	bitsMap := &bM
	var err error
	tip.DiffMx.Lock()
	defer tip.DiffMx.Unlock()
	if tip.Diffs == nil ||
		len(*tip.Diffs) != len(fork.List[1].AlgoVers) {
		bitsMap, err = cx.RealNode.Chain.
			CalcNextRequiredDifficultyPlan9Controller(tip)
		if err != nil {
			log.ERROR(err)
			return
		}
	} else {
		bitsMap = tip.Diffs
	}
	bitses := NewBitses()
	bitses.Put(*bitsMap)
	msg = append(msg, bitses)
	txs := mB.MsgBlock().Transactions
	for i := range txs {
		t := (&Transaction{}).Put(txs[i])
		msg = append(msg, t)
	}
	return MinerContainer{*msg.CreateContainer(WorkMagic)}
}

func LoadMinerContainer(b []byte) (out MinerContainer) {
	out.Data = b
	return
}

func (mC *MinerContainer) GetIPs() []*net.IP {
	return NewIPs().DecodeOne(mC.Get(0)).Get()
}

func (mC *MinerContainer) GetP2PListenersPort() uint16 {
	return NewPort().DecodeOne(mC.Get(1)).Get()
}

func (mC *MinerContainer) GetRPCListenersPort() uint16 {
	return NewPort().DecodeOne(mC.Get(2)).Get()
}

func (mC *MinerContainer) GetControllerListenerPort() uint16 {
	return NewPort().DecodeOne(mC.Get(3)).Get()
}

func (mC *MinerContainer) GetNewHeight() (out int32) {
	return NewInt32().DecodeOne(mC.Get(4)).Get()
	return
}

func (mC *MinerContainer) GetPrevBlockHash() (out *chainhash.Hash) {
	return NewHash().DecodeOne(mC.Get(5)).Get()
}

func (mC *MinerContainer) GetBitses() map[int32]uint32 {
	return NewBitses().DecodeOne(mC.Get(6)).Get()
}

func (mC *MinerContainer) GetTxs() (out []*wire.MsgTx) {
	count := mC.Count()
	i := count
	// there has to be at least one transaction so we won't check if there is
	for i = 7; i < count; i++ {
		out = append(out, NewTransaction().DecodeOne(mC.Get(i)).Get())
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

type PauseContainer struct {
	Container
}

func LoadPauseContainer(b []byte) (out PauseContainer) {
	out.Data = b
	return
}

func (mC *PauseContainer) GetIPs() []*net.IP {
	return NewIPs().DecodeOne(mC.Get(0)).Get()
}
