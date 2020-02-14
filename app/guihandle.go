// +build !headless

package app

import (
	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/wallet"
)

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		//sys := view.DuOSboot()
		rc := rcd.RcInit(cx)
		// sys.Components = mvc.LoadComponents(duo,rc)
		
		// var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		duo, err := duoui.DuOuI()

		// sys.Components["logger"].Controller()
		rc.DuoUIloggerController()
		
		interrupt.AddHandler(func() {
			close(duo.Quit)
		})
		
		log.INFO("IsFirstRun? ", rc.Boot.IsFirstRun)

		// signal the GUI that the back end is ready
		log.DEBUG("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls
		// immediately the GUI starts
		go func() {
			nodeChan := make(chan *rpc.Server)
			// Start Node
			err = gui.DuoUInode(rc.Cx, nodeChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for nodeChan")
			rc.Cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			rc.Cx.Node.Store(true)
			
			walletChan := make(chan *wallet.Wallet)
			// Start wallet
			err = gui.Services(rc.Cx, walletChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for walletChan")
			rc.Cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			rc.Cx.Wallet.Store(true)
			rc.Boot.IsBoot = false
			duo.Ready <- struct{}{}
		}()
		
		// Start up GUI
		log.DEBUG("starting up GUI")
		// go func() {
		err = gui.WalletGUI(duo,rc)
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
		if !rc.Cx.Node.Load().(bool) {
			close(rc.Cx.WalletKill)
		}
		log.DEBUG("shutting down wallet")
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return
	}
}
