package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"image/color"
)

func DuoUIbody(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Flexed(1, func() {
			helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xe4, G: 0xe4, B: 0xe4}, 0, 0, 0, 0)
			// Body <<<
			DuoUIsidebar(duo)
			DuoUIcontent(duo)
			// Body >>>
		}),
	)
}
