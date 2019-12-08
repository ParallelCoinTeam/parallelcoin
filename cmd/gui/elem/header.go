package elem

import (
	"github.com/p9c/pod/pkg/conte"
	"image/color"
)

func Header(cx *conte.Xt) {

	header := flv.Rigid(gtx, func() {
		drawRect(gtx, cs.Width.Max, 64, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
		drawRect(gtx, 64, 64, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30})
		th.IconButton(u.ico.logo).Layout(gtx, u.fab)
		//th.Image(Duo)

	})
}
