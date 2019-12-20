package loader

import (
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

func DuoUIloader(duo *duoui.DuoUI, cx *conte.Xt, firstRun bool) (err error) {
	go func() {
		if err := DuoUIloaderLoop(firstRun, cx); err != nil {
			log.FATAL(err)
		}
	}()
	app.Main()
	return
}
