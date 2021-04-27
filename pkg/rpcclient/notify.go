package rpcclient

import (
	"bytes"
	"encoding/hex"
	js "encoding/json"
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"time"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

var (
	// ErrWebsocketsRequired is an error to describe the condition where the caller is trying to use a websocket-only
	// feature, such as requesting notifications or other websocket requests when the client is configured to run in
	// HTTP POST mode.
	ErrWebsocketsRequired = errors.New(
		"a websocket connection is required " +
			"to use this feature",
	)
)

// notificationState is used to track the current state of successfully registered notification so the state can be
// automatically re-established on reconnect.
type notificationState struct {
	notifyBlocks       bool
	notifyNewTx        bool
	notifyNewTxVerbose bool
	notifyReceived     map[string]struct{}
	notifySpent        map[btcjson.OutPoint]struct{}
}

// Copy returns a deep copy of the receiver.
func (s *notificationState) Copy() *notificationState {
	var stateCopy notificationState
	stateCopy.notifyBlocks = s.notifyBlocks
	stateCopy.notifyNewTx = s.notifyNewTx
	stateCopy.notifyNewTxVerbose = s.notifyNewTxVerbose
	stateCopy.notifyReceived = make(map[string]struct{})
	for addr := range s.notifyReceived {
		stateCopy.notifyReceived[addr] = struct{}{}
	}
	stateCopy.notifySpent = make(map[btcjson.OutPoint]struct{})
	for op := range s.notifySpent {
		stateCopy.notifySpent[op] = struct{}{}
	}
	return &stateCopy
}

func // newNotificationState returns a new notification state ready to be
// populated.
newNotificationState() *notificationState {
	return &notificationState{
		notifyReceived: make(map[string]struct{}),
		notifySpent:    make(map[btcjson.OutPoint]struct{}),
	}
}

func // newNilFutureResult returns a new future result channel that already
// has the result waiting on the channel with the reply set to nil.
// This is useful to ignore things such as notifications when the caller didn't
// specify any notification handlers.
newNilFutureResult() chan *response {
	responseChan := make(chan *response, 1)
	responseChan <- &response{result: nil, err: nil}
	return responseChan
}

// NotificationHandlers defines callback function pointers to invoke with
// notifications. Since all of the functions are nil by default, all notifications are effectively ignored until
// their handlers are set to a concrete callback.
//
// NOTE: Unless otherwise documented, these handlers must NOT directly call any blocking calls on the client
// instance since the input reader goroutine blocks until the callback has completed. Doing so will result in a
// deadlock situation.
type NotificationHandlers struct {
	// OnClientConnected is invoked when the client connects or reconnects to the RPC server. This callback is run async
	// with the rest of the notification handlers, and is safe for blocking client requests.
	OnClientConnected func()
	// OnBlockConnected is invoked when a block is connected to the longest (best) chain. It will only be invoked if a
	// preceding call to NotifyBlocks has been made to register for the notification and the function is non-nil. NOTE:
	// Deprecated. Use OnFilteredBlockConnected instead.
	OnBlockConnected func(hash *chainhash.Hash, height int32, t time.Time)
	// OnFilteredBlockConnected is invoked when a block is connected to the longest (best) chain. It will only be
	// invoked if a preceding call to NotifyBlocks has been made to register for the notification and the function is
	// non-nil. Its parameters differ from OnBlockConnected: it receives the block's height, header, and relevant
	// transactions.
	OnFilteredBlockConnected func(height int32, header *wire.BlockHeader, txs []*util.Tx)
	// OnBlockDisconnected is invoked when a block is disconnected from the longest (best) chain. It will only be
	// invoked if a preceding call to NotifyBlocks has been made to register for the notification and the function is
	// non-nil. NOTE: Deprecated. Use OnFilteredBlockDisconnected instead.
	OnBlockDisconnected func(hash *chainhash.Hash, height int32, t time.Time)
	// OnFilteredBlockDisconnected is invoked when a block is disconnected from the longest (best) chain. It will only
	// be invoked if a preceding NotifyBlocks has been made to register for the notification and the call to function is
	// non-nil. Its parameters differ from OnBlockDisconnected: it receives the block's height and header.
	OnFilteredBlockDisconnected func(height int32, header *wire.BlockHeader)
	// OnRecvTx is invoked when a transaction that receives funds to a registered address is received into the memory
	// pool and also connected to the longest (best) chain. It will only be invoked if a preceding call to
	// NotifyReceived, Rescan, or RescanEndHeight has been made to register for the notification and the function is
	// non-nil. NOTE: Deprecated. Use OnRelevantTxAccepted instead.
	OnRecvTx func(transaction *util.Tx, details *btcjson.BlockDetails)
	// OnRedeemingTx is invoked when a transaction that spends a registered outpoint is received into the memory pool
	// and also connected to the longest (best) chain.
	//
	// It will only be invoked if a preceding call to NotifySpent, Rescan, or RescanEndHeight has been made to register
	// for the notification and the function is non-nil.
	//
	// NOTE: The NotifyReceived will automatically register notifications for the outpoints that are now "owned" as a
	// result of receiving funds to the registered addresses.
	//
	// This means it is possible for this to invoked indirectly as the result of a NotifyReceived call. NOTE:
	// Deprecated. Use OnRelevantTxAccepted instead.
	OnRedeemingTx func(transaction *util.Tx, details *btcjson.BlockDetails)
	// OnRelevantTxAccepted is invoked when an unmined transaction passes the client's transaction filter.
	//
	// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
	OnRelevantTxAccepted func(transaction []byte)
	// OnRescanFinished is invoked after a rescan finishes due to a previous call to Rescan or RescanEndHeight. Finished
	// rescans should be signaled on this notification, rather than relying on the return result of a rescan request,
	// due to how pod may send various rescan notifications after the rescan request has already returned.
	//
	// NOTE: Deprecated. Not used with RescanBlocks.
	OnRescanFinished func(hash *chainhash.Hash, height int32, blkTime time.Time)
	// OnRescanProgress is invoked periodically when a rescan is underway. It will only be invoked if a preceding call
	// to Rescan or RescanEndHeight has been made and the function is non-nil.
	//
	// NOTE: Deprecated. Not used with RescanBlocks.
	OnRescanProgress func(hash *chainhash.Hash, height int32, blkTime time.Time)
	// OnTxAccepted is invoked when a transaction is accepted into the memory pool. It will only be invoked if a
	// preceding call to NotifyNewTransactions with the verbose flag set to false has been made to register for the
	// notification and the function is non-nil.
	OnTxAccepted func(hash *chainhash.Hash, amount amt.Amount)
	// OnTxAccepted is invoked when a transaction is accepted into the memory pool. It will only be invoked if a
	// preceding call to NotifyNewTransactions with the verbose flag set to true has been made to register for the
	// notification and the function is non-nil.
	OnTxAcceptedVerbose func(txDetails *btcjson.TxRawResult)
	// OnPodConnected is invoked when a wallet connects or disconnects from pod. This will only be available when client
	// is connected to a wallet server such as btcwallet.
	OnPodConnected func(connected bool)
	// OnAccountBalance is invoked with account balance updates. This will only be available when speaking to a wallet
	// server such as btcwallet.
	OnAccountBalance func(account string, balance amt.Amount, confirmed bool)
	// OnWalletLockState is invoked when a wallet is locked or unlocked. This will only be available when client is
	// connected to a wallet server such as btcwallet.
	OnWalletLockState func(locked bool)
	// OnUnknownNotification is invoked when an unrecognized notification is received. This typically means the
	// notification handling code for this package needs to be updated for a new notification type or the caller is
	// using a custom notification this package does not know about.
	OnUnknownNotification func(method string, params []js.RawMessage)
}

// handleNotification examines the passed notification type, performs conversions to get the raw notification types into
// higher level types and delivers the notification to the appropriate On<X> handler registered with the client.
func (c *Client) handleNotification(ntfn *rawNotification) {
	D.Ln("<<<Handling Notification>>>", ntfn.Method)
	// Ignore the notification if the client is not interested in any notifications.
	if c.ntfnHandlers == nil {
		D.Ln("<<<no notification handlers registered>>>")
		return
	}
	switch ntfn.Method {
	// OnBlockConnected
	case btcjson.BlockConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnBlockConnected == nil {
			D.Ln("<<<no OnBlockConnected callback registered>>>")
			return
		}
		blockHash, blockHeight, blockTime, e := parseChainNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid block connected notification:", e)
			return
		}
		c.ntfnHandlers.OnBlockConnected(blockHash, blockHeight, blockTime)
	// OnFilteredBlockConnected
	case btcjson.FilteredBlockConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnFilteredBlockConnected == nil {
			D.Ln("<<<no OnFilteredBlockConnected callback registered>>>")
			return
		}
		blockHeight, blockHeader, transactions, e :=
			parseFilteredBlockConnectedParams(ntfn.Params)
		if e != nil {
			W.Ln(
				"received invalid filtered block connected notification:",
				e,
			)
			return
		}
		c.ntfnHandlers.OnFilteredBlockConnected(
			blockHeight,
			blockHeader, transactions,
		)
	// OnBlockDisconnected
	case btcjson.BlockDisconnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnBlockDisconnected == nil {
			D.Ln("<<<no OnBlockDisconnected callback registered>>>")
			return
		}
		blockHash, blockHeight, blockTime, e := parseChainNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid block connected notification:", e)
			return
		}
		c.ntfnHandlers.OnBlockDisconnected(blockHash, blockHeight, blockTime)
	// OnFilteredBlockDisconnected
	case btcjson.FilteredBlockDisconnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnFilteredBlockDisconnected == nil {
			D.Ln("<<<no OnFilteredBlockDisconnected callback registered>>>")
			return
		}
		blockHeight, blockHeader, e := parseFilteredBlockDisconnectedParams(ntfn.Params)
		if e != nil {
			W.Ln(
				"received invalid filtered block disconnected"+
					" notification"+
					":", e,
			)
			return
		}
		c.ntfnHandlers.OnFilteredBlockDisconnected(blockHeight, blockHeader)
	// OnRecvTx
	case btcjson.RecvTxNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRecvTx == nil {
			D.Ln("<<<no OnRecvTx callback registered>>>")
			return
		}
		tx, block, e := parseChainTxNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid recvtx notification:", e)
			return
		}
		c.ntfnHandlers.OnRecvTx(tx, block)
	// OnRedeemingTx
	case btcjson.RedeemingTxNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRedeemingTx == nil {
			D.Ln("<<<no OnRedeemingTx callback registered>>>")
			return
		}
		tx, block, e := parseChainTxNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid redeemingtx notification:", e)
			return
		}
		c.ntfnHandlers.OnRedeemingTx(tx, block)
	// OnRelevantTxAccepted
	case btcjson.RelevantTxAcceptedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRelevantTxAccepted == nil {
			D.Ln("<<<no OnRelevantTxAccepted callback registered>>>")
			return
		}
		transaction, e := parseRelevantTxAcceptedParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid relevanttxaccepted notification:", e)
			return
		}
		c.ntfnHandlers.OnRelevantTxAccepted(transaction)
	// OnRescanFinished
	case btcjson.RescanFinishedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRescanFinished == nil {
			D.Ln("<<<no OnRescanFinished callback registered>>>")
			return
		}
		hash, height, blkTime, e := parseRescanProgressParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid rescanfinished notification:", e)
			return
		}
		c.ntfnHandlers.OnRescanFinished(hash, height, blkTime)
	// OnRescanProgress
	case btcjson.RescanProgressNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRescanProgress == nil {
			D.Ln("<<<no OnRescanProgress callback registered>>>")
			return
		}
		hash, height, blkTime, e := parseRescanProgressParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid rescanprogress notification:", e)
			return
		}
		c.ntfnHandlers.OnRescanProgress(hash, height, blkTime)
	// OnTxAccepted
	case btcjson.TxAcceptedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnTxAccepted == nil {
			D.Ln("<<<no OnTxAccepted callback registered>>>")
			return
		}
		hash, amt, e := parseTxAcceptedNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid tx accepted notification:", e)
			return
		}
		c.ntfnHandlers.OnTxAccepted(hash, amt)
	// OnTxAcceptedVerbose
	case btcjson.TxAcceptedVerboseNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnTxAcceptedVerbose == nil {
			D.Ln("<<<no OnTxAcceptedVerbose callback registered>>>")
			return
		}
		rawTx, e := parseTxAcceptedVerboseNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid tx accepted verbose notification:", e)
			return
		}
		c.ntfnHandlers.OnTxAcceptedVerbose(rawTx)
	// OnPodConnected
	case btcjson.PodConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnPodConnected == nil {
			D.Ln("<<<no OnPodConnected callback registered>>>")
			return
		}
		connected, e := parsePodConnectedNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid pod connected notification:", e)
			return
		}
		c.ntfnHandlers.OnPodConnected(connected)
	// OnAccountBalance
	case btcjson.AccountBalanceNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnAccountBalance == nil {
			D.Ln("<<<no OnAccountBalance callback registered>>>")
			return
		}
		account, bal, conf, e := parseAccountBalanceNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid account balance notification:", e)
			return
		}
		c.ntfnHandlers.OnAccountBalance(account, bal, conf)
	// OnWalletLockState
	case btcjson.WalletLockStateNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnWalletLockState == nil {
			D.Ln("<<<no OnWalletLockState callback registered>>>")
			return
		}
		// The account name is not notified, so the return value is discarded.
		_, locked, e := parseWalletLockStateNtfnParams(ntfn.Params)
		if e != nil {
			W.Ln("received invalid wallet lock state notification:", e)
			return
		}
		c.ntfnHandlers.OnWalletLockState(locked)
	// OnUnknownNotification
	default:
		if c.ntfnHandlers.OnUnknownNotification == nil {
			D.Ln("<<<no OnUnknownNotification callback registered>>>")
			return
		}
		c.ntfnHandlers.OnUnknownNotification(ntfn.Method, ntfn.Params)
	}
}

// wrongNumParams is an error type describing an unparseable JSON-RPC notification due to an incorrect number of
// parameters for the expected notification type.
//
// The value is the number of parameters of the invalid notification.
type wrongNumParams int

// BTCJSONError satisfies the builtin error interface.
func (e wrongNumParams) Error() string {
	return fmt.Sprintf("wrong number of parameters (%d)", e)
}

// parseChainNtfnParams parses out the block hash and height from the parameters of blockconnected and blockdisconnected
// notifications.
func parseChainNtfnParams(params []js.RawMessage) (
	*chainhash.Hash,
	int32, time.Time, error,
) {
	if len(params) != 3 {
		return nil, 0, time.Time{}, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a string.
	var blockHashStr string
	e := js.Unmarshal(params[0], &blockHashStr)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Unmarshal second parameter as an integer.
	var blockHeight int32
	e = js.Unmarshal(params[1], &blockHeight)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Unmarshal third parameter as unix time.
	var blockTimeUnix int64
	e = js.Unmarshal(params[2], &blockTimeUnix)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Create hash from block hash string.
	blockHash, e := chainhash.NewHashFromStr(blockHashStr)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Create time.Time from unix time.
	blockTime := time.Unix(blockTimeUnix, 0)
	return blockHash, blockHeight, blockTime, nil
}

// parseFilteredBlockConnectedParams parses out the parameters included in a filteredblockconnected notification.
//
// NOTE: This is a pod extension ported from github. com/decred/dcrrpcclient and requires a websocket connection.
func parseFilteredBlockConnectedParams(params []js.RawMessage) (
	int32,
	*wire.BlockHeader, []*util.Tx, error,
) {
	if len(params) < 3 {
		return 0, nil, nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an integer.
	var blockHeight int32
	e := js.Unmarshal(params[0], &blockHeight)
	if e != nil {
		return 0, nil, nil, e
	}
	// Unmarshal second parameter as a slice of bytes.
	blockHeaderBytes, e := parseHexParam(params[1])
	if e != nil {
		return 0, nil, nil, e
	}
	// Deserialize block header from slice of bytes.
	var blockHeader wire.BlockHeader
	e = blockHeader.Deserialize(bytes.NewReader(blockHeaderBytes))
	if e != nil {
		return 0, nil, nil, e
	}
	// Unmarshal third parameter as a slice of hex-encoded strings.
	var hexTransactions []string
	e = js.Unmarshal(params[2], &hexTransactions)
	if e != nil {
		return 0, nil, nil, e
	}
	// Create slice of transactions from slice of strings by hex-decoding.
	transactions := make([]*util.Tx, len(hexTransactions))
	for i, hexTx := range hexTransactions {
		transaction, e := hex.DecodeString(hexTx)
		if e != nil {
			return 0, nil, nil, e
		}
		transactions[i], e = util.NewTxFromBytes(transaction)
		if e != nil {
			return 0, nil, nil, e
		}
	}
	return blockHeight, &blockHeader, transactions, nil
}

// parseFilteredBlockDisconnectedParams parses out the parameters included in a filteredblockdisconnected notification.
//
// NOTE: This is a pod extension ported from github. com/decred/dcrrpcclient and requires a websocket connection.
func parseFilteredBlockDisconnectedParams(params []js.RawMessage) (
	int32,
	*wire.BlockHeader, error,
) {
	if len(params) < 2 {
		return 0, nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an integer.
	var blockHeight int32
	e := js.Unmarshal(params[0], &blockHeight)
	if e != nil {
		return 0, nil, e
	}
	// Unmarshal second parameter as a slice of bytes.
	blockHeaderBytes, e := parseHexParam(params[1])
	if e != nil {
		return 0, nil, e
	}
	// Deserialize block header from slice of bytes.
	var blockHeader wire.BlockHeader
	e = blockHeader.Deserialize(bytes.NewReader(blockHeaderBytes))
	if e != nil {
		return 0, nil, e
	}
	return blockHeight, &blockHeader, nil
}

func parseHexParam(param js.RawMessage) ([]byte, error) {
	var s string
	e := js.Unmarshal(param, &s)
	if e != nil {
		return nil, e
	}
	return hex.DecodeString(s)
}

// parseRelevantTxAcceptedParams parses out the parameter included in a relevanttxaccepted notification.
func parseRelevantTxAcceptedParams(params []js.RawMessage) (transaction []byte, e error) {
	if len(params) < 1 {
		return nil, wrongNumParams(len(params))
	}
	return parseHexParam(params[0])
}

// parseChainTxNtfnParams parses out the transaction and optional details about the block it's mined in from the
// parameters of recvtx and redeemingtx notifications.
func parseChainTxNtfnParams(params []js.RawMessage) (
	*util.Tx,
	*btcjson.BlockDetails, error,
) {
	if len(params) == 0 || len(params) > 2 {
		return nil, nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a string.
	var txHex string
	e := js.Unmarshal(params[0], &txHex)
	if e != nil {
		return nil, nil, e
	}
	// If present, unmarshal second optional parameter as the block details JSON object.
	var block *btcjson.BlockDetails
	if len(params) > 1 {
		e = js.Unmarshal(params[1], &block)
		if e != nil {
			return nil, nil, e
		}
	}
	// Hex decode and deserialize the transaction.
	serializedTx, e := hex.DecodeString(txHex)
	if e != nil {
		return nil, nil, e
	}
	var msgTx wire.MsgTx
	e = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if e != nil {
		return nil, nil, e
	}
	// TODO: Change recvtx and redeemingtx callback signatures to use nicer
	//  types for details about the block (block hash as a chainhash.Hash,
	//  block time as a time.Time, etc.).
	return util.NewTx(&msgTx), block, nil
}

// parseRescanProgressParams parses out the height of the last rescanned block from the parameters of rescanfinished and
// rescanprogress notifications.
func parseRescanProgressParams(params []js.RawMessage) (*chainhash.Hash, int32, time.Time, error) {
	if len(params) != 3 {
		return nil, 0, time.Time{}, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an string.
	var hashStr string
	e := js.Unmarshal(params[0], &hashStr)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Unmarshal second parameter as an integer.
	var height int32
	e = js.Unmarshal(params[1], &height)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Unmarshal third parameter as an integer.
	var blkTime int64
	e = js.Unmarshal(params[2], &blkTime)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	// Decode string encoding of block hash.
	hash, e := chainhash.NewHashFromStr(hashStr)
	if e != nil {
		return nil, 0, time.Time{}, e
	}
	return hash, height, time.Unix(blkTime, 0), nil
}

// parseTxAcceptedNtfnParams parses out the transaction hash and total amount from the parameters of a txaccepted
// notification.
func parseTxAcceptedNtfnParams(params []js.RawMessage) (
	*chainhash.Hash,
	amt.Amount, error,
) {
	if len(params) != 2 {
		return nil, 0, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a string.
	var txHashStr string
	e := js.Unmarshal(params[0], &txHashStr)
	if e != nil {
		return nil, 0, e
	}
	// Unmarshal second parameter as a floating point number.
	var fAmt float64
	e = js.Unmarshal(params[1], &fAmt)
	if e != nil {
		return nil, 0, e
	}
	// Bounds check amount.
	amt, e := amt.NewAmount(fAmt)
	if e != nil {
		return nil, 0, e
	}
	// Decode string encoding of transaction sha.
	txHash, e := chainhash.NewHashFromStr(txHashStr)
	if e != nil {
		return nil, 0, e
	}
	return txHash, amt, nil
}

// parseTxAcceptedVerboseNtfnParams parses out details about a raw transaction from the parameters of a
// txacceptedverbose notification.
func parseTxAcceptedVerboseNtfnParams(params []js.RawMessage) (
	*btcjson.TxRawResult,
	error,
) {
	if len(params) != 1 {
		return nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a raw transaction result object.
	var rawTx btcjson.TxRawResult
	e := js.Unmarshal(params[0], &rawTx)
	if e != nil {
		return nil, e
	}
	// TODO: change txacceptedverbose notification callbacks to use nicer
	//  types for all details about the transaction (i.e.
	//  decoding hashes from their string encoding).
	return &rawTx, nil
}

// parsePodConnectedNtfnParams parses out the connection status of pod and btcwallet from the parameters of a
// podconnected notification.
func parsePodConnectedNtfnParams(params []js.RawMessage) (bool, error) {
	if len(params) != 1 {
		return false, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a boolean.
	var connected bool
	e := js.Unmarshal(params[0], &connected)
	if e != nil {
		return false, e
	}
	return connected, nil
}

// parseAccountBalanceNtfnParams parses out the account name, total balance, and whether or not the balance is confirmed
// or unconfirmed from the parameters of an accountbalance notification.
func parseAccountBalanceNtfnParams(params []js.RawMessage) (
	account string,
	balance amt.Amount, confirmed bool, e error,
) {
	if len(params) != 3 {
		return "", 0, false, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a string.
	e = js.Unmarshal(params[0], &account)
	if e != nil {
		return "", 0, false, e
	}
	// Unmarshal second parameter as a floating point number.
	var fBal float64
	e = js.Unmarshal(params[1], &fBal)
	if e != nil {
		return "", 0, false, e
	}
	// Unmarshal third parameter as a boolean.
	e = js.Unmarshal(params[2], &confirmed)
	if e != nil {
		return "", 0, false, e
	}
	// Bounds check amount.
	bal, e := amt.NewAmount(fBal)
	if e != nil {
		return "", 0, false, e
	}
	return account, bal, confirmed, nil
}

// parseWalletLockStateNtfnParams parses out the account name and locked state of an account from the parameters of a
// walletlockstate notification.
func parseWalletLockStateNtfnParams(params []js.RawMessage) (
	account string,
	locked bool, e error,
) {
	if len(params) != 2 {
		return "", false, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as a string.
	e = js.Unmarshal(params[0], &account)
	if e != nil {
		return "", false, e
	}
	// Unmarshal second parameter as a boolean.
	e = js.Unmarshal(params[1], &locked)
	if e != nil {
		return "", false, e
	}
	return account, locked, nil
}

// FutureNotifyBlocksResult is a future promise to deliver the result of a NotifyBlocksAsync RPC invocation (or an
// applicable error).
type FutureNotifyBlocksResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
func (r FutureNotifyBlocksResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// NotifyBlocksAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance.
//
// See NotifyBlocks for the blocking version and more details.
//
// NOTE: This is a pod extension and requires a websocket connection.
func (c *Client) NotifyBlocksAsync() FutureNotifyBlocksResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	cmd := btcjson.NewNotifyBlocksCmd()
	return c.sendCmd(cmd)
}

// NotifyBlocks registers the client to receive notifications when blocks are connected and disconnected from the main
// chain.
//
// The notifications are delivered to the notification handlers associated with the client. Calling this function has no
// effect if there are no notification handlers and will result in an error if the client is configured to run in HTTP
// POST mode.
//
// The notifications delivered as a result of this call will be via one of or OnBlockDisconnected. NOTE: This is a pod
// extension and requires a websocket connection.
func (c *Client) NotifyBlocks() (e error) {
	return c.NotifyBlocksAsync().Receive()
}

// FutureNotifySpentResult is a future promise to deliver the result of a NotifySpentAsync RPC invocation (or an
// applicable error).
//
// NOTE: Deprecated. Use FutureLoadTxFilterResult instead.
type FutureNotifySpentResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
func (r FutureNotifySpentResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// notifySpentInternal is the same as notifySpentAsync except it accepts the converted outpoints as a parameter so the
// client can more efficiently recreate the previous notification state on reconnect.
func (c *Client) notifySpentInternal(outpoints []btcjson.OutPoint) FutureNotifySpentResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	cmd := btcjson.NewNotifySpentCmd(outpoints)
	return c.sendCmd(cmd)
}

// newOutPointFromWire constructs the json representation of a transaction outpoint from the wire type.
func newOutPointFromWire(op *wire.OutPoint) btcjson.OutPoint {
	return btcjson.OutPoint{
		Hash:  op.Hash.String(),
		Index: op.Index,
	}
}

// NotifySpentAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance.
//
// See NotifySpent for the blocking version and more details.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use LoadTxFilterAsync instead.
func (c *Client) NotifySpentAsync(outpoints []*wire.OutPoint) FutureNotifySpentResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	ops := make([]btcjson.OutPoint, 0, len(outpoints))
	for _, outpoint := range outpoints {
		ops = append(ops, newOutPointFromWire(outpoint))
	}
	cmd := btcjson.NewNotifySpentCmd(ops)
	return c.sendCmd(cmd)
}

// NotifySpent registers the client to receive notifications when the passed transaction outputs are spent.
//
// The notifications are delivered to the notification handlers associated with the client. Calling this function has no
// effect if there are no notification handlers and will result in an error if the client is configured to run in HTTP
// POST mode.
//
// The notifications delivered as a result of this call will be via OnRedeemingTx.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use LoadTxFilter instead.
func (c *Client) NotifySpent(outpoints []*wire.OutPoint) (e error) {
	return c.NotifySpentAsync(outpoints).Receive()
}

// FutureNotifyNewTransactionsResult is a future promise to deliver the result of a NotifyNewTransactionsAsync RPC
// invocation (or an applicable error).
type FutureNotifyNewTransactionsResult chan *response

// Receive waits for the response promised by the future and returns an error if the registration was not successful.
func (r FutureNotifyNewTransactionsResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// NotifyNewTransactionsAsync returns an instance of a type that can be used to get the result of the RPC at some future
// time by invoking the Receive function on the returned instance.
//
// See NotifyNewTransactionsAsync for the blocking version and more details.
//
// NOTE: This is a pod extension and requires a websocket connection.
func (c *Client) NotifyNewTransactionsAsync(verbose bool) FutureNotifyNewTransactionsResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	cmd := btcjson.NewNotifyNewTransactionsCmd(&verbose)
	return c.sendCmd(cmd)
}

// NotifyNewTransactions registers the client to receive notifications every time a new transaction is accepted to the
// memory pool.
//
// The notifications are delivered to the notification handlers associated with the client. Calling this function has no
// effect if there are no notification handlers and will result in an error if the client is configured to run in HTTP
// POST mode.
//
// The notifications delivered as a result of this call will be via one of OnTxAccepted (when verbose is false) or
// OnTxAcceptedVerbose ( when verbose is true). NOTE: This is a pod extension and requires a websocket connection.
func (c *Client) NotifyNewTransactions(verbose bool) (e error) {
	return c.NotifyNewTransactionsAsync(verbose).Receive()
}

// FutureNotifyReceivedResult is a future promise to deliver the result of a NotifyReceivedAsync RPC invocation (or an
// applicable error).
//
// NOTE: Deprecated. Use FutureLoadTxFilterResult instead.
type FutureNotifyReceivedResult chan *response

// Receive waits for the response promised by the future and returns an error if the registration was not successful.
func (r FutureNotifyReceivedResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// notifyReceivedInternal is the same as notifyReceivedAsync except it accepts the converted addresses as a parameter so
// the client can more efficiently recreate the previous notification state on reconnect.
func (c *Client) notifyReceivedInternal(addresses []string) FutureNotifyReceivedResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	// Convert addresses to strings.
	cmd := btcjson.NewNotifyReceivedCmd(addresses)
	return c.sendCmd(cmd)
}

// NotifyReceivedAsync returns an instance of a type that can be used to get the result of the RPC at some future time
// by invoking the Receive function on the returned instance.
//
// See NotifyReceived for the blocking version and more details.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use LoadTxFilterAsync instead.
func (c *Client) NotifyReceivedAsync(addresses []btcaddr.Address) FutureNotifyReceivedResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	// Convert addresses to strings.
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}
	cmd := btcjson.NewNotifyReceivedCmd(addrs)
	return c.sendCmd(cmd)
}

// NotifyReceived registers the client to receive notifications every time a new transaction which pays to one of the
// passed addresses is accepted to memory pool or in a block connected to the block chain.
//
// In addition, when one of these transactions is detected, the client is also automatically registered for
// notifications when the new transaction outpoints the address now has available are spent (
//
// See NotifySpent). The notifications are delivered to the notification handlers associated with the client.
//
// Calling this function has no effect if there are no notification handlers and will result in an error if the client
// is configured to run in HTTP POST mode.
//
// The notifications delivered as a result of this call will be via one of *OnRecvTx (for transactions that receive
// funds to one of the passed addresses) or OnRedeemingTx ( for transactions which spend from one of the outpoints which
// are automatically registered upon receipt of funds to the address).
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use LoadTxFilter instead.
func (c *Client) NotifyReceived(addresses []btcaddr.Address) (e error) {
	return c.NotifyReceivedAsync(addresses).Receive()
}

// FutureRescanResult is a future promise to deliver the result of a RescanAsync or RescanEndHeightAsync RPC invocation
// ( or an applicable error).
//
// NOTE: Deprecated. Use FutureRescanBlocksResult instead.
type FutureRescanResult chan *response

// Receive waits for the response promised by the future and returns an error if the rescan was not successful.
func (r FutureRescanResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// RescanAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance.
//
// See Rescan for the blocking version and more details.
//
// NOTE: Rescan requests are not issued on client reconnect and must be performed manually (ideally with a new start
// height based on the last rescan progress notification).
//
// See the OnClientConnected notification callback for a good call site to reissue rescan requests on connect and
// reconnect.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use RescanBlocksAsync instead.
func (c *Client) RescanAsync(
	startBlock *chainhash.Hash,
	addresses []btcaddr.Address,
	outpoints []*wire.OutPoint,
) FutureRescanResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	// Convert block hashes to strings.
	var startBlockHashStr string
	if startBlock != nil {
		startBlockHashStr = startBlock.String()
	}
	// Convert addresses to strings.
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}
	// Convert outpoints.
	ops := make([]btcjson.OutPoint, 0, len(outpoints))
	for _, op := range outpoints {
		ops = append(ops, newOutPointFromWire(op))
	}
	cmd := btcjson.NewRescanCmd(startBlockHashStr, addrs, ops, nil)
	return c.sendCmd(cmd)
}

// Rescan rescans the block chain starting from the provided starting block to the end of the longest chain for
// transactions that pay to the passed addresses and transactions which spend the passed outpoints.
//
// The notifications of found transactions are delivered to the notification handlers associated with client and this
// call will not return until the rescan has completed. Calling this function has no effect if there are no notification
// handlers and will result in an error if the client is configured to run in HTTP POST mode.
//
// The notifications delivered as a result of this call will be via one of OnRedeemingTx (for transactions which spend
// from the one of the passed outpoints), OnRecvTx (for transactions that receive funds to one of the passed addresses),
// and OnRescanProgress (for rescan progress updates).
//
// See RescanEndBlock to also specify an ending block to finish the rescan without continuing through the best block on
// the main chain.
//
// NOTE: Rescan requests are not issued on client reconnect and must be performed manually (ideally with a new start
// height based on the last rescan progress notification).
//
// See the OnClientConnected notification callback for a good call site to reissue rescan requests on connect and
// reconnect.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use RescanBlocks instead.
func (c *Client) Rescan(
	startBlock *chainhash.Hash,
	addresses []btcaddr.Address,
	outpoints []*wire.OutPoint,
) (e error) {
	return c.RescanAsync(startBlock, addresses, outpoints).Receive()
}

// RescanEndBlockAsync returns an instance of a type that can be used to get the result of the RPC at some future time
// by invoking the Receive function on the returned instance.
//
// See RescanEndBlock for the blocking version and more details.
//
// NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use RescanBlocksAsync instead.
func (c *Client) RescanEndBlockAsync(
	startBlock *chainhash.Hash,
	addresses []btcaddr.Address, outpoints []*wire.OutPoint,
	endBlock *chainhash.Hash,
) FutureRescanResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	// Ignore the notification if the client is not interested in notifications.
	if c.ntfnHandlers == nil {
		return newNilFutureResult()
	}
	// Convert block hashes to strings.
	var startBlockHashStr, endBlockHashStr string
	if startBlock != nil {
		startBlockHashStr = startBlock.String()
	}
	if endBlock != nil {
		endBlockHashStr = endBlock.String()
	}
	// Convert addresses to strings.
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.String())
	}
	// Convert outpoints.
	ops := make([]btcjson.OutPoint, 0, len(outpoints))
	for _, op := range outpoints {
		ops = append(ops, newOutPointFromWire(op))
	}
	cmd := btcjson.NewRescanCmd(
		startBlockHashStr, addrs, ops,
		&endBlockHashStr,
	)
	return c.sendCmd(cmd)
}

// RescanEndHeight rescans the block chain starting from the provided starting block up to the provided ending block for
// transactions that pay to the passed addresses and transactions which spend the passed outpoints.
//
// The notifications of found transactions are delivered to the notification handlers associated with client and this
// call will not return until the rescan has completed.
//
// Calling this function has no effect if there are no notification handlers and will result in an error if the client
// is configured to run in HTTP POST mode.
//
// The notifications delivered as a result of this call will be via one of OnRedeemingTx (for transactions which spend
// from the one of the passed outpoints), OnRecvTx (for transactions that receive funds to one of the passed addresses),
// and OnRescanProgress (for rescan progress updates).
//
// See Rescan to also perform a rescan through current end of the longest
// chain. NOTE: This is a pod extension and requires a websocket connection.
//
// NOTE: Deprecated. Use RescanBlocks instead.
func (c *Client) RescanEndHeight(
	startBlock *chainhash.Hash,
	addresses []btcaddr.Address, outpoints []*wire.OutPoint,
	endBlock *chainhash.Hash,
) (e error) {
	return c.RescanEndBlockAsync(
		startBlock, addresses, outpoints,
		endBlock,
	).Receive()
}

// FutureLoadTxFilterResult is a future promise to deliver the result of a LoadTxFilterAsync RPC invocation (or an
// applicable error).
//
// NOTE: This is a pod extension ported from github.com/decred/dcrrpcclient and requires a websocket connection.
type FutureLoadTxFilterResult chan *response

// Receive waits for the response promised by the future and returns an error if the registration was not successful.
//
// NOTE: This is a pod extension ported from github.com/decred/dcrrpcclient and requires a websocket connection.
func (r FutureLoadTxFilterResult) Receive() (e error) {
	_, e = receiveFuture(r)
	return e
}

// LoadTxFilterAsync returns an instance of a type that can be used to get the result of the RPC at some future time by
// invoking the Receive function on the returned instance.
//
// See LoadTxFilter for the blocking version and more details.
//
// NOTE: This is a pod extension ported from github. com/decred/dcrrpcclient and requires a websocket connection.
func (c *Client) LoadTxFilterAsync(
	reload bool, addresses []btcaddr.Address,
	outPoints []wire.OutPoint,
) FutureLoadTxFilterResult {
	addrStrs := make([]string, len(addresses))
	for i, a := range addresses {
		addrStrs[i] = a.EncodeAddress()
	}
	outPointObjects := make([]btcjson.OutPoint, len(outPoints))
	for i := range outPoints {
		outPointObjects[i] = btcjson.OutPoint{
			Hash:  outPoints[i].Hash.String(),
			Index: outPoints[i].Index,
		}
	}
	cmd := btcjson.NewLoadTxFilterCmd(reload, addrStrs, outPointObjects)
	return c.sendCmd(cmd)
}

// LoadTxFilter loads reloads or adds data to a websocket client's transaction filter.
//
// The filter is consistently updated based on inspected transactions during mempool acceptance, block acceptance, and
// for all rescanned blocks.
//
// NOTE: This is a pod extension ported from github. com/decred/dcrrpcclient and requires a websocket connection.
func (c *Client) LoadTxFilter(reload bool, addresses []btcaddr.Address, outPoints []wire.OutPoint) (e error) {
	return c.LoadTxFilterAsync(reload, addresses, outPoints).Receive()
}
