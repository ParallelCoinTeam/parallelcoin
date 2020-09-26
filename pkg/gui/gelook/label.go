// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"github.com/p9c/pod/pkg/gui/gel"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type DuoUILabel struct {
	// Face defines the text style.
	Font text.Font
	// Color is the text color.
	Color string
	// Alignment specify the text alignment.
	Alignment text.Alignment
	// MaxLines limits the number of lines. Zero means no limit.
	MaxLines int
	Text     string
	TextSize unit.Value

	shaper text.Shaper
}

func (t *DuoUITheme) H1(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(96.0/16.0), txt)
}

func (t *DuoUITheme) H2(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(60.0/16.0), txt)
}

func (t *DuoUITheme) H3(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(48.0/16.0), txt)
}

func (t *DuoUITheme) H4(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(34.0/16.0), txt)
}

func (t *DuoUITheme) H5(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(24.0/16.0), txt)
}

func (t *DuoUITheme) H6(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(20.0/16.0), txt)
}

func (t *DuoUITheme) Body1(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize, txt)
}

func (t *DuoUITheme) Body2(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(14.0/16.0), txt)
}

func (t *DuoUITheme) Caption(txt string) DuoUILabel {
	return t.DuoUILabel(t.TextSize.Scale(12.0/16.0), txt)
}

func (t *DuoUITheme) DuoUILabel(size unit.Value, txt string) DuoUILabel {
	return DuoUILabel{
		Text:     txt,
		Color:    "Dark",
		TextSize: size,
		shaper:   t.Shaper,
	}
}

func (l DuoUILabel) Layout(gtx *layout.Context) {
	paint.ColorOp{Color: HexARGB(l.Color)}.Add(gtx.Ops)
	tl := gel.Label{Alignment: l.Alignment, MaxLines: l.MaxLines}
	tl.Layout(gtx, l.shaper, l.Font, l.TextSize, l.Text)
}
