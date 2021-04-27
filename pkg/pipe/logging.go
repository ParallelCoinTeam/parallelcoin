package pipe

import (
	"github.com/niubaoshu/gotiny"
	"github.com/p9c/interrupt"
	"github.com/p9c/qu"
	"go.uber.org/atomic"

	"github.com/p9c/log"
)

func LogConsume(
	quit qu.C, handler func(ent *log.Entry) (e error,),
	filter func(pkg string) (out bool), args ...string,
) *Worker {
	D.Ln("starting log consumer")
	return Consume(
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

func Start(w *Worker) {
	D.Ln("sending start signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("run ")); n < 1 || E.Chk(e) {
		D.Ln("failed to write", w.Args)
	}
}

// Stop running the worker
func Stop(w *Worker) {
	D.Ln("sending stop signal")
	var n int
	var e error
	if n, e = w.StdConn.Write([]byte("stop")); n < 1 || E.Chk(e) {
		D.Ln("failed to write", w.Args)
	}
}

// Kill sends a kill signal via the pipe logger
func Kill(w *Worker) {
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
	if e = w.Cmd.Wait(); E.Chk(e) {
	}
	D.Ln("sent kill signal")
}

// SetLevel sets the level of logging from the worker
func SetLevel(w *Worker, level string) {
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

// LogServe starts up a handler to listen to logs from the child process worker
func LogServe(quit qu.C, appName string) {
	D.Ln("starting log server")
	lc := log.AddLogChan()
	var logOn atomic.Bool
	logOn.Store(false)
	p := Serve(
		quit, func(b []byte) (e error) {
			// listen for commands to enable/disable logging
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "run ":
					D.Ln("setting to run")
					logOn.Store(true)
				case "stop":
					D.Ln("stopping")
					logOn.Store(false)
				case "slvl":
					D.Ln("setting level", log.Levels[b[4]])
					log.SetLogLevel(log.Levels[b[4]])
				case "kill":
					D.Ln("received kill signal from pipe, shutting down", appName)
					interrupt.Request()
					quit.Q()
				}
			}
			return
		},
	)
	go func() {
	out:
		for {
			select {
			case <-quit.Wait():
				if !log.LogChanDisabled.Load() {
					log.LogChanDisabled.Store(true)
				}
				D.Ln("quitting pipe logger")
				interrupt.Request()
				logOn.Store(false)
			out2:
				// drain log channel
				for {
					select {
					case <-lc:
						break
					default:
						break out2
					}
				}
				break out
			case ent := <-lc:
				if !logOn.Load() {
					break out
				}
				var n int
				var e error
				if n, e = p.Write(gotiny.Marshal(&ent)); !E.Chk(e) {
					if n < 1 {
						E.Ln("short write")
					}
				} else {
					break out
				}
			}
		}
		<-interrupt.HandlersDone
		D.Ln("finished pipe logger")
	}()
}

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
