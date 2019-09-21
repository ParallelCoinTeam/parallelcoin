package app

import (
	"fmt"
	"sync"

	"github.com/urfave/cli"

	"github.com/parallelcointeam/parallelcoin/app/apputil"
	"github.com/parallelcointeam/parallelcoin/cmd/node"
	"github.com/parallelcointeam/parallelcoin/cmd/node/rpc"
	"github.com/parallelcointeam/parallelcoin/cmd/walletmain"
	"github.com/parallelcointeam/parallelcoin/pkg/conte"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
	"github.com/parallelcointeam/parallelcoin/pkg/wallet"
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
					ERROR("error starting node ", err)
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
