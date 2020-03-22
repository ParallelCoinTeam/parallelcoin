// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"image/color"

	"github.com/p9c/pod/pkg/gel"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type DuoUIeditor struct {
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA
	// Width is the number of monospaced characters that will display
	Width  int
	shaper text.Shaper
}

func (t *DuoUItheme) DuoUIeditor(hint, color string, width int) DuoUIeditor {
	return DuoUIeditor{
		TextSize:  t.TextSize,
		Color:     HexARGB(color),
		shaper:    t.Shaper,
		Hint:      hint,
		HintColor: HexARGB(t.Colors["Hint"]),
		Width:     width,
	}
}

func (e DuoUIeditor) Layout(gtx *layout.Context, editor *gel.Editor) {
	var stack op.StackOp
	stack.Push(gtx.Ops)
	var macro op.MacroOp
	macro.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := gel.Label{Alignment: editor.Alignment}
	tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	macro.Stop()
	//if w := gtx.Dimensions.Size.X; gtx.Constraints.Width.Min < w {
	//	gtx.Constraints.Width.Min = w
	//}
	//if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
	//	gtx.Constraints.Height.Min = h
	//}
	//if w := gtx.Constraints.Width.Min; 16*e.Width < w {
	gtx.Constraints.Width.Max = 16 * e.Width
	gtx.Constraints.Width.Min = 16 * e.Width
	//}
	//if h := gtx.Dimensions.Size.Y; gtx.Constraints.Height.Min < h {
	gtx.Constraints.Height.Max = 40
	gtx.Constraints.Height.Min = 40
	//}
	editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	if editor.Len() > 0 {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		editor.PaintText(gtx)
	} else {
		macro.Add()
	}
	paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
	editor.PaintCaret(gtx)
	stack.Pop()
}
