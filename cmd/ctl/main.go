package ctl

import (
	"bufio"
	"bytes"
	js "encoding/json"
	"fmt"
	"github.com/stalker-loki/app/slog"
	"io"
	"os"
	"strings"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// HelpPrint is the uninitialized help print function
var HelpPrint = func() {
	fmt.Println("help has not been overridden")
}

// Main is the entry point for the pod.Ctl component
func Main(args []string, cx *conte.Xt) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	//
	method := args[0]
	usageFlags, err := btcjson.MethodUsageFlags(method)
	if err != nil {
		slog.Error(err)
		HelpPrint()
		os.Exit(1)
	}
	if usageFlags&unusableFlags != 0 {
		slog.Errorf("The '%s' command can only be used via websockets\n", method)
		HelpPrint()
		os.Exit(1)
	}
	// Convert remaining command line args to a slice of interface values to
	// be passed along as parameters to new command creation function.
	// Since some commands, such as submitblock,
	// can involve data which is too large for the Operating System to allow
	// as a normal command line parameter,
	// support using '-' as an argument to allow the argument to be read from
	// a stdin pipe.
	bio := bufio.NewReader(os.Stdin)
	params := make([]interface{}, 0, len(args[1:]))
	for _, arg := range args[1:] {
		if arg == "-" {
			param, err := bio.ReadString('\n')
			if err != nil && err != io.EOF {
				slog.Errorf("Failed to read data from stdin: %v\n", err)
				os.Exit(1)
			}
			if err == io.EOF && len(param) == 0 {
				slog.Errorf("Not enough lines provided on stdin")
				os.Exit(1)
			}
			param = strings.TrimRight(param, "\r\n")
			params = append(params, param)
			continue
		}
		params = append(params, arg)
	}
	// Attempt to create the appropriate command using the arguments provided
	// by the user.
	var cmd interface{}
	if cmd, err = btcjson.NewCmd(method, params...); slog.Check(err) {
		// Show the error along with its error code when it's a json.
		// BTCJSONError as it realistically will always be since the NewCmd function
		// is only supposed to return errors of that type.
		var ok bool
		var e btcjson.Error
		if e, ok = err.(btcjson.Error); ok {
			slog.Errorf("%s command: %v (code: %s)\n", method, err, e.ErrorCode)
			CommandUsage(method)
			os.Exit(1)
		}
		// The error is not a json.BTCJSONError and this really should not happen.
		// Nevertheless fall back to just showing the error if it should
		// happen due to a bug in the package.
		slog.Errorf("%s command: %v\n", method, err)
		CommandUsage(method)
		os.Exit(1)
	}
	// Marshal the command into a JSON-RPC byte slice in preparation for sending
	// it to the RPC server.
	var marshalledJSON []byte
	if marshalledJSON, err = btcjson.MarshalCmd(1, cmd); slog.Check(err) {
		os.Exit(1)
	}
	// Send the JSON-RPC request to the server using the user-specified
	// connection configuration.
	var result []byte
	if result, err = sendPostRequest(marshalledJSON, cx); slog.Check(err) {
		os.Exit(1)
	}
	// Choose how to display the result based on its type.
	strResult := string(result)
	switch {
	case strings.HasPrefix(strResult, "{") || strings.HasPrefix(strResult, "["):
		var dst bytes.Buffer
		if err = js.Indent(&dst, result, "", "  "); slog.Check(err) {
			slog.Errorf("Failed to format result: %v", err)
			os.Exit(1)
		}
		fmt.Println(dst.String())
	case strings.HasPrefix(strResult, `"`):
		var str string
		if err := js.Unmarshal(result, &str); err != nil {
			slog.Errorf("Failed to unmarshal result: %v", err)
			os.Exit(1)
		}
		fmt.Println(str)
	case strResult != "null":
		fmt.Println(strResult)
	}
}

// CommandUsage display the usage for a specific command.
func CommandUsage(method string) {
	usage, err := btcjson.MethodUsageText(method)
	if err != nil {
		slog.Error(err)
		// This should never happen since the method was already checked
		// before calling this function, but be safe.
		fmt.Println("Failed to obtain command usage:", err)
		return
	}
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", usage)
}
