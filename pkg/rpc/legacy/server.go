package legacy

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	js "encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/btcsuite/websocket"

	"github.com/parallelcointeam/parallelcoin/pkg/rpc/json"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet/chain"
)

type WebsocketClient struct {
	conn          *websocket.Conn
	authenticated bool
	remoteAddr    string
	allRequests   chan []byte
	responses     chan []byte
	quit          chan struct{} // closed on disconnect
	wg            sync.WaitGroup
}

func NewWebsocketClient(c *websocket.Conn, authenticated bool, remoteAddr string) *WebsocketClient {
	return &WebsocketClient{
		conn:          c,
		authenticated: authenticated,
		remoteAddr:    remoteAddr,
		allRequests:   make(chan []byte),
		responses:     make(chan []byte),
		quit:          make(chan struct{}),
	}
}
func (c *WebsocketClient) Send(b []byte) error {
	select {
	case c.responses <- b:
		return nil
	case <-c.quit:
		return errors.New("websocket client disconnected")
	}
}

// Server holds the items the RPC server may need to access (auth,
// config, shutdown, etc.)
type Server struct {
	HTTPServer   http.Server
	Wallet       *wallet.Wallet
	WalletLoader *wallet.Loader
	ChainClient  chain.Interface
	// handlerLookup       func(string) (requestHandler, bool)
	HandlerMutex        sync.Mutex
	Listeners           []net.Listener
	AuthSHA             [sha256.Size]byte
	Upgrader            websocket.Upgrader
	MaxPostClients      int64 // Max concurrent HTTP POST clients.
	MaxWebsocketClients int64 // Max concurrent websocket clients.
	WG                  sync.WaitGroup
	Quit                chan struct{}
	QuitMutex           sync.Mutex
	RequestShutdownChan chan struct{}
}

// JSONAuthFail sends a message back to the client if the http auth is rejected.
func JSONAuthFail(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="pod RPC"`)
	http.Error(w, "401 Unauthorized.", http.StatusUnauthorized)
}

// NewServer creates a new server for serving legacy RPC client connections,
// both HTTP POST and websocket.
func NewServer(opts *Options, walletLoader *wallet.Loader, listeners []net.Listener) *Server {
	serveMux := http.NewServeMux()
	const rpcAuthTimeoutSeconds = 10
	server := &Server{
		HTTPServer: http.Server{
			Handler: serveMux,
			// Timeout connections which don't complete the initial
			// handshake within the allowed timeframe.
			ReadTimeout: time.Second * rpcAuthTimeoutSeconds,
		},
		WalletLoader:        walletLoader,
		MaxPostClients:      opts.MaxPOSTClients,
		MaxWebsocketClients: opts.MaxWebsocketClients,
		Listeners:           listeners,
		// A hash of the HTTP basic auth string is used for a constant
		// time comparison.
		AuthSHA: sha256.Sum256(HTTPBasicAuth(opts.Username, opts.Password)),
		Upgrader: websocket.Upgrader{
			// Allow all origins.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		Quit:                make(chan struct{}),
		RequestShutdownChan: make(chan struct{}, 1),
	}
	serveMux.Handle("/", ThrottledFn(opts.MaxPOSTClients,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Connection", "close")
			w.Header().Set("Content-Type", "application/json")
			r.Close = true
			if err := server.CheckAuthHeader(r); err != nil {
				log <- cl.Wrn("unauthorized client connection attempt")
				JSONAuthFail(w)
				return
			}
			server.WG.Add(1)
			server.POSTClientRPC(w, r)
			server.WG.Done()
		}))
	serveMux.Handle("/ws", ThrottledFn(opts.MaxWebsocketClients,
		func(w http.ResponseWriter, r *http.Request) {
			authenticated := false
			switch server.CheckAuthHeader(r) {
			case nil:
				authenticated = true
			case ErrNoAuth:
				// nothing
			default:
				// If auth was supplied but incorrect, rather than simply
				// being missing, immediately terminate the connection.
				log <- cl.Wrn("disconnecting improperly authorized websocket client")
				JSONAuthFail(w)
				return
			}
			conn, err := server.Upgrader.Upgrade(w, r, nil)
			if err != nil {
				log <- cl.Warnf{
					"cannot websocket upgrade client %s: %v",
					r.RemoteAddr, err,
				}
				return
			}
			wsc := NewWebsocketClient(conn, authenticated, r.RemoteAddr)
			server.WebsocketClientRPC(wsc)
		}))
	for _, lis := range listeners {
		server.Serve(lis)
	}
	return server
}

// HTTPBasicAuth returns the UTF-8 bytes of the HTTP Basic authentication
// string:
//
//   "Basic " + base64(username + ":" + password)
func HTTPBasicAuth(username, password string) []byte {
	const header = "Basic "
	base64 := base64.StdEncoding
	b64InputLen := len(username) + len(":") + len(password)
	b64Input := make([]byte, 0, b64InputLen)
	b64Input = append(b64Input, username...)
	b64Input = append(b64Input, ':')
	b64Input = append(b64Input, password...)
	output := make([]byte, len(header)+base64.EncodedLen(b64InputLen))
	copy(output, header)
	base64.Encode(output[len(header):], b64Input)
	return output
}

// Serve serves HTTP POST and websocket RPC for the legacy JSON-RPC RPC server.
// This function does not block on lis.Accept.
func (s *Server) Serve(lis net.Listener) {
	s.WG.Add(1)
	go func() {
		log <- cl.Info{"wallet RPC server listening on ", lis.Addr(), cl.Ine()}
		err := s.HTTPServer.Serve(lis)
		log <- cl.Trace{"finished serving wallet RPC:", err, cl.Ine()}
		s.WG.Done()
	}()
}

// RegisterWallet associates the legacy RPC server with the wallet.  This
// function must be called before any wallet RPCs can be called by clients.
func (s *Server) RegisterWallet(w *wallet.Wallet) {
	s.HandlerMutex.Lock()
	s.Wallet = w
	s.HandlerMutex.Unlock()
}

// Stop gracefully shuts down the rpc server by stopping and disconnecting all clients, disconnecting the chain server connection, and closing the wallet's account files.  This blocks until shutdown completes.
func (s *Server) Stop() {
	s.QuitMutex.Lock()
	select {
	case <-s.Quit:
		s.QuitMutex.Unlock()
		return
	default:
	}
	// Stop the connected wallet and chain server, if any.
	s.HandlerMutex.Lock()
	wallet := s.Wallet
	chainClient := s.ChainClient
	s.HandlerMutex.Unlock()
	if wallet != nil {
		wallet.Stop()
	}
	if chainClient != nil {
		chainClient.Stop()
	}
	// Stop all the listeners.
	for _, listener := range s.Listeners {
		err := listener.Close()
		if err != nil {
			log <- cl.Errorf{
				"cannot close listener `%s`: %v %s",
				listener.Addr(), err, cl.Ine()}
		}
	}
	// Signal the remaining goroutines to stop.
	close(s.Quit)
	s.QuitMutex.Unlock()
	// First wait for the wallet and chain server to stop, if they
	// were ever set.
	if wallet != nil {
		wallet.WaitForShutdown()
	}
	if chainClient != nil {
		chainClient.WaitForShutdown()
	}
	// Wait for all remaining goroutines to exit.
	s.WG.Wait()
}

// SetChainServer sets the chain server client component needed to run a fully
// functional bitcoin wallet RPC server.  This can be called to enable RPC
// passthrough even before a loaded wallet is set, but the wallet's RPC client
// is preferred.
func (s *Server) SetChainServer(chainClient chain.Interface) {
	s.HandlerMutex.Lock()
	s.ChainClient = chainClient
	s.HandlerMutex.Unlock()
}

// HandlerClosure creates a closure function for handling requests of the given
// method.  This may be a request that is handled directly by btcwallet, or
// a chain server request that is handled by passing the request down to pod.
//
// NOTE: These handlers do not handle special cases, such as the authenticate
// method.  Each of these must be checked beforehand (the method is already
// known) and handled accordingly.
func (s *Server) HandlerClosure(request *json.Request) LazyHandler {
	s.HandlerMutex.Lock()
	// With the lock held, make copies of these pointers for the closure.
	wallet := s.Wallet
	chainClient := s.ChainClient
	if wallet != nil && chainClient == nil {
		chainClient = wallet.ChainClient()
		s.ChainClient = chainClient
	}
	s.HandlerMutex.Unlock()
	return LazyApplyHandler(request, wallet, chainClient)
}

// ErrNoAuth represents an error where authentication could not succeed
// due to a missing Authorization HTTP header.
var ErrNoAuth = errors.New("no auth")

// CheckAuthHeader checks the HTTP Basic authentication supplied by a client
// in the HTTP request r.  It errors with ErrNoAuth if the request does not
// contain the Authorization header, or another non-nil error if the
// authentication was provided but incorrect.
//
// This check is time-constant.
func (s *Server) CheckAuthHeader(r *http.Request) error {
	authhdr := r.Header["Authorization"]
	if len(authhdr) == 0 {
		return ErrNoAuth
	}
	authsha := sha256.Sum256([]byte(authhdr[0]))
	cmp := subtle.ConstantTimeCompare(authsha[:], s.AuthSHA[:])
	if cmp != 1 {
		return errors.New("bad auth")
	}
	return nil
}

// ThrottledFn wraps an http.HandlerFunc with throttling of concurrent active
// clients by responding with an HTTP 429 when the threshold is crossed.
func ThrottledFn(threshold int64, f http.HandlerFunc) http.Handler {
	return Throttled(threshold, f)
}

// Throttled wraps an http.Handler with throttling of concurrent active
// clients by responding with an HTTP 429 when the threshold is crossed.
func Throttled(threshold int64, h http.Handler) http.Handler {
	var active int64
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt64(&active, 1)
		defer atomic.AddInt64(&active, -1)
		if current-1 >= threshold {
			log <- cl.Warnf{
				"reached threshold of %d concurrent active clients", threshold,
			}
			http.Error(w, "429 Too Many Requests", 429)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// // sanitizeRequest returns a sanitized string for the request which may be
// // safely logged.  It is intended to strip private keys, passphrases, and any
// // other secrets from request parameters before they may be saved to a log file.
// func sanitizeRequest(// 	r *json.Request) string {
// 	// These are considered unsafe to log, so sanitize parameters.
// 	switch r.Method {
// 	case "encryptwallet", "importprivkey", "importwallet",
// 		"signrawtransaction", "walletpassphrase",
// 		"walletpassphrasechange":
// 		return fmt.Sprintf(`{"id":%v,"method":"%s","netparams":SANITIZED %d parameters}`,
// 			r.ID, r.Method, len(r.Params))
// 	}
// 	return fmt.Sprintf(`{"id":%v,"method":"%s","netparams":%v}`, r.ID,
// 		r.Method, r.Params)
// }

// IDPointer returns a pointer to the passed ID, or nil if the interface is nil.
// Interface pointers are usually a red flag of doing something incorrectly,
// but this is only implemented here to work around an oddity with json,
// which uses empty interface pointers for response IDs.
func IDPointer(id interface{}) (p *interface{}) {
	if id != nil {
		p = &id
	}
	return
}

// InvalidAuth checks whether a websocket request is a valid (parsable)
// authenticate request and checks the supplied username and passphrase
// against the server auth.
func (s *Server) InvalidAuth(req *json.Request) bool {
	cmd, err := json.UnmarshalCmd(req)
	if err != nil {
		return false
	}
	authCmd, ok := cmd.(*json.AuthenticateCmd)
	if !ok {
		return false
	}
	// Check credentials.
	login := authCmd.Username + ":" + authCmd.Passphrase
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
	authSha := sha256.Sum256([]byte(auth))
	return subtle.ConstantTimeCompare(authSha[:], s.AuthSHA[:]) != 1
}
func (s *Server) WebsocketClientRead(wsc *WebsocketClient) {
	for {
		_, request, err := wsc.conn.ReadMessage()
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log <- cl.Warnf{
					"websocket receive failed from client %s: %v",
					wsc.remoteAddr, err,
				}
			}
			close(wsc.allRequests)
			break
		}
		wsc.allRequests <- request
	}
}
func (s *Server) WebsocketClientRespond(wsc *WebsocketClient) {
	// A for-select with a read of the quit channel is used instead of a
	// for-range to provide clean shutdown.  This is necessary due to
	// WebsocketClientRead (which sends to the allRequests chan) not closing
	// allRequests during shutdown if the remote websocket client is still
	// connected.
out:
	for {
		select {
		case reqBytes, ok := <-wsc.allRequests:
			if !ok {
				// client disconnected
				break out
			}
			var req json.Request
			err := js.Unmarshal(reqBytes, &req)
			if err != nil {
				if !wsc.authenticated {
					// Disconnect immediately.
					break out
				}
				resp := MakeResponse(req.ID, nil,
					json.ErrRPCInvalidRequest)
				mresp, err := js.Marshal(resp)
				// We expect the marshal to succeed.  If it
				// doesn't, it indicates some non-marshalable
				// type in the response.
				if err != nil {
					panic(err)
				}
				err = wsc.Send(mresp)
				if err != nil {
					break out
				}
				continue
			}
			if req.Method == "authenticate" {
				if wsc.authenticated || s.InvalidAuth(&req) {
					// Disconnect immediately.
					break out
				}
				wsc.authenticated = true
				resp := MakeResponse(req.ID, nil, nil)
				// Expected to never fail.
				mresp, err := js.Marshal(resp)
				if err != nil {
					panic(err)
				}
				err = wsc.Send(mresp)
				if err != nil {
					break out
				}
				continue
			}
			if !wsc.authenticated {
				// Disconnect immediately.
				break out
			}
			switch req.Method {
			case "stop":
				resp := MakeResponse(req.ID,
					"wallet stopping.", nil)
				mresp, err := js.Marshal(resp)
				// Expected to never fail.
				if err != nil {
					panic(err)
				}
				err = wsc.Send(mresp)
				if err != nil {
					break out
				}
				s.RequestProcessShutdown()
				// break
			default:
				req := req // Copy for the closure
				f := s.HandlerClosure(&req)
				wsc.wg.Add(1)
				go func() {
					resp, jsonErr := f()
					mresp, err := json.MarshalResponse(req.ID, resp, jsonErr)
					if err != nil {
						log <- cl.Error{
							"unable to marshal response:", err, cl.Ine()}
					} else {
						_ = wsc.Send(mresp)
					}
					wsc.wg.Done()
				}()
			}
		case <-s.Quit:
			break out
		}
	}
	// allow client to disconnect after all handler goroutines are done
	wsc.wg.Wait()
	close(wsc.responses)
	s.WG.Done()
}
func (s *Server) WebsocketClientSend(wsc *WebsocketClient) {
	const deadline = 2 * time.Second
out:
	for {
		select {
		case response, ok := <-wsc.responses:
			if !ok {
				// client disconnected
				break out
			}
			err := wsc.conn.SetWriteDeadline(time.Now().Add(deadline))
			if err != nil {
				log <- cl.Warnf{
					"cannot set write deadline on client %s: %v",
					wsc.remoteAddr, err,
				}
			}
			err = wsc.conn.WriteMessage(websocket.TextMessage,
				response)
			if err != nil {
				log <- cl.Warnf{
					"failed websocket send to client %s: %v", wsc.remoteAddr, err,
				}
				break out
			}
		case <-s.Quit:
			break out
		}
	}
	close(wsc.quit)
	log <- cl.Info{"disconnected websocket client", wsc.remoteAddr, cl.Ine()}
	s.WG.Done()
}

// WebsocketClientRPC starts the goroutines to serve JSON-RPC requests over a
// websocket connection for a single client.
func (s *Server) WebsocketClientRPC(wsc *WebsocketClient) {
	log <- cl.Infof{
		"new websocket client %v %s", wsc.remoteAddr, cl.Ine(),
	}
	// Clear the read deadline set before the websocket hijacked
	// the connection.
	if err := wsc.conn.SetReadDeadline(time.Time{}); err != nil {
		log <- cl.Warn{
			"cannot remove read deadline:", err,
		}
	}
	// WebsocketClientRead is intentionally not run with the waitgroup
	// so it is ignored during shutdown.  This is to prevent a hang during
	// shutdown where the goroutine is blocked on a read of the
	// websocket connection if the client is still connected.
	go s.WebsocketClientRead(wsc)
	s.WG.Add(2)
	go s.WebsocketClientRespond(wsc)
	go s.WebsocketClientSend(wsc)
	<-wsc.quit
}

// MaxRequestSize specifies the maximum number of bytes in the request body
// that may be read from a client.  This is currently limited to 4MB.
const MaxRequestSize = 1024 * 1024 * 4

// POSTClientRPC processes and replies to a JSON-RPC client request.
func (s *Server) POSTClientRPC(w http.ResponseWriter, r *http.Request) {
	body := http.MaxBytesReader(w, r.Body, MaxRequestSize)
	rpcRequest, err := ioutil.ReadAll(body)
	if err != nil {
		// TODO: what if the underlying reader errored?
		http.Error(w, "413 Request Too Large.",
			http.StatusRequestEntityTooLarge)
		return
	}
	// First check whether wallet has a handler for this request's method.
	// If unfound, the request is sent to the chain server for further
	// processing.  While checking the methods, disallow authenticate
	// requests, as they are invalid for HTTP POST clients.
	var req json.Request
	err = js.Unmarshal(rpcRequest, &req)
	if err != nil {
		resp, err := json.MarshalResponse(req.ID, nil, json.ErrRPCInvalidRequest)
		if err != nil {
			log <- cl.Error{
				"Unable to marshal response:", err, cl.Ine()}
			http.Error(w, "500 Internal Server Error",
				http.StatusInternalServerError)
			return
		}
		_, err = w.Write(resp)
		if err != nil {
			log <- cl.Warn{
				"cannot write invalid request request to client:", err,
			}
		}
		return
	}
	// Create the response and error from the request.  Two special cases
	// are handled for the authenticate and stop request methods.
	var res interface{}
	var jsonErr *json.RPCError
	var stop bool
	switch req.Method {
	case "authenticate":
		// Drop it.
		return
	case "stop":
		stop = true
		res = "pod/wallet stopping"
	default:
		res, jsonErr = s.HandlerClosure(&req)()
	}
	// Marshal and send.
	mresp, err := json.MarshalResponse(req.ID, res, jsonErr)
	if err != nil {
		log <- cl.Error{
			"unable to marshal response:", err, cl.Ine()}
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(mresp)
	if err != nil {
		log <- cl.Warn{
			"unable to respond to client:", err,
		}
	}
	if stop {
		s.RequestProcessShutdown()
	}
}
func (s *Server) RequestProcessShutdown() {
	select {
	case s.RequestShutdownChan <- struct{}{}:
	default:
	}
}

// RequestProcessShutdownChan returns a channel that is sent to when an authorized
// client requests remote shutdown.
func (s *Server) RequestProcessShutdownChan() <-chan struct{} {
	return s.RequestShutdownChan
}
