package duoui

import (
	"gioui.org/layout"

	"github.com/stalker-loki/pod/cmd/gui/component"
)

func (ui *DuoUI) DuoUImenu() func() {
	nav := ui.ly.Navigation
	return func() {
		nav.Width = 48
		nav.Height = 48
		nav.TextSize = 0
		nav.IconSize = 24
		nav.PaddingVertical = 4
		nav.PaddingHorizontal = 0
		if ui.ly.Viewport > 740 {
			nav.Width = 96
			nav.Height = 72
			nav.TextSize = 48
			nav.IconSize = 36
			nav.PaddingVertical = 8
			nav.PaddingHorizontal = 0
		}
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
			Spacing:   layout.SpaceEvenly}.
			Layout(ui.ly.Context, layout.Rigid(
				component.MainNavigation(ui.rc, ui.ly.Context,
					ui.ly.Theme, ui.ly.Pages, nav)),
			)
	}
}
