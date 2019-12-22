// +build !headless

package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/urfave/cli"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		duo := duoui.DuOuI(cx)

		var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			firstRun = true
		}

		log.INFO("ima", firstRun)

		//loader.DuoUIloader(duo, cx, firstRun)

		Configure(cx, c)
		// Start Node
		err = gui.DuOSnode(cx)
		if err != nil {
			log.ERROR(err)
		}


		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}

		// We open up wallet creation

		gui.WalletGUI(duo)

		//b.IsBootLogo = false
		//b.IsBoot = false

		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return
	}
}
