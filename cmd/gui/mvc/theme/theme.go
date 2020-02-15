// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/cmd/gui/ico"
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/font"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type DuoUItheme struct {
	Shaper text.Shaper
	Color  struct {
		Light     string
		Dark      string
		Bg        string
		Text      string
		Primary   string
		Secondary string
		Success   string
		Danger    string
		Warning   string
		Info      string
		Hint      string
		InvText   string
	}
	Font struct {
		Primary, Secondary, Mono text.Typeface
	}
	TextSize unit.Value
	Icons    map[string]*DuoUIicon
}

func NewDuoUItheme() *DuoUItheme {
	t := &DuoUItheme{
		Shaper: font.Default(),
	}
	t.Color.Light = "ffcfcfcf"
	t.Color.Dark = "ff303030"
	t.Color.Bg = "ffcfcfcf"
	t.Color.Text = "ff303030"
	t.Color.Primary = "ff308080"
	t.Color.Secondary = "ff803080"
	t.Color.Success = "ff30cf30"
	t.Color.Danger = "ffcf3030"
	t.Color.Warning = "ffcfcf30"
	t.Color.Info = "ff3080cf"
	t.Color.Hint = "ff888888"
	t.Color.InvText = "0xcfcfcf"
	t.Font.Primary = "bariol"
	t.Font.Secondary = "plan9"
	t.Font.Mono = "go"

	t.TextSize = unit.Sp(16)

	i := make(map[string]*DuoUIicon)
	i["checkBoxCheckedIcon"] = mustIcon(NewDuoUIicon(icons.ToggleCheckBox))
	i["checkBoxUncheckedIcon"] = mustIcon(NewDuoUIicon(icons.ToggleCheckBoxOutlineBlank))
	i["radioCheckedIcon"] = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonChecked))
	i["radioUncheckedIcon"] = mustIcon(NewDuoUIicon(icons.ToggleRadioButtonUnchecked))

	i["iconCancel"] = mustIcon(NewDuoUIicon(icons.NavigationCancel))
	i["iconOK"] = mustIcon(NewDuoUIicon(icons.NavigationCheck))
	i["iconClose"] = mustIcon(NewDuoUIicon(icons.NavigationClose))

	i["logo"] = mustIcon(NewDuoUIicon(ico.ParallelCoin))

	i["overviewIcon"] = mustIcon(NewDuoUIicon(icons.ActionHome))
	i["sendIcon"] = mustIcon(NewDuoUIicon(icons.ActionStarRate))
	i["receiveIcon"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))
	i["addressBookIcon"] = mustIcon(NewDuoUIicon(icons.ActionBook))
	i["historyIcon"] = mustIcon(NewDuoUIicon(icons.ActionHistory))


	i["closeIcon"] = mustIcon(NewDuoUIicon(icons.NavigationClose))
	i["settingsIcon"] = mustIcon(NewDuoUIicon(icons.ActionSettings))
	i["blocksIcon"] = mustIcon(NewDuoUIicon(icons.ActionExplore))
	i["networkIcon"] = mustIcon(NewDuoUIicon(icons.ActionFingerprint))
	i["traceIcon"] = mustIcon(NewDuoUIicon(icons.ActionTrackChanges))
	i["consoleIcon"] = mustIcon(NewDuoUIicon(icons.ActionInput))
	i["helpIcon"] = mustIcon(NewDuoUIicon(icons.NavigationArrowDropDown))

	t.Icons = i
	return t
}

func mustIcon(ic *DuoUIicon, err error) *DuoUIicon {
	if err != nil {
		panic(err)
	}
	return ic
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
