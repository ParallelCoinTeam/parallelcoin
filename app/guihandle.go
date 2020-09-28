// +build !headless

package app

import (
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pkg/app/slog"
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
		config.Configure(cx, c.Command.Name, true)
		slog.Warn("starting GUI")
		rc := rcd.RcInit(cx)
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		var duo *model.DuoUI
		duo, err = duoui.DuOuI(rc)
		rc.DuoUIloggerController()
		interrupt.AddHandler(func() {
			slog.Debug("guiHandle interrupt")
			close(rc.Quit)
		})
		slog.Debug("IsFirstRun? ", rc.Boot.IsFirstRun)
		// signal the GUI that the back end is ready
		slog.Debug("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this
		// way it falls immediately the GUI starts
		if !rc.Boot.IsFirstRun {
			go rc.StartServices()
		}
		// Start up GUI
		slog.Debug("starting up GUI")
		cx.WaitGroup.Add(1)
		if err = gui.WalletGUI(duo, rc); slog.Check(err) {}
		cx.WaitGroup.Done()
		slog.Debug("wallet GUI finished")
		// wait for stop signal
		<-rc.Quit
		cx.WaitGroup.Wait()
		slog.Debug("shutting down node")
		if !cx.Node.Load() {
			close(cx.WalletKill)
		}
		slog.Debug("shutting down wallet")
		if !cx.Wallet.Load() {
			close(cx.NodeKill)
		}
		return
	}
}
