package p9

import (
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
	Icons         map[string]*_icon
	scrollBarSize int
}

func NewFonts() (f map[string]text.Typeface) {
	f = make(map[string]text.Typeface)
	f["Primary"] = "bariol"
	f["Secondary"] = "plan9"
	f["Mono"] = "go"
	return f
}

func NewTheme(fontCollection []text.FontFace, quit chan struct{}) (th *Theme) {
	th = &Theme{
		quit:          quit,
		Shaper:        text.NewCache(fontCollection),
		Collection:    fontCollection,
		TextSize:      unit.Sp(16),
		Colors:        NewColors(),
		Fonts:         NewFonts(),
		scrollBarSize: 0,
	}
	th.Icons = th.NewIcons()
	return
}
