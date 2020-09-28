package rpcclient

import (
	"bytes"
	"encoding/hex"
	js "encoding/json"
	"github.com/stalker-loki/app/slog"

	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// FutureGetBestBlockHashResult is a future promise to deliver the result of a GetBestBlockAsync RPC invocation (or an applicable error).
type FutureGetBestBlockHashResult chan *response

// Receive waits for the response promised by the future and returns the hash of the best block in the longest block chain.
func (r FutureGetBestBlockHashResult) Receive() (ch *chainhash.Hash, err error) {
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

// GetBestBlockHashAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBestBlockHash for the blocking version and more details.
func (c *Client) GetBestBlockHashAsync() FutureGetBestBlockHashResult {
	cmd := btcjson.NewGetBestBlockHashCmd()
	return c.sendCmd(cmd)
}

// GetBestBlockHash returns the hash of the best block in the longest block chain.
func (c *Client) GetBestBlockHash() (ch *chainhash.Hash, err error) {
	return c.GetBestBlockHashAsync().Receive()
}

// FutureGetBlockResult is a future promise to deliver the result of a GetBlockAsync RPC invocation (or an applicable error).
type FutureGetBlockResult chan *response

// Receive waits for the response promised by the future and returns the raw block requested from the server given its hash.
func (r FutureGetBlockResult) Receive() (mb *wire.MsgBlock, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var blockHex string
	if err = js.Unmarshal(res, &blockHex); slog.Check(err) {
		return
	}
	// Decode the serialized block hex to raw bytes.
	var serializedBlock []byte
	if serializedBlock, err = hex.DecodeString(blockHex); slog.Check(err) {
		return
	}
	// Deserialize the block and return it.
	mb = &wire.MsgBlock{}
	if err = mb.Deserialize(bytes.NewReader(serializedBlock)); slog.Check(err) {
		return
	}
	return
}

// GetBlockAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlock for the blocking version and more details.
func (c *Client) GetBlockAsync(blockHash *chainhash.Hash) FutureGetBlockResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetBlockCmd(hash, btcjson.Bool(false), nil)
	return c.sendCmd(cmd)
}

// GetBlock returns a raw block from the server given its hash. GetBlockVerbose to retrieve a data structure with information about the block instead.
func (c *Client) GetBlock(blockHash *chainhash.Hash) (mb *wire.MsgBlock, err error) {
	return c.GetBlockAsync(blockHash).Receive()
}

// FutureGetBlockVerboseResult is a future promise to deliver the result of a GetBlockVerboseAsync RPC invocation (or an applicable error).
type FutureGetBlockVerboseResult chan *response

// Receive waits for the response promised by the future and returns the data structure from the server with information about the requested block.
func (r FutureGetBlockVerboseResult) Receive() (blockResult *btcjson.GetBlockVerboseResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the raw result into a BlockResult.
	blockResult = &btcjson.GetBlockVerboseResult{}
	if err = js.Unmarshal(res, blockResult); slog.Check(err) {
		return
	}
	return
}

// GetBlockVerboseAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockVerbose for the blocking version and more details.
func (c *Client) GetBlockVerboseAsync(blockHash *chainhash.Hash) FutureGetBlockVerboseResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetBlockCmd(hash, btcjson.Bool(true), nil)
	return c.sendCmd(cmd)
}

// GetBlockVerbose returns a data structure from the server with information about a block given its hash. See GetBlockVerboseTx to retrieve transaction data structures as well. See GetBlock to retrieve a raw block instead.
func (c *Client) GetBlockVerbose(blockHash *chainhash.Hash) (gbvr *btcjson.GetBlockVerboseResult, err error) {
	return c.GetBlockVerboseAsync(blockHash).Receive()
}

// GetBlockVerboseTxAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockVerboseTx or the blocking version and more details.
func (c *Client) GetBlockVerboseTxAsync(blockHash *chainhash.Hash) FutureGetBlockVerboseResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetBlockCmd(hash, btcjson.Bool(true), btcjson.Bool(true))
	return c.sendCmd(cmd)
}

// GetBlockVerboseTx returns a data structure from the server with information about a block and its transactions given its hash. See GetBlockVerbose if only transaction hashes are preferred. See GetBlock to retrieve a raw block instead.
func (c *Client) GetBlockVerboseTx(blockHash *chainhash.Hash) (gbvr *btcjson.GetBlockVerboseResult, err error) {
	return c.GetBlockVerboseTxAsync(blockHash).Receive()
}

// FutureGetBlockCountResult is a future promise to deliver the result of a GetBlockCountAsync RPC invocation (or an applicable error).
type FutureGetBlockCountResult chan *response

// Receive waits for the response promised by the future and returns the number of blocks in the longest block chain.
func (r FutureGetBlockCountResult) Receive() (count int64, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as an int64.
	if err = js.Unmarshal(res, &count); slog.Check(err) {
		return
	}
	return
}

// GetBlockCountAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockCount for the blocking version and more details.
func (c *Client) GetBlockCountAsync() FutureGetBlockCountResult {
	cmd := btcjson.NewGetBlockCountCmd()
	return c.sendCmd(cmd)
}

// GetBlockCount returns the number of blocks in the longest block chain.
func (c *Client) GetBlockCount() (i int64, err error) {
	return c.GetBlockCountAsync().Receive()
}

// FutureGetDifficultyResult is a future promise to deliver the result of a GetDifficultyAsync RPC invocation (or an applicable error).
type FutureGetDifficultyResult chan *response

// Receive waits for the response promised by the future and returns the proof-of-work difficulty as a multiple of the minimum difficulty.
func (r FutureGetDifficultyResult) Receive() (difficulty float64, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as a float64.
	if err = js.Unmarshal(res, &difficulty); slog.Check(err) {
		return
	}
	return
}

// GetDifficultyAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetDifficulty for the blocking version and more details.
func (c *Client) GetDifficultyAsync(algo string) FutureGetDifficultyResult {
	cmd := btcjson.NewGetDifficultyCmd(algo)
	return c.sendCmd(cmd)
}

// GetDifficulty returns the proof-of-work difficulty as a multiple of the minimum difficulty.
func (c *Client) GetDifficulty(algo string) (f float64, err error) {
	return c.GetDifficultyAsync(algo).Receive()
}

// FutureGetBlockChainInfoResult is a promise to deliver the result of a GetBlockChainInfoAsync RPC invocation (or an applicable error).
type FutureGetBlockChainInfoResult chan *response

// Receive waits for the response promised by the future and returns chain info result provided by the server.
func (r FutureGetBlockChainInfoResult) Receive() (chainInfo *btcjson.GetBlockChainInfoResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	chainInfo = &btcjson.GetBlockChainInfoResult{}
	if err = js.Unmarshal(res, chainInfo); slog.Check(err) {
		return
	}
	return
}

// GetBlockChainInfoAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. GetBlockChainInfo for the blocking version and more details.
func (c *Client) GetBlockChainInfoAsync() FutureGetBlockChainInfoResult {
	cmd := btcjson.NewGetBlockChainInfoCmd()
	return c.sendCmd(cmd)
}

// GetBlockChainInfo returns information related to the processing state of various chain-specific details such as the current difficulty from the tip of the main chain.
func (c *Client) GetBlockChainInfo() (gbcir *btcjson.GetBlockChainInfoResult, err error) {
	return c.GetBlockChainInfoAsync().Receive()
}

// FutureGetBlockHashResult is a future promise to deliver the result of a GetBlockHashAsync RPC invocation (or an applicable error).
type FutureGetBlockHashResult chan *response

// Receive waits for the response promised by the future and returns the hash of the block in the best block chain at the given height.
func (r FutureGetBlockHashResult) Receive() (h *chainhash.Hash, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as a string-encoded sha.
	var txHashStr string
	if err = js.Unmarshal(res, &txHashStr); slog.Check(err) {
		return
	}
	return chainhash.NewHashFromStr(txHashStr)
}

// GetBlockHashAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockHash for the blocking version and more details.
func (c *Client) GetBlockHashAsync(blockHeight int64) FutureGetBlockHashResult {
	cmd := btcjson.NewGetBlockHashCmd(blockHeight)
	return c.sendCmd(cmd)
}

// GetBlockHash returns the hash of the block in the best block chain at the given height.
func (c *Client) GetBlockHash(blockHeight int64) (h *chainhash.Hash, err error) {
	return c.GetBlockHashAsync(blockHeight).Receive()
}

// FutureGetBlockHeaderResult is a future promise to deliver the result of a GetBlockHeaderAsync RPC invocation (or an applicable error).
type FutureGetBlockHeaderResult chan *response

// Receive waits for the response promised by the future and returns the blockheader requested from the server given its hash.
func (r FutureGetBlockHeaderResult) Receive() (bh *wire.BlockHeader, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var bhHex string
	if err = js.Unmarshal(res, &bhHex); slog.Check(err) {
		return
	}
	var serializedBH []byte
	if serializedBH, err = hex.DecodeString(bhHex); slog.Check(err) {
		return
	}
	// Deserialize the blockheader and return it.
	bh = &wire.BlockHeader{}
	if err = bh.Deserialize(bytes.NewReader(serializedBH)); slog.Check(err) {
		return
	}
	return
}

// GetBlockHeaderAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockHeader for the blocking version and more details.
func (c *Client) GetBlockHeaderAsync(blockHash *chainhash.Hash) FutureGetBlockHeaderResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetBlockHeaderCmd(hash, btcjson.Bool(false))
	return c.sendCmd(cmd)
}

// GetBlockHeader returns the blockheader from the server given its hash. See GetBlockHeaderVerbose to retrieve a data structure with information about the block instead.
func (c *Client) GetBlockHeader(blockHash *chainhash.Hash) (bh *wire.BlockHeader, err error) {
	return c.GetBlockHeaderAsync(blockHash).Receive()
}

// FutureGetBlockHeaderVerboseResult is a future promise to deliver the result of a GetBlockAsync RPC invocation (or an applicable error).
type FutureGetBlockHeaderVerboseResult chan *response

// Receive waits for the response promised by the future and returns the data structure of the blockheader requested from the server given its hash.
func (r FutureGetBlockHeaderVerboseResult) Receive() (bh *btcjson.GetBlockHeaderVerboseResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	bh = &btcjson.GetBlockHeaderVerboseResult{}
	if err = js.Unmarshal(res, &bh); slog.Check(err) {
		return
	}
	return bh, nil
}

// GetBlockHeaderVerboseAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetBlockHeader for the blocking version and more details.
func (c *Client) GetBlockHeaderVerboseAsync(blockHash *chainhash.Hash) FutureGetBlockHeaderVerboseResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetBlockHeaderCmd(hash, btcjson.Bool(true))
	return c.sendCmd(cmd)
}

// GetBlockHeaderVerbose returns a data structure with information about the blockheader from the server given its hash. See GetBlockHeader to retrieve a blockheader instead.
func (c *Client) GetBlockHeaderVerbose(blockHash *chainhash.Hash) (gbhvr *btcjson.GetBlockHeaderVerboseResult, err error) {
	return c.GetBlockHeaderVerboseAsync(blockHash).Receive()
}

// FutureGetMempoolEntryResult is a future promise to deliver the result of a GetMempoolEntryAsync RPC invocation (or an applicable error).
type FutureGetMempoolEntryResult chan *response

// Receive waits for the response promised by the future and returns a data structure with information about the transaction in the memory pool given its hash.
func (r FutureGetMempoolEntryResult) Receive() (mempoolEntryResult *btcjson.GetMempoolEntryResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as an array of strings.
	mempoolEntryResult = &btcjson.GetMempoolEntryResult{}
	if err = js.Unmarshal(res, &mempoolEntryResult); slog.Check(err) {
		return
	}
	return
}

// GetMempoolEntryAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetMempoolEntry for the blocking version and more details.
func (c *Client) GetMempoolEntryAsync(txHash string) FutureGetMempoolEntryResult {
	cmd := btcjson.NewGetMempoolEntryCmd(txHash)
	return c.sendCmd(cmd)
}

// GetMempoolEntry returns a data structure with information about the transaction in the memory pool given its hash.
func (c *Client) GetMempoolEntry(txHash string) (gmer *btcjson.GetMempoolEntryResult, err error) {
	return c.GetMempoolEntryAsync(txHash).Receive()
}

// FutureGetRawMempoolResult is a future promise to deliver the result of a GetRawMempoolAsync RPC invocation (or an applicable error).
type FutureGetRawMempoolResult chan *response

// Receive waits for the response promised by the future and returns the hashes of all transactions in the memory pool.
func (r FutureGetRawMempoolResult) Receive() (h []*chainhash.Hash, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as an array of strings.
	var txHashStrs []string
	if err = js.Unmarshal(res, &txHashStrs); slog.Check(err) {
		return
	}
	// Create a slice of ShaHash arrays from the string slice.
	txHashes := make([]*chainhash.Hash, 0, len(txHashStrs))
	var txHash *chainhash.Hash
	for _, hashStr := range txHashStrs {
		if txHash, err = chainhash.NewHashFromStr(hashStr); slog.Check(err) {
			return
		}
		txHashes = append(txHashes, txHash)
	}
	return
}

// GetRawMempoolAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetRawMempool for the blocking version and more details.
func (c *Client) GetRawMempoolAsync() FutureGetRawMempoolResult {
	cmd := btcjson.NewGetRawMempoolCmd(btcjson.Bool(false))
	return c.sendCmd(cmd)
}

// GetRawMempool returns the hashes of all transactions in the memory pool. See GetRawMempoolVerbose to retrieve data structures with information about the transactions instead.
func (c *Client) GetRawMempool() (h []*chainhash.Hash, err error) {
	return c.GetRawMempoolAsync().Receive()
}

// FutureGetRawMempoolVerboseResult is a future promise to deliver the result of a GetRawMempoolVerboseAsync RPC invocation (or an applicable error).
type FutureGetRawMempoolVerboseResult chan *response

// Receive waits for the response promised by the future and returns a map of transaction hashes to an associated data structure with information about the transaction for all transactions in the memory pool.
func (r FutureGetRawMempoolVerboseResult) Receive() (mempoolItems map[string]btcjson.GetRawMempoolVerboseResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as a map of strings (tx SHAs) to their detailed results.
	if err = js.Unmarshal(res, &mempoolItems); slog.Check(err) {
		return
	}
	return
}

// GetRawMempoolVerboseAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetRawMempoolVerbose for the blocking version and more details.
func (c *Client) GetRawMempoolVerboseAsync() FutureGetRawMempoolVerboseResult {
	cmd := btcjson.NewGetRawMempoolCmd(btcjson.Bool(true))
	return c.sendCmd(cmd)
}

// GetRawMempoolVerbose returns a map of transaction hashes to an associated data structure with information about the transaction for all transactions in the memory pool. See GetRawMempool to retrieve only the transaction hashes instead.
func (c *Client) GetRawMempoolVerbose() (grmvr map[string]btcjson.GetRawMempoolVerboseResult, err error) {
	return c.GetRawMempoolVerboseAsync().Receive()
}

// FutureEstimateFeeResult is a future promise to deliver the result of a EstimateFeeAsync RPC invocation (or an applicable error).
type FutureEstimateFeeResult chan *response

// Receive waits for the response promised by the future and returns the info provided by the server.
func (r FutureEstimateFeeResult) Receive() (fee float64, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a getinfo result object.
	if err = js.Unmarshal(res, &fee); slog.Check(err) {
		return
	}
	return
}

// EstimateFeeAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See EstimateFee for the blocking version and more details.
func (c *Client) EstimateFeeAsync(numBlocks int64) FutureEstimateFeeResult {
	cmd := btcjson.NewEstimateFeeCmd(numBlocks)
	return c.sendCmd(cmd)
}

// EstimateFee provides an estimated fee  in bitcoins per kilobyte.
func (c *Client) EstimateFee(numBlocks int64) (u float64, err error) {
	return c.EstimateFeeAsync(numBlocks).Receive()
}

// FutureVerifyChainResult is a future promise to deliver the result of a VerifyChainAsync, VerifyChainLevelAsyncRPC, or VerifyChainBlocksAsync invocation (or an applicable error).
type FutureVerifyChainResult chan *response

// Receive waits for the response promised by the future and returns whether or not the chain verified based on the check level and number of blocks to verify specified in the original call.
func (r FutureVerifyChainResult) Receive() (verified bool, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal the result as a boolean.
	if err = js.Unmarshal(res, &verified); slog.Check(err) {
		return
	}
	return
}

// VerifyChainAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See VerifyChain for the blocking version and more details.
func (c *Client) VerifyChainAsync() FutureVerifyChainResult {
	cmd := btcjson.NewVerifyChainCmd(nil, nil)
	return c.sendCmd(cmd)
}

// VerifyChain requests the server to verify the block chain database using the default check level and number of blocks to verify. See VerifyChainLevel and VerifyChainBlocks to override the defaults.
func (c *Client) VerifyChain() (b bool, err error) {
	return c.VerifyChainAsync().Receive()
}

// VerifyChainLevelAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See VerifyChainLevel for the blocking version and more details.
func (c *Client) VerifyChainLevelAsync(checkLevel int32) FutureVerifyChainResult {
	cmd := btcjson.NewVerifyChainCmd(&checkLevel, nil)
	return c.sendCmd(cmd)
}

// VerifyChainLevel requests the server to verify the block chain database using the passed check level and default number of blocks to verify. The check level controls how thorough the verification is with higher numbers increasing the amount of checks done as consequently how long the verification takes. See VerifyChain to use the default check level and VerifyChainBlocks to override the number of blocks to verify.
func (c *Client) VerifyChainLevel(checkLevel int32) (b bool, err error) {
	return c.VerifyChainLevelAsync(checkLevel).Receive()
}

// VerifyChainBlocksAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See VerifyChainBlocks for the blocking version and more details.
func (c *Client) VerifyChainBlocksAsync(checkLevel, numBlocks int32) FutureVerifyChainResult {
	cmd := btcjson.NewVerifyChainCmd(&checkLevel, &numBlocks)
	return c.sendCmd(cmd)
}

// VerifyChainBlocks requests the server to verify the block chain database using the passed check level and number of blocks to verify. The check level controls how thorough the verification is with higher numbers increasing the amount of checks done as consequently how long the verification takes. The number of blocks refers to the number of blocks from the end of the current longest chain. See VerifyChain and VerifyChainLevel to use defaults.
func (c *Client) VerifyChainBlocks(checkLevel, numBlocks int32) (b bool, err error) {
	return c.VerifyChainBlocksAsync(checkLevel, numBlocks).Receive()
}

// FutureGetTxOutResult is a future promise to deliver the result of a GetTxOutAsync RPC invocation (or an applicable error).
type FutureGetTxOutResult chan *response

// Receive waits for the response promised by the future and returns a transaction given its hash.
func (r FutureGetTxOutResult) Receive() (txOutInfo *btcjson.GetTxOutResult, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// take care of the special case where the output has been spent already it should return the string "null"
	if string(res) == "null" {
		return
	}
	// Unmarshal result as an gettxout result object.
	if err = js.Unmarshal(res, &txOutInfo); slog.Check(err) {
		return
	}
	return
}

// GetTxOutAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetTxOut for the blocking version and more details.
func (c *Client) GetTxOutAsync(txHash *chainhash.Hash, index uint32, mempool bool) FutureGetTxOutResult {
	hash := ""
	if txHash != nil {
		hash = txHash.String()
	}
	cmd := btcjson.NewGetTxOutCmd(hash, index, &mempool)
	return c.sendCmd(cmd)
}

// GetTxOut returns the transaction output info if it's unspent and nil, otherwise.
func (c *Client) GetTxOut(txHash *chainhash.Hash, index uint32, mempool bool) (gtor *btcjson.GetTxOutResult, err error) {
	return c.GetTxOutAsync(txHash, index, mempool).Receive()
}

// FutureRescanBlocksResult is a future promise to deliver the result of a RescanBlocksAsync RPC invocation (or an
// applicable error). NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
type FutureRescanBlocksResult chan *response

// Receive waits for the response promised by the future and returns the discovered rescanblocks data. NOTE: This is a
// btcsuite extension ported from github.com/decred/dcrrpcclient.
func (r FutureRescanBlocksResult) Receive() (rescanBlocksResult []btcjson.RescannedBlock, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	if err = js.Unmarshal(res, &rescanBlocksResult); slog.Check(err) {
		return
	}
	return rescanBlocksResult, nil
}

// RescanBlocksAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See RescanBlocks for the blocking version and more details. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) RescanBlocksAsync(blockHashes []chainhash.Hash) FutureRescanBlocksResult {
	strBlockHashes := make([]string, len(blockHashes))
	for i := range blockHashes {
		strBlockHashes[i] = blockHashes[i].String()
	}
	cmd := btcjson.NewRescanBlocksCmd(strBlockHashes)
	return c.sendCmd(cmd)
}

// RescanBlocks rescans the blocks identified by blockHashes, in order, using the client's loaded transaction filter.  The blocks do not need to be on the main chain, but they do need to be adjacent to each other. NOTE: This is a btcsuite extension ported from github.com/decred/dcrrpcclient.
func (c *Client) RescanBlocks(blockHashes []chainhash.Hash) (rb []btcjson.RescannedBlock, err error) {
	return c.RescanBlocksAsync(blockHashes).Receive()
}

// FutureInvalidateBlockResult is a future promise to deliver the result of a InvalidateBlockAsync RPC invocation (or an applicable error).
type FutureInvalidateBlockResult chan *response

// Receive waits for the response promised by the future and returns the raw block requested from the server given its hash.
func (r FutureInvalidateBlockResult) Receive() (err error) {
	_, err = receiveFuture(r)
	return err
}

// InvalidateBlockAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See InvalidateBlock for the blocking version and more details.
func (c *Client) InvalidateBlockAsync(blockHash *chainhash.Hash) FutureInvalidateBlockResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewInvalidateBlockCmd(hash)
	return c.sendCmd(cmd)
}

// InvalidateBlock invalidates a specific block.
func (c *Client) InvalidateBlock(blockHash *chainhash.Hash) (err error) {
	return c.InvalidateBlockAsync(blockHash).Receive()
}

// FutureGetCFilterResult is a future promise to deliver the result of a GetCFilterAsync RPC invocation (or an applicable error).
type FutureGetCFilterResult chan *response

// Receive waits for the response promised by the future and returns the raw filter requested from the server given its block hash.
func (r FutureGetCFilterResult) Receive() (msgCFilter *wire.MsgCFilter, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var filterHex string
	if err = js.Unmarshal(res, &filterHex); slog.Check(err) {
		return
	}
	// Decode the serialized cf hex to raw bytes.
	var serializedFilter []byte
	if res, err = hex.DecodeString(filterHex); slog.Check(err) {
		return
	}
	// Assign the filter bytes to the correct field of the wire message. We aren't going to set the block hash or extended flag, since we don't actually get that back in the RPC response.
	msgCFilter = &wire.MsgCFilter{Data: serializedFilter}
	return
}

// GetCFilterAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetCFilter for the blocking version and more details.
func (c *Client) GetCFilterAsync(blockHash *chainhash.Hash,
	filterType wire.FilterType) FutureGetCFilterResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetCFilterCmd(hash, filterType)
	return c.sendCmd(cmd)
}

// GetCFilter returns a raw filter from the server given its block hash.
func (c *Client) GetCFilter(blockHash *chainhash.Hash,
	filterType wire.FilterType) (mcf *wire.MsgCFilter, err error) {
	return c.GetCFilterAsync(blockHash, filterType).Receive()
}

// FutureGetCFilterHeaderResult is a future promise to deliver the result of a GetCFilterHeaderAsync RPC invocation (or an applicable error).
type FutureGetCFilterHeaderResult chan *response

// Receive waits for the response promised by the future and returns the raw filter header requested from the server given its block hash.
func (r FutureGetCFilterHeaderResult) Receive() (msgCFHeaders *wire.MsgCFHeaders, err error) {
	var res []byte
	if res, err = receiveFuture(r); slog.Check(err) {
		return
	}
	// Unmarshal result as a string.
	var headerHex string
	if err = js.Unmarshal(res, &headerHex); slog.Check(err) {
		return
	}
	// Assign the decoded header into a hash
	headerHash, err := chainhash.NewHashFromStr(headerHex)
	if err != nil {
		slog.Error(err)
		return nil, err
	}
	// Assign the hash to a headers message and return it.
	msgCFHeaders = &wire.MsgCFHeaders{PrevFilterHeader: *headerHash}
	return
}

// GetCFilterHeaderAsync returns an instance of a type that can be used to get the result of the RPC at some future time by invoking the Receive function on the returned instance. See GetCFilterHeader for the blocking version and more details.
func (c *Client) GetCFilterHeaderAsync(blockHash *chainhash.Hash,
	filterType wire.FilterType) FutureGetCFilterHeaderResult {
	hash := ""
	if blockHash != nil {
		hash = blockHash.String()
	}
	cmd := btcjson.NewGetCFilterHeaderCmd(hash, filterType)
	return c.sendCmd(cmd)
}

// GetCFilterHeader returns a raw filter header from the server given its block hash.
func (c *Client) GetCFilterHeader(blockHash *chainhash.Hash,
	filterType wire.FilterType) (mcfh *wire.MsgCFHeaders, err error) {
	return c.GetCFilterHeaderAsync(blockHash, filterType).Receive()
}
