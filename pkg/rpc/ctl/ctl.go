package ctl

import (
	"fmt"
	"os"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// Call uses settings in the context to call the method with the given parameters and returns the raw json bytes
func Call(cx *conte.Xt, method string, params ...interface{}) (result []byte){
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	usageFlags, err := btcjson.MethodUsageFlags(method)
	if err != nil {
		Error(err)
		fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", method)
		// HelpPrint()
		// os.Exit(1)
		return
	}
	if usageFlags&unusableFlags != 0 {
		fmt.Fprintf(
			os.Stderr,
			"The '%s' command can only be used via websockets\n", method)
		// HelpPrint()
		// os.Exit(1)
		return
	}
	// Attempt to create the appropriate command using the arguments provided by the user.
	cmd, err := btcjson.NewCmd(method, params...)
	if err != nil {
		Error(err)
		// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
		// since the NewCmd function is only supposed to return errors of that type.
		if jerr, ok := err.(btcjson.BTCJSONError); ok {
			fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n", method, err, jerr.ErrorCode)
			// CommandUsage(method)
			// os.Exit(1)
			return
		}
		// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
		// showing the error if it should happen due to a bug in the package.
		fmt.Fprintf(os.Stderr, "%s command: %v\n", method, err)
		// CommandUsage(method)
		// os.Exit(1)
		return
	}
	// Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	marshalledJSON, err := btcjson.MarshalCmd(1, cmd)
	if err != nil {
		Error(err)
		// fmt.Println(err)
		// os.Exit(1)
		return
	}
	// Send the JSON-RPC request to the server using the user-specified connection configuration.
	result, err = sendPostRequest(marshalledJSON, cx)
	if err != nil {
		Error(err)
		// os.Exit(1)
		return
	}
	return
}
