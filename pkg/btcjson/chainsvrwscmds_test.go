package btcjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestChainSvrWsCmds tests all of the chain server websocket-specific commands marshal and unmarshal into valid results
// include handling of optional fields being omitted in the marshalled command, while optional fields with defaults have
// the default assigned on unmarshalled commands.
func TestChainSvrWsCmds(t *testing.T) {
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
			name: "authenticate",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("authenticate", "user", "pass")
			},
			staticCmd: func() interface{} {
				return btcjson.NewAuthenticateCmd("user", "pass")
			},
			marshalled:   `{"jsonrpc":"1.0","method":"authenticate","netparams":["user","pass"],"id":1}`,
			unmarshalled: &btcjson.AuthenticateCmd{Username: "user", Passphrase: "pass"},
		},
		{
			name: "notifyblocks",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("notifyblocks")
			},
			staticCmd: func() interface{} {
				return btcjson.NewNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"notifyblocks","netparams":[],"id":1}`,
			unmarshalled: &btcjson.NotifyBlocksCmd{},
		},
		{
			name: "stopnotifyblocks",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("stopnotifyblocks")
			},
			staticCmd: func() interface{} {
				return btcjson.NewStopNotifyBlocksCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifyblocks","netparams":[],"id":1}`,
			unmarshalled: &btcjson.StopNotifyBlocksCmd{},
		},
		{
			name: "notifynewtransactions",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("notifynewtransactions")
			},
			staticCmd: func() interface{} {
				return btcjson.NewNotifyNewTransactionsCmd(nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","netparams":[],"id":1}`,
			unmarshalled: &btcjson.NotifyNewTransactionsCmd{
				Verbose: btcjson.Bool(false),
			},
		},
		{
			name: "notifynewtransactions optional",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("notifynewtransactions", true)
			},
			staticCmd: func() interface{} {
				return btcjson.NewNotifyNewTransactionsCmd(btcjson.Bool(true))
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifynewtransactions","netparams":[true],"id":1}`,
			unmarshalled: &btcjson.NotifyNewTransactionsCmd{
				Verbose: btcjson.Bool(true),
			},
		},
		{
			name: "stopnotifynewtransactions",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("stopnotifynewtransactions")
			},
			staticCmd: func() interface{} {
				return btcjson.NewStopNotifyNewTransactionsCmd()
			},
			marshalled:   `{"jsonrpc":"1.0","method":"stopnotifynewtransactions","netparams":[],"id":1}`,
			unmarshalled: &btcjson.StopNotifyNewTransactionsCmd{},
		},
		{
			name: "notifyreceived",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("notifyreceived", []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return btcjson.NewNotifyReceivedCmd([]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifyreceived","netparams":[["1Address"]],"id":1}`,
			unmarshalled: &btcjson.NotifyReceivedCmd{
				Addresses: []string{"1Address"},
			},
		},
		{
			name: "stopnotifyreceived",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("stopnotifyreceived", []string{"1Address"})
			},
			staticCmd: func() interface{} {
				return btcjson.NewStopNotifyReceivedCmd([]string{"1Address"})
			},
			marshalled: `{"jsonrpc":"1.0","method":"stopnotifyreceived","netparams":[["1Address"]],"id":1}`,
			unmarshalled: &btcjson.StopNotifyReceivedCmd{
				Addresses: []string{"1Address"},
			},
		},
		{
			name: "notifyspent",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("notifyspent", `[{"hash":"123","index":0}]`)
			},
			staticCmd: func() interface{} {
				ops := []btcjson.OutPoint{{Hash: "123", Index: 0}}
				return btcjson.NewNotifySpentCmd(ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"notifyspent","netparams":[[{"hash":"123","index":0}]],"id":1}`,
			unmarshalled: &btcjson.NotifySpentCmd{
				OutPoints: []btcjson.OutPoint{{Hash: "123", Index: 0}},
			},
		},
		{
			name: "stopnotifyspent",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("stopnotifyspent", `[{"hash":"123","index":0}]`)
			},
			staticCmd: func() interface{} {
				ops := []btcjson.OutPoint{{Hash: "123", Index: 0}}
				return btcjson.NewStopNotifySpentCmd(ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"stopnotifyspent","netparams":[[{"hash":"123","index":0}]],"id":1}`,
			unmarshalled: &btcjson.StopNotifySpentCmd{
				OutPoints: []btcjson.OutPoint{{Hash: "123", Index: 0}},
			},
		},
		{
			name: "rescan",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("rescan", "123", `["1Address"]`,
					`[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]`,
				)
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []btcjson.OutPoint{{
					Hash:  "0000000000000000000000000000000000000000000000000000000000000123",
					Index: 0,
				},
				}
				return btcjson.NewRescanCmd("123", addrs, ops, nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescan","netparams":["123",["1Address"],[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]],"id":1}`,
			unmarshalled: &btcjson.RescanCmd{
				BeginBlock: "123",
				Addresses:  []string{"1Address"},
				OutPoints:  []btcjson.OutPoint{{Hash: "0000000000000000000000000000000000000000000000000000000000000123",
												   Index: 0,
				},
				},
				EndBlock:   nil,
			},
		},
		{
			name: "rescan optional",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("rescan", "123", `["1Address"]`, `[{"hash":"123","index":0}]`, "456")
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []btcjson.OutPoint{{Hash: "123", Index: 0}}
				return btcjson.NewRescanCmd("123", addrs, ops, btcjson.String("456"))
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescan","netparams":["123",["1Address"],[{"hash":"123","index":0}],"456"],"id":1}`,
			unmarshalled: &btcjson.RescanCmd{
				BeginBlock: "123",
				Addresses:  []string{"1Address"},
				OutPoints:  []btcjson.OutPoint{{Hash: "123", Index: 0}},
				EndBlock:   btcjson.String("456"),
			},
		},
		{
			name: "loadtxfilter",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("loadtxfilter", false, `["1Address"]`,
					`[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]`,
				)
			},
			staticCmd: func() interface{} {
				addrs := []string{"1Address"}
				ops := []btcjson.OutPoint{{
					Hash:  "0000000000000000000000000000000000000000000000000000000000000123",
					Index: 0,
				},
				}
				return btcjson.NewLoadTxFilterCmd(false, addrs, ops)
			},
			marshalled: `{"jsonrpc":"1.0","method":"loadtxfilter","netparams":[false,["1Address"],[{"hash":"0000000000000000000000000000000000000000000000000000000000000123","index":0}]],"id":1}`,
			unmarshalled: &btcjson.LoadTxFilterCmd{
				Reload:    false,
				Addresses: []string{"1Address"},
				OutPoints: []btcjson.OutPoint{{Hash: "0000000000000000000000000000000000000000000000000000000000000123",
												  Index: 0,
				},
				},
			},
		},
		{
			name: "rescanblocks",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("rescanblocks",
					`["0000000000000000000000000000000000000000000000000000000000000123"]`,
				)
			},
			staticCmd: func() interface{} {
				blockhashes := []string{"0000000000000000000000000000000000000000000000000000000000000123"}
				return btcjson.NewRescanBlocksCmd(blockhashes)
			},
			marshalled: `{"jsonrpc":"1.0","method":"rescanblocks","netparams":[["0000000000000000000000000000000000000000000000000000000000000123"]],"id":1}`,
			unmarshalled: &btcjson.RescanBlocksCmd{
				BlockHashes: []string{"0000000000000000000000000000000000000000000000000000000000000123"},
			},
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
