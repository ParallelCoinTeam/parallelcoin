package duoui

import (
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
)

func DuoUIcontent(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	// Content <<<
	return func() {
		var page = parallel.DuoUIpage{
			TxColor:      parallel.HexARGB("ff303030"),
			Width:        0,
			Height:       0,
			BgColor:      parallel.HexARGB("ffcfcfcf"),
			CornerRadius: unit.Dp(16),
		}
		switch duo.CurrentPage {
		case "Overview":
			page.Layout(duo.DuoUIcontext,DuoUIoverview(duo, cx, rc) )
		case "Send":
			page.Layout(duo.DuoUIcontext,DuoUIsend(duo, cx, rc) )
		case "Receive":
			page.Layout(duo.DuoUIcontext,DuoUIreceive(duo, cx, rc) )
		case "History":
			page.Layout(duo.DuoUIcontext,DuoUIhistory(duo, cx, rc) )
		case "AddressBook":
			page.Layout(duo.DuoUIcontext,DuoUIaddressbook(duo, cx, rc) )
		case "Explorer":
			page.Layout(duo.DuoUIcontext,DuoUIexplorer(duo, cx, rc) )
		case "Network":
			page.Layout(duo.DuoUIcontext,DuoUInetwork(duo, cx, rc) )
		case "Console":
			page.Layout(duo.DuoUIcontext,DuoUIconsole(duo, cx, rc) )
		case "Settings":
			page.Layout(duo.DuoUIcontext,DuoUIsettings(duo, cx, rc) )
		default:
			page.Layout(duo.DuoUIcontext,DuoUIoverview(duo, cx, rc) )
		}
	}
	// Content >>>
}
