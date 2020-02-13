// +build !headless

package app

import (
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/mvc/view"
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
		sys := view.DuOSboot()
		sys.Rc = rcd.RcInit(cx)
		// sys.Components = mvc.LoadComponents(duo,rc)
		
		// var firstRun bool
		if !apputil.FileExists(*cx.Config.WalletFile) {
			sys.Rc.Boot.IsFirstRun = true
		}
		duo, err := duoui.DuOuI(sys.Rc)
		sys.Duo = duo
		// sys.Components["logger"].Controller()
		
		interrupt.AddHandler(func() {
			close(sys.Duo.Quit)
		})
		
		log.INFO("IsFirstRun? ", sys.Rc.Boot.IsFirstRun)
		
		// (rc, cx, firstRun)
		// rcd.ListenInit(*rc)
		go func() {
			for {
				select {
				case <-sys.Duo.Quit:
					break
				}
			}
		}()
		
		// signal the GUI that the back end is ready
		log.DEBUG("sending ready signal")
		// we can do this without blocking because the channel has 1 buffer this way it falls
		// immediately the GUI starts
		go func() {
			nodeChan := make(chan *rpc.Server)
			// Start Node
			err = gui.DuoUInode(sys.Rc.Cx, nodeChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for nodeChan")
			sys.Rc.Cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			sys.Rc.Cx.Node.Store(true)
			
			walletChan := make(chan *wallet.Wallet)
			// Start wallet
			err = gui.Services(sys.Rc.Cx, walletChan)
			if err != nil {
				log.ERROR(err)
			}
			log.DEBUG("waiting for walletChan")
			sys.Rc.Cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			sys.Rc.Cx.Wallet.Store(true)
			sys.Rc.Boot.IsBoot = false
			sys.Duo.Ready <- struct{}{}
		}()
		
		// Start up GUI
		log.DEBUG("starting up GUI")
		// go func() {
		err = gui.WalletGUI(sys)
		if err != nil {
			log.ERROR(err)
		}
		
		log.DEBUG("wallet GUI finished")
		// }()
		// wait for stop signal
		<-sys.Duo.Quit
		// b.IsBootLogo = false
		// b.IsBoot = false
		log.DEBUG("shutting down node")
		if !sys.Rc.Cx.Node.Load().(bool) {
			close(sys.Rc.Cx.WalletKill)
		}
		log.DEBUG("shutting down wallet")
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return
	}
}
