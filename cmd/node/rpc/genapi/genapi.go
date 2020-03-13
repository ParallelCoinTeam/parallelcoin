package main

import (
	"os"
	"text/template"

	log "github.com/p9c/logi"
)

type handler struct {
	Method, Handler, Cmd, Res, ResType string
}

var handlers = []handler{
	{
		Method:  "addnode",
		Handler: "HandleAddNode",
		Cmd:     "btcjson.AddNodeCmd",
		Res:     "AddNodeRes",
		ResType: "None",
	},
	{
		Method:  "createrawtransaction",
		Handler: "HandleCreateRawTransaction",
		Cmd:     "btcjson.CreateRawTransactionCmd",
		Res:     "CreateRawTransactionRes",
		ResType: "string",
	},
	{
		Method:  "decoderawtransaction",
		Handler: "HandleDecodeRawTransaction",
		Cmd:     "btcjson.DecodeRawTransactionCmd",
		Res:     "DecodeRawTransactionRes",
		ResType: "btcjson.TxRawDecodeResult",
	},
	{
		Method:  "decodescript",
		Handler: "HandleDecodeScript",
		Cmd:     "btcjson.DecodeScriptCmd",
		Res:     "DecodeScriptRes",
		ResType: "btcjson.DecodeScriptResult",
	},
	{
		Method:  "estimatefee",
		Handler: "HandleEstimateFee",
		Cmd:     "btcjson.EstimateFeeCmd",
		Res:     "EstimateFeeRes",
		ResType: "float64",
	},
	{
		Method:  "generate",
		Handler: "HandleGenerate",
		Cmd:     "nil",
		Res:     "GenerateRes",
		ResType: "[]string",
	},
	{
		Method:  "getaddednodeinfo",
		Handler: "HandleGetAddedNodeInfo",
		Cmd:     "btcjson.GetAddedNodeInfoCmd",
		Res:     "GetAddedNodeInfoRes",
		ResType: "[]btcjson.GetAddedNodeInfoResultAddr",
	},
	{
		Method:  "getbestblock",
		Handler: "HandleGetBestBlock",
		Cmd:     "nil",
		Res:     "GetBestBlockRes",
		ResType: "btcjson.GetBestBlockResult",
	},
	{
		Method:  "getbestblockhash",
		Handler: "HandleGetBestBlockHash",
		Cmd:     "nil",
		Res:     "GetBestBlockHashRes",
		ResType: "string",
	},
	{
		Method:  "getblock",
		Handler: "HandleGetBlock",
		Cmd:     "btcjson.GetBlockCmd",
		Res:     "GetBlockRes",
		ResType: "btcjson.GetBlockVerboseResult",
	},
	{
		Method:  "getblockchaininfo",
		Handler: "HandleGetBlockChainInfo",
		Cmd:     "nil",
		Res:     "GetBlockChainInfoRes",
		ResType: "btcjson.GetBlockChainInfoResult",
	},
	{
		Method:  "getblockcount",
		Handler: "HandleGetBlockCount",
		Cmd:     "nil",
		Res:     "GetBlockCountRes",
		ResType: "int64",
	},
	{
		Method:  "getblockhash",
		Handler: "HandleGetBlockHash",
		Cmd:     "btcjson.GetBlockHashCmd",
		Res:     "GetBlockHashRes",
		ResType: "string",
	},
	{
		Method:  "getblockheader",
		Handler: "HandleGetBlockHeader",
		Cmd:     "btcjson.GetBlockHeaderCmd",
		Res:     "GetBlockHeaderRes",
		ResType: "btcjson.GetBlockHeaderVerboseResult",
	},
	{
		Method:  "getblocktemplate",
		Handler: "HandleGetBlockTemplate",
		Cmd:     "btcjson.GetBlockTemplateCmd",
		Res:     "GetBlockTemplateRes",
		ResType: "string",
	},
	{
		Method:  "getcfilter",
		Handler: "HandleGetCFilter",
		Cmd:     "btcjson.GetCFilterCmd",
		Res:     "GetCFilterRes",
		ResType: "string",
	},
	{
		Method:  "getcfilterheader",
		Handler: "HandleGetCFilterHeader",
		Cmd:     "btcjson.GetCFilterHeaderCmd",
		Res:     "GetCFilterHeaderRes",
		ResType: "string",
	},
	{
		Method:  "getconnectioncount",
		Handler: "HandleGetConnectionCount",
		Cmd:     "nil",
		Res:     "GetConnectionCountRes",
		ResType: "int32",
	},
	{
		Method:  "getcurrentnet",
		Handler: "HandleGetCurrentNet",
		Cmd:     "nil",
		Res:     "GetCurrentNetRes",
		ResType: "string",
	},
	{
		Method:  "getdifficulty",
		Handler: "HandleGetDifficulty",
		Cmd:     "btcjson.GetDifficultyCmd",
		Res:     "GetDifficultyRes",
		ResType: "float64",
	},
	{
		Method:  "getgenerate",
		Handler: "HandleGetGenerate",
		Cmd:     "btcjson.GetHeadersCmd",
		Res:     "GetGenerateRes",
		ResType: "bool",
	},
	{
		Method:  "gethashespersec",
		Handler: "HandleGetHashesPerSec",
		Cmd:     "nil",
		Res:     "GetHashesPerSecRes",
		ResType: "float64",
	},
	{
		Method:  "getheaders",
		Handler: "HandleGetHeaders",
		Cmd:     "btcjson.GetHeadersCmd",
		Res:     "GetHeadersRes",
		ResType: "[]string",
	},
	{
		Method:  "getinfo",
		Handler: "HandleGetInfo",
		Cmd:     "nil",
		Res:     "GetInfoRes",
		ResType: "btcjson.InfoChainResult0",
	},
	{
		Method:  "getmempoolinfo",
		Handler: "HandleGetMempoolInfo",
		Cmd:     "nil",
		Res:     "GetMempoolInfoRes",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "getmininginfo",
		Handler: "HandleGetMiningInfo",
		Cmd:     "nil",
		Res:     "GetMiningInfoRes",
		ResType: "btcjson.GetMiningInfoResult",
	},
	{
		Method:  "getnettotals",
		Handler: "HandleGetNetTotals",
		Cmd:     "nil",
		Res:     "GetNetTotalsRes",
		ResType: "btcjson.GetNetTotalsResult",
	},
	{
		Method:  "getnetworkhashps",
		Handler: "HandleGetNetworkHashPS",
		Cmd:     "btcjson.GetNetworkHashPSCmd",
		Res:     "GetNetworkHashPSRes",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getpeerinfo",
		Handler: "HandleGetPeerInfo",
		Cmd:     "nil",
		Res:     "GetPeerInfoRes",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getrawmempool",
		Handler: "HandleGetRawMempool",
		Cmd:     "btcjson.GetRawMempoolCmd",
		Res:     "GetRawMempoolRes",
		ResType: "[]string",
	},
	{
		Method:  "getrawtransaction",
		Handler: "HandleGetRawTransaction",
		Cmd:     "btcjson.GetRawTransactionCmd",
		Res:     "GetRawTransactionRes",
		ResType: "string",
	},
	{
		Method:  "gettxout",
		Handler: "HandleGetTxOut",
		Cmd:     "btcjson.GetTxOutCmd",
		Res:     "GetTxOutRes",
		ResType: "string",
	},
	{
		Method:  "help",
		Handler: "HandleHelp",
		Cmd:     "btcjson.HelpCmd",
		Res:     "HelpRes",
		ResType: "string",
	},
	{
		Method:  "node",
		Handler: "HandleNode",
		Cmd:     "btcjson.NodeCmd",
		Res:     "NodeRes",
		ResType: "None",
	},
	{
		Method:  "ping",
		Handler: "HandlePing",
		Cmd:     "nil",
		Res:     "PingRes",
		ResType: "None",
	},
	{
		Method:  "searchrawtransactions",
		Handler: "HandleSearchRawTransactions",
		Cmd:     "btcjson.SearchRawTransactionsCmd",
		Res:     "SearchRawTransactionsRes",
		ResType: "[]btcjson.SearchRawTransactionsResult",
	},
	{
		Method:  "sendrawtransaction",
		Handler: "HandleSendRawTransaction",
		Cmd:     "btcjson.SendRawTransactionCmd",
		Res:     "SendRawTransactionRes",
		ResType: "None",
	},
	{
		Method:  "setgenerate",
		Handler: "HandleSetGenerate",
		Cmd:     "btcjson.SetGenerateCmd",
		Res:     "SetGenerateRes",
		ResType: "None",
	},
	{
		Method:  "stop",
		Handler: "HandleStop",
		Cmd:     "nil",
		Res:     "StopRes",
		ResType: "None",
	},
	{
		Method:  "restart",
		Handler: "HandleRestart",
		Cmd:     "nil",
		Res:     "RestartRes",
		ResType: "None",
	},
	{
		Method:  "resetchain",
		Handler: "HandleResetChain",
		Cmd:     "nil",
		Res:     "ResetChainRes",
		ResType: "None",
	},
	{
		Method:  "submitblock",
		Handler: "HandleSubmitBlock",
		Cmd:     "btcjson.SubmitBlockCmd",
		Res:     "SubmitBlockRes",
		ResType: "string",
	},
	{
		Method:  "uptime",
		Handler: "HandleUptime",
		Cmd:     "nil",
		Res:     "UptimeRes",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "validateaddress",
		Handler: "HandleValidateAddress",
		Cmd:     "btcjson.ValidateAddressCmd",
		Res:     "ValidateAddressRes",
		ResType: "btcjson.ValidateAddressChainResult",
	},
	{
		Method:  "verifychain",
		Handler: "HandleVerifyChain",
		Cmd:     "btcjson.VerifyChainCmd",
		Res:     "VerifyChainRes",
		ResType: "bool",
	},
	{
		Method:  "verifymessage",
		Handler: "HandleVerifyMessage",
		Cmd:     "btcjson.VerifyMessageCmd",
		Res:     "VerifyMessageRes",
		ResType: "bool",
	},
	{
		Method:  "version",
		Handler: "HandleVersion",
		Cmd:     "btcjson.VersionCmd",
		Res:     "VersionRes",
		ResType: "map[string]btcjson.VersionResult",
	},
}

func main() {
	t := template.Must(template.New("noderpc").Parse(NodeRPCHandlerTpl))
	if err := t.Execute(os.Stdout, handlers); log.L.Check(err) {
	}
}

var NodeRPCHandlerTpl = `package rpc

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type API struct {
	Ch     interface{}
	Params interface{}
}

RPCHandlersBeforeInit = map[string]CommandHandler{
{{range .}}	"{{ .Method }}":{ 
		{{ .Handler }}, make(chan API), func() API {
			return API{
				{{ .Cmd }}{},
				make(chan {{ .Res }}),
			}
		},
	}, 
{{end}}
}

type (
	None struct{} {{range .}}
	{{.Res}} struct {
		Res {{.ResType}}
		Err error
	}
	{{end}}
)
`
