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
		rc := rcd.RcInit()
		// var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			rc.Boot.IsFirstRun = true
		}
		
		duo, err := duoui.DuOuI(rc, cx)
		interrupt.AddHandler(func() {
			close(duo.Quit)
		})
		
		log.INFO("IsFirstRun? ", rc.Boot.IsFirstRun)
		
		// loader.DuoUIloader(rc, cx, firstRun)
		// rcd.ListenInit()
		// go func() {
		//	for {
		//		select {
		//		case <-duo.Quit:
		//			break
		//			case
		//		}
		//	}
		// }()
		
		// signal the GUI that the back end is ready
		log.DEBUG("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls
		// immediately the GUI starts
		go func() {
			nodeChan := make(chan *rpc.Server)
			// Start Node
			err = gui.DuoUInode(cx, nodeChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for nodeChan")
			cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			cx.Node.Store(true)
			
			walletChan := make(chan *wallet.Wallet)
			// Start wallet
			err = gui.Services(cx, walletChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for walletChan")
			cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			cx.Wallet.Store(true)
			duo.Ready <- struct{}{}
			rc.Boot.IsBoot = false
		}()
		
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
