package gui

import (
	"os"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (wg *WalletGUI) Runner() (err error) {
	wg.ShellRunCommandChan = make(chan string)
	wg.MinerRunCommandChan = make(chan string)
	interrupt.AddHandler(func() {
		if wg.running {
			// 		wg.ShellRunCommandChan <- "stop"
			consume.Kill(wg.Shell)
		}
		// close(wg.quit)
	})
	go func() {
		Debug("starting node run controller")
	out:
		for {
			select {
			case cmd := <-wg.ShellRunCommandChan:
				switch cmd {
				case "run":
					Debug("run called")
					if wg.running {
						Debug("already running...")
						break
					}
					*wg.cx.Config.NodeOff = false
					*wg.cx.Config.WalletOff = !*wg.walletLocked
					*wg.cx.Config.Network = wg.cx.ActiveNet.Name
					save.Pod(wg.cx.Config)
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
						"--rpclisten", *wg.cx.Config.RPCConnect,
						"-n", wg.cx.ActiveNet.Name,
						"--servertls=true", "--clienttls=true",
						// "--noinitialload",
						"--notty",
						"--pipelog", "shell"}
					// args = apputil.PrependForWindows(args)
					wg.runnerQuit = make(chan struct{})
					wg.Shell = consume.Log(wg.runnerQuit, func(ent *logi.Entry) (err error) {
						// TODO: make a log view for this
						// Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(wg.Shell)
					wg.running = true
				case "stop":
					Debug("stop called")
					if !wg.running {
						Debug("wasn't running...")
						break
					}
					consume.Kill(wg.Shell)
					*wg.cx.Config.NodeOff = true
					*wg.cx.Config.WalletOff = true
					save.Pod(wg.cx.Config)
					// if wg.mining {
					// 	go func() {
					// 		wg.MinerRunCommandChan <- "stop"
					// 	}()
					// }
					wg.running = false
				case "restart":
					Debug("restart called")
					go func() {
						wg.ShellRunCommandChan <- "stop"
						wg.running = false
						wg.ShellRunCommandChan <- "run"
						wg.running = true
					}()
				}
			case cmd := <-wg.MinerRunCommandChan:
				switch cmd {
				case "run":
					Debug("run called for miner")
					// if wg.running == false {
					// 	Debug("not running because shell is not running")
					// 	wg.mining = false
					// 	break
					// }
					if *wg.cx.Config.GenThreads == 0 {
						wg.mining = false
						break
					}
					*wg.cx.Config.Generate = true
					save.Pod(wg.cx.Config)
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir, "--pipelog", "kopach"}
					// args = apputil.PrependForWindows(args)
					wg.minerQuit = make(chan struct{})
					wg.Miner = consume.Log(wg.minerQuit, func(ent *logi.Entry) (err error) {
						// TODO: make a log view for this
						// Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(wg.Miner)
					wg.mining = true
				case "stop":
					Debug("stop called for miner")
					consume.Kill(wg.Miner)
					*wg.cx.Config.Generate = false
					save.Pod(wg.cx.Config)
					wg.mining = false
				case "restart":
					Debug("restart called for miner")
					go func() {
						wg.MinerRunCommandChan <- "stop"
						wg.mining = false
						wg.MinerRunCommandChan <- "run"
						wg.mining = true
					}()
				}
			case <-wg.quit:
				Debug("runner received quit signal")
				consume.Kill(wg.Shell)
				break out
			}
		}
	}()
	// if wg.running {
	// 	Debug("starting shell")
	// 	wg.ShellRunCommandChan <- "run"
	// }
	// if wg.mining {
	// 	Debug("starting miner")
	// 	wg.MinerRunCommandChan <- "run"
	// }
	return nil
}
