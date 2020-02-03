package components

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"github.com/p9c/pod/pkg/gui/widget/material"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
	"image/color"
)

var (
	itemsList = &layout.List{
		Axis: layout.Vertical,
	}
	singleItem = &layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}
	Icon, _ = material.NewIcon(icons.EditorMonetizationOn)
)

func listItem(duo *models.DuoUI, name, value string) func() {
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(duo.DuoUIcontext,
			layout.Rigid(func() {
				layout.Flex{}.Layout(duo.DuoUIcontext,
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(duo.DuoUIcontext, func() {
							if Icon != nil {
								Icon.Color = hexARGB("ff303030")
								Icon.Layout(duo.DuoUIcontext, unit.Px(float32(32)))
							}
							duo.DuoUIcontext.Dimensions = layout.Dimensions{
								Size: image.Point{X: 32, Y: 32},
							}
						})
					}),
					layout.Rigid(func() {
						txt := duo.DuoUItheme.H6(name)
						txt.Color = duo.DuoUIconfiguration.SecondaryTextColor
						txt.Layout(duo.DuoUIcontext)
					}),
				)
			}),
			layout.Rigid(func() {
				value := duo.DuoUItheme.H5(value)
				value.Color = color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}
				value.Alignment = text.End
				value.Layout(duo.DuoUIcontext)
			}),
		)
	}
}

func DuoUIbalanceWidget(duo *models.DuoUI, rc *rcd.RcVar) {
	in := layout.UniformInset(unit.Dp(16))
	in.Layout(duo.DuoUIcontext, func() {
		cs := duo.DuoUIcontext.Constraints
		navButtons := []func(){
			listItem(duo, "Balance :", rc.Balance+" "+duo.DuoUIconfiguration.Abbrevation),
			func() { helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, "ffbdbdbd", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0}) },
			listItem(duo, "Unconfirmed :", rc.Unconfirmed),
			func() { helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, 1, "ffbdbdbd", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0}) },
			listItem(duo, "Transactions :", fmt.Sprint(rc.Transactions.TxsNumber)),
		}
		itemsList.Layout(duo.DuoUIcontext, len(navButtons), func(i int) {
			layout.UniformInset(unit.Dp(0)).Layout(duo.DuoUIcontext, navButtons[i])
		})
	})
}

func hexARGB(s string) (c color.RGBA) {
	_, _ = fmt.Sscanf(s, "%02x%02x%02x%02x", &c.A, &c.R, &c.G, &c.B)
	return
}
