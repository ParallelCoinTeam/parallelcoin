package serve

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/pkg/logg"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/comm/pipe"
)

// Log starts up a handler to listen to logs from the child process worker
func Log(quit qu.C, appName string) {
	D.Ln("starting log server")
	lc := logg.AddLogChan()
	// interrupt.AddHandler(func(){
	// 	// logi.L.RemoveLogChan(lc)
	// })
	// pkgChan := make(chan Pk.Package)
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(
		quit, func(b []byte) (e error) {
			// listen for commands to enable/disable logging
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "run ":
					D.Ln("setting to run")
					logOn.Store(true)
				case "stop":
					D.Ln("stopping")
					logOn.Store(false)
				case "slvl":
					D.Ln("setting level", logg.Levels[b[4]])
					logg.SetLogLevel(logg.Levels[b[4]])
				case "kill":
					D.Ln("received kill signal from pipe, shutting down", appName)
					interrupt.Request()
					quit.Q()
				}
			}
			return
		},
	)
	go func() {
	out:
		for {
			select {
			case <-quit.Wait():
				// interrupt.Request()
				if !logg.LogChanDisabled.Load() {
					logg.LogChanDisabled.Store(true)
				}
				D.Ln("quitting pipe logger") // , interrupt.GoroutineDump())
				interrupt.Request()
				logOn.Store(false)
				// <-interrupt.HandlersDone
			out2:
				// drain log channel
				for {
					select {
					case <-lc:
						break
					default:
						break out2
					}
				}
				break out
			case ent := <-lc:
				if !logOn.Load() {
					break out
				}
				var n int
				var e error
				if n, e = p.Write(gotiny.Marshal(&ent)); !E.Chk(e) {
					// D.Ln(interrupt.GoroutineDump())
					if n < 1 {
						E.Ln("short write")
					}
				} else {
					break out
					// 	quit.Q()
				}
			}
		}
		<-interrupt.HandlersDone
		D.Ln("finished pipe logger")
	}()
}
