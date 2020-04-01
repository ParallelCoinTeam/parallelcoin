package duoui

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/component"
)

func (ui *DuoUI) DuoUImenu() func() {
	return func() {
		ui.ly.Navigation.Width = 48
		ui.ly.Navigation.Height = 48
		ui.ly.Navigation.TextSize = 0
		ui.ly.Navigation.IconSize = 24
		ui.ly.Navigation.PaddingVertical = 4
		ui.ly.Navigation.PaddingHorizontal = 0
		if ui.ly.Viewport > 740 {
			ui.ly.Navigation.Width = 96
			ui.ly.Navigation.Height = 72
			ui.ly.Navigation.TextSize = 48
			ui.ly.Navigation.IconSize = 36
			ui.ly.Navigation.PaddingVertical = 8
			ui.ly.Navigation.PaddingHorizontal = 0
		}

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.Layout(ui.ly.Context,
			layout.Rigid(component.MainNavigation(ui.rc, ui.ly.Context, ui.ly.Theme, ui.ly.Pages, ui.ly.Navigation)),
		)
	}
}
