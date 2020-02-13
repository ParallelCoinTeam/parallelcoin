package rcd

import (
	"sync"
	"time"
	
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var (
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
)

type RcVar struct {
	Boot             *Boot
	Events           chan Event
	UpdateTrigger    chan struct{}
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
	Balance          string
	Unconfirmed      string
	TxsNumber        int
	CommandsHistory  models.DuoUIcommandsHistory
	Transactions     models.DuoUItransactions
	Txs              models.DuoUItransactionsExcerpts
	LastTxs          models.DuoUItransactions
	Settings         models.DuoUIsettings
	Sent             bool
	ShowDialog       bool
	Toasts           []func()
	Localhost        models.DuoUIlocalHost
	Uptime           int
	Peers            []*btcjson.GetPeerInfoResult `json:"peers"`
	Blocks           []models.DuoUIblock
	screen           string `json:"screen"`
	mutex            sync.Mutex
}

type Boot struct {
	IsBoot     bool   `json:"boot"`
	IsFirstRun bool   `json:"firstrun"`
	IsBootMenu bool   `json:"menu"`
	IsBootLogo bool   `json:"logo"`
	IsLoading  bool   `json:"loading"`
	IsScreen   string `json:"screen"`
}

// type rcVar interface {
//	GetDuoUItransactions(sfrom, count int, cat string)
//	GetDuoUIbalance()
//	GetDuoUItransactionsExcerpts()
//	DuoSend(wp string, ad string, am float64)
//	GetDuoUItatus()
//	PushDuoUIalert(t string, m interface{}, at string)
//	GetDuoUIblockHeight()
//	GetDuoUIblockCount()
//	GetDuoUIconnectionCount()
// }

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
					
					// Out: input(duo),
				},
			},
			CommandsNumber: 1,
		},
		Transactions: models.DuoUItransactions{},
		Txs:          models.DuoUItransactionsExcerpts{},
		LastTxs:      models.DuoUItransactions{},
		Sent:         false,
		ShowDialog:   true,
		Localhost:    models.DuoUIlocalHost{},
		screen:       "",
	}
}

// func input(duo models.DuoUI) func() {
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
// }

func (rc *RcVar) RCtoast() {
	// tickerChannel := time.NewTicker(3 * time.Second)
	// go func() {
	//	for {
	//		select {
	//		case <-tickerChannel.C:
	//			for i := range rc.Toasts {
	//				log.DEBUG("RRRRRR")
	//				if i < len(rc.Toasts)-1 {
	//					copy(rc.Toasts[i:], rc.Toasts[i+1:])
	//				}
	//				rc.Toasts[len(rc.Toasts)-1] = nil // or the zero value of T
	//				rc.Toasts = rc.Toasts[:len(rc.Toasts)-1]				}
	//		}
	//	}
	// }()
	// time.Sleep(6 * time.Second)
}
