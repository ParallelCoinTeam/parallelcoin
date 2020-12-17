package pipe

import (
	"fmt"
	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/pod/pkg/util/logi"
	qu "github.com/p9c/pod/pkg/util/quit"
	"io"
	"os"
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
			// if n, err = os.Stderr.Write(append([]byte("BACKUP:\n"), data[:n]...)); Check(err) {
			// }
			n, err = w.StdConn.Read(data)
			if n == 0 {
				Trace("read zero from stdconn", args)
				onBackup = true
				logi.L.LogChanDisabled = true
				// break out
				// close(quit)
			}
			if err != nil && err != io.EOF {
				// Probably the child process has died, so quit
				Error("err:", err)
				// if err = w.Interrupt(); Check(err) {
				// }
				onBackup = true
				// break out
			} else if n > 0 {
				if err := handler(data[:n]); Check(err) {
				}
			}
			if n, err = w.StdPipe.Read(data); Check(err) {
			}
			// when the child stops sending over RPC, fall back to the also working but not printing stderr
			if n > 0 {
				prefix := "[" + args[len(args)-1] + "]"
				if onBackup {
					prefix += "b"
				}
				printIt := true
				if logi.L.LogChanDisabled {
					printIt = false
					// prefix += "l"
				}
				// switch {
				// case onBackup:
				// 	prefix = "onBackup"
				// 	// fallthrough
				// case logi.L.LogChanDisabled:
				// 	// printIt = false
				// 	// prefix += "LogChanDisabled"
				// }
				// if onBackup || logi.L.LogChanDisabled {
				// 	prefix +=
				// printIt = false
				// }
				if printIt {
					fmt.Fprint(os.Stderr, prefix+" "+string(data[:n]))
					// Debug(prefix, onBackup, logi.L.LogChanDisabled, strings.TrimSpace(string(data[:n])))
				}
			}
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
				Debug(interrupt.GoroutineDump())
				break out
			// case l := <-logi.L.LogChan:
			// 	if l.CodeLocation == "" && l.Level == "" && l.Package == "" && l.Text == "" && l.Time.Equal(time.Time{}) {
			// 		Debug("log chan has closed")
			// 		break out
			// 	} else {
			// 		// send it back
			// 		logi.L.LogChan <- l
			// 	}
			default:
			}
			n, err = os.Stdin.Read(data)
			if err != nil && err != io.EOF {
				Debug("err: ", err)
			}
			if n > 0 {
				if err := handler(data[:n]); Check(err) {
					break out
				}
			}
		}
		Debug(interrupt.GoroutineDump())
		Debug("pipe server shut down")
	}()
	// si, _ := os.Stdin.Stat()
	// imod := si.Mode()
	// os.Stdin.Chmod(imod &^ syscall.O_NONBLOCK)
	// so, _ := os.Stdin.Stat()
	// omod := so.Mode()
	// os.Stdin.Chmod(omod &^ syscall.O_NONBLOCK)
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
