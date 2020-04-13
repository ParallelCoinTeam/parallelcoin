package app

import (
	"fmt"
	"github.com/p9c/pod/app/config"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name)
		Debug("starting shell")
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
			// log.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				Error("failed to create wallet", err)
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(1)
		}
		Warn("starting node")
		if !*cx.Config.NodeOff {
			go func() {
				err = node.Main(cx, shutdownChan)
				if err != nil {
					Error("error starting node ", err)
				}
			}()
			cx.RPCServer = <-cx.NodeChan
		}
		Warn("starting wallet")
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx)
				if err != nil {
					fmt.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-cx.WalletChan
		}
		Debug("shell started")
		cx.WaitGroup.Wait()
		return nil
	}
}
