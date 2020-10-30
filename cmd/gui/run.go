package gui

import (
	"os"
	"time"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (wg *WalletGUI) Runner() (err error) {
	interrupt.AddHandler(func() {
		if wg.running {
			// 		wg.RunCommandChan <- "stop"
			consume.Kill(wg.Shell)
		}
		close(wg.quit)
	})
	go func() {
		Debug("starting node run controller")
	out:
		for {
			select {
			case cmd := <-wg.RunCommandChan:
				switch cmd {
				case "run":
					if wg.running {
						break
					}
					Debug("run called")
					args := []string{os.Args[0], "-D", *wg.cx.Config.DataDir, "--pipelog", "shell"}
					args = apputil.PrependForWindows(args)
					wg.Shell = consume.Log(wg.quit, func(ent *logi.Entry) (err error) {
						Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					},args...)
					consume.Start(wg.Shell)
					wg.running = true
				case "stop":
					if !wg.running {
						break
					}
					Debug("stop called")
					consume.Kill(wg.Shell)
					wg.running = false
				case "restart":
					Debug("restart called")
					go func() {
						wg.RunCommandChan <- "stop"
						time.Sleep(time.Second)
						wg.RunCommandChan <- "run"
					}()
				}
			case <-wg.quit:
				Debug("runner received quit signal")
				consume.Kill(wg.Shell)
				break out
			}
		}
	}()
	return nil
}
