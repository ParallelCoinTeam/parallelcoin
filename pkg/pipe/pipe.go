package pipe

import (
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"io"
	"os"
)

func Consume(quit chan struct{}, handler func([]byte) error, args ...string) *worker.Worker {
	var n int
	var err error
	Debug("spawning worker process")
	w := worker.Spawn(args...)
	data := make([]byte, 8192)
	go func() {
	out:
		for {
			select {
			case <-quit:
				break out
			default:
			}
			n, err = w.StdConn.Read(data)
			if err != nil && err != io.EOF {
				Error("err:", err)
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
				}
			}
		}
	}()
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
