package consume

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/log"
	"github.com/p9c/pod/pkg/pipe"
	"github.com/p9c/pod/pkg/pipe/stdconn/worker"
	"github.com/p9c/qu"
)

// FilterNone is a filter that doesn't
func FilterNone(string) bool {
	return false
}

// SimpleLog is a very simple log printer
func SimpleLog(name string) func(ent *log.Entry) (e error) {
	return func(ent *log.Entry) (e error) {
		D.F(
			"%s[%s] %s %s",
			name,
			ent.Level,
			// ent.Time.Format(time.RFC3339),
			ent.Text,
			ent.CodeLocation,
		)
		return
	}
}

func Log(
	quit qu.C, handler func(ent *log.Entry) (e error,), filter func(pkg string) (out bool), args ...string,
) *worker.Worker {
	D.Ln("starting log consumer")
	return pipe.Consume(
		quit, func(b []byte) (e error) {
			// we are only listening for entries
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "entr":
					var ent log.Entry
					n := gotiny.Unmarshal(b, &ent)
					D.Ln("consume", n)
					if filter(ent.Package) {
						// if the worker filter is out of sync this stops it printing
						return
					}
					switch ent.Level {
					case log.Fatal:
					case log.Error:
					case log.Warn:
					case log.Info:
					case log.Check:
					case log.Debug:
					case log.Trace:
					default:
						D.Ln("got an empty log entry")
						return
					}
					if e = handler(&ent); E.Chk(e) {
					}
				}
			}
			return
		}, args...,
	)
}

func Start(w *worker.Worker) {
	D.Ln("sending start signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("run ")); n < 1 || E.Chk(e) {
		D.Ln("failed to write", w.Args)
	}
}

// Stop running the worker
func Stop(w *worker.Worker) {
	D.Ln("sending stop signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("stop")); n < 1 || E.Chk(e) {
		D.Ln("failed to write", w.Args)
	}
}

// Kill sends a kill signal via the pipe logger
func Kill(w *worker.Worker) {
	var e error
	if w == nil {
		D.Ln("asked to kill worker that is already nil")
		return
	}
	var n int
	D.Ln("sending kill signal")
	if n, e = w.StdConn.Write([]byte("kill")); n < 1 || E.Chk(e) {
		D.Ln("failed to write")
		return
	}
	// close(w.Quit)
	// w.StdConn.Quit.Q()
	if e = w.Cmd.Wait(); E.Chk(e) {
	}
	D.Ln("sent kill signal")
}

// SetLevel sets the level of logging from the worker
func SetLevel(w *worker.Worker, level string) {
	if w == nil {
		return
	}
	D.Ln("sending set level", level)
	lvl := 0
	for i := range log.Levels {
		if level == log.Levels[i] {
			lvl = i
		}
	}
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("slvl" + string(byte(lvl)))); n < 1 ||
		E.Chk(e) {
		D.Ln("failed to write")
	}
}

//
// func SetFilter(w *worker.Worker, pkgs Pk.Package) {
// 	if w == nil {
// 		return
// 	}
// 	I.Ln("sending set filter")
// 	if n, e= w.StdConn.Write(Pkg.Get(pkgs).Data); n < 1 ||
// 		E.Chk(e) {
// 		D.Ln("failed to write")
// 	}
// }
