package gui

import (
	"io/ioutil"
	"os"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

const slash = string(os.PathSeparator)

func (wg *WalletGUI) Runner() (err error) {
	wg.NodeRunCommandChan = make(chan string)
	wg.WalletRunCommandChan = make(chan string)
	wg.MinerRunCommandChan = make(chan string)
	interrupt.AddHandler(func() {
		if wg.runningNode.Load() {
					wg.NodeRunCommandChan <- "stop"
			// consume.Kill(wg.Node)
		}
		if wg.runningWallet.Load() {
			wg.WalletRunCommandChan <- "stop"
			// consume.Kill(wg.Node)
		}
		// close(wg.quit)
	})
	go func() {
		Debug("starting node run controller")
	out:
		for {
			select {
			case cmd := <-wg.NodeRunCommandChan:
				switch cmd {
				case "run":
					wg.nodeQuit = make(chan struct{})
					Debug("run called")
					if wg.runningNode.Load() {
						Debug("already running...")
						break
					}
					// wp := *wg.cx.Config.WalletPass
					*wg.cx.Config.NodeOff = false
					// *wg.cx.Config.WalletOff = *wg.walletLocked
					// todo: if locked shouldn't pass be zeroed?
					*wg.cx.Config.Network = wg.cx.ActiveNet.Name
					save.Pod(wg.cx.Config)
					// if !*wg.cx.Config.WalletOff {
					// 	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
					// 	walletPassPath := *wg.cx.Config.DataDir + slash + wg.cx.ActiveNet.Params.Name + slash + "wp.txt"
					// 	Debug("runner", walletPassPath)
					// 	b := []byte(wp)
					// 	if err = ioutil.WriteFile(walletPassPath, b, 0700); Check(err) {
					// 	}
					// 	Debug("created password cookie")
					// }
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
						// "--rpclisten", *wg.cx.Config.RPCConnect,
						// "-n", wg.cx.ActiveNet.Name,
						"--servertls=true", "--clienttls=true",
						// "--noinitialload",
						// "--runasservice",
						// "--notty",
						"--pipelog", "node"}
					// args = apputil.PrependForWindows(args)
					wg.Node = consume.Log(wg.nodeQuit, func(ent *logi.Entry) (err error) {
						// TODO: make a log view for this
						Debugf("NODE[%s] %s %s",
							ent.Level,
							// ent.Time.Format(time.RFC3339),
							ent.Text,
							ent.CodeLocation,
						)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(wg.Node)
					wg.runningNode.Store(true)
				case "stop":
					Debug("stop called")
					// if !wg.runningNode.Load() {
					// 	Debug("wasn't running...")
					// 	break
					// }
					if wg.runningWallet.Load() {
						wg.WalletRunCommandChan <- "stop"
						wg.walletLocked.Store(true)
					}
					consume.Kill(wg.Node)
					*wg.cx.Config.NodeOff = true
					save.Pod(wg.cx.Config)
					wg.runningNode.Store(false)
				case "restart":
					Debug("restart called")
					go func() {
						wg.NodeRunCommandChan <- "stop"
						// wg.running = false
						wg.NodeRunCommandChan <- "run"
						// wg.running = true
					}()
				}
			case cmd := <-wg.WalletRunCommandChan:
				switch cmd {
				case "run":
					wg.walletQuit = make(chan struct{})
					Debug("run called")
					if wg.runningWallet.Load() {
						Debug("already running...")
						break
					}
					wp := *wg.cx.Config.WalletPass
					// *wg.cx.Config.WalletOff = false
					// *wg.cx.Config.WalletOff = *wg.walletLocked
					// todo: if locked shouldn't pass be zeroed?
					*wg.cx.Config.Network = wg.cx.ActiveNet.Name
					save.Pod(wg.cx.Config)
					if !*wg.cx.Config.WalletOff {
						// for security with apps launching the wallet, the public password can be set with a file that is deleted after
						walletPassPath := *wg.cx.Config.DataDir + slash + wg.cx.ActiveNet.Params.Name + slash + "wp.txt"
						Debug("runner", walletPassPath)
						b := []byte(wp)
						if err = ioutil.WriteFile(walletPassPath, b, 0700); Check(err) {
						}
						Debug("created password cookie")
					}
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
						// "--rpcconnect", *wg.cx.Config.RPCConnect,
						// "-n", wg.cx.ActiveNet.Name,
						"--servertls=true", "--clienttls=true",
						// "--noinitialload",
						// "--runasservice",
						// "--notty",
						"--pipelog", "wallet"}
					// args = apputil.PrependForWindows(args)
					wg.Wallet = consume.Log(wg.walletQuit, func(ent *logi.Entry) (err error) {
						// TODO: make a log view for this
						Debugf("WALLET[%s] %s %s",
							ent.Level,
							// ent.Time.Format(time.RFC3339),
							ent.Text,
							ent.CodeLocation,
						)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(wg.Wallet)
					wg.walletLocked.Store(false)
					wg.runningWallet.Store(true)
				case "stop":
					Debug("stop called")
					if !wg.runningWallet.Load() {
						Debug("wasn't running...")
						break
					}
					consume.Kill(wg.Wallet)
					*wg.cx.Config.WalletOff = true
					save.Pod(wg.cx.Config)
					wg.walletLocked.Store(true)
					wg.runningWallet.Store(false)
				case "restart":
					Debug("restart called")
					go func() {
						wg.WalletRunCommandChan <- "stop"
						// wg.running = false
						wg.WalletRunCommandChan <- "run"
						// wg.running = true
					}()
				}
			case cmd := <-wg.MinerRunCommandChan:
				switch cmd {
				case "run":
					Debug("run called for miner")
					if *wg.cx.Config.GenThreads == 0 {
						wg.mining.Store(false)
						break
					}
					*wg.cx.Config.Generate = true
					save.Pod(wg.cx.Config)
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir, "--pipelog", "kopach"}
					// args = apputil.PrependForWindows(args)
					wg.minerQuit = make(chan struct{})
					wg.Miner = consume.Log(wg.minerQuit, func(ent *logi.Entry) (err error) {
						// TODO: make a log view for this
						Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(wg.Miner)
					wg.mining.Store(true)
				case "stop":
					Debug("stop called for miner")
					consume.Kill(wg.Miner)
					*wg.cx.Config.Generate = false
					save.Pod(wg.cx.Config)
					wg.mining.Store(false)
				case "restart":
					Debug("restart called for miner")
					go func() {
						wg.MinerRunCommandChan <- "stop"
						wg.mining.Store(false)
						wg.MinerRunCommandChan <- "run"
						wg.mining.Store(true)
					}()
				}
			case <-wg.quit:
				Debug("runner received quit signal")
				wg.NodeRunCommandChan <- "stop"
				wg.MinerRunCommandChan <- "stop"
				// consume.Kill(wg.Miner)
				// consume.Kill(wg.Node)
				break out
			}
		}
	}()
	return nil
}
