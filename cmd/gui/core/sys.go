package core

import (
	"encoding/base64"
	"fmt"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/gui/vue/lib"
	"github.com/p9c/pod/cmd/gui/vue/lib/css"
	"github.com/p9c/pod/cmd/gui/vue/lib/html"
	"github.com/p9c/pod/cmd/gui/vue/pnl"
	"github.com/p9c/pod/pkg/conte"
	"github.com/robfig/cron"

	"github.com/zserge/webview"
	"time"

	"github.com/p9c/pod/cmd/gui/db"
)

const (
	windowWidth  = 960
	windowHeight = 800
)

func MountDuOS(cx *conte.Xt, cr *cron.Cron) (d *DuOS) {
	d = &DuOS{Cx: cx, Cr: cr}
	d.Data = mod.DuOSdata{
		Status:               d.GetDuOSstatus(),
		TransactionsExcerpts: d.GetTransactionsExcertps(),
		Addressbook:          d.GetAddressBook(),
	}
	d.Config = GetCoreCofig(d.Cx)
	d.Wv = webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "ParallelCoin - DUO - True Story",
		Resizable:              false,
		Debug:                  true,
		URL:                    `data:text/html,` + string(html.VUEHTML(html.VUEx(lib.VUElogo(), html.VUEheader(), html.VUEnav(lib.ICO()), html.VUEoverview()), css.CSS(css.ROOT(), css.GRID(), css.NAV()))),
		ExternalInvokeCallback: d.HandleRPC,
	})

	//d.Components = comp.Components(d.db)
	d.db.DuoVueDbInit(d.Cx.DataDir)

	return d
}

func RunDuOS(d DuOS) {
	var err error
	// eval vue lib
	evalB(&d, string(lib.VUE))
	// eval vfg lib
	evalB(&d, string(lib.VFG))
	evalB(&d, string(lib.VEJ))
	// eval ejs lib

	// init duOS variable
	evalD(&d, duOSjs(d.db))

	// init alert variable
	a := DuOSalert{
		Time:      time.Now(),
		Title:     "Welcome",
		Message:   "to ParallelCoin",
		AlertType: "success",
	}

	d.Alert = a
	_, err = d.Wv.Bind("duos", &DuOS{
		Cx: d.Cx,
		db: d.db,
		//Components: d.Components,
		Config: d.Config,
		//Repo:       conf.GetParallelCoinRepo,
		Data: d.Data,
	})
	// init duOS status
	d.Render("status", d.GetDuOSstatus())

	// Db inteface
	_, err = d.Wv.Bind("db", &db.DuOSdb{})
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	//injectCss(d)

	evalD(&d, CoreJs(d.db))

	evalD(&d, pnl.PanelsJs(pnl.Panels(d.db)))

	d.Cr.AddFunc("@every 1s", func() {
		d.Wv.Dispatch(func() {
			d.Render("status", d.GetDuOSstatus())
			//d.Wv.Eval(`document.getElementById("balanceValue").innerHTML = "` + d.GetDuOSstatus().Balance.Balance + `";`)
			//d.Wv.Eval(`document.getElementById("unconfirmedalue").innerHTML = "` + d.GetDuOSstatus().Balance.Unconfirmed + `";`)
			//d.Wv.Eval(`document.getElementById("txsnumberValue").innerHTML = "` + fmt.Sprint(d.GetDuOSstatus().TxsNumber) + `";`)
			//d.Wv.Eval(`document.getElementById("txsnumberValue").innerHTML = "` + fmt.Sprint(d.GetDuOSstatus().TxsNumber) + `";`)
			//
			//b, err := json.Marshal(tasks)
			//if err == nil {
			//	w.Eval(fmt.Sprintf("rpc.render(%s)", string(b)))
			//}

		})
	})
	//// Css

	// Js

	//d.GetPeerInfo()

	//	d.Cr.AddFunc("@every 1s", func() {
	//		d.Wv.Dispatch(func() {
	//			//d.Render("status", d.GetDuOSstatus())
	//		})
	//	})
	//	d.Cr.AddFunc("@every 10s", func() {
	//		d.Wv.Dispatch(func() {
	//			//d.Render("alert", d.GetPeerInfo())
	//		})
	//	})
	//
}

func injectCss(d DuOS) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	//d.Wv.InjectCSS(string(lib.GetMaterial))

	// Core Css
	//d.Wv.InjectCSS(string(css.CSS(css.ROOT(), css.GRID(), css.NAV())))
	//
	//for _, alj := range comp.Apps(d.db) {
	//	d.Wv.InjectCSS(string(alj.Css))
	//}
	// comp
	//for _, c := range comp.Components(d.db) {
	//	d.Wv.InjectCSS(string(c.Css))
	//}
}

func evalD(d *DuOS, l string) {
	err := d.Wv.Eval(l)
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
}

func evalB(d *DuOS, l string) {
	lib, err := base64.StdEncoding.DecodeString(l)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
	}
	evalD(d, string(lib))
}
