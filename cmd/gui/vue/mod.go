package vue

import (
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/db"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/mod"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/rpc/json"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type DuoVUE struct {
	cx         *conte.Xt
	db         db.DuoVUEdb       `json:"db"`
	Core       DuoVUEcore        `json:"core"`
	Config     DuoVUEConfig      `json:"conf"`
	Components []mod.DuoVUEcomp  `json:"comp"`
	Repo       []mod.DuoVUEcomp  `json:"repo"`
	Status     DuoVUEstatus      `json:"stat"`
	Icons      map[string]string `json:"ico"`
}

type DuoVUEcore struct {
	*mod.DuoGuiItem
	VUE      string     `json:"vue"`
	CoreHtml string     `json:"html"`
	CoreJs   []byte     `json:"js"`
	CoreCss  []byte     `json:"css"`
	Node     DuoVUEnode `json:"node"`
}
type DuoVUEnode struct {
	rpc              *rpc.Server
	IsCurrent        bool                             `json:"iscurrent"`
	NetworkLastBlock int32                            `json:"networklastblock"`
	PeerInfo         []*json.GetPeerInfoResult        `json:"peerinfo"`
	Tx               json.GetTransactionResult        `json:"tx"`
	TxDetail         json.GetTransactionDetailsResult `json:"txdetail"`
	BlockChainInfo   *json.GetBlockChainInfoResult    `json:"blockchaininfo"`
	BlockCount       int32                            `json:"blockcount"`
	BlockHash        string                           `json:"blockhash"`
	Block            json.GetBlockVerboseResult       `json:"block"`
	ConnectionCount  int32                            `json:"connectioncount"`
}

type DuoVUEbalance struct {
	Balance     string `json:"balance"`
	Unconfirmed string `json:"unconfirmed"`
}

type DuoVUEtransactions struct {
	Txs       []json.ListTransactionsResult `json:"txs"`
	TxsNumber int                           `json:"txsnumber"`
}

type DuoVUEAddressBook struct {
	Num       int           `json:"num"`
	Addresses []mod.Address `json:"addresses"`
}

type DuoVUEblock struct {
	rpc        *rpc.Server
	Height     int64   `json:"height"`
	PowAlgoID  uint32  `json:"pow"`
	Difficulty float64 `json:"diff"`
	Amount     float64 `json:"amount"`
	TxNum      int     `json:"txnum"`
	Time       int64   `json:"time"`
}

type DuoVUEchain struct {
	rpc        *rpc.Server
	LastPushed int64         `json:"lastpushed"`
	Blocks     []DuoVUEblock `json:"blocks"`
}

// System Ststus
type DuoVUEstatus struct {
	dv *DuoVUE
	//BCI           interface{}                   `json:"bci"`
	Version       string                        `json:"ver"`
	WalletVersion map[string]json.VersionResult `json:"walletver"`
	UpTime        int64                         `json:"uptime"`
	Cpu           []cpu.InfoStat                `json:"cpu"`
	CpuPercent    []float64                     `json:"cpupercent"`
	Memory        mem.VirtualMemoryStat         `json:"mem"`
	Disk          disk.UsageStat                `json:"disk"`
	CurrentNet    string                        `json:"net"`
	Chain         string                        `json:"chain"`
	HashesPerSec  int64                         `json:"hashrate"`
	Height        int32                         `json:"height"`
	BestBlockHash string                        `json:"bestblockhash"`
	NetworkHashPS int64                         `json:"networkhashrate"`
	Peers         []*json.GetPeerInfoResult     `json:"peers"`
	//MempoolInfo      string                        `json:"ver"`
	Difficulty       float64       `json:"diff"`
	Balance          DuoVUEbalance `json:"balance"`
	BlockCount       int64         `json:"blockcount"`
	ConnectionCount  int32         `json:"connectioncount"`
	NetworkLastBlock int32         `json:"networklastblock"`
	TxsNumber        int           `json:"txsnumber"`
}

func (d *DuoVUE) GetStatus() DuoVUEstatus {
	//getBlockChainInfo, _ := rpc.HandleGetBlockChainInfo(d.cx.RPCServer, nil, nil)
	//d.Status.BCI = getBlockChainInfo
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	d.Status.Cpu = sc
	d.Status.CpuPercent = sp
	d.Status.Memory = *sm
	d.Status.Disk = *sd

	params := d.cx.RPCServer.Cfg.ChainParams
	chain := d.cx.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	gnhpsCmd := json.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(d.cx.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	v, err := rpc.HandleVersion(d.cx.RPCServer, nil, nil)
	if err != nil {
	}
	d.Status.Version = "0.0.1"
	d.Status.WalletVersion = v.(map[string]json.VersionResult)
	d.Status.UpTime = time.Now().Unix() - d.cx.RPCServer.Cfg.StartupTime
	d.Status.CurrentNet = d.cx.RPCServer.Cfg.ChainParams.Net.String()
	d.Status.NetworkHashPS = networkHashesPerSec
	//s.MempoolInfo =
	d.Status.HashesPerSec = int64(d.cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	d.Status.Chain = params.Name
	d.Status.Height = chainSnapshot.Height
	d.Status.Peers = d.GetPeerInfo()
	//s.Headers = chainSnapshot.Height
	d.Status.BestBlockHash = chainSnapshot.Hash.String()
	d.Status.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	d.Status.Balance.Balance = d.GetBalance().Balance
	d.Status.Balance.Unconfirmed = d.GetBalance().Unconfirmed
	d.Status.BlockCount = d.GetBlockCount()
	d.Status.ConnectionCount = d.GetConnectionCount()
	d.Status.NetworkLastBlock = d.GetNetworkLastBlock()
	return d.Status
}
