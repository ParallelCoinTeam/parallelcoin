package rcd

import (
	"time"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
)

// System Ststus

func (r *RcVar) GetDuoUIstatus() {
	v, err := chainrpc.HandleVersion(r.cx.RPCServer, nil, nil)
	if err != nil {
	}
	r.Status.Version = "0.0.1"
	r.Status.Wallet.WalletVersion = v.(map[string]btcjson.VersionResult)
	r.Status.StartTime = time.Unix(0, r.cx.RPCServer.Cfg.StartupTime)
	r.Status.CurrentNet = r.cx.RPCServer.Cfg.ChainParams.Net.String()
	r.Status.Chain = r.cx.RPCServer.Cfg.ChainParams.Name
	return
}

func (r *RcVar) GetDuoUIlocalLost() {
	r.Localhost = *new(model.DuoUIlocalHost)
	// sm, _ := mem.VirtualMemory()
	// sc, _ := cpu.Info()
	// sp, _ := cpu.Percent(0, true)
	// sd, _ := disk.Usage("/")
	// r.Localhost.Cpu = sc
	// r.Localhost.CpuPercent = sp
	// r.Localhost.Memory = *sm
	// r.Localhost.Disk = *sd
	return
}
