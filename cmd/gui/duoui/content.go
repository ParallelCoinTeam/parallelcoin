package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
)

func DuoUIcontent(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.Gc,
		layout.Flexed(1, func() {
			duo.Comp.Content.Inset.Layout(duo.Gc, func() {
				// Content <<<
				switch duo.Menu.Current {
				case "overview":
					DuoUIoverview(duo,cx,rc)
				case "history":
					DuoUIhistory(duo,cx,rc)
				case "addressbook":
					DuoUIaddressbook(duo,cx,rc)
				case "explorer":
					DuoUIexplorer(duo,cx,rc)
				case "network":
					DuoUInetwork(duo,cx,rc)
				case "console":
					DuoUIconsole(duo,cx,rc)
				case "settings":
					DuoUIsettings(duo,cx,rc)
				default:
					DuoUIoverview(duo,cx,rc)
				}
				// Content >>>
			})
		}),
	)
}
