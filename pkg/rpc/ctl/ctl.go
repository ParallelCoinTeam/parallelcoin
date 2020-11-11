package ctl

import (
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// Call uses settings in the context to call the method with the given parameters and returns the raw json bytes
func Call(cx *conte.Xt, method string, params ...interface{}) (result []byte, err error) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	var usageFlags btcjson.UsageFlag
	usageFlags, err = btcjson.MethodUsageFlags(method)
	if err != nil {
		Errorf("Unrecognized command '%s'\n", method)
		// HelpPrint()
		return
	}
	if usageFlags&unusableFlags != 0 {
		Errorf("The '%s' command can only be used via websockets\n", method)
		// HelpPrint()
		return
	}
	// Attempt to create the appropriate command using the arguments provided by the user.
	var cmd interface{}
	cmd, err = btcjson.NewCmd(method, params...)
	if err != nil {
		// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
		// since the NewCmd function is only supposed to return errors of that type.
		if jerr, ok := err.(btcjson.BTCJSONError); ok {
			Errorf("%s command: %v (code: %s)\n", method, err, jerr.ErrorCode)
			// CommandUsage(method)
			return
		}
		// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
		// showing the error if it should happen due to a bug in the package.
		Errorf("%s command: %v\n", method, err)
		// CommandUsage(method)
		return
	}
	// Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	var marshalledJSON []byte
	marshalledJSON, err = btcjson.MarshalCmd(1, cmd)
	if err != nil {
		Error(err)
		return
	}
	// Send the JSON-RPC request to the server using the user-specified connection configuration.
	result, err = sendPostRequest(marshalledJSON, cx)
	if err != nil {
		Error(err)
		return
	}
	return
}
