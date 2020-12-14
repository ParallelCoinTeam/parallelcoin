package serve

import (
	"github.com/p9c/pod/pkg/util/interrupt"
	"go.uber.org/atomic"
	
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/Entry"
	"github.com/p9c/pod/pkg/util/logi/Pkg"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
)

func Log(quit chan struct{}, saveFunc func(p Pk.Package) (success bool), appName string) {
	Debug("starting log server")
	lc := logi.L.AddLogChan()
	logQuit := make(chan struct{})
	interrupt.AddHandler(
		func() {
			close(logQuit)
		},
	)
	pkgChan := make(chan Pk.Package)
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(
		logQuit, func(b []byte) (err error) {
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
				case "pkgs":
					pkgs := Pkg.LoadContainer(b).GetPackages()
					for i := range pkgs {
						(*logi.L.Packages)[i] = pkgs[i]
					}
					// save settings
					if !saveFunc(pkgs) {
						Error("failed to save log filter configuration")
					}
				case "kill":
					Debug("received kill signal from pipe, shutting down", appName)
					// time.Sleep(time.Second*5)
					// close(quit)
					// time.Sleep(time.Second * 3)
					interrupt.Request()
					// pprof.Lookup("goroutine").WriteTo(os.Stderr, 2)
					
					// os.Exit(0)
					// break
					// os.Exit(0)
				}
			}
			return
		},
	)
	go func() {
	out:
		for {
			select {
			case <-logQuit:
				Debug("quitting pipe logger")
				break out
			case e := <-lc:
				if logOn.Load() {
					if n, err := p.Write(Entry.Get(&e).Data); !Check(err) {
						if n < 1 {
							Error("short write")
						}
					}
				}
			case pk := <-pkgChan:
				if n, err := p.Write(Pkg.Get(pk).Data); !Check(err) {
					if n < 1 {
						Error("short write")
					}
				}
			}
		}
	}()
}
