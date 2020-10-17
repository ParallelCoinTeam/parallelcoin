package dap

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/pkg/gui/wallet/appdata"
	"github.com/p9c/pod/pkg/gui/wallet/dap/mod"
	"github.com/p9c/pod/pkg/gui/wallet/dap/res"
	"github.com/p9c/pod/pkg/gui/wallet/nav"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	noReturn = func(gtx C) D { return D{} }
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)
type dap struct {
	boot mod.Dap
}

func NewDap(cx *conte.Xt, title string) dap {
	//if cfg.Initial {
	//	fmt.Println("running initial setup")
	//}
	d := mod.Dap{
		Rc:   RcInit(cx),
		Apps: make(map[string]mod.Sap),
	}

	d.UI = &mod.UserInterface{
		Theme: theme.NewTheme(),
		//mob:   make(chan bool),
	}

	d.UI.Window = app.NewWindow(
		app.Size(unit.Dp(1024), unit.Dp(800)),
		app.Title(title),
	)

	n := &nav.Navigation{
		Name:         "Navigacion",
		Bg:           d.UI.Theme.Colors["NavBg"],
		ItemIconSize: unit.Px(24),
	}
	d.UI.N = n

	s := &mod.Settings{
		Dir: appdata.Dir("dap", false),
	}
	d.S = s

	d.UI.R = res.Resposnsivity(0, 0)

	return dap{boot: d}
}

func (d *dap) NewSap(s mod.Sap) {
	d.boot.Apps[s.Title] = s
	return
}
func (d *dap) BOOT() *mod.Dap {
	return &d.boot
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func RcInit(cx *conte.Xt) (r *mod.RcVar) {
	b := mod.Boot{
		IsBoot:     true,
		IsFirstRun: false,
		IsBootMenu: false,
		IsBootLogo: false,
		IsLoading:  false,
	}
	// d := models.DuoUIdialog{
	//	Show:   true,
	//	Ok:     func() { r.Dialog.Show = false },
	//	Cancel: func() { r.Dialog.Show = false },
	//	Title:  "Dialog!",
	//	Text:   "Dialog text",
	// }
	//l := new(model.DuoUIlog)

	r = &mod.RcVar{
		Cx: cx,
		//db:          new(DuoUIdb),
		Boot: &b,
		//AddressBook: new(model.DuoUIaddressBook),
		//Status: &model.DuoUIstatus{
		//	Node: &model.NodeStatus{},
		//	Wallet: &model.WalletStatus{
		//		WalletVersion: make(map[string]btcjson.VersionResult),
		//		LastTxs:       &model.DuoUItransactionsExcerpts{},
		//	},
		//	Kopach: &model.KopachStatus{},
		//},
		//Dialog:   &model.DuoUIdialog{},
		//Settings: settings(cx),
		//Log:      l,
		Quit:  make(chan struct{}),
		Ready: make(chan struct{}),
	}
	//r.db.DuoUIdbInit(r.cx.DataDir)
	return
}
