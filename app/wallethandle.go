package app

import (
	"sync"

	"github.com/urfave/cli"

	"git.parallelcoin.io/dev/pod/app/apputil"
	"git.parallelcoin.io/dev/pod/cmd/walletmain"
	"git.parallelcoin.io/dev/pod/pkg/conte"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
	"git.parallelcoin.io/dev/pod/pkg/wallet"
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
				cx.Log <- cl.Error{"failed to create wallet",
					err, cl.Ine()}
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
				cx.Log <- cl.Error{"failed to start up wallet",
					 err, cl.Ine()}
			}
		}()
		cx.WalletServer = <-walletChan
		wg.Wait()
		return
	}
}
