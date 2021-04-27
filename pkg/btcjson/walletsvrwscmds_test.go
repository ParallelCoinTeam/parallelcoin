package btcjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestWalletSvrWsCmds tests all of the wallet server websocket-specific commands marshal and unmarshal into valid
// results include handling of optional fields being omitted in the marshalled command, while optional fields with
// defaults have the default assigned on unmarshalled commands.
func TestWalletSvrWsCmds(t *testing.T) {
	t.Parallel()
	testID := 1
	tests := []struct {
		name         string
		newCmd       func() (interface{}, error)
		staticCmd    func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "createencryptedwallet",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("createencryptedwallet", "pass")
			},
			staticCmd: func() interface{} {
				return btcjson.NewCreateEncryptedWalletCmd("pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"createencryptedwallet","netparams":["pass"],"id":1}`,
			unmarshalled: &btcjson.CreateEncryptedWalletCmd{Passphrase: "pass"},
		},
		{
			name: "exportwatchingwallet",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("exportwatchingwallet")
			},
			staticCmd: func() interface{} {
				return btcjson.NewExportWatchingWalletCmd(nil, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","netparams":[],"id":1}`,
			unmarshalled: &btcjson.ExportWatchingWalletCmd{
				Account:  nil,
				Download: btcjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional1",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("exportwatchingwallet", "acct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewExportWatchingWalletCmd(btcjson.String("acct"), nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","netparams":["acct"],"id":1}`,
			unmarshalled: &btcjson.ExportWatchingWalletCmd{
				Account:  btcjson.String("acct"),
				Download: btcjson.Bool(false),
			},
		},
		{
			name: "exportwatchingwallet optional2",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("exportwatchingwallet", "acct", true)
			},
			staticCmd: func() interface{} {
				return btcjson.NewExportWatchingWalletCmd(btcjson.String("acct"),
					btcjson.Bool(true),
				)
			},
			marshalled: `{"jsonrpc":"1.0","method":"exportwatchingwallet","netparams":["acct",true],"id":1}`,
			unmarshalled: &btcjson.ExportWatchingWalletCmd{
				Account:  btcjson.String("acct"),
				Download: btcjson.Bool(true),
			},
		},
		{
			name: "getunconfirmedbalance",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("getunconfirmedbalance")
			},
			staticCmd: func() interface{} {
				return btcjson.NewGetUnconfirmedBalanceCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","netparams":[],"id":1}`,
			unmarshalled: &btcjson.GetUnconfirmedBalanceCmd{
				Account: nil,
			},
		},
		{
			name: "getunconfirmedbalance optional1",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("getunconfirmedbalance", "acct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewGetUnconfirmedBalanceCmd(btcjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"getunconfirmedbalance","netparams":["acct"],"id":1}`,
			unmarshalled: &btcjson.GetUnconfirmedBalanceCmd{
				Account: btcjson.String("acct"),
			},
		},
		{
			name: "listaddresstransactions",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("listaddresstransactions", `["1Address"]`)
			},
			staticCmd: func() interface{} {
				return btcjson.NewListAddressTransactionsCmd([]string{"1Address"}, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","netparams":[["1Address"]],"id":1}`,
			unmarshalled: &btcjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   nil,
			},
		},
		{
			name: "listaddresstransactions optional1",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("listaddresstransactions", `["1Address"]`, "acct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewListAddressTransactionsCmd([]string{"1Address"},
					btcjson.String("acct"),
				)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listaddresstransactions","netparams":[["1Address"],"acct"],"id":1}`,
			unmarshalled: &btcjson.ListAddressTransactionsCmd{
				Addresses: []string{"1Address"},
				Account:   btcjson.String("acct"),
			},
		},
		{
			name: "listalltransactions",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("listalltransactions")
			},
			staticCmd: func() interface{} {
				return btcjson.NewListAllTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","netparams":[],"id":1}`,
			unmarshalled: &btcjson.ListAllTransactionsCmd{
				Account: nil,
			},
		},
		{
			name: "listalltransactions optional",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("listalltransactions", "acct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewListAllTransactionsCmd(btcjson.String("acct"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"listalltransactions","netparams":["acct"],"id":1}`,
			unmarshalled: &btcjson.ListAllTransactionsCmd{
				Account: btcjson.String("acct"),
			},
		},
		{
			name: "recoveraddresses",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("recoveraddresses", "acct", 10)
			},
			staticCmd: func() interface{} {
				return btcjson.NewRecoverAddressesCmd("acct", 10)
			},
			marshalled: `{"jsonrpc":"1.0","method":"recoveraddresses","netparams":["acct",10],"id":1}`,
			unmarshalled: &btcjson.RecoverAddressesCmd{
				Account: "acct",
				N:       10,
			},
		},
		{
			name: "walletislocked",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("walletislocked")
			},
			staticCmd: func() interface{} {
				return btcjson.NewWalletIsLockedCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"walletislocked","netparams":[],"id":1}`,
			unmarshalled: &btcjson.WalletIsLockedCmd{},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the command as created by the new static command creation function.
		marshalled, e := btcjson.MarshalCmd(testID, test.staticCmd())
		if e != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, e,
			)
			continue
		}
		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled,
			)
			continue
		}
		// Ensure the command is created without error via the generic new command creation function.
		cmd, e := test.newCmd()
		if e != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, e,
			)
		}
		// Marshal the command as created by the generic new command creation function.
		marshalled, e = btcjson.MarshalCmd(testID, cmd)
		if e != nil {
			t.Errorf("MarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, e,
			)
			continue
		}
		if !bytes.Equal(marshalled, []byte(test.marshalled)) {
			t.Errorf("Test #%d (%s) unexpected marshalled data - "+
				"got %s, want %s", i, test.name, marshalled,
				test.marshalled,
			)
			continue
		}
		var request btcjson.Request
		if e = json.Unmarshal(marshalled, &request); E.Chk(e) {
			t.Errorf("Test #%d (%s) unexpected error while "+
				"unmarshalling JSON-RPC request: %v", i,
				test.name, e,
			)
			continue
		}
		cmd, e = btcjson.UnmarshalCmd(&request)
		if e != nil {
			t.Errorf("UnmarshalCmd #%d (%s) unexpected error: %v", i,
				test.name, e,
			)
			continue
		}
		if !reflect.DeepEqual(cmd, test.unmarshalled) {
			t.Errorf("Test #%d (%s) unexpected unmarshalled command "+
				"- got %s, want %s", i, test.name,
				fmt.Sprintf("(%T) %+[1]v", cmd),
				fmt.Sprintf("(%T) %+[1]v\n", test.unmarshalled),
			)
			continue
		}
	}
}
