package chainrpc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/p9c/matrjoska/pkg/amt"
	block2 "github.com/p9c/matrjoska/pkg/block"
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/fork"
	"github.com/p9c/log"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/matrjoska/pkg/blockchain"
	"github.com/p9c/matrjoska/pkg/btcjson"
	"github.com/p9c/matrjoska/pkg/chainhash"
	"github.com/p9c/matrjoska/pkg/database"
	"github.com/p9c/matrjoska/pkg/ecc"
	"github.com/p9c/interrupt"
	"github.com/p9c/matrjoska/pkg/mempool"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/util"
	"github.com/p9c/matrjoska/pkg/wire"
)

// HandleAddNode handles addnode commands.
func HandleAddNode(s *Server, cmd interface{}, closeChan qu.C) (ifc interface{}, e error) {
	var msg string
	c, ok := cmd.(*btcjson.AddNodeCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("addnode")
		D.Ln(h, e)
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	addr := NormalizeAddress(c.Addr, s.Cfg.ChainParams.DefaultPort)
	switch c.SubCmd {
	case "add":
		e = s.Cfg.ConnMgr.Connect(addr, true)
	case "remove":
		e = s.Cfg.ConnMgr.RemoveByAddr(addr)
	case "onetry":
		e = s.Cfg.ConnMgr.Connect(addr, false)
	default:
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "invalid subcommand for addnode",
		}
	}
	if e != nil {
		E.Ln(e)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: e.Error(),
		}
	}
	// no data returned unless an error.
	return nil, nil
}

// HandleAskWallet is the handler for commands that are recognized as valid, but are unable to answer correctly since it
// involves wallet state.
func HandleAskWallet(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	return nil, ErrRPCNoWallet
}

// HandleCreateRawTransaction handles createrawtransaction commands.
func HandleCreateRawTransaction(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.CreateRawTransactionCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("createrawtransaction")
		D.Ln(h, e)
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Validate the locktime, if given.
	if c.LockTime != nil &&
		(*c.LockTime < 0 || *c.LockTime > int64(wire.MaxTxInSequenceNum)) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Locktime out of range",
		}
	}
	// Add all transaction inputs to a new transaction after performing some validity checks.
	mtx := wire.NewMsgTx(wire.TxVersion)
	for _, input := range c.Inputs {
		var txHash *chainhash.Hash
		txHash, e = chainhash.NewHashFromStr(input.Txid)
		if e != nil {
			E.Ln(e)
			return nil, DecodeHexError(input.Txid)
		}
		prevOut := wire.NewOutPoint(txHash, input.Vout)
		txIn := wire.NewTxIn(prevOut, []byte{}, nil)
		if c.LockTime != nil && *c.LockTime != 0 {
			txIn.Sequence = wire.MaxTxInSequenceNum - 1
		}
		mtx.AddTxIn(txIn)
	}
	// Add all transaction outputs to the transaction after performing some validity checks.
	params := s.Cfg.ChainParams
	for encodedAddr, amount := range c.Amounts {
		// Ensure amount is in the valid range for monetary amounts.
		if amount <= 0 || amount > amt.MaxSatoshi.ToDUO() {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCType,
				Message: "Invalid amount",
			}
		}
		// Decode the provided address.
		var addr btcaddr.Address
		addr, e = btcaddr.Decode(encodedAddr, params)
		if e != nil {
			E.Ln(e)
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address or key: " + e.Error(),
			}
		}
		// Ensure the address is one of the supported types and that the network encoded with the address matches the
		// network the Server is currently on.
		switch addr.(type) {
		case *btcaddr.PubKeyHash:
		case *btcaddr.ScriptHash:
		default:
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address or key",
			}
		}
		if !addr.IsForNet(params) {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidAddressOrKey,
				Message: "Invalid address: " + encodedAddr +
					" is for the wrong network",
			}
		}
		// Create a new script which pays to the provided address.
		var pkScript []byte
		pkScript, e = txscript.PayToAddrScript(addr)
		if e != nil {
			E.Ln(e)
			context := "Failed to generate pay-to-address script"
			return nil, InternalRPCError(e.Error(), context)
		}
		// Convert the amount to satoshi.
		var satoshi amt.Amount
		satoshi, e = amt.NewAmount(amount)
		if e != nil {
			E.Ln(e)
			context := "Failed to convert amount"
			return nil, InternalRPCError(e.Error(), context)
		}
		txOut := wire.NewTxOut(int64(satoshi), pkScript)
		mtx.AddTxOut(txOut)
	}
	// Set the Locktime, if given.
	if c.LockTime != nil {
		mtx.LockTime = uint32(*c.LockTime)
	}
	// Return the serialized and hex-encoded transaction. Note that this is intentionally not directly returning because
	// the first return value is a string and it would result in returning an empty string to the client instead of
	// nothing (nil) in the case of an error.
	mtxHex, e := MessageToHex(mtx)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	return mtxHex, nil
}

// HandleDecodeRawTransaction handles decoderawtransaction commands.
func HandleDecodeRawTransaction(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.DecodeRawTransactionCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("decoderawtransaction")
		D.Ln(h, e)
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Deserialize the transaction.
	hexStr := c.HexTx
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, e := hex.DecodeString(hexStr)
	if e != nil {
		E.Ln(e)
		return nil, DecodeHexError(hexStr)
	}
	var mtx wire.MsgTx
	e = mtx.Deserialize(bytes.NewReader(serializedTx))
	if e != nil {
		E.Ln(e)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX decode failed: " + e.Error(),
		}
	}
	// Create and return the result.
	txReply := btcjson.TxRawDecodeResult{
		Txid:     mtx.TxHash().String(),
		Version:  mtx.Version,
		Locktime: mtx.LockTime,
		Vin:      CreateVinList(&mtx),
		Vout:     CreateVoutList(&mtx, s.Cfg.ChainParams, nil),
	}
	return txReply, nil
}

// HandleDecodeScript handles decodescript commands.
func HandleDecodeScript(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.DecodeScriptCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("decodescript")
		D.Ln(h, e)
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Convert the hex script to bytes.
	hexStr := c.HexScript
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	script, e := hex.DecodeString(hexStr)
	if e != nil {
		E.Ln(e)
		return nil, DecodeHexError(hexStr)
	}
	// The disassembled string will contain [error] inline if the script doesn't fully parse, so ignore the error here.
	disbuf, _ := txscript.DisasmString(script)
	// Get information about the script. Ignore the error here since an error means the script couldn't parse and there
	// is no additinal information about it anyways.
	scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(
		script,
		s.Cfg.ChainParams,
	)
	addresses := make([]string, len(addrs))
	for i, addr := range addrs {
		addresses[i] = addr.EncodeAddress()
	}
	// Convert the script itself to a pay-to-script-hash address.
	p2sh, e := btcaddr.NewScriptHash(script, s.Cfg.ChainParams)
	if e != nil {
		E.Ln(e)
		context := "Failed to convert script to pay-to-script-hash"
		return nil, InternalRPCError(e.Error(), context)
	}
	// Generate and return the reply.
	reply := btcjson.DecodeScriptResult{
		Asm:       disbuf,
		ReqSigs:   int32(reqSigs),
		Type:      scriptClass.String(),
		Addresses: addresses,
	}
	if scriptClass != txscript.ScriptHashTy {
		reply.P2sh = p2sh.EncodeAddress()
	}
	return reply, nil
}

// HandleEstimateFee handles estimatefee commands.
func HandleEstimateFee(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.EstimateFeeCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("estimatefee")
		D.Ln(h, e)
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	if s.Cfg.FeeEstimator == nil {
		return nil, errors.New("fee estimation disabled")
	}
	if c.NumBlocks <= 0 {
		return -1.0, errors.New("parameter NumBlocks must be positive")
	}
	feeRate, e := s.Cfg.FeeEstimator.EstimateFee(uint32(c.NumBlocks))
	if e != nil {
		E.Ln(e)
		return -1.0, e
	}
	// Convert to satoshis per kb.
	return float64(feeRate), nil
}

// HandleGenerate handles generate commands.
func HandleGenerate(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	// Respond with an error if there are no addresses to pay the created blocks to.
	if len(s.StateCfg.ActiveMiningAddrs) == 0 {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInternal.Code,
			Message: "No payment addresses specified via --miningaddr",
		}
	}
	// Respond with an error if there's virtually 0 chance of mining a block with the CPU.
	if !s.Cfg.ChainParams.GenerateSupported {
		return nil, &btcjson.RPCError{
			Code: btcjson.ErrRPCDifficulty,
			Message: fmt.Sprintf(
				"No support for `generate` on the current"+
					" network, %s, as it's unlikely to be possible to mine a block"+
					" with the CPU.", s.Cfg.ChainParams.Net,
			),
		}
	}
	D.Ln("cpu miner stuff is missing here")
	// Set the algorithm according to the port we were called on
	// s.Cfg.CPUMiner.SetAlgo(s.Cfg.Algo)
	// c := cmd.(*btcjson.GenerateCmd)
	// // Respond with an error if the client is requesting 0 blocks to be
	// // generated.
	// if c.NumBlocks == 0 {
	// 	return nil, &btcjson.RPCError{
	// 		Code:    btcjson.ErrRPCInternal.Code,
	// 		Message: "Please request a nonzero number of blocks to generate.",
	// 	}
	// }
	// // Create a reply
	// reply := make([]string, c.NumBlocks)
	// blockHashes, e := s.Cfg.CPUMiner.GenerateNBlocks(0, c.NumBlocks,
	// 	s.Cfg.Algo)
	// if e != nil  {
	// 	L.Script	// 	return nil, &btcjson.RPCError{
	// 		Code:    btcjson.ErrRPCInternal.Code,
	// 		Message: err.ScriptError(),
	// 	}
	// }
	// // Mine the correct number of blocks, assigning the hex representation of
	// // the hash of each one to its place in the reply.
	// for i, hash := range blockHashes {
	// 	reply[i] = hash.String()
	// }
	// return reply, nil
	return nil, nil
}

// HandleGetAddedNodeInfo handles getaddednodeinfo commands.
func HandleGetAddedNodeInfo(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetAddedNodeInfoCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getaddednodeinfo")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Retrieve a list of persistent (added) peers from the Server and filter the list of peers per the specified
	// address (if any).
	peers := s.Cfg.ConnMgr.PersistentPeers()
	if c.Node != nil {
		node := *c.Node
		found := false
		for i, peer := range peers {
			if peer.ToPeer().Addr() == node {
				peers = peers[i : i+1]
				found = true
			}
		}
		if !found {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCClientNodeNotAdded,
				Message: "Node has not been added",
			}
		}
	}
	// Without the dns flag, the result is just a slice of the addresses as strings.
	if !c.DNS {
		results := make([]string, 0, len(peers))
		for _, peer := range peers {
			results = append(results, peer.ToPeer().Addr())
		}
		return results, nil
	}
	// With the dns flag, the result is an array of JSON objects which include the result of DNS lookups for each peer.
	results := make([]*btcjson.GetAddedNodeInfoResult, 0, len(peers))
	for _, rpcPeer := range peers {
		// Set the "address" of the peer which could be an ip address or a domain name.
		peer := rpcPeer.ToPeer()
		var result btcjson.GetAddedNodeInfoResult
		result.AddedNode = peer.Addr()
		result.Connected = btcjson.Bool(peer.Connected())
		// Split the address into host and port portions so we can do a DNS lookup against the host. When no port is
		// specified in the address, just use the address as the host.
		var host string
		host, _, e = net.SplitHostPort(peer.Addr())
		if e != nil {
			host = peer.Addr()
		}
		var ipList []string
		switch {
		case net.ParseIP(host) != nil, strings.HasSuffix(host, ".onion"):
			ipList = make([]string, 1)
			ipList[0] = host
		default:
			// Do a DNS lookup for the address. If the lookup fails, just use the host.
			ips, e := Lookup(s.StateCfg)(host)
			if e != nil {
				ipList = make([]string, 1)
				ipList[0] = host
				break
			}
			ipList = make([]string, 0, len(ips))
			for _, ip := range ips {
				ipList = append(ipList, ip.String())
			}
		}
		// Add the addresses and connection info to the result.
		addrs := make([]btcjson.GetAddedNodeInfoResultAddr, 0, len(ipList))
		for _, ip := range ipList {
			var addr btcjson.GetAddedNodeInfoResultAddr
			addr.Address = ip
			addr.Connected = "false"
			if ip == host && peer.Connected() {
				addr.Connected = log.DirectionString(peer.Inbound())
			}
			addrs = append(addrs, addr)
		}
		result.Addresses = &addrs
		results = append(results, &result)
	}
	return results, nil
}

// HandleGetBestBlock implements the getbestblock command.
func HandleGetBestBlock(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	// All other "get block" commands give either the height, the hash, or both but require the block SHA. This gets
	// both for the best block.
	best := s.Cfg.Chain.BestSnapshot()
	result := &btcjson.GetBestBlockResult{
		Hash:   best.Hash.String(),
		Height: best.Height,
	}
	return result, nil
}

// HandleGetBestBlockHash implements the getbestblockhash command.
func HandleGetBestBlockHash(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	best := s.Cfg.Chain.BestSnapshot()
	return best.Hash.String(), nil
}

// HandleGetBlock implements the getblock command.
func HandleGetBlock(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetBlockCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getblock")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Load the raw block bytes from the database.
	hash, e := chainhash.NewHashFromStr(c.Hash)
	if e != nil {
		return nil, DecodeHexError(c.Hash)
	}
	var blkBytes []byte
	e = s.Cfg.DB.View(
		func(dbTx database.Tx) (e error) {
			blkBytes, e = dbTx.FetchBlock(hash)
			return e
		},
	)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCBlockNotFound,
			Message: "Block not found",
		}
	}
	// When the verbose flag isn't set, simply return the serialized block as a hex-encoded string.
	if c.Verbose != nil && !*c.Verbose {
		return hex.EncodeToString(blkBytes), nil
	}
	// The verbose flag is set, so generate the JSON object and return it. Deserialize the block.
	blk, e := block2.NewFromBytes(blkBytes)
	if e != nil {
		context := "Failed to deserialize block"
		return nil, InternalRPCError(e.Error(), context)
	}
	// Get the block height from chain.
	blockHeight, e := s.Cfg.Chain.BlockHeightByHash(hash)
	if e != nil {
		context := blockheightfail
		return nil, InternalRPCError(e.Error(), context)
	}
	blk.SetHeight(blockHeight)
	best := s.Cfg.Chain.BestSnapshot()
	// Get next block hash unless there are none.
	var nextHashString string
	if blockHeight < best.Height {
		nextHash, e := s.Cfg.Chain.BlockHashByHeight(blockHeight + 1)
		if e != nil {
			context := "No next block"
			return nil, InternalRPCError(e.Error(), context)
		}
		nextHashString = nextHash.String()
	}
	params := s.Cfg.ChainParams
	blockHeader := &blk.WireBlock().Header
	algoname := fork.GetAlgoName(blockHeader.Version, blockHeight)
	a := fork.GetAlgoVer(algoname, blockHeight)
	algoid := fork.GetAlgoID(algoname, blockHeight)
	blockReply := btcjson.GetBlockVerboseResult{
		Hash:          c.Hash,
		Version:       blockHeader.Version,
		VersionHex:    fmt.Sprintf("%08x", blockHeader.Version),
		PowAlgoID:     algoid,
		PowAlgo:       algoname,
		PowHash:       blk.WireBlock().BlockHashWithAlgos(blockHeight).String(),
		MerkleRoot:    blockHeader.MerkleRoot.String(),
		PreviousHash:  blockHeader.PrevBlock.String(),
		Nonce:         blockHeader.Nonce,
		Time:          blockHeader.Timestamp.Unix(),
		Confirmations: int64(1 + best.Height - blockHeight),
		Height:        int64(blockHeight),
		TxNum:         len(blk.Transactions()),
		Size:          int32(len(blkBytes)),
		StrippedSize:  int32(blk.WireBlock().SerializeSizeStripped()),
		Weight:        int32(blockchain.GetBlockWeight(blk)),
		Bits:          strconv.FormatInt(int64(blockHeader.Bits), 16),
		Difficulty:    GetDifficultyRatio(blockHeader.Bits, params, a),
		NextHash:      nextHashString,
	}
	if c.VerboseTx == nil || !*c.VerboseTx {
		transactions := blk.Transactions()
		txNames := make([]string, len(transactions))
		for i, tx := range transactions {
			txNames[i] = tx.Hash().String()
		}
		blockReply.Tx = txNames
	} else {
		txns := blk.Transactions()
		rawTxns := make([]btcjson.TxRawResult, len(txns))
		for i, tx := range txns {
			rawTxn, e := CreateTxRawResult(
				params, tx.MsgTx(),
				tx.Hash().String(), blockHeader, hash.String(),
				blockHeight, best.Height,
			)
			if e != nil {
				return nil, e
			}
			rawTxns[i] = *rawTxn
		}
		blockReply.RawTx = rawTxns
	}
	return blockReply, nil
}

// HandleGetBlockChainInfo implements the getblockchaininfo command.
func HandleGetBlockChainInfo(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	// Obtain a snapshot of the current best known blockchain state. We'll populate the response to this call primarily
	// from this snapshot.
	params := s.Cfg.ChainParams
	chain := s.Cfg.Chain
	chainSnapshot := chain.BestSnapshot()
	chainInfo := &btcjson.GetBlockChainInfoResult{
		Chain:         params.Name,
		Blocks:        chainSnapshot.Height,
		Headers:       chainSnapshot.Height,
		BestBlockHash: chainSnapshot.Hash.String(),
		Difficulty:    GetDifficultyRatio(chainSnapshot.Bits, params, 2),
		MedianTime:    chainSnapshot.MedianTime.Unix(),
		Pruned:        false,
		// Bip9SoftForks: make(map[string]*btcjson.Bip9SoftForkDescription),
	}
	// Next, populate the response with information describing the current status of soft-forks deployed via the
	// super-majority block signalling mechanism.
	// height := chainSnapshot.Height
	// chainInfo.SoftForks = []*btcjson.SoftForkDescription{
	// 	{
	// 		ID:      "bip34",
	// 		Version: 2,
	// 		Reject: struct {
	// 			Status bool `json:"status"`
	// 		}{
	// 			Status: height >= params.BIP0034Height,
	// 		},
	// 	},
	// 	{
	// 		ID:      "bip66",
	// 		Version: 3,
	// 		Reject: struct {
	// 			Status bool `json:"status"`
	// 		}{
	// 			Status: height >= params.BIP0066Height,
	// 		},
	// 	},
	// 	{
	// 		ID:      "bip65",
	// 		Version: 4,
	// 		Reject: struct {
	// 			Status bool `json:"status"`
	// 		}{
	// 			Status: height >= params.BIP0065Height,
	// 		},
	// 	},
	// }
	// // Finally, query the BIP0009 version bits state for all currently defined BIP0009 soft-fork deployments.
	// for deployment, deploymentDetails := range params.Deployments {
	// 	// Map the integer deployment ID into a human readable fork-name.
	// 	var forkName string
	// 	switch deployment {
	// 	case chaincfg.DeploymentTestDummy:
	// 		forkName = "dummy"
	// 	case chaincfg.DeploymentCSV:
	// 		forkName = "csv"
	// 	case chaincfg.DeploymentSegwit:
	// 		forkName = "segwit"
	// 	default:
	// 		return nil, &btcjson.RPCError{
	// 			Code: btcjson.ErrRPCInternal.Code,
	// 			Message: fmt.Sprintf(
	// 				"Unknown deployment %v "+
	// 					"detected", deployment,
	// 			),
	// 		}
	// 	}
	// 	// Query the chain for the current status of the deployment as identified by its deployment ID.
	// 	deploymentStatus, e := chain.ThresholdState(uint32(deployment))
	// 	if e != nil  {
	// 		General	// 		context := "Failed to obtain deployment status"
	// 		return nil, InternalRPCError(err.GeneralError(), context)
	// 	}
	// 	// Attempt to convert the current deployment status into a human readable string. If the status is unrecognized,
	// 	// then a non-nil error is returned.
	// 	statusString, e := SoftForkStatus(deploymentStatus)
	// 	if e != nil  {
	// 		General	// 		return nil, &btcjson.RPCError{
	// 			Code: btcjson.ErrRPCInternal.Code,
	// 			Message: fmt.Sprintf(
	// 				"unknown deployment status: %v",
	// 				deploymentStatus,
	// 			),
	// 		}
	// 	}
	// 	// Finally, populate the soft-fork description with all the information gathered above.
	// 	chainInfo.Bip9SoftForks[forkName] = &btcjson.Bip9SoftForkDescription{
	// 		Status:    strings.ToLower(statusString),
	// 		Bit:       deploymentDetails.BitNumber,
	// 		StartTime: int64(deploymentDetails.StartTime),
	// 		Timeout:   int64(deploymentDetails.ExpireTime),
	// 	}
	// }
	return chainInfo, nil
}

// HandleGetBlockCount implements the getblockcount command.
func HandleGetBlockCount(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	best := s.Cfg.Chain.BestSnapshot()
	return int64(best.Height), nil
}

// HandleGetBlockHash implements the getblockhash command.
func HandleGetBlockHash(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetBlockHashCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getblockhash")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	hash, e := s.Cfg.Chain.BlockHashByHeight(int32(c.Index))
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCOutOfRange,
			Message: "Block number out of range",
		}
	}
	return hash.String(), nil
}

// HandleGetBlockHeader implements the getblockheader command.
func HandleGetBlockHeader(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetBlockHeaderCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getblockheader")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Fetch the header from chain.
	hash, e := chainhash.NewHashFromStr(c.Hash)
	if e != nil {
		return nil, DecodeHexError(c.Hash)
	}
	blockHeader, e := s.Cfg.Chain.HeaderByHash(hash)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCBlockNotFound,
			Message: "Block not found",
		}
	}
	// When the verbose flag isn't set, simply return the serialized block header as a hex-encoded string.
	if c.Verbose != nil && !*c.Verbose {
		var headerBuf bytes.Buffer
		e = blockHeader.Serialize(&headerBuf)
		if e != nil {
			context := "Failed to serialize block header"
			return nil, InternalRPCError(e.Error(), context)
		}
		return hex.EncodeToString(headerBuf.Bytes()), nil
	}
	// The verbose flag is set, so generate the JSON object and return it. Get the block height from chain.
	blockHeight, e := s.Cfg.Chain.BlockHeightByHash(hash)
	if e != nil {
		context := blockheightfail
		return nil, InternalRPCError(e.Error(), context)
	}
	best := s.Cfg.Chain.BestSnapshot()
	// Get next block hash unless there are none.
	var nextHashString string
	if blockHeight < best.Height {
		nextHash, e := s.Cfg.Chain.BlockHashByHeight(blockHeight + 1)
		if e != nil {
			context := "No next block"
			return nil, InternalRPCError(e.Error(), context)
		}
		nextHashString = nextHash.String()
	}
	var a int32 = 2
	if blockHeader.Version == 514 {
		a = 514
	}
	params := s.Cfg.ChainParams
	blockHeaderReply := btcjson.GetBlockHeaderVerboseResult{
		Hash:          c.Hash,
		Confirmations: int64(1 + best.Height - blockHeight),
		Height:        blockHeight,
		Version:       blockHeader.Version,
		VersionHex:    fmt.Sprintf("%08x", blockHeader.Version),
		MerkleRoot:    blockHeader.MerkleRoot.String(),
		NextHash:      nextHashString,
		PreviousHash:  blockHeader.PrevBlock.String(),
		Nonce:         uint64(blockHeader.Nonce),
		Time:          blockHeader.Timestamp.Unix(),
		Bits:          strconv.FormatInt(int64(blockHeader.Bits), 16),
		Difficulty:    GetDifficultyRatio(blockHeader.Bits, params, a),
	}
	return blockHeaderReply, nil
}

// HandleGetBlockTemplate implements the getblocktemplate command. See https:// en.bitcoin.it/wiki/BIP_0022 and
// https://en.bitcoin.it/wiki/BIP_0023 for more details.
func HandleGetBlockTemplate(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetBlockTemplateCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getblocktemplate")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	request := c.Request
	// Set the default mode and override it if supplied.
	mode := "template"
	if request != nil && request.Mode != "" {
		mode = request.Mode
	}
	switch mode {
	case "template":
		return HandleGetBlockTemplateRequest(s, request, closeChan)
	case "proposal":
		return HandleGetBlockTemplateProposal(s, request)
	}
	return nil, &btcjson.RPCError{
		Code:    btcjson.ErrRPCInvalidParameter,
		Message: "Invalid mode",
	}
}

// HandleGetBlockTemplateLongPoll is a helper for handleGetBlockTemplateRequest which deals with handling long polling
// for block templates. When a caller sends a request with a long poll ID that was previously returned, a response is
// not sent until the caller should stop working on the previous block template in favor of the new one.
//
// In particular, this is the case when the old block template is no longer valid due to a solution already being found
// and added to the block chain, or new transactions have shown up and some time has passed without finding a solution.
// See https://en.bitcoin.it/wiki/ BIP_0022 for more details.
func HandleGetBlockTemplateLongPoll(
	s *Server,
	longPollID string,
	useCoinbaseValue bool, closeChan qu.C,
) (interface{}, error) {
	state := s.GBTWorkState
	state.Lock()
	// The state unlock is intentionally not deferred here since it needs to be manually unlocked before waiting for a
	// notification about block template changes.
	if e := state.UpdateBlockTemplate(s, useCoinbaseValue); E.Chk(e) {
		state.Unlock()
		return nil, e
	}
	// Just return the current block template if the long poll ID provided by the caller is invalid.
	prevHash, lastGenerated, e := DecodeTemplateID(longPollID)
	var result *btcjson.GetBlockTemplateResult
	if e != nil {
		result, e = state.BlockTemplateResult(useCoinbaseValue, nil)
		if e != nil {
			state.Unlock()
			return nil, e
		}
		state.Unlock()
		return result, nil
	}
	// Return the block template now if the specific block template/ identified by the long poll ID no longer matches
	// the current block template as this means the provided template is stale.
	prevTemplateHash := &state.Template.Block.Header.PrevBlock
	if !prevHash.IsEqual(prevTemplateHash) ||
		lastGenerated != state.LastGenerated.Unix() {
		// Include whether or not it is valid to submit work against the old block template depending on whether or not
		// a solution has already been found and added to the block chain.
		submitOld := prevHash.IsEqual(prevTemplateHash)
		result, e = state.BlockTemplateResult(
			useCoinbaseValue,
			&submitOld,
		)
		if e != nil {
			state.Unlock()
			return nil, e
		}
		state.Unlock()
		return result, nil
	}
	// Register the previous hash and last generated time for notifications Get a channel that will be notified when the
	// template associated with the provided ID is stale and a new block template should be returned to the caller.
	longPollChan := state.TemplateUpdateChan(prevHash, lastGenerated)
	state.Unlock()
	select {
	// When the client closes before it's time to send a reply, just return now so the goroutine doesn't hang around.
	case <-closeChan.Wait():
		return nil, ErrClientQuit
	// Wait until signal received to send the reply.
	case <-longPollChan.Wait():
		// Fallthrough
	}
	// Get the lastest block template
	state.Lock()
	defer state.Unlock()
	if e = state.UpdateBlockTemplate(s, useCoinbaseValue); E.Chk(e) {
		return nil, e
	}
	// Include whether or not it is valid to submit work against the old block template depending on whether or not a
	// solution has already been found and added to the block chain.
	submitOld := prevHash.IsEqual(&state.Template.Block.Header.PrevBlock)
	result, e = state.BlockTemplateResult(useCoinbaseValue, &submitOld)
	if e != nil {
		return nil, e
	}
	return result, nil
}

// HandleGetBlockTemplateProposal is a helper for handleGetBlockTemplate which deals with block proposals. See
// https://en.bitcoin.it/wiki/BIP_0023 for more details.
func HandleGetBlockTemplateProposal(
	s *Server,
	request *btcjson.TemplateRequest,
) (interface{}, error) {
	hexData := request.Data
	if hexData == "" {
		return false, &btcjson.RPCError{
			Code:    btcjson.ErrRPCType,
			Message: fmt.Sprintf("Data must contain the hex-encoded serialized block that is being proposed"),
		}
	}
	// Ensure the provided data is sane and deserialize the proposed block.
	if len(hexData)%2 != 0 {
		hexData = "0" + hexData
	}
	dataBytes, e := hex.DecodeString(hexData)
	if e != nil {
		return false, &btcjson.RPCError{
			Code: btcjson.ErrRPCDeserialization,
			Message: fmt.Sprintf(
				"data must be hexadecimal string (not %q)",
				hexData,
			),
		}
	}
	var msgBlock wire.Block
	if e := msgBlock.Deserialize(bytes.NewReader(dataBytes)); E.Chk(e) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "block decode failed: " + e.Error(),
		}
	}
	block := block2.NewBlock(&msgBlock)
	// Ensure the block is building from the expected previous block.
	expectedPrevHash := s.Cfg.Chain.BestSnapshot().Hash
	prevHash := &block.WireBlock().Header.PrevBlock
	if !expectedPrevHash.IsEqual(prevHash) {
		return "bad-prevblk", nil
	}
	if e := s.Cfg.Chain.CheckConnectBlockTemplate(block); E.Chk(e) {
		if _, ok := e.(blockchain.RuleError); !ok {
			errStr := fmt.Sprintf("failed to process block proposal: %v", e)
			E.Ln(errStr)
			
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCVerify,
				Message: errStr,
			}
		}
		I.Ln("rejected block proposal:", e)
		
		return ChainErrToGBTErrString(e), nil
	}
	return nil, nil
}

// HandleGetBlockTemplateRequest is a helper for handleGetBlockTemplate which deals with generating and returning block
// templates to the caller. It handles both long poll requests as specified by BIP 0022 as well as regular requests.
//
// In addition, it detects the capabilities reported by the caller in regards to whether or not it supports creating its
// own coinbase (the coinbasetxn and coinbasevalue capabilities) and modifies the returned block template accordingly.
func HandleGetBlockTemplateRequest(
	s *Server,
	request *btcjson.TemplateRequest,
	closeChan qu.C,
) (interface{}, error) {
	// Extract the relevant passed capabilities and restrict the result to either a coinbase value or a coinbase
	// transaction object depending on the request. Default to only providing a coinbase value.
	useCoinbaseValue := true
	if request != nil {
		var hasCoinbaseValue, hasCoinbaseTxn bool
		for _, capability := range request.Capabilities {
			switch capability {
			case "coinbasetxn":
				hasCoinbaseTxn = true
			case "coinbasevalue":
				hasCoinbaseValue = true
			}
		}
		if hasCoinbaseTxn && !hasCoinbaseValue {
			useCoinbaseValue = false
		}
	}
	// When a coinbase transaction has been requested, respond with an error if there are no addresses to pay the
	// created block template to.
	if !useCoinbaseValue && len(s.StateCfg.ActiveMiningAddrs) == 0 {
		return nil, &btcjson.RPCError{
			Code: btcjson.ErrRPCInternal.Code,
			Message: "A coinbase transaction has been requested, " +
				"but the Server has not been configured with " +
				"any payment addresses via --miningaddr",
		}
	}
	// Return an error if there are no peers connected since there is no way to relay a found block or receive
	// transactions to work on. However, allow this workState when running in the regression test or simulation test
	// mode.
	netwk := (s.Config.Network.V())[0]
	if !(netwk == 'r' || netwk == 's') &&
		s.Cfg.ConnMgr.ConnectedCount() == 0 {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCClientNotConnected,
			Message: "Pod is not connected to network",
		}
	}
	// No point in generating or accepting work before the chain is synced.
	currentHeight := s.Cfg.Chain.BestSnapshot().Height
	if currentHeight != 0 && !s.Cfg.SyncMgr.IsCurrent() {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCClientInInitialDownload,
			Message: "Pod is not yet synchronised...",
		}
	}
	// When a long poll ID was provided, this is a long poll request by the client to be notified when block template
	// referenced by the ID should be replaced with a new one.
	if request != nil && request.LongPollID != "" {
		return HandleGetBlockTemplateLongPoll(
			s, request.LongPollID,
			useCoinbaseValue, closeChan,
		)
	}
	// Protect concurrent access when updating block templates.
	workState := s.GBTWorkState
	workState.Lock()
	defer workState.Unlock()
	// Get and return a block template. A new block template will be generated when the current best block has changed
	// or the transactions in the memory pool have been updated and it has been at least five seconds since the last
	// template was generated.
	//
	// Otherwise, the timestamp for the existing block template is updated (and possibly the difficulty on testnet per
	// the consesus rules).
	if e := workState.UpdateBlockTemplate(s, useCoinbaseValue); E.Chk(e) {
		return nil, e
	}
	return workState.BlockTemplateResult(useCoinbaseValue, nil)
}

// HandleGetCFilter implements the getcfilter command.
func HandleGetCFilter(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	if s.Cfg.CfIndex == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoCFIndex,
			Message: "The CF index must be enabled for this command",
		}
	}
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetCFilterCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getfilter")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	hash, e := chainhash.NewHashFromStr(c.Hash)
	if e != nil {
		return nil, DecodeHexError(c.Hash)
	}
	filterBytes, e := s.Cfg.CfIndex.FilterByBlockHash(hash, c.FilterType)
	if e != nil {
		D.F("could not find committed filter for %v: %v", hash, e)
		
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCBlockNotFound,
			Message: "block not found",
		}
	}
	T.Ln("found committed filter for", hash)
	return hex.EncodeToString(filterBytes), nil
}

// HandleGetCFilterHeader implements the getcfilterheader command.
func HandleGetCFilterHeader(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	if s.Cfg.CfIndex == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoCFIndex,
			Message: "The CF index must be enabled for this command",
		}
	}
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetCFilterHeaderCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getcfilterheader")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	hash, e := chainhash.NewHashFromStr(c.Hash)
	if e != nil {
		return nil, DecodeHexError(c.Hash)
	}
	headerBytes, e := s.Cfg.CfIndex.FilterHeaderByBlockHash(hash, c.FilterType)
	if len(headerBytes) > 0 {
		D.Ln("found header of committed filter for", hash)
		
	} else {
		D.F(
			"could not find header of committed filter for %v: %v",
			hash,
			e,
		)
		
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCBlockNotFound,
			Message: "Block not found",
		}
	}
	
	e = hash.SetBytes(headerBytes)
	if e != nil {
		D.Ln(e)
		
	}
	return hash.String(), nil
}

// HandleGetConnectionCount implements the getconnectioncount command.
func HandleGetConnectionCount(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	return s.Cfg.ConnMgr.ConnectedCount(), nil
}

// HandleGetCurrentNet implements the getcurrentnet command.
func HandleGetCurrentNet(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	return s.Cfg.ChainParams.Net, nil
}

// HandleGetDifficulty implements the getdifficulty command.
// TODO: This command should default to the configured algo for cpu mining
//  and take an optional parameter to query by algo
func HandleGetDifficulty(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetDifficultyCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getdifficulty")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	best := s.Cfg.Chain.BestSnapshot()
	prev, e := s.Cfg.Chain.BlockByHash(&best.Hash)
	if e != nil {
		E.Ln("ERROR", e)
		
	}
	var algo = prev.WireBlock().Header.Version
	if algo != 514 {
		algo = 2
	}
	bestbits := best.Bits
	if c.Algo == fork.Scrypt && algo != 514 {
		algo = 514
		for {
			if prev.WireBlock().Header.Version != 514 {
				ph := prev.WireBlock().Header.PrevBlock
				prev, e = s.Cfg.Chain.BlockByHash(&ph)
				if e != nil {
					E.Ln("ERROR", e)
					
				}
				continue
			}
			bestbits = prev.WireBlock().Header.Bits
			break
		}
	}
	if c.Algo == fork.SHA256d && algo != 2 {
		algo = 2
		for {
			if prev.WireBlock().Header.Version == 514 {
				ph := prev.WireBlock().Header.PrevBlock
				prev, e = s.Cfg.Chain.BlockByHash(&ph)
				if e != nil {
					E.Ln("ERROR", e)
					
				}
				continue
			}
			bestbits = prev.WireBlock().Header.Bits
			break
		}
	}
	return GetDifficultyRatio(bestbits, s.Cfg.ChainParams, algo), nil
}

// HandleGetGenerate implements the getgenerate command.
func HandleGetGenerate(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) { // cpuminer
	_, ok := cmd.(*btcjson.GetGenerateCmd)
	if ok {
		result := *s.Config.Controller
		return &result, nil
	}
	// generating := s.StateCfg.Miner != nil
	// if generating {
	//	D.Ln("miner is running internally")
	// } else {
	//	D.Ln("miner is not running")
	// }
	// return nil, nil
	// return s.Cfg.CPUMiner.IsMining(), nil
	return false, errors.New("command was not a btcjson.GetGenerateCmd")
}

// var startTime = time.Now()

// HandleGetHashesPerSec implements the gethashespersec command.
func HandleGetHashesPerSec(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) { // cpuminer
	// return int64(s.,
	// Cfg.CPUMiner.HashesPerSecond()), nil
	// TODO: finish this - needs generator for momentary rate (ewma)
	D.Ln("miner hashes per second - multicast thing TODO")
	// simple average for now
	return int(s.Cfg.Hashrate.Load()), nil
}

// HandleGetHeaders implements the getheaders command.
//
// NOTE: This is a btcsuite extension originally ported from github.com/decred/dcrd.
func HandleGetHeaders(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetHeadersCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getheaders")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Fetch the requested headers from chain while respecting the provided block locators and stop hash.
	blockLocators := make([]*chainhash.Hash, len(c.BlockLocators))
	for i := range c.BlockLocators {
		blockLocator, e := chainhash.NewHashFromStr(c.BlockLocators[i])
		if e != nil {
			return nil, DecodeHexError(c.BlockLocators[i])
		}
		blockLocators[i] = blockLocator
	}
	var hashStop chainhash.Hash
	if c.HashStop != "" {
		e := chainhash.Decode(&hashStop, c.HashStop)
		if e != nil {
			return nil, DecodeHexError(c.HashStop)
		}
	}
	headers := s.Cfg.SyncMgr.LocateHeaders(blockLocators, &hashStop)
	// Return the serialized block headers as hex-encoded strings.
	hexBlockHeaders := make([]string, len(headers))
	var buf bytes.Buffer
	for i, h := range headers {
		e := h.Serialize(&buf)
		if e != nil {
			return nil, InternalRPCError(
				e.Error(),
				"Failed to serialize block header",
			)
		}
		hexBlockHeaders[i] = hex.EncodeToString(buf.Bytes())
		buf.Reset()
	}
	return hexBlockHeaders, nil
}

// HandleGetInfo implements the getinfo command. We only return the fields that are not related to wallet functionality.
// TODO: simplify this, break it up
func HandleGetInfo(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (ret interface{}, e error) {
	var Difficulty, dBlake2b, dBlake14lr, dBlake2s, dKeccak, dScrypt, dSHA256D,
	dSkein, dStribog, dX11 float64
	var lastbitsScrypt, lastbitsSHA256D uint32
	best := s.
		Cfg.
		Chain.
		BestSnapshot()
	v := s.Cfg.Chain.Index.LookupNode(&best.Hash)
	foundcount, height := 0, best.Height
	switch fork.GetCurrent(height) {
	case 0:
		for foundcount < 9 && height > 0 {
			switch fork.GetAlgoName(v.Header().Version, height) {
			case fork.SHA256d:
				if lastbitsSHA256D == 0 {
					foundcount++
					lastbitsSHA256D = v.Header().Bits
					dSHA256D = GetDifficultyRatio(
						lastbitsSHA256D,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			case fork.Scrypt:
				if lastbitsScrypt == 0 {
					foundcount++
					lastbitsScrypt = v.Header().Bits
					dScrypt = GetDifficultyRatio(
						lastbitsScrypt,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			default:
			}
			v = v.RelativeAncestor(1)
			height--
		}
		switch s.Cfg.Algo {
		case fork.SHA256d:
			Difficulty = dSHA256D
		case fork.Scrypt:
			Difficulty = dScrypt
		default:
		}
		ret = &btcjson.InfoChainResult0{
			// Version: int32(
			// version.Tag,
			// 1000000*version.AppMajor +
			// 	10000*version.AppMinor +
			// 	100*version.AppPatch,
			// ),
			ProtocolVersion:   int32(MaxProtocolVersion),
			Blocks:            best.Height,
			TimeOffset:        int64(s.Cfg.TimeSource.Offset().Seconds()),
			Connections:       s.Cfg.ConnMgr.ConnectedCount(),
			Proxy:             s.Config.ProxyAddress.V(),
			PowAlgoID:         fork.GetAlgoID(s.Cfg.Algo, height),
			PowAlgo:           s.Cfg.Algo,
			Difficulty:        Difficulty,
			DifficultySHA256D: dSHA256D,
			DifficultyScrypt:  dScrypt,
			TestNet:           (s.Config.Network.V())[0] == 't',
			RelayFee:          s.StateCfg.ActiveMinRelayTxFee.ToDUO(),
		}
	case 1:
		foundcount, height := 0, best.Height
		for foundcount < 9 &&
			height > fork.List[fork.GetCurrent(height)].ActivationHeight-512 {
			switch fork.GetAlgoName(v.Header().Version, height) {
			case fork.Scrypt:
				if lastbitsScrypt == 0 {
					foundcount++
					lastbitsScrypt = v.Header().Bits
					dScrypt = GetDifficultyRatio(
						lastbitsScrypt,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			case fork.SHA256d:
				if lastbitsSHA256D == 0 {
					foundcount++
					lastbitsSHA256D = v.Header().Bits
					dSHA256D = GetDifficultyRatio(
						lastbitsSHA256D,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			default:
			}
			v = v.RelativeAncestor(1)
			height--
		}
		switch s.Cfg.Algo {
		case fork.Scrypt:
			Difficulty = dScrypt
		case fork.SHA256d:
			Difficulty = dSHA256D
		default:
		}
		ret = &btcjson.InfoChainResult{
			// Version: int32(
			// 	1000000*version.AppMajor +
			// 		10000*version.AppMinor +
			// 		100*version.AppPatch,
			// ),
			ProtocolVersion:     int32(MaxProtocolVersion),
			Blocks:              best.Height,
			TimeOffset:          int64(s.Cfg.TimeSource.Offset().Seconds()),
			Connections:         s.Cfg.ConnMgr.ConnectedCount(),
			Proxy:               s.Config.ProxyAddress.V(),
			PowAlgoID:           fork.GetAlgoID(s.Cfg.Algo, height),
			PowAlgo:             s.Cfg.Algo,
			Difficulty:          Difficulty,
			DifficultyBlake2b:   dBlake2b,
			DifficultyBlake14lr: dBlake14lr,
			DifficultyBlake2s:   dBlake2s,
			DifficultyKeccak:    dKeccak,
			DifficultyScrypt:    dScrypt,
			DifficultySHA256D:   dSHA256D,
			DifficultySkein:     dSkein,
			DifficultyStribog:   dStribog,
			DifficultyX11:       dX11,
			TestNet:             (s.Config.Network.V())[0] == 't',
			RelayFee:            s.StateCfg.ActiveMinRelayTxFee.ToDUO(),
		}
	}
	return ret, nil
}

// HandleGetMempoolInfo implements the getmempoolinfo command.
func HandleGetMempoolInfo(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	mempoolTxns := s.Cfg.TxMemPool.TxDescs()
	var numBytes int64
	for _, txD := range mempoolTxns {
		numBytes += int64(txD.Tx.MsgTx().SerializeSize())
	}
	ret := &btcjson.GetMempoolInfoResult{
		Size:  int64(len(mempoolTxns)),
		Bytes: numBytes,
	}
	return ret, nil
}

// HandleGetMiningInfo implements the getmininginfo command. We only return the fields that are not related to wallet
// functionality. This function returns more information than parallelcoind. TODO: simplify this, break it up
func HandleGetMiningInfo(
	s *Server, cmd interface{},
	closeChan qu.C,
) (ret interface{}, e error) {
	// cpuminer
	// Create a default getnetworkhashps command to use defaults and make use of the existing getnetworkhashps handler.
	gnhpsCmd := btcjson.NewGetNetworkHashPSCmd(nil, nil)
	networkHashesPerSecIface, e := HandleGetNetworkHashPS(s, gnhpsCmd, closeChan)
	if e != nil {
		return nil, e
	}
	networkHashesPerSec, ok := networkHashesPerSecIface.(int64)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInternal.Code,
			Message: "networkHashesPerSec is not an int64",
		}
	}
	var Difficulty, dScrypt, dSHA256D float64
	var lastbitsScrypt, lastbitsSHA256D uint32
	best := s.Cfg.Chain.BestSnapshot()
	v := s.Cfg.Chain.Index.LookupNode(&best.Hash)
	foundCount, height := 0, best.Height
	switch fork.GetCurrent(height) {
	case 0:
		for foundCount < 2 && height > 0 {
			switch fork.GetAlgoName(v.Header().Version, height) {
			case fork.SHA256d:
				if lastbitsSHA256D == 0 {
					foundCount++
					lastbitsSHA256D = v.Header().Bits
					dSHA256D = GetDifficultyRatio(
						lastbitsSHA256D,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			case fork.Scrypt:
				if lastbitsScrypt == 0 {
					foundCount++
					lastbitsScrypt = v.Header().Bits
					dScrypt = GetDifficultyRatio(
						lastbitsScrypt,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			default:
			}
			v = v.RelativeAncestor(1)
			height--
		}
		switch s.Cfg.Algo {
		case fork.SHA256d:
			Difficulty = dSHA256D
		case fork.Scrypt:
			Difficulty = dScrypt
		default:
		}
		D.Ln("missing generate stats in here")
		ret = &btcjson.GetMiningInfoResult0{
			Blocks:             int64(best.Height),
			CurrentBlockSize:   best.BlockSize,
			CurrentBlockWeight: best.BlockWeight,
			CurrentBlockTx:     best.NumTxns,
			PowAlgoID:          fork.GetAlgoID(s.Cfg.Algo, height),
			PowAlgo:            s.Cfg.Algo,
			Difficulty:         Difficulty,
			DifficultySHA256D:  dSHA256D,
			DifficultyScrypt:   dScrypt,
			// Generate:           s.Cfg.CPUMiner.IsMining(),
			// GenProcLimit:       s.Cfg.CPUMiner.NumWorkers(),
			// HashesPerSec:       int64(s.Cfg.CPUMiner.HashesPerSecond()),
			NetworkHashPS: networkHashesPerSec,
			PooledTx:      uint64(s.Cfg.TxMemPool.Count()),
			TestNet:       (s.Config.Network.V())[0] == 't',
		}
	case 1:
		fc, height := 0, best.Height
		for fc < 9 && height > fork.List[fork.GetCurrent(height)].ActivationHeight-512 {
			switch fork.GetAlgoName(v.Header().Version, height) {
			case fork.Scrypt:
				if lastbitsScrypt == 0 {
					fc++
					lastbitsScrypt = v.Header().Bits
					dScrypt = GetDifficultyRatio(
						lastbitsScrypt,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			case fork.SHA256d:
				if lastbitsSHA256D == 0 {
					fc++
					lastbitsSHA256D = v.Header().Bits
					dSHA256D = GetDifficultyRatio(
						lastbitsSHA256D,
						s.Cfg.ChainParams, v.Header().Version,
					)
				}
			default:
			}
			v = v.RelativeAncestor(1)
			height--
		}
		switch s.Cfg.Algo {
		case fork.Scrypt:
			Difficulty = dScrypt
		case fork.SHA256d:
			Difficulty = dSHA256D
		default:
		}
		D.Ln("missing cpu miner stuff in here") // cpuminer
		ret = &btcjson.GetMiningInfoResult{
			Blocks:             int64(best.Height),
			CurrentBlockSize:   best.BlockSize,
			CurrentBlockWeight: best.BlockWeight,
			CurrentBlockTx:     best.NumTxns,
			PowAlgoID:          fork.GetAlgoID(s.Cfg.Algo, height),
			PowAlgo:            s.Cfg.Algo,
			Difficulty:         Difficulty,
			DifficultyScrypt:   dScrypt,
			DifficultySHA256D:  dSHA256D,
			NetworkHashPS:      networkHashesPerSec,
			PooledTx:           uint64(s.Cfg.TxMemPool.Count()),
			TestNet:            (s.Config.Network.V())[0] == 't',
		}
	}
	return ret, nil
}

// HandleGetNetTotals implements the getnettotals command.
func HandleGetNetTotals(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	totalBytesRecv, totalBytesSent := s.Cfg.ConnMgr.NetTotals()
	reply := &btcjson.GetNetTotalsResult{
		TotalBytesRecv: totalBytesRecv,
		TotalBytesSent: totalBytesSent,
		TimeMillis:     time.Now().UTC().UnixNano() / int64(time.Millisecond),
	}
	return reply, nil
}

// HandleGetNetworkHashPS implements the getnetworkhashps command. This command does not default to the same end block
// as the parallelcoind. TODO: Really this needs to be expanded to show per-algorithm hashrates
func HandleGetNetworkHashPS(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetNetworkHashPSCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getnetworkhashps")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Note: All valid error return paths should return an int64. Literal zeros are inferred as int, and won't coerce to
	// int64 because the return value is an interface{}.
	//
	// When the passed height is too high or zero, just return 0 now since we can't reasonably calculate the number of
	// network hashes per second from invalid values. When it's negative, use the current best block height.
	best := s.Cfg.Chain.BestSnapshot()
	endHeight := int32(-1)
	if c.Height != nil {
		endHeight = int32(*c.Height)
	}
	if endHeight > best.Height || endHeight == 0 {
		return int64(0), nil
	}
	if endHeight < 0 {
		endHeight = best.Height
	}
	// Calculate the number of blocks per retarget interval based on the chain parameters.
	blocksPerRetarget := int32(
		s.Cfg.ChainParams.TargetTimespan / s.Cfg.
			ChainParams.TargetTimePerBlock,
	)
	// Calculate the starting block height based on the passed number of blocks.
	//
	// When the passed value is negative, use the last block the difficulty changed as the starting height. Also make
	// sure the starting height is not before the beginning of the chain.
	numBlocks := int32(120)
	if c.Blocks != nil {
		numBlocks = int32(*c.Blocks)
	}
	var startHeight int32
	if numBlocks <= 0 {
		startHeight = endHeight - ((endHeight % blocksPerRetarget) + 1)
	} else {
		startHeight = endHeight - numBlocks
	}
	if startHeight < 0 {
		startHeight = 0
	}
	T.F(
		"calculating network hashes per second from %d to %d",
		startHeight,
		endHeight,
	)
	
	// Find the min and max block timestamps as well as calculate the total amount of work that happened between the
	// start and end blocks.
	var minTimestamp, maxTimestamp time.Time
	totalWork := big.NewInt(0)
	for curHeight := startHeight; curHeight <= endHeight; curHeight++ {
		hash, e := s.Cfg.Chain.BlockHashByHeight(curHeight)
		if e != nil {
			context := "Failed to fetch block hash"
			return nil, InternalRPCError(e.Error(), context)
		}
		// Fetch the header from chain.
		header, e := s.Cfg.Chain.HeaderByHash(hash)
		if e != nil {
			context := "Failed to fetch block header"
			return nil, InternalRPCError(e.Error(), context)
		}
		if curHeight == startHeight {
			minTimestamp = header.Timestamp
			maxTimestamp = minTimestamp
		} else {
			totalWork.Add(
				totalWork, blockchain.CalcWork(
					header.Bits,
					best.Height+1, header.Version,
				),
			)
			if minTimestamp.After(header.Timestamp) {
				minTimestamp = header.Timestamp
			}
			if maxTimestamp.Before(header.Timestamp) {
				maxTimestamp = header.Timestamp
			}
		}
	}
	// Calculate the difference in seconds between the min and max block timestamps and avoid division by zero in the
	// case where there is no time difference.
	timeDiff := int64(maxTimestamp.Sub(minTimestamp) / time.Second)
	if timeDiff == 0 {
		return int64(0), nil
	}
	hashesPerSec := new(big.Int).Div(totalWork, big.NewInt(timeDiff))
	return hashesPerSec.Int64(), nil
}

// HandleGetPeerInfo implements the getpeerinfo command.
func HandleGetPeerInfo(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	peers := s.Cfg.ConnMgr.ConnectedPeers()
	syncPeerID := s.Cfg.SyncMgr.SyncPeerID()
	infos := make([]*btcjson.GetPeerInfoResult, 0, len(peers))
	for _, p := range peers {
		statsSnap := p.ToPeer().StatsSnapshot()
		var addr, addrLocal string
		if statsSnap.Inbound {
			addr =
				statsSnap.Addr
			addrLocal =
				p.ToPeer().LocalAddr().String()
			// (*s.Config.P2PConnect)[0]
		} else {
			addr =
				statsSnap.Addr
			addrLocal =
			// (*s.Config.P2PConnect)[0]
				p.ToPeer().LocalAddr().String()
		}
		info := &btcjson.GetPeerInfoResult{
			ID:             statsSnap.ID,
			Addr:           addr,
			AddrLocal:      addrLocal,
			Services:       fmt.Sprintf("%08d", uint64(statsSnap.Services)),
			RelayTxes:      !p.IsTxRelayDisabled(),
			LastSend:       statsSnap.LastSend.Unix(),
			LastRecv:       statsSnap.LastRecv.Unix(),
			BytesSent:      statsSnap.BytesSent,
			BytesRecv:      statsSnap.BytesRecv,
			ConnTime:       statsSnap.ConnTime.Unix(),
			PingTime:       float64(statsSnap.LastPingMicros),
			TimeOffset:     statsSnap.TimeOffset,
			Version:        statsSnap.Version,
			SubVer:         statsSnap.UserAgent,
			Inbound:        statsSnap.Inbound,
			StartingHeight: statsSnap.StartingHeight,
			CurrentHeight:  statsSnap.LastBlock,
			BanScore:       int32(p.GetBanScore()),
			FeeFilter:      p.GetFeeFilter(),
			SyncNode:       statsSnap.ID == syncPeerID,
		}
		if p.ToPeer().LastPingNonce() != 0 {
			wait := float64(time.Since(statsSnap.LastPingTime).Nanoseconds())
			// We actually want microseconds.
			info.PingWait = wait / 1000
		}
		infos = append(infos, info)
	}
	return infos, nil
}

// HandleGetRawMempool implements the getrawmempool command.
func HandleGetRawMempool(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	c := cmd.(*btcjson.GetRawMempoolCmd)
	mp := s.Cfg.TxMemPool
	if c.Verbose != nil && *c.Verbose {
		return mp.RawMempoolVerbose(), nil
	}
	// The response is simply an array of the transaction hashes if the verbose flag is not set.
	descs := mp.TxDescs()
	hashStrings := make([]string, len(descs))
	for i := range hashStrings {
		hashStrings[i] = descs[i].Tx.Hash().String()
	}
	return hashStrings, nil
}

// HandleGetRawTransaction implements the getrawtransaction command.
func HandleGetRawTransaction(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("getrawtransaction")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Convert the provided transaction hash hex to a Hash.
	txHash, e := chainhash.NewHashFromStr(c.Txid)
	if e != nil {
		return nil, DecodeHexError(c.Txid)
	}
	verbose := false
	if c.Verbose != nil {
		verbose = *c.Verbose != 0
	}
	// Try to fetch the transaction from the memory pool and if that fails, try the block database.
	var mtx *wire.MsgTx
	var blkHash *chainhash.Hash
	var blkHeight int32
	tx, e := s.Cfg.TxMemPool.FetchTransaction(txHash)
	if e != nil {
		if s.Cfg.TxIndex == nil {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCNoTxInfo,
				Message: "The transaction index must be " +
					"enabled to query the blockchain " +
					"(specify --txindex)",
			}
		}
		// Look up the location of the transaction.
		var blockRegion *database.BlockRegion
		blockRegion, e = s.Cfg.TxIndex.TxBlockRegion(txHash)
		if e != nil {
			context := "Failed to retrieve transaction location"
			return nil, InternalRPCError(e.Error(), context)
		}
		if blockRegion == nil {
			return nil, NoTxInfoError(txHash)
		}
		// Load the raw transaction bytes from the database.
		var txBytes []byte
		e = s.Cfg.DB.View(
			func(dbTx database.Tx) (e error) {
				txBytes, e = dbTx.FetchBlockRegion(blockRegion)
				return e
			},
		)
		if e != nil {
			return nil, NoTxInfoError(txHash)
		}
		// When the verbose flag isn't set, simply return the serialized transaction as a hex-encoded string. This is
		// done here to avoid deserializing it only to reserialize it again later.
		if !verbose {
			return hex.EncodeToString(txBytes), nil
		}
		// Grab the block height.
		blkHash = blockRegion.Hash
		blkHeight, e = s.Cfg.Chain.BlockHeightByHash(blkHash)
		if e != nil {
			context := "Failed to retrieve block height"
			return nil, InternalRPCError(e.Error(), context)
		}
		// Deserialize the transaction
		var msgTx wire.MsgTx
		e = msgTx.Deserialize(bytes.NewReader(txBytes))
		if e != nil {
			context := deserialfail
			return nil, InternalRPCError(e.Error(), context)
		}
		mtx = &msgTx
	} else {
		// When the verbose flag isn't set, simply return the network-serialized transaction as a hex-encoded string.
		if !verbose {
			// Note that this is intentionally not directly returning because the first return value is a string and it
			// would result in returning an empty string to the client instead of nothing (nil) in the case of an error.
			var mtxHex string
			mtxHex, e = MessageToHex(tx.MsgTx())
			if e != nil {
				return nil, e
			}
			return mtxHex, nil
		}
		mtx = tx.MsgTx()
	}
	// The verbose flag is set, so generate the JSON object and return it.
	var blkHeader *wire.BlockHeader
	var blkHashStr string
	var chainHeight int32
	if blkHash != nil {
		// Fetch the header from chain.
		var header wire.BlockHeader
		header, e = s.Cfg.Chain.HeaderByHash(blkHash)
		if e != nil {
			context := "Failed to fetch block header"
			return nil, InternalRPCError(e.Error(), context)
		}
		blkHeader = &header
		blkHashStr = blkHash.String()
		chainHeight = s.Cfg.Chain.BestSnapshot().Height
	}
	rawTxn, e := CreateTxRawResult(
		s.Cfg.ChainParams, mtx, txHash.String(),
		blkHeader, blkHashStr, blkHeight, chainHeight,
	)
	if e != nil {
		return nil, e
	}
	return *rawTxn, nil
}

// HandleGetTxOut handles gettxout commands.
func HandleGetTxOut(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.GetTxOutCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("gettxout")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Convert the provided transaction hash hex to a Hash.
	txHash, e := chainhash.NewHashFromStr(c.Txid)
	if e != nil {
		return nil, DecodeHexError(c.Txid)
	}
	// If requested and the tx is available in the mempool try to fetch it from there, otherwise attempt to fetch from
	// the block database.
	var bestBlockHash string
	var confirmations int32
	var value int64
	var pkScript []byte
	var isCoinbase bool
	includeMempool := true
	if c.IncludeMempool != nil {
		includeMempool = *c.IncludeMempool
	}
	// TODO: This is racy.  It should attempt to fetch it directly and check the error.
	if includeMempool && s.Cfg.TxMemPool.HaveTransaction(txHash) {
		tx, e := s.Cfg.TxMemPool.FetchTransaction(txHash)
		if e != nil {
			return nil, NoTxInfoError(txHash)
		}
		mtx := tx.MsgTx()
		if c.Vout > uint32(len(mtx.TxOut)-1) {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInvalidTxVout,
				Message: "Output index number (vout) does not " +
					"exist for transaction.",
			}
		}
		txOut := mtx.TxOut[c.Vout]
		if txOut == nil {
			errStr := fmt.Sprintf(
				"Output index: %d for txid: %s "+
					"does not exist", c.Vout, txHash,
			)
			return nil, InternalRPCError(errStr, "")
		}
		best := s.Cfg.Chain.BestSnapshot()
		bestBlockHash = best.Hash.String()
		confirmations = 0
		value = txOut.Value
		pkScript = txOut.PkScript
		isCoinbase = blockchain.IsCoinBaseTx(mtx)
	} else {
		out := wire.OutPoint{Hash: *txHash, Index: c.Vout}
		entry, e := s.Cfg.Chain.FetchUtxoEntry(out)
		if e != nil {
			return nil, NoTxInfoError(txHash)
		}
		// To match the behavior of the reference client, return nil (JSON null) if the transaction output is spent by
		// another transaction already in the main chain. Mined transactions that are spent by a mempool transaction are
		// not affected by this.
		if entry == nil || entry.IsSpent() {
			return nil, nil
		}
		best := s.Cfg.Chain.BestSnapshot()
		bestBlockHash = best.Hash.String()
		confirmations = 1 + best.Height - entry.BlockHeight()
		value = entry.Amount()
		pkScript = entry.PkScript()
		isCoinbase = entry.IsCoinBase()
	}
	// Disassemble script into single line printable format. The disassembled string will contain [error] inline if the
	// script doesn't fully parse, so ignore the error here.
	disbuf, _ := txscript.DisasmString(pkScript)
	// Get further info about the script. Ignore the error here since an error means the script couldn't parse and there
	// is no additional information about it anyways.
	scriptClass, addrs, reqSigs, _ := txscript.ExtractPkScriptAddrs(pkScript, s.Cfg.ChainParams)
	addresses := make([]string, len(addrs))
	for i, addr := range addrs {
		addresses[i] = addr.EncodeAddress()
	}
	txOutReply := &btcjson.GetTxOutResult{
		BestBlock:     bestBlockHash,
		Confirmations: int64(confirmations),
		Value:         amt.Amount(value).ToDUO(),
		ScriptPubKey: btcjson.ScriptPubKeyResult{
			Asm:       disbuf,
			Hex:       hex.EncodeToString(pkScript),
			ReqSigs:   int32(reqSigs),
			Type:      scriptClass.String(),
			Addresses: addresses,
		},
		Coinbase: isCoinbase,
	}
	return txOutReply, nil
}

// HandleHelp implements the help command.
func HandleHelp(s *Server, cmd interface{}, closeChan qu.C) (
	interface{}, error,
) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.HelpCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCUsage(true)
		if e != nil {
			msg = e.Error() + "\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Provide a usage overview of all commands when no specific command was specified.
	var command string
	if c.Command != nil {
		command = *c.Command
	}
	if command == "" {
		var usage string
		usage, e = s.HelpCacher.RPCUsage(false)
		if e != nil {
			context := "Failed to generate RPC usage"
			return nil, InternalRPCError(e.Error(), context)
		}
		return usage, nil
	}
	// Chk that the command asked for is supported and implemented. Only search the main list of handlers since help
	// should not be provided for commands that are unimplemented or related to wallet functionality.
	if _, ok := RPCHandlers[command]; !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Unknown command: " + command,
		}
	}
	// Get the help for the command.
	help, e := s.HelpCacher.RPCMethodHelp(command)
	if e != nil {
		context := "Failed to generate help"
		return nil, InternalRPCError(e.Error(), context)
	}
	return help, nil
}

// HandleNode handles node commands.
func HandleNode(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.NodeCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("help")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	var addr string
	var nodeID uint64
	var errN error
	params := s.Cfg.ChainParams
	switch c.SubCmd {
	case "disconnect":
		// If we have a valid uint disconnect by node id. Otherwise, attempt to disconnect by address, returning an
		// error if a valid IP address is not supplied.
		if nodeID, errN = strconv.ParseUint(c.Target, 10, 32); errN == nil {
			e = s.Cfg.ConnMgr.DisconnectByID(int32(nodeID))
		} else {
			if _, _, errP := net.SplitHostPort(c.Target); errP == nil || net.ParseIP(c.Target) != nil {
				addr = NormalizeAddress(c.Target, params.DefaultPort)
				e = s.Cfg.ConnMgr.DisconnectByAddr(addr)
			} else {
				return nil, &btcjson.RPCError{
					Code:    btcjson.ErrRPCInvalidParameter,
					Message: "invalid address or node ID",
				}
			}
		}
		if e != nil && PeerExists(s.Cfg.ConnMgr, addr, int32(nodeID)) {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCMisc,
				Message: "can't disconnect a permanent peer, use remove",
			}
		}
	case "remove":
		// If we have a valid uint disconnect by node id. Otherwise, attempt to disconnect by address, returning an
		// error if a valid IP address is not supplied.
		if nodeID, errN = strconv.ParseUint(c.Target, 10, 32); errN == nil {
			e = s.Cfg.ConnMgr.RemoveByID(int32(nodeID))
		} else {
			if _, _, errP := net.SplitHostPort(c.Target); errP == nil || net.ParseIP(c.Target) != nil {
				addr = NormalizeAddress(c.Target, params.DefaultPort)
				e = s.Cfg.ConnMgr.RemoveByAddr(addr)
			} else {
				return nil, &btcjson.RPCError{
					Code:    btcjson.ErrRPCInvalidParameter,
					Message: "invalid address or node ID",
				}
			}
		}
		if e != nil && PeerExists(s.Cfg.ConnMgr, addr, int32(nodeID)) {
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCMisc,
				Message: "can't remove a temporary peer, use disconnect",
			}
		}
	case "connect":
		addr = NormalizeAddress(c.Target, params.DefaultPort)
		// Default to temporary connections.
		subCmd := "temp"
		if c.ConnectSubCmd != nil {
			subCmd = *c.ConnectSubCmd
		}
		switch subCmd {
		case "perm", "temp":
			e = s.Cfg.ConnMgr.Connect(addr, subCmd == "perm")
		default:
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInvalidParameter,
				Message: "invalid subcommand for node connect",
			}
		}
	default:
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "invalid subcommand for node",
		}
	}
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: e.Error(),
		}
	}
	// no data returned unless an error.
	return nil, nil
}

// HandlePing implements the ping command.
func HandlePing(s *Server, cmd interface{}, closeChan qu.C) (
	interface{}, error,
) {
	// Ask Server to ping \o_
	nonce, e := wire.RandomUint64()
	if e != nil {
		return nil, InternalRPCError(
			"Not sending ping - failed to generate nonce: "+e.Error(), "",
		)
	}
	s.Cfg.ConnMgr.BroadcastMessage(wire.NewMsgPing(nonce))
	return nil, nil
}

// HandleSearchRawTransactions implements the searchrawtransactions command.
// TODO: simplify this, break it up
func HandleSearchRawTransactions(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	// Respond with an error if the address index is not enabled.
	addrIndex := s.Cfg.AddrIndex
	if addrIndex == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCMisc,
			Message: "Address index must be enabled (--addrindex)",
		}
	}
	// Override the flag for including extra previous output information in each input if needed.
	c := cmd.(*btcjson.SearchRawTransactionsCmd)
	vinExtra := false
	if c.VinExtra != nil {
		vinExtra = *c.VinExtra != 0
	}
	// Including the extra previous output information requires the transaction index. Currently the address index
	// relies on the transaction index, so this check is redundant, but it's better to be safe in case the address index
	// is ever changed to not rely on it.
	if vinExtra && s.Cfg.TxIndex == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCMisc,
			Message: "Transaction index must be enabled (--txindex)",
		}
	}
	// Attempt to decode the supplied address.
	params := s.Cfg.ChainParams
	addr, e := btcaddr.Decode(c.Address, params)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Invalid address or key: " + e.Error(),
		}
	}
	// Override the default number of requested entries if needed. Also, just return now if the number of requested
	// entries is zero to avoid extra work.
	numRequested := 100
	if c.Count != nil {
		numRequested = *c.Count
		if numRequested < 0 {
			numRequested = 1
		}
	}
	if numRequested == 0 {
		return nil, nil
	}
	// Override the default number of entries to skip if needed.
	var numToSkip int
	if c.Skip != nil {
		numToSkip = *c.Skip
		if numToSkip < 0 {
			numToSkip = 0
		}
	}
	// Override the reverse flag if needed.
	var reverse bool
	if c.Reverse != nil {
		reverse = *c.Reverse
	}
	// Add transactions from mempool first if client asked for reverse order. Otherwise, they will be added last (as
	// needed depending on the requested counts). NOTE: This code doesn't txsort by dependency. This might be something to
	// do in the future for the client's convenience, or leave it to the client.
	numSkipped := uint32(0)
	addressTxns := make([]RetrievedTx, 0, numRequested)
	if reverse {
		// Transactions in the mempool are not in a block header yet, so the block header field in the retrieved
		// transaction struct is left nil.
		mpTxns, mpSkipped := FetchMempoolTxnsForAddress(
			s, addr,
			uint32(numToSkip), uint32(numRequested),
		)
		numSkipped += mpSkipped
		for _, tx := range mpTxns {
			addressTxns = append(addressTxns, RetrievedTx{Tx: tx})
		}
	}
	// Fetch transactions from the database in the desired order if more are needed.
	if len(addressTxns) < numRequested {
		e = s.Cfg.DB.View(
			func(dbTx database.Tx) (e error) {
				regions, dbSkipped, e := addrIndex.TxRegionsForAddress(
					dbTx, addr,
					uint32(numToSkip)-numSkipped, uint32(numRequested-len(addressTxns)),
					reverse,
				)
				if e != nil {
					return e
				}
				// Load the raw transaction bytes from the database.
				serializedTxns, e := dbTx.FetchBlockRegions(regions)
				if e != nil {
					return e
				}
				// Add the transaction and the hash of the block it is contained in to the list. Note that the transaction
				// is left serialized here since the caller might have requested non-verbose output and hence there would
				// be/ no point in deserializing it just to reserialize it later.
				for i, serializedTx := range serializedTxns {
					addressTxns = append(
						addressTxns, RetrievedTx{
							TxBytes: serializedTx,
							BlkHash: regions[i].Hash,
						},
					)
				}
				numSkipped += dbSkipped
				return nil
			},
		)
		if e != nil {
			context := "Failed to load address index entries"
			return nil, InternalRPCError(e.Error(), context)
		}
	}
	// Add transactions from mempool last if client did not request reverse order and the number of results is still
	// under the number requested.
	if !reverse && len(addressTxns) < numRequested {
		// Transactions in the mempool are not in a block header yet, so the block header field in the retrieved
		// transaction struct is left nil.
		mpTxns, mpSkipped := FetchMempoolTxnsForAddress(
			s, addr,
			uint32(numToSkip)-numSkipped, uint32(numRequested-len(addressTxns)),
		)
		numSkipped += mpSkipped
		for _, tx := range mpTxns {
			addressTxns = append(addressTxns, RetrievedTx{Tx: tx})
		}
	}
	// Address has never been used if neither source yielded any results.
	if len(addressTxns) == 0 {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoTxInfo,
			Message: "No information available about address",
		}
	}
	// Serialize all of the transactions to hex.
	hexTxns := make([]string, len(addressTxns))
	for i := range addressTxns {
		// Simply encode the raw bytes to hex when the retrieved transaction is already in serialized form.
		rtx := &addressTxns[i]
		if rtx.TxBytes != nil {
			hexTxns[i] = hex.EncodeToString(rtx.TxBytes)
			continue
		}
		// Serialize the transaction first and convert to hex when the retrieved transaction is the deserialized
		// structure.
		hexTxns[i], e = MessageToHex(rtx.Tx.MsgTx())
		if e != nil {
			return nil, e
		}
	}
	// When not in verbose mode, simply return a list of serialized txns.
	if c.Verbose != nil && *c.Verbose == 0 {
		return hexTxns, nil
	}
	// Normalize the provided filter addresses (if any) to ensure there are no duplicates.
	filterAddrMap := make(map[string]struct{})
	if c.FilterAddrs != nil && len(*c.FilterAddrs) > 0 {
		for _, addr := range *c.FilterAddrs {
			filterAddrMap[addr] = struct{}{}
		}
	}
	// The verbose flag is set, so generate the JSON object and return it.
	best := s.Cfg.Chain.BestSnapshot()
	srtList := make([]btcjson.SearchRawTransactionsResult, len(addressTxns))
	for i := range addressTxns {
		// The deserialized transaction is needed, so deserialize the retrieved transaction if it's in serialized form
		// (which will be the case when it was lookup up from the database). Otherwise, use the existing deserialized
		// transaction.
		var rtx *RetrievedTx
		rtx = &addressTxns[i]
		var mtx *wire.MsgTx
		if rtx.Tx == nil {
			// Deserialize the transaction.
			mtx = new(wire.MsgTx)
			e = mtx.Deserialize(bytes.NewReader(rtx.TxBytes))
			if e != nil {
				context := deserialfail
				return nil, InternalRPCError(e.Error(), context)
			}
		} else {
			mtx = rtx.Tx.MsgTx()
		}
		result := &srtList[i]
		result.Hex = hexTxns[i]
		result.TxID = mtx.TxHash().String()
		result.Vin, e = CreateVinListPrevOut(
			s, mtx, params, vinExtra,
			filterAddrMap,
		)
		if e != nil {
			return nil, e
		}
		result.VOut = CreateVoutList(mtx, params, filterAddrMap)
		result.Version = mtx.Version
		result.LockTime = mtx.LockTime
		// Transactions grabbed from the mempool aren't yet in a block, so conditionally fetch block details here. This
		// will be reflected in the final JSON output (mempool won't have confirmations or block information).
		var blkHeader *wire.BlockHeader
		var blkHashStr string
		var blkHeight int32
		if blkHash := rtx.BlkHash; blkHash != nil {
			// Fetch the header from chain.
			header, e := s.Cfg.Chain.HeaderByHash(blkHash)
			if e != nil {
				return nil, &btcjson.RPCError{
					Code:    btcjson.ErrRPCBlockNotFound,
					Message: "Block not found",
				}
			}
			// Get the block height from chain.
			height, e := s.Cfg.Chain.BlockHeightByHash(blkHash)
			if e != nil {
				context := blockheightfail
				return nil, InternalRPCError(e.Error(), context)
			}
			blkHeader = &header
			blkHashStr = blkHash.String()
			blkHeight = height
		}
		// Add the block information to the result if there is any.
		if blkHeader != nil {
			// This is not a typo, they are identical in Bitcoin Core as well.
			result.Time = blkHeader.Timestamp.Unix()
			result.Blocktime = blkHeader.Timestamp.Unix()
			result.BlockHash = blkHashStr
			result.Confirmations = uint64(1 + best.Height - blkHeight)
		}
	}
	return srtList, nil
}

// HandleSendRawTransaction implements the sendrawtransaction command.
func HandleSendRawTransaction(
	s *Server,
	cmd interface{},
	closeChan qu.C,
) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.SendRawTransactionCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("sendrawtransaction")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Deserialize and send off to tx relay
	hexStr := c.HexTx
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	serializedTx, e := hex.DecodeString(hexStr)
	if e != nil {
		return nil, DecodeHexError(hexStr)
	}
	var msgTx wire.MsgTx
	e = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX decode failed: " + e.Error(),
		}
	}
	// Use 0 for the tag to represent local node.
	tx := util.NewTx(&msgTx)
	acceptedTxs, e := s.Cfg.TxMemPool.ProcessTransaction(s.Cfg.Chain, tx, false, false, 0)
	if e != nil {
		// When the error is a rule error, it means the transaction was simply rejected as opposed to something actually
		// going wrong, so log such. Otherwise, something really did go wrong, so log an actual error. In both cases, a
		// JSON-RPC error is returned to the client with the deserialization error code (to match bitcoind behavior).
		if _, ok := e.(mempool.RuleError); ok {
			D.F("rejected transaction %v: %v", tx.Hash(), e)
			
		} else {
			E.F(
				"failed to process transaction %v: %v", tx.Hash(), e,
			)
		}
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "TX rejected: " + e.Error(),
		}
	}
	// When the transaction was accepted it should be the first item in the returned array of accepted transactions.
	//
	// The only way this will not be true is if the API for ProcessTransaction changes and this code is not properly
	// updated, but ensure the condition holds as a safeguard.
	//
	// Also, since an error is being returned to the caller, ensure the transaction is removed from the memory pool.
	if len(acceptedTxs) == 0 || !acceptedTxs[0].Tx.Hash().IsEqual(tx.Hash()) {
		s.Cfg.TxMemPool.RemoveTransaction(tx, true)
		errStr := fmt.Sprintf("transaction %v is not in accepted list", tx.Hash())
		return nil, InternalRPCError(errStr, "")
	}
	// Generate and relay inventory vectors for all newly accepted transactions into the memory pool due to the original
	// being accepted.
	s.Cfg.ConnMgr.RelayTransactions(acceptedTxs)
	// Notify both websocket and getblocktemplate long poll clients of all newly accepted transactions.
	s.NotifyNewTransactions(acceptedTxs)
	// Keep track of all the sendrawtransaction request txns so that they can be rebroadcast if they don't make their
	// way into a block.
	txD := acceptedTxs[0]
	iv := wire.NewInvVect(wire.InvTypeTx, txD.Tx.Hash())
	s.Cfg.ConnMgr.AddRebroadcastInventory(iv, txD)
	return tx.Hash().String(), nil
}

// HandleSetGenerate implements the setgenerate command.
func HandleSetGenerate(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) { // cpuminer
	c, ok := cmd.(*btcjson.SetGenerateCmd)
	if ok {
		if c.Generate {
			s.Cfg.StartController.Signal()
		} else {
			s.Cfg.StopController.Signal()
		}
		return &struct{}{}, nil
	}
	return nil, errors.New("command was not a btcjson.SetGenerateCmd")
	// var msg string
	// var e error
	// // c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	// c, ok := cmd.(*btcjson.SetGenerateCmd)
	// if !ok {
	//	var h string
	//	h, e = s.HelpCacher.RPCMethodHelp("setgenerate")
	//	if e != nil  {
	//		msg = err.E.Ln() + "\n\n"
	//	}
	//	msg += h
	//	return nil, &btcjson.RPCError{
	//		Code:    btcjson.ErrRPCInvalidParameter,
	//		Message: msg,
	//		// "invalid subcommand for addnode",
	//	}
	// }
	// D.S(c)
	// // Disable generation regardless of the provided generate flag if the maximum number of threads (goroutines for our
	// // purposes) is 0. Otherwise enable or disable it depending on the provided flag. l.ScriptError(*c.GenProcLimit,
	// // c.Generate)
	// generate := c.Generate
	// genProcLimit := *s.Config.GenThreads
	// if c.GenProcLimit != nil {
	//	genProcLimit = *c.GenProcLimit
	//	if !generate {
	//		*c.GenProcLimit = 0
	//	}
	//	if *c.GenProcLimit == 0 {
	//		generate = false
	//	}
	// }
	// D.Ln("generating", generate, "threads", genProcLimit)
	// // if s.Cfg.CPUMiner.IsMining() {
	// // 	// if s.cfg.CPUMiner.GetAlgo() != s.cfg.Algo {
	// // 	s.Cfg.CPUMiner.Stop()
	// // 	generate = true
	// // 	// }
	// // }
	// // if !generate {
	// // 	s.Cfg.CPUMiner.Stop()
	// // } else {
	// // 	// Respond with an error if there are no addresses to pay the created
	// // 	// blocks to.
	// // 	if len(s.StateCfg.ActiveMiningAddrs) == 0 {
	// // 		return nil, &btcjson.RPCError{
	// // 			Code:    btcjson.ErrRPCInternal.Code,
	// // 			Message: "no payment addresses specified via --miningaddr",
	// // 		}
	// // 	}
	// // 	// It's safe to call start even if it's already started.
	// // 	s.Cfg.CPUMiner.SetNumWorkers(int32(genProcLimit))
	// // 	s.Cfg.CPUMiner.Start()
	// // }
	// //*s.Config.Generate = generate
	// //*s.Config.GenThreads = genProcLimit
	// //if s.StateCfg.Miner != nil {
	// //	D.Ln("stopping existing miner")
	// //	consume.Kill(s.StateCfg.Miner)
	// //	s.StateCfg.Miner = nil
	// //}
	// D.Ln("saving configuration")
	// save.Pod(s.Config)
	// //if *s.Config.Generate && *s.Config.GenThreads != 0 {
	// //	D.Ln("starting miner")
	// //	args := []string{os.Args[0], "-D", *s.Config.DataDir}
	// //	if *s.Config.KopachGUI {
	// //		args = append(args, "--kopachgui")
	// //	}
	// //	args = append(args, "kopach")
	// //	// args = apputil.PrependForWindows(args)
	// //	s.StateCfg.Miner = consume.Log(s.Quit, func(ent *log.Entry) (e error) {
	// //		D.Ln(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
	// //		return
	// //	}, func(pkg string) (out bool) {
	// //		return false
	// //	}, args...)
	// //	consume.Start(s.StateCfg.Miner)
	// //} else {
	// //	consume.Kill(s.StateCfg.Miner)
	// //}
	// return nil, nil
}

// HandleStop implements the stop command.
func HandleStop(s *Server, cmd interface{}, closeChan qu.C) (
	interface{}, error,
) {
	interrupt.Request()
	return nil, nil
}

// HandleRestart implements the restart command.
func HandleRestart(s *Server, cmd interface{}, closeChan qu.C) (
	interface{}, error,
) {
	// select {
	// case s.RequestProcessShutdown <- struct{}{}:
	// default:
	// }
	interrupt.RequestRestart()
	return nil, nil
}

// HandleSubmitBlock implements the submitblock command.
func HandleSubmitBlock(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.SubmitBlockCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("submitblock")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Deserialize the submitted block.
	hexStr := c.HexBlock
	if len(hexStr)%2 != 0 {
		hexStr = "0" + c.HexBlock
	}
	serializedBlock, e := hex.DecodeString(hexStr)
	if e != nil {
		return nil, DecodeHexError(hexStr)
	}
	block, e := block2.NewFromBytes(serializedBlock)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDeserialization,
			Message: "Block decode failed: " + e.Error(),
		}
	}
	// Process this block using the same rules as blocks coming from other nodes. This will in turn relay it to the
	// network like normal.
	_, e = s.Cfg.SyncMgr.SubmitBlock(block, blockchain.BFNone)
	if e != nil {
		return fmt.Sprintf("rejected: %s", e.Error()), nil
	}
	I.F(
		"accepted block %s via submitblock", block.Hash(),
	)
	
	return nil, nil
}

// HandleUnimplemented is the handler for commands that should ultimately be supported but are not yet implemented.
func HandleUnimplemented(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	return nil, ErrRPCUnimplemented
}

// HandleUptime implements the uptime command.
func HandleUptime(s *Server, cmd interface{}, closeChan qu.C) (
	interface{}, error,
) {
	return time.Now().Unix() - s.Cfg.StartupTime, nil
}

// HandleValidateAddress implements the validateaddress command.
func HandleValidateAddress(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.ValidateAddressCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("validateaddress")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	result := btcjson.ValidateAddressChainResult{}
	addr, e := btcaddr.Decode(c.Address, s.Cfg.ChainParams)
	if e != nil {
		// Return the default value (false) for IsValid.
		return result, nil
	}
	result.Address = addr.EncodeAddress()
	result.IsValid = true
	return result, nil
}

// HandleVerifyChain implements the verifychain command.
func HandleVerifyChain(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.VerifyChainCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("verifychain")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	var checkLevel, checkDepth int32
	if c.CheckLevel != nil {
		checkLevel = *c.CheckLevel
	}
	if c.CheckDepth != nil {
		checkDepth = *c.CheckDepth
	}
	e = VerifyChain(s, checkLevel, checkDepth)
	return e == nil, nil
}

// HandleResetChain deletes the existing chain database and restarts
func HandleResetChain(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	// dbName := blockdb.NamePrefix + "_" + *s.Config.DbType
	// if *s.Config.DbType == "sqlite" {
	// 	dbName += ".db"
	// }
	// dbPath := filepath.Join(filepath.Join(*s.Config.DataDir, s.Cfg.ChainParams.Name), dbName)
	// select {
	// case s.RequestProcessShutdown <- struct{}{}:
	// default:
	// }
	// defer interrupt.RequestRestart()
	// os.RemoveAll(dbPath)
	return nil, nil
}

// HandleVerifyMessage implements the verifymessage command.
func HandleVerifyMessage(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	var msg string
	var e error
	// c, ok := cmd.(*btcjson.GetRawTransactionCmd)
	c, ok := cmd.(*btcjson.VerifyMessageCmd)
	if !ok {
		var h string
		h, e = s.HelpCacher.RPCMethodHelp("verifymessage")
		if e != nil {
			msg = e.Error() + "\n\n"
		}
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// Decode the provided address.
	params := s.Cfg.ChainParams
	addr, e := btcaddr.Decode(c.Address, params)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Invalid address or key: " + e.Error(),
		}
	}
	// Only P2PKH addresses are valid for signing.
	if _, ok := addr.(*btcaddr.PubKeyHash); !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCType,
			Message: "Address is not a pay-to-pubkey-hash address",
		}
	}
	// Decode base64 signature.
	sig, e := base64.StdEncoding.DecodeString(c.Signature)
	if e != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCParse.Code,
			Message: "Malformed base64 encoding: " + e.Error(),
		}
	}
	// Validate the signature - this just shows that it was valid at all. we will compare it with the key next.
	var buf bytes.Buffer
	e = wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	if e != nil {
		D.Ln(e)
		
	}
	e = wire.WriteVarString(&buf, 0, c.Message)
	if e != nil {
		D.Ln(e)
		
	}
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	pk, wasCompressed, e := ecc.RecoverCompact(
		ecc.S256(), sig,
		expectedMessageHash,
	)
	if e != nil {
		// Mirror Bitcoin Core behavior, which treats error in RecoverCompact as invalid signature.
		return false, nil
	}
	// Reconstruct the pubkey hash.
	var serializedPK []byte
	if wasCompressed {
		serializedPK = pk.SerializeCompressed()
	} else {
		serializedPK = pk.SerializeUncompressed()
	}
	address, e := btcaddr.NewPubKey(serializedPK, params)
	if e != nil {
		// Again mirror Bitcoin Core behavior, which treats error in public key reconstruction as invalid signature.
		return false, nil
	}
	// Return boolean if addresses match.
	return address.EncodeAddress() == c.Address, nil
}

// HandleVersion implements the version command. NOTE: This is a btcsuite extension ported from github.com/decred/dcrd.
func HandleVersion(s *Server, cmd interface{}, closeChan qu.C) (interface{}, error) {
	result := map[string]btcjson.VersionResult{
		"podjsonrpcapi": {
			VersionString: JSONRPCSemverString,
			Major:         JSONRPCSemverMajor,
			Minor:         JSONRPCSemverMinor,
			Patch:         JSONRPCSemverPatch,
		},
	}
	return result, nil
}
