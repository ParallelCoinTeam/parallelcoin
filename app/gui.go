package app

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/urfave/cli"
	"os"
	"sync"
	"sync/atomic"

	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
)

//func
//guiHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
//	return func(c *cli.Context) (err error) {
//		var wg sync.WaitGroup
//		nodeChan := make(chan *rpc.Server)
//		walletChan := make(chan *wallet.Wallet)
//		kill := make(chan struct{})
//		Configure(cx)
//		if *cx.Config.TLS || *cx.Config.ServerTLS {
//			// generate the tls certificate if configured
//			_, _ = walletmain.GenerateRPCKeyPair(cx.Config, true)
//		}
//		shutdownChan := make(chan struct{})
//		dbFilename :=
//			*cx.Config.DataDir + slash +
//				cx.ActiveNet.Params.Name + slash +
//				wallet.WalletDbName
//		if !apputil.FileExists(dbFilename) {
//			log.L.SetLevel("off", false)
//			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
//				log.ERROR("failed to create wallet", err)
//			}
//			fmt.Println("restart to complete initial setup")
//			os.Exit(1)
//			log.L.SetLevel(*cx.Config.LogLevel, true)
//		}
//		if !*cx.Config.WalletOff {
//			go func() {
//				err = walletmain.Main(cx.Config, cx.StateCfg,
//					cx.ActiveNet, walletChan, kill, &wg)
//				if err != nil {
//					log.ERROR(err)
//					fmt.Println("error running wallet:", err)
//				}
//			}()
//			save.Pod(cx.Config)
//		}
//		if !*cx.Config.NodeOff {
//			go func() {
//				Configure(cx)
//				err = node.Main(cx, shutdownChan, kill, nodeChan, &wg)
//				if err != nil {
//					log.ERROR("error starting node ", err)
//				}
//			}()
//			cx.RPCServer = <-nodeChan
//		}
//		cx.WalletServer = <-walletChan
//		gui.GUI(cx)
//		wg.Wait()
//		return nil
//	}
//}

var guiHandle = func(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		log.WARN("starting GUI")
		Configure(cx)
		shutdownChan := make(chan struct{})
		walletChan := make(chan *wallet.Wallet)
		nodeChan := make(chan *rpc.Server)
		cx.WalletKill = make(chan struct{})
		cx.NodeKill = make(chan struct{})
		cx.Wallet = &atomic.Value{}
		cx.Wallet.Store(false)
		cx.Node = &atomic.Value{}
		cx.Node.Store(false)
		var err error
		var wg sync.WaitGroup
		if !*cx.Config.NodeOff {
			go func() {
				log.INFO(cx.Language.RenderText("goApp_STARTINGNODE"))
				err = node.Main(cx, shutdownChan, cx.NodeKill, nodeChan, &wg)
				if err != nil {
					fmt.Println("error running node:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for nodeChan")
			cx.RPCServer = <-nodeChan
			log.DEBUG("nodeChan sent")
			cx.Node.Store(true)
		}
		if !*cx.Config.WalletOff {
			go func() {
				log.INFO("starting wallet")
				err = walletmain.Main(cx.Config, cx.StateCfg,
					cx.ActiveNet, walletChan, cx.WalletKill, &wg)
				if err != nil {
					fmt.Println("error running wallet:", err)
					os.Exit(1)
				}
			}()
			log.DEBUG("waiting for walletChan")
			cx.WalletServer = <-walletChan
			log.DEBUG("walletChan sent")
			cx.Wallet.Store(true)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(cx.WalletKill)
			close(cx.NodeKill)
		})
		gui.GUI(cx)
		if !cx.Node.Load().(bool) {
			close(cx.WalletKill)
		}
		if !cx.Wallet.Load().(bool) {
			close(cx.NodeKill)
		}
		return err
	}
}
