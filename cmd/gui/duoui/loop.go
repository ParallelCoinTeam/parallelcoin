package duoui

import (
	"github.com/p9c/gio-parallel/io/system"
)

func DuoUIloop(duo *DuoUI) error {

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
	duo.comp.view.l.Layout(duo.gc, DuoUIheader(duo), DuoUIbody(duo))
	// END View >>>
}

// END OMIT
