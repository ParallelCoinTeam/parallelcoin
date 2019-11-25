package gui

import (
	"encoding/json"
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
func
evalJs(r *rcvar) (err error) {
	err = r.evalJsFile([]string{
		"libs/js/vue.js",
		"libs/js/ej2-vue.min.js",
		"libs/js/vfg.js",
		"js/vue/duos.js",
		"js/vue/ico/logo.js",
		"js/vue/ico/overview.js",
		"js/vue/ico/history.js",
		"js/vue/ico/addressbook.js",
		"js/vue/ico/explorer.js",
		"js/vue/ico/settings.js",
		"js/vue/panels/balance.js",
		"js/vue/panels/send.js",
		"js/vue/panels/peers.js",
		"js/vue/panels/status.js",
		"js/vue/panels/networkhashrate.js",
		"js/vue/panels/localhashrate.js",
		"js/vue/panels/history.js",
		"js/vue/panels/latestxs.js",
		"js/vue/panels/addressbook.js",
		"js/vue/panels/settings.js",
		"js/vue/pages/overview.js",
		"js/vue/pages/history.js",
		"js/vue/pages/addressbook.js",
		"js/vue/pages/explorer.js",
		"js/vue/pages/settings.js",
		"js/vue/layout/header.js",
		"js/vue/layout/nav.js",
		"js/vue/layout/xorg.js",
		"js/vue/dui.js"})
	return
}
func
injectCss(r *rcvar) {
	r.injectCssFile([]string{
		"libs/css/material.css",
		"css/theme/root.css",
		"css/theme/colors.css",
		"css/theme/grid.css",
		"css/theme/helpers.css",
		"css/duistyle.css",
		"css/dui.css"})
}
func
(r *rcvar) evalJsFile(fls []string) (err error) {
	for _, f := range fls {
		err = r.cx.Gui.Wv.Eval(getFile(f, *r.cx.Gui.Fs))
		if err != nil {
			log.ERROR("error binding to webview:", err)
		}
	}
	return
}
func
(r *rcvar) injectCssFile(fls []string) {
	for _, f := range fls {
		r.cx.Gui.Wv.InjectCSS(getFile(f, *r.cx.Gui.Fs))
	}
}



func
(r *rcvar) Render(cmd string, data interface{}) (err error) {
	var b []byte
	b, err = json.Marshal(data)
	if err == nil {
		r.cx.Gui.Wv.Dispatch(func() {
			//r.cx.Gui.Wv.Eval(cmd + "=" + string(b) + ";")
			r.cx.Gui.Wv.Eval("rcvar." + cmd + "=" + string(b) + ";")
			//log.INFO("WORKS:VAR->", cmd + "=" + string(b) + ";")
		})
	}
	return
}
