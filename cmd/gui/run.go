package gui

import (
	"os"
	"time"

	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (wg *WalletGUI) Runner() (err error) {
	interrupt.AddHandler(func() {
		if wg.running {
			// 		wg.RunCommandChan <- "stop"
			consume.Kill(wg.Worker)
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
					wg.Worker = consume.Log(wg.quit, func(ent *logi.Entry) (err error) {
						Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					}, os.Args[0], "-D", *wg.cx.Config.DataDir, "--pipelog", "shell")
					consume.Start(wg.Worker)
					wg.running = true
					// go func() {
					// 	if err = wg.Worker.Wait(); !Check(err) {
					// 		wg.running = false
					// 	}
					// }()
				case "stop":
					if !wg.running {
						break
					}
					Debug("stop called")
					// consume.Stop(wg.Worker)
					consume.Kill(wg.Worker)
					// if err = wg.Worker.Stop(); !Check(err) {
					// 	Debug("stopped worker")
					// }
					// if err = wg.Worker.Interrupt(); !Check(err) {
					// 	Debug("interrupted worker")
					// } else {
					// 	Debug(err)
					// }
					// if err = wg.Worker.Kill(); !Check(err) {
					// 	Debug("killed worker")
					// } else {
					// 	Debug(err)
					// }
					// if err = wg.Worker.StdConn.Close(); !Check(err) {
					// 	Debug("closed worker connection")
					// } else {
					// 	Debug(err)
					// }
					// go func() {
					// time.Sleep(time.Second * 4)
					wg.running = false
					// }()
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
				consume.Kill(wg.Worker)
				break out
			}
		}
	}()
	return nil
}
