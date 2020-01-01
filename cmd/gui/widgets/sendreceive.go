package widgets

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"

)

var (
	topLabel   = "testtopLabel"
	//lineEditor = &widget.Editor{
	//	SingleLine: true,
	//	Submit:     true,
	//}
	list = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(0))
)

func DuoUIsendreceive(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar){
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Flexed(1, func() {
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, 180, helpers.HexARGB("ff30cf30"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

			layout.Flex{
				Axis:layout.Vertical,
			}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {
					//helpers.DuoUIinputField(duo, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
					//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.Cs.Width.Max,duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}, 0, 0, 0, 0, unit.Dp(0))
				}),
			)



			//address := duo.comp.sendReceive.l.Flex(duo.DuoUIcontext, 0.3, func() {
			//ln.Layout(duo.DuoUIcontext, func() {
			//	DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
			//	in.Layout(duo.DuoUIcontext, func() {
			//		DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
			//		e := duo.Th.Editor("DUO address")
			//		e.Font.Style = text.Italic
			//		e.Font.Size = unit.Dp(24)
			//		e.Layout(duo.DuoUIcontext, lineEditor)
			//		for _, e := range lineEditor.Events(duo.DuoUIcontext) {
			//			if e, ok := e.(widget.SubmitEvent); ok {
			//				topLabel = e.Text
			//				lineEditor.SetText("")
			//			}
			//		}
			//	})
			//})
			//})
			//amount := duo.comp.sendReceive.l.Flex(duo.DuoUIcontext, 0.3, func() {
			//	DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
			//	in := layout.UniformInset(unit.Dp(8))
			//
			//	in.Layout(duo.DuoUIcontext, func() {
			//		DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
			//
			//		e := duo.Th.Editor("DUO amount")
			//		e.Font.Style = text.Italic
			//		e.Font.Size = unit.Dp(24)
			//		e.Layout(duo.DuoUIcontext, lineEditor)
			//		for _, e := range lineEditor.Events(duo.DuoUIcontext) {
			//			if e, ok := e.(widget.SubmitEvent); ok {
			//				topLabel = e.Text
			//				lineEditor.SetText("")
			//			}
			//		}
			//	})
			//
			//})
			//buttons := duo.comp.sendReceive.l.Flex(duo.DuoUIcontext, 0.4, func() {
			//	DuoUIdrawRect(duo.DuoUIcontext, duo.Cs.Width.Max, duo.Cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}, 0, 0, 0, 0, unit.Dp(0))
			//})
			//duo.comp.sendReceive.l.Layout(duo.DuoUIcontext, address)

		}),
	)
}
