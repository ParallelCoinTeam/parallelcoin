package gui

import (
	"gioui.org/app"

	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func WalletGUI(duo *model.DuoUI, rc *rcd.RcVar) (err error) {
	go func() {
		L.Debug("starting UI main loop")
		if rc.IsReady != false {
		}
		if err := duoui.DuoUImainLoop(duo, rc); err != nil {
			L.Fatal(err.Error(), "- shutting down")
		}
	}()
	L.Debug("starting up gio app main")
	app.Main()
	L.Debug("GUI shut down")
	return
}
