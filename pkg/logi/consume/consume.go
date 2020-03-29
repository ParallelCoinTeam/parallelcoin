package consume

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/Entry"
	"github.com/p9c/pod/pkg/pipe"
	"github.com/p9c/pod/pkg/stdconn/worker"
)

func Log(quit chan struct{}, handler func(ent *logi.Entry) (
	err error), args ...string) *worker.Worker {
	Debug("starting log consumer")
	return pipe.Consume(quit, func(b []byte) (err error) {
		// we are only listening for entries
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "entr":
				//Debug(b)
				e := Entry.LoadContainer(b).Struct()
				//Debugs(e)
				if err := handler(e); Check(
					err) {
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
