// +build !headless

package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/pkg/logi/serve"
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		serve.Log(cx.KillAll)
		config.Configure(cx, c.Command.Name)
		L.Warn("starting GUI")
		rc := rcd.RcInit(cx)
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		duo, err := duoui.DuOuI(rc)
		rc.DuoUIloggerController()
		interrupt.AddHandler(func() {
			L.Debug("guiHandle interrupt")
			close(rc.Quit)
		})
		L.Info("IsFirstRun? ", rc.Boot.IsFirstRun)
		// signal the GUI that the back end is ready
		L.Debug("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls immediately the GUI starts
		if !rc.Boot.IsFirstRun {
			go rc.StartServices()
		}
		// Start up GUI
		L.Debug("starting up GUI")
		err = gui.WalletGUI(duo, rc)
		if err != nil {
			L.Error(err)
		}
		L.Debug("wallet GUI finished")
		// wait for stop signal
		<-rc.Quit
		L.Debug("shutting down node")
		if !cx.Node.Load() {
			close(cx.WalletKill)
		}
		L.Debug("shutting down wallet")
		if !cx.Wallet.Load() {
			close(cx.NodeKill)
		}
		return
	}
}
