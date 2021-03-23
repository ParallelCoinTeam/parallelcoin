package ctl

import (
	"bufio"
	"bytes"
	js "encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/pod"
	"io"
	"os"
	"strings"
	
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/ctl"
)

// HelpPrint is the uninitialized help print function
var HelpPrint = func() {
	fmt.Println("help has not been overridden")
}

// Main is the entry point for the pod.Ctl component
func Main(args []string, cx *pod.State) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	method := args[0]
	var usageFlags btcjson.UsageFlag
	var e error
	if usageFlags, e = btcjson.MethodUsageFlags(method); E.Chk(e) {
		_, _ = fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", method)
		HelpPrint()
		os.Exit(1)
	}
	if usageFlags&unusableFlags != 0 {
		_, _ = fmt.Fprintf(os.Stderr, "The '%s' command can only be used via websockets\n", method)
		HelpPrint()
		os.Exit(1)
	}
	// Convert remaining command line args to a slice of interface values to be passed along as parameters to new
	// command creation function. Since some commands, such as submitblock, can involve data which is too large for the
	// Operating System to allow as a normal command line parameter, support using '-' as an argument to allow the
	// argument to be read from a stdin pipe.
	bio := bufio.NewReader(os.Stdin)
	params := make([]interface{}, 0, len(args[1:]))
	for _, arg := range args[1:] {
		if arg == "-" {
			var param string
			if param, e = bio.ReadString('\n'); E.Chk(e) && e != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to read data from stdin: %v\n", e)
				os.Exit(1)
			}
			if e == io.EOF && len(param) == 0 {
				_, _ = fmt.Fprintln(os.Stderr, "Not enough lines provided on stdin")
				os.Exit(1)
			}
			param = strings.TrimRight(param, "\r\n")
			params = append(params, param)
			continue
		}
		params = append(params, arg)
	}
	var result []byte
	if result, e = ctl.Call(cx, *cx.Config.Wallet, method, params...); E.Chk(e) {
		return
	}
	// // Attempt to create the appropriate command using the arguments provided by the user.
	// cmd, e := btcjson.NewCmd(method, params...)
	// if e != nil  {
	// 	E.Ln(e)
	// 	// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
	// 	// since the NewCmd function is only supposed to return errors of that type.
	// 	if jerr, ok := err.(btcjson.BTCJSONError); ok {
	// 		fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n", method, e, jerr.ErrorCode)
	// 		CommandUsage(method)
	// 		os.Exit(1)
	// 	}
	// 	// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
	// 	// showing the error if it should happen due to a bug in the package.
	// 	fmt.Fprintf(os.Stderr, "%s command: %v\n", method, e)
	// 	CommandUsage(method)
	// 	os.Exit(1)
	// }
	// // Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
	// marshalledJSON, e := btcjson.MarshalCmd(1, cmd)
	// if e != nil  {
	// 	E.Ln(e)
	// 	fmt.Println(e)
	// 	os.Exit(1)
	// }
	// // Send the JSON-RPC request to the server using the user-specified connection configuration.
	// result, e := sendPostRequest(marshalledJSON, cx)
	// if e != nil  {
	// 	E.Ln(e)
	// 	os.Exit(1)
	// }
	// Choose how to display the result based on its type.
	strResult := string(result)
	switch {
	case strings.HasPrefix(strResult, "{") || strings.HasPrefix(strResult, "["):
		var dst bytes.Buffer
		if e = js.Indent(&dst, result, "", "  "); E.Chk(e) {
			fmt.Printf("Failed to format result: %v", e)
			os.Exit(1)
		}
		fmt.Println(dst.String())
	case strings.HasPrefix(strResult, `"`):
		var str string
		if e = js.Unmarshal(result, &str); E.Chk(e) {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to unmarshal result: %v", e)
			os.Exit(1)
		}
		fmt.Println(str)
	case strResult != "null":
		fmt.Println(strResult)
	}
}

// CommandUsage display the usage for a specific command.
func CommandUsage(method string) {
	var usage string
	var e error
	if usage, e = btcjson.MethodUsageText(method); E.Chk(e) {
		// This should never happen since the method was already checked before calling this function, but be safe.
		fmt.Println("Failed to obtain command usage:", e)
		return
	}
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", usage)
}
