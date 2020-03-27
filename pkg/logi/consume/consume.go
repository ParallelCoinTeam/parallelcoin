package consume

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/Entry"
	"github.com/p9c/pod/pkg/pipe"
	"github.com/p9c/pod/pkg/stdconn/worker"
)

func Log(quit chan struct{}, handler func(ent *logi.Entry) (
	err error), args ...string) *worker.Worker {
	L.Debug("starting log consumer")
	return pipe.Consume(quit, func(b []byte) (err error) {
		// we are only listening for entries
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "entr":
				e := Entry.LoadContainer(b)
				//L.Debugs(e)
				if err := handler((&e).Struct()); L.Check(
					err) {
				}
			}
		}
		return
	}, args...)
}

func Start(w *worker.Worker) {
	L.Debug("sending start signal")
	if n, err := w.StdConn.Write([]byte("run ")); n < 1 || L.Check(err) {
		L.Debug("failed to write")
	}
}

func Stop(w *worker.Worker) {
	L.Debug("sending stop signal")
	if n, err := w.StdConn.Write([]byte("stop")); n < 1 || L.Check(err) {
		L.Debug("failed to write")
	}
}
