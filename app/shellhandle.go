package app

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/pod/pkg/logg"
	"io/ioutil"
	"os"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/config"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/app/conte"
	"github.com/p9c/pod/cmd/node"
	"github.com/p9c/pod/cmd/walletmain"
	"github.com/p9c/pod/pkg/wallet"
)

func ShellHandle(cx *conte.Xt) func(c *cli.Context) (e error) {
	return func(c *cli.Context) (e error) {
		logg.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
		logg.App = " shell"
		config.Configure(cx, c.Command.Name, true)
		D.Ln("starting shell")
		if *cx.Config.TLS || *cx.Config.ServerTLS {
			// generate the tls certificate if configured
			if apputil.FileExists(*cx.Config.RPCCert) && apputil.FileExists(*cx.Config.RPCKey) &&
				apputil.FileExists(*cx.Config.CAFile) {
				
			} else {
				_, _ = walletmain.GenerateRPCKeyPair(cx.Config, true)
			}
		}
		dbFilename :=
			*cx.Config.DataDir + slash +
				cx.ActiveNet.Params.Name + slash +
				wallet.DbName
		if !apputil.FileExists(dbFilename) && !cx.IsGUI {
			// log.SetLevel("off", false)
			if e := walletmain.CreateWallet(cx.ActiveNet, cx.Config); E.Chk(e) {
				E.Ln("failed to create wallet", e)
			}
			fmt.Println("restart to complete initial setup")
			os.Exit(1)
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
		if !*cx.Config.NodeOff {
			go func() {
				e = node.Main(cx)
				if e != nil {
					E.Ln("error starting node ", e)
				}
			}()
			I.Ln("starting node")
			if !*cx.Config.DisableRPC {
				cx.RPCServer = <-cx.NodeChan
			}
			I.Ln("node started")
		}
		if !*cx.Config.WalletOff {
			go func() {
				e = walletmain.Main(cx)
				if e != nil {
					fmt.Println("error running wallet:", e)
				}
			}()
			// I.Ln("starting wallet")
			// if !*cx.Config.DisableRPC {
			// 	cx.WalletServer = <-cx.WalletChan
			// }
			// I.Ln("wallet started")
		}
		D.Ln("shell started")
		// cx.WaitGroup.Wait()
		cx.WaitWait()
		return nil
	}
}
