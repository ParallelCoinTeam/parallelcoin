package duoui

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
var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
)
func DuoUIexplorer(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {

	in := layout.UniformInset(unit.Dp(60))
	in.Layout(duo.DuoUIcontext, func() {

		transList.Layout(duo.DuoUIcontext, len(rc.Txs.Txs), func(i int) {
			// Invert list
			//i = len(txs.Txs) - 1 - i
			//
			t := rc.Txs.Txs[i]
			a := 1.0
			//const duration = 5
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, helpers.HexARGB("ff535353"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(duo.DuoUIcontext,
						layout.Rigid(func() {
							num := duo.DuoUItheme.Body1(fmt.Sprint(i))
							num.Color = helpers.Alpha(a, num.Color)
							num.Layout(duo.DuoUIcontext)
						}),
						layout.Rigid(func() {
							tim := duo.DuoUItheme.Body1(t.TxID)
							tim.Color = helpers.Alpha(a, tim.Color)
							tim.Layout(duo.DuoUIcontext)
						}),
						layout.Rigid(func() {
							amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", t.Amount))
							amount.Color = helpers.RGB(0x003300)
							amount.Color = helpers.Alpha(a, amount.Color)
							amount.Alignment = text.End
							amount.Font.Variant = "Mono"
							amount.Font.Weight = text.Bold
							amount.Layout(duo.DuoUIcontext)
						}),
						layout.Rigid(func() {
							sat := duo.DuoUItheme.Body1(t.Category)
							sat.Color = helpers.Alpha(a, sat.Color)
							sat.Layout(duo.DuoUIcontext)
						}),
						layout.Rigid(func() {
							l := duo.DuoUItheme.Body2(t.Time)
							l.Color = duo.DuoUItheme.Color.Hint
							l.Color = helpers.Alpha(a, l.Color)
							l.Layout(duo.DuoUIcontext)
						}),
					)
				}),
				layout.Rigid(func() {
					sat := duo.DuoUItheme.Body1(fmt.Sprintf("%0.8f", t.Amount))
					sat.Color = helpers.Alpha(a, sat.Color)
					sat.Layout(duo.DuoUIcontext)
				}),
			)
		})

	})
}
