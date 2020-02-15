// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"image/color"
)

type DuoUIeditor struct {
	Font text.Font
	// Color is the text color.
	Color string
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA
	Text      string
	shaper    text.Shaper
}

func (t *DuoUItheme) DuoUIeditor(hint, txt string) DuoUIeditor {
	return DuoUIeditor{
		Font: text.Font{
			Size: t.TextSize,
		},
		Color:     t.Color.Text,
		shaper:    t.Shaper,
		Hint:      hint,
		HintColor: HexARGB(t.Color.Hint),
		Text:      txt,
	}
}

//
//func fill(gtx *layout.Context, col color.RGBA) {
//	cs := gtx.Constraints
//	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
//	dr := f32.Rectangle{
//		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
//	}
//	paint.ColorOp{Color: col}.Add(gtx.Ops)
//	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
//	gtx.Dimensions = layout.Dimensions{Size: d}
//}
//
//var typeRegistry = make(map[string]reflect.Type)
//
//func makeInstance(name string) interface{} {
//	v := reflect.New(typeRegistry[name]).Elem()
//	// Maybe fill in fields here if necessary
//	return v.Interface()
//}
//


func (e DuoUIeditor) Layout(gtx *layout.Context, editor *controller.Editor) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := controller.Label{Alignment: editor.Alignment}
	tl.Layout(gtx, e.shaper, e.Font, e.Hint)
	macro.Stop()
	if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
		gtx.Constraints.Width.Min = w
	}
	if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
		gtx.Constraints.Height.Min = h
	}
	editor.Layout(gtx, e.shaper, e.Font)
	if editor.Len() > 0 {
		paint.ColorOp{Color: HexARGB(e.Color)}.Add(gtx.Ops)
		editor.PaintText(gtx)
	} else {
		macro.Add()
	}
	paint.ColorOp{Color: HexARGB(e.Color)}.Add(gtx.Ops)
	editor.PaintCaret(gtx)
	stack.Pop()
}


//func (e DuoUIeditor) Layout(gtx *layout.Context, editor *controller.DuoUIeditor) {
//
//	bgcol := helpers.HexARGB("ffcfcfcf")
//	brcol := helpers.HexARGB("ff303030")
//	hmin := gtx.Constraints.Width.Min
//	vmin := gtx.Constraints.Height.Min
//	layout.Stack{Alignment: layout.Center}.Layout(gtx,
//		layout.Expanded(func() {
//			rr := float32(gtx.Px(unit.Dp(4)))
//			clip.Rect{
//				Rect: f32.Rectangle{Max: f32.Point{
//					X: float32(gtx.Constraints.Width.Min),
//					Y: float32(gtx.Constraints.Height.Min),
//				}},
//				NE: rr, NW: rr, SE: rr, SW: rr,
//			}.Op(gtx.Ops).Add(gtx.Ops)
//			fill(gtx, brcol)
//		}),
//		layout.Stacked(func() {
//			gtx.Constraints.Width.Min = hmin
//			gtx.Constraints.Height.Min = vmin
//			layout.Align(layout.Center).Layout(gtx, func() {
//				layout.Inset{Top: unit.Dp(1), Bottom: unit.Dp(1), Left: unit.Dp(1), Right: unit.Dp(1)}.Layout(gtx, func() {
//
//					layout.Stack{Alignment: layout.Center}.Layout(gtx,
//						layout.Expanded(func() {
//							rr := float32(gtx.Px(unit.Dp(4)))
//							clip.Rect{
//								Rect: f32.Rectangle{Max: f32.Point{
//									X: float32(gtx.Constraints.Width.Min),
//									Y: float32(gtx.Constraints.Height.Min),
//								}},
//								NE: rr, NW: rr, SE: rr, SW: rr,
//							}.Op(gtx.Ops).Add(gtx.Ops)
//							fill(gtx, bgcol)
//						}),
//						layout.Stacked(func() {
//							gtx.Constraints.Width.Min = hmin
//							gtx.Constraints.Height.Min = vmin
//							layout.Align(layout.Center).Layout(gtx, func() {
//								layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func() {
//									//paint.ColorOp{Color: col}.Add(gtx.Ops)
//									//widget.Label{}.Layout(gtx, btn.shaper, btn.Font, btn.Text)
//									edit := DuoUIeditor{
//										Font:      text.Font{},
//										Color:     color.RGBA{},
//										Hint:      "",
//										HintColor: color.RGBA{},
//										Text:      "",
//										shaper:    nil,
//									}
//									edit.Font.Style = text.Italic
//
//									edit.Layout(gtx, editor)
//									for _, e := range editor.Events(gtx) {
//										if _, ok := e.(widget.SubmitEvent); ok {
//											//topLabel = e.Text
//											editor.SetText("")
//
//											var stack op.StackOp
//											stack.Push(gtx.Ops)
//											var macro op.MacroOp
//											macro.Record(gtx.Ops)
//											paint.ColorOp{Color: edit.HintColor}.Add(gtx.Ops)
//											tl := widget.Label{Alignment: editor.Alignment}
//											tl.Layout(gtx, edit.shaper, edit.Font, edit.Hint)
//											macro.Stop()
//											if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
//												gtx.Constraints.Width.Min = w
//											}
//											if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
//												gtx.Constraints.Height.Min = h
//											}
//											editor.Layout(gtx, edit.shaper, edit.Font)
//											if editor.Len() > 0 {
//												paint.ColorOp{Color: edit.Color}.Add(gtx.Ops)
//												editor.PaintText(gtx)
//											} else {
//												macro.Add()
//											}
//											paint.ColorOp{Color: edit.Color}.Add(gtx.Ops)
//											editor.PaintCaret(gtx)
//											stack.Pop()
//
//										}
//									}
//								})
//							})
//						}),
//					)
//				})
//			})
//		}),
//	)
//
//}
