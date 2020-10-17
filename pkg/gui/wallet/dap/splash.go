package dap

import (
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
	"image"
)

func (d *dap) SplashScreen(gtx C) D {
	return box.BoxPanel(d.boot.UI.Theme, func(gtx C) D {
		if d.boot.Rc.IsReady {

		}
		return lyt.Format(gtx,
			"max(vflex(middle,r(inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp0dp,_)),f(1,inset(0dp0dp0dp0dp,_))))",
			logo(d.boot.UI.Theme),
			headline(d.boot.UI.Theme),
			liveLog(d.boot.UI.Theme),
		)
	})(gtx)
}

func logo(th *theme.Theme) func(gtx C) D {
	return func(gtx C) D {
		logo := th.Icons["Logo"]
		var dim D
		size := gtx.Px(unit.Dp((float32(gtx.Constraints.Max.Y) * 0.236)))
		logo.Color = helper.HexARGB(th.Colors["Bg"])
		logo.Layout(gtx, unit.Px(float32(size)))
		dim = D{
			Size: image.Point{X: size, Y: size},
		}
		return dim
	}
}

func headline(th *theme.Theme) func(gtx C) D {
	return func(gtx C) D {
		txt := theme.H1(th, "PLAN NINE FROM FAR, FAR AWAY SPACE")
		txt.Color = helper.HexARGB(th.Colors["Silver"])
		txt.TextSize = unit.Dp((float32(gtx.Constraints.Max.Y) * 0.118))

		txt.Alignment = text.Middle
		return txt.Layout(gtx)
	}
}

func liveLog(th *theme.Theme) func(gtx C) D {
	return noReturn
	//return func(gtx C) D {
	//layout.Flexed(1, component.DuoUIlogger(ui.rc, ui.ly.Context, ui.ly.Theme)),
	//}
}
