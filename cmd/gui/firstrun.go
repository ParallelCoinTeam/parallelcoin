package gui

import (
	"encoding/hex"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/webview"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
	"github.com/shurcooL/vfsgen"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CreateForm struct {
	cx *conte.Xt
}

func FirstRun(cx *conte.Xt) {
	log.INFO("FFFFFstarting GUI")

	var fs http.FileSystem = http.Dir("./pkg/gui/assets/firstrun")
	err := vfsgen.Generate(fs, vfsgen.Options{})
	if err != nil {
		log.FATAL(err)
	}
	w := webview.New(webview.Settings{
		Width:                  480,
		Height:                 640,
		Debug:                  true,
		Title:                  "ParallelCoin - DUO - True Story",
		URL:                    "data:text/html," + url.PathEscape(getFile("/index.html", fs)),
		ExternalInvokeCallback: handleRPCfirstrun,
	})

	log.INFO("starting GUI")

	defer w.Exit()
	w.Dispatch(func() {

		w.Bind("cw", &CreateForm{cx: cx})

		// Load CSS
	})
	w.Run()

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
	case data == "fullscreen":
		w.SetFullscreen(true)
	case data == "unfullscreen":
		w.SetFullscreen(false)
	case data == "open":
		log.Println("open", w.Dialog(webview.DialogTypeOpen, 0, "Open file", ""))
	case strings.HasPrefix(data, "changeTitle:"):
		w.SetTitle(strings.TrimPrefix(data, "changeTitle:"))
	}
}

func (c *CreateForm) CreateWallet(p, s, b, f string) {
	var err error
	var seed []byte
	if f == "" {
		f = *c.cx.Config.WalletFile
	}
	l := wallet.NewLoader(c.cx.ActiveNet, *c.cx.Config.WalletFile, 250)

	if s == "" {
		seed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			log.ERROR(err)
			panic(err)
		}
	} else {
		seed, err = hex.DecodeString(s)
		if err != nil {
			// Need to make JS invocation to embed
			log.ERROR(err)
		}
	}

	_, err = l.CreateNewWallet([]byte(b), []byte(p), seed, time.Now(), true)
	if err != nil {
		log.ERROR(err)
		panic(err)
	}

	log.INFO("ratattaa", b, p, seed)
	*c.cx.Config.WalletPass = b
	*c.cx.Config.WalletFile = f

	save.Pod(c.cx.Config)

	log.INFO(c)
}
