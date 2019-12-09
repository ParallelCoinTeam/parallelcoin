package elm

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/mod"
	"image/color"
)

func DuoUIheader(d *mod.DuoUI) layout.FlexChild {
	return d.Layouts.View.Rigid(&d.Gtx , func() {
		hlp.DuoUIdrawRect(&d.Gtx , d.Gtx.Constraints.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
		hlp.DuoUIdrawRect(&d.Gtx , 64, 64, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30})
		d.Theme.IconButton(d.Ico.Logo).Layout(&d.Gtx , d.Buttons.Logo)
	})
}
