package consume

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/Entry"
	"github.com/p9c/pod/pkg/pipe"
	"github.com/p9c/pod/pkg/stdconn"
)

func Log(quit chan struct{}, handler func(ent *logi.Entry) (
	err error)) stdconn.StdConn {
	return pipe.Serve(quit, func(b []byte) (err error) {
		// we are only listening for entries
		if len(b) >= 4 {
			magic := string(b[:4])
			switch magic {
			case "entr":
				if err := handler(Entry.LoadContainer(b).Struct()); L.Check(
					err) {
				}
			}
		}
		return
	})
}

func Start(conn stdconn.StdConn) {
	conn.Write([]byte("run "))
}

func Stop(conn stdconn.StdConn) {
	conn.Write([]byte("stop"))
}