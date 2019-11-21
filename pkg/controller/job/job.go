package job

import (
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/Bitses"
	"github.com/p9c/pod/pkg/simplebuffer/Hash"
	"github.com/p9c/pod/pkg/simplebuffer/IPs"
	"github.com/p9c/pod/pkg/simplebuffer/Int32"
	"github.com/p9c/pod/pkg/simplebuffer/Transaction"
	"github.com/p9c/pod/pkg/simplebuffer/Uint16"
	"github.com/p9c/pod/pkg/util"
	"net"
)

var WorkMagic = []byte{'w', 'o', 'r', 'k'}

type Job struct {
	simplebuffer.Container
}

// Get returns a message broadcast by a node and each field is decoded
// where possible avoiding memory allocation (slicing the data). Yes,
// this is not concurrent safe, put a mutex in to share it.
// Using the same principles as used in FlatBuffers,
// we define a message type that instead of using a reflect based encoder,
// there is a creation function,
// and a set of methods that extracts the individual requested field without
// copying memory, or deserialize their contents which will be concurrent safe
// All of the fields are in the same order that they will be serialized to
func Get(cx *conte.Xt, mB *util.Block,
	msg simplebuffer.Serializers) (out Job) {
	//msg := append(Serializers{}, GetMessageBase(cx)...)
	bH := cx.RealNode.Chain.BestSnapshot().Height + 1
	nBH := Int32.New().Put(bH)
	msg = append(msg, nBH)
	mH := Hash.New().Put(*mB.Hash())
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
	if tip.Diffs == nil ||
		len(*tip.Diffs) != len(fork.List[1].AlgoVers) {
		bitsMap, err = cx.RealNode.Chain.
			CalcNextRequiredDifficultyPlan9Controller(tip)
		if err != nil {
			log.ERROR(err)
			return
		}
		tip.DiffMx.Lock()
		tip.Diffs = bitsMap
		tip.DiffMx.Unlock()
	} else {
		bitsMap = tip.Diffs
	}
	bitses := Bitses.NewBitses()
	bitses.Put(*bitsMap)
	msg = append(msg, bitses)
	txs := mB.MsgBlock().Transactions
	for i := range txs {
		t := (&Transaction.Transaction{}).Put(txs[i])
		msg = append(msg, t)
	}
	return Job{*msg.CreateContainer(WorkMagic)}
}

func LoadMinerContainer(b []byte) (out Job) {
	out.Data = b
	return
}

func (mC *Job) GetIPs() []*net.IP {
	return IPs.New().DecodeOne(mC.Get(0)).Get()
}

func (mC *Job) GetP2PListenersPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(1)).Get()
}

func (mC *Job) GetRPCListenersPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(2)).Get()
}

func (mC *Job) GetControllerListenerPort() uint16 {
	return Uint16.New().DecodeOne(mC.Get(3)).Get()
}

func (mC *Job) GetNewHeight() (out int32) {
	return Int32.New().DecodeOne(mC.Get(4)).Get()
}

func (mC *Job) GetPrevBlockHash() (out *chainhash.Hash) {
	return Hash.New().DecodeOne(mC.Get(5)).Get()
}

func (mC *Job) GetBitses() map[int32]uint32 {
	return Bitses.NewBitses().DecodeOne(mC.Get(6)).Get()
}

func (mC *Job) GetTxs() (out []*wire.MsgTx) {
	count := mC.Count()
	i := count
	// there has to be at least one transaction so we won't check if there is
	for i = 7; i < count; i++ {
		out = append(out, Transaction.NewTransaction().DecodeOne(mC.Get(i)).Get())
	}
	return
}
