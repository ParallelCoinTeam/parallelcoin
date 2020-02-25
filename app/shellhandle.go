package app

import (
	"os"
	"sync"
	
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
		var wg sync.WaitGroup
		Configure(cx, c)
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
		if !*cx.Config.NodeOff {
			go func() {
				Configure(cx, c)
				err = node.Main(cx, shutdownChan)
				if err != nil {
					log.ERROR("error starting node ", err)
				}
			}()
			cx.RPCServer = <-cx.NodeChan
		}
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
		// interrupt.AddHandler(func() {
		// 	log.WARN("interrupt received, " +
		// 		"shutting down shell modules")
		// 	close(cx.WalletKill)
		// 	close(cx.NodeKill)
		// })
		wg.Wait()
		return nil
	}
}
