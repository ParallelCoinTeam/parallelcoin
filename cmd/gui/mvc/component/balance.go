package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

var (
	itemsList = &layout.List{
		Axis: layout.Vertical,
	}
	singleItem = &layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}
	icon, _ = theme.NewDuoUIicon(icons.EditorMonetizationOn)
)

func DuoUIbalance(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		in := layout.UniformInset(unit.Dp(16))
		in.Layout(gtx, func() {
			//cs := gtx.Constraints
			navButtons := []func(){
				listItem(gtx, th, "BALANCE :", rc.Status.Wallet.Balance+" "+rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, "UNCNFIRMED :", rc.Status.Wallet.Unconfirmed+" "+rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, "TRANSACTIONS :", fmt.Sprint(rc.Status.Wallet.TxsNumber)),
			}
			itemsList.Layout(gtx, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, navButtons[i])
			})
		})
	}
}

func listItem(gtx *layout.Context, th *theme.DuoUItheme, name, value string) func() {
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(gtx, func() {
							if icon != nil {
								icon.Color = theme.HexARGB(th.Color.Dark)
								icon.Layout(gtx, unit.Px(float32(32)))
							}
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: 32, Y: 32},
							}
						})
					}),
					layout.Rigid(func() {
						txt := th.H6(name)
						txt.Font.Typeface = th.Font.Primary
						txt.Color = theme.HexARGB(th.Color.Primary)
						txt.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func() {
				value := th.H5(value)
				value.Font.Typeface = th.Font.Primary
				value.Color = theme.HexARGB(th.Color.Dark)
				value.Alignment = text.End
				value.Layout(gtx)
			}),
		)
	}
}
