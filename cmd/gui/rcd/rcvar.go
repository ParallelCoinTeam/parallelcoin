package rcd

type RcVar struct {
	Alert        DuOSalert
	Status       DuOStatus
	Hashes       int64
	NetHash      int64
	BlockHeight  int32
	BestBlock    string
	Difficulty   float64
	BlockCount   int64
	NetLastBlock int32
	Connections  int32

	Balance      string
	Unconfirmed  string
	TxsNumber    int
	Transactions DuOStransactions
	Txs          DuOStransactionsExcerpts
	LastTxs      DuOStransactions

	Sent       bool
	IsFirstRun bool
	Localhost  DuOSlocalHost

	screen string `json:"screen"`
}

type rcVar interface {
	GetDuOStransactions(sfrom, count int, cat string)
	GetDuOSbalance()
	GetDuOStransactionsExcerpts()
	DuoSend(wp string, ad string, am float64)
	GetDuOStatus()
	PushDuOSalert(t string, m interface{}, at string)
	GetDuOSblockHeight()
	GetDuOSblockCount()
	GetDuOSconnectionCount()
}
