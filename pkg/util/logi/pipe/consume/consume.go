package consume

import (
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/Entry"
	"github.com/p9c/pod/pkg/util/logi/Pkg"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
	"github.com/p9c/pod/pkg/util/quit"
)

func FilterNone(string) bool {
	return false
}

func SimpleLog(name string) func(ent *logi.Entry) (err error) {
	return func(ent *logi.Entry) (err error) {
		Debugf(
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
	quit qu.C, handler func(ent *logi.Entry) (
	err error,
), filter func(pkg string) (out bool),
	args ...string,
) *worker.Worker {
	Debug("starting log consumer")
	return pipe.Consume(
		quit, func(b []byte) (err error) {
			// we are only listening for entries
			if len(b) >= 4 {
				magic := string(b[:4])
				switch magic {
				case "entr":
					// Debug(b)
					e := Entry.LoadContainer(b).Struct()
					if filter(e.Package) {
						// if the worker filter is out of sync this stops it printing
						return
					}
					switch e.Level {
					case logi.Fatal:
					case logi.Error:
					case logi.Warn:
					case logi.Info:
					case logi.Check:
					case logi.Debug:
					case logi.Trace:
					default:
						Debug("got an empty log entry")
						return
					}
					// Debugf("%s%s %s%s", color, e.Text, logi.ColorOff, e.CodeLocation)
					if err := handler(e); Check(err) {
					}
				}
			}
			return
		}, args...,
	)
}

func Start(w *worker.Worker) {
	Debug("sending start signal")
	if n, err := w.StdConn.Write([]byte("run ")); n < 1 || Check(err) {
		Debug("failed to write")
	}
}

func Stop(w *worker.Worker) {
	Debug("sending stop signal")
	if n, err := w.StdConn.Write([]byte("stop")); n < 1 || Check(err) {
		Debug("failed to write")
	}
}

func Kill(w *worker.Worker) {
	var err error
	if w == nil {
		Debug("asked to kill worker that is already nil")
		return
	}
	var n int
	Debug("sending kill signal")
	if n, err = w.StdConn.Write([]byte("kill")); n < 1 || Check(err) {
		Debug("failed to write")
		return
	}
	// close(w.Quit)
	// w.StdConn.Quit.Q()
	Debug("sent kill signal")
}

func SetLevel(w *worker.Worker, level string) {
	if w == nil {
		return
	}
	Debug("sending set level", level)
	lvl := 0
	for i := range logi.Levels {
		if level == logi.Levels[i] {
			lvl = i
		}
	}
	if n, err := w.StdConn.Write([]byte("slvl" + string(byte(lvl)))); n < 1 ||
		Check(err) {
		Debug("failed to write")
	}
}

func SetFilter(w *worker.Worker, pkgs Pk.Package) {
	if w == nil {
		return
	}
	Info("sending set filter")
	if n, err := w.StdConn.Write(Pkg.Get(pkgs).Data); n < 1 ||
		Check(err) {
		Debug("failed to write")
	}
}
