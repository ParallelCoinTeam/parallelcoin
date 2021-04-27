package pipe

import (
	"io"
	"os"

	"github.com/p9c/log"

	"github.com/p9c/interrupt"
	"github.com/p9c/qu"
)

// Consume listens for messages from a child process over a stdio pipe.
func Consume(
	quit qu.C,
	handler func([]byte) error,
	args ...string,
) *Worker {
	var n int
	var e error
	D.Ln("spawning worker process", args)
	var w *Worker
	if w, e = Spawn(quit, args...); E.Chk(e) {
	}
	data := make([]byte, 8192)
	go func() {
	out:
		for {
			select {
			case <-interrupt.HandlersDone.Wait():
				D.Ln("quitting log consumer")
				break out
			case <-quit.Wait():
				D.Ln("breaking on quit signal")
				break out
			default:
			}
			n, e = w.StdConn.Read(data)
			if n == 0 {
				F.Ln("read zero from pipe", args)
				log.LogChanDisabled.Store(true)
				break out
			}
			if E.Chk(e) && e != io.EOF {
				// Probably the child process has died, so quit
				E.Ln("err:", e)
				break out
			} else if n > 0 {
				if e = handler(data[:n]); E.Chk(e) {
				}
			}
		}
	}()
	return w
}

// Serve runs a goroutine processing the FEC encoded packets, gathering them and
// decoding them to be delivered to a handler function
func Serve(quit qu.C, handler func([]byte) error) *StdConn {
	var n int
	var e error
	data := make([]byte, 8192)
	go func() {
		D.Ln("starting pipe server")
	out:
		for {
			select {
			case <-quit.Wait():
				break out
			default:
			}
			n, e = os.Stdin.Read(data)
			if e != nil && e != io.EOF {
				break out
			}
			if n > 0 {
				if e = handler(data[:n]); E.Chk(e) {
					break out
				}
			}
		}
		D.Ln("pipe server shut down")
	}()
	return New(os.Stdin, os.Stdout, quit)
}
