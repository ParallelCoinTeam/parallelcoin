package duoui

type RcVar struct {
	Alert DuOSalert
	Status       DuOStatus
	Hashes       int64
	NetHash      int64
	BlockHeight  int32
	BestBlock    string
	Difficulty   float64
	BlockCount   int64
	NetLastBlock int32
	Connections  int32

	Balance     string
	Unconfirmed string
	TxsNumber   int
	Transactions DuOStransactions
	Txs          DuOStransactionsExcerpts
	LastTxs      DuOStransactions

	Sent       bool
	IsFirstRun bool
	Localhost  DuOSlocalHost

	screen string `json:"screen"`
}

//type rcVar interface {
//	GetDuOStransactions(sfrom, count int, cat string)
//	GetDuOSbalance()
//	GetDuOStransactionsExcerpts()
//	DuoSend(wp string, ad string, am float64)
//	GetDuOStatus()
//	PushDuOSalert(t string, m interface{}, at string)
//	GetDuOSblockHeight()
//	GetDuOSblockCount()
//	GetDuOSconnectionCount()
//}


func RcInit() *RcVar{
	return &RcVar{
		Alert:        DuOSalert{},
		Status:       DuOStatus{},
		Hashes:       0,
		NetHash:      0,
		BlockHeight:  0,
		BestBlock:    "",
		Difficulty:   0,
		BlockCount:   0,
		NetLastBlock: 0,
		Connections:  0,
		Balance:      "",
		Unconfirmed:  "",
		TxsNumber:    0,
		Transactions: DuOStransactions{},
		Txs:          DuOStransactionsExcerpts{},
		LastTxs:      DuOStransactions{},
		Sent:         false,
		IsFirstRun:   false,
		Localhost:    DuOSlocalHost{},
		screen:       "",
	}
}