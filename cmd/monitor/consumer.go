package monitor

import (
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/consume"
	"github.com/p9c/pod/pkg/stdconn"
	"os"
)

func (s *State) Consume() {
	consume.NewAPI(stdconn.New(os.Stdin, os.Stdout, s.Ctx.KillAll), func(ent logi.Entry) error {
		L.Debugs(ent)
		return nil
	})
out:
	for {
		select {
		case <-s.Ctx.KillAll:
			break out
		}
	}
}
