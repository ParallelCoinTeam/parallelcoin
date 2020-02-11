package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget"
	"github.com/p9c/pod/pkg/log"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	logMessages []log.Entry
	logChan     = make(chan log.Entry, 111)
	stopLogger = make(chan struct{})
	passPhrase        = ""
	confirmPassPhrase = ""
	passEditor        = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	confirmPassEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	encryption         = new(widget.CheckBox)
	seed               = new(widget.CheckBox)
	buttonCreateWallet = new(widget.Button)
	list               = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(8))
)


func DuoUIloader(duo *models.DuoUI) {
	//const buflen = 9
	layout.UniformInset(unit.Dp(10)).Layout(duo.DuoUIcontext, func() {
		//const n = 1e6
			cs := duo.DuoUIcontext.Constraints
		helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, duo.DuoUItheme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		logOutputList.Layout(duo.DuoUIcontext, len(logMessages), func(i int) {
			t := logMessages[i]
			col := "ff3080cf"
			if t.Level == "TRC" {
				col = "ff3080cf"
			}
			if t.Level == "DBG" {
				col = "ffcfcf30"
			}
			if t.Level == "INF" {
				col = "ff30cf30"
			}
			if t.Level == "WRN" {
				col = "ffcfcf30"
			}
			if t.Level == "ERROR" {
				col = "ffcf8030"
			}
			if t.Level == "FTL" {
				col = "ffcf3030"
			}


			logText := duo.DuoUItheme.Caption(fmt.Sprint(i) + "->" + fmt.Sprint(t.Text))
			logText.Font.Typeface = "bariol"
			logText.Color = helpers.HexARGB(col)
			logText.Layout(duo.DuoUIcontext)
			op.InvalidateOp{}.Add(duo.DuoUIcontext.Ops)

		})
	})
}
