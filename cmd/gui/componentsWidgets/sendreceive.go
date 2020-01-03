package componentsWidgets

import (
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/widget"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/unit"
)

var (
	topLabel   = "testtopLabel"
	lineEditor = &widget.DuoUIeditor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(0))
)

func DuoUIsend(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) {
	layout.Flex{}.Layout(duo.DuoUIcontext,
		layout.Rigid(func() {
			helpers.DuoUIdrawRectangle(duo.DuoUIcontext, duo.DuoUIconstraints.Width.Max, 180, helpers.HexARGB("ff30cf30"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(duo.DuoUIcontext,
				layout.Rigid(func() {

					//cs := duo.DuoUIcontext.Constraints
					//DuoUIinputField(duo, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xcf, G: 0x30, B: 0xcf}, 9.9, 9.9, 9.9, 9.9, unit.Dp(0))
					ln.Layout(duo.DuoUIcontext, func() {
						cs := duo.DuoUIcontext.Constraints
						helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
						in.Layout(duo.DuoUIcontext, func() {
							helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//e := duo.DuoUItheme.DuoUIeditor("DUO address", "DUO dva")
							//e.Font.Style = text.Italic
							//e.Font.Size = unit.Dp(24)
							//e.Layout(duo.DuoUIcontext, lineEditor)
							//for _, e := range lineEditor.Events(duo.DuoUIcontext) {
							//	if e, ok := e.(widget.SubmitEvent); ok {
							//		topLabel = e.Text
							//		lineEditor.SetText("")
							//	}
							//}
						})
					})
					//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max,cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}, 0, 0, 0, 0, unit.Dp(0))
					//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ffbdbdbd"), [4]float32{0, 0, 0, 0}, unit.Dp(0))

					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(duo.DuoUIcontext,


						layout.Rigid(func() {
							//ln.Layout(duo.DuoUIcontext, func() {
							//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//in.Layout(duo.DuoUIcontext, func() {
							//	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//	e := duo.DuoUItheme.DuoUIeditor("DUO address", "DUO dva")
							//	e.Font.Style = text.Italic
							//	e.Font.Size = unit.Dp(24)
							//	e.Layout(duo.DuoUIcontext, lineEditor)
							//for _, e := range lineEditor.Events(duo.DuoUIcontext) {
							//	if e, ok := e.(widget.SubmitEvent); ok {
							//		topLabel = e.Text
							//		lineEditor.SetText("")
							//	}
							//}
							//})
							//})
						}),

						layout.Rigid(func() {
							//DuoUIdrawRect(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
							//helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//
							//in := layout.UniformInset(unit.Dp(8))
							//
							//in.Layout(duo.DuoUIcontext, func() {
							//	//DuoUIdrawRect(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("ff303030"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//	helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, helpers.HexARGB("fff4f4f4"), [4]float32{0, 0, 0, 0}, unit.Dp(0))
							//
							//	e := duo.DuoUItheme.DuoUIeditor("DUO amount", "DUO dva")
							//	e.Font.Style = text.Italic
							//	e.Font.Size = unit.Dp(24)
							//	e.Layout(duo.DuoUIcontext, lineEditor)
							//	for _, e := range lineEditor.Events(duo.DuoUIcontext) {
							//		if e, ok := e.(widget.SubmitEvent); ok {
							//			topLabel = e.Text
							//			lineEditor.SetText("")
							//		}
							//	}
							//})

						}))
				}))
		}),
	)
}
