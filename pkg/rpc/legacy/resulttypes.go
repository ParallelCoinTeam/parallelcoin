package legacy

//
// import (
// 	"github.com/stalker-loki/pod/pkg/rpc/btcjson"
// )
//
// // errors are returned as *btcjson.RPCError type
// type (
// 	None                   struct{}
// 	AddMultiSigAddressRes struct {
// 		Err error
// 		Res string
// 	}
// 	CreateMultiSigRes struct {
// 		Err error
// 		Res btcjson.CreateMultiSigResult
// 	}
// 	DumpPrivKeyRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetAccountRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetAccountAddressRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetAddressesByAccountRes struct {
// 		Err error
// 		Res []string
// 	}
// 	GetBalanceRes struct {
// 		Err error
// 		Res float64
// 	}
// 	GetBestBlockHashRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetBlockCountRes struct {
// 		Err error
// 		Res int32
// 	}
// 	GetInfoRes struct {
// 		Err error
// 		Res btcjson.InfoWalletResult
// 	}
// 	GetNewAddressRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetRawChangeAddressRes struct {
// 		Err error
// 		Res string
// 	}
// 	GetReceivedByAccountRes struct {
// 		Err error
// 		Res float64
// 	}
// 	GetReceivedByAddressRes struct {
// 		Err error
// 		Res float64
// 	}
// 	GetTransactionRes struct {
// 		Err error
// 		Res btcjson.GetTransactionResult
// 	}
// 	HelpWithChainRPCRes struct {
// 		Err error
// 		Res string
// 	}
// 	HelpNoChainRPCRes struct {
// 		Err error
// 		Res string
// 	}
// 	ImportPrivKeyRes struct {
// 		Err error
// 		Res None
// 	}
// 	KeypoolRefillRes struct {
// 		Err error
// 		Res None
// 	}
// 	ListAccountsRes struct {
// 		Err error
// 		Res map[string]float64
// 	}
// 	ListLockUnspentRes struct {
// 		Err error
// 		Res []btcjson.TransactionInput
// 	}
// 	ListReceivedByAccountRes struct {
// 		Err error
// 		Res []btcjson.ListReceivedByAccountResult
// 	}
// 	ListReceivedByAddressRes struct {
// 		Err error
// 		Res btcjson.ListReceivedByAddressResult
// 	}
// 	ListSinceBlockRes struct {
// 		Err error
// 		Res btcjson.ListSinceBlockResult
// 	}
// 	ListTransactionsRes struct {
// 		Err error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	ListUnspentRes struct {
// 		Err error
// 		Res []btcjson.ListUnspentResult
// 	}
// 	LockUnspentRes struct {
// 		Err error
// 		Res bool
// 	}
// 	SendFromRes struct {
// 		Err error
// 		Res string
// 	}
// 	SendManyRes struct {
// 		Err error
// 		Res string
// 	}
// 	SendToAddressRes struct {
// 		Err error
// 		Res string
// 	}
// 	SetTxFeeRes struct {
// 		Err error
// 		Res bool
// 	}
// 	SignMessageRes struct {
// 		Err error
// 		Res string
// 	}
// 	SignRawTransactionRes struct {
// 		Err error
// 		Res btcjson.SignRawTransactionResult
// 	}
// 	ValidateAddressRes struct {
// 		Err error
// 		Res btcjson.ValidateAddressWalletResult
// 	}
// 	VerifyMessageRes struct {
// 		Err error
// 		Res bool
// 	}
// 	WalletLockRes struct {
// 		Err error
// 		Res None
// 	}
// 	WalletPassphraseRes struct {
// 		Err error
// 		Res None
// 	}
// 	WalletPassphraseChangeRes struct {
// 		Err error
// 		Res None
// 	}
// 	CreateNewAccountRes struct {
// 		Err error
// 		Res None
// 	}
// 	GetBestBlockRes struct {
// 		Err error
// 		Res btcjson.GetBestBlockResult
// 	}
// 	GetUnconfirmedBalanceRes struct {
// 		Err error
// 		Res float64
// 	}
// 	ListAddressTransactionsRes struct {
// 		Err error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	ListAllTransactionsRes struct {
// 		Err error
// 		Res []btcjson.ListTransactionsResult
// 	}
// 	RenameAccountRes struct {
// 		Err error
// 		Res None
// 	}
// 	WalletIsLockedRes struct {
// 		Err error
// 		Res bool
// 	}
// 	DropWalletHistoryRes struct {
// 		Err error
// 		Res string
// 	}
// )
