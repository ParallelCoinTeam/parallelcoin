package duoui

import (
	"github.com/p9c/gio-parallel/layout"
	"image/color"
)

func DuoUIbody(duo *DuoUI) layout.FlexChild {
	return duo.comp.view.l.Flex(duo.gc, 1, func() {
		drawRect(duo.gc, color.RGBA{A: 0xff, R: 0xe4, B: 0xe4, G: 0xe4})
		// Body <<<
		duo.comp.body.l.Layout(duo.gc, DuoUIsidebar(duo), DuoUIcontent(duo))
		// Body >>>
	})
}
