package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/clipboard"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	blocksList = &layout.List{
		Axis: layout.Vertical,
	}
	blockNumber = &controller.DuoUIcounter{
		Value:        11,
		OperateValue: 1,
		From:         0,
		To:           15,
	}
)

func (ui *DuoUI) DuoUIexplorer() func() {
	return func() {
		ui.rc.GetBlocksExcerpts(0, 11)
		in := layout.UniformInset(unit.Dp(0))
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
							Axis: layout.Horizontal,
						}.Layout(ui.ly.Context,
							layout.Rigid(func() {
								num := ui.ly.Theme.Body1(fmt.Sprint(i))
								num.Font.Typeface = ui.ly.Theme.Font.Primary
								num.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								num.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								var linkButton theme.DuoUIbutton
								linkButton = ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Mono, fmt.Sprint(b.Height), ui.ly.Theme.Color.Light, ui.ly.Theme.Color.Dark, "", ui.ly.Theme.Color.Light, 16, 0, 60, 24, 0, 0)
								for b.Link.Clicked(ui.ly.Context) {
									clipboard.Set(b.BlockHash)
								}
								linkButton.Layout(ui.ly.Context, b.Link)
							}),
							layout.Rigid(func() {
								amount := ui.ly.Theme.H5(fmt.Sprintf("%0.8f", b.Amount))
								amount.Font.Typeface = ui.ly.Theme.Font.Primary
								amount.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								amount.Alignment = text.End
								amount.Font.Variant = "Mono"
								amount.Font.Weight = text.Bold
								amount.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								sat := ui.ly.Theme.Body1(fmt.Sprint(b.TxNum))
								sat.Font.Typeface = ui.ly.Theme.Font.Primary
								sat.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								sat.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								sat := ui.ly.Theme.Body1(fmt.Sprint(b.BlockHash))
								sat.Font.Typeface = ui.ly.Theme.Font.Mono
								sat.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								sat.Layout(ui.ly.Context)
							}),
							layout.Rigid(func() {
								l := ui.ly.Theme.Body2(b.Time)
								l.Font.Typeface = ui.ly.Theme.Font.Primary
								l.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
								l.Layout(ui.ly.Context)
							}),
						)
					}),
					layout.Rigid(func() {
						sat := ui.ly.Theme.Body1(fmt.Sprintf("%0.8f", b.Amount))
						sat.Color = theme.HexARGB(ui.ly.Theme.Color.Hint)
						sat.Layout(ui.ly.Context)
					}),
				)
			})
		})
	}
}
