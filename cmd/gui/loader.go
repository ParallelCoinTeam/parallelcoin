package gui

import (
	"encoding/hex"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
	"net/url"
	"time"
)

func Loader(b *Bios, cx *conte.Xt) {
	var err error
	b.cx = cx
	log.INFO("FFFFFstarting GUI")
	//var fs http.FileSystem = http.Dir("./pkg/gui/assets/f/assets")
	//err := vfsgen.Generate(fs, vfsgen.Options{
	//	PackageName:  "guiLibsLoader",
	//	BuildTags:    "dev",
	//	VariableName: "Loader",
	//})
	//if err != nil {
	//	log.FATAL(err)
	//}
	//b.Fs = &fs

	b.Wv = webview.New(webview.Settings{
		Width:     600,
		Height:    800,
		Debug:     true,
		Resizable: false,
		Title:     "ParallelCoin - DUO - True Story",
		URL:       "data:text/html," + url.PathEscape(getFile("loader.html", *b.Fs)),
		//ExternalInvokeCallback: handleRPCfirstrun,
	})

	log.INFO("starting GUI")

	//defer b.Wv.Exit()
	b.Wv.Dispatch(func() {

		_, err = b.Wv.Bind("duos", &Bios{
			cx:         b.cx,
			IsFirstRun: b.IsFirstRun,
		})

		err = b.Wv.Eval(getFile("js/svelte.js", *b.Fs))
		if err != nil {
			log.DEBUG("error binding to webview:", err)
		}

		b.Wv.InjectCSS(getFile("css/theme/root.css", *b.Fs))
		b.Wv.InjectCSS(getFile("css/theme/colors.css", *b.Fs))
		b.Wv.InjectCSS(getFile("css/theme/helpers.css", *b.Fs))
		b.Wv.InjectCSS(getFile("css/loader.css", *b.Fs))
		b.Wv.InjectCSS(getFile("css/svelte.css", *b.Fs))

		// Load CSS
	})
	b.Wv.Run()

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

func handleRPCfirstrun(w webview.WebView, data string) {
	switch {
	case data == "close":
		w.Terminate()
	case data == "open":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
	}
}

func (b *Bios) CreateWallet(pr, sd, pb, fl string) {
	var err error
	var seed []byte
	if fl == "" {
		fl = *b.cx.Config.WalletFile
	}
	l := wallet.NewLoader(b.cx.ActiveNet, *b.cx.Config.WalletFile, 250)

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

	b.IsFirstRun = false
	*b.cx.Config.WalletPass = pb
	*b.cx.Config.WalletFile = fl

	save.Pod(b.cx.Config)
	log.INFO(b)
}

func (c *Bios) CloseLoader() {
	c.Wv.Exit()
}
