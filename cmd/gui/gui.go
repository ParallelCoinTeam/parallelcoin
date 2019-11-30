package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/log"
	"net/url"
)

func WalletGUI(cx *conte.Xt) (err error) {
	rc := &rcvar{}
	cx.Gui.Wv = webview.New(webview.Settings{
		Width:     1024,
		Height:    760,
		Debug:     true,
		Resizable: false,
		Title:     "ParallelCoin - DUO - True Story",
		URL:       "data:text/html," + url.PathEscape(getFile("vue.html", *cx.Gui.Fs)),
	})
	cx.Gui.Wv.SetColor(68, 68, 68, 255)

	//conf := pod.Config{}

	//log.INFO(cx.Config)
	rc.cx = cx

	//_, err = cx.Gui.Wv.Bind("alert", &DuOSalert{})
	//
	//_, err = cx.Gui.Wv.Bind("status", &DuOStatus{})
	//
	//_, err = cx.Gui.Wv.Bind("hashes", &DuOShashes{})
	//_, err = cx.Gui.Wv.Bind("nethash", &DuOSnetworkHash{})
	//_, err = cx.Gui.Wv.Bind("height", &DuOSheight{})
	//_, err = cx.Gui.Wv.Bind("bestblock", &DuOSbestBlockHash{})
	//
	//_, err = cx.Gui.Wv.Bind("blockcount", &DuOSblockCount{})
	//_, err = cx.Gui.Wv.Bind("netlastblock", &DuOSnetLastBlock{})
	//_, err = cx.Gui.Wv.Bind("connections", &DuOSconnections{})
	//
	//_, err = cx.Gui.Wv.Bind("balance", &DuOSbalance{})
	//_, err = cx.Gui.Wv.Bind("unconfirmed", &DuOSunconfirmed{})
	//_, err = cx.Gui.Wv.Bind("txsnumber", &DuOStransactionsNumber{})
	//
	//_, err = cx.Gui.Wv.Bind("transactions", &DuOStransactions{})
	//_, err = cx.Gui.Wv.Bind("txs", &DuOStransactionsExcerpts{})
	//_, err = cx.Gui.Wv.Bind("lastxs", &DuOStransactions{})
	//
	//_, err = cx.Gui.Wv.Bind("localhost", &DuOSlocalHost{})

	defer cx.Gui.Wv.Exit()

	cx.Gui.Wv.Dispatch(func() {

		_, err = cx.Gui.Wv.Bind("rcvar", &rcvar{})

		daemon := DaemonConfig{
			Config: cx.Config,
			Schema: pod.GetConfigSchema(),
		}
		// Bind configuration
		_, err = cx.Gui.Wv.Bind("duOSsettings", &DuOSsettings{
			Daemon: daemon,
		})

		// Bind navigation
		_, err = cx.Gui.Wv.Bind("duOSnav", &DuOSnav{
			cx:     cx,
			Screen: "PageOverview",
		})
		if err != nil {
			log.ERROR("error binding to webview:", err)
		}

		// Load CSS files
		rc.injectCss()
		// Load JavaScript Files
		err = rc.evalJs()
		if err != nil {
			log.ERROR("error binding to webview:", err)
		}
	})

	rc.DuOSgatherer()
	//cx.Gui.Wv.Dispatch(func() {

	//log.INFO("ssasasass", rc)
	//})

	cx.Gui.Wv.Run()

	return
}
