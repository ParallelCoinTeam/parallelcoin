package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func WalletHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		config.Configure(cx, c.Command.Name, true)
		*cx.Config.WalletFile = *cx.Config.DataDir + string(os.PathSeparator) +
			cx.ActiveNet.Name + string(os.PathSeparator) + wallet.WalletDbName
		// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
		// 	Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(*cx.Config.WalletFile) && !cx.IsGUI {
			// Debug(cx.ActiveNet.Name, *cx.Config.WalletFile)
			if err := walletmain.CreateWallet(cx.ActiveNet, cx.Config); err != nil {
				Error("failed to create wallet", err)
				return err
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
		}
		// for security with apps launching the wallet, the public password can be set with a file that is deleted after
		walletPassPath := *cx.Config.DataDir + slash + cx.ActiveNet.Params.Name + slash + "wp.txt"
		Debug("reading password from", walletPassPath)
		if apputil.FileExists(walletPassPath) {
			var b []byte
			if b, err = ioutil.ReadFile(walletPassPath); !Check(err) {
				*cx.Config.WalletPass = string(b)
				Debug("read password '" + string(b) + "'")
				for i := range b {
					b[i] = 0
				}
				if err = ioutil.WriteFile(walletPassPath, b, 0700); Check(err) {
				}
				if err = os.Remove(walletPassPath); Check(err) {
				}
				Debug("wallet cookie deleted", *cx.Config.WalletPass)
			}
		}
		cx.WalletKill = make(chan struct{})
		go func() {
			err = walletmain.Main(cx)
			if err != nil {
				Error("failed to start up wallet", err)
			}
		}()
		if !*cx.Config.DisableRPC {
			cx.WalletServer = <-cx.WalletChan
		}
		cx.WaitGroup.Wait()
		return
	}
}
