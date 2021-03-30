package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/opts"
	"github.com/p9c/pod/pkg/pod"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/p9c/qu"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/podconfig"
)

func WalletHandle(cx *pod.State) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		log.AppColorizer = color.Bit24(255, 255, 128, false).Sprint
		log.App = "wallet"
		podconfig.Configure(cx, c.Command.Name, true)
		cx.Config.WalletFile.Set(filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, opts.DbName))
		// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
		// 	Params.Name + slash + wallet.WalletDbName
		if !apputil.FileExists(cx.Config.WalletFile.V()) && !cx.IsGUI {
			// D.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
			if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); E.Chk(e) {
				E.Ln("failed to create wallet", e)
				return e
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(0)
		}
		// for security with apps launching the wallet, the public password can be set with a file that is deleted after
		walletPassPath := filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, "wp.txt")
		D.Ln("reading password from", walletPassPath)
		if apputil.FileExists(walletPassPath) {
			var b []byte
			if b, e = ioutil.ReadFile(walletPassPath); !E.Chk(e) {
				cx.Config.WalletPass.SetBytes(b)
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
