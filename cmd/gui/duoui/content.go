package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIcontent(duo *DuoUI) {
	layout.Flex{}.Layout(duo.gc,
		layout.Flexed(1, func() {
			duo.comp.Content.Inset.Layout(duo.gc, func() {
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
