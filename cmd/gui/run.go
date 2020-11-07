package gui

import (
	"os"
	"time"

	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (wg *WalletGUI) Runner() (err error) {
	wg.NodeRunCommandChan = make(chan string)
	interrupt.AddHandler(func() {
		if wg.running {
			// 		wg.NodeRunCommandChan <- "stop"
			consume.Kill(wg.Shell)
		}
		close(wg.quit)
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
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir,
						"--rpclisten", *wg.cx.Config.RPCConnect,
						"--servertls=false", "--clienttls=false",
						"--pipelog", "shell"}
					// args = apputil.PrependForWindows(args)
					wg.runnerQuit = make(chan struct{})
					wg.Shell = consume.Log(wg.runnerQuit, func(ent *logi.Entry) (err error) {
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
				case "restart":
					Debug("restart called")
					go func() {
						wg.NodeRunCommandChan <- "stop"
						time.Sleep(time.Second)
						wg.NodeRunCommandChan <- "run"
					}()
				}
				case cmd := <-wg.MinerRunCommandChan:
					switch cmd {
					case "run":
						Debug("run called for miner")

					case "stop":
						Debug("stop called for miner")

					case "restart":
						Debug("restart called for miner")
						
					}
			case <-wg.quit:
				Debug("runner received quit signal")
				consume.Kill(wg.Shell)
				break out
			}
		}
	}()
	wg.NodeRunCommandChan <- "run"
	if *wg.cx.Config.Generate && *wg.cx.Config.GenThreads > 0 {
		wg.MinerRunCommandChan <- "run"
	}
	return nil
}
