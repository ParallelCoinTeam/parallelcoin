package component

import (
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/controller"
	"github.com/p9c/pod/pkg/gui/theme"
	"image"
	"image/color"
)

func SetPage(rc *rcd.RcVar, page *theme.DuoUIpage) {
	rc.CurrentPage = page
}

func CurrentCurrentPageColor(showPage, page, color, currentPageColor string) (c string) {
	if showPage == page {
		c = currentPageColor
	} else {
		c = color
	}
	return
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func Editor(gtx *layout.Context, th *theme.DuoUItheme, editorControler *controller.Editor, label string, handler func(controller.SubmitEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			cs := gtx.Constraints
			theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 32, "fff4f4f4", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				theme.DuoUIdrawRectangle(gtx, cs.Width.Max, 30, "ffffffff", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
				e := th.DuoUIeditor(label)
				e.Font.Typeface = th.Fonts["Primary"]
				e.Font.Style = text.Italic
				e.Layout(gtx, editorControler)
				for _, e := range editorControler.Events(gtx) {
					if e, ok := e.(controller.SubmitEvent); ok {
						handler(e)
						editorControler.SetText("")
					}
				}
			})
		})
	}
}

func Button(gtx *layout.Context, th *theme.DuoUItheme, buttonController *controller.Button, font text.Typeface, textSize int, color, bgColor, label string, handler func()) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			var button theme.DuoUIbutton
			button = th.DuoUIbutton(font, label, color, bgColor, "", "", "", th.Colors["Light"], textSize, 0, 128, 48, 0, 0)
			for buttonController.Clicked(gtx) {
				handler()
			}
			button.Layout(gtx, buttonController)
		})

	}
}

func Label(gtx *layout.Context, th *theme.DuoUItheme, font text.Typeface, size float32, color, label string) func() {
	return func() {
		l := th.DuoUIlabel(unit.Dp(size), label)
		l.Font.Typeface = font
		l.Color = theme.HexARGB(color)
		l.Layout(gtx)
	}
}
