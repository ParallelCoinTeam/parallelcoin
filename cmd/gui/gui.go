package gui

import (
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(duo *model.DuoUI, rc *rcd.RcVar) (err error) {
	go func() {
		log.DEBUG("starting UI main loop")
		if duo.IsReady != false {
		}
		if err := duoui.DuoUImainLoop(duo,rc); err != nil {
			log.FATAL(err.Error(), "- shutting down")
		}
	}()
	log.DEBUG("starting up gio app main")
	app.Main()
	log.DEBUG("GUI shut down")
	return
}
