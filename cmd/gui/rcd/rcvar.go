package rcd

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"time"
)

var (
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
)

type RcVar struct {
	Boot             *Boot
	Alert            models.DuoUIalert
	Status           models.DuoUIstatus
	Hashes           int64
	NetHash          int64
	BlockHeight      int32
	BestBlock        string
	Difficulty       float64
	BlockCount       int64
	NetworkLastBlock int32
	ConnectionCount  int32

	Balance         string
	Unconfirmed     string
	TxsNumber       int
	CommandsHistory models.DuoUIcommandsHistory
	Transactions    models.DuoUItransactions
	Txs             models.DuoUItransactionsExcerpts
	LastTxs         models.DuoUItransactions
	Settings        models.DuoUIsettings

	Sent              bool
	IsNotificationRun bool
	Localhost         models.DuoUIlocalHost

	Uptime int
	Peers  []*btcjson.GetPeerInfoResult `json:"peers"`
	Blocks []models.DuoUIblock
	screen string `json:"screen"`
}

type Boot struct {
	IsBoot     bool   `json:"boot"`
	IsFirstRun bool   `json:"firstrun"`
	IsBootMenu bool   `json:"menu"`
	IsBootLogo bool   `json:"logo"`
	IsLoading  bool   `json:"loading"`
	IsScreen   string `json:"screen"`
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
	b := Boot{
		IsBoot:     true,
		IsFirstRun: false,
		IsBootMenu: false,
		IsBootLogo: false,
		IsLoading:  false,
		IsScreen:   "",
	}
	return &RcVar{
		Boot:             &b,
		Alert:            models.DuoUIalert{},
		Status:           models.DuoUIstatus{},
		Hashes:           0,
		NetHash:          0,
		BlockHeight:      0,
		BestBlock:        "",
		Difficulty:       0,
		BlockCount:       0,
		NetworkLastBlock: 0,
		ConnectionCount:  0,
		Balance:          "",
		Unconfirmed:      "",
		TxsNumber:        0,
		CommandsHistory: models.DuoUIcommandsHistory{
			Commands: []models.DuoUIcommand{
				models.DuoUIcommand{
					ComID:    "input",
					Category: "input",
					Time:     time.Now(),

					//Out: input(duo),
				},
			},
			CommandsNumber: 1,
		},
		Transactions:      models.DuoUItransactions{},
		Txs:               models.DuoUItransactionsExcerpts{},
		LastTxs:           models.DuoUItransactions{},
		Sent:              false,
		IsNotificationRun: true,
		Localhost:         models.DuoUIlocalHost{},
		screen:            "",
	}
}

//func input(duo models.DuoUI) func() {
//	return func() {
//		e := duo.DuoUItheme.DuoUIeditor("Run command", "Run txt")
//		e.Font.Style = text.Regular
//		e.Font.Size = unit.Dp(16)
//		e.Layout(duo.DuoUIcontext, consoleInputField)
//		for _, e := range consoleInputField.Events(duo.DuoUIcontext) {
//			if e, ok := e.(widget.SubmitEvent); ok {
//				rc.CommandsHistory.Commands = append(rc.CommandsHistory.Commands, models.DuoUIcommand{
//					CommandsID: e.Text,
//					Time:       time.Time{},
//				})
//				consoleInputField.SetText("")
//			}
//		}
//	}
//}
