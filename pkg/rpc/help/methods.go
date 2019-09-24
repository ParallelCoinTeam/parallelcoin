//+build !generate

package rpchelp

import (
	"github.com/parallelcointeam/parallelcoin/pkg/rpc/json"
)

// HelpDescs contains the locale-specific help strings along with the locale.
var HelpDescs = []struct {
	Locale   string // Actual locale, e.g. en_US
	GoLocale string // Locale used in Go names, e.g. EnUS
	Descs    map[string]string
}{
	{"en_US", "EnUS", helpDescsEnUS}, // helpdescs_en_US.go
}

// Methods contains all methods and result types that help is generated for, for every locale.
var Methods = []struct {
	Method      string
	ResultTypes []interface{}
}{
	{"addmultisigaddress", returnsString},
	{"createmultisig", []interface{}{(*json.CreateMultiSigResult)(nil)}},
	{"dumpprivkey", returnsString},
	{"getaccount", returnsString},
	{"getaccountaddress", returnsString},
	{"getaddressesbyaccount", returnsStringArray},
	{"getbalance", append(returnsNumber, returnsNumber[0])},
	{"getbestblockhash", returnsString},
	{"getblockcount", returnsNumber},
	{"getinfo", []interface{}{(*json.InfoWalletResult)(nil)}},
	{"getnewaddress", returnsString},
	{"getrawchangeaddress", returnsString},
	{"getreceivedbyaccount", returnsNumber},
	{"getreceivedbyaddress", returnsNumber},
	{"gettransaction", []interface{}{(*json.GetTransactionResult)(nil)}},
	{"help", append(returnsString, returnsString[0])},
	{"importprivkey", nil},
	{"keypoolrefill", nil},
	{"listaccounts", []interface{}{(*map[string]float64)(nil)}},
	{"listlockunspent", []interface{}{(*[]json.TransactionInput)(nil)}},
	{"listreceivedbyaccount", []interface{}{(*[]json.ListReceivedByAccountResult)(nil)}},
	{"listreceivedbyaddress", []interface{}{(*[]json.ListReceivedByAddressResult)(nil)}},
	{"listsinceblock", []interface{}{(*json.ListSinceBlockResult)(nil)}},
	{"listtransactions", returnsLTRArray},
	{"listunspent", []interface{}{(*json.ListUnspentResult)(nil)}},
	{"lockunspent", returnsBool},
	{"sendfrom", returnsString},
	{"sendmany", returnsString},
	{"sendtoaddress", returnsString},
	{"settxfee", returnsBool},
	{"signmessage", returnsString},
	{"signrawtransaction", []interface{}{(*json.SignRawTransactionResult)(nil)}},
	{"validateaddress", []interface{}{(*json.ValidateAddressWalletResult)(nil)}},
	{"verifymessage", returnsBool},
	{"walletlock", nil},
	{"walletpassphrase", nil},
	{"walletpassphrasechange", nil},
	{"createnewaccount", nil},
	{"exportwatchingwallet", returnsString},
	{"getbestblock", []interface{}{(*json.GetBestBlockResult)(nil)}},
	{"getunconfirmedbalance", returnsNumber},
	{"listaddresstransactions", returnsLTRArray},
	{"listalltransactions", returnsLTRArray},
	{"renameaccount", nil},
	{"walletislocked", returnsBool},
}

// Common return types.
var returnsBool = []interface{}{(*bool)(nil)}

// Common return types.
var returnsLTRArray = []interface{}{(*[]json.ListTransactionsResult)(nil)}

// Common return types.
var returnsNumber = []interface{}{(*float64)(nil)}

// Common return types.
var returnsString = []interface{}{(*string)(nil)}

// Common return types.
var returnsStringArray = []interface{}{(*[]string)(nil)}
