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

func Loader(cx *conte.Xt) {
	log.INFO("FFFFFstarting GUI")
	var fs http.FileSystem = http.Dir("./pkg/gui/assets/f/assets")
	err := vfsgen.Generate(fs, vfsgen.Options{
		PackageName:  "guiLibsFirstRun",
		BuildTags:    "!dev",
		VariableName: "FirstRun",
	})
	if err != nil {
		log.FATAL(err)
	}
	cx.FileSystem = &fs

	cx.WebView = webview.New(webview.Settings{
		Width:                  600,
		Height:                 800,
		Debug:                  true,
		Resizable:              false,
		Title:                  "ParallelCoin - DUO - True Story",
		URL:                    "data:text/html," + url.PathEscape(getFile("index.html", *cx.FileSystem)),
		ExternalInvokeCallback: handleRPCfirstrun,
	})

	log.INFO("starting GUI")

	defer cx.WebView.Exit()
	cx.WebView.Dispatch(func() {

		err = cx.WebView.Eval(getFile("dui.js", fs))
		if err != nil {
			log.DEBUG("error binding to webview:", err)
		}

		cx.WebView.InjectCSS(getFile("css/theme/root.css", fs))
		cx.WebView.InjectCSS(getFile("css/theme/colors.css", fs))
		cx.WebView.InjectCSS(getFile("css/theme/helpers.css", fs))
		cx.WebView.InjectCSS(getFile("css/style.css", fs))
		cx.WebView.InjectCSS(getFile("dui.css", fs))

		cx.WebView.Bind("cw", &CreateForm{cx: cx})

		// Load CSS
	})
	cx.WebView.Run()

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

func (c *CreateForm) CloseFirstRun() {
	c.cx.WebView.Terminate()
}
