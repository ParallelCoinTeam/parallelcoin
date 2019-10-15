package ini

import (
	"fmt"
	"github.com/p9c/gui"
	"github.com/p9c/pod/pkg/bundler"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/core"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/p9c/pod/pkg/log"
	"github.com/robfig/cron"
	"net"
	"net/http"
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

func InitDuOS() core.DuOS {
	log.L.SetLevel("trace", false)
	d := core.DuOS{
		CtX: conte.GetNewContext(appName, "main"),
		CrN: cron.New(),
		//GuI: initGUI(),
		//Data:   DuOSdata{},
		//Alert:  alert.DuOSalert{},
	}
	d.CrN.Start()
	// Needs separation and org
	//cx.App = getApp(cx)
	log.DEBUG("running App")

	//d.Config = d.Config.GetCoreCofig(d.Cx)

	d.GuI = initGUI()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	d.GuI.Bind("start", func() {
		log.TRACE("UI is ready")
	})

	// Create and bind Go object to the UI

	duos := &core.DuOS{}

	d.GuI.Bind("duOS", duos.GetDuOS())

	// Create and bind Go object to the UI



	// Call JS that calls Go and so on and so on...
	m := d.GuI.Eval(fmt.Sprint(duos.GetDuOS()))
	fmt.Println(m)

	d.SrV.Data = &srv.DuOSdata{
		//Status: d.Services.Status.GetDuOSstatus(),
		//TransactionsExcerpts: d.GetTransactionsExcertps(),
		//Addressbook:          d.GetAddressBook(),
	}

	//d.Components = comp.Components(d.db)
	d.DbS.DuOSdbInit(d.CtX.DataDir)
	bnd.DuOSsveBundler()
	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.ERROR(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	http.HandleFunc("pipe.js", bnd.PipeJsHandler)
	http.HandleFunc("svelte.js", bnd.BndJsHandler)
	http.HandleFunc("svelte.css", bnd.BndCssHandler)


	//fmt.Println("asdasasas", bnd.DuOSsveBundler())
	//http.Handle("/", http.FileServer(bnd.DuOSsveBundler()))
	//log.ERROR(http.ListenAndServe(":0", nil))


	//http.ListenAndServe(":8080", nil)

	d.GuI.Load(fmt.Sprintf("http://%s", ln.Addr()))
	//d.GuI.Load("http://127.0.0.1:0")

	//http.Handle("/echo", websocket.Handler(EchoServer))
	//err := http.ListenAndServe(":12345", nil)
	//if err != nil {
	//	panic("ListenAndServe: " + err.Error())
	//}

	//d.GuI.Eval(string(bnd.DuOSsveBundler()["svelte.js"]))
	return d
}

func initGUI() gui.UI {
	ui, err := gui.New("", "", 800, 600)
	if err != nil {
		log.ERROR("running App", err)
	}
	return ui
}
