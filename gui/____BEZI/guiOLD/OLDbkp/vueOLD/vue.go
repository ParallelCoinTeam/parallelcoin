package vueOLD

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/core"
	"github.com/robfig/cron"
	"time"

	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp/conf"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/comp/lib"
	"github.com/p9c/pod/cmd/gui/OLDbkp/vue/db"
	"github.com/p9c/pod/pkg/conte"
)

const (
	windowWidth  = 960
	windowHeight = 800
)

func GetDuoVUE(cx *conte.Xt, cr *cron.Cron) *DuoVUE {
	dV := &DuoVUE{}
	dV.cx = cx
	dV.Config = GetCoreCofig(cx)
	dV.cr = cr
	dV.Components = comp.Components(dV.db)
	dV.db.DuoVueDbInit(dV.cx.DataDir)

	return dV
}

func RunVue(dV DuoVUE) {
	var err error
	a := DuoVUEalert{
		Time:      time.Now(),
		Title:     "Welcome",
		Message:   "to ParallelCoin",
		AlertType: "success",
	}
	d := DuoVUEdata{
		Alert:                a,
		Status:               dV.GetDuoVUEstatus(),
		TransactionsExcerpts: dV.GetTransactionsExcertps(),
		Addressbook:          dV.GetAddressBook(),
	}
	_, err = dV.Web.Bind("system", &DuoVUE{
		cx:         dV.cx,
		db:         dV.db,
		Components: dV.Components,
		Config:     dV.Config,
		Repo:       conf.GetParallelCoinRepo,
		Icons:      lib.GetIcons(),
		Data:       d,
	})
	// Db inteface
	_, err = dV.Web.Bind("db", &db.DuoVUEdb{})
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	// Css
	injectCss(dV)

	// Js

	//dV.GetPeerInfo()

	evalJs(dV)
	dV.cr.AddFunc("@every 1s", func() {
		dV.Web.Dispatch(func() {
			dV.Render("status", dV.GetDuoVUEstatus())
		})
	})
	dV.cr.AddFunc("@every 10s", func() {
		dV.Web.Dispatch(func() {
			dV.Render("alert", dV.GetPeerInfo())
		})
	})

}

func evalJs(dV DuoVUE) {
	// vue
	//vueLib, err := base64.StdEncoding.DecodeString(lib.GetLibVue)
	//if err != nil {
	//	fmt.Printf("Error decoding string: %s ", err.Error())
	//	return
	//}
	//err = dV.Web.Eval(string(vueLib))
	// ej2
	//getEj2Vue, err := base64.StdEncoding.DecodeString(lib.GetEjs2Vue)
	//if err != nil {
	//	fmt.Printf("Error decoding string: %s ", err.Error())
	//	return
	//}
	//err = dV.Web.Eval(string(getEj2Vue))
	//// libs
	//for _, lib := range lib.GetLibs() {
	//	lb, err := base64.StdEncoding.DecodeString(string(lib))
	//	err = dV.Web.Eval(string(lb))
	//	if err != nil {
	//		fmt.Printf("Error decoding string: %s ", err.Error())
	//		return
	//	}
	//}
	// for _, js := range t.Data["js"] {
	// 	err = w.Eval(string(js))
	// }

	err = dV.Web.Eval(core.CoreHeadJs)
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = dV.Web.Eval(core.CompLoopJs(dV.db))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = dV.Web.Eval(core.AppsLoopJs(dV.db))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	err = dV.Web.Eval(core.CoreJs)
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	//fmt.Println("MIkaaaaaaaaaa:", core.CoreJs(d))
}

func injectCss(dV DuoVUE) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	dV.Web.InjectCSS(string(lib.GetMaterial))

	// Core Css
	dV.Web.InjectCSS(string(comp.GetCoreCss))

	for _, alj := range comp.Apps(dV.db) {
		dV.Web.InjectCSS(string(alj.Css))
	}
	// comp
	for _, c := range comp.Components(dV.db) {
		dV.Web.InjectCSS(string(c.Css))
	}
}
