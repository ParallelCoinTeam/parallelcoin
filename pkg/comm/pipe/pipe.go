package pipe

import (
	"io"
	"os"
	
	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
)

func Consume(quit qu.C, handler func([]byte) error, args ...string) *worker.Worker {
	var n int
	var err error
	Debug("spawning worker process", args)
	w, _ := worker.Spawn(quit, args...)
	data := make([]byte, 8192)
	onBackup := false
	go func() {
	out:
		for {
			// Debug("readloop")
			select {
			case <-interrupt.HandlersDone:
				Debug("quitting log consumer")
				break out
			case <-quit:
				Debug("breaking on quit signal")
				break out
			default:
			}
			n, err = w.StdConn.Read(data)
			if n == 0 {
				Trace("read zero from stdconn", args)
				onBackup = true
				logi.L.LogChanDisabled.Store(true)
				break out
			}
			if err != nil && err != io.EOF {
				// Probably the child process has died, so quit
				Error("err:", err)
				onBackup = true
				break out
			} else if n > 0 {
				if err := handler(data[:n]); Check(err) {
				}
			}
			// if n, err = w.StdPipe.Read(data); Check(err) {
			// }
			// // when the child stops sending over RPC, fall back to the also working but not printing stderr
			// if n > 0 {
			// 	prefix := "[" + args[len(args)-1] + "]"
			// 	if onBackup {
			// 		prefix += "b"
			// 	}
			// 	printIt := true
			// 	if logi.L.LogChanDisabled {
			// 		printIt = false
			// 	}
			// 	if printIt {
			// 		fmt.Fprint(os.Stderr, prefix+" "+string(data[:n]))
			// 	}
			// }
		}
	}()
	return w
}

func Serve(quit qu.C, handler func([]byte) error) *stdconn.StdConn {
	var n int
	var err error
	data := make([]byte, 8192)
	go func() {
		Debug("starting pipe server")
	out:
		for {
			select {
			case <-quit:
				// Debug(interrupt.GoroutineDump())
				break out
			default:
			}
			n, err = os.Stdin.Read(data)
			if err != nil && err != io.EOF {
				Debug("err: ", err)
				break out
			}
			if n > 0 {
				if err := handler(data[:n]); Check(err) {
					break out
				}
			}
		}
		// Debug(interrupt.GoroutineDump())
		Debug("pipe server shut down")
	}()
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
