package btcjson_test

import (
	"encoding/json"
	"fmt"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// This example demonstrates how to create and marshal a command into a JSON-RPC request.
func ExampleMarshalCmd() {
	// Create a new getblock command.  Notice the nil parameter indicates to use the default parameter for that fields.  This is a common pattern used in all of the New<Foo>Cmd functions in this package for optional fields.  Also, notice the call to json.Bool which is a convenience function for creating a pointer out of a primitive for optional parameters.
	blockHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	gbCmd := btcjson.NewGetBlockCmd(blockHash, btcjson.Bool(false), nil)
	// Marshal the command to the format suitable for sending to the RPC server.  Typically the client would increment the id here which is request so the response can be identified.
	id := 1
	marshalledBytes, e := btcjson.MarshalCmd(id, gbCmd)
	if e != nil {
		btcjson.E.Ln(e)
		return
	}
	// Display the marshalled command.  Ordinarily this would be sent across the wire to the RPC server, but for this example, just display it.
	fmt.Printf("%s\n", marshalledBytes)
	// Output: {"jsonrpc":"1.0","method":"getblock","netparams":["000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",false],"id":1}
}

// This example demonstrates how to unmarshal a JSON-RPC request and then unmarshal the concrete request into a concrete command.
func ExampleUnmarshalCmd() {
	// Ordinarily this would be read from the wire, but for this example, it is hard coded here for clarity.
	data := []byte(`{"jsonrpc":"1.0","method":"getblock","netparams":["000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",false],"id":1}`)
	// Unmarshal the raw bytes from the wire into a JSON-RPC request.
	var request btcjson.Request
	if e := json.Unmarshal(data, &request); E.Chk(e) {
		btcjson.E.Ln(e)
		return
	}
	// Typically there isn't any need to examine the request fields directly like this as the caller already knows what response to expect based on the command it sent.  However, this is done here to demonstrate why the unmarshal process is two steps.
	if request.ID == nil {
		btcjson.E.Ln("Unexpected notification")
		return
	}
	if request.Method != "getblock" {
		btcjson.E.Ln("Unexpected method")
		return
	}
	// Unmarshal the request into a concrete command.
	cmd, e := btcjson.UnmarshalCmd(&request)
	if e != nil {
		btcjson.E.Ln(e)
		return
	}
	// Type assert the command to the appropriate type.
	gbCmd, ok := cmd.(*btcjson.GetBlockCmd)
	if !ok {
		E.F("Incorrect command type: %T\n", cmd)
		return
	}
	// Display the fields in the concrete command.
	fmt.Println("Hash:", gbCmd.Hash)
	fmt.Println("Verbose:", *gbCmd.Verbose)
	fmt.Println("VerboseTx:", *gbCmd.VerboseTx)
	// Output:
	// Hash: 000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f
	// Verbose: false
	// VerboseTx: false
}

// This example demonstrates how to marshal a JSON-RPC response.
func ExampleMarshalResponse() {
	// Marshal a new JSON-RPC response.  For example, this is a response to a getblockheight request.
	marshalledBytes, e := btcjson.MarshalResponse(1, 350001, nil)
	if e != nil {
		btcjson.E.Ln(e)
		return
	}
	// Display the marshalled response.  Ordinarily this would be sent across the wire to the RPC client, but for this example, just display it.
	fmt.Printf("%s\n", marshalledBytes)
	// Output: {"result":350001,"error":null,"id":1}
}

// This example demonstrates how to unmarshal a JSON-RPC response and then unmarshal the result field in the response to a concrete type.
func Example_unmarshalResponse() {
	// Ordinarily this would be read from the wire, but for this example, it is hard coded here for clarity.  This is an example response to a getblockheight request.
	data := []byte(`{"result":350001,"error":null,"id":1}`)
	// Unmarshal the raw bytes from the wire into a JSON-RPC response.
	var response btcjson.Response
	if e := json.Unmarshal(data, &response); E.Chk(e) {
		btcjson.E.Ln("Malformed JSON-RPC response:", e)
		return
	}
	// Chk the response for an error from the server.  For example, the server might return an error if an invalid/unknown block hash is requested.
	if response.Error != nil {
		btcjson.E.Ln(response.Error)
		return
	}
	// Unmarshal the result into the expected type for the response.
	var blockHeight int32
	if e := json.Unmarshal(response.Result, &blockHeight); E.Chk(e) {
		E.F("Unexpected result type: %T\n", response.Result)
		return
	}
	fmt.Println("Block height:", blockHeight)
	// Output: Block height: 350001
}
