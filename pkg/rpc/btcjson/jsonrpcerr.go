package btcjson

// General application defined JSON errors.
const (
	ErrRPCMisc                RPCErrorCode = -1
	ErrRPCForbiddenBySafeMode RPCErrorCode = -2
	ErrRPCType                RPCErrorCode = -3
	ErrRPCInvalidAddressOrKey RPCErrorCode = -5
	ErrRPCOutOfMemory         RPCErrorCode = -7
	ErrRPCInvalidParameter    RPCErrorCode = -8
	ErrRPCDatabase            RPCErrorCode = -20
	ErrRPCDeserialization     RPCErrorCode = -22
	ErrRPCVerify              RPCErrorCode = -25
	// Peer-to-peer client errors.
	ErrRPCClientNotConnected      RPCErrorCode = -9
	ErrRPCClientInInitialDownload RPCErrorCode = -10
	ErrRPCClientNodeNotAdded      RPCErrorCode = -24
	// Wallet JSON errors
	ErrRPCWallet                    RPCErrorCode = -4
	ErrRPCWalletInsufficientFunds   RPCErrorCode = -6
	ErrRPCWalletInvalidAccountName  RPCErrorCode = -11
	ErrRPCWalletKeypoolRanOut       RPCErrorCode = -12
	ErrRPCWalletUnlockNeeded        RPCErrorCode = -13
	ErrRPCWalletPassphraseIncorrect RPCErrorCode = -14
	ErrRPCWalletWrongEncState       RPCErrorCode = -15
	ErrRPCWalletEncryptionFailed    RPCErrorCode = -16
	ErrRPCWalletAlreadyUnlocked     RPCErrorCode = -17
	
	// Specific Errors related to commands. These are the ones a user of the RPC server are most likely to see.
	// Generally, the codes should match one of the more general errors above.
	
	ErrRPCBlockNotFound     RPCErrorCode = -5
	ErrRPCBlockCount        RPCErrorCode = -5
	ErrRPCBestBlockHash     RPCErrorCode = -5
	ErrRPCDifficulty        RPCErrorCode = -5
	ErrRPCOutOfRange        RPCErrorCode = -1
	ErrRPCNoTxInfo          RPCErrorCode = -5
	ErrRPCNoCFIndex         RPCErrorCode = -5
	ErrRPCNoNewestBlockInfo RPCErrorCode = -5
	ErrRPCInvalidTxVout     RPCErrorCode = -5
	ErrRPCRawTxString       RPCErrorCode = -32602
	ErrRPCDecodeHexString   RPCErrorCode = -22
	// Errors that are specific to pod.
	ErrRPCNoWallet      RPCErrorCode = -1
	ErrRPCNoChain       RPCErrorCode = -1
	ErrRPCUnimplemented RPCErrorCode = -1
)

// Standard JSON-RPC 2.0 errors.
var (
	ErrRPCInternal = &RPCError{
		Code:    -32603,
		Message: "Internal error",
	}
	// Standard JSON-RPC 2.0 errors.
	ErrRPCInvalidParams = &RPCError{
		Code:    -32602,
		Message: "Invalid parameters",
	}
	// Standard JSON-RPC 2.0 errors.
	ErrRPCInvalidRequest = &RPCError{
		Code:    -32600,
		Message: "Invalid request",
	}
	// Standard JSON-RPC 2.0 errors.
	ErrRPCMethodNotFound = &RPCError{
		Code:    -32601,
		Message: "Method not found",
	}
	// Standard JSON-RPC 2.0 errors.
	ErrRPCParse = &RPCError{
		Code:    -32700,
		Message: "Parse error",
	}
)
