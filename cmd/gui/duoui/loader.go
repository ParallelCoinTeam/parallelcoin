package duoui

import (
	"encoding/hex"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
	"time"
)

func (duo *DuoUI)DuoUIloader(firstRun bool) (err error) {

	// create loader window
	//cx.Gui.Wv = webview.New(webview.Settings{
	//	Width:     600,
	//	Height:    800,
	//	Debug:     true,
	//	Resizable: false,
	//	Title:     "ParallelCoin - DUO - True Story",
	//	URL:       "data:text/html," + url.PathEscape(getFile("loader.html", *cx.Gui.Fs)),
	//	//ExternalInvokeCallback: handleRPCfirstrun,
	//})
	//cx.Gui.Wv.SetColor(68, 68, 68, 255)

	//defer cx.Gui.Wv.Exit()
	//cx.Gui.Wv.Dispatch(func() {
	//
	//	_, err = cx.Gui.Wv.Bind("duos", &rcvar{
	//		cx:         cx,
	//		IsFirstRun: firstRun,
	//	})
	//
	//	err = cx.Gui.Wv.Eval(getFile("js/svelte.js", *cx.Gui.Fs))
	//	if err != nil {
	//		log.DEBUG("error binding to webview:", err)
	//	}
	//
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/root.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/colors.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/theme/helpers.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/loader.css", *cx.Gui.Fs))
	//	cx.Gui.Wv.InjectCSS(getFile("css/svelte.css", *cx.Gui.Fs))

	// Load CSS
	//})
	//cx.Gui.Wv.Run()

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
	return
}

//func handleRPCfirstrun(w webview.WebView, data string) {
//	switch {
//	case data == "close":
//		w.Terminate()
//	case data == "open":
//		//log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
//	}
//}

func (duo *DuoUI) CreateWallet(pr, sd, pb, fl string) {
	var err error
	var seed []byte
	if fl == "" {
		fl = *duo.cx.Config.WalletFile
	}
	l := wallet.NewLoader(duo.cx.ActiveNet, *duo.cx.Config.WalletFile, 250)

	if sd == "" {
		seed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			log.ERROR(err)
			panic(err)
		}
	} else {
		seed, err = hex.DecodeString(sd)
		if err != nil {
			// Need to make JS invocation to embed
			log.ERROR(err)
		}
	}

	_, err = l.CreateNewWallet([]byte(pb), []byte(pr), seed, time.Now(), true)
	if err != nil {
		log.ERROR(err)
		panic(err)
	}

	//duo.Boot.IsFirstRun = false
	*duo.cx.Config.WalletPass = pb
	*duo.cx.Config.WalletFile = fl

	save.Pod(duo.cx.Config)
	//log.INFO(rc)
}

func (r *RcVar) CloseDuOSloader() {
	//r.cx.Gui.Wv.Exit()
}
