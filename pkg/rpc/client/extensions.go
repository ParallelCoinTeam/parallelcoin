package rpcclient

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	js "encoding/json"
	"fmt"
	"github.com/stalker-loki/app/slog"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
)

// FutureDebugLevelResult is a future promise to deliver the result of a DebugLevelAsync RPC invocation (or an applicable error).
type FutureDebugLevelResult chan *response

// Receive waits for the response promised by the future and returns the result of setting the debug logging level to the passed level specification or the list of of the available subsystems for the special keyword 'show'.
func (r FutureDebugLevelResult) Receive() (result string, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as a string.
	if err = js.Unmarshal(res, &result); slog.Check(err) {
		return
	}
	return
}

// DebugLevelAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See DebugLevel for the blocking version and more details. NOTE: This is a pod extension.
func (c *Client) DebugLevelAsync(levelSpec string) FutureDebugLevelResult {
	cmd := btcjson.NewDebugLevelCmd(levelSpec)
	return c.sendCmd(cmd)
}

// DebugLevel dynamically sets the debug logging level to the passed level specification.
// The levelspec can be either a debug level or of the form:
// 	<subsystem>=<level>,<subsystem2>=<level2>,...
// Additionally, the special keyword 'show' can be used to get a list of the
// available subsystems.  NOTE: This is a pod extension.
func (c *Client) DebugLevel(levelSpec string) (dl string, err error) {
	return c.DebugLevelAsync(levelSpec).Receive()
}

// FutureCreateEncryptedWalletResult is a future promise to deliver the error result of a CreateEncryptedWalletAsync RPC invocation.
type FutureCreateEncryptedWalletResult chan *response

// Receive waits for and returns the error response promised by the future.
func (r FutureCreateEncryptedWalletResult) Receive() (err error) {
	if _, err = receiveFuture(r); slog.Check(err) {
	}
	return err
}

// CreateEncryptedWalletAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See CreateEncryptedWallet for the blocking version and more details. NOTE: This is a btcwallet extension.
func (c *Client) CreateEncryptedWalletAsync(passphrase string) FutureCreateEncryptedWalletResult {
	cmd := btcjson.NewCreateEncryptedWalletCmd(passphrase)
	return c.sendCmd(cmd)
}

// CreateEncryptedWallet requests the creation of an encrypted wallet.  Wallets managed by btcwallet are only written to disk with encrypted private keys, and generating wallets on the fly is impossible as it requires user input for the encryption passphrase.  This RPC specifies the passphrase and instructs the wallet creation.  This may error if a wallet is already opened, or the new wallet cannot be written to disk. NOTE: This is a btcwallet extension.
func (c *Client) CreateEncryptedWallet(passphrase string) (err error) {
	return c.CreateEncryptedWalletAsync(passphrase).Receive()
}

// FutureListAddressTransactionsResult is a future promise to deliver the result of a ListAddressTransactionsAsync RPC invocation (or an applicable error).
type FutureListAddressTransactionsResult chan *response

// Receive waits for the response promised by the future and returns information about all transactions associated with the provided addresses.
func (r FutureListAddressTransactionsResult) Receive() (transactions []btcjson.ListTransactionsResult, err error) {
	res, err := receiveFuture(r)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	// Unmarshal the result as an array of listtransactions objects.
	err = js.Unmarshal(res, &transactions)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	return transactions, nil
}

// ListAddressTransactionsAsync returns an instance of a type that can be used get the result of the RPC at some future time by invoking the Receive function on the returned instance. See ListAddressTransactions for the blocking version and more details. NOTE: This is a pod extension.
func (c *Client) ListAddressTransactionsAsync(addresses []util.Address, account string) FutureListAddressTransactionsResult {
	// Convert addresses to strings.
	addrs := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		addrs = append(addrs, addr.EncodeAddress())
	}
	cmd := btcjson.NewListAddressTransactionsCmd(addrs, &account)
	return c.sendCmd(cmd)
}

// ListAddressTransactions returns information about all transactions associated with the provided addresses. NOTE: This is a btcwallet extension.
func (c *Client) ListAddressTransactions(addresses []util.Address, account string) (ltr []btcjson.ListTransactionsResult, err error) {
	return c.ListAddressTransactionsAsync(addresses, account).Receive()
}

// FutureGetBestBlockResult is a future promise to deliver the result of a GetBestBlockAsync RPC invocation (or an applicable error).
type FutureGetBestBlockResult chan *response

// Receive waits for the response promised by the future and returns the hash and height of the block in the longest (best) chain.
func (r FutureGetBestBlockResult) Receive() (h *chainhash.Hash, bh int32, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a getbestblock result object.
	var bestBlock btcjson.GetBestBlockResult
	if err = js.Unmarshal(res, &bestBlock); slog.Check(err) {
		return
	}
	// Convert to hash from string.
	if h, err = chainhash.NewHashFromStr(bestBlock.Hash); slog.Check(err) {
		return
	}
	bh = bestBlock.Height
	return
}

// GetBestBlockAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBestBlock for the blocking version and more details. NOTE: This is a pod extension.
func (c *Client) GetBestBlockAsync() FutureGetBestBlockResult {
	cmd := btcjson.NewGetBestBlockCmd()
	return c.sendCmd(cmd)
}

// GetBestBlock returns the hash and height of the block in the longest (best) chain. NOTE: This is a pod extension.
func (c *Client) GetBestBlock() (h *chainhash.Hash, bh int32, err error) {
	return c.GetBestBlockAsync().Receive()
}

// FutureGetCurrentNetResult is a future promise to deliver the result of a GetCurrentNetAsync RPC invocation (or an applicable error).
type FutureGetCurrentNetResult chan *response

// Receive waits for the response promised by the future and returns the network the server is running on.
func (r FutureGetCurrentNetResult) Receive() (n wire.BitcoinNet, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as an int64.
	var net int64
	if err = js.Unmarshal(res, &net); slog.Check(err) {
		return
	}
	return wire.BitcoinNet(net), nil
}

// GetCurrentNetAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetCurrentNet for the blocking version and more details. NOTE: This is a pod extension.
func (c *Client) GetCurrentNetAsync() FutureGetCurrentNetResult {
	cmd := btcjson.NewGetCurrentNetCmd()
	return c.sendCmd(cmd)
}

// GetCurrentNet returns the network the server is running on. NOTE: This is a pod extension.
func (c *Client) GetCurrentNet() (bn wire.BitcoinNet, err error) {
	return c.GetCurrentNetAsync().Receive()
}

// FutureGetHeadersResult is a future promise to deliver the result of a getheaders RPC invocation (or an applicable error). NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
type FutureGetHeadersResult chan *response

// Receive waits for the response promised by the future and returns the getheaders result. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (r FutureGetHeadersResult) Receive() (headers []wire.BlockHeader, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a slice of strings.
	var result []string
	if err = js.Unmarshal(res, &result); slog.Check(err) {
		return
	}
	// Deserialize the []string into []wire.BlockHeader.
	headers = make([]wire.BlockHeader, len(result))
	var serialized []byte
	for i, headerHex := range result {
		if serialized, err = hex.DecodeString(headerHex); slog.Check(err) {
			return
		}
		if err = headers[i].Deserialize(bytes.NewReader(serialized)); slog.Check(err) {
			return
		}
	}
	return
}

// GetHeadersAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetHeaders for the blocking version and more details. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) GetHeadersAsync(blockLocators []chainhash.Hash, hashStop *chainhash.Hash) FutureGetHeadersResult {
	locators := make([]string, len(blockLocators))
	for i := range blockLocators {
		locators[i] = blockLocators[i].String()
	}
	hash := ""
	if hashStop != nil {
		hash = hashStop.String()
	}
	cmd := btcjson.NewGetHeadersCmd(locators, hash)
	return c.sendCmd(cmd)
}

// GetHeaders mimics the wire protocol getheaders and headers messages by returning all headers on the main chain after the first known block in the locators, up until a block hash matches hashStop. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) GetHeaders(blockLocators []chainhash.Hash, hashStop *chainhash.Hash) (bh []wire.BlockHeader, err error) {
	return c.GetHeadersAsync(blockLocators, hashStop).Receive()
}

// FutureExportWatchingWalletResult is a future promise to deliver the result of an ExportWatchingWalletAsync RPC invocation (or an applicable error).
type FutureExportWatchingWalletResult chan *response

// Receive waits for the response promised by the future and returns the exported wallet.
func (r FutureExportWatchingWalletResult) Receive() (walletBytes []byte, txStoreBytes []byte, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a JSON object.
	var obj map[string]interface{}
	if err = js.Unmarshal(res, &obj); slog.Check(err) {
		return
	}
	// Check for the wallet and tx string fields in the object.
	base64Wallet, ok := obj["wallet"].(string)
	if !ok {
		err = fmt.Errorf("unexpected response type for exportwatchingwallet 'wallet' field: %T\n",
			obj["wallet"])
		slog.Debug(err)
		return
	}
	base64TxStore, ok := obj["tx"].(string)
	if !ok {
		err = fmt.Errorf("unexpected response type for exportwatchingwallet 'tx' field: %T\n", obj["tx"])
		slog.Debug(err)
		return
	}
	if walletBytes, err = base64.StdEncoding.DecodeString(base64Wallet); slog.Check(err) {
		return
	}
	if txStoreBytes, err = base64.StdEncoding.DecodeString(base64TxStore); slog.Check(err) {
	}
	return
}

// ExportWatchingWalletAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See ExportWatchingWallet for the blocking version and more details. NOTE: This is a btcwallet extension.
func (c *Client) ExportWatchingWalletAsync(account string) FutureExportWatchingWalletResult {
	cmd := btcjson.NewExportWatchingWalletCmd(&account, btcjson.Bool(true))
	return c.sendCmd(cmd)
}

// ExportWatchingWallet returns the raw bytes for a watching-only version of wallet.bin and tx.bin, respectively, for the specified account that can be used by btcwallet to enable a wallet which does not have the private keys necessary to spend funds. NOTE: This is a btcwallet extension.
func (c *Client) ExportWatchingWallet(account string) (b []byte, bb []byte, err error) {
	return c.ExportWatchingWalletAsync(account).Receive()
}

// FutureSessionResult is a future promise to deliver the result of a SessionAsync RPC invocation (or an applicable error).
type FutureSessionResult chan *response

// Receive waits for the response promised by the future and returns the session result.
func (r FutureSessionResult) Receive() (session *btcjson.SessionResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a session result object.
	session = &btcjson.SessionResult{}
	if err = js.Unmarshal(res, session); slog.Check(err) {
		return
	}
	return
}

// SessionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See Session for the blocking version and more details. NOTE: This is a btcsuite extension.
func (c *Client) SessionAsync() FutureSessionResult {
	// Not supported in HTTP POST mode.
	if c.config.HTTPPostMode {
		return newFutureError(ErrWebsocketsRequired)
	}
	cmd := btcjson.NewSessionCmd()
	return c.sendCmd(cmd)
}

// Session returns details regarding a websocket client's current connection. This RPC requires the client to be running in websocket mode. NOTE: This is a btcsuite extension.
func (c *Client) Session() (sr *btcjson.SessionResult, err error) {
	return c.SessionAsync().Receive()
}

// FutureVersionResult is a future promise to deliver the result of a version RPC invocation (or an applicable error). NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
type FutureVersionResult chan *response

// Receive waits for the response promised by the future and returns the version result. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (r FutureVersionResult) Receive() (map[string]btcjson.VersionResult,
	error) {
	res, err := receiveFuture(r)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	// Unmarshal result as a version result object.
	var vr map[string]btcjson.VersionResult
	err = js.Unmarshal(res, &vr)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	return vr, nil
}

// VersionAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See Version for the blocking version and more details. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) VersionAsync() FutureVersionResult {
	cmd := btcjson.NewVersionCmd()
	return c.sendCmd(cmd)
}

// Version returns information about the server's JSON-RPC API versions. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) Version() (vr map[string]btcjson.VersionResult, err error) {
	return c.VersionAsync().Receive()
}
