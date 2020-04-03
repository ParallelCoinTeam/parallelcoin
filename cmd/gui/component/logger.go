package component

import (
	"fmt"
	"time"

	"github.com/p9c/pod/cmd/gui/rcd"
	log "github.com/p9c/pod/pkg/logi"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/gelook"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
)

var StartupTime = time.Now()

func DuoUIlogger(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		// const buflen = 9
		layout.UniformInset(unit.Dp(10)).Layout(gtx, func() {
			// const n = 1e6
			cs := gtx.Constraints
			gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Colors["Dark"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			lm := rc.Log.LogMessages.Load().([]log.Entry)
			logOutputList.Layout(gtx, len(lm), func(i int) {
				t := lm[i]
				logText := th.Caption(fmt.Sprintf("%-12s", t.Time.Sub(StartupTime)/time.Second*time.Second) + " " + fmt.Sprint(t.Text))
				logText.Font.Typeface = th.Fonts["Mono"]

				logText.Color = th.Colors["Primary"]
				if t.Level == "TRC" {
					logText.Color = th.Colors["Success"]
				}
				if t.Level == "DBG" {
					logText.Color = th.Colors["Secondary"]
				}
				if t.Level == "INF" {
					logText.Color = th.Colors["Info"]
				}
				if t.Level == "WRN" {
					logText.Color = th.Colors["Warning"]
				}
				if t.Level == "ERROR" {
					logText.Color = th.Colors["Danger"]
				}
				if t.Level == "FTL" {
					logText.Color = th.Colors["Primary"]
				}

				logText.Layout(gtx)

			})
			op.InvalidateOp{}.Add(gtx.Ops)
		})
	}
}
