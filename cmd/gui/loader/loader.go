package loader

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/log"
)

var (
	logOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	logMessages []log.Entry
	logChan     = make(chan log.Entry, 111)
)

func init() {
	log.L.LogChan = logChan
	log.L.SetLevel("Info", false)
	go func() {
		for {
			select {
			case n := <-log.L.LogChan:
				logMessages = append(logMessages, n)
			}
		}
	}()
}
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
		})
	})
}
