package explorer

import (
	"os"
	"time"

	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/consume"
)

func (ex *Explorer) Runner() (err error) {
	interrupt.AddHandler(func() {
		if ex.running {
			// 		ex.RunCommandChan <- "stop"
			consume.Kill(ex.Shell)
		}
		close(ex.quit)
	})
	go func() {
		Debug("starting node run controller")
	out:
		for {
			select {
			case cmd := <-ex.RunCommandChan:
				switch cmd {
				case "run":
					if ex.running {
						break
					}
					Debug("run called")
					args := []string{os.Args[0], "-D", *ex.cx.Config.DataDir,
						"--rpclisten", *ex.cx.Config.RPCConnect,
						"--servertls=false", "--clienttls=false",
						"--pipelog", "shell"}
					// args = apputil.PrependForWindows(args)
					ex.Shell = consume.Log(ex.quit, func(ent *logi.Entry) (err error) {
						// Debug(ent.Level, ent.Time, ent.Text, ent.CodeLocation)
						return
					}, func(pkg string) (out bool) {
						return false
					}, args...)
					consume.Start(ex.Shell)
					ex.running = true

				case "stop":
					if !ex.running {
						break
					}
					Debug("stop called")
					consume.Kill(ex.Shell)
					ex.running = false
				case "restart":
					Debug("restart called")
					go func() {
						ex.RunCommandChan <- "stop"
						time.Sleep(time.Second)
						ex.RunCommandChan <- "run"
					}()
				}
			case <-ex.quit:
				Debug("runner received quit signal")
				consume.Kill(ex.Shell)
				break out
			}
		}
	}()
	return nil
}
