package duoui

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"
)

func DuoUIcontent(duo *DuoUI) layout.FlexChild {
	return duo.comp.Body.Layout.Flex(duo.gc, 1, func() {
		duo.comp.Content.Inset.Layout(duo.gc, func() {
			helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, B: 0xf4, G: 0xf4}, 0, 0, 0, 0)
			// Content <<<
			var content layout.FlexChild
			switch duo.menu.Current {
			case "overview":
				content = DuoUIoverview(duo)
			case "history":
				content = DuoUIhistory(duo)
			case "addressbook":
				content = DuoUIaddressbook(duo)
			case "explorer":
				content = DuoUIexplorer(duo)
			case "network":
				content = DuoUInetwork(duo)
			case "settings":
				content = DuoUIsettings(duo)
			default:
				content = DuoUIoverview(duo)
			}

			duo.comp.Content.Layout.Layout(duo.gc, content)
			// Content >>>
		})
	})
}
