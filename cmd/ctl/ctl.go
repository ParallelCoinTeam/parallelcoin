package ctl

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/btcsuite/go-socks/socks"

	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pod/config"
)

// Call uses settings in the context to call the method with the given parameters and returns the raw json bytes
func Call(
	cx *config.Config, wallet bool, method string, params ...interface{},
) (result []byte, e error) {
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	var usageFlags btcjson.UsageFlag
	usageFlags, e = btcjson.MethodUsageFlags(method)
	if e != nil {
		e = errors.New("Unrecognized command '" + method + "' : " + e.Error())
		// HelpPrint()
		return
	}
	if usageFlags&btcjson.UnusableFlags != 0 {
		E.F("The '%s' command can only be used via websockets\n", method)
		// HelpPrint()
		return
	}
	// Attempt to create the appropriate command using the arguments provided by the user.
	var cmd interface{}
	cmd, e = btcjson.NewCmd(method, params...)
	if e != nil {
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
	if e != nil {
		return
	}
	// Send the JSON-RPC request to the server using the user-specified connection configuration.
	result, e = sendPostRequest(marshalledJSON, cx, wallet)
	if e != nil {
		return
	}
	return
}

// newHTTPClient returns a new HTTP client that is configured according to the proxy and TLS settings in the associated
// connection configuration.
func newHTTPClient(cfg *config.Config) (*http.Client, func(), error) {
	var dial func(ctx context.Context, network string,
		addr string) (net.Conn, error)
	ctx, cancel := context.WithCancel(context.Background())
	// Configure proxy if needed.
	if cfg.ProxyAddress.V() != "" {
		proxy := &socks.Proxy{
			Addr:     cfg.ProxyAddress.V(),
			Username: cfg.ProxyUser.V(),
			Password: cfg.ProxyPass.V(),
		}
		dial = func(_ context.Context, network string, addr string) (
			net.Conn, error,
		) {
			c, e := proxy.Dial(network, addr)
			if e != nil {
				return nil, e
			}
			go func() {
			out:
				for {
					select {
					case <-ctx.Done():
						if e := c.Close(); E.Chk(e) {
						}
						break out
					}
				}
			}()
			return c, nil
		}
	}
	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if cfg.ClientTLS.True() && cfg.RPCCert.V() != "" {
		pem, e := ioutil.ReadFile(cfg.RPCCert.V())
		if e != nil {
			cancel()
			return nil, nil, e
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: cfg.TLSSkipVerify.True(),
		}
	}
	// Create and return the new HTTP client potentially configured with a proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Proxy:                  nil,
			DialContext:            dial,
			TLSClientConfig:        tlsConfig,
			TLSHandshakeTimeout:    0,
			DisableKeepAlives:      false,
			DisableCompression:     false,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    0,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        0,
			ResponseHeaderTimeout:  0,
			ExpectContinueTimeout:  0,
			TLSNextProto:           nil,
			ProxyConnectHeader:     nil,
			MaxResponseHeaderBytes: 0,
			WriteBufferSize:        0,
			ReadBufferSize:         0,
			ForceAttemptHTTP2:      false,
		},
	}
	return &client, cancel, nil
}

// sendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode to the server described in the passed
// config struct. It also attempts to unmarshal the response as a JSON-RPC response and returns either the result field
// or the error field depending on whether or not there is an error.
func sendPostRequest(
	marshalledJSON []byte, cx *config.Config, wallet bool,
) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if cx.ClientTLS.True() {
		protocol = "https"
	}
	serverAddr := cx.RPCConnect.V()
	if wallet {
		serverAddr = cx.WalletServer.V()
		_, _ = fmt.Fprintln(os.Stderr, "using wallet server", serverAddr)
	}
	url := protocol + "://" + serverAddr
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, e := http.NewRequest("POST", url, bodyReader)
	if e != nil {
		return nil, e
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")
	// Configure basic access authorization.
	httpRequest.SetBasicAuth(cx.Username.V(), cx.Password.V())
	// T.Ln(cx.Username.V(), cx.Password.V())
	// Create the new HTTP client that is configured according to the user - specified options and submit the request.
	var httpClient *http.Client
	var cancel func()
	httpClient, cancel, e = newHTTPClient(cx)
	if e != nil {
		return nil, e
	}
	httpResponse, e := httpClient.Do(httpRequest)
	if e != nil {
		return nil, e
	}
	// close connection
	cancel()
	// Read the raw bytes and close the response.
	respBytes, e := ioutil.ReadAll(httpResponse.Body)
	if e := httpResponse.Body.Close(); E.Chk(e) {
		e = fmt.Errorf("error reading json reply: %v", e)
		return nil, e
	}
	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is empty. This should not happen very often, but it's
		// better than showing nothing in case the target server has a poor implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode),
			)
		}
		return nil, fmt.Errorf("%s", respBytes)
	}
	// Unmarshal the response.
	var resp btcjson.Response
	if e := json.Unmarshal(respBytes, &resp); E.Chk(e) {
		return nil, e
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

// ListCommands categorizes and lists all of the usable commands along with their one-line usage.
func ListCommands() (s string) {
	const (
		categoryChain uint8 = iota
		categoryWallet
		numCategories
	)
	// Get a list of registered commands and categorize and filter them.
	cmdMethods := btcjson.RegisteredCmdMethods()
	categorized := make([][]string, numCategories)
	for _, method := range cmdMethods {
		var e error
		var flags btcjson.UsageFlag
		if flags, e = btcjson.MethodUsageFlags(method); E.Chk(e) {
			continue
		}
		// Skip the commands that aren't usable from this utility.
		if flags&btcjson.UnusableFlags != 0 {
			continue
		}
		var usage string
		if usage, e = btcjson.MethodUsageText(method); E.Chk(e) {
			continue
		}
		// Categorize the command based on the usage flags.
		category := categoryChain
		if flags&btcjson.UFWalletOnly != 0 {
			category = categoryWallet
		}
		categorized[category] = append(categorized[category], usage)
	}
	// Display the command according to their categories.
	categoryTitles := make([]string, numCategories)
	categoryTitles[categoryChain] = "Chain Server Commands:"
	categoryTitles[categoryWallet] = "Wallet Server Commands (--wallet):"
	for category := uint8(0); category < numCategories; category++ {
		s += categoryTitles[category]
		s += "\n"
		for _, usage := range categorized[category] {
			s += "\t" + usage + "\n"
		}
		s += "\n"
	}
	return
}

// HelpPrint is the uninitialized help print function
var HelpPrint = func() {
	fmt.Println("help has not been overridden")
}

// CtlMain is the entry point for the pod.Ctl component
func CtlMain(cx *config.Config) {
	args := cx.ExtraArgs
	if len(args) < 1 {
		ListCommands()
		os.Exit(1)
	}
	// Ensure the specified method identifies a valid registered command and is one of the usable types.
	method := args[0]
	var usageFlags btcjson.UsageFlag
	var e error
	if usageFlags, e = btcjson.MethodUsageFlags(method); E.Chk(e) {
		_, _ = fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", method)
		HelpPrint()
		os.Exit(1)
	}
	if usageFlags&btcjson.UnusableFlags != 0 {
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
	if result, e = Call(cx, cx.UseWallet.True(), method, params...); E.Chk(e) {
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
		if e = json.Indent(&dst, result, "", "  "); E.Chk(e) {
			fmt.Printf("Failed to format result: %v", e)
			os.Exit(1)
		}
		fmt.Println(dst.String())
	case strings.HasPrefix(strResult, `"`):
		var str string
		if e = json.Unmarshal(result, &str); E.Chk(e) {
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
