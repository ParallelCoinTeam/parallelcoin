package parts

import (
	"gioui.org/text"
	"github.com/gioapp/gel/helper"

	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

func (g *GioWallet) GetSettings() {
	// f, err := g.sh.ID()
	// checkError(err)
	// fv, ft, err := g.sh.Version()
	// checkError(err)
	// g.Status = Status{
	//	Title: "string",
	//	//HostingSize: "uint",
	//	PeerId:  f.ID,
	//	Version: fv + " " + ft,
	//	//Gateway:   ,
	//	//Api:       "string",
	//	Addresses: f.Addresses,
	//	Pub:       f.PublicKey,
	// }
}

func (g *GioWallet) settingsHeader() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		helper.Fill(gtx, helper.HexARGB(g.ui.Theme.Colors["PanelBg"]))
		return lyt.Format(gtx, "vflex(middle,r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_)),r(inset(5dp0dp5dp0dp,_))))",
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H5(g.ui.Theme, "LANGUAGE")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H4(g.ui.Theme, "BUTTONS")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H6(g.ui.Theme, "ANALYTICS")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.Body(g.ui.Theme, "Help improve this app by sending anonymous usage data")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.Body(g.ui.Theme, "No CIDs, filenames, or other personal information are collected. We want metrics to show us which features are useful to help us prioritise what to work on next, and system configuration information to guide our testing.")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				title := theme.H6(g.ui.Theme, "Configure what is collected")
				title.Alignment = text.Start
				return title.Layout(gtx)
			})
	})
}

func (g *GioWallet) settingsBody() func(gtx C) D {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		title := theme.H5(g.ui.Theme, "Duo CONFIG")
		title.Alignment = text.Start
		return title.Layout(gtx)
	}
}
