package control

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/kopach/control/templates"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/util"
	"time"
)

// GetMsgBlockTemplate gets a Message for the current chain paying to a
// given address
func (c *Controller) GetMsgBlockTemplate(addr util.Address) (mbt *templates.Message, err error) {
	mbt = &templates.Message{
		Nonce:     c.uuid,
		PrevBlock: c.cx.RealNode.Chain.BestSnapshot().Hash,
		Height:    c.height.Load(),
		Bits:      make(map[int32]uint32),
		Merkles:   make(map[int32]chainhash.Hash),
	}
	mbt.ResetCoinbases()
	next, curr, more := fork.AlgoVerIterator(c.height.Load())
	for ; more(); next() {
		var templateX *mining.BlockTemplate
		if templateX, err = c.blockTemplateGenerator.NewBlockTemplate(
			0, addr, fork.GetAlgoName(curr(), c.height.Load()),
		); Check(err) {
		} else {
			mbt.SetCoinbase(curr(), templateX.Block.Transactions[len(templateX.Block.Transactions)-1])
			mbt.Bits[curr()] = templateX.Block.Header.Bits
			mbt.Merkles[curr()] = templateX.Block.Header.MerkleRoot
			Debugf(
				"))))))))))))))))))) %d %d %0.8f %08x %v",
				mbt.Height,
				curr(),
				util.Amount(mbt.GetCoinbase(curr()).TxOut[0].Value).ToDUO(),
				mbt.Bits[curr()],
				mbt.Merkles[curr()],
			)
			mbt.Timestamp = templateX.Block.Header.Timestamp.Add(time.Second)
			mbt.SetTxs(templateX.Block.Transactions[:len(templateX.Block.Transactions)-1])
			Debugs(mbt.GetTxs())
			Debugs(mbt.GetCoinbase(curr()))
		}
	}
	return
}

func getBlkTemplateGenerator(cx *conte.Xt) *mining.BlkTmplGenerator {
	policy := mining.Policy{
		BlockMinWeight:    uint32(*cx.Config.BlockMinWeight),
		BlockMaxWeight:    uint32(*cx.Config.BlockMaxWeight),
		BlockMinSize:      uint32(*cx.Config.BlockMinSize),
		BlockMaxSize:      uint32(*cx.Config.BlockMaxSize),
		BlockPrioritySize: uint32(*cx.Config.BlockPrioritySize),
		TxMinFreeFee:      cx.StateCfg.ActiveMinRelayTxFee,
	}
	s := cx.RealNode
	return mining.NewBlkTmplGenerator(
		&policy,
		s.ChainParams, s.TxMemPool, s.Chain, s.TimeSource,
		s.SigCache, s.HashCache,
	)
}
