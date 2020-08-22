package duoui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/stalker-loki/pod/cmd/gui/component"
	"github.com/stalker-loki/pod/pkg/gui/gel"
	"github.com/stalker-loki/pod/pkg/gui/gelook"
)

var (
	logoButton = new(gel.Button)
	headerList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
)

func (ui *DuoUI) DuoUIheader() func() {
	th := ui.ly.Theme
	ctx := ui.ly.Context
	return func() {
		iSize := 32
		iWidth := 48
		iHeight := 48
		iPadV := 3
		iPadH := 3
		if ui.ly.Viewport > 740 {
			iSize = 64
			iWidth = 96
			iHeight = 96
			iPadV = 6
			iPadH = 6
		}
		th.DuoUIcontainer(0, th.Colors["Dark"]).Layout(ctx, layout.NW, func() {
			layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(ctx,
				layout.Rigid(func() {
					var logoMeniItem gelook.DuoUIbutton
					logoMeniItem = th.DuoUIbutton("", "", "",
						th.Colors["Dark"], "", "", "logo",
						th.Colors["Light"], 0, iSize, iWidth,
						iHeight, iPadV, iPadH, iPadV, iPadH)
					for logoButton.Clicked(ctx) {
						th.ChangeLightDark()
					}
					logoMeniItem.IconLayout(ctx, logoButton)
				}),
				layout.Flexed(1, component.HeaderMenu(ui.rc,
					ctx, th, ui.ly.Pages)),
				layout.Rigid(component.Label(ctx, th,
					th.Fonts["Primary"], 12, th.Colors["Light"],
					ui.rc.Status.Wallet.Balance.Load()+" "+ui.rc.Settings.Abbrevation)),
				layout.Rigid(component.Label(ctx, th,
					th.Fonts["Primary"], 12, th.Colors["Light"],
					fmt.Sprint(ui.ly.Viewport))),
			)
		})
	}
}

func renderIcon(gtx *layout.Context, icon *gelook.DuoUIicon) func() {
	return func() {
		icon.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0x55, B: 0x30}
		icon.Layout(gtx, unit.Dp(float32(48)))
		pointer.Rect(image.Rectangle{Max: image.Point{
			X: 64,
			Y: 64,
		}}).Add(gtx.Ops)
	}
}
