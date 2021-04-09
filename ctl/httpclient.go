package ctl

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	js "encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	
	"github.com/btcsuite/go-socks/socks"
	
	"github.com/p9c/monorepo/pkg/btcjson"
)

// newHTTPClient returns a new HTTP client that is configured according to the proxy and TLS settings in the associated
// connection configuration.
func newHTTPClient(
	proxyAddress,
	proxyUser,
	proxyPass,
	rpcCert string,
	tlsEnabled,
	skipVerify bool,
) (client *http.Client, e error) {
	// Configure proxy if needed.
	var dial func(network, addr string) (net.Conn, error)
	if proxyAddress != "" {
		proxy := &socks.Proxy{
			Addr:     proxyAddress,
			Username: proxyUser,
			Password: proxyPass,
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
	if tlsEnabled && rpcCert != "" {
		var pem []byte
		if pem, e = ioutil.ReadFile(rpcCert); E.Chk(e) {
			return nil, e
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: skipVerify,
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
func sendPostRequest(marshalledJSON []byte, ctls, useWallet bool, nodeAddr, walletAddr, username, password string,
	proxyAddress,
	proxyUser,
	proxyPass,
	rpcCert string,
	tlsEnabled,
	skipVerify bool,
) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if ctls {
		protocol = "https"
	}
	serverAddr := nodeAddr
	if useWallet {
		serverAddr = walletAddr
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
	httpRequest.SetBasicAuth(username, password)
	// Create the new HTTP client that is configured according to the user - specified options and submit the request.
	var httpClient *http.Client
	if httpClient, e = newHTTPClient(
		proxyAddress,
		proxyUser,
		proxyPass,
		rpcCert,
		tlsEnabled,
		skipVerify,
	); E.Chk(e) {
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
