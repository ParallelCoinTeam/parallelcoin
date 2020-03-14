package main

import (
	"os"
	"sort"
	"text/template"

	log "github.com/p9c/logi"
)

type handler struct {
	Method, Handler, Cmd, ResType string
}

type handlersT []handler

func (h handlersT) Len() int {
	return len(h)
}

func (h handlersT) Less(i, j int) bool {
	return h[i].Method < h[j].Method
}

func (h handlersT) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

var handlers = handlersT{
	{
		Method:  "addmultisigaddress",
		Handler: "AddMultiSigAddress",
		Cmd:     "btcjson.AddMultisigAddressCmd",
		ResType: "string",
	},
	{
		Method:  "createmultisig",
		Handler: "CreateMultiSig",
		Cmd:     "btcjson.CreateMultisigCmd",
		ResType: "btcjson.CreateMultiSigResult",
	},
	{
		Method:  "dumpprivkey",
		Handler: "DumpPrivKey",
		Cmd:     "btcjson.DumpPrivKeyCmd",
		ResType: "string",
	},
	{
		Method:  "getaccount",
		Handler: "GetAccount",
		Cmd:     "btcjson.GetAccountCmd",
		ResType: "string",
	},
	{
		Method:  "getaccountaddress",
		Handler: "GetAccountAddress",
		Cmd:     "btcjson.GetAccountAddressCmd",
		ResType: "string",
	},
	{
		Method:  "getaddressesbyaccount",
		Handler: "GetAddressesByAccount",
		Cmd:     "btcjson.GetAddressesByAccountCmd",
		ResType: "[]string",
	},
	{
		Method:  "getbalance",
		Handler: "GetBalance",
		Cmd:     "btcjson.GetBalanceCmd",
		ResType: "float64",
	},
	{
		Method:  "getbestblockhash",
		Handler: "GetBalance",
		Cmd:     "btcjson.GetBalanceCmd",
		ResType: "float64",
	},
	{
		Method:  "getbestblockhash",
		Handler: "GetBestBlockHash",
		Cmd:     "None",
		ResType: "string",
	},
	{
		Method:  "getblockcount",
		Handler: "GetBlockCount",
		Cmd:     "None",
		ResType: "int32",
	},
	{
		Method:  "getinfo",
		Handler: "GetInfo",
		Cmd:     "None",
		ResType: "btcjson.InfoWalletResult",
	},
	{
		Method:  "getnewaddress",
		Handler: "GetNewAddress",
		Cmd:     "btcjson.GetNewAddressCmd",
		ResType: "string",
	},
	{
		Method:  "getrawchangeaddress",
		Handler: "GetRawChangeAddress",
		Cmd:     "btcjson.GetRawChangeAddressCmd",
		ResType: "string",
	},
	{
		Method:  "getreceivedbyaccount",
		Handler: "GetReceivedByAccount",
		Cmd:     "btcjson.GetReceivedByAccountCmd",
		ResType: "float64",
	},
	{
		Method:  "getreceivedbyaddress",
		Handler: "GetReceivedByAddress",
		Cmd:     "btcjson.GetReceivedByAddressCmd",
		ResType: "float64",
	},
	{
		Method:  "gettransaction",
		Handler: "GetTransaction",
		Cmd:     "btcjson.GetTransactionCmd",
		ResType: "btcjson.GetTransactionResult",
	},
	{
		Method:  "help",
		Handler: "HelpWithChainRPC",
		Cmd:     "btcjson.HelpCmd",
		ResType: "string",
	},
	{
		Method:  "importprivkey",
		Handler: "ImportPrivKey",
		Cmd:     "btcjson.ImportPrivKeyCmd",
		ResType: "None",
	},
	{
		Method:  "keypoolrefill",
		Handler: "KeypoolRefill",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "listaccounts",
		Handler: "ListAccounts",
		Cmd:     "btcjson.ListAccountsCmd",
		ResType: "map[string]float64",
	},
	{
		Method:  "listlockunspent",
		Handler: "ListLockUnspent",
		Cmd:     "None",
		ResType: "[]btcjson.TransactionInput",
	},
	{
		Method:  "listreceivedbyaccount",
		Handler: "ListReceivedByAccount",
		Cmd:     "btcjson.ListReceivedByAccountCmd",
		ResType: "[]btcjson.ListReceivedByAccountResult",
	},
	{
		Method:  "listreceivedbyaddress",
		Handler: "ListReceivedByAddress",
		Cmd:     "btcjson.ListReceivedByAddressCmd",
		ResType: "btcjson.ListReceivedByAddressResult",
	},
	{
		Method:  "listsinceblock",
		Handler: "ListSinceBlock",
		Cmd:     "btcjson.ListSinceBlockCmd",
		ResType: "btcjson.ListSinceBlockResult",
	},
	{
		Method:  "listtransactions",
		Handler: "ListTransactions",
		Cmd:     "btcjson.ListTransactionsCmd",
		ResType: "[]btcjson.ListTransactionsResult",
	},
	{
		Method:  "listunspent",
		Handler: "ListUnspent",
		Cmd:     "btcjson.ListUnspentCmd",
		ResType: "[]btcjson.ListUnspentResult",
	},
	{
		Method:  "lockunspent",
		Handler: "LockUnspent",
		Cmd:     "btcjson.LockUnspentCmd",
		ResType: "bool",
	},
	{
		Method:  "sendfrom",
		Handler: "LockUnspent",
		Cmd:     "btcjson.LockUnspentCmd",
		ResType: "bool",
	},
	{
		Method:  "sendmany",
		Handler: "SendMany",
		Cmd:     "btcjson.SendManyCmd",
		ResType: "string",
	},
	{
		Method:  "sendtoaddress",
		Handler: "SendToAddress",
		Cmd:     "btcjson.SendToAddressCmd",
		ResType: "string",
	},
	{
		Method:  "settxfee",
		Handler: "SetTxFee",
		Cmd:     "btcjson.SetTxFeeCmd",
		ResType: "bool",
	},
	{
		Method:  "signmessage",
		Handler: "SignMessage",
		Cmd:     "btcjson.SignMessageCmd",
		ResType: "string",
	},
	{
		Method:  "signrawtransaction",
		Handler: "SignRawTransaction",
		Cmd:     "btcjson.SignRawTransactionCmd",
		ResType: "btcjson.SignRawTransactionResult",
	},
	{
		Method:  "validateaddress",
		Handler: "ValidateAddress",
		Cmd:     "btcjson.ValidateAddressCmd",
		ResType: "btcjson.ValidateAddressWalletResult",
	},
	{
		Method:  "verifymessage",
		Handler: "VerifyMessage",
		Cmd:     "btcjson.VerifyMessageCmd",
		ResType: "bool",
	},
	{
		Method:  "walletlock",
		Handler: "WalletLock",
		Cmd:     "None",
		ResType: "None",
	},
	{
		Method:  "walletpassphrase",
		Handler: "WalletPassphrase",
		Cmd:     "btcjson.WalletPassphraseCmd",
		ResType: "None",
	},
	{
		Method:  "walletpassphrasechange",
		Handler: "WalletPassphraseChange",
		Cmd:     "btcjson.WalletPassphraseChangeCmd",
		ResType: "None",
	},
	{
		Method:  "createnewaccount",
		Handler: "CreateNewAccount",
		Cmd:     "btcjson.CreateNewAccountCmd",
		ResType: "None",
	},
	{
		Method:  "getbestblock",
		Handler: "GetBestBlock",
		Cmd:     "None",
		ResType: "btcjson.GetBestBlockResult",
	},
	{
		Method:  "getunconfirmedbalance",
		Handler: "GetUnconfirmedBalance",
		Cmd:     "btcjson.GetUnconfirmedBalanceCmd",
		ResType: "float64",
	},
	{
		Method:  "listaddresstransactions",
		Handler: "GetUnconfirmedBalance",
		Cmd:     "btcjson.GetUnconfirmedBalanceCmd",
		ResType: "float64",
	},
	{
		Method:  "listaddresstransactions",
		Handler: "ListAddressTransactions",
		Cmd:     "btcjson.ListAddressTransactionsCmd",
		ResType: "[]btcjson.ListTransactionsResult",
	},
	{
		Method:  "listalltransactions",
		Handler: "ListAllTransactions",
		Cmd:     "btcjson.ListAllTransactionsCmd",
		ResType: "[]btcjson.ListTransactionsResult",
	},
	{
		Method:  "renameaccount",
		Handler: "RenameAccount",
		Cmd:     "btcjson.RenameAccountCmd",
		ResType: "None",
	},
	{
		Method:  "walletislocked",
		Handler: "WalletIsLocked",
		Cmd:     "None",
		ResType: "bool",
	},
	{
		Method:  "dropwallethistory",
		Handler: "HandleDropWalletHistory",
		Cmd:     "None",
		ResType: "string",
	},
}

func main() {
	log.L.SetLevel("trace", true, "pod")
	if fd, err := os.Create("../rpchandlers.go"); log.L.Check(err) {
	} else {
		defer fd.Close()
		t := template.Must(template.New("noderpc").Parse(NodeRPCHandlerTpl))
		sort.Sort(handlers)
		if err = t.Execute(fd, handlers); log.L.Check(err) {
		}
	}
}

const (
	RPCMapName = "RPCHandlers"
	Worker     = "CAPI"
)

var NodeRPCHandlerTpl = `// generated by go run github.com/p9c/pod/cmd/node/rpc/genapi/gen.go; DO NOT EDIT

package rpc

import (
	"io"
	"net/rpc"
	"time"

	log "github.com/p9c/logi"

	"github.com/p9c/pod/pkg/rpc/btcjson"
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
	quit    chan struct{}
}

// NewCAPI returns a new CAPI 
func NewCAPI(quit chan struct{}, timeout ...time.Duration) (c *CAPI) {
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
func (a API) {{.Handler}}(cmd {{.Cmd}}) (err error) {
	` + RPCMapName + `["{{.Method}}"].Call <- API{a.Ch, cmd, nil}
	return
}

// {{.Handler}}Check checks if a new message arrived on the result channel and 
// returns true if it does, as well as storing the value in the Result field
func (a API) {{.Handler}}Check() (isNew bool) {
	select {
	case o := <- a.Ch.(chan {{.Handler}}Res):
		a.Result = o.Res
		isNew = true
	default:
	}
	return
}

// {{.Handler}}GetRes returns a pointer to the value in the Result field
func (a API) {{.Handler}}GetRes() (out *{{.ResType}}) {
	ar := a.Result.({{.ResType}})
	return &ar
}

// {{.Handler}}Wait calls the method and blocks until it returns or 5 seconds passes
func (a API) {{.Handler}}Wait() (out *{{.ResType}}, err error) {
	select {
	case <-time.After(time.Second*5):
		break
	case o := <- a.Ch.(chan {{.Handler}}Res):
		out, err = o.Res, o.Err
	}
	return
}
{{end}}

// RunAPI starts up the api handler server that receives rpc.API messages and runs the handler and returns the result
// Note that the parameters are type asserted to prevent the consumer of the API from sending wrong message types not
// because it's necessary since they are interfaces end to end
func RunAPI(server *Server, quit chan struct{}) {
	nrh := ` + RPCMapName + `
	go func() {
		var err error
		var res interface{}
		for {
			select { {{range .}}
			case msg := <-nrh["{{.Method}}"].Call:
				if res, err = nrh["{{.Method}}"].
					Fn(server, msg.Params.({{.Cmd}}), nil); log.L.Check(err) {
				}
				if r, ok := res.({{.ResType}}); ok { 
					msg.Ch.(chan {{.Handler}}Res) <- {{.Handler}}Res{&r, err} } {{end}}
			case <-quit:
				return
			}
		}
	}()
}

// RPC API functions to use with net/rpc
{{range .}}
func (c *CAPI) {{.Handler}}(req *{{.Cmd}}, resp *{{.ResType}}) (err error) {
	nrh := ` + RPCMapName + `
	res := nrh["{{.Method}}"].Result()
	res.Params = req
	nrh["{{.Method}}"].Call <- res
	select {
	case *resp = <-res.Ch.(chan {{.ResType}}):
	case <-time.After(c.Timeout):
	case <- c.quit:
	} 
	return 
}
{{end}}
// Client call wrappers for a CAPI client with a given Conn
{{range .}}
func (r *CAPIClient) {{.Handler}}(cmd ...{{.Cmd}}) (res {{.ResType}}, err error) {
	var c {{.Cmd}}
	if len(cmd) > 0 {
		c = cmd[0]
	}
	if err = r.Call("` + Worker + `.{{.Handler}}", c, &res); log.L.Check(err) {
	}
	return
}
{{end}}
`
