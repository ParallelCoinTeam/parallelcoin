package app

import (
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/conte"
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/wallet"
)

func shellHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		Configure(cx, c)
		L.Debug("starting shell")
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
			// L.SetLevel("off", false)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				L.Error("failed to create wallet", err)
			}
			log.Println("restart to complete initial setup")
			os.Exit(1)
		}
		L.Warn("starting node")
		if !*cx.Config.NodeOff {
			go func() {
				err = node.Main(cx, shutdownChan)
				if err != nil {
					L.Error("error starting node ", err)
				}
			}()
			cx.RPCServer = <-cx.NodeChan
		}
		L.Warn("starting wallet")
		if !*cx.Config.WalletOff {
			go func() {
				err = walletmain.Main(cx)
				if err != nil {
					log.Println("error running wallet:", err)
				}
			}()
			cx.WalletServer = <-cx.WalletChan
		}
		L.Debug("shell started")
		cx.WaitGroup.Wait()
		return nil
	}
}
