package worker

import (
	"encoding/binary"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/interrupt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
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
// Main the main thread of the kopach miner
func Main(cx *conte.Xt, quit chan struct{}) {
	printlnE("[worker] starting up")
	interrupt.AddHandler(func(){
		close(quit)
	})
	// we only want one thread
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(0)
	if os.Getenv("POD_TRACE") == "on" {
		if f, err := os.Create("testtrace.out"); err != nil {
			printlnE("[worker] tracing env POD_TRACE=on but we can't write to it",
				err)
		} else {
			log.DEBUG("[worker] tracing started")
			err = trace.Start(f)
			if err != nil {
				printlnE("[worker] could not start tracing", err)
			} else {
				interrupt.AddHandler(
					func() {
						printlnE("[worker] stopping trace")
						trace.Stop()
						err := f.Close()
						if err != nil {
							log.ERROR(err)
						}
					},
				)
			}
		}
	}
	// listen to stdin
	go func() {
		// allocate 2mb for max block size
		b := make([]byte, 1<<21)
	out:
		for {
			printlnE("[worker] OK")
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
			printlnE("[worker] message decoded\n",spew.Sdump(blk))
		}
	}()
out:
	for {
		select {
		case <-quit:
			printE("[worker] shutting down\n\r")
			break out
		}
	}
	<-interrupt.HandlersDone
}
