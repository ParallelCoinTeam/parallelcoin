package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/constant"
	"github.com/p9c/pod/pkg/opts"
	"github.com/p9c/pod/pkg/pod"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/p9c/qu"
	
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/apputil"
	"github.com/p9c/pod/pkg/podconfig"
)

// WalletHandle runs the wallet server
func WalletHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(255, 255, 128, false).Sprint
	log.App = "wallet"
	podconfig.Configure(cx, true)
	cx.Config.WalletFile.Set(filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, constant.DbName))
	// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
	// 	Params.Name + slash + wallet.WalletDbName
	if !apputil.FileExists(cx.Config.WalletFile.V()) && !cx.IsGUI {
		// D.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
		if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); opts.E.Chk(e) {
			opts.E.Ln("failed to create wallet", e)
			return e
		}
		fmt.Println("restart to complete initial setup")
		os.Exit(0)
	}
	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
	walletPassPath := filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, "wp.txt")
	opts.D.Ln("reading password from", walletPassPath)
	if apputil.FileExists(walletPassPath) {
		var b []byte
		if b, e = ioutil.ReadFile(walletPassPath); !opts.E.Chk(e) {
			cx.Config.WalletPass.SetBytes(b)
			opts.D.Ln("read password '" + string(b) + "'")
			for i := range b {
				b[i] = 0
			}
			if e = ioutil.WriteFile(walletPassPath, b, 0700); opts.E.Chk(e) {
			}
			if e = os.Remove(walletPassPath); opts.E.Chk(e) {
			}
			opts.D.Ln("wallet cookie deleted", *cx.Config.WalletPass)
		}
	}
	cx.WalletKill = qu.T()
	if e = walletmain.Main(cx); opts.E.Chk(e) {
		opts.E.Ln("failed to start up wallet", e)
	}
	// if !*cx.Config.DisableRPC {
	// 	cx.WalletServer = <-cx.WalletChan
	// }
	// cx.WaitGroup.Wait()
	cx.WaitWait()
	return
}
