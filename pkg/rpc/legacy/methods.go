package legacy

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	js "encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/p9c/pod/pkg/logi"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txrules "github.com/p9c/pod/pkg/chain/tx/rules"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
	ec "github.com/p9c/pod/pkg/util/elliptic"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	"github.com/p9c/pod/pkg/wallet/chain"
)

// // confirmed checks whether a transaction at height txHeight has met minconf
// // confirmations for a blockchain at height curHeight.
// func confirmed(// 	minconf, txHeight, curHeight int32) bool {
// 	return confirms(txHeight, curHeight) >= minconf
// }

// Confirms returns the number of confirmations for a transaction in a block at
// height txHeight (or -1 for an unconfirmed tx) given the chain height
// curHeight.
func Confirms(txHeight, curHeight int32) int32 {
	switch {
	case txHeight == -1, txHeight > curHeight:
		return 0
	default:
		return curHeight - txHeight + 1
	}
}

// var RPCHandlers = map[string]struct {
// 	Handler          RequestHandler
// 	HandlerWithChain RequestHandlerChainRequired
// 	// Function variables cannot be compared against anything but nil, so
// 	// use a boolean to record whether help generation is necessary.  This
// 	// is used by the tests to ensure that help can be generated for every
// 	// implemented method.
// 	//
// 	// A single map and this bool is here is used rather than several maps
// 	// for the unimplemented handlers so every method has exactly one
// 	// handler function.
// 	//
// 	// The Return field returns a new channel of the type returned by this function. This makes it possible to
// 	// use this for callers to receive a response in the `cpc` library which implements the functions as channel pipes
// 	NoHelp bool
// 	Params interface{}
// 	Return func() interface{}
// }{
// 	// Reference implementation wallet methods (implemented)
// 	"addmultisigaddress": {
// 		Handler: AddMultiSigAddress,
// 		Params:  make(chan btcjson.AddMultisigAddressCmd),
// 		Return:  func() interface{} { return make(chan AddMultiSigAddressRes) },
// 	},
// 	"createmultisig": {
// 		Handler: CreateMultiSig,
// 		Params:  make(chan btcjson.CreateMultisigCmd),
// 		Return:  func() interface{} { return make(chan CreateMultiSigRes) },
// 	},
// 	"dumpprivkey": {
// 		Handler: DumpPrivKey,
// 		Params:  make(chan btcjson.DumpPrivKeyCmd),
// 		Return:  func() interface{} { return make(chan DumpPrivKeyRes) },
// 	},
// 	"getaccount": {
// 		Handler: GetAccount,
// 		Params:  make(chan btcjson.GetAccountCmd),
// 		Return:  func() interface{} { return make(chan GetAccountRes) },
// 	},
// 	"getaccountaddress": {
// 		Handler: GetAccountAddress,
// 		Params:  make(chan btcjson.GetAccountAddressCmd),
// 		Return:  func() interface{} { return make(chan GetAccountAddressRes) },
// 	},
// 	"getaddressesbyaccount": {
// 		Handler: GetAddressesByAccount,
// 		Params:  make(chan btcjson.GetAddressesByAccountCmd),
// 		Return:  func() interface{} { return make(chan GetAddressesByAccountRes) },
// 	},
// 	"getbalance": {
// 		Handler: GetBalance,
// 		Params:  make(chan btcjson.GetBalanceCmd),
// 		Return:  func() interface{} { return make(chan GetBalanceRes) },
// 	},
// 	"getbestblockhash": {
// 		Handler: GetBestBlockHash,
// 		Return:  func() interface{} { return make(chan GetBestBlockHashRes) },
// 	},
// 	"getblockcount": {
// 		Handler: GetBlockCount,
// 		Return:  func() interface{} { return make(chan GetBlockCountRes) },
// 	},
// 	"getinfo": {
// 		HandlerWithChain: GetInfo,
// 		Return:           func() interface{} { return make(chan GetInfoRes) },
// 	},
// 	"getnewaddress": {
// 		Handler: GetNewAddress,
// 		Params:  make(chan btcjson.GetNewAddressCmd),
// 		Return:  func() interface{} { return make(chan GetNewAddressRes) },
// 	},
// 	"getrawchangeaddress": {
// 		Handler: GetRawChangeAddress,
// 		Params:  make(chan btcjson.GetRawChangeAddressCmd),
// 		Return:  func() interface{} { return make(chan GetRawChangeAddressRes) },
// 	},
// 	"getreceivedbyaccount": {
// 		Handler: GetReceivedByAccount,
// 		Params:  make(chan btcjson.GetReceivedByAccountCmd),
// 		Return:  func() interface{} { return make(chan GetReceivedByAccountRes) },
// 	},
// 	"getreceivedbyaddress": {
// 		Handler: GetReceivedByAddress,
// 		Params:  make(chan btcjson.GetReceivedByAddressCmd),
// 		Return:  func() interface{} { return make(chan GetReceivedByAddressRes) },
// 	},
// 	"gettransaction": {
// 		Handler: GetTransaction,
// 		Params:  make(chan btcjson.GetTransactionCmd),
// 		Return:  func() interface{} { return make(chan GetTransactionRes) },
// 	},
// 	"help": {
// 		Handler:          HelpNoChainRPC,
// 		HandlerWithChain: HelpWithChainRPC,
// 		Params:           make(chan btcjson.HelpCmd),
// 		Return:           func() interface{} { return make(chan HelpNoChainRPCRes) },
// 	},
// 	"importprivkey": {
// 		Handler: ImportPrivKey,
// 		Params:  make(chan btcjson.ImportPrivKeyCmd),
// 		Return:  func() interface{} { return make(chan ImportPrivKeyRes) },
// 	},
// 	"keypoolrefill": {
// 		Handler: KeypoolRefill,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan KeypoolRefillRes) },
// 	},
// 	"listaccounts": {
// 		Handler: ListAccounts,
// 		Params:  make(chan btcjson.ListAccountsCmd),
// 		Return:  func() interface{} { return make(chan ListAccountsRes) },
// 	},
// 	"listlockunspent": {
// 		Handler: ListLockUnspent,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan ListLockUnspentRes) },
// 	},
// 	"listreceivedbyaccount": {
// 		Handler: ListReceivedByAccount,
// 		Params:  make(chan btcjson.ListReceivedByAccountCmd),
// 		Return:  func() interface{} { return make(chan ListReceivedByAccountRes) },
// 	},
// 	"listreceivedbyaddress": {
// 		Handler: ListReceivedByAddress,
// 		Params:  make(chan btcjson.ListReceivedByAddressCmd),
// 		Return:  func() interface{} { return make(chan ListReceivedByAddressRes) },
// 	},
// 	"listsinceblock": {
// 		HandlerWithChain: ListSinceBlock,
// 		Params:           make(chan btcjson.ListSinceBlockCmd),
// 		Return:           func() interface{} { return make(chan ListSinceBlockRes) },
// 	},
// 	"listtransactions": {
// 		Handler: ListTransactions,
// 		Params:  make(chan btcjson.ListTransactionsCmd),
// 		Return:  func() interface{} { return make(chan ListTransactionsRes) },
// 	},
// 	"listunspent": {
// 		Handler: ListUnspent,
// 		Params:  make(chan btcjson.ListUnspentCmd),
// 		Return:  func() interface{} { return make(chan ListUnspentRes) },
// 	},
// 	"lockunspent": {
// 		Handler: LockUnspent,
// 		Params:  make(chan btcjson.LockUnspentCmd),
// 		Return:  func() interface{} { return make(chan LockUnspentRes) },
// 	},
// 	"sendfrom": {
// 		HandlerWithChain: SendFrom,
// 		Params:           make(chan btcjson.SendFromCmd),
// 		Return:           func() interface{} { return make(chan SendFromRes) },
// 	},
// 	"sendmany": {
// 		Handler: SendMany,
// 		Params:  make(chan btcjson.SendManyCmd),
// 		Return:  func() interface{} { return make(chan SendManyRes) },
// 	},
// 	"sendtoaddress": {
// 		Handler: SendToAddress,
// 		Params:  make(chan btcjson.SendToAddressCmd),
// 		Return:  func() interface{} { return make(chan SendToAddressRes) },
// 	},
// 	"settxfee": {
// 		Handler: SetTxFee,
// 		Params:  make(chan btcjson.SetTxFeeCmd),
// 		Return:  func() interface{} { return make(chan SetTxFeeRes) },
// 	},
// 	"signmessage": {
// 		Handler: SignMessage,
// 		Params:  make(chan btcjson.SignMessageCmd),
// 		Return:  func() interface{} { return make(chan SignMessageRes) },
// 	},
// 	"signrawtransaction": {
// 		HandlerWithChain: SignRawTransaction,
// 		Params:           make(chan btcjson.SignRawTransactionCmd),
// 		Return:           func() interface{} { return make(chan SignRawTransactionRes) },
// 	},
// 	"validateaddress": {
// 		Handler: ValidateAddress,
// 		Params:  make(chan btcjson.ValidateAddressCmd),
// 		Return:  func() interface{} { return make(chan ValidateAddressRes) },
// 	},
// 	"verifymessage": {
// 		Handler: VerifyMessage,
// 		Params:  make(chan btcjson.VerifyMessageCmd),
// 		Return:  func() interface{} { return make(chan VerifyMessageRes) },
// 	},
// 	"walletlock": {
// 		Handler: WalletLock,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan WalletLockRes) },
// 	},
// 	"walletpassphrase": {
// 		Handler: WalletPassphrase,
// 		Params:  make(chan btcjson.WalletPassphraseCmd),
// 		Return:  func() interface{} { return make(chan WalletPassphraseRes) },
// 	},
// 	"walletpassphrasechange": {
// 		Handler: WalletPassphraseChange,
// 		Params:  make(chan btcjson.WalletPassphraseChangeCmd),
// 		Return:  func() interface{} { return make(chan WalletPassphraseChangeRes) },
// 	},
// 	// Reference implementation methods (still unimplemented)
// 	"backupwallet":         {Handler: Unimplemented, NoHelp: true},
// 	"dumpwallet":           {Handler: Unimplemented, NoHelp: true},
// 	"getwalletinfo":        {Handler: Unimplemented, NoHelp: true},
// 	"importwallet":         {Handler: Unimplemented, NoHelp: true},
// 	"listaddressgroupings": {Handler: Unimplemented, NoHelp: true},
// 	// Reference methods which can't be implemented by btcwallet due to
// 	// design decision differences
// 	"encryptwallet": {Handler: Unsupported, NoHelp: true},
// 	"move":          {Handler: Unsupported, NoHelp: true},
// 	"setaccount":    {Handler: Unsupported, NoHelp: true},
// 	// Extensions to the reference client JSON-RPC API
// 	"createnewaccount": {
// 		Handler: CreateNewAccount,
// 		Params:  make(chan btcjson.CreateNewAccountCmd),
// 		Return:  func() interface{} { return make(chan CreateNewAccountRes) },
// 	},
// 	"getbestblock": {
// 		Handler: GetBestBlock,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan GetBestBlockRes) },
// 	},
// 	// This was an extension but the reference implementation added it as
// 	// well, but with a different API (no account parameter).  It's listed
// 	// here because it hasn't been update to use the reference
// 	// implemenation's API.
// 	"getunconfirmedbalance": {
// 		Handler: GetUnconfirmedBalance,
// 		Params:  make(chan btcjson.GetUnconfirmedBalanceCmd),
// 		Return:  func() interface{} { return make(chan GetUnconfirmedBalanceRes) },
// 	},
// 	"listaddresstransactions": {
// 		Handler: ListAddressTransactions,
// 		Params:  make(chan btcjson.ListAddressTransactionsCmd),
// 		Return:  func() interface{} { return make(chan ListAddressTransactionsRes) },
// 	},
// 	"listalltransactions": {
// 		Handler: ListAllTransactions,
// 		Params:  make(chan btcjson.ListAllTransactionsCmd),
// 		Return:  func() interface{} { return make(chan ListAllTransactionsRes) },
// 	},
// 	"renameaccount": {
// 		Handler: RenameAccount,
// 		Params:  make(chan btcjson.RenameAccountCmd),
// 		Return:  func() interface{} { return make(chan RenameAccountRes) },
// 	},
// 	"walletislocked": {
// 		Handler: WalletIsLocked,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan WalletIsLockedRes) },
// 	},
// 	"dropwallethistory": {
// 		Handler: HandleDropWalletHistory,
// 		Params:  make(chan struct{}),
// 		Return:  func() interface{} { return make(chan DropWalletHistoryRes) },
// 	},
// }

// Unimplemented handles an Unimplemented RPC request with the
// appropiate error.
func Unimplemented(interface{}, *wallet.Wallet) (interface{}, error) {
	return nil, &btcjson.RPCError{
		Code:    btcjson.ErrRPCUnimplemented,
		Message: "Method unimplemented",
	}
}

// Unsupported handles a standard bitcoind RPC request which is
// Unsupported by btcwallet due to design differences.
func Unsupported(interface{}, *wallet.Wallet) (interface{}, error) {
	return nil, &btcjson.RPCError{
		Code:    -1,
		Message: "Request unsupported by wallet",
	}
}

// LazyHandler is a closure over a requestHandler or passthrough request with
// the RPC server's wallet and chain server variables as part of the closure
// context.
type LazyHandler func() (interface{}, *btcjson.RPCError)

// LazyApplyHandler looks up the best request handler func for the method,
// returning a closure that will execute it with the (required) wallet and
// (optional) consensus RPC server.  If no handlers are found and the
// chainClient is not nil, the returned handler performs RPC passthrough.
func LazyApplyHandler(request *btcjson.Request, w *wallet.Wallet, chainClient chain.Interface) LazyHandler {
	handlerData, ok := RPCHandlers[request.Method]
	if ok && handlerData.Handler != nil && w != nil && chainClient != nil {
		return func() (interface{}, *btcjson.RPCError) {
			cmd, err := btcjson.UnmarshalCmd(request)
			if err != nil {
				log.L.Error(err)
				return nil, btcjson.ErrRPCInvalidRequest
			}
			switch client := chainClient.(type) {
			case *chain.RPCClient:
				resp, err := handlerData.Handler(cmd,
					w, client)
				if err != nil {
					log.L.Error(err)
					return nil, JSONError(err)
				}
				return resp, nil
			default:
				return nil, &btcjson.RPCError{
					Code:    -1,
					Message: "Chain RPC is inactive",
				}
			}
		}
	}
	// log.L.Info("handler", handlerData.Handler, "wallet", w)
	if ok && handlerData.Handler != nil && w != nil {
		log.L.Info("handling", request.Method)
		return func() (interface{}, *btcjson.RPCError) {
			cmd, err := btcjson.UnmarshalCmd(request)
			if err != nil {
				log.L.Error(err)
				return nil, btcjson.ErrRPCInvalidRequest
			}
			resp, err := handlerData.Handler(cmd, w)
			if err != nil {
				log.L.Error(err)
				return nil, JSONError(err)
			}
			return resp, nil
		}
	}
	// Fallback to RPC passthrough
	return func() (interface{}, *btcjson.RPCError) {
		log.L.Info("passing to node", request.Method)
		if chainClient == nil {
			return nil, &btcjson.RPCError{
				Code:    -1,
				Message: "Chain RPC is inactive",
			}
		}
		switch client := chainClient.(type) {
		case *chain.RPCClient:
			resp, err := client.RawRequest(request.Method,
				request.Params)
			if err != nil {
				log.L.Error(err)
				return nil, JSONError(err)
			}
			return &resp, nil
		default:
			return nil, &btcjson.RPCError{
				Code:    -1,
				Message: "Chain RPC is inactive",
			}
		}
	}
}

// MakeResponse makes the JSON-RPC response struct for the result and error
// returned by a requestHandler.  The returned response is not ready for
// marshaling and sending off to a client, but must be
func MakeResponse(id, result interface{}, err error) btcjson.Response {
	idPtr := IDPointer(id)
	if err != nil {
		log.L.Error(err)
		return btcjson.Response{
			ID:    idPtr,
			Error: JSONError(err),
		}
	}
	resultBytes, err := js.Marshal(result)
	if err != nil {
		log.L.Error(err)
		return btcjson.Response{
			ID: idPtr,
			Error: &btcjson.RPCError{
				Code:    btcjson.ErrRPCInternal.Code,
				Message: "Unexpected error marshalling result",
			},
		}
	}
	return btcjson.Response{
		ID:     idPtr,
		Result: js.RawMessage(resultBytes),
	}
}

// JSONError creates a JSON-RPC error from the Go error.
func JSONError(err error) *btcjson.RPCError {
	if err == nil {
		return nil
	}
	code := btcjson.ErrRPCWallet
	switch e := err.(type) {
	case btcjson.RPCError:
		return &e
	case *btcjson.RPCError:
		return e
	case DeserializationError:
		code = btcjson.ErrRPCDeserialization
	case InvalidParameterError:
		code = btcjson.ErrRPCInvalidParameter
	case ParseError:
		code = btcjson.ErrRPCParse.Code
	case waddrmgr.ManagerError:
		switch e.ErrorCode {
		case waddrmgr.ErrWrongPassphrase:
			code = btcjson.ErrRPCWalletPassphraseIncorrect
		}
	}
	return &btcjson.RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

// MakeMultiSigScript is a helper function to combine common logic for
// AddMultiSig and CreateMultiSig.
func MakeMultiSigScript(w *wallet.Wallet, keys []string, nRequired int) ([]byte, error) {
	keysesPrecious := make([]*util.AddressPubKey, len(keys))
	// The address list will made up either of addreseses (pubkey hash), for
	// which we need to look up the keys in wallet, straight pubkeys, or a
	// mixture of the two.
	for i, a := range keys {
		// try to parse as pubkey address
		a, err := DecodeAddress(a, w.ChainParams())
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		switch addr := a.(type) {
		case *util.AddressPubKey:
			keysesPrecious[i] = addr
		default:
			pubKey, err := w.PubKeyForAddress(addr)
			if err != nil {
				log.L.Error(err)
				return nil, err
			}
			pubKeyAddr, err := util.NewAddressPubKey(
				pubKey.SerializeCompressed(), w.ChainParams())
			if err != nil {
				log.L.Error(err)
				return nil, err
			}
			keysesPrecious[i] = pubKeyAddr
		}
	}
	return txscript.MultiSigScript(keysesPrecious, nRequired)
}

// AddMultiSigAddress handles an addmultisigaddress request by adding a
// multisig address to the given wallet.
func AddMultiSigAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	var msg string
	cmd, ok := icmd.(*btcjson.AddMultisigAddressCmd)
	// cmd, ok := icmd.(*btcjson.ListTransactionsCmd)
	if !ok {
		var h string
		h = HelpDescsEnUS()["addmultisigaddress"]
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	// If an account is specified, ensure that is the imported account.
	if cmd.Account != nil && *cmd.Account != waddrmgr.ImportedAddrAccountName {
		return nil, &ErrNotImportedAccount
	}
	secp256k1Addrs := make([]util.Address, len(cmd.Keys))
	for i, k := range cmd.Keys {
		addr, err := DecodeAddress(k, w.ChainParams())
		if err != nil {
			log.L.Error(err)
			return nil, ParseError{err}
		}
		secp256k1Addrs[i] = addr
	}
	script, err := w.MakeMultiSigScript(secp256k1Addrs, cmd.NRequired)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	p2shAddr, err := w.ImportP2SHRedeemScript(script)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return p2shAddr.EncodeAddress(), nil
}

// CreateMultiSig handles an createmultisig request by returning a
// multisig address for the given inputs.
func CreateMultiSig(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	var msg string
	cmd, ok := icmd.(*btcjson.CreateMultisigCmd)
	if !ok {
		var h string
		h = HelpDescsEnUS()["createmultisig"]
		msg += h
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	script, err := MakeMultiSigScript(w, cmd.Keys, cmd.NRequired)
	if err != nil {
		log.L.Error(err)
		return nil, ParseError{err}
	}
	address, err := util.NewAddressScriptHash(script, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		// above is a valid script, shouldn't happen.
		return nil, err
	}
	return btcjson.CreateMultiSigResult{
		Address:      address.EncodeAddress(),
		RedeemScript: hex.EncodeToString(script),
	}, nil
}

// DumpPrivKey handles a dumpprivkey request with the private key
// for a single address, or an appropiate error if the wallet
// is locked.
func DumpPrivKey(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	var msg string
	cmd, ok := icmd.(*btcjson.DumpPrivKeyCmd)
	if !ok {
		msg = HelpDescsEnUS()["dumpprivkey"]
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: msg,
			// "invalid subcommand for addnode",
		}
	}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	key, err := w.DumpWIFPrivateKey(addr)
	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {
		// Address was found, but the private key isn't
		// accessible.
		return nil, &ErrWalletUnlockNeeded
	}
	return key, err
}

// // dumpWallet handles a dumpwallet request by returning  all private
// // keys in a wallet, or an appropiate error if the wallet is locked.
// // TODO: finish this to match bitcoind by writing the dump to a file.
// func dumpWallet(// 	icmd interface{}, w *wallet.Wallet) (interface{}, error) {
// 	keys, err := w.DumpPrivKeys()
// 	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {
// 		return nil, &ErrWalletUnlockNeeded
// 	}
// 	return keys, err
// }

// GetAddressesByAccount handles a getaddressesbyaccount request by returning
// all addresses for an account, or an error if the requested account does
// not exist.
func GetAddressesByAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetAddressesByAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["dumpprivkey"],
			// "invalid subcommand for addnode",
		}
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	addrs, err := w.AccountAddresses(account)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	addrStrs := make([]string, len(addrs))
	for i, a := range addrs {
		addrStrs[i] = a.EncodeAddress()
	}
	return addrStrs, nil
}

// GetBalance handles a getbalance request by returning the balance for an
// account (wallet), or an error if the requested account does not
// exist.
func GetBalance(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetBalanceCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getbalance"],
			// "invalid subcommand for addnode",
		}
	}
	var balance util.Amount
	var err error
	accountName := "*"
	if cmd.Account != nil {
		accountName = *cmd.Account
	}
	if accountName == "*" {
		balance, err = w.CalculateBalance(int32(*cmd.MinConf))
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
	} else {
		var account uint32
		account, err = w.AccountNumber(waddrmgr.KeyScopeBIP0044, accountName)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		bals, err := w.CalculateAccountBalances(account, int32(*cmd.MinConf))
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		balance = bals.Spendable
	}
	return balance.ToDUO(), nil
}

// GetBestBlock handles a getbestblock request by returning a JSON object
// with the height and hash of the most recently processed block.
func GetBestBlock(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	blk := w.Manager.SyncedTo()
	result := &btcjson.GetBestBlockResult{
		Hash:   blk.Hash.String(),
		Height: blk.Height,
	}
	return result, nil
}

// GetBestBlockHash handles a getbestblockhash request by returning the hash
// of the most recently processed block.
func GetBestBlockHash(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	blk := w.Manager.SyncedTo()
	return blk.Hash.String(), nil
}

// GetBlockCount handles a getblockcount request by returning the chain height
// of the most recently processed block.
func GetBlockCount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	blk := w.Manager.SyncedTo()
	return blk.Height, nil
}

// GetInfo handles a getinfo request by returning the a structure containing
// information about the current state of btcwallet.
// exist.
func GetInfo(icmd interface{}, w *wallet.Wallet, chainClient ...*chain.RPCClient) (interface{}, error) {
	if len(chainClient) < 1 || chainClient[0] == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoChain,
			Message: "there is currently no chain client to get this response",
		}
	}
	// Call down to pod for all of the information in this command known
	// by them.
	info, err := chainClient[0].GetInfo()
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	bal, err := w.CalculateBalance(1)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// TODO(davec): This should probably have a database version as opposed
	//  to using the manager version.
	info.WalletVersion = int32(waddrmgr.LatestMgrVersion)
	info.Balance = bal.ToDUO()
	info.PaytxFee = float64(txrules.DefaultRelayFeePerKb)
	// We don't set the following since they don't make much sense in the
	// wallet architecture:
	//  - unlocked_until
	//  - errors
	return info, nil
}

func DecodeAddress(s string, params *netparams.Params) (util.Address, error) {
	addr, err := util.DecodeAddress(s, params)
	if err != nil {
		log.L.Error(err)
		msg := fmt.Sprintf("Invalid address %q: decode failed with %#q", s, err)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}
	}
	if !addr.IsForNet(params) {
		msg := fmt.Sprintf("Invalid address %q: not intended for use on %s",
			addr, params.Name)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: msg,
		}
	}
	return addr, nil
}

// GetAccount handles a getaccount request by returning the account name
// associated with a single address.
func GetAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getaccount"],
			// "invalid subcommand for addnode",
		}
	}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Fetch the associated account
	account, err := w.AccountOfAddress(addr)
	if err != nil {
		log.L.Error(err)
		return nil, &ErrAddressNotInWallet
	}
	acctName, err := w.AccountName(waddrmgr.KeyScopeBIP0044, account)
	if err != nil {
		log.L.Error(err)
		return nil, &ErrAccountNameNotFound
	}
	return acctName, nil
}

// GetAccountAddress handles a getaccountaddress by returning the most
// recently-created chained address that has not yet been used (does not yet
// appear in the blockchain, or any tx that has arrived in the pod mempool).
// If the most recently-requested address has been used, a new address (the
// next chained address in the keypool) is used.  This can fail if the keypool
// runs out (and will return json.ErrRPCWalletKeypoolRanOut if that happens).
func GetAccountAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetAccountAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getaccountaddress"],
			// "invalid subcommand for addnode",
		}
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	addr, err := w.CurrentAddress(account, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return addr.EncodeAddress(), err
}

// GetUnconfirmedBalance handles a getunconfirmedbalance extension request
// by returning the current unconfirmed balance of an account.
func GetUnconfirmedBalance(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetUnconfirmedBalanceCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getunconfirmedbalance"],
			// "invalid subcommand for addnode",
		}
	}
	acctName := "default"
	if cmd.Account != nil {
		acctName = *cmd.Account
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	bals, err := w.CalculateAccountBalances(account, 1)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return (bals.Total - bals.Spendable).ToDUO(), nil
}

// ImportPrivKey handles an importprivkey request by parsing
// a WIF-encoded private key and adding it to an account.
func ImportPrivKey(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ImportPrivKeyCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["importprivkey"],
			// "invalid subcommand for addnode",
		}
	}
	// Ensure that private keys are only imported to the correct account.
	//
	// Yes, Label is the account name.
	if cmd.Label != nil && *cmd.Label != waddrmgr.ImportedAddrAccountName {
		return nil, &ErrNotImportedAccount
	}
	wif, err := util.DecodeWIF(cmd.PrivKey)
	if err != nil {
		log.L.Error(err)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "WIF decode failed: " + err.Error(),
		}
	}
	if !wif.IsForNet(w.ChainParams()) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidAddressOrKey,
			Message: "Key is not intended for " + w.ChainParams().Name,
		}
	}
	// Import the private key, handling any errors.
	_, err = w.ImportPrivateKey(waddrmgr.KeyScopeBIP0044, wif, nil, *cmd.Rescan)
	switch {
	case waddrmgr.IsError(err, waddrmgr.ErrDuplicateAddress):
		// Do not return duplicate key errors to the client.
		return nil, nil
	case waddrmgr.IsError(err, waddrmgr.ErrLocked):
		return nil, &ErrWalletUnlockNeeded
	}
	return nil, err
}

// KeypoolRefill handles the keypoolrefill command. Since we handle the keypool
// automatically this does nothing since refilling is never manually required.
func KeypoolRefill(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	return nil, nil
}

// CreateNewAccount handles a createnewaccount request by creating and
// returning a new account. If the last account has no transaction history
// as per BIP 0044 a new account cannot be created so an error will be returned.
func CreateNewAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.CreateNewAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["createnewaccount"],
			// "invalid subcommand for addnode",
		}
	}
	// The wildcard * is reserved by the rpc server with the special meaning
	// of "all accounts", so disallow naming accounts to this string.
	if cmd.Account == "*" {
		return nil, &ErrReservedAccountName
	}
	_, err := w.NextAccount(waddrmgr.KeyScopeBIP0044, cmd.Account)
	if waddrmgr.IsError(err, waddrmgr.ErrLocked) {
		return nil, &btcjson.RPCError{
			Code: btcjson.ErrRPCWalletUnlockNeeded,
			Message: "Creating an account requires the wallet to be unlocked. " +
				"Enter the wallet passphrase with walletpassphrase to unlock",
		}
	}
	return nil, err
}

// RenameAccount handles a renameaccount request by renaming an account.
// If the account does not exist an appropiate error will be returned.
func RenameAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.RenameAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["renameaccount"],
			// "invalid subcommand for addnode",
		}
	}
	// The wildcard * is reserved by the rpc server with the special meaning
	// of "all accounts", so disallow naming accounts to this string.
	if cmd.NewAccount == "*" {
		return nil, &ErrReservedAccountName
	}
	// Check that given account exists
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.OldAccount)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return nil, w.RenameAccount(waddrmgr.KeyScopeBIP0044, account, cmd.NewAccount)
}

// GetNewAddress handles a getnewaddress request by returning a new
// address for an account.  If the account does not exist an appropiate
// error is returned.
// TODO: Follow BIP 0044 and warn if number of unused addresses exceeds
// the gap limit.
func GetNewAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetNewAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getnewaddress"],
			// "invalid subcommand for addnode",
		}
	}
	acctName := "default"
	if cmd.Account != nil {
		acctName = *cmd.Account
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	addr, err := w.NewAddress(account, waddrmgr.KeyScopeBIP0044, false)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Return the new payment address string.
	return addr.EncodeAddress(), nil
}

// GetRawChangeAddress handles a getrawchangeaddress request by creating
// and returning a new change address for an account.
//
// Note: bitcoind allows specifying the account as an optional parameter,
// but ignores the parameter.
func GetRawChangeAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetRawChangeAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getrawchangeaddress"],
			// "invalid subcommand for addnode",
		}
	}
	acctName := "default"
	if cmd.Account != nil {
		acctName = *cmd.Account
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, acctName)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	addr, err := w.NewChangeAddress(account, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Return the new payment address string.
	return addr.EncodeAddress(), nil
}

// GetReceivedByAccount handles a getreceivedbyaccount request by returning
// the total amount received by addresses of an account.
func GetReceivedByAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetReceivedByAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getreceivedbyaccount"],
			// "invalid subcommand for addnode",
		}
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.Account)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// TODO: This is more inefficient that it could be, but the entire
	// algorithm is already dominated by reading every transaction in the
	// wallet's history.
	results, err := w.TotalReceivedForAccounts(
		waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf),
	)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	acctIndex := int(account)
	if account == waddrmgr.ImportedAddrAccount {
		acctIndex = len(results) - 1
	}
	return results[acctIndex].TotalReceived.ToDUO(), nil
}

// GetReceivedByAddress handles a getreceivedbyaddress request by returning
// the total amount received by a single address.
func GetReceivedByAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetReceivedByAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["getreceivedbyaddress"],
			// "invalid subcommand for addnode",
		}
	}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	total, err := w.TotalReceivedForAddr(addr, int32(*cmd.MinConf))
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return total.ToDUO(), nil
}

// GetTransaction handles a gettransaction request by returning details about
// a single transaction saved by wallet.
func GetTransaction(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.GetTransactionCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["gettransaction"],
			// "invalid subcommand for addnode",
		}
	}
	txHash, err := chainhash.NewHashFromStr(cmd.Txid)
	if err != nil {
		log.L.Error(err)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDecodeHexString,
			Message: "Transaction hash string decode failed: " + err.Error(),
		}
	}
	details, err := wallet.ExposeUnstableAPI(w).TxDetails(txHash)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	if details == nil {
		return nil, &ErrNoTransactionInfo
	}
	syncBlock := w.Manager.SyncedTo()
	// TODO: The serialized transaction is already in the DB, so
	// reserializing can be avoided here.
	var txBuf bytes.Buffer
	txBuf.Grow(details.MsgTx.SerializeSize())
	err = details.MsgTx.Serialize(&txBuf)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// TODO: Add a "generated" field to this result type.  "generated":true
	// is only added if the transaction is a coinbase.
	ret := btcjson.GetTransactionResult{
		TxID:            cmd.Txid,
		Hex:             hex.EncodeToString(txBuf.Bytes()),
		Time:            details.Received.Unix(),
		TimeReceived:    details.Received.Unix(),
		WalletConflicts: []string{}, // Not saved
		// Generated:     blockchain.IsCoinBaseTx(&details.MsgTx),
	}
	if details.Block.Height != -1 {
		ret.BlockHash = details.Block.Hash.String()
		ret.BlockTime = details.Block.Time.Unix()
		ret.Confirmations = int64(Confirms(details.Block.Height, syncBlock.Height))
	}
	var (
		debitTotal  util.Amount
		creditTotal util.Amount // Excludes change
		fee         util.Amount
		feeF64      float64
	)
	for _, deb := range details.Debits {
		debitTotal += deb.Amount
	}
	for _, cred := range details.Credits {
		if !cred.Change {
			creditTotal += cred.Amount
		}
	}
	// Fee can only be determined if every input is a debit.
	if len(details.Debits) == len(details.MsgTx.TxIn) {
		var outputTotal util.Amount
		for _, output := range details.MsgTx.TxOut {
			outputTotal += util.Amount(output.Value)
		}
		fee = debitTotal - outputTotal
		feeF64 = fee.ToDUO()
	}
	if len(details.Debits) == 0 {
		// Credits must be set later, but since we know the full length
		// of the details slice, allocate it with the correct cap.
		ret.Details = make([]btcjson.GetTransactionDetailsResult, 0, len(details.Credits))
	} else {
		ret.Details = make([]btcjson.GetTransactionDetailsResult, 1, len(details.Credits)+1)
		ret.Details[0] = btcjson.GetTransactionDetailsResult{
			// Fields left zeroed:
			//   InvolvesWatchOnly
			//   Account
			//   Address
			//   Vout
			//
			// TODO(jrick): Address and Vout should always be set,
			// but we're doing the wrong thing here by not matching
			// core.  Instead, gettransaction should only be adding
			// details for transaction outputs, just like
			// listtransactions (but using the short result format).
			Category: "send",
			Amount:   (-debitTotal).ToDUO(), // negative since it is a send
			Fee:      &feeF64,
		}
		ret.Fee = feeF64
	}
	credCat := wallet.RecvCategory(details, syncBlock.Height, w.ChainParams()).String()
	for _, cred := range details.Credits {
		// Change is ignored.
		if cred.Change {
			continue
		}
		var address string
		var accountName string
		_, addrs, _, err := txscript.ExtractPkScriptAddrs(
			details.MsgTx.TxOut[cred.Index].PkScript, w.ChainParams())
		if err == nil && len(addrs) == 1 {
			addr := addrs[0]
			address = addr.EncodeAddress()
			account, err := w.AccountOfAddress(addr)
			if err == nil {
				name, err := w.AccountName(waddrmgr.KeyScopeBIP0044, account)
				if err == nil {
					accountName = name
				}
			}
		}
		ret.Details = append(ret.Details, btcjson.GetTransactionDetailsResult{
			// Fields left zeroed:
			//   InvolvesWatchOnly
			//   Fee
			Account:  accountName,
			Address:  address,
			Category: credCat,
			Amount:   cred.Amount.ToDUO(),
			Vout:     cred.Index,
		})
	}
	ret.Amount = creditTotal.ToDUO()
	return ret, nil
}

func HandleDropWalletHistory(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (out interface{}, err error) {
	log.L.Debug("dropping wallet history")
	if err = DropWalletHistory(w)(nil); log.L.Check(err) {
	}
	log.L.Debug("dropped wallet history")
	// go func() {
	// 	rwt, err := w.Database().BeginReadWriteTx()
	// 	if err != nil {
	// 		log.L.Error(err)
	// 	}
	// 	ns := rwt.ReadWriteBucket([]byte("waddrmgr"))
	// 	w.Manager.SetSyncedTo(ns, nil)
	// 	if err = rwt.Commit(); log.L.Check(err) {
	// 	}
	// }()
	defer interrupt.RequestRestart()
	return nil, err
}

// These generators create the following global variables in this package:
//
//   var localeHelpDescs map[string]func() map[string]string
//   var requestUsages string
//
// localeHelpDescs maps from locale strings (e.g. "en_US") to a function that
// builds a map of help texts for each RPC server method.  This prevents help
// text maps for every locale map from being rooted and created during init.
// Instead, the appropiate function is looked up when help text is first needed
// using the current locale and saved to the global below for futher reuse.
//
// requestUsages contains single line usages for every supported request,
// separated by newlines.  It is set during init.  These usages are used for all
// locales.
//
//go:generate go run ../../internal/rpchelp/genrpcserverhelp.go legacyrpc
//go:generate gofmt -w rpcserverhelp.go
var HelpDescs map[string]string
var HelpDescsMutex sync.Mutex // Help may execute concurrently, so synchronize access.

// HelpWithChainRPC handles the help request when the RPC server has been
// associated with a consensus RPC client.  The additional RPC client is used to
// include help messages for methods implemented by the consensus server via RPC
// passthrough.
func HelpWithChainRPC(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	return Help(icmd, w, chainClient[0])
}

// HelpNoChainRPC handles the help request when the RPC server has not been
// associated with a consensus RPC client.  No help messages are included for
// passthrough requests.
func HelpNoChainRPC(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	return Help(icmd, w, nil)
}

// Help handles the Help request by returning one line usage of all available
// methods, or full Help for a specific method.  The chainClient is optional,
// and this is simply a helper function for the HelpNoChainRPC and
// HelpWithChainRPC handlers.
func Help(icmd interface{}, w *wallet.Wallet, chainClient *chain.RPCClient) (interface{}, error) {
	cmd := icmd.(*btcjson.HelpCmd)
	// pod returns different help messages depending on the kind of
	// connection the client is using.  Only methods availble to HTTP POST
	// clients are available to be used by wallet clients, even though
	// wallet itself is a websocket client to pod.  Therefore, create a
	// POST client as needed.
	//
	// Returns nil if chainClient is currently nil or there is an error
	// creating the client.
	//
	// This is hacky and is probably better handled by exposing help usage
	// texts in a non-internal pod package.
	postClient := func() *rpcclient.Client {
		if chainClient == nil {
			return nil
		}
		c, err := chainClient.POSTClient()
		if err != nil {
			log.L.Error(err)
			return nil
		}
		return c
	}
	if cmd.Command == nil || *cmd.Command == "" {
		// Prepend chain server usage if it is available.
		usages := RequestUsages
		client := postClient()
		if client != nil {
			rawChainUsage, err := client.RawRequest("help", nil)
			var chainUsage string
			if err == nil {
				_ = js.Unmarshal([]byte(rawChainUsage), &chainUsage)
			}
			if chainUsage != "" {
				usages = "Chain server usage:\n\n" + chainUsage + "\n\n" +
					"Wallet server usage (overrides chain requests):\n\n" +
					RequestUsages
			}
		}
		return usages, nil
	}
	defer HelpDescsMutex.Unlock()
	HelpDescsMutex.Lock()
	if HelpDescs == nil {
		// TODO: Allow other locales to be set via config or detemine
		// this from environment variables.  For now, hardcode US
		// English.
		HelpDescs = LocaleHelpDescs["en_US"]()
	}
	helpText, ok := HelpDescs[*cmd.Command]
	if ok {
		return helpText, nil
	}
	// Return the chain server's detailed help if possible.
	var chainHelp string
	client := postClient()
	if client != nil {
		param := make([]byte, len(*cmd.Command)+2)
		param[0] = '"'
		copy(param[1:], *cmd.Command)
		param[len(param)-1] = '"'
		rawChainHelp, err := client.RawRequest("help", []js.RawMessage{param})
		if err == nil {
			_ = js.Unmarshal([]byte(rawChainHelp), &chainHelp)
		}
	}
	if chainHelp != "" {
		return chainHelp, nil
	}
	return nil, &btcjson.RPCError{
		Code:    btcjson.ErrRPCInvalidParameter,
		Message: fmt.Sprintf("No help for method '%s'", *cmd.Command),
	}
}

// ListAccounts handles a listaccounts request by returning a map of account
// names to their balances.
func ListAccounts(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListAccountsCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listaccounts"],
			// "invalid subcommand for addnode",
		}
	}
	accountBalances := map[string]float64{}
	results, err := w.AccountBalances(waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf))
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	for _, result := range results {
		accountBalances[result.AccountName] = result.AccountBalance.ToDUO()
	}
	// Return the map.  This will be marshaled into a JSON object.
	return accountBalances, nil
}

// ListLockUnspent handles a listlockunspent request by returning an slice of
// all locked outpoints.
func ListLockUnspent(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	return w.LockedOutpoints(), nil
}

// ListReceivedByAccount handles a listreceivedbyaccount request by returning
// a slice of objects, each one containing:
//  "account": the receiving account;
//  "amount": total amount received by the account;
//  "confirmations": number of confirmations of the most recent transaction.
// It takes two parameters:
//  "minconf": minimum number of confirmations to consider a transaction -
//             default: one;
//  "includeempty": whether or not to include addresses that have no transactions -
//                  default: false.
func ListReceivedByAccount(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListReceivedByAccountCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listreceivedbyaccount"],
			// "invalid subcommand for addnode",
		}
	}
	results, err := w.TotalReceivedForAccounts(
		waddrmgr.KeyScopeBIP0044, int32(*cmd.MinConf),
	)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	jsonResults := make([]btcjson.ListReceivedByAccountResult, 0, len(results))
	for _, result := range results {
		jsonResults = append(jsonResults, btcjson.ListReceivedByAccountResult{
			Account:       result.AccountName,
			Amount:        result.TotalReceived.ToDUO(),
			Confirmations: uint64(result.LastConfirmation),
		})
	}
	return jsonResults, nil
}

// ListReceivedByAddress handles a listreceivedbyaddress request by returning
// a slice of objects, each one containing:
//  "account": the account of the receiving address;
//  "address": the receiving address;
//  "amount": total amount received by the address;
//  "confirmations": number of confirmations of the most recent transaction.
// It takes two parameters:
//  "minconf": minimum number of confirmations to consider a transaction -
//             default: one;
//  "includeempty": whether or not to include addresses that have no transactions -
//                  default: false.
func ListReceivedByAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListReceivedByAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listreceivedbyaddress"],
			// "invalid subcommand for addnode",
		}
	}
	// Intermediate data for each address.
	type AddrData struct {
		// Total amount received.
		amount util.Amount
		// Number of confirmations of the last transaction.
		confirmations int32
		// Hashes of transactions which include an output paying to the address
		tx []string
		// Account which the address belongs to
		// account string
	}
	syncBlock := w.Manager.SyncedTo()
	// Intermediate data for all addresses.
	allAddrData := make(map[string]AddrData)
	// Create an AddrData entry for each active address in the account.
	// Otherwise we'll just get addresses from transactions later.
	sortedAddrs, err := w.SortedActivePaymentAddresses()
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	for _, address := range sortedAddrs {
		// There might be duplicates, just overwrite them.
		allAddrData[address] = AddrData{}
	}
	minConf := *cmd.MinConf
	var endHeight int32
	if minConf == 0 {
		endHeight = -1
	} else {
		endHeight = syncBlock.Height - int32(minConf) + 1
	}
	err = wallet.ExposeUnstableAPI(w).RangeTransactions(0, endHeight, func(details []wtxmgr.TxDetails) (bool, error) {
		confirmations := Confirms(details[0].Block.Height, syncBlock.Height)
		for _, tx := range details {
			for _, cred := range tx.Credits {
				pkScript := tx.MsgTx.TxOut[cred.Index].PkScript
				_, addrs, _, err := txscript.ExtractPkScriptAddrs(
					pkScript, w.ChainParams())
				if err != nil {
					log.L.Error(err)
					// Non standard script, skip.
					continue
				}
				for _, addr := range addrs {
					addrStr := addr.EncodeAddress()
					addrData, ok := allAddrData[addrStr]
					if ok {
						addrData.amount += cred.Amount
						// Always overwrite confirmations with newer ones.
						addrData.confirmations = confirmations
					} else {
						addrData = AddrData{
							amount:        cred.Amount,
							confirmations: confirmations,
						}
					}
					addrData.tx = append(addrData.tx, tx.Hash.String())
					allAddrData[addrStr] = addrData
				}
			}
		}
		return false, nil
	})
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Massage address data into output format.
	numAddresses := len(allAddrData)
	ret := make([]btcjson.ListReceivedByAddressResult, numAddresses)
	idx := 0
	for address, addrData := range allAddrData {
		ret[idx] = btcjson.ListReceivedByAddressResult{
			Address:       address,
			Amount:        addrData.amount.ToDUO(),
			Confirmations: uint64(addrData.confirmations),
			TxIDs:         addrData.tx,
		}
		idx++
	}
	return ret, nil
}

// ListSinceBlock handles a listsinceblock request by returning an array of maps
// with details of sent and received wallet transactions since the given block.
func ListSinceBlock(icmd interface{}, w *wallet.Wallet,
	cc ...*chain.RPCClient) (interface{}, error) {
	if len(cc) < 1 || cc[0] == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoChain,
			Message: "there is currently no chain client to get this response",
		}
	}
	chainClient := cc[0]
	cmd, ok := icmd.(*btcjson.ListSinceBlockCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listsinceblock"],
			// "invalid subcommand for addnode",
		}
	}
	syncBlock := w.Manager.SyncedTo()
	targetConf := int64(*cmd.TargetConfirmations)
	// For the result we need the block hash for the last block counted
	// in the blockchain due to confirmations. We send this off now so that
	// it can arrive asynchronously while we figure out the rest.
	gbh := chainClient.GetBlockHashAsync(int64(syncBlock.Height) + 1 - targetConf)
	var start int32
	if cmd.BlockHash != nil {
		hash, err := chainhash.NewHashFromStr(*cmd.BlockHash)
		if err != nil {
			log.L.Error(err)
			return nil, DeserializationError{err}
		}
		block, err := chainClient.GetBlockVerboseTx(hash)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		start = int32(block.Height) + 1
	}
	txInfoList, err := w.ListSinceBlock(start, -1, syncBlock.Height)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Done with work, get the response.
	blockHash, err := gbh.Receive()
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	res := btcjson.ListSinceBlockResult{
		Transactions: txInfoList,
		LastBlock:    blockHash.String(),
	}
	return res, nil
}

// ListTransactions handles a listtransactions request by returning an
// array of maps with details of sent and recevied wallet transactions.
func ListTransactions(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (txs interface{}, err error) {
	if len(chainClient) < 1 || chainClient[0] == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoChain,
			Message: "there is currently no chain client to get this response",
		}
	}
	cmd, ok := icmd.(*btcjson.ListTransactionsCmd)
	if !ok || cmd.From == nil || cmd.Count == nil || cmd.Account != nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listtransactions"],
			// "invalid subcommand for addnode",
		}
	}
	// TODO: ListTransactions does not currently understand the difference
	//  between transactions pertaining to one account from another.  This
	//  will be resolved when wtxmgr is combined with the waddrmgr namespace.
	if *cmd.Account != "*" {
		// For now, don't bother trying to continue if the user
		// specified an account, since this can't be (easily or
		// efficiently) calculated.
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCWallet,
			Message: "Transactions are not yet grouped by account",
		}
	}
	txs, err = w.ListTransactions(*cmd.From, *cmd.Count)
	return txs, err
}

// ListAddressTransactions handles a listaddresstransactions request by
// returning an array of maps with details of spent and received wallet
// transactions.  The form of the reply is identical to listtransactions,
// but the array elements are limited to transaction details which are
// about the addresess included in the request.
func ListAddressTransactions(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListAddressTransactionsCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listaddresstransactions"],
			// "invalid subcommand for addnode",
		}
	}
	if cmd.Account != nil && *cmd.Account != "*" {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Listing transactions for addresses may only be done for all accounts",
		}
	}
	// Decode addresses.
	hash160Map := make(map[string]struct{})
	for _, addrStr := range cmd.Addresses {
		addr, err := DecodeAddress(addrStr, w.ChainParams())
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		hash160Map[string(addr.ScriptAddress())] = struct{}{}
	}
	return w.ListAddressTransactions(hash160Map)
}

// ListAllTransactions handles a listalltransactions request by returning
// a map with details of sent and recevied wallet transactions.  This is
// similar to ListTransactions, except it takes only a single optional
// argument for the account name and replies with all transactions.
func ListAllTransactions(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListAllTransactionsCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listalltransactions"],
			// "invalid subcommand for addnode",
		}
	}
	if cmd.Account != nil && *cmd.Account != "*" {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: "Listing all transactions may only be done for all accounts",
		}
	}
	return w.ListAllTransactions()
}

// ListUnspent handles the listunspent command.
func ListUnspent(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ListUnspentCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["listunspent"],
			// "invalid subcommand for addnode",
		}
	}
	var addresses map[string]struct{}
	if cmd.Addresses != nil {
		addresses = make(map[string]struct{})
		// confirm that all of them are good:
		for _, as := range *cmd.Addresses {
			a, err := DecodeAddress(as, w.ChainParams())
			if err != nil {
				log.L.Error(err)
				return nil, err
			}
			addresses[a.EncodeAddress()] = struct{}{}
		}
	}
	return w.ListUnspent(int32(*cmd.MinConf), int32(*cmd.MaxConf), addresses)
}

// LockUnspent handles the lockunspent command.
func LockUnspent(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.LockUnspentCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["lockunspent"],
			// "invalid subcommand for addnode",
		}
	}
	switch {
	case cmd.Unlock && len(cmd.Transactions) == 0:
		w.ResetLockedOutpoints()
	default:
		for _, input := range cmd.Transactions {
			txHash, err := chainhash.NewHashFromStr(input.Txid)
			if err != nil {
				log.L.Error(err)
				return nil, ParseError{err}
			}
			op := wire.OutPoint{Hash: *txHash, Index: input.Vout}
			if cmd.Unlock {
				w.UnlockOutpoint(op)
			} else {
				w.LockOutpoint(op)
			}
		}
	}
	return true, nil
}

// MakeOutputs creates a slice of transaction outputs from a pair of address
// strings to amounts.  This is used to create the outputs to include in newly
// created transactions from a JSON object describing the output destinations
// and amounts.
func MakeOutputs(pairs map[string]util.Amount, chainParams *netparams.Params) ([]*wire.TxOut, error) {
	outputs := make([]*wire.TxOut, 0, len(pairs))
	for addrStr, amt := range pairs {
		addr, err := util.DecodeAddress(addrStr, chainParams)
		if err != nil {
			log.L.Error(err)
			return nil, fmt.Errorf("cannot decode address: %s", err)
		}
		pkScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			log.L.Error(err)
			return nil, fmt.Errorf("cannot create txout script: %s", err)
		}
		outputs = append(outputs, wire.NewTxOut(int64(amt), pkScript))
	}
	return outputs, nil
}

// SendPairs creates and sends payment transactions.
// It returns the transaction hash in string format upon success
// All errors are returned in json.RPCError format
func SendPairs(w *wallet.Wallet, amounts map[string]util.Amount,
	account uint32, minconf int32, feeSatPerKb util.Amount) (string, error) {
	outputs, err := MakeOutputs(amounts, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return "", err
	}
	txHash, err := w.SendOutputs(outputs, account, minconf, feeSatPerKb)
	if err != nil {
		log.L.Error(err)
		if err == txrules.ErrAmountNegative {
			return "", ErrNeedPositiveAmount
		}
		if waddrmgr.IsError(err, waddrmgr.ErrLocked) {
			return "", &ErrWalletUnlockNeeded
		}
		switch err.(type) {
		case btcjson.RPCError:
			return "", err
		}
		return "", &btcjson.RPCError{
			Code:    btcjson.ErrRPCInternal.Code,
			Message: err.Error(),
		}
	}
	txHashStr := txHash.String()
	log.L.Info("successfully sent transaction", txHashStr)
	return txHashStr, nil
}
func IsNilOrEmpty(s *string) bool {
	return s == nil || *s == ""
}

// SendFrom handles a sendfrom RPC request by creating a new transaction
// spending unspent transaction outputs for a wallet to another payment
// address.  Leftover inputs not sent to the payment address or a fee for
// the miner are sent back to a new address in the wallet.  Upon success,
// the TxID for the created transaction is returned.
func SendFrom(icmd interface{}, w *wallet.Wallet, chainClient *chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.SendFromCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["sendfrom"],
			// "invalid subcommand for addnode",
		}
	}
	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.
	if !IsNilOrEmpty(cmd.Comment) || !IsNilOrEmpty(cmd.CommentTo) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}
	}
	account, err := w.AccountNumber(
		waddrmgr.KeyScopeBIP0044, cmd.FromAccount,
	)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Check that signed integer parameters are positive.
	if cmd.Amount < 0 {
		return nil, ErrNeedPositiveAmount
	}
	minConf := int32(*cmd.MinConf)
	if minConf < 0 {
		return nil, ErrNeedPositiveMinconf
	}
	// Create map of address and amount pairs.
	amt, err := util.NewAmount(cmd.Amount)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	pairs := map[string]util.Amount{
		cmd.ToAddress: amt,
	}
	return SendPairs(w, pairs, account, minConf,
		txrules.DefaultRelayFeePerKb)
}

// SendMany handles a sendmany RPC request by creating a new transaction
// spending unspent transaction outputs for a wallet to any number of
// payment addresses.  Leftover inputs not sent to the payment address
// or a fee for the miner are sent back to a new address in the wallet.
// Upon success, the TxID for the created transaction is returned.
func SendMany(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.SendManyCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["sendmany"],
			// "invalid subcommand for addnode",
		}
	}
	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.
	if !IsNilOrEmpty(cmd.Comment) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}
	}
	account, err := w.AccountNumber(waddrmgr.KeyScopeBIP0044, cmd.FromAccount)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Check that minconf is positive.
	minConf := int32(*cmd.MinConf)
	if minConf < 0 {
		return nil, ErrNeedPositiveMinconf
	}
	// Recreate address/amount pairs, using dcrutil.Amount.
	pairs := make(map[string]util.Amount, len(cmd.Amounts))
	for k, v := range cmd.Amounts {
		amt, err := util.NewAmount(v)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		pairs[k] = amt
	}
	return SendPairs(w, pairs, account, minConf, txrules.DefaultRelayFeePerKb)
}

// SendToAddress handles a sendtoaddress RPC request by creating a new
// transaction spending unspent transaction outputs for a wallet to another
// payment address.  Leftover inputs not sent to the payment address or a fee
// for the miner are sent back to a new address in the wallet.  Upon success,
// the TxID for the created transaction is returned.
func SendToAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.SendToAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["sendtoaddress"],
			// "invalid subcommand for addnode",
		}
	}
	// Transaction comments are not yet supported.  Error instead of
	// pretending to save them.
	if !IsNilOrEmpty(cmd.Comment) || !IsNilOrEmpty(cmd.CommentTo) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCUnimplemented,
			Message: "Transaction comments are not yet supported",
		}
	}
	amt, err := util.NewAmount(cmd.Amount)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Check that signed integer parameters are positive.
	if amt < 0 {
		return nil, ErrNeedPositiveAmount
	}
	// Mock up map of address and amount pairs.
	pairs := map[string]util.Amount{
		cmd.Address: amt,
	}
	// sendtoaddress always spends from the default account, this matches bitcoind
	return SendPairs(w, pairs, waddrmgr.DefaultAccountNum, 1,
		txrules.DefaultRelayFeePerKb)
}

// SetTxFee sets the transaction fee per kilobyte added to transactions.
func SetTxFee(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.SetTxFeeCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["settxfee"],
			// "invalid subcommand for addnode",
		}
	}
	// Check that amount is not negative.
	if cmd.Amount < 0 {
		return nil, ErrNeedPositiveAmount
	}
	// A boolean true result is returned upon success.
	return true, nil
}

// SignMessage signs the given message with the private key for the given
// address
func SignMessage(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.SignMessageCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["signmessage"],
			// "invalid subcommand for addnode",
		}
	}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	privKey, err := w.PrivKeyForAddress(addr)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	var buf bytes.Buffer
	err = wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	if err != nil {
		log.L.Error(err)
		log.L.Debug(err)
	}
	err = wire.WriteVarString(&buf, 0, cmd.Message)
	if err != nil {
		log.L.Error(err)
		log.L.Debug(err)
	}
	messageHash := chainhash.DoubleHashB(buf.Bytes())
	sigbytes, err := ec.SignCompact(ec.S256(), privKey,
		messageHash, true)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	return base64.StdEncoding.EncodeToString(sigbytes), nil
}

// SignRawTransaction handles the signrawtransaction command.
func SignRawTransaction(icmd interface{}, w *wallet.Wallet,
	cc ...*chain.RPCClient) (interface{}, error) {
	if len(cc) < 1 || cc[0] == nil {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCNoChain,
			Message: "there is currently no chain client to get this response",
		}
	}
	chainClient := cc[0]
	cmd, ok := icmd.(*btcjson.SignRawTransactionCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["signrawtransaction"],
			// "invalid subcommand for addnode",
		}
	}
	serializedTx, err := DecodeHexStr(cmd.RawTx)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewBuffer(serializedTx))
	if err != nil {
		log.L.Error(err)
		e := errors.New("TX decode failed")
		return nil, DeserializationError{e}
	}
	var hashType txscript.SigHashType
	switch *cmd.Flags {
	case "ALL":
		hashType = txscript.SigHashAll
	case "NONE":
		hashType = txscript.SigHashNone
	case "SINGLE":
		hashType = txscript.SigHashSingle
	case "ALL|ANYONECANPAY":
		hashType = txscript.SigHashAll | txscript.SigHashAnyOneCanPay
	case "NONE|ANYONECANPAY":
		hashType = txscript.SigHashNone | txscript.SigHashAnyOneCanPay
	case "SINGLE|ANYONECANPAY":
		hashType = txscript.SigHashSingle | txscript.SigHashAnyOneCanPay
	default:
		e := errors.New("Invalid sighash parameter")
		return nil, InvalidParameterError{e}
	}
	// TODO: really we probably should look these up with pod anyway to
	// make sure that they match the blockchain if present.
	inputs := make(map[wire.OutPoint][]byte)
	scripts := make(map[string][]byte)
	var cmdInputs []btcjson.RawTxInput
	if cmd.Inputs != nil {
		cmdInputs = *cmd.Inputs
	}
	for _, rti := range cmdInputs {
		inputHash, err := chainhash.NewHashFromStr(rti.Txid)
		if err != nil {
			log.L.Error(err)
			return nil, DeserializationError{err}
		}
		script, err := DecodeHexStr(rti.ScriptPubKey)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		// redeemScript is only actually used iff the user provided
		// private keys. In which case, it is used to get the scripts
		// for signing. If the user did not provide keys then we always
		// get scripts from the wallet.
		// Empty strings are ok for this one and hex.DecodeString will
		// DTRT.
		if cmd.PrivKeys != nil && len(*cmd.PrivKeys) != 0 {
			redeemScript, err := DecodeHexStr(rti.RedeemScript)
			if err != nil {
				log.L.Error(err)
				return nil, err
			}
			addr, err := util.NewAddressScriptHash(redeemScript,
				w.ChainParams())
			if err != nil {
				log.L.Error(err)
				return nil, DeserializationError{err}
			}
			scripts[addr.String()] = redeemScript
		}
		inputs[wire.OutPoint{
			Hash:  *inputHash,
			Index: rti.Vout,
		}] = script
	}
	// Now we go and look for any inputs that we were not provided by
	// querying pod with getrawtransaction. We queue up a bunch of async
	// requests and will wait for replies after we have checked the rest of
	// the arguments.
	requested := make(map[wire.OutPoint]rpcclient.FutureGetTxOutResult)
	for _, txIn := range tx.TxIn {
		// Did we get this outpoint from the arguments?
		if _, ok := inputs[txIn.PreviousOutPoint]; ok {
			continue
		}
		// Asynchronously request the output script.
		requested[txIn.PreviousOutPoint] = chainClient.GetTxOutAsync(
			&txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index,
			true)
	}
	// Parse list of private keys, if present. If there are any keys here
	// they are the keys that we may use for signing. If empty we will
	// use any keys known to us already.
	var keys map[string]*util.WIF
	if cmd.PrivKeys != nil {
		keys = make(map[string]*util.WIF)
		for _, key := range *cmd.PrivKeys {
			wif, err := util.DecodeWIF(key)
			if err != nil {
				log.L.Error(err)
				return nil, DeserializationError{err}
			}
			if !wif.IsForNet(w.ChainParams()) {
				s := "key network doesn't match wallet's"
				return nil, DeserializationError{errors.New(s)}
			}
			addr, err := util.NewAddressPubKey(wif.SerializePubKey(),
				w.ChainParams())
			if err != nil {
				log.L.Error(err)
				return nil, DeserializationError{err}
			}
			keys[addr.EncodeAddress()] = wif
		}
	}
	// We have checked the rest of the args. now we can collect the async
	// txs. TODO: If we don't mind the possibility of wasting work we could
	// move waiting to the following loop and be slightly more asynchronous.
	for outPoint, resp := range requested {
		result, err := resp.Receive()
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		script, err := hex.DecodeString(result.ScriptPubKey.Hex)
		if err != nil {
			log.L.Error(err)
			return nil, err
		}
		inputs[outPoint] = script
	}
	// All args collected. Now we can sign all the inputs that we can.
	// `complete' denotes that we successfully signed all outputs and that
	// all scripts will run to completion. This is returned as part of the
	// reply.
	signErrs, err := w.SignTransaction(&tx, hashType, inputs, keys, scripts)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	var buf bytes.Buffer
	buf.Grow(tx.SerializeSize())
	// All returned errors (not OOM, which panics) encounted during
	// bytes.Buffer writes are unexpected.
	if err = tx.Serialize(&buf); err != nil {
		panic(err)
	}
	signErrors := make([]btcjson.SignRawTransactionError, 0, len(signErrs))
	for _, e := range signErrs {
		input := tx.TxIn[e.InputIndex]
		signErrors = append(signErrors, btcjson.SignRawTransactionError{
			TxID:      input.PreviousOutPoint.Hash.String(),
			Vout:      input.PreviousOutPoint.Index,
			ScriptSig: hex.EncodeToString(input.SignatureScript),
			Sequence:  input.Sequence,
			Error:     e.Error.Error(),
		})
	}
	return btcjson.SignRawTransactionResult{
		Hex:      hex.EncodeToString(buf.Bytes()),
		Complete: len(signErrors) == 0,
		Errors:   signErrors,
	}, nil
}

// ValidateAddress handles the validateaddress command.
func ValidateAddress(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.ValidateAddressCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["validateaddress"],
			// "invalid subcommand for addnode",
		}
	}
	result := btcjson.ValidateAddressWalletResult{}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		// Use result zero value (IsValid=false).
		return result, nil
	}
	// We could put whether or not the address is a script here,
	// by checking the type of "addr", however, the reference
	// implementation only puts that information if the script is
	// "ismine", and we follow that behaviour.
	result.Address = addr.EncodeAddress()
	result.IsValid = true
	ainfo, err := w.AddressInfo(addr)
	if err != nil {
		log.L.Error(err)
		if waddrmgr.IsError(err, waddrmgr.ErrAddressNotFound) {
			// No additional information available about the address.
			return result, nil
		}
		return nil, err
	}
	// The address lookup was successful which means there is further
	// information about it available and it is "mine".
	result.IsMine = true
	acctName, err := w.AccountName(waddrmgr.KeyScopeBIP0044, ainfo.Account())
	if err != nil {
		log.L.Error(err)
		return nil, &ErrAccountNameNotFound
	}
	result.Account = acctName
	switch ma := ainfo.(type) {
	case waddrmgr.ManagedPubKeyAddress:
		result.IsCompressed = ma.Compressed()
		result.PubKey = ma.ExportPubKey()
	case waddrmgr.ManagedScriptAddress:
		result.IsScript = true
		// The script is only available if the manager is unlocked, so
		// just break out now if there is an error.
		script, err := ma.Script()
		if err != nil {
			log.L.Error(err)
			break
		}
		result.Hex = hex.EncodeToString(script)
		// This typically shouldn't fail unless an invalid script was
		// imported.  However, if it fails for any reason, there is no
		// further information available, so just set the script type
		// a non-standard and break out now.
		class, addrs, reqSigs, err := txscript.ExtractPkScriptAddrs(
			script, w.ChainParams())
		if err != nil {
			log.L.Error(err)
			result.Script = txscript.NonStandardTy.String()
			break
		}
		addrStrings := make([]string, len(addrs))
		for i, a := range addrs {
			addrStrings[i] = a.EncodeAddress()
		}
		result.Addresses = addrStrings
		// Multi-signature scripts also provide the number of required
		// signatures.
		result.Script = class.String()
		if class == txscript.MultiSigTy {
			result.SigsRequired = int32(reqSigs)
		}
	}
	return result, nil
}

// VerifyMessage handles the verifymessage command by verifying the provided
// compact signature for the given address and message.
func VerifyMessage(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.VerifyMessageCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["verifymessage"],
			// "invalid subcommand for addnode",
		}
	}
	addr, err := DecodeAddress(cmd.Address, w.ChainParams())
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// decode base64 signature
	sig, err := base64.StdEncoding.DecodeString(cmd.Signature)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	// Validate the signature - this just shows that it was valid at all.
	// we will compare it with the key next.
	var buf bytes.Buffer
	err = wire.WriteVarString(&buf, 0, "Parallelcoin Signed Message:\n")
	if err != nil {
		log.L.Error(err)
		log.L.Debug(err)
	}
	err = wire.WriteVarString(&buf, 0, cmd.Message)
	if err != nil {
		log.L.Error(err)
		log.L.Debug(err)
	}
	expectedMessageHash := chainhash.DoubleHashB(buf.Bytes())
	pk, wasCompressed, err := ec.RecoverCompact(ec.S256(), sig,
		expectedMessageHash)
	if err != nil {
		log.L.Error(err)
		return nil, err
	}
	var serializedPubKey []byte
	if wasCompressed {
		serializedPubKey = pk.SerializeCompressed()
	} else {
		serializedPubKey = pk.SerializeUncompressed()
	}
	// Verify that the signed-by address matches the given address
	switch checkAddr := addr.(type) {
	case *util.AddressPubKeyHash: // ok
		return bytes.Equal(util.Hash160(serializedPubKey), checkAddr.Hash160()[:]), nil
	case *util.AddressPubKey: // ok
		return string(serializedPubKey) == checkAddr.String(), nil
	default:
		return nil, errors.New("address type not supported")
	}
}

// WalletIsLocked handles the walletislocked extension request by
// returning the current lock state (false for unlocked, true for locked)
// of an account.
func WalletIsLocked(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	return w.Locked(), nil
}

// WalletLock handles a walletlock request by locking the all account
// wallets, returning an error if any wallet is not encrypted (for example,
// a watching-only wallet).
func WalletLock(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	w.Lock()
	return nil, nil
}

// WalletPassphrase responds to the walletpassphrase request by unlocking
// the wallet.  The decryption key is saved in the wallet until timeout
// seconds expires, after which the wallet is locked.
func WalletPassphrase(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.WalletPassphraseCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["walletpassphrase"],
			// "invalid subcommand for addnode",
		}
	}
	timeout := time.Second * time.Duration(cmd.Timeout)
	var unlockAfter <-chan time.Time
	if timeout != 0 {
		unlockAfter = time.After(timeout)
	}
	err := w.Unlock([]byte(cmd.Passphrase), unlockAfter)
	return nil, err
}

// WalletPassphraseChange responds to the walletpassphrasechange request
// by unlocking all accounts with the provided old passphrase, and
// re-encrypting each private key with an AES key derived from the new
// passphrase.
//
// If the old passphrase is correct and the passphrase is changed, all
// wallets will be immediately locked.
func WalletPassphraseChange(icmd interface{}, w *wallet.Wallet,
	chainClient ...*chain.RPCClient) (interface{}, error) {
	cmd, ok := icmd.(*btcjson.WalletPassphraseChangeCmd)
	if !ok {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCInvalidParameter,
			Message: HelpDescsEnUS()["walletpassphrasechange"],
			// "invalid subcommand for addnode",
		}
	}
	err := w.ChangePrivatePassphrase([]byte(cmd.OldPassphrase),
		[]byte(cmd.NewPassphrase))
	if waddrmgr.IsError(err, waddrmgr.ErrWrongPassphrase) {
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCWalletPassphraseIncorrect,
			Message: "Incorrect passphrase",
		}
	}
	return nil, err
}

// DecodeHexStr decodes the hex encoding of a string, possibly prepending a
// leading '0' character if there is an odd number of bytes in the hex string.
// This is to prevent an error for an invalid hex string when using an odd
// number of bytes when calling hex.Decode.
func DecodeHexStr(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		log.L.Error(err)
		return nil, &btcjson.RPCError{
			Code:    btcjson.ErrRPCDecodeHexString,
			Message: "Hex string decode failed: " + err.Error(),
		}
	}
	return decoded, nil
}
