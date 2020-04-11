// +build !headless

package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/logi/serve"
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		serve.Log(cx.KillAll, save.Filters(*cx.Config.DataDir))
		config.Configure(cx, c.Command.Name)
		Warn("starting GUI")
		rc := rcd.RcInit(cx)
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		duo, err := duoui.DuOuI(rc)
		rc.DuoUIloggerController()
		interrupt.AddHandler(func() {
			Debug("guiHandle interrupt")
			close(rc.Quit)
		})
		Info("IsFirstRun? ", rc.Boot.IsFirstRun)
		// signal the GUI that the back end is ready
		Debug("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls immediately the GUI starts
		if !rc.Boot.IsFirstRun {
			go rc.StartServices()
		}
		// Start up GUI
		Debug("starting up GUI")
		err = gui.WalletGUI(duo, rc)
		if err != nil {
			Error(err)
		}
		Debug("wallet GUI finished")
		// wait for stop signal
		<-rc.Quit
		Debug("shutting down node")
		if !cx.Node.Load() {
			close(cx.WalletKill)
		}
		Debug("shutting down wallet")
		if !cx.Wallet.Load() {
			close(cx.NodeKill)
		}
		return
	}
}
