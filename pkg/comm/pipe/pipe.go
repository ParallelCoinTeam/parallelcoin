package pipe

import (
	"github.com/p9c/pod/pkg/comm/stdconn"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/stalker-loki/app/slog"
	"io"
	"os"
)

func Consume(quit chan struct{}, handler func([]byte) error, args ...string) *worker.Worker {
	var n int
	var err error
	slog.Debug("spawning worker process", args)
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
				// Probably the child process has died, so quit
				slog.Error("err:", err)
				break out
			} else if n > 0 {
				if err := handler(data[:n]); slog.Check(err) {
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
		slog.Debug("starting pipe server")
	out:
		for {
			select {
			case <-quit:
				break out
			default:
			}
			n, err = os.Stdin.Read(data)
			if err != nil && err != io.EOF {
				slog.Debug("err: ", err)
			}
			if n > 0 {
				if err := handler(data[:n]); slog.Check(err) {
				}
			}
		}
	}()
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
