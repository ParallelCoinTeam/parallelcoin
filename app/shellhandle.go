package app

import (
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/cli"
	"os"
	"sync"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		nodeChan := make(chan *rpc.Server)
		walletChan := make(chan *wallet.Wallet)
		kill := make(chan struct{})
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
			//log.L.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				log.ERROR("failed to create wallet", err)
			}
			log.Println("restart to complete initial setup")
			os.Exit(1)
			//log.L.SetLevel(*cx.Config.LogLevel, true)
		}
		if !*cx.Config.NodeOff {
			go func() {
				Configure(cx, c)
				err = node.Main(cx, shutdownChan, kill, nodeChan, &wg)
				if err != nil {
					log.ERROR("error starting node ", err)
				}
			}()
			cx.RPCServer = <-nodeChan
		}
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx.Config, cx.StateCfg,
					cx.ActiveNet, walletChan, kill, &wg)
				if err != nil {
					log.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-walletChan
			//save.Pod(cx.Config)
		}
		interrupt.AddHandler(func() {
			log.WARN("interrupt received, " +
				"shutting down shell modules")
			close(cx.WalletKill)
			close(cx.NodeKill)
		})
		wg.Wait()
		return nil
	}
}
