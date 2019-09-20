package app

import (
	"fmt"
	"sync"

	"github.com/urfave/cli"

	"git.parallelcoin.io/dev/pod/app/apputil"
	"git.parallelcoin.io/dev/pod/cmd/node"
	"git.parallelcoin.io/dev/pod/cmd/node/rpc"
	"git.parallelcoin.io/dev/pod/cmd/walletmain"
	"git.parallelcoin.io/dev/pod/pkg/conte"
	"git.parallelcoin.io/dev/pod/pkg/util/cl"
	"git.parallelcoin.io/dev/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		var wg sync.WaitGroup
		Configure(cx)
		shutdownChan := make(chan struct{})
		dbFilename :=
			*cx.Config.DataDir + slash +
				cx.ActiveNet.Params.Name + slash +
				wallet.WalletDbName
		if !apputil.FileExists(dbFilename) {
			cl.Register.SetAllLevels("off")
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				cx.Log <- cl.Error{"failed to create wallet", err}
			}
			cl.Register.SetAllLevels(*cx.Config.LogLevel)
		}
		nodeChan := make(chan *rpc.Server)
		walletChan := make(chan *wallet.Wallet)
		kill := make(chan struct{})
		if !*cx.Config.NodeOff {
			go func() {
				err = node.Main(cx, shutdownChan, kill, nodeChan, &wg)
				if err != nil {
					log <- cl.Error{"error starting node ", err,
						cl.Ine()}
				}
			}()
			cx.RPCServer = <-nodeChan
		}
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx.Config, cx.StateCfg,
					cx.ActiveNet, walletChan, kill, &wg)
				if err != nil {
					fmt.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-walletChan
		}
		wg.Wait()
		return nil
	}
}
