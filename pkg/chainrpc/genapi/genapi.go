// +build generate

package main

import (
	"github.com/p9c/log"
	"os"
	"sort"
	"text/template"
)

const (
	RPCMapName = "RPCHandlers"
	Worker     = "CAPI"
)

type handler struct {
	Method, Handler, Cmd, ResType string
}

type handlersT []handler

func (h handlersT) Len() int           { return len(h) }
func (h handlersT) Less(i, j int) bool { return h[i].Method < h[j].Method }
func (h handlersT) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

var handlers = handlersT{
	{
		Method:  "addnode",
		Handler: "AddNode",
		Cmd:     "*btcjson.AddNodeCmd",
		ResType: "None",
	},
	{
		Method:  "createrawtransaction",
		Handler: "CreateRawTransaction",
		Cmd:     "*btcjson.CreateRawTransactionCmd",
		ResType: "string",
	},
	{
		Method:  "decoderawtransaction",
		Handler: "DecodeRawTransaction",
		Cmd:     "*btcjson.DecodeRawTransactionCmd",
		ResType: "btcjson.TxRawDecodeResult",
	},
	{
		Method:  "decodescript",
		Handler: "DecodeScript",
		Cmd:     "*btcjson.DecodeScriptCmd",
		ResType: "btcjson.DecodeScriptResult",
	},
	{
		Method:  "estimatefee",
		Handler: "EstimateFee",
		Cmd:     "*btcjson.EstimateFeeCmd",
		ResType: "float64",
	},
	{
		Method:  "generate",
		Handler: "Generate",
		Cmd:     "*None",
		ResType: "[]string",
	},
	{
		Method:  "getaddednodeinfo",
		Handler: "GetAddedNodeInfo",
		Cmd:     "*btcjson.GetAddedNodeInfoCmd",
		ResType: "[]btcjson.GetAddedNodeInfoResultAddr",
	},
	{
		Method:  "getbestblock",
		Handler: "GetBestBlock",
		Cmd:     "*None",
		ResType: "btcjson.GetBestBlockResult",
	},
	{
		Method:  "getbestblockhash",
		Handler: "GetBestBlockHash",
		Cmd:     "*None",
		ResType: "string",
	},
	{
		Method:  "getblock",
		Handler: "GetBlock",
		Cmd:     "*btcjson.GetBlockCmd",
		ResType: "btcjson.GetBlockVerboseResult",
	},
	{
		Method:  "getblockchaininfo",
		Handler: "GetBlockChainInfo",
		Cmd:     "*None",
		ResType: "btcjson.GetBlockChainInfoResult",
	},
	{
		Method:  "getblockcount",
		Handler: "GetBlockCount",
		Cmd:     "*None",
		ResType: "int64",
	},
	{
		Method:  "getblockhash",
		Handler: "GetBlockHash",
		Cmd:     "*btcjson.GetBlockHashCmd",
		ResType: "string",
	},
	{
		Method:  "getblockheader",
		Handler: "GetBlockHeader",
		Cmd:     "*btcjson.GetBlockHeaderCmd",
		ResType: "btcjson.GetBlockHeaderVerboseResult",
	},
	{
		Method:  "getblocktemplate",
		Handler: "GetBlockTemplate",
		Cmd:     "*btcjson.GetBlockTemplateCmd",
		ResType: "string",
	},
	{
		Method:  "getcfilter",
		Handler: "GetCFilter",
		Cmd:     "*btcjson.GetCFilterCmd",
		ResType: "string",
	},
	{
		Method:  "getcfilterheader",
		Handler: "GetCFilterHeader",
		Cmd:     "*btcjson.GetCFilterHeaderCmd",
		ResType: "string",
	},
	{
		Method:  "getconnectioncount",
		Handler: "GetConnectionCount",
		Cmd:     "*None",
		ResType: "int32",
	},
	{
		Method:  "getcurrentnet",
		Handler: "GetCurrentNet",
		Cmd:     "*None",
		ResType: "string",
	},
	{
		Method:  "getdifficulty",
		Handler: "GetDifficulty",
		Cmd:     "*btcjson.GetDifficultyCmd",
		ResType: "float64",
	},
	{
		Method:  "getgenerate",
		Handler: "GetGenerate",
		Cmd:     "*btcjson.GetHeadersCmd",
		ResType: "bool",
	},
	{
		Method:  "gethashespersec",
		Handler: "GetHashesPerSec",
		Cmd:     "*None",
		ResType: "float64",
	},
	{
		Method:  "getheaders",
		Handler: "GetHeaders",
		Cmd:     "*btcjson.GetHeadersCmd",
		ResType: "[]string",
	},
	{
		Method:  "getinfo",
		Handler: "GetInfo",
		Cmd:     "*None",
		ResType: "btcjson.InfoChainResult0",
	},
	{
		Method:  "getmempoolinfo",
		Handler: "GetMempoolInfo",
		Cmd:     "*None",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "getmininginfo",
		Handler: "GetMiningInfo",
		Cmd:     "*None",
		ResType: "btcjson.GetMiningInfoResult",
	},
	{
		Method:  "getnettotals",
		Handler: "GetNetTotals",
		Cmd:     "*None",
		ResType: "btcjson.GetNetTotalsResult",
	},
	{
		Method:  "getnetworkhashps",
		Handler: "GetNetworkHashPS",
		Cmd:     "*btcjson.GetNetworkHashPSCmd",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getpeerinfo",
		Handler: "GetPeerInfo",
		Cmd:     "*None",
		ResType: "[]btcjson.GetPeerInfoResult",
	},
	{
		Method:  "getrawmempool",
		Handler: "GetRawMempool",
		Cmd:     "*btcjson.GetRawMempoolCmd",
		ResType: "[]string",
	},
	{
		Method:  "getrawtransaction",
		Handler: "GetRawTransaction",
		Cmd:     "*btcjson.GetRawTransactionCmd",
		ResType: "string",
	},
	{
		Method:  "gettxout",
		Handler: "GetTxOut",
		Cmd:     "*btcjson.GetTxOutCmd",
		ResType: "string",
	},
	{
		Method:  "help",
		Handler: "Help",
		Cmd:     "*btcjson.HelpCmd",
		ResType: "string",
	},
	{
		Method:  "node",
		Handler: "Node",
		Cmd:     "*btcjson.NodeCmd",
		ResType: "None",
	},
	{
		Method:  "ping",
		Handler: "Ping",
		Cmd:     "*None",
		ResType: "None",
	},
	{
		Method:  "searchrawtransactions",
		Handler: "SearchRawTransactions",
		Cmd:     "*btcjson.SearchRawTransactionsCmd",
		ResType: "[]btcjson.SearchRawTransactionsResult",
	},
	{
		Method:  "sendrawtransaction",
		Handler: "SendRawTransaction",
		Cmd:     "*btcjson.SendRawTransactionCmd",
		ResType: "None",
	},
	{
		Method:  "setgenerate",
		Handler: "SetGenerate",
		Cmd:     "*btcjson.SetGenerateCmd",
		ResType: "None",
	},
	{
		Method:  "stop",
		Handler: "Stop",
		Cmd:     "*None",
		ResType: "None",
	},
	{
		Method:  "restart",
		Handler: "Restart",
		Cmd:     "*None",
		ResType: "None",
	},
	{
		Method:  "resetchain",
		Handler: "ResetChain",
		Cmd:     "*None",
		ResType: "None",
	},
	{
		Method:  "submitblock",
		Handler: "SubmitBlock",
		Cmd:     "*btcjson.SubmitBlockCmd",
		ResType: "string",
	},
	{
		Method:  "uptime",
		Handler: "Uptime",
		Cmd:     "*None",
		ResType: "btcjson.GetMempoolInfoResult",
	},
	{
		Method:  "validateaddress",
		Handler: "ValidateAddress",
		Cmd:     "*btcjson.ValidateAddressCmd",
		ResType: "btcjson.ValidateAddressChainResult",
	},
	{
		Method:  "verifychain",
		Handler: "VerifyChain",
		Cmd:     "*btcjson.VerifyChainCmd",
		ResType: "bool",
	},
	{
		Method:  "verifymessage",
		Handler: "VerifyMessage",
		Cmd:     "*btcjson.VerifyMessageCmd",
		ResType: "bool",
	},
	{
		Method:  "version",
		Handler: "Version",
		Cmd:     "*btcjson.VersionCmd",
		ResType: "map[string]btcjson.VersionResult",
	},
}

func main() {
	var NodeRPCHandlerTpl = `//go:generate go run -tags generate ./genapi/.
// generated; DO NOT EDIT

package chainrpc

import (
	"io"
	"net/rpc"
	"time"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/btcjson"
)

// API stores the channel, parameters and result values from calls via
// the channel
type API struct {
	Ch     interface{}
	Params interface{}
	Result interface{}
}

// CAPI is the central structure for configuration and access to a 
// net/rpc API access endpoint for this RPC API
type CAPI struct {
	Timeout time.Duration
	quit    qu.C
}

// NewCAPI returns a new CAPI 
func NewCAPI(quit qu.C, timeout ...time.Duration) (c *CAPI) {
	c = &CAPI{quit: quit}
	if len(timeout)>0 {
		c.Timeout = timeout[0]
	} else {
		c.Timeout = time.Second * 5
	}
	return 
}

// Wrappers around RPC calls
type CAPIClient struct {
	*rpc.Client
}

// New creates a new client for a kopach_worker.
// Note that any kind of connection can be used here, other than the StdConn
func NewCAPIClient(conn io.ReadWriteCloser) *CAPIClient {
	return &CAPIClient{rpc.NewClient(conn)}
}

type (
	// None means no parameters it is not checked so it can be nil
	None struct{} {{range .}}
	// {{.Handler}}Res is the result from a call to {{.Handler}}
	{{.Handler}}Res struct { Res *{{.ResType}}; Err error }{{end}}
)

// ` + RPCMapName + `BeforeInit are created first and are added to the main list 
// when the init runs.
//
// - Fn is the handler function
// 
// - Call is a channel carrying a struct containing parameters and error that is 
// listened to in RunAPI to dispatch the calls
// 
// - Result is a bundle of command parameters and a channel that the result will be sent 
// back on
//
// Get and save the Result function's return, and you can then call the call functions
// check, result and wait functions for asynchronous and synchronous calls to RPC functions
var ` + RPCMapName + `BeforeInit = map[string]CommandHandler{
{{range .}}	"{{.Method}}":{ 
		Fn: Handle{{.Handler}}, Call: make(chan API, 32), 
		Result: func() API { return API{Ch: make(chan {{.Handler}}Res)} }}, 
{{end}}
}

// API functions
//
// The functions here provide access to the RPC through a convenient set of functions
// generated for each call in the RPC API to request, check for, access the results and
// wait on results

{{range .}}
// {{.Handler}} calls the method with the given parameters
func (a API) {{.Handler}}(cmd {{.Cmd}}) (e error) {
	` + RPCMapName + `["{{.Method}}"].Call <-API{a.Ch, cmd, nil}
	return
}

// {{.Handler}}Chk checks if a new message arrived on the result channel and
// returns true if it does, as well as storing the value in the Result field
func (a API) {{.Handler}}Chk() (isNew bool) {
	select {
	case o := <-a.Ch.(chan {{.Handler}}Res):
		if o.Err != nil {
			a.Result = o.Err
		} else {
			a.Result = o.Res
		}
		isNew = true
	default:
	}
	return
}

// {{.Handler}}GetRes returns a pointer to the value in the Result field
func (a API) {{.Handler}}GetRes() (out *{{.ResType}}, e error) {
	out, _ = a.Result.(*{{.ResType}})
	e, _ = a.Result.(error)
	return 
}

// {{.Handler}}Wait calls the method and blocks until it returns or 5 seconds passes
func (a API) {{.Handler}}Wait(cmd {{.Cmd}}) (out *{{.ResType}}, e error) {
	` + RPCMapName + `["{{.Method}}"].Call <-API{a.Ch, cmd, nil}
	select {
	case <-time.After(time.Second*5):
		break
	case o := <-a.Ch.(chan {{.Handler}}Res):
		out, e = o.Res, o.Err
	}
	return
}
{{end}}

// RunAPI starts up the api handler server that receives rpc.API messages and runs the handler and returns the result
// Note that the parameters are type asserted to prevent the consumer of the API from sending wrong message types not
// because it's necessary since they are interfaces end to end
func RunAPI(server *Server, quit qu.C) {
	nrh := ` + RPCMapName + `
	go func() {
		D.Ln("starting up node cAPI")
		var e error
		var res interface{}
		for {
			select { {{range .}}
			case msg := <-nrh["{{.Method}}"].Call:
				if res, e = nrh["{{.Method}}"].
					Fn(server, msg.Params.({{.Cmd}}), nil); E.Chk(e) {
				}
				if r, ok := res.({{.ResType}}); ok { 
					msg.Ch.(chan {{.Handler}}Res) <-{{.Handler}}Res{&r, e} } {{end}}
			case <-quit.Wait():
				D.Ln("stopping wallet cAPI")
				return
			}
		}
	}()
}

// RPC API functions to use with net/rpc
{{range .}}
func (c *CAPI) {{.Handler}}(req {{.Cmd}}, resp {{.ResType}}) (e error) {
	nrh := ` + RPCMapName + `
	res := nrh["{{.Method}}"].Result()
	res.Params = req
	nrh["{{.Method}}"].Call <- res
	select {
	case resp = <-res.Ch.(chan {{.ResType}}):
	case <-time.After(c.Timeout):
	case <-c.quit.Wait():
	} 
	return 
}
{{end}}
// Client call wrappers for a CAPI client with a given Conn
{{range .}}
func (r *CAPIClient) {{.Handler}}(cmd ...{{.Cmd}}) (res {{.ResType}}, e error) {
	var c {{.Cmd}}
	if len(cmd) > 0 {
		c = cmd[0]
	}
	if e = r.Call("` + Worker + `.{{.Handler}}", c, &res); E.Chk(e) {
	}
	return
}
{{end}}
`
	log.SetLogLevel("trace")
	if fd, e := os.Create("rpchandlers.go"); E.Chk(e) {
		if fd, e := os.OpenFile("rpchandlers.go", os.O_RDWR|os.O_CREATE, 0755); E.Chk(e) {
		} else {
			defer fd.Close()
			t := template.Must(template.New("noderpc").Parse(NodeRPCHandlerTpl))
			sort.Sort(handlers)
			if e = t.Execute(fd, handlers); E.Chk(e) {
			}
		}
	} else {
		defer fd.Close()
		t := template.Must(template.New("noderpc").Parse(NodeRPCHandlerTpl))
		sort.Sort(handlers)
		if e = t.Execute(fd, handlers); E.Chk(e) {
		}
	}
}
