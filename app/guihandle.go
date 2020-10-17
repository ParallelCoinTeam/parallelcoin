// +build !headless

package app

import (
	"gioui.org/app"
	"github.com/urfave/cli"

	"github.com/p9c/pod/pkg/util/interrupt"

	"github.com/p9c/pod/app/config"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"

	_ "gioui.org/app/permission/storage"

	gwallet "github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/gui/wallet/dap"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name, true)
		Warn("starting GUI")
		d := dap.NewDap(cx, "Duo App Plan9")
		Debug("wallet file", *cx.Config.WalletFile)
		if !apputil.FileExists(*cx.Config.WalletFile) {
			d.BOOT().Rc.Boot.IsFirstRun = true
		}
		// duo, err := duoui.DuOuI(rc)
		// rc.DuoUIloggerController()
		interrupt.AddHandler(func() {
			Debug("guiHandle interrupt")
			close(d.BOOT().Rc.Quit)
		})
		// Debug("IsFirstRun? ", rc.Boot.IsFirstRun)
		// signal the GUI that the back end is ready
		Debug("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls immediately the GUI starts
		if !d.BOOT().Rc.Boot.IsFirstRun {
			// go d.BOOT().Rc.StartServices()
		}
		// Start up GUI
		Debug("starting up GUI")
		cx.WaitGroup.Add(1)
		// err = gui.WalletGUI(duo, rc)
		// if err != nil {
		//	Error(err)
		// }
		d.NewSap(gwallet.NewGioWallet(d.BOOT()))
		go d.DAP()
		app.Main()

		cx.WaitGroup.Done()
		Debug("wallet GUI finished")
		// wait for stop signal
		<-d.BOOT().Rc.Quit
		cx.WaitGroup.Wait()
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
