package cmd

import (
	"fmt"
	"github.com/p9c/pod/pkg/util/interrupt"
	//// This enables pprof
	//_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/p9c/pod/app"
	"github.com/p9c/pod/pkg/util/limits"
)

// Main is the main entry point for pod
func Main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	debug.SetGCPercent(10)
	if err := limits.SetLimits(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to set limits: %v\n", err)
		os.Exit(1)
	}
	//
	//f, err := os.Create("testtrace.out")
	//if err != nil {
	//	panic(err)
	//}
	//err = trace.Start(f)
	//if err != nil {
	//	panic(err)
	//}
	////mf, err := os.Create("testmem.prof")
	////if err != nil {
	////	log.FATAL("could not create memory profile: ", err)
	////}
	////go func() {
	////	time.Sleep(time.Minute)
	////	runtime.GC() // get up-to-date statistics
	////	if err := pprof.WriteHeapProfile(mf); err != nil {
	////		log.FATAL("could not write memory profile: ", err)
	////	}
	////}()
	////cf, err := os.Create("testcpu.prof")
	////if err != nil {
	////	log.FATAL("could not create CPU profile: ", err)
	////}
	////if err := pprof.StartCPUProfile(cf); err != nil {
	////	log.FATAL("could not start CPU profile: ", err)
	////}
	//go func() {
	//	log.INFO(http.ListenAndServe("localhost:6060", nil))
	//}()
	//interrupt.AddHandler(
	//	func() {
	//		fmt.Println("stopping trace")
	//		trace.Stop()
	//		//pprof.StopCPUProfile()
	//		err := f.Close()
	//		if err != nil {
	//			log.ERROR(err)
	//		}
	//		//err = mf.Close()
	//		//if err != nil {
	//		//	log.ERROR(err)
	//		//}
	//	},
	//)
	app.Main()
	<-interrupt.HandlersDone
}
