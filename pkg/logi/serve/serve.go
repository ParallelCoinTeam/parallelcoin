package serve

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/Entry"
	"github.com/p9c/pod/pkg/pipe"
	"go.uber.org/atomic"
)

func Log(quit chan struct{}) {
	L.Debug("starting log server")
	lc := logi.L.AddLogChan()
	var logOn atomic.Bool
	logOn.Store(false)
	p := pipe.Serve(quit, func(b []byte) (err error) {
		// listen for commands to enable/disable logging
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "run ":
				L.Debug("setting to run")
				logOn.Store(true)
			case "stop":
				L.Debug("stopping")
				logOn.Store(false)
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
					if n, err := p.Write(Entry.Get(&e).Data); !L.Check(err) {
						if n < 1 {
							L.Error("short write")
						}
					}
				}
			}
		}
	}()
}
