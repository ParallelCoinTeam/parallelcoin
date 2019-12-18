package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"
)

func DuoUIbody(duo *DuoUI) layout.FlexChild {
	return duo.comp.View.Layout.Flex(duo.gc, 1, func() {
		helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xe4, G: 0xe4, B: 0xe4}, 0, 0, 0, 0)
		// Body <<<
		duo.comp.Body.Layout.Layout(duo.gc, DuoUIsidebar(duo), DuoUIcontent(duo))
		// Body >>>
	})
}
