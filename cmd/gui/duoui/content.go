package duoui

import (
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
)

func DuoUIcontent(duo *DuoUI, cx *conte.Xt, rc *rcd.RcVar) func() {
	// Content <<<
	return func() {
		var page = parallel.DuoUIpage{
			TxColor: parallel.HexARGB("ff303030"),
			BgColor: parallel.HexARGB("ffcfcfcf"),
		}
		switch duo.m.CurrentPage {
		case "OVERVIEW":
			page.Layout(DuoUIoverview(duo, cx, rc))
		case "SEND":
			page.Layout(DuoUIsend(duo, cx, rc))
		case "RECEIVE":
			page.Layout(DuoUIreceive(duo))
		case "HISTORY":
			page.Layout(DuoUIhistory(duo, cx, rc))
		case "ADDRESSBOOK":
			page.Layout(DuoUIaddressbook(duo))
		case "EXPLORER":
			page.Layout(duo.DuoUIexplorer(cx, rc))
		case "NETWORK":
			page.Layout(DuoUInetwork(duo))
		case "CONSOLE":
			page.Layout(DuoUIconsole(duo, cx, rc))
		case "TRACE":
			page.Layout(DuoUItrace(duo, cx, rc))
		case "SETTINGS":
			page.Layout(DuoUIsettings(duo, cx, rc))
		default:
			page.Layout(DuoUIoverview(duo, cx, rc))
		}
	}
	// Content >>>
}
