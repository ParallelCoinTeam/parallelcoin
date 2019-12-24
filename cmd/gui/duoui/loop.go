package duoui

import (
	"errors"
	"image/color"
	
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/io/system"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func DuoUImainLoop(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) error {

	for {
		e := <-duo.Ww.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			duo.Gc.Reset(e.Config, e.Size)
			DuoUIgrid(duo,cx,rc)
			e.Frame(duo.Gc.Ops)
		}
	}
}

// START OMIT
func DuoUIgrid(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	duo.Comp.View.Layout.Layout(duo.Gc,
		layout.Rigid(func() {
			cs := duo.Gc.Constraints
			helpers.DuoUIdrawRectangle(duo.Gc, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}, 0, 0, 0, 0, unit.Dp(0))
			DuoUIheader(duo,rc)
		}),
		layout.Flexed(1, func() {
			cs := duo.Gc.Constraints
			helpers.DuoUIdrawRectangle(duo.Gc, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 0, 0, 0, 0, unit.Dp(0))
			DuoUIbody(duo,cx,rc)
		}),
	)
	// END View >>>
}
