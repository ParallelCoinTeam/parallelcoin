package ini

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/core"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/p9c/pod/pkg/log"
	"github.com/robfig/cron"
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
		Cx: conte.GetNewContext(appName, "main"),
		Cr: cron.New(),
		//Data:   DuOSdata{},
		//Alert:  alert.DuOSalert{},
	}
	d.Cr.Start()
	// Needs separation and org
	//cx.App = getApp(cx)
	log.DEBUG("running App")

	//d.Config = d.Config.GetCoreCofig(d.Cx)

	d.Services.Data = &srv.DuOSdata{
		//Status: d.Services.Status.GetDuOSstatus(),
		//TransactionsExcerpts: d.GetTransactionsExcertps(),
		//Addressbook:          d.GetAddressBook(),
	}

	//d.Components = comp.Components(d.db)
	d.DB.DuOSdbInit(d.Cx.DataDir)

	return d
}
