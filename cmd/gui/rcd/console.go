package rcd

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/go-socks/socks"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"github.com/p9c/pod/pkg/wallet/chain"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

func (r *RcVar) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	//params := split[1:]
	params := make([]interface{}, 0, len(split[1:]))

	log.INFO(len(params))
	//if len(params) < 1 {
	//	params = nil
	//}

	//cmd, err := btcjson.NewCmd(split[0], params...)
	//if err != nil {
	//	log.ERROR(err)
	//	// Show the error along with its error code when it's a json.
	//	// Error as it realistically will always be since the NewCmd function
	//	// is only supposed to return errors of that type.
	//	if jerr, ok := err.(btcjson.Error); ok {
	//		fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n",
	//			split[0], err, jerr.ErrorCode)
	//		commandUsage(split[0])
	//		os.Exit(1)
	//	}
	//	// The error is not a json.Error and this really should not happen.
	//	// Nevertheless fall back to just showing the error if it should
	//	// happen due to a bug in the package.
	//	fmt.Fprintf(os.Stderr, "%s command: %v\n", split[0], err)
	//	commandUsage(split[0])
	//	os.Exit(1)
	//}
	//// Marshal the command into a JSON-RPC byte slice in preparation for sending
	//// it to the RPC server.
	//marshalledJSON, err := btcjson.MarshalCmd(1, cmd)
	//if err != nil {
	//	log.ERROR(err)
	//	log.Println(err)
	//	os.Exit(1)
	//}
	//// Send the JSON-RPC request to the server using the user-specified
	//// connection configuration.
	////prevState := *r.cx.Config.Wallet
	//
	//result, err := sendPostRequest(marshalledJSON, r.cx)
	////*r.cx.Config.Wallet = prevState
	//if err != nil {
	//	log.ERROR(err)
	//	os.Exit(1)
	//}
	//// Choose how to display the result based on its type.
	//o = string(result)
	//

	//strResult := string(result)
	//switch {
	//case strings.HasPrefix(strResult, "{") || strings.HasPrefix(strResult, "["):
	//	var dst bytes.Buffer
	//	if err := js.Indent(&dst, result, "", "  "); err != nil {
	//		log.Printf("Failed to format result: %v", err)
	//		os.Exit(1)
	//	}
	//	log.Println(dst.String())
	//case strings.HasPrefix(strResult, `"`):
	//	var str string
	//	if err := js.Unmarshal(result, &str); err != nil {
	//		fmt.Fprintf(os.Stderr, "Failed to unmarshal result: %v",
	//			err)
	//		os.Exit(1)
	//	}
	//	log.Println(str)
	//case strResult != "null":
	//	log.Println(strResult)
	//}

	//c, err := btcjson.NewCmd(split[0], strings.Join(params, " "))
	c, err := btcjson.NewCmd(split[0], params...)
	if err != nil {
		o = fmt.Sprint(err)
	}
	handler, ok := legacy.RPCHandlers[split[0]]
	if ok {
		var out interface{}
		if handler.HandlerWithChain != nil {
			chainClient, err := startChainRPC(r.cx.Config, r.cx.ActiveNet, walletmain.ReadCAFile(r.cx.Config))
			if err != nil {
				log.ERROR(
					"unable to open connection to consensus RPC server:", err)
			}
			out, err = handler.HandlerWithChain(
				c,
				r.cx.WalletServer,
				chainClient)
			log.DEBUG("HandlerWithChain")
		}
		if handler.Handler != nil {
			out, err = handler.Handler(
				c,
				r.cx.WalletServer)
			if err != nil {
				log.ERROR(
					"unable to open connection to consensus RPC server:", err)
			}
			log.DEBUG("Handler")
		}
		if err != nil {
			o = fmt.Sprint(err)
		} else {
			if split[0] == "help" {
				o = out.(string)
			} else {
				j, _ := json.MarshalIndent(out, "", "  ")
				o = fmt.Sprint(string(j))
			}
		}
	} else {
		if split[0] == "" {
			o = ""
		} else {
			o = "Command does not exist"
		}
	}
	return
}

// startChainRPC opens a RPC client connection to a pod server for blockchain
// services.  This function uses the RPC options from the global config and
// there is no recovery in case the server is not available or if there is an
// authentication error.  Instead, all requests to the client will simply error.
func startChainRPC(config *pod.Config, activeNet *netparams.Params, certs []byte) (*chain.RPCClient, error) {
	log.TRACEF(
		"attempting RPC client connection to %v, TLS: %s",
		*config.RPCConnect, fmt.Sprint(*config.TLS),
	)
	rpcC, err := chain.NewRPCClient(activeNet, *config.RPCConnect,
		*config.Username, *config.Password, certs, !*config.TLS, 0)
	if err != nil {
		log.ERROR(err)
		return nil, err
	}
	err = rpcC.Start()
	return rpcC, err
}

// CommandUsage display the usage for a specific command.
func commandUsage(method string) {
	usage, err := btcjson.MethodUsageText(method)
	if err != nil {
		log.ERROR(err)
		// This should never happen since the method was already checked
		// before calling this function, but be safe.
		log.Println("Failed to obtain command usage:", err)
		return
	}
	log.Println("Usage:")
	log.Printf("  %s\n", usage)
}

// sendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode
// to the server described in the passed config struct.  It also attempts to
// unmarshal the response as a JSON-RPC response and returns either the result
// field or the error field depending on whether or not there is an error.
//func sendPostRequest(marshalledJSON []byte, cx *conte.Xt) ([]byte, error) {
//	// Generate a request to the configured RPC server.
//	protocol := "http"
//	if *cx.Config.TLS {
//		protocol = "https"
//	}
//	serverAddr := *cx.Config.RPCConnect
//	//serverAddr := *cx.Config.WalletServer
//	//if *cx.Config.Wallet {
//	//	serverAddr = *cx.Config.WalletServer
//	//	log.Println("using wallet server", serverAddr)
//	//}
//	url := protocol + "://" + serverAddr
//	bodyReader := bytes.NewReader(marshalledJSON)
//	httpRequest, err := http.NewRequest("POST", url, bodyReader)
//	if err != nil {
//		log.ERROR(err)
//		return nil, err
//	}
//	httpRequest.Close = true
//	httpRequest.Header.Set("Content-Type", "application/json")
//	// Configure basic access authorization.
//	httpRequest.SetBasicAuth(*cx.Config.Username, *cx.Config.Password)
//	// Create the new HTTP client that is configured according to the user
//	// - specified options and submit the request.
//	httpClient, err := newHTTPClient(cx.Config)
//	if err != nil {
//		log.ERROR(err)
//		return nil, err
//	}
//	httpResponse, err := httpClient.Do(httpRequest)
//	if err != nil {
//		log.ERROR(err)
//		return nil, err
//	}
//	// Read the raw bytes and close the response.
//	respBytes, err := ioutil.ReadAll(httpResponse.Body)
//	httpResponse.Body.Close()
//	if err != nil {
//		log.ERROR(err)
//		err = fmt.Errorf("error reading json reply: %v", err)
//		log.ERROR(err)
//		return nil, err
//	}
//	// Handle unsuccessful HTTP responses
//	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
//		// Generate a standard error to return if the server body is empty.
//		// This should not happen very often,
//		// but it's better than showing nothing in case the target server has
//		// a poor implementation.
//		if len(respBytes) == 0 {
//			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
//				http.StatusText(httpResponse.StatusCode))
//		}
//		return nil, fmt.Errorf("%s", respBytes)
//	}
//	// Unmarshal the response.
//	var resp btcjson.Response
//	if err := js.Unmarshal(respBytes, &resp); err != nil {
//		return nil, err
//	}
//	if resp.Error != nil {
//		return nil, resp.Error
//	}
//	return resp.Result, nil
//}

// newHTTPClient returns a new HTTP client that is configured according to the
// proxy and TLS settings in the associated connection configuration.
func newHTTPClient(cfg *pod.Config) (*http.Client, error) {
	// Configure proxy if needed.
	var dial func(network, addr string) (net.Conn, error)
	if *cfg.Proxy != "" {
		proxy := &socks.Proxy{
			Addr:     *cfg.Proxy,
			Username: *cfg.ProxyUser,
			Password: *cfg.ProxyPass,
		}
		dial = func(network, addr string) (net.Conn, error) {
			c, err := proxy.Dial(network, addr)
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
			return c, nil
		}
	}
	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if *cfg.TLS && *cfg.RPCCert != "" {
		pem, err := ioutil.ReadFile(*cfg.RPCCert)
		if err != nil {
			log.ERROR(err)
			return nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs: pool,
			// nolint
			InsecureSkipVerify: *cfg.TLSSkipVerify,
		}
	}
	// Create and return the new HTTP client potentially configured with a
	// proxy and TLS.
	client := http.Client{
		Transport: &http.Transport{
			Dial:            dial,
			TLSClientConfig: tlsConfig,
		},
	}
	return &client, nil
}
