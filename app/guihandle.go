// +build !headless

package app

import (
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	
	"github.com/urfave/cli"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		duo := duoui.DuOuI(cx)
		interrupt.AddHandler(func(){
			close(duo.Quit)
		})
		var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			firstRun = true
		}

		log.INFO("ima", firstRun)

		//loader.DuoUIloader(duo, cx, firstRun)

		Configure(cx, c)
		
		
		// Start node
		err = gui.DuOSnode(cx)
		if err != nil {
			log.ERROR(err)
		}

		// Start wallet
		err = gui.Services(cx)
		if err != nil {
			log.ERROR(err)
		}
		
		// signal the GUI that the back end is ready
		log.DEBUG("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls immediately the GUI starts
		duo.Ready <-struct{}{}
		// Start up GUI
		go func() {
			gui.WalletGUI(duo)
			log.DEBUG("wallet GUI finished")
		}()
		// wait for stop signal
		<-duo.Quit
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
