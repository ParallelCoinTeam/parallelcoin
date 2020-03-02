package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/cmd/gui/theme"
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

func DuoUIstatus(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		cs := gtx.Constraints
		theme.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		in := layout.UniformInset(unit.Dp(16))
		in.Layout(gtx, func() {
			//cs := gtx.Constraints
			bigStatus := []func(){
				listItem(gtx, th, 32, "BALANCE :", rc.Status.Wallet.Balance+" "+rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, 32, "UNCNFIRMED :", rc.Status.Wallet.Unconfirmed+" "+rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, 32, "TRANSACTIONS :", fmt.Sprint(rc.Status.Wallet.TxsNumber)),

				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, 16, "Block Count :", fmt.Sprint(rc.Status.Node.BlockCount)),

				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, 16, "Difficulty :", fmt.Sprint(rc.Status.Node.Difficulty)),

				HorizontalLine(gtx, 1, th.Color.LightGrayII),
				listItem(gtx, th, 16, "Connections :", fmt.Sprint(rc.Status.Node.ConnectionCount)),
			}
			itemsList.Layout(gtx, len(bigStatus), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, bigStatus[i])
			})
		})
	}
}

func listItem(gtx *layout.Context, th *theme.DuoUItheme, size int, name, value string) func() {
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
								icon.Layout(gtx, unit.Px(float32(size)))
							}
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: size, Y: size},
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
				value.TextSize = unit.Dp(float32(size))
				value.Font.Typeface = th.Font.Primary
				value.Color = theme.HexARGB(th.Color.Dark)
				value.Alignment = text.End
				value.Layout(gtx)
			}),
		)
	}
}
