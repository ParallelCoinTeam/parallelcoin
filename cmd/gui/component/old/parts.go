package component

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/gel"
	"github.com/p9c/pod/pkg/gui/gelook"
)

func SetPage(rc *rcd.RcVar, page *gelook.DuoUIpage) {
	page.Command()
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

func Editor(gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		th.DuoUIcontainer(8, "ffffffff").
			Layout(gtx, layout.NW, func() {
				width := gtx.Constraints.Width.Max
				e := th.DuoUIeditor(label,
					th.Colors["Black"], th.Colors["White"], width)
				e.Font.Typeface = th.Fonts["Mono"]
				e.TextSize = unit.Dp(12)
				layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
					e.Layout(gtx, editorController)
				})
				for _, e := range editorController.Events(gtx) {
					switch e.(type) {
					case gel.ChangeEvent:
						handler(e)
					}
				}
			})
	}
}

func StringsArrayEditor(gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.EditorEvent)) func() {
	return func() {
		th.DuoUIcontainer(8, th.Colors["White"]).Layout(gtx, layout.NW, func() {
			e := th.DuoUIeditor(label,
				th.Colors["Black"], th.Colors["White"], 16)
			e.Font.Typeface = th.Fonts["Mono"]
			layout.UniformInset(unit.Dp(4)).Layout(gtx, func() {
				e.Layout(gtx, editorController)
			})
			for _, e := range editorController.Events(gtx) {
				switch e.(type) {
				case gel.ChangeEvent:
					handler(e)
				}
			}
		})
	}
}

func ConsoleInput(gtx *layout.Context, th *gelook.DuoUItheme, editorController *gel.Editor, label string, handler func(gel.SubmitEvent)) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			e := th.DuoUIeditor(label, "Dark", "Light", 120)
			e.Font.Typeface = th.Fonts["Primary"]
			e.Color = gelook.HexARGB(th.Colors["Light"])
			e.Font.Style = text.Italic
			e.Layout(gtx, editorController)
			for _, e := range editorController.Events(gtx) {
				if e, ok := e.(gel.SubmitEvent); ok {
					handler(e)
					editorController.SetText("")
				}
			}
		})
	}
}

func Button(gtx *layout.Context, th *gelook.DuoUItheme, buttonController *gel.Button, font text.Typeface, textSize, pt, pr, pb, pl int, color, bgColor, label string, handler func()) func() {
	return func() {
		th.DuoUIcontainer(0, th.Colors[""])
		button := th.DuoUIbutton(font, label, color, bgColor,
			"", "", "", "",
			textSize, 0, 128, 48, pt, pr, pb, pl)
		for buttonController.Clicked(gtx) {
			handler()
		}
		button.Layout(gtx, buttonController)
	}
}

func MonoButton(gtx *layout.Context, th *gelook.DuoUItheme, buttonController *gel.Button, textSize int, color, bgColor, font, label string, handler func()) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
			//var button gelook.Button
			button := th.Button(label)
			if font != "" {
				button.Font.Typeface = th.Fonts[font]
			}
			if color != "" {
				button.Color = gelook.HexARGB(th.Colors[color])
			}
			if textSize != 0 {
				button.TextSize = unit.Dp(float32(textSize))
			}
			if bgColor != "" {
				button.Background = gelook.HexARGB(th.Colors[bgColor])
			}
			for buttonController.Clicked(gtx) {
				handler()
			}
			button.Layout(gtx, buttonController)
		})
	}
}

func Label(gtx *layout.Context, th *gelook.DuoUItheme, font text.Typeface, size float32, color, label string) func() {
	return func() {
		l := th.DuoUIlabel(unit.Dp(size), label)
		l.Font.Typeface = font
		l.Color = color
		l.Layout(gtx)
	}
}
