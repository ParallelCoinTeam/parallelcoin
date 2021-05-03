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
	var e error
	if runtime.GOOS != "darwin" {
		if e = limits.SetLimits(); E.Chk(e) { // todo: doesn't work on non-linux
			_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", e)
			os.Exit(1)
		}
	}
	var f *os.File
	if os.Getenv("POD_TRACE") == "on" {
		D.Ln("starting trace")
		if f, e = os.Create(fmt.Sprintf("%v.trace", fmt.Sprint(os.Args))); E.Chk(e) {
			E.Ln(
				"tracing env POD_TRACE=on but we can't write to it",
				e,
			)
		} else {
			e = trace.Start(f)
			if e != nil {
				E.Ln("could not start tracing", e)
			} else {
				D.Ln("tracing started")
				defer trace.Stop()
				defer func() {
					if e := f.Close(); E.Chk(e) {
					}
				}()
				interrupt.AddHandler(
					func() {
						D.Ln("stopping trace")
						trace.Stop()
						e := f.Close()
						if e != nil {
						}
					},
				)
			}
		}
	}
	res := app.Main()
	D.Ln("returning value", res, os.Args)
	if os.Getenv("POD_TRACE") == "on" {
		D.Ln("stopping trace")
		trace.Stop()
		defer func() {
			if e := f.Close(); E.Chk(e) {
			}
		}()
	}
	if res != 0 {
		E.Ln("quitting with error")
		// D.Ln(interrupt.GoroutineDump())
		os.Exit(res)
	}
}
