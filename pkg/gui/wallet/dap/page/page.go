package page

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/wallet/container"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)

type Page struct {
	Title  string
	Header W
	Body   W
	Footer W
}

func (p *Page) P(th *theme.Theme, ly string) func(gtx C) D {
	return container.C().
		//OutsideColor(th.Colors["PanelBg"]).
		//BorderColor(th.Colors["Border"]).
		InsideColor(th.Colors["PanelBg"]).
		Margin(0).
		Border(0).
		Padding(4).
		Layout(func(gtx C) D {
			return lyt.Format(gtx, ly, p.Header, p.Body, p.Footer)
		})
}

//
//func pageButton(th *theme.Theme, b *widget.Clickable, f func(), icon, page string) func(gtx C) D {
//	return func(gtx C) D {
//		btn := material.IconButton(th.T, b, th.Icons[icon])
//		btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
//		btn.Size = unit.Dp(21)
//		btn.Background = helper.HexARGB(th.Colors["Secondary"])
//		for b.Clicked() {
//			f()
//			//d.UI.N.CurrentPage = page
//		}
//		return btn.Layout(gtx)
//	}
//}
