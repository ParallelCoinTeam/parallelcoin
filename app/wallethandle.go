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
		logg.App=color.Bit24(255,255,0,false).Sprint("wallet")
		config.Configure(cx, c.Command.Name, true)
		*cx.Config.WalletFile = *cx.Config.DataDir + string(os.PathSeparator) +
			cx.ActiveNet.Name + string(os.PathSeparator) + wallet.DbName
		// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
		// 	Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(*cx.Config.WalletFile) && !cx.IsGUI {
			// dbg.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
			if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); err.Chk(e) {
				err.Ln("failed to create wallet", e)
				return e
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
		}
		// for security with apps launching the wallet, the public password can be set with a file that is deleted after
		walletPassPath := *cx.Config.DataDir + slash + cx.ActiveNet.Params.Name + slash + "wp.txt"
		dbg.Ln("reading password from", walletPassPath)
		if apputil.FileExists(walletPassPath) {
			var b []byte
			if b, e = ioutil.ReadFile(walletPassPath); !err.Chk(e) {
				*cx.Config.WalletPass = string(b)
				dbg.Ln("read password '" + string(b) + "'")
				for i := range b {
					b[i] = 0
				}
				if e = ioutil.WriteFile(walletPassPath, b, 0700); err.Chk(e) {
				}
				if e = os.Remove(walletPassPath); err.Chk(e) {
				}
				dbg.Ln("wallet cookie deleted", *cx.Config.WalletPass)
			}
		}
		cx.WalletKill = qu.T()
		if e = walletmain.Main(cx); err.Chk(e) {
			err.Ln("failed to start up wallet", e)
		}
		// if !*cx.Config.DisableRPC {
		// 	cx.WalletServer = <-cx.WalletChan
		// }
		// cx.WaitGroup.Wait()
		cx.WaitWait()
		return
	}
}
