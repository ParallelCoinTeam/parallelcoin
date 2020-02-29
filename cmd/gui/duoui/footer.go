package duoui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/component"

	"github.com/p9c/pod/cmd/gui/mvc/theme"
)

var (
	cornerNav = &layout.List{
		Axis: layout.Horizontal,
	}
	footerNav = &layout.List{
		Axis: layout.Horizontal,
	}
)

func (ui *DuoUI) DuoUIfooter() func() {
	return func() {
		cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 64, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

		layout.Flex{Spacing: layout.SpaceBetween}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
					cornerButtons := []func(){
						component.QuitButton(ui.rc, ui.ly.Context, ui.ly.Theme),

						component.LogButton(ui.rc, ui.ly.Context, ui.ly.Theme),
					}
					cornerNav.Layout(ui.ly.Context, len(cornerButtons), func(i int) {
						layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, cornerButtons[i])
					})
				})
			}),
			layout.Rigid(component.FooterRightMenu(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages)),
		)
	}
}
