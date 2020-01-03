package duoui

import (
	"errors"

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
		select {
		case <- duo.Ready:
			duo.IsReady = true
		case <-duo.Quit:
			log.DEBUG("quit signal received")
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-duo.DuoUIwindow.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.DEBUG("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				if duo.IsReady {
					duo.DuoUIcontext.Reset(e.Config, e.Size)
					DuoUIgrid(duo, cx, rc)
					e.Frame(duo.DuoUIcontext.Ops)
				}else {
					duo.DuoUIcontext.Reset(e.Config, e.Size)
					DuoUImainMenu(duo, cx, rc)
					e.Frame(duo.DuoUIcontext.Ops)
				}
			}
		}
	}
}

// Main wallet screen
func DuoUImainMenu(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	duo.DuoUIcomponents.View.Layout.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 64, helpers.HexARGB("ffcfcfcf"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			//DuoUIheader(duo,rc)
		}),
		layout.Flexed(1, func() {
			cs := duo.DuoUIcontext.Constraints
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
			//DuoUIbody(duo,cx,rc)
		}),
	)
}

// Main wallet screen
func DuoUIgrid(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// START View <<<
	cs := duo.DuoUIcontext.Constraints
	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

	layout.Flex{Axis:layout.Vertical}.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			DuoUIheader(duo,rc)
		}),
		layout.Flexed(1, func() {
			DuoUIbody(duo,cx,rc)
		}),
		layout.Rigid(func() {
			DuoUIfooter(duo,rc)
		}),
	)
}
