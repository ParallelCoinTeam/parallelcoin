package chainrpc

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	amount2 "github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/chaincfg"

	"github.com/p9c/qu"

	"github.com/btcsuite/websocket"
	"golang.org/x/crypto/ripemd160"

	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

// Notification types

type NotificationBlockConnected block.Block
type NotificationBlockDisconnected block.Block
type NotificationRegisterAddr struct {
	WSC   *WSClient
	Addrs []string
}
type NotificationRegisterBlocks WSClient

// Notification control requests

type NotificationRegisterClient WSClient
type NotificationRegisterNewMempoolTxs WSClient
type NotificationRegisterSpent struct {
	WSC *WSClient
	OPs []*wire.OutPoint
}
type NotificationTxAcceptedByMempool struct {
	IsNew bool
	Tx    *util.Tx
}
type NotificationUnregisterAddr struct {
	WSC  *WSClient
	Addr string
}
type NotificationUnregisterBlocks WSClient
type NotificationUnregisterClient WSClient
type NotificationUnregisterNewMempoolTxs WSClient
type NotificationUnregisterSpent struct {
	WSC *WSClient
	OP  *wire.OutPoint
}
type RescanKeys struct {
	Fallbacks           map[string]struct{}
	PubKeyHashes        map[[ripemd160.Size]byte]struct{}
	ScriptHashes        map[[ripemd160.Size]byte]struct{}
	CompressedPubKeys   map[[33]byte]struct{}
	UncompressedPubKeys map[[65]byte]struct{}
	Unspent             map[wire.OutPoint]struct{}
}
type Semaphore qu.C

// WSClient provides an abstraction for handling a websocket client. The overall data flow is split into 3 main
// goroutines, a possible 4th goroutine for long-running operations (only started if request is made), and a websocket
// manager which is used to allow things such as broadcasting requested notifications to all connected websocket
// clients.
//
// Inbound messages are read via the inHandler goroutine and generally dispatched to their own handler. However, certain
// potentially long-running operations such as rescans, are sent to the asyncHander goroutine and are limited to one at
// a time.
//
// There are two outbound message types - one for responding to client requests and another for async notifications.
// Responses to client requests use SendMessage which employs a buffered channel thereby limiting the number of
// outstanding requests that can be made.
//
// Notifications are sent via QueueNotification which implements a queue via notificationQueueHandler to ensure sending
// notifications from other subsystems can't block.
//
// Ultimately, all messages are sent via the outHandler.
type WSClient struct {
	// ws client requires a mutex lock for access
	sync.Mutex
	// Server is the RPC Server that is servicing the client.
	Server *Server
	// Conn is the underlying websocket connection.
	Conn *websocket.Conn
	// Addr is the remote address of the client.
	Addr string
	// SessionID is a random ID generated for each client when connected. These IDs may be queried by a client using the
	// session RPC. A change to the session ID indicates that the client reconnected.
	SessionID    uint64
	AddrRequests map[string]struct{}
	// SpentRequests is a set of unspent Outpoints a wallet has requested notifications for when they are spent by a
	// processed transaction. Owned by the notification manager.
	SpentRequests map[wire.OutPoint]struct{}
	// FilterData is the new generation transaction filter backported from github.com/decred/dcrd for the new backported
	// `loadtxfilter` and `rescanblocks` methods.
	FilterData *WSClientFilter
	// Networking infrastructure.
	ServiceRequestSem Semaphore
	NtfnChan          chan []byte
	SendChan          chan WSResponse
	Quit              qu.C
	WG                sync.WaitGroup
	// Disconnected indicated whether or not the websocket client is Disconnected.
	Disconnected bool
	// Authenticated specifies whether a client has been Authenticated and therefore is allowed to communicated over the
	// websocket.
	Authenticated bool
	// IsAdmin specifies whether a client may change the state of the Server; false means its access is only to the
	// limited set of RPC calls.
	IsAdmin bool
	// VerboseTxUpdates specifies whether a client has requested verbose information about all new transactions.
	VerboseTxUpdates bool
	// AddrRequests is a set of addresses the caller has requested to be notified about. It is maintained here so all
	// requests can be removed when a wallet disconnects. Owned by the notification manager.
}

// WSClientFilter tracks relevant addresses for each websocket client for the `rescanblocks` extension. It is modified
// by the `loadtxfilter` command. NOTE: This extension was ported from github.com/decred/dcrd
type WSClientFilter struct {
	mu sync.Mutex
	// Implemented fast paths for address lookup.
	PubKeyHashes        map[[ripemd160.Size]byte]struct{}
	ScriptHashes        map[[ripemd160.Size]byte]struct{}
	CompressedPubKeys   map[[33]byte]struct{}
	UncompressedPubKeys map[[65]byte]struct{}
	// A fallback address lookup map in case a fast path doesn't exist. Only exists for completeness. If using this
	// shows up in a profile, there's a good chance a fast path should be added.
	OtherAddresses map[string]struct{}
	// Outpoints of Unspent outputs.
	Unspent map[wire.OutPoint]struct{}
}

// WSCommandHandler describes a callback function used to handle a specific command.
type WSCommandHandler func(*WSClient, interface{}) (interface{}, error)

// WSNtfnMgr is a connection and notification manager used for websockets. It allows websocket clients to register for
// notifications they are interested in.
//
// When an event happens elsewhere in the code such as transactions being added to the memory pool or block
// connects/disconnects, the notification manager is provided with the relevant details needed to figure out which
// websocket clients need to be notified based on what they have registered for and notifies them accordingly.
//
// It is also used to keep track of all connected websocket clients.
type WSNtfnMgr struct {
	// Server is the RPC Server the notification manager is associated with.
	Server *Server
	// QueueNotification queues a notification for handling.
	QueueNotification chan interface{}
	// NotificationMsgs feeds notificationHandler with notifications and client (un)registration requests from a queue
	// as well as registration and unregistering requests from clients.
	NotificationMsgs chan interface{}
	// Access channel for current number of connected clients.
	NumClients chan int
	// Shutdown handling
	WG   sync.WaitGroup
	Quit qu.C
}

// WSResponse houses a message to send to a connected websocket client as well as a channel to reply on when the message
// is sent.
type WSResponse struct {
	Msg      []byte
	DoneChan chan bool
}

const (
	// WebsocketSendBufferSize is the number of elements the send channel can queue before blocking. Note that this only
	// applies to requests handled directly in the websocket client input handler or the async handler since
	// notifications have their own queuing mechanism independent of the send channel buffer.
	WebsocketSendBufferSize = 50
)

// ErrClientQuit describes the error where a client send is not processed due to the client having already been
// disconnected or dropped.
var ErrClientQuit = errors.New("client quit")

// ErrRescanReorg defines the error that is returned when an unrecoverable reorganize is detected during a rescan.
var ErrRescanReorg = btcjson.RPCError{
	Code:    btcjson.ErrRPCDatabase,
	Message: "Reorganize",
}

// TimeZeroVal is simply the zero value for a time.Time and is used to avoid creating multiple instances.
var TimeZeroVal = time.Time{}

// WSHandlers maps RPC command strings to appropriate websocket handler functions. This is set by init because help
// references WSHandlers and thus causes a dependency loop.
var WSHandlers map[string]WSCommandHandler

var WSHandlersBeforeInit = map[string]WSCommandHandler{
	"loadtxfilter":              HandleLoadTxFilter,
	"help":                      HandleWebsocketHelp,
	"notifyblocks":              HandleNotifyBlocks,
	"notifynewtransactions":     HandleNotifyNewTransactions,
	"notifyreceived":            HandleNotifyReceived,
	"notifyspent":               HandleNotifySpent,
	"session":                   HandleSession,
	"stopnotifyblocks":          HandleStopNotifyBlocks,
	"stopnotifynewtransactions": HandleStopNotifyNewTransactions,
	"stopnotifyspent":           HandleStopNotifySpent,
	"stopnotifyreceived":        HandleStopNotifyReceived,
	"rescan":                    HandleRescan,
	"rescanblocks":              HandleRescanBlocks,
}

// UnspentSlice returns a slice of currently-unspent outpoints for the rescan lookup keys. This is primarily intended to
// be used to register outpoints for continuous notifications after a rescan has completed.
func (r *RescanKeys) UnspentSlice() []*wire.OutPoint {
	ops := make([]*wire.OutPoint, 0, len(r.Unspent))
	for op := range r.Unspent {
		opCopy := op
		ops = append(ops, &opCopy)
	}
	return ops
}

// WebsocketHandler handles a new websocket client by creating a new wsClient, starting it, and blocking until the
// connection closes. Since it blocks, it must be run in a separate goroutine.
//
// It should be invoked from the websocket Server handler which runs each new connection in a new goroutine thereby
// satisfying the requirement.
func (s *Server) WebsocketHandler(
	conn *websocket.Conn, remoteAddr string,
	authenticated bool, isAdmin bool,
) {
	// Clear the read deadline that was set before the websocket hijacked the connection.
	e := conn.SetReadDeadline(TimeZeroVal)
	if e != nil {
		D.Ln(e)
	}
	// Limit max number of websocket clients.
	T.Ln("new websocket client", remoteAddr)
	if s.NtfnMgr.GetNumClients()+1 > s.Config.RPCMaxWebsockets.V() {
		I.F(
			"max websocket clients exceeded [%d] - disconnecting client"+
				" %s",
			s.Config.RPCMaxWebsockets,
			remoteAddr,
		)
		if e = conn.Close(); E.Chk(e) {
		}
		return
	}
	// Create a new websocket client to handle the new websocket connection and wait for it to shutdown.
	//
	// Once it has shutdown (and hence disconnected), remove it and any notifications it registered for.
	client, e := NewWebsocketClient(s, conn, remoteAddr, authenticated, isAdmin)
	if e != nil {
		E.F("failed to serve client %s: %v %s", remoteAddr, e)
		if e := conn.Close(); E.Chk(e) {
		}
		return
	}
	s.NtfnMgr.AddClient(client)
	client.Start()
	client.WaitForShutdown()
	s.NtfnMgr.RemoveClient(client)
	T.Ln("disconnected websocket client", remoteAddr)
}

// Disconnect disconnects the websocket client.
func (c *WSClient) Disconnect() {
	c.Lock()
	defer c.Unlock()
	// Nothing to do if already disconnected.
	if c.Disconnected {
		return
	}
	T.Ln("disconnecting websocket client", c.Addr)
	c.Quit.Q()
	if e := c.Conn.Close(); E.Chk(e) {
	}
	c.Disconnected = true
}

// IsDisconnected returns whether or not the websocket client is disconnected.
func (c *WSClient) IsDisconnected() bool {
	c.Lock()
	isDisconnected := c.Disconnected
	c.Unlock()
	return isDisconnected
}

// QueueNotification queues the passed notification to be sent to the websocket client.
//
// This function, as the name implies, is only intended for notifications since it has additional logic to prevent other
// subsystems, such as the memory pool and block manager, from blocking even when the send channel is full.
//
// If the client is in the process of shutting down, this function returns ErrClientQuit. This is intended to be checked
// by long-running notification handlers to stop processing if there is no more work needed to be done.
func (c *WSClient) QueueNotification(marshalledJSON []byte) (e error) {
	// Don't queue the message if disconnected.
	if c.IsDisconnected() {
		return ErrClientQuit
	}
	c.NtfnChan <- marshalledJSON
	return nil
}

// SendMessage sends the passed json to the websocket client. It is backed by a buffered channel, so it will not block
// until the send channel is full.
//
// Note however that QueueNotification must be used for sending async notifications instead of the this function. This
// approach allows a limit to the number of outstanding requests a client can make without preventing or blocking on
// async notifications.
func (c *WSClient) SendMessage(marshalledJSON []byte, doneChan chan bool) {
	// Don't send the message if disconnected.
	if c.IsDisconnected() {
		if doneChan != nil {
			doneChan <- false
		}
		return
	}
	c.SendChan <- WSResponse{Msg: marshalledJSON, DoneChan: doneChan}
}

// Start begins processing input and output messages.
func (c *WSClient) Start() {
	T.Ln("starting websocket client", c.Addr)
	// Start processing input and output.
	c.WG.Add(3)
	go c.InHandler()
	go c.NotificationQueueHandler()
	go c.OutHandler()
}

// WaitForShutdown blocks until the websocket client goroutines are stopped and the connection is closed.
func (c *WSClient) WaitForShutdown() {
	c.WG.Wait()
}

// InHandler handles all incoming messages for the websocket connection. It must be run as a goroutine.
func (c *WSClient) InHandler() {
out:
	for {
		// Break out of the loop once the quit channel has been closed. Use a non-blocking select here so we fall
		// through otherwise.
		select {
		case <-c.Quit.Wait():
			break out
		default:
		}
		_, msg, e := c.Conn.ReadMessage()
		if e != nil {
			// Log the error if it's not due to disconnecting.
			if e != io.EOF && e != io.ErrUnexpectedEOF {
				T.F(
					"websocket receive error from %s: %v",
					c.Addr, e,
				)
			}
			break out
		}
		var request btcjson.Request
		e = json.Unmarshal(msg, &request)
		if e != nil {
			if !c.Authenticated {
				break out
			}
			jsonErr := &btcjson.RPCError{
				Code:    btcjson.ErrRPCParse.Code,
				Message: "Failed to parse request: " + e.Error(),
			}
			reply, e := CreateMarshalledReply(nil, nil, jsonErr)
			if e != nil {
				E.Ln(
					"failed to marshal parse failure reply:", e,
				)
				continue
			}
			c.SendMessage(reply, nil)
			continue
		}
		// The JSON-RPC 1.0 spec defines that notifications must have their "id" set to null and states that
		// notifications do not have a response.
		//
		// A JSON-RPC 2.0 notification is a request with "json-rpc":"2.0", and without an "id" member. The specification
		// states that notifications must not be responded to. JSON-RPC 2.0 permits the null value as a valid request
		// id, therefore such requests are not notifications.
		//
		// Bitcoin Core serves requests with "id":null or even an absent "id", and responds to such requests with
		// "id":null in the response.
		//
		// Pod does not respond to any request without and "id" or "id":null, regardless the indicated JSON-RPC protocol
		// version unless RPC quirks are enabled. With RPC quirks enabled, such requests will be responded to if the
		// reqeust does not indicate JSON-RPC version.
		//
		// RPC quirks can be enabled by the user to avoid compatibility issues with software relying on Core's behavior.
		if request.ID == nil && !(c.Server.Config.RPCQuirks.True() && request.Jsonrpc == "") {
			if !c.Authenticated {
				break out
			}
			continue
		}
		cmd := ParseCmd(&request)
		if cmd.Err != nil {
			if !c.Authenticated {
				break out
			}
			reply, e := CreateMarshalledReply(cmd.ID, nil, cmd.Err)
			if e != nil {
				E.F("failed to marshal parse failure reply:", e)
				continue
			}
			c.SendMessage(reply, nil)
			continue
		}
		T.F("received command <%s> from %s", cmd.Method, c.Addr)
		// T.Ln(interrupt.GoroutineDump())
		// Chk auth. The client is immediately disconnected if the first request of an unauthenticated websocket
		// client is not the authenticate request, an authenticate request is received when the client is already
		// authenticated, or incorrect authentication credentials are provided in the request.
		switch authCmd, ok := cmd.Cmd.(*btcjson.AuthenticateCmd); {
		case c.Authenticated && ok:
			W.F(
				"websocket client %s is already authenticated", c.Addr,
			)
			break out
		case !c.Authenticated && !ok:
			W.F("unauthenticated websocket message received")
			break out
		case !c.Authenticated:
			// Chk credentials.
			login := authCmd.Username + ":" + authCmd.Passphrase
			auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(login))
			authSha := sha256.Sum256([]byte(auth))
			cmp := subtle.ConstantTimeCompare(authSha[:], c.Server.AuthSHA[:])
			limitcmp := subtle.ConstantTimeCompare(authSha[:], c.Server.LimitAuthSHA[:])
			if cmp != 1 && limitcmp != 1 {
				W.Ln("authentication failure from", c.Addr)
				break out
			}
			c.Authenticated = true
			c.IsAdmin = cmp == 1
			// Marshal and send response.
			reply, e := CreateMarshalledReply(cmd.ID, nil, nil)
			if e != nil {
				E.Ln("failed to marshal authenticate reply:", e)
				continue
			}
			c.SendMessage(reply, nil)
			continue
		}
		// Chk if the client is using limited RPC credentials and error when not authorized to call this RPC.
		if !c.IsAdmin {
			if _, ok := RPCLimited[request.Method]; !ok {
				jsonErr := &btcjson.RPCError{
					Code:    btcjson.ErrRPCInvalidParams.Code,
					Message: "limited user not authorized for this method",
				}
				// Marshal and send response.
				reply, e := CreateMarshalledReply(request.ID, nil, jsonErr)
				if e != nil {
					E.Ln("failed to marshal parse failure reply:", e)
					continue
				}
				c.SendMessage(reply, nil)
				continue
			}
		}
		// Asynchronously handle the request. A semaphore is used to limit the number of concurrent requests currently
		// being serviced. If the semaphore can not be acquired, simply wait until a request finished before reading the
		// next RPC request from the websocket client.
		//
		// This could be a little fancier by timing out and failing when it takes too long to service the request, but
		// if that is done, the read of the next request should not be blocked by this semaphore, otherwise the next
		// request will be read and will probably sit here for another few seconds before timing out as well. This will
		// cause the total timeout duration for later requests to be much longer than the check here would imply.
		//
		// If a timeout is added, the semaphore acquiring should be moved inside of the new goroutine with a select
		// statement that also reads a time.After channel. This will unblock the read of the next request from the
		// websocket client and allow many requests to be waited on concurrently.
		c.ServiceRequestSem.Acquire()
		go func() {
			c.ServiceRequest(cmd)
			c.ServiceRequestSem.Release()
		}()
	}
	// Ensure the connection is closed.
	c.Disconnect()
	c.WG.Done()
	T.Ln("websocket client input handler done for", c.Addr)
}

// NotificationQueueHandler handles the queuing of outgoing notifications for the websocket client. This runs as a muxer
// for various sources of input to ensure that queuing up notifications to be sent will not block. Otherwise, slow
// clients could bog down the other systems (such as the mempool or block manager) which are queuing the data. The data
// is passed on to outHandler to actually be written. It must be run as a goroutine.
func (c *WSClient) NotificationQueueHandler() {
	ntfnSentChan := make(chan bool, 1) // nonblocking sync
	// pendingNtfns is used as a queue for notifications that are ready to be sent once there are no outstanding
	// notifications currently being sent.
	//
	// The waiting flag is used over simply checking for items in the pending list to ensure cleanup knows what has and
	// hasn't been sent to the outHandler.
	//
	// Currently no special cleanup is needed, however if something like a done channel is added to notifications in the
	// future, not knowing what has and hasn't been sent to the outHandler (and thus who should respond to the done
	// channel) would be problematic without using this approach.
	pendingNtfns := list.New()
	waiting := false
out:
	for {
		select {
		// This channel is notified when a message is being queued to be sent across the network socket.
		//
		// It will either send the message immediately if a send is not already in progress, or queue the message to
		// be sent once the other pending messages are sent.
		case msg := <-c.NtfnChan:
			if !waiting {
				c.SendMessage(msg, ntfnSentChan)
			} else {
				pendingNtfns.PushBack(msg)
			}
			waiting = true
			// This channel is notified when a notification has been sent across the
			// network socket.
		case <-ntfnSentChan:
			// No longer waiting if there are no more messages in the pending
			// messages queue.
			next := pendingNtfns.Front()
			if next == nil {
				waiting = false
				continue
			}
			// Notify the outHandler about the next item to asynchronously send.
			msg := pendingNtfns.Remove(next).([]byte)
			c.SendMessage(msg, ntfnSentChan)
		case <-c.Quit.Wait():
			break out
		}
	}
	// Drain any wait channels before exiting so nothing is left waiting around to send.
cleanup:
	for {
		select {
		case <-c.NtfnChan:
		case <-ntfnSentChan:
		default:
			break cleanup
		}
	}
	c.WG.Done()
	T.Ln("websocket client notification queue handler done for", c.Addr)
}

// OutHandler handles all outgoing messages for the websocket connection. It must be run as a goroutine.
//
// It uses a buffered channel to serialize output messages while allowing the sender to continue running asynchronously.
//
// It must be run as a goroutine.
func (c *WSClient) OutHandler() {
out:
	for {
		// Send any messages ready for send until the quit channel is closed.
		select {
		case r := <-c.SendChan:
			e := c.Conn.WriteMessage(websocket.TextMessage, r.Msg)
			if e != nil {
				c.Disconnect()
				break out
			}
			if r.DoneChan != nil {
				r.DoneChan <- true
			}
		case <-c.Quit.Wait():
			break out
		}
	}
	// Drain any wait channels before exiting so nothing is left waiting around to send.
cleanup:
	for {
		select {
		case r := <-c.SendChan:
			if r.DoneChan != nil {
				r.DoneChan <- false
			}
		default:
			break cleanup
		}
	}
	c.WG.Done()
	T.Ln("websocket client output handler done for", c.Addr)
}

// ServiceRequest services a parsed RPC request by looking up and executing the appropriate RPC handler. The response is
// marshalled and sent to the websocket client.
func (c *WSClient) ServiceRequest(r *ParsedRPCCmd) {
	var (
		result interface{}
		e      error
	)
	// Lookup the websocket extension for the command and if it doesn't exist fallback to handling the command as a
	// standard command.
	wsHandler, ok := WSHandlers[r.Method]
	if ok {
		result, e = wsHandler(c, r.Cmd)
	} else {
		result, e = c.Server.StandardCmdResult(r, nil)
	}
	reply, e := CreateMarshalledReply(r.ID, result, e)
	if e != nil {
		E.F("failed to marshal reply for <%s> command: %v", r.Method, e)
		return
	}
	c.SendMessage(reply, nil)
}

// AddAddress adds an address to a wsClientFilter, treating it correctly based on the type of address passed as an
// argument. NOTE: This extension was ported from github.com/decred/dcrd
func (f *WSClientFilter) AddAddress(a btcaddr.Address) {
	switch a := a.(type) {
	case *btcaddr.PubKeyHash:
		f.PubKeyHashes[*a.Hash160()] = struct{}{}
		return
	case *btcaddr.ScriptHash:
		f.ScriptHashes[*a.Hash160()] = struct{}{}
		return
	case *btcaddr.PubKey:
		serializedPubKey := a.ScriptAddress()
		switch len(serializedPubKey) {
		case 33: // compressed
			var compressedPubKey [33]byte
			copy(compressedPubKey[:], serializedPubKey)
			f.CompressedPubKeys[compressedPubKey] = struct{}{}
			return
		case 65: // uncompressed
			var uncompressedPubKey [65]byte
			copy(uncompressedPubKey[:], serializedPubKey)
			f.UncompressedPubKeys[uncompressedPubKey] = struct{}{}
			return
		}
	}
	f.OtherAddresses[a.EncodeAddress()] = struct{}{}
}

// AddAddressStr parses an address from a string and then adds it to the wsClientFilter using addAddress.
//
// NOTE: This extension was ported from github.com/decred/dcrd
func (f *WSClientFilter) AddAddressStr(s string, params *chaincfg.Params) {
	// If address can't be decoded, no point in saving it since it should also impossible to create the address from an
	// inspected transaction output script.
	a, e := btcaddr.Decode(s, params)
	if e != nil {
		return
	}
	f.AddAddress(a)
}

// AddUnspentOutPoint adds an outpoint to the wsClientFilter.
//
// NOTE: This extension was ported from github.com/decred/dcrd
func (f *WSClientFilter) AddUnspentOutPoint(op *wire.OutPoint) {
	f.Unspent[*op] = struct{}{}
}

// ExistsAddress returns true if the passed address has been added to the wsClientFilter. NOTE: This extension was
// ported from github.com/decred/dcrd
func (f *WSClientFilter) ExistsAddress(a btcaddr.Address) bool {
	switch a := a.(type) {
	case *btcaddr.PubKeyHash:
		_, ok := f.PubKeyHashes[*a.Hash160()]
		return ok
	case *btcaddr.ScriptHash:
		_, ok := f.ScriptHashes[*a.Hash160()]
		return ok
	case *btcaddr.PubKey:
		serializedPubKey := a.ScriptAddress()
		switch len(serializedPubKey) {
		case 33: // compressed
			var compressedPubKey [33]byte
			copy(compressedPubKey[:], serializedPubKey)
			_, ok := f.CompressedPubKeys[compressedPubKey]
			if !ok {
				_, ok = f.PubKeyHashes[*a.PubKeyHash().Hash160()]
			}
			return ok
		case 65: // uncompressed
			var uncompressedPubKey [65]byte
			copy(uncompressedPubKey[:], serializedPubKey)
			_, ok := f.UncompressedPubKeys[uncompressedPubKey]
			if !ok {
				_, ok = f.PubKeyHashes[*a.PubKeyHash().Hash160()]
			}
			return ok
		}
	}
	_, ok := f.OtherAddresses[a.EncodeAddress()]
	return ok
}

// ExistsUnspentOutPoint returns true if the passed outpoint has been added to the wsClientFilter.
//
// NOTE: This extension was ported from github.com/decred/dcrd
func (f *WSClientFilter) ExistsUnspentOutPoint(op *wire.OutPoint) bool {
	_, ok := f.Unspent[*op]
	return ok
}

// // removeAddress removes the passed address, if it exists, from the
// wsClientFilter. NOTE: This extension was ported from github.com/decred/dcrd
// func (f *wsClientFilter) removeAddress(a util.Address) {
// 	switch a := a.(type) {
// 	case *util.PubKeyHash:
// 		delete(f.pubKeyHashes, *a.Hash160())
// 		return
// 	case *util.ScriptHash:
// 		delete(f.scriptHashes, *a.Hash160())
// 		return
// 	case *util.PubKey:
// 		serializedPubKey := a.ScriptAddress()
// 		switch len(serializedPubKey) {
// 		case 33: // compressed
// 			var compressedPubKey [33]byte
// 			copy(compressedPubKey[:], serializedPubKey)
// 			delete(f.compressedPubKeys, compressedPubKey)
// 			return
// 		case 65: // uncompressed
// 			var uncompressedPubKey [65]byte
// 			copy(uncompressedPubKey[:], serializedPubKey)
// 			delete(f.uncompressedPubKeys, uncompressedPubKey)
// 			return
// 		}
// 	}
// 	delete(f.otherAddresses, a.EncodeAddress())
// }
// // removeAddressStr parses an address from a string and then removes it from
// // the wsClientFilter using removeAddress. NOTE: This extension was ported
// // from github.com/decred/dcrd
// func (f *wsClientFilter) removeAddressStr(s string, netparams *chaincfg.Params) {
// 	a, e := util.Decode(s, netparams)
// 	if e ==  nil {
// 		f.removeAddress(a)
// 	} else {
// 		delete(f.otherAddresses, s)
// 	}
// }
// // removeUnspentOutPoint removes the passed outpoint, if it exists, from the
// wsClientFilter. NOTE: This extension was ported from github.com/decred/dcrd
// func (f *wsClientFilter) removeUnspentOutPoint(op *wire.OutPoint) {
// 	delete(f.unspent, *op)
// }

// AddClient adds the passed websocket client to the notification manager.
func (m *WSNtfnMgr) AddClient(wsc *WSClient) {
	m.QueueNotification <- (*NotificationRegisterClient)(wsc)
}

// SendNotifyBlockConnected passes a block newly-connected to the best chain to the notification manager for block and
// transaction notification processing.
func (m *WSNtfnMgr) SendNotifyBlockConnected(block *block.Block) {
	// As NotifyBlockConnected will be called by the block manager and the RPC Server may no longer be running, use a
	// select statement to unblock enqueuing the notification once the RPC Server has begun shutting down.
	select {
	case m.QueueNotification <- (*NotificationBlockConnected)(block):
	case <-m.Quit.Wait():
	}
}

// SendNotifyBlockDisconnected passes a block disconnected from the best chain to the notification manager for block
// notification processing.
func (m *WSNtfnMgr) SendNotifyBlockDisconnected(block *block.Block) {
	// As NotifyBlockDisconnected will be called by the block manager and the RPC Server may no longer be running, use a
	// select statement to unblock enqueuing the notification once the RPC Server has begun shutting down.
	select {
	case m.QueueNotification <- (*NotificationBlockDisconnected)(block):
	case <-m.Quit.Wait():
	}
}

// SendNotifyMempoolTx passes a transaction accepted by mempool to the notification manager for transaction notification
// processing. If isNew is true, the tx is is a new transaction, rather than one added to the mempool during a reorg.
func (m *WSNtfnMgr) SendNotifyMempoolTx(tx *util.Tx, isNew bool) {
	n := &NotificationTxAcceptedByMempool{
		IsNew: isNew,
		Tx:    tx,
	}
	// As NotifyMempoolTx will be called by mempool and the RPC Server may no longer be running, use a select statement
	// to unblock enqueuing the notification once the RPC Server has begun shutting down.
	select {
	case m.QueueNotification <- n:
	case <-m.Quit.Wait():
	}
}

// GetNumClients returns the number of clients actively being served.
func (m *WSNtfnMgr) GetNumClients() (n int) {
	select {
	case n = <-m.NumClients:
	case <-m.Quit.Wait(): // Use default n (0) if Server has shut down.
	}
	return
}

// RegisterBlockUpdates requests block update notifications to the passed websocket client.
func (m *WSNtfnMgr) RegisterBlockUpdates(wsc *WSClient) {
	m.QueueNotification <- (*NotificationRegisterBlocks)(wsc)
}

// RegisterNewMempoolTxsUpdates requests notifications to the passed websocket client when new transactions are added to
// the memory pool.
func (m *WSNtfnMgr) RegisterNewMempoolTxsUpdates(wsc *WSClient) {
	m.QueueNotification <- (*NotificationRegisterNewMempoolTxs)(wsc)
}

// RegisterSpentRequests requests a notification when each of the passed outpoints is confirmed spent (contained in a
// block connected to the main chain) for the passed websocket client. The request is automatically removed once the
// notification has been sent.
func (m *WSNtfnMgr) RegisterSpentRequests(wsc *WSClient, ops []*wire.OutPoint) {
	m.QueueNotification <- &NotificationRegisterSpent{
		WSC: wsc,
		OPs: ops,
	}
}

// RegisterTxOutAddressRequests requests notifications to the passed websocket client when a transaction output spends
// to the passed address.
func (m *WSNtfnMgr) RegisterTxOutAddressRequests(
	wsc *WSClient, addrs []string,
) {
	m.QueueNotification <- &NotificationRegisterAddr{
		WSC:   wsc,
		Addrs: addrs,
	}
}

// RemoveClient removes the passed websocket client and all notifications registered for it.
func (m *WSNtfnMgr) RemoveClient(wsc *WSClient) {
	select {
	case m.QueueNotification <- (*NotificationUnregisterClient)(wsc):
	case <-m.Quit.Wait():
	}
}

// Shutdown shuts down the manager, stopping the notification queue and notification handler goroutines.
func (m *WSNtfnMgr) Shutdown() {
	m.Quit.Q()
}

// Start starts the goroutines required for the manager to queue and process websocket client notifications.
func (m *WSNtfnMgr) Start() {
	go m.QueueHandler()
	go m.NotificationHandler()
}

// UnregisterBlockUpdates removes block update notifications for the passed websocket client.
func (m *WSNtfnMgr) UnregisterBlockUpdates(wsc *WSClient) {
	m.QueueNotification <- (*NotificationUnregisterBlocks)(wsc)
}

// UnregisterNewMempoolTxsUpdates removes notifications to the passed websocket client when new transaction are added to
// the memory pool.
func (m *WSNtfnMgr) UnregisterNewMempoolTxsUpdates(wsc *WSClient) {
	m.QueueNotification <- (*NotificationUnregisterNewMempoolTxs)(wsc)
}

// UnregisterSpentRequest removes a request from the passed websocket client to be notified when the passed outpoint is
// confirmed spent (contained in a block connected to the main chain).
func (m *WSNtfnMgr) UnregisterSpentRequest(
	wsc *WSClient,
	op *wire.OutPoint,
) {
	m.QueueNotification <- &NotificationUnregisterSpent{
		WSC: wsc,
		OP:  op,
	}
}

// UnregisterTxOutAddressRequest removes a request from the passed websocket client to be notified when a transaction
// spends to the passed address.
func (m *WSNtfnMgr) UnregisterTxOutAddressRequest(wsc *WSClient, addr string) {
	m.QueueNotification <- &NotificationUnregisterAddr{
		WSC:  wsc,
		Addr: addr,
	}
}

// WaitForShutdown blocks until all notification manager goroutines have finished.
func (m *WSNtfnMgr) WaitForShutdown() {
	m.WG.Wait()
}

// AddAddrRequests adds the websocket client wsc to the address to client set addrMap so wsc will be notified for any
// mempool or block transaction outputs spending to any of the addresses in addrs.
func (*WSNtfnMgr) AddAddrRequests(
	addrMap map[string]map[qu.C]*WSClient, wsc *WSClient, addrs []string,
) {
	for _, addr := range addrs {
		// Track the request in the client as well so it can be quickly be removed on disconnect.
		wsc.AddrRequests[addr] = struct{}{}
		// Add the client to the set of clients to notify when the outpoint is seen. Create map as needed.
		cmap, ok := addrMap[addr]
		if !ok {
			cmap = make(map[qu.C]*WSClient)
			addrMap[addr] = cmap
		}
		cmap[wsc.Quit] = wsc
	}
}

// AddSpentRequests modifies a map of watched outpoints to sets of websocket clients to add a new request watch all of
// the outpoints in ops and create and send a notification when spent to the websocket client wsc.
func (m *WSNtfnMgr) AddSpentRequests(
	opMap map[wire.
OutPoint]map[qu.C]*WSClient, wsc *WSClient, ops []*wire.OutPoint,
) {
	for _, op := range ops {
		// Track the request in the client as well so it can be quickly be removed on disconnect.
		wsc.SpentRequests[*op] = struct{}{}
		// Add the client to the list to notify when the outpoint is seen. Create the list as needed.
		cmap, ok := opMap[*op]
		if !ok {
			cmap = make(map[qu.C]*WSClient)
			opMap[*op] = cmap
		}
		cmap[wsc.Quit] = wsc
	}
	// Chk if any transactions spending these outputs already exists in the mempool, if so send the notification
	// immediately.
	spends := make(map[chainhash.Hash]*util.Tx)
	for _, op := range ops {
		spend := m.Server.Cfg.TxMemPool.CheckSpend(*op)
		if spend != nil {
			D.F(
				"found existing mempool spend for outpoint<%v>: %v",
				op, spend.Hash(),
			)
			spends[*spend.Hash()] = spend
		}
	}
	for _, spend := range spends {
		m.NotifyForTx(opMap, nil, spend, nil)
	}
}

// NotificationHandler reads notifications and control messages from the queue handler and processes one at a time.
func (m *WSNtfnMgr) NotificationHandler() {
	// clients is a map of all currently connected websocket clients.
	clients := make(map[qu.C]*WSClient)
	// Maps used to hold lists of websocket clients to be notified on certain events. Each websocket client also keeps
	// maps for the events which have multiple triggers to make removal from these lists on connection close less
	// horrendously.
	//
	// Where possible, the quit channel is used as the unique id for a client since it is quite a bit more efficient
	// than using the entire struct.
	blockNotifications := make(map[qu.C]*WSClient)
	txNotifications := make(map[qu.C]*WSClient)
	watchedOutPoints := make(map[wire.OutPoint]map[qu.C]*WSClient)
	watchedAddrs := make(map[string]map[qu.C]*WSClient)
out:
	for {
		select {
		case n, ok := <-m.NotificationMsgs:
			if !ok {
				// queueHandler quit.
				break out
			}
			switch n := n.(type) {
			case *NotificationBlockConnected:
				block := (*block.Block)(n)
				// Skip iterating through all txs if no tx notification requests exist.
				if len(watchedOutPoints) != 0 || len(watchedAddrs) != 0 {
					for _, tx := range block.Transactions() {
						m.NotifyForTx(
							watchedOutPoints,
							watchedAddrs, tx, block,
						)
					}
				}
				if len(blockNotifications) != 0 {
					m.NotifyBlockConnected(
						blockNotifications,
						block,
					)
					m.NotifyFilteredBlockConnected(
						blockNotifications,
						block,
					)
				}
			case *NotificationBlockDisconnected:
				block := (*block.Block)(n)
				if len(blockNotifications) != 0 {
					m.NotifyBlockDisconnected(
						blockNotifications,
						block,
					)
					m.NotifyFilteredBlockDisconnected(
						blockNotifications,
						block,
					)
				}
			case *NotificationTxAcceptedByMempool:
				if n.IsNew && len(txNotifications) != 0 {
					m.NotifyForNewTx(txNotifications, n.Tx)
				}
				m.NotifyForTx(watchedOutPoints, watchedAddrs, n.Tx, nil)
				m.NotifyRelevantTxAccepted(n.Tx, clients)
			case *NotificationRegisterBlocks:
				wsc := (*WSClient)(n)
				blockNotifications[wsc.Quit] = wsc
			case *NotificationUnregisterBlocks:
				wsc := (*WSClient)(n)
				delete(blockNotifications, wsc.Quit)
			case *NotificationRegisterClient:
				wsc := (*WSClient)(n)
				clients[wsc.Quit] = wsc
			case *NotificationUnregisterClient:
				wsc := (*WSClient)(n)
				// Remove any requests made by the client as well as the client itself.
				delete(blockNotifications, wsc.Quit)
				delete(txNotifications, wsc.Quit)
				for k := range wsc.SpentRequests {
					op := k
					m.RemoveSpentRequest(watchedOutPoints, wsc, &op)
				}
				for addr := range wsc.AddrRequests {
					m.RemoveAddrRequest(watchedAddrs, wsc, addr)
				}
				delete(clients, wsc.Quit)
			case *NotificationRegisterSpent:
				m.AddSpentRequests(watchedOutPoints, n.WSC, n.OPs)
			case *NotificationUnregisterSpent:
				m.RemoveSpentRequest(watchedOutPoints, n.WSC, n.OP)
			case *NotificationRegisterAddr:
				m.AddAddrRequests(watchedAddrs, n.WSC, n.Addrs)
			case *NotificationUnregisterAddr:
				m.RemoveAddrRequest(watchedAddrs, n.WSC, n.Addr)
			case *NotificationRegisterNewMempoolTxs:
				wsc := (*WSClient)(n)
				txNotifications[wsc.Quit] = wsc
			case *NotificationUnregisterNewMempoolTxs:
				wsc := (*WSClient)(n)
				delete(txNotifications, wsc.Quit)
			default:
				W.Ln("unhandled notification type")
			}
		case m.NumClients <- len(clients):
		case <-m.Quit.Wait():
			// RPC Server shutting down.
			break out
		}
	}
	for _, c := range clients {
		c.Disconnect()
	}
	m.WG.Done()
}

// NotifyBlockConnected notifies websocket clients that have registered for block updates when a block is connected to
// the main chain.
func (*WSNtfnMgr) NotifyBlockConnected(
	clients map[qu.C]*WSClient, block *block.Block,
) {
	// Notify interested websocket clients about the connected block.
	ntfn := btcjson.NewBlockConnectedNtfn(
		block.Hash().String(), block.Height(),
		block.WireBlock().Header.Timestamp.Unix(),
	)
	marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
	if e != nil {
		E.Ln("failed to marshal block connected notification:", e)
		return
	}
	for _, wsc := range clients {
		e := wsc.QueueNotification(marshalledJSON)
		if e != nil {
		}
	}
}

// NotifyBlockDisconnected notifies websocket clients that have registered for block updates when a block is
// disconnected from the main chain (due to a reorganize).
func (*WSNtfnMgr) NotifyBlockDisconnected(
	clients map[qu.C]*WSClient, block *block.Block,
) {
	// Skip notification creation if no clients have requested block connected/ disconnected notifications.
	if len(clients) == 0 {
		return
	}
	// Notify interested websocket clients about the disconnected block.
	ntfn := btcjson.NewBlockDisconnectedNtfn(
		block.Hash().String(),
		block.Height(), block.WireBlock().Header.Timestamp.Unix(),
	)
	marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
	if e != nil {
		E.Ln(
			"failed to marshal block disconnected notification:",
			e,
		)
		return
	}
	for _, wsc := range clients {
		e := wsc.QueueNotification(marshalledJSON)
		if e != nil {
		}
	}
}

// NotifyFilteredBlockConnected notifies websocket clients that have registered for block updates when a block is
// connected to the main chain.
func (m *WSNtfnMgr) NotifyFilteredBlockConnected(
	clients map[qu.C]*WSClient, block *block.Block,
) {
	// Create the common portion of the notification that is the same for every client.
	var w bytes.Buffer
	e := block.WireBlock().Header.Serialize(&w)
	if e != nil {
		E.Ln(
			"failed to serialize header for filtered block connected"+
				" notification:", e,
		)
		return
	}
	ntfn := btcjson.NewFilteredBlockConnectedNtfn(
		block.Height(),
		hex.EncodeToString(w.Bytes()), nil,
	)
	// Search for relevant transactions for each client and save them serialized in hex encoding for the notification.
	subscribedTxs := make(map[qu.C][]string)
	for _, tx := range block.Transactions() {
		var txHex string
		for quitChan := range m.GetSubscribedClients(tx, clients) {
			if txHex == "" {
				txHex = TxHexString(tx.MsgTx())
			}
			subscribedTxs[quitChan] = append(subscribedTxs[quitChan], txHex)
		}
	}
	for quitChan, wsc := range clients {
		// Add all discovered transactions for this client. For clients that have no new-style filter, add the empty
		// string slice.
		ntfn.SubscribedTxs = subscribedTxs[quitChan]
		// Marshal and queue notification.
		marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
		if e != nil {
			E.F("failed to marshal filtered block connected notification:", e)
			return
		}
		e = wsc.QueueNotification(marshalledJSON)
		if e != nil {
			D.Ln(e)
		}
	}
}

// NotifyFilteredBlockDisconnected notifies websocket clients that have registered for block updates when a block is
// disconnected from the main chain (due to a reorganize).
func (*WSNtfnMgr) NotifyFilteredBlockDisconnected(
	clients map[qu.C]*WSClient, block *block.Block,
) {
	// Skip notification creation if no clients have requested block connected/ disconnected notifications.
	if len(clients) == 0 {
		return
	}
	// Notify interested websocket clients about the disconnected block.
	var w bytes.Buffer
	e := block.WireBlock().Header.Serialize(&w)
	if e != nil {
		E.Ln("failed to serialize header for filtered block disconnected notification:", e)
		return
	}
	ntfn := btcjson.NewFilteredBlockDisconnectedNtfn(block.Height(), hex.EncodeToString(w.Bytes()))
	marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
	if e != nil {
		E.Ln("failed to marshal filtered block disconnected notification:", e)
		return
	}
	for _, wsc := range clients {
		e := wsc.QueueNotification(marshalledJSON)
		if e != nil {
			D.Ln(e)
		}
	}
}

// NotifyForNewTx notifies websocket clients that have registered for updates when a new transaction is added to the
// memory pool.
func (m *WSNtfnMgr) NotifyForNewTx(
	clients map[qu.C]*WSClient,
	tx *util.Tx,
) {
	txHashStr := tx.Hash().String()
	mtx := tx.MsgTx()
	var amount int64
	for _, txOut := range mtx.TxOut {
		amount += txOut.Value
	}
	ntfn := btcjson.NewTxAcceptedNtfn(txHashStr, amount2.Amount(amount).ToDUO())
	marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
	if e != nil {
		E.Ln("failed to marshal tx notification:", e)
		return
	}
	var verboseNtfn *btcjson.TxAcceptedVerboseNtfn
	var marshalledJSONVerbose []byte
	for _, wsc := range clients {
		if wsc.VerboseTxUpdates {
			net := m.Server.Cfg.ChainParams
			var rawTx *btcjson.TxRawResult
			rawTx, e = CreateTxRawResult(
				net, mtx, txHashStr, nil,
				"", 0, 0,
			)
			if e != nil {
				return
			}
			verboseNtfn = btcjson.NewTxAcceptedVerboseNtfn(*rawTx)
			marshalledJSONVerbose, e = btcjson.MarshalCmd(
				nil,
				verboseNtfn,
			)
			if e != nil {
				E.Ln("failed to marshal verbose tx notification:", e)
			}
			if marshalledJSONVerbose != nil {
				e = wsc.QueueNotification(marshalledJSONVerbose)
				if e != nil {
					D.Ln(e)
				}
				continue
			}
			return
		}
		e = wsc.QueueNotification(marshalledJSONVerbose)
		if e != nil {
		} else {
			e := wsc.QueueNotification(marshalledJSON)
			if e != nil {
			}
		}
	}
}

// NotifyForTx examines the inputs and outputs of the passed transaction, notifying websocket clients of outputs
// spending to a watched address and inputs spending a watched outpoint.
func (m *WSNtfnMgr) NotifyForTx(
	ops map[wire.OutPoint]map[qu.C]*WSClient,
	addrs map[string]map[qu.C]*WSClient, tx *util.Tx, block *block.Block,
) {
	if len(ops) != 0 {
		m.NotifyForTxIns(ops, tx, block)
	}
	if len(addrs) != 0 {
		m.NotifyForTxOuts(ops, addrs, tx, block)
	}
}

// NotifyForTxIns examines the inputs of the passed transaction and sends interested websocket clients a redeemingtx
// notification if any inputs spend a watched output. If block is non-nil, any matching spent requests are removed.
func (m *WSNtfnMgr) NotifyForTxIns(
	ops map[wire.
OutPoint]map[qu.C]*WSClient, tx *util.Tx, block *block.Block,
) {
	// Nothing to do if nobody is watching outpoints.
	if len(ops) == 0 {
		return
	}
	txHex := ""
	wscNotified := make(map[qu.C]struct{})
	for _, txIn := range tx.MsgTx().TxIn {
		prevOut := &txIn.PreviousOutPoint
		if cmap, ok := ops[*prevOut]; ok {
			if txHex == "" {
				txHex = TxHexString(tx.MsgTx())
			}
			marshalledJSON, e := NewRedeemingTxNotification(txHex, tx.Index(), block)
			if e != nil {
				E.Ln(
					"failed to marshal redeemingtx notification:", e,
				)
				continue
			}
			for wscQuit, wsc := range cmap {
				if block != nil {
					m.RemoveSpentRequest(ops, wsc, prevOut)
				}
				if _, ok := wscNotified[wscQuit]; !ok {
					wscNotified[wscQuit] = struct{}{}
					e := wsc.QueueNotification(marshalledJSON)
					if e != nil {
						D.Ln(e)
					}
				}
			}
		}
	}
}

// NotifyForTxOuts examines each transaction output, notifying interested websocket clients of the transaction if an
// output spends to a watched address. A spent notification request is automatically registered for the client for each
// matching output.
func (m *WSNtfnMgr) NotifyForTxOuts(
	ops map[wire.OutPoint]map[qu.C]*WSClient,
	addrs map[string]map[qu.C]*WSClient, tx *util.Tx, block *block.Block,
) {
	// Nothing to do if nobody is listening for address notifications.
	if len(addrs) == 0 {
		return
	}
	txHex := ""
	wscNotified := make(map[qu.C]struct{})
	for i, txOut := range tx.MsgTx().TxOut {
		var txAddrs []btcaddr.Address
		var e error
		_, txAddrs, _, e = txscript.ExtractPkScriptAddrs(
			txOut.PkScript, m.Server.Cfg.ChainParams,
		)
		if e != nil {
			continue
		}
		for _, txAddr := range txAddrs {
			cmap, ok := addrs[txAddr.EncodeAddress()]
			if !ok {
				continue
			}
			if txHex == "" {
				txHex = TxHexString(tx.MsgTx())
			}
			ntfn := btcjson.NewRecvTxNtfn(
				txHex, BlockDetails(
					block,
					tx.Index(),
				),
			)
			marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
			if e != nil {
				E.Ln("Failed to marshal processedtx notification:", e)
				continue
			}
			op := []*wire.OutPoint{wire.NewOutPoint(tx.Hash(), uint32(i))}
			for wscQuit, wsc := range cmap {
				m.AddSpentRequests(ops, wsc, op)
				if _, ok := wscNotified[wscQuit]; !ok {
					wscNotified[wscQuit] = struct{}{}
					e := wsc.QueueNotification(marshalledJSON)
					if e != nil {
					}
				}
			}
		}
	}
}

// NotifyRelevantTxAccepted examines the inputs and outputs of the passed transaction, notifying websocket clients of
// outputs spending to a watched address and inputs spending a watched outpoint.
//
// Any outputs paying to a watched address result in the output being watched as well for future notifications.
func (m *WSNtfnMgr) NotifyRelevantTxAccepted(
	tx *util.Tx, clients map[qu.C]*WSClient,
) {
	clientsToNotify := m.GetSubscribedClients(tx, clients)
	if len(clientsToNotify) != 0 {
		n := btcjson.NewRelevantTxAcceptedNtfn(TxHexString(tx.MsgTx()))
		marshalled, e := btcjson.MarshalCmd(nil, n)
		if e != nil {
			E.Ln("failed to marshal notification:", e)
			return
		}
		for quitChan := range clientsToNotify {
			e := clients[quitChan].QueueNotification(marshalled)
			if e != nil {
			}
		}
	}
}

// QueueHandler maintains a queue of notifications and notification handler control messages.
func (m *WSNtfnMgr) QueueHandler() {
	QueueHandler(m.QueueNotification, m.NotificationMsgs, m.Quit)
	m.WG.Done()
}

// RemoveAddrRequest removes the websocket client wsc from the address to client set addrs so it will no longer receive
// notification updates for any transaction outputs send to addr.
func (*WSNtfnMgr) RemoveAddrRequest(
	addrs map[string]map[qu.C]*WSClient, wsc *WSClient, addr string,
) {
	// Remove the request tracking from the client.
	delete(wsc.AddrRequests, addr)
	// Remove the client from the list to notify.
	cmap, ok := addrs[addr]
	if !ok {
		W.F(
			"attempt to remove nonexistent addr request <%s> for websocket client %s",
			addr, wsc.Addr,
		)
		return
	}
	delete(cmap, wsc.Quit)
	// Remove the map entry altogether if there are no more clients interested in it.
	if len(cmap) == 0 {
		delete(addrs, addr)
	}
}

// RemoveSpentRequest modifies a map of watched outpoints to remove the websocket client wsc from the set of clients to
// be notified when a watched outpoint is spent. If wsc is the last client, the outpoint key is removed from the map.
func (*WSNtfnMgr) RemoveSpentRequest(
	ops map[wire.
OutPoint]map[qu.C]*WSClient, wsc *WSClient, op *wire.OutPoint,
) {
	// Remove the request tracking from the client.
	delete(wsc.SpentRequests, *op)
	// Remove the client from the list to notify.
	notifyMap, ok := ops[*op]
	if !ok {
		W.Ln(
			"attempt to remove nonexistent spent request for"+
				" websocket client", wsc.Addr,
		)
		return
	}
	delete(notifyMap, wsc.Quit)
	// Remove the map entry altogether if there are no more clients interested in it.
	if len(notifyMap) == 0 {
		delete(ops, *op)
	}
}

// GetSubscribedClients returns the set of all websocket client quit channels that are registered to receive
// notifications regarding tx, either due to tx spending a watched output or outputting to a watched address.
//
// Matching client's filters are updated based on this transaction's outputs and output addresses that may be relevant
// for a client.
func (m *WSNtfnMgr) GetSubscribedClients(
	tx *util.Tx,
	clients map[qu.C]*WSClient,
) map[qu.C]struct{} {
	// Use a map of client quit channels as keys to prevent duplicates when multiple inputs and/or outputs are relevant
	// to the client.
	subscribed := make(map[qu.C]struct{})
	msgTx := tx.MsgTx()
	for _, input := range msgTx.TxIn {
		for quitChan, wsc := range clients {
			wsc.Lock()
			filter := wsc.FilterData
			wsc.Unlock()
			if filter == nil {
				continue
			}
			filter.mu.Lock()
			if filter.ExistsUnspentOutPoint(&input.PreviousOutPoint) {
				subscribed[quitChan] = struct{}{}
			}
			filter.mu.Unlock()
		}
	}
	var e error
	for i, output := range msgTx.TxOut {
		var addrs []btcaddr.Address
		_, addrs, _, e = txscript.ExtractPkScriptAddrs(
			output.PkScript, m.Server.Cfg.ChainParams,
		)
		if e != nil {
			// Clients are not able to subscribe to nonstandard or non-address outputs.
			continue
		}
		for quitChan, wsc := range clients {
			wsc.Lock()
			filter := wsc.FilterData
			wsc.Unlock()
			if filter == nil {
				continue
			}
			filter.mu.Lock()
			for _, a := range addrs {
				if filter.ExistsAddress(a) {
					subscribed[quitChan] = struct{}{}
					op := wire.OutPoint{
						Hash:  *tx.Hash(),
						Index: uint32(i),
					}
					filter.AddUnspentOutPoint(&op)
				}
			}
			filter.mu.Unlock()
		}
	}
	return subscribed
}
func (s Semaphore) Acquire() {
	s <- struct{}{}
}
func (s Semaphore) Release() {
	<-s
}

// BlockDetails creates a BlockDetails struct to include in btcws notifications from a block and a transaction's block
// index.
func BlockDetails(block *block.Block, txIndex int) *btcjson.BlockDetails {
	if block == nil {
		return nil
	}
	return &btcjson.BlockDetails{
		Height: block.Height(),
		Hash:   block.Hash().String(),
		Index:  txIndex,
		Time:   block.WireBlock().Header.Timestamp.Unix(),
	}
}

// CheckAddressValidity checks the validity of each address in the passed string slice. It does this by attempting to
// decode each address using the current active network parameters.
//
// If any single address fails to decode properly, the function returns an error. Otherwise, nil is returned.
func CheckAddressValidity(addrs []string, params *chaincfg.Params) (e error) {
	for _, addr := range addrs {
		_, e := btcaddr.Decode(addr, params)
		if e != nil {
			return &btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidAddressOrKey,
				Message: fmt.Sprintf(
					"Invalid address or key: %v",
					addr,
				),
			}
		}
	}
	return nil
}

// DescendantBlock returns the appropriate JSON-RPC error if a current block fetched during a reorganize is not a direct
// child of the parent block hash.
func DescendantBlock(
	prevHash *chainhash.Hash, curBlock *block.Block,
) (e error) {
	curHash := &curBlock.WireBlock().Header.PrevBlock
	if !prevHash.IsEqual(curHash) {
		E.F(
			"stopping rescan for reorged block %v (replaced by block %v)",
			prevHash, curHash,
		)
		return &ErrRescanReorg
	}
	return nil
}

// DeserializeOutpoints deserializes each serialized outpoint.
func DeserializeOutpoints(serializedOuts []btcjson.OutPoint) (
	[]*wire.OutPoint,
	error,
) {
	outpoints := make([]*wire.OutPoint, 0, len(serializedOuts))
	for i := range serializedOuts {
		blockHash, e := chainhash.NewHashFromStr(serializedOuts[i].Hash)
		if e != nil {
			return nil, DecodeHexError(serializedOuts[i].Hash)
		}
		index := serializedOuts[i].Index
		outpoints = append(outpoints, wire.NewOutPoint(blockHash, index))
	}
	return outpoints, nil
}

// HandleLoadTxFilter implements the loadtxfilter command extension for websocket connections. NOTE: This extension is
// ported from github.com/decred/dcrd
func HandleLoadTxFilter(wsc *WSClient, icmd interface{}) (interface{}, error) {
	cmd := icmd.(*btcjson.LoadTxFilterCmd)
	outPoints := make([]wire.OutPoint, len(cmd.OutPoints))
	for i := range cmd.OutPoints {
		hash, e := chainhash.NewHashFromStr(cmd.OutPoints[i].Hash)
		if e != nil {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidParameter,
				Message: e.Error(),
			}
		}
		outPoints[i] = wire.OutPoint{
			Hash:  *hash,
			Index: cmd.OutPoints[i].Index,
		}
	}
	params := wsc.Server.Cfg.ChainParams
	wsc.Lock()
	if cmd.Reload || wsc.FilterData == nil {
		wsc.FilterData = NewWSClientFilter(
			cmd.Addresses, outPoints,
			params,
		)
		wsc.Unlock()
	} else {
		wsc.Unlock()
		wsc.FilterData.mu.Lock()
		for _, a := range cmd.Addresses {
			wsc.FilterData.AddAddressStr(a, params)
		}
		for i := range outPoints {
			wsc.FilterData.AddUnspentOutPoint(&outPoints[i])
		}
		wsc.FilterData.mu.Unlock()
	}
	return nil, nil
}

// HandleNotifyBlocks implements the notifyblocks command extension for websocket connections.
func HandleNotifyBlocks(wsc *WSClient, icmd interface{}) (interface{}, error) {
	wsc.Server.NtfnMgr.RegisterBlockUpdates(wsc)
	return nil, nil
}

// HandleNotifyNewTransactions implements the notifynewtransactions command
// extension for websocket connections.
func HandleNotifyNewTransactions(
	wsc *WSClient,
	icmd interface{},
) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.NotifyNewTransactionsCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	wsc.VerboseTxUpdates = cmd.Verbose != nil && *cmd.Verbose
	wsc.Server.NtfnMgr.RegisterNewMempoolTxsUpdates(wsc)
	return nil, nil
}

// HandleNotifyReceived implements the notifyreceived command extension for websocket connections.
func HandleNotifyReceived(wsc *WSClient, icmd interface{}) (
	interface{},
	error,
) {
	cmd, ok := icmd.(*btcjson.NotifyReceivedCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	// Decode addresses to validate input, but the strings slice is used directly if these are all ok.
	e := CheckAddressValidity(cmd.Addresses, wsc.Server.Cfg.ChainParams)
	if e != nil {
		return nil, e
	}
	wsc.Server.NtfnMgr.RegisterTxOutAddressRequests(wsc, cmd.Addresses)
	return nil, nil
}

// HandleNotifySpent implements the notifyspent command extension for websocket connections.
func HandleNotifySpent(wsc *WSClient, icmd interface{}) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.NotifySpentCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	outpoints, e := DeserializeOutpoints(cmd.OutPoints)
	if e != nil {
		return nil, e
	}
	wsc.Server.NtfnMgr.RegisterSpentRequests(wsc, outpoints)
	return nil, nil
}

// HandleRescan implements the rescan command extension for websocket connections.
//
// NOTE: This does not smartly handle reorgs, and fixing requires database changes (for safe, concurrent access to full
// block ranges, and support for other chains than the best chain).
//
// It will, however, detect whether a reorg removed a block that was previously processed, and result in the handler
// erroring.
//
// Clients must handle this by finding a block still in the chain (perhaps from a rescanprogress notification) to resume
// their rescan. TODO: simplify, modularise this
func HandleRescan(wsc *WSClient, icmd interface{}) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.RescanCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	outpoints := make([]*wire.OutPoint, 0, len(cmd.OutPoints))
	for i := range cmd.OutPoints {
		cmdOutpoint := &cmd.OutPoints[i]
		blockHash, e := chainhash.NewHashFromStr(cmdOutpoint.Hash)
		if e != nil {
			return nil, DecodeHexError(cmdOutpoint.Hash)
		}
		outpoint := wire.NewOutPoint(blockHash, cmdOutpoint.Index)
		outpoints = append(outpoints, outpoint)
	}
	numAddrs := len(cmd.Addresses)
	if numAddrs == 1 {
		D.Ln("beginning rescan for 1 address")
	} else {
		D.F("beginning rescan for %d addresses", numAddrs)
	}
	// Build lookup maps.
	lookups := RescanKeys{
		Fallbacks:           map[string]struct{}{},
		PubKeyHashes:        map[[ripemd160.Size]byte]struct{}{},
		ScriptHashes:        map[[ripemd160.Size]byte]struct{}{},
		CompressedPubKeys:   map[[33]byte]struct{}{},
		UncompressedPubKeys: map[[65]byte]struct{}{},
		Unspent:             map[wire.OutPoint]struct{}{},
	}
	var compressedPubkey [33]byte
	var uncompressedPubkey [65]byte
	params := wsc.Server.Cfg.ChainParams
	for _, addrStr := range cmd.Addresses {
		addr, e := btcaddr.Decode(addrStr, params)
		if e != nil {
			jsonErr := btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Rescan address " + addrStr + ": " +
					e.Error(),
			}
			return nil, &jsonErr
		}
		switch a := addr.(type) {
		case *btcaddr.PubKeyHash:
			lookups.PubKeyHashes[*a.Hash160()] = struct{}{}
		case *btcaddr.ScriptHash:
			lookups.ScriptHashes[*a.Hash160()] = struct{}{}
		case *btcaddr.PubKey:
			pubkeyBytes := a.ScriptAddress()
			switch len(pubkeyBytes) {
			case 33: // Compressed
				copy(compressedPubkey[:], pubkeyBytes)
				lookups.CompressedPubKeys[compressedPubkey] = struct{}{}
			case 65: // Uncompressed
				copy(uncompressedPubkey[:], pubkeyBytes)
				lookups.UncompressedPubKeys[uncompressedPubkey] = struct{}{}
			default:
				jsonErr := btcjson.RPCError{
					Code:    btcjson.ErrRPCInvalidAddressOrKey,
					Message: "Pubkey " + addrStr + " is of unknown length",
				}
				return nil, &jsonErr
			}
		default:
			// A new address type must have been added. Use encoded payment address string as a fallback until a fast
			// path is added.
			lookups.Fallbacks[addrStr] = struct{}{}
		}
	}
	for _, outpoint := range outpoints {
		lookups.Unspent[*outpoint] = struct{}{}
	}
	chain := wsc.Server.Cfg.Chain
	minBlockHash, e := chainhash.NewHashFromStr(cmd.BeginBlock)
	if e != nil {
		return nil, DecodeHexError(cmd.BeginBlock)
	}
	minBlock, e := chain.BlockHeightByHash(minBlockHash)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCBlockNotFound,
			Message: "ScriptError getting block: " + e.Error(),
		}
	}
	maxBlock := int32(math.MaxInt32)
	if cmd.EndBlock != nil {
		maxBlockHash, e := chainhash.NewHashFromStr(*cmd.EndBlock)
		if e != nil {
			return nil, DecodeHexError(*cmd.EndBlock)
		}
		maxBlock, e = chain.BlockHeightByHash(maxBlockHash)
		if e != nil {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCBlockNotFound,
				Message: "ScriptError getting block: " + e.Error(),
			}
		}
	}
	// lastBlock and lastBlockHash track the previously-rescanned block. They equal nil when no previous blocks have
	// been rescanned.
	var lastBlock *block.Block
	var lastBlockHash *chainhash.Hash
	// A ticker is created to wait at least 10 seconds before notifying the websocket client of the current progress
	// completed by the rescan.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// Instead of fetching all block shas at once, fetch in smaller chunks to ensure large rescans consume a limited
	// amount of memory.
fetchRange:
	for minBlock < maxBlock {
		// Limit the max number of hashes to fetch at once to the maximum number of items allowed in a single inventory.
		// This value could be higher since it's not creating inventory messages, but this mirrors the limiting logic
		// used in the peer-to-peer protocol.
		maxLoopBlock := maxBlock
		if maxLoopBlock-minBlock > wire.MaxInvPerMsg {
			maxLoopBlock = minBlock + wire.MaxInvPerMsg
		}
		hashList, e := chain.HeightRange(minBlock, maxLoopBlock)
		if e != nil {
			E.Ln("error looking up block range:", e)
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCDatabase,
				Message: "Database error: " + e.Error(),
			}
		}
		if len(hashList) == 0 {
			// The rescan is finished if no blocks hashes for this range were successfully fetched and a stop block was
			// provided.
			if maxBlock != math.MaxInt32 {
				break
			}
			// If the rescan is through the current block, set up the client to continue to receive notifications
			// regarding all rescanned addresses and the current set of unspent outputs.
			//
			// This is done safely by temporarily grabbing exclusive access of the block manager. If no more blocks have
			// been attached between this pause and the fetch above, then it is safe to register the websocket client
			// for continuous notifications if necessary.
			//
			// Otherwise, continue the fetch loop again to rescan the new blocks (or error due to an irrecoverable
			// reorganize).
			pauseGuard := wsc.Server.Cfg.SyncMgr.Pause()
			best := wsc.Server.Cfg.Chain.BestSnapshot()
			curHash := &best.Hash
			again := true
			if lastBlockHash == nil || *lastBlockHash == *curHash {
				again = false
				n := wsc.Server.NtfnMgr
				n.RegisterSpentRequests(wsc, lookups.UnspentSlice())
				n.RegisterTxOutAddressRequests(wsc, cmd.Addresses)
			}
			close(pauseGuard)
			// this err value is nil if it got to here from above at 1577
			// if e != nil  {
			//  Errorf(
			// 		"ScriptError fetching best block hash:", e,
			// 	}
			// 	return nil, &json.RPCError{
			// 		Code:    json.ErrRPCDatabase,
			// 		Message: "Database error: " + err.ScriptError(),
			// 	}
			// }
			if again {
				continue
			}
			break
		}
	loopHashList:
		for i := range hashList {
			var blk *block.Block
			blk, e = chain.BlockByHash(&hashList[i])
			if e != nil {
				// Only handle reorgs if a block could not be found for the hash.
				if dbErr, ok := e.(database.DBError); !ok ||
					dbErr.ErrorCode != database.ErrBlockNotFound {
					E.Ln("error looking up block:", e)
					return nil, &btcjson.RPCError{
						Code: btcjson.ErrRPCDatabase,
						Message: "Database error: " +
							e.Error(),
					}
				}
				// If an absolute max block was specified, don't attempt to handle the reorg.
				if maxBlock != math.MaxInt32 {
					E.Ln("stopping rescan for reorged block", cmd.EndBlock)
					return nil, &ErrRescanReorg
				}
				// If the lookup for the previously valid block hash failed, there may have been a reorg. Fetch a new
				// range of block hashes and verify that the previously processed block (if there was any) still exists
				// in the database. If it doesn't, we error.
				//
				// A goto is used to branch execution back to before the range was evaluated, as it must be
				// reevaluated for the new hashList.
				minBlock += int32(i)
				hashList, e = RecoverFromReorg(
					chain,
					minBlock, maxBlock, lastBlockHash,
				)
				if e != nil {
					return nil, e
				}
				if len(hashList) == 0 {
					break fetchRange
				}
				goto loopHashList
			}
			if i == 0 && lastBlockHash != nil {
				// Ensure the new hashList is on the same fork as the last block from the old hashList.
				jsonErr := DescendantBlock(lastBlockHash, blk)
				if jsonErr != nil {
					return nil, jsonErr
				}
			}
			// A select statement is used to stop rescans if the client requesting the rescan has disconnected.
			select {
			case <-wsc.Quit.Wait():
				D.F("stopped rescan at height %v for disconnected client", blk.Height())
				return nil, nil
			default:
				RescanBlock(wsc, &lookups, blk)
				lastBlock = blk
				lastBlockHash = blk.Hash()
			}
			// Periodically notify the client of the progress completed. Continue with next block if no progress
			// notification is needed yet.
			select {
			case <-ticker.C: // fallthrough
			default:
				continue
			}
			n := btcjson.NewRescanProgressNtfn(
				hashList[i].String(),
				blk.Height(), blk.WireBlock().Header.Timestamp.Unix(),
			)
			mn, e := btcjson.MarshalCmd(nil, n)
			if e != nil {
				E.F(
					"failed to marshal rescan progress notification: %v",
					e,
				)
				continue
			}
			if e = wsc.QueueNotification(mn); e == ErrClientQuit {
				// Finished if the client disconnected.
				D.F(
					"stopped rescan at height %v for disconnected client",
					blk.Height(),
				)
				return nil, nil
			}
		}
		minBlock += int32(len(hashList))
	}
	// Notify websocket client of the finished rescan. Due to how pod asynchronously queues notifications to not block
	// calling code, there is no guarantee that any of the notifications created during rescan (such as rescanprogress,
	// recvtx and redeemingtx) will be received before the rescan RPC returns.
	//
	// Therefore, another method is needed to safely inform clients that all rescan notifications have been sent.
	n := btcjson.NewRescanFinishedNtfn(
		lastBlockHash.String(),
		lastBlock.Height(),
		lastBlock.WireBlock().Header.Timestamp.Unix(),
	)
	if mn, e := btcjson.MarshalCmd(nil, n); E.Chk(e) {
		E.F(
			"failed to marshal rescan finished notification: %v", e,
		)
	} else {
		// The rescan is finished, so we don't care whether the client has disconnected at this point, so discard error.
		_ = wsc.QueueNotification(mn)
	}
	I.Ln("finished rescan")
	return nil, nil
}

// HandleRescanBlocks implements the rescanblocks command extension for websocket connections.
//
// NOTE: This extension is ported from github.com/decred/dcrd
func HandleRescanBlocks(wsc *WSClient, icmd interface{}) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.RescanBlocksCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	// Load client's transaction filter.  Must exist in order to continue.
	wsc.Lock()
	filter := wsc.FilterData
	wsc.Unlock()
	if filter == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCMisc,
			Message: "Transaction filter must be loaded before rescanning",
		}
	}
	blockHashes := make([]*chainhash.Hash, len(cmd.BlockHashes))
	for i := range cmd.BlockHashes {
		hash, e := chainhash.NewHashFromStr(cmd.BlockHashes[i])
		if e != nil {
			return nil, e
		}
		blockHashes[i] = hash
	}
	discoveredData := make([]btcjson.RescannedBlock, 0, len(blockHashes))
	// Iterate over each block in the request and rescan. When a block contains relevant transactions, add it to the
	// response.
	bc := wsc.Server.Cfg.Chain
	params := wsc.Server.Cfg.ChainParams
	var lastBlockHash *chainhash.Hash
	for i := range blockHashes {
		block, e := bc.BlockByHash(blockHashes[i])
		if e != nil {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCBlockNotFound,
				Message: "Failed to fetch block: " + e.Error(),
			}
		}
		if lastBlockHash != nil && block.WireBlock().Header.
			PrevBlock != *lastBlockHash {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidParameter,
				Message: fmt.Sprintf(
					"Block %v is not a child of %v",
					blockHashes[i], lastBlockHash,
				),
			}
		}
		lastBlockHash = blockHashes[i]
		transactions := RescanBlockFilter(filter, block, params)
		if len(transactions) != 0 {
			discoveredData = append(
				discoveredData, btcjson.RescannedBlock{
					Hash:         cmd.BlockHashes[i],
					Transactions: transactions,
				},
			)
		}
	}
	return &discoveredData, nil
}

// HandleSession implements the session command extension for websocket connections.
func HandleSession(wsc *WSClient, icmd interface{}) (interface{}, error) {
	return &btcjson.SessionResult{SessionID: wsc.SessionID}, nil
}

// HandleStopNotifyBlocks implements the stopnotifyblocks command extension for websocket connections.
func HandleStopNotifyBlocks(wsc *WSClient, icmd interface{}) (
	interface{},
	error,
) {
	wsc.Server.NtfnMgr.UnregisterBlockUpdates(wsc)
	return nil, nil
}

// HandleStopNotifyNewTransactions implements the stopnotifynewtransactions command extension for websocket connections.
func HandleStopNotifyNewTransactions(
	wsc *WSClient, icmd interface{},
) (interface{}, error) {
	wsc.Server.NtfnMgr.UnregisterNewMempoolTxsUpdates(wsc)
	return nil, nil
}

// HandleStopNotifyReceived implements the stopnotifyreceived command extension for websocket connections.
func HandleStopNotifyReceived(wsc *WSClient, icmd interface{}) (
	interface{},
	error,
) {
	cmd, ok := icmd.(*btcjson.StopNotifyReceivedCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	// Decode addresses to validate input, but the strings slice is used directly if these are all ok.
	e := CheckAddressValidity(cmd.Addresses, wsc.Server.Cfg.ChainParams)
	if e != nil {
		return nil, e
	}
	for _, addr := range cmd.Addresses {
		wsc.Server.NtfnMgr.UnregisterTxOutAddressRequest(wsc, addr)
	}
	return nil, nil
}

// HandleStopNotifySpent implements the stopnotifyspent command extension for websocket connections.
func HandleStopNotifySpent(wsc *WSClient, icmd interface{}) (
	interface{}, error,
) {
	cmd, ok := icmd.(*btcjson.StopNotifySpentCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	outpoints, e := DeserializeOutpoints(cmd.OutPoints)
	if e != nil {
		return nil, e
	}
	for _, outpoint := range outpoints {
		wsc.Server.NtfnMgr.UnregisterSpentRequest(wsc, outpoint)
	}
	return nil, nil
}

// HandleWebsocketHelp implements the help command for websocket connections.
func HandleWebsocketHelp(wsc *WSClient, icmd interface{}) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.HelpCmd)
	if !ok {
		return nil, btcjson.ErrRPCInternal
	}
	// Provide a usage overview of all commands when no specific command was specified.
	var command string
	if cmd.Command != nil {
		command = *cmd.Command
	}
	if command == "" {
		usage, e := wsc.Server.HelpCacher.RPCUsage(true)
		if e != nil {
			context := "Failed to generate RPC usage"
			return nil, InternalRPCError(e.Error(), context)
		}
		return usage, nil
	}
	// Chk that the command asked for is supported and implemented. Search the list of websocket handlers as well as
	// the main list of handlers since help should only be provided for those cases.
	valid := true
	if _, ok := RPCHandlers[command]; !ok {
		if _, ok := WSHandlers[command]; !ok {
			valid = false
		}
	}
	if !valid {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Unknown command: " + command,
		}
	}
	// Get the help for the command.
	help, e := wsc.Server.HelpCacher.RPCMethodHelp(command)
	if e != nil {
		context := "Failed to generate help"
		return nil, InternalRPCError(e.Error(), context)
	}
	return help, nil
}

func init() {

	WSHandlers = WSHandlersBeforeInit
}
func MakeSemaphore(n int) Semaphore {
	return Semaphore(qu.Ts(n))
}

// NewRedeemingTxNotification returns a new marshalled redeemingtx notification with the passed parameters.
func NewRedeemingTxNotification(
	txHex string, index int,
	block *block.Block,
) ([]byte, error) {
	// Create and marshal the notification.
	ntfn := btcjson.NewRedeemingTxNtfn(txHex, BlockDetails(block, index))
	return btcjson.MarshalCmd(nil, ntfn)
}

// NewWSClientFilter creates a new, empty wsClientFilter struct to be used for a websocket client.
//
// NOTE: This extension was ported from github.com/decred/ dcrd
func NewWSClientFilter(
	addresses []string, unspentOutPoints []wire.OutPoint,
	params *chaincfg.Params,
) *WSClientFilter {
	filter := &WSClientFilter{
		PubKeyHashes:        map[[ripemd160.Size]byte]struct{}{},
		ScriptHashes:        map[[ripemd160.Size]byte]struct{}{},
		CompressedPubKeys:   map[[33]byte]struct{}{},
		UncompressedPubKeys: map[[65]byte]struct{}{},
		OtherAddresses:      map[string]struct{}{},
		Unspent: make(
			map[wire.OutPoint]struct{},
			len(unspentOutPoints),
		),
	}
	for _, s := range addresses {
		filter.AddAddressStr(s, params)
	}
	for i := range unspentOutPoints {
		filter.AddUnspentOutPoint(&unspentOutPoints[i])
	}
	return filter
}

// NewWebsocketClient returns a new websocket client given the notification manager, websocket connection, remote
// address, and whether or not the client has already been authenticated (via HTTP Basic access authentication). The
// returned client is ready to start.
//
// Once started, the client will process incoming and outgoing messages in separate goroutines complete with queuing and
// asynchrous handling for long-running operations.
func NewWebsocketClient(
	server *Server, conn *websocket.Conn,
	remoteAddr string, authenticated bool, isAdmin bool,
) (*WSClient, error) {
	sessionID, e := wire.RandomUint64()
	if e != nil {
		return nil, e
	}
	client := &WSClient{
		Conn:              conn,
		Addr:              remoteAddr,
		Authenticated:     authenticated,
		IsAdmin:           isAdmin,
		SessionID:         sessionID,
		Server:            server,
		AddrRequests:      make(map[string]struct{}),
		SpentRequests:     make(map[wire.OutPoint]struct{}),
		ServiceRequestSem: MakeSemaphore(server.Config.RPCMaxConcurrentReqs.V()),
		NtfnChan:          make(chan []byte, 1), // nonblocking sync
		SendChan:          make(chan WSResponse, WebsocketSendBufferSize),
		Quit:              qu.T(),
	}
	return client, nil
}

// NewWSNotificationManager returns a new notification manager ready for use. See wsNotificationManager for more
// details.
func NewWSNotificationManager(server *Server) *WSNtfnMgr {
	return &WSNtfnMgr{
		Server:            server,
		QueueNotification: make(chan interface{}),
		NotificationMsgs:  make(chan interface{}),
		NumClients:        make(chan int),
		Quit:              qu.T(),
	}
}

// QueueHandler manages a queue of empty interfaces, reading from in and sending the oldest unsent to out.
//
// This handler stops when either of the in or quit channels are closed, and closes out before returning, without
// waiting to send any variables still remaining in the queue.
func QueueHandler(in <-chan interface{}, out chan<- interface{}, quit qu.C) {
	var q []interface{}
	var dequeue chan<- interface{}
	skipQueue := out
	var next interface{}
out:
	for {
		select {
		case n, ok := <-in:
			if !ok {
				// Sender closed input channel.
				break out
			}
			// Either send to out immediately if skipQueue is non-nil (queue is empty) and reader is ready, or append to
			// the queue and send later.
			select {
			case skipQueue <- n:
			default:
				q = append(q, n)
				dequeue = out
				skipQueue = nil
				next = q[0]
			}
		case dequeue <- next:
			copy(q, q[1:])
			q[len(q)-1] = nil // avoid leak
			q = q[:len(q)-1]
			if len(q) == 0 {
				dequeue = nil
				skipQueue = out
			} else {
				next = q[0]
			}
		case <-quit.Wait():
			break out
		}
	}
	close(out)
}

// RecoverFromReorg attempts to recover from a detected reorganize during a rescan.
//
// It fetches a new range of block shas from the database and verifies that the new range of blocks is on the same fork
// as a previous range of blocks.
//
// If this condition does not hold true, the JSON-RPC error for an unrecoverable reorganize is returned.
func RecoverFromReorg(
	chain *blockchain.BlockChain, minBlock, maxBlock int32,
	lastBlock *chainhash.Hash,
) ([]chainhash.Hash, error) {
	hashList, e := chain.HeightRange(minBlock, maxBlock)
	if e != nil {
		E.Ln("error looking up block range:", e)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDatabase,
			Message: "Database error: " + e.Error(),
		}
	}
	if lastBlock == nil || len(hashList) == 0 {
		return hashList, nil
	}
	blk, e := chain.BlockByHash(&hashList[0])
	if e != nil {
		E.Ln("error looking up possibly reorged block:", e)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDatabase,
			Message: "Database error: " + e.Error(),
		}
	}
	jsonErr := DescendantBlock(lastBlock, blk)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return hashList, nil
}

// RescanBlock rescans all transactions in a single block. This is a helper function for handleRescan.
func RescanBlock(wsc *WSClient, lookups *RescanKeys, blk *block.Block) {
	for _, tx := range blk.Transactions() {
		// Hexadecimal representation of this tx. Only created if needed, and reused for later notifications if already
		// made.
		var txHex string
		// All inputs and outputs must be iterated through to correctly modify the unspent map, however, just a single
		// notification for any matching transaction inputs or outputs should be created and sent.
		spentNotified := false
		recvNotified := false
		for _, txin := range tx.MsgTx().TxIn {
			if _, ok := lookups.Unspent[txin.PreviousOutPoint]; ok {
				delete(lookups.Unspent, txin.PreviousOutPoint)
				if spentNotified {
					continue
				}
				if txHex == "" {
					txHex = TxHexString(tx.MsgTx())
				}
				marshalledJSON, e := NewRedeemingTxNotification(
					txHex,
					tx.Index(), blk,
				)
				if e != nil {
					E.Ln(
						"failed to marshal redeemingtx notification:",
						e,
					)
					continue
				}
				e = wsc.QueueNotification(marshalledJSON)
				// Stop the rescan early if the websocket client disconnected.
				if e == ErrClientQuit {
					return
				}
				spentNotified = true
			}
		}
		for txOutIdx, txout := range tx.MsgTx().TxOut {
			_, addrs, _, _ := txscript.ExtractPkScriptAddrs(
				txout.PkScript, wsc.Server.Cfg.ChainParams,
			)
			for _, addr := range addrs {
				switch a := addr.(type) {
				case *btcaddr.PubKeyHash:
					if _, ok := lookups.PubKeyHashes[*a.Hash160()]; !ok {
						continue
					}
				case *btcaddr.ScriptHash:
					if _, ok := lookups.ScriptHashes[*a.Hash160()]; !ok {
						continue
					}
				case *btcaddr.PubKey:
					found := false
					switch sa := a.ScriptAddress(); len(sa) {
					case 33: // Compressed
						var key [33]byte
						copy(key[:], sa)
						if _, ok := lookups.CompressedPubKeys[key]; ok {
							found = true
						}
					case 65: // Uncompressed
						var key [65]byte
						copy(key[:], sa)
						if _, ok := lookups.UncompressedPubKeys[key]; ok {
							found = true
						}
					default:
						W.F("skipping rescanned pubkey of unknown serialized length", len(sa))
						continue
					}
					// If the transaction output pays to the pubkey of a rescanned P2PKH address, include it as well.
					if !found {
						pkh := a.PubKeyHash()
						if _, ok := lookups.PubKeyHashes[*pkh.Hash160()]; !ok {
							continue
						}
					}
				default:
					// A new address type must have been added. Encode as a payment address string and check the
					// fallback map.
					addrStr := addr.EncodeAddress()
					_, ok := lookups.Fallbacks[addrStr]
					if !ok {
						continue
					}
				}
				outpoint := wire.OutPoint{
					Hash:  *tx.Hash(),
					Index: uint32(txOutIdx),
				}
				lookups.Unspent[outpoint] = struct{}{}
				if recvNotified {
					continue
				}
				if txHex == "" {
					txHex = TxHexString(tx.MsgTx())
				}
				ntfn := btcjson.NewRecvTxNtfn(
					txHex,
					BlockDetails(blk, tx.Index()),
				)
				marshalledJSON, e := btcjson.MarshalCmd(nil, ntfn)
				if e != nil {
					E.Ln("failed to marshal recvtx notification:", e)
					return
				}
				e = wsc.QueueNotification(marshalledJSON)
				// Stop the rescan early if the websocket client disconnected.
				if e == ErrClientQuit {
					return
				}
				recvNotified = true
			}
		}
	}
}

// RescanBlockFilter rescans a block for any relevant transactions for the passed lookup keys. Any discovered
// transactions are returned hex encoded as a string slice.
//
// NOTE: This extension is ported from github.com/decred/dcrd
func RescanBlockFilter(
	filter *WSClientFilter, block *block.Block,
	params *chaincfg.Params,
) []string {
	var transactions []string
	filter.mu.Lock()
	for _, tx := range block.Transactions() {
		msgTx := tx.MsgTx()
		// Keep track of whether the transaction has already been added to the result. It shouldn't be added twice.
		added := false
		// Scan inputs if not a coinbase transaction.
		if !blockchain.IsCoinBaseTx(msgTx) {
			for _, input := range msgTx.TxIn {
				if !filter.ExistsUnspentOutPoint(&input.PreviousOutPoint) {
					continue
				}
				if !added {
					transactions = append(
						transactions,
						TxHexString(msgTx),
					)
					added = true
				}
			}
		}
		var e error
		// Scan outputs.
		for i, output := range msgTx.TxOut {
			var addrs []btcaddr.Address
			_, addrs, _, e = txscript.ExtractPkScriptAddrs(
				output.PkScript, params,
			)
			if e != nil {
				continue
			}
			for _, a := range addrs {
				if !filter.ExistsAddress(a) {
					continue
				}
				op := wire.OutPoint{
					Hash:  *tx.Hash(),
					Index: uint32(i),
				}
				filter.AddUnspentOutPoint(&op)
				if !added {
					transactions = append(
						transactions,
						TxHexString(msgTx),
					)
					added = true
				}
			}
		}
	}
	filter.mu.Unlock()
	return transactions
}

// TxHexString returns the serialized transaction encoded in hexadecimal.
func TxHexString(tx *wire.MsgTx) string {
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	// Ignore Serialize's error, as writing to a bytes.buffer cannot fail.
	e := tx.Serialize(buf)
	if e != nil {
	}
	return hex.EncodeToString(buf.Bytes())
}
