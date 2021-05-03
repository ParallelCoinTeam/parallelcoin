package gel

import (
	"image/color"

	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/text"
	"github.com/p9c/gio/unit"

	"github.com/p9c/gel/f32color"
)

// TextInput is a simple text input widget
type TextInput struct {
	*Window
	font     text.Font
	textSize unit.Value
	// Color is the text color.
	color color.NRGBA
	// Hint contains the text displayed when the editor is empty.
	hint string
	// HintColor is the color of hint text.
	hintColor color.NRGBA
	// SelectionColor is the color of the background for selected text.
	selectionColor color.NRGBA
	editor         *Editor
	shaper         text.Shaper
}

// TextInput creates a simple text input widget
func (w *Window) TextInput(editor *Editor, hint string) *TextInput {
	var fon text.Font
	var e error
	if fon, e = w.collection.Font("bariol regular"); E.Chk(e) {
		panic(e)
	}
	ti := &TextInput{
		Window:         w,
		editor:         editor,
		textSize:       w.TextSize,
		font:           fon,
		color:          w.Colors.GetNRGBAFromName("DocText"),
		shaper:         w.shaper,
		hint:           hint,
		hintColor:      w.Colors.GetNRGBAFromName("Hint"),
		selectionColor: w.Colors.GetNRGBAFromName("scrim"),
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

// SelectionColor sets the color to render the text
func (ti *TextInput) SelectionColor(color string) *TextInput {
	ti.selectionColor = ti.Theme.Colors.GetNRGBAFromName(color)
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
func (ti *TextInput) Fn(gtx l.Context) l.Dimensions {
	defer op.Save(gtx.Ops).Load()
	macro := op.Record(gtx.Ops)
	paint.ColorOp{Color: ti.hintColor}.Add(gtx.Ops)
	var maxlines int
	if ti.editor.singleLine {
		maxlines = 1
	}
	tl := Label{
		Window:    ti.Window,
		font:      ti.font,
		color:     ti.hintColor,
		alignment: ti.editor.alignment,
		maxLines:  maxlines,
		text:      ti.hint,
		textSize:  ti.textSize,
		shaper:    ti.shaper,
	}
	dims := tl.Fn(gtx)
	call := macro.Stop()
	if w := dims.Size.X; gtx.Constraints.Min.X < w {
		gtx.Constraints.Min.X = w
	}
	if h := dims.Size.Y; gtx.Constraints.Min.Y < h {
		gtx.Constraints.Min.Y = h
	}
	dims = ti.editor.Layout(gtx, ti.shaper, ti.font, ti.textSize)
	disabled := gtx.Queue == nil
	if ti.editor.Len() > 0 {
		paint.ColorOp{Color: blendDisabledColor(disabled, ti.selectionColor)}.Add(gtx.Ops)
		ti.editor.PaintSelection(gtx)
		paint.ColorOp{Color: blendDisabledColor(disabled, ti.color)}.Add(gtx.Ops)
		ti.editor.PaintText(gtx)
	} else {
		call.Add(gtx.Ops)
	}
	if !disabled {
		paint.ColorOp{Color: ti.color}.Add(gtx.Ops)
		ti.editor.PaintCaret(gtx)
	}
	return dims
}

func blendDisabledColor(disabled bool, c color.NRGBA) color.NRGBA {
	if disabled {
		return f32color.Disabled(c)
	}
	return c
}
