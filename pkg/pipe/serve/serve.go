package serve

import (
	"github.com/niubaoshu/gotiny"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/util/logi"
)

// Log starts up a handler to listen to logs from the child process worker
func Log(quit qu.C, appName string) {
	dbg.Ln("starting log server")
	lc := logi.L.AddLogChan()
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
					dbg.Ln("setting to run")
					logOn.Store(true)
				case "stop":
					dbg.Ln("stopping")
					logOn.Store(false)
				case "slvl":
					dbg.Ln("setting level", logi.Levels[b[4]])
					logi.L.SetLevel(logi.Levels[b[4]], false, "pod")
				case "kill":
					dbg.Ln("received kill signal from pipe, shutting down", appName)
					// time.Sleep(time.Second*5)
					// time.Sleep(time.Second * 3)
					// logi.L.LogChanDisabled = true
					// logi.L.LogChan = nil
					interrupt.Request()
					quit.Q()
					// <-interrupt.HandlersDone
					
					// quit.Q()
					// goroutineDump()
					// dbg.Ln(interrupt.GoroutineDump())
					// pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
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
				if !logi.L.LogChanDisabled.Load() {
					logi.L.LogChanDisabled.Store(true)
				}
				logi.L.Writer.Write.Store(true)
				dbg.Ln("quitting pipe logger") // , interrupt.GoroutineDump())
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
				if n, e = p.Write(gotiny.Marshal(&ent)); !err.Chk(e) {
					// dbg.Ln(interrupt.GoroutineDump())
					if n < 1 {
						err.Ln("short write")
					}
				} else {
					break out
					// 	quit.Q()
				}
			}
		}
		<-interrupt.HandlersDone
		dbg.Ln("finished pipe logger")
	}()
}
