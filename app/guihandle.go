// +build !headless

package app

import (
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/gcx"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/urfave/cli"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		//		//var firstRun bool
		//		//if !apputil.FileExists(*cx.Config.WalletFile) {
		//		//	firstRun = true
		//		//}
		//		//utils.GetBiosMessage(view, "starting GUI")

		//err := gui.Services(cx)
		Configure(cx, c)
		// Start Node
		err = gui.DuOSnode(cx)
		if err != nil {
			log.ERROR(err)
		}
		cx.Gui = &gcx.GUI{
			Cf: &gcx.Configuration{
				Assets:"./pkg/gui/assets/filesystem",
			},
		}
		err = gui.DuOSfileSystem(cx)


		//gui.DuOSloader(cx, firstRun)

		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}

		//gui.DuOSgatherer(cx)
		// We open up wallet creation
		gui.WalletGUI(cx)

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
