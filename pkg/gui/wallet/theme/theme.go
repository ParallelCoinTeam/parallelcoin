package theme

import (
	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/ico"
)

type Theme struct {
	T *material.Theme
	//Shaper   text.Shaper
	//TextSize unit.Value
	TextSize float32
	//Color    struct {
	//	Primary color.RGBA
	//	Text    color.RGBA
	//	Hint    color.RGBA
	//	InvText color.RGBA
	//}
	Colors map[string]string
	Fonts  map[string]text.Typeface
	Icons  map[string]*widget.Icon

	scrollBarSize int
}

func register(fnt text.Font, ttf []byte) {
	//face, err := opentype.Parse(ttf)
	//if err != nil {
	//	panic(fmt.Sprintf("failed to parse font: %v", err))
	//}
	//fnt.Typeface = "Go"
	//font.Register(fnt, face)
}

func NewTheme() *Theme {
	t := &Theme{
		T: material.NewTheme(gofont.Collection()),
		//Shaper: font.Default(),
		TextSize: 80,
	}
	t.Colors = NewColors()
	//t.TextSize = unit.Sp(16)
	t.Icons = ico.NewDuoUIicons()

	t.T.Color.Primary = helper.HexARGB(t.Colors["Primary"])
	t.T.Color.Text = helper.HexARGB(t.Colors["Charcoal"])
	t.T.Color.Hint = helper.HexARGB(t.Colors["Silver"])
	return t
}

func (t *Theme) ChangeLightDark() {
	light := t.Colors["Light"]
	dark := t.Colors["Dark"]
	lightGray := t.Colors["LightGrayIII"]
	darkGray := t.Colors["DarkGrayII"]
	t.Colors["Light"] = dark
	t.Colors["Dark"] = light
	t.Colors["LightGrayIII"] = darkGray
	t.Colors["DarkGrayII"] = lightGray
}
