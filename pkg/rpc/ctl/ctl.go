package ctl

import (
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/pod"
	
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// Call uses settings in the context to call the method with the given parameters and returns the raw json bytes
func Call(cx *pod.State, wallet bool, method string, params ...interface{}) (result []byte, e error) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	var usageFlags btcjson.UsageFlag
	usageFlags, e = btcjson.MethodUsageFlags(method)
	if e != nil  {
		e = errors.New("Unrecognized command '" + method + "' : " + e.Error())
		// HelpPrint()
		return
	}
	if usageFlags&unusableFlags != 0 {
		E.F("The '%s' command can only be used via websockets\n", method)
		// HelpPrint()
		return
	}
	// Attempt to create the appropriate command using the arguments provided by the user.
	var cmd interface{}
	cmd, e = btcjson.NewCmd(method, params...)
	if e != nil  {
		// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
		// since the NewCmd function is only supposed to return errors of that type.
		if jerr, ok := e.(btcjson.GeneralError); ok {
			errText := fmt.Sprintf("%s command: %v (code: %s)\n", method, e, jerr.ErrorCode)
			e = errors.New(errText)
			// CommandUsage(method)
			return
		}
		// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
		// showing the error if it should happen due to a bug in the package.
		errText := fmt.Sprintf("%s command: %v\n", method, e)
		e = errors.New(errText)
		// CommandUsage(method)
		return
	}
	// Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	var marshalledJSON []byte
	marshalledJSON, e = btcjson.MarshalCmd(1, cmd)
	if e != nil  {
				return
	}
	// Send the JSON-RPC request to the server using the user-specified connection configuration.
	result, e = sendPostRequest(marshalledJSON, cx, wallet)
	if e != nil  {
				return
	}
	return
}
