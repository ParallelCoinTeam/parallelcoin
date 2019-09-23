//+build !nogui
// +build !headless

package vue

import (
	enjs "encoding/json"
	"github.com/p9c/pod/cmd/gui/vue/alert"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/zserge/webview"
	"log"
	"strings"
)

// func (v *DuoVUEnode) Addnode(a *json.AddNodeCmd) {
// 	r, err := v.cx.RPCServer.HandleAddNode(v.cx.RPCServer, a, nil)
// 	return
// }
// func (v *DuoVUEnode) Createrawtransaction(a *json.CreateRawTransactionCmd) {
// 	r, err := v.cx.RPCServer.HandleCreateRawTransaction(v.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
// func (v *DuoVUEnode) Decoderawtransaction(a *json.DecodeRawTransactionCmd) {
// 	r, err := v.cx.RPCServer.HandleDecodeRawTransaction(v.cx.RPCServer, a, nil)
// 	r = json.TxRawDecodeResult{}
// 	return
// }
// func (v *DuoVUEnode) Decodescript(a *json.DecodeScriptCmd) {
// 	r, err := v.cx.RPCServer.HandleDecodeScript(v.cx.RPCServer, a, nil)
// 	return
// }
// func (v *DuoVUEnode) Estimatefee(a *json.EstimateFeeCmd) {
// 	r, err := v.cx.RPCServer.HandleEstimateFee(v.cx.RPCServer, a, nil)
// 	r = 0.0
// 	return
// }
// func (v *DuoVUEnode) Generate(a *json.GenerateCmd) {
// 	r, err := v.cx.RPCServer.HandleGenerate(v.cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuoVUEnode) Getaddednodeinfo(a *json.GetAddedNodeInfoCmd) {
// 	r, err := v.cx.RPCServer.HandleGetAddedNodeInfo(v.cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuoVUEnode) Getbestblock() {
// 	r, err := v.cx.RPCServer.HandleGetBestBlock(v.cx.RPCServer, a, nil)
// 	r = json.GetBestBlockResult{}
// 	return
// }
// func (v *DuoVUEnode) Getbestblockhash() {
// 	r, err := v.cx.RPCServer.HandleGetBestBlockHash(v.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
// func (v *DuoVUEnode) Getblock(a *json.GetBlockCmd) {
// 	r, err := v.cx.RPCServer.HandleGetBlock(v.cx.RPCServer, a, nil)
// 	r = json.GetBlockVerboseResult{}
// 	return
// }
func (d *DuoVUE) GetBlockChainInfo() {
	getBlockChainInfo, err := rpc.HandleGetBlockChainInfo(d.cx.RPCServer, nil, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	var ok bool
	d.Core.Node.BlockChainInfo, ok = getBlockChainInfo.(*btcjson.GetBlockChainInfoResult)
	if !ok {
		d.Core.Node.BlockChainInfo = &btcjson.GetBlockChainInfoResult{}
	}

}

func (d *DuoVUE) GetBlockCount() int64 {
	getBlockCount, err := rpc.HandleGetBlockCount(d.cx.RPCServer, nil, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	d.Status.BlockCount = getBlockCount.(int64)
	return d.Status.BlockCount
}
func (d *DuoVUE) GetBlockHash(blockHeight int) {
	hcmd := btcjson.GetBlockHashCmd{
		Index: int64(blockHeight),
	}
	hash, err := rpc.HandleGetBlockHash(d.cx.RPCServer, &hcmd, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	d.Core.Node.BlockHash = hash.(string)
}
func (d *DuoVUE) GetBlock(hash string) {
	verbose, verbosetx := true, true
	bcmd := btcjson.GetBlockCmd{
		Hash:      hash,
		Verbose:   &verbose,
		VerboseTx: &verbosetx,
	}
	bl, err := rpc.HandleGetBlock(d.cx.RPCServer, &bcmd, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	d.Core.Node.Block = bl.(btcjson.GetBlockVerboseResult)
}

// func (v *DuoVUEnode) Getblockheader(a *json.GetBlockHeaderCmd) {
// 	r, err := v.cx.RPCServer.HandleGetBlockHeader(v.cx.RPCServer, a, nil)
// 	r = json.GetBlockHeaderVerboseResult{}
// 	return
// }

func (d *DuoVUE) GetConnectionCount() int32 {
	d.Status.ConnectionCount = d.cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	return d.Status.ConnectionCount
}

func (d *DuoVUE) GetDifficulty() float64 {
	c := btcjson.GetDifficultyCmd{}
	r, err := rpc.HandleGetDifficulty(d.cx.RPCServer, c, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	d.Status.Difficulty = r.(float64)
	return d.Status.Difficulty
}

// func (v *DuoVUEnode) Gethashespersec() {
// 	r, err := v.cx.RPCServer.HandleGetHashesPerSec(v.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
// func (v *DuoVUEnode) Getheaders(a *json.GetHeadersCmd) {
// 	r, err := v.cx.RPCServer.HandleGetHeaders(v.cx.RPCServer, a, nil)
// 	r = []string{}
// 	return
// }
// func (v *DuoVUEnode) Getinfo() {
// 	r, err := v.cx.RPCServer.HandleGetInfo(v.cx.RPCServer, a, nil)
// 	r = json.InfoChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmempoolinfo() {
// 	r, err := v.cx.RPCServer.HandleGetMempoolInfo(v.cx.RPCServer, a, nil)
// 	r = json.GetMempoolInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getmininginfo() {
// 	r, err := v.cx.RPCServer.HandleGetMiningInfo(v.cx.RPCServer, a, nil)
// 	r = json.GetMiningInfoResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnettotals() {
// 	r, err := v.cx.RPCServer.HandleGetNetTotals(v.cx.RPCServer, a, nil)
// 	r = json.GetNetTotalsResult{}
// 	return
// }
// func (v *DuoVUEnode) Getnetworkhashps(a *json.GetNetworkHashPSCmd) {
// 	r, err := v.cx.RPCServer.HandleGetNetworkHashPS(v.cx.RPCServer, a, nil)
// 	r = int64(0)
// 	return
// }
func (d *DuoVUE) GetPeerInfo() []*btcjson.GetPeerInfoResult {
	getPeers, err := rpc.HandleGetPeerInfo(d.cx.RPCServer, nil, nil)
	if err != nil {
		alert.PushAlert(err.Error(), "error")
	}
	d.Status.Peers = getPeers.([]*btcjson.GetPeerInfoResult)
	//fmt.Println("ssssssssssssssssss", d.Status.Peers)
	return d.Status.Peers
}

// func (v *DuoVUEnode) Stop() {
// 	r, err := v.cx.RPCServer.HandleStop(v.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (d *DuoVUE) Uptime() (r int64) {
	rRaw, err := rpc.HandleUptime(d.cx.RPCServer, nil, nil)
	if err != nil {
	}
	//rRaw = int64(0)
	d.Status.UpTime = rRaw.(int64)
	return d.Status.UpTime
}

// func (v *DuoVUEnode) Validateaddress(a *json.ValidateAddressCmd) {
// 	r, err := v.cx.RPCServer.HandleValidateAddress(v.cx.RPCServer, a, nil)
// 	r = json.ValidateAddressChainResult{}
// 	return
// }
// func (v *DuoVUEnode) Verifychain(a *json.VerifyChainCmd) {
// 	r, err := v.cx.RPCServer.HandleVerifyChain(v.cx.RPCServer, a, nil)
// }
// func (v *DuoVUEnode) Verifymessage(a *json.VerifyMessageCmd) {
// 	r, err := v.cx.RPCServer.HandleVerifyMessage(v.cx.RPCServer, a, nil)
// 	r = ""
// 	return
// }
func (s *DuoVUEstatus) GetWalletVersion(d DuoVUE) map[string]btcjson.VersionResult {
	v, err := rpc.HandleVersion(d.cx.RPCServer, nil, nil)
	if err != nil {
	}
	return v.(map[string]btcjson.VersionResult)
}

func render(w webview.WebView, cmd string, data interface{}) {
	b, err := enjs.Marshal(data)
	if err == nil {
		w.Eval("duoSystem." + cmd + "=" + string(b) + ";")
	}
}

func (d *DuoVUE) HandleRPC(w webview.WebView, vc string) {
	switch {
	case vc == "close":
		w.Terminate()
	case vc == "fullscreen":
		w.SetFullscreen(true)
	case vc == "unfullscreen":
		w.SetFullscreen(false)
	case strings.HasPrefix(vc, "changeTitle:"):
		w.SetTitle(strings.TrimPrefix(vc, "changeTitle:"))
	case vc == "addressBook":
		render(w, vc, d.GetAddressBook())
	case vc == "balance":
		render(w, vc, d.GetBalance())
	case vc == "status":
		render(w, vc, d.GetStatus())
	case strings.HasPrefix(vc, "transactions:"):
		t := strings.TrimPrefix(vc, "transactions:")
		cmd := struct {
			From  int    `json:"from"`
			Count int    `json:"count"`
			C     string `json:"c"`
		}{}
		if err := enjs.Unmarshal([]byte(t), &cmd); err != nil {
			log.Println(err)
		}
		render(w, "transactions", d.GetTransactions(cmd.From, cmd.Count, cmd.C))
	case strings.HasPrefix(vc, "send:"):
		s := strings.TrimPrefix(vc, "send:")
		cmd := struct {
			Wp string  `json:"wp"`
			Ad string  `json:"ad"`
			Am float64 `json:"am"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}

		render(w, "send", d.DuoSend(cmd.Wp, cmd.Ad, cmd.Am))
	case strings.HasPrefix(vc, "createAddress:"):
		s := strings.TrimPrefix(vc, "createAddress:")
		cmd := struct {
			Account string `json:"account"`
			Label   string `json:"label"`
		}{}
		if err := enjs.Unmarshal([]byte(s), &cmd); err != nil {
			log.Println(err)
		}
		render(w, "send", d.CreateNewAddress(cmd.Account, cmd.Label))

	}
}
