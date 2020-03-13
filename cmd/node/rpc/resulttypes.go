package rpc

//
// type API struct {
// 	Ch     interface{}
// 	Params interface{}
// }
//
// // errors are returned as *btcjson.RPCError type
// type (
// 	None        struct{}
// 	AddNodeRes struct {
// 		Res None
// 		Err error
// 	}
// 	CreateRawTransactionRes struct {
// 		Res string
// 		Err error
// 	}
// 	DebugLevelRes struct {
// 		Res None
// 		Err error
// 	}
// 	DecodeRawTransactionRes struct {
// 		Res btcjson.TxRawDecodeResult
// 		Err error
// 	}
// 	DecodeScriptRes struct {
// 		Res btcjson.DecodeScriptResult
// 		Err error
// 	}
// 	EstimateFeeRes struct {
// 		Res float64
// 		Err error
// 	}
// 	GenerateRes struct {
// 		Res []string
// 		Err error
// 	}
// 	GetAddedNodeInfoRes struct {
// 		Res []btcjson.GetAddedNodeInfoResultAddr
// 		Err error
// 	}
// 	GetBestBlockRes struct {
// 		Res btcjson.GetBestBlockResult
// 		Err error
// 	}
// 	GetBestBlockHashRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetBlockRes struct {
// 		Res btcjson.GetBlockVerboseResult // can be string if not verbose
// 		Err error
// 	}
// 	GetBlockChainInfoRes struct {
// 		Res btcjson.GetBlockChainInfoResult
// 		Err error
// 	}
// 	GetBlockCountRes struct {
// 		Res int64
// 		Err error
// 	}
// 	GetBlockHashRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetBlockHeaderRes struct {
// 		Res btcjson.GetBlockHeaderVerboseResult // can be string if not verbose
// 		Err error
// 	}
// 	GetBlockTemplateRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetCFilterRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetCFilterHeaderRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetConnectionCountRes struct {
// 		Res int32
// 		Err error
// 	}
// 	GetCurrentNetRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetDifficultyRes struct {
// 		Res float64
// 		Err error
// 	}
// 	GetGenerateRes struct {
// 		Res bool
// 		Err error
// 	}
// 	GetHashesPerSecRes struct {
// 		Res float64
// 		Err error
// 	}
// 	GetHeadersRes struct {
// 		Res []string
// 		Err error
// 	}
// 	GetInfoRes struct {
// 		Res btcjson.InfoChainResult0
// 		Err error
// 	}
// 	GetMempoolInfoRes struct {
// 		Res btcjson.GetMempoolInfoResult
// 		Err error
// 	}
// 	GetMiningInfo0Res struct {
// 		Res btcjson.GetMiningInfoResult0
// 		Err error
// 	}
// 	GetMiningInfoRes struct {
// 		Res btcjson.GetMiningInfoResult
// 		Err error
// 	}
// 	GetNetTotalsRes struct {
// 		Res btcjson.GetNetTotalsResult
// 		Err error
// 	}
// 	GetNetworkHashPSRes struct {
// 		Res int64
// 		Err error
// 	}
// 	GetPeerInfoRes struct {
// 		Res []btcjson.GetPeerInfoResult
// 		Err error
// 	}
// 	GetRawMempoolRes struct {
// 		Res []string
// 		Err error
// 	}
// 	GetRawTransactionRes struct {
// 		Res string
// 		Err error
// 	}
// 	GetTxOutRes struct {
// 		Res btcjson.GetTxOutResult
// 		Err error
// 	}
// 	HelpRes struct {
// 		Res string
// 		Err error
// 	}
// 	NodeRes struct {
// 		Res None
// 		Err error
// 	}
// 	PingRes struct {
// 		Res None
// 		Err error
// 	}
// 	SearchRawTransactionsRes struct {
// 		Res []btcjson.SearchRawTransactionsResult
// 		Err error
// 	}
// 	SendRawTransactionRes struct {
// 		Res None
// 		Err error
// 	}
// 	SetGenerateRes struct {
// 		Res None
// 		Err error
// 	}
// 	StopRes struct {
// 		Res None
// 		Err error
// 	}
// 	RestartRes struct {
// 		Res None
// 		Err error
// 	}
// 	ResetChainRes struct {
// 		Res None
// 		Err error
// 	}
// 	SubmitBlockRes struct {
// 		Res string
// 		Err error
// 	}
// 	UptimeRes struct {
// 		Res int64
// 		Err error
// 	}
// 	ValidateAddressRes struct {
// 		Res btcjson.ValidateAddressChainResult
// 		Err error
// 	}
// 	VerifyChainRes struct {
// 		Res bool
// 		Err error
// 	}
// 	VerifyMessageRes struct {
// 		Res bool
// 		Err error
// 	}
// 	VersionRes struct {
// 		Res map[string]btcjson.VersionResult
// 		Err error
// 	}
// )
