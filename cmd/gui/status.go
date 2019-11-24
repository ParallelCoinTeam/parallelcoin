package gui

import (
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
)

// System Ststus
type DuOStatus struct {
	Version       string                           `json:"ver"`
	WalletVersion map[string]btcjson.VersionResult `json:"walletver"`
	UpTime        int64                            `json:"uptime"`
	CurrentNet    string                           `json:"net"`
	Chain         string                           `json:"chain"`
	HashesPerSec  int64                            `json:"hashrate"`
	Height        int32                            `json:"height"`
	BestBlockHash string                           `json:"bestblockhash"`
	NetworkHashPS int64                            `json:"networkhashrate"`
	Difficulty    float64                          `json:"diff"`
	//MempoolInfo      string                        `json:"ver"`
}

type DuOSblockCount struct {
	int64
}
type DuOSconnections struct {
	int32
}
type DuOSnetLastBlock struct {
	int32
}
type DuOSnetworkHashPS struct {
	int64
}
type DuOShashesPerSec struct {
	int64
}



func
(r *rcvar) GetDuOStatus() {
	r.status = *new(DuOStatus)
	params := r.cx.RPCServer.Cfg.ChainParams
	chain := r.cx.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	//gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	//params := r.cx.RPCServer.Cfg.ChainParams
	//chain := r.cx.RPCServer.Cfg.Chain
	//chainSnapshot := chain.BestSnapshot()
	v, err := rpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.status.Version = "0.0.1"
	r.status.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.status.UpTime = time.Now().Unix() - r.cx.RPCServer.Cfg.StartupTime
	r.status.CurrentNet = r.cx.RPCServer.Cfg.ChainParams.Net.String()
	r.status.Chain = params.Name
	r.status.Height = chainSnapshot.Height
	//r.status.Headers = chainSnapshot.Height
	r.status.BestBlockHash = chainSnapshot.Hash.String()
	r.status.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	return
}
func
(r *rcvar) GetDuOSnetworkHashesPerSec() {
	gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(r.cx.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	r.nethash = networkHashesPerSec
}
func
(r *rcvar) GetDuOShashesPerSec() {
	r.hashes = int64(r.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
}
func
(r *rcvar) GetDuOSnetworkLastBlock() {
	for _, g := range rcv.cx.RPCServer.Cfg.ConnMgr.ConnectedPeers() {
		l := g.ToPeer().StatsSnapshot().LastBlock
		if l > r.netlastblock {
			r.netlastblock = l
		}
	}
	return
}
func
(r *rcvar) GetDuOSblockCount() {
	getBlockCount, err := rpc.HandleGetBlockCount(r.cx.RPCServer, nil, nil)
	if err != nil {
		r.PushDuOSalert("Error", err.Error(), "error")
	}
	r.blockcount = getBlockCount.(int64)
	return
}
func
(r *rcvar) GetDuOSconnectionCount() {
	r.connections = rcv.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return
}
