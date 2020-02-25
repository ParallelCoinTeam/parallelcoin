package app

import (
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		log.DEBUG("starting shell")
		if *cx.Config.TLS || *cx.Config.ServerTLS {
			// generate the tls certificate if configured
			_, _ = walletmain.GenerateRPCKeyPair(cx.Config, true)
		}
		shutdownChan := make(chan struct{})
		dbFilename :=
			*cx.Config.DataDir + slash +
				cx.ActiveNet.Params.Name + slash +
				wallet.WalletDbName
		if !apputil.FileExists(dbFilename) {
			// log.L.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				log.ERROR("failed to create wallet", err)
			}
			log.Println("restart to complete initial setup")
			os.Exit(1)
			// log.L.SetLevel(*cx.Config.LogLevel, true)
		}
		log.WARN("starting node")
		if !*cx.Config.NodeOff {
			go func() {
				err = node.Main(cx, shutdownChan)
				if err != nil {
					log.ERROR("error starting node ", err)
				}
			}()
			cx.RPCServer = <-cx.NodeChan
		}
		log.WARN("starting wallet")
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx)
				if err != nil {
					log.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-cx.WalletChan
			// save.Pod(cx.Config)
		}
		log.DEBUG("shell started")
		// interrupt.AddHandler(func() {
		// 	log.WARN("interrupt received, " +
		// 		"shutting down shell modules")
		// 	close(cx.WalletKill)
		// 	close(cx.NodeKill)
		// })
		cx.WaitGroup.Wait()
		return nil
	}
}
