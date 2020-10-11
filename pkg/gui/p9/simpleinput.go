package p9

import (
	"image/color"

	l "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/f32color"
)

type _textInput struct {
	th       *Theme
	font     text.Font
	textSize unit.Value
	// Color is the text color.
	color color.RGBA
	// Hint contains the text displayed when the editor is empty.
	hint string
	// HintColor is the color of hint text.
	hintColor color.RGBA
	editor    *Editor
	shaper    text.Shaper
}

// SimpleInput creates a simple text input widget
func (th *Theme) SimpleInput(editor *Editor) *_textInput {
	e := &_textInput{
		th:        th,
		editor:    editor,
		textSize:  th.textSize,
		color:     th.Colors.Get("DocText"),
		shaper:    th.shaper,
		hint:      "hint",
		hintColor: th.Colors.Get("Hint"),
	}
	e.Font("bariol regular")
	return e
}

// Font sets the font for the text input widget
func (e *_textInput) Font(font string) *_textInput {
	for i := range e.th.collection {
		if e.th.collection[i].Font.Typeface == text.Typeface(font) {
			e.font = e.th.collection[i].Font
			break
		}
	}
	return e
}

// TextScale sets the size of the text relative to the base font size
func (e *_textInput) TextScale(scale float32) *_textInput {
	e.textSize = e.th.textSize.Scale(scale)
	return e
}

// Color sets the color to render the text
func (e *_textInput) Color(color string) *_textInput {
	e.color = e.th.Colors.Get(color)
	return e
}

// Hint sets the text to show when the box is empty
func (e *_textInput) Hint(hint string) *_textInput {
	e.hint = hint
	return e
}

// HintColor sets the color of the hint text
func (e *_textInput) HintColor(color string) *_textInput {
	e.hintColor = e.th.Colors.Get(color)
	return e
}

// Fn renders the text input widget
func (e *_textInput) Fn(c l.Context) l.Dimensions {
	defer op.Push(c.Ops).Pop()
	macro := op.Record(c.Ops)
	paint.ColorOp{Color: e.hintColor}.Add(c.Ops)
	tl := _text{alignment: e.editor.alignment}
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
