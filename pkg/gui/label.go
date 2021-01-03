package gui

import (
	"image/color"
	
	l "gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type Label struct {
	*Window
	// Face defines the text style.
	font text.Font
	// Color is the text color.
	color color.NRGBA
	// Alignment specify the text alignment.
	alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	maxLines int
	text     string
	textSize unit.Value
	shaper   text.Shaper
}

// Label creates a label that prints a block of text
func (w *Window) Label() (l *Label) {
	var f text.Font
	if fon, err := w.Theme.collection.Font("plan9"); !Check(err) {
		f = fon
	}
	return &Label{
		Window:   w,
		text:     "",
		font:     f,
		color:    w.Colors.GetNRGBAFromName("DocText"),
		textSize: unit.Sp(1),
		shaper:   w.shaper,
	}
}

// Text sets the text to render in the label
func (l *Label) Text(text string) *Label {
	l.text = text
	return l
}

// TextScale sets the size of the text relative to the base font size
func (l *Label) TextScale(scale float32) *Label {
	l.textSize = l.Theme.TextSize.Scale(scale)
	return l
}

// MaxLines sets the maximum number of lines to render
func (l *Label) MaxLines(maxLines int) *Label {
	l.maxLines = maxLines
	return l
}

// Alignment sets the text alignment, left, right or centered
func (l *Label) Alignment(alignment text.Alignment) *Label {
	l.alignment = alignment
	return l
}

// Color sets the color of the label font
func (l *Label) Color(color string) *Label {
	l.color = l.Theme.Colors.GetNRGBAFromName(color)
	return l
}

// Font sets the font out of the available font collection
func (l *Label) Font(font string) *Label {
	if fon, err := l.Theme.collection.Font(font); !Check(err) {
		l.font = fon
	}
	return l
}

type ScaleType map[string]float32

var (
	Scales = ScaleType{
		"H1":      96.0 / 16.0,
		"H2":      60.0 / 16.0,
		"H3":      48.0 / 16.0,
		"H4":      34.0 / 16.0,
		"H5":      24.0 / 16.0,
		"H6":      20.0 / 16.0,
		"Body1":   1,
		"Body2":   14.0 / 16.0,
		"Caption": 12.0 / 16.0,
	}
)

func (w *Window) H1(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H1"]).Font("plan9").Text(txt)
	return
}

func (w *Window) H2(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H2"]).Font("plan9").Text(txt)
	return
}

func (w *Window) H3(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H3"]).Font("plan9").Text(txt)
	return
}

func (w *Window) H4(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H4"]).Font("plan9").Text(txt)
	return
}

func (w *Window) H5(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H5"]).Font("plan9").Text(txt)
	return
}

func (w *Window) H6(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["H6"]).Font("plan9").Text(txt)
	return
}

func (w *Window) Body1(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Body1"]).Font("bariol regular").Text(txt)
	return
}

func (w *Window) Body2(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Body2"]).Font("bariol regular").Text(txt)
	return
}

func (w *Window) Caption(txt string) (l *Label) {
	l = w.Label().TextScale(Scales["Caption"]).Font("bariol regular").Text(txt)
	return
}

// Fn renders the label as specified
func (l *Label) Fn(gtx l.Context) l.Dimensions {
	paint.ColorOp{Color: l.color}.Add(gtx.Ops)
	tl := Text{alignment: l.alignment, maxLines: l.maxLines}
	return tl.Fn(gtx, l.shaper, l.font, l.textSize, l.text)
}
