package ini

import (
	"fmt"
	"github.com/p9c/pod/pkg/bundler"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/core"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/p9c/pod/pkg/log"
	"github.com/robfig/cron"
	qtc "github.com/therecipe/qt/core"
	"github.com/therecipe/qt/webkit"
	"github.com/therecipe/qt/widgets"
	"net"
	"net/http"
	"net/url"
	"os"
)

const (
	winWidth  = 800
	winHeight = 550

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

const (
	appName           = "pod"
	confExt           = ".json"
	podConfigFilename = appName + confExt
	// ctlAppName           = "ctl"
	// ctlConfigFilename    = ctlAppName + confExt
	// nodeAppName          = "node"
	// nodeConfigFilename   = nodeAppName + confExt
	// walletAppName        = "wallet"
	// walletConfigFilename = walletAppName + confExt
)

func JsHandler(w http.ResponseWriter, r *http.Request) {
	// Getting the headers so we can set the correct mime type
	headers := w.Header()
	headers["Content-Type"] = []string{"application/javascript"}
	fmt.Fprint(w, `alert('Hello World!');`)
}

func InitDuOS() core.DuOS {
	log.L.SetLevel("trace", false)
	d := core.DuOS{
		CtX: conte.GetNewContext(appName, "main"),
		CrN: cron.New(),
		GuI: &core.GuI{},
		//Data:   DuOSdata{},
		//Alert:  alert.DuOSalert{},
		BnD: bnd.DuOSsveBundler(),
	}
	d.CrN.Start()
	// Needs separation and org
	//cx.App = getApp(cx)
	log.DEBUG("running App")

	d.SrV.Data = &srv.DuOSdata{
		//Status: d.Services.Status.GetDuOSstatus(),
		//TransactionsExcerpts: d.GetTransactionsExcertps(),
		//Addressbook:          d.GetAddressBook(),
	}

	//d.Components = comp.Components(d.db)
	d.DbS.DuOSdbInit(d.CtX.DataDir)
	//d.GuI.Load("data:text/html," + url.PathEscape(bnd.DecompressHexString(d.BnD["index.html"].Data)))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.ERROR(err)
	}
	defer ln.Close()

	//d.Config = d.Config.GetCoreCofig(d.Cx)
	d.GuI.QT = widgets.NewQApplication(len(os.Args), os.Args)
	var window = new(widgets.QWidget)
	window = widgets.NewQWidget(nil, 0)
	window.SetWindowTitle("ParallelCoin")
	window.Resize2(800, 600)
	d.GuI.WebView = webkit.NewQWebView(window)
	url := "data:text/html," + url.PathEscape(bnd.DecompressHexString(d.BnD["index.html"].Data))
	qurl := qtc.NewQUrl3(url, qtc.QUrl__TolerantMode)
	d.GuI.WebView.Load(qurl)

	window.Show()

	d.GuI.QT.Exec()

	//d.GuI = initGUI()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	//d.GuI.Bind("start", func() {
	//	log.TRACE("UI is ready")
	//})

	//

	// Create and bind Go object to the UI

	//duos := &core.DuOS{}

	//d.GuI.Bind("duOS", duos.GetDuOS())

	// Create and bind Go object to the UI

	// Call JS that calls Go and so on and so on...
	//m := d.GuI.Eval(fmt.Sprint(duos.GetDuOS()))
	//fmt.Println(m)

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	headers := w.Header()
	//	headers["Content-Type"] = []string{"application/javascript"}
	//	fmt.Fprint(w, "jhjhjhjhjhjhjhjhjh")
	//})

	//http.HandleFunc("/file.js", JsHandler)

	//d.GuI.Load(fmt.Sprintf("http://%s", ln.Addr()))
	//fmt.Println("asdasasas", bnd.DuOSsveBundler())
	//http.Handle("/", http.FileServer())
	//log.ERROR(http.ListenAndServe(":0", nil))

	//http.ListenAndServe(":8080", nil)
	//
	//for _, a := range d.BnD {
	//	http.HandleFunc(bnd.Path(a), func(w http.ResponseWriter, r *http.Request) {
	//		headers := w.Header()
	//		headers["Content-Type"] = []string{a.ContentType}
	//	})
	//}
	//ContentType
	// Load HTML after Go functions are bound to JS

	// d.BnD.DuOSassetsHandler()
	//http.Handle("/echo", websocket.Handler(EchoServer))
	//err := http.ListenAndServe(":12345", nil)
	//if err != nil {
	//	panic("ListenAndServe: " + err.Error())
	//}
	//d.GuI.Eval(bnd.DecompressHexString(d.BnD["svelte.js"].Data))

	return d
}
