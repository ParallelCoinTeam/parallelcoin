// +build !headless

package app

import (
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	log "github.com/p9c/logi"
	"github.com/p9c/util/interrupt"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		log.L.Warn("starting GUI")
		rc := rcd.RcInit(cx)
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		duo, err := duoui.DuOuI(rc)
		rc.DuoUIloggerController()
		interrupt.AddHandler(func() {
			log.L.Debug("guiHandle interrupt")
			close(rc.Quit)
		})
		log.L.Info("IsFirstRun? ", rc.Boot.IsFirstRun)
		// signal the GUI that the back end is ready
		log.L.Debug("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls immediately the GUI starts
		go rc.StartServices()
		// Start up GUI
		log.L.Debug("starting up GUI")
		err = gui.WalletGUI(duo, rc)
		if err != nil {
			log.L.Error(err)
		}
		log.L.Debug("wallet GUI finished")
		// wait for stop signal
		<-rc.Quit
		log.L.Debug("shutting down node")
		if !cx.Node.Load() {
			close(cx.WalletKill)
		}
		log.L.Debug("shutting down wallet")
		if !cx.Wallet.Load() {
			close(cx.NodeKill)
		}
		return
	}
}
