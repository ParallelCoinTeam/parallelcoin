package stat

import (
	"github.com/p9c/pod/cmd/node/rpc"
	conte2 "github.com/p9c/pod/gui/____BEZI/test/pkg/conte"
	mod2 "github.com/p9c/pod/gui/____BEZI/test/pkg/duos/mod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

// System Ststus
type DuOSstatus struct {
	Version          string                           `json:"ver"`
	WalletVersion    map[string]btcjson.VersionResult `json:"walletver"`
	UpTime           int64                            `json:"uptime"`
	Cpu              []cpu.InfoStat                   `json:"cpu"`
	CpuPercent       []float64                        `json:"cpupercent"`
	Memory           mem.VirtualMemoryStat            `json:"mem"`
	Disk             disk.UsageStat                   `json:"disk"`
	CurrentNet       string                           `json:"net"`
	Chain            string                           `json:"chain"`
	HashesPerSec     int64                            `json:"hashrate"`
	Height           int32                            `json:"height"`
	BestBlockHash    string                           `json:"bestblockhash"`
	NetworkHashPS    int64                            `json:"networkhashrate"`
	Difficulty       float64                          `json:"diff"`
	Balance          mod2.DuOSbalance                 `json:"balance"`
	BlockCount       int64                            `json:"blockcount"`
	ConnectionCount  int32                            `json:"connectioncount"`
	NetworkLastBlock int32                            `json:"networklastblock"`
	TxsNumber        int                              `json:"txsnumber"`
	//MempoolInfo      string                        `json:"ver"`
}

func (s *DuOSstatus) GetDuOSstatus(cx *conte2.Xt) *DuOSstatus {
	s = new(DuOSstatus)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	s.Cpu = sc
	s.CpuPercent = sp
	s.Memory = *sm
	s.Disk = *sd
	params := cx.RPCServer.Cfg.ChainParams
	chain := cx.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(cx.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	v, err := rpc.HandleVersion(cx.RPCServer, nil, nil)
	if err != nil {
	}
	s.Version = "0.0.1"
	s.WalletVersion = v.(map[string]btcjson.VersionResult)
	s.UpTime = time.Now().Unix() - cx.RPCServer.Cfg.StartupTime
	s.CurrentNet = cx.RPCServer.Cfg.ChainParams.Net.String()
	s.NetworkHashPS = networkHashesPerSec
	s.HashesPerSec = int64(cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	s.Chain = params.Name
	s.Height = chainSnapshot.Height
	//s.Headers = chainSnapshot.Height
	s.BestBlockHash = chainSnapshot.Hash.String()
	s.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	//s.Balance.Balance = s.Balance.GetBalance().Balance
	//s.Balance.Unconfirmed = s.GetBalance().Unconfirmed
	//s.BlockCount = s.GetBlockCount()
	//s.ConnectionCount = s.GetConnectionCount()
	//s.NetworkLastBlock = s.GetNetworkLastBlock()
	return s
}
