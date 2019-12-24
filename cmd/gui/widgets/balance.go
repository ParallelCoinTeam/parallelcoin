package widgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIbalanceWidget(duo *models.DuoUI, rc *rcd.RcVar) {
	in := layout.Inset{
		Top:    unit.Dp(15),
		Right:  unit.Dp(30),
		Bottom: unit.Dp(15),
		Left:   unit.Dp(30),
	}
	in.Layout(duo.Gc, func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(duo.Gc,
			layout.Rigid(func() {
				balanceTxt := duo.Th.H6("Balance :")
				balanceTxt.Color = duo.Conf.StatusTextColor
				balanceTxt.Layout(duo.Gc)
			}),
			layout.Rigid(func() {

				balanceVal := duo.Th.H4(rc.Balance + " " + duo.Conf.Abbrevation)
				balanceVal.Color = duo.Conf.StatusTextColor
				balanceVal.Alignment = text.End
				balanceVal.Layout(duo.Gc)

			}),
			layout.Rigid(func() {
				balanceUnconfirmed := duo.Th.H6("Unconfirmed :" + rc.Unconfirmed)
				balanceUnconfirmed.Color = duo.Conf.StatusTextColor
				balanceUnconfirmed.Alignment = text.End
				balanceUnconfirmed.Layout(duo.Gc)
			}),
			layout.Rigid(func() {
				txsNumber := duo.Th.H6("Transactions :" + fmt.Sprint(rc.Transactions.TxsNumber))
				txsNumber.Color = duo.Conf.StatusTextColor
				txsNumber.Alignment = text.End
				txsNumber.Layout(duo.Gc)
			}),

		)

	})
}
