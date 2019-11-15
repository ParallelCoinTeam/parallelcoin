package main

import (
	"encoding/binary"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/sem"
	"github.com/p9c/pod/pkg/util"
	"os"
	"runtime"
	"runtime/debug"
)

var quitMessage = string([]byte{255, 255, 255, 255})

func printE(a ...interface{}) {
	_, _ = fmt.Fprint(os.Stderr, a...)
}

func printlnE(a ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, a...)
	_, _ = fmt.Fprint(os.Stderr, "\r")
}

func printfE(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}

func printErr(err error, fn func()) {
	if err != nil {
		printlnE(err)
		if fn != nil {
			fn()
		}
	}
}

func main() {
	//printlnE("starting up")
	quit := make(sem.T)
	// we only want one thread
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(0)
	//if err := limits.SetLimits(); err != nil {
	//	printfE("failed to set limits: %v\n", err)
	//	os.Exit(1)
	//}
	////
	//f, err := os.Create("kopachtrace.out")
	//if err != nil {
	//	panic(err)
	//}
	//err = trace.Start(f)
	//printlnE("tracing started")
	//if err != nil {
	//	panic(err)
	//}

	//oldState, err := terminal.MakeRaw(0)
	//printErr(err, func() { os.Exit(1) })
	//defer func() {
	//	err := terminal.Restore(0, oldState)
	//	if err != nil {
	//		printlnE(err)
	//	}
	//}()
	// listen to stdin
	go func() {
		// allocate 2mb for max block size
		b := make([]byte, 1<<21)
	out:
		for {
			printlnE("OK")
			_, err := os.Stdin.Read(b[:4])
			if err != nil {
				printlnE(err)
			}
			out := string(b[:4])
			if out == quitMessage {
				break out
			}
			bLen := binary.BigEndian.Uint32(b[:4])
			_, err = os.Stdin.Read(b[:bLen])
			if err != nil {
				printlnE(err)
			}
			blk, err := util.NewBlockFromBytes(b[:bLen])
			if err != nil {
				printlnE(err)
			}
			printlnE(spew.Sdump(blk))
		}
		close(quit)
	}()
out:
	for {
		select {
		case <-quit:
			printE("\n\rshutting down\n\r")
			break out
		}
	}
}
