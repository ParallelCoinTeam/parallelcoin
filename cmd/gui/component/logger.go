package component

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/rcd"
	"time"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/theme"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

var StartupTime = time.Now()

func DuoUIlogger(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		// const buflen = 9
		layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
			// const n = 1e6
			cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			logOutputList.Layout(gtx, len(rc.Log.LogMessages), func(i int) {
				t := rc.Log.LogMessages[i]
				logText := th.Caption(fmt.Sprintf("%-12s", t.Time.Sub(StartupTime)/time.Second*time.Second) + " " + fmt.Sprint(t.Text))
				logText.Font.Typeface = th.Font.Mono

				logText.Color = theme.HexARGB(th.Color.Primary)
				if t.Level == "TRC" {
					logText.Color = theme.HexARGB(th.Color.Success)
				}
				if t.Level == "DBG" {
					logText.Color = theme.HexARGB(th.Color.Secondary)
				}
				if t.Level == "INF" {
					logText.Color = theme.HexARGB(th.Color.Info)
				}
				if t.Level == "WRN" {
					logText.Color = theme.HexARGB(th.Color.Warning)
				}
				if t.Level == "ERROR" {
					logText.Color = theme.HexARGB(th.Color.Danger)
				}
				if t.Level == "FTL" {
					logText.Color = theme.HexARGB(th.Color.Primary)
				}

				logText.Layout(gtx)
				op.InvalidateOp{}.Add(gtx.Ops)

			})
		})
	}
}
