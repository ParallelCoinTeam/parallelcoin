package duoui

import (
	"github.com/p9c/gio-parallel/layout"
	"image/color"
)

func DuoUIcontent(duo *DuoUI) layout.FlexChild {
	return duo.comp.body.l.Flex(duo.gc, 1, func() {
		duo.comp.content.i.Layout(duo.gc, func() {
			drawRect(duo.gc, color.RGBA{A: 0xff, R: 0xf4, B: 0xf4, G: 0xf4})
			// Content <<<

			var content layout.FlexChild
			switch duo.menu.current {
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

			duo.comp.content.l.Layout(duo.gc, content)
			// Content >>>
		})
	})
}
