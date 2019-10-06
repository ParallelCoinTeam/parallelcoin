package core

import (
	"github.com/p9c/pod/cmd/gui/db"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func CoreJs(d db.DuOSdb) string {
	return `
Vue.config.devtools = true;
Vue.use(VueFormGenerator);

Vue.prototype.$eventHub = new Vue();
 
const duoSystem = {
	config:null,
	node:null,
	wallet:null,
	status:null,
	theme:false,
	isBoot:false,
	isLoading:false,
	isDev:true,
	isScreen:'overview',
	timer: '',
};

const core = new Vue({ 
	data () { return { 
	duoSystem }},
});`
}

// GetMsg loads the message variable
func (d *DuOS) PushDuOSalert(t string, m interface{}, at string) {
	a := new(DuOSalert)
	a.Time = time.Now()
	a.Title = t
	a.Message = m
	a.AlertType = at
	//d.Render("alert", a)
}

func (d *DuOS) GetDuOSstatus() mod.DuOSstatus {
	status := *new(mod.DuOSstatus)
	sm, _ := mem.VirtualMemory()
	sc, _ := cpu.Info()
	sp, _ := cpu.Percent(0, true)
	sd, _ := disk.Usage("/")
	status.Cpu = sc
	status.CpuPercent = sp
	status.Memory = *sm
	status.Disk = *sd
	params := d.Cx.RPCServer.Cfg.ChainParams
	chain := d.Cx.RPCServer.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, err := rpc.HandleGetNetworkHashPS(d.Cx.RPCServer, gnhpsCmd, nil)
	if err != nil {
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
	}
	v, err := rpc.HandleVersion(d.Cx.RPCServer, nil, nil)
	if err != nil {
	}
	status.Version = "0.0.1"
	status.WalletVersion = v.(map[string]btcjson.VersionResult)
	status.UpTime = time.Now().Unix() - d.Cx.RPCServer.Cfg.StartupTime
	status.CurrentNet = d.Cx.RPCServer.Cfg.ChainParams.Net.String()
	status.NetworkHashPS = networkHashesPerSec
	status.HashesPerSec = int64(d.Cx.RPCServer.Cfg.CPUMiner.HashesPerSecond())
	status.Chain = params.Name
	status.Height = chainSnapshot.Height
	//s.Headers = chainSnapshot.Height
	status.BestBlockHash = chainSnapshot.Hash.String()
	status.Difficulty = rpc.GetDifficultyRatio(chainSnapshot.Bits, params, 2)
	status.Balance.Balance = d.GetBalance().Balance
	status.Balance.Unconfirmed = d.GetBalance().Unconfirmed
	status.BlockCount = d.GetBlockCount()
	status.ConnectionCount = d.GetConnectionCount()
	status.NetworkLastBlock = d.GetNetworkLastBlock()
	return status
}
