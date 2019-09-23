package app

import (
	"sync"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/util"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/wallet"
)

func walletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		Configure(cx)
		dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
			Params.Name + slash + wallet.WalletDbName
		if !util.FileExists(dbFilename) {
			log.L.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				log.ERROR("failed to create wallet", err)
				return err
			}
			log.L.SetLevel(*cx.Config.LogLevel, true)
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
