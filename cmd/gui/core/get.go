package core

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/amp/lib"
	"github.com/p9c/pod/cmd/gui/amp/lib/css"
	"github.com/p9c/pod/cmd/gui/amp/lib/html"
	"github.com/p9c/pod/cmd/gui/mod"
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
		URL:                    `data:text/html,` + string(html.AMPHTML("ssss", html.AMPx(lib.AMPlogo(), html.AMPheader(), html.AMPnav()), lib.AMPlib(), css.CSS(css.AMProot(), css.AMPgrid()), lib.AMPsw())),
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
		Icons:                lib.GetIcons(),
	}
	d.Alert = a
	_, err = d.Wv.Bind("system", &DuOS{
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
}

//
//// Css
//injectCss(d)

// Js

//d.GetPeerInfo()

//	evalJs(d)
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
//}

//func evalJs(d DuOS) {
//	// vue
//	vueLib, err := base64.StdEncoding.DecodeString(lib.AMP)
//	if err != nil {
//		fmt.Printf("Error decoding string: %s ", err.Error())
//		return
//	}
//	err = d.Wv.Eval(string(vueLib))
//	// ej2
//	getAMP, err := base64.StdEncoding.DecodeString(lib.AMP)
//	if err != nil {
//		fmt.Printf("Error decoding string: %s ", err.Error())
//		return
//	}
//	err = d.Wv.Eval(string(getAMP))
//	// libs
//	for _, lb := range lib.AMPLIBS {
//		l, err := base64.StdEncoding.DecodeString(string(lb))
//		err = d.Wv.Eval(string(l))
//		if err != nil {
//			fmt.Printf("Error decoding string: %s ", err.Error())
//			return
//		}
//	}
//
//}

//func injectCss(dV DuOS) {
//	// material
//	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
//	// if err != nil {
//	// 	fmt.Printf("Error decoding string: %s ", err.Error())
//	// 	return
//	// }
//	d.Wv.InjectCSS(string(lib.GetMaterial))
//
//	// Core Css
//	d.Wv.InjectCSS(string(comp.GetCoreCss))
//
//	for _, alj := range comp.Apps(d.db) {
//		d.Wv.InjectCSS(string(alj.Css))
//	}
//	// comp
//	for _, c := range comp.Components(d.db) {
//		d.Wv.InjectCSS(string(c.Css))
//	}
//}
