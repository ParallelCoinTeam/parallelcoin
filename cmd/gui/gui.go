package gui

import (
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/app"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (err error) {
	go func() {
		log.DEBUG("starting UI main loop")

		if duo.IsReady != false {

			rc.GetDuoUIbalance(duo, cx)
			rc.GetDuoUIunconfirmedBalance(duo, cx)
			rc.GetDuoUIblockHeight(duo, cx)
			rc.GetDuoUIstatus(duo, cx)
			rc.GetDuoUIlocalLost(duo)
			rc.GetDuoUIdifficulty(duo, cx)

			rc.GetDuoUIlastTxs(duo, cx)

		}
		if err := duoui.DuoUImainLoop(duo, cx, rc); err != nil {
			log.FATAL(err.Error(), "- shutting down")
		}
	}()
	log.DEBUG("starting up gio app main")
	app.Main()
	log.DEBUG("GUI shut down")
	return
}
