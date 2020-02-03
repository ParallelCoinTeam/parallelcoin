package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
)

func DuoUIexplorer(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	return func() {
		rc.GetBlocksExcerpts(cx, 0, 5)

		in := layout.UniformInset(unit.Dp(60))
		in.Layout(duo.DuoUIcontext, func() {

			blocksList.Layout(duo.DuoUIcontext, len(rc.Blocks), func(i int) {
				// Invert list
				//i = len(txs.Txs) - 1 - i
				//
				b := rc.Blocks[i]
				a := 1.0
				//const duration = 5
				cs := duo.DuoUIcontext.Constraints
				helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, "ff535353", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

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
								tim := duo.DuoUItheme.Body1(fmt.Sprint(b.Height))
								tim.Color = helpers.Alpha(a, tim.Color)
								tim.Layout(duo.DuoUIcontext)
							}),
							layout.Rigid(func() {
								amount := duo.DuoUItheme.H5(fmt.Sprintf("%0.8f", b.Amount))
								amount.Color = helpers.RGB(0x003300)
								amount.Color = helpers.Alpha(a, amount.Color)
								amount.Alignment = text.End
								amount.Font.Variant = "Mono"
								amount.Font.Weight = text.Bold
								amount.Layout(duo.DuoUIcontext)
							}),
							layout.Rigid(func() {
								sat := duo.DuoUItheme.Body1(fmt.Sprint(b.TxNum))
								sat.Color = helpers.Alpha(a, sat.Color)
								sat.Layout(duo.DuoUIcontext)
							}),
							layout.Rigid(func() {
								sat := duo.DuoUItheme.Body1(fmt.Sprint(b.BlockHash))
								sat.Color = helpers.Alpha(a, sat.Color)
								sat.Layout(duo.DuoUIcontext)
							}),
							layout.Rigid(func() {
								l := duo.DuoUItheme.Body2(b.Time)
								l.Color = duo.DuoUItheme.Color.Hint
								l.Color = helpers.Alpha(a, l.Color)
								l.Layout(duo.DuoUIcontext)
							}),
						)
					}),
					layout.Rigid(func() {
						sat := duo.DuoUItheme.Body1(fmt.Sprintf("%0.8f", b.Amount))
						sat.Color = helpers.Alpha(a, sat.Color)
						sat.Layout(duo.DuoUIcontext)
					}),
				)
			})

		})
	}
}
