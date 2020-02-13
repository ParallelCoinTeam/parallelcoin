package gui

import (
	"github.com/p9c/pod/cmd/gui/mvc/view"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(sys *view.DuOS) (err error) {
	go func() {

		log.DEBUG("starting UI main loop")

		if sys.Duo.IsReady != false {

		}
		if err := duoui.DuoUImainLoop(sys); err != nil {
			log.FATAL(err.Error(), "- shutting down")
		}
	}()
	log.DEBUG("starting up gio app main")
	app.Main()
	log.DEBUG("GUI shut down")
	return
}
