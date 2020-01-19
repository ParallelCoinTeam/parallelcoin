package rcd

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type RcVar struct {
	Alert            models.DuoUIalert
	Status           models.DuoUIstatus
	Hashes           int64
	NetHash          int64
	BlockHeight      int32
	BestBlock        string
	Difficulty       float64
	BlockCount       int64
	NetworkLastBlock int32
	ConnectionCount int32

	Balance      string
	Unconfirmed  string
	TxsNumber    int
	Transactions models.DuoUItransactions
	Txs          models.DuoUItransactionsExcerpts
	LastTxs      models.DuoUItransactions
	Settings     models.DuoUIsettings

	Sent       bool
	IsFirstRun bool
	Localhost  models.DuoUIlocalHost

	Uptime int
	Peers  []*btcjson.GetPeerInfoResult `json:"peers"`
	Blocks []models.DuoUIblock
	screen string `json:"screen"`
}


//type rcVar interface {
//	GetDuoUItransactions(sfrom, count int, cat string)
//	GetDuoUIbalance()
//	GetDuoUItransactionsExcerpts()
//	DuoSend(wp string, ad string, am float64)
//	GetDuoUItatus()
//	PushDuoUIalert(t string, m interface{}, at string)
//	GetDuoUIblockHeight()
//	GetDuoUIblockCount()
//	GetDuoUIconnectionCount()
//}

func RcInit() *RcVar {
	return &RcVar{
		Alert:            models.DuoUIalert{},
		Status:           models.DuoUIstatus{},
		Hashes:           0,
		NetHash:          0,
		BlockHeight:      0,
		BestBlock:        "",
		Difficulty:       0,
		BlockCount:       0,
		NetworkLastBlock: 0,
		ConnectionCount: 0,
		Balance:          "",
		Unconfirmed:      "",
		TxsNumber:        0,
		Transactions:     models.DuoUItransactions{},
		Txs:              models.DuoUItransactionsExcerpts{},
		LastTxs:          models.DuoUItransactions{},
		Sent:             false,
		IsFirstRun:       false,
		Localhost:        models.DuoUIlocalHost{},
		screen:           "",
	}
}
