package gui

import (
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/app"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (err error) {
	go func() {
		log.DEBUG("starting UI main loop")
		if err := duoui.DuoUImainLoop(duo, cx, rc); err != nil {
			log.FATAL(err.Error(), "- shutting down")
		}
	}()
	log.DEBUG("starting up gio app main")
	app.Main()
	log.DEBUG("GUI shut down")
	return
}
