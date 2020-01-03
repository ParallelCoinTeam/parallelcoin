package duoui

import (
	"github.com/p9c/pod/cmd/gui/componentsWidgets"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
)

func DuoUIhistory(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	componentsWidgets.DuoUItransactionsWidget(duo,cx,rc)
}
