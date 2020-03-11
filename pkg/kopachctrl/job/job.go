package job

import (
	"fmt"
	"net"
	"sort"
	"time"

	log "github.com/p9c/logi"
	"github.com/p9c/simplebuffer"
	"github.com/p9c/simplebuffer/Bitses"
	"github.com/p9c/simplebuffer/Hash"
	"github.com/p9c/simplebuffer/Hashes"
	"github.com/p9c/simplebuffer/IPs"
	"github.com/p9c/simplebuffer/Int32"
	"github.com/p9c/simplebuffer/Uint16"

	"github.com/p9c/chainhash"
	"github.com/p9c/fork"
	"github.com/p9c/util"
	"github.com/p9c/wire"

	blockchain "github.com/p9c/chain"
	"github.com/p9c/pod/pkg/conte"
)

var Magic = []byte{'w', 'o', 'r', 'k'}

type Container struct {
	simplebuffer.Container
}

type Job struct {
	IPs             []*net.IP
	P2PListenerPort uint16
	RPCListenerPort uint16
	SubmitPort      uint16
	Height          int32
	PrevBlockHash   *chainhash.Hash
	Bitses          blockchain.TargetBits
	Hashes          map[int32]*chainhash.Hash
	CoinBases       map[int32]*util.Tx
}

// Get returns a message broadcast by a node and each field is decoded
// where possible avoiding memory allocation (slicing the data). Yes,
// this is not concurrent safe, put a mutex in to share it.
// Using the same principles as used in FlatBuffers,
// we define a message type that instead of using a reflect based encoder,
// there is a creation function,
// and a set of methods that extracts the individual requested field without
// copying memory, or deserialize their contents which will be concurrent safe
// The varying coinbase payment values are in transaction 0 last output,
// the individual varying transactions are stored separately and will be
// reassembled at the end
func Get(cx *conte.Xt, mB *util.Block, msg simplebuffer.Serializers, cbs *map[int32]*util.Tx) (out Container, txr []*util.Tx) {
	// msg := append(Serializers{}, GetMessageBase(cx)...)
	if txr == nil {
		txr = []*util.Tx{}
	}
	bH := cx.RealNode.Chain.BestSnapshot().Height + 1
	nBH := Int32.New().Put(bH)
	msg = append(msg, nBH)
	mH := Hash.New().Put(mB.MsgBlock().Header.PrevBlock)
	msg = append(msg, mH)
	tip := cx.RealNode.Chain.BestChain.Tip()
	// // this should be the same as the block in the notification
	// tth := tip.Header()
	// tH := &tth
	// tbh := tH.BlockHash()
	// if tbh.IsEqual(mB.Hash()) {
	//	log.L.Debug("notification block is tip block")
	// } else {
	//	log.L.Debug("notification block is not tip block")
	// }
	bitsMap := make(blockchain.TargetBits)
	var err error
	df := tip.Diffs.Load().(blockchain.TargetBits)
	if df == nil ||
		len(df) != len(fork.List[1].AlgoVers) {
		bitsMap, err = cx.RealNode.Chain.
			CalcNextRequiredDifficultyPlan9Controller(tip)
		if err != nil {
			log.L.Error(err)
			return
		}
		tip.Diffs.Store(bitsMap)
	} else {
		bitsMap = tip.Diffs.Load().(blockchain.TargetBits)
	}
	// log.L.Traces(*bitsMap)
	bitses := Bitses.NewBitses()
	bitses.Put(bitsMap)
	msg = append(msg, bitses)
	// Now we need to get the values for coinbase for each algorithm
	// then regenerate the merkle roots
	// To mine this block a miner only needs the matching merkle roots for the version number
	// but to get them first get the values
	var val int64
	mTS := make(map[int32]*chainhash.Hash)
	txs := mB.Transactions()[0]
	for i, v := range mB.Transactions()[1:] {
		txr[i] = v
	}
	nbH := bH
	if (cx.ActiveNet.Net == wire.MainNet &&
		nbH == fork.List[1].ActivationHeight) ||
		(cx.ActiveNet.Net == wire.TestNet3 &&
			nbH == fork.List[1].TestnetStart) {
		nbH++
	}
	for i := range bitsMap {
		val = blockchain.CalcBlockSubsidy(nbH, cx.ActiveNet, i)
		txc := txs.MsgTx().Copy()
		txc.TxOut[len(txc.TxOut)-1].Value = val
		txx := util.NewTx(txc.Copy())
		// log.L.Traces(txs)
		(*cbs)[i] = txx
		// log.L.Trace("coinbase for version", i, txx.MsgTx().TxOut[len(txx.MsgTx().TxOut)-1].Value)
		mTree := blockchain.BuildMerkleTreeStore(
			append([]*util.Tx{txx}, txr...), false)
		// log.L.Traces(mTree[len(mTree)-1].CloneBytes())
		mTS[i] = &chainhash.Hash{}
		mTS[i].SetBytes(mTree[len(mTree)-1].CloneBytes())
	}
	// log.L.Traces(mTS)
	mHashes := Hashes.NewHashes()
	mHashes.Put(mTS)
	msg = append(msg, mHashes)
	// previously were sending blocks, no need for that really miner only needs
	// valid block headers
	// txs := mB.MsgBlock().Transactions
	// for i := range txs {
	// 	t := (&Transaction.Transaction{}).Put(txs[i])
	// 	msg = append(msg, t)
	// }
	// log.L.Traces(msg)
	return Container{*msg.CreateContainer(Magic)}, txr
}

// LoadContainer takes a message byte slice payload and loads it into a container
// ready to be decoded
func LoadContainer(b []byte) (out Container) {
	out.Data = b
	return
}

func (j *Container) GetIPs() []*net.IP {
	return IPs.New().DecodeOne(j.Get(0)).Get()
}

func (j *Container) GetP2PListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(1)).Get()
}

func (j *Container) GetRPCListenersPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(2)).Get()
}

func (j *Container) GetControllerListenerPort() uint16 {
	return Uint16.New().DecodeOne(j.Get(3)).Get()
}

func (j *Container) GetNewHeight() (out int32) {
	return Int32.New().DecodeOne(j.Get(4)).Get()
}

func (j *Container) GetPrevBlockHash() (out *chainhash.Hash) {
	return Hash.New().DecodeOne(j.Get(5)).Get()
}

func (j *Container) GetBitses() blockchain.TargetBits {
	return Bitses.NewBitses().DecodeOne(j.Get(6)).Get()
}

// GetHashes returns the merkle roots per version
func (j *Container) GetHashes() (out map[int32]*chainhash.Hash) {
	return Hashes.NewHashes().DecodeOne(j.Get(7)).Get()
}

func (j *Container) String() (s string) {
	s += fmt.Sprint("\ntype '"+string(Magic)+"' elements:", j.Count())
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
	s += fmt.Sprintf("6 Previous Block Hash (sha256d): %064x",
		j.GetPrevBlockHash().CloneBytes())
	s += "\n"
	bitses := j.GetBitses()
	s += fmt.Sprint("7 Difficulty targets:\n")
	var sortedBitses []int
	for i := range bitses {
		sortedBitses = append(sortedBitses, int(i))
	}
	sort.Ints(sortedBitses)
	for i := range sortedBitses {
		s += fmt.Sprintf("  %2d %-10v %d %064x", sortedBitses[i],
			fork.List[fork.GetCurrent(h)].
				AlgoVers[int32(sortedBitses[i])],
			bitses[int32(sortedBitses[i])],
			fork.CompactToBig(bitses[int32(sortedBitses[i])]).Bytes())
		s += "\n"
	}
	s += "8 Merkles:\n"
	hashes := j.GetHashes()
	for i := range sortedBitses {
		s += fmt.Sprintf("  %2d %s\n", sortedBitses[i],
			hashes[int32(sortedBitses[i])].String())
	}

	// s += spew.Sdump(j.GetHashes())
	return
}

// Struct returns a handy Go struct version
// This can be used at the start of a new block to get a handy struct,
// the first work received triggers startup and locks the worker into sending
// solutions there, until there is a new PrevBlockHash,
// the work controller (kopach) only responds to updates from this first one
// (or if it stops sending) - the controller keeps track of individual
// controller servers multicasting and when it deletes a newly gone dark
// controller when it comes to send if it isn't found it falls back to the
// next available to submit
func (j *Container) Struct() (out Job) {
	out = Job{
		IPs:             j.GetIPs(),
		P2PListenerPort: j.GetP2PListenersPort(),
		RPCListenerPort: j.GetRPCListenersPort(),
		SubmitPort:      j.GetControllerListenerPort(),
		Height:          j.GetNewHeight(),
		PrevBlockHash:   j.GetPrevBlockHash(),
		Bitses:          j.GetBitses(),
		Hashes:          j.GetHashes(),
	}
	return
}

// GetMsgBlock takes the handy go struct version and returns a wire.MsgBlock
// ready for giving nonce extranonce and computing the merkel root based on
// the extranonce in the coinbase as needs to be done when mining,
// so this would be called for each round for each algorithm to start.
func (j *Job) GetMsgBlock(version int32) (out *wire.MsgBlock) {
	found := false
	for i := range j.Bitses {
		if i == version {
			found = true
		}
	}
	if found {
		out = &wire.MsgBlock{
			Header: wire.BlockHeader{
				Version:    version,
				PrevBlock:  *j.PrevBlockHash,
				MerkleRoot: *j.Hashes[version],
				Timestamp:  time.Now(),
			},
			// Transactions: j.Txs,
		}
	}
	return
}
