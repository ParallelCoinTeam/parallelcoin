package vue

import (
	"encoding/base64"
	"fmt"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/core"

	"github.com/zserge/webview"

	"git.parallelcoin.io/dev/pod/cmd/gui/vue/alert"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/conf"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/comp/lib"
	"git.parallelcoin.io/dev/pod/cmd/gui/vue/db"
	"git.parallelcoin.io/dev/pod/pkg/conte"
)

func GetDuoVUE(cx *conte.Xt) *DuoVUE {
	v := &DuoVUE{}
	v.cx = cx
	// v.Cfg = DuoGuiCfg{}
	// cr := DuoVUEcore{}
	// cr.mod.DuoGuiItem.Name = "Sys"
	// cr.mod.DuoGuiItem.Slug = "sys"
	// cr.mod.DuoGuiItem.Version = "0.0.1"
	// cr.mod.DuoGuiItem.CompType = "System"
	// v.Core = cr
	//v.Status.GetDuoVUEstatus()
	// v.initVueSystem()
	v.db.DuoVueDbInit(v.cx.DataDir)
	// v.getVUEduoNode()
	v.Components = comp.Components(v.db)
	v.Config = GetCoreCofig(v.cx)
	return v
}

func RunVue(w webview.WebView, v DuoVUE) {
	// cx.RPCServer.Cfg.SyncMgr.IsCurrent()
	var err error
	_, err = w.Bind("system", &DuoVUE{
		cx:         v.cx,
		db:         v.db,
		Components: v.Components,
		Config:     v.Config,
		Repo:       conf.GetParallelCoinRepo,
		Icons:      lib.GetIcons(),
	})

	// Db inteface
	_, err = w.Bind("db", &db.DuoVUEdb{})
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	_, err = w.Bind("alert", &alert.DuoVUEalert{})
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}

	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	// Css
	injectCss(w, v.db)

	// Js
	evalJs(w, v.db)
	//fmt.Println("MIkaaaaaaaaaa:", v.Icons)

}

func evalJs(w webview.WebView, d db.DuoVUEdb) {
	// vue
	vueLib, err := base64.StdEncoding.DecodeString(lib.GetLibVue)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}
	err = w.Eval(string(vueLib))
	// ej2
	getEj2Vue, err := base64.StdEncoding.DecodeString(lib.GetEjs2Vue)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}
	err = w.Eval(string(getEj2Vue))
	// libs
	for _, lib := range lib.GetLibs() {
		lb, err := base64.StdEncoding.DecodeString(string(lib))
		err = w.Eval(string(lb))
		if err != nil {
			fmt.Printf("Error decoding string: %s ", err.Error())
			return
		}
	}
	// for _, js := range t.Data["js"] {
	// 	err = w.Eval(string(js))
	// }
	err = w.Eval(core.CoreJs(d))
	if err != nil {
		fmt.Println("error binding to webview:", err)
	}
	//fmt.Println("MIkaaaaaaaaaa:", core.CoreJs(d))
}

func injectCss(w webview.WebView, d db.DuoVUEdb) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	w.InjectCSS(string(lib.GetMaterial))
	// comp
	for _, c := range comp.Components(d) {
		w.InjectCSS(string(c.Css))
	}
	// Core Css
	w.InjectCSS(string(comp.GetCoreCss))
}
