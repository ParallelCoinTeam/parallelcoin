package rpc

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/conformal/fastsha256"

	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/util"
)

// Uint256Size is the number of bytes needed to represent an unsigned 256-bit
// integer.
const Uint256Size = 32

// GetworkDataLen is the length of the data field of the getwork RPC.
// It consists of the serialized block header plus the internal sha256
// padding. The internal sha256 padding consists of a single 1 bit followed
// by enough zeros to pad the message out to 56 bytes followed by length of
// the message in bits encoded as a big-endian uint64 (8 bytes).  Thus,
// the resulting length is a multiple of the sha256 block size (64 bytes).
var GetworkDataLen = (1 + ((wire.MaxBlockHeaderPayload + 8) / fastsha256.
	BlockSize)) * fastsha256.BlockSize

// Hash1Len is the length of the hash1 field of the getwork RPC.
// It consists of a zero hash plus the internal sha256 padding.
// See the getworkDataLen comment for details about the internal sha256
// padding format.
var Hash1Len = (1 + ((chainhash.HashSize + 8) / fastsha256.
	BlockSize)) * fastsha256.BlockSize

// BigToLEUint256 returns the passed big integer as an unsigned 256-bit
// integer encoded as little-endian bytes.
// Numbers which are larger than the max unsigned 256-bit integer are truncated.
func BigToLEUint256(n *big.Int) [Uint256Size]byte {
	// Pad or truncate the big-endian big int to correct number of bytes.
	nBytes := n.Bytes()
	nlen := len(nBytes)
	pad := 0
	start := 0
	if nlen <= Uint256Size {
		pad = Uint256Size - nlen
	} else {
		start = nlen - Uint256Size
	}
	var buf [Uint256Size]byte
	copy(buf[pad:], nBytes[start:])
	// Reverse the bytes to little endian and return them.
	for i := 0; i < Uint256Size/2; i++ {
		buf[i], buf[Uint256Size-1-i] = buf[Uint256Size-1-i], buf[i]
	}
	return buf
}

// HandleGetWork handles the getwork call
func HandleGetWork(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*btcjson.GetWorkCmd)
	if len(s.StateCfg.ActiveMiningAddrs) == 0 {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInternal.Code,
			Message: "No payment addresses specified via --miningaddr",
		}
	}
	netwk := (*s.Config.Network)[0]
	if !((netwk == 'r') || (netwk == 's')) &&
		s.Cfg.ConnMgr.ConnectedCount() == 0 {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCClientNotConnected,
			Message: "Pod is not connected to network",
		}
	}
	// No point in generating or accepting work before the chain is synced.
	latestHeight := s.Cfg.Chain.BestSnapshot().Height
	if latestHeight != 0 && !s.Cfg.SyncMgr.IsCurrent() {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCClientInInitialDownload,
			Message: "Pod is not yet synchronised...",
		}
	}
	state := s.GBTWorkState
	state.Lock()
	defer state.Unlock()
	if c.Data != nil {
		return HandleGetWorkSubmission(s, *c.Data)
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	payToAddr := s.StateCfg.ActiveMiningAddrs[rand.Intn(len(s.StateCfg.ActiveMiningAddrs))]
	lastTxUpdate := s.GBTWorkState.LastTxUpdate
	latestHash := &s.Cfg.Chain.BestSnapshot().Hash
	generator := s.Cfg.Generator
	if state.Template == nil {
		var err error
		state.Template, err = generator.NewBlockTemplate(0, payToAddr,
			s.Cfg.Algo)
		if err != nil {
			log.ERROR(err)
			return nil, err
		}
	}
	msgBlock := state.Template.Block
	if msgBlock == nil || state.prevHash == nil ||
		!state.prevHash.IsEqual(latestHash) ||
		(state.LastTxUpdate != lastTxUpdate &&
			time.Now().After(state.LastGenerated.Add(time.Minute))) {
		//	Reset the extra nonce and clear all cached template variations if
		//	the best block changed.
		if state.prevHash != nil && !state.prevHash.IsEqual(latestHash) {
			e := state.UpdateBlockTemplate(s, false)
			if e != nil {
				log.WARN("failed to update block template", e)
			}
		}
		//	Reset the previous best hash the block template was generated
		//	against so any errors below cause the next invocation to try
		//	again.
		state.prevHash = nil
		var err error
		state.Template, err = generator.NewBlockTemplate(0, payToAddr,
			s.Cfg.Algo)
		if err != nil {
			errStr := fmt.Sprintf("Failed to create new block template: %v", err)
			log.ERROR(errStr)
			return nil, &btcjson.RPCError{
				Code:    btcjson.ErrRPCInternal.Code,
				Message: errStr,
			}
		}
		msgBlock = state.Template.Block
		// Update work state to ensure another block template isn't generated
		// until needed.
		state.Template.Block = msgBlock
		state.LastGenerated = time.Now()
		state.LastTxUpdate = lastTxUpdate
		state.prevHash = latestHash
		log.DEBUGC(func() string {
			return fmt.Sprintf(
				"generated block template (timestamp %v, target %064x, "+
					"merkle root %s, signature script %x)",
				msgBlock.Header.Timestamp, fork.CompactToBig(msgBlock.Header.
					Bits), msgBlock.Header.MerkleRoot,
				msgBlock.Transactions[0].TxIn[0].SignatureScript,
			)
		})
	} else {
		//	At this point, there is a saved block template and a new request for
		// work was made but either the available transactions haven't change
		// or it hasn't been long enough to trigger a new block template to be
		// generated.
		// So, update the existing block template and track the variations so
		// each variation can be regenerated if a caller finds an answer and
		// makes a submission against it. Update the time of the block template
		// to the current time while accounting for the median time of the past
		// several blocks per the chain consensus rules.
		e := generator.UpdateBlockTime(0, msgBlock)
		if e != nil {
			log.WARN("failed to update block time", e)
		}
		// Increment the extra nonce and update the block template with the new
		// value by regenerating the coinbase script and setting the merkle root
		// to the new value.
		log.DEBUGF(
			"updated block template (timestamp %v, target %064x, "+
				"merkle root %s, signature script %x)",
			msgBlock.Header.Timestamp, fork.CompactToBig(msgBlock.Header.Bits),
			msgBlock.Header.MerkleRoot, msgBlock.Transactions[0].TxIn[0].
				SignatureScript)
	}
	//	In order to efficiently store the variations of block templates that
	//	have been provided to callers save a pointer to the block as well as the
	// modified signature script keyed by the merkle root.  This information,
	// along with the data that is included in a work submission, is used to
	// rebuild the block before checking the submitted solution.
	/*
		coinbaseTx := msgBlock.Transactions[0]
		state.blockInfo[msgBlock.Header.MerkleRoot] = &workStateBlockInfo{
			msgBlock:        msgBlock,
			signatureScript: coinbaseTx.TxIn[0].SignatureScript,
		}
	*/
	// Serialize the block header into a buffer large enough to hold the the
	// block header and the internal sha256 padding that is added and returned
	// as part of the data below.
	data := make([]byte, 0, GetworkDataLen)
	buf := bytes.NewBuffer(data)
	err := msgBlock.Header.Serialize(buf)
	if err != nil {
		log.ERROR(err)
		errStr := fmt.Sprintf("Failed to serialize data: %v", err)
		log.WARN(errStr)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInternal.Code,
			Message: errStr,
		}
	}
	// Calculate the midstate for the block header.  The midstate here is the
	// internal state of the sha256 algorithm for the first chunk of the block
	// header (sha256 operates on 64-byte chunks) which is before the nonce.
	// This allows sophisticated callers to avoid hashing the first chunk over
	// and over while iterating the nonce range.
	data = data[:buf.Len()]
	midstate := fastsha256.MidState256(data)
	// Expand the data slice to include the full data buffer and apply the
	// internal sha256 padding which consists of a single 1 bit followed by
	// enough zeros to pad the message out to 56 bytes followed by the length of
	// the message in bits encoded as a big-endian uint64 (8 bytes).  Thus, the
	// resulting length is a multiple of the sha256 block size (64 bytes).  This
	// makes the data ready for sophisticated caller to make use of only the
	// second chunk along with the midstate for the first chunk.
	data = data[:GetworkDataLen]
	data[wire.MaxBlockHeaderPayload] = 0x80
	binary.BigEndian.PutUint64(data[len(data)-8:],
		wire.MaxBlockHeaderPayload*8)
	//	Create the hash1 field which is a zero hash along with the internal
	// sha256 padding as described above.  This field is really quite useless,
	// but it is required for compatibility with the reference implementation.
	var hash1 = make([]byte, Hash1Len)
	hash1[chainhash.HashSize] = 0x80
	binary.BigEndian.PutUint64(hash1[len(hash1)-8:], chainhash.HashSize*8)
	// The final result reverses the each of the fields to little endian. In
	// particular, the data, hash1, and midstate fields are treated as arrays of
	// uint32s (per the internal sha256 hashing state) which are in big endian,
	// and thus each 4 bytes is byte swapped.  The target is also in big endian,
	// but it is treated as a uint256 and byte swapped to little endian
	// accordingly. The fact the fields are reversed in this way is rather odd
	// and likely an artifact of some legacy internal state in the reference
	// implementation, but it is required for compatibility.
	ReverseUint32Array(data)
	ReverseUint32Array(hash1)
	ReverseUint32Array(midstate[:])
	target := BigToLEUint256(fork.CompactToBig(msgBlock.Header.Bits))
	reply := &btcjson.GetWorkResult{
		Data:     hex.EncodeToString(data),
		Hash1:    hex.EncodeToString(hash1),
		Midstate: hex.EncodeToString(midstate[:]),
		Target:   hex.EncodeToString(target[:]),
	}
	return reply, nil
}

//	HandleGetWorkSubmission is a helper for handleGetWork which deals with the
// calling submitting work to be verified and processed. This function MUST be
// called with the RPC workstate locked.
func HandleGetWorkSubmission(s *Server, hexData string) (interface{}, error) {
	// Ensure the provided data is sane.
	if len(hexData)%2 != 0 {
		hexData = "0" + hexData
	}
	data, err := hex.DecodeString(hexData)
	if err != nil {
		log.ERROR(err)
		return nil, &btcjson.RPCError{
			Code: btcjson.ErrRPCInvalidParameter,
			Message: fmt.Sprintf("argument must be "+
				"hexadecimal string (not %q)", hexData),
		}
	}
	if len(data) != GetworkDataLen {
		return false, &btcjson.RPCError{
			Code: btcjson.ErrRPCInvalidParameter,
			Message: fmt.Sprintf("argument must be "+
				"%d bytes (not %d)", GetworkDataLen,
				len(data)),
		}
	}
	// Reverse the data as if it were an array of 32-bit unsigned integers. The
	// fact the getwork request and submission data is reversed in this way is
	// rather odd and likey an artifact of some legacy internal state in the
	// reference implementation, but it is required for compatibility.
	ReverseUint32Array(data)
	// Deserialize the block header from the data.
	var submittedHeader wire.BlockHeader
	bhBuf := bytes.NewBuffer(data[0:wire.MaxBlockHeaderPayload])
	err = submittedHeader.Deserialize(bhBuf)
	if err != nil {
		log.ERROR(err)
		return false, &btcjson.RPCError{
			Code: btcjson.ErrRPCInvalidParameter,
			Message: fmt.Sprintf("argument does not "+
				"contain a valid block header: %v", err),
		}
	}
	// Look up the full block for the provided data based on the merkle root.
	// Return false to indicate the solve failed if it's not available.
	state := s.GBTWorkState

	if state.Template.Block.Header.MerkleRoot.String() == "" {
		log.DEBUG(
			"Block submitted via getwork has no matching template for merkle root",
			submittedHeader.MerkleRoot)
		return false, nil
	}
	// Reconstruct the block using the submitted header stored block info.
	msgBlock := state.Template.Block
	block := util.NewBlock(msgBlock)
	msgBlock.Header.Timestamp = submittedHeader.Timestamp
	msgBlock.Header.Nonce = submittedHeader.Nonce
	msgBlock.Transactions[0].TxIn[0].SignatureScript = state.Template.Block.
		Transactions[0].TxIn[0].SignatureScript
	merkles := blockchain.BuildMerkleTreeStore(block.Transactions(), false)
	msgBlock.Header.MerkleRoot = *merkles[len(merkles)-1]
	// Ensure the submitted block hash is less than the target difficulty.
	pl := fork.GetMinDiff(s.Cfg.Algo, s.Cfg.Chain.BestSnapshot().Height)
	err = blockchain.CheckProofOfWork(block, pl, s.Cfg.Chain.BestSnapshot().Height)
	if err != nil {
		log.ERROR(err)
		// Anything other than a rule violation is an unexpected error, so return
		// that error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInternal.Code,
				Message: fmt.Sprintf("Unexpected error while checking proof"+
					" of work: %v", err),
			}
		}
		log.DEBUG("block submitted via getwork does not meet the required proof of"+
			" work:", err)
		return false, nil
	}
	latestHash := &s.Cfg.Chain.BestSnapshot().Hash
	if !msgBlock.Header.PrevBlock.IsEqual(latestHash) {
		log.DEBUGF(
			"block submitted via getwork with previous block %s is stale",
			msgBlock.Header.PrevBlock)
		return false, nil
	}
	// Process this block using the same rules as blocks coming from other
	// nodes.  This will in turn relay it to the network like normal.
	_, isOrphan, err := s.Cfg.Chain.ProcessBlock(0, block, 0,
		s.Cfg.Chain.BestSnapshot().Height)
	if err != nil || isOrphan {
		log.ERROR(err)
		// Anything other than a rule violation is an unexpected error, so return
		// that error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			return nil, &btcjson.RPCError{
				Code: btcjson.ErrRPCInternal.Code,
				Message: fmt.Sprintf("Unexpected error while processing block"+
					": %v", err),
			}
		}
		log.INFO("block submitted via getwork rejected:", err)
		return false, nil
	}
	// The block was accepted.
	blockSha := block.Hash()
	log.INFO("block submitted via getwork accepted:", blockSha)
	return true, nil
}

// ReverseUint32Array treats the passed bytes as a series of uint32s and
// reverses the byte order of each uint32.  The passed byte slice must be a
// multiple of 4 for a correct result.  The passed bytes slice is modified.
func ReverseUint32Array(b []byte) {
	blen := len(b)
	for i := 0; i < blen; i += 4 {
		b[i], b[i+3] = b[i+3], b[i]
		b[i+1], b[i+2] = b[i+2], b[i+1]
	}
}
