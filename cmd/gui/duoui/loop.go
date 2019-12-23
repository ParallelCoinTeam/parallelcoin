package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/io/system"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
	"image/color"
)

func DuoUImainLoop(duo *DuoUI) error {
	// wait for ready signal, before this comment the splash screen code should be triggered with a
	// system.FrameEvent select
	log.DEBUG("waiting for back end to become ready")
ready:
	for {
		select {
		case <-duo.Quit:
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case <-duo.Ready:
			log.DEBUG("starting main loop")
			break ready
		case e := <-duo.ww.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				duo.gc.Reset(e.Config, e.Size)
			}
		}
	}
	for {
		select {
		case <-duo.Quit:
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-duo.ww.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				duo.gc.Reset(e.Config, e.Size)
				DuoUIgrid(duo)
				e.Frame(duo.gc.Ops)
			}
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
