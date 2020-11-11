package ctl

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	js "encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/btcsuite/go-socks/socks"

	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// newHTTPClient returns a new HTTP client that is configured according to the proxy and TLS settings in the associated
// connection configuration.
func newHTTPClient(cfg *pod.Config) (*http.Client, func(), error) {
	var dial func(ctx context.Context, network string, addr string) (net.Conn, error)
	ctx, cancel := context.WithCancel(context.Background())
	// Configure proxy if needed.
	if *cfg.Proxy != "" {
		proxy := &socks.Proxy{
			Addr:     *cfg.Proxy,
			Username: *cfg.ProxyUser,
			Password: *cfg.ProxyPass,
		}
		dial = func(_ context.Context, network string, addr string) (net.Conn, error) {
			c, err := proxy.Dial(network, addr)
			if err != nil {
				Error(err)
				return nil, err
			}
			go func() {
			out:
				for {
					select {
					case <-ctx.Done():
						if err := c.Close(); Check(err) {
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
	if *cfg.TLS && *cfg.RPCCert != "" {
		pem, err := ioutil.ReadFile(*cfg.RPCCert)
		if err != nil {
			Error(err)
			return nil, nil, err
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(pem)
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: *cfg.TLSSkipVerify,
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
func sendPostRequest(marshalledJSON []byte, cx *conte.Xt) ([]byte, error) {
	// Generate a request to the configured RPC server.
	protocol := "http"
	if *cx.Config.TLS {
		protocol = "https"
	}
	serverAddr := *cx.Config.RPCConnect
	if *cx.Config.Wallet {
		serverAddr = *cx.Config.WalletServer
		fmt.Println("using wallet server", serverAddr)
	}
	url := protocol + "://" + serverAddr
	bodyReader := bytes.NewReader(marshalledJSON)
	httpRequest, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		Error(err)
		return nil, err
	}
	httpRequest.Close = true
	httpRequest.Header.Set("Content-Type", "application/json")
	// Configure basic access authorization.
	httpRequest.SetBasicAuth(*cx.Config.Username, *cx.Config.Password)
	// Create the new HTTP client that is configured according to the user - specified options and submit the request.
	var httpClient *http.Client
	var cancel func()
	httpClient, cancel, err = newHTTPClient(cx.Config)
	if err != nil {
		Error(err)
		return nil, err
	}
	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		Error(err)
		return nil, err
	}
	// close connection
	cancel()
	// Read the raw bytes and close the response.
	respBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err := httpResponse.Body.Close(); Check(err) {
		err = fmt.Errorf("error reading json reply: %v", err)
		Error(err)
		return nil, err
	}
	// Handle unsuccessful HTTP responses
	if httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		// Generate a standard error to return if the server body is empty. This should not happen very often, but it's
		// better than showing nothing in case the target server has a poor implementation.
		if len(respBytes) == 0 {
			return nil, fmt.Errorf("%d %s", httpResponse.StatusCode,
				http.StatusText(httpResponse.StatusCode))
		}
		return nil, fmt.Errorf("%s", respBytes)
	}
	// Unmarshal the response.
	var resp btcjson.Response
	if err := js.Unmarshal(respBytes, &resp); err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
