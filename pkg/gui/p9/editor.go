package p9

import (
	"image/color"

	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _editor struct {
	Font     text.Font
	TextSize unit.Value
	// Color is the text color.
	Color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	Hint string
	// HintColor is the color of hint text.
	HintColor color.RGBA
	Editor    *widget.Editor

	shaper text.Shaper
}

func (th *Theme) Editor(editor *widget.Editor, hint string) _editor {
	return _editor{
		Editor:    editor,
		TextSize:  th.TextSize,
		Color:     th.Colors.Get("Text"),
		shaper:    th.Shaper,
		Hint:      hint,
		HintColor: th.Colors.Get("Hint"),
	}
}

func (e _editor) Fn(gtx l.Context) l.Dimensions {
	defer op.Push(gtx.Ops).Pop()
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: e.HintColor}.Add(gtx.Ops)
	tl := widget.Label{Alignment: e.Editor.Alignment}
	dims := tl.Layout(gtx, e.shaper, e.Font, e.TextSize, e.Hint)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = e.Editor.Layout(gtx, e.shaper, e.Font, e.TextSize)
	disabled := gtx.Queue == nil
	if e.Editor.Len() > 0 {
		textColor := e.Color
		if disabled {
			textColor = f32color.MulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(gtx.Ops)
		e.Editor.PaintText(gtx)
	} else {
		call.Add(gtx.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: e.Color}.Add(gtx.Ops)
		e.Editor.PaintCaret(gtx)
	}
	return dims
}
