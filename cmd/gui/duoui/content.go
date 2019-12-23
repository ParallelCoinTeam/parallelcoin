package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"

	"image/color"
)

func DuoUIcontent(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Flexed(1, func() {
			duo.comp.Content.Inset.Layout(duo.gc, func() {
				helpers.DuoUIdrawRectangle(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, B: 0xf4, G: 0xf4}, 0, 0, 0, 0, unit.Dp(0))
				// Content <<<
				switch duo.menu.Current {
				case "overview":
					DuoUIoverview(duo)
				case "history":
					DuoUIhistory(duo)
				case "addressbook":
					DuoUIaddressbook(duo)
				case "explorer":
					DuoUIexplorer(duo)
				case "network":
					DuoUInetwork(duo)
				case "console":
					DuoUIconsole(duo)
				case "settings":
					DuoUIsettings(duo)
				default:
					DuoUIoverview(duo)
				}
				// Content >>>
			})
		}),
	)
}
