package rpcclient

import (
	"bytes"
	"encoding/hex"
	js "encoding/json"
	"github.com/stalker-loki/app/slog"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
)

// SigHashType enumerates the available signature hashing types that the function accepts.
type SigHashType string

// Constants used to indicate the signature hash type for SignRawTransaction.
const (
	// SigHashAll indicates ALL of the outputs should be signed.
	SigHashAll SigHashType = "ALL"
	// SigHashNone indicates NONE of the outputs should be signed.  This can be thought of as specifying the signer does not care where the bitcoins go.
	SigHashNone SigHashType = "NONE"
	// SigHashSingle indicates that a SINGLE output should be signed.  This can be thought of specifying the signer only cares about where ONE of the outputs goes, but not any of the others.
	SigHashSingle SigHashType = "SINGLE"
	// SigHashAllAnyoneCanPay indicates that signer does not care where the other inputs to the transaction come from, so it allows other people to add inputs.  In addition, it uses the SigHashAll signing method for outputs.
	SigHashAllAnyoneCanPay SigHashType = "ALL|ANYONECANPAY"
	// SigHashNoneAnyoneCanPay indicates that signer does not care where the other inputs to the transaction come from, so it allows other people to add inputs.  In addition, it uses the SigHashNone signing method for outputs.
	SigHashNoneAnyoneCanPay SigHashType = "NONE|ANYONECANPAY"
	// SigHashSingleAnyoneCanPay indicates that signer does not care where the other inputs to the transaction come from, so it allows other people to add inputs.  In addition, it uses the SigHashSingle signing method for outputs.
	SigHashSingleAnyoneCanPay SigHashType = "SINGLE|ANYONECANPAY"
)

// String returns the SighHashType in human-readable form.
func (s SigHashType) String() string {
	return string(s)
}

// FutureGetRawTransactionResult is a future promise to deliver the result of a GetRawTransactionAsync RPC invocation (or an applicable error).
type FutureGetRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns a transaction given its hash.
func (r FutureGetRawTransactionResult) Receive() (t *util.Tx, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var txHex string
	if err = js.Unmarshal(res, &txHex); slog.Check(err) {
		return
	}
	// Decode the serialized transaction hex to raw bytes.
	var serializedTx []byte
	if serializedTx, err = hex.DecodeString(txHex); slog.Check(err) {
		return
	}
	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); slog.Check(err) {
		return
	}
	return util.NewTx(&msgTx), nil
}

// GetRawTransactionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetRawTransaction for the blocking version and more details.
func (c *Client) GetRawTransactionAsync(txHash *chainhash.Hash) FutureGetRawTransactionResult {
	hash := ""
	if txHash != nil {
		hash = txHash.String()
	}
	cmd := btcjson.NewGetRawTransactionCmd(hash, btcjson.Int(0))
	return c.sendCmd(cmd)
}

// GetRawTransaction returns a transaction given its hash. See GetRawTransactionVerbose to obtain additional information about the transaction.
func (c *Client) GetRawTransaction(txHash *chainhash.Hash) (t *util.Tx, err error) {
	return c.GetRawTransactionAsync(txHash).Receive()
}

// FutureGetRawTransactionVerboseResult is a future promise to deliver the result of a GetRawTransactionVerboseAsync RPC invocation (or an applicable error).
type FutureGetRawTransactionVerboseResult chan *response

// Receive waits for the response promised by the future and returns information about a transaction given its hash.
func (r FutureGetRawTransactionVerboseResult) Receive() (rawTxResult *btcjson.TxRawResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a gettrawtransaction result object.
	if err = js.Unmarshal(res, rawTxResult); slog.Check(err) {
		return
	}
	return
}

// GetRawTransactionVerboseAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetRawTransactionVerbose for the blocking version and more details.
func (c *Client) GetRawTransactionVerboseAsync(txHash *chainhash.Hash) FutureGetRawTransactionVerboseResult {
	hash := ""
	if txHash != nil {
		hash = txHash.String()
	}
	cmd := btcjson.NewGetRawTransactionCmd(hash, btcjson.Int(1))
	return c.sendCmd(cmd)
}

// GetRawTransactionVerbose returns information about a transaction given its hash. See GetRawTransaction to obtain only the transaction already deserialized.
func (c *Client) GetRawTransactionVerbose(txHash *chainhash.Hash) (tr *btcjson.TxRawResult, err error) {
	return c.GetRawTransactionVerboseAsync(txHash).Receive()
}

// FutureDecodeRawTransactionResult is a future promise to deliver the result of a DecodeRawTransactionAsync RPC invocation (or an applicable error).
type FutureDecodeRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns information about a transaction given its serialized bytes.
func (r FutureDecodeRawTransactionResult) Receive() (rawTxResult *btcjson.TxRawResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a decoderawtransaction result object.
	if err = js.Unmarshal(res, rawTxResult); slog.Check(err) {
	}
	return
}

// DecodeRawTransactionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See DecodeRawTransaction for the blocking version and more details.
func (c *Client) DecodeRawTransactionAsync(serializedTx []byte) FutureDecodeRawTransactionResult {
	txHex := hex.EncodeToString(serializedTx)
	cmd := btcjson.NewDecodeRawTransactionCmd(txHex)
	return c.sendCmd(cmd)
}

// DecodeRawTransaction returns information about a transaction given its serialized bytes.
func (c *Client) DecodeRawTransaction(serializedTx []byte) (tr *btcjson.TxRawResult, err error) {
	return c.DecodeRawTransactionAsync(serializedTx).Receive()
}

// FutureCreateRawTransactionResult is a future promise to deliver the result of a CreateRawTransactionAsync RPC invocation (or an applicable error).
type FutureCreateRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns a new transaction spending the provided inputs and sending to the provided addresses.
func (r FutureCreateRawTransactionResult) Receive() (msgTx *wire.MsgTx, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var txHex string
	if err = js.Unmarshal(res, &txHex); slog.Check(err) {
		return
	}
	// Decode the serialized transaction hex to raw bytes.
	var serializedTx []byte
	if serializedTx, err = hex.DecodeString(txHex); slog.Check(err) {
		return
	}
	// Deserialize the transaction and return it.
	msgTx = &wire.MsgTx{}
	if err := msgTx.Deserialize(bytes.NewReader(serializedTx)); err != nil {
		return nil, err
	}
	return
}

// CreateRawTransactionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See CreateRawTransaction for the blocking version and more details.
func (c *Client) CreateRawTransactionAsync(inputs []btcjson.TransactionInput,
	amounts map[util.Address]util.Amount, lockTime *int64) FutureCreateRawTransactionResult {
	convertedAmts := make(map[string]float64, len(amounts))
	for addr, amount := range amounts {
		convertedAmts[addr.String()] = amount.ToDUO()
	}
	cmd := btcjson.NewCreateRawTransactionCmd(inputs, convertedAmts, lockTime)
	return c.sendCmd(cmd)
}

// CreateRawTransaction returns a new transaction spending the provided inputs and sending to the provided addresses.
func (c *Client) CreateRawTransaction(inputs []btcjson.TransactionInput, amounts map[util.Address]util.Amount,
	lockTime *int64) (t *wire.MsgTx, err error) {
	return c.CreateRawTransactionAsync(inputs, amounts, lockTime).Receive()
}

// FutureSendRawTransactionResult is a future promise to deliver the result of a SendRawTransactionAsync RPC invocation (or an applicable error).
type FutureSendRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns the result of submitting the encoded transaction to the server which then relays it to the network.
func (r FutureSendRawTransactionResult) Receive() (ch *chainhash.Hash, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var txHashStr string
	if err = js.Unmarshal(res, &txHashStr); slog.Check(err) {
		return
	}
	return chainhash.NewHashFromStr(txHashStr)
}

// SendRawTransactionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See SendRawTransaction for the blocking version and more details.
func (c *Client) SendRawTransactionAsync(tx *wire.MsgTx, allowHighFees bool) FutureSendRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	cmd := btcjson.NewSendRawTransactionCmd(txHex, &allowHighFees)
	return c.sendCmd(cmd)
}

// SendRawTransaction submits the encoded transaction to the server which will then relay it to the network.
func (c *Client) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (h *chainhash.Hash, err error) {
	return c.SendRawTransactionAsync(tx, allowHighFees).Receive()
}

// FutureSignRawTransactionResult is a future promise to deliver the result of one of the SignRawTransactionAsync family
// of RPC invocations (or an applicable error).
type FutureSignRawTransactionResult chan *response

// Receive waits for the response promised by the future and returns the signed transaction as well as whether or not
// all inputs are now signed.
func (r FutureSignRawTransactionResult) Receive() (msgTx *wire.MsgTx, cpl bool, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal as a signrawtransaction result.
	var signRawTxResult btcjson.SignRawTransactionResult
	if err = js.Unmarshal(res, &signRawTxResult); slog.Check(err) {
		return
	}
	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(signRawTxResult.Hex)
	if err != nil {
		slog.Error(err)
		return
	}
	// Deserialize the transaction and return it.
	msgTx = &wire.MsgTx{}
	if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); slog.Check(err) {
		return
	}
	cpl = signRawTxResult.Complete
	return
}

// SignRawTransactionAsync returns an instance of a type that can be used to get the result of the RPC at some future
// time by invoking the Receive function on returned instance. See SignRawTransaction for the blocking version and more
// details.
func (c *Client) SignRawTransactionAsync(tx *wire.MsgTx) FutureSignRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	cmd := btcjson.NewSignRawTransactionCmd(txHex, nil, nil, nil)
	return c.sendCmd(cmd)
}

// SignRawTransaction signs inputs for the passed transaction and returns the signed transaction as well as whether or
// not all inputs are now signed. This function assumes the RPC server already knows the input transactions and private
// keys for the passed transaction which needs to be signed and uses the default signature hash type.
//
// Use one of the SignRawTransaction# variants to specify that information if needed.
func (c *Client) SignRawTransaction(tx *wire.MsgTx) (t *wire.MsgTx, is bool, err error) {
	return c.SignRawTransactionAsync(tx).Receive()
}

// SignRawTransaction2Async returns an instance of a type that can be used to get the result of the RPC at some future
// time by invoking the Receive on the returned instance. See SignRawTransaction2 for the blocking version and more
// details.
func (c *Client) SignRawTransaction2Async(tx *wire.MsgTx, inputs []btcjson.RawTxInput) FutureSignRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	cmd := btcjson.NewSignRawTransactionCmd(txHex, &inputs, nil, nil)
	return c.sendCmd(cmd)
}

// SignRawTransaction2 signs inputs for the passed transaction given the list information about the input transactions needed to perform the signing process. This only input transactions that need to be specified are ones the RPC server does not already know.  Already known input transactions will be merged with the specified transactions. See SignRawTransaction if the RPC server already knows the input transactions.
func (c *Client) SignRawTransaction2(tx *wire.MsgTx, inputs []btcjson.RawTxInput) (msgTx *wire.MsgTx, is bool, err error) {
	return c.SignRawTransaction2Async(tx, inputs).Receive()
}

// SignRawTransaction3Async returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See SignRawTransaction3 for the blocking version and more details.
func (c *Client) SignRawTransaction3Async(tx *wire.MsgTx, inputs []btcjson.RawTxInput, privKeysWIF []string,
) FutureSignRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	cmd := btcjson.NewSignRawTransactionCmd(txHex, &inputs, &privKeysWIF,
		nil)
	return c.sendCmd(cmd)
}

// SignRawTransaction3 signs inputs for the passed transaction given the list of information about extra input
// transactions and a list of private keys needed to perform the signing process.
//
// The private keys must be in wallet import format (WIF). This only input transactions that need to be specified are
// ones the RPC server does not already know.  Already known input transactions will be merged with the specified
// transactions.  This means the list of transaction inputs can be nil if the RPC server already knows them all.
//
// NOTE: Unlike the merging functionality of the input transactions, ONLY the specified private keys will be used, so
// even if the server already knows some of the private keys, they will NOT be used.
//
// See SignRawTransaction if the RPC server already knows the input transactions and private keys or SignRawTransaction2
// if it already knows the private keys.
func (c *Client) SignRawTransaction3(tx *wire.MsgTx, inputs []btcjson.RawTxInput, privKeysWIF []string,
) (msgTx *wire.MsgTx, is bool, err error) {
	return c.SignRawTransaction3Async(tx, inputs, privKeysWIF).Receive()
}

// SignRawTransaction4Async returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See SignRawTransaction4 for the blocking version and more details.
func (c *Client) SignRawTransaction4Async(tx *wire.MsgTx, inputs []btcjson.RawTxInput, privKeysWIF []string,
	hashType SigHashType) FutureSignRawTransactionResult {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return newFutureError(err)
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	cmd := btcjson.NewSignRawTransactionCmd(txHex, &inputs, &privKeysWIF,
		btcjson.String(string(hashType)))
	return c.sendCmd(cmd)
}

// SignRawTransaction4 signs inputs for the passed transaction using the the specified signature hash type given the list of information about extra input transactions and a potential list of private keys needed to perform the signing process.  The private keys, if specified, must be in wallet import format (WIF). The only input transactions that need to be specified are ones the RPC server does not already know.  This means the list of transaction inputs can be nil if the RPC server already knows them all. NOTE: Unlike the merging functionality of the input transactions, ONLY the specified private keys will be used, so even if the server already knows some of the private keys, they will NOT be used.  The list of private keys can be nil in which case any private keys the RPC server knows will be used. This function should only used if a non-default signature hash type is desired.  Otherwise, see SignRawTransaction if the RPC server already knows the input transactions and private keys, SignRawTransaction2 if it already knows the private keys, or SignRawTransaction3 if it does not know both.
func (c *Client) SignRawTransaction4(tx *wire.MsgTx, inputs []btcjson.RawTxInput, privKeysWIF []string,
	hashType SigHashType) (msgTx *wire.MsgTx, is bool, err error) {
	return c.SignRawTransaction4Async(tx, inputs, privKeysWIF, hashType).Receive()
}

// FutureSearchRawTransactionsResult is a future promise to deliver the result of the SearchRawTransactionsAsync RPC invocation (or an applicable error).
type FutureSearchRawTransactionsResult chan *response

// Receive waits for the response promised by the future and returns the found raw transactions.
func (r FutureSearchRawTransactionsResult) Receive() (msgTxns []*wire.MsgTx, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal as an array of strings.
	var searchRawTxnsResult []string
	if err = js.Unmarshal(res, &searchRawTxnsResult); slog.Check(err) {
		return
	}
	// Decode and deserialize each transaction.
	msgTxns = make([]*wire.MsgTx, 0, len(searchRawTxnsResult))
	for _, hexTx := range searchRawTxnsResult {
		// Decode the serialized transaction hex to raw bytes.
		var serializedTx []byte
		if serializedTx, err = hex.DecodeString(hexTx); slog.Check(err) {
			return
		}
		// Deserialize the transaction and add it to the result slice.
		var msgTx wire.MsgTx
		if err = msgTx.Deserialize(bytes.NewReader(serializedTx)); slog.Check(err) {
			return
		}
		msgTxns = append(msgTxns, &msgTx)
	}
	return
}

// SearchRawTransactionsAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See SearchRawTransactions for the blocking version and more details.
func (c *Client) SearchRawTransactionsAsync(address util.Address, skip, count int, reverse bool, filterAddrs []string) FutureSearchRawTransactionsResult {
	addr := address.EncodeAddress()
	verbose := btcjson.Int(0)
	cmd := btcjson.NewSearchRawTransactionsCmd(addr, verbose, &skip, &count,
		nil, &reverse, &filterAddrs)
	return c.sendCmd(cmd)
}

// SearchRawTransactions returns transactions that involve the passed address. NOTE: Chain servers do not typically provide this capability unless it has specifically been enabled. See SearchRawTransactionsVerbose to retrieve a list of data structures with information about the transactions instead of the transactions themselves.
func (c *Client) SearchRawTransactions(address util.Address, skip, count int, reverse bool, filterAddrs []string,
) (msgTx []*wire.MsgTx, err error) {
	return c.SearchRawTransactionsAsync(address, skip, count, reverse, filterAddrs).Receive()
}

// FutureSearchRawTransactionsVerboseResult is a future promise to deliver the result of the SearchRawTransactionsVerboseAsync RPC invocation (or an applicable error).
type FutureSearchRawTransactionsVerboseResult chan *response

// Receive waits for the response promised by the future and returns the found raw transactions.
func (r FutureSearchRawTransactionsVerboseResult) Receive() (result []*btcjson.SearchRawTransactionsResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal as an array of raw transaction results.
	if err = js.Unmarshal(res, &result); slog.Check(err) {
	}
	return
}

// SearchRawTransactionsVerboseAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See SearchRawTransactionsVerbose for the blocking version and more details.
func (c *Client) SearchRawTransactionsVerboseAsync(address util.Address, skip, count int, includePrevOut, reverse bool,
	filterAddrs *[]string) FutureSearchRawTransactionsVerboseResult {
	addr := address.EncodeAddress()
	verbose := btcjson.Int(1)
	var prevOut *int
	if includePrevOut {
		prevOut = btcjson.Int(1)
	}
	cmd := btcjson.NewSearchRawTransactionsCmd(addr, verbose, &skip, &count, prevOut, &reverse, filterAddrs)
	return c.sendCmd(cmd)
}

// SearchRawTransactionsVerbose returns a list of data structures that describe transactions which involve the passed address. NOTE: Chain servers do not typically provide this capability unless it has specifically been enabled. See SearchRawTransactions to retrieve a list of raw transactions instead.
func (c *Client) SearchRawTransactionsVerbose(address util.Address, skip, count int, includePrevOut, reverse bool,
	filterAddrs []string) (result []*btcjson.SearchRawTransactionsResult, err error) {
	return c.SearchRawTransactionsVerboseAsync(address, skip, count, includePrevOut, reverse, &filterAddrs).Receive()
}

// FutureDecodeScriptResult is a future promise to deliver the result of a DecodeScriptAsync RPC invocation (or an applicable error).
type FutureDecodeScriptResult chan *response

// Receive waits for the response promised by the future and returns information about a script given its serialized bytes.
func (r FutureDecodeScriptResult) Receive() (decodeScriptResult *btcjson.DecodeScriptResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a decodescript result object.
	decodeScriptResult = &btcjson.DecodeScriptResult{}
	if err = js.Unmarshal(res, decodeScriptResult); slog.Check(err) {
		return
	}
	return
}

// DecodeScriptAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See DecodeScript for the blocking version and more details.
func (c *Client) DecodeScriptAsync(serializedScript []byte) FutureDecodeScriptResult {
	scriptHex := hex.EncodeToString(serializedScript)
	cmd := btcjson.NewDecodeScriptCmd(scriptHex)
	return c.sendCmd(cmd)
}

// DecodeScript returns information about a script given its serialized bytes.
func (c *Client) DecodeScript(serializedScript []byte) (result *btcjson.DecodeScriptResult, err error) {
	return c.DecodeScriptAsync(serializedScript).Receive()
}
