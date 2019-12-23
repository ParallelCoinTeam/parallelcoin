package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/io/system"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

func DuoUImainLoop(duo *DuoUI) error {

	for {
		e := <-duo.ww.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			duo.gc.Reset(e.Config, e.Size)
			DuoUIgrid(duo)
			e.Frame(duo.gc.Ops)
		}
	}
}

// START OMIT
func DuoUIgrid(duo *DuoUI) {
	// START View <<<
	duo.comp.View.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			cs := duo.gc.Constraints
			helpers.DuoUIdrawRectangle(duo.gc, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
			DuoUIheader(duo)
		}),
		layout.Flexed(1, func() {
			cs := duo.gc.Constraints
			helpers.DuoUIdrawRectangle(duo.gc, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 0, 0, 0, 0, unit.Dp(0))
			DuoUIbody(duo)
		}),
	)
	// END View >>>
}
