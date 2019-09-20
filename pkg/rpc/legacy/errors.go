package legacy

import (
	"errors"

	"git.parallelcoin.io/dev/pod/pkg/rpc/json"
)

// TODO(jrick): There are several error paths which 'replace' various errors
// with a more appropiate error from the json package.  Create a map of
// these replacements so they can be handled once after an RPC handler has
// returned and before the error is marshaled.
// Error types to simplify the reporting of specific categories of
// errors, and their *json.RPCError creation.
type (
	// DeserializationError describes a failed deserializaion due to bad
	// user input.  It corresponds to json.ErrRPCDeserialization.
	DeserializationError struct {
		error
	}
	// InvalidParameterError describes an invalid parameter passed by
	// the user.  It corresponds to json.ErrRPCInvalidParameter.
	InvalidParameterError struct {
		error
	}
	// ParseError describes a failed parse due to bad user input.  It
	// corresponds to json.ErrRPCParse.
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
	ErrAddressNotInWallet = json.RPCError{
		Code:    json.ErrRPCWallet,
		Message: "address not found in wallet",
	}
	ErrAccountNameNotFound = json.RPCError{
		Code:    json.ErrRPCWalletInvalidAccountName,
		Message: "account name not found",
	}
	ErrUnloadedWallet = json.RPCError{
		Code:    json.ErrRPCWallet,
		Message: "Request requires a wallet but wallet has not loaded yet",
	}
	ErrWalletUnlockNeeded = json.RPCError{
		Code:    json.ErrRPCWalletUnlockNeeded,
		Message: "Enter the wallet passphrase with walletpassphrase first",
	}
	ErrNotImportedAccount = json.RPCError{
		Code:    json.ErrRPCWallet,
		Message: "imported addresses must belong to the imported account",
	}
	ErrNoTransactionInfo = json.RPCError{
		Code:    json.ErrRPCNoTxInfo,
		Message: "No information for transaction",
	}
	ErrReservedAccountName = json.RPCError{
		Code:    json.ErrRPCInvalidParameter,
		Message: "Account name is reserved by RPC server",
	}
)
