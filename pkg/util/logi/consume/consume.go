package consume

import (
	"github.com/p9c/pod/pkg/comm/pipe"
	"github.com/p9c/pod/pkg/comm/stdconn/worker"
	"github.com/p9c/pod/pkg/util/logi"
	"github.com/p9c/pod/pkg/util/logi/Entry"
	"github.com/p9c/pod/pkg/util/logi/Pkg"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
)

func Log(quit chan struct{}, handler func(ent *logi.Entry) (
	err error), filter func(pkg string) (out bool),
	args ...string) *worker.Worker {
	Debug("starting log consumer")
	return pipe.Consume(quit, func(b []byte) (err error) {
		// we are only listening for entries
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "entr":
				//Debug(b)
				e := Entry.LoadContainer(b).Struct()
				if filter(e.Package) {
					// if the worker filter is out of sync this stops it
					// printing
					return
				}
				// Debugs(e)
				color := logi.ColorYellow
				//Debug(e.Level)
				switch e.Level {
				case logi.Fatal:
					color = logi.ColorRed
				case logi.Error:
					color = logi.ColorOrange
				case logi.Warn:
					color = logi.ColorYellow
				case logi.Info:
					color = logi.ColorGreen
				case logi.Check:
					color = logi.ColorCyan
				case logi.Debug:
					color = logi.ColorBlue
				case logi.Trace:
					color = logi.ColorViolet
				default:
					Debug("got an empty log entry")
					return
				}
				Debugf("%s%s %s%s", color, e.Text,
					logi.ColorOff, e.CodeLocation)
				if err := handler(e); Check(err) {
				}
			}
		}
		return
	}, args...)
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
