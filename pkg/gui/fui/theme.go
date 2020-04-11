package fui

import (
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gui/fonts"
	"image/color"
)

func init() {
	fonts.Register()
}

type Active struct {
	Colors    string
	Icons     string
	Primary   string
	Secondary string
	Mono      string
	Language  string
}
type Theme struct {
	Shaper   text.Shaper
	TextSize unit.Value
	Current  Active
	Colors   map[string]map[string]color.RGBA
	Icons    map[string]map[string]Icon
	Fonts    []string
	Lexicon  map[string]map[string]string
}

type Icon struct {
	Color color.RGBA
	src   []byte
	// Cached values.
	op       paint.ImageOp
	imgSize  int
	imgColor color.RGBA
}

func GetDefaultTheme() (out *Theme) {
	out = &Theme{}
	out.Colors["default"] = LightTheme
	out.Current.Colors = "default"
	out.Current.Primary = "bariol"
	out.Current.Secondary = "plan9"
	out.Current.Mono = "go"
	out.Fonts = []string{
		"bariol",
		"plan9",
		"go",
	}
	out.Current.Mono = "go"
	return
}
