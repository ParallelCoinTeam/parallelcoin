package serve

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/Entry"
	"github.com/p9c/pod/pkg/pipe"
	"go.uber.org/atomic"
)

func Log(quit chan struct{}) {
	Debug("starting log server")
	lc := logi.L.AddLogChan()
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(quit, func(b []byte) (err error) {
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
				logi.L.SetLevel(logi.Tags[logi.Levels[b[4]]], false, "pod")
			}
		}
		return
	})
	go func() {
	out:
		for {
			select {
			case <-quit:
				break out
			case e := <-lc:
				if logOn.Load() {
					if n, err := p.Write(Entry.Get(&e).Data); !Check(err) {
						if n < 1 {
							Error("short write")
						}
					}
				}
			}
		}
	}()
}
