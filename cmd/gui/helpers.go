package gui

import (
	"github.com/p9c/pod/pkg/log"
	"io/ioutil"
	"net/http"
)

func
getFile(f string, fs http.FileSystem) string {
	file, err := fs.Open(f)
	if err != nil {
		log.FATAL(err)
	}
	defer file.Close()
	body, err := ioutil.ReadAll(file)
	return string(body)
}

func evalJs(rc *rcvar) {
	var err error
	err = rc.w.Eval(getFile("libs/js/vue.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("libs/js/ej2-vue.min.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("libs/js/vfg.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/duos.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/logo.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/overview.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/history.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/addressbook.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/explorer.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/ico/settings.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/balance.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/send.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/peers.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/status.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/networkhashrate.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/localhashrate.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/history.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/latestxs.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/addressbook.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/panels/settings.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/pages/overview.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/pages/history.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/pages/addressbook.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/pages/explorer.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/pages/settings.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/layout/header.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/layout/nav.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/layout/xorg.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}

	err = rc.w.Eval(getFile("js/vue/dui.js",*rc.fs))
	if err != nil {
		log.ERROR("error binding to webview:", err)
	}



}

func injectCss(rc *rcvar) {
	// material
	// getMaterial, err := base64.StdEncoding.DecodeString(lib.GetMaterial)
	// if err != nil {
	// 	fmt.Printf("Error decoding string: %s ", err.Error())
	// 	return
	// }
	rc.w.InjectCSS(getFile("libs/css/material.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/theme/root.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/theme/colors.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/theme/grid.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/theme/helpers.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/duistyle.css",*rc.fs))
	rc.w.InjectCSS(getFile("css/dui.css",*rc.fs))

}
