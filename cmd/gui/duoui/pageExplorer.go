package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
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

func (duo *DuoUI) DuoUIexplorer(cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.Model.DuoUIcontext, func() {
		rc.GetBlocksExcerpts(cx, 0, 5)

		in := layout.UniformInset(unit.Dp(60))
		in.Layout(duo.Model.DuoUIcontext, func() {

			blocksList.Layout(duo.Model.DuoUIcontext, len(rc.Blocks), func(i int) {
				// Invert list
				//i = len(txs.Txs) - 1 - i
				//
				b := rc.Blocks[i]
				a := 1.0
				//const duration = 5
				cs := duo.Model.DuoUIcontext.Constraints
				helpers.DuoUIdrawRectangle(duo.Model.DuoUIcontext, cs.Width.Max, 1, helpers.HexARGB("ff535353"), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(duo.Model.DuoUIcontext,
					layout.Rigid(func() {
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(duo.Model.DuoUIcontext,
							layout.Rigid(func() {
								num := duo.Model.DuoUItheme.Body1(fmt.Sprint(i))
								num.Color = helpers.Alpha(a, num.Color)
								num.Layout(duo.Model.DuoUIcontext)
							}),
							layout.Rigid(func() {
								tim := duo.Model.DuoUItheme.Body1(fmt.Sprint(b.Height))
								tim.Color = helpers.Alpha(a, tim.Color)
								tim.Layout(duo.Model.DuoUIcontext)
							}),
							layout.Rigid(func() {
								amount := duo.Model.DuoUItheme.H5(fmt.Sprintf("%0.8f", b.Amount))
								amount.Color = helpers.RGB(0x003300)
								amount.Color = helpers.Alpha(a, amount.Color)
								amount.Alignment = text.End
								amount.Font.Variant = "Mono"
								amount.Font.Weight = text.Bold
								amount.Layout(duo.Model.DuoUIcontext)
							}),
							layout.Rigid(func() {
								sat := duo.Model.DuoUItheme.Body1(fmt.Sprint(b.TxNum))
								sat.Color = helpers.Alpha(a, sat.Color)
								sat.Layout(duo.Model.DuoUIcontext)
							}),
							layout.Rigid(func() {
								sat := duo.Model.DuoUItheme.Body1(fmt.Sprint(b.BlockHash))
								sat.Color = helpers.Alpha(a, sat.Color)
								sat.Layout(duo.Model.DuoUIcontext)
							}),
							layout.Rigid(func() {
								l := duo.Model.DuoUItheme.Body2(b.Time)
								l.Color = duo.Model.DuoUItheme.Color.Hint
								l.Color = helpers.Alpha(a, l.Color)
								l.Layout(duo.Model.DuoUIcontext)
							}),
						)
					}),
					layout.Rigid(func() {
						sat := duo.Model.DuoUItheme.Body1(fmt.Sprintf("%0.8f", b.Amount))
						sat.Color = helpers.Alpha(a, sat.Color)
						sat.Layout(duo.Model.DuoUIcontext)
					}),
				)
			})

		})
	}
}
