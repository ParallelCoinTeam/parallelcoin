// SPDX-License-Identifier: Unlicense OR MIT

package gelook

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gui/fonts"
)

type DuoUITheme struct {
	Shaper        text.Shaper
	TextSize      unit.Value
	Colors        map[string]string
	Fonts         map[string]text.Typeface
	Icons         map[string]*DuoUIIcon
	scrollBarSize int
}

func init() {
	fonts.Register()
}

func NewDuoUITheme() *DuoUITheme {
	t := &DuoUITheme{
		Shaper: font.Default(),
	}
	t.Colors = NewDuoUIcolors()
	t.Fonts = NewDuoUIFonts()
	t.TextSize = unit.Sp(16)
	t.Icons = NewDuoUIIcons()
	return t
}

func NewDuoUIFonts() (f map[string]text.Typeface) {
	f = make(map[string]text.Typeface)
	f["Primary"] = "bariol"
	f["Secondary"] = "plan9"
	f["Mono"] = "go"
	return f
}

func (t *DuoUITheme) ChangeLightDark() {
	light := t.Colors["Light"]
	dark := t.Colors["Dark"]
	lightGray := t.Colors["LightGrayIII"]
	darkGray := t.Colors["DarkGrayII"]
	t.Colors["Light"] = dark
	t.Colors["Dark"] = light
	t.Colors["LightGrayIII"] = darkGray
	t.Colors["DarkGrayII"] = lightGray
}
