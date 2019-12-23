package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIbody(duo *DuoUI) {
	duo.comp.Body.Layout.Layout(duo.gc,
		layout.Rigid(func() {
			DuoUIsidebar(duo)
		}),
		layout.Flexed(1, func() {
			DuoUIcontent(duo)
		}),
	)
}
