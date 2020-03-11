package gui

import (
	"gioui.org/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	log "github.com/p9c/logi"
)

func WalletGUI(duo *model.DuoUI, rc *rcd.RcVar) (err error) {
	go func() {
		log.L.Debug("starting UI main loop")
		if rc.IsReady != false {
		}
		if err := duoui.DuoUImainLoop(duo, rc); err != nil {
			log.L.Fatal(err.Error(), "- shutting down")
		}
	}()
	log.L.Debug("starting up gio app main")
	app.Main()
	log.L.Debug("GUI shut down")
	return
}
