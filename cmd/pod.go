package cmd

import (
	"fmt"
	// This enables pprof
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"

	"github.com/stalker-loki/pod/pkg/util/interrupt"

	"github.com/stalker-loki/pod/app"
	"github.com/stalker-loki/pod/pkg/util/limits"
)

// Main is the main entry point for pod
func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 3)
	debug.SetGCPercent(10)
	if runtime.GOOS != "darwin" {
		if err := limits.SetLimits(); err != nil { // todo: doesn't work on non-linux
			_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
			os.Exit(1)
		}
	}
	if os.Getenv("POD_TRACE") == "on" {
		if f, err := os.Create("testtrace.out"); err != nil {
			Error("tracing env POD_TRACE=on but we can't write to it",
				err)
		} else {
			Debug("tracing started")
			err = trace.Start(f)
			if err != nil {
				Error("could not start tracing", err)
			} else {
				interrupt.AddHandler(func() {
					Debug("stopping trace")
					trace.Stop()
					err := f.Close()
					if err != nil {
						Error(err)
					}
				},
				)
			}
		}
	}
	app.Main()
}
