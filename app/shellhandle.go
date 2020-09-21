package app

import (
	"fmt"
	"github.com/stalker-loki/app/slog"
	"github.com/stalker-loki/pod/app/config"
	"os"

	"github.com/urfave/cli"

	"github.com/stalker-loki/pod/app/apputil"
	"github.com/stalker-loki/pod/app/conte"
	"github.com/stalker-loki/pod/cmd/node"
	"github.com/stalker-loki/pod/cmd/walletmain"
	"github.com/stalker-loki/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name, true)
		slog.Debug("starting shell")
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
				slog.Error("failed to create wallet", err)
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(1)
		}
		slog.Warn("starting node")
		if !*cx.Config.NodeOff {
			go func() {
				err = node.Main(cx, shutdownChan)
				if err != nil {
					slog.Error("error starting node ", err)
				}
			}()
			cx.RPCServer = <-cx.NodeChan
		}
		slog.Warn("starting wallet")
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx)
				if err != nil {
					fmt.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-cx.WalletChan
		}
		slog.Debug("shell started")
		cx.WaitGroup.Wait()
		return nil
	}
}
