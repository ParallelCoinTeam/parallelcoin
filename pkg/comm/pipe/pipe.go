package pipe

import (
	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/interrupt"
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
	go func() {
	out:
		for {
			select {
			case <-quit:
				Debug("quitting log consumer")
				break out
			default:
			}
			n, err = w.StdConn.Read(data)
			// Trace("read from stdconn", n, args)
			if n == 0 {
				break out
				// close(quit)
			}
			if err != nil && err != io.EOF {
				// Probably the child process has died, so quit
				Error("err:", err)
				// if err = w.Interrupt(); Check(err) {
				// }
				break out
			} else if n > 0 {
				if err := handler(data[:n]); Check(err) {
				}
			}

		}
	}()
	return w
}

func Serve(quit qu.C, handler func([]byte) error) stdconn.StdConn {
	var n int
	var err error
	qChan := qu.T()
	data := make([]byte, 8192)
	go func() {
		Debug("starting pipe server")
	out:
		for {
			select {
			case <-qChan:
				break out
			case <-quit:
				qChan.Quit()
				break out
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
	return stdconn.New(os.Stdin, os.Stdout, qChan)
}
