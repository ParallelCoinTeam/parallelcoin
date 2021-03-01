package control

import (
	"errors"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/util"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)

// MsgBlockTemplate describes a message that a mining worker can use to
// construct a block to mine on. Two methods are exported that allow a
// controlling node to
type MsgBlockTemplate struct {
	Height    int32
	PrevBlock chainhash.Hash
	Diffs     map[int32]uint32
	Merkles   map[int32]chainhash.Hash
	coinbases map[int32]*wire.MsgTx
	txs       []*wire.MsgTx
	Timestamp time.Time
}

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (c *Controller) GetNewAddressFromWallet() (addr util.Address, err error) {
	if c.walletClient != nil {
		if !c.walletClient.Disconnected() {
			Debug("have access to a wallet, generating address")
			if addr, err = c.walletClient.GetNewAddress("default"); Check(err) {
			} else {
				Debug("-------- found address", addr)
			}
		}
	} else {
		err = errors.New("no wallet available for new address")
		Debug(err)
	}
	return
}

// GetNewAddressFromMiningAddrs tries to get an address from the mining
// addresses list in the configuration file
func (c *Controller) GetNewAddressFromMiningAddrs() (addr util.Address, err error) {
	if c.cx.Config.MiningAddrs == nil {
		err = errors.New("mining addresses is nil")
		Debug(err)
		return
	}
	if len(*c.cx.Config.MiningAddrs) < 1 {
		err = errors.New("no mining addresses")
		Debug(err)
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*c.cx.Config.MiningAddrs))
	addr = c.cx.StateCfg.ActiveMiningAddrs[p2a]
	// remove the address from the state
	if p2a == 0 {
		c.cx.StateCfg.ActiveMiningAddrs = c.cx.StateCfg.ActiveMiningAddrs[1:]
	} else {
		c.cx.StateCfg.ActiveMiningAddrs = append(
			c.cx.StateCfg.ActiveMiningAddrs[:p2a],
			c.cx.StateCfg.ActiveMiningAddrs[p2a+1:]...,
		)
	}
	// update the config
	var ma cli.StringSlice
	for i := range c.cx.StateCfg.ActiveMiningAddrs {
		ma = append(ma, c.cx.StateCfg.ActiveMiningAddrs[i].String())
	}
	*c.cx.Config.MiningAddrs = ma
	save.Pod(c.cx.Config)
	return
}

// GetMsgBlockTemplate gets a MsgBlockTemplate for the current chain paying to a
// given address
func (c *Controller) GetMsgBlockTemplate(addr util.Address) (mbt *MsgBlockTemplate, err error) {
	mbt = &MsgBlockTemplate{
		PrevBlock: c.cx.RealNode.Chain.BestSnapshot().Hash,
		Height:    c.height.Load(),
		Diffs:     make(map[int32]uint32),
		Merkles:   make(map[int32]chainhash.Hash),
		coinbases: make(map[int32]*wire.MsgTx),
	}
	next, curr, more := fork.AlgoVerIterator(c.height.Load())
	for ; more(); next() {
		var templateX *mining.BlockTemplate
		if templateX, err = c.blockTemplateGenerator.NewBlockTemplate(
			0, addr, fork.GetAlgoName(
				curr(), c.height.Load(),
			),
		); Check(err) {
		} else {
			mbt.coinbases[curr()] = templateX.Block.Transactions[len(templateX.Block.Transactions)-1]
			mbt.Diffs[curr()] = templateX.Block.Header.Bits
			mbt.Merkles[curr()] = templateX.Block.Header.MerkleRoot
			Debugf(
				"))))))))))))))))))) %d %d %0.8f %08x %v",
				mbt.Height,
				curr(),
				util.Amount(mbt.coinbases[curr()].TxOut[0].Value).ToDUO(),
				mbt.Diffs[curr()],
				mbt.Merkles[curr()],
			)
			mbt.Timestamp = templateX.Block.Header.Timestamp.Add(time.Second)
			mbt.txs = templateX.Block.Transactions[:len(templateX.Block.Transactions)-1]
			Debugs(mbt.txs)
			Debugs(mbt.coinbases[curr()])
		}
	}
	return
}

// Serialize the MsgBlockTemplate for the wire
func (m *MsgBlockTemplate) Serialize() []byte {
	return gotiny.Marshal(m)
}

// DeserializeMsgBlockTemplate takes a message expected to be a MsgBlockTemplate
// and reconstitutes it
func DeserializeMsgBlockTemplate(b []byte) (m *MsgBlockTemplate) {
	m = &MsgBlockTemplate{}
	gotiny.Unmarshal(b, m)
	return
}

// GenBlockHeader generate a block given a version number to use for mining
// (nonce is empty, date can be updated, version changes merkle and target bits.
// All the data required for this is in the exported fields that are serialized
// for over the wire
func (m *MsgBlockTemplate) GenBlockHeader(vers int32) *wire.BlockHeader {
	return &wire.BlockHeader{
		Version:    vers,
		PrevBlock:  m.PrevBlock,
		MerkleRoot: m.Merkles[vers],
		Timestamp:  m.Timestamp,
		Bits:       m.Diffs[vers],
	}
}

// Reconstruct takes a received block from the wire and reattaches the transactions
func (m *MsgBlockTemplate) Reconstruct(hdr *wire.BlockHeader) *wire.MsgBlock {
	if hdr.PrevBlock != m.PrevBlock {
		Error("block is not for same parent block")
		return nil
	}
	msgBlock := &wire.MsgBlock{Header: *hdr}
	// the coinbase is the last transaction
	txs := append(m.txs, m.coinbases[msgBlock.Header.Version])
	for _, tx := range txs {
		if err := msgBlock.AddTransaction(tx); err != nil {
			return nil
		}
	}
	return msgBlock
}

