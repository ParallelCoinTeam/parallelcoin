package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI) DuoUIexplorer() (*layout.Context, func()) {
	return ui.ly.Context, func() {
		//ui.rc.GetBlocksExcerpts(ui.rc.Cx, 0, 5)
		in := layout.UniformInset(unit.Dp(60))
		in.Layout(ui.ly.Context, func() {
			blocksList.Layout(ui.ly.Context, len(ui.rc.Blocks), func(i int) {
				b := ui.rc.Blocks[i]
				cs := ui.ly.Context.Constraints
				theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, "ff535353", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				layout.Flex{
					Spacing: layout.SpaceBetween,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						layout.Flex{
							Axis: layout.Vertical,
						}.Layout(ui.ly.Context,
							layout.Rigid(func() {
								num := ui.ly.Theme.Body1(fmt.Sprint(i))
								num.Color = ui.ly.Theme.Color.Hint
								num.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								tim := ui.ly.Theme.Body1(fmt.Sprint(b.Height))
								tim.Color = ui.ly.Theme.Color.Hint
								tim.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								amount := ui.ly.Theme.H5(fmt.Sprintf("%0.8f", b.Amount))
								amount.Color = ui.ly.Theme.Color.Hint
								amount.Alignment = text.End
								amount.Font.Variant = "Mono"
								amount.Font.Weight = text.Bold
								amount.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								sat := ui.ly.Theme.Body1(fmt.Sprint(b.TxNum))
								sat.Color = ui.ly.Theme.Color.Hint
								sat.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								sat := ui.ly.Theme.Body1(fmt.Sprint(b.BlockHash))
								sat.Color = ui.ly.Theme.Color.Hint
								sat.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								l := ui.ly.Theme.Body2(b.Time)
								l.Color = ui.ly.Theme.Color.Hint
								l.Layout(ui.ly.Context)
							}),
						)
					}),
					layout.Rigid(func() {
						sat := ui.ly.Theme.Body1(fmt.Sprintf("%0.8f", b.Amount))
						sat.Color = ui.ly.Theme.Color.Hint
						sat.Layout(ui.ly.Context)
					}),
				)
			})
		})
	}
}
