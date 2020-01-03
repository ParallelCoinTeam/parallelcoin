package duoui

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/unit"
)

func DuoUIcontent(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	// Content <<<
	var page = components.DuoUIpage{
		TxColor:      helpers.HexARGB("ff303030"),
		Width:        0,
		Height:       0,
		BgColor:      helpers.HexARGB("ffcfcfcf"),
		CornerRadius: unit.Dp(16),
	}
	switch duo.CurrentPage {
	case "Overview":
		page.Layout(duo.DuoUIcontext, func() { DuoUIoverview(duo, cx, rc) })
	case "Send":
		page.Layout(duo.DuoUIcontext, func() { DuoUIsend(duo, cx, rc) })
	case "Receive":
		page.Layout(duo.DuoUIcontext, func() { DuoUIreceive(duo, cx, rc) })
	case "History":
		page.Layout(duo.DuoUIcontext, func() { DuoUIhistory(duo, cx, rc) })
	case "AddressBook":
		page.Layout(duo.DuoUIcontext, func() { DuoUIaddressbook(duo, cx, rc) })
	case "Explorer":
		page.Layout(duo.DuoUIcontext, func() { DuoUIexplorer(duo, cx, rc) })
	case "Network":
		page.Layout(duo.DuoUIcontext, func() { DuoUInetwork(duo, cx, rc) })
	case "Console":
		page.Layout(duo.DuoUIcontext, func() { DuoUIconsole(duo, cx, rc) })
	case "Settings":
		page.Layout(duo.DuoUIcontext, func() { DuoUIsettings(duo, cx, rc) })
	default:
		page.Layout(duo.DuoUIcontext, func() { DuoUIoverview(duo, cx, rc) })
	}
	// Content >>>
}
