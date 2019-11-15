package app

import (
	"fmt"
	"os"
	"sync"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
)

func walletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		Configure(cx, c)
		dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
			Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(dbFilename) {
			log.L.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				log.ERROR("failed to create wallet", err)
				return err
			}
			log.INFO("quitting, restart to initialize")
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
			//log.L.SetLevel(*cx.Config.LogLevel, true)
		}
		walletChan := make(chan *wallet.Wallet)
		cx.WalletKill = make(chan struct{})
		go func() {
			err = walletmain.Main(cx.Config, cx.StateCfg,
				cx.ActiveNet, walletChan, cx.WalletKill, &wg)
			if err != nil {
				log.ERROR("failed to start up wallet", err)
			}
		}()
		cx.WalletServer = <-walletChan
		wg.Wait()
		return
	}
}
