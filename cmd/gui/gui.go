package gui

import (
	"gioui.org/app"

	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
)

func WalletGUI(duo *model.DuoUI, rc *rcd.RcVar) (err error) {
	go func() {
		Debug("starting UI main loop")
		if rc.IsReady != false {
		}
		if err := duoui.DuoUImainLoop(duo, rc); err != nil {
			Fatal(err.Error(), "- shutting down")
		}
	}()
	Debug("starting up gio app main")
	app.Main()
	Debug("GUI shut down")
	return
}
