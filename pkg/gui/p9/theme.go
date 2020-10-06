package p9

import (
	"gioui.org/text"
	"gioui.org/unit"
)

type Theme struct {
	quit          chan struct{}
	shaper        text.Shaper
	collection    []text.FontFace
	textSize      unit.Value
	Colors        Colors
	fonts         map[string]text.Typeface
	icons         map[string]*_icon
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
		shaper:        text.NewCache(fontCollection),
		collection:    fontCollection,
		textSize:      unit.Sp(16),
		Colors:        NewColors(),
		fonts:         NewFonts(),
		scrollBarSize: 0,
	}
	th.icons = th.NewIcons()
	return
}
