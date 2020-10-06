// SPDX-License-Identifier: Unlicense OR MIT

package p9

import (
	"image/color"

	"gioui.org/f32"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	quit          chan struct{}
	Shaper        text.Shaper
	Collection    []text.FontFace
	TextSize      unit.Value
	Colors        Colors
	Fonts         map[string]text.Typeface
	Icons         map[string]*Icon
	scrollBarSize int
}
//
// func NewTheme() *Theme {
// 	t := &Theme{
// 		Shaper: font.Default(),
// 	}
// 	t.Colors = NewDuoUIcolors()
// 	// t.Fonts = NewDuoUIfonts()
// 	t.TextSize = unit.Sp(16)
// 	t.Icons = NewIcons()
// 	return t
// }

func NewFonts() (f map[string]text.Typeface) {
	f = make(map[string]text.Typeface)
	f["Primary"] = "bariol"
	f["Secondary"] = "plan9"
	f["Mono"] = "go"
	return f
}

func (th *Theme) ChangeLightDark() {
	light := th.Colors["Light"]
	dark := th.Colors["Dark"]
	lightGray := th.Colors["LightGrayIII"]
	darkGray := th.Colors["DarkGrayII"]
	th.Colors["Light"] = dark
	th.Colors["Dark"] = light
	th.Colors["LightGrayIII"] = darkGray
	th.Colors["DarkGrayII"] = lightGray
}

//
// type Theme struct {
// 	quit       chan struct{}
// 	Collection []text.FontFace
// 	Shaper     text.Shaper
// 	Color      struct {
// 		Primary, Text, Hint, InvText color.RGBA
// 	}
// 	TextSize unit.Value
// 	Icon     struct {
// 		CheckBoxChecked, CheckBoxUnchecked, RadioChecked, RadioUnchecked *widget.Icon
// 	}
// }

func NewTheme(fontCollection []text.FontFace, quit chan struct{}) *Theme {
	t := &Theme{
		quit:          quit,
		Shaper:        text.NewCache(fontCollection),
		Collection:    fontCollection,
		TextSize:      unit.Sp(16),
		Colors:        NewColors(),
		Fonts:         NewFonts(),
		Icons:         NewIcons(),
		scrollBarSize: 0,
	}
	return t
}

func fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints
	d := cs.Min
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: d}
}
