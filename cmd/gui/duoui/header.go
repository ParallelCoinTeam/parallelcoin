package duoui

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
)

var (
	logoButton = new(gel.Button)
	headerList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
)

func (ui *DuoUI) DuoUIheader() func() {
	return func() {
		ui.ly.Theme.DuoUIitem(0, ui.ly.Theme.Colors["Dark"]).Layout(ui.ly.Context, layout.NW, func() {
			layout.Flex{
				Axis:      layout.Horizontal,
				Spacing:   layout.SpaceBetween,
				Alignment: layout.Middle,
			}.Layout(ui.ly.Context,
				layout.Rigid(func() {
					var logoMeniItem gelook.DuoUIbutton
					logoMeniItem = ui.ly.Theme.DuoUIbutton("", "",
						"", ui.ly.Theme.Colors["Dark"], "", "",
						"logo", ui.ly.Theme.Colors["Light"], 16, 64, 96, 96, 8, 8)
					for logoButton.Clicked(ui.ly.Context) {
						ui.ly.Theme.ChangeLightDark()
					}
					logoMeniItem.IconLayout(ui.ly.Context, logoButton)
				}),
				layout.Flexed(1, component.HeaderMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
				layout.Rigid(component.Label(ui.ly.Context, ui.ly.Theme, ui.ly.Theme.Fonts["Primary"], 12, ui.ly.Theme.Colors["Light"], ui.rc.Status.Wallet.Balance.Load()+" "+ui.rc.Settings.Abbrevation)),
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
