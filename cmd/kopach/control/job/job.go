package job

import (
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	
	"github.com/niubaoshu/gotiny"
	
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	
	"github.com/p9c/pod/app/conte"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/util"
)

var Magic = []byte{'j', 'o', 'b', 1}

type Job struct {
	// IPs             map[string]struct{}
	// P2PListenerPort uint16
	// RPCListenerPort uint16
	ControllerNonce uint64
	Height          int32
	PrevBlockHash   *chainhash.Hash
	Bitses          blockchain.TargetBits
	MerkleRoots     map[int32]*chainhash.Hash
	CoinBases       map[int32]*wire.MsgTx
}

// Get returns a message broadcast by a controller containing all the necessary
// data to construct blocks to mine
func Get(cx *conte.Xt, templates []*mining.BlockTemplate) (cbs *map[int32]*wire.MsgTx, out []byte, txr []*util.Tx) {
	_temp := make(map[int32]*wire.MsgTx)
	cbs = &_temp
	bH := cx.RealNode.Chain.BestSnapshot().Height + 1
	tip := cx.RealNode.Chain.BestChain.Tip()
	bitsMap := make(blockchain.TargetBits)
	var err error
	df, ok := tip.Diffs.Load().(blockchain.TargetBits)
	if df == nil || !ok ||
		len(df) != len(fork.List[1].AlgoVers) {
		if bitsMap, err = cx.RealNode.Chain.CalcNextRequiredDifficultyPlan9Controller(tip); Check(err) {
			return
		}
		tip.Diffs.Store(bitsMap)
	} else {
		bitsMap = tip.Diffs.Load().(blockchain.TargetBits)
	}
	Traces(bitsMap)
	mTS := make(map[int32]*chainhash.Hash)
	for i := range templates {
		mTS[templates[i].Block.Header.Version] = &templates[i].Block.Header.MerkleRoot
		(*cbs)[templates[i].Block.Header.Version] = templates[i].Block.Transactions[0]
	}
	for _, x := range templates[0].Block.Transactions[1:] {
		txr = append(txr, util.NewTx(x))
	}
	prevBlock := templates[0].Block.Header.PrevBlock
	adv := p2padvt.GetAdvt(cx)
	jrb := Job{
		ControllerNonce: adv.UUID,
		Height:          bH,
		PrevBlockHash:   &prevBlock,
		Bitses:          bitsMap,
		MerkleRoots:     mTS,
		// CoinBases:       *cbs,
	}
	out = gotiny.Marshal(&jrb)
	return cbs, out, txr
}
