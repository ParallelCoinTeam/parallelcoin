package p9

import (
	"image/color"

	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type _label struct {
	th *Theme
	// Face defines the text style.
	font text.Font
	// Color is the text color.
	color color.RGBA
	// Alignment specify the text alignment.
	alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	maxLines int
	text     string
	textSize unit.Value

	shaper text.Shaper
}

// _text creates a label that prints a block of text
func (th *Theme) Label() (l *_label) {
	var f text.Font
	for i := range th.collection {
		if th.collection[i].Font.Typeface == "plan9" {
			f = th.collection[i].Font
		}
	}
	return &_label{
		th:       th,
		text:     "",
		font:     f,
		color:    th.Colors.Get("DocText"),
		textSize: unit.Sp(1),
		shaper:   th.shaper,
	}
}

// Text sets the text to render in the label
func (l *_label) Text(text string) *_label {
	l.text = text
	return l
}

// TextScale sets the size of the text relative to the base font size
func (l *_label) TextScale(scale float32) *_label {
	l.textSize = l.th.textSize.Scale(scale)
	return l
}

// MaxLines sets the maximum number of lines to render
func (l *_label) MaxLines(maxLines int) *_label {
	l.maxLines = maxLines
	return l
}

// Alignment sets the text alignment, left, right or centered
func (l *_label) Alignment(alignment text.Alignment) *_label {
	l.alignment = alignment
	return l
}

// Color sets the color of the label font
func (l *_label) Color(color string) *_label {
	l.color = l.th.Colors.Get(color)
	return l
}

// Font sets the font out of the available font collection
func (l *_label) Font(font string) *_label {
	var f text.Font
	for i := range l.th.collection {
		// Debug(th.Collection[i].Font)
		if l.th.collection[i].Font.Typeface == text.Typeface(font) {
			f = l.th.collection[i].Font
		}
	}
	l.font = f
	return l
}

func (th *Theme) H1(txt string) (l *_label) {
	l = th.Label().TextScale(96.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) H2(txt string) (l *_label) {
	l = th.Label().TextScale(60.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) H3(txt string) (l *_label) {
	l = th.Label().TextScale(48.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) H4(txt string) (l *_label) {
	l = th.Label().TextScale(34.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) H5(txt string) (l *_label) {
	l = th.Label().TextScale(24.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) H6(txt string) (l *_label) {
	l = th.Label().TextScale(20.0 / 16.0).Font("plan9").Text(txt)
	return
}

func (th *Theme) Body1(txt string) (l *_label) {
	l = th.Label().TextScale(1).Font("bariol regular").Text(txt)
	return
}

func (th *Theme) Body2(txt string) (l *_label) {
	l = th.Label().TextScale(14.0 / 16.0).Font("bariol regular").Text(txt)
	return
}

func (th *Theme) Caption(txt string) (l *_label) {
	l = th.Label().TextScale(12.0 / 16.0).Font("bariol regular").Text(txt)
	return
}

// Fn renders the label as specified
func (l *_label) Fn(gtx l.Context) l.Dimensions {
	paint.ColorOp{Color: l.color}.Add(gtx.Ops)
	tl := _text{alignment: l.alignment, maxLines: l.maxLines}
	return tl.Fn(gtx, l.shaper, l.font, l.textSize, l.text)
}
