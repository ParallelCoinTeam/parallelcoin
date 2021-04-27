package wallet

import (
	"errors"
	
	"github.com/p9c/matrjoska/pkg/btcjson"
)

// TODO(jrick): There are several error paths which 'replace' various errors
//  with a more appropiate error from the json package.  Create a map of
//  these replacements so they can be handled once after an RPC handler has
//  returned and before the error is marshaled.
//
// BTCJSONError types to simplify the reporting of specific categories of
// errors, and their *json.RPCError creation.
type (
	// DeserializationError describes a failed deserializaion due to bad user input. It corresponds to
	// json.ErrRPCDeserialization.
	DeserializationError struct {
		error
	}
	// InvalidParameterError describes an invalid parameter passed by the user. It corresponds to
	// json.ErrRPCInvalidParameter.
	InvalidParameterError struct {
		error
	}
	// ParseError describes a failed parse due to bad user input. It corresponds to json.ErrRPCParse.
	ParseError struct {
		error
	}
)

// Errors variables that are defined once here to avoid duplication below.
var (
	ErrNeedPositiveAmount = InvalidParameterError{
		errors.New("amount must be positive"),
	}
	ErrNeedPositiveMinconf = InvalidParameterError{
		errors.New("minconf must be positive"),
	}
	ErrAddressNotInWallet = btcjson.RPCError{
		Code:    btcjson.ErrRPCWallet,
		Message: "address not found in wallet",
	}
	ErrAccountNameNotFound = btcjson.RPCError{
		Code:    btcjson.ErrRPCWalletInvalidAccountName,
		Message: "account name not found",
	}
	ErrUnloadedWallet = btcjson.RPCError{
		Code:    btcjson.ErrRPCWallet,
		Message: "Request requires a wallet but wallet has not loaded yet",
	}
	ErrWalletUnlockNeeded = btcjson.RPCError{
		Code:    btcjson.ErrRPCWalletUnlockNeeded,
		Message: "Enter the wallet passphrase with walletpassphrase first",
	}
	ErrNotImportedAccount = btcjson.RPCError{
		Code:    btcjson.ErrRPCWallet,
		Message: "imported addresses must belong to the imported account",
	}
	ErrNoTransactionInfo = btcjson.RPCError{
		Code:    btcjson.ErrRPCNoTxInfo,
		Message: "No information for transaction",
	}
	ErrReservedAccountName = btcjson.RPCError{
		Code:    btcjson.ErrRPCInvalidParameter,
		Message: "Account name is reserved by RPC server",
	}
)
