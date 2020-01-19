// +build !headless

package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"

	"github.com/urfave/cli"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		rc := rcd.RcInit()
		//var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.IsFirstRun = true
		}

		duo, err := duoui.DuOuI(rc, cx)
		interrupt.AddHandler(func() {
			close(duo.Quit)
		})

		log.INFO("ima", rc.IsFirstRun)

		//loader.DuoUIloader(rc, cx, firstRun)

		// signal the GUI that the back end is ready
		log.DEBUG("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls
		// immediately the GUI starts
		duo.Ready <- struct{}{}
		// Start Node
		err = gui.DuoUInode(cx)
		if err != nil {
			log.ERROR(err)
		}
		// Start wallet
		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}

		// Start up GUI
		log.DEBUG("starting up GUI")
		// go func() {
		err = gui.WalletGUI(duo, cx, rc)
		if err != nil {
			log.ERROR(err)
		}

		log.DEBUG("wallet GUI finished")
		// }()
		// wait for stop signal
		<-duo.Quit
		// b.IsBootLogo = false
		// b.IsBoot = false
		log.DEBUG("shutting down node")
		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		log.DEBUG("shutting down wallet")
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return
	}
}
