package app

import (
	"fmt"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/logi/serve"
	"os"
	"sync"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func WalletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		serve.Log(cx.KillAll, save.Filters(*cx.Config.DataDir))
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
		wg.Wait()
		return
	}
}
