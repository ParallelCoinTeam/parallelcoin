package worker

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/ipc"
	"github.com/p9c/pod/pkg/util/interrupt"
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
)

// Main the main thread of the kopach miner worker module
func Main(cx *conte.Xt, quit chan struct{}) {
	printlnE("starting up")
	// we only want one thread
	runtime.GOMAXPROCS(1)
	debug.SetGCPercent(10)
	if os.Getenv("POD_TRACE") == "on" {
		if f, err := os.Create("testtrace.out"); err != nil {
			printlnE("tracing env POD_TRACE=on but we can't write to it",
				err)
		} else {
			printlnE("tracing started")
			err = trace.Start(f)
			if err != nil {
				printlnE("could not start tracing", err)
			} else {
				interrupt.AddHandler(func() {
					printlnE("stopping trace")
					trace.Stop()
					err := f.Close()
					if err != nil {
						printlnE(err)
					}
				})
			}
		}
	}
	w, err := ipc.NewWorker()
	if err != nil {
		printlnE(err)
		close(quit)
		return
	}
	interrupt.AddHandler(func() {
		// interrupt will receive a sigquit signal (
		// ctrl-c) when the close method is invoked on the controller IPC
		// via the Close() method
		//close(quit)
	})
	// listen to stdin
	go func() {
		// allocate 2mb for max block size
		b := make([]byte, 1<<21)
	out:
		for {
			printlnE("OK")
			w.Write([]byte("testing"))
			n, err := w.Read(b)
			if err != nil {
				printlnE(err)
				continue
			}
			if string(b[:n]) == string(ipc.QuitCommand) {
				printlnE("received quit message")
				close(quit)
			}
			payload := b[:n]
			printlnE(spew.Sdump(payload))
			select {
			case <-quit:
				break out
			default:
			}
		}
		printlnE("worker message handler finished")
	}()
out:
	for {
		select {
		case <-quit:
			printlnE("shutting down")
			break out
		}
	}
	// don't jump the gun!
	//<-interrupt.HandlersDone
	printlnE("worker is finished")
}

func printE(a ...interface{}) {
	out := append([]interface{}{"[worker]"}, a...)
	_, _ = fmt.Fprint(os.Stderr, out...)
}

func printlnE(a ...interface{}) {
	out := append([]interface{}{"[worker]"}, a...)
	_, _ = fmt.Fprintln(os.Stderr, out...)
	//_, _ = fmt.Fprint(os.Stderr, "\r")
}

func printfE(format string, a ...interface{}) {
	out := append([]interface{}{"[worker]"}, a...)
	_, _ = fmt.Fprintf(os.Stderr, format, out...)
}

func printErr(err error, fn func()) {
	if err != nil {
		printlnE(err)
		if fn != nil {
			fn()
		}
	}
}
