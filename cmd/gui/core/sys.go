package core

import (
	"encoding/base64"
	"fmt"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/cmd/gui/vue/lib"
	"github.com/p9c/pod/cmd/gui/vue/lib/css"
	"github.com/p9c/pod/cmd/gui/vue/lib/html"
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

func MountDuOS(cx *conte.Xt, cr *cron.Cron) *DuOS {
	d := &DuOS{Cx: cx, Cr: cr}
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
	a := DuOSalert{
		Time:      time.Now(),
		Title:     "Welcome",
		Message:   "to ParallelCoin",
		AlertType: "success",
	}
	dD := mod.DuOSdata{
		Status:               d.GetDuOSstatus(),
		TransactionsExcerpts: d.GetTransactionsExcertps(),
		Addressbook:          d.GetAddressBook(),
	}
	d.Alert = a
	_, err = d.Wv.Bind("duos", &DuOS{
		Cx: d.Cx,
		db: d.db,
		//Components: d.Components,
		Config: d.Config,
		//Repo:       conf.GetParallelCoinRepo,
		Data: dD,
	})
	// Db inteface
	_, err = d.Wv.Bind("db", &db.DuOSdb{})
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	//injectCss(d)
	evalJs(d)

	//
	//
	//d.Cr.AddFunc("@every 1s", func() {
	//	d.Wv.Dispatch(func() {
	//		//d.Render("status", d.GetDuOSstatus())
	//		d.Wv.Eval(`document.getElementById("balance").innerHTML = "` + d.GetDuOSstatus().Balance.Balance + `";`)
	//	})
	//})
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

func evalJs(dV DuOS) {
	// vue
	vueLib, err := base64.StdEncoding.DecodeString(lib.VUE)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}
	err = dV.Wv.Eval(string(vueLib))

	// ej2
	//getAMP, err := base64.StdEncoding.DecodeString(lib.AMP)
	//if err != nil {
	//	fmt.Printf("Error decoding string: %s ", err.Error())
	//	return
	//}
	//err = d.Wv.Eval(string(getAMP))

	// libs
	//for _, lb := range lib.AMPLIBS {
	//	l, err := base64.StdEncoding.DecodeString(string(lb))
	//	err = d.Wv.Eval(string(l))
	//	if err != nil {
	//		fmt.Printf("Error decoding string: %s ", err.Error())
	//		return
	//	}
	//}

	//err = dV.Wv.Eval(core.CoreHeadJs)
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//
	//err = dV.Wv.Eval(core.CompLoopJs(dV.db))
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//
	//err = dV.Wv.Eval(core.AppsLoopJs(dV.db))
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//
	//err = dV.Wv.Eval(core.CoreJs)
	//if err != nil {
	//	fmt.Println("error binding to webview:", err)
	//}
	//fmt.Println("MIkaaaaaaaaaa:", core.CoreJs(d))
}

func injectCss(dV DuOS) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	//dV.Wv.InjectCSS(string(lib.GetMaterial))

	// Core Css
	//dV.Wv.InjectCSS(string(css.CSS(css.ROOT(), css.GRID(), css.NAV())))
	//
	//for _, alj := range comp.Apps(dV.db) {
	//	dV.Wv.InjectCSS(string(alj.Css))
	//}
	// comp
	//for _, c := range comp.Components(d.db) {
	//	dV.Wv.InjectCSS(string(c.Css))
	//}
}
