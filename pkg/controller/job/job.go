package job

import (
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"

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
	//log.SPEW(msg)
	return Job{*msg.CreateContainer(WorkMagic)}
}

func LoadMinerContainer(b []byte) (out Job) {
	out.Data = b
	return
}

func (j *Job) GetIPs() []*net.IP {
	return IPs.New().DecodeOne(j.Get(0)).Get()
}

func (j *Job) GetP2PListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(1)).Get()
}

func (j *Job) GetRPCListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(2)).Get()
}

func (j *Job) GetControllerListenerPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(3)).Get()
}

func (j *Job) GetNewHeight() (out int32) {
	return Int32.New().DecodeOne(j.Get(4)).Get()
}

func (j *Job) GetPrevBlockHash() (out *chainhash.Hash) {
	return Hash.New().DecodeOne(j.Get(5)).Get()
}

func (j *Job) GetBitses() map[int32]uint32 {
	return Bitses.NewBitses().DecodeOne(j.Get(6)).Get()
}

func (j *Job) GetTxs() (out []*wire.MsgTx) {
	count := j.Count()
	i := count
	// there has to be at least one transaction so we won't check if there is
	for i = 7; i < count; i++ {
		out = append(out, Transaction.NewTransaction().DecodeOne(j.Get(i)).Get())
	}
	return
}

func (j *Job) String() (s string) {
	s += fmt.Sprint("type '"+string(WorkMagic)+"' elements:", j.Count())
	s += "\n"
	ips := j.GetIPs()
	s += "1 IPs:"
	for i := range ips {
		s += fmt.Sprint(" ", ips[i].String())
	}
	s += "\n"
	s += fmt.Sprint("2 P2PListenersPort: ", j.GetP2PListenersPort())
	s += "\n"
	s += fmt.Sprint("3 RPCListenersPort: ", j.GetRPCListenersPort())
	s += "\n"
	s += fmt.Sprint("4 ControllerListenerPort: ",
		j.GetControllerListenerPort())
	s += "\n"
	h := j.GetNewHeight()
	s += fmt.Sprint("5 Block height: ", h)
	s += "\n"
	s += fmt.Sprint("6 Previous Block Hash (sha256d): ", j.GetPrevBlockHash())
	s += "\n"
	bitses := j.GetBitses()
	s += fmt.Sprint("7 Difficulty targets:\n")
	for i := range bitses {
		s += fmt.Sprintf("  version: %3d %-10v %64x", i,
			fork.List[fork.GetCurrent(h)].
			AlgoVers[i],  fork.CompactToBig(bitses[i]).Bytes())
		s += "\n"
	}
	s+="8 Transactions:\n"
	s += spew.Sdump(j.GetTxs())
	return
}
