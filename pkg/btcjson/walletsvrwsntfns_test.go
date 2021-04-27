package btcjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestWalletSvrWsNtfns tests all of the chain server websocket-specific notifications marshal and unmarshal into valid
// results include handling of optional fields being omitted in the marshalled command, while optional fields with
// defaults have the default assigned on unmarshalled commands.
func TestWalletSvrWsNtfns(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		newNtfn      func() (interface{}, error)
		staticNtfn   func() interface{}
		marshalled   string
		unmarshalled interface{}
	}{
		{
			name: "accountbalance",
			newNtfn: func() (interface{}, error) {
				return btcjson.NewCmd("accountbalance", "acct", 1.25, true)
			},
			staticNtfn: func() interface{} {
				return btcjson.NewAccountBalanceNtfn("acct", 1.25, true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"accountbalance","netparams":["acct",1.25,true],"id":null}`,
			unmarshalled: &btcjson.AccountBalanceNtfn{
				Account:   "acct",
				Balance:   1.25,
				Confirmed: true,
			},
		},
		{
			name: "podconnected",
			newNtfn: func() (interface{}, error) {
				return btcjson.NewCmd("podconnected", true)
			},
			staticNtfn: func() interface{} {
				return btcjson.NewPodConnectedNtfn(true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"podconnected","netparams":[true],"id":null}`,
			unmarshalled: &btcjson.PodConnectedNtfn{
				Connected: true,
			},
		},
		{
			name: "walletlockstate",
			newNtfn: func() (interface{}, error) {
				return btcjson.NewCmd("walletlockstate", true)
			},
			staticNtfn: func() interface{} {
				return btcjson.NewWalletLockStateNtfn(true)
			},
			marshalled: `{"jsonrpc":"1.0","method":"walletlockstate","netparams":[true],"id":null}`,
			unmarshalled: &btcjson.WalletLockStateNtfn{
				Locked: true,
			},
		},
		{
			name: "newtx",
			newNtfn: func() (interface{}, error) {
				return btcjson.NewCmd("newtx",
					"acct",
					`{"account":"acct","address":"1Address","category":"send","amount":1.5,"bip125-replaceable":"unknown","fee":0.0001,"confirmations":1,"trusted":true,"txid":"456","walletconflicts":[],"time":12345678,"timereceived":12345876,"vout":789,"otheraccount":"otheracct"}`,
				)
			},
			staticNtfn: func() interface{} {
				result := btcjson.ListTransactionsResult{
					Abandoned:       false,
					Account:         "acct",
					Address:         "1Address",
					Category:        "send",
					Amount:          1.5,
					Fee:             *btcjson.Float64(0.0001),
					Confirmations:   1,
					TxID:            "456",
					WalletConflicts: []string{},
					Time:            12345678,
					TimeReceived:    12345876,
					Trusted:         true,
					Vout:            789,
					OtherAccount:    "otheracct",
				}
				return btcjson.NewNewTxNtfn("acct", result)
			},
			marshalled: `{"jsonrpc":"1.0","method":"newtx","netparams":["acct",{"abandoned":false,"account":"acct","address":"1Address","amount":1.5,"bip125-replaceable":"unknown","category":"send","confirmations":1,"fee":0.0001,"time":12345678,"timereceived":12345876,"trusted":true,"txid":"456","vout":789,"walletconflicts":[],"otheraccount":"otheracct"}],"id":null}`,
			unmarshalled: &btcjson.NewTxNtfn{
				Account: "acct",
				Details: btcjson.ListTransactionsResult{
					Abandoned:       false,
					Account:         "acct",
					Address:         "1Address",
					Category:        "send",
					Amount:          1.5,
					Fee:             *btcjson.Float64(0.0001),
					Confirmations:   1,
					TxID:            "456",
					WalletConflicts: []string{},
					Time:            12345678,
					TimeReceived:    12345876,
					Trusted:         true,
					Vout:            789,
					OtherAccount:    "otheracct",
				},
			},
		},
	}
	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Marshal the notification as created by the new static creation function.  The ID is nil for notifications.
		marshalled, e := btcjson.MarshalCmd(nil, test.staticNtfn())
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
		// Ensure the notification is created without error via the generic new notification creation function.
		cmd, e := test.newNtfn()
		if e != nil {
			t.Errorf("Test #%d (%s) unexpected NewCmd error: %v ",
				i, test.name, e,
			)
		}
		// Marshal the notification as created by the generic new notification creation function. The ID is nil for
		// notifications.
		marshalled, e = btcjson.MarshalCmd(nil, cmd)
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
