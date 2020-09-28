package rpcclient

import (
	"bytes"
	"encoding/hex"
	js "encoding/json"
	"errors"
	"fmt"
	"github.com/stalker-loki/app/slog"
	"time"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
)

var ( // ErrWebsocketsRequired is an error to describe the condition where
	// the caller is trying to use a websocket-only feature,
	// such as requesting notifications or other websocket requests when the
	// client is configured to run in HTTP POST mode.
	ErrWebsocketsRequired = errors.New("a websocket connection is required " +
		"to use this feature")
)

type // notificationState is used to track the current state of successfully
	// registered notification so the state can be automatically re-established
	// on reconnect.
	notificationState struct {
		notifyBlocks       bool
		notifyNewTx        bool
		notifyNewTxVerbose bool
		notifyReceived     map[string]struct{}
		notifySpent        map[btcjson.OutPoint]struct{}
	}

func // Copy returns a deep copy of the receiver.
(s *notificationState) Copy() *notificationState {
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

type // NotificationHandlers defines callback function pointers to invoke with
	// notifications.  Since all of the functions are nil by default,
	// all notifications are effectively ignored until their handlers are set to
	// a concrete callback.
	// NOTE: Unless otherwise documented,
	// these handlers must NOT directly call any blocking calls on the client
	// instance since the input reader goroutine blocks until the callback has
	// completed.  Doing so will result in a deadlock situation.
	NotificationHandlers struct {
		// OnClientConnected is invoked when the client connects or reconnects to the RPC server.  This callback is run async with the rest of the notification handlers, and is safe for blocking client requests.
		OnClientConnected func()
		// OnBlockConnected is invoked when a block is connected to the longest (best) chain.  It will only be invoked if a preceding call to NotifyBlocks has been made to register for the notification and the function is non-nil. NOTE: Deprecated. Use OnFilteredBlockConnected instead.
		OnBlockConnected func(hash *chainhash.Hash, height int32, t time.Time)
		// OnFilteredBlockConnected is invoked when a block is connected to the longest (best) chain.  It will only be invoked if a preceding call to NotifyBlocks has been made to register for the notification and the function is non-nil.  Its parameters differ from OnBlockConnected: it receives the block's height, header, and relevant transactions.
		OnFilteredBlockConnected func(height int32, header *wire.BlockHeader,
			txs []*util.Tx)
		// OnBlockDisconnected is invoked when a block is disconnected from the longest (best) chain.  It will only be invoked if a preceding call to NotifyBlocks has been made to register for the notification and the function is non-nil. NOTE: Deprecated. Use OnFilteredBlockDisconnected instead.
		OnBlockDisconnected func(hash *chainhash.Hash, height int32, t time.Time)
		// OnFilteredBlockDisconnected is invoked when a block is disconnected from the longest (best) chain.  It will only be invoked if a preceding NotifyBlocks has been made to register for the notification and the call to function is non-nil.  Its parameters differ from OnBlockDisconnected: it receives the block's height and header.
		OnFilteredBlockDisconnected func(height int32, header *wire.BlockHeader)
		// OnRecvTx is invoked when a transaction that receives funds to a registered address is received into the memory pool and also connected to the longest (best) chain.  It will only be invoked if a preceding call to NotifyReceived, Rescan, or RescanEndHeight has been made to register for the notification and the function is non-nil. NOTE: Deprecated. Use OnRelevantTxAccepted instead.
		OnRecvTx func(transaction *util.Tx, details *btcjson.BlockDetails)
		// OnRedeemingTx is invoked when a transaction that spends a registered outpoint is received into the memory pool and also connected to the longest (best) chain.  It will only be invoked if a preceding call to NotifySpent, Rescan, or RescanEndHeight has been made to register for the notification and the function is non-nil.
		// NOTE: The NotifyReceived will automatically register notifications for the outpoints that are now "owned" as a result of receiving funds to the registered addresses.  This means it is possible for this to invoked indirectly as the result of a NotifyReceived call. NOTE: Deprecated. Use OnRelevantTxAccepted instead.
		OnRedeemingTx func(transaction *util.Tx, details *btcjson.BlockDetails)
		// OnRelevantTxAccepted is invoked when an unmined transaction passes the client's transaction filter.
		// NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
		OnRelevantTxAccepted func(transaction []byte)
		// OnRescanFinished is invoked after a rescan finishes due to a previous call to Rescan or RescanEndHeight.  Finished rescans should be signaled on this notification, rather than relying on the return result of a rescan request, due to how pod may send various rescan notifications after the rescan request has already returned. NOTE: Deprecated. Not used with RescanBlocks.
		OnRescanFinished func(hash *chainhash.Hash, height int32, blkTime time.Time)
		// OnRescanProgress is invoked periodically when a rescan is underway. It will only be invoked if a preceding call to Rescan or RescanEndHeight has been made and the function is non-nil. NOTE: Deprecated. Not used with RescanBlocks.
		OnRescanProgress func(hash *chainhash.Hash, height int32, blkTime time.Time)
		// OnTxAccepted is invoked when a transaction is accepted into the memory pool.  It will only be invoked if a preceding call to NotifyNewTransactions with the verbose flag set to false has been made to register for the notification and the function is non-nil.
		OnTxAccepted func(hash *chainhash.Hash, amount util.Amount)
		// OnTxAccepted is invoked when a transaction is accepted into the memory pool.  It will only be invoked if a preceding call to NotifyNewTransactions with the verbose flag set to true has been made to register for the notification and the function is non-nil.
		OnTxAcceptedVerbose func(txDetails *btcjson.TxRawResult)
		// OnPodConnected is invoked when a wallet connects or disconnects from pod.
		// This will only be available when client is connected to a wallet server such as btcwallet.
		OnPodConnected func(connected bool)
		// OnAccountBalance is invoked with account balance updates.
		// This will only be available when speaking to a wallet server such as btcwallet.
		OnAccountBalance func(account string, balance util.Amount, confirmed bool)
		// OnWalletLockState is invoked when a wallet is locked or unlocked.
		// This will only be available when client is connected to a wallet server such as btcwallet.
		OnWalletLockState func(locked bool)
		// OnUnknownNotification is invoked when an unrecognized notification is received.  This typically means the notification handling code for this package needs to be updated for a new notification type or the caller is using a custom notification this package does not know about.
		OnUnknownNotification func(method string, params []js.RawMessage)
	}

func // handleNotification examines the passed notification type,
// performs conversions to get the raw notification types into higher level
// types and delivers the notification to the appropriate On<X> handler
// registered with the client.
(c *Client) handleNotification(ntfn *rawNotification) {
	// Ignore the notification if the client is not interested in any notifications.
	if c.ntfnHandlers == nil {
		return
	}
	switch ntfn.Method {
	// OnBlockConnected
	case btcjson.BlockConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnBlockConnected == nil {
			return
		}
		blockHash, blockHeight, blockTime, err := parseChainNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid block connected notification:", err)
			return
		}
		c.ntfnHandlers.OnBlockConnected(blockHash, blockHeight, blockTime)
	// OnFilteredBlockConnected
	case btcjson.FilteredBlockConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnFilteredBlockConnected == nil {
			return
		}
		blockHeight, blockHeader, transactions, err :=
			parseFilteredBlockConnectedParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid filtered block connected notification:",
				err)
			return
		}
		c.ntfnHandlers.OnFilteredBlockConnected(blockHeight,
			blockHeader, transactions)
	// OnBlockDisconnected
	case btcjson.BlockDisconnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnBlockDisconnected == nil {
			return
		}
		blockHash, blockHeight, blockTime, err := parseChainNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid block connected notification:", err)
			return
		}
		c.ntfnHandlers.OnBlockDisconnected(blockHash, blockHeight, blockTime)
	// OnFilteredBlockDisconnected
	case btcjson.FilteredBlockDisconnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnFilteredBlockDisconnected == nil {
			return
		}
		blockHeight, blockHeader, err := parseFilteredBlockDisconnectedParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid filtered block disconnected"+
				" notification"+
				":", err)
			return
		}
		c.ntfnHandlers.OnFilteredBlockDisconnected(blockHeight, blockHeader)
	// OnRecvTx
	case btcjson.RecvTxNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRecvTx == nil {
			return
		}
		tx, block, err := parseChainTxNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid recvtx notification:", err)
			return
		}
		c.ntfnHandlers.OnRecvTx(tx, block)
	// OnRedeemingTx
	case btcjson.RedeemingTxNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRedeemingTx == nil {
			return
		}
		tx, block, err := parseChainTxNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid redeemingtx notification:", err)
			return
		}
		c.ntfnHandlers.OnRedeemingTx(tx, block)
	// OnRelevantTxAccepted
	case btcjson.RelevantTxAcceptedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRelevantTxAccepted == nil {
			return
		}
		transaction, err := parseRelevantTxAcceptedParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid relevanttxaccepted notification:", err)
			return
		}
		c.ntfnHandlers.OnRelevantTxAccepted(transaction)
	// OnRescanFinished
	case btcjson.RescanFinishedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRescanFinished == nil {
			return
		}
		hash, height, blkTime, err := parseRescanProgressParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid rescanfinished notification:", err)
			return
		}
		c.ntfnHandlers.OnRescanFinished(hash, height, blkTime)
	// OnRescanProgress
	case btcjson.RescanProgressNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnRescanProgress == nil {
			return
		}
		hash, height, blkTime, err := parseRescanProgressParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid rescanprogress notification:", err)
			return
		}
		c.ntfnHandlers.OnRescanProgress(hash, height, blkTime)
	// OnTxAccepted
	case btcjson.TxAcceptedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnTxAccepted == nil {
			return
		}
		hash, amt, err := parseTxAcceptedNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid tx accepted notification:", err)
			return
		}
		c.ntfnHandlers.OnTxAccepted(hash, amt)
	// OnTxAcceptedVerbose
	case btcjson.TxAcceptedVerboseNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnTxAcceptedVerbose == nil {
			return
		}
		rawTx, err := parseTxAcceptedVerboseNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid tx accepted verbose notification:", err)
			return
		}
		c.ntfnHandlers.OnTxAcceptedVerbose(rawTx)
	// OnPodConnected
	case btcjson.PodConnectedNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnPodConnected == nil {
			return
		}
		connected, err := parsePodConnectedNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid pod connected notification:", err)
			return
		}
		c.ntfnHandlers.OnPodConnected(connected)
	// OnAccountBalance
	case btcjson.AccountBalanceNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnAccountBalance == nil {
			return
		}
		account, bal, conf, err := parseAccountBalanceNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid account balance notification:", err)
			return
		}
		c.ntfnHandlers.OnAccountBalance(account, bal, conf)
	// OnWalletLockState
	case btcjson.WalletLockStateNtfnMethod:
		// Ignore the notification if the client is not interested in it.
		if c.ntfnHandlers.OnWalletLockState == nil {
			return
		}
		// The account name is not notified, so the return value is discarded.
		_, locked, err := parseWalletLockStateNtfnParams(ntfn.Params)
		if err != nil {
			slog.Error(err)
			slog.Warn("received invalid wallet lock state notification:", err)
			return
		}
		c.ntfnHandlers.OnWalletLockState(locked)
	// OnUnknownNotification
	default:
		if c.ntfnHandlers.OnUnknownNotification == nil {
			return
		}
		c.ntfnHandlers.OnUnknownNotification(ntfn.Method, ntfn.Params)
	}
}

type // wrongNumParams is an error type describing an unparseable JSON-RPC
	// notification due to an incorrect number of parameters for the expected
	// notification type.  The value is the number of parameters of the invalid
	// notification.
	wrongNumParams int

// BTCJSONError satisfies the builtin error interface.
func (e wrongNumParams) Error() string {
	return fmt.Sprintf("wrong number of parameters (%d)", e)
}

// parseChainNtfnParams parses out the block hash and height from the
// parameters of blockconnected and blockdisconnected notifications.
func parseChainNtfnParams(params []js.RawMessage) (blockHash *chainhash.Hash, blockHeight int32, blockTime time.Time,
	err error) {
	if len(params) != 3 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a string.
	var blockHashStr string
	if err = js.Unmarshal(params[0], &blockHashStr); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as an integer.
	if err = js.Unmarshal(params[1], &blockHeight); slog.Check(err) {
		return
	}
	// Unmarshal third parameter as unix time.
	var blockTimeUnix int64
	if err = js.Unmarshal(params[2], &blockTimeUnix); slog.Check(err) {
		return
	}
	// Create hash from block hash string.
	if blockHash, err = chainhash.NewHashFromStr(blockHashStr); slog.Check(err) {
		return
	}
	// Create time.Time from unix time.
	blockTime = time.Unix(blockTimeUnix, 0)
	return
}

// parseFilteredBlockConnectedParams parses out the parameters included in a
// filteredblockconnected notification.
// NOTE: This is a pod extension ported from github.
// com/decred/dcrrpcclient and requires a websocket connection.
func parseFilteredBlockConnectedParams(params []js.RawMessage) (blockHeight int32, blockHeader *wire.BlockHeader,
	transactions []*util.Tx, err error) {
	if len(params) < 3 {
		return 0, nil, nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an integer.
	if err = js.Unmarshal(params[0], &blockHeight); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as a slice of bytes.
	var blockHeaderBytes []byte
	if blockHeaderBytes, err = parseHexParam(params[1]); slog.Check(err) {
		return
	}
	// Deserialize block header from slice of bytes.
	if err = blockHeader.Deserialize(bytes.NewReader(blockHeaderBytes)); slog.Check(err) {
		return
	}
	// Unmarshal third parameter as a slice of hex-encoded strings.
	var hexTransactions []string
	if err = js.Unmarshal(params[2], &hexTransactions); slog.Check(err) {
		return
	}
	// Create slice of transactions from slice of strings by hex-decoding.
	transactions = make([]*util.Tx, len(hexTransactions))
	var transaction []byte
	for i, hexTx := range hexTransactions {
		if transaction, err = hex.DecodeString(hexTx); slog.Check(err) {
			return
		}
		if transactions[i], err = util.NewTxFromBytes(transaction); slog.Check(err) {
			return
		}
	}
	return
}

// parseFilteredBlockDisconnectedParams parses out the parameters
// included in a filteredblockdisconnected notification.
// : This is a pod extension ported from github.
// com/decred/dcrrpcclient and requires a websocket connection.
func parseFilteredBlockDisconnectedParams(params []js.RawMessage) (blockHeight int32, blockHeader *wire.BlockHeader, err error) {
	if len(params) < 2 {
		return 0, nil, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an integer.
	if err = js.Unmarshal(params[0], &blockHeight); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as a slice of bytes.
	var blockHeaderBytes []byte
	if blockHeaderBytes, err = parseHexParam(params[1]); slog.Check(err) {
		return
	}
	// Deserialize block header from slice of bytes.
	if err = blockHeader.Deserialize(bytes.NewReader(blockHeaderBytes)); slog.Check(err) {
		return
	}
	return
}

func parseHexParam(param js.RawMessage) (b []byte, err error) {
	var s string
	if err = js.Unmarshal(param, &s); slog.Check(err) {
		return
	}
	return hex.DecodeString(s)
}

// parseRelevantTxAcceptedParams parses out the parameter included in a
// relevanttxaccepted notification.
func parseRelevantTxAcceptedParams(params []js.RawMessage) (transaction []byte, err error) {
	if len(params) < 1 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	return parseHexParam(params[0])
}

func // parseChainTxNtfnParams parses out the transaction and optional details
// about the block it's mined in from the parameters of recvtx and
// redeemingtx notifications.
parseChainTxNtfnParams(params []js.RawMessage) (tx *util.Tx, block *btcjson.BlockDetails, err error) {
	if len(params) == 0 || len(params) > 2 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a string.
	var txHex string
	if err = js.Unmarshal(params[0], &txHex); slog.Check(err) {
		return
	}
	// If present, unmarshal second optional parameter as the block details JSON object.
	if len(params) > 1 {
		if err = js.Unmarshal(params[1], &block); slog.Check(err) {
			return
		}
	}
	// Hex decode and deserialize the transaction.
	var serializedTx []byte
	if serializedTx, err = hex.DecodeString(txHex); slog.Check(err) {
		return
	}
	var msgTx wire.MsgTx
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); slog.Check(err) {
		return
	}
	// TODO: Change recvtx and redeemingtx callback signatures to use nicer
	//  types for details about the block (block hash as a chainhash.Hash,
	//  block time as a time.Time, etc.).
	return util.NewTx(&msgTx), block, nil
}

// parseRescanProgressParams parses out the height of the last
// rescanned block from the parameters of rescanfinished and rescanprogress
// notifications.
func parseRescanProgressParams(params []js.RawMessage) (hash *chainhash.Hash, height int32, t time.Time, err error) {
	if len(params) != 3 {
		return nil, 0, time.Time{}, wrongNumParams(len(params))
	}
	// Unmarshal first parameter as an string.
	var hashStr string
	if err = js.Unmarshal(params[0], &hashStr); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as an integer.
	if err = js.Unmarshal(params[1], &height); slog.Check(err) {
		return
	}
	// Unmarshal third parameter as an integer.
	var blkTime int64
	if err = js.Unmarshal(params[2], &blkTime); slog.Check(err) {
		return
	}
	// Decode string encoding of block hash.
	if hash, err = chainhash.NewHashFromStr(hashStr); slog.Check(err) {
		return
	}
	t = time.Unix(blkTime, 0)
	return
}

// parseTxAcceptedNtfnParams parses out the transaction hash and total
// amount from the parameters of a txaccepted notification.
func parseTxAcceptedNtfnParams(params []js.RawMessage) (txHash *chainhash.Hash, amt util.Amount, err error) {
	if len(params) != 2 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a string.
	var txHashStr string
	if err = js.Unmarshal(params[0], &txHashStr); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as a floating point number.
	var fAmt float64
	if err = js.Unmarshal(params[1], &fAmt); slog.Check(err) {
		return
	}
	// Bounds check amount.
	if amt, err = util.NewAmount(fAmt); slog.Check(err) {
		return
	}
	// Decode string encoding of transaction sha.
	if txHash, err = chainhash.NewHashFromStr(txHashStr); slog.Check(err) {
		return
	}
	return
}

// parseTxAcceptedVerboseNtfnParams parses out details about a raw
// transaction from the parameters of a txacceptedverbose notification.
func parseTxAcceptedVerboseNtfnParams(params []js.RawMessage) (rawTx *btcjson.TxRawResult, err error) {
	if len(params) != 1 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a raw transaction result object.
	if err = js.Unmarshal(params[0], rawTx); slog.Check(err) {
		return
	}
	// TODO: change txacceptedverbose notification callbacks to use nicer
	//  types for all details about the transaction (i.e.
	//  decoding hashes from their string encoding).
	return
}

// parsePodConnectedNtfnParams parses out the connection status of pod
// and btcwallet from the parameters of a podconnected notification.
func parsePodConnectedNtfnParams(params []js.RawMessage) (connected bool, err error) {
	if len(params) != 1 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a boolean.
	if err = js.Unmarshal(params[0], &connected); slog.Check(err) {
		return
	}
	return
}

// parseAccountBalanceNtfnParams parses out the account name,
// total balance, and whether or not the balance is confirmed or unconfirmed
// from the parameters of an accountbalance notification.
func parseAccountBalanceNtfnParams(params []js.RawMessage) (
	account string, bal util.Amount, confirmed bool, err error) {
	if len(params) != 3 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a string.
	if err = js.Unmarshal(params[0], &account); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as a floating point number.
	var fBal float64
	if err = js.Unmarshal(params[1], &fBal); slog.Check(err) {
		return
	}
	// Unmarshal third parameter as a boolean.
	if err = js.Unmarshal(params[2], &confirmed); slog.Check(err) {
		return
	}
	// Bounds check amount.
	if bal, err = util.NewAmount(fBal); slog.Check(err) {
		return
	}
	return account, bal, confirmed, nil
}

func // parseWalletLockStateNtfnParams parses out the account name and locked
// state of an account from the parameters of a walletlockstate notification.
parseWalletLockStateNtfnParams(params []js.RawMessage) (account string, locked bool, err error) {
	if len(params) != 2 {
		err = wrongNumParams(len(params))
		slog.Debug(err)
		return
	}
	// Unmarshal first parameter as a string.
	if err = js.Unmarshal(params[0], &account); slog.Check(err) {
		return
	}
	// Unmarshal second parameter as a boolean.
	if err = js.Unmarshal(params[1], &locked); slog.Check(err) {
		return
	}
	return
}

// FutureNotifyBlocksResult is a future promise to deliver the result of
// a NotifyBlocksAsync RPC invocation (or an applicable error).
type FutureNotifyBlocksResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
func (r FutureNotifyBlocksResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return
}

// NotifyBlocksAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
// See NotifyBlocks for the blocking version and more details.
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

// NotifyBlocks registers the client to receive notifications when
// blocks are connected and disconnected from the main chain.
// The notifications are delivered to the notification handlers associated
// with the client.  Calling this function has no effect if there are no
// notification handlers and will result in an error if the client is
// configured to run in HTTP POST mode.
// The notifications delivered as a result of this call will be via one of or
// OnBlockDisconnected. NOTE: This is a pod extension and requires a
// websocket connection.
func (c *Client) NotifyBlocks() (err error) {
	return c.NotifyBlocksAsync().Receive()
}

// FutureNotifySpentResult is a future promise to deliver the result of
// a NotifySpentAsync RPC invocation (or an applicable error).
// NOTE: Deprecated. Use FutureLoadTxFilterResult instead.
type FutureNotifySpentResult chan *response

func // Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
(r FutureNotifySpentResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// notifySpentInternal is the same as notifySpentAsync except it accepts
// the converted outpoints as a parameter so the client can more efficiently
// recreate the previous notification state on reconnect.
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

// newOutPointFromWire constructs the json representation of a transaction
// outpoint from the wire type.
func newOutPointFromWire(op *wire.OutPoint) btcjson.OutPoint {
	return btcjson.OutPoint{
		Hash:  op.Hash.String(),
		Index: op.Index,
	}
}

// NotifySpentAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
// See NotifySpent for the blocking version and more details.
// NOTE: This is a pod extension and requires a websocket connection.
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

// NotifySpent registers the client to receive notifications when the
// passed transaction outputs are spent.
// The notifications are delivered to the notification handlers associated
// with the client.  Calling this function has no effect if there are no
// notification handlers and will result in an error if the client is
// configured to run in HTTP POST mode.
// The notifications delivered as a result of this call will be via
// OnRedeemingTx. NOTE: This is a pod extension and requires a websocket
// connection. NOTE: Deprecated. Use LoadTxFilter instead.
func (c *Client) NotifySpent(outpoints []*wire.OutPoint) (err error) {
	return c.NotifySpentAsync(outpoints).Receive()
}

type // FutureNotifyNewTransactionsResult is a future promise to deliver the
	// result of a NotifyNewTransactionsAsync RPC invocation (
	// or an applicable error).
	FutureNotifyNewTransactionsResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
func (r FutureNotifyNewTransactionsResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// NotifyNewTransactionsAsync returns an instance of a type that can be
// used to get the result of the RPC at some future time by invoking the
// Receive function on the returned instance.
// See NotifyNewTransactionsAsync for the blocking version and more details.
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

// NotifyNewTransactions registers the client to receive notifications
// every time a new transaction is accepted to the memory pool.
// The notifications are delivered to the notification handlers associated
// with the client.  Calling this function has no effect if there are no
// notification handlers and will result in an error if the client is
// configured to run in HTTP POST mode.
// The notifications delivered as a result of this call will be via one of
// OnTxAccepted (when verbose is false) or OnTxAcceptedVerbose (
// when verbose is true). NOTE: This is a pod extension and requires a
// websocket connection.
func (c *Client) NotifyNewTransactions(verbose bool) (err error) {
	return c.NotifyNewTransactionsAsync(verbose).Receive()
}

// FutureNotifyReceivedResult is a future promise to deliver the result
// of a NotifyReceivedAsync RPC invocation (or an applicable error).
// NOTE: Deprecated. Use FutureLoadTxFilterResult instead.
type FutureNotifyReceivedResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
func (r FutureNotifyReceivedResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// notifyReceivedInternal is the same as notifyReceivedAsync except it
// accepts the converted addresses as a parameter so the client can more
// efficiently recreate the previous notification state on reconnect.
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

// NotifyReceivedAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
// See NotifyReceived for the blocking version and more details.
// NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use LoadTxFilterAsync instead.
func (c *Client) NotifyReceivedAsync(addresses []util.Address) FutureNotifyReceivedResult {
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

// NotifyReceived registers the client to receive notifications every
// time a new transaction which pays to one of the passed addresses is
// accepted to memory pool or in a block connected to the block chain.
// In addition, when one of these transactions is detected,
// the client is also automatically registered for notifications when the new
// transaction outpoints the address now has available are spent (
// See NotifySpent).  The notifications are delivered to the notification
// handlers associated with the client.
// Calling this function has no effect if there are no notification handlers
// and will result in an error if the client is configured to run in HTTP
// POST mode. The notifications delivered as a result of this call will be
// via one of *OnRecvTx (for transactions that receive funds to one of the
// passed addresses) or OnRedeemingTx (
// for transactions which spend from one of the outpoints which are
// automatically registered upon receipt of funds to the address).
// NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use LoadTxFilter instead.
func (c *Client) NotifyReceived(addresses []util.Address) (err error) {
	return c.NotifyReceivedAsync(addresses).Receive()
}

// FutureRescanResult is a future promise to deliver the result of a
// RescanAsync or RescanEndHeightAsync RPC invocation (
// or an applicable error). NOTE: Deprecated.
// Use FutureRescanBlocksResult instead.
type FutureRescanResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the rescan was not successful.
func (r FutureRescanResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// RescanAsync returns an instance of a type that can be used to get the
// result of the RPC at some future time by invoking the Receive function on
// the returned instance. See Rescan for the blocking version and more
// details. NOTE: Rescan requests are not issued on client reconnect and must
// be performed manually (ideally with a new start height based on the last
// rescan progress notification).
// See the OnClientConnected notification callback for a good call site to
// reissue rescan requests on connect and reconnect.
// NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use RescanBlocksAsync instead.
func (c *Client) RescanAsync(startBlock *chainhash.Hash,
	addresses []util.Address,
	outpoints []*wire.OutPoint) FutureRescanResult {
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

// Rescan rescans the block chain starting from the provided starting
// block to the end of the longest chain for transactions that pay to the
// passed addresses and transactions which spend the passed outpoints.
// The notifications of found transactions are delivered to the notification
// handlers associated with client and this call will not return until the
// rescan has completed.  Calling this function has no effect if there are no
// notification handlers and will result in an error if the client is
// configured to run in HTTP POST mode.
// The notifications delivered as a result of this call will be via one of
// OnRedeemingTx (for transactions which spend from the one of the passed
// outpoints), OnRecvTx (for transactions that receive funds to one of the
// passed addresses), and OnRescanProgress (for rescan progress updates).
// See RescanEndBlock to also specify an ending block to finish the rescan
// without continuing through the best block on the main chain.
// NOTE: Rescan requests are not issued on client reconnect and must be
// performed manually (ideally with a new start height based on the last
// rescan progress notification).
// See the OnClientConnected notification callback for a good call site to
// reissue rescan requests on connect and reconnect.
// NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use RescanBlocks instead.
func (c *Client) Rescan(startBlock *chainhash.Hash,
	addresses []util.Address,
	outpoints []*wire.OutPoint) (err error) {
	return c.RescanAsync(startBlock, addresses, outpoints).Receive()
}

// RescanEndBlockAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
// See RescanEndBlock for the blocking version and more details.
// NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use RescanBlocksAsync instead.
func (c *Client) RescanEndBlockAsync(startBlock *chainhash.Hash,
	addresses []util.Address, outpoints []*wire.OutPoint,
	endBlock *chainhash.Hash) FutureRescanResult {
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
	cmd := btcjson.NewRescanCmd(startBlockHashStr, addrs, ops,
		&endBlockHashStr)
	return c.sendCmd(cmd)
}

// RescanEndHeight rescans the block chain starting from the provided
// starting block up to the provided ending block for transactions that pay
// to the passed addresses and transactions which spend the passed outpoints.
// The notifications of found transactions are delivered to the notification
// handlers associated with client and this call will not return until the
// rescan has completed.  Calling this function has no effect if there are no
// notification handlers and will result in an error if the client is
// configured to run in HTTP POST mode.
// The notifications delivered as a result of this call will be via one of
// OnRedeemingTx (for transactions which spend from the one of the passed
// outpoints), OnRecvTx (for transactions that receive funds to one of the
// passed addresses), and OnRescanProgress (for rescan progress updates).
// See Rescan to also perform a rescan through current end of the longest
// chain. NOTE: This is a pod extension and requires a websocket connection.
// NOTE: Deprecated. Use RescanBlocks instead.
func (c *Client) RescanEndHeight(startBlock *chainhash.Hash,
	addresses []util.Address, outpoints []*wire.OutPoint,
	endBlock *chainhash.Hash) (err error) {
	return c.RescanEndBlockAsync(startBlock, addresses, outpoints,
		endBlock).Receive()
}

// FutureLoadTxFilterResult is a future promise to deliver the result of
// a LoadTxFilterAsync RPC invocation (or an applicable error).
// NOTE: This is a pod extension ported from github.com/decred/dcrrpcclient
// and requires a websocket connection.
type FutureLoadTxFilterResult chan *response

// Receive waits for the response promised by the future and returns an
// error if the registration was not successful.
// NOTE: This is a pod extension ported from github.com/decred/dcrrpcclient
// and requires a websocket connection.
func (r FutureLoadTxFilterResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// LoadTxFilterAsync returns an instance of a type that can be used to
// get the result of the RPC at some future time by invoking the Receive
// function on the returned instance.
// See LoadTxFilter for the blocking version and more details.
// NOTE: This is a pod extension ported from github.
// com/decred/dcrrpcclient and requires a websocket connection.
func (c *Client) LoadTxFilterAsync(reload bool, addresses []util.Address, outPoints []wire.OutPoint,
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

// LoadTxFilter loads reloads or adds data to a websocket client's
// transaction filter.
// The filter is consistently updated based on inspected transactions during
// mempool acceptance, block acceptance, and for all rescanned blocks.
// NOTE: This is a pod extension ported from github.
// com/decred/dcrrpcclient and requires a websocket connection.
func (c *Client) LoadTxFilter(reload bool, addresses []util.Address, outPoints []wire.OutPoint) (err error) {
	return c.LoadTxFilterAsync(reload, addresses, outPoints).Receive()
}
