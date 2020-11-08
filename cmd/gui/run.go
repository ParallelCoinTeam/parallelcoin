package gui

import (
	"os"
	"runtime"

	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (wg *WalletGUI) Runner() (err error) {
	wg.NodeRunCommandChan = make(chan string)
	wg.MinerRunCommandChan = make(chan string)
	wg.MinerThreadsChan = make(chan int)
	interrupt.AddHandler(func() {
		if wg.running {
			// 		wg.NodeRunCommandChan <- "stop"
			consume.Kill(wg.Shell)
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
					Debug("run called")
					if wg.running {
						Debug("already running...")
						break
					}
					*wg.cx.Config.NodeOff = false
					*wg.cx.Config.WalletOff = false
					save.Pod(wg.cx.Config)
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
						"--rpclisten", *wg.cx.Config.RPCConnect,
						"--servertls=false", "--clienttls=false",
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
					wg.running = false
					*wg.cx.Config.NodeOff = true
					*wg.cx.Config.WalletOff = true
					save.Pod(wg.cx.Config)
					if wg.mining {
						go func(){
							wg.MinerRunCommandChan <- "stop"
						}()
					}
				case "restart":
					Debug("restart called")
					go func() {
						wg.NodeRunCommandChan <- "stop"
						wg.NodeRunCommandChan <- "run"
					}()
				}
			case cmd := <-wg.MinerRunCommandChan:
				switch cmd {
				case "run":
					Debug("run called for miner")
					if wg.running == false {
						Debug("not running because shell is not running")
						wg.mining = false
						break
					}
					wg.mining = true
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
				case "stop":
					Debug("stop called for miner")
					wg.mining = false
					consume.Kill(wg.Miner)
					*wg.cx.Config.Generate = false
					save.Pod(wg.cx.Config)
				case "restart":
					Debug("restart called for miner")
					go func(){
						wg.MinerRunCommandChan <- "stop"
						wg.MinerRunCommandChan <- "start"
					}()
				}
			case *wg.cx.Config.GenThreads = <-wg.MinerThreadsChan:
				Debug("setting threads to", *wg.cx.Config.GenThreads)
				if *wg.cx.Config.GenThreads == 0 {
					go func() {
						wg.MinerRunCommandChan <- "stop"
					}()
					break
				}
				if *wg.cx.Config.GenThreads < 0 {
					*wg.cx.Config.GenThreads = runtime.NumCPU()
				}
				if wg.mining {
					go func() {
						wg.MinerRunCommandChan <- "restart"
					}()
				}
			case <-wg.quit:
				Debug("runner received quit signal")
				consume.Kill(wg.Shell)
				break out
			}
		}
	}()
	if !(*wg.cx.Config.NodeOff && *wg.cx.Config.WalletOff) {
		wg.NodeRunCommandChan <- "run"
	}
	if *wg.cx.Config.Generate && *wg.cx.Config.GenThreads > 0 {
		wg.MinerRunCommandChan <- "run"
	}
	return nil
}
