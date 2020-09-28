// Package monitor is a log viewer and filter and configuration interface
//
// +build !headless

package gui

import (
	"gioui.org/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/stalker-loki/app/slog"
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
