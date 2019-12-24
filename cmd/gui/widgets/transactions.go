package widgets

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIlatestTxsWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {

	rc.GetDuoUIlastTxs(duo, cx)
	in := layout.UniformInset(unit.Dp(30))
	in.Layout(duo.Gc, func() {
		duo.Comp.Status.Layout.Layout(duo.Gc,
			// Balance status item
			layout.Rigid(func() {

				//const n = 5
				//list.Layout(duo.Gc, n, func(i int) {
				//	txt := fmt.Sprintf("List element #%d", i)
				//
				//	duo.Th.H3(txt).Layout(duo.Gc)
				//})
				transList := &layout.List{
					Axis: layout.Vertical,
				}
				transList.Layout(duo.Gc, len(rc.Transactions.Txs), func(i int) {
					// Invert list
					//i = len(txs.Txs) - 1 - i
					t := rc.Transactions.Txs[i]
					a := 1.0
					//const duration = 5

					layout.Flex{
						Spacing:layout.SpaceBetween,
					}.Layout(duo.Gc,
						layout.Rigid(func() {
							tim := duo.Th.Body1(t.TxID)
							tim.Color = helpers.Alpha(a, tim.Color)
							tim.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							amount := duo.Th.H5(fmt.Sprintf("%0.8f", t.Amount))
							amount.Color = helpers.RGB(0x003300)
							amount.Color = helpers.Alpha(a, amount.Color)
							amount.Alignment = text.End
							amount.Font.Variant = "Mono"
							amount.Font.Weight = text.Bold
							amount.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							sat := duo.Th.Body1(t.Category)
							sat.Color = helpers.Alpha(a, sat.Color)
							sat.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							sat := duo.Th.Body1(fmt.Sprintf("%0.8f", t.Amount))
							sat.Color = helpers.Alpha(a, sat.Color)
							sat.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							l := duo.Th.Body2(helpers.FormatTime(t.Time))
							l.Color = duo.Th.Color.Hint
							l.Color = helpers.Alpha(a, l.Color)
							l.Layout(duo.Gc)
						}),
					)
				})

			}))
	})
}


func DuoUItransactionsWidget(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {

	rc.GetDuoUITransactionsExcertps(duo, cx)
	in := layout.UniformInset(unit.Dp(30))
	in.Layout(duo.Gc, func() {
		duo.Comp.Status.Layout.Layout(duo.Gc,
			// Balance status item
			layout.Rigid(func() {

				//const n = 5
				//list.Layout(duo.Gc, n, func(i int) {
				//	txt := fmt.Sprintf("List element #%d", i)
				//
				//	duo.Th.H3(txt).Layout(duo.Gc)
				//})
				transList := &layout.List{
					Axis: layout.Vertical,
				}
				transList.Layout(duo.Gc, len(rc.Transactions.Txs), func(i int) {
					// Invert list
					//i = len(txs.Txs) - 1 - i
					t := rc.Transactions.Txs[i]
					a := 1.0
					//const duration = 5

					layout.Flex{
						Spacing:layout.SpaceBetween,
					}.Layout(duo.Gc,
						layout.Rigid(func() {
							tim := duo.Th.Body1(t.TxID)
							tim.Color = helpers.Alpha(a, tim.Color)
							tim.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							amount := duo.Th.H5(fmt.Sprintf("%0.8f", t.Amount))
							amount.Color = helpers.RGB(0x003300)
							amount.Color = helpers.Alpha(a, amount.Color)
							amount.Alignment = text.End
							amount.Font.Variant = "Mono"
							amount.Font.Weight = text.Bold
							amount.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							sat := duo.Th.Body1(t.Category)
							sat.Color = helpers.Alpha(a, sat.Color)
							sat.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							sat := duo.Th.Body1(fmt.Sprintf("%0.8f", t.Amount))
							sat.Color = helpers.Alpha(a, sat.Color)
							sat.Layout(duo.Gc)
						}),
						layout.Rigid(func() {
							l := duo.Th.Body2(helpers.FormatTime(t.Time))
							l.Color = duo.Th.Color.Hint
							l.Color = helpers.Alpha(a, l.Color)
							l.Layout(duo.Gc)
						}),
					)
				})

			}))
	})
}
