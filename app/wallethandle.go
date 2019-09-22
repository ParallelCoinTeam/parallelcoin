package app

import (
	"sync"

	"github.com/urfave/cli"

	"github.com/parallelcointeam/parallelcoin/app/apputil"
	"github.com/parallelcointeam/parallelcoin/cmd/walletmain"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/log"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet"
)

func walletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		Configure(cx)
		dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
			Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(dbFilename) {
			cl.Register.SetAllLevels("off")
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				log.ERROR("failed to create wallet", err)
				return err
			}
			cl.Register.SetAllLevels(*cx.Config.LogLevel)
		}
		walletChan := make(chan *wallet.Wallet)
		cx.WalletKill = make(chan struct{})
		go func() {
			err = walletmain.Main(cx.Config, cx.StateCfg,
				cx.ActiveNet, walletChan, cx.WalletKill, &wg)
			if err != nil {
				log.ERROR("failed to start up wallet",err)
			}
		}()
		cx.WalletServer = <-walletChan
		wg.Wait()
		return
	}
}
