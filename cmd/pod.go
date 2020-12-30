package cmd

import (
	"fmt"
	// This enables pprof
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	
	"github.com/p9c/pod/app"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/limits"
)

// Main is the main entry point for pod
func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 3)
	debug.SetGCPercent(10)
	var err error
	if runtime.GOOS != "darwin" {
		if err = limits.SetLimits(); err != nil { // todo: doesn't work on non-linux
			_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
			os.Exit(1)
		}
	}
	var f *os.File
	if os.Getenv("POD_TRACE") == "on" {
		Debug("starting trace")
		if f, err = os.Create(fmt.Sprintf("%v.trace", fmt.Sprint(os.Args))); err != nil {
			Error(
				"tracing env POD_TRACE=on but we can't write to it",
				err,
			)
		} else {
			err = trace.Start(f)
			if err != nil {
				Error("could not start tracing", err)
			} else {
				Debug("tracing started")
				defer trace.Stop()
				defer f.Close()
				interrupt.AddHandler(
					func() {
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
	res := app.Main()
	Debug("returning value", res, os.Args)
	if os.Getenv("POD_TRACE") == "on" {
		Debug("stopping trace")
		trace.Stop()
		f.Close()
	}
	if res != 0 {
		Warn("quitting with error")
		Debug(interrupt.GoroutineDump())
		os.Exit(res)
	}
}
