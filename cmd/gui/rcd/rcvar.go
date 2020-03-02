package rcd

import (
	"github.com/p9c/pod/cmd/gui/controller"
	"github.com/p9c/pod/cmd/gui/theme"
	"time"

	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type RcVar struct {
	cx              *conte.Xt
	db              *DuoUIdb
	Boot            *Boot
	Events          chan Event
	UpdateTrigger   chan struct{}
	Status          *model.DuoUIstatus
	Dialog          *model.DuoUIdialog
	Log             *model.DuoUIlog
	CommandsHistory *model.DuoUIcommandsHistory

	Settings  *model.DuoUIsettings
	Sent      bool
	Toasts    []model.DuoUItoast
	Localhost model.DuoUIlocalHost
	Uptime    int
	Peers     []*btcjson.GetPeerInfoResult `json:"peers"`

	AddressBook model.DuoUIaddressBook
	ShowPage    string
	CurrentPage *theme.DuoUIpage
	// NodeChan   chan *rpc.Server
	// WalletChan chan *wallet.Wallet
	Explorer *model.Explorer
	Quit     chan struct{}
	Ready    chan struct{}
	IsReady  bool
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

func RcInit(cx *conte.Xt) (r *RcVar) {
	b := Boot{
		IsBoot:     false,
		IsFirstRun: false,
		IsBootMenu: false,
		IsBootLogo: false,
		IsLoading:  false,
		IsScreen:   "",
	}
	// d := models.DuoUIdialog{
	//	Show:   true,
	//	Ok:     func() { r.Dialog.Show = false },
	//	Cancel: func() { r.Dialog.Show = false },
	//	Title:  "Dialog!",
	//	Text:   "Dialog text",
	// }
	l := new(model.DuoUIlog)

	r = &RcVar{
		cx:   cx,
		Boot: &b,
		Status: &model.DuoUIstatus{
			Node: &model.NodeStatus{},
			Wallet: &model.WalletStatus{
				WalletVersion: make(map[string]btcjson.VersionResult),
				Transactions:  &model.DuoUItransactions{},
				Txs:           &model.DuoUItransactionsExcerpts{},
				LastTxs:       &model.DuoUItransactions{},
			},
			Kopach: &model.KopachStatus{},
		},
		Dialog:   &model.DuoUIdialog{},
		Settings: settings(cx),
		Log:      l,
		CommandsHistory: &model.DuoUIcommandsHistory{
			Commands: []model.DuoUIcommand{
				model.DuoUIcommand{
					ComID:    "input",
					Category: "input",
					Time:     time.Now(),

					// Out: input(duo),
				},
			},
			CommandsNumber: 1,
		},
		Sent:      false,
		Localhost: model.DuoUIlocalHost{},
		ShowPage:  "OVERVIEW",
		Explorer: &model.Explorer{
			PerPage: &controller.DuoUIcounter{
				Value:           20,
				OperateValue:    1,
				From:            0,
				To:              50,
				CounterIncrease: new(controller.Button),
				CounterDecrease: new(controller.Button),
				CounterReset:    new(controller.Button),
			},
			Page: &controller.DuoUIcounter{
				Value:           0,
				OperateValue:    1,
				From:            0,
				To:              50,
				CounterIncrease: new(controller.Button),
				CounterDecrease: new(controller.Button),
				CounterReset:    new(controller.Button),
			},
			Blocks:      []model.DuoUIblock{},
			SingleBlock: btcjson.GetBlockVerboseResult{},
		},
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}, 1),
	}
	return
}
