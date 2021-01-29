package serve

import (
	"github.com/niubaoshu/gotiny"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/util/interrupt"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/util/logi"
)

func Log(quit qu.C, appName string) {
	Debug("starting log server")
	lc := logi.L.AddLogChan()
	// interrupt.AddHandler(func(){
	// 	// logi.L.RemoveLogChan(lc)
	// })
	// pkgChan := make(chan Pk.Package)
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(
		quit, func(b []byte) (err error) {
			// listen for commands to enable/disable logging
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "run ":
					Debug("setting to run")
					logOn.Store(true)
				case "stop":
					Debug("stopping")
					logOn.Store(false)
				case "slvl":
					Debug("setting level", logi.Levels[b[4]])
					logi.L.SetLevel(logi.Levels[b[4]], false, "pod")
				// case "pkgs":
				// 	pkgs := Pkg.LoadContainer(b).GetPackages()
				// 	for i := range pkgs {
				// 		(*logi.L.Packages)[i] = pkgs[i]
				// 	}
				// 	save settings
				// if !saveFunc(pkgs) {
				// 	Error("failed to save log filter configuration")
				// }
				case "kill":
					Debug("received kill signal from pipe, shutting down", appName)
					// time.Sleep(time.Second*5)
					// time.Sleep(time.Second * 3)
					quit.Q()
					// logi.L.LogChanDisabled = true
					// logi.L.LogChan = nil
					interrupt.Request()
					// <-interrupt.HandlersDone
					
					// quit.Q()
					// goroutineDump()
					// Debug(interrupt.GoroutineDump())
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
			case <-quit:
				// interrupt.Request()
				if !logi.L.LogChanDisabled.Load() {
					logi.L.LogChanDisabled.Store(true)
				}
				logi.L.Writer.Write.Store(true)
				Debug("quitting pipe logger") // , interrupt.GoroutineDump())
				interrupt.Request()
				logOn.Store(false)
				// <-interrupt.HandlersDone
				break out
			case e := <-lc:
				if !logOn.Load() {
					break out
				}
				if n, err := p.Write(gotiny.Marshal(e)); !Check(err) {
					// Debug(interrupt.GoroutineDump())
					if n < 1 {
						Error("short write")
					}
				} else {
					break out
					// 	quit.Q()
				}
				// case pk := <-pkgChan:
				// 	if !logOn.Load() {
				// 		break out
				// 	}
				// 	if n, err := p.Write(Pkg.Get(pk).Data); !Check(err) {
				// 		if n < 1 {
				// 			Error("short write")
				// 		}
				// 	} else {
				// 		break out
				// 	}
			}
		}
		
		<-interrupt.HandlersDone
		Debug("finished pipe logger")
	}()
}
