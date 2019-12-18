package duoui

import (
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/text"
	"github.com/p9c/pod/pkg/gio/unit"
	"github.com/p9c/pod/pkg/gio/widget"
	"github.com/p9c/pod/cmd/gui/helpers"
	"image/color"
)

var (
	topLabel   = "testtopLabel"
	lineEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	list = &layout.List{
		Axis: layout.Vertical,
	}
	ln = layout.UniformInset(unit.Dp(1))
	in = layout.UniformInset(unit.Dp(8))
)

func DuoUIsendreceive(duo *DuoUI) layout.FlexChild {
	return duo.comp.OverviewTop.Layout.Flex(duo.gc, 0.6, func() {
		helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, 180, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}, 0, 0, 0, 0)

		widgets := []func(){
			func() {

				helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				ln.Layout(duo.gc, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
					in.Layout(duo.gc, func() {
						e := duo.th.Editor("Hint")
						e.Font.Style = text.Italic
						e.Font.Size = unit.Dp(24)
						e.Layout(duo.gc, lineEditor)
						for _, e := range lineEditor.Events(duo.gc) {
							if e, ok := e.(widget.SubmitEvent); ok {
								topLabel = e.Text
								lineEditor.SetText("")
							}
						}
					})
				})
			},
			func() {

				helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
				ln.Layout(duo.gc, func() {
					helpers.DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
					in.Layout(duo.gc, func() {
						e := duo.th.Editor("Hint")
						e.Font.Style = text.Italic
						e.Font.Size = unit.Dp(24)
						e.Layout(duo.gc, lineEditor)
						for _, e := range lineEditor.Events(duo.gc) {
							if e, ok := e.(widget.SubmitEvent); ok {
								topLabel = e.Text
								lineEditor.SetText("")
							}
						}
					})
				})
			},
		}
		list.Layout(duo.gc, len(widgets), func(i int) {
			layout.UniformInset(unit.Dp(16)).Layout(duo.gc, widgets[i])
		})
		//address := duo.comp.sendReceive.l.Flex(duo.gc, 0.3, func() {
		//ln.Layout(duo.gc, func() {
		//	DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
		//	in.Layout(duo.gc, func() {
		//		DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
		//		e := duo.th.Editor("DUO address")
		//		e.Font.Style = text.Italic
		//		e.Font.Size = unit.Dp(24)
		//		e.Layout(duo.gc, lineEditor)
		//		for _, e := range lineEditor.Events(duo.gc) {
		//			if e, ok := e.(widget.SubmitEvent); ok {
		//				topLabel = e.Text
		//				lineEditor.SetText("")
		//			}
		//		}
		//	})
		//})
		//})
		//amount := duo.comp.sendReceive.l.Flex(duo.gc, 0.3, func() {
		//	DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}, 9.9, 9.9, 9.9, 9.9)
		//	in := layout.UniformInset(unit.Dp(8))
		//
		//	in.Layout(duo.gc, func() {
		//		DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0xf4, G: 0xf4, B: 0xf4}, 9.9, 9.9, 9.9, 9.9)
		//
		//		e := duo.th.Editor("DUO amount")
		//		e.Font.Style = text.Italic
		//		e.Font.Size = unit.Dp(24)
		//		e.Layout(duo.gc, lineEditor)
		//		for _, e := range lineEditor.Events(duo.gc) {
		//			if e, ok := e.(widget.SubmitEvent); ok {
		//				topLabel = e.Text
		//				lineEditor.SetText("")
		//			}
		//		}
		//	})
		//
		//})
		//buttons := duo.comp.sendReceive.l.Flex(duo.gc, 0.4, func() {
		//	DuoUIdrawRect(duo.gc, duo.cs.Width.Max, duo.cs.Height.Max, color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}, 0, 0, 0, 0)
		//})
		//duo.comp.sendReceive.l.Layout(duo.gc, address)

	})
}
