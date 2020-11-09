package nav

import (
	"image"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"

	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

var (
	logoBtn = new(widget.Clickable)
)

type (
	D = layout.Dimensions
	C = layout.Context
	W = layout.Widget
)

type Logo struct {
	Title string
	Logo  string
}

func (n *Navigation) LogoLayout(th *theme.Theme) func(gtx C) D {
	return func(gtx C) D {
		return material.Clickable(gtx, logoBtn, func(gtx C) D {
			gtx.Constraints.Min.X = gtx.Constraints.Max.X
			for logoBtn.Clicked() {
				switch n.Wide {
				case false:
					n.Wide = true
				default:
					n.Wide = false
				}
			}
			inset := "inset(8dp18dp8dp8dp,_)"
			if n.Wide {
				inset = "inset(8dp8dp8dp8dp,_)"
			}
			return lyt.Format(gtx, inset, func(gtx C) D {
				logo := th.Icons["Logo"]
				if n.Wide {
					logo = th.Icons["Logo"]
				}
				var d D
				size := gtx.Px(unit.Dp(48))
				logo.Color = helper.HexARGB(th.Colors["PanelBg"])
				logo.Layout(gtx, unit.Px(float32(size)))
				d = D{
					Size: image.Point{X: size, Y: size},
				}
				return d
			})
		})
	}
}
