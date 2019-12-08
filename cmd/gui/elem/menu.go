package elem

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/conte"
	"image/color"
)

func DuoUImenu(cx *conte.Xt) (err error) {
	main := flv.Rigid(gtx, func() {
		//balance := flh.Rigid(gtx, func() {
		//	in.Layout(gtx, func() {
		//		th.H3("balance :" + r.balance).Layout(gtx)
		//	})
		//})
		sidebar := flh.Rigid(gtx, func() {

			drawRect(gtx, 64, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})

			flm := layout.Flex{Axis: layout.Vertical}
			overview := flm.Rigid(gtx, func() {
				th.IconButton(u.ico.overview).Layout(gtx, u.fab)
			})
			history := flm.Rigid(gtx, func() {
				th.IconButton(u.ico.history).Layout(gtx, u.fab)
			})
			network := flm.Rigid(gtx, func() {
				th.IconButton(u.ico.network).Layout(gtx, u.fab)
			})
			settings := flm.Rigid(gtx, func() {
				th.IconButton(u.ico.settings).Layout(gtx, u.fab)
			})
			flm.Layout(gtx, overview, history, network, settings)

		})
		content := flh.Rigid(gtx, func() {
			in.Layout(gtx, func() {
				drawRect(gtx, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf})
				th.H3("balance :" + r.balance).Layout(gtx)
			})
		})

		flh.Layout(gtx, sidebar, content)
	}
