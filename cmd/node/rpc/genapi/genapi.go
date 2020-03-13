package main

import (
	"os"
	"text/template"

	log "github.com/p9c/logi"
)

type handler struct {
	Method, Handler, Cmd, ResType string
}

var handlers = []handler{
	{
		Method:  "addnode",
		Handler: "AddNode",
		Cmd:     "btcjson.AddNodeCmd",
		ResType: "None",
	},
	{
		Method:  "createrawtransaction",
		Handler: "CreateRawTransaction",
		Cmd:     "btcjson.CreateRawTransactionCmd",
		ResType: "string",
	},
	{
		Method:  "decoderawtransaction",
		Handler: "DecodeRawTransaction",
		Cmd:     "btcjson.DecodeRawTransactionCmd",
		ResType: "btcjson.TxRawDecodeResult",
	},
	{
		Method:  "decodescript",
		Handler: "DecodeScript",
		Cmd:     "btcjson.DecodeScriptCmd",
		ResType: "btcjson.DecodeScriptResult",
	},
	{
		Method:  "estimatefee",
		Handler: "EstimateFee",
		Cmd:     "btcjson.EstimateFeeCmd",
		ResType: "float64",
	},
	{
		Method:  "generate",
		Handler: "Generate",
		Cmd:     "None",
		ResType: "[]string",
	},
	{
		Method:  "getaddednodeinfo",
		Handler: "GetAddedNodeInfo",
		Cmd:     "btcjson.GetAddedNodeInfoCmd",
		ResType: "[]btcjson.GetAddedNodeInfoResultAddr",
	},
	{
		Method:  "getbestblock",
		Handler: "GetBestBlock",
		Cmd:     "None",
		ResType: "btcjson.GetBestBlockResult",
	},
	{
		Method:  "getbestblockhash",
		Handler: "GetBestBlockHash",
		Cmd:     "None",
		ResType: "string",
	},
	{
		Method:  "getblock",
		Handler: "GetBlock",
		Cmd:     "btcjson.GetBlockCmd",
		ResType: "btcjson.GetBlockVerboseResult",
	},
	{
		Method:  "getblockchaininfo",
		Handler: "GetBlockChainInfo",
		Cmd:     "None",
		ResType: "btcjson.GetBlockChainInfoResult",
	},
	{
		Method:  "getblockcount",
		Handler: "GetBlockCount",
		Cmd:     "None",
		ResType: "int64",
	},
	{
		Method:  "getblockhash",
		Handler: "GetBlockHash",
		Cmd:     "btcjson.GetBlockHashCmd",
		ResType: "string",
	},
	{
		Method:  "getblockheader",
		Handler: "GetBlockHeader",
		Cmd:     "btcjson.GetBlockHeaderCmd",
		ResType: "btcjson.GetBlockHeaderVerboseResult",
	},
	{
		Method:  "getblocktemplate",
		Handler: "GetBlockTemplate",
		Cmd:     "btcjson.GetBlockTemplateCmd",
		ResType: "string",
	},
	{
		Method:  "getcfilter",
		Handler: "GetCFilter",
		Cmd:     "btcjson.GetCFilterCmd",
		ResType: "string",
	},
	{
		Method:  "getcfilterheader",
		Handler: "GetCFilterHeader",
		Cmd:     "btcjson.GetCFilterHeaderCmd",
		ResType: "string",
	},
	{
		Method:  "getconnectioncount",
		Handler: "GetConnectionCount",
		Cmd:     "None",
		ResType: "int32",
	},
	{
		Method:  "getcurrentnet",
		Handler: "GetCurrentNet",
		Cmd:     "None",
		ResType: "string",
	},
	{
		Method:  "getdifficulty",
		Handler: "GetDifficulty",
		Cmd:     "btcjson.GetDifficultyCmd",
		ResType: "float64",
	},
	{
		Method:  "getgenerate",
		Handler: "GetGenerate",
		Cmd:     "btcjson.GetHeadersCmd",
		ResType: "bool",
	},
	{
		Method:  "gethashespersec",
		Handler: "GetHashesPerSec",
		Cmd:     "None",
		ResType: "float64",
	},
	{
		Method:  "getheaders",
		Handler: "GetHeaders",
		Cmd:     "btcjson.GetHeadersCmd",
		ResType: "[]string",
	},
	{
		Method:  "getinfo",
		Handler: "GetInfo",
		Cmd:     "None",
		ResType: "btcjson.InfoChainResult0",
	},
	{
		Method:  "getmempoolinfo",
		Handler: "GetMempoolInfo",
		Cmd:     "None",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "getmininginfo",
		Handler: "GetMiningInfo",
		Cmd:     "None",
		ResType: "btcjson.GetMiningInfoResult",
	},
	{
		Method:  "getnettotals",
		Handler: "GetNetTotals",
		Cmd:     "None",
		ResType: "btcjson.GetNetTotalsResult",
	},
	{
		Method:  "getnetworkhashps",
		Handler: "GetNetworkHashPS",
		Cmd:     "btcjson.GetNetworkHashPSCmd",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getpeerinfo",
		Handler: "GetPeerInfo",
		Cmd:     "None",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getrawmempool",
		Handler: "GetRawMempool",
		Cmd:     "btcjson.GetRawMempoolCmd",
		ResType: "[]string",
	},
	{
		Method:  "getrawtransaction",
		Handler: "GetRawTransaction",
		Cmd:     "btcjson.GetRawTransactionCmd",
		ResType: "string",
	},
	{
		Method:  "gettxout",
		Handler: "GetTxOut",
		Cmd:     "btcjson.GetTxOutCmd",
		ResType: "string",
	},
	{
		Method:  "help",
		Handler: "Help",
		Cmd:     "btcjson.HelpCmd",
		ResType: "string",
	},
	{
		Method:  "node",
		Handler: "Node",
		Cmd:     "btcjson.NodeCmd",
		ResType: "None",
	},
	{
		Method:  "ping",
		Handler: "Ping",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "searchrawtransactions",
		Handler: "SearchRawTransactions",
		Cmd:     "btcjson.SearchRawTransactionsCmd",
		ResType: "[]btcjson.SearchRawTransactionsResult",
	},
	{
		Method:  "sendrawtransaction",
		Handler: "SendRawTransaction",
		Cmd:     "btcjson.SendRawTransactionCmd",
		ResType: "None",
	},
	{
		Method:  "setgenerate",
		Handler: "SetGenerate",
		Cmd:     "btcjson.SetGenerateCmd",
		ResType: "None",
	},
	{
		Method:  "stop",
		Handler: "Stop",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "restart",
		Handler: "Restart",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "resetchain",
		Handler: "ResetChain",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "submitblock",
		Handler: "SubmitBlock",
		Cmd:     "btcjson.SubmitBlockCmd",
		ResType: "string",
	},
	{
		Method:  "uptime",
		Handler: "Uptime",
		Cmd:     "None",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "validateaddress",
		Handler: "ValidateAddress",
		Cmd:     "btcjson.ValidateAddressCmd",
		ResType: "btcjson.ValidateAddressChainResult",
	},
	{
		Method:  "verifychain",
		Handler: "VerifyChain",
		Cmd:     "btcjson.VerifyChainCmd",
		ResType: "bool",
	},
	{
		Method:  "verifymessage",
		Handler: "VerifyMessage",
		Cmd:     "btcjson.VerifyMessageCmd",
		ResType: "bool",
	},
	{
		Method:  "version",
		Handler: "Version",
		Cmd:     "btcjson.VersionCmd",
		ResType: "map[string]btcjson.VersionResult",
	},
}

func main() {
	log.L.SetLevel("trace", true, "pod")
	if fd, err := os.Create("../rpchandlers.go"); log.L.Check(err) {
	} else {
		defer fd.Close()
		t := template.Must(template.New("noderpc").Parse(NodeRPCHandlerTpl))
		if err = t.Execute(fd, handlers); log.L.Check(err) {
		}
	}
}

var NodeRPCHandlerTpl = `// generated by go run gen.go; DO NOT EDIT

package rpc

import (
	log "github.com/p9c/logi"

	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type API struct {
	Ch     interface{}
	Params interface{}
}

var RPCHandlersBeforeInit = map[string]CommandHandler{
{{range .}}	"{{.Method}}":{ 
		Fn: Handle{{.Handler}}, 
		Call: make(chan API), 
		Result: func() API {
			return API{
				Ch: make(chan {{.Handler}}Res),
			}
		},
	}, 
{{end}}
}

type (
	None struct{} {{range .}}
	{{.Handler}}Res struct {
		Res {{.ResType}}
		Err error
	}
	{{end}}
)
{{range .}}
func (a API) {{.Handler}}(cmd {{.Cmd}}) (err error) {
		RPCHandlers["{{.Method}}"].Call <- API{a.Ch, cmd}
		return
}
{{end}}

// StartAPI starts up the api handler server that receives rpc.API messages and runs the handler and returns the result
// Note that the parameters are type asserted to prevent the consumer of the API from sending wrong message types not
// because it's necessary since they are interfaces end to end
func StartAPI(server *Server, quit chan struct{}) {
	nrh := RPCHandlers
	go func() {
		var err error
		var res interface{}
		for {
			select { {{range .}}
			case msg := <-nrh["{{.Method}}"].Call:
				if res, err = nrh["{{.Method}}"].
					Fn(server, msg.Params.({{.Cmd}}), nil); log.L.Check(err) {
				}
				msg.Ch.(chan {{.Handler}}Res) <- {{.Handler}}Res{
					Res: res.({{.ResType}}), Err: err} {{end}}
			case <-quit:
				return
			}
		}
	}()
}
`
