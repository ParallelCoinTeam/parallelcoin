package walletrpc

//
// import (
// 	"github.com/p9c/pod/pkg/btcjson"
// )
//
// // errors are returned as *btcjson.RPCError type
// type (
// 	None                   struct{}
// 	AddMultiSigAddressRes struct {
// 		e error
// 		Res string
// 	}
// 	CreateMultiSigRes struct {
// 		e error
// 		Res btcjson.CreateMultiSigResult
// 	}
// 	DumpPrivKeyRes struct {
// 		e error
// 		Res string
// 	}
// 	GetAccountRes struct {
// 		e error
// 		Res string
// 	}
// 	GetAccountAddressRes struct {
// 		e error
// 		Res string
// 	}
// 	GetAddressesByAccountRes struct {
// 		e error
// 		Res []string
// 	}
// 	GetBalanceRes struct {
// 		e error
// 		Res float64
// 	}
// 	GetBestBlockHashRes struct {
// 		e error
// 		Res string
// 	}
// 	GetBlockCountRes struct {
// 		e error
// 		Res int32
// 	}
// 	GetInfoRes struct {
// 		e error
// 		Res btcjson.InfoWalletResult
// 	}
// 	GetNewAddressRes struct {
// 		e error
// 		Res string
// 	}
// 	GetRawChangeAddressRes struct {
// 		e error
// 		Res string
// 	}
// 	GetReceivedByAccountRes struct {
// 		e error
// 		Res float64
// 	}
// 	GetReceivedByAddressRes struct {
// 		e error
// 		Res float64
// 	}
// 	GetTransactionRes struct {
// 		e error
// 		Res btcjson.GetTransactionResult
// 	}
// 	HelpWithChainRPCRes struct {
// 		e error
// 		Res string
// 	}
// 	HelpNoChainRPCRes struct {
// 		e error
// 		Res string
// 	}
// 	ImportPrivKeyRes struct {
// 		e error
// 		Res None
// 	}
// 	KeypoolRefillRes struct {
// 		e error
// 		Res None
// 	}
// 	ListAccountsRes struct {
// 		e error
// 		Res map[string]float64
// 	}
// 	ListLockUnspentRes struct {
// 		e error
// 		Res []btcjson.TransactionInput
// 	}
// 	ListReceivedByAccountRes struct {
// 		e error
// 		Res []btcjson.ListReceivedByAccountResult
// 	}
// 	ListReceivedByAddressRes struct {
// 		e error
// 		Res btcjson.ListReceivedByAddressResult
// 	}
// 	ListSinceBlockRes struct {
// 		e error
// 		Res btcjson.ListSinceBlockResult
// 	}
// 	ListTransactionsRes struct {
// 		e error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	ListUnspentRes struct {
// 		e error
// 		Res []btcjson.ListUnspentResult
// 	}
// 	LockUnspentRes struct {
// 		e error
// 		Res bool
// 	}
// 	SendFromRes struct {
// 		e error
// 		Res string
// 	}
// 	SendManyRes struct {
// 		e error
// 		Res string
// 	}
// 	SendToAddressRes struct {
// 		e error
// 		Res string
// 	}
// 	SetTxFeeRes struct {
// 		e error
// 		Res bool
// 	}
// 	SignMessageRes struct {
// 		e error
// 		Res string
// 	}
// 	SignRawTransactionRes struct {
// 		e error
// 		Res btcjson.SignRawTransactionResult
// 	}
// 	ValidateAddressRes struct {
// 		e error
// 		Res btcjson.ValidateAddressWalletResult
// 	}
// 	VerifyMessageRes struct {
// 		e error
// 		Res bool
// 	}
// 	WalletLockRes struct {
// 		e error
// 		Res None
// 	}
// 	WalletPassphraseRes struct {
// 		e error
// 		Res None
// 	}
// 	WalletPassphraseChangeRes struct {
// 		e error
// 		Res None
// 	}
// 	CreateNewAccountRes struct {
// 		e error
// 		Res None
// 	}
// 	GetBestBlockRes struct {
// 		e error
// 		Res btcjson.GetBestBlockResult
// 	}
// 	GetUnconfirmedBalanceRes struct {
// 		e error
// 		Res float64
// 	}
// 	ListAddressTransactionsRes struct {
// 		e error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	ListAllTransactionsRes struct {
// 		e error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	RenameAccountRes struct {
// 		e error
// 		Res None
// 	}
// 	WalletIsLockedRes struct {
// 		e error
// 		Res bool
// 	}
// 	DropWalletHistoryRes struct {
// 		e error
// 		Res string
// 	}
// )
