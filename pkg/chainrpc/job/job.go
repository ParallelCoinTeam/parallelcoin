package job

//
// import (
// 	"errors"
// 	"github.com/p9c/matrjoska/pkg/chaincfg"
// 	"time"
//
// 	"github.com/niubaoshu/gotiny"
//
// 	"github.com/p9c/matrjoska/pkg/blockchain"
// 	"github.com/p9c/matrjoska/pkg/chainhash"
// 	"github.com/p9c/matrjoska/pkg/blockchain/fork"
// 	"github.com/p9c/matrjoska/pkg/wire"
// 	"github.com/p9c/matrjoska/pkg/util"
// )
//

var Magic = []byte{'j', 'o', 'b', 1}

//
// type Job struct {
// 	// IPs             map[string]struct{}
// 	// P2PListenerPort uint16
// 	// RPCListenerPort uint16
// 	ControllerNonce uint64
// 	Height          int32
// 	PrevBlockHash   *chainhash.Hash
// 	MinTimestamp    time.Time
// 	Diffs           blockchain.Diffs
// 	Merkles         blockchain.Merkles
// 	// CoinBases       map[int32]*util.Tx
// }
//
// // Get returns a message broadcast by a node and each field is decoded where
// // possible avoiding memory allocation (slicing the data). Yes, this is not
// // concurrent safe, put a mutex in to share it. Using the same principles as
// // used in FlatBuffers, we define a message type that instead of using a reflect
// // based encoder, there is a creation function, and a set of methods that
// // extracts the individual requested field without copying memory, or
// // deserialize their contents which will be concurrent safe The varying coinbase
// // payment values are in transaction 0 last output, the individual varying
// // transactions are stored separately and will be reassembled at the end
// func Get(
// 	node *blockchain.BlockChain,
// 	activeNet *chaincfg.Params,
// 	uuid uint64,
// 	mB *util.Block,
// ) (cbs *map[int32]*util.Tx, out []byte, txr []*util.Tx,) {
// 	_temp := make(map[int32]*util.Tx)
// 	cbs = &_temp
// 	bH := node.BestSnapshot().Height + 1
// 	tip := node.BestChain.Tip()
// 	bitsMap := make(blockchain.Diffs)
// 	var e error
// 	df, ok := tip.Diffs.Load().(blockchain.Diffs)
// 	if df == nil || !ok ||
// 		len(df) != len(fork.List[1].AlgoVers) {
// 		if bitsMap, e = node.CalcNextRequiredDifficultyPlan9Controller(tip); E.Chk(e) {
// 			return
// 		}
// 		tip.Diffs.Store(bitsMap)
// 	} else {
// 		bitsMap = tip.Diffs.Load().(blockchain.Diffs)
// 	}
// 	// Now we need to get the values for coinbase for each algorithm then regenerate
// 	// the merkle roots To mine this block a miner only needs the matching merkle
// 	// roots for the version number but to get them first get the values
// 	var val int64
// 	mTS := make(map[int32]*chainhash.Hash)
// 	mBtx := mB.Transactions()
// 	root := len(mBtx) - 1
// 	coinbase := mBtx[root]
// 	transactions := mB.Transactions()[:root]
// 	txr = make([]*util.Tx, len(transactions))
// 	for i, v := range transactions {
// 		txr[i] = v
// 	}
// 	nbH := bH
// 	if (activeNet.Net == wire.MainNet &&
// 		nbH == fork.List[1].ActivationHeight) ||
// 		(activeNet.Net == wire.TestNet3 &&
// 			nbH == fork.List[1].TestnetStart) {
// 		nbH++
// 	}
//
// 	for i := range bitsMap {
// 		// set value according to version and block height
// 		val = blockchain.CalcBlockSubsidy(nbH, activeNet, i)
// 		txc := coinbase.MsgTx().Copy()
// 		txc.TxOut[len(txc.TxOut)-1].Value = val
// 		txx := util.NewTx(txc.Copy())
// 		// D.S(coinbase)
// 		(*cbs)[i] = txx
// 		// D.Ln("coinbase for version", i, txx.MsgTx().TxOut[len(txx.MsgTx().TxOut)-1].Value)
// 		mTree := blockchain.BuildMerkleTreeStore(
// 			append(txr, txx), false,
// 		)
// 		// D.S(mTree)
// 		mr := mTree.GetRoot()
// 		if mr == nil {
// 			e = errors.New("got a nil merkle root")
// 			panic(e)
// 			return
// 		}
// 		mTS[i] = mr
// 	}
// 	jrb := Job{
// 		ControllerNonce: uuid,
// 		Height:          bH,
// 		PrevBlockHash:   &mB.WireBlock().Header.PrevBlock,
// 		Diffs:           bitsMap,
// 		Merkles:         mTS,
// 		MinTimestamp:    tip.Header().Timestamp.Truncate(time.Second).Add(time.Second),
// 	}
// 	// jrb.CoinBases= make(map[int32]*util.Tx)
// 	// for i := range *cbs {
// 	// 	jrb.CoinBases[i] = (*cbs)[i]
// 	// }
// 	out = gotiny.Marshal(&jrb)
// 	// D.S(jrb)
// 	// D.S(out)
// 	// var testy []byte
// 	// for i := range out {
// 	// 	testy = append(testy, out[i])
// 	// }
// 	// var jr Job
// 	// D.Ln(gotiny.Unmarshal(testy, &jr))
// 	// D.S(jr)
// 	// D.Ln("job size", len(jobber))
// 	// return Container{*msg.CreateContainer(Magic)}, txr
// 	return cbs, out, txr
// }
//
// //
// // // LoadContainer takes a message byte slice payload and loads it into a
// // // container ready to be decoded
// // func LoadContainer(b []byte) (out Container) {
// // 	out.Data = b
// // 	return
// // }
// //
// // func (j *Container) GetIPs() []*net.IP {
// // 	return IPs.New().DecodeOne(j.Get(0)).Get()
// // }
// //
// // func (j *Container) GetP2PListenersPort() uint16 {
// // 	return Uint16.New().DecodeOne(j.Get(1)).Get()
// // }
// //
// // func (j *Container) GetRPCListenersPort() uint16 {
// // 	return Uint16.New().DecodeOne(j.Get(2)).Get()
// // }
// //
// // func (j *Container) GetControllerListenerPort() uint16 {
// // 	return Uint16.New().DecodeOne(j.Get(3)).Get()
// // }
// //
// // func (j *Container) GetNewHeight() (out int32) {
// // 	return Int32.New().DecodeOne(j.Get(4)).Get()
// // }
// //
// // func (j *Container) GetPrevBlockHash() (out *chainhash.Hash) {
// // 	return Hash.New().DecodeOne(j.Get(5)).Get()
// // }
// //
// // func (j *Container) GetBitses() blockchain.Bits {
// // 	return Bits.NewBitses().DecodeOne(j.Get(6)).Get()
// // }
// //
// // // GetHashes returns the merkle roots per version
// // func (j *Container) GetHashes() (out map[int32]*chainhash.Hash) {
// // 	return Merkles.NewHashes().DecodeOne(j.Get(7)).Get()
// // }
// //
// // func (j *Container) String() (s string) {
// // 	s += fmt.Sprint("\ntype '"+string(Magic)+"' elements:", j.Count())
// // 	s += "\n"
// // 	ips := j.GetIPs()
// // 	s += "1 IPs:"
// // 	for i := range ips {
// // 		s += fmt.Sprint(" ", ips[i].String())
// // 	}
// // 	s += "\n"
// // 	s += fmt.Sprint("2 P2PListenersPort: ", j.GetP2PListenersPort())
// // 	s += "\n"
// // 	s += fmt.Sprint("3 RPCListenersPort: ", j.GetRPCListenersPort())
// // 	s += "\n"
// // 	s += fmt.Sprint(
// // 		"4 ControllerListenerPort: ",
// // 		j.GetControllerListenerPort(),
// // 	)
// // 	s += "\n"
// // 	h := j.GetNewHeight()
// // 	s += fmt.Sprint("5 Block height: ", h)
// // 	s += "\n"
// // 	s += fmt.Sprintf(
// // 		"6 Previous Block Hash (sha256d): %064x",
// // 		j.GetPrevBlockHash().CloneBytes(),
// // 	)
// // 	s += "\n"
// // 	bitses := j.GetBitses()
// // 	s += fmt.Sprint("7 Difficulty targets:\n")
// // 	var sortedBitses []int
// // 	for i := range bitses {
// // 		sortedBitses = append(sortedBitses, int(i))
// // 	}
// // 	txsort.Ints(sortedBitses)
// // 	for i := range sortedBitses {
// // 		s += fmt.Sprintf(
// // 			"  %2d %-10v %d %064x",
// // 			sortedBitses[i],
// // 			fork.List[fork.GetCurrent(h)].AlgoVers[int32(sortedBitses[i])],
// // 			bitses[int32(sortedBitses[i])],
// // 			fork.CompactToBig(bitses[int32(sortedBitses[i])]).Bytes(),
// // 		)
// // 		s += "\n"
// // 	}
// // 	s += "8 Merkles:\n"
// // 	hashes := j.GetHashes()
// // 	for i := range sortedBitses {
// // 		s += fmt.Sprintf(
// // 			"  %2d %s\n", sortedBitses[i],
// // 			hashes[int32(sortedBitses[i])].String(),
// // 		)
// // 	}
// //
// // 	// s += spew.Sdump(j.GetHashes())
// // 	return
// // }
// //
// // //
// // // // Struct returns a handy Go struct version This can be used at the start of a
// // // // new block to get a handy struct, the first work received triggers startup and
// // // // locks the worker into sending solutions there, until there is a new
// // // // PrevBlockHash, the work controller (kopach) only responds to updates from
// // // // this first one (or if it stops sending) - the controller keeps track of
// // // // individual controller servers multicasting and when it deletes a newly gone
// // // // dark controller when it comes to send if it isn't found it falls back to the
// // // // next available to submit
// // // func (j *Container) Struct() (out Job) {
// // // 	out = Job{
// // // 		IPs:             j.GetIPs(),
// // // 		P2PListenerPort: j.GetP2PListenersPort(),
// // // 		RPCListenerPort: j.GetRPCListenersPort(),
// // // 		ControllerPort:      j.GetControllerListenerPort(),
// // // 		Height:          j.GetNewHeight(),
// // // 		PrevBlockHash:   j.GetPrevBlockHash(),
// // // 		Bits:          j.GetBitses(),
// // // 		Merkles:          j.GetHashes(),
// // // 	}
// // // 	return
// // // }
//
// // GetMsgBlock takes the handy go struct version and returns a wire.WireBlock
// // ready for giving nonce extranonce and computing the merkel root based on the
// // extranonce in the coinbase as needs to be done when mining, so this would be
// // called for each round for each algorithm to start.
// func (j *Job) GetMsgBlock(version int32) (out *wire.WireBlock) {
// 	found := false
// 	for i := range j.Diffs {
// 		if i == version {
// 			found = true
// 		}
// 	}
// 	if found {
// 		tn := time.Now().Truncate(time.Second)
// 		if tn.Sub(j.MinTimestamp) < time.Second {
// 			tn = j.MinTimestamp
// 		}
// 		out = &wire.WireBlock{
// 			Header: wire.BlockHeader{
// 				Version:    version,
// 				PrevBlock:  *j.PrevBlockHash,
// 				MerkleRoot: *j.Merkles[version],
// 				Timestamp:  tn,
// 			},
// 			// Transactions: j.Txs,
// 		}
// 	}
// 	return
// }
