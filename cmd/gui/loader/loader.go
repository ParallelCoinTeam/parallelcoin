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
		logOutputList.Layout(duo.DuoUIcontext, len(logMessages), func(i int) {
			t := logMessages[i]
			cs := duo.DuoUIcontext.Constraints
			col := "ff3030cf"

			if t.Level == "TRC" {
				col = "ff3030cf"
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
			if t.Level == "Error" {
				col = "ffcf8030"
			}
			if t.Level == "FTL" {
				col = "ffcf3030"
			}

			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB(col), [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

			logText := duo.DuoUItheme.H6(fmt.Sprint(i) + "->" + fmt.Sprint(t.Text))
			logText.Layout(duo.DuoUIcontext)
		})
	})
}
