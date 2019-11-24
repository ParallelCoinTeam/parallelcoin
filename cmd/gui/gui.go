package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"net/url"
)

func WalletGUI(cx *conte.Xt) (err error) {

	cx.Gui.Wv = webview.New(webview.Settings{
		Width:  1024,
		Height: 760,
		Debug:  true,
		Resizable: false,
		Title:  "ParallelCoin - DUO - True Story",
		URL:    "data:text/html," + url.PathEscape(getFile("vue.html", *cx.Gui.Fs)),
	})
	cx.Gui.Wv.SetColor(68, 68, 68, 255)

	_, err = cx.Gui.Wv.Bind("alert", &DuOSalert{})
	_, err = cx.Gui.Wv.Bind("balance", &DuOSbalance{})
	_, err = cx.Gui.Wv.Bind("lastxs", &DuOStransactions{})
	_, err = cx.Gui.Wv.Bind("blockcount", &DuOSblockCount{})
	_, err = cx.Gui.Wv.Bind("connections", &DuOSconnections{})
	_, err = cx.Gui.Wv.Bind("netlastblock", &DuOSnetLastBlock{})
	_, err = cx.Gui.Wv.Bind("status", &DuOStatus{})
	_, err = cx.Gui.Wv.Bind("txs", &DuOStransactionsExcerpts{})

	defer cx.Gui.Wv.Exit()
	cx.Gui.Wv.Dispatch(func() {
		// Load CSS files
		injectCss(rcv)
		// Load JavaScript Files
		err = evalJs(rcv)
	})
	cx.Gui.Wv.Run()

	return
}
