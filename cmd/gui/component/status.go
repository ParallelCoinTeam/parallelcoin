package component

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/gelook"
	"github.com/p9c/pod/cmd/gui/rcd"
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
)

func DuoUIstatus(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		cs := gtx.Constraints
		gelook.DuoUIdrawRectangle(gtx, cs.Width.Max, cs.Height.Max, th.Colors["Light"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		in := layout.UniformInset(unit.Dp(16))
		in.Layout(gtx, func() {
			//cs := gtx.Constraints
			bigStatus := []func(){
				listItem(gtx, th, 22, 6, "EditorMonetizationOn", "BALANCE :", rc.Status.Wallet.Balance.Load()+" "+rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Colors["LightGrayII"]),
				listItem(gtx, th, 22, 6, "MapsLayersClear", "UNCONFIRMED :", rc.Status.Wallet.Unconfirmed.Load()+" "+
					rc.Settings.Abbrevation),
				HorizontalLine(gtx, 1, th.Colors["LightGrayII"]),
				listItem(gtx, th, 22, 6, "CommunicationImportExport", "TRANSACTIONS :", fmt.Sprint(rc.Status.Wallet.TxsNumber)),

				HorizontalLine(gtx, 1, th.Colors["LightGrayII"]),
				listItem(gtx, th, 16, 4, "DeviceWidgets", "Block Count :", fmt.Sprint(rc.Status.Node.BlockCount)),

				HorizontalLine(gtx, 1, th.Colors["LightGrayII"]),
				listItem(gtx, th, 16, 4, "ImageTimer", "Difficulty :", fmt.Sprint(rc.Status.Node.Difficulty)),

				HorizontalLine(gtx, 1, th.Colors["LightGrayII"]),
				listItem(gtx, th, 16, 4, "NotificationVPNLock", "Connections :", fmt.Sprint(rc.Status.Node.ConnectionCount)),
			}
			itemsList.Layout(gtx, len(bigStatus), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, bigStatus[i])
			})
		})
	}
}

func listItem(gtx *layout.Context, th *gelook.DuoUItheme, size, top int, iconName, name, value string) func() {
	return func() {
		icon := th.Icons[iconName]
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(float32(top)), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(gtx, func() {
							if icon != nil {
								icon.Color = gelook.HexARGB(th.Colors["Dark"])
								icon.Layout(gtx, unit.Px(float32(size)))
							}
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: size, Y: size},
							}
						})
					}),
					layout.Rigid(func() {
						txt := th.DuoUIlabel(unit.Dp(float32(size)), name)
						txt.Font.Typeface = th.Fonts["Primary"]
						txt.Color = gelook.HexARGB(th.Colors["Primary"])
						txt.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func() {
				value := th.H5(value)
				value.TextSize = unit.Dp(float32(size))
				value.Font.Typeface = th.Fonts["Primary"]
				value.Color = gelook.HexARGB(th.Colors["Dark"])
				value.Alignment = text.End
				value.Layout(gtx)
			}),
		)
	}
}
