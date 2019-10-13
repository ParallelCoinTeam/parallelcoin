package ini

import (
	"github.com/p9c/lorca"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/core"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/svelte/__OLDvue/lib/html"
	"github.com/robfig/cron"
	"net/url"
	"runtime"
)

const (
	winWidth  = 800
	winHeight = 550

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	runtime.LockOSThread()
}

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

	d.SrV.Data = &srv.DuOSdata{
		//Status: d.Services.Status.GetDuOSstatus(),
		//TransactionsExcerpts: d.GetTransactionsExcertps(),
		//Addressbook:          d.GetAddressBook(),
	}

	//d.Components = comp.Components(d.db)
	d.DbS.DuOSdbInit(d.CtX.DataDir)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	//ln, err := net.Listen("tcp", "127.0.0.1:0")
	//if err != nil {
	//	log.ERROR(err)
	//}
	//defer ln.Close()
	//go http.Serve(ln, http.FileServer(FS))
	//d.GuI.Load(fmt.Sprintf("http://%s", ln.Addr()))
	d.GuI.Load("http://127.0.0.1:5000")

	return d
}

func initGUI() lorca.UI {
	ui, err := lorca.New("data:text/html,"+url.PathEscape(html.HTML), "", 800, 600)
	if err != nil {
		log.ERROR("running App", err)
	}
	return ui
}
