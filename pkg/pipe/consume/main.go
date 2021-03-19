package consume

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/logg"
	"github.com/p9c/pod/pkg/util/qu"
)

// FilterNone is a filter that doesn't
func FilterNone(string) bool {
	return false
}

// SimpleLog is a very simple log printer
func SimpleLog(name string) func(ent *logg.Entry) (e error) {
	return func(ent *logg.Entry) (e error) {
		dbg.F(
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
	quit qu.C, handler func(ent *logg.Entry) (e error,), filter func(pkg string) (out bool), args ...string,
) *worker.Worker {
	dbg.Ln("starting log consumer")
	return pipe.Consume(
		quit, func(b []byte) (e error) {
			// we are only listening for entries
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "entr":
					var ent logg.Entry
					n := gotiny.Unmarshal(b, &ent)
					dbg.Ln("consume", n)
					if filter(ent.Package) {
						// if the worker filter is out of sync this stops it printing
						return
					}
					switch ent.Level {
					case logg.Fatal:
					case logg.Error:
					case logg.Warn:
					case logg.Info:
					case logg.Check:
					case logg.Debug:
					case logg.Trace:
					default:
						dbg.Ln("got an empty log entry")
						return
					}
					if e = handler(&ent); err.Chk(e) {
					}
				}
			}
			return
		}, args...,
	)
}

func Start(w *worker.Worker) {
	dbg.Ln("sending start signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("run ")); n < 1 || err.Chk(e) {
		dbg.Ln("failed to write", w.Args)
	}
}

// Stop running the worker
func Stop(w *worker.Worker) {
	dbg.Ln("sending stop signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("stop")); n < 1 || err.Chk(e) {
		dbg.Ln("failed to write", w.Args)
	}
}

// Kill sends a kill signal via the pipe logger
func Kill(w *worker.Worker) {
	var e error
	if w == nil {
		dbg.Ln("asked to kill worker that is already nil")
		return
	}
	var n int
	dbg.Ln("sending kill signal")
	if n, e = w.StdConn.Write([]byte("kill")); n < 1 || err.Chk(e) {
		dbg.Ln("failed to write")
		return
	}
	// close(w.Quit)
	// w.StdConn.Quit.Q()
	if e = w.Cmd.Wait(); err.Chk(e) {
	}
	dbg.Ln("sent kill signal")
}

// SetLevel sets the level of logging from the worker
func SetLevel(w *worker.Worker, level string) {
	if w == nil {
		return
	}
	dbg.Ln("sending set level", level)
	lvl := 0
	for i := range logg.Levels {
		if level == logg.Levels[i] {
			lvl = i
		}
	}
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("slvl" + string(byte(lvl)))); n < 1 ||
		err.Chk(e) {
		dbg.Ln("failed to write")
	}
}

//
// func SetFilter(w *worker.Worker, pkgs Pk.Package) {
// 	if w == nil {
// 		return
// 	}
// 	inf.Ln("sending set filter")
// 	if n, e= w.StdConn.Write(Pkg.Get(pkgs).Data); n < 1 ||
// 		err.Chk(e) {
// 		dbg.Ln("failed to write")
// 	}
// }
