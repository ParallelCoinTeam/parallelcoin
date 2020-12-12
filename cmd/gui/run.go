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
	wg.node.RunCommandChan = make(chan string)
	wg.wallet.RunCommandChan = make(chan string)
	wg.miner.RunCommandChan = make(chan string)
	interrupt.AddHandler(func() {
		if wg.node.running.Load() {
					wg.node.RunCommandChan <- "stop"
			// consume.Kill(wg.Node)
		}
		if wg.wallet.running.Load() {
			wg.wallet.RunCommandChan <- "stop"
			// consume.Kill(wg.Node)
		}
		// close(wg.quit)
	})
	go func() {
		Debug("starting node run controller")
	out:
		for {
			select {
			case cmd := <-wg.node.RunCommandChan:
				switch cmd {
				case "run":
					wg.nodeQuit = make(chan struct{})
					Debug("run called")
					if wg.node.running.Load() {
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
					wg.node.running.Store(true)
				case "stop":
					Debug("stop called")
					// if !wg.runningNode.Load() {
					// 	Debug("wasn't running...")
					// 	break
					// }
					if wg.node.running.Load() {
						wg.wallet.RunCommandChan <- "stop"
						wg.walletLocked.Store(true)
					}
					consume.Kill(wg.Node)
					*wg.cx.Config.NodeOff = true
					save.Pod(wg.cx.Config)
					wg.node.running.Store(false)
				case "restart":
					Debug("restart called")
					go func() {
						wg.node.RunCommandChan <- "stop"
						// wg.running = false
						wg.node.RunCommandChan <- "run"
						// wg.running = true
					}()
				}
			case cmd := <-wg.wallet.RunCommandChan:
				switch cmd {
				case "run":
					wg.walletQuit = make(chan struct{})
					Debug("run called")
					if wg.wallet.running.Load() {
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
					wg.wallet.running.Store(true)
				case "stop":
					Debug("stop called")
					if !wg.wallet.running.Load() {
						Debug("wasn't running...")
						break
					}
					consume.Kill(wg.Wallet)
					*wg.cx.Config.WalletOff = true
					save.Pod(wg.cx.Config)
					wg.walletLocked.Store(true)
					wg.wallet.running.Store(false)
				case "restart":
					Debug("restart called")
					go func() {
						wg.wallet.RunCommandChan <- "stop"
						// wg.running = false
						wg.wallet.RunCommandChan <- "run"
						// wg.running = true
					}()
				}
			case cmd := <-wg.miner.RunCommandChan:
				switch cmd {
				case "run":
					Debug("run called for miner")
					if *wg.cx.Config.GenThreads == 0 {
						wg.miner.running.Store(false)
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
					wg.miner.running.Store(true)
				case "stop":
					Debug("stop called for miner")
					consume.Kill(wg.Miner)
					*wg.cx.Config.Generate = false
					save.Pod(wg.cx.Config)
					wg.miner.running.Store(false)
				case "restart":
					Debug("restart called for miner")
					go func() {
						wg.miner.RunCommandChan <- "stop"
						wg.miner.running.Store(false)
						wg.miner.RunCommandChan <- "run"
						wg.miner.running.Store(true)
					}()
				}
			case <-wg.quit:
				Debug("runner received quit signal")
				wg.node.RunCommandChan <- "stop"
				wg.miner.RunCommandChan <- "stop"
				// consume.Kill(wg.Miner)
				// consume.Kill(wg.Node)
				break out
			}
		}
	}()
	return nil
}
