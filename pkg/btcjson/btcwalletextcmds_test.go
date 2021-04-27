package btcjson_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TestBtcWalletExtCmds tests all of the btcwallet extended commands marshal and unmarshal into valid results include
// handling of optional fields being omitted in the marshalled command, while optional fields with defaults have the
// default assigned on unmarshalled commands.
func TestBtcWalletExtCmds(t *testing.T) {
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
			name: "createnewaccount",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("createnewaccount", "acct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewCreateNewAccountCmd("acct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"createnewaccount","netparams":["acct"],"id":1}`,
			unmarshalled: &btcjson.CreateNewAccountCmd{
				Account: "acct",
			},
		},
		{
			name: "dumpwallet",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("dumpwallet", "filename")
			},
			staticCmd: func() interface{} {
				return btcjson.NewDumpWalletCmd("filename")
			},
			marshalled: `{"jsonrpc":"1.0","method":"dumpwallet","netparams":["filename"],"id":1}`,
			unmarshalled: &btcjson.DumpWalletCmd{
				Filename: "filename",
			},
		},
		{
			name: "importaddress",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("importaddress", "1Address", "")
			},
			staticCmd: func() interface{} {
				return btcjson.NewImportAddressCmd("1Address", "", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","netparams":["1Address",""],"id":1}`,
			unmarshalled: &btcjson.ImportAddressCmd{
				Address: "1Address",
				Rescan:  btcjson.Bool(true),
			},
		},
		{
			name: "importaddress optional",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("importaddress", "1Address", "acct", false)
			},
			staticCmd: func() interface{} {
				return btcjson.NewImportAddressCmd("1Address", "acct", btcjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importaddress","netparams":["1Address","acct",false],"id":1}`,
			unmarshalled: &btcjson.ImportAddressCmd{
				Address: "1Address",
				Account: "acct",
				Rescan:  btcjson.Bool(false),
			},
		},
		{
			name: "importpubkey",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("importpubkey", "031234")
			},
			staticCmd: func() interface{} {
				return btcjson.NewImportPubKeyCmd("031234", nil)
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","netparams":["031234"],"id":1}`,
			unmarshalled: &btcjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: btcjson.Bool(true),
			},
		},
		{
			name: "importpubkey optional",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("importpubkey", "031234", false)
			},
			staticCmd: func() interface{} {
				return btcjson.NewImportPubKeyCmd("031234", btcjson.Bool(false))
			},
			marshalled: `{"jsonrpc":"1.0","method":"importpubkey","netparams":["031234",false],"id":1}`,
			unmarshalled: &btcjson.ImportPubKeyCmd{
				PubKey: "031234",
				Rescan: btcjson.Bool(false),
			},
		},
		{
			name: "importwallet",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("importwallet", "filename")
			},
			staticCmd: func() interface{} {
				return btcjson.NewImportWalletCmd("filename")
			},
			marshalled: `{"jsonrpc":"1.0","method":"importwallet","netparams":["filename"],"id":1}`,
			unmarshalled: &btcjson.ImportWalletCmd{
				Filename: "filename",
			},
		},
		{
			name: "renameaccount",
			newCmd: func() (interface{}, error) {
				return btcjson.NewCmd("renameaccount", "oldacct", "newacct")
			},
			staticCmd: func() interface{} {
				return btcjson.NewRenameAccountCmd("oldacct", "newacct")
			},
			marshalled: `{"jsonrpc":"1.0","method":"renameaccount","netparams":["oldacct","newacct"],"id":1}`,
			unmarshalled: &btcjson.RenameAccountCmd{
				OldAccount: "oldacct",
				NewAccount: "newacct",
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
