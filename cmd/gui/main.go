// +build !nogui

package gui

import (
	"sync"
	"sync/atomic"

	"github.com/zserge/webview"

	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue"
	"github.com/parallelcointeam/parallelcoin/cmd/gui/vue/comp"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/util/interrupt"
)

const (
	windowWidth  = 1440
	windowHeight = 900
)

var getWebView = func(v vue.DuoVUE, t string) webview.WebView {
	w := webview.New(webview.Settings{
		Width:                  windowWidth,
		Height:                 windowHeight,
		Title:                  "ParallelCoin - DUO - True Story",
		Resizable:              false,
		Debug:                  true,
		URL:                    `data:text/html,` + t,
		ExternalInvokeCallback: v.HandleRPC,
	})
	return w
}

func Main(cx *conte.Xt, wg *sync.WaitGroup) {
	WARN("starting gui")
	v := vue.GetDuoVUE(cx)
	w := getWebView(*v, comp.GetAppHtml)
	cleaned := &atomic.Value{}
	cleaned.Store(false)
	cleanup := func() {
		if !cleaned.Load().(bool) {
			cleaned.Store(true)
			DEBUG("terminating webview")
			w.Terminate()
			interrupt.Request()
			DEBUG("waiting for waitgroup")
			wg.Wait()
			DEBUG("exiting webview")
			w.Exit()
		}
	}
	interrupt.AddHandler(func() {
		cleanup()
	})
	defer cleanup()
	w.Dispatch(func() {

		// dec, err := base64.StdEncoding.DecodeString(lib.GetLibVue)
		// if err != nil {
		// 	fmt.Printf("Error decoding string: %s ", err.Error())
		// 	return
		// }
		// w.Eval(string(dec))
		vue.RunVue(w, *v)

	})
	w.Run()
}
