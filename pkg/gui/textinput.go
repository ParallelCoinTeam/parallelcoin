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
	var err error
	var fon text.Font
	if fon, err = w.collection.Font("bariol regular"); Check(err) {
		panic(err)
	}
	e := &TextInput{
		Window:    w,
		editor:    editor,
		textSize:  w.TextSize,
		font:      fon,
		color:     w.Colors.GetNRGBAFromName("DocText"),
		shaper:    w.shaper,
		hint:      hint,
		hintColor: w.Colors.GetNRGBAFromName("Hint"),
	}
	e.Font("bariol regular")
	return e
}

// Font sets the font for the text input widget
func (e *TextInput) Font(font string) *TextInput {
	if fon, err := e.Theme.collection.Font(font); !Check(err) {
		e.editor.font = fon
	}
	return e
}

// TextScale sets the size of the text relative to the base font size
func (e *TextInput) TextScale(scale float32) *TextInput {
	e.textSize = e.Theme.TextSize.Scale(scale)
	return e
}

// Color sets the color to render the text
func (e *TextInput) Color(color string) *TextInput {
	e.color = e.Theme.Colors.GetNRGBAFromName(color)
	return e
}

// Hint sets the text to show when the box is empty
func (e *TextInput) Hint(hint string) *TextInput {
	e.hint = hint
	return e
}

// HintColor sets the color of the hint text
func (e *TextInput) HintColor(color string) *TextInput {
	e.hintColor = e.Theme.Colors.GetNRGBAFromName(color)
	return e
}

// Fn renders the text input widget
func (e *TextInput) Fn(c l.Context) l.Dimensions {
	defer op.Push(c.Ops).Pop()
	macro := op.Record(c.Ops)
	paint.ColorOp{Color: e.hintColor}.Add(c.Ops)
	tl := Text{alignment: e.editor.alignment}
	dims := tl.Fn(c, e.shaper, e.font, e.textSize, e.hint)
	call := macro.Stop()
	if w := dims.Size.X; c.Constraints.Min.X < w {
		c.Constraints.Min.X = w
	}
	if h := dims.Size.Y; c.Constraints.Min.Y < h {
		c.Constraints.Min.Y = h
	}
	dims = e.editor.Layout(c, e.shaper, e.font, e.textSize)
	disabled := c.Queue == nil
	if e.editor.Len() > 0 {
		textColor := e.color
		if disabled {
			textColor = f32color.MulAlpha(textColor, 150)
		}
		paint.ColorOp{Color: textColor}.Add(c.Ops)
		e.editor.PaintText(c)
	} else {
		call.Add(c.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: e.color}.Add(c.Ops)
		e.editor.PaintCaret(c)
	}
	return dims
}
