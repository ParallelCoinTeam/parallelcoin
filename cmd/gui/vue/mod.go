package vue

import (
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/db"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/rpc/json"
	"github.com/robfig/cron"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/zserge/webview"
	"time"
)


type DuoVUE struct {
	cx         *conte.Xt
	cr         *cron.Cron
	Web        webview.WebView
	db         db.DuoVUEdb       `json:"db"`
	Core       DuoVUEcore        `json:"core"`
	Config     DuoVUEConfig      `json:"conf"`
	Components []mod.DuoVUEcomp  `json:"comp"`
	Repo       []mod.DuoVUEcomp  `json:"repo"`
	Data       DuoVUEdata        `json:"d"`
	Icons      map[string]string `json:"ico"`
}

type DuoVUEdata struct {
	Alert        DuoVUEalert               `json:"alert"`
	Status       mod.DuoVUEstatus          `json:"status"`
	Addressbook  mod.DuoVUEaddressBook     `json:"addressbook"`
	Peers        []*json.GetPeerInfoResult `json:"peers"`
	Transactions mod.DuoVUEtransactions    `json:"transactions"`
	Blocks       mod.DuoVUEblocks          `json:"blocks"`
	Send         mod.Send                  `json:"send"`
}

type DuoVUEcore struct {
	*mod.DuoGuiItem
	VUE      string `json:"vue"`
	CoreHtml string `json:"html"`
	CoreJs   []byte `json:"js"`
	CoreCss  []byte `json:"css"`
}

func (dv *DuoVUE) GetDuoVUEstatus() (mod.DuoVUEstatus) {
	status := *new(mod.DuoVUEstatus)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	status.Cpu = sc
	status.CpuPercent = sp
	status.Memory = *sm
	status.Disk = *sd
	params := dv.cx.RPCServer.Cfg.ChainParams
	chain := dv.cx.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	gnhpsCmd := json.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(dv.cx.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	v, err := rpc.HandleVersion(dv.cx.RPCServer, nil, nil)
	if err != nil {
	}
	status.Version = "0.0.1"
	status.WalletVersion = v.(map[string]json.VersionResult)
	status.UpTime = time.Now().Unix() - dv.cx.RPCServer.Cfg.StartupTime
	status.CurrentNet = dv.cx.RPCServer.Cfg.ChainParams.Net.String()
	status.NetworkHashPS = networkHashesPerSec
	status.HashesPerSec = int64(dv.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	status.Chain = params.Name
	status.Height = chainSnapshot.Height
	//s.Headers = chainSnapshot.Height
	status.BestBlockHash = chainSnapshot.Hash.String()
	status.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	status.Balance.Balance = dv.GetBalance().Balance
	status.Balance.Unconfirmed = dv.GetBalance().Unconfirmed
	status.BlockCount = dv.GetBlockCount()
	status.ConnectionCount = dv.GetConnectionCount()
	status.NetworkLastBlock = dv.GetNetworkLastBlock()
	return status
}
