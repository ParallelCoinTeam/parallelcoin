package pipe

import (
	"github.com/p9c/pod/pkg/stdconn"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"io"
	"os"
)

func Consume(handler func([]byte) error, quit chan struct{}, args ...string) stdconn.StdConn {
	var n int
	var err error
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
				L.Error("err:", err)
			} else if n > 0 {
				if err := handler(data[:n]); L.Check(err) {
				}
			}

		}
	}()
	return w.StdConn
}

func Serve(handler func([]byte) error, quit chan struct{}) stdconn.StdConn {
	var n int
	var err error
	data := make([]byte, 8192)
	go func() {
	out:
		for {
			select {
			case <-quit:
				break out
			default:
			}
			n, err = os.Stdin.Read(data)
			if err != nil && err != io.EOF {
				L.Debug("err: ", err)
			}
			if n > 0 {
				if err := handler(data[:n]); L.Check(err) {
				}
			}
		}
	}()
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
