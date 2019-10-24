package cmd

import (
	"fmt"
	// This enables pprof
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/p9c/pod/app"
	"github.com/p9c/pod/pkg/util/limits"
)

// Main is the main entry point for pod
func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(50)
	if err := limits.SetLimits(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	os.Exit(app.Main())
	/*
		_=func() {
			f, err := os.Create("trace.out")
			if err != nil {
				panic(err)
			}
			err = trace.Start(f)
			if err != nil {
				panic(err)
			}
			mf, err := os.Create("mem.prof")
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
		}
		go func() {
			time.Sleep(time.Minute)
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(mf); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}()
		cf, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(cf); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
		interrupt.AddHandler(
			func() {
				fmt.Println("stopping trace")
				trace.Stop()
				pprof.StopCPUProfile()
				f.Close()
				mf.Close()
			},
		)
	*/
}
