package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"github.com/p9c/pod/pkg/log"
	"net/url"
	"os"
)

const slash = string(os.PathSeparator)


func GUI(cx *conte.Xt) {
	rc := rcvar{
		Xt:     cx,
		alert:  DuOSalert{},
		status: DuOStatus{},
		txs:    DuOStransactionsExcerpts{},
		lastxs: DuOStransactions{},
	}


	rc.WebView = webview.New(webview.Settings{
		Width:  1024,
		Height: 760,
		Debug:  true,
		Title:  "ParallelCoin - DUO - True Story",
		URL:    "data:text/html," + url.PathEscape(getFile("/w/index.html", *cx.FileSystem)),
	})


	//b := Bios{
	//	Theme:      false,
	//	IsBoot:     true,
	//	IsBootMenu: true,
	//	IsBootLogo: true,
	//	IsLoading:  false,
	//	IsDev:      true,
	//	IsScreen:   "overview",
	//}

	log.INFO("starting GUI")

	defer rc.WebView.Exit()
	rc.WebView.Dispatch(func() {

		// Load CSS files
		injectCss(&rc)
		// Load JavaScript Files
		evalJs(&rc)
	})
	rc.WebView.Run()

	//
	//go func() {
	//	for _ = range time.NewTicker(time.Second * 1).C {
	//
	//
	//		//status, err := json.Marshal(rc.GetDuOStatus())
	//		//if err != nil {
	//		//}
	//		//transactions, err := json.Marshal(rc.GetTransactions(0, 555, ""))
	//		//if err != nil {
	//		//}
	//}
	//}()

}
