package gui

import (
	"image/color"
	
	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	
	"github.com/p9c/pod/pkg/gui/f32color"
)

// TextInput is a simple text input widget
type TextInput struct {
	*Window
	// Theme    *Theme
	font     text.Font
	textSize unit.Value
	// Color is the text color.
	color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	hint string
	// HintColor is the color of hint text.
	hintColor color.NRGBA
	editor    *Editor
	shaper    text.Shaper
}

// TextInput creates a simple text input widget
func (w *Window) TextInput(editor *Editor, hint string) *TextInput {
	var fon text.Font
	var e error
	if fon, e = w.collection.Font("bariol regular"); E.Chk(e) {
		panic(e)
	}
	ti := &TextInput{
		Window:    w,
		editor:    editor,
		textSize:  w.TextSize,
		font:      fon,
		color:     w.Colors.GetNRGBAFromName("DocText"),
		shaper:    w.shaper,
		hint:      hint,
		hintColor: w.Colors.GetNRGBAFromName("Hint"),
	}
	ti.Font("bariol regular")
	return ti
}

// Font sets the font for the text input widget
func (ti *TextInput) Font(font string) *TextInput {
	var fon text.Font
	var e error
	if fon, e = ti.Theme.collection.Font(font); !E.Chk(e) {
		ti.editor.font = fon
	}
	return ti
}

// TextScale sets the size of the text relative to the base font size
func (ti *TextInput) TextScale(scale float32) *TextInput {
	ti.textSize = ti.Theme.TextSize.Scale(scale)
	return ti
}

// Color sets the color to render the text
func (ti *TextInput) Color(color string) *TextInput {
	ti.color = ti.Theme.Colors.GetNRGBAFromName(color)
	return ti
}

// Hint sets the text to show when the box is empty
func (ti *TextInput) Hint(hint string) *TextInput {
	ti.hint = hint
	return ti
}

// HintColor sets the color of the hint text
func (ti *TextInput) HintColor(color string) *TextInput {
	ti.hintColor = ti.Theme.Colors.GetNRGBAFromName(color)
	return ti
}

// Fn renders the text input widget
func (ti *TextInput) Fn(c l.Context) l.Dimensions {
	defer op.Push(c.Ops).Pop()
	macro := op.Record(c.Ops)
	paint.ColorOp{Color: ti.hintColor}.Add(c.Ops)
	tl := Text{alignment: ti.editor.alignment}
	dims := tl.Fn(c, ti.shaper, ti.font, ti.textSize, ti.hint)
	call := macro.Stop()
	if w := dims.Size.X; c.Constraints.Min.X < w {
		c.Constraints.Min.X = w
	}
	if h := dims.Size.Y; c.Constraints.Min.Y < h {
		c.Constraints.Min.Y = h
	}
	dims = ti.editor.Layout(c, ti.shaper, ti.font, ti.textSize)
	disabled := c.Queue == nil
	if ti.editor.Len() > 0 {
		textColor := ti.color
		if disabled {
			textColor = f32color.MulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(c.Ops)
		ti.editor.PaintText(c)
	} else {
		call.Add(c.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: ti.color}.Add(c.Ops)
		ti.editor.PaintCaret(c)
	}
	return dims
}
