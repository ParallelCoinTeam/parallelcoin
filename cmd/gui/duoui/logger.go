package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/unit"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

func (ui *DuoUI) DuoUIlogger() func() {
	return func() {
		//const buflen = 9
		layout.UniformInset(unit.Dp(10)).Layout(ui.ly.Context, func() {
			//const n = 1e6
			cs := ui.ly.Context.Constraints
			theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			logOutputList.Layout(ui.ly.Context, len(ui.rc.Log.LogMessages), func(i int) {
				t := ui.rc.Log.LogMessages[i]
				col := ui.ly.Theme.Color.Primary
				if t.Level == "TRC" {
					col = ui.ly.Theme.Color.Success
				}
				if t.Level == "DBG" {
					col = ui.ly.Theme.Color.Secondary
				}
				if t.Level == "INF" {
					col = ui.ly.Theme.Color.Info
				}
				if t.Level == "WRN" {
					col = ui.ly.Theme.Color.Warning
				}
				if t.Level == "ERROR" {
					col = ui.ly.Theme.Color.Danger
				}
				if t.Level == "FTL" {
					col = ui.ly.Theme.Color.Primary
				}

				logText := ui.ly.Theme.Caption(fmt.Sprint(i) + "->" + fmt.Sprint(t.Text))
				logText.Font.Typeface = "bariol"
				logText.Color = theme.HexARGB(col)
				logText.Layout(ui.ly.Context)
				op.InvalidateOp{}.Add(ui.ly.Context.Ops)

			})
		})
	}
}

