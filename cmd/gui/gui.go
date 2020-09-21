package gui

import (
	"gioui.org/app"
	"github.com/stalker-loki/app/slog"
	"github.com/stalker-loki/pod/cmd/gui/duoui"
	"github.com/stalker-loki/pod/cmd/gui/model"
	"github.com/stalker-loki/pod/cmd/gui/rcd"
	"os"
)

func WalletGUI(duo *model.DuoUI, rc *rcd.RcVar) (err error) {
	go func() {
		slog.Debug("starting UI main loop")
		if rc.IsReady != false {
		}
		if err := duoui.DuoUImainLoop(duo, rc); slog.Check(err) {
			slog.Fatal("shutting down")
			//close(rc.Quit)
			//time.Sleep(time.Second * 2)
			os.Exit(1)
		}
	}()
	slog.Debug("starting up gio app main")
	app.Main()
	slog.Debug("GUI shut down")
	return
}
