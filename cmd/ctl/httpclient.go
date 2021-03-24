package ctl

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	js "encoding/json"
	"fmt"
	"github.com/p9c/pod/pkg/podcfg"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	
	"github.com/btcsuite/go-socks/socks"
	
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/btcjson"
)

// newHTTPClient returns a new HTTP client that is configured according to the proxy and TLS settings in the associated
// connection configuration.
func newHTTPClient(cfg *podcfg.Config) (client *http.Client, e error) {
	// Configure proxy if needed.
	var dial func(network, addr string) (net.Conn, error)
	if *cfg.Proxy != "" {
		proxy := &socks.Proxy{
			Addr:     *cfg.Proxy,
			Username: *cfg.ProxyUser,
			Password: *cfg.ProxyPass,
		}
		dial = func(network, addr string) (c net.Conn, e error) {
			if c, e = proxy.Dial(network, addr); E.Chk(e) {
				return nil, e
			}
			return c, nil
		}
	}
	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if *cfg.TLS && *cfg.RPCCert != "" {
		var pem []byte
		if pem, e = ioutil.ReadFile(*cfg.RPCCert); E.Chk(e) {
			return nil, e
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: *cfg.TLSSkipVerify,
		}
	}
	// Create and return the new HTTP client potentially configured with a proxy and TLS.
	client = &http.Client{
		Transport: &http.Transport{
			Dial:            dial,
			TLSClientConfig: tlsConfig,
		},
	}
	return
}

// sendPostRequest sends the marshalled JSON-RPC command using HTTP-POST mode to the server described in the passed
// config struct. It also attempts to unmarshal the response as a JSON-RPC response and returns either the result field
// or the error field depending on whether or not there is an error.
func sendPostRequest(marshalledJSON []byte, cx *pod.State) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if *cx.Config.TLS {
		protocol = "https"
	}
	serverAddr := *cx.Config.RPCConnect
	if *cx.Config.Wallet {
		serverAddr = *cx.Config.WalletServer
		_, _ = fmt.Fprintln(os.Stderr, "ctl: using wallet server", serverAddr)
	}
	url := protocol + "://" + serverAddr
	bodyReader := bytes.NewReader(marshalledJSON)
	var httpRequest *http.Request
	var e error
	if httpRequest, e = http.NewRequest("POST", url, bodyReader); E.Chk(e) {
		return nil, e
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")
	// Configure basic access authorization.
	httpRequest.SetBasicAuth(*cx.Config.Username, *cx.Config.Password)
	// Create the new HTTP client that is configured according to the user - specified options and submit the request.
	var httpClient *http.Client
	if httpClient, e = newHTTPClient(cx.Config); E.Chk(e) {
		return nil, e
	}
	var httpResponse *http.Response
	if httpResponse, e = httpClient.Do(httpRequest); E.Chk(e) {
		return nil, e
	}
	// Read the raw bytes and close the response.
	var respBytes []byte
	if respBytes, e = ioutil.ReadAll(httpResponse.Body); E.Chk(e) {
	}
	if e = httpResponse.Body.Close(); E.Chk(e) {
		e = fmt.Errorf("error reading json reply: %v", e)
		return nil, e
	}
	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is empty. This should not happen very often, but it's
		// better than showing nothing in case the target server has a poor implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf(
				"%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode),
			)
		}
		return nil, fmt.Errorf("%s", respBytes)
	}
	// Unmarshal the response.
	var resp btcjson.Response
	if e := js.Unmarshal(respBytes, &resp); E.Chk(e) {
		return nil, e
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
