package pipe

import (
	"github.com/p9c/log"
	"io"
	"os"
	
	"github.com/p9c/pod/pkg/pipe/stdconn"
	"github.com/p9c/pod/pkg/pipe/stdconn/worker"
	"github.com/p9c/pod/pkg/util/interrupt"
	"github.com/p9c/qu"
)

// Consume listens for messages from a child process over a stdio pipe.
func Consume(quit qu.C, handler func([]byte) error, args ...string) *worker.Worker {
	var n int
	var e error
	D.Ln("spawning worker process", args)
	w, _ := worker.Spawn(quit, args...)
	data := make([]byte, 8192)
	// onBackup := false
	go func() {
	out:
		for {
			// D.Ln("readloop")
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
				F.Ln("read zero from stdconn", args)
				// onBackup = true
				log.LogChanDisabled.Store(true)
				break out
			}
			if E.Chk(e) && e != io.EOF {
				// Probably the child process has died, so quit
				E.Ln("err:", e)
				// onBackup = true
				break out
			} else if n > 0 {
				if e = handler(data[:n]); E.Chk(e) {
				}
			}
			// if n, e = w.StdPipe.Read(data); E.Chk(e) {
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

// Serve runs a goroutine processing the FEC encoded packets, gathering them and
// decoding them to be delivered to a handler function
func Serve(quit qu.C, handler func([]byte) error) *stdconn.StdConn {
	var n int
	var e error
	data := make([]byte, 8192)
	go func() {
		D.Ln("starting pipe server")
	out:
		for {
			select {
			case <-quit.Wait():
				// D.Ln(interrupt.GoroutineDump())
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
		// D.Ln(interrupt.GoroutineDump())
		D.Ln("pipe server shut down")
	}()
	return stdconn.New(os.Stdin, os.Stdout, quit)
}
