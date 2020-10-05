// SPDX-License-Identifier: Unlicense OR MIT

package plan9

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
)

type LabelStyle struct {
	// Face defines the text style.
	Font text.Font
	// Color is the text color.
	Color color.RGBA
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Value

	shaper text.Shaper
}

func H1(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(96.0/16.0), "plan9", txt)
	return
}

func H2(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(60.0/16.0), "plan9", txt)
	return
}

func H3(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(48.0/16.0), "plan9", txt)
	return
}

func H4(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(34.0/16.0), "plan9", txt)
	return
}

func H5(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(24.0/16.0), "plan9", txt)
	return
}

func H6(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(20.0/16.0), "plan9", txt)
	return
}

func Body1(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize, "bariol regular", txt)
	return
}

func Body2(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(14.0/16.0), "bariol regular", txt)
	return
}

func Caption(th *Theme, txt string) (l LabelStyle) {
	l = Label(th, th.TextSize.Scale(12.0/16.0), "bariol regular", txt)
	return
}

func Label(th *Theme, size unit.Value, font, txt string) (l LabelStyle) {
	var f text.Font
	for i := range th.Collection {
		// Debug(th.Collection[i].Font)
		if th.Collection[i].Font.Typeface == text.Typeface(font) {
			f = th.Collection[i].Font
		}
	}
	return LabelStyle{
		Text:     txt,
		Font:     f,
		Color:    th.Color.Text,
		TextSize: size,
		shaper:   th.Shaper,
	}
}

func (l LabelStyle) Layout(gtx layout.Context) layout.Dimensions {
	paint.ColorOp{Color: l.Color}.Add(gtx.Ops)
	tl := widget.Label{Alignment: l.Alignment, MaxLines: l.MaxLines}
	return tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
