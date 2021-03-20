package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/logg"
	"io/ioutil"
	"os"
	
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/config"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func WalletHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.AppColorizer = color.Bit24(255, 255, 128, false).Sprint
		logg.App = "wallet"
		config.Configure(cx, c.Command.Name, true)
		*cx.Config.WalletFile = *cx.Config.DataDir + string(os.PathSeparator) +
			cx.ActiveNet.Name + string(os.PathSeparator) + wallet.DbName
		// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
		// 	Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(*cx.Config.WalletFile) && !cx.IsGUI {
			// D.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
			if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); E.Chk(e) {
				E.Ln("failed to create wallet", e)
				return e
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
		}
		// for security with apps launching the wallet, the public password can be set with a file that is deleted after
		walletPassPath := *cx.Config.DataDir + slash + cx.ActiveNet.Params.Name + slash + "wp.txt"
		D.Ln("reading password from", walletPassPath)
		if apputil.FileExists(walletPassPath) {
			var b []byte
			if b, e = ioutil.ReadFile(walletPassPath); !E.Chk(e) {
				*cx.Config.WalletPass = string(b)
				D.Ln("read password '" + string(b) + "'")
				for i := range b {
					b[i] = 0
				}
				if e = ioutil.WriteFile(walletPassPath, b, 0700); E.Chk(e) {
				}
				if e = os.Remove(walletPassPath); E.Chk(e) {
				}
				D.Ln("wallet cookie deleted", *cx.Config.WalletPass)
			}
		}
		cx.WalletKill = qu.T()
		if e = walletmain.Main(cx); E.Chk(e) {
			E.Ln("failed to start up wallet", e)
		}
		// if !*cx.Config.DisableRPC {
		// 	cx.WalletServer = <-cx.WalletChan
		// }
		// cx.WaitGroup.Wait()
		cx.WaitWait()
		return
	}
}
