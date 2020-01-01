package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
)

func DuoUIcontent(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// Content <<<
		switch duo.CurrentPage {
		case "Overview":
			DuoUIoverview(duo, cx, rc)
		case "History":
			DuoUIhistory(duo, cx, rc)
		case "AddressBook":
			DuoUIaddressbook(duo, cx, rc)
		case "Explorer":
			DuoUIexplorer(duo, cx, rc)
		case "Network":
			DuoUInetwork(duo, cx, rc)
		case "Console":
			DuoUIconsole(duo, cx, rc)
		case "Settings":
			DuoUIsettings(duo, cx, rc)
		default:
			DuoUIoverview(duo, cx, rc)
		}
	// Content >>>
}
