package control

import (
	"fmt"
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/control/hashrate"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/kopach/control/pause"
	"github.com/p9c/pod/cmd/kopach/control/sol"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/comm/transport"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	"github.com/urfave/cli"
	"net"
	"time"
)

var handlersMulticast = transport.Handlers{
	string(sol.Magic):      processSolMsg,
	string(p2padvt.Magic):  processAdvtMsg,
	string(hashrate.Magic): processHashrateMsg,
}

func processAdvtMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	Debug("processing advertisment message", src, dst)
	c := ctx.(*Controller)
	var j p2padvt.Advertisment
	gotiny.Unmarshal(b, &j)
	Trace(j.IPs)
	uuid := j.UUID
	if _, ok := c.otherNodes[uuid]; !ok {
		// if we haven't already added it to the permanent peer list, we can add it now
		Debug("uuid", j.UUID, "P2P", j.P2P)
		Info("connecting to lan peer with same PSK", j.IPs, j.UUID)
		// try all IPs
		if *c.cx.Config.AutoListen {
			c.cx.Config.P2PConnect = &cli.StringSlice{}
		}
		for addr := range j.IPs {
			peerIP := net.JoinHostPort(addr, fmt.Sprint(j.P2P))
			_, addresses := routeable.GetAllInterfacesAndAddresses()
			for i := range addresses {
				addrS := net.JoinHostPort(addresses[i].IP.String(), fmt.Sprint(j.P2P))
				if addrS == peerIP {
					Debug("not connecting to self")
					continue
				}
				if *c.cx.Config.AutoListen {
					*c.cx.Config.P2PConnect = append(*c.cx.Config.P2PConnect, addrS)
				}
			}
			if *c.cx.Config.AutoListen {
				save.Pod(c.cx.Config)
			}
			if err = c.cx.RPCServer.Cfg.ConnMgr.Connect(
				peerIP,
				false,
			); Check(err) {
				continue
			}
			Debug("connected to peer via address", peerIP)
			c.otherNodes[uuid] = &nodeSpec{}
			c.otherNodes[uuid].addr = addr
			break
		}
	}
	// update last seen time for uuid for garbage collection of stale disconnected
	// nodes
	c.otherNodes[uuid].Time = time.Now()
	// If we lose connection for more than 9 seconds we delete and if the node
	// reappears it can be reconnected
	for i := range c.otherNodes {
		if time.Now().Sub(c.otherNodes[i].Time) > time.Second*9 {
			// also remove from connection manager
			if err = c.cx.RPCServer.Cfg.ConnMgr.RemoveByAddr(c.otherNodes[i].addr); Check(err) {
			}
			Debug("deleting", c.otherNodes[i])
			delete(c.otherNodes, i)
		}
	}
	on := int32(len(c.otherNodes))
	Trace("other nodes", on)
	c.cx.OtherNodes.Store(on)
	return
}

// Solutions submitted by workers
func processSolMsg(ctx interface{}, src net.Addr, dst string, b []byte,) (err error) {
	Debug("received solution", src, dst)
	c := ctx.(*Controller)
	if !c.active.Load() { // || !c.cx.Node.Load() {
		Debug("not active yet")
		return
	}
	var s sol.Solution
	gotiny.Unmarshal(b, &s)
	uuid := s.UUID
	if uuid != c.uuid {
		Debug("solution not from current controller")
		return
	}
	var newHeader *wire.BlockHeader
	if newHeader, err = s.Decode(); Check(err) {
		return
	}
	var msgBlock *wire.MsgBlock
	if msgBlock, err = c.msgBlockTemplate.Reconstruct(newHeader); Check(err) {
		return
	}
	// msgBlock := wire.NewMsgBlock(newHeader)
	Debug("-------------------------------------------------------")
	Debugs(msgBlock)
	if !msgBlock.Header.PrevBlock.IsEqual(&c.cx.RPCServer.Cfg.Chain.BestSnapshot().Hash) {
		Debug("block submitted by kopach miner worker is stale")
		if err := c.updateAndSendWork(); Check(err) {
		}
		return
	}
	// Warn(msgBlock.Header.Version)
	// cb, ok := c.coinbases.Load().(map[int32]*util.Tx)[msgBlock.Header.Version]
	cbRaw := c.coinbases.Load()
	cbrs, ok := cbRaw.(*map[int32]*util.Tx)
	if !ok {
		Debug("coinbases not correct type", cbrs)
		return
	}
	var cb *util.Tx
	cb, ok = (*cbrs)[msgBlock.Header.Version]
	if !ok {
		Debug("coinbase not found")
		return
	}
	Debug("copying over transactions")
	t := c.transactions.Load()
	var rtx []*util.Tx
	rtx, ok = t.([]*util.Tx)
	var txs []*util.Tx
	// copy merkle root
	txs = append(rtx, cb)
	for i := range txs {
		msgBlock.Transactions = append(msgBlock.Transactions, txs[i].MsgTx())
	}
	// set old blocks to pause and send pause directly as block is probably a
	// solution
	Debug("sending pause to workers")
	if err = c.multiConn.SendMany(pause.Magic, c.pauseShards); Check(err) {
		return
	}
	block := util.NewBlock(msgBlock)
	var isOrphan bool
	Debug("submitting block for processing")
	if isOrphan, err = c.cx.RealNode.SyncManager.ProcessBlock(block, blockchain.BFNone); Check(err) {
		// Anything other than a rule violation is an unexpected error, so log that
		// error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			Warnf(
				"Unexpected error while processing block submitted via kopach miner:", err,
			)
			return
		} else {
			Warn("block submitted via kopach miner rejected:", err)
			if isOrphan {
				Debug("block is an orphan")
				return
			}
			return
		}
	}
	Trace("the block was accepted")
	c.height.Store(block.Height())
	Tracec(
		func() string {
			bmb := block.MsgBlock()
			coinbaseTx := bmb.Transactions[0].TxOut[0]
			prevHeight := block.Height() - 1
			prevBlock, _ := c.cx.RealNode.Chain.BlockByHeight(prevHeight)
			prevTime := prevBlock.MsgBlock().Header.Timestamp.Unix()
			since := bmb.Header.Timestamp.Unix() - prevTime
			bHash := bmb.BlockHashWithAlgos(block.Height())
			return fmt.Sprintf(
				"new block height %d %08x %s%10d %08x %v %s %ds since prev",
				block.Height(),
				prevBlock.MsgBlock().Header.Bits,
				bHash,
				bmb.Header.Timestamp.Unix(),
				bmb.Header.Bits,
				util.Amount(coinbaseTx.Value),
				fork.GetAlgoName(
					bmb.Header.Version,
					block.Height(),
				), since,
			)
		},
	)
	return
}

// hashrate reports from workers
func processHashrateMsg(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
	c := ctx.(*Controller)
	var hr hashrate.Hashrate
	gotiny.Unmarshal(b, &hr)
	if c.lastNonce == hr.Nonce {
		return
	}
	c.lastNonce = hr.Nonce
	// add to total hash counts
	c.hashCount.Store(c.hashCount.Load() + uint64(hr.Count))
	return
}
