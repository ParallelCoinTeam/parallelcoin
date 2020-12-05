package pipe

import (
	"io"
	"os"
	"syscall"

	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
)

func Consume(quit chan struct{}, handler func([]byte) error, args ...string) *worker.Worker {
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
				close(quit)
			}
			if err != nil && err != io.EOF {
				// Probably the child process has died, so quit
				Error("err:", err)
				break out
			} else if n > 0 {
				if err := handler(data[:n]); Check(err) {
				}
			}

		}
	}()
	return w
}

func Serve(quit chan struct{}, handler func([]byte) error) stdconn.StdConn {
	var n int
	var err error
	data := make([]byte, 8192)
	go func() {
		Debug("starting pipe server")
	out:
		for {
			select {
			case <-quit:
				break out
			default:
			}
			n, err = os.Stdin.Read(data)
			if err != nil && err != io.EOF {
				Debug("err: ", err)
			}
			if n > 0 {
				if err := handler(data[:n]); Check(err) {
					// break out
				}
			}
		}
	}()
	si, _ := os.Stdin.Stat()
	imod := si.Mode()
	os.Stdin.Chmod(imod &^ syscall.O_NONBLOCK)
	so, _ := os.Stdin.Stat()
	omod := so.Mode()
	os.Stdin.Chmod(omod &^ syscall.O_NONBLOCK)
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
