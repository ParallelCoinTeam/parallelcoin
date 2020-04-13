package app

import (
	"fmt"
	"github.com/p9c/pod/app/config"
	"github.com/urfave/cli"
	"os"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func WalletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name)
		dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
			Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(dbFilename) {
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				Error("failed to create wallet", err)
				return err
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
		}
		walletChan := make(chan *wallet.Wallet)
		cx.WalletKill = make(chan struct{})
		go func() {
			err = walletmain.Main(cx)
			if err != nil {
				Error("failed to start up wallet", err)
			}
		}()
		cx.WalletServer = <-walletChan
		cx.WaitGroup.Wait()
		return
	}
}
